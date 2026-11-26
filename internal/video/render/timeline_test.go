package render

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

func cenaTeste(i int, cenario string) ativos.Cena {
	return ativos.Cena{IndiceLinha: i, CenarioID: cenario, MetaforaVisual: "m", ElementosMoveis: []ativos.ElementoMovel{
		{ID: "p", Tipo: "personagem", SpriteID: "narrador_pose_1", Acao: "gesto_enfatiza"},
		{ID: "o", Tipo: "objeto", Acao: "destaque_foco"},
		{ID: "c", Tipo: "camera", Acao: "travelling_suave"},
		{ID: "e", Tipo: "efeito", Acao: "particulas_profundidade"},
	}}
}

func contratosTimeline() (ContratosRender, Timings, map[string]string) {
	c := ContratosRender{
		Script:     roteiro.Script{Lines: []roteiro.LinhaRoteiro{{Text: "a", Natureza: "fato_consolidado"}, {Text: "b", Natureza: "alegacao_terceiro"}}},
		Storyboard: ativos.Storyboard{SchemaVersion: 1, Cenas: []ativos.Cena{cenaTeste(0, "cen_01"), cenaTeste(1, "cen_02")}},
		Sprites: map[string]ativos.DescritorSprite{
			"narrador_pose_1":      {ID: "narrador_pose_1", PersonagemID: "narrador", Tipo: "pose", ArquivoPNG: "sprites/narrador_pose_1.png", Largura: 167, Altura: 616},
			"narrador_expressao_2": {ID: "narrador_expressao_2", PersonagemID: "narrador", Tipo: "expressao", ArquivoPNG: "sprites/narrador_expressao_2.png", Largura: 299, Altura: 621},
		},
	}
	return c, timingsValido(), map[string]string{"cen_01": "flow/cen_01.png", "cen_02": "flow/cen_02.png"}
}

func camadasPorTipo(b Batida) map[string][]CamadaItem {
	m := map[string][]CamadaItem{}
	for _, c := range b.Camadas {
		m[c.Camada] = append(m[c.Camada], c)
	}
	return m
}

func TestTimelineCamadasOrdenadasPorProfundidade(t *testing.T) {
	c, tm, arqs := contratosTimeline()
	tl, err := Compor(c, tm, arqs)
	if err != nil {
		t.Fatal(err)
	}
	if err := tl.Validar(); err != nil {
		t.Fatalf("timeline deve validar: %v", err)
	}
	if len(tl.Batidas) != 2 || tl.DuracaoS != tm.Linhas[1].FimS {
		t.Fatalf("batidas ou duracao incorretas: %+v", tl)
	}
	b := tl.Batidas[0]
	if b.InicioS != tm.Linhas[0].InicioS || b.FimS != tm.Linhas[0].FimS || b.CenarioID != "cen_01" {
		t.Fatalf("janela da batida deve ser a da linha: %+v", b)
	}
	ordem := []string{"cenario", "aderecos", "sprite", "efeito", "legenda"}
	por := camadasPorTipo(b)
	for i, nome := range ordem {
		if len(por[nome]) == 0 {
			t.Fatalf("camada %q ausente: %+v", nome, b.Camadas)
		}
		if b.Camadas[i].Camada != nome {
			t.Fatalf("posicao %d deveria ser %q, e %q", i, nome, b.Camadas[i].Camada)
		}
	}
	if por["cenario"][0].Arquivo != "flow/cen_01.png" || por["sprite"][0].Arquivo != "sprites/narrador_pose_1.png" {
		t.Fatalf("arquivos de cenario e sprite devem constar na camada: %+v", b.Camadas)
	}
}

