package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"legoloop/internal/dag"
)

type leakyExec struct {
	recordedWorkdir string
}

func (l *leakyExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	l.recordedWorkdir = workdir

	// Grava o artefato declarado
	artPath := filepath.Join(workdir, n.Output.Artifact)
	if err := os.WriteFile(artPath, []byte("conteudo do artefato"), 0o644); err != nil {
		return Result{}, err
	}

	// Grava um arquivo extra (vazamento / lixo que NÃO deve ir para o workspace do projeto)
	leakPath := filepath.Join(workdir, "leak.tmp")
	if err := os.WriteFile(leakPath, []byte("arquivo lixo"), 0o644); err != nil {
		return Result{}, err
	}

	return Result{
		ExitCode: 0,
		Tokens:   TokenUsage{Input: 10, Output: 20},
		Model:    "test/isolated",
	}, nil
}

func TestRunIsolatedPromotesOnlyDeclaredArtifact(t *testing.T) {
	workdir := t.TempDir()

	src := `{"name":"test-iso","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"node-iso",
		"type":"script",
		"command":["echo","hello"],
		"execution":{"mode":"isolated"},
		"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"declared.md","gate":"file_exists('declared.md') && exit_code == 0"}
	}]}`

	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("falha ao parsear dag: %v", err)
	}

	exec := &leakyExec{}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  workdir,
		Executor: exec,
	})
	if err != nil {
		t.Fatalf("Run falhou inesperadamente: %v", err)
	}

	if !rep.OK {
		t.Fatalf("relatório reportou falha: %+v", rep)
	}

	// O executor não pode ter rodado diretamente no workspace principal
	if exec.recordedWorkdir == workdir {
		t.Fatalf("executor rodou no workdir principal; deveria rodar em ambiente isolado")
	}

	// 1. O artefato declarado DEVE estar presente no workspace principal
	promotedArtifact := filepath.Join(workdir, "declared.md")
	content, err := os.ReadFile(promotedArtifact)
	if err != nil {
		t.Fatalf("artefato declarado não foi promovido para o workspace: %v", err)
	}
	if string(content) != "conteudo do artefato" {
		t.Fatalf("conteúdo do artefato incorreto: %q", string(content))
	}

	// 2. O arquivo lixo/vazado NÃO DEVE existir no workspace principal
	leakInMain := filepath.Join(workdir, "leak.tmp")
	if _, err := os.Stat(leakInMain); !os.IsNotExist(err) {
		t.Fatalf("arquivo não declarado vazou para o workspace principal!")
	}

	// 3. O diretório temporário do isolamento deve ter sido limpo
	if _, err := os.Stat(exec.recordedWorkdir); !os.IsNotExist(err) {
		t.Fatalf("diretório isolado não foi limpo após execução: %s", exec.recordedWorkdir)
	}
}

func TestRunIsolatedWithGitWorktree(t *testing.T) {
	// Cria um diretório temporário configurado como repositório Git
	repoDir := t.TempDir()

	// git init
	cmd := exec.Command("git", "init", repoDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git não disponível ou erro no init: %v (%s)", err, string(out))
	}

	// Configura user/email local para poder comitar
	_ = exec.Command("git", "-C", repoDir, "config", "user.name", "Test").Run()
	_ = exec.Command("git", "-C", repoDir, "config", "user.email", "test@example.com").Run()

	// Cria e comita um arquivo inicial
	readmePath := filepath.Join(repoDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Repo Teste"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = exec.Command("git", "-C", repoDir, "add", "README.md").Run()
	if out, err := exec.Command("git", "-C", repoDir, "commit", "-m", "initial commit").CombinedOutput(); err != nil {
		t.Skipf("não foi possível criar commit inicial: %v (%s)", err, string(out))
	}

	src := `{"name":"test-git-iso","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"node-git-iso",
		"type":"script",
		"command":["echo","git-ok"],
		"execution":{"mode":"isolated"},
		"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"output.txt","gate":"file_exists('output.txt') && exit_code == 0"}
	}]}`

	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("falha ao parsear dag: %v", err)
	}

	execTracker := &leakyExec{}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  repoDir,
		Executor: execTracker,
	})
	if err != nil {
		t.Fatalf("Run falhou com worktree: %v", err)
	}

	if !rep.OK {
		t.Fatalf("relatório reportou erro: %+v", rep)
	}

	// Verifica se artefato declarado foi promovido
	if _, err := os.Stat(filepath.Join(repoDir, "output.txt")); err != nil {
		t.Fatalf("artefato declarado output.txt não promovido: %v", err)
	}

	// Verifica se lixo não vazou
	if _, err := os.Stat(filepath.Join(repoDir, "leak.tmp")); !os.IsNotExist(err) {
		t.Fatalf("arquivo lixo leak.tmp vazou para o repositório principal!")
	}
}

type contextInCheckingExec struct {
	foundInputs []string
}

func (c *contextInCheckingExec) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	for _, in := range n.ContextIn {
		p := filepath.Join(workdir, in)
		if _, err := os.Stat(p); err == nil {
			c.foundInputs = append(c.foundInputs, in)
		}
	}
	artPath := filepath.Join(workdir, n.Output.Artifact)
	_ = os.WriteFile(artPath, []byte("ok"), 0o644)
	return Result{ExitCode: 0, Model: "test/context-check"}, nil
}

func TestRunIsolatedInjectsContextIn(t *testing.T) {
	repoDir := t.TempDir()
	_ = exec.Command("git", "init", repoDir).Run()
	_ = exec.Command("git", "-C", repoDir, "config", "user.name", "Test").Run()
	_ = exec.Command("git", "-C", repoDir, "config", "user.email", "test@example.com").Run()
	_ = os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# init"), 0o644)
	_ = exec.Command("git", "-C", repoDir, "add", "README.md").Run()
	_ = exec.Command("git", "-C", repoDir, "commit", "-m", "init").Run()

	scratchDir := filepath.Join(repoDir, "scratch")
	_ = os.MkdirAll(scratchDir, 0o755)
	untrackedInput := filepath.Join(scratchDir, "untracked_input.json")
	if err := os.WriteFile(untrackedInput, []byte(`{"hello":"world"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	src := `{"name":"test-context-in","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"node-consumer",
		"type":"script",
		"command":["echo","done"],
		"execution":{"mode":"isolated"},
		"context_in":["scratch/untracked_input.json"],
		"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"out.txt","gate":"file_exists('out.txt')"}
	}]}`

	d, err := dag.Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse falhou: %v", err)
	}

	checker := &contextInCheckingExec{}
	rep, err := Run(context.Background(), d, Options{
		Workdir:  repoDir,
		Executor: checker,
	})
	if err != nil {
		t.Fatalf("Run falhou: %v", err)
	}
	if !rep.OK {
		t.Fatalf("rep falhou: %+v", rep)
	}
	if len(checker.foundInputs) != 1 || checker.foundInputs[0] != "scratch/untracked_input.json" {
		t.Fatalf("context_in untracked_input.json NÃO foi injetado no ambiente isolado! Encontrados: %v", checker.foundInputs)
	}
}
