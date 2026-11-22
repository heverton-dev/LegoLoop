package atualizar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestShaDivergenteMantemBinario(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(dir, "legoloop-bin")
	conteudoOriginal := []byte("binario-original-v1")
	if err := os.WriteFile(binPath, conteudoOriginal, 0o755); err != nil {
		t.Fatal(err)
	}
	hashAntes := hashBytes(conteudoOriginal)

	artefato := NomeArtefato(runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/estavel/SHA256SUMS":
			hashEsperado := hashBytes([]byte("binario-esperado-v2"))
			fmt.Fprintf(w, "%s  %s\n", hashEsperado, artefato)
		case "/estavel/" + artefato:
			_, _ = w.Write([]byte("binario-adulterado-v2"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	_, err := Executar(context.Background(), Opcao{
		Canal:   CanalEstavel,
		BaseURL: srv.URL,
		Binario: binPath,
		Home:    home,
		Cliente: srv.Client(),
	})
	if !errors.Is(err, ErrShaDivergente) {
		t.Fatalf("esperava ErrShaDivergente, obteve %v", err)
	}

	atual, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatalf("falha ao ler binario atual: %v", err)
	}
	if hashBytes(atual) != hashAntes {
		t.Fatalf("binario alterado apos SHA divergente: antes=%s depois=%s", hashAntes, hashBytes(atual))
	}
	if _, err := os.Stat(binPath + ".prev"); !os.IsNotExist(err) {
		t.Fatalf("nao deveria criar .prev quando SHA diverge")
	}
}

func TestCanalInvalidoRejeitado(t *testing.T) {
	dir := t.TempDir()
	_, err := Executar(context.Background(), Opcao{
		Canal:   "nightly",
		BaseURL: "https://exemplo.invalido",
		Binario: filepath.Join(dir, "legoloop"),
		Home:    dir,
	})
	if !errors.Is(err, ErrCanalInvalido) {
		t.Fatalf("esperava ErrCanalInvalido, obteve %v", err)
	}
}

func hashBytes(dados []byte) string {
	soma := sha256.Sum256(dados)
	return hex.EncodeToString(soma[:])
}
