package studio

import (
	"context"
	"encoding/json"
	"io"
	"legoloop/internal/caminhos"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/provider"
)

const host = "127.0.0.1:7878"

func newServer(t *testing.T) (*Server, *atomic.Int32) {
	t.Helper()
	idx, err := graph.NewIndex(graph.DefaultCatalog())
	if err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	s := &Server{
		Index:      idx,
		ProjectDir: t.TempDir(),
		Discover: func(context.Context) provider.Inventory {
			calls.Add(1)
			return provider.Inventory{
				RouterURL: "http://127.0.0.1:20128", RouterModels: []string{"code-pro"},
				Ollama: []string{"qwen2.5-coder:7b"}, CLIs: []provider.CLI{{Name: "claude", Path: `C:\secret\path\claude.exe`}},
				APIKeys: []string{"anthropic"},
			}
		},
		Host: host,
	}
	return s, &calls
}

func do(t *testing.T, h http.Handler, method, path, body string, mut ...func(*http.Request)) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "http://"+host+path, strings.NewReader(body))
	req.Host = host
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, m := range mut {
		m(req)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), dst); err != nil {
		t.Fatalf("json: %v — corpo: %s", err, rec.Body.String())
	}
}

const validComp = `{"name":"greenfield","nodes":[
 {"id":"n1","block":"interview","routing":{"mode":"suggested","model":"sonnet","fallback":"qwen2.5-coder:7b"}},
 {"id":"n2","block":"aidd-spec","routing":{"mode":"auto","combo":"code-pro"}}],
 "edges":[{"from":"n1","to":"n2"}]}`

func TestBlocksEndpoint(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "GET", "/api/blocks", "")
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	var body struct {
		Blocks []graph.Block       `json:"blocks"`
		Types  []string            `json:"types"`
		Matrix map[string][]string `json:"matrix"`
	}
	decode(t, rec, &body)
	if len(body.Blocks) != 8 || len(body.Types) != 9 || len(body.Matrix["interview"]) != 1 {
		t.Fatalf("payload = %+v", body)
	}
}

func TestSecurityHeadersAndStatic(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "GET", "/", "")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "LegoLoop Studio") {
		t.Fatalf("index: %d", rec.Code)
	}
	h := rec.Header()
	if !strings.Contains(h.Get("Content-Security-Policy"), "default-src 'self'") ||
		h.Get("X-Content-Type-Options") != "nosniff" || h.Get("X-Frame-Options") != "DENY" ||
		h.Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("cabeçalhos de segurança ausentes: %v", h)
	}
	for _, asset := range []string{"/app.js", "/style.css", "/icons.js"} {
		if r := do(t, s.Handler(), "GET", asset, ""); r.Code != 200 {
			t.Errorf("%s: %d", asset, r.Code)
		}
	}
}

func TestDNSRebindingBlocked(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "GET", "/api/blocks", "", func(r *http.Request) { r.Host = "evil.example.com" })
	if rec.Code != http.StatusMisdirectedRequest {
		t.Fatalf("Host forjado aceito: %d", rec.Code)
	}
}

func TestInventoryHidesPathsAndCaches(t *testing.T) {
	s, calls := newServer(t)
	h := s.Handler()
	rec := do(t, h, "GET", "/api/inventory", "")
	if strings.Contains(rec.Body.String(), "secret") {
		t.Fatal("caminho local vazou no inventário")
	}
	var inv struct {
		Models []provider.Model `json:"models"`
		CLIs   []string         `json:"clis"`
		Router bool             `json:"router_online"`
	}
	decode(t, rec, &inv)
	if len(inv.Models) != 5 || inv.Models[0].ID != "code-pro" || !inv.Router || inv.CLIs[0] != "claude" {
		t.Fatalf("inventário = %+v", inv)
	}
	do(t, h, "GET", "/api/inventory", "")
	if calls.Load() != 1 {
		t.Fatalf("inventário deveria ser cacheado, chamadas=%d", calls.Load())
	}
}

func TestRouteEndpoint(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "GET", "/api/route?from=Artifact:URL&to=Artifact:Spec", "")
	var body struct {
		Path []string `json:"path"`
		Cost int      `json:"cost"`
	}
	decode(t, rec, &body)
	if strings.Join(body.Path, ",") != "obscura,rea,impeccable,aidd-spec" || body.Cost != 17 {
		t.Fatalf("rota = %+v", body)
	}
	if r := do(t, s.Handler(), "GET", "/api/route?from=Artifact:Audit&to=Artifact:URL", ""); r.Code != 422 {
		t.Fatalf("rota impossível: %d", r.Code)
	}
}

