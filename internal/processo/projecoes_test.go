package processo_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/processo"
	"legoloop/internal/studiodoc"
)

var agora = time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)

// quatro é um processo de 4 etapas encadeadas, com títulos e peças distintos.
func quatro() processo.Processo {
	p := valido()
	titulos := map[string]string{"e01": "Clarificar o pedido", "e02": "Construir o módulo", "e03": "Verificar o módulo", "e04": "Documentar o módulo"}
	pecas := map[string]string{"e01": "aidd-grill", "e02": "aidd-tdd", "e03": "aidd-qa", "e04": "aidd-docs"}
	p.Etapas = nil
	for i, id := range []string{"e01", "e02", "e03", "e04"} {
		e := etapa(id)
		e.Titulo, e.Peca = titulos[id], pecas[id]
		if i > 0 {
			e.Depende = []string{p.Etapas[i-1].ID}
		}
		p.Etapas = append(p.Etapas, e)
	}
	return p
}

func projetar(t *testing.T, p processo.Processo) *processo.Projecoes {
	t.Helper()
	pr, err := processo.Projetar(p, agora)
	if err != nil {
		t.Fatalf("Projetar: %v", err)
	}
	return pr
}

// Extratores independentes do código de produção: cada arquivo é lido do seu jeito, de propósito.
func idsDag(t *testing.T, b []byte) []string {
	t.Helper()
	var d struct{ Nodes []struct{ ID string } }
	if err := json.Unmarshal(b, &d); err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, n := range d.Nodes {
		out = append(out, n.ID)
	}
	return out
}

func idsRoute(t *testing.T, b []byte) []string {
	t.Helper()
	var d struct {
		Nodes []struct{ Type, Label string }
	}
	if err := json.Unmarshal(b, &d); err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile("^`([^`]+)` — ")
	var out []string
	for _, n := range d.Nodes {
		if n.Type == "unit_task" {
			out = append(out, re.FindStringSubmatch(n.Label)[1])
		}
	}
	return out
}

func idsPorRegex(md, re string) []string {
	var out []string
	for _, m := range regexp.MustCompile(re).FindAllStringSubmatch(md, -1) {
		out = append(out, m[1])
	}
	return out
}

func idsPlan(md string) []string {
	return idsPorRegex(md, "(?m)^#### Ticket `TICK-[^`]+` — `([^`]+)` — ")
}
func idsSkill(md string) []string { return idsPorRegex(md, "(?m)^\\d+\\. Etapa `([^`]+)`: ") }
func idsSpec(md string) []string {
	return idsPorRegex(md, "(?m)^\\* \\[ \\] `AC-[^`]+`: \\*\\*Etapa `([^`]+)`")
}

func TestProjecoesCoerentesSeisArquivosContamAMesmaHistoria(t *testing.T) {
	pr := projetar(t, quatro())
	quero := []string{"e01", "e02", "e03", "e04"}
	vistas := map[string][]string{
		"dag": idsDag(t, pr.DAG), "loop": idsDag(t, pr.Loop), "route": idsRoute(t, pr.Route),
		"plan": idsPlan(pr.Plan), "skill": idsSkill(pr.Skill), "spec": idsSpec(pr.Spec),
	}
	for nome, got := range vistas {
		if !slices.Equal(got, quero) {
			t.Errorf("%s: ids %v, queria %v", nome, got, quero)
		}
	}
	if err := pr.Verificar(quatro()); err != nil {
		t.Errorf("Verificar: %v", err)
	}
}

