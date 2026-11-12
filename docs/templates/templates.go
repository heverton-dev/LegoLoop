// Package templates embute os templates canônicos (SSOT) de INTENT, SPEC, PLAN e ROUTE.
//
// Os motores que geram esses artefatos (interview, specgen, plangen) leem os templates
// daqui, e não de cópias em strings Go. Assim, o arquivo em docs/templates é a única fonte
// de verdade, e os testes de conformidade dos motores falham quando um deles diverge.
package templates

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Version é a versão dos templates (campo template_version do ROUTE e linha "Template:" dos documentos).
const Version = "2.2.0"

//go:embed TEMPLATE_INTENT.md
var intentMD string

//go:embed TEMPLATE_SPEC.md
var specMD string

//go:embed TEMPLATE_PLAN.md
var planMD string

//go:embed TEMPLATE_ROUTE.json
var routeJSON []byte

//go:embed TEMPLATE_REPORT.md
var reportMD string

//go:embed TEMPLATE_INTERVIEW.md
var interviewMD string

//go:embed TEMPLATE_HANDOFF.md
var handoffMD string

//go:embed TEMPLATE_DAG_VIDEOREA.json
var dagVideoReaJSON []byte

//go:embed TEMPLATE_DAG.json
var dagJSON []byte

// Intent devolve o TEMPLATE_INTENT.md.
func Intent() string { return intentMD }

// Spec devolve o TEMPLATE_SPEC.md.
func Spec() string { return specMD }

// Plan devolve o TEMPLATE_PLAN.md.
func Plan() string { return planMD }

// Route devolve o TEMPLATE_ROUTE.json.
func Route() []byte { return slices.Clone(routeJSON) }

// Report devolve o TEMPLATE_REPORT.md.
func Report() string { return reportMD }

// Interview devolve o TEMPLATE_INTERVIEW.md.
func Interview() string { return interviewMD }

// Handoff devolve o TEMPLATE_HANDOFF.md.
func Handoff() string { return handoffMD }

// DAGVideoREA devolve o TEMPLATE_DAG_VIDEOREA.json.
func DAGVideoREA() []byte { return slices.Clone(dagVideoReaJSON) }

// DAG devolve o TEMPLATE_DAG.json universal para qualquer processo.
func DAG() []byte { return slices.Clone(dagJSON) }

var (
	headingRe  = regexp.MustCompile(`(?m)^## (\d+\. .+?)\s*$`)
	guidanceRe = regexp.MustCompile(`(?s)<!--.*?-->\n?`)
)

// Headings lista, em ordem, os títulos de seção numerados (nível 2) de um documento.
func Headings(md string) []string {
	var out []string
	for _, m := range headingRe.FindAllStringSubmatch(md, -1) {
		out = append(out, m[1])
	}
	return out
}

// StripGuidance remove os comentários HTML de orientação ao autor.
func StripGuidance(md string) string {
	return strings.TrimSpace(guidanceRe.ReplaceAllString(md, "")) + "\n"
}

// RouteShape é o conjunto de chaves (ordenado) de cada nível do ROUTE v2.
type RouteShape struct {
	Top, Strategy, Node, Isolation, Routing, Tooling, Inputs, Output []string
}

// RouteKeys extrai do template as chaves que todo ROUTE gerado deve ter.
// Falha se os nós do template não forem isomórficos.
func RouteKeys() (RouteShape, error) {
	var doc map[string]any
	if err := json.Unmarshal(routeJSON, &doc); err != nil {
		return RouteShape{}, fmt.Errorf("templates: ROUTE inválido: %w", err)
	}
	nodes, _ := doc["nodes"].([]any)
	if len(nodes) == 0 {
		return RouteShape{}, fmt.Errorf("templates: ROUTE sem nós")
	}
	first, _ := nodes[0].(map[string]any)
	for i, n := range nodes {
		m, _ := n.(map[string]any)
		if !slices.Equal(keysOf(m), keysOf(first)) {
			return RouteShape{}, fmt.Errorf("templates: o nó %d do ROUTE não é isomórfico ao primeiro", i)
		}
	}
	return RouteShape{
		Top:       keysOf(doc),
		Strategy:  keysOf(asMap(doc["execution_strategy"])),
		Node:      keysOf(first),
		Isolation: keysOf(asMap(first["isolation"])),
		Routing:   keysOf(asMap(first["routing"])),
		Tooling:   keysOf(asMap(first["tooling"])),
		Inputs:    keysOf(asMap(first["inputs"])),
		Output:    keysOf(asMap(first["output"])),
	}, nil
}

func asMap(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}

func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
