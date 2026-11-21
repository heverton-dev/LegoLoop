package cli

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func rodarBinarioMatrix(bin string, args ...string) (int, string, string) {
	cmd := exec.Command(bin, args...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err == nil {
		return 0, out.String(), errb.String()
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), out.String(), errb.String()
	}
	return -1, out.String(), err.Error()
}

// TestVerboMatrixRemovido prova o AC-E06.2: o verbo matrix foi removido da CLI,
// retornando exit 2 por comando desconhecido, e não aparece na ajuda do binário.
func TestVerboMatrixRemovido(t *testing.T) {
	bin := buildBinary(t)

	code, _, errb := rodarBinarioMatrix(bin, "matrix")
	if code != 2 {
		t.Fatalf("legoloop matrix: esperava exit 2, obteve %d (stderr=%q)", code, errb)
	}
	if !strings.Contains(errb, `comando desconhecido "matrix"`) {
		t.Errorf("legoloop matrix: esperava mensagem de comando desconhecido, obteve %q", errb)
	}

	_, helpOut, helpErr := rodarBinarioMatrix(bin)
	tudo := helpOut + "\n" + helpErr
	if strings.Contains(tudo, "matrix") {
		t.Errorf("ajuda do legoloop não deve conter 'matrix', obteve:\n%s", tudo)
	}

	codeHelp, outHelp, _ := rodarBinarioMatrix(bin, "--help")
	if codeHelp != 0 {
		t.Fatalf("legoloop --help: esperava exit 0, obteve %d", codeHelp)
	}
	if strings.Contains(outHelp, "matrix") {
		t.Errorf("legoloop --help não deve conter 'matrix', obteve:\n%s", outHelp)
	}
}
