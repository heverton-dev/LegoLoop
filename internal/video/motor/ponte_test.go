package motor

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
	"github.com/heverton-dev/LegoLoop/internal/diorama/render"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

// entradaDemo monta a cena 1 do demo: 2 linhas, 1 adulto, 1 criança e 2 fundos do inventário.
func entradaDemo(t *testing.T) (Entrada, string) {
	t.Helper()
	dir := t.TempDir()
	for _, f := range []string{"plates/banco.png", "plates/rua.png"} {
		escreverArquivo(t, filepath.Join(dir, f))
	}
	e := Entrada{
		Cenas: ativos.Storyboard{SchemaVersion: 2, Cenas: []ativos.Cena{
			{IndiceLinha: 0, CenarioID: "banco", ElementosMoveis: []ativos.ElementoMovel{
				{ID: "m1", Tipo: "personagem", PersonagemID: "analista", Pose: "point", Humor: "worry", Fala: true},
			}},
			{IndiceLinha: 1, CenarioID: "rua", ElementosMoveis: []ativos.ElementoMovel{
				{ID: "m2", Tipo: "personagem", PersonagemID: "analista", Pose: "idle", Humor: "smile"},
				{ID: "m3", Tipo: "personagem", PersonagemID: "crianca-teste", Pose: "cheer", Humor: "grin", Fala: true},
			}},
		}},
		Elenco: ElencoPonte{Personagens: []PersonagemPonte{
			{ID: "analista", Cast: map[string]any{"skin": "#E6BE8E", "suit": "#3B4660", "glasses": true}},
			{ID: "crianca-teste", Cast: map[string]any{"skin": "#EBC39A", "proporcao": map[string]any{"cabeca": 1.3, "corpo": 0.62}}},
		}},
		Timings: render.Timings{SchemaVersion: 1, Provedor: "teste", Voz: "v", Linhas: []render.LinhaTempo{
			{IndiceLinha: 0, InicioS: 0.2, FimS: 2.9}, {IndiceLinha: 1, InicioS: 3.1, FimS: 6.4},
		}},
		Fundos: []inventario.Item{
			{ID: "banco", Tipo: "cenario", Arquivo: "plates/banco.png"},
			{ID: "rua", Tipo: "cenario", Arquivo: "plates/rua.png"},
		},
		Estilo:      inventario.Estilo{ID: "molde-de-pilula", Luz: &inventario.Luz{Rim: true, Cor: "#FFB347"}},
		DuracaoVozS: 6.71,
		FPS:         30,
	}
	return e, dir
}

