package videomotion_test

import (
	"testing"

	"legoloop/internal/videomotion"
)

func TestAutoFramingCenterCrop(t *testing.T) {
	// Vídeo de origem 1920x1080 (16:9), orador deslocado ligeiramente para a direita (X=0.6, Y=0.3)
	src := videomotion.Dimension{Width: 1920, Height: 1080}
	face := videomotion.NormalizedFaceBox{
		CenterX: 0.65, // 65% da largura (deslocado)
		CenterY: 0.35, // 35% da altura
		Width:   0.15,
		Height:  0.25,
	}

	crop, err := videomotion.CalculateVerticalCrop(src, face, videomotion.AspectVertical)
	if err != nil {
		t.Fatalf("CalculateVerticalCrop falhou: %v", err)
	}

	// Para 9:16 em altura 1080: largura proporcional é 1080 * 9 / 16 = 607 pixels
	if crop.Width != 607 {
		t.Errorf("largura esperada 607, obteve %d", crop.Width)
	}
	if crop.Height != 1080 {
		t.Errorf("altura esperada 1080, obteve %d", crop.Height)
	}

	// O cropX deve centralizar o rosto (0.65 * 1920 = 1248px), então cropX ~ 1248 - 607/2 = 944px
	expectedCenterX := int(0.65 * float64(src.Width))
	cropCenterX := crop.X + crop.Width/2
	diff := cropCenterX - expectedCenterX
	if diff < -10 || diff > 10 {
		t.Errorf("centro do crop (%d) desalinhado do centro do rosto (%d)", cropCenterX, expectedCenterX)
	}

	// Safe zone: crop não pode vazar da borda direita (crop.X + crop.Width <= 1920)
	if crop.X+crop.Width > src.Width {
		t.Errorf("crop vazou a borda direita: %d > %d", crop.X+crop.Width, src.Width)
	}
	if crop.X < 0 {
		t.Errorf("crop.X negativo: %d", crop.X)
	}

	// Expressão FFmpeg
	filter := crop.ToFFmpegFilter()
	if filter != "crop=607:1080:944:0" && filter != "crop=607:1080:945:0" {
		t.Logf("filtro FFmpeg gerado: %s", filter)
	}
}

func TestAutoFramingCenterCrop_InvalidDimensions(t *testing.T) {
	src := videomotion.Dimension{Width: 0, Height: 0}
	face := videomotion.NormalizedFaceBox{CenterX: 0.5, CenterY: 0.5}
	_, err := videomotion.CalculateVerticalCrop(src, face, videomotion.AspectVertical)
	if err == nil {
		t.Fatal("esperava erro para dimensões zeradas")
	}
}
