package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestKraftenColetaInterativa prova que parâmetros ausentes são perguntados um a
// um via stdin (3 vias: terminal, agente/chat e Studio alimentam o mesmo fluxo).
func TestKraftenColetaInterativa(t *testing.T) {
	zipPath := zipSaudavelKraften(t)
	entradas := strings.Join([]string{
		zipPath,
		"app01",
		"sub01",
		"server.exemplo.org",
		"exemplo.org",
		"",
		"",
	}, "\n")
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdin: strings.NewReader(entradas), Stdout: &stdout, Stderr: &stderr}
	code := cmdKraftenDeploy(context.Background(), []string{
		"preflight", "--dest", t.TempDir(),
	}, deps)
	if code != 0 {
		t.Fatalf("coleta interativa deve concluir com exit 0, obteve %d (stderr: %s)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "--zip") {
		t.Fatalf("esperava prompt do parâmetro --zip na coleta: %s", stderr.String())
	}
}

// escreverArquivo grava um artefato de etapa no diretório de destino.
func escreverArquivo(t *testing.T, dir, nome, corpo string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, nome), []byte(corpo), 0o644); err != nil {
		t.Fatalf("gravar %s: %v", nome, err)
	}
}

// TestKraftenColetaPedeUmPorVez prova a ordem e o número de perguntas.
func TestKraftenColetaPedeUmPorVez(t *testing.T) {
	var prompts []string
	perguntar = func(rotulo string, d Deps) (string, error) {
		prompts = append(prompts, rotulo)
		return "x", nil
	}
	defer func() { perguntar = perguntarStdin }()
	dest := t.TempDir()
	escreverArquivo(t, dest, "deploy.json", `{"status":"done"}`)
	escreverArquivo(t, dest, "gates.json", `{"cf_old":0,"cf_new":1,"tls_verdict":0,"http_code":200}`)
	code := cmdKraftenDeploy(context.Background(), []string{
		"receipt", "--dest", dest,
	}, Deps{Stdin: strings.NewReader(""), Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}})
	if code != 0 {
		t.Fatalf("receipt com coleta deve concluir, obteve %d", code)
	}
	if len(prompts) != 1 || prompts[0] != "--app" {
		t.Fatalf("esperava exatamente o prompt --app, obteve %v", prompts)
	}
}
