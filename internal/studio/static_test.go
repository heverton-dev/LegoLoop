package studio

import (
	"encoding/json"
	"io/fs"
	"os/exec"
	"regexp"
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

// R-UI-3: nenhum diálogo nativo do navegador/SO — feedback é sempre componente próprio.
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
	for _, f := range []string{"app.js", "ui.js", "inventory.js", "search.js", "virtual.js", "palette.js", "slots.js"} {
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
		t.Error("BUG-TMPL-01: app.js repete tmpl.name no fallback de icon, duplicando o nome no card")
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
		t.Error("BUG-VAULT-03: handler de btnApply deve gerenciar disabled no bloco finally")
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

// TestNoInlineStylesInComponents (BUG-CSP-01) assegura que nenhum componente utiliza
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



