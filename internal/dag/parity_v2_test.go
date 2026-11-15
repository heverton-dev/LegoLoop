package dag_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/combo"
	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/taxonomy"
)

func TestSchemaV2AndTriggerParity(t *testing.T) {
	t.Run("SchemaV2EParseParidade", testSchemaV2EParseParidade)
	t.Run("SensibilidadeCredenciaisEEnv", testSensibilidadeCredenciais)
	t.Run("TriggerUnificacaoEPortas", testTriggerUnificacaoEPortas)
	t.Run("QuatroParesRouteEDagAlinhados", testQuatroParesAlinhados)
}

func testSchemaV2EParseParidade(t *testing.T) {
	validarPropriedadesSchemaV2(t)
	doc := montarDocV2Completo()
	d, err := dag.Parse(doc)
	if err != nil {
		t.Fatalf("DAG v2 completo com mcp/role/combo deveria validar: %v", err)
	}
	if d.Nodes[0].Role != "engenheiro-backend" || d.Nodes[0].Papel != "engenheiro-backend" {
		t.Fatalf("role/papel nao carregados: %+v", d.Nodes[0])
	}
	if d.Nodes[0].Routing.Via != "9router" || d.Nodes[0].Routing.Combo != "t-codigo-128k" {
		t.Fatalf("via ou combo do n1 nao carregados: %+v", d.Nodes[0].Routing)
	}
	if d.Nodes[1].Type != dag.TypeMCP || d.Nodes[1].Routing.Via != "oficial" {
		t.Fatalf("no mcp ou via do n2 nao carregados: %+v", d.Nodes[1])
	}
	verificarRejeicoesV2(t)
}

