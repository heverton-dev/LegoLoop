// Package render implementa o Ciclo 3 do macroprocesso legoloop-video-diorama:
// voz neural, timeline em camadas, mixagem e renderizacao multiformato.
package render

import (
	"errors"
	"fmt"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
)

const (
	ResolucaoPadrao   = "1080p"
	ReprovasWERPadrao = 3
	LimiteWERLinha    = 0.15
	LimiteWERGlobal   = 0.05
)

var (
	ErrPreVooBloqueado      = errors.New("diorama-render: pre-voo bloqueado")
	ErrProvedorDesconhecido = errors.New("diorama-render: provedor de voz nao registrado")
	ErrCreditosNegados      = errors.New("diorama-render: gate de creditos negou chamada paga")
	ErrWERAcimaLimite       = errors.New("diorama-render: WER acima do limite")
	ErrTimelineRitmo        = errors.New("diorama-render: timeline fora do ritmo")
	ErrLoudnessReprovado    = errors.New("diorama-render: loudness fora do alvo")
	ErrRenderInvalido       = errors.New("diorama-render: render reprovado")
	ErrCorteVertical        = errors.New("diorama-render: corte 9:16 reprovado")
	ErrSchemaInvalido       = errors.New("diorama-render: schema invalido")
)

func resolucaoValida(n string) bool { return n == "720p" || n == "1080p" || n == "4k" }

func cameraValida(n string) bool {
	return n == "zoom" || n == "pan" || n == "shake" || n == "estatica"
}

func camadaValida(n string) bool {
	switch n {
	case "cenario", "aderecos", "sprite", "efeito", "legenda":
		return true
	}
	return false
}

func sfxValido(n string) bool {
	return n == "troca_cena" || n == "movimento_objeto" || n == "destaque"
}

func invalido(formato string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrSchemaInvalido, fmt.Sprintf(formato, args...))
}

// ConfigRender e o render.json opcional do Ciclo 3.
type ConfigRender struct {
	SchemaVersion  int      `json:"schema_version"`
	Resolucoes     []string `json:"resolucoes"`
	BGM            string   `json:"bgm,omitempty"`
	BGMGanhoDB     *float64 `json:"bgm_ganho_db,omitempty"`
	SFXCatalogo    string   `json:"sfx_catalogo,omitempty"`
	ReprovasWERMax int      `json:"reprovas_wer_max"`
}

// ConfigPadrao devolve a configuracao usada na ausencia de render.json.
func ConfigPadrao() ConfigRender {
	return ConfigRender{SchemaVersion: 1, Resolucoes: []string{ResolucaoPadrao}, ReprovasWERMax: ReprovasWERPadrao}
}

// Normalizar preenche apenas os campos ausentes com os padroes.
func (c ConfigRender) Normalizar() ConfigRender {
	if len(c.Resolucoes) == 0 {
		c.Resolucoes = []string{ResolucaoPadrao}
	}
	if c.ReprovasWERMax == 0 {
		c.ReprovasWERMax = ReprovasWERPadrao
	}
	if c.BGMGanhoDB == nil {
		padrao := BGMGanhoPadraoDB
		c.BGMGanhoDB = &padrao
	}
	return c
}

// Validar confere versao, resolucoes unicas e limite de regravacoes.
func (c ConfigRender) Validar() error {
	if c.SchemaVersion != 1 {
		return invalido("render.json schema_version deve ser 1")
	}
	if c.ReprovasWERMax < 0 {
		return invalido("reprovas_wer_max negativo")
	}
	if g := c.BGMGanhoDB; g != nil && (*g < BGMGanhoMinDB || *g > BGMGanhoMaxDB) {
		return invalido("bgm_ganho_db fora de %.0f a %.0f dB", BGMGanhoMinDB, BGMGanhoMaxDB)
	}
	vistas := map[string]bool{}
	for _, r := range c.Resolucoes {
		if !resolucaoValida(r) || vistas[r] {
			return invalido("resolucao %q invalida ou repetida", r)
		}
		vistas[r] = true
	}
	return nil
}

