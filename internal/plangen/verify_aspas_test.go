package plangen

import (
	"reflect"
	"strings"
	"testing"
)

// TestSplitVerifyAspasDeArgumento cobre o A21: argumento entre aspas duplas (regex do go test -run)
// é um argv único e não exige gate manual.
func TestSplitVerifyAspasDeArgumento(t *testing.T) {
	argv, ok := splitVerify(`go test -count=1 ./internal/video/inventario/ -run "TestBusca|TestHash"`)
	if !ok {
		t.Fatal("comando com aspas de argumento deveria virar argv")
	}
	want := []string{"go", "test", "-count=1", "./internal/video/inventario/", "-run", "TestBusca|TestHash"}
	if !reflect.DeepEqual(argv[0], want) {
		t.Fatalf("argv = %q, quer %q", argv[0], want)
	}
}

func TestSplitVerifyAspasPerigosasContinuamManuais(t *testing.T) {
	for _, cmd := range []string{
		`go test -run "A;rm -rf x"`,
		`go test -run "$HOME"`,
		`go test -run "A`,
		`go test -run 'A|B'`,
	} {
		if _, ok := splitVerify(cmd); ok {
			t.Errorf("%s deveria exigir gate manual", cmd)
		}
	}
	if _, g := verifyAndGate(`go test -run "A;B"`, nil); !strings.HasPrefix(g, "prompt(") {
		t.Errorf("gate = %s", g)
	}
}

func TestSplitVerifyAspasColadas(t *testing.T) {
	if _, ok := splitVerify(`go test -run "A|B"x`); ok {
		t.Error("texto colado após aspas deveria exigir gate manual")
	}
}
