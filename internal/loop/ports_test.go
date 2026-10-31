package loop_test

import (
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/graph"
	"legoloop/internal/loop"
)

func TestDerivePortsClassicLoop(t *testing.T) {
	// Loop clássico: interview (Consumes Idea, Produces Briefing) -> aidd-spec (Consumes Briefing, Produces Spec)
	lp := &loop.Loop{
		ID:      "idea-to-spec",
		Name:    "Idea to Spec",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:   "n1",
				Type: dag.TypeSkill,
				Ref:  "skills/interview",
				Routing: dag.RoutingConfig{
					Mode:  dag.ModeManual,
					Model: "llama3.1:8b",
				},
				Output: dag.HandoffContract{
					Artifact: "BRIEFING.md",
					Gate:     "file_exists('BRIEFING.md')",
				},
			},
			{
				ID:        "n2",
				Type:      dag.TypeSkill,
				Ref:       "skills/aidd-spec",
				DependsOn: []string{"n1"},
				Routing: dag.RoutingConfig{
					Mode:  dag.ModeManual,
					Model: "llama3.1:8b",
				},
				Output: dag.HandoffContract{
					Artifact: "SPEC.md",
					Gate:     "has_binary_criteria('SPEC.md')",
				},
			},
		},
	}

	ports, err := loop.DerivePorts(lp)
	if err != nil {
		t.Fatalf("DerivePorts falhou: %v", err)
	}

	if ports.Produces != string(graph.Spec) {
		t.Errorf("Produces esperado %s, obteve %s", graph.Spec, ports.Produces)
	}
	if ports.Artifact != "SPEC.md" {
		t.Errorf("Artifact esperado SPEC.md, obteve %s", ports.Artifact)
	}
	if ports.Gate != "has_binary_criteria('SPEC.md')" {
		t.Errorf("Gate esperado tem has_binary_criteria, obteve %s", ports.Gate)
	}
	if len(ports.Consumes) == 0 || ports.Consumes[0] != string(graph.Idea) {
		t.Errorf("Consumes esperado [Artifact:Idea], obteve %v", ports.Consumes)
	}
}

func TestDerivePortsCustomGenericLoop(t *testing.T) {
	// Loop genérico com artefato customizado
	lp := &loop.Loop{
		ID:      "custom-pipeline",
		Name:    "Custom Pipeline",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:   "n1",
				Type: dag.TypeScript,
				Command: []string{"python", "process.py"},
				Routing: dag.RoutingConfig{
					Mode:  dag.ModeManual,
					Model: "local",
				},
				Output: dag.HandoffContract{
					Artifact: "CUSTOM_OUTPUT.json",
					Gate:     "file_exists('CUSTOM_OUTPUT.json')",
				},
			},
		},
	}

	ports, err := loop.DerivePorts(lp)
	if err != nil {
		t.Fatalf("DerivePorts falhou: %v", err)
	}

	if ports.Produces != "Artifact:Any" {
		t.Errorf("Produces esperado Artifact:Any, obteve %s", ports.Produces)
	}
	if ports.Artifact != "CUSTOM_OUTPUT.json" {
		t.Errorf("Artifact esperado CUSTOM_OUTPUT.json, obteve %s", ports.Artifact)
	}
	if len(ports.Consumes) == 0 || ports.Consumes[0] != "Artifact:Any" {
		t.Errorf("Consumes esperado [Artifact:Any], obteve %v", ports.Consumes)
	}
}

func TestDerivePortsEmptyLoop(t *testing.T) {
	lp := &loop.Loop{
		ID:      "empty",
		Name:    "Empty",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes:   nil,
	}

	_, err := loop.DerivePorts(lp)
	if err == nil {
		t.Fatal("esperava erro para loop sem nós")
	}
}
