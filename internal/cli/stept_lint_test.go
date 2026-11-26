package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/cli"
)

func rodarLint(t *testing.T, texto string, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	d := cli.Deps{Stdin: strings.NewReader(texto), Stdout: &out, Stderr: &errb, Home: t.TempDir()}
	code := cli.Run(context.Background(), append([]string{"stept", "lint"}, args...), d)
	return code, out.String(), errb.String()
}

// TestSteptLintReprovaEvasao cobre a lei da prova no gancho de parada: texto com evasao sai com 1.
func TestSteptLintReprovaEvasao(t *testing.T) {
	code, _, errb := rodarLint(t, "O teste passou (go test: ok).\nNão verificado: falta rodar o gate.\n")
	if code != 1 {
		t.Fatalf("exit = %d; quer 1", code)
	}
	if !strings.Contains(errb, "evasao-de-prova") {
		t.Errorf("stderr sem a regra: %q", errb)
	}
}

// TestSteptLintAprovaTextoComProva garante exit 0 quando nao ha evasao.
func TestSteptLintAprovaTextoComProva(t *testing.T) {
	code, _, errb := rodarLint(t, "O worker terminou com exit 0 (swarm_run).\n")
	if code != 0 {
		t.Fatalf("exit = %d; quer 0 (stderr %q)", code, errb)
	}
}

// TestSteptLintJSON garante a saida estruturada para ganchos de outros harnesses.
func TestSteptLintJSON(t *testing.T) {
	code, out, _ := rodarLint(t, "Ainda falta validar o schema.\n", "--json")
	if code != 1 || !strings.Contains(out, `"regra":"evasao-de-prova"`) {
		t.Fatalf("exit=%d out=%q", code, out)
	}
}

func rodarGancho(t *testing.T, payload string) (int, string) {
	t.Helper()
	code, _, errb := rodarLint(t, payload, "--hook")
	return code, errb
}

// TestSteptLintGanchoBloqueiaEvasao cobre o gancho de parada: exit 2 bloqueia a resposta no harness.
func TestSteptLintGanchoBloqueiaEvasao(t *testing.T) {
	code, errb := rodarGancho(t, `{"last_assistant_message":"Ainda falta validar o schema."}`)
	if code != 2 || !strings.Contains(errb, "evasao-de-prova") {
		t.Fatalf("exit=%d stderr=%q", code, errb)
	}
}

func TestSteptLintGanchoAprovaTextoComProva(t *testing.T) {
	if code, errb := rodarGancho(t, `{"last_assistant_message":"O teste passou (go test: ok)."}`); code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, errb)
	}
}

// TestSteptLintGanchoLeTranscript cobre o payload sem a mensagem pronta: le o ultimo texto do assistente.
func TestSteptLintGanchoLeTranscript(t *testing.T) {
	arq := filepath.Join(t.TempDir(), "t.jsonl")
	linhas := `{"type":"user","message":{"content":"oi"}}` + "\n" +
		`{"type":"assistant","message":{"content":[{"type":"text","text":"Não verificado."}]}}` + "\n"
	if err := os.WriteFile(arq, []byte(linhas), 0o644); err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]string{"transcript_path": arq})
	if code, errb := rodarGancho(t, string(payload)); code != 2 {
		t.Fatalf("exit=%d stderr=%q", code, errb)
	}
}

// TestSteptLintGanchoPayloadInvalidoNaoQuebra garante que payload malformado nao trava a sessao.
func TestSteptLintGanchoPayloadInvalidoNaoQuebra(t *testing.T) {
	if code, _ := rodarGancho(t, "isto nao e json"); code != 0 {
		t.Fatalf("exit=%d; quer 0", code)
	}
}
