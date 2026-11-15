package runner_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/brain"
	"legoloop/internal/combo"
	"legoloop/internal/dag"
	"legoloop/internal/interview"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

func TestN10EnginesConnectedInProduction(t *testing.T) {
	t.Run("ProvenanceAndRedBlueBoundary", testProvenanceAndRedBlue)
	t.Run("AndonCordRetryAndSBAR", testAndonCordAndSBAR)
	t.Run("STE100AndFlightRecordTelemetry", testSTE100AndFlightRecord)
	t.Run("IsolatedASTMerge", testIsolatedASTMerge)
	t.Run("OrchestrixComboAntibodiesSpeculativeAndHTTP413429", testOrchestrixN10AndHTTP)
	t.Run("SubprocessoPathJailAndFlagsPropagation", testSubprocessoPathJailAndFlags)
	t.Run("HarnessGatewayEnvAndStderr", testHarnessGatewayAndStderr)
}

func testProvenanceAndRedBlue(t *testing.T) {
	dir := t.TempDir()
	provDAG := mustParseDAG(t, `{"name":"prov","nodes":[{"id":"n1","type":"script","command":["x"],
		"context_in":["inexistente.txt"],"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"out.txt","gate":"exit_code == 0"}}]}`)
	if _, err := runner.Run(context.Background(), provDAG, runner.Options{Workdir: dir, Executor: &staticExec{}}); err == nil || !strings.Contains(err.Error(), "arquivo fantasma") {
		t.Fatalf("esperava erro de ValidateProvenance, obteve: %v", err)
	}

	rbDAG := mustParseDAG(t, `{"name":"rb","nodes":[{"id":"n_blue","type":"script","command":["x"],
		"context_out":["calc_test.go"],"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"out.txt","gate":"exit_code == 0"}}]}`)
	rbDAG.Nodes[0].Role = "blue_team"
	rep, err := runner.Run(context.Background(), rbDAG, runner.Options{Workdir: dir, Executor: &staticExec{}})
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}
	if rep.OK || len(rep.Nodes) != 1 || !strings.Contains(rep.Nodes[0].Error, "reward_hacking") {
		t.Fatalf("esperava bloqueio ValidateRedBlueBoundary (reward_hacking), obteve: %+v", rep.Nodes)
	}
}

type flakyExec struct {
	calls     int
	failUntil int
}

func (f *flakyExec) Execute(_ context.Context, n dag.Node, workdir string) (runner.Result, error) {
	f.calls++
	if f.calls <= f.failUntil {
		return runner.Result{ExitCode: 1}, errors.New("falha transiente no executor")
	}
	_ = os.WriteFile(filepath.Join(workdir, n.Output.Artifact), []byte("ok"), 0o644)
	return runner.Result{ExitCode: 0, Model: "test/andon"}, nil
}

func testAndonCordAndSBAR(t *testing.T) {
	d := mustParseDAG(t, `{"name":"andon","nodes":[{"id":"n_andon","type":"script","command":["x"],
		"routing":{"mode":"manual","model":"none"},"output":{"artifact":"out.txt","gate":"exit_code == 0"}}]}`)
	recoverExec := &flakyExec{failUntil: 1}
	repOK, err := runner.Run(context.Background(), d, runner.Options{Workdir: t.TempDir(), Executor: recoverExec})
	if err != nil || !repOK.OK || recoverExec.calls != 2 {
		t.Fatalf("esperava recuperação na 2a tentativa via ExecuteWithAndon (calls=%d, ok=%v, err=%v)", recoverExec.calls, repOK.OK, err)
	}

	alwaysFail := &flakyExec{failUntil: 5}
	repFail, err := runner.Run(context.Background(), d, runner.Options{Workdir: t.TempDir(), Executor: alwaysFail})
	if err != nil || repFail.OK || alwaysFail.calls != 2 {
		t.Fatalf("esperava parada Andon Cord em 2 tentativas (calls=%d, ok=%v)", alwaysFail.calls, repFail.OK)
	}
	joinedNotes := strings.Join(repFail.Nodes[0].Notes, "\n") + "\n" + repFail.Nodes[0].Error
	if !strings.Contains(joinedNotes, "andon_cord") || !strings.Contains(joinedNotes, "[S] SITUATION:") {
		t.Fatalf("esperava diagnóstico SBAR e andon_cord no relatório, obteve: %+v", repFail.Nodes[0])
	}
}

