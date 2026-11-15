package orchestrix

import (
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/studiodoc"
)

func fixtureProcessoV2() *studiodoc.Processo {
	return &studiodoc.Processo{
		ID:            "sub-v2",
		Name:          "Sub Processo V2",
		Description:   "Descricao preservada do subprocesso",
		Version:       "2.0.0",
		SchemaVersion: 2,
		Date:          "2026-11-14",
		ExecutionStrategy: map[string]any{
			"supervision_level": "epic",
		},
		UpstreamSignatures: map[string]any{
			"route": map[string]any{"sha256": "abc123"},
		},
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   []dag.Node{fixtureNodeV2()},
	}
}

func fixtureNodeV2() dag.Node {
	return dag.Node{
		ID:     "n1",
		Name:   "No V2",
		Type:   dag.TypeSkill,
		Nature: "micro",
		Role:   "engenheiro-backend",
		Ref:    "skills/exemplo",
		Routing: dag.RoutingConfig{
			Mode: dag.ModeSuggested, Harness: "9router", Via: "9router",
			Model: "t-codigo-128k", Effort: "high",
		},
		Isolation:        map[string]any{"write_boundary": "strict"},
		Limits:           map[string]any{"max_output_tokens": 300},
		RollbackStrategy: "abort_processo",
		Tooling:          map[string]any{"runtime": "powershell"},
		Verification:     map[string]any{"gate": "exit_code == 0"},
		SpecRef:          "AC-1",
		Status:           "PENDING",
		ContextIn:        []string{"in.txt"},
		ContextOut:       []string{"out.txt"},
		Output:           dag.HandoffContract{Artifact: "out.txt", Gate: "file_exists('out.txt')"},
	}
}

func TestParseSubDAG_PreservaMetadados(t *testing.T) {
	proc := fixtureProcessoV2()
	parsed, err := parseSubDAG(proc)
	if err != nil {
		t.Fatalf("parseSubDAG falhou: %v", err)
	}
	if parsed.Description != proc.Description {
		t.Errorf("Description = %q, want %q", parsed.Description, proc.Description)
	}
	if parsed.SchemaVersion != proc.SchemaVersion {
		t.Errorf("SchemaVersion = %d, want %d", parsed.SchemaVersion, proc.SchemaVersion)
	}
	if parsed.Date != proc.Date {
		t.Errorf("Date = %q, want %q", parsed.Date, proc.Date)
	}
	if parsed.ExecutionStrategy == nil || parsed.ExecutionStrategy["supervision_level"] != "epic" {
		t.Errorf("ExecutionStrategy nao preservada: %+v", parsed.ExecutionStrategy)
	}
	if parsed.UpstreamSignatures == nil || parsed.UpstreamSignatures["route"] == nil {
		t.Errorf("UpstreamSignatures nao preservadas: %+v", parsed.UpstreamSignatures)
	}
}
