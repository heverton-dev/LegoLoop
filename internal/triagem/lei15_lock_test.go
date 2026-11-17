package triagem

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLei15UmTermoParaLock prova o débito EINT (Lei 15: um termo, um significado):
// o erro de lock tem um único nome — ErrLockOcupado — coerente entre código e as
// especificações em docs/specs/. ErrLockAtivo não pode aparecer em nenhuma SPEC.
func TestLei15UmTermoParaLock(t *testing.T) {
	if ErrLockOcupado.Error() != "triagem: lock ocupado" {
		t.Fatalf("ErrLockOcupado mudou de mensagem: %q", ErrLockOcupado.Error())
	}
	if !errors.Is(ErrLockOcupado, ErrLockOcupado) {
		t.Fatal("ErrLockOcupado não é comparável por errors.Is")
	}
	raiz := filepath.Join("..", "..", "docs", "specs")
	achados, err := filepath.Glob(filepath.Join(raiz, "*.md"))
	if err != nil || len(achados) == 0 {
		t.Fatalf("specs não encontradas em %s: %v", raiz, err)
	}
	for _, spec := range achados {
		dados, err := os.ReadFile(spec)
		if err != nil {
			t.Fatalf("lendo %s: %v", spec, err)
		}
		for i, linha := range strings.Split(string(dados), "\n") {
			if strings.Contains(linha, "ErrLockAtivo") {
				t.Errorf("%s:%d usa o termo ErrLockAtivo; o canônico é ErrLockOcupado", filepath.Base(spec), i+1)
			}
		}
	}
}
