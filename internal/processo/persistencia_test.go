package processo_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"legoloop/internal/processo"
)

// TestPersistenciaIdaEVolta: o processo gravado em loops/<nome>.processo.json volta igual e válido.
func TestPersistenciaIdaEVolta(t *testing.T) {
	p := quatro()
	data, err := p.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	path := processo.CaminhoNoCofre(t.TempDir(), p.Nome)
	if filepath.Base(path) != p.Nome+".processo.json" || filepath.Base(filepath.Dir(path)) != "loops" {
		t.Fatalf("caminho no cofre errado: %s", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	lido, err := processo.Ler(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(lido, p) {
		t.Errorf("o processo lido difere do gravado:\n%+v\n%+v", lido, p)
	}
}

// TestPersistenciaRecusaInvalido: processo inválido não é serializado, e arquivo quebrado ou inválido não é lido.
func TestPersistenciaRecusaInvalido(t *testing.T) {
	ruim := quatro()
	ruim.Versao = "1.0"
	if _, err := ruim.Marshal(); !errors.Is(err, processo.ErrInvalido) {
		t.Errorf("Marshal deve recusar versão inválida: %v", err)
	}
	dir := t.TempDir()
	for nome, conteudo := range map[string]string{"quebrado.json": "{", "invalido.json": `{"nome":"x","versao":"1.0.0"}`} {
		path := filepath.Join(dir, nome)
		if err := os.WriteFile(path, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := processo.Ler(path); !errors.Is(err, processo.ErrInvalido) {
			t.Errorf("%s: esperava ErrInvalido, veio %v", nome, err)
		}
	}
}
