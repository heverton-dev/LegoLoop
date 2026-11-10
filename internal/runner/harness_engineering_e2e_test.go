package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/dag"
)

type e2eHarnessExecutor struct {
	payloads map[string]string
}

func (e *e2eHarnessExecutor) Execute(ctx context.Context, n dag.Node, workdir string) (Result, error) {
	if payload, ok := e.payloads[n.ID]; ok {
		p := filepath.Join(workdir, filepath.FromSlash(n.Output.Artifact))
		if err := os.WriteFile(p, []byte(payload), 0o644); err != nil {
			return Result{}, err
		}
	}
	notes := []string{}
	if n.Routing.Effort != "" {
		notes = append(notes, "effort "+string(n.Routing.Effort)+" aplicado")
	}
	return Result{ExitCode: 0, Model: "harness-e2e/" + n.ID, Notes: notes}, nil
}

func TestHarnessEngineeringE2E(t *testing.T) {
	dagJSON := `{
		"name": "harness-engineering-processo-e2e",
		"version": "1.0.0",
		"nodes": [
			{
				"id": "node_spec",
				"type": "script",
				"command": ["go", "version"],
				"routing": {
					"mode": "manual",
					"model": "flash",
					"effort": "low"
				},
				"output": {
					"artifact": "spec.json",
					"schema": "schemas/spec.schema.json",
					"gate": "exit_code == 0"
				}
			},
			{
				"id": "node_deploy",
				"type": "script",
				"command": ["go", "version"],
				"depends_on": ["node_spec"],
				"routing": {
					"mode": "manual",
					"model": "sonnet",
					"effort": "high"
				},
				"output": {
					"artifact": "release.txt",
					"gate": "schema_valid('schemas/spec.schema.json', 'spec.json') && interactive('Aprovar promoção do artefato para prod?') && exit_code == 0"
				}
			}
		]
	}`

	schemaContent := `{
		"type": "object",
		"required": ["version", "status", "features"],
		"properties": {
			"version": {"type": "string"},
			"status": {"type": "string"},
			"features": {"type": "array"}
		}
	}`

	validSpecPayload := `{"version": "1.0.0", "status": "approved", "features": ["schema", "effort", "interactive"]}`
	invalidSpecPayload := `{"version": "1.0.0", "status": "draft"}` // faltando "features"

	setupWorkdir := func(t *testing.T) string {
		t.Helper()
		dir := t.TempDir()
		schemaDir := filepath.Join(dir, "schemas")
		if err := os.MkdirAll(schemaDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(schemaDir, "spec.schema.json"), []byte(schemaContent), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	d, err := dag.Parse([]byte(dagJSON))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	// Cenário 1: Sucesso Completo (Schema Válido + Interactive Gate Aprovado com --yes)
	t.Run("E2E_Success_With_Schema_And_AutoApprove", func(t *testing.T) {
		workdir := setupWorkdir(t)
		exec := &e2eHarnessExecutor{
			payloads: map[string]string{
				"node_spec":   validSpecPayload,
				"node_deploy": "Release v1.0.0 gerada com sucesso",
			},
		}

		rep, err := Run(context.Background(), d, Options{
			Workdir:     workdir,
			Executor:    exec,
			AutoApprove: true,
		})
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if !rep.OK {
			t.Fatalf("Processo deveria passar com sucesso total, relatório: %+v", rep)
		}
		if len(rep.Nodes) != 2 {
			t.Fatalf("esperava 2 nós reportados, obteve %d", len(rep.Nodes))
		}
		if rep.Nodes[0].Status != StatusPassed || rep.Nodes[1].Status != StatusPassed {
			t.Fatalf("ambos nós deveriam ter StatusPassed: %+v", rep.Nodes)
		}
		if len(rep.Nodes[0].Notes) == 0 || !strings.Contains(rep.Nodes[0].Notes[0], "effort low") {
			t.Errorf("telemetria de effort low não encontrada no node_spec: %v", rep.Nodes[0].Notes)
		}
		if len(rep.Nodes[1].Notes) == 0 || !strings.Contains(rep.Nodes[1].Notes[0], "effort high") {
			t.Errorf("telemetria de effort high não encontrada no node_deploy: %v", rep.Nodes[1].Notes)
		}
	})

	// Cenário 2: Falha no Schema Enforcement (node_spec gera JSON incompatível -> bloqueia node_deploy)
	t.Run("E2E_SchemaEnforcement_Failure_Blocks_Downstream", func(t *testing.T) {
		workdir := setupWorkdir(t)
		exec := &e2eHarnessExecutor{
			payloads: map[string]string{
				"node_spec":   invalidSpecPayload,
				"node_deploy": "Release não deve ser executada",
			},
		}

		rep, err := Run(context.Background(), d, Options{
			Workdir:     workdir,
			Executor:    exec,
			AutoApprove: true,
		})
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if rep.OK {
			t.Fatal("Processo deveria reprovar por Schema Enforcement")
		}
		if rep.Nodes[0].Status != StatusGateFailed {
			t.Fatalf("node_spec deveria ter StatusGateFailed, obteve: %v", rep.Nodes[0].Status)
		}
		if !strings.Contains(rep.Nodes[0].Error, "schema enforcement") {
			t.Fatalf("mensagem de erro deve citar schema enforcement: %q", rep.Nodes[0].Error)
		}
		if rep.Nodes[1].Status != StatusSkipped {
			t.Fatalf("node_deploy deveria ter sido pulado (StatusSkipped), obteve: %v", rep.Nodes[1].Status)
		}
	})

	// Cenário 3: Rejeição no Stop-Gate Interativo
	t.Run("E2E_InteractiveGate_Rejected_By_Operator", func(t *testing.T) {
		workdir := setupWorkdir(t)
		exec := &e2eHarnessExecutor{
			payloads: map[string]string{
				"node_spec":   validSpecPayload,
				"node_deploy": "Release abortada",
			},
		}

		prompterReject := &runnerMockPrompter{answer: false}
		rep, err := Run(context.Background(), d, Options{
			Workdir:  workdir,
			Executor: exec,
			Prompter: prompterReject,
		})
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		if rep.OK {
			t.Fatal("Processo deveria reprovar após rejeição humana no gate")
		}
		if rep.Nodes[0].Status != StatusPassed {
			t.Fatalf("node_spec deveria passar: %v", rep.Nodes[0].Status)
		}
		if rep.Nodes[1].Status != StatusGateFailed {
			t.Fatalf("node_deploy deveria ter StatusGateFailed por recusa humana: %v", rep.Nodes[1].Status)
		}
		if prompterReject.prompt != "Aprovar promoção do artefato para prod?" {
			t.Fatalf("mensagem do prompt inesperada: %q", prompterReject.prompt)
		}
	})
}
