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
	URL        string            `json:"url"`
	RawVideo   string            `json:"raw_video,omitempty"`
	AuthCookie string            `json:"auth_cookie,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
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
	ReconstructedSnippet string        `json:"reconstructed_snippet,omitempty"`
}

// OriginalMatrix consolida a stack tecnológica original do criador antes da curadoria.
type OriginalMatrix struct {
	VideoTitle       string         `json:"video_title"`
	TotalMonthlyCost string         `json:"total_monthly_cost"`
	Tools            []OriginalTool `json:"tools"`
	OriginalSteps    []string       `json:"original_steps"`
}

// CurationReport consolida a curadoria contra o cofre e economia de soberania.
type CurationReport struct {
	SuggestedName string            `json:"suggested_name"`
	OriginalCost  string            `json:"original_cost"`
	SovereignCost string            `json:"sovereign_cost"`
	Replacements  map[string]string `json:"replacements"`
	AdoptedTools  []OriginalTool    `json:"adopted_tools,omitempty"`
	Approved      bool              `json:"approved"`
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

// Validate valida a integridade do input de vídeo.
func (v VideoInput) Validate() error {
	if strings.TrimSpace(v.URL) == "" && strings.TrimSpace(v.RawVideo) == "" {
		return ErrInvalidURL
	}
	return nil
}