type staticExec struct {
	content string
}

func (s *staticExec) Execute(_ context.Context, n dag.Node, workdir string) (runner.Result, error) {
	body := s.content
	if body == "" {
		body = "Texto curto e direto sem termos dubios."
	}
	_ = os.WriteFile(filepath.Join(workdir, filepath.FromSlash(n.Output.Artifact)), []byte(body), 0o644)
	return runner.Result{ExitCode: 0, Tokens: runner.TokenUsage{Input: 25, Output: 10}, Model: "test/ste"}, nil
}

func testSTE100AndFlightRecord(t *testing.T) {
	d := mustParseDAG(t, `{"name":"ste-flight","nodes":[{"id":"n_ste","type":"script","command":["x"],
		"rule_tags":["ste100"],"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"doc.md","gate":"exit_code == 0"}}]}`)
	badRep, err := runner.Run(context.Background(), d, runner.Options{
		Workdir:  t.TempDir(),
		Executor: &staticExec{content: "Talvez este documento funcione se tudo correr bem."},
	})
	if err != nil || badRep.OK || badRep.Nodes[0].Status != runner.StatusGateFailed || !strings.Contains(badRep.Nodes[0].Error, "hedging") {
		t.Fatalf("esperava falha ValidateSTE100 por hedging, obteve: %+v", badRep.Nodes)
	}

	goodDir := t.TempDir()
	flightPath := filepath.Join(goodDir, "custom_flight.json")
	goodRep, err := runner.Run(context.Background(), d, runner.Options{
		Workdir:          goodDir,
		FlightRecordPath: flightPath,
		Executor:         &staticExec{content: "O sistema valida a regra com precisao."},
	})
	if err != nil || !goodRep.OK {
		t.Fatalf("esperava aprovação STE100, obteve: %+v (err=%v)", goodRep.Nodes, err)
	}
	data, readErr := os.ReadFile(flightPath)
	if readErr != nil || !strings.Contains(string(data), `"first_pass_yield": 1`) {
		t.Fatalf("FlightRecord não salvo corretamente em %s: err=%v data=%s", flightPath, readErr, string(data))
	}
}

func testIsolatedASTMerge(t *testing.T) {
	workdir := t.TempDir()
	baseCode := "package sample\n\nfunc First() int {\n\treturn 1\n}\n"
	incomingCode := "package sample\n\nfunc Second() int {\n\treturn 2\n}\n"
	if err := os.WriteFile(filepath.Join(workdir, "sample.go"), []byte(baseCode), 0o644); err != nil {
		t.Fatal(err)
	}
	d := mustParseDAG(t, `{"name":"ast-merge","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{
		"id":"n_ast","type":"script","command":["x"],"execution":{"mode":"isolated"},
		"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"sample.go","gate":"file_exists('sample.go') && exit_code == 0"}}]}`)
	rep, err := runner.Run(context.Background(), d, runner.Options{
		Workdir:  workdir,
		Executor: &staticExec{content: incomingCode},
	})
	if err != nil || !rep.OK {
		t.Fatalf("Run isolated AST merge falhou: err=%v rep=%+v", err, rep)
	}
	merged, err := os.ReadFile(filepath.Join(workdir, "sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(merged), "func First()") || !strings.Contains(string(merged), "func Second()") {
		t.Fatalf("MergeAST não fundiu First e Second:\n%s", string(merged))
	}
}

func testOrchestrixN10AndHTTP(t *testing.T) {
	skills := setupTestSkill(t)
	abDB := brain.NewAntibodiesDB()
	_ = abDB.Record("skills/interview", "erro-antigo", "Causa Raiz X", "Solucao Y")

	var promptsReceived []string
	srv := newSpeculativeTestServer(t, &promptsReceived)
	defer srv.Close()

	specCalled := false
	ox := &orchestrix.Orchestrix{
		Inventory:    provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"t-chat-128k", "t-raciocinio-128k", "t-codigo-128k"}},
		HTTP:         srv.Client(),
		RouterKey:    "k9",
		RouterKeyURL: srv.URL,
		SkillsDir:    skills,
		Antibodies:   abDB,
		Speculative:  makeTestSpeculativeRunner(&specCalled),
	}
	verifyComboValidationAnd413(t, ox)
	verifySpeculativeAndAntibodies(t, ox, abDB, &promptsReceived, &specCalled)
}