func TestProjecoesCoerentesOrdemSegueAsDependencias(t *testing.T) {
	p := quatro()
	p.Etapas[0], p.Etapas[3] = p.Etapas[3], p.Etapas[0] // declaradas fora de ordem: e04, e02, e03, e01
	ord, err := p.Ordem()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, e := range ord {
		ids = append(ids, e.ID)
	}
	if !slices.Equal(ids, []string{"e01", "e02", "e03", "e04"}) {
		t.Fatalf("Ordem = %v", ids)
	}
	pr := projetar(t, p)
	for nome, got := range map[string][]string{"dag": idsDag(t, pr.DAG), "route": idsRoute(t, pr.Route), "plan": idsPlan(pr.Plan), "skill": idsSkill(pr.Skill), "spec": idsSpec(pr.Spec)} {
		if !slices.Equal(got, ids) {
			t.Errorf("%s fora da ordem de dependência: %v", nome, got)
		}
	}
}

func TestProjecoesCoerentesOrdemEhEstavelEntreIndependentes(t *testing.T) {
	p := quatro()
	for i := range p.Etapas {
		p.Etapas[i].Depende = nil
	}
	p.Etapas[3].Depende = []string{"e01", "e02", "e03"}
	ord, err := p.Ordem()
	if err != nil || ord[0].ID != "e01" || ord[1].ID != "e02" || ord[2].ID != "e03" || ord[3].ID != "e04" {
		t.Errorf("etapas independentes mantêm a ordem declarada: %v %v", ord, err)
	}
}

// corromper remove a etapa e03 de uma só projeção, do jeito próprio de cada formato.
func corromper(t *testing.T, pr *processo.Projecoes, qual string) {
	t.Helper()
	tiraNo := func(b []byte) []byte {
		var d map[string]any
		if err := json.Unmarshal(b, &d); err != nil {
			t.Fatal(err)
		}
		nodes := d["nodes"].([]any)
		var keep []any
		for _, n := range nodes {
			m := n.(map[string]any)
			if m["id"] == "e03" || strings.Contains(asString(m["label"]), "`e03`") {
				continue
			}
			keep = append(keep, n)
		}
		d["nodes"] = keep
		out, _ := json.Marshal(d)
		return out
	}
	tiraBloco := func(md string) string { // remove as linhas que citam `e03`
		var keep []string
		for _, l := range strings.Split(md, "\n") {
			if !strings.Contains(l, "`e03`") {
				keep = append(keep, l)
			}
		}
		return strings.Join(keep, "\n")
	}
	switch qual {
	case "dag":
		pr.DAG = tiraNo(pr.DAG)
	case "loop":
		pr.Loop = tiraNo(pr.Loop)
	case "route":
		pr.Route = tiraNo(pr.Route)
	case "plan":
		pr.Plan = tiraBloco(pr.Plan)
	case "skill":
		pr.Skill = tiraBloco(pr.Skill)
	case "spec":
		pr.Spec = tiraBloco(pr.Spec)
	}
}

func asString(v any) string { s, _ := v.(string); return s }

func TestProjecoesCoerentesRemoverEtapaDeUmaSoProjecaoFalha(t *testing.T) {
	for _, qual := range []string{"dag", "loop", "route", "plan", "skill", "spec"} {
		t.Run(qual, func(t *testing.T) {
			pr := projetar(t, quatro())
			corromper(t, pr, qual)
			err := pr.Verificar(quatro())
			if !errors.Is(err, processo.ErrIncoerente) {
				t.Fatalf("esperava ErrIncoerente, veio %v", err)
			}
			if !strings.Contains(err.Error(), qual) {
				t.Errorf("o erro deve dizer qual projeção divergiu (%s): %v", qual, err)
			}
		})
	}
}

func TestProjecoesCoerentesTrocarAOrdemDeUmaProjecaoFalha(t *testing.T) {
	pr := projetar(t, quatro())
	var d map[string]any
	_ = json.Unmarshal(pr.DAG, &d)
	n := d["nodes"].([]any)
	n[1], n[2] = n[2], n[1]
	pr.DAG, _ = json.Marshal(d)
	if err := pr.Verificar(quatro()); !errors.Is(err, processo.ErrIncoerente) {
		t.Errorf("ordem trocada no DAG deve falhar: %v", err)
	}
}

