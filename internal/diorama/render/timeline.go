package render

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

// Quadro base da composicao: a timeline e sempre 16:9; os cortes 9:16 reusam as caixas.
const (
	LarguraBase = 1920
	AlturaBase  = 1080
)

var ErrEntradaTimeline = errors.New("diorama-render: entrada da timeline invalida")

// Posicao horizontal (fracao da largura) de cada personagem; o primeiro, o narrador, fica a direita.
func slotsPersonagem() []float64 { return []float64{0.65, 0.35, 0.5, 0.8, 0.2} }

func entradaInvalida(formato string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrEntradaTimeline, fmt.Sprintf(formato, args...))
}

func arredondar(v float64) int { return int(math.Round(v)) }

type gruposElementos struct {
	objetos, personagens, efeitos []ativos.ElementoMovel
	cameraAcao                    string
}

func agruparElementos(cena ativos.Cena) (gruposElementos, error) {
	var g gruposElementos
	for _, el := range cena.ElementosMoveis {
		switch el.Tipo {
		case "objeto":
			g.objetos = append(g.objetos, el)
		case "personagem":
			g.personagens = append(g.personagens, el)
		case "efeito":
			g.efeitos = append(g.efeitos, el)
		case "camera":
			if g.cameraAcao == "" {
				g.cameraAcao = el.Acao
			}
		default:
			return g, entradaInvalida("linha %d com elemento %q de tipo %q desconhecido", cena.IndiceLinha, el.ID, el.Tipo)
		}
	}
	return g, nil
}

// cameraDaAcao traduz a acao do storyboard em um movimento de camera conhecido.
func cameraDaAcao(acao string) string {
	a := strings.ToLower(acao)
	switch {
	case strings.Contains(a, "zoom"):
		return "zoom"
	case strings.Contains(a, "shake") || strings.Contains(a, "tremor"):
		return "shake"
	case strings.Contains(a, "travelling") || strings.Contains(a, "pan"):
		return "pan"
	}
	return "estatica"
}

func quadroCompleto() ativos.BoundingBox {
	return ativos.BoundingBox{X: 0, Y: 0, W: LarguraBase, H: AlturaBase}
}

func caixaLegenda() ativos.BoundingBox {
	return ativos.BoundingBox{X: arredondar(0.05 * LarguraBase), Y: arredondar(0.78 * AlturaBase), W: arredondar(0.9 * LarguraBase), H: arredondar(0.18 * AlturaBase)}
}

// caixaObjeto distribui ate 3 objetos lado a lado na metade esquerda do palco.
func caixaObjeto(k int) ativos.BoundingBox {
	slot := k % 3
	return ativos.BoundingBox{
		X: arredondar(0.05*LarguraBase) + slot*arredondar(0.28*LarguraBase), Y: arredondar(0.30 * AlturaBase),
		W: arredondar(0.25 * LarguraBase), H: arredondar(0.45 * AlturaBase),
	}
}

// caixaSprite usa 70% da altura, mantem a proporcao e apoia os pes em 95% do quadro.
func caixaSprite(d ativos.DescritorSprite, slot int) (ativos.BoundingBox, error) {
	if d.Largura <= 0 || d.Altura <= 0 {
		return ativos.BoundingBox{}, entradaInvalida("sprite %q sem dimensoes", d.ID)
	}
	h := arredondar(0.70 * AlturaBase)
	w := arredondar(float64(h) * float64(d.Largura) / float64(d.Altura))
	if w > LarguraBase {
		return ativos.BoundingBox{}, entradaInvalida("sprite %q mais largo que o quadro", d.ID)
	}
	slots := slotsPersonagem()
	x := arredondar(slots[slot%len(slots)]*LarguraBase) - w/2
	x = max(0, min(x, LarguraBase-w))
	return ativos.BoundingBox{X: x, Y: arredondar(0.95*AlturaBase) - h, W: w, H: h}, nil
}