func TestTimelineCaixasDentroDoQuadro(t *testing.T) {
	c, tm, arqs := contratosTimeline()
	tl, _ := Compor(c, tm, arqs)
	for _, b := range tl.Batidas {
		for _, cm := range b.Camadas {
			x := cm.Caixa
			if x.W <= 0 || x.H <= 0 || x.X < 0 || x.Y < 0 || x.X+x.W > LarguraBase || x.Y+x.H > AlturaBase {
				t.Fatalf("camada %s fora do quadro %dx%d: %+v", cm.ID, LarguraBase, AlturaBase, x)
			}
		}
	}
	spr := camadasPorTipo(tl.Batidas[0])["sprite"][0].Caixa
	if spr.H != 756 || spr.W != 205 {
		t.Fatalf("sprite deve ter 70%% da altura mantendo proporcao: %+v", spr)
	}
}

func TestTimelineTrocaDeExpressaoPorNatureza(t *testing.T) {
	c, tm, arqs := contratosTimeline()
	tl, _ := Compor(c, tm, arqs)
	if got := camadasPorTipo(tl.Batidas[0])["sprite"][0].SpriteID; got != "narrador_pose_1" {
		t.Fatalf("fato consolidado mantem a pose do storyboard: %s", got)
	}
	if got := camadasPorTipo(tl.Batidas[1])["sprite"][0].SpriteID; got != "narrador_expressao_2" {
		t.Fatalf("alegacao de terceiro usa a expressao do mesmo personagem: %s", got)
	}
	delete(c.Sprites, "narrador_expressao_2")
	tl, err := Compor(c, tm, arqs)
	if err != nil {
		t.Fatalf("sem sprite de expressao nao pode falhar: %v", err)
	}
	if got := camadasPorTipo(tl.Batidas[1])["sprite"][0].SpriteID; got != "narrador_pose_1" {
		t.Fatalf("sem expressao mantem a pose: %s", got)
	}
}

func TestTimelineCameraPorAcao(t *testing.T) {
	casos := map[string]string{
		"travelling_suave": "pan", "pan_lento": "pan", "zoom_suave": "zoom", "zoom_in": "zoom",
		"shake_leve": "shake", "tremor": "shake", "parada": "estatica", "": "estatica",
	}
	for acao, esperado := range casos {
		if got := cameraDaAcao(acao); got != esperado {
			t.Errorf("acao %q: obteve %q, quer %q", acao, got, esperado)
		}
	}
	c, tm, arqs := contratosTimeline()
	c.Storyboard.Cenas[0].ElementosMoveis = c.Storyboard.Cenas[0].ElementosMoveis[:3]
	c.Storyboard.Cenas[1].ElementosMoveis[2].Acao = "zoom_suave"
	tl, _ := Compor(c, tm, arqs)
	if tl.Batidas[0].Camera != "pan" || tl.Batidas[1].Camera != "zoom" {
		t.Fatalf("camera incorreta: %s %s", tl.Batidas[0].Camera, tl.Batidas[1].Camera)
	}
	c.Storyboard.Cenas[0].ElementosMoveis = c.Storyboard.Cenas[0].ElementosMoveis[:1]
	if tl, _ := Compor(c, tm, arqs); tl.Batidas[0].Camera != "estatica" {
		t.Fatalf("sem elemento de camera deve ser estatica: %s", tl.Batidas[0].Camera)
	}
}

func TestTimelineEntradasInvalidas(t *testing.T) {
	casos := map[string]func(*ContratosRender, *Timings, map[string]string){
		"tipo desconhecido": func(c *ContratosRender, _ *Timings, _ map[string]string) {
			c.Storyboard.Cenas[0].ElementosMoveis[1].Tipo = "hologram"
		},
		"personagem sem sprite": func(c *ContratosRender, _ *Timings, _ map[string]string) {
			c.Storyboard.Cenas[0].ElementosMoveis[0].SpriteID = ""
		},
		"sprite inexistente": func(c *ContratosRender, _ *Timings, _ map[string]string) {
			c.Storyboard.Cenas[0].ElementosMoveis[0].SpriteID = "fantasma"
		},
		"cenario sem arquivo": func(_ *ContratosRender, _ *Timings, a map[string]string) { delete(a, "cen_02") },
		"linha sem timing":    func(_ *ContratosRender, tm *Timings, _ map[string]string) { tm.Linhas = tm.Linhas[:1] },
		"indice fora":         func(c *ContratosRender, _ *Timings, _ map[string]string) { c.Storyboard.Cenas[1].IndiceLinha = 7 },
	}
	for nome, mut := range casos {
		t.Run(nome, func(t *testing.T) {
			c, tm, arqs := contratosTimeline()
			mut(&c, &tm, arqs)
			if _, err := Compor(c, tm, arqs); !errors.Is(err, ErrEntradaTimeline) {
				t.Fatalf("esperava ErrEntradaTimeline, obteve %v", err)
			}
		})
	}
}