func TestProjecoesCoerentesPlanNaoTemAplicarASkill(t *testing.T) {
	pr := projetar(t, quatro())
	if strings.Contains(pr.Plan, "Aplicar a skill") {
		t.Error("o PLAN não pode ter 'Aplicar a skill': o título é o da etapa")
	}
	for _, e := range quatro().Etapas {
		if !strings.Contains(pr.Plan, e.Titulo) || !strings.Contains(pr.Plan, "peça `"+e.Peca+"`") {
			t.Errorf("o PLAN deve trazer o título e a peça da etapa %s", e.ID)
		}
		if !strings.Contains(pr.Skill, "aplique a skill `"+e.Peca+"`") {
			t.Errorf("a skill deve mapear a etapa %s à peça %s", e.ID, e.Peca)
		}
	}
}

func TestProjecoesCoerentesDeterministicas(t *testing.T) {
	a, b := projetar(t, quatro()), projetar(t, quatro())
	if string(a.DAG) != string(b.DAG) || string(a.Route) != string(b.Route) || a.Plan != b.Plan || a.Spec != b.Spec || a.Skill != b.Skill || string(a.Loop) != string(b.Loop) {
		t.Error("a mesma entrada deve gerar os mesmos bytes")
	}
}

func TestProjecoesCoerentesDagELoopSaoValidos(t *testing.T) {
	pr := projetar(t, quatro())
	d, err := dag.Parse(pr.DAG)
	if err != nil {
		t.Fatalf("dag.Parse: %v", err)
	}
	if d.Nodes[1].Ref != "skills/aidd-tdd" || !slices.Equal(d.Nodes[1].DependsOn, []string{"e01"}) {
		t.Errorf("nó e02: %+v", d.Nodes[1])
	}
	if d.Nodes[1].Output.Artifact != "out/e02.md" || d.Nodes[1].Output.Gate != "file_exists('out/e02.md')" {
		t.Errorf("handoff do nó e02: %+v", d.Nodes[1].Output)
	}
	var lp studiodoc.Loop
	if err := json.Unmarshal(pr.Loop, &lp); err != nil || lp.ID != "legoloop-build-engine" || len(lp.Nodes) != 4 {
		t.Errorf("loop: %v %+v", err, lp)
	}
}

func TestProjecoesCoerentesRouteLigaAsDependenciasDasEtapas(t *testing.T) {
	var d struct {
		Nodes []struct {
			ID, Type  string
			DependsOn []string `json:"depends_on"`
		}
	}
	if err := json.Unmarshal(projetar(t, quatro()).Route, &d); err != nil {
		t.Fatal(err)
	}
	var tickets []string
	deps := map[string][]string{}
	for _, n := range d.Nodes {
		if n.Type == "unit_task" {
			tickets = append(tickets, n.ID)
			deps[n.ID] = n.DependsOn
		}
	}
	if !slices.Equal(tickets, []string{"TICK-01.1", "TICK-01.2", "TICK-01.3", "TICK-01.4"}) {
		t.Fatalf("tickets = %v", tickets)
	}
	if !slices.Contains(deps["TICK-01.2"], "TICK-01.1") || !slices.Contains(deps["TICK-01.4"], "TICK-01.3") {
		t.Errorf("a ordem das etapas vira dependência entre tickets: %v", deps)
	}
}

func TestProjecoesCoerentesSpecSegueOContratoDoTemplate(t *testing.T) {
	pr := projetar(t, quatro())
	for _, h := range []string{"1. Premissas Técnicas Assumidas", "2. Fronteiras de Ação & Regras de Governança", "3. Visão Geral, Arquitetura & Costuras de Teste",
		"4. Contratos de Dados & Schemas Formais", "5. Critérios Binários de Aceite por Nó", "6. Comandos Operacionais & Critério de Conclusão Global", "STOP-GATE 02"} {
		if !strings.Contains(pr.Spec, h) {
			t.Errorf("a SPEC não tem a seção %q", h)
		}
	}
	if regexp.MustCompile(`\{\{\s*[A-Za-z_][^}\n]*\}\}`).MatchString(pr.Spec) {
		t.Error("a SPEC deixou placeholder")
	}
	if len(pr.Avisos) != 0 {
		t.Errorf("todo AC da SPEC deve estar coberto por um ticket: %v", pr.Avisos)
	}
}

