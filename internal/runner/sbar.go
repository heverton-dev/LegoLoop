package runner

import "fmt"

// SBAR representa a estrutura formal do protocolo militar/industrial SBAR.
type SBAR struct {
	Situation      string `json:"situation"`
	Background     string `json:"background"`
	Assessment     string `json:"assessment"`
	Recommendation string `json:"recommendation"`
}

// NewSBAR cria uma nova instância validada de SBAR.
func NewSBAR(situation, background, assessment, recommendation string) SBAR {
	return SBAR{
		Situation:      situation,
		Background:     background,
		Assessment:     assessment,
		Recommendation: recommendation,
	}
}

// Format serializa o SBAR no formato canônico textual de escalonamento.
func (s SBAR) Format() string {
	return fmt.Sprintf(
		"[S] SITUATION: %s\n[B] BACKGROUND: %s\n[A] ASSESSMENT: %s\n[R] RECOMMENDATION: %s",
		s.Situation,
		s.Background,
		s.Assessment,
		s.Recommendation,
	)
}

// FormatNodeError sintetiza uma falha de execução de nó em formato SBAR.
func FormatNodeError(nodeID string, exitCode int, errStr string) SBAR {
	sit := fmt.Sprintf("Falha mecânica de execução no nó %s", nodeID)
	bg := fmt.Sprintf("Nó executado com exit code %d sob o kernel Orchestrix", exitCode)
	ass := fmt.Sprintf("Erro capturado: %s", errStr)
	rec := "Puxar Andon Cord e aguardar auditoria de engenharia (Zero Mutação Cega)"
	return NewSBAR(sit, bg, ass, rec)
}
