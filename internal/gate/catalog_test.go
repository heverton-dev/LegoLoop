package gate

import (
	"sort"
	"strings"
	"testing"
)

// O catálogo é a fonte do inventário de gates: cobre exatamente os builtins, sem sobra nem falta.
func TestCatalogCoversEveryBuiltin(t *testing.T) {
	cat := Catalog()
	var got, want []string
	for _, c := range cat {
		got = append(got, c.Name)
		if !strings.HasPrefix(c.Signature, c.Name+"(") || strings.TrimSpace(c.Description) == "" {
			t.Errorf("entrada incompleta: %+v", c)
		}
	}
	for name := range builtins() {
		want = append(want, name)
	}
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("catálogo = %v, builtins = %v", got, want)
	}
}

func TestCatalogSignaturesMatchArity(t *testing.T) {
	for _, c := range Catalog() {
		b := builtins()[c.Name]
		params := strings.Count(c.Signature, ",") + 1
		if len(b.args) == 0 {
			t.Fatalf("%s sem argumentos?", c.Name)
		}
		if params != len(b.args) {
			t.Errorf("%s: assinatura %q tem %d parâmetros, builtin exige %d", c.Name, c.Signature, params, len(b.args))
		}
	}
}