func TestProjecoesCoerentesLacunaNaoViraNoExecutavel(t *testing.T) {
	p := quatro()
	p.Etapas[2].Peca = ""
	if _, err := processo.Projetar(p, agora); !errors.Is(err, processo.ErrLacuna) || !errors.Is(err, processo.ErrInvalido) {
		t.Errorf("etapa sem peça e sem composição: %v", err)
	}
}

func TestProjecoesCoerentesEtapaComposta(t *testing.T) {
	p := quatro()
	p.Etapas[2].Peca, p.Etapas[2].Compoe = "", "legoloop-outro@1.2.0"
	pr := projetar(t, p)
	d, err := dag.Parse(pr.DAG)
	if err != nil {
		t.Fatal(err)
	}
	if d.Nodes[2].Type != dag.TypeSubloop || d.Nodes[2].Ref != "legoloop-outro" {
		t.Errorf("etapa composta vira subloop: %+v", d.Nodes[2])
	}
	if !strings.Contains(pr.Plan, "processo `legoloop-outro@1.2.0`") {
		t.Error("o PLAN deve dizer qual processo a etapa compõe")
	}
}

func TestProjecoesCoerentesRecusaIdOuPecaQueNaoViramReferenciaSegura(t *testing.T) {
	for nome, muda := range map[string]func(*processo.Processo){
		"id com ponto":     func(p *processo.Processo) { p.Etapas[0].ID = "e.01"; p.Etapas[1].Depende = []string{"e.01"} },
		"id com espaço":    func(p *processo.Processo) { p.Etapas[0].ID = "e 01"; p.Etapas[1].Depende = []string{"e 01"} },
		"peça com barra":   func(p *processo.Processo) { p.Etapas[1].Peca = "../fora" },
		"peça com espaço":  func(p *processo.Processo) { p.Etapas[1].Peca = "aidd tdd" },
		"peça com crase":   func(p *processo.Processo) { p.Etapas[1].Peca = "a`b" },
		"compõe sem nome":  func(p *processo.Processo) { p.Etapas[1].Peca, p.Etapas[1].Compoe = "", "@1.0.0" },
		"compõe com barra": func(p *processo.Processo) { p.Etapas[1].Peca, p.Etapas[1].Compoe = "", "a/b@1.0.0" },
	} {
		t.Run(nome, func(t *testing.T) {
			p := quatro()
			muda(&p)
			if _, err := processo.Projetar(p, agora); !errors.Is(err, processo.ErrInvalido) {
				t.Errorf("esperava ErrInvalido, veio %v", err)
			}
		})
	}
}

func TestProjecoesCoerentesTituloMalicioso(t *testing.T) {
	p := quatro()
	p.Etapas[1].Titulo = "Construir\n## 7. Seção falsa\n* [ ] `AC-X.9`: **Etapa `e99`** `perigo`"
	pr := projetar(t, p)
	secao := regexp.MustCompile(`(?m)^#+ 7\.`)
	if secao.MatchString(pr.Plan) || secao.MatchString(pr.Skill) || secao.MatchString(pr.Spec) {
		t.Error("o título não pode abrir seção nova")
	}
	if !slices.Equal(idsSpec(pr.Spec), []string{"e01", "e02", "e03", "e04"}) || !slices.Equal(idsPlan(pr.Plan), []string{"e01", "e02", "e03", "e04"}) {
		t.Errorf("o título não pode forjar etapa: spec=%v plan=%v", idsSpec(pr.Spec), idsPlan(pr.Plan))
	}
	if err := pr.Verificar(p); err != nil {
		t.Errorf("Verificar: %v", err)
	}
}

