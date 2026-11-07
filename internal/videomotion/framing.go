package videomotion

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidSourceDimension indica resolução de entrada inválida.
	ErrInvalidSourceDimension = errors.New("videomotion: dimensões de origem inválidas")
)

// Dimension representa largura e altura em pixels.
type Dimension struct {
	Width  int
	Height int
}

// NormalizedFaceBox descreve a posição e proporção do rosto normalizada (0.0 a 1.0).
type NormalizedFaceBox struct {
	CenterX float64
	CenterY float64
	Width   float64
	Height  float64
}

// CropWindow define a janela de recorte final para composição.
type CropWindow struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ToFFmpegFilter formata o recorte no filtro nativo do ffmpeg.
func (c CropWindow) ToFFmpegFilter() string {
	return fmt.Sprintf("crop=%d:%d:%d:%d", c.Width, c.Height, c.X, c.Y)
}

// CalculateVerticalCrop calcula a janela 9:16 centralizando o orador com clamping de borda.
func CalculateVerticalCrop(src Dimension, face NormalizedFaceBox, aspect AspectRatio) (*CropWindow, error) {
	if src.Width <= 0 || src.Height <= 0 {
		return nil, ErrInvalidSourceDimension
	}
	targetWidth, targetHeight := computeTargetAspectDimensions(src, aspect)
	facePixelX := int(face.CenterX * float64(src.Width))
	cropX := facePixelX - (targetWidth / 2)
	cropX = clampInt(cropX, 0, src.Width-targetWidth)
	cropY := 0
	if targetHeight < src.Height {
		cropY = (src.Height - targetHeight) / 2
	}
	return &CropWindow{X: cropX, Y: cropY, Width: targetWidth, Height: targetHeight}, nil
}

func computeTargetAspectDimensions(src Dimension, aspect AspectRatio) (int, int) {
	if aspect == AspectVertical {
		targetW := (src.Height * 9) / 16
		if targetW > src.Width {
			targetW = src.Width
		}
		return targetW, src.Height
	}
	return src.Width, src.Height
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
