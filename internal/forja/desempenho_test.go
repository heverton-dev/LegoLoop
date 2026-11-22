package forja

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDesempenhoDryRun valida AC-NF.3: dry-run com 24 harnesses em ate 30 s.
func TestDesempenhoDryRun(t *testing.T) {
	home := t.TempDir()
	harnesses := make([]string, 0, 24)
	for i := 1; i <= 24; i++ {
		nome := fmt.Sprintf("harness-%02d", i)
		harnesses = append(harnesses, nome)
		dir := filepath.Join(home, "."+nome)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	inicio := time.Now()
	rel, err := Executar(context.Background(), Opcao{
		Home:      home,
		Perfil:    "padrao",
		Aplicar:   false,
		Harnesses: harnesses,
	})
	duracao := time.Since(inicio)
	if err != nil {
		t.Fatalf("Executar dry-run falhou: %v", err)
	}
	if duracao > 30*time.Second {
		t.Fatalf("dry-run com 24 harnesses excedeu 30s: %v", duracao)
	}
	if len(rel.Etapas) < 24 {
		t.Fatalf("esperava ao menos 24 etapas, obtido %d", len(rel.Etapas))
	}
}