func TestProjecoesCoerentesSkillTemFrontmatterValido(t *testing.T) {
	sk := projetar(t, quatro()).Skill
	if !strings.HasPrefix(sk, "---\nname: legoloop-build-engine\n") {
		t.Errorf("frontmatter: %q", sk[:60])
	}
	desc := regexp.MustCompile(`(?m)^description: (.+)$`).FindStringSubmatch(sk)
	if desc == nil || len([]rune(desc[1])) > 300 {
		t.Errorf("description ausente ou acima de 300 caracteres: %v", desc)
	}
}

func TestProjecoesCoerentesProcessoInvalidoNaoProjeta(t *testing.T) {
	p := quatro()
	p.Nome = "sem-prefixo"
	if _, err := processo.Projetar(p, agora); !errors.Is(err, processo.ErrInvalido) {
		t.Errorf("Validate deve barrar antes de projetar: %v", err)
	}
	c := quatro()
	c.Etapas[0].Depende = []string{"e04"}
	if _, err := c.Ordem(); !errors.Is(err, processo.ErrCiclo) {
		t.Errorf("ciclo: %v", err)
	}
}

func TestProjecoesCoerentesOrdemEscolheAPrimeiraDeclaradaQueEstaPronta(t *testing.T) {
	p := quatro()
	a, b, c := etapa("a", "b"), etapa("b"), etapa("c", "b")
	p.Etapas = []processo.Etapa{a, b, c}
	ord, err := p.Ordem()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, e := range ord {
		ids = append(ids, e.ID)
	}
	if !slices.Equal(ids, []string{"b", "a", "c"}) {
		t.Errorf("depois de b, a (declarada antes) vem antes de c: %v", ids)
	}
}

func TestProjecoesCoerentesPecaEComposicaoJuntasSaoRecusadas(t *testing.T) {
	p := quatro()
	p.Etapas[1].Compoe = "legoloop-outro@1.0.0" // já tem Peca
	_, err := processo.Projetar(p, agora)
	if !errors.Is(err, processo.ErrInvalido) || !strings.Contains(err.Error(), "peça ou processo composto") {
		t.Errorf("peça e composição ao mesmo tempo: %v", err)
	}
}

func TestProjecoesCoerentesMensagensDeReferenciaInsegura(t *testing.T) {
	for nome, c := range map[string]struct {
		muda func(*processo.Processo)
		quer string
	}{
		"id com ponto":       {func(p *processo.Processo) { p.Etapas[0].ID = "e.01"; p.Etapas[1].Depende = []string{"e.01"} }, "o id deve ter"},
		"caminho com espaço": {func(p *processo.Processo) { p.Etapas[0].Produz.Caminho = "out/a b.md" }, "não pode ter espaço"},
		"caminho com crase":  {func(p *processo.Processo) { p.Etapas[0].Produz.Caminho = "out/a`b.md" }, "não pode ter espaço"},
		"peça insegura":      {func(p *processo.Processo) { p.Etapas[1].Peca = "../x" }, "não é seguro"},
	} {
		t.Run(nome, func(t *testing.T) {
			p := quatro()
			c.muda(&p)
			if _, err := processo.Projetar(p, agora); err == nil || !strings.Contains(err.Error(), c.quer) {
				t.Errorf("a mensagem deve dizer o motivo (%q): %v", c.quer, err)
			}
		})
	}
}

