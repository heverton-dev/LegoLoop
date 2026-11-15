package studiodoc

import (
	"errors"
	"regexp"
	"time"

	"legoloop/internal/dag"
)

var (
	ErrNotFound  = errors.New("Processo não encontrado")
	ErrInvalidID = errors.New("identificador de Processo inválido")
	ErrExists    = errors.New("Processo já existe")
)

var idRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,63}$`)

func ValidateID(id string) error {
	if !idRegex.MatchString(id) {
		return ErrInvalidID
	}
	return nil
}

type Processo struct {
	Schema             string         `json:"$schema,omitempty"`
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	Version            string         `json:"version,omitempty"`
	SchemaVersion      int            `json:"schema_version,omitempty"`
	Date               string         `json:"date,omitempty"`
	UpstreamSignatures map[string]any `json:"upstream_signatures,omitempty"`
	ExecutionStrategy  map[string]any `json:"execution_strategy,omitempty"`
	Trigger            *dag.Trigger   `json:"trigger,omitempty"`
	Nodes              []dag.Node     `json:"nodes"`
	CreatedAt          time.Time      `json:"created_at,omitempty"`
	UpdatedAt          time.Time      `json:"updated_at,omitempty"`
}

type Summary struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	NodesCount    int       `json:"nodes_count"`
	FinalArtifact string    `json:"final_artifact,omitempty"`
	TriggerType   string    `json:"trigger_type,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
