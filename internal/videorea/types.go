package videorea

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidURL sinaliza URL malformada ou plataforma não suportada.
	ErrInvalidURL = errors.New("videorea: url inválida ou plataforma não suportada")

	// ErrEmptyDossier sinaliza ausência de conteúdo na extração forense.
	ErrEmptyDossier = errors.New("videorea: dossiê forense vazio")

	// ErrIncompleteDossier sinaliza dossiê sem nenhum achado forense válido.
	ErrIncompleteDossier = errors.New("dossier sem achados forenses válidos")

	// ErrGateRejected sinaliza rejeição no stop-gate humano.
	ErrGateRejected = errors.New("videorea: operação rejeitada no stop-gate")
)

// Platform enumera as plataformas de vídeo suportadas.
type Platform string

const (
	PlatformYouTube   Platform = "youtube"
	PlatformInstagram Platform = "instagram"
	PlatformTikTok    Platform = "tiktok"
	PlatformLocal     Platform = "local"
)

// VideoInput representa a entrada para o crawler de engenharia reversa.
type VideoInput struct {
	URL         string            `json:"url"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	RawVideo    string            `json:"raw_video,omitempty"`
	Duration    float64           `json:"duration,omitempty"`
	AudioPath   string            `json:"audio_path,omitempty"`
	Platform    Platform          `json:"platform,omitempty"`
	AuthCookie  string            `json:"auth_cookie,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// ForensicsDossier consolida as evidências forenses extraídas do vídeo.
type ForensicsDossier struct {
	VideoURL        string   `json:"video_url"`
	Platform        Platform `json:"platform"`
	Title           string   `json:"title"`
	TranscribedText string   `json:"transcribed_text"`
	ExplicitPrompts []string `json:"explicit_prompts"`
	HiddenPrompts   []string `json:"hidden_prompts"`
	ToolsIdentified []string `json:"tools_identified"`
	HiddenTools     []string `json:"hidden_tools"`
	OCRFindings     []string `json:"ocr_findings"`
	SuggestedSlugs  []string `json:"suggested_slugs"`
}

// ValidateStrict rejeita o dossiê que não traz nenhum achado forense.
func (d *ForensicsDossier) ValidateStrict() error {
	if len(d.ToolsIdentified) == 0 && len(d.HiddenTools) == 0 && len(d.ExplicitPrompts) == 0 &&
		len(d.HiddenPrompts) == 0 && len(d.OCRFindings) == 0 {
		return ErrIncompleteDossier
	}
	return nil
}

// DetectionType categoriza o método de descoberta de uma ferramenta.
type DetectionType string

const (
	DetectionExplicitMention  DetectionType = "explicit_mention"
	DetectionExplicitLink     DetectionType = "explicit_link"
	DetectionImplicitTerminal DetectionType = "implicit_terminal"
	DetectionImplicitCode     DetectionType = "implicit_code"
)

// OriginalTool representa uma ferramenta ou script na stack original do autor.
type OriginalTool struct {
	Name                 string        `json:"name"`
	Category             string        `json:"category"`
	IsFree               bool          `json:"is_free"`
	LicenseType          string        `json:"license_type"`
	CostTier             string        `json:"cost_tier"`
	Detection            DetectionType `json:"detection"`
	SourceEvidence       string        `json:"source_evidence"`
	Timestamp            string        `json:"timestamp,omitempty"`
	RepositoryURL        string        `json:"repository_url,omitempty"`
	LocalArtifactPath    string        `json:"local_artifact_path,omitempty"`
	ReconstructedSnippet string        `json:"reconstructed_snippet,omitempty"`
}

// SovereignReplacement detalha a peça soberana do cofre que substitui uma comercial.
type SovereignReplacement struct {
	Name      string `json:"name"`
	Tier      string `json:"tier"`
	Cost      string `json:"cost"`
	License   string `json:"license"`
	Rationale string `json:"rationale"`
}

// ToolDiagnosis consolida a situação e o encaixe real da ferramenta no sistema do operador.
type ToolDiagnosis struct {
	InstalledOnHost bool          `json:"installed_on_host"`
	HostLocation    string        `json:"host_location,omitempty"`
	FitStatus       LegoFitStatus `json:"fit_status"`
	FitSolution     string        `json:"fit_solution"`
	GoldStandard    string        `json:"gold_standard"`
	InstallCommand  string        `json:"install_command,omitempty"`
}

// SubstitutedTool vincula a ferramenta proprietária original ao seu substituto e diagnóstico.
type SubstitutedTool struct {
	Original    OriginalTool         `json:"original"`
	Replacement SovereignReplacement `json:"replacement"`
	Diagnosis   ToolDiagnosis        `json:"diagnosis"`
}

// AdoptedTool vincula a ferramenta gratuita adotada ao seu diagnóstico no host do operador.
type AdoptedTool struct {
	Adopted   OriginalTool  `json:"adopted"`
	Diagnosis ToolDiagnosis `json:"diagnosis"`
}

// OriginalMatrix consolida a stack tecnológica original do criador antes da curadoria.
type OriginalMatrix struct {
	VideoTitle       string         `json:"video_title"`
	TotalMonthlyCost string         `json:"total_monthly_cost"`
	Tools            []OriginalTool `json:"tools"`
	OriginalSteps    []string       `json:"original_steps"`
}

// LegoFitStatus qualifica o encaixe de uma ferramenta local contra as exigências do pipeline.
type LegoFitStatus string

const (
	FitOuro           LegoFitStatus = "FIT_OURO"
	FitFrouxoDefasado LegoFitStatus = "FIT_FROUXO_DEFASADO"
	MissingAquisicao  LegoFitStatus = "MISSING_AQUISICAO"
)

// CurationReport consolida a curadoria contra o cofre e economia de soberania.
type CurationReport struct {
	SuggestedName    string            `json:"suggested_name"`
	VideoTitle       string            `json:"video_title"`
	OriginalCost     string            `json:"original_cost"`
	SovereignCost    string            `json:"sovereign_cost"`
	Replacements     map[string]string `json:"replacements"`
	SubstitutedTools []SubstitutedTool `json:"substituted_tools,omitempty"`
	AdoptedTools     []AdoptedTool     `json:"adopted_tools,omitempty"`
	OriginalMatrix   *OriginalMatrix   `json:"original_matrix,omitempty"`
	Approved         bool              `json:"approved"`
}

// GovernanceDraft representa os arquivos gerados pela fábrica.
type GovernanceDraft struct {
	Slug          string `json:"slug"`
	IntentPath    string `json:"intent_path"`
	SpecPath      string `json:"spec_path"`
	PlanPath      string `json:"plan_path"`
	RoutePath     string `json:"route_path"`
	BookDir       string `json:"book_dir"`
	GoldenFixture string `json:"golden_fixture"`
}

// IntentDraft representa o rascunho de intenção materializado.
type IntentDraft struct {
	Slug       string `json:"slug"`
	IntentPath string `json:"intent_path"`
	Confidence int    `json:"confidence"`
	Status     string `json:"status"`
}

// SpecDraft representa o rascunho de especificação formal materializado.
type SpecDraft struct {
	Slug     string `json:"slug"`
	SpecPath string `json:"spec_path"`
	Status   string `json:"status"`
}

// ExecutionPackageDraft reúne os 4 artefatos técnicos consolidados da sub-etapa 3.4.
type ExecutionPackageDraft struct {
	Slug        string `json:"slug"`
	PlanPath    string `json:"plan_path"`
	RoutePath   string `json:"route_path"`
	DAGPath     string `json:"dag_path"`
	FixturePath string `json:"fixture_path"`
}

// PipelineStage enumera as fases sequenciais do pipeline de engenharia reversa.
type PipelineStage string

const (
	StageForensics PipelineStage = "forensics"
	StageCurated   PipelineStage = "curated"
	StageInterview PipelineStage = "interview"
	StageIntent    PipelineStage = "intent"
	StageSpec      PipelineStage = "spec"
	StageRoute     PipelineStage = "route"
	StageDAG       PipelineStage = "dag_fixture"
	StageSkill     PipelineStage = "skill"
	StageBook      PipelineStage = "book"
	StageClosed    PipelineStage = "closed"
)

// StopGateType identifica cada portão de auditoria humana obrigatório.
type StopGateType string

const (
	GateCuradoria StopGateType = "SG_CURADORIA"
	GateIntent    StopGateType = "SG_INTENT"
	GateSpec      StopGateType = "SG_SPEC"
	GateRoute     StopGateType = "SG_ROUTE"
	GateDAG       StopGateType = "SG_DAG"
)

// StopGateDecision encapsula o veredito humano e observações de um portão.
type StopGateDecision struct {
	Gate     StopGateType `json:"gate"`
	Approved bool         `json:"approved"`
	Notes    string       `json:"notes,omitempty"`
}

// Validate valida a integridade do input de vídeo.
func (v VideoInput) Validate() error {
	if strings.TrimSpace(v.URL) == "" && strings.TrimSpace(v.RawVideo) == "" {
		return ErrInvalidURL
	}
	return nil
}
