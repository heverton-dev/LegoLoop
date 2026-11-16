package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWebQA_HelpEUsoInvalido(t *testing.T) {
	var out, errb bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errb}

	if code := Run(context.Background(), []string{"web-qa", "--help"}, d); code != 0 {
		t.Fatalf("esperava exit 0 em --help, obteve %d", code)
	}
	if !strings.Contains(out.String(), "web-qa") {
		t.Fatalf("help sem menção a web-qa: %s", out.String())
	}

	for _, args := range [][]string{
		{"web-qa"},
		{"web-qa", "invalido"},
		{"web-qa", "run"},
		{"web-qa", "corrigir"},
	} {
		out.Reset()
		errb.Reset()
		if code := Run(context.Background(), args, d); code != 2 {
			t.Errorf("%v: esperava exit 2, obteve %d (stderr=%s)", args, code, errb.String())
		}
	}
}

func TestWebQA_PreflightJSON(t *testing.T) {
	var gotArgs []string
	var out, errb bytes.Buffer
	d := Deps{
		Stdout: &out,
		Stderr: &errb,
		WebQARunner: func(_ context.Context, args []string, stdout, _ io.Writer) int {
			gotArgs = append([]string(nil), args...)
			_, _ = stdout.Write([]byte(`{"status":"PASSED"}`))
			return 0
		},
	}

	code := Run(context.Background(), []string{"web-qa", "preflight", "--json"}, d)
	if code != 0 {
		t.Fatalf("preflight falhou: exit=%d stderr=%s", code, errb.String())
	}
	if len(gotArgs) != 3 || !strings.HasSuffix(gotArgs[0], filepath.FromSlash("tools/webqa/webqa.mjs")) {
		t.Fatalf("script inesperado: %v", gotArgs)
	}
	if !reflect.DeepEqual(gotArgs[1:], []string{"preflight", "--json"}) {
		t.Fatalf("argumentos encaminhados inesperados: %v", gotArgs[1:])
	}
}

func TestWebQA_RunEncaminhaFlagsESemVerbo(t *testing.T) {
	cases := [][]string{
		{"web-qa", "run", "--project", ".", "--port", "7893", "--alvo", "studio"},
		{"web-qa", "--project", ".", "--port", "7893", "--alvo", "studio"},
	}
	for _, tc := range cases {
		var gotArgs []string
		var out, errb bytes.Buffer
		d := Deps{
			Stdout: &out,
			Stderr: &errb,
			WebQARunner: func(_ context.Context, args []string, _, _ io.Writer) int {
				gotArgs = append([]string(nil), args...)
				return 0
			},
		}
		if code := Run(context.Background(), tc, d); code != 0 {
			t.Fatalf("%v: exit=%d stderr=%s", tc, code, errb.String())
		}
		want := []string{"run", "--project", ".", "--alvo", "studio", "--port", "7893"}
		if !reflect.DeepEqual(gotArgs[1:], want) {
			t.Errorf("%v: encaminhou %v, queria %v", tc, gotArgs[1:], want)
		}
	}
}

func TestWebQA_EncaminhaExitCodeRealDoRunner(t *testing.T) {
	for _, wantCode := range []int{0, 1, 2, 7} {
		var out, errb bytes.Buffer
		d := Deps{
			Stdout: &out,
			Stderr: &errb,
			WebQARunner: func(_ context.Context, _ []string, _, _ io.Writer) int {
				return wantCode
			},
		}
		got := Run(context.Background(), []string{"web-qa", "preflight"}, d)
		if got != wantCode {
			t.Errorf("exit code: obteve %d, queria %d", got, wantCode)
		}
	}
}

func TestWebQA_CorrigirBloqueiaPathTraversal(t *testing.T) {
	var out, errb bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errb}
	outside := filepath.Join(t.TempDir(), "RELATORIO_QA_WEB.json")
	if err := os.WriteFile(outside, []byte(`{"alvo":"x","findings":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, bad := range []string{outside, filepath.Join("..", "..", "..", "fora.json")} {
		out.Reset()
		errb.Reset()
		code := Run(context.Background(), []string{"web-qa", "corrigir", "--relatorio", bad}, d)
		if code != 1 {
			t.Errorf("path traversal %q deveria falhar com 1, obteve %d", bad, code)
		}
		if !strings.Contains(errb.String(), "path jail") {
			t.Errorf("mensagem deveria citar path jail, obteve: %s", errb.String())
		}
	}
}

func TestWebQA_CorrigirMontaPromptDeRelatorioReal(t *testing.T) {
	rig := newInterRig(t, "", lowReply)
	spy := &sysLLM{inner: rig.llm}
	rig.d.Fabricar.LLM = spy

	if err := os.MkdirAll(rig.doc, 0o755); err != nil {
		t.Fatal(err)
	}
	relPath := filepath.Join(rig.doc, "RELATORIO_QA_WEB.json")
	payload := map[string]any{
		"alvo": "studio",
		"findings": []map[string]string{
			{"tool": "playwright", "severity": "critical", "route": "/painel", "message": "HTTP 500"},
			{"tool": "axe-core", "severity": "serious", "route": "/painel", "message": "contraste baixo"},
			{"tool": "lychee", "severity": "critical", "route": "/docs", "message": "link quebrado"},
		},
	}
	raw, _ := json.Marshal(payload)
	if err := os.WriteFile(relPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	code := run(t, rig.d, "web-qa", "corrigir", "--relatorio", relPath, "--vault", rig.vault, "--docs", rig.doc, "--model", "sonnet")
	if code != 0 {
		t.Fatalf("corrigir falhou: exit=%d stderr=%s", code, rig.errb.String())
	}
	var p passo
	if err := json.Unmarshal(rig.out.Bytes(), &p); err != nil {
		t.Fatalf("saída inválida: %v\n%s", err, rig.out.String())
	}
	if p.Sessao == "" {
		t.Fatalf("sessão não criada: %+v", p)
	}
	for _, want := range []string{"studio", "Total de achados: 3", "playwright=1", "axe-core=1", "lychee=1", "/painel=2", "critical=2", "serious=1", "HTTP 500"} {
		if !strings.Contains(spy.system, want) {
			t.Errorf("prompt enviado ao LLM não contém %q:\n%s", want, spy.system)
		}
	}
}

func TestWebQA_ErroExplicitoSemScript(t *testing.T) {
	t.Chdir(t.TempDir())
	var out, errb bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errb}

	code := Run(context.Background(), []string{"web-qa", "run", "--url", "http://127.0.0.1:7878"}, d)
	if code != 1 {
		t.Fatalf("esperava exit 1 sem webqa.mjs, obteve %d", code)
	}
	if !strings.Contains(errb.String(), "webqa.mjs") {
		t.Fatalf("erro deveria mencionar webqa.mjs, obteve: %s", errb.String())
	}
}
