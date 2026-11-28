package motor

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/ativos"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

// Erros tipados da ponte cenas.json para timeline.json (AC-E10.3, AC-E10.6, AC-E10.8).
var (
	ErrFundoAusente        = errors.New("motor-ponte: fundo ausente do inventario ou do disco")
	ErrPersonagemSemElenco = errors.New("motor-ponte: personagem_id ausente do elenco")
	ErrLinhaSemTiming      = errors.New("motor-ponte: linha sem timing")
	ErrDuracaoVoz          = errors.New("motor-ponte: duracao da voz menor que os timings")
	ErrProporcaoInvalida   = errors.New("motor-ponte: proporcao exige cabeca > 1 e corpo < 1")
	ErrCheckTimeline       = errors.New("motor-ponte: check_timeline reprovou a timeline")
)

// fpsPadrao é a taxa de quadros do motor quando a entrada não define uma.
const fpsPadrao = 30

// PersonagemPonte lê um personagem de elenco.json com o cast completo (cores, booleanos e proporcao).
type PersonagemPonte struct {
	ID   string         `json:"id"`
	Cast map[string]any `json:"cast"`
}

// ElencoPonte lê elenco.json sem perder campos que não são texto.
type ElencoPonte struct {
	Personagens []PersonagemPonte `json:"personagens"`
}

// Entrada reúne cenas.json v2, elenco.json, timings.json, fundos e estilo da cena.
type Entrada struct {
	Cenas       ativos.Storyboard
	Elenco      ElencoPonte
	Timings     render.Timings
	Fundos      []inventario.Item
	Estilo      inventario.Estilo
	DuracaoVozS float64
	FPS         int
	Audio       string
}

// Luz é o bloco light do corte do motor.
type Luz struct {
	Rim     string `json:"rim"`
	RimSide string `json:"rimSide"`
}

// Boneco é um personagem dentro do corte do motor.
type Boneco struct {
	Puppet  string      `json:"puppet"`
	X       float64     `json:"x"`
	Y       float64     `json:"y"`
	S       float64     `json:"s"`
	Mood    string      `json:"mood,omitempty"`
	Pose    string      `json:"pose"`
	Face    int         `json:"face"`
	Talk    bool        `json:"talk,omitempty"`
	Enter   string      `json:"enter,omitempty"`
	SpeakAt *[2]float64 `json:"speakAt,omitempty"`
}

