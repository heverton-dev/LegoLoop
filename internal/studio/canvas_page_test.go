package studio

import (
	"regexp"
	"strings"
	"testing"
)

// TICK-905/907/908/910: a página do canvas carrega os scripts na ordem das dependências e expõe a
// paleta (busca, categorias, lista virtualizada) numa lateral fixa.
func TestCanvasPageWiring(t *testing.T) {
	html := readStatic(t, "index.html")
	order := []string{"ui.js", "search.js", "virtual.js", "palette.js", "slots.js", "inspector.js", "app.js"}
	last := -1
	for _, f := range order {
		i := strings.Index(html, `src="`+f+`"`)
		if i < 0 {
			t.Fatalf("index.html não carrega %s", f)
		}
		if i < last {
			t.Errorf("%s fora de ordem: depende dos anteriores (%v)", f, order)
		}
		last = i
	}
	for _, id := range []string{"palette-search", "palette-cats", "palette-list", "palette-spacer", "palette-rows",
		"palette-empty", "palette-empty-title", "palette-empty-text", "palette-banner", "palette-hint", "inventory"} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Errorf("index.html sem id=%q", id)
		}
	}
	if !regexp.MustCompile(`id="palette-search"[^>]*type="search"|type="search"[^>]*id="palette-search"`).MatchString(html) {
		t.Error("a busca da paleta deve ser input type=search")
	}
}

func TestCanvasScriptsAreServed(t *testing.T) {
	s, _ := newServer(t)
	for _, f := range []string{"/virtual.js", "/palette.js", "/slots.js", "/inspector.js"} {
		if r := do(t, s.Handler(), "GET", f, ""); r.Code != 200 {
			t.Errorf("%s: %d", f, r.Code)
		}
	}
}

// R-UI-8: laterais fixas, só a lista interna rola; a página nunca ganha rolagem própria.
func TestSidebarsAreFixedWithInternalScroll(t *testing.T) {
	css := readStatic(t, "style.css")
	rule := func(selector string) string {
		m := regexp.MustCompile(`(?ms)^` + regexp.QuoteMeta(selector) + `\s*\{([^}]*)\}`).FindStringSubmatch(css)
		if m == nil {
			t.Fatalf("regra %s ausente no CSS", selector)
		}
		return m[1]
	}
	if r := rule(".sidebar"); !strings.Contains(r, "overflow: hidden") || !strings.Contains(r, "flex-direction: column") {
		t.Errorf(".sidebar deve ser coluna fixa sem rolagem própria: %s", r)
	}
	if r := rule(".palette-list"); !strings.Contains(r, "overflow-y: auto") || !strings.Contains(r, "flex: 1") || !strings.Contains(r, "min-height: 0") {
		t.Errorf(".palette-list deve ser a área que rola: %s", r)
	}
	if r := rule(".inspector"); !strings.Contains(r, "overflow-y: auto") {
		t.Errorf(".inspector deve rolar por dentro: %s", r)
	}
	if r := rule("html, body"); !strings.Contains(r, "overflow: hidden") {
		t.Errorf("a página não pode ter rolagem própria: %s", r)
	}
	if r := rule(".palette__item, .palette__header"); !strings.Contains(r, "height: 48px") {
		t.Errorf("as linhas da paleta precisam de altura fixa (virtualização): %s", r)
	}
}
