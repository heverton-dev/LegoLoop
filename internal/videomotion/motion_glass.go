package videomotion

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidCardDimensions indica largura ou altura inválida.
	ErrInvalidCardDimensions = errors.New("videomotion: dimensões do cartão de vidro inválidas")
)

// GlassCardSpec especifica os parâmetros visuais e temporais do cartão 3D.
type GlassCardSpec struct {
	WidthPx        int
	HeightPx       int
	BorderRadiusPx int
	BackdropBlurPx int
	BorderGlowHex  string
	ScaleStart     float64
	ScaleEnd       float64
	DurationMs     int64
}

// GlassCardComposition encapsula as propriedades renderizáveis do cartão.
type GlassCardComposition struct {
	WidthPx        int
	HeightPx       int
	BorderRadiusPx int
	BackdropBlurPx int
	BorderGlowHex  string
	ScaleStart     float64
	ScaleEnd       float64
	DurationMs     int64
}

// BuildGlassCard valida as dimensões e constrói a composição do cartão de vidro.
func BuildGlassCard(spec GlassCardSpec) (*GlassCardComposition, error) {
	if spec.WidthPx <= 0 || spec.HeightPx <= 0 {
		return nil, ErrInvalidCardDimensions
	}
	if spec.DurationMs <= 0 {
		spec.DurationMs = 2000
	}
	if spec.ScaleStart == 0 && spec.ScaleEnd == 0 {
		spec.ScaleStart = 0.95
		spec.ScaleEnd = 1.05
	}
	return &GlassCardComposition{
		WidthPx:        spec.WidthPx,
		HeightPx:       spec.HeightPx,
		BorderRadiusPx: spec.BorderRadiusPx,
		BackdropBlurPx: spec.BackdropBlurPx,
		BorderGlowHex:  spec.BorderGlowHex,
		ScaleStart:     spec.ScaleStart,
		ScaleEnd:       spec.ScaleEnd,
		DurationMs:     spec.DurationMs,
	}, nil
}

// GenerateCSS produz o bloco de estilos CSS puro para Remotion/Web-Canvas.
func (c *GlassCardComposition) GenerateCSS() string {
	return fmt.Sprintf(
		"width: %dpx; height: %dpx; border-radius: %dpx; backdrop-filter: blur(%dpx); border: 1px solid %s; box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);",
		c.WidthPx, c.HeightPx, c.BorderRadiusPx, c.BackdropBlurPx, c.BorderGlowHex,
	)
}

// CalculateScaleAt calcula a escala do push-in suave em função do tempo decorrido.
func (c *GlassCardComposition) CalculateScaleAt(elapsedMs int64) float64 {
	if elapsedMs <= 0 {
		return c.ScaleStart
	}
	if elapsedMs >= c.DurationMs {
		return c.ScaleEnd
	}
	progress := float64(elapsedMs) / float64(c.DurationMs)
	// Interpolação suave ease-in-out clássica
	smoothProgress := progress * progress * (3.0 - 2.0*progress)
	return c.ScaleStart + (c.ScaleEnd-c.ScaleStart)*smoothProgress
}