// QuadroCamera define um quadro-chave da câmera do motor (t, x, y, z).
type QuadroCamera struct {
	T float64 `json:"t"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Corte é um shot do timeline.json do motor (scene ou card, com layers, cam, fx e shake).
type Corte struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Lines    [2]int           `json:"lines"`
	Text     string           `json:"text,omitempty"`
	CPS      float64          `json:"cps,omitempty"`
	BG       string           `json:"bg,omitempty"`
	BGFilter string           `json:"bgFilter,omitempty"`
	Set      map[string]any   `json:"set,omitempty"`
	Duration float64          `json:"duration"`
	Chars    []Boneco         `json:"chars"`
	Layers   []map[string]any `json:"layers,omitempty"`
	Cam      []QuadroCamera   `json:"cam,omitempty"`
	FX       []string         `json:"fx,omitempty"`
	Shake    []map[string]any `json:"shake,omitempty"`
	Grade    any              `json:"grade,omitempty"`
	Tin      string           `json:"tin,omitempty"`
	Tout     string           `json:"tout,omitempty"`
	Light    *Luz             `json:"light,omitempty"`
}

// Timeline é o contrato timeline.json lido pelo motor e pelo check_timeline.py.
type Timeline struct {
	FPS   int                       `json:"fps"`
	Audio string                    `json:"audio,omitempty"`
	Music string                    `json:"music,omitempty"`
	Cast  map[string]map[string]any `json:"cast"`
	Shots []Corte                   `json:"shots"`
}

// Checador executa o check_timeline do motor sobre a pasta do projeto.
type Checador interface {
	Checar(dir string) (string, error)
}

// ChecadorPython roda o check_timeline.py real do motor.
type ChecadorPython struct {
	Script string
}

// Checar roda o script e devolve a saída combinada.
func (c ChecadorPython) Checar(dir string) (string, error) {
	saida, err := exec.Command("python", c.Script, dir).CombinedOutput()
	return string(saida), err
}

// Montar converte a entrada em timeline; dir resolve os caminhos relativos dos fundos.
func Montar(e Entrada, dir string) (Timeline, error) {
	if e.FPS <= 0 {
		e.FPS = fpsPadrao
	}
	cast, err := montarCast(e.Elenco)
	if err != nil {
		return Timeline{}, err
	}
	fins, err := finsEmQuadros(e)
	if err != nil {
		return Timeline{}, err
	}
	tl := Timeline{FPS: e.FPS, Audio: e.Audio, Music: e.Cenas.Musica, Cast: cast}
	inicio := 0
	for i, c := range e.Cenas.Cenas {
		cortes, err := montarCortesDaCena(e, dir, c, cast, fins[i]-inicio)
		if err != nil {
			return Timeline{}, err
		}
		inicio = fins[i]
		tl.Shots = append(tl.Shots, cortes...)
	}
	return tl, nil
}

// Gerar grava dir/timeline.json e exige CHECK OK do check_timeline do motor.
func Gerar(e Entrada, dir string, c Checador) error {
	tl, err := Montar(e, dir)
	if err != nil {
		return err
	}
	dados, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "timeline.json"), dados, 0o644); err != nil {
		return err
	}
	saida, err := c.Checar(dir)
	if err != nil || !strings.Contains(saida, "CHECK OK") {
		return fmt.Errorf("%w: %v: %s", ErrCheckTimeline, err, strings.TrimSpace(saida))
	}
	return nil
}

// montarCast copia o cast de cada personagem e valida a proporcao de criança.
func montarCast(el ElencoPonte) (map[string]map[string]any, error) {
	cast := map[string]map[string]any{}
	for _, p := range el.Personagens {
		copia := map[string]any{}
		maps.Copy(copia, p.Cast)
		if prop, tem := copia["proporcao"]; tem {
			if err := validarProporcao(prop); err != nil {
				return nil, fmt.Errorf("%w: personagem %q", err, p.ID)
			}
		}
		cast[p.ID] = copia
	}
	return cast, nil
}

// validarProporcao exige cabeca maior que 1 e corpo entre 0 e 1 (AC-E10.8).
func validarProporcao(prop any) error {
	m, ok := prop.(map[string]any)
	if !ok {
		return ErrProporcaoInvalida
	}
	cabeca, okC := m["cabeca"].(float64)
	corpo, okP := m["corpo"].(float64)
	if !okC || !okP || cabeca <= 1 || corpo <= 0 || corpo >= 1 {
		return ErrProporcaoInvalida
	}
	return nil
}

// finsEmQuadros devolve o quadro final de cada cena; a última termina no fim da voz.
func finsEmQuadros(e Entrada) ([]int, error) {
	fps := float64(e.FPS)
	fins := make([]int, len(e.Cenas.Cenas))
	anterior := 0
	for i, c := range e.Cenas.Cenas {
		if c.IndiceLinha < 0 || c.IndiceLinha >= len(e.Timings.Linhas) {
			return nil, fmt.Errorf("%w: linha %d", ErrLinhaSemTiming, c.IndiceLinha)
		}
		fins[i] = int(math.Round(e.Timings.Linhas[c.IndiceLinha].FimS * fps))
		if i == len(fins)-1 {
			voz := int(math.Round(e.DuracaoVozS * fps))
			if voz < fins[i] {
				return nil, fmt.Errorf("%w: voz %.3fs, timings ate %.3fs", ErrDuracaoVoz, e.DuracaoVozS, float64(fins[i])/fps)
			}
			fins[i] = voz
		}
		if fins[i] <= anterior {
			return nil, fmt.Errorf("%w: linha %d sem quadros", ErrLinhaSemTiming, c.IndiceLinha)
		}
		anterior = fins[i]
	}
	return fins, nil
}

// resolverFundo acha o cenario no inventario: item com Arquivo vira bg-imagem (confere o arquivo
// em disco); item sem Arquivo (codigo-desenhado, sets.js) vira set a partir de Parametros.
func resolverFundo(fundos []inventario.Item, dir, id string) (string, map[string]any, error) {
	for _, f := range fundos {
		if f.ID != id || f.Tipo != "cenario" {
			continue
		}
		if f.Arquivo == "" {
			return "", f.Parametros, nil
		}
		caminho := f.Arquivo
		if !filepath.IsAbs(caminho) {
			caminho = filepath.Join(dir, caminho)
		}
		if _, err := os.Stat(caminho); err != nil {
			return "", nil, fmt.Errorf("%w: %q (%v)", ErrFundoAusente, id, err)
		}
		return filepath.ToSlash(f.Arquivo), nil, nil
	}
	return "", nil, fmt.Errorf("%w: cenario %q fora do inventario", ErrFundoAusente, id)
}

// montarBonecos posiciona os personagens da cena usando X/Y/Escala explicitos ou faixas iguais.
func montarBonecos(elementos []ativos.ElementoMovel, cast map[string]map[string]any) ([]Boneco, error) {
	var ids []ativos.ElementoMovel
	for _, el := range elementos {
		if el.PersonagemID != "" {
			ids = append(ids, el)
		}
	}
	bonecos := []Boneco{}
	for k, el := range ids {
		if _, ok := cast[el.PersonagemID]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrPersonagemSemElenco, el.PersonagemID)
		}
		x, y, s := posicaoBoneco(el, k, len(ids))
		bonecos = append(bonecos, Boneco{
			Puppet: el.PersonagemID, X: x, Y: y, S: s,
			Mood: el.Humor, Pose: poseOuIdle(el.Pose), Face: direcao(el.Olhar, x),
			Talk: el.Fala, Enter: el.Entrada,
		})
	}
	return bonecos, nil
}

func posicaoBoneco(el ativos.ElementoMovel, k, total int) (float64, float64, float64) {
	x := el.X
	if x <= 0 {
		x = 1920 * float64(k+1) / float64(total+1)
	}
	y := el.Y
	if y <= 0 {
		y = 1000
	}
	s := el.Escala
	if s <= 0 {
		s = 1
	}
	return x, y, s
}

// poseOuIdle usa idle quando a cena não define pose.
func poseOuIdle(pose string) string {
	if pose == "" {
		return "idle"
	}
	return pose
}

// direcao vira o boneco pelo olhar; sem olhar, ele encara o centro do quadro.
func direcao(olhar string, x float64) int {
	switch olhar {
	case "esquerda":
		return -1
	case "direita":
		return 1
	}
	if x < 960 {
		return 1
	}
	return -1
}

// luzDoEstilo gera light.rim só quando o estilo declara luz com rim (AC-E10.6).
func luzDoEstilo(est inventario.Estilo) *Luz {
	if est.Luz == nil || !est.Luz.Rim {
		return nil
	}
	cor := est.Luz.Cor
	if cor == "" {
		cor = "#FFD27A"
	}
	return &Luz{Rim: cor, RimSide: "left"}
}
