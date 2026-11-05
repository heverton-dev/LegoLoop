package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// buildBinary compila o binário real uma vez por teste. O teste roda o executável, não a função Run:
// é a única forma de provar que o main liga o stdin (defeito A1 da auditoria).
func buildBinary(t *testing.T) string {
	t.Helper()
	name := "legoloop"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(t.TempDir(), name)
	out, err := exec.Command("go", "build", "-o", bin, "legoloop/cmd/legoloop").CombinedOutput()
	if err != nil {
		t.Fatalf("go build falhou: %v\n%s", err, out)
	}
	return bin
}

// binRouter é um 9Router em loopback: lista um modelo e responde ao chat com a confiança roteirizada.
// Cada resposta repete a pergunta com a confiança do turno; a última se repete depois de esgotadas.
func binRouter(t *testing.T, confidences ...int) *httptest.Server {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/models":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]string{{"id": "code-fast"}}})
		case "/v1/chat/completions":
			i := min(int(calls.Add(1))-1, len(confidences)-1)
			reply := "HIPÓTESE: extrator\nCONFIANÇA: " + binItoa(confidences[i]) + "%\nPERGUNTA: Qual o formato?\nPALPITE: mp3"
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": reply}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func binItoa(n int) string { b, _ := json.Marshal(n); return string(b) }

type binRun struct {
	out, errb       string
	code            int
	vault, docsRoot string
}

// runBinary executa o binário com HOME, cofre e projeto descartáveis e o stdin dado.
func runBinary(t *testing.T, bin string, router *httptest.Server, stdin string, args ...string) binRun {
	t.Helper()
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	full := append([]string{"fabricar", "extrair audio", "--goal", "extrair audio de video", "--vault", vault, "--docs", proj}, args...)
	return runBinaryStdin(t, bin, router, root, proj, stdin, full...)
}

// runBinaryEm executa o binário com os argumentos dados, HOME em root e diretório de trabalho em proj, sem stdin.
func runBinaryEm(t *testing.T, bin string, router *httptest.Server, root, proj string, args ...string) binRun {
	t.Helper()
	return runBinaryStdin(t, bin, router, root, proj, "", args...)
}

func runBinaryStdin(t *testing.T, bin string, router *httptest.Server, root, proj, stdin string, args ...string) binRun {
	t.Helper()
	r := binRun{vault: filepath.Join(root, "vault"), docsRoot: proj}
	_ = os.MkdirAll(proj, 0o755)
	cmd := exec.Command(bin, args...)
	cmd.Dir = proj
	cmd.Env = append(os.Environ(), "NINEROUTER_URL="+router.URL, "HOME="+root, "USERPROFILE="+root)
	cmd.Stdin = strings.NewReader(stdin)
	var so, se bytes.Buffer
	cmd.Stdout, cmd.Stderr = &so, &se
	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if ee, ok := err.(*exec.ExitError); ok {
			r.code = ee.ExitCode()
		}
	case <-time.After(90 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("o binário não terminou em 90 s\nstdout:\n%s\nstderr:\n%s", so.String(), se.String())
	}
	r.out, r.errb = so.String(), se.String()
	return r
}

func filesUnder(root string) []string {
	var out []string
	_ = filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err == nil && !i.IsDir() {
			out = append(out, p)
		}
		return nil
	})
	return out
}

// TestFabricarBinaryInterativo prova AC-E02.1 pelo executável: stdin em pipe, sem --yes.
func TestFabricarBinaryInterativo(t *testing.T) {
	bin := buildBinary(t)

	t.Run("a entrevista começa: o erro de terminal não aparece", func(t *testing.T) {
		r := runBinary(t, bin, binRouter(t, 40), "0\n")
		all := r.out + r.errb
		if strings.Contains(all, "exige um terminal") {
			t.Fatalf("o main não ligou o stdin:\n%s", all)
		}
		if !strings.Contains(r.out, "Qual o formato?") {
			t.Errorf("a pergunta da entrevista não apareceu:\n%s", all)
		}
	})

	t.Run("a confiança é impressa a cada turno, em ordem", func(t *testing.T) {
		r := runBinary(t, bin, binRouter(t, 40, 70, 90), "1\n1\n0\n")
		last := -1
		for _, want := range []string{"Confiança: 40%", "Confiança: 70%", "Confiança: 90%"} {
			i := strings.Index(r.out, want)
			if i < 0 || i < last {
				t.Fatalf("faltou %q, ou fora de ordem:\n%s", want, r.out)
			}
			last = i
		}
	})

	t.Run("cancelar na entrevista não grava nada", func(t *testing.T) {
		r := runBinary(t, bin, binRouter(t, 40), "0\n")
		if r.code == 0 {
			t.Errorf("cancelar deveria sair com código diferente de 0\n%s", r.out)
		}
		if got := append(filesUnder(r.vault), filesUnder(r.docsRoot)...); len(got) != 1 || filepath.Base(got[0]) != "registry_matrix.json" {
			t.Errorf("gravou antes da confirmação: %v", got)
		}
	})
}