func TestValidateComposition(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "POST", "/api/validate", validComp)
	var res ValidationResult
	decode(t, rec, &res)
	if !res.OK || len(res.Issues) != 0 {
		t.Fatalf("composição válida rejeitada: %+v", res)
	}
	var d dag.DAG
	if err := json.Unmarshal(res.DAG, &d); err != nil {
		t.Fatal(err)
	}
	n2 := d.Nodes[1]
	if n2.Ref != "skills/aidd-spec" || n2.DependsOn[0] != "n1" || n2.Inputs["n1"] != "file('BRIEFING.md')" || n2.Output.Artifact != "SPEC.md" {
		t.Fatalf("nó composto incorretamente: %+v", n2)
	}
	if d.Nodes[0].Inputs["idea"] != "var.idea" {
		t.Fatalf("nó raiz deveria receber var.idea: %+v", d.Nodes[0].Inputs)
	}
}

func TestValidateDetectsProblems(t *testing.T) {
	s, _ := newServer(t)
	cases := map[string]string{
		"incompatível":              `{"name":"x","nodes":[{"id":"a","block":"interview","routing":{"mode":"manual","model":"sonnet"}},{"id":"b","block":"aidd-tdd","routing":{"mode":"manual","model":"sonnet"}}],"edges":[{"from":"a","to":"b"}]}`,
		"ciclo":                     `{"name":"x","nodes":[{"id":"a","block":"aidd-tdd","routing":{"mode":"manual","model":"sonnet"}},{"id":"b","block":"code-review","routing":{"mode":"manual","model":"sonnet"}}],"edges":[{"from":"a","to":"b"},{"from":"b","to":"a"}]}`,
		"Peça inexistente":          `{"name":"x","nodes":[{"id":"a","block":"virus","routing":{"mode":"manual","model":"sonnet"}}],"edges":[]}`,
		"artefato colidindo":        `{"name":"x","nodes":[{"id":"a","block":"interview","routing":{"mode":"manual","model":"sonnet"}},{"id":"b","block":"aidd-grill","routing":{"mode":"manual","model":"sonnet"}}],"edges":[]}`,
		"aresta órfã":               `{"name":"x","nodes":[{"id":"a","block":"interview","routing":{"mode":"manual","model":"sonnet"}}],"edges":[{"from":"a","to":"zz"}]}`,
		"aresta duplicada":          `{"name":"x","nodes":[{"id":"a","block":"interview","routing":{"mode":"manual","model":"sonnet"}},{"id":"b","block":"aidd-spec","routing":{"mode":"manual","model":"sonnet"}}],"edges":[{"from":"a","to":"b"},{"from":"a","to":"b"}]}`,
		"modelo fora do inventário": `{"name":"x","nodes":[{"id":"a","block":"interview","routing":{"mode":"manual","model":"gpt-9"}}],"edges":[]}`,
		"vazio":                     `{"name":"x","nodes":[],"edges":[]}`,
	}
	for name, body := range cases {
		rec := do(t, s.Handler(), "POST", "/api/validate", body)
		var res ValidationResult
		decode(t, rec, &res)
		if res.OK || len(res.Issues) == 0 {
			t.Errorf("%s: deveria reprovar: %+v", name, res)
		}
	}
}

func TestMaterializeWritesProcesso(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "POST", "/api/materialize", validComp)
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	data, err := os.ReadFile(caminhos.CaminhoDAG(s.ProjectDir))
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(data)
	if err != nil || d.Name != "greenfield" || len(d.Nodes) != 2 {
		t.Fatalf("Processo gravado inválido: %v", err)
	}
	bad := do(t, s.Handler(), "POST", "/api/materialize", `{"name":"x","nodes":[],"edges":[]}`)
	if bad.Code != 422 {
		t.Fatalf("composição inválida materializada: %d", bad.Code)
	}
}

