package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
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

func writeTestFile(t *testing.T, root, rel, body string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func createTestLink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err == nil {
		return
	}
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Fatalf("falha ao criar symlink/junction: %v (%s)", err, string(out))
	}
}

func TestIsolatedExcludesCredenciaisAndSymlinks(t *testing.T) {
	t.Run("ExcluiCredenciaisDoTemp", func(t *testing.T) {
		workdir := t.TempDir()
		writeTestFile(t, workdir, "main.go", "package main\n")
		writeTestFile(t, workdir, "sub/helper.go", "package sub\n")
		sensiveis := []string{
			".env", ".env.local", ".env-empresa", ".env-full", ".env-pessoal",
			"credential_gemini_backup.json", "sub/credential_prod.json",
			"credentials.json", "cert.pem", "priv.key", ".npmrc", ".netrc",
			"id_rsa", ".aws/credentials", "secret_app.json",
		}
		for _, rel := range sensiveis {
			writeTestFile(t, workdir, rel, "SEGREDO=123\n")
		}
		isoDir, cleanup, err := createIsolationEnv(context.Background(), workdir)
		if err != nil {
			t.Fatalf("createIsolationEnv falhou: %v", err)
		}
		defer cleanup()
		if _, err := os.Stat(filepath.Join(isoDir, "main.go")); err != nil {
			t.Fatalf("main.go deveria ser copiado: %v", err)
		}
		for _, rel := range sensiveis {
			p := filepath.Join(isoDir, filepath.FromSlash(rel))
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				t.Errorf("arquivo sensível %q vazou para ambiente isolado", rel)
			}
		}
	})

	t.Run("RejeitaSymlinkComLstat", func(t *testing.T) {
		outside := t.TempDir()
		writeTestFile(t, outside, "segredo.txt", "TOP_SECRET")

		workdir := t.TempDir()
		writeTestFile(t, workdir, "interno/dado.txt", "OK")
		createTestLink(t, outside, filepath.Join(workdir, "escape"))
		createTestLink(t, filepath.Join(workdir, "interno"), filepath.Join(workdir, "atalho"))
		if _, cleanup, err := createTempCopy(workdir); err == nil || !strings.Contains(err.Error(), "symlink") {
			if cleanup != nil {
				cleanup()
			}
			t.Fatalf("createTempCopy deveria rejeitar symlink/junction com mensagem de symlink, obteve: %v", err)
		}

		dst := t.TempDir()
		if err := injectContextIn(workdir, dst, []string{"atalho/dado.txt"}); err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("injectContextIn deveria rejeitar symlink/junction em context_in, obteve: %v", err)
		}
		if err := copyPath(filepath.Join(workdir, "atalho"), filepath.Join(dst, "atalho")); err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("copyPath deveria rejeitar symlink/junction com erro, obteve: %v", err)
		}
		if err := promoteArtifact(workdir, dst, "atalho/dado.txt"); err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("promoteArtifact deveria rejeitar symlink/junction na origem, obteve: %v", err)
		}

		iso := t.TempDir()
		writeTestFile(t, iso, "atalho/out.txt", "payload")
		if err := promoteArtifact(iso, workdir, "atalho/out.txt"); err == nil || !strings.Contains(err.Error(), "symlink") {
			t.Fatalf("promoteArtifact deveria rejeitar symlink/junction no destino, obteve: %v", err)
		}
	})

	t.Run("FallbackGitSemCommitExcluiCredenciais", func(t *testing.T) {
		workdir := t.TempDir()
		if out, err := exec.Command("git", "init", workdir).CombinedOutput(); err != nil {
			t.Skipf("git indisponível: %v (%s)", err, string(out))
		}
		writeTestFile(t, workdir, "app.go", "package main\n")
		writeTestFile(t, workdir, ".env-empresa", "TOKEN=segredo\n")
		writeTestFile(t, workdir, "credential_gemini_backup.json", `{"key":"segredo"}`)

		isoDir, cleanup, err := createIsolationEnv(context.Background(), workdir)
		if err != nil {
			t.Fatalf("createIsolationEnv falhou no fallback: %v", err)
		}
		defer cleanup()
		if _, err := os.Stat(filepath.Join(isoDir, "app.go")); err != nil {
			t.Fatalf("app.go deveria existir no fallback: %v", err)
		}
		for _, bad := range []string{".env-empresa", "credential_gemini_backup.json"} {
			if _, err := os.Stat(filepath.Join(isoDir, bad)); !os.IsNotExist(err) {
				t.Errorf("credencial %q vazou no fallback de git worktree", bad)
			}
		}
	})

	t.Run("BloqueiaCredenciaisEmContextIn", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()
		writeTestFile(t, src, ".env-empresa", "TOKEN=abc")
		writeTestFile(t, src, "credential_gemini_backup.json", "{}")
		for _, bad := range []string{".env-empresa", "credential_gemini_backup.json", "."} {
			if err := injectContextIn(src, dst, []string{bad}); err == nil {
				t.Errorf("injectContextIn deveria bloquear %q", bad)
			}
		}
	})

	t.Run("PropagaErrosCopyPathEPromoteArtifact", func(t *testing.T) {
		src := t.TempDir()
		dst := t.TempDir()
		err := copyPath(filepath.Join(src, "nao-existe.txt"), filepath.Join(dst, "nao-existe.txt"))
		if err == nil || !os.IsNotExist(err) {
			t.Fatalf("copyPath deveria propagar erro de arquivo ausente, obteve: %v", err)
		}
		for _, badArt := range []string{"../fora.txt", "C:drive.txt", ".env-vazamento", "credential_x.json", "ausente.txt"} {
			if err := promoteArtifact(src, dst, badArt); err == nil {
				t.Errorf("promoteArtifact deveria falhar para %q", badArt)
			}
		}
		writeTestFile(t, src, "sub/ok.txt", "ok")
		if err := promoteArtifact(src, dst, "sub"); err == nil || !strings.Contains(err.Error(), "diretório") {
			t.Errorf("promoteArtifact deveria falhar para diretório, obteve: %v", err)
		}
	})
}
