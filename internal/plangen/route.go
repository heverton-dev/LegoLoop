package plangen

import (
	"encoding/json"
	"fmt"
	"strings"

	"legoloop/docs/templates"
)

const (
	routeSchemaURL = "https://legoloop.dev/schemas/route.v2.json"
	defaultTimeout = 180

	contractVersion = 2
	finalGateID     = "SG-FINAL"
	finalBinary     = "bin/legoloop.exe"
	tipoUnidade     = "unit_task"
	rollbackAborta  = "abort_pipeline"
)

type routeFallback struct {
	Harness string `json:"harness"`
	Model   string `json:"model"`
}

type routeRouting struct {
	Mode     string         `json:"mode"`
	Harness  string         `json:"harness"`
	Model    string         `json:"model"`
	Effort   *string        `json:"effort"`
	Fallback *routeFallback `json:"fallback"`
}

type routeIsolation struct {
	Workspace        string `json:"workspace"`
	Executor         string `json:"executor"`
	SubagentsAllowed any    `json:"subagents_allowed"` // false ou "ephemeral" (AGENTS.md, regra 4)
	WriteBoundary    string `json:"write_boundary"`
}

type routeRetryPolicy struct {
	MaxRetries           int  `json:"max_retries"`
	BackoffMS            int  `json:"backoff_ms"`
	FallbackModelOnRetry bool `json:"fallback_model_on_retry"`
}

type routeTooling struct {
	Runtime      string   `json:"runtime"`
	McpsRequired bool     `json:"mcps_required"`
	McpsAllowed  []string `json:"mcps_allowed"`
}

type routeInputs struct {
	SourceFiles      []string `json:"source_files"`
	SpecAnchor       *string  `json:"spec_anchor"`
	ContextArtifacts []string `json:"context_artifacts"`
}

type routeOutput struct {
	Artifacts []string   `json:"artifacts"`
	Verify    [][]string `json:"verify"`
	Gate      string     `json:"gate"`
}

type tokenBudget struct {
	MaxInputTokens  int `json:"max_input_tokens"`
	MaxOutputTokens int `json:"max_output_tokens"`
}

type routeTelemetry struct {
	StartedAt      *string        `json:"started_at"`
	FinishedAt     *string        `json:"finished_at"`
	DurationMS     *int64         `json:"duration_ms"`
	TokensConsumed map[string]int `json:"tokens_consumed"`
	ExitCode       *int           `json:"exit_code"`
	GitCommit      *string        `json:"git_commit"`
	ErrorMessage   *string        `json:"error_message"`
}

// routeNode segue, campo a campo e na mesma ordem, o nó do TEMPLATE_ROUTE.json.
type routeNode struct {
	ID               string           `json:"id"`
	InitiativeID     *string          `json:"initiative_id"`
	EpicID           *string          `json:"epic_id"`
	Label            string           `json:"label"`
	Nature           string           `json:"nature"`
	Extreme          bool             `json:"extreme"` // true libera Opus, só com aprovação humana
	ClassRule        string           `json:"class_rule"`
	Type             string           `json:"type"`
	Ref              *string          `json:"ref"`
	Status           string           `json:"status"`
	ExecutionMode    string           `json:"execution_mode"`
	Isolation        routeIsolation   `json:"isolation"`
	TimeoutSec       int              `json:"timeout_sec"`
	RetryPolicy      routeRetryPolicy `json:"retry_policy"`
	TokenBudget      tokenBudget      `json:"token_budget"`
	DependsOn        []string         `json:"depends_on"`
	Routing          routeRouting     `json:"routing"`
	Tooling          routeTooling     `json:"tooling"`
	Inputs           routeInputs      `json:"inputs"`
	SpecRef          *string          `json:"spec_ref"`
	TestTarget       *string          `json:"test_target"`
	Defects          []any            `json:"defects"`
	RollbackStrategy string           `json:"rollback_strategy"`
	Output           routeOutput      `json:"output"`
	Telemetry        *routeTelemetry  `json:"telemetry"`

	epic string // agrupamento no PLAN; não vai para o JSON
}

// ConcurrencyStrategy governa a execução paralela de subagentes efêmeros limpos.
type ConcurrencyStrategy struct {
	Mode                 string `json:"mode"`
	MaxParallelSubagents int    `json:"max_parallel_subagents"`
	HarnessCapability    string `json:"harness_capability"`
	RateLimitThrottleMS  int    `json:"rate_limit_throttle_ms"`
}

// RollbackPolicy define as ações mecânicas em caso de falha de um nó.
type RollbackPolicy struct {
	Action            string `json:"action"`
	GitStrategy       string `json:"git_strategy"`
	PreserveArtifacts bool   `json:"preserve_artifacts"`
}