func validarPropriedadesSchemaV2(t *testing.T) {
	t.Helper()
	raw, err := os.ReadFile(filepath.FromSlash("../../schemas/dag.v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	var sc map[string]any
	if err := json.Unmarshal(raw, &sc); err != nil {
		t.Fatal(err)
	}
	nodes := sc["properties"].(map[string]any)["nodes"].(map[string]any)
	props := nodes["items"].(map[string]any)["properties"].(map[string]any)
	enumTipos := props["type"].(map[string]any)["enum"].([]any)
	if !contemStringAny(enumTipos, "mcp") {
		t.Fatalf("schemas/dag.v2.json deve aceitar type mcp; enum=%v", enumTipos)
	}
	for _, k := range []string{"role", "papel", "rule_tags", "token_budget", "test_seam"} {
		if _, ok := props[k]; !ok {
			t.Errorf("schemas/dag.v2.json ausente nodes[].properties.%s", k)
		}
	}
	rProps := props["routing"].(map[string]any)["properties"].(map[string]any)
	for _, k := range []string{"via", "combo"} {
		if _, ok := rProps[k]; !ok {
			t.Errorf("schemas/dag.v2.json ausente routing.properties.%s", k)
		}
	}
}

func contemStringAny(xs []any, alvo string) bool {
	for _, x := range xs {
		if s, ok := x.(string); ok && s == alvo {
			return true
		}
	}
	return false
}

func montarDocV2Completo() []byte {
	return []byte(`{
	  "$schema":"https://legoloop.dev/schemas/dag.v2.json",
	  "schema_version":2,
	  "name":"paridade-v2",
	  "version":"2.0.0",
	  "date":"2026-11-14",
	  "execution_strategy":{"supervision_level":"micro"},
	  "upstream_signatures":{"route":{"path":"docs/routes/x.json","sha256":"abc"}},
	  "trigger":{"type":"manual"},
	  "nodes":[
	    {
	      "id":"n1","name":"n1","type":"unit_no","nature":"micro",
	      "role":"engenheiro-backend","papel":"engenheiro-backend",
	      "rule_tags":["tdd","ste"],
	      "token_budget":{"max_input":8000,"max_output":2000},
	      "test_seam":{"target":"internal/dag","verify":"go test ./internal/dag/..."},
	      "routing":{"mode":"auto","harness":"9router","via":"9router","combo":"t-codigo-128k","effort":"medium"},
	      "execution":{"mode":"isolated"},"isolation":{"write_boundary":"strict"},
	      "limits":{"max_output_tokens":2000},"tooling":{"runtime":"go"},
	      "verification":{"gate":"file_exists('A.md')"},"rollback_strategy":"stash_changes",
	      "spec_ref":"AC-1","status":"PENDING","context_in":[],"context_out":["A.md"],
	      "output":{"artifact":"A.md","gate":"file_exists('A.md')"}
	    },
	    {
	      "id":"n2","name":"n2","type":"mcp","nature":"micro","ref":"rea",
	      "command":["rea","inspect"],"depends_on":["n1"],
	      "routing":{"mode":"manual","harness":"claude","via":"oficial","model":"sonnet","effort":"low"},
	      "execution":{"mode":"inline"},"isolation":{"write_boundary":"strict"},
	      "limits":{"max_output_tokens":500},"tooling":{"runtime":"mcp"},
	      "verification":{"gate":"file_exists('B.md')"},"rollback_strategy":"stash_changes",
	      "spec_ref":"AC-2","status":"PENDING","context_in":["A.md"],"context_out":["B.md"],
	      "output":{"artifact":"B.md","gate":"file_exists('B.md')"}
	    }
	  ]
	}`)
}

func verificarRejeicoesV2(t *testing.T) {
	t.Helper()
	casos := map[string]string{
		"role invalido":    `"role":"backend-developer","papel":"backend-developer"`,
		"role papel dif":   `"role":"arquiteto","papel":"engenheiro-backend"`,
		"mcp sem ref":      `__NO_REF__`,
		"mcp ref ruim":     `"ref":"../evil"`,
		"mcp cmd relativo": `"ref":"rea","command":["./evil"]`,
	}
	for nome, sub := range casos {
		var js string
		if sub == "__NO_REF__" {
			js = casoMcpSemRef()
		} else if strings.Contains(sub, "role") {
			js = casoRoleInvalido(sub)
		} else {
			js = casoMcpInvalido(sub)
		}
		if _, err := dag.Parse([]byte(js)); err == nil {
			t.Errorf("%s deveria falhar em dag.Parse", nome)
		}
	}
}

func casoRoleInvalido(campos string) string {
	return `{"name":"x","version":"1","trigger":{"type":"manual"},"nodes":[{` +
		`"id":"a","type":"unit_no",` + campos + `,"routing":{"mode":"manual","model":"m"},` +
		`"output":{"artifact":"A.md","gate":"true"}}]}`
}

func casoMcpSemRef() string {
	return `{"name":"x","version":"1","trigger":{"type":"manual"},"nodes":[{` +
		`"id":"a","type":"mcp","routing":{"mode":"manual","model":"m"},` +
		`"output":{"artifact":"A.md","gate":"true"}}]}`
}

func casoMcpInvalido(extra string) string {
	return `{"name":"x","version":"1","trigger":{"type":"manual"},"nodes":[{` +
		`"id":"a","type":"mcp",` + extra + `,"routing":{"mode":"manual","model":"m"},` +
		`"output":{"artifact":"A.md","gate":"true"}}]}`
}

func testSensibilidadeCredenciais(t *testing.T) {
	caminhos := []string{
		".env",
		".env.local",
		".env-empresa",
		"config/.env.prod",
		"credential_gemini.json",
		"sub/credential_backup.json",
	}
	for _, p := range caminhos {
		js := `{"name":"x","version":"1","trigger":{"type":"manual"},"nodes":[{` +
			`"id":"a","type":"unit_no","routing":{"mode":"manual","model":"m"},` +
			`"output":{"artifact":"` + p + `","gate":"true"}}]}`
		if _, err := dag.Parse([]byte(js)); err == nil {
			t.Errorf("caminho sensivel %q deveria ser bloqueado em artifact", p)
		}
	}
}

func testTriggerUnificacaoEPortas(t *testing.T) {
	idx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		t.Fatal(err)
	}
	aliasKind := "trigger_" + string(dag.TriggerFimLegado)
	if got := graph.CanonicalTriggerKind(aliasKind); got != graph.KindTriggerProcessoFim {
		t.Fatalf("CanonicalTriggerKind(%q) = %q, esperado %q", aliasKind, got, graph.KindTriggerProcessoFim)
	}
	aliasID := "trigger-" + strings.ReplaceAll(string(dag.TriggerFimLegado), "_", "-")
	bAlias, ok := idx.Block(aliasID)
	bCanon, okCanon := idx.Block("trigger-processo-fim")
	if !ok || !okCanon || bAlias.ID != bCanon.ID {
		t.Fatalf("alias %q (%v) e canonico %q (%v) devem resolver para a mesma peca", aliasID, ok, bCanon.ID, okCanon)
	}
	validarPortasTriggerEMatrix(t, idx)
}

func validarPortasTriggerEMatrix(t *testing.T, idx *graph.Index) {
	t.Helper()
	mat := idx.Matrix()
	pecas := idx.Blocks()
	for _, a := range pecas {
		trigA := graph.IsTriggerKind(a.Kind)
		var esperados []string
		for _, b := range pecas {
			trigB := graph.IsTriggerKind(b.Kind)
			pode, err := idx.CanConnect(a.ID, b.ID)
			if err != nil {
				t.Fatalf("CanConnect(%s, %s): %v", a.ID, b.ID, err)
			}
			if trigB && pode {
				t.Errorf("gatilho %s nao deve aceitar entrada de %s", b.ID, a.ID)
			}
			if trigA && !trigB && !pode {
				t.Errorf("gatilho %s deve conectar a %s", a.ID, b.ID)
			}
			if pode {
				esperados = append(esperados, b.ID)
			}
		}
		if !slices.Equal(mat[a.ID], esperados) {
			t.Errorf("Matrix[%s] = %v diverge de CanConnect = %v", a.ID, mat[a.ID], esperados)
		}
	}
}

type parRouteDag struct {
	route string
	dag   string
}

func testQuatroParesAlinhados(t *testing.T) {
	pares := []parRouteDag{
		{"../../docs/routes/2026-11-10_ROUTE_DISSECAR.json", "../../docs/dags/2026-11-10_DAG_DISSECAR.json"},
		{"../../docs/routes/2026-11-13_ROUTE_STUDIO_PAPEIS_E_COMBOS.json", "../../docs/dags/2026-11-13_DAG_STUDIO_PAPEIS_E_COMBOS.json"},
		{"../../docs/routes/2026-11-14_ROUTE_DESIGN_UI.json", "../../docs/dags/2026-11-14_DAG_DESIGN_UI.json"},
		{"../../docs/routes/2026-11-14_ROUTE_SWARM.json", "../../docs/dags/2026-11-14_DAG_SWARM.json"},
	}
	for _, p := range pares {
		validarParRouteDag(t, p)
	}
}

type routeParityDoc struct {
	UpstreamSignatures map[string]struct {
		SHA256 string `json:"sha256"`
	} `json:"upstream_signatures"`
	Nodes []routeParityNode `json:"nodes"`
}

type routeParityNode struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Role      string   `json:"role"`
	DependsOn []string `json:"depends_on"`
	Routing   struct {
		Via   string `json:"via"`
		Combo string `json:"combo"`
		Model string `json:"model"`
	} `json:"routing"`
}

