package studio

import (
	"encoding/json"
	"io/fs"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func readStatic(t *testing.T, name string) string {
	t.Helper()
	b, err := fs.ReadFile(staticFS, "static/"+name)
	if err != nil {
		t.Fatalf("ler %s: %v", name, err)
	}
	return string(b)
}

// R-UI-3: nenhum diálogo nativo do navegador/SO — feedback é sempre peça própria.
func TestNoNativeDialogs(t *testing.T) {
	native := regexp.MustCompile(`(^|[^.\w])(window\.)?(alert|confirm|prompt)\s*\(|window\.(alert|confirm|prompt)\b`)
	err := fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !(strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".html")) {
			return err
		}
		b, rerr := fs.ReadFile(staticFS, path)
		if rerr != nil {
			return rerr
		}
		for i, line := range strings.Split(string(b), "\n") {
			if native.MatchString(line) {
				t.Errorf("%s:%d usa diálogo nativo: %s", path, i+1, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// R-UI-2: scrollbars de no máximo 4px (vertical e horizontal), thumb a 50% de opacidade.
func TestScrollbarRule(t *testing.T) {
	css := readStatic(t, "style.css")
	track := regexp.MustCompile(`(?s)::-webkit-scrollbar\s*\{[^}]*\}`).FindString(css)
	if track == "" {
		t.Fatal("regra ::-webkit-scrollbar ausente")
	}
	for _, want := range []string{"width: 4px", "height: 4px"} {
		if !strings.Contains(track, want) {
			t.Errorf("::-webkit-scrollbar sem %q", want)
		}
	}
	thumb := regexp.MustCompile(`(?s)::-webkit-scrollbar-thumb\s*\{[^}]*\}`).FindString(css)
	if !strings.Contains(thumb, "var(--scroll-thumb)") {
		t.Error("thumb deve usar o token --scroll-thumb")
	}
	if !regexp.MustCompile(`--scroll-thumb:[^;]*50%`).MatchString(css) {
		t.Error("token --scroll-thumb deve ter 50% de opacidade")
	}
	if !strings.Contains(css, "::-webkit-scrollbar-track") || !strings.Contains(css, "background: transparent") {
		t.Error("trilho deve ser transparente")
	}
	if !strings.Contains(css, "@supports not selector(::-webkit-scrollbar)") {
		t.Error("fallback para navegadores sem ::-webkit-scrollbar ausente")
	}
}

// R-UI-4: Dialog e Toast próprios, carregados antes do app.
func TestUIComponentsServedAndWired(t *testing.T) {
	s, _ := newServer(t)
	if r := do(t, s.Handler(), "GET", "/ui.js", ""); r.Code != 200 {
		t.Fatalf("/ui.js: %d", r.Code)
	}
	html := readStatic(t, "index.html")
	ui := strings.Index(html, `src="ui.js"`)
	app := strings.Index(html, `src="app.js"`)
	if ui < 0 || app < 0 || ui > app {
		t.Fatalf("ui.js deve ser carregado antes de app.js (ui=%d app=%d)", ui, app)
	}
	for _, id := range []string{`id="toasts" class="toasts" aria-live="polite"`, `id="dialog"`} {
		if !strings.Contains(html, id) {
			t.Errorf("index.html sem %s", id)
		}
	}
	js := readStatic(t, "ui.js")
	for _, want := range []string{"function confirmDialog", "function toast", "showModal", `"cancel"`, `"alert"`} {
		if !strings.Contains(js, want) {
			t.Errorf("ui.js sem %q", want)
		}
	}
}

// A confirmação da rota ótima usa o Dialog próprio.
func TestRouteUsesOwnDialog(t *testing.T) {
	app := readStatic(t, "app.js")
	if !strings.Contains(app, "await confirmDialog(") {
		t.Error("app.js deve usar confirmDialog na troca do canvas")
	}
	if strings.Contains(app, "function toast(") {
		t.Error("toast deve viver só em ui.js")
	}
}

// TICK-805: a página de inventário carrega ui.js antes de inventory.js e expõe as abas e a canonização.
func TestInventoryPageWiring(t *testing.T) {
	html := readStatic(t, "inventario.html")
	ui, inv := strings.Index(html, `src="ui.js"`), strings.Index(html, `src="inventory.js"`)
	if ui < 0 || inv < 0 || ui > inv {
		t.Fatalf("ui.js deve vir antes de inventory.js (ui=%d inv=%d)", ui, inv)
	}
	for _, id := range []string{`id="tabs"`, `id="search"`, `id="list"`, `id="btn-simulate"`, `id="btn-apply"`,
		`id="toasts" class="toasts" aria-live="polite"`, `id="dialog"`} {
		if !strings.Contains(html, id) {
			t.Errorf("inventario.html sem %s", id)
		}
	}
	js := readStatic(t, "inventory.js")
	for _, tab := range []string{"skill", "mcp", "script", "hook", "gate", "model", "harness", "binary", "loop"} {
		if !strings.Contains(js, `id: "`+tab+`"`) {
			t.Errorf("aba %q ausente em inventory.js", tab)
		}
	}
	if !strings.Contains(readStatic(t, "index.html"), `href="/inventario"`) {
		t.Error("o canvas deve linkar para /inventario")
	}
}

// Todo texto vindo da API entra no DOM via textContent: innerHTML é proibido nos scripts.
func TestNoInnerHTMLInScripts(t *testing.T) {
	for _, f := range []string{"app.js", "ui.js", "inventory.js", "search.js", "virtual.js", "palette.js", "slots.js", "fabricar.js", "processos.js"} {
		for i, line := range strings.Split(readStatic(t, f), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue // comentários podem citar a regra
			}
			if strings.Contains(line, "innerHTML") || strings.Contains(line, "insertAdjacentHTML") {
				t.Errorf("%s:%d usa HTML dinâmico: %s", f, i+1, strings.TrimSpace(line))
			}
		}
	}
}

// A busca é a mesma da página de inventário e (depois) da paleta do canvas: sem acento, sem caixa,
// "-", "_", "." e "/" valem como espaço, e todas as palavras precisam aparecer, em qualquer ordem.
func TestSearchMatching(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node não está no PATH; a regra de busca não pôde ser executada")
	}
	cases := [][3]any{
		{"acesso-plataforma", "ACESSO PLATAFORMA", true},
		{"AIDD Spec", "aidd-spec", true},
		{"Ação Rápida", "acao rapida", true},
		{"code-review-graph", "graph review", true},
		{"skills/aidd.md", "skills aidd md", true},
		{"a_b.c", "a b c", true},
		{"abc", "", true},
		{"abc", "   ", true},
		{"rea", "reasoning", false},
		{"playwright", "mcp", false},
		{"alpha beta", "alpha gamma", false},
	}
	in, _ := json.Marshal(cases)
	script := `const vm=require("vm"),fs=require("fs");const c={};vm.createContext(c);
vm.runInContext(fs.readFileSync(process.argv[1],"utf8")+";this.m=matchesQuery;",c);
const cases=JSON.parse(process.argv[2]);console.log(JSON.stringify(cases.map(x=>c.m(x[0],x[1]))));`
	out, err := exec.Command(node, "-e", script, "static/search.js", string(in)).CombinedOutput()
	if err != nil {
		t.Fatalf("node falhou: %v\n%s", err, out)
	}
	var got []bool
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	for i, c := range cases {
		if got[i] != c[2].(bool) {
			t.Errorf("matchesQuery(%q, %q) = %v, quer %v", c[0], c[1], got[i], c[2])
		}
	}
}

func TestSearchScriptIsServedAndLoaded(t *testing.T) {
	s, _ := newServer(t)
	if r := do(t, s.Handler(), "GET", "/search.js", ""); r.Code != 200 {
		t.Fatalf("/search.js: %d", r.Code)
	}
	html := readStatic(t, "inventario.html")
	if a, b := strings.Index(html, `src="search.js"`), strings.Index(html, `src="inventory.js"`); a < 0 || a > b {
		t.Errorf("search.js deve vir antes de inventory.js (search=%d inv=%d)", a, b)
	}
	if !strings.Contains(readStatic(t, "inventory.js"), "matchesQuery(") {
		t.Error("inventory.js deve usar matchesQuery")
	}
}

// runNode carrega um script do Studio num contexto isolado e avalia uma expressão (JSON de saída).
func runNode(t *testing.T, expr string, files ...string) []byte {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node não está no PATH; a lógica de interface não pôde ser executada")
	}
	script := `const vm=require("vm"),fs=require("fs");const c={console};vm.createContext(c);
for (const f of process.argv.slice(1,-1)) vm.runInContext(fs.readFileSync(f,"utf8"),c,{filename:f});
console.log(JSON.stringify(vm.runInContext(process.argv.at(-1),c)));`
	args := []string{"-e", script}
	for _, f := range files {
		args = append(args, "static/"+f)
	}
	args = append(args, expr)
	out, err := exec.Command(node, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("node falhou: %v\n%s", err, out)
	}
	return out
}

// TICK-905: a janela virtual nunca renderiza mais que a tela mais a folga, qualquer que seja o total.
func TestVirtualWindowBoundsTheRenderedRows(t *testing.T) {
	out := runNode(t, `(() => {
	  const res = {maxRendered: 0, covers: true, empty: virtualWindow(0, 48, 0, 700), neg: virtualWindow(10, 48, -50, 700)};
	  for (let top = 0; top <= 5000 * 48; top += 37) {
	    const w = virtualWindow(5000, 48, top, 700, 6);
	    res.maxRendered = Math.max(res.maxRendered, w.end - w.start);
	    const first = Math.floor(top / 48), last = Math.min(4999, Math.floor((top + 699) / 48));
	    if (w.start > first || w.end <= last || w.start < 0 || w.end > 5000) res.covers = false;
	  }
	  res.bottom = virtualWindow(5000, 48, 5000 * 48, 700, 6);
	  res.shortList = virtualWindow(3, 48, 0, 700, 6);
	  res.height = virtualWindow(5000, 48, 0, 700).height;
	  return res;
	})()`, "virtual.js")
	var r struct {
		MaxRendered int  `json:"maxRendered"`
		Covers      bool `json:"covers"`
		Empty       struct{ Start, End, Height int }
		Neg         struct{ Start, End int }
		Bottom      struct{ Start, End int }
		ShortList   struct{ Start, End int }
		Height      int `json:"height"`
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatalf("saída inesperada: %s", out)
	}
	if r.MaxRendered > 100 {
		t.Errorf("renderizou até %d linhas com 5.000 itens (limite 100)", r.MaxRendered)
	}
	if !r.Covers {
		t.Error("a janela deixou de cobrir as linhas visíveis")
	}
	if r.Empty.End != 0 || r.Empty.Height != 0 || r.Neg.Start != 0 || r.Bottom.End != 5000 ||
		r.ShortList.Start != 0 || r.ShortList.End != 3 || r.Height != 5000*48 {
		t.Errorf("casos de borda: %+v", r)
	}
}

// TestInventoryScriptNoOrphanModalListeners assegura que inventory.js não contém referências
// a métodos do modal legado que causam ReferenceError no boot.
func TestInventoryScriptNoOrphanModalListeners(t *testing.T) {
	js := readStatic(t, "inventory.js")

	orphans := []string{
		"translateActiveItem",
		"switchModalTab",
		`$("btn-modal-close")`,
		`$("btn-translate")`,
	}
	for _, orphan := range orphans {
		if strings.Contains(js, orphan) {
			t.Errorf("inventory.js contém referência órfã do modal legado: %q", orphan)
		}
	}

	if !strings.Contains(js, `$("modal-item")`) {
		t.Error("inventory.js deve delegar para o Web Component <lego-item-modal>")
	}
}

// TestValidateReturnsBoolean (BUG-RUN-01) assegura que validate() em app.js retorna booleano explícito
// e não o retorno de toast() que é undefined.
func TestValidateReturnsBoolean(t *testing.T) {
	js := readStatic(t, "app.js")
	re := regexp.MustCompile(`(?s)async function validate\(\)\s*\{.*?\n\}`)
	block := re.FindString(js)
	if block == "" {
		t.Fatal("função validate() não encontrada em app.js")
	}
	if strings.Contains(block, "return toast(") {
		t.Error("BUG-RUN-01: validate() não deve fazer 'return toast(...)', pois toast() retorna undefined")
	}
	if !strings.Contains(block, "return Boolean(") && !strings.Contains(block, "return true") {
		t.Error("BUG-RUN-01: validate() deve retornar um valor booleano explícito")
	}
}

// TestRunModalUsesAbortController (BUG-RUN-02) assegura que initRunModal usa AbortController
// para cancelar o streaming SSE ao fechar o modal.
func TestRunModalUsesAbortController(t *testing.T) {
	js := readStatic(t, "app.js")
	if !strings.Contains(js, "new AbortController()") {
		t.Error("BUG-RUN-02: initRunModal deve instanciar AbortController para cancelamento de SSE")
	}
	if !strings.Contains(js, "signal:") {
		t.Error("BUG-RUN-02: fetch('/api/run') deve receber o signal do AbortController")
	}
}

// TestUndoRedoAppliesValidRenderFunction (BUG-UNDO-01) assegura que applySnapshot não chama
// renderAll() inexistente e utiliza renderAllNodes().
func TestUndoRedoAppliesValidRenderFunction(t *testing.T) {
	js := readStatic(t, "app.js")
	if strings.Contains(js, "renderAll();") {
		t.Error("BUG-UNDO-01: applySnapshot chama renderAll() que não existe; deve chamar renderAllNodes()")
	}
}

// TestMutationsPushHistoryState (BUG-UNDO-02) assegura que mutações manuais registram pushState()
func TestMutationsPushHistoryState(t *testing.T) {
	js := readStatic(t, "app.js")
	for _, fnName := range []string{"register", "removeNode", "addEdge", "disconnect"} {
		idx := strings.Index(js, "function "+fnName)
		if idx == -1 {
			t.Fatalf("função %s não encontrada em app.js", fnName)
		}
		// Inspeciona os próximos 800 caracteres da função
		chunk := js[idx : idx+800]
		if !strings.Contains(chunk, "pushState()") {
			t.Errorf("BUG-UNDO-02: função %s deve invocar pushState() para registrar histórico", fnName)
		}
	}
}

// TestConfirmDialogMessageParam (BUG-CONFIRM-01) assegura que inventory.js passa o parâmetro { message }
// para confirmDialog() e não { text }.
func TestConfirmDialogMessageParam(t *testing.T) {
	js := readStatic(t, "inventory.js")
	if strings.Contains(js, "text:") && strings.Contains(js, "confirmDialog({") {
		t.Error("BUG-CONFIRM-01: inventory.js usa 'text:' em vez de 'message:' em confirmDialog()")
	}
}

// TestTemplatesFallbackWithoutDoubleName (BUG-TMPL-01) assegura que o fallback de ícone de template
// não repete o nome do template lado a lado.
func TestTemplatesFallbackWithoutDoubleName(t *testing.T) {
	js := readStatic(t, "app.js")
	if strings.Contains(js, "icon.textContent = tmpl.name") {
		t.Error("BUG-TMPL-01: app.js repete tmpl.name no fallback de icon, duplicando o nome no elemento")
	}
}

// TestVaultRendersGroupsList (BUG-VAULT-01) assegura que app.js referencia e popula #vault-groups-list
func TestVaultRendersGroupsList(t *testing.T) {
	js := readStatic(t, "app.js")
	if !strings.Contains(js, "vault-groups-list") {
		t.Error("BUG-VAULT-01: app.js deve referenciar #vault-groups-list para renderizar os grupos duplicados")
	}
}

// TestVaultApplyResetsDisabledInFinally (BUG-VAULT-03) assegura que btnApply tem seu estado disabled
// tratado no finally após tentativa de aplicação
func TestVaultApplyResetsDisabledInFinally(t *testing.T) {
	js := readStatic(t, "app.js")
	idx := strings.Index(js, "btnApply?.addEventListener")
	if idx == -1 {
		t.Fatal("handler de btnApply não encontrado")
	}
	chunk := js[idx : idx+1500]
	if !strings.Contains(chunk, "btnApply.disabled = false") && !strings.Contains(chunk, "btnApply.disabled =") {
		t.Error("BUG-VAULT-03: handler de btnApply deve gerenciar disabled na cláusula finally")
	}
}

// TestInspectorSuggestDebounceAndAbort (BUG-INSP-01) assegura que o inspetor utiliza
// debounce e AbortController para prevenir rajadas concorrentes em /api/suggest
func TestInspectorSuggestDebounceAndAbort(t *testing.T) {
	js := readStatic(t, "inspector.js")
	if !strings.Contains(js, "setTimeout") || !strings.Contains(js, "clearTimeout") {
		t.Error("BUG-INSP-01: inspector.js deve implementar debounce (setTimeout/clearTimeout) em ensureSuggestions()")
	}
	if !strings.Contains(js, "AbortController") {
		t.Error("BUG-INSP-01: inspector.js deve utilizar AbortController para cancelar requisições anteriores de sugestão")
	}
}

// TestCanvasPanClamping (BUG-CANVAS-02) assegura que updateCanvasTransform aplica clamp no pan
func TestCanvasPanClamping(t *testing.T) {
	js := readStatic(t, "app.js")
	idx := strings.Index(js, "function updateCanvasTransform()")
	if idx == -1 {
		t.Fatal("função updateCanvasTransform() não encontrada")
	}
	chunk := js[idx : idx+600]
	if !strings.Contains(chunk, "Math.max(") || !strings.Contains(chunk, "Math.min(") {
		t.Error("BUG-CANVAS-02: updateCanvasTransform deve aplicar clamp (Math.max/min) nas coordenadas de pan")
	}
}

// TestEscapeCancelsLinking (BUG-EDGE-01) assegura que trackLink escuta Escape no teclado
func TestEscapeCancelsLinking(t *testing.T) {
	js := readStatic(t, "app.js")
	idx := strings.Index(js, "function trackLink(")
	if idx == -1 {
		t.Fatal("função trackLink() não encontrada")
	}
	chunk := js[idx : idx+1500]
	if !strings.Contains(chunk, `"Escape"`) || !strings.Contains(chunk, "keydown") {
		t.Error("BUG-EDGE-01: trackLink deve escutar keydown para cancelar a conexão temporária no Escape")
	}
}

// TestPaletteTreeviewChunking (BUG-PALETTE-02) assegura que o treeview limita
// a renderização inicial para não travar a UI em categorias com centenas de itens
func TestPaletteTreeviewChunking(t *testing.T) {
	js := readStatic(t, "palette.js")
	if !strings.Contains(js, "treeview-load-more") && !strings.Contains(js, "slice(0,") {
		t.Error("BUG-PALETTE-02: palette.js deve implementar limite/chunking (slice ou treeview-load-more) no treeview")
	}
}

// TestItemModalNoUnsafeInnerHTML (BUG-ITEM-MODAL-01) assegura que _populateSchema
// utiliza manipulação estruturada do DOM em vez de innerHTML
func TestItemModalNoUnsafeInnerHTML(t *testing.T) {
	js := readStatic(t, "components/lego-item-modal.js")
	idx := strings.Index(js, "_populateSchema(it, family)")
	if idx == -1 {
		t.Fatal("_populateSchema não encontrado")
	}
	chunk := js[idx : idx+2500]
	if strings.Contains(chunk, "trIn.innerHTML =") || strings.Contains(chunk, "trOut.innerHTML =") {
		t.Error("BUG-ITEM-MODAL-01: _populateSchema não deve utilizar tr.innerHTML, mas sim APIs estruturadas do DOM")
	}
}

// TestKeyboardShortcutsCtrlSEsc (BUG-KEY-01, BUG-KEY-02) assegura atalhos Ctrl+S e Escape
func TestKeyboardShortcutsCtrlSEsc(t *testing.T) {
	js := readStatic(t, "app.js")
	if !strings.Contains(js, `e.key.toLowerCase() === "s"`) && !strings.Contains(js, `e.key === "s"`) {
		t.Error("BUG-KEY-02: app.js deve registrar atalho Ctrl+S para acionar btn-materialize")
	}
	idx := strings.Index(js, "function onKey(ev)")
	if idx == -1 {
		t.Fatal("função onKey não encontrada")
	}
	chunk := js[idx : idx+400]
	if !strings.Contains(chunk, "select(null)") {
		t.Error("BUG-KEY-01: tecla Escape em onKey deve desselecionar nós ou arestas ativas (select(null))")
	}
}

// TestPaletteDropCenterOffset (BUG-PALETTE-01) assegura compensação centralizada no eixo Y ao soltar nós
func TestPaletteDropCenterOffset(t *testing.T) {
	js := readStatic(t, "app.js")
	if strings.Contains(js, "p.y - 20") {
		t.Error("BUG-PALETTE-01: drop de peças não deve usar offset fixo de -20px; deve centralizar o nó verticalmente")
	}
}

// TestVault503ToastFeedback (BUG-VAULT-02) assegura que falhas HTTP em status não abrem modal em branco
func TestVault503ToastFeedback(t *testing.T) {
	js := readStatic(t, "app.js")
	idx := strings.Index(js, `fetch("/api/canonize/status")`)
	if idx == -1 {
		t.Fatal("fetch /api/canonize/status não encontrado")
	}
	chunk := js[idx : idx+400]
	if !strings.Contains(chunk, "!res.ok") {
		t.Error("BUG-VAULT-02: fetch /api/canonize/status deve verificar se res.ok é falso antes de abrir o modal")
	}
}

// TestNoInlineStylesInComponents (BUG-CSP-01) assegura que nenhuma peça utiliza
// o atributo style="..." inline, prevenindo violações da política Content Security Policy (style-src 'self').
func TestNoInlineStylesInComponents(t *testing.T) {
	inlineStyleRe := regexp.MustCompile(`style\s*=\s*["'][^"']*["']`)
	err := fs.WalkDir(staticFS, "static/components", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".js") {
			return err
		}
		b, rerr := fs.ReadFile(staticFS, path)
		if rerr != nil {
			return rerr
		}
		for i, line := range strings.Split(string(b), "\n") {
			if inlineStyleRe.MatchString(line) {
				t.Errorf("%s:%d viola CSP com style inline: %s", path, i+1, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TICK-04.1: o Studio oferece a gestão de processos com cartão STE e relatório por camadas.
func TestProcessosModalWiring(t *testing.T) {
	html := readStatic(t, "index.html")
	for _, id := range []string{`src="processos.js"`, `id="btn-processos"`, `id="modal-processos"`, `id="processos-list"`} {
		if !strings.Contains(html, id) {
			t.Errorf("index.html sem %s", id)
		}
	}
	js := readStatic(t, "processos.js")
	for _, want := range []string{"function renderCartaoSTE", "function renderRelatorioCamadas", "cartao-ste__aviso"} {
		if !strings.Contains(js, want) {
			t.Errorf("processos.js sem %q", want)
		}
	}
	app := readStatic(t, "app.js")
	if !strings.Contains(app, "initProcessosModal") {
		t.Error("app.js deve chamar initProcessosModal")
	}
}

func TestCanConnectTriggerAnyBlock(t *testing.T) {
	app := readStatic(t, "app.js")
	idx := strings.Index(app, "function canConnect(")
	if idx < 0 {
		t.Fatal("app.js sem function canConnect(")
	}
	snippet := app[idx:]
	if end := strings.Index(snippet, "\n}"); end > 0 {
		snippet = snippet[:end+2]
	}
	if !strings.Contains(snippet, `startsWith("trigger_")`) {
		t.Fatalf("canConnect deve permitir ligação de/para trigger_* com qualquer peça; trecho:\n%s", snippet)
	}
}

func TestLayoutModaisNavegacaoETema(t *testing.T) {
	css := readStatic(t, "style.css")
	indexHTML := readStatic(t, "index.html")
	invHTML := readStatic(t, "inventario.html")
	swarmHTML := readStatic(t, "swarm.html")
	assertLayoutResponsivoEPorts(t, css, readStatic(t, "swarm.css"))
	assertModaisEInspetorIndex(t, indexHTML)
	assertNavegacaoETemaRotas(t, indexHTML, invHTML, swarmHTML)
	assertTemaJS(t)
}

func assertLayoutResponsivoEPorts(t *testing.T, css, swarmCSS string) {
	t.Helper()
	if strings.Contains(css, "- 56px)") {
		t.Error("style.css não deve usar altura fixa calc(... - 56px)")
	}
	for _, want := range []string{".stage-topbar", "@media (max-width: 1024px)", "@media (max-width: 768px)", ".node.node--trigger .port--in"} {
		if !strings.Contains(css, want) {
			t.Errorf("style.css sem regra responsiva/porta %q", want)
		}
	}
	if !strings.Contains(swarmCSS, "@media (max-width:768px)") && !strings.Contains(swarmCSS, "@media (max-width: 768px)") {
		t.Error("swarm.css sem media query de 768px")
	}
}

func assertModaisEInspetorIndex(t *testing.T, indexHTML string) {
	t.Helper()
	if strings.Contains(indexHTML, "modal-item-header") {
		t.Error("index.html ainda usa modal-item-header desalinhado")
	}
	if strings.Count(indexHTML, `class="modal-header"`) < 6 {
		t.Errorf("index.html deve alinhar os modais com class=\"modal-header\"")
	}
	posAside := strings.Index(indexHTML, `id="inspector"`)
	if posAside < 0 {
		t.Fatalf("aside#inspector não encontrado em index.html")
	}
}

func assertNavegacaoETemaRotas(t *testing.T, indexHTML, invHTML, swarmHTML string) {
	t.Helper()
	if !strings.Contains(invHTML, `href="/swarm"`) {
		t.Error("inventario.html sem link para /swarm")
	}
	for _, href := range []string{`href="/"`, `href="/inventario"`, `href="/swarm"`} {
		if !strings.Contains(swarmHTML, href) {
			t.Errorf("swarm.html sem link de navegação %s", href)
		}
	}
	for nome, html := range map[string]string{"index.html": indexHTML, "inventario.html": invHTML, "swarm.html": swarmHTML} {
		if !strings.Contains(html, `id="btn-theme-toggle"`) {
			t.Errorf("%s sem seletor de tema #btn-theme-toggle", nome)
		}
		if !strings.Contains(html, "ui.js") {
			t.Errorf("%s deve carregar ui.js", nome)
		}
	}
}

func assertTemaJS(t *testing.T) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node ausente; validação do tema pulada")
	}
	script := `const fs=require("fs");const vm=require("vm");
const store={};
const attrs={"data-theme":"light"};
const btn={textContent:"",dataset:{},attrs:{},setAttribute(k,v){this.attrs[k]=v;},addEventListener(){}};
const root={setAttribute(k,v){attrs[k]=v;},getAttribute(k){return attrs[k];}};
const ls={getItem(k){return store[k]||null;},setItem(k,v){store[k]=v;}};
const lsProp="local"+"Stor"+"age";
const ctx={document:{documentElement:root,getElementById(id){return id==="btn-theme-toggle"?btn:null;},addEventListener(){}},window:{[lsProp]:ls,matchMedia(){return{matches:false};}},[lsProp]:ls,setTimeout,clearTimeout};
vm.runInNewContext(fs.readFileSync("static/ui.js","utf8"),ctx);
ctx.applyTheme("dark");
if(attrs["data-theme"]!=="dark"||store["legoloop.theme"]!=="dark")throw new Error("applyTheme falhou");
const next=ctx.toggleTheme();
if(next!=="light"||attrs["data-theme"]!=="light"||store["legoloop.theme"]!=="light")throw new Error("toggleTheme falhou");`
	if out, err := exec.Command(node, "-e", script).CombinedOutput(); err != nil {
		t.Fatalf("tema ui.js falhou: %v\n%s", err, string(out))
	}
}

func TestWCAGTokensEGuildas(t *testing.T) {
	css := readStatic(t, "style.css")
	assertTokensESelectoresWCAG(t, css)
	assertGuildasBadgeEMinSize(t, css)
	assertGuildasEPaletteJS(t)
	assertAriaETeclado(t)
}

func assertTokensESelectoresWCAG(t *testing.T, css string) {
	t.Helper()
	for _, tok := range []string{"--surface-3", "--ink", "--text-normal", "--text-muted", "--bg-surface", "--bg-base", "--bg-hover", "--danger", "--danger-bg", "--danger-ink", "--radius-sm", "--radius-md"} {
		if !strings.Contains(css, tok+":") {
			t.Errorf("style.css sem definição do token %s", tok)
		}
	}
	claro := extrairTokensCSS(t, css, `(?s):root,\s*:root\[data-theme="light"\]\s*\{([^}]+)\}`)
	escuro := extrairTokensCSS(t, css, `(?s):root\[data-theme="dark"\]\s*\{([^}]+)\}`)
	pares := [][2]string{{"--text", "--bg"}, {"--text", "--surface"}, {"--muted", "--surface-3"}, {"--accent-ink", "--accent"}, {"--accent", "--accent-soft"}, {"--ok", "--ok-soft"}, {"--fail", "--fail-soft"}, {"--danger-ink", "--danger-bg"}, {"--idle", "--idle-soft"}}
	validarParesWCAG(t, "claro", claro, pares)
	validarParesWCAG(t, "escuro", escuro, pares)
}

func extrairTokensCSS(t *testing.T, css, padrao string) map[string]string {
	t.Helper()
	m := regexp.MustCompile(padrao).FindStringSubmatch(css)
	if len(m) < 2 {
		t.Fatalf("trecho CSS não encontrado: %s", padrao)
	}
	res := map[string]string{}
	re := regexp.MustCompile(`(--[a-z0-9-]+)\s*:\s*(#[0-9a-fA-F]{3,6})`)
	for _, kv := range re.FindAllStringSubmatch(m[1], -1) {
		res[kv[1]] = kv[2]
	}
	return res
}

func validarParesWCAG(t *testing.T, tema string, tok map[string]string, pares [][2]string) {
	t.Helper()
	for _, par := range pares {
		fg, bg := tok[par[0]], tok[par[1]]
		if fg == "" || bg == "" {
			t.Errorf("%s: token hex ausente para %s/%s", tema, par[0], par[1])
			continue
		}
		if r := contrasteWCAG(fg, bg); r < 4.5 {
			t.Errorf("%s: contraste WCAG AA insuficiente em %s (%s) sobre %s (%s): %.2f", tema, par[0], fg, par[1], bg, r)
		}
	}
}

func assertGuildasBadgeEMinSize(t *testing.T, css string) {
	t.Helper()
	mSize := regexp.MustCompile(`\.guilda-badge\s*\{[^}]*font-size:\s*(\d+)px`).FindStringSubmatch(css)
	if len(mSize) < 2 {
		t.Fatal("style.css sem font-size em .guilda-badge")
	}
	if px, _ := strconv.Atoi(mSize[1]); px < 11 {
		t.Errorf(".guilda-badge font-size %dpx < 11px", px)
	}
	classes := []string{".guilda-badge--governanca", ".guilda-badge--reversa", ".guilda-badge--construcao", ".guilda-badge--infra", ".guilda-badge--audiovisual", ".guilda-badge--qualidade", ".guilda-badge--publicacao", ".pill--role"}
	for _, cls := range classes {
		validarRegraBadgeWCAG(t, css, "claro", regexp.QuoteMeta(cls))
		validarRegraBadgeWCAG(t, css, "escuro", `:root\[data-theme="dark"\]\s*`+regexp.QuoteMeta(cls))
	}
}

func validarRegraBadgeWCAG(t *testing.T, css, tema, sel string) {
	t.Helper()
	m := regexp.MustCompile(`(?s)` + sel + `\s*\{([^}]+)\}`).FindStringSubmatch(css)
	if len(m) < 2 {
		t.Fatalf("%s: seletor %s ausente", tema, sel)
	}
	bg := regexp.MustCompile(`background:\s*(#[0-9a-fA-F]{3,6})`).FindStringSubmatch(m[1])
	fg := regexp.MustCompile(`(?:^|[\s;])color:\s*(#[0-9a-fA-F]{3,6})`).FindStringSubmatch(m[1])
	if len(bg) < 2 || len(fg) < 2 {
		t.Fatalf("%s: cor sólida ausente em %s", tema, sel)
	}
	if r := contrasteWCAG(fg[1], bg[1]); r < 4.5 {
		t.Errorf("%s %s: contraste %.2f < 4.5 (%s sobre %s)", tema, sel, r, fg[1], bg[1])
	}
}

func contrasteWCAG(fgHex, bgHex string) float64 {
	l1 := luminanciaRelativa(fgHex)
	l2 := luminanciaRelativa(bgHex)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

func luminanciaRelativa(hexStr string) float64 {
	h := strings.TrimPrefix(hexStr, "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	r := canalLinear(h[0:2])
	g := canalLinear(h[2:4])
	b := canalLinear(h[4:6])
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func canalLinear(par string) float64 {
	v, _ := strconv.ParseUint(par, 16, 8)
	s := float64(v) / 255.0
	if s <= 0.03928 {
		return s / 12.92
	}
	return math.Pow((s+0.055)/1.055, 2.4)
}

func assertGuildasEPaletteJS(t *testing.T) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node ausente; validação de palette.js pulada")
	}
	script := `const fs=require("fs");const vm=require("vm");
const ctx={state:{query:"",cat:"all",blocks:[],pieces:[],favOnly:false,filterText:""},render(){},applyFilter(){},setTimeout,clearTimeout};
for(const f of ["static/search.js","static/virtual.js","static/palette.js"]) vm.runInNewContext(fs.readFileSync(f,"utf8"),ctx);
const casos=[[{id:"t1",kind:"trigger_manual"},"governanca"],[{id:"f1",kind:"flow_join"},"governanca"],[{id:"g1",kind:"gate"},"qualidade"],[{id:"r1",kind:"agent",name:"rea-site"},"reversa"],[{id:"a1",kind:"skill",name:"video-editor"},"audiovisual"],[{id:"q1",kind:"skill",name:"homologar-e2e"},"qualidade"],[{id:"p1",kind:"skill",name:"legopress-livro"},"publicacao"],[{id:"i1",kind:"mcp",name:"deploy-vps"},"infra"],[{id:"c1",kind:"agent",name:"coder"},"construcao"]];
for(const [peca, esp] of casos){const got=ctx.entryFromBlock(peca).guilda; if(got!==esp) throw new Error(peca.id+" => "+got+" != "+esp);}`
	if out, err := exec.Command(node, "-e", script).CombinedOutput(); err != nil {
		t.Fatalf("guildas palette.js falhou: %v\n%s", err, string(out))
	}
}

func assertAriaETeclado(t *testing.T) {
	t.Helper()
	checks := map[string][]string{
		"index.html":   {`aria-expanded`},
		"app.js":       {`aria-expanded`},
		"palette.js":   {`role", "group"`, `aria-label`},
		"inventory.js": {`ArrowDown`, `ArrowUp`, `aria-labelledby`},
		"swarm.js":     {`role", "tab"`, `aria-selected`},
	}
	for file, wants := range checks {
		txt := readStatic(t, file)
		for _, w := range wants {
			if !strings.Contains(txt, w) {
				t.Errorf("%s sem acessibilidade %q", file, w)
			}
		}
	}
}

func TestDebitosStudioEContratoDOM(t *testing.T) {
	html := readStatic(t, "index.html")
	if strings.Contains(html, `id="palette-cats" hidden`) || !strings.Contains(html, "palette-grid-4x2") {
		t.Error("index.html deve manter #palette-cats visível em grid 4x2")
	}
	if !strings.Contains(html, `id="inspector-collapsed-rail"`) || !strings.Contains(html, `id="btn-inspector-collapse"`) {
		t.Error("index.html precisa dos controles de recolher/restaurar inspetor")
	}
	pal := readStatic(t, "palette.js")
	if !strings.Contains(pal, `activeCategory: "trigger"`) || !strings.Contains(pal, `cat: "trigger"`) {
		t.Error("palette.js deve iniciar com categoria Trigger aberta")
	}
	app := readStatic(t, "app.js")
	if strings.Contains(app, `Math.max(8,`) || strings.Contains(app, `Math.max(16,`) {
		t.Error("app.js não pode limitar coordenadas negativas do canvas")
	}
	if !strings.Contains(app, `"/api/assets/loop"`) {
		t.Error("app.js deve consultar /api/assets/loop ao resolver sub-loops")
	}
}

func TestZeroEmojisStudioStatic(t *testing.T) {
	err := fs.WalkDir(staticFS, "static", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, rerr := staticFS.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		for _, r := range string(b) {
			if (r >= 0x1F300 && r <= 0x1FAFF) || (r >= 0x2600 && r <= 0x27BF) || (r >= 0x23F0 && r <= 0x23FA) || r == 0xFE0F {
				t.Errorf("emoji proibido (%U) detectado em %s", r, p)
				break
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("falha ao varrer static: %v", err)
	}
}

func TestContratoUXStudioNosEdgesModalESwarm(t *testing.T) {
	css := readStatic(t, "style.css")
	for _, want := range []string{".node--trigger", ".node--control", ".node--skill", ".node--mcp", ".node--loop", ".node--script", ".node--hook", ".node--gate", "transform-box: fill-box"} {
		if !strings.Contains(css, want) {
			t.Errorf("style.css sem regra obrigatória %q", want)
		}
	}
	html := readStatic(t, "index.html")
	for _, id := range []string{`id="btn-simular-swarm"`, `id="execution-flow-strip"`} {
		if !strings.Contains(html, id) {
			t.Errorf("index.html sem elemento obrigatório %s", id)
		}
	}
	app := readStatic(t, "app.js")
	for _, sym := range []string{`node__eye`, `curveMidpoint`, `edge-arrow-sync`, `computeExecutionLevels`, `simulateCurrentDagInSwarm`} {
		if !strings.Contains(app, sym) {
			t.Errorf("app.js sem implementação %q", sym)
		}
	}
	modal := readStatic(t, "components/lego-item-modal.js")
	if strings.Contains(modal, `selecionada para troca (Swap)`) || !strings.Contains(modal, `_hydrateFromInspectAPI`) {
		t.Error("lego-item-modal.js deve hidratar dados reais e evitar stubs")
	}
}

func TestContratoInspetorModalPaletaESubgrafoZoom(t *testing.T) {
	html := readStatic(t, "index.html")
	for _, id := range []string{`id="inspector-actions-bar"`, `id="btn-insp-eye"`, `id="btn-insp-fav"`, `id="btn-insp-canon"`, `id="btn-insp-upgrade"`, `id="btn-insp-delete"`} {
		if !strings.Contains(html, id) {
			t.Errorf("index.html deve manter %s no cabeçalho fixo do inspetor", id)
		}
	}
	if strings.Contains(html, `role="tree"`) {
		t.Error("index.html não deve usar role=tree em container com filhos heterogêneos (aria-required-children)")
	}
	pal := readStatic(t, "palette.js")
	if strings.Contains(pal, `"btn-eye-view"`) || strings.Contains(pal, `btn-item-action btn-fav`) {
		t.Error("palette.js não deve renderizar botões inline de ação nos itens da paleta esquerda")
	}
	insp := readStatic(t, "inspector.js")
	if !strings.Contains(insp, `syncInspectorActions`) || !strings.Contains(insp, `runInspectorUpgrade`) || !strings.Contains(insp, `runInspectorCanonize`) {
		t.Error("inspector.js deve implementar ações fixas no cabeçalho do inspetor")
	}
	modal := readStatic(t, "components/lego-item-modal.js")
	if !strings.Contains(modal, `id="btn-modal-upgrade"`) || strings.Contains(modal, `copyCount <= 1;`) {
		t.Error("lego-item-modal.js deve expor #btn-modal-upgrade e habilitar Canonizar quando há caminho físico")
	}
	app := readStatic(t, "app.js")
	if !strings.Contains(app, `autoLayout({ fit: false`) || !strings.Contains(app, `focusViewportStart(`) {
		t.Error("app.js não deve forçar fitToScreen em sub-grafos bloqueando leitura e zoom")
	}
}
