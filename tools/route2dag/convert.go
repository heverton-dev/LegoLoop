package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type route struct {
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	Date               string               `json:"date"`
	UpstreamSignatures map[string]signature `json:"upstream_signatures"`
	ExecutionStrategy  map[string]any       `json:"execution_strategy"`
	Nodes              []routeNode          `json:"nodes"`
}

type signature struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type routeNode struct {
	ID         string         `json:"id"`
	Nature     string         `json:"nature"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	DependsOn  []string       `json:"depends_on"`
	Isolation  map[string]any `json:"isolation"`
	Tooling    map[string]any `json:"tooling"`
	TimeoutSec int            `json:"timeout_sec"`
	Budget     struct {
		In  int `json:"max_input_tokens"`
		Out int `json:"max_output_tokens"`
	} `json:"token_budget"`
	Routing struct {
		Mode     string  `json:"mode"`
		Harness  string  `json:"harness"`
		Model    string  `json:"model"`
		Effort   *string `json:"effort"`
		Fallback *struct {
			Model string `json:"model"`
		} `json:"fallback"`
	} `json:"routing"`
	Inputs struct {
		SourceFiles []string `json:"source_files"`
		SpecAnchor  *string  `json:"spec_anchor"`
	} `json:"inputs"`
	SpecRef    *string `json:"spec_ref"`
	TestTarget *string `json:"test_target"`
	Output     struct {
		Artifacts []string   `json:"artifacts"`
		Verify    [][]string `json:"verify"`
		Gate      string     `json:"gate"`
	} `json:"output"`
	RollbackStrategy string `json:"rollback_strategy"`
}

// estratégia: só estas chaves da execution_strategy da ROUTE vão para o DAG.
var chavesEstrategia = []string{"supervision_level", "concurrency", "default_isolation", "max_parallel_workers",
	"tdd_enforced", "default_timeout_sec", "rollback_on_failure", "rollback_policy"}

// Converter projeta o JSON de uma ROUTE em JSON de DAG. Devolve também o nome
// de arquivo canônico (<data>_DAG_<NOME>.json).
func Converter(raw []byte) ([]byte, string, error) {
	var r route
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := dec.Decode(&r); err != nil {
		return nil, "", fmt.Errorf("ROUTE inválida: %w", err)
	}
	if r.Name == "" || len(r.Nodes) == 0 {
		return nil, "", errors.New("ROUTE sem name ou sem nodes")
	}
	nodes := make([]map[string]any, 0, len(r.Nodes))
	for _, n := range r.Nodes {
		nd, err := projetarNo(n)
		if err != nil {
			return nil, "", err
		}
		nodes = append(nodes, nd)
	}
	d := map[string]any{
		"$schema":             "https://legoloop.dev/schemas/dag.v1.json",
		"name":                "legoloop-" + strings.ToLower(strings.ReplaceAll(r.Name, "_", "-")),
		"version":             r.Version,
		"date":                r.Date,
		"description":         "Espelho executável da ROUTE " + r.Name,
		"upstream_signatures": assinaturas(r),
		"execution_strategy":  estrategia(r.ExecutionStrategy),
		"trigger":             map[string]any{"type": "manual"},
		"nodes":               nodes,
	}
	out, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, "", err
	}
	return append(out, '\n'), fmt.Sprintf("%s_DAG_%s.json", r.Date, r.Name), nil
}

func assinaturas(r route) map[string]any {
	out := map[string]any{}
	for k, v := range r.UpstreamSignatures {
		out[k] = v.SHA256
	}
	return out
}

func estrategia(src map[string]any) map[string]any {
	out := map[string]any{}
	for _, k := range chavesEstrategia {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}