func TestTimelineDeterministicaEOrdenaPorLinha(t *testing.T) {
	c, tm, arqs := contratosTimeline()
	c.Storyboard.Cenas[0], c.Storyboard.Cenas[1] = c.Storyboard.Cenas[1], c.Storyboard.Cenas[0]
	a, err := Compor(c, tm, arqs)
	if err != nil || a.Batidas[0].IndiceLinha != 0 {
		t.Fatalf("cenas fora de ordem devem sair ordenadas por linha: %v %+v", err, a.Batidas)
	}
	b, _ := Compor(c, tm, arqs)
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Fatal("duas composicoes devem gerar JSON identico")
	}
}

func TestTimelineMuitosObjetosFicamNoQuadro(t *testing.T) {
	c, tm, arqs := contratosTimeline()
	el := c.Storyboard.Cenas[0].ElementosMoveis
	for i := 0; i < 5; i++ {
		el = append(el, ativos.ElementoMovel{ID: "extra" + string(rune('a'+i)), Tipo: "objeto", Acao: "destaque_foco"})
	}
	c.Storyboard.Cenas[0].ElementosMoveis = el
	tl, err := Compor(c, tm, arqs)
	if err != nil {
		t.Fatal(err)
	}
	if err := tl.Validar(); err != nil {
		t.Fatalf("profundidade deve seguir estritamente crescente com varios objetos: %v", err)
	}
}

func projetoTimeline(t *testing.T, extra map[string]string) string {
	t.Helper()
	timings, _ := json.Marshal(timingsLongo())
	arqs := map[string]string{
		"data/timings.json":  string(timings),
		"flow/cenarios.json": `[{"id":"cen_01","nome":"c","modo":"codigo","arquivo":"flow/cen_01.svg","hash_entrada":"h"}]`,
		"flow/cen_01.svg":    "<svg/>",
	}
	for k, v := range extra {
		arqs[k] = v
	}
	return projetoFixture(t, arqs)
}

