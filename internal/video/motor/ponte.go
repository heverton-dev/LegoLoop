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

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
	"github.com/heverton-dev/LegoLoop/internal/diorama/render"
	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
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
	Puppet string  `json:"puppet"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	S      float64 `json:"s"`
	Mood   string  `json:"mood,omitempty"`
	Pose   string  `json:"pose"`
	Face   int     `json:"face"`
	Talk   bool    `json:"talk,omitempty"`
}

// Corte é um shot do timeline.json do motor.
type Corte struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Lines    [2]int   `json:"lines"`
	BG       string   `json:"bg"`
	Duration float64  `json:"duration"`
	Chars    []Boneco `json:"chars"`
	Light    *Luz     `json:"light,omitempty"`
}

// Timeline é o contrato timeline.json lido pelo motor e pelo check_timeline.py.
type Timeline struct {
	FPS   int                       `json:"fps"`
	Audio string                    `json:"audio,omitempty"`
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
	tl := Timeline{FPS: e.FPS, Audio: e.Audio, Cast: cast}
	inicio := 0
	for i, c := range e.Cenas.Cenas {
		corte, err := montarCorte(e, dir, c, cast)
		if err != nil {
			return Timeline{}, err
		}
		corte.Duration = float64(fins[i]-inicio) / float64(e.FPS)
		inicio = fins[i]
		tl.Shots = append(tl.Shots, corte)
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

// montarCorte resolve o fundo, os bonecos e a luz de contorno de uma cena.
func montarCorte(e Entrada, dir string, c ativos.Cena, cast map[string]map[string]any) (Corte, error) {
	bg, err := resolverFundo(e.Fundos, dir, c.CenarioID)
	if err != nil {
		return Corte{}, err
	}
	bonecos, err := montarBonecos(c.ElementosMoveis, cast)
	if err != nil {
		return Corte{}, err
	}
	return Corte{
		ID:    fmt.Sprintf("l%03d", c.IndiceLinha),
		Type:  "scene",
		Lines: [2]int{c.IndiceLinha, c.IndiceLinha},
		BG:    bg,
		Chars: bonecos,
		Light: luzDoEstilo(e.Estilo),
	}, nil
}

// resolverFundo acha o cenário no inventário e confere o arquivo em disco.
func resolverFundo(fundos []inventario.Item, dir, id string) (string, error) {
	for _, f := range fundos {
		if f.ID != id || f.Tipo != "cenario" || f.Arquivo == "" {
			continue
		}
		caminho := f.Arquivo
		if !filepath.IsAbs(caminho) {
			caminho = filepath.Join(dir, caminho)
		}
		if _, err := os.Stat(caminho); err != nil {
			return "", fmt.Errorf("%w: %q (%v)", ErrFundoAusente, id, err)
		}
		return filepath.ToSlash(f.Arquivo), nil
	}
	return "", fmt.Errorf("%w: cenario %q fora do inventario", ErrFundoAusente, id)
}

// montarBonecos posiciona os personagens da cena em faixas iguais sobre o piso.
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
		x := 1920 * float64(k+1) / float64(len(ids)+1)
		bonecos = append(bonecos, Boneco{
			Puppet: el.PersonagemID, X: x, Y: 1000, S: 1,
			Mood: el.Humor, Pose: poseOuIdle(el.Pose), Face: direcao(el.Olhar, x), Talk: el.Fala,
		})
	}
	return bonecos, nil
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
