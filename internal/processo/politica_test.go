package processo_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/processo"
)

// TestPoliticaLazy prova a regra: por padrão só o índice vai para a pasta do harness.
func TestPoliticaLazy(t *testing.T) {
	p := processo.Politica{}
	if p.Vincula("legoloop-fabricar") {
		t.Error("skill comum não pode ser vinculada em modo lazy")
	}
	if !p.Vincula(processo.IndiceSkills) {
		t.Error("a skill-índice sempre é vinculada")
	}
	if !(processo.Politica{Excecoes: map[string]bool{"legoloop-fechar": true}}).Vincula("legoloop-fechar") {
		t.Error("exceção listada deve ser vinculada")
	}
	if !(processo.Politica{Eager: true}).Vincula("qualquer-uma") {
		t.Error("modo eager vincula tudo")
	}
}

func TestLerExcecoes(t *testing.T) {
	home := t.TempDir()
	if got := processo.LerExcecoes(home); len(got) != 0 {
		t.Fatalf("sem arquivo deve voltar vazio: %v", got)
	}
	dir := filepath.Join(home, ".legoloop")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	txt := "# skills que ficam sempre visíveis\nlegoloop-fechar\n\n  legoloop-listar  \n"
	if err := os.WriteFile(filepath.Join(dir, "skills-eager.txt"), []byte(txt), 0o644); err != nil {
		t.Fatal(err)
	}
	got := processo.LerExcecoes(home)
	if len(got) != 2 || !got["legoloop-fechar"] || !got["legoloop-listar"] {
		t.Fatalf("exceções = %v, quer legoloop-fechar e legoloop-listar", got)
	}
}
