package videomotion

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"strings"
)

var (
	// ErrEmptyPrompt indica prompt vazio para a bridge.
	ErrEmptyPrompt = errors.New("videomotion: prompt de ativo sintético não pode ser vazio")
)

// AssetResult descreve o ativo resolvido pela bridge.
type AssetResult struct {
	FilePath     string
	FromCache    bool
	UsedFallback bool
}

// StudioAssetBridge orquestra a geração de B-roll com cache e fallback local.
type StudioAssetBridge struct {
	apiKey string
	cache  *AssetCache
}

// NewStudioAssetBridge inicializa o conector.
func NewStudioAssetBridge(apiKey string, cache *AssetCache) *StudioAssetBridge {
	return &StudioAssetBridge{apiKey: apiKey, cache: cache}
}

// ResolveAsset busca no cache ou gera o ativo via fallback procedural local.
func (b *StudioAssetBridge) ResolveAsset(prompt, seed string, aspect AspectRatio) (*AssetResult, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, ErrEmptyPrompt
	}
	key := b.cache.ComputeKey(prompt, seed, aspect)
	filePath := b.cache.GetFilePath(key, ".png")
	if _, found := b.cache.Get(key); found {
		return &AssetResult{FilePath: filePath, FromCache: true, UsedFallback: false}, nil
	}
	imgData, err := generateProceduralHalftone(aspect)
	if err != nil {
		return nil, err
	}
	if err := b.cache.Put(key, imgData, ".png"); err != nil {
		return nil, err
	}
	return &AssetResult{FilePath: filePath, FromCache: false, UsedFallback: true}, nil
}

func generateProceduralHalftone(aspect AspectRatio) ([]byte, error) {
	w, h := 1080, 1920
	if aspect == AspectHorizontal {
		w, h = 1920, 1080
	} else if aspect == AspectSquare {
		w, h = 1080, 1080
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	fillEditorialBackground(img, w, h)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func fillEditorialBackground(img *image.RGBA, w, h int) {
	darkBg := color.RGBA{R: 15, G: 16, B: 18, A: 255}
	dotColor := color.RGBA{R: 229, G: 9, B: 20, A: 40} // Vermelho editorial suave
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x%24 == 0) && (y%24 == 0) {
				img.Set(x, y, dotColor)
			} else {
				img.Set(x, y, darkBg)
			}
		}
	}
}
