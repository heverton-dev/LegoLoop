package main

import (
	"fmt"
	"maps"
	"strings"
)

// projetarNo converte um nó da ROUTE em nó de DAG.
func projetarNo(n routeNode) (map[string]any, error) {
	if n.ID == "" {
		return nil, fmt.Errorf("nó sem id")
	}
	if n.Type != "unit_task" && n.Type != "script" && n.Type != "stop_gate" {
		return nil, fmt.Errorf("nó %q: tipo %q não suportado", n.ID, n.Type)
	}
	nd := map[string]any{
		"id":         n.ID,
		"name":       strings.ToLower(strings.NewReplacer("-", "_", ".", "_").Replace(n.ID)),
		"type":       n.Type,
		"nature":     n.Nature,
		"depends_on": n.DependsOn,
		"isolation":  n.Isolation,
		"output":     map[string]any{"artifact": "scratch/gates/" + n.ID + ".json", "gate": n.Output.Gate},
		"status":     n.Status,
	}
	if n.Type == "stop_gate" {
		nd["routing"] = map[string]any{"mode": "manual", "model": "human/operator"}
		if n.SpecRef != nil {
			nd["spec_ref"] = *n.SpecRef
		}
		return nd, nil
	}
	if n.Routing.Harness == "direto" && n.Routing.Via != "9router" {
		return nil, fmt.Errorf("nó %q: harness direto exige via 9router, recebeu %q", n.ID, n.Routing.Via)
	}
	nd["routing"] = roteamento(n)
	unidade(nd, n)
	return nd, nil
}

func roteamento(n routeNode) map[string]any {
	rt := map[string]any{"mode": n.Routing.Mode, "model": n.Routing.Model}
	if n.Routing.Via != "" {
		rt["via"] = n.Routing.Via
	}
	if n.Routing.Effort != nil {
		rt["effort"] = *n.Routing.Effort
	}
	if n.Routing.Fallback != nil {
		rt["fallback"] = n.Routing.Fallback.Model
	}
	return rt
}

// unidade preenche os campos próprios do nó de código.
func unidade(nd map[string]any, n routeNode) {
	tooling := map[string]any{"harness": n.Routing.Harness}
	maps.Copy(tooling, n.Tooling)
	nd["tooling"] = tooling
	nd["context_in"] = contextoEntrada(n)
	nd["context_out"] = n.Output.Artifacts
	nd["limits"] = map[string]any{"max_duration_sec": n.TimeoutSec,
		"max_prompt_tokens": n.Budget.In, "max_completion_tokens": n.Budget.Out}
	if n.SpecRef != nil {
		nd["spec_ref"] = *n.SpecRef
	}
	if n.RollbackStrategy != "" {
		nd["rollback_strategy"] = n.RollbackStrategy
	}
	if len(n.Output.Verify) > 0 {
		nd["command"] = n.Output.Verify[0]
		nd["verification"] = map[string]any{"test_command": strings.Join(n.Output.Verify[0], " "), "exit_code": 0}
	}
	if n.TestTarget != nil {
		file, sig, _ := strings.Cut(*n.TestTarget, ":")
		nd["test_seam"] = map[string]string{"type": "highest_seam", "signature": sig, "file": file}
	}
}

func contextoEntrada(n routeNode) []string {
	in := append([]string{}, n.Inputs.SourceFiles...)
	if n.Inputs.SpecAnchor != nil {
		path, _, _ := strings.Cut(*n.Inputs.SpecAnchor, "#")
		in = append(in, path)
	}
	return in
}
