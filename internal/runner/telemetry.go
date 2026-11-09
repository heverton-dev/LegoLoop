package runner

import (
	"encoding/json"
	"fmt"
	"os"
)

// FlightRecord consolida a telemetria unificada no padrão OpenTelemetry GenAI (Motor V).
type FlightRecord struct {
	Timestamp          string     `json:"timestamp"`
	PipelineName       string     `json:"pipeline_name"`
	TotalTokens        TokenUsage `json:"total_tokens"`
	NodesPassed        int        `json:"nodes_passed"`
	NodesFailed        int        `json:"nodes_failed"`
	FirstPassYield     float64    `json:"first_pass_yield"`
	ContextUtilization float64    `json:"context_utilization,omitempty"`
}

// SaveFlightRecord grava o relatório de voo em formato JSON indentado no caminho informado.
func (f FlightRecord) SaveFlightRecord(path string) error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("flight_record: erro ao serializar json: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("flight_record: erro ao gravar arquivo %q: %w", path, err)
	}

	return nil
}