type testBrainAdapter struct {
	cand  orchestrix.SpeculativeCandidate
	usage runner.TokenUsage
}

func (b *testBrainAdapter) Chat(ctx context.Context, msgs []interview.Message) (string, error) {
	return b.Infer(ctx, msgs[len(msgs)-1].Content)
}

func (b *testBrainAdapter) Infer(ctx context.Context, prompt string) (string, error) {
	txt, tok, err := b.cand.Run(ctx, prompt)
	b.usage = tok
	return txt, err
}

func (b *testBrainAdapter) Name() string           { return b.cand.Name }
func (b *testBrainAdapter) Model() string          { return b.cand.Name }
func (b *testBrainAdapter) Timeout() time.Duration { return b.cand.Timeout }
func (b *testBrainAdapter) Effort() brain.Effort   { return brain.EffortMedium }

func makeTestSpeculativeRunner(called *bool) orchestrix.SpeculativeRunner {
	return func(ctx context.Context, fast, heavy orchestrix.SpeculativeCandidate, prompt string) (string, string, runner.TokenUsage, error) {
		*called = true
		fb, hb := &testBrainAdapter{cand: fast}, &testBrainAdapter{cand: heavy}
		txt, used, err := brain.ExecuteSpeculative(ctx, fb, hb, prompt, fast.Timeout)
		if err != nil {
			return "", "", runner.TokenUsage{}, err
		}
		if used == hb.Name() {
			return txt, used, hb.usage, nil
		}
		return txt, used, fb.usage, nil
	}
}

