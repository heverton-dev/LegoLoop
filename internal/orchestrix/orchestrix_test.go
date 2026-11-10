package orchestrix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
)

// TestMain transforma o binário de teste num harness real quando LEGOLOOP_HELPER está definido.
func TestMain(m *testing.M) {
	switch os.Getenv("LEGOLOOP_HELPER") {
	case "claude":
		helperClaude()
	case "claude-unit":
		helperClaudeUnit()
	case "claude-tools":
		helperClaudeTools()
	case "script":
		helperScript()
	case "substep":
		helperSubstep()
	case "agy", "agy-error", "omp", "omp-error", "mimo", "mimo-error":
		helperHarness(os.Getenv("LEGOLOOP_HELPER"))
	case "claude-api-error":
		// formato real observado no claude 2.1.314 quando a API recusa a chamada
		os.Stdout.Write([]byte(`{"type":"result","is_error":true,"terminal_reason":"api_error","result":"Credit balance is too low","usage":{"input_tokens":0,"output_tokens":0}}`))
		os.Exit(1)
	case "flood":
		os.Stdout.Write([]byte(strings.Repeat("x", 4096)))
		os.Exit(0)
	case "claude-retry":
		helperClaudeRetry()
	case "fail-out":
		fmt.Print("SAIDA-DO-COMMAND\n")
		os.Exit(1)
	case "echo-args":
		fmt.Print(strings.Join(os.Args[1:], " "))
		os.Exit(0)
	default:
		os.Exit(m.Run())
	}
}

func helperClaude() {
	prompt, _ := io.ReadAll(os.Stdin)
	args := strings.Join(os.Args[1:], " ")
	for _, want := range []string{"-p", "--model sonnet", "--output-format json", "--no-session-persistence",
		"--strict-mcp-config", "--safe-mode", "--setting-sources user", "--tools "} {
		if !strings.Contains(args, want) {
			fmt.Fprintln(os.Stderr, "flag ausente:", want)
			os.Exit(2)
		}
	}
	leak := os.Getenv("LEGOLOOP_SECRET_LEAK")
	cwd, _ := os.Getwd()
	entries, _ := os.ReadDir(cwd)
	leak += fmt.Sprintf("\ncwd=%s\ncwd_entries=%d\nnodefault=%s", cwd, len(entries), os.Getenv("NoDefaultCurrentDirectoryInExePath"))
	out, _ := json.Marshal(map[string]any{
		"type":   "result",
		"result": "# BRIEFING\nleak=" + leak + "\n" + string(prompt),
		"usage":  map[string]int{"input_tokens": 120, "output_tokens": 45},
	})
	os.Stdout.Write(out)
	os.Exit(0)
}

func helperScript() {
	os.WriteFile("made.txt", []byte("feito"), 0o644)
	fmt.Print("saida-do-script")
	os.Exit(3)
}

func helperSubstep() {
	target := "step.txt"
	if len(os.Args) > 1 && os.Args[len(os.Args)-1] != "" {
		target = os.Args[len(os.Args)-1]
	}
	os.WriteFile(target, []byte("ok "+target), 0o644)
	os.Exit(0)
}

func skillsDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "skills", "interview")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p, "SKILL.md"), []byte("Entreviste o usuário."), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func skillNode(r dag.RoutingConfig, inputs map[string]string) dag.Node {
	return dag.Node{ID: "n1", Type: dag.TypeSkill, Ref: "skills/interview", Routing: r, Inputs: inputs,
		Output: dag.HandoffContract{Artifact: "out/BRIEFING.md", Gate: "true"}}
}