type PalavraTempo struct {
	Texto   string  `json:"texto"`
	InicioS float64 `json:"inicio_s"`
	FimS    float64 `json:"fim_s"`
}

type VisemaTempo struct {
	Viseme  string  `json:"viseme"`
	InicioS float64 `json:"inicio_s"`
	FimS    float64 `json:"fim_s"`
}

type LinhaTempo struct {
	IndiceLinha int            `json:"indice_linha"`
	InicioS     float64        `json:"inicio_s"`
	FimS        float64        `json:"fim_s"`
	WER         float64        `json:"wer"`
	Regravacoes int            `json:"regravacoes"`
	Palavras    []PalavraTempo `json:"palavras"`
	Visemas     []VisemaTempo  `json:"visemas,omitempty"`
}

// Timings e o contrato data/timings.json, alinhado sempre pelo Whisper.
type Timings struct {
	SchemaVersion int          `json:"schema_version"`
	Provedor      string       `json:"provedor"`
	Voz           string       `json:"voz"`
	WERGlobal     float64      `json:"wer_global"`
	Linhas        []LinhaTempo `json:"linhas"`
}

// Validar confere ordem, janelas e ausencia de sobreposicao.
func (t Timings) Validar() error {
	if t.SchemaVersion != 1 || strings.TrimSpace(t.Provedor) == "" || strings.TrimSpace(t.Voz) == "" {
		return invalido("timings exige schema_version 1, provedor e voz")
	}
	if len(t.Linhas) == 0 {
		return invalido("timings sem linhas")
	}
	if t.WERGlobal < 0 || t.WERGlobal > 1 {
		return invalido("wer_global fora de 0 a 1")
	}
	fim := 0.0
	for i, l := range t.Linhas {
		if l.IndiceLinha != i {
			return invalido("linha %d com indice_linha %d", i, l.IndiceLinha)
		}
		if l.InicioS < fim {
			return invalido("linha %d sobrepoe a anterior", i)
		}
		if err := l.validar(); err != nil {
			return err
		}
		fim = l.FimS
	}
	return nil
}

func (l LinhaTempo) validar() error {
	if l.InicioS < 0 || l.FimS <= l.InicioS {
		return invalido("linha %d com janela invalida", l.IndiceLinha)
	}
	if l.WER < 0 || l.WER > 1 || l.Regravacoes < 0 {
		return invalido("linha %d com WER ou regravacoes invalidos", l.IndiceLinha)
	}
	if len(l.Palavras) == 0 {
		return invalido("linha %d sem palavras", l.IndiceLinha)
	}
	cursor := l.InicioS
	for _, p := range l.Palavras {
		if strings.TrimSpace(p.Texto) == "" || p.InicioS < cursor || p.FimS < p.InicioS || p.FimS > l.FimS {
			return invalido("linha %d com palavra %q fora da janela ou sobreposta", l.IndiceLinha, p.Texto)
		}
		cursor = p.FimS
	}
	return l.validarVisemas()
}

func (l LinhaTempo) validarVisemas() error {
	cursor := l.InicioS
	for _, v := range l.Visemas {
		if strings.TrimSpace(v.Viseme) == "" || v.InicioS < cursor || v.FimS < v.InicioS || v.FimS > l.FimS {
			return invalido("linha %d com viseme fora da janela", l.IndiceLinha)
		}
		cursor = v.InicioS
	}
	return nil
}

type CamadaItem struct {
	ID           string             `json:"id"`
	Camada       string             `json:"camada"`
	Profundidade int                `json:"profundidade"`
	SpriteID     string             `json:"sprite_id,omitempty"`
	Arquivo      string             `json:"arquivo,omitempty"`
	Palavras     []PalavraTempo     `json:"palavras,omitempty"`
	Caixa        ativos.BoundingBox `json:"caixa"`
}