func newSpeculativeTestServer(t *testing.T, prompts *[]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("erro ao ler body: %v", err)
		}
		raw := string(bodyBytes)
		*prompts = append(*prompts, raw)
		if strings.Contains(raw, `"model":"t-codigo-128k"`) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_, _ = w.Write([]byte(`{"error":"context window exceeded"}`))
			return
		}
		if strings.Contains(raw, `"model":"t-chat-128k"`) {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"Artefato final validado."}}],"usage":{"prompt_tokens":2100,"completion_tokens":20}}`))
	}))
}

func verifyComboValidationAnd413(t *testing.T, ox *orchestrix.Orchestrix) {
	t.Helper()
	nTool := dag.Node{
		ID:      "n_combo_tool",
		Type:    dag.TypeSkill,
		Ref:     "skills/interview",
		Tooling: map[string]any{"require_tools": true},
		Routing: dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "t-codigo-128k"},
		Output:  dag.HandoffContract{Artifact: "out.md", Gate: "true"},
	}
	if _, err := ox.Execute(context.Background(), nTool, t.TempDir()); !errors.Is(err, combo.ErrFerramentaNaoSuportada) {
		t.Fatalf("esperava ErrFerramentaNaoSuportada de ValidarRequisicao, obteve: %v", err)
	}

	n413 := dag.Node{
		ID:      "n_413",
		Type:    dag.TypeSkill,
		Ref:     "skills/interview",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "t-codigo-128k"},
		Output:  dag.HandoffContract{Artifact: "out.md", Gate: "true"},
	}
	if _, err := ox.Execute(context.Background(), n413, t.TempDir()); err == nil || !strings.Contains(err.Error(), "413") {
		t.Fatalf("esperava erro explícito HTTP 413, obteve: %v", err)
	}
}

func verifySpeculativeAndAntibodies(t *testing.T, ox *orchestrix.Orchestrix, abDB *brain.AntibodiesDB, prompts *[]string, specCalled *bool) {
	t.Helper()
	*prompts = nil
	nSpec := dag.Node{
		ID:      "n_spec",
		Type:    dag.TypeSkill,
		Ref:     "skills/interview",
		Routing: dag.RoutingConfig{Mode: dag.ModeAuto, Combo: "t-chat-128k", Fallback: "t-raciocinio-128k"},
		Output:  dag.HandoffContract{Artifact: "out.md", Gate: "true"},
	}
	res, err := ox.Execute(context.Background(), nSpec, t.TempDir())
	if err != nil {
		t.Fatalf("esperava sucesso via hot-swap especulativo em ModeAuto após 429, obteve: %v", err)
	}
	if !*specCalled || !strings.Contains(res.Model, "t-raciocinio-128k") || len(*prompts) < 2 {
		t.Fatalf("esperava ExecuteSpeculative -> t-raciocinio-128k (spec=%v model=%s calls=%d)", *specCalled, res.Model, len(*prompts))
	}
	secondPrompt := (*prompts)[len(*prompts)-1]
	if !strings.Contains(secondPrompt, "ANTICORPO PREVENTIVO") || !strings.Contains(secondPrompt, "CONTEXTO NEGATIVO") {
		t.Fatalf("prompt do segundo modelo deveria conter ANTICORPO PREVENTIVO e CONTEXTO NEGATIVO:\n%s", secondPrompt)
	}
	if !strings.Contains(abDB.InjectPreventiveContext("skills/interview"), "429") {
		t.Fatalf("AntibodiesDB deveria ter registrado a falha 429 do modelo rápido")
	}
}

func testSubprocessoPathJailAndFlags(t *testing.T) {
	store, err := studiodoc.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	skills := setupTestSkill(t)
	srv := newSpeculativeTestServer(t, new([]string))
	defer srv.Close()
	saveSubprocessoTestSpecs(t, store)

	ox := &orchestrix.Orchestrix{
		Inventory: provider.Inventory{RouterURL: srv.URL, RouterModels: []string{"t-raciocinio-128k"}},
		HTTP:      srv.Client(), RouterKey: "k9", RouterKeyURL: srv.URL, SkillsDir: skills,
		AllowScripts: true, Vars: map[string]string{"TARGET_FILE": "OK.txt"},
	}
	ox.SetAcervo(store)
	escapeNode := dag.Node{
		ID: "sub_escape", Type: dag.TypeSubprocesso, Ref: "sub-vars",
		Inputs:  map[string]string{"bad": "../segredo.txt"},
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "sub"},
		Output:  dag.HandoffContract{Artifact: "OK.txt", Gate: "true"},
	}
	if _, err := ox.Execute(context.Background(), escapeNode, t.TempDir()); err == nil || !strings.Contains(strings.ToLower(err.Error()), "path jail") {
		t.Fatalf("esperava bloqueio de Path Jail em transferInputs do subprocesso, obteve: %v", err)
	}

	okNode := dag.Node{
		ID: "sub_ok", Type: dag.TypeSubprocesso, Ref: "sub-vars",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "sub"},
		Output:  dag.HandoffContract{Artifact: "OK.txt", Gate: "true"},
	}
	if _, err := ox.Execute(context.Background(), okNode, t.TempDir()); err != nil {
		t.Fatalf("esperava propagação de Vars no subprocesso, obteve: %v", err)
	}
	ox.AllowScripts = false
	scriptSub := dag.Node{
		ID: "sub_no_scripts", Type: dag.TypeSubprocesso, Ref: "sub-script-slot",
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "sub"},
		Output:  dag.HandoffContract{Artifact: "OK.txt", Gate: "true"},
	}
	if _, err := ox.Execute(context.Background(), scriptSub, t.TempDir()); err == nil || !strings.Contains(err.Error(), "--allow-scripts") {
		t.Fatalf("esperava propagação de AllowScripts=false no subprocesso, obteve: %v", err)
	}
}

func saveSubprocessoTestSpecs(t *testing.T, store *studiodoc.Store) {
	t.Helper()
	subVars := &studiodoc.Processo{
		ID: "sub-vars", Name: "Sub Vars", Version: "2.0.0", Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{{
			ID: "n1", Type: dag.TypeSkill, Ref: "skills/interview",
			Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "t-raciocinio-128k"},
			Output:  dag.HandoffContract{Artifact: "OK.txt", Gate: "file_exists('${TARGET_FILE}')"},
		}},
	}
	subScript := &studiodoc.Processo{
		ID: "sub-script-slot", Name: "Sub Script Slot", Version: "2.0.0", Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{{
			ID: "n1", Type: dag.TypeSkill, Ref: "skills/interview",
			Routing: dag.RoutingConfig{Mode: dag.ModeManual, Model: "t-raciocinio-128k"},
			Slots:   &dag.Slots{Scripts: [][]string{{"echo", "oi"}}},
			Output:  dag.HandoffContract{Artifact: "OK.txt", Gate: "true"},
		}},
	}
	if err := store.Save(subVars); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(subScript); err != nil {
		t.Fatal(err)
	}
}

func testHarnessGatewayAndStderr(t *testing.T) {
	repo := initTestGitRepo(t)
	defer os.RemoveAll(repo)
	ox := &orchestrix.Orchestrix{
		Inventory: provider.Inventory{
			RouterURL: "https://untrusted.local",
			CLIs:      []provider.CLI{{Name: "omp", Path: os.Args[0]}},
		},
		RouterKey:     "chave-secreta-gw",
		RouterKeyURL:  "https://trusted.local",
		EstadoHarness: func(_, _ string) string { return "homologado" },
	}
	ox.PermitirUnitNos()
	unitNode := dag.Node{
		ID:      "u_gw",
		Type:    dag.NodeType("unit_no"),
		Role:    "autor",
		Slots:   &dag.Slots{Gates: []string{"exit_code == 0"}},
		Routing: dag.RoutingConfig{Mode: dag.ModeManual, Harness: "omp", Via: "9router", Model: "omp:x"},
		Output:  dag.HandoffContract{Artifact: "out.json", Gate: "exit_code == 0"},
	}
	_, err := ox.Execute(context.Background(), unitNode, repo)
	if err == nil || !strings.Contains(err.Error(), "base confiável") {
		t.Fatalf("esperava validação envGateway em unit_harness, obteve: %v", err)
	}
	ox.Inventory.RouterURL = "https://trusted.local"
	_, err = ox.Execute(context.Background(), unitNode, repo)
	if err == nil || !strings.Contains(err.Error(), "stderr:") {
		t.Fatalf("esperava captura de stderr no backend unit_harness, obteve: %v", err)
	}
	if strings.Contains(err.Error(), "chave-secreta-gw") {
		t.Fatalf("chave vazou no erro: %v", err)
	}
}

func setupTestSkill(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "skills", "interview")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p, "SKILL.md"), []byte("Instrucao de teste."), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func mustParseDAG(t *testing.T, raw string) *dag.DAG {
	t.Helper()
	d, err := dag.Parse([]byte(raw))
	if err != nil {
		t.Fatalf("Parse DAG: %v", err)
	}
	return d
}