func readArtifact(t *testing.T, workdir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(workdir, "out", "BRIEFING.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestClaudeCLIIsolatedProcess(t *testing.T) {
	t.Setenv("LEGOLOOP_SECRET_LEAK", "vazou")
	work := t.TempDir()
	os.WriteFile(filepath.Join(work, "IDEIA.md"), []byte("app de tarefas"), 0o644)
	o := &Orchestrix{
		Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		SkillsDir: skillsDir(t),
		Vars:      map[string]string{"idea": "clone do trello"},
		ExtraEnv:  []string{"LEGOLOOP_HELPER=claude"},
	}
	n := skillNode(dag.RoutingConfig{Mode: dag.ModeSuggested, Model: "gemini-x", Fallback: "sonnet"},
		map[string]string{"prompt": "var.idea", "doc": "file('IDEIA.md')", "nota": "literal"})
	res, err := o.Execute(context.Background(), n, work)
	if err != nil {
		t.Fatal(err)
	}
	got := readArtifact(t, work)
	for _, want := range []string{"# BRIEFING", "Entreviste o usuário.", "clone do trello", "app de tarefas", "literal", "out/BRIEFING.md"} {
		if !strings.Contains(got, want) {
			t.Errorf("artefato sem %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "leak=vazou") {
		t.Fatal("ambiente do processo pai vazou para a Orchestrix")
	}
	// SEC-02: o harness roda num diretório vazio, nunca no repositório (sem .claude/ nem CLAUDE.md).
	if strings.Contains(got, "cwd="+work) || !strings.Contains(got, "cwd_entries=0") || !strings.Contains(got, "nodefault=1") {
		t.Fatalf("subprocesso não isolado do repositório:\n%s", got)
	}
	if res.Tokens.Input != 120 || res.Tokens.Output != 45 || res.ExitCode != 0 {
		t.Fatalf("result = %+v", res)
	}
}

func TestRouterAutoCombo(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer k9" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Write([]byte(`{"choices":[{"message":{"content":"# SPEC\n- [ ] ok"}}],"usage":{"prompt_tokens":2300,"completion_tokens":50}}`))
	}))
	defer srv.Close()
	o := &Orchestrix{
		Inventory: provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"code-pro"}},
		HTTP:      srv.Client(), RouterKey: "k9", RouterKeyURL: srv.URL, SkillsDir: skillsDir(t),
	}
	work := t.TempDir()
	res, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "code-pro"}, nil), work)
	if err != nil {
		t.Fatal(err)
	}
	if gotBody["model"] != "code-pro" || gotBody["stream"] != false {
		t.Fatalf("payload = %v", gotBody)
	}
	if readArtifact(t, work) != "# SPEC\n- [ ] ok" {
		t.Fatal("artefato incorreto")
	}
	// 9Router soma +2000 fixos em prompt_tokens (docs/references/achados_9router.md).
	// Auditoria: o valor bruto reportado pelo 9Router fica registrado ao lado do real.
	if res.Tokens.ReportedInput != 2300 {
		t.Fatalf("valor bruto do 9Router não registrado: %+v", res.Tokens)
	}
	if res.Tokens.Input != 300 || res.Tokens.Output != 50 {
		t.Fatalf("tokens = %+v", res.Tokens)
	}
}

func TestOllamaManual(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if r.URL.Path != "/api/generate" || body["model"] != "qwen2.5-coder:7b" || body["stream"] != false {
			w.WriteHeader(400)
			return
		}
		w.Write([]byte(`{"response":"codigo","prompt_eval_count":11,"eval_count":7}`))
	}))
	defer srv.Close()
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"qwen2.5-coder:7b"}}, HTTP: srv.Client(), OllamaURL: srv.URL, SkillsDir: skillsDir(t)}
	work := t.TempDir()
	res, err := o.Execute(context.Background(), skillNode(dag.RoutingConfig{Mode: dag.ModeManual, Model: "qwen2.5-coder:7b"}, nil), work)
	if err != nil {
		t.Fatal(err)
	}
	if readArtifact(t, work) != "codigo" || res.Tokens.Input != 11 || res.Tokens.Output != 7 {
		t.Fatalf("res=%+v", res)
	}
}

func TestRegraDeOuroBlocksUnavailableModels(t *testing.T) {
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"llama3"}}, SkillsDir: skillsDir(t)}
	cases := []dag.RoutingConfig{
		{Mode: dag.ModeManual, Model: "opus", Fallback: "llama3"}, // manual ignora fallback
		{Mode: dag.ModeSuggested, Model: "opus", Fallback: "sonnet"},
		{Mode: dag.ModeAuto, Combo: "code-pro"},
	}
	for _, r := range cases {
		if _, err := o.Execute(context.Background(), skillNode(r, nil), t.TempDir()); err == nil || !strings.Contains(err.Error(), "inventário") {
			t.Errorf("%+v: esperado bloqueio por inventário, err=%v", r, err)
		}
	}
}

func TestPromptInputErrors(t *testing.T) {
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"m"}}, SkillsDir: skillsDir(t)}
	r := dag.RoutingConfig{Mode: dag.ModeManual, Model: "m"}
	bad := []dag.Node{
		skillNode(r, map[string]string{"x": "var.inexistente"}),
		skillNode(r, map[string]string{"x": "file('../fora.md')"}),
		skillNode(r, map[string]string{"x": "file('nao-existe.md')"}),
	}
	ref := skillNode(r, nil)
	ref.Ref = "../../etc"
	missing := skillNode(r, nil)
	missing.Ref = "skills/nao-existe"
	bad = append(bad, ref, missing)
	for _, n := range bad {
		if _, err := o.Execute(context.Background(), n, t.TempDir()); err == nil {
			t.Errorf("nó %+v deveria falhar", n.Inputs)
		}
	}
}