func TestProjecoesCoerentesNoDoDagTemEntradasEVersao(t *testing.T) {
	pr := projetar(t, quatro())
	d, err := dag.Parse(pr.DAG)
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Nodes[1].Inputs["e01"]; got != "file('out/e01.md')" {
		t.Errorf("o nó consome o artefato da etapa anterior: %q", got)
	}
	if d.Nodes[0].Inputs != nil {
		t.Errorf("etapa sem dependência não consome nada: %v", d.Nodes[0].Inputs)
	}
	var lp studiodoc.Loop
	_ = json.Unmarshal(pr.Loop, &lp)
	if d.Version != "1.0.0" || lp.Version != "1.0.0" || !lp.CreatedAt.Equal(agora) {
		t.Errorf("a versão do processo vai ao DAG e ao loop, e a data vem do relógio injetado: %q %q %v", d.Version, lp.Version, lp.CreatedAt)
	}
}

func TestProjecoesCoerentesTituloNaoQuebraTrechosDeCodigo(t *testing.T) {
	p := quatro()
	p.Etapas[1].Titulo = "Construir `com crase` aqui"
	pr := projetar(t, p)
	var head, ac string
	for _, l := range strings.Split(pr.Plan, "\n") {
		if strings.HasPrefix(l, "#### Ticket") && strings.Contains(l, "`e02`") {
			head = l
		}
	}
	for _, l := range strings.Split(pr.Spec, "\n") {
		if strings.HasPrefix(l, "* [ ]") && strings.Contains(l, "`e02`") {
			ac = l
		}
	}
	if strings.Count(head, "`") != 6 || strings.Count(ac, "`") != 4 {
		t.Errorf("crase do título deve virar apóstrofo (plan=%d, spec=%d crases):\n%s\n%s", strings.Count(head, "`"), strings.Count(ac, "`"), head, ac)
	}
}

func TestProjecoesCoerentesEtapaComposta_SkillMandaExecutarOProcesso(t *testing.T) {
	p := quatro()
	p.Etapas[2].Peca, p.Etapas[2].Compoe = "", "legoloop-outro@1.2.0"
	sk := projetar(t, p).Skill
	if !strings.Contains(sk, "Etapa `e03`: execute o processo `legoloop-outro@1.2.0`") {
		t.Errorf("a skill deve mandar executar o processo composto:\n%s", sk)
	}
}

func TestProjecoesCoerentesDescricaoDaSkillRespeitaOLimite(t *testing.T) {
	p := quatro()
	p.Objetivo = strings.Repeat("objetivo muito longo ", 60)
	sk := projetar(t, p).Skill
	desc := regexp.MustCompile(`(?m)^description: (.+)$`).FindStringSubmatch(sk)
	if desc == nil || len([]rune(desc[1])) > 300 {
		t.Errorf("description passou de 300 caracteres: %d", len([]rune(desc[1])))
	}
}

func TestProjetarEmiteDAGNome(t *testing.T) {
	p := quatro()
	pr := projetar(t, p)
	if pr.DAGNome == "" {
		t.Fatal("DAGNome não pode ser vazio")
	}
	if !strings.Contains(pr.DAGNome, "_DAG_") || !strings.HasSuffix(pr.DAGNome, ".json") {
		t.Errorf("DAGNome fora do padrão ISO YYYY-MM-DD_DAG_<SLUG>.json: %q", pr.DAGNome)
	}
}

func TestCheckParity(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.json")
	f2 := filepath.Join(dir, "b.json")
	f3 := filepath.Join(dir, "c.json")
	if err := os.WriteFile(f1, []byte("conteudo identical"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("conteudo identical"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f3, []byte("conteudo diferente"), 0o644); err != nil {
		t.Fatal(err)
	}

	igual, err := processo.CheckParity(f1, f2)
	if err != nil || !igual {
		t.Errorf("CheckParity deve retornar true para arquivos idênticos, obteve %v, err=%v", igual, err)
	}
	diferente, err := processo.CheckParity(f1, f3)
	if err != nil || diferente {
		t.Errorf("CheckParity deve retornar false para arquivos distintos, obteve %v, err=%v", diferente, err)
	}
}
