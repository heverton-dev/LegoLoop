package studiodoc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/graph"
)

// Ports descreve os contratos de entrada e saída derivados de um loop.
type Ports struct {
	Consumes []string `json:"consumes"`
	Produces string   `json:"produces"`
	Artifact string   `json:"artifact"`
	Gate     string   `json:"gate"`
}

const AnyArtifact = "Artifact:Any"

// DerivePorts infere as portas do sub-loop a partir do artefato final e das raízes.
func DerivePorts(l *Processo) (*Ports, error) {
	if len(l.Nodes) == 0 {
		return nil, errors.New("loop não possui nós")
	}
	raw, err := json.Marshal(struct {
		Name    string       `json:"name"`
		Version string       `json:"version,omitempty"`
		Trigger *dag.Trigger `json:"trigger,omitempty"`
		Nodes   []dag.Node   `json:"nodes"`
	}{Name: l.Name, Version: l.Version, Trigger: l.Trigger, Nodes: l.Nodes})
	if err != nil {
		return nil, fmt.Errorf("serializar loop: %w", err)
	}
	d, err := dag.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("DAG inválido: %w", err)
	}
	return buildPorts(d), nil
}

func buildPorts(d *dag.DAG) *Ports {
	layers := d.Layers()
	lastLayer := layers[len(layers)-1]
	finalNodeID := lastLayer[len(lastLayer)-1]
	finalNode := d.Node(finalNodeID)

	produces := matchProduces(finalNode.Output.Artifact)
	consumes := deriveConsumes(d, layers[0])

	return &Ports{
		Consumes: consumes,
		Produces: produces,
		Artifact: finalNode.Output.Artifact,
		Gate:     finalNode.Output.Gate,
	}
}

func matchProduces(artifact string) string {
	for _, b := range graph.DefaultCatalog() {
		if strings.EqualFold(b.Artifact, artifact) {
			return string(b.Produces)
		}
	}
	return AnyArtifact
}

func deriveConsumes(d *dag.DAG, rootIDs []string) []string {
	var list []string
	seen := make(map[string]bool)
	for _, id := range rootIDs {
		n := d.Node(id)
		for _, t := range matchConsumes(n) {
			if !seen[t] {
				seen[t] = true
				list = append(list, t)
			}
		}
	}
	if len(list) == 0 {
		return []string{AnyArtifact}
	}
	return list
}

func matchConsumes(n *dag.Node) []string {
	for _, b := range graph.DefaultCatalog() {
		if b.ID == n.Ref || b.Ref == n.Ref {
			var out []string
			for _, c := range b.Consumes {
				out = append(out, string(c))
			}
			return out
		}
	}
	return []string{AnyArtifact}
}