func TestPostHardening(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()
	cross := do(t, h, "POST", "/api/validate", validComp, func(r *http.Request) { r.Header.Set("Origin", "https://evil.example.com") })
	if cross.Code != 403 {
		t.Fatalf("Origin cruzada aceita: %d", cross.Code)
	}
	same := do(t, h, "POST", "/api/validate", validComp, func(r *http.Request) { r.Header.Set("Origin", "http://"+host) })
	if same.Code != 200 {
		t.Fatalf("mesma origem rejeitada: %d", same.Code)
	}
	ctype := do(t, h, "POST", "/api/validate", validComp, func(r *http.Request) { r.Header.Set("Content-Type", "text/plain") })
	if ctype.Code != 415 {
		t.Fatalf("content-type simples aceito (CSRF): %d", ctype.Code)
	}
	big := `{"name":"` + strings.Repeat("a", maxBodyBytes) + `","nodes":[],"edges":[]}`
	if r := do(t, h, "POST", "/api/validate", big); r.Code != 413 {
		t.Fatalf("payload gigante: %d", r.Code)
	}
	if r := do(t, h, "POST", "/api/validate", `{"name":"x","nodes":[],"edges":[],"__proto__":{}}`); r.Code != 400 {
		t.Fatalf("campo desconhecido aceito: %d", r.Code)
	}
	if r := do(t, h, "GET", "/api/validate", ""); r.Code != 405 {
		t.Fatalf("método errado: %d", r.Code)
	}
}

func TestRateLimit(t *testing.T) {
	s, _ := newServer(t)
	s.RatePerSecond, s.Burst = 1, 3
	h := s.Handler()
	codes := map[int]int{}
	for range 10 {
		codes[do(t, h, "GET", "/api/blocks", "").Code]++
	}
	if codes[429] == 0 || codes[200] != 3 {
		t.Fatalf("rate limit não aplicado: %v", codes)
	}
}

func TestListenRejectsNonLoopback(t *testing.T) {
	s, _ := newServer(t)
	for _, addr := range []string{"0.0.0.0:7878", "192.168.0.10:7878", ":7878", "exemplo.com:80"} {
		if err := s.ListenAndServe(context.Background(), addr, nil); err == nil {
			t.Errorf("%s deveria ser recusado", addr)
		}
	}
}

