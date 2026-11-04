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
	finalGateID    = "SG-FINAL"
	finalBinary    = "bin/legoloop.exe"
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
	SubagentsAllowed bool   `json:"subagents_allowed"`
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
	Max int `json:"max_output_tokens"`
}

// routeNode segue, campo a campo e na mesma ordem, o nó do TEMPLATE_ROUTE.json.
type routeNode struct {
	ID               string         `json:"id"`
	Label            string         `json:"label"`
	Nature           string         `json:"nature"`
	Type             string         `json:"type"`
	Ref              *string        `json:"ref"`
	Status           string         `json:"status"`
	ExecutionMode    string         `json:"execution_mode"`
	Isolation        routeIsolation `json:"isolation"`
	TimeoutSec       int            `json:"timeout_sec"`
	TokenBudget      tokenBudget    `json:"token_budget"`
	DependsOn        []string       `json:"depends_on"`
	Routing          routeRouting   `json:"routing"`
	Tooling          routeTooling   `json:"tooling"`
	Inputs           routeInputs    `json:"inputs"`
	SpecRef          *string        `json:"spec_ref"`
	TestTarget       *string        `json:"test_target"`
	Defects          []any          `json:"defects"`
	RollbackStrategy string         `json:"rollback_strategy"`
	Output           routeOutput    `json:"output"`

	epic string // agrupamento no PLAN; não vai para o JSON
}

type routeStrategy struct {
	SupervisionLevel   string         `json:"supervision_level"`
	DefaultIsolation   routeIsolation `json:"default_isolation"`
	MaxParallelWorkers int            `json:"max_parallel_workers"`
	TDDEnforced        bool           `json:"tdd_enforced"`
	DefaultTimeoutSec  int            `json:"default_timeout_sec"`
	RollbackOnFailure  string         `json:"rollback_on_failure"`
}

type routeDoc struct {
	Schema            string        `json:"$schema"`
	TemplateVersion   string        `json:"template_version"`
	Name              string        `json:"name"`
	Version           string        `json:"version"`
	Date              string        `json:"date"`
	SessionID         *string       `json:"session_id"`
	IntentRef         *string       `json:"intent_ref"`
	SpecRef           string        `json:"spec_ref"`
	PlanRef           string        `json:"plan_ref"`
	ExecutionStrategy routeStrategy `json:"execution_strategy"`
	Nodes             []routeNode   `json:"nodes"`
}

func str(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func newRouteDoc(m planMeta, nodes []routeNode) routeDoc {
	return routeDoc{
		Schema: routeSchemaURL, TemplateVersion: templates.Version, Name: m.RouteName, Version: "1.0.0", Date: m.Date,
		IntentRef: str(m.IntentRef), SpecRef: m.SpecRel, PlanRef: m.PlanRel,
		ExecutionStrategy: routeStrategy{
			SupervisionLevel: string(m.Level), DefaultIsolation: localIsolation("orchestrix_hermetic"),
			MaxParallelWorkers: 2, TDDEnforced: true, DefaultTimeoutSec: defaultTimeout, RollbackOnFailure: "abort_pipeline",
		},
		Nodes: nodes,
	}
}

func (d routeDoc) marshal() ([]byte, error) { return json.MarshalIndent(d, "", "  ") }

// localIsolation: subagentes de harness são proibidos (AGENTS.md, regra 4), em qualquer nó.
func localIsolation(executor string) routeIsolation {
	return routeIsolation{Workspace: "local", Executor: executor, SubagentsAllowed: false}
}

// routingFor devolve roteamento e teto de saída por natureza, como no exemplo do template.
func routingFor(nature string) (routeRouting, int) {
	effort, model, budget := "high", "code-pro", 600
	if nature == "nano" {
		effort, model, budget = "low", "code-nano", 300
	}
	return routeRouting{Mode: "suggested", Harness: "9router", Model: model, Effort: &effort,
		Fallback: &routeFallback{Harness: "antigravity", Model: "flash"}}, budget
}

func ticketNode(t Ticket, m planMeta) routeNode {
	routing, budget := routingFor(t.Nature)
	verify, g := verifyAndGate(t.ValidationCmd, t.TargetFiles)
	if t.Gate != "" {
		g = t.Gate
	}
	return routeNode{
		ID: t.ID, Label: t.Title, Nature: t.Nature, Type: "unit_task", Status: "PENDING", ExecutionMode: t.ExecutionMode,
		Isolation: localIsolation("orchestrix_hermetic"), TimeoutSec: defaultTimeout, TokenBudget: tokenBudget{budget},
		DependsOn: nonNil(t.DependsOn), Routing: routing,
		Tooling: routeTooling{Runtime: "powershell", McpsAllowed: []string{}},
		Inputs:  routeInputs{SourceFiles: sourceFiles(t.TargetFiles), SpecAnchor: specAnchor(m, t), ContextArtifacts: []string{}},
		SpecRef: str(strings.Join(t.SpecCriteria, ", ")), TestTarget: testTarget(t.TargetFiles),
		Defects: []any{}, RollbackStrategy: "stash_changes",
		Output: routeOutput{Artifacts: nonNil(t.TargetFiles), Verify: verify, Gate: g},
		epic:   t.EpicID,
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
	return routeNode{
		ID: id, Label: label, Nature: "nano", Type: "stop_gate", Status: "PENDING", ExecutionMode: "sequential",
		Isolation: localIsolation("local_process"), TimeoutSec: defaultTimeout, TokenBudget: tokenBudget{0},
		DependsOn: nonNil(deps), Routing: routing, Tooling: routeTooling{Runtime: "powershell", McpsAllowed: []string{}},
		Inputs: in, Defects: []any{}, RollbackStrategy: "abort_pipeline", Output: out, epic: epic,
	}
}