// expressaoSeAplica troca a pose pela expressao do mesmo personagem nas falas incertas.
func expressaoSeAplica(sprites map[string]ativos.DescritorSprite, base ativos.DescritorSprite, natureza string) ativos.DescritorSprite {
	if base.Tipo == "expressao" || (natureza != "alegacao_terceiro" && natureza != "evidencia_preliminar") {
		return base
	}
	var ids []string
	for id, s := range sprites {
		if s.PersonagemID == base.PersonagemID && s.Tipo == "expressao" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return base
	}
	sort.Strings(ids)
	return sprites[ids[0]]
}

func itemSprite(c ContratosRender, el ativos.ElementoMovel, slot int, natureza string) (CamadaItem, error) {
	base, ok := c.Sprites[el.SpriteID]
	if el.SpriteID == "" || !ok {
		return CamadaItem{}, entradaInvalida("personagem %q com sprite_id %q ausente", el.ID, el.SpriteID)
	}
	d := expressaoSeAplica(c.Sprites, base, natureza)
	caixa, err := caixaSprite(d, slot)
	if err != nil {
		return CamadaItem{}, err
	}
	return CamadaItem{ID: el.ID, Camada: "sprite", SpriteID: d.ID, Arquivo: d.ArquivoPNG, Caixa: caixa}, nil
}

func comporBatida(c ContratosRender, tm Timings, arquivos map[string]string, cena ativos.Cena) (Batida, error) {
	idx := cena.IndiceLinha
	if idx < 0 || idx >= len(tm.Linhas) || idx >= len(c.Script.Lines) {
		return Batida{}, entradaInvalida("linha %d sem timing ou sem roteiro", idx)
	}
	arquivo, ok := arquivos[cena.CenarioID]
	if !ok {
		return Batida{}, entradaInvalida("cenario %q sem arquivo", cena.CenarioID)
	}
	g, err := agruparElementos(cena)
	if err != nil {
		return Batida{}, err
	}
	camadas := []CamadaItem{{ID: "cenario:" + cena.CenarioID, Camada: "cenario", Arquivo: arquivo, Caixa: quadroCompleto()}}
	for k, el := range g.objetos {
		camadas = append(camadas, CamadaItem{ID: el.ID, Camada: "aderecos", Caixa: caixaObjeto(k)})
	}
	for k, el := range g.personagens {
		item, err := itemSprite(c, el, k, c.Script.Lines[idx].Natureza)
		if err != nil {
			return Batida{}, err
		}
		camadas = append(camadas, item)
	}
	for _, el := range g.efeitos {
		camadas = append(camadas, CamadaItem{ID: el.ID, Camada: "efeito", Caixa: quadroCompleto()})
	}
	camadas = append(camadas, CamadaItem{ID: fmt.Sprintf("legenda:linha_%02d", idx), Camada: "legenda", Caixa: caixaLegenda()})
	for i := range camadas {
		camadas[i].Profundidade = i
	}
	l := tm.Linhas[idx]
	return Batida{IndiceLinha: idx, InicioS: l.InicioS, FimS: l.FimS, CenarioID: cena.CenarioID, Camera: cameraDaAcao(g.cameraAcao), Camadas: camadas}, nil
}

// Compor monta uma batida por cena, em ordem de linha, a partir de entradas ja carregadas.
func Compor(c ContratosRender, tm Timings, arquivosCenario map[string]string) (Timeline, error) {
	cenas := append([]ativos.Cena(nil), c.Storyboard.Cenas...)
	sort.SliceStable(cenas, func(i, j int) bool { return cenas[i].IndiceLinha < cenas[j].IndiceLinha })
	tl := Timeline{SchemaVersion: 1}
	for _, cena := range cenas {
		b, err := comporBatida(c, tm, arquivosCenario, cena)
		if err != nil {
			return Timeline{}, err
		}
		tl.Batidas = append(tl.Batidas, b)
	}
	if n := len(tm.Linhas); n > 0 {
		tl.DuracaoS = tm.Linhas[n-1].FimS
	}
	return tl, nil
}

