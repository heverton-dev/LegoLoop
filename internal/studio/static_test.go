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