func TestScriptNodeNoShell(t *testing.T) {
	work := t.TempDir()
	o := &Orchestrix{ExtraEnv: []string{"LEGOLOOP_HELPER=script"}, AllowScripts: true}
	n := dag.Node{ID: "s", Type: dag.TypeScript, Command: []string{os.Args[0], "; rm -rf /", "$(whoami)"},
		Output: dag.HandoffContract{Artifact: "log.txt", Gate: "true"}}
	res, err := o.Execute(context.Background(), n, work)
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 3 {
		t.Fatalf("exit code = %d", res.ExitCode)
	}
	if b, _ := os.ReadFile(filepath.Join(work, "made.txt")); string(b) != "feito" {
		t.Fatal("script não rodou no workdir")
	}
	if b, _ := os.ReadFile(filepath.Join(work, "log.txt")); string(b) != "saida-do-script" {
		t.Fatalf("stdout não virou artefato: %q", b)
	}
	n.Command = []string{"comando-que-nao-existe-legoloop"}
	if _, err := o.Execute(context.Background(), n, work); err == nil {
		t.Fatal("comando inexistente deveria falhar")
	}
}

func TestBackendHTTPErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "chat") {
			w.Write([]byte(`{"choices":[]}`))
			return
		}
		w.WriteHeader(500)
	}))
	defer srv.Close()
	o := &Orchestrix{Inventory: provider.Inventory{Ollama: []string{"m"}, RouterURL: srv.URL, RouterModels: []string{"c"}},
		HTTP: srv.Client(), OllamaURL: srv.URL, SkillsDir: skillsDir(t)}
	for _, r := range []dag.RoutingConfig{{Mode: dag.ModeManual, Model: "m"}, {Mode: dag.ModeAuto, Combo: "c"}} {
		if _, err := o.Execute(context.Background(), skillNode(r, nil), t.TempDir()); err == nil {
			t.Errorf("%+v: esperado erro de backend", r)
		}
	}
}

func TestClaudeCLIFailureAndPlainOutput(t *testing.T) {
	o := &Orchestrix{Inventory: provider.Inventory{CLIs: []provider.CLI{{Name: "claude", Path: os.Args[0]}}},
		SkillsDir: skillsDir(t), ExtraEnv: []string{"LEGOLOOP_HELPER=script"}}
	r := dag.RoutingConfig{Mode: dag.ModeManual, Model: "haiku"}
	if _, err := o.Execute(context.Background(), skillNode(r, nil), t.TempDir()); err == nil {
		t.Fatal("exit code != 0 do claude deveria falhar")
	}
	if r := parseClaudeOutput([]byte("texto puro")); r.text != "texto puro" || r.usage.Input != 0 || r.isError {
		t.Fatal("saída não-JSON deve ser tratada como texto")
	}
}

func TestUnknownNodeType(t *testing.T) {
	o := &Orchestrix{}
	if _, err := o.Execute(context.Background(), dag.Node{Type: "mcp"}, t.TempDir()); err == nil {
		t.Fatal("tipo desconhecido deveria falhar")
	}
}

func TestScriptVariableInterpolation(t *testing.T) {
	work := t.TempDir()
	o := &Orchestrix{
		ExtraEnv:     []string{"LEGOLOOP_HELPER=echo-args"},
		AllowScripts: true,
		Vars:         map[string]string{"VIDEO_URL": "https://youtube.com/watch?v=123", "NAME": "shorts"},
	}
	n := dag.Node{
		ID:      "ingest",
		Type:    dag.TypeScript,
		Command: []string{os.Args[0], "--url", "${VIDEO_URL}", "--name", "${NAME}"},
		Output:  dag.HandoffContract{Artifact: "out.txt", Gate: "true"},
	}
	res, err := o.Execute(context.Background(), n, work)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit code esperado 0, obteve %d", res.ExitCode)
	}
	content, err := os.ReadFile(filepath.Join(work, "out.txt"))
	if err != nil {
		t.Fatalf("falha ao ler artefato: %v", err)
	}
	want := "--url https://youtube.com/watch?v=123 --name shorts"
	if string(content) != want {
		t.Fatalf("interpolação falhou!\nquer: %q\nobteve: %q", want, string(content))
	}
}