func lerTimings(projDir string, nLinhas int) (Timings, error) {
	var tm Timings
	if err := lerJSON(projDir, "data/timings.json", &tm, true); err != nil {
		return tm, entradaInvalida("data/timings.json: %v", err)
	}
	if err := tm.Validar(); err != nil {
		return tm, fmt.Errorf("%w: %w", ErrEntradaTimeline, err)
	}
	if len(tm.Linhas) != nLinhas {
		return tm, entradaInvalida("timings tem %d linhas e script.json tem %d", len(tm.Linhas), nLinhas)
	}
	return tm, nil
}

// arquivosDosCenarios resolve cada cenario usado ao arquivo de flow/ e exige que ele exista.
func arquivosDosCenarios(projDir string, sb ativos.Storyboard) (map[string]string, error) {
	var lista []ativos.Cenario
	if err := lerJSON(projDir, "flow/cenarios.json", &lista, false); err != nil {
		return nil, entradaInvalida("flow/cenarios.json: %v", err)
	}
	porID := map[string]string{}
	for _, cen := range lista {
		porID[cen.ID] = cen.Arquivo
	}
	saida := map[string]string{}
	for _, cena := range sb.Cenas {
		rel, ok := porID[cena.CenarioID]
		if !ok || strings.TrimSpace(rel) == "" {
			return nil, entradaInvalida("cenario %q ausente em flow/cenarios.json", cena.CenarioID)
		}
		caminho, err := resolverSeguro(projDir, rel)
		if err != nil {
			return nil, entradaInvalida("cenario %q: %v", cena.CenarioID, err)
		}
		if _, err := os.Stat(caminho); err != nil {
			return nil, entradaInvalida("cenario %q sem arquivo %s", cena.CenarioID, rel)
		}
		saida[cena.CenarioID] = rel
	}
	return saida, nil
}

// construirTimeline le as entradas de forma estrita, compoe, segmenta, aplica o karaoke e confere os gates.
func construirTimeline(projDir string, c ContratosRender) (Timeline, string, error) {
	tm, err := lerTimings(projDir, len(c.Script.Lines))
	if err != nil {
		return Timeline{}, "", err
	}
	arquivos, err := arquivosDosCenarios(projDir, c.Storyboard)
	if err != nil {
		return Timeline{}, "", err
	}
	tl, err := Compor(c, tm, arquivos)
	if err != nil {
		return Timeline{}, "", err
	}
	tmL := AlinharLegendas(c.Script, tm)
	tl = AplicarLegendas(SegmentarBatidas(tl, tmL), tmL)
	if err := tl.Validar(); err != nil {
		return Timeline{}, "", err
	}
	if _, err := ChecarRitmo(tl); err != nil {
		return Timeline{}, "", err
	}
	srt, err := GerarSRT(tl)
	return tl, srt, err
}

func gravarSaidasTimeline(projDir string, tl Timeline, srt string) error {
	b, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return gravarEmLote(projDir, []saidaArquivo{
		{"renders/legenda.srt", []byte(srt)},
		{"timeline.json", append(b, '\n')},
	})
}

// ExecutarTimeline so grava timeline.json e renders/legenda.srt se a composicao inteira for
// valida e passar no gate de ritmo.
func ExecutarTimeline(ctx context.Context, projDir string) (Timeline, error) {
	if err := ctx.Err(); err != nil {
		return Timeline{}, err
	}
	c, err := CarregarContratos(projDir)
	if err != nil {
		return Timeline{}, fmt.Errorf("%w: %w", ErrEntradaTimeline, err)
	}
	tl, srt, err := construirTimeline(projDir, c)
	if err != nil {
		return Timeline{}, err
	}
	return tl, gravarSaidasTimeline(projDir, tl, srt)
}