func escreverArquivo(t *testing.T, caminho string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminho, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// AC-E10.3: cada corte tem bg existente e a duração total iguala a voz (± 1 quadro).
func TestPonteMontaCortesComFundoEDuracaoDaVoz(t *testing.T) {
	e, dir := entradaDemo(t)
	tl, err := Montar(e, dir)
	if err != nil {
		t.Fatalf("Montar: %v", err)
	}
	if len(tl.Shots) != 2 || tl.FPS != 30 {
		t.Fatalf("quer 2 cortes a 30 fps, obteve %d a %d", len(tl.Shots), tl.FPS)
	}
	total := 0.0
	for _, s := range tl.Shots {
		if _, err := os.Stat(filepath.Join(dir, s.BG)); err != nil || s.Type != "scene" {
			t.Fatalf("corte %s com bg %q inexistente ou tipo %q", s.ID, s.BG, s.Type)
		}
		total += s.Duration
	}
	if math.Abs(total-e.DuracaoVozS) > 1.0/30 {
		t.Fatalf("duração total %.3f difere da voz %.3f em mais de 1 quadro", total, e.DuracaoVozS)
	}
	c := tl.Shots[1].Chars
	if len(c) != 2 || c[0].Puppet != "analista" || c[1].Mood != "grin" || c[1].Pose != "cheer" || !c[1].Talk {
		t.Fatalf("personagens do corte 2 incorretos: %+v", c)
	}
}

// AC-E10.3: fundo, personagem ou timing ausente retorna erro tipado.
func TestPonteRecusaEntradaIncompleta(t *testing.T) {
	casos := map[string]struct {
		mudar func(*Entrada)
		quer  error
	}{
		"fundo fora do inventario":  {func(e *Entrada) { e.Fundos = e.Fundos[:1] }, ErrFundoAusente},
		"arquivo do fundo ausente":  {func(e *Entrada) { e.Fundos[1].Arquivo = "plates/nao.png" }, ErrFundoAusente},
		"personagem fora do elenco": {func(e *Entrada) { e.Elenco.Personagens = e.Elenco.Personagens[:1] }, ErrPersonagemSemElenco},
		"linha sem timing":          {func(e *Entrada) { e.Timings.Linhas = e.Timings.Linhas[:1] }, ErrLinhaSemTiming},
		"voz mais curta":            {func(e *Entrada) { e.DuracaoVozS = 5 }, ErrDuracaoVoz},
	}
	for nome, c := range casos {
		e, dir := entradaDemo(t)
		c.mudar(&e)
		if _, err := Montar(e, dir); !errors.Is(err, c.quer) {
			t.Errorf("%s: erro %v, quer %v", nome, err, c.quer)
		}
	}
}

// AC-E10.6: cada corte recebe light.rim do estilo; estilo sem luz gera corte sem contorno.
func TestPonteLuzDeContornoDoEstilo(t *testing.T) {
	e, dir := entradaDemo(t)
	tl, err := Montar(e, dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range tl.Shots {
		if s.Light == nil || s.Light.Rim != "#FFB347" {
			t.Fatalf("corte %s sem light.rim do estilo: %+v", s.ID, s.Light)
		}
	}
	e.Estilo.Luz = nil
	tl, err = Montar(e, dir)
	if err != nil {
		t.Fatal(err)
	}
	dados, err := json.Marshal(tl)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(dados), `"light"`) || strings.Contains(string(dados), `"rim"`) {
		t.Fatalf("estilo sem luz gerou contorno: %s", dados)
	}
}

// AC-E10.8: a criança leva cast.proporcao ao motor; o adulto segue sem proporcao.
func TestPonteProporcaoDeCrianca(t *testing.T) {
	e, dir := entradaDemo(t)
	tl, err := Montar(e, dir)
	if err != nil {
		t.Fatal(err)
	}
	prop, ok := tl.Cast["crianca-teste"]["proporcao"].(map[string]any)
	if !ok || prop["cabeca"] != 1.3 || prop["corpo"] != 0.62 {
		t.Fatalf("cast da criança sem proporcao: %+v", tl.Cast["crianca-teste"])
	}
	if _, tem := tl.Cast["analista"]["proporcao"]; tem {
		t.Fatalf("adulto ganhou proporcao: %+v", tl.Cast["analista"])
	}
	e.Elenco.Personagens[1].Cast["proporcao"] = map[string]any{"cabeca": 0.9, "corpo": 0.62}
	if _, err := Montar(e, dir); !errors.Is(err, ErrProporcaoInvalida) {
		t.Fatalf("proporcao com cabeca <= 1 deveria falhar, erro %v", err)
	}
}

// checadorFalso simula o check_timeline.py do motor.
type checadorFalso struct {
	saida string
	dir   string
}

func (c *checadorFalso) Checar(dir string) (string, error) {
	c.dir = dir
	return c.saida, nil
}

// AC-E10.3: Gerar grava timeline.json e exige CHECK OK do motor.
func TestPonteGerarGravaEChecaTimeline(t *testing.T) {
	e, dir := entradaDemo(t)
	ok := &checadorFalso{saida: "2 shots\nCHECK OK\n"}
	if err := Gerar(e, dir, ok); err != nil || ok.dir != dir {
		t.Fatalf("Gerar: %v (dir checado %q)", err, ok.dir)
	}
	dados, err := os.ReadFile(filepath.Join(dir, "timeline.json"))
	if err != nil || !strings.Contains(string(dados), `"shots"`) {
		t.Fatalf("timeline.json ausente ou sem shots: %v", err)
	}
	if err := Gerar(e, dir, &checadorFalso{saida: "CHECK FAILED"}); !errors.Is(err, ErrCheckTimeline) {
		t.Fatalf("CHECK FAILED deveria retornar ErrCheckTimeline, erro %v", err)
	}
}

// AC-E10.3: roda o check_timeline.py real do motor (LEGOLOOP_CHECK_TIMELINE ou motor-local/).
func TestPonteCheckTimelineReal(t *testing.T) {
	script := os.Getenv("LEGOLOOP_CHECK_TIMELINE")
	if script == "" {
		script = filepath.Join("..", "..", "..", "motor-local", "check_timeline.py")
	}
	if _, err := os.Stat(script); err != nil {
		t.Skipf("check_timeline.py do motor ausente (%v)", err)
	}
	if _, err := exec.LookPath("python"); err != nil {
		t.Skipf("python ausente: %v", err)
	}
	e, dir := entradaDemo(t)
	if err := Gerar(e, dir, ChecadorPython{Script: script}); err != nil {
		t.Fatalf("check_timeline real: %v", err)
	}
}
