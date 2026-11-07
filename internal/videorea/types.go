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

// CurationReport consolida a curadoria contra o cofre e economia de soberania.
type CurationReport struct {
	SuggestedName string            `json:"suggested_name"`
	OriginalCost  string            `json:"original_cost"`
	SovereignCost string            `json:"sovereign_cost"`
	Replacements  map[string]string `json:"replacements"`
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
