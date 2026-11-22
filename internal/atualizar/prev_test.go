package atualizar

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGuardaPrevERevalida(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(home, ".claude", "skills", "legoloop-skills"), 0o755); err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(dir, "legoloop-bin")
	binAnterior := []byte("legoloop-v1-anterior")
	binNovo := []byte("legoloop-v2-beta-novo")
	if err := os.WriteFile(binPath, binAnterior, 0o755); err != nil {
		t.Fatal(err)
	}

	artefato := NomeArtefato(runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/beta/SHA256SUMS":
			fmt.Fprintf(w, "%s  %s\n", hashBytes(binNovo), artefato)
		case "/beta/" + artefato:
			_, _ = w.Write(binNovo)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	res, err := Executar(context.Background(), Opcao{
		Canal:   CanalBeta,
		BaseURL: srv.URL,
		Binario: binPath,
		Home:    home,
		Cliente: srv.Client(),
	})
	if err != nil {
		t.Fatalf("Executar falhou: %v", err)
	}
	if !res.Revalidado {
		t.Fatalf("esperava Revalidado=true")
	}
	if res.BinarioPrev != binPath+".prev" {
		t.Fatalf("BinarioPrev inesperado: %q", res.BinarioPrev)
	}

	atual, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(atual, binNovo) {
		t.Fatalf("binario novo nao gravado corretamente")
	}

	prev, err := os.ReadFile(binPath + ".prev")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prev, binAnterior) {
		t.Fatalf(".prev nao preservou o binario anterior byte a byte")
	}
}

func TestRevalidacaoFalhaComDeriva(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(home, ".claude", "skills", "skill-intrusa"), 0o755); err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(dir, "legoloop-bin")
	binAnterior := []byte("legoloop-v1")
	binNovo := []byte("legoloop-v2")
	if err := os.WriteFile(binPath, binAnterior, 0o755); err != nil {
		t.Fatal(err)
	}

	artefato := NomeArtefato(runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/beta/SHA256SUMS":
			fmt.Fprintf(w, "%s  %s\n", hashBytes(binNovo), artefato)
		case "/beta/" + artefato:
			_, _ = w.Write(binNovo)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	_, err := Executar(context.Background(), Opcao{
		Canal:   CanalBeta,
		BaseURL: srv.URL,
		Binario: binPath,
		Home:    home,
		Cliente: srv.Client(),
	})
	if !errors.Is(err, ErrDerivaDetectada) {
		t.Fatalf("esperava ErrDerivaDetectada, obteve %v", err)
	}
}
