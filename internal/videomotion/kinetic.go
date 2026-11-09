package videomotion

import (
	"errors"
)

var (
	// ErrNoWordsForKinetic indica ausência de palavras para a timeline cinética.
	ErrNoWordsForKinetic = errors.New("videomotion: nenhuma palavra fornecida para o motor cinético")
)

// KineticThemeConfig define as fontes e cores da tipografia reativa.
type KineticThemeConfig struct {
	PunchFontFamily   string
	RegularFontFamily string
	AccentColorHex    string
	DefaultColorHex   string
}

// KineticState captura as propriedades visuais de um instante t.
type KineticState struct {
	ActiveWordIndex int
	ActiveWord      WordItem
	FontFamily      string
	ColorHex        string
	FontSizePx      int
	IsPunch         bool
}

// KineticEngine orquestra a animação de texto sincronizada à fala.
type KineticEngine struct {
	words []WordItem
	theme KineticThemeConfig
}

// NewKineticEngine inicializa o motor com temas padrão ou customizados.
func NewKineticEngine(words []WordItem, theme KineticThemeConfig) (*KineticEngine, error) {
	if len(words) == 0 {
		return nil, ErrNoWordsForKinetic
	}
	if theme.PunchFontFamily == "" {
		theme.PunchFontFamily = "Anton-Regular"
	}
	if theme.RegularFontFamily == "" {
		theme.RegularFontFamily = "Didot-Serif"
	}
	if theme.AccentColorHex == "" {
		theme.AccentColorHex = "#E50914"
	}
	if theme.DefaultColorHex == "" {
		theme.DefaultColorHex = "#FFFFFF"
	}
	return &KineticEngine{words: words, theme: theme}, nil
}

// GetKineticStateAt resolve a palavra ativa e o estilo no instante t.
func (k *KineticEngine) GetKineticStateAt(timeMs int64) KineticState {
	idx := k.findActiveWordIndex(timeMs)
	if idx < 0 {
		return k.buildIdleState()
	}
	w := k.words[idx]
	if w.IsPunchWord {
		return KineticState{
			ActiveWordIndex: idx,
			ActiveWord:      w,
			FontFamily:      k.theme.PunchFontFamily,
			ColorHex:        k.theme.AccentColorHex,
			FontSizePx:      96,
			IsPunch:         true,
		}
	}
	return KineticState{
		ActiveWordIndex: idx,
		ActiveWord:      w,
		FontFamily:      k.theme.RegularFontFamily,
		ColorHex:        k.theme.DefaultColorHex,
		FontSizePx:      48,
		IsPunch:         false,
	}
}

func (k *KineticEngine) findActiveWordIndex(timeMs int64) int {
	for i, w := range k.words {
		if timeMs >= w.StartMs && timeMs <= w.EndMs {
			return i
		}
	}
	for i := len(k.words) - 1; i >= 0; i-- {
		if timeMs >= k.words[i].EndMs {
			return i
		}
	}
	return 0
}

func (k *KineticEngine) buildIdleState() KineticState {
	return KineticState{
		ActiveWordIndex: -1,
		FontFamily:      k.theme.RegularFontFamily,
		ColorHex:        k.theme.DefaultColorHex,
		FontSizePx:      48,
		IsPunch:         false,
	}
}

// PerspectiveGridOffsetAt calcula o deslocamento suave do grid em perspectiva 3D.
func (k *KineticEngine) PerspectiveGridOffsetAt(timeMs int64) float64 {
	speedPixelsPerSec := 60.0
	return float64(timeMs) / 1000.0 * speedPixelsPerSec
}