type routeStrategy struct {
	SupervisionLevel   string              `json:"supervision_level"`
	Concurrency        ConcurrencyStrategy `json:"concurrency"`
	DefaultIsolation   routeIsolation      `json:"default_isolation"`
	MaxParallelWorkers int                 `json:"max_parallel_workers"`
	TDDEnforced        bool                `json:"tdd_enforced"`
	DefaultTimeoutSec  int                 `json:"default_timeout_sec"`
	RollbackOnFailure  string              `json:"rollback_on_failure"`
	RollbackPolicy     RollbackPolicy      `json:"rollback_policy"`
}

type routeDoc struct {
	Schema             string              `json:"$schema"`
	SchemaVersion      int                 `json:"schema_version"`
	TemplateVersion    string              `json:"template_version"`
	Name               string              `json:"name"`
	Version            string              `json:"version"`
	Date               string              `json:"date"`
	SessionID          *string             `json:"session_id"`
	IntentRef          *string             `json:"intent_ref"`
	SpecRef            string              `json:"spec_ref"`
	PlanRef            string              `json:"plan_ref"`
	InitiativeRef      *string             `json:"initiative_ref"`
	UpstreamSignatures *UpstreamSignatures `json:"upstream_signatures"`
	ExecutionStrategy  routeStrategy       `json:"execution_strategy"`
	Nodes              []routeNode         `json:"nodes"`
}

