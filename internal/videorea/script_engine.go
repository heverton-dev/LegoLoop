package videorea

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/brain"
)

var (
	ErrInvalidCartoonTheme    = errors.New("videorea: tema do cartoon short não pode ser vazio")
	ErrScriptGenerationFailed = errors.New("videorea: falha ao gerar roteiro do cartoon short")
	ErrTTSFailed              = errors.New("videorea: falha na síntese de voz Edge-TTS")
	ErrRenderFailed           = errors.New("videorea: falha na renderização visual do cartoon short")
)

// CartoonShortRequest define os parâmetros para a geração do short.
type CartoonShortRequest struct {
	Theme        string `json:"theme"`
	DurationSec  int    `json:"duration_sec"`
	Voice        string `json:"voice"`
	OutputFolder string `json:"output_folder"`
	GenerateMP4  bool   `json:"generate_mp4"`
}

// Validate verifica a consistência dos parâmetros de entrada.
func (r CartoonShortRequest) Validate() error {
	if strings.TrimSpace(r.Theme) == "" {
		return ErrInvalidCartoonTheme
	}
	if r.DurationSec <= 0 {
		r.DurationSec = 30
	}
	if strings.TrimSpace(r.Voice) == "" {
		r.Voice = "pt-BR-AntonioNeural"
	}
	return nil
}

// CartoonScene representa uma cena individual sincronizada.
type CartoonScene struct {
	Index       int      `json:"index"`
	Timestamp   string   `json:"timestamp"`
	Narration   string   `json:"narration"`
	EmotionTag  string   `json:"emotion_tag"`
	VisualCue   string   `json:"visual_cue"`
	Keywords    []string `json:"keywords"`
	AccentColor string   `json:"accent_color"`
}

// CartoonScript representa a estrutura canônica do roteiro.
type CartoonScript struct {
	Title        string         `json:"title"`
	Theme        string         `json:"theme"`
	Hook         string         `json:"hook"`
	FactOrLegend string         `json:"fact_or_legend"`
	Scenes       []CartoonScene `json:"scenes"`
	TotalWords   int            `json:"total_words"`
}

// BuildCartoonScript orquestra a geração do roteiro com gancho, fatos e tags de emoção.
func BuildCartoonScript(ctx context.Context, theme string, b brain.Brain) (*CartoonScript, error) {
	trimmed := strings.TrimSpace(theme)
	if trimmed == "" {
		return nil, ErrInvalidCartoonTheme
	}
	prompt := buildScriptPrompt(trimmed)
	raw, err := b.Infer(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrScriptGenerationFailed, err)
	}
	return parseCartoonScriptJSON(raw, trimmed)
}

func buildScriptPrompt(theme string) string {
	return fmt.Sprintf(`Você é o roteirista mestre de Cartoon Shorts.
Gere um roteiro viral para o tema "%s" seguindo as regras da skill cartoon-audio-short:
1. Duração total: 20 a 30 segundos (ritmo de 3.7 palavras/segundo, total 70 a 110 palavras).
2. Estrutura: Hook paradoxal sem introdução, Escalada com fatos e Twist irônico final.
3. Classificação: Verdade ou Lenda documentada.
4. Adicione tags curtas de emoção [reflexivo], [animado], [pausa curta] na narração.
Retorne EXCLUSIVAMENTE um objeto JSON válido com os campos:
title, theme, hook, fact_or_legend, scenes (array de index, timestamp, narration, emotion_tag, visual_cue, keywords, accent_color), total_words.`, theme)
}

func parseCartoonScriptJSON(raw, theme string) (*CartoonScript, error) {
	clean := stripJSONMarkdown(raw)
	var s CartoonScript
	if err := json.Unmarshal([]byte(clean), &s); err != nil {
		return nil, fmt.Errorf("%w: erro ao decodificar JSON: %v", ErrScriptGenerationFailed, err)
	}
	if s.Theme == "" {
		s.Theme = theme
	}
	if s.TotalWords == 0 {
		s.TotalWords = countScriptWords(s)
	}
	return &s, nil
}

func stripJSONMarkdown(s string) string {
	trimmed := strings.TrimSpace(s)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	return strings.TrimSpace(trimmed)
}

func countScriptWords(s CartoonScript) int {
	total := 0
	for _, sc := range s.Scenes {
		words := strings.Fields(sc.Narration)
		total += len(words)
	}
	return total
}
