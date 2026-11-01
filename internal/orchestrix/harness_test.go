package orchestrix

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

// helperHarness imita a saída REAL de cada CLI (formatos capturados em 31/10/2026).
// Ecoa no texto de resposta o prompt recebido, para o teste verificar a entrega.
func helperHarness(kind string) {
	args := os.Args[1:]
	switch kind {
	case "agy":
		prompt := ""
		for _, a := range args {
			if p, ok := strings.CutPrefix(a, "-p="); ok {
				prompt = p
			}
		}
		fmt.Printf(`{"status":"SUCCESS","response":%q,"usage":{"input_tokens":18693,"output_tokens":12}}`, "AGY "+strings.Join(args, " ")+"\n"+prompt)
	case "agy-error":
		fmt.Print(`{"status":"ERROR","response":"model not found","usage":{"input_tokens":0,"output_tokens":0}}`)
		os.Exit(1)
	case "omp":
		body, _ := os.ReadFile("prompt.md") // anexo @prompt.md relativo ao cwd
		text := fmt.Sprintf("OMP %s\n%s", strings.Join(args, " "), body)
		fmt.Printf(`{"type":"agent_start"}`+"\n"+`{"type":"agent_end","messages":[{"role":"user","content":[]},{"role":"assistant","content":[{"type":"text","text":%q}],"usage":{"input":40,"output":7},"stopReason":"stop"}]}`+"\n", text)
	case "omp-error":
		fmt.Print(`{"type":"agent_end","messages":[{"role":"assistant","content":[],"usage":{"input":0,"output":0},"stopReason":"error","errorMessage":"401 API key required for remote API access"}]}` + "\n")
		os.Exit(0) // o omp real sai com 0 mesmo em erro
	case "mimo":
		body, _ := os.ReadFile("prompt.md")
		fmt.Printf(`{"type":"step_start","part":{"type":"step-start"}}`+"\n"+`{"type":"text","part":{"type":"text","text":%q}}`+"\n"+`{"type":"text","part":{"type":"text","text":"!"}}`+"\n"+`{"type":"step_finish","part":{"type":"step-finish","tokens":{"input":98581,"output":3}}}`+"\n", "MIMO "+strings.Join(args, " ")+"\n"+string(body))
	case "mimo-error":
		fmt.Print(`{"type":"error","error":{"name":"ProviderModelNotFound","data":{"message":"modelo x/y não existe"}}}` + "\n")
		os.Exit(1)
	}
	os.Exit(0)
}

func harnessOrchestrix(t *testing.T, cli, helper string) *Orchestrix {
	t.Helper()
	return &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: cli, Path: os.Args[0]}}},
		SkillsDir: skillsDir(t), ExtraEnv: []string{"LEGOLOOP_HELPER=" + helper},
	}
}

func TestHarnessExecutors(t *testing.T) {
	cases := []struct {
		cli, model, helper string
		wantArgs           []string
		in, out            int
	}{
		{"agy", "agy:gemini-3.8-flash-low", "agy", []string{"--output-format json", "--model gemini-3.8-flash-low", "--disable-slash-commands"}, 18693, 12},
		{"agy", "agy:default", "agy", []string{"--output-format json"}, 18693, 12},
		{"omp", "omp:code-fast", "omp", []string{"-p", "--no-session", "--no-tools", "--mode json", "--model code-fast", "@prompt.md"}, 40, 7},
		{"mimo", "mimo:xiaomi/mimo-v2", "mimo", []string{"run", "--format json", "--pure", "-m xiaomi/mimo-v2", "-f prompt.md"}, 98581, 3},
		{"opencode", "opencode:default", "mimo", []string{"run", "--format json", "--pure", "-f prompt.md"}, 98581, 3},
	}
	for _, c := range cases {
		o := harnessOrchestrix(t, c.cli, c.helper)
		work := t.TempDir()
		res, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: c.model}, nil), work)
		if err != nil {
			t.Fatalf("%s: %v", c.model, err)
		}
		got := readArtifact(t, work)
		for _, w := range append(c.wantArgs, "Entreviste o usuário.") { // a skill chegou ao harness
			if !strings.Contains(got, w) {
				t.Errorf("%s: saída sem %q:\n%s", c.model, w, got)
			}
		}
		if strings.HasSuffix(c.model, ":default") && strings.Contains(got, "--model") {
			t.Errorf("%s: modelo padrão não deve passar --model/-m", c.model)
		}
		if res.Tokens.Input != c.in || res.Tokens.Output != c.out || res.Model != c.cli+"/"+c.model {
			t.Errorf("%s: telemetria %+v model=%q", c.model, res.Tokens, res.Model)
		}
	}
}

func TestHarnessErrorsAreSurfaced(t *testing.T) {
	cases := map[string][2]string{
		"agy-error":  {"agy", "model not found"},
		"omp-error":  {"omp", "401 API key required"},
		"mimo-error": {"mimo", "modelo x/y não existe"},
	}
	for helper, c := range cases {
		o := harnessOrchestrix(t, c[0], helper)
		_, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: c[0] + ":m"}, nil), t.TempDir())
		if err == nil || !strings.Contains(err.Error(), c[1]) {
			t.Errorf("%s: erro não reportado: %v", helper, err)
		}
	}
}

// agy só aceita o prompt como argumento: acima do limite seguro, recusa com clareza.
func TestAgyPromptTooLong(t *testing.T) {
	_, _, err := agyArgs("m", strings.Repeat("x", maxArgPromptBytes+1), "")
	if err == nil || !strings.Contains(err.Error(), "excede") {
		t.Fatalf("prompt gigante aceito: %v", err)
	}
}

func TestHarnessParsersRejectEmptyOutput(t *testing.T) {
	for name, parse := range map[string]func([]byte) (string, int, int, error){
		"agy": parseAgy, "omp": parseOmp, "mimo": parseOpencodeEvents,
	} {
		if _, _, _, err := parse([]byte("")); err == nil {
			t.Errorf("%s: saída vazia aceita", name)
		}
	}
}

func TestHarnessEffortClamping(t *testing.T) {
	o := harnessOrchestrix(t, "agy", "agy")
	work := t.TempDir()
	cfg := dag.RoutingConfig{Mode: dag.ModeManual, Model: "agy:gemini-3.8-flash-low", Effort: dag.EffortLow}
	res, err := o.Execute(context.Background(), skillNode(cfg, nil), work)
	if err != nil {
		t.Fatalf("Execute com effort falhou: %v", err)
	}
	got := readArtifact(t, work)
	if !strings.Contains(got, "--effort low") && !strings.Contains(got, "--effort=low") {
		t.Errorf("saída do agy deveria conter flag de effort: %s", got)
	}
	hasEffortNote := false
	for _, n := range res.Notes {
		if strings.Contains(n, "effort low") || strings.Contains(n, "effort 'low'") {
			hasEffortNote = true
			break
		}
	}
	if !hasEffortNote {
		t.Errorf("res.Notes deveria registrar o esforço aplicado: %v", res.Notes)
	}
}

