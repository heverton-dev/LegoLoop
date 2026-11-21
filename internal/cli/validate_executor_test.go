package cli

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func buildBinarioTeste(t *testing.T) string {
	t.Helper()
	nome := "legoloop"
	if runtime.GOOS == "windows" {
		nome += ".exe"
	}
	bin := filepath.Join(t.TempDir(), nome)
	cmd := exec.Command("go", "build", "-o", bin, "../../cmd/legoloop")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha ao compilar binário de teste: %v\n%s", err, out)
	}
	return bin
}

func gravarDAGIncoerente(t *testing.T, dir string) string {
	t.Helper()
	conteudo := `{
  "name": "dag-executor-incoerente",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "no-stop-errado",
      "type": "stop_gate",
      "isolation": {"executor": "orchestrix_hermetic"},
      "routing": {"mode": "manual", "model": "m1"},
      "output": {"artifact": "saida-stop.md", "gate": "file_exists('saida-stop.md')"}
    },
    {
      "id": "no-script-errado",
      "type": "script",
      "command": ["echo", "ok"],
      "isolation": {"executor": "local_process"},
      "routing": {"mode": "manual", "model": "m1"},
      "output": {"artifact": "saida-script.md", "gate": "file_exists('saida-script.md')"}
    }
  ]
}`
	arq := filepath.Join(dir, "incoerente.dag.json")
	if err := os.WriteFile(arq, []byte(conteudo), 0o600); err != nil {
		t.Fatalf("falha ao gravar DAG incoerente: %v", err)
	}
	return arq
}

func executarBinario(bin string, args ...string) (int, string) {
	cmd := exec.Command(bin, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err == nil {
		return 0, buf.String()
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), buf.String()
	}
	return -1, buf.String() + "\n" + err.Error()
}

func verificarLinhasIncoerentes(t *testing.T, subcmd, saida string, ids []string) {
	t.Helper()
	linhas := strings.Split(saida, "\n")
	for _, id := range ids {
		encontrou := false
		for _, linha := range linhas {
			if strings.Contains(linha, id) && strings.Contains(linha, "incoerente") {
				encontrou = true
				break
			}
		}
		if !encontrou {
			t.Errorf("%s: saída não contém linha de incoerência para nó %q:\n%s", subcmd, id, saida)
		}
	}
}

func TestValidateERunRecusamIncoerenteComExit2(t *testing.T) {
	bin := buildBinarioTeste(t)
	dagPath := gravarDAGIncoerente(t, t.TempDir())
	nosEsperados := []string{"no-stop-errado", "no-script-errado"}

	casos := []struct {
		nome string
		args []string
	}{
		{nome: "validate", args: []string{"validate", dagPath}},
		{nome: "run", args: []string{"run", dagPath, "--allow-scripts", "--workdir", t.TempDir()}},
	}

	for _, tc := range casos {
		t.Run(tc.nome, func(t *testing.T) {
			code, saida := executarBinario(bin, tc.args...)
			if code != 2 {
				t.Errorf("%s: esperado exit code 2, obteve %d\nsaída:\n%s", tc.nome, code, saida)
			}
			verificarLinhasIncoerentes(t, tc.nome, saida, nosEsperados)
		})
	}
}