func TestListenAndServeLifecycle(t *testing.T) {
	s, _ := newServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	ready := make(chan string, 1)
	done := make(chan error, 1)
	go func() { done <- s.ListenAndServe(ctx, "127.0.0.1:0", ready) }()
	addr := <-ready
	resp, err := http.Get("http://" + addr + "/api/blocks")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestJSONErrorShape(t *testing.T) {
	s, _ := newServer(t)
	rec := do(t, s.Handler(), "POST", "/api/validate", `{`)
	var body map[string]string
	decode(t, rec, &body)
	if rec.Code != 400 || body["error"] == "" {
		t.Fatalf("erro sem formato JSON: %d %s", rec.Code, rec.Body.String())
	}
}

// SEC-11: dados extras após o objeto JSON são rejeitados.
func TestTrailingJSONRejected(t *testing.T) {
	s, _ := newServer(t)
	if r := do(t, s.Handler(), "POST", "/api/validate", validComp+`{"name":"x"}`); r.Code != 400 {
		t.Fatalf("JSON extra aceito: %d", r.Code)
	}
}

// R2-5: lixo começando com '}' após o objeto também é rejeitado.
func TestTrailingBraceRejected(t *testing.T) {
	s, _ := newServer(t)
	if r := do(t, s.Handler(), "POST", "/api/validate", validComp+`}`); r.Code != 400 {
		t.Fatalf("'}' extra aceito: %d", r.Code)
	}
}

// SPEC-2026-11-01: Entrega e unificação do Web Component <lego-item-modal>.
func TestLegoItemModalWebComponentDelivery(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	// 1. Script do Web Component entregue via HTTP
	recJS := do(t, h, "GET", "/components/lego-item-modal.js", "")
	if recJS.Code != 200 {
		t.Fatalf("falha ao servir lego-item-modal.js: code %d", recJS.Code)
	}
	jsBody := recJS.Body.String()
	if !strings.Contains(jsBody, "class LegoItemModal extends HTMLElement") {
		t.Fatalf("lego-item-modal.js não contém definição da classe: %s", jsBody)
	}
	if !strings.Contains(jsBody, `customElements.define("lego-item-modal", LegoItemModal)`) {
		t.Fatalf("lego-item-modal.js não registra o custom element")
	}

	// 2. Canvas index.html inclui o web component e não tem o dialog duplicado
	recIndex := do(t, h, "GET", "/", "")
	if recIndex.Code != 200 {
		t.Fatalf("falha ao servir /: code %d", recIndex.Code)
	}
	idxBody := recIndex.Body.String()
	if !strings.Contains(idxBody, "components/lego-item-modal.js") {
		t.Fatalf("index.html não inclui script do web component")
	}
	if !strings.Contains(idxBody, `<lego-item-modal id="modal-item">`) {
		t.Fatalf("index.html não instancia <lego-item-modal id=\"modal-item\">")
	}
	if strings.Contains(idxBody, `<dialog id="modal-item" class="dialog modal-item">`) {
		t.Fatalf("index.html ainda contém dialog legado duplicado")
	}

	// 3. Inventário inventario.html inclui o web component e não tem o dialog duplicado
	recInv := do(t, h, "GET", "/inventario", "")
	if recInv.Code != 200 {
		t.Fatalf("falha ao servir /inventario: code %d", recInv.Code)
	}
	invBody := recInv.Body.String()
	if !strings.Contains(invBody, "components/lego-item-modal.js") {
		t.Fatalf("inventario.html não inclui script do web component")
	}
	if !strings.Contains(invBody, `<lego-item-modal id="modal-item">`) {
		t.Fatalf("inventario.html não instancia <lego-item-modal id=\"modal-item\">")
	}
	if strings.Contains(invBody, `<dialog id="modal-item" class="dialog modal-item">`) {
		t.Fatalf("inventario.html ainda contém dialog legado duplicado")
	}

	// 4. Verificação de fidelidade de dados: zero frases estáticas mentirosas e presença de auditoria de Gate
	if strings.Contains(jsBody, "Executa tarefas técnicas determinísticas com isolamento estrito de ambiente.") {
		t.Fatalf("lego-item-modal.js ainda contém texto genérico hardcoded")
	}
	if !strings.Contains(jsBody, "_getAuditedArtifact") || !strings.Contains(jsBody, "gate-audit-card") {
		t.Fatalf("lego-item-modal.js não implementa auditoria explícita de artefatos para Gates")
	}
	if !strings.Contains(jsBody, "triad-filters") || !strings.Contains(jsBody, "Substitutas (Swap)") {
		t.Fatalf("lego-item-modal.js não implementa a tríade ontológica de compatibilidade")
	}

	// 5. Paleta preserva campos de contrato
	recPal := do(t, h, "GET", "/palette.js", "")
	if recPal.Code != 200 {
		t.Fatalf("falha ao servir palette.js")
	}
	palBody := recPal.Body.String()
	if !strings.Contains(palBody, "consumes: b.consumes") || !strings.Contains(palBody, "artifact: b.artifact") {
		t.Fatalf("palette.js descarta campos de contrato da Peça")
	}

	fullIdx, err := graph.NewIndex(graph.FullCatalog())
	if err != nil {
		t.Fatalf("falha ao criar index FullCatalog: %v", err)
	}
	sFull := &Server{Index: fullIdx, Host: host}
	recBlocks := do(t, sFull.Handler(), "GET", "/api/blocks", "")
	if recBlocks.Code != 200 {
		t.Fatalf("falha em /api/blocks")
	}
	var blocksData struct {
		Blocks []struct {
			ID       string   `json:"id"`
			Kind     string   `json:"kind"`
			Artifact string   `json:"artifact"`
			Gate     string   `json:"gate"`
			Consumes []string `json:"consumes"`
			Produces string   `json:"produces"`
		} `json:"blocks"`
	}
	decode(t, recBlocks, &blocksData)
	var foundGithubTrigger bool
	for _, b := range blocksData.Blocks {
		if b.ID == "github-board-trigger" {
			foundGithubTrigger = true
			if b.Artifact != "WAVE_PLAN.json" {
				t.Fatalf("artefato incorreto para github-board-trigger: %s", b.Artifact)
			}
			if b.Gate != "file_exists('WAVE_PLAN.json')" {
				t.Fatalf("gate incorreto para github-board-trigger: %s", b.Gate)
			}
			if len(b.Consumes) != 0 {
				t.Fatalf("gatilho não deve consumir nada, consumiu %v", b.Consumes)
			}
		}
	}
	if !foundGithubTrigger {
		t.Fatalf("github-board-trigger não encontrado em /api/blocks com FullCatalog")
	}

	// 7. Validação das novas peças de Segurança, Git e Deploy em /api/blocks
	expectedEp21 := map[string]string{
		"env-vault":     "security_vault",
		"git-commit":    "vcs_git_commit",
		"deploy-docker": "deploy_docker",
	}
	for id, expectedKind := range expectedEp21 {
		var found bool
		for _, b := range blocksData.Blocks {
			if b.ID == id {
				found = true
				if b.Kind != expectedKind {
					t.Fatalf("Peça %s tem kind incorreto: got %s, want %s", id, b.Kind, expectedKind)
				}
				break
			}
		}
		if !found {
			t.Fatalf("Peça %s do Épico 21 não encontrado em /api/blocks", id)
		}
	}
}
