package orchestrix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

// SEC-01: nós script só rodam com consentimento explícito.
func TestScriptRequiresConsent(t *testing.T) {
	n := dag.Node{ID: "s", Type: dag.TypeScript, Command: []string{"go", "version"},
		Output: dag.HandoffContract{Artifact: "x.txt", Gate: "true"}}
	_, err := (&Orchestrix{}).Execute(context.Background(), n, t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "--allow-scripts") {
		t.Fatalf("script sem consentimento executou: %v", err)
	}
}

// SEC-04: gravar o artefato não pode atravessar hardlink para outros projetos.
func TestArtifactWriteBreaksHardlinks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"response":"novo"}`))
	}))
	defer srv.Close()
	work, other := t.TempDir(), t.TempDir()
	shared := filepath.Join(other, "SPEC.md")
	os.WriteFile(shared, []byte("projeto B"), 0o644)
	os.MkdirAll(filepath.Join(work, "out"), 0o755)
	if err := os.Link(shared, filepath.Join(work, "out", "BRIEFING.md")); err != nil {
		t.Fatalf("hardlink: %v", err)
	}
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"m"}}, HTTP: srv.Client(), OllamaURL: srv.URL, SkillsDir: skillsDir(t)}
	if _, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "m"}, nil), work); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(shared); string(b) != "projeto B" {
		t.Fatalf("escrita vazou pelo hardlink: %q", b)
	}
	if readArtifact(t, work) != "novo" {
		t.Fatal("artefato não gravado")
	}
}

// SEC-06: a credencial do 9Router só vai para a base confiável configurada.
func TestRouterCredentialOnlyForTrustedBase(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer srv.Close()
	o := &Orchestrix{Inventory: provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"c"}},
		HTTP: srv.Client(), RouterKey: "k-teste", RouterKeyURL: "http://127.0.0.1:20128", SkillsDir: skillsDir(t)}
	if _, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "c"}, nil), t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatal("credencial enviada para base não confiável")
	}
}

// SEC-07: saída acima do limite falha rápido em vez de travar o runner.
func TestOutputFloodFailsFast(t *testing.T) {
	o := &Orchestrix{Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		SkillsDir: skillsDir(t), ExtraEnv: []string{"LEGOLOOP_HELPER=flood"}, MaxOutputBytes: 1024}
	done := make(chan error, 1)
	go func() {
		_, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "haiku"}, nil), t.TempDir())
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "limite") {
			t.Fatalf("estouro não reportado: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("runner travou com saída gigante")
	}
}

// R2-1: artefato não pode atravessar symlink/junction commitado no repositório (ex.: out -> .git).
func TestArtifactRefusesLinkedParent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"response":"hook malicioso"}`))
	}))
	defer srv.Close()
	work := t.TempDir()
	target := filepath.Join(work, "alvo")
	os.MkdirAll(target, 0o755)
	if err := os.Symlink("alvo", filepath.Join(work, "out")); err != nil {
		// Windows sem privilégio de symlink: junction via mklink /J não exige admin.
		if out, jerr := exec.Command("cmd", "/c", "mklink", "/J", filepath.Join(work, "out"), target).CombinedOutput(); jerr != nil {
			t.Fatalf("não foi possível criar link de diretório: %v %s", jerr, out)
		}
	}
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"m"}}, HTTP: srv.Client(), OllamaURL: srv.URL, SkillsDir: skillsDir(t)}
	_, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "m"}, nil), work)
	if err == nil {
		t.Fatal("artefato gravado através de diretório linkado")
	}
	if _, statErr := os.Stat(filepath.Join(target, "BRIEFING.md")); !os.IsNotExist(statErr) {
		t.Fatal("arquivo apareceu no destino do link")
	}
}

// R3-3: a recusa acontece ANTES de criar qualquer diretório através do link.
func TestLinkedParentCreatesNothing(t *testing.T) {
	work := t.TempDir()
	target := filepath.Join(work, "alvo")
	os.MkdirAll(target, 0o755)
	if err := os.Symlink("alvo", filepath.Join(work, "out")); err != nil {
		if out, jerr := exec.Command("cmd", "/c", "mklink", "/J", filepath.Join(work, "out"), target).CombinedOutput(); jerr != nil {
			t.Fatalf("link de diretório: %v %s", jerr, out)
		}
	}
	root, err := os.OpenRoot(work)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := writeArtifact(root, "out/sub/X.md", []byte("x")); err == nil {
		t.Fatal("gravação através de link aceita")
	}
	if _, err := os.Stat(filepath.Join(target, "sub")); !os.IsNotExist(err) {
		t.Fatal("diretório criado através do link antes da recusa")
	}
	if err := writeArtifact(root, "novo/a/b/Y.md", []byte("y")); err != nil {
		t.Fatalf("diretórios reais aninhados deveriam ser criados: %v", err)
	}
}
