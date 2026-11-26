package inventario

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func itemComArquivo(t *testing.T, raiz string) Item {
	t.Helper()
	caminho := filepath.Join(raiz, "arquivos", "trilha.mp3")
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte("conteudo-original"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := HashArquivo(caminho)
	if err != nil {
		t.Fatal(err)
	}
	return Item{ID: "trilha-epica", Tipo: "trilha", Arquivo: "arquivos/trilha.mp3", Hash: h}
}

func TestHashArquivoSHA256(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "x.bin")
	if err := os.WriteFile(caminho, []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := HashArquivo(caminho)
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if err != nil || h != want {
		t.Fatalf("HashArquivo = %q, %v; quer %q", h, err, want)
	}
}

func TestHashVerificarIntegro(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.Verificar(itemComArquivo(t, s.Raiz)); err != nil {
		t.Fatalf("Verificar = %v, quer nil", err)
	}
}

// AC-E07.5: arquivo alterado retorna ErrHashDivergente apontando o item.
func TestHashVerificarAlterado(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	it := itemComArquivo(t, s.Raiz)
	if err := os.WriteFile(filepath.Join(s.Raiz, it.Arquivo), []byte("mexido"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := s.Verificar(it)
	if !errors.Is(err, ErrHashDivergente) || !strings.Contains(err.Error(), it.ID) {
		t.Fatalf("Verificar = %v, quer ErrHashDivergente com %q", err, it.ID)
	}
}

// AC-E07.5: arquivo ausente retorna ErrHashDivergente apontando o item.
func TestHashVerificarAusente(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	it := itemComArquivo(t, s.Raiz)
	if err := os.Remove(filepath.Join(s.Raiz, it.Arquivo)); err != nil {
		t.Fatal(err)
	}
	err := s.Verificar(it)
	if !errors.Is(err, ErrHashDivergente) || !strings.Contains(err.Error(), it.ID) {
		t.Fatalf("Verificar = %v, quer ErrHashDivergente com %q", err, it.ID)
	}
}

func TestHashVerificarSemArquivo(t *testing.T) {
	s := Store{Raiz: t.TempDir()}
	if err := s.Verificar(Item{ID: "so-json", Tipo: "pele"}); err != nil {
		t.Fatalf("Verificar sem arquivo = %v, quer nil", err)
	}
}
