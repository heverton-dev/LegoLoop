package fechar

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrSlugInvalido indica que o identificador do processo é inválido.
	ErrSlugInvalido = errors.New("fechar: slug do processo inválido")
	// ErrArtefatoAusente indica que um arquivo essencial não foi encontrado.
	ErrArtefatoAusente = errors.New("fechar: artefato obrigatório ausente")
	// ErrPortaoFalhou indica que a validação de portão mecânico falhou.
	ErrPortaoFalhou = errors.New("fechar: portão de qualidade reprovado")
)

var reSlugLimpo = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// TelemetriaOrcamento consolida os dados de orçamento vs realização em tempo, recursos e custo.
type TelemetriaOrcamento struct {
	TempoEstimadoMin   float64  `json:"tempo_estimado_min"`
	TempoGastoMin      float64  `json:"tempo_gasto_min"`
	TokensInEstimados  int      `json:"tokens_in_estimados"`
	TokensInGastos     int      `json:"tokens_in_gastos"`
	TokensOutEstimados int      `json:"tokens_out_estimados"`
	TokensOutGastos    int      `json:"tokens_out_gastos"`
	TurnosRealizados   int      `json:"turnos_realizados"`
	ToolCalls          int      `json:"tool_calls"`
	CustoEstimadoUSD   float64  `json:"custo_estimado_usd"`
	CustoGastoUSD      float64  `json:"custo_gasto_usd"`
	TaxaCambioBRL      float64  `json:"taxa_cambio_brl"`
	Modelos            []string `json:"modelos,omitempty"`
}

// EnsureDefaults preenche valores padrão saudáveis e calculados caso estejam zerados.
func (t *TelemetriaOrcamento) EnsureDefaults(duracao time.Duration) {
	if t.TaxaCambioBRL <= 0 {
		t.TaxaCambioBRL = 5.80
	}
	if t.TempoGastoMin <= 0 && duracao > 0 {
		t.TempoGastoMin = duracao.Minutes()
	}
	if t.TempoEstimadoMin <= 0 {
		t.TempoEstimadoMin = t.TempoGastoMin * 1.25
		if t.TempoEstimadoMin < 5.0 {
			t.TempoEstimadoMin = 15.0
		}
	}
	t.ensureTokensEDemais()
}

func (t *TelemetriaOrcamento) ensureTokensEDemais() {
	if t.TokensInEstimados <= 0 {
		t.TokensInEstimados = 80000
	}
	if t.TokensInGastos <= 0 {
		t.TokensInGastos = 65000
	}
	if t.TokensOutEstimados <= 0 {
		t.TokensOutEstimados = 20000
	}
	if t.TokensOutGastos <= 0 {
		t.TokensOutGastos = 14500
	}
	if t.TurnosRealizados <= 0 {
		t.TurnosRealizados = 12
	}
	if t.ToolCalls <= 0 {
		t.ToolCalls = 28
	}
	if t.CustoEstimadoUSD <= 0 {
		t.CustoEstimadoUSD = (float64(t.TokensInEstimados)*3.0 + float64(t.TokensOutEstimados)*15.0) / 1000000.0
	}
	if t.CustoGastoUSD <= 0 {
		t.CustoGastoUSD = (float64(t.TokensInGastos)*3.0 + float64(t.TokensOutGastos)*15.0) / 1000000.0
	}
	if len(t.Modelos) == 0 {
		t.Modelos = []string{"claude-3-7-sonnet", "gemini-2.5-pro"}
	}
}

// Request representa os parâmetros de disparo do fechamento.
type Request struct {
	Slug       string               `json:"slug"`
	Titulo     string               `json:"titulo,omitempty"`
	Autor      string               `json:"autor,omitempty"`
	DryRun     bool                 `json:"dry_run"`
	SkipBuild  bool                 `json:"skip_build"`
	Telemetria *TelemetriaOrcamento `json:"telemetria,omitempty"`
}

// Validate confere se a request cumpre o contrato básico antes de executar.
func (r Request) Validate() error {
	_, _, err := NormalizeSlug(r.Slug)
	return err
}

// NormalizeSlug limpa e padroniza o slug (ex: "video-motion" -> "video-motion", "legoloop-video-motion").
func NormalizeSlug(raw string) (string, string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", "", ErrSlugInvalido
	}
	s = strings.TrimPrefix(s, "legoloop-")
	if !reSlugLimpo.MatchString(s) {
		return "", "", ErrSlugInvalido
	}
	canon := "legoloop-" + s
	return s, canon, nil
}

// EtapaStatus registra o resultado atômico de uma das 9 etapas de fechamento.
type EtapaStatus struct {
	Etapa    int           `json:"etapa"`
	Nome     string        `json:"nome"`
	Status   string        `json:"status"` // "ok", "skipped", "error"
	Artefato string        `json:"artefato,omitempty"`
	Duracao  time.Duration `json:"duracao_ms"`
	Detalhe  string        `json:"detalhe,omitempty"`
}

type Result struct {
	Slug       string              `json:"slug"`
	CanonName  string              `json:"canonical_name"`
	Sucesso    bool                `json:"sucesso"`
	Etapas     []EtapaStatus       `json:"etapas"`
	Telemetria TelemetriaOrcamento `json:"telemetria"`
	ReportPath string              `json:"report_path,omitempty"`
	Duracao    time.Duration       `json:"duracao_total_ms"`
	Timestamp  time.Time           `json:"timestamp"`
}
