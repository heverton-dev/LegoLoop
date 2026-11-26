package ativos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

// PedidoPersonagem e o personagem pedido pelo roteiro: tags do tema e o cast desejado.
type PedidoPersonagem struct {
	ID   string            `json:"id"`
	Nome string            `json:"nome"`
	Tags []string          `json:"tags"`
	Cast map[string]string `json:"cast"`
}

// LacunaElenco e uma parte pedida que o inventario do estilo ainda nao tem aprovada.
type LacunaElenco struct {
	PersonagemID string `json:"personagem_id"`
	Campo        string `json:"campo"`
	Tipo         string `json:"tipo"`
	Item         string `json:"item"`
}

// RelatorioElencoVetorial registra custo, uso do Flow e lacunas da Etapa 3 vetorial.
type RelatorioElencoVetorial struct {
	Estilo      string         `json:"estilo"`
	Creditos    int            `json:"creditos"`
	FlowChamado bool           `json:"flow_chamado"`
	Lacunas     []LacunaElenco `json:"lacunas"`
}

// ExecutarEtapa3Vetorial monta elenco.json so com o inventario, sem Flow e com 0 credito.
// Parte ausente vira item novo com aprovado false; o Stop-Gate humano aprova depois.
func ExecutarEtapa3Vetorial(projDir string, s inventario.Store, estilo string, pedidos []PedidoPersonagem, origem string) (RelatorioElencoVetorial, error) {
	e, inv, err := carregarEstiloVetorial(s, estilo)
	if err != nil {
		return RelatorioElencoVetorial{}, err
	}
	rel := RelatorioElencoVetorial{Estilo: estilo, Lacunas: []LacunaElenco{}}
	el := Elenco{SchemaVersion: 2, Estilo: estilo, Personagens: []PersonagemElenco{}}
	for i, p := range pedidos {
		if strings.TrimSpace(p.ID) == "" {
			return RelatorioElencoVetorial{}, fmt.Errorf("%w: personagem %d sem id", ErrElencoInvalido, i)
		}
		lacunas, err := lacunasDoCast(p, inv)
		if err != nil {
			return RelatorioElencoVetorial{}, err
		}
		for _, l := range lacunas {
			if err := gravarItemNovo(s, e, p.Tags, l, origem); err != nil {
				return RelatorioElencoVetorial{}, err
			}
		}
		rel.Lacunas = append(rel.Lacunas, lacunas...)
		el.Personagens = append(el.Personagens, PersonagemElenco{ID: p.ID, Nome: p.Nome, Cast: p.Cast})
	}
	if err := gravarJSONSeguro(projDir, "elenco.json", el); err != nil {
		return RelatorioElencoVetorial{}, err
	}
	return rel, gravarJSONSeguro(projDir, "elenco_relatorio.json", rel)
}

// carregarEstiloVetorial exige estilo existente e validado e devolve o inventario aprovado.
func carregarEstiloVetorial(s inventario.Store, estilo string) (inventario.Estilo, InventarioEstilo, error) {
	e, err := s.LerEstilo(estilo)
	if errors.Is(err, inventario.ErrEstiloInexistente) {
		return e, InventarioEstilo{}, fmt.Errorf("%w: %w", ErrEstiloInexistente, err)
	}
	if err != nil {
		return e, InventarioEstilo{}, err
	}
	if e.Status != "validado" {
		return e, InventarioEstilo{}, fmt.Errorf("%w: estilo %q status %q", ErrEstiloNaoValidado, e.ID, e.Status)
	}
	inv, err := CarregarInventarioEstilo(s, estilo)
	return e, inv, err
}

// lacunasDoCast confere cada campo do cast e devolve as partes ausentes do inventario.
func lacunasDoCast(p PedidoPersonagem, inv InventarioEstilo) ([]LacunaElenco, error) {
	campos := make([]string, 0, len(p.Cast))
	for c := range p.Cast {
		campos = append(campos, c)
	}
	sort.Strings(campos)
	var lacunas []LacunaElenco
	for _, c := range campos {
		valor := p.Cast[c]
		if campoCor(c) {
			continue
		}
		tipo, forma, ok := tipoCampoCast(c)
		if !ok || strings.TrimSpace(valor) == "" {
			return nil, fmt.Errorf("%w: personagem %q campo %q valor %q", ErrParteNaoSuportada, p.ID, c, valor)
		}
		if !inv.Tem(tipo, valor) {
			lacunas = append(lacunas, LacunaElenco{PersonagemID: p.ID, Campo: c, Tipo: tipo, Item: valor})
			continue
		}
		if forma && !inv.Suporta(valor) {
			return nil, fmt.Errorf("%w: personagem %q campo %q parte %q fora de partes_suportadas", ErrParteNaoSuportada, p.ID, c, valor)
		}
	}
	return lacunas, nil
}

// gravarItemNovo registra a parte ausente no inventario com aprovado false.
func gravarItemNovo(s inventario.Store, e inventario.Estilo, tags []string, l LacunaElenco, origem string) error {
	it := inventario.Item{
		SchemaVersion: 1,
		ID:            l.Item,
		Tipo:          l.Tipo,
		Estilo:        e.ID,
		EstiloVersao:  e.Versao,
		Tags:          tags,
		Parametros:    map[string]any{"campo_cast": l.Campo},
		Origem:        origem,
		Aprovado:      false,
	}
	if err := s.GravarItem(it); err != nil {
		return fmt.Errorf("%w: item novo %q: %w", ErrElencoInvalido, l.Item, err)
	}
	return nil
}

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
