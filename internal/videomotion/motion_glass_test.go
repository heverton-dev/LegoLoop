package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestGlassCardComposition(t *testing.T) {
	spec := videomotion.GlassCardSpec{
		WidthPx:        800,
		HeightPx:       500,
		BorderRadiusPx: 24,
		BackdropBlurPx: 16,
		BorderGlowHex:  "#FFFFFF33",
		ScaleStart:     0.95,
		ScaleEnd:       1.05,
		DurationMs:     2400,
	}

	comp, err := videomotion.BuildGlassCard(spec)
	if err != nil {
		t.Fatalf("BuildGlassCard falhou: %v", err)
	}

	if comp.WidthPx != 800 || comp.HeightPx != 500 {
		t.Errorf("dimensões esperadas 800x500, obteve %dx%d", comp.WidthPx, comp.HeightPx)
	}

	// Validação de CSS e transformações 3D
	css := comp.GenerateCSS()
	if !containsStr(css, "backdrop-filter: blur(16px)") {
		t.Errorf("esperava backdrop-filter com blur de 16px no CSS, obteve:\n%s", css)
	}
	if !containsStr(css, "border-radius: 24px") {
		t.Errorf("esperava border-radius de 24px no CSS")
	}

	// Curva de push-in nos frames
	scaleAtMid := comp.CalculateScaleAt(1200) // meio da animação
	if scaleAtMid <= 0.95 || scaleAtMid >= 1.05 {
		t.Errorf("scale no meio da animação deve estar interpolado entre 0.95 e 1.05, obteve %f", scaleAtMid)
	}
}

func TestGlassCardComposition_InvalidDimensions(t *testing.T) {
	spec := videomotion.GlassCardSpec{
		WidthPx:  0,
		HeightPx: -100,
	}
	_, err := videomotion.BuildGlassCard(spec)
	if err == nil {
		t.Fatal("esperava erro para dimensões inválidas do cartão")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (hasSubstring(s, sub))))
}

func hasSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
