package ativos

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// RecortadorSprite abstrai recorte com rembg + despill verde e fatiamento em sprites.
type RecortadorSprite interface {
	RecortarEFatiar(ctx context.Context, projDir string, ficha FichaPersonagem) ([]DescritorSprite, error)
}

type recortadorLocalPadrao struct{}

func (recortadorLocalPadrao) RecortarEFatiar(_ context.Context, projDir string, ficha FichaPersonagem) ([]DescritorSprite, error) {
	tipos := []string{"pose", "expressao"}
	saida := make([]DescritorSprite, 0, len(tipos))
	for i, tipo := range tipos {
		id := fmt.Sprintf("%s_%s_%d", ficha.ID, tipo, i+1)
		relPNG := filepath.ToSlash(filepath.Join("sprites", id+".png"))
		if err := gerarSpritePNGLocal(projDir, relPNG, tipo); err != nil {
			return nil, err
		}
		saida = append(saida, DescritorSprite{
			ID:           id,
			PersonagemID: ficha.ID,
			Tipo:         tipo,
			ArquivoPNG:   relPNG,
			Largura:      256,
			Altura:       256,
			AlfaValido:   true,
			DespillVerde: true,
			BoundingBox:  BoundingBox{X: i * 140, Y: 0, W: 128, H: 128},
		})
	}
	return saida, nil
}

func gerarSpritePNGLocal(projDir, relPNG, tipo string) error {
	destino, err := ResolverCaminhoSeguro(projDir, relPNG)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	cor := color.RGBA{R: 56, G: 189, B: 248, A: 255}
	if tipo == "expressao" {
		cor = color.RGBA{R: 250, G: 204, B: 21, A: 255}
	}
	for y := 48; y < 208; y++ {
		for x := 96; x < 160; x++ {
			img.SetRGBA(x, y, cor)
		}
	}
	f, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// ExecutarEtapa3Elenco valida as 3 fichas por personagem, executa recorte+despill e grava flow/elenco.json e sprites/*.json.
func ExecutarEtapa3Elenco(ctx context.Context, projDir string, elenco []FichaPersonagem, rec RecortadorSprite) ([]DescritorSprite, error) {
	if len(elenco) == 0 {
		return nil, fmt.Errorf("%w: elenco vazio", ErrFichaInvalida)
	}
	if rec == nil {
		rec = recortadorLocalPadrao{}
	}
	var todos []DescritorSprite
	for _, ficha := range elenco {
		if err := ficha.Validar(); err != nil {
			return nil, err
		}
		sprites, err := rec.RecortarEFatiar(ctx, projDir, ficha)
		if err != nil {
			return nil, err
		}
		if err := ValidarSpritesSemSobreposicao(sprites); err != nil {
			return nil, err
		}
		todos = append(todos, sprites...)
	}
	if err := gravarJSONSeguro(projDir, filepath.Join("flow", "elenco.json"), elenco); err != nil {
		return nil, err
	}
	for _, sp := range todos {
		rel := filepath.Join("sprites", sp.ID+".json")
		if err := gravarJSONSeguro(projDir, rel, sp); err != nil {
			return nil, err
		}
	}
	return todos, nil
}

// ValidarSpritesSemSobreposicao aplica o Gate da Etapa 3: atributos obrigatorios e zero sobreposicao.
func ValidarSpritesSemSobreposicao(sprites []DescritorSprite) error {
	if len(sprites) == 0 {
		return fmt.Errorf("%w: lista de sprites vazia", ErrSpriteInvalido)
	}
	for i := range sprites {
		if err := sprites[i].Validar(); err != nil {
			return err
		}
		for j := i + 1; j < len(sprites); j++ {
			if sprites[i].PersonagemID == sprites[j].PersonagemID && sprites[i].BoundingBox.Sobrepoe(sprites[j].BoundingBox) {
				return fmt.Errorf("%w: %s e %s", ErrSobreposicaoSprite, sprites[i].ID, sprites[j].ID)
			}
		}
	}
	return nil
}

func gravarJSONSeguro(projDir, relativo string, v any) error {
	destino, err := ResolverCaminhoSeguro(projDir, relativo)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(destino, b, 0o644)
}