func validarParRouteDag(t *testing.T, p parRouteDag) {
	t.Helper()
	rawR, err := os.ReadFile(filepath.FromSlash(p.route))
	if err != nil {
		t.Fatal(err)
	}
	rawD, err := os.ReadFile(filepath.FromSlash(p.dag))
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(rawD)
	if err != nil {
		t.Fatalf("%s nao passa em dag.Parse: %v", filepath.Base(p.dag), err)
	}
	if d.SchemaVersion != 2 {
		t.Fatalf("%s: schema_version = %d, esperado 2", filepath.Base(p.dag), d.SchemaVersion)
	}
	var r, dg routeParityDoc
	if err := json.Unmarshal(rawR, &r); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawD, &dg); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(rawR)
	wantHash := hex.EncodeToString(sum[:])
	if got := dg.UpstreamSignatures["route"].SHA256; got != wantHash {
		t.Errorf("%s: upstream route sha256 = %s, esperado %s", filepath.Base(p.dag), got, wantHash)
	}
	compararNosRouteDag(t, filepath.Base(p.dag), r.Nodes, d.Nodes)
}

func compararNosRouteDag(t *testing.T, nome string, rNos []routeParityNode, dNos []dag.Node) {
	t.Helper()
	if len(rNos) != len(dNos) {
		t.Fatalf("%s: %d nos no DAG vs %d na ROUTE", nome, len(dNos), len(rNos))
	}
	for i, rn := range rNos {
		dn := dNos[i]
		if rn.ID != dn.ID || !slices.Equal(rn.DependsOn, dn.DependsOn) {
			t.Errorf("%s[%d]: divergencia estrutural ROUTE(%s,%v) != DAG(%s,%v)",
				nome, i, rn.ID, rn.DependsOn, dn.ID, dn.DependsOn)
		}
		if dn.Type == dag.TypeStopGate {
			continue
		}
		validarComboEPapelNo(t, nome, dn)
	}
}

func validarComboEPapelNo(t *testing.T, nome string, dn dag.Node) {
	t.Helper()
	if dn.Routing.Via == "" {
		t.Errorf("%s[%s]: routing.via vazio", nome, dn.ID)
	}
	c := dn.Routing.Combo
	if c == "" {
		c = dn.Routing.Model
	}
	if _, ok := combo.CatalogoCombos[c]; !ok {
		t.Errorf("%s[%s]: combo/model %q nao certificado em CatalogoCombos", nome, dn.ID, c)
	}
	if dn.Role != "" {
		if err := taxonomy.ValidarPapel(dn.Role); err != nil {
			t.Errorf("%s[%s]: role %q invalido: %v", nome, dn.ID, dn.Role, err)
		}
	}
}