func TestTimelineExecutarGravaTimelineJSON(t *testing.T) {
	dir := projetoTimeline(t, nil)
	tl, err := ExecutarTimeline(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if err != nil {
		t.Fatalf("timeline.json ausente: %v", err)
	}
	var lido Timeline
	if err := json.Unmarshal(b, &lido); err != nil || lido.Validar() != nil || len(lido.Batidas) != len(tl.Batidas) {
		t.Fatalf("timeline.json deve reler e validar: %v %v", err, lido.Validar())
	}
	srt, err := os.ReadFile(filepath.Join(dir, "renders", "legenda.srt"))
	if err != nil || !strings.Contains(string(srt), " --> ") {
		t.Fatalf("renders/legenda.srt deve ser gravado com a timeline: %v", err)
	}
}

func TestTimelineExecutarNaoGravaComEntradaInvalida(t *testing.T) {
	casos := map[string]map[string]string{
		"campo desconhecido em cenas": {"cenas.json": jCenasComExtra},
		"indice duplicado":            {"cenas.json": jCenasDuplicadas},
		"cenario fora de flow":        {"flow/cenarios.json": `[]`},
		"arquivo de cenario ausente":  {"flow/cenarios.json": `[{"id":"cen_01","nome":"c","modo":"codigo","arquivo":"flow/nao_existe.svg","hash_entrada":"h"}]`},
		"timings invalido":            {"data/timings.json": "{"},
		"timings com campo extra":     {"data/timings.json": `{"schema_version":1,"x":1}`},
		"timings com outra contagem":  {"data/timings.json": `{"schema_version":1,"provedor":"p","voz":"v","linhas":[{"indice_linha":0,"inicio_s":0,"fim_s":1,"wer":0,"regravacoes":0,"palavras":[{"texto":"a","inicio_s":0,"fim_s":1}]}]}`},
	}
	for nome, extra := range casos {
		t.Run(nome, func(t *testing.T) {
			dir := projetoTimeline(t, extra)
			_, err := ExecutarTimeline(context.Background(), dir)
			if !errors.Is(err, ErrEntradaTimeline) {
				t.Fatalf("esperava ErrEntradaTimeline, obteve %v", err)
			}
			if _, e := os.Stat(filepath.Join(dir, "renders", "legenda.srt")); e == nil {
				t.Fatal("renders/legenda.srt nao pode existir apos falha")
			}
			if _, e := os.Stat(filepath.Join(dir, "timeline.json")); e == nil {
				t.Fatal("timeline.json nao pode existir apos falha")
			}
		})
	}
}

func TestTimelineExecutarSegmentaEAprovaNoGate(t *testing.T) {
	dir := projetoTimeline(t, nil)
	tl, err := ExecutarTimeline(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(tl.Batidas) != 6 {
		t.Fatalf("2 linhas de 9,75 s devem virar 6 batidas de 3,25 s: %d", len(tl.Batidas))
	}
	if m, err := ChecarRitmo(tl); err != nil || m.MedianaS < 3.2 || m.MedianaS > 3.3 {
		t.Fatalf("ritmo da timeline gravada: %+v %v", m, err)
	}
}

func TestTimelineExecutarGateDeRitmoBloqueiaEnaoGrava(t *testing.T) {
	curto, _ := json.Marshal(timingsValido())
	dir := projetoTimeline(t, map[string]string{"data/timings.json": string(curto)})
	if _, err := ExecutarTimeline(context.Background(), dir); !errors.Is(err, ErrTimelineRitmo) {
		t.Fatalf("4 s com 2 batidas (30 cortes/min) deve reprovar o gate: %v", err)
	}
	if _, e := os.Stat(filepath.Join(dir, "renders", "legenda.srt")); e == nil {
		t.Fatal("renders/legenda.srt nao pode existir apos reprovar o ritmo")
	}
	if _, e := os.Stat(filepath.Join(dir, "timeline.json")); e == nil {
		t.Fatal("timeline.json nao pode existir apos reprovar o ritmo")
	}
}

func TestTimelineExecutarSemContratos(t *testing.T) {
	if _, err := ExecutarTimeline(context.Background(), t.TempDir()); !errors.Is(err, ErrEntradaTimeline) || !errors.Is(err, ErrPreVooBloqueado) {
		t.Fatalf("projeto vazio deve bloquear com os dois erros tipados: %v", err)
	}
	ctx, cancela := context.WithCancel(context.Background())
	cancela()
	if _, err := ExecutarTimeline(ctx, projetoTimeline(t, nil)); !errors.Is(err, context.Canceled) {
		t.Fatalf("contexto cancelado deve interromper: %v", err)
	}
}

const (
	jCenasComExtra   = `{"schema_version":1,"projeto_slug":"demo","total_linhas":2,"extra":1,"cenas":[]}`
	jCenasDuplicadas = `{"schema_version":1,"projeto_slug":"demo","total_linhas":2,"cenas":[
{"indice_linha":0,"texto":"a","cenario_id":"cen_01","metafora_visual":"m","elementos_moveis":[]},
{"indice_linha":0,"texto":"b","cenario_id":"cen_01","metafora_visual":"m","elementos_moveis":[]}]}`
)
