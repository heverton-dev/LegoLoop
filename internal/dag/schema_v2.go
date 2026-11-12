package dag

import (
	"errors"
	"fmt"
	"strings"
)

// ErrCampoV2 indica um campo obrigatório ausente num DAG com schema_version 2.
var ErrCampoV2 = errors.New("dag: campo obrigatório do schema_version 2 ausente")

const avisoLegado = "DAG sem schema_version: modo legado, campos do contrato v2 não exigidos"

// AvisoLegado devolve o aviso do modo legado, ou vazio quando o DAG declara schema_version.
func (d *DAG) AvisoLegado() string {
	if d.SchemaVersion == 0 {
		return avisoLegado
	}
	return ""
}

func (d *DAG) checkSchemaV2() error {
	switch d.SchemaVersion {
	case 0:
		return nil
	case 2:
	default:
		return fmt.Errorf("dag: schema_version %d não suportado (use 2)", d.SchemaVersion)
	}
	if campo := campoTopoAusente(d); campo != "" {
		return fmt.Errorf("%w: topo: %s", ErrCampoV2, campo)
	}
	for _, n := range d.Nodes {
		if campo := campoNoAusente(n); campo != "" {
			return fmt.Errorf("%w: nó %q: %s", ErrCampoV2, n.ID, campo)
		}
	}
	return nil
}

func campoTopoAusente(d *DAG) string {
	switch {
	case strings.TrimSpace(d.Date) == "":
		return "date"
	case d.ExecutionStrategy == nil:
		return "execution_strategy"
	case d.UpstreamSignatures == nil:
		return "upstream_signatures"
	}
	return ""
}

func campoNoAusente(n Node) string {
	campos := []struct {
		nome  string
		vazio bool
	}{
		{"name", strings.TrimSpace(n.Name) == ""},
		{"nature", n.Nature == ""},
		{"context_in", n.ContextIn == nil},
		{"context_out", n.ContextOut == nil},
		{"isolation.write_boundary", n.Isolation["write_boundary"] == nil},
		{"limits", n.Limits == nil},
		{"rollback_strategy", n.RollbackStrategy == ""},
		{"tooling", n.Tooling == nil},
		{"verification", n.Verification == nil},
		{"spec_ref", n.SpecRef == ""},
		{"status", n.Status == ""},
		{"routing.harness", n.Routing.Harness == ""},
		{"routing.model", n.Routing.Model == ""},
		{"routing.effort", n.Routing.Effort == ""},
	}
	for _, c := range campos {
		if c.vazio {
			return c.nome
		}
	}
	return ""
}