type EventoSFX struct {
	Tipo    string  `json:"tipo"`
	TempoS  float64 `json:"tempo_s"`
	Arquivo string  `json:"arquivo"`
}

type Batida struct {
	IndiceLinha int          `json:"indice_linha"`
	InicioS     float64      `json:"inicio_s"`
	FimS        float64      `json:"fim_s"`
	CenarioID   string       `json:"cenario_id"`
	Camera      string       `json:"camera"`
	Camadas     []CamadaItem `json:"camadas"`
	SFX         []EventoSFX  `json:"sfx,omitempty"`
}

// Timeline e o contrato timeline.json consumido por mix e render.
type Timeline struct {
	SchemaVersion int      `json:"schema_version"`
	DuracaoS      float64  `json:"duracao_s"`
	Batidas       []Batida `json:"batidas"`
}

// Validar confere continuidade das batidas, camadas e eventos de SFX.
func (t Timeline) Validar() error {
	if t.SchemaVersion != 1 || t.DuracaoS <= 0 || len(t.Batidas) == 0 {
		return invalido("timeline exige schema_version 1, duracao e batidas")
	}
	cursor := 0.0
	for i, b := range t.Batidas {
		if b.InicioS < cursor || b.FimS <= b.InicioS || b.FimS > t.DuracaoS {
			return invalido("batida %d com janela invalida", i)
		}
		if err := b.validar(i); err != nil {
			return err
		}
		cursor = b.FimS
	}
	return nil
}

func (b Batida) validar(i int) error {
	if strings.TrimSpace(b.CenarioID) == "" || !cameraValida(b.Camera) {
		return invalido("batida %d sem cenario ou com camera %q invalida", i, b.Camera)
	}
	if err := b.validarCamadas(i); err != nil {
		return err
	}
	for _, s := range b.SFX {
		if !sfxValido(s.Tipo) || strings.TrimSpace(s.Arquivo) == "" || s.TempoS < b.InicioS || s.TempoS > b.FimS {
			return invalido("batida %d com SFX %q invalido", i, s.Tipo)
		}
	}
	return nil
}

func (b Batida) validarCamadas(i int) error {
	prof := -1
	for _, c := range b.Camadas {
		if !camadaValida(c.Camada) || c.Profundidade <= prof {
			return invalido("batida %d com camada %q invalida ou profundidade nao crescente", i, c.Camada)
		}
		prof = c.Profundidade
		if err := b.validarPalavrasLegenda(i, c); err != nil {
			return err
		}
	}
	return nil
}

// validarPalavrasLegenda exige palavras de karaoke dentro da batida, em ordem e sem sobreposicao.
func (b Batida) validarPalavrasLegenda(i int, c CamadaItem) error {
	cursor := b.InicioS
	for _, p := range c.Palavras {
		if strings.TrimSpace(p.Texto) == "" || p.InicioS < cursor || p.FimS < p.InicioS || p.FimS > b.FimS {
			return invalido("batida %d com palavra de legenda %q fora da janela ou sobreposta", i, p.Texto)
		}
		cursor = p.FimS
	}
	return nil
}

// RelatorioRender registra o que foi usado e medido no render.
type RelatorioRender struct {
	Encoder      string            `json:"encoder"`
	ProvedorVoz  string            `json:"provedor_voz"`
	LUFS         float64           `json:"lufs"`
	TruePeakDBTP float64           `json:"true_peak_dbtp"`
	Renders      map[string]string `json:"renders,omitempty"`
	Shorts       []string          `json:"shorts,omitempty"`
}

// Validar exige encoder e provedor de voz registrados.
func (r RelatorioRender) Validar() error {
	if strings.TrimSpace(r.Encoder) == "" || strings.TrimSpace(r.ProvedorVoz) == "" {
		return invalido("relatorio exige encoder e provedor_voz")
	}
	return nil
}
