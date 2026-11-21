package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestNomeHostilEhRecusado prova AC-E03.3: os 6 verbos que recebem nome/slug
// de Processo recusam nomes hostis com exit 2 e não gravam nada fora do Cofre.
func TestNomeHostilEhRecusado(t *testing.T) {
	bin := buildBinary(t)
	hostis := []string{"../x", "a/b", `..\x`}
	for _, hostil := range hostis {
		for _, caso := range casosNomeHostil(hostil) {
			t.Run(caso.verbo+"/"+hostil, func(t *testing.T) {
				verificarRecusaNomeHostil(t, bin, caso)
			})
		}
	}
}

type casoHostil struct {
	verbo string
	args  func(vault, proj, procDir string) []string
}

func casosNomeHostil(hostil string) []casoHostil {
	return []casoHostil{
		{verbo: "fabricar", args: func(v, p, _ string) []string {
			return []string{"fabricar", hostil, "--yes", "--sem-entrevista", "--vault", v, "--docs", p}
		}},
		{verbo: "compor", args: func(v, p, _ string) []string {
			return []string{"compor", hostil, "legoloop-alvo@1.0.0", "--vault", v, "--docs", p, "--yes"}
		}},
		{verbo: "remover", args: func(v, _, _ string) []string {
			return []string{"remover", hostil, "--vault", v, "--yes"}
		}},
		{verbo: "refinar", args: func(v, p, _ string) []string {
			return []string{"refinar", hostil, "--etapa", "e1", "--titulo", "x", "--vault", v, "--docs", p, "--yes"}
		}},
		{verbo: "exec", args: func(_, p, d string) []string {
			return []string{"exec", hostil, "--processos-dir", d, "--workdir", p}
		}},
		{verbo: "sessoes", args: func(_, _, d string) []string {
			return []string{"sessoes", hostil, "--processos-dir", d}
		}},
	}
}

func verificarRecusaNomeHostil(t *testing.T, bin string, caso casoHostil) {
	t.Helper()
	root := t.TempDir()
	vault := filepath.Join(root, "vault")
	proj := filepath.Join(root, "proj")
	procDir := filepath.Join(root, "processos")
	for _, dir := range []string{vault, proj, procDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	code, out, errb := rodarBinarioHermetico(t, bin, root, proj, caso.args(vault, proj, procDir)...)
	if code != 2 {
		t.Fatalf("%s: exit = %d; quero 2\nstdout: %s\nstderr: %s", caso.verbo, code, out, errb)
	}
	if _, err := os.Stat(filepath.Join(root, "x")); !os.IsNotExist(err) {
		t.Fatalf("%s gravou artefato fora do Cofre em %s", caso.verbo, filepath.Join(root, "x"))
	}
	if entries, err := os.ReadDir(vault); err != nil || len(entries) != 0 {
		t.Fatalf("%s alterou o Cofre: entries=%d err=%v", caso.verbo, len(entries), err)
	}
}

func rodarBinarioHermetico(t *testing.T, bin, root, workdir string, args ...string) (int, string, string) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = workdir
	cmd.Env = append(os.Environ(), "HOME="+root, "USERPROFILE="+root, "NINEROUTER_URL=http://127.0.0.1:1")
	var so, se bytes.Buffer
	cmd.Stdout, cmd.Stderr = &so, &se
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatalf("executar binário falhou: %v", err)
	}
	return code, so.String(), se.String()
}