func str(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func defaultConcurrency() ConcurrencyStrategy {
	return ConcurrencyStrategy{
		Mode:                 "hybrid",
		MaxParallelSubagents: 2,
		HarnessCapability:    "subagent_cli_isolated",
		RateLimitThrottleMS:  500,
	}
}

func defaultRollbackPolicy() RollbackPolicy {
	return RollbackPolicy{
		Action:            "abort_and_clean",
		GitStrategy:       "discard_uncommitted",
		PreserveArtifacts: true,
	}
}

func defaultStrategy(level SupervisionLevel) routeStrategy {
	return routeStrategy{
		SupervisionLevel:   string(level),
		Concurrency:        defaultConcurrency(),
		DefaultIsolation:   localIsolation("orchestrix_hermetic", "strict"),
		MaxParallelWorkers: 2,
		TDDEnforced:        true,
		DefaultTimeoutSec:  defaultTimeout,
		RollbackOnFailure:  rollbackAborta,
		RollbackPolicy:     defaultRollbackPolicy(),
	}
}

func newRouteDoc(m planMeta, nodes []routeNode) routeDoc {
	sig := &UpstreamSignatures{
		Intent: FileSignatureFor(m.BaseDir, m.IntentRef),
		Spec:   *FileSignatureFor(m.BaseDir, m.SpecRel),
		Plan:   *FileSignatureFor(m.BaseDir, m.PlanRel),
	}
	return routeDoc{
		Schema: routeSchemaURL, SchemaVersion: contractVersion, TemplateVersion: templates.Version, Name: m.RouteName, Version: "1.0.0", Date: m.Date,
		IntentRef: str(m.IntentRef), SpecRef: m.SpecRel, PlanRef: m.PlanRel, InitiativeRef: str(m.InitiativeRef),
		UpstreamSignatures: sig,
		ExecutionStrategy:  defaultStrategy(m.Level),
		Nodes:              nodes,
	}
}

func (d routeDoc) marshal() ([]byte, error) { return json.MarshalIndent(d, "", "  ") }

// localIsolation: subagentes de harness são proibidos (AGENTS.md, regra 4), em qualquer nó.
func localIsolation(executor, writeBoundary string) routeIsolation {
	if writeBoundary == "" {
		writeBoundary = "strict"
	}
	return routeIsolation{Workspace: "local", Executor: executor, SubagentsAllowed: false, WriteBoundary: writeBoundary}
}

func defaultRetryPolicy(isStopGate bool) routeRetryPolicy {
	if isStopGate {
		return routeRetryPolicy{MaxRetries: 0, BackoffMS: 0, FallbackModelOnRetry: false}
	}
	return routeRetryPolicy{MaxRetries: 1, BackoffMS: 1000, FallbackModelOnRetry: true}
}

func tokenBudgetFor(nature string) (tokenBudget, routeRouting) {
	effort, model, out := "high", "code-pro", 600
	in := 16000
	if nature == "nano" {
		effort, model, out = "low", "code-fast", 300
		in = 8000
	} else if nature == "macro" {
		effort, model, out = "medium", "code-fast", 1200
		in = 32000
	}
	return tokenBudget{MaxInputTokens: in, MaxOutputTokens: out}, routeRouting{
		Mode: "suggested", Harness: "9router", Model: model, Effort: &effort,
		Fallback: &routeFallback{Harness: "antigravity", Model: "flash"},
	}
}

func resolveInitiative(t Ticket, m planMeta) *string {
	if t.InitiativeID != "" {
		return str(t.InitiativeID)
	}
	if iNum, ok := m.EpicInitiative[epicNumber(t.EpicID)]; ok && iNum > 0 {
		s := fmt.Sprintf("I%02d", iNum)
		return &s
	}
	return nil
}

func ticketNode(t Ticket, m planMeta) routeNode {
	budget, routing := tokenBudgetFor(t.Nature)
	verify, g := verifyAndGate(t.ValidationCmd, t.TargetFiles)
	if t.Gate != "" {
		g = t.Gate
	}
	return routeNode{
		ID: t.ID, InitiativeID: resolveInitiative(t, m), EpicID: str(t.EpicID), Label: t.Title, Nature: t.Nature,
		ClassRule: "nature=" + t.Nature, Type: tipoUnidade, Status: "PENDING", ExecutionMode: t.ExecutionMode,
		Isolation:   localIsolation("orchestrix_hermetic", "strict"),
		TimeoutSec:  defaultTimeout,
		RetryPolicy: defaultRetryPolicy(false),
		TokenBudget: budget,
		DependsOn:   nonNil(t.DependsOn), Routing: routing,
		Tooling:          routeTooling{Runtime: "powershell", McpsAllowed: []string{}},
		Inputs:           routeInputs{SourceFiles: sourceFiles(t.TargetFiles), SpecAnchor: specAnchor(m, t), ContextArtifacts: []string{}},
		SpecRef:          str(strings.Join(t.SpecCriteria, ", ")),
		TestTarget:       testTarget(t.TargetFiles),
		Defects:          []any{},
		RollbackStrategy: "stash_changes",
		Output:           routeOutput{Artifacts: nonNil(t.TargetFiles), Verify: verify, Gate: g},
		Telemetry:        nil,
		epic:             t.EpicID,
	}
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return append([]string{}, s...)
}

func isTestFile(p string) bool { return strings.HasSuffix(p, "_test.go") }

func sourceFiles(files []string) []string {
	var src []string
	for _, f := range files {
		if !isTestFile(f) {
			src = append(src, f)
		}
	}
	if len(src) == 0 {
		return nonNil(files)
	}
	return src
}

func testTarget(files []string) *string {
	for _, f := range files {
		if isTestFile(f) {
			return str(f)
		}
	}
	return nil
}

func specAnchor(m planMeta, t Ticket) *string {
	if len(t.SpecCriteria) == 0 {
		return nil
	}
	a := fmt.Sprintf("%s#%s", m.SpecRel, t.SpecCriteria[0])
	return &a
}

// humanGate é um Stop-Gate humano: executor local, roteamento manual, gate DSL prompt(...).
func humanGate(id, question string, deps []string, epic string) routeNode {
	return stopGate(id, "Stop-Gate Humano: "+id, deps, epic, routeOutput{
		Artifacts: []string{}, Verify: [][]string{}, Gate: "prompt('" + question + "')"},
		routeRouting{Mode: "manual", Harness: "local", Model: "human/operator"}, routeInputs{SourceFiles: []string{}, ContextArtifacts: []string{}})
}

// finalGate é o SG-FINAL do template: regressão global e recompilação fresca do binário.
func finalGate(deps []string) routeNode {
	return stopGate(finalGateID, "Homologação Final (Regressão Global & Recompilação Fresca)", deps, "", routeOutput{
		Artifacts: []string{finalBinary},
		Verify:    [][]string{{"go", "test", "./..."}, {"go", "vet", "./..."}, {"go", "build", "-o", finalBinary, "./cmd/legoloop"}},
		Gate:      "exit_code == 0 && file_exists('" + finalBinary + "')"},
		routeRouting{Mode: "manual", Harness: "local", Model: "script/local"}, routeInputs{SourceFiles: []string{"cmd/legoloop/main.go"}, ContextArtifacts: []string{}})
}

func stopGate(id, label string, deps []string, epic string, out routeOutput, routing routeRouting, in routeInputs) routeNode {
	wb := "none"
	if len(out.Artifacts) > 0 {
		wb = "strict"
	}
	return routeNode{
		ID: id, InitiativeID: nil, EpicID: str(epic), Label: label, Nature: "nano", ClassRule: "stop_gate", Type: "stop_gate", Status: "PENDING", ExecutionMode: "sequential",
		Isolation:   localIsolation("local_process", wb),
		TimeoutSec:  defaultTimeout,
		RetryPolicy: defaultRetryPolicy(true),
		TokenBudget: tokenBudget{MaxInputTokens: 0, MaxOutputTokens: 0},
		DependsOn:   nonNil(deps), Routing: routing,
		Tooling:          routeTooling{Runtime: "powershell", McpsAllowed: []string{}},
		Inputs:           in,
		Defects:          []any{},
		RollbackStrategy: rollbackAborta,
		Output:           out,
		Telemetry:        nil,
		epic:             epic,
	}
}
