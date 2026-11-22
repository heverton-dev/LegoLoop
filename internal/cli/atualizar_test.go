package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/atualizar"
)

func TestCmdAtualizarShaDivergente(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(dir, "legoloop")
	orig := []byte("original-cli")
	if err := os.WriteFile(binPath, orig, 0o755); err != nil {
		t.Fatal(err)
	}

	artefato := atualizar.NomeArtefato(runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/estavel/SHA256SUMS":
			fmt.Fprintf(w, "%s  %s\n", sha256Hex([]byte("esperado")), artefato)
		case "/estavel/" + artefato:
			_, _ = w.Write([]byte("adulterado"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	var out, errBuf bytes.Buffer
	code := Run(context.Background(), []string{
		"atualizar", "--canal", "estavel", "--base-url", srv.URL,
		"--binario", binPath, "--home", home,
	}, Deps{Stdout: &out, Stderr: &errBuf})
	if code == 0 {
		t.Fatalf("esperava exit != 0 para SHA divergente; stdout=%s stderr=%s", out.String(), errBuf.String())
	}
	depois, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(depois, orig) {
		t.Fatalf("binario alterado mesmo com SHA divergente")
	}
}

func TestCmdAtualizarBetaComPrevEJSON(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(filepath.Join(home, ".claude", "skills", "legoloop-skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(dir, "legoloop")
	anterior := []byte("v1-cli")
	novo := []byte("v2-beta-cli")
	if err := os.WriteFile(binPath, anterior, 0o755); err != nil {
		t.Fatal(err)
	}

	artefato := atualizar.NomeArtefato(runtime.GOOS, runtime.GOARCH)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/beta/SHA256SUMS":
			fmt.Fprintf(w, "%s  %s\n", sha256Hex(novo), artefato)
		case "/beta/" + artefato:
			_, _ = w.Write(novo)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	var out, errBuf bytes.Buffer
	code := Run(context.Background(), []string{
		"atualizar", "--canal", "beta", "--base-url", srv.URL,
		"--binario", binPath, "--home", home, "--json",
	}, Deps{Stdout: &out, Stderr: &errBuf})
	if code != 0 {
		t.Fatalf("esperava exit 0, obteve %d; stderr=%s", code, errBuf.String())
	}

	var res atualizar.Resultado
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatalf("JSON invalido: %v\nsaida=%s", err, out.String())
	}
	if !res.Revalidado || res.Canal != "beta" {
		t.Fatalf("resultado inesperado: %+v", res)
	}
	prev, err := os.ReadFile(binPath + ".prev")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prev, anterior) {
		t.Fatalf(".prev nao guarda o binario anterior")
	}
}

func TestCmdAtualizarCanalInvalidoExit2(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), []string{
		"atualizar", "--canal", "invalido", "--base-url", "https://exemplo", "--binario", "x", "--home", t.TempDir(),
	}, Deps{Stdout: &out, Stderr: &errBuf})
	if code != 2 {
		t.Fatalf("esperava exit 2 para canal invalido, obteve %d (stderr=%s)", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "canal invalido") {
		t.Fatalf("esperava mensagem de canal invalido, obteve %q", errBuf.String())
	}
}

func sha256Hex(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
