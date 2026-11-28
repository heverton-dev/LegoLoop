package ativos

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
)

const (
	DimensaoMinimaSprite = 64
	MinElementosMoveis   = 4
)

var (
	ErrPreVooBloqueado    = errors.New("diorama-ativos: pre-voo bloqueado")
	ErrFichaInvalida      = errors.New("diorama-ativos: ficha de elenco invalida")
	ErrSpriteInvalido     = errors.New("diorama-ativos: descritor de sprite invalido")
	ErrSobreposicaoSprite = errors.New("diorama-ativos: sobreposicao de bounding_box detectada")
	ErrFlowTravaAusente   = errors.New("diorama-ativos: trava flow_projeto obrigatoria")
	ErrDensidadeCenarios  = errors.New("diorama-ativos: densidade insuficiente de cenarios")
	ErrStoryboardInvalido = errors.New("diorama-ativos: gate de storyboard reprovado")

	ErrElencoInvalido         = errors.New("diorama-ativos: elenco.json invalido")
	ErrParteNaoSuportada      = fmt.Errorf("diorama-ativos: %w", inventario.ErrParteNaoSuportada)
	ErrPersonagemForaDoElenco = errors.New("diorama-ativos: personagem_id ausente do elenco.json")
	ErrItemForaDoInventario   = errors.New("diorama-ativos: item ausente do inventario do estilo")
)

// Elenco e o contrato de elenco.json v2: cada personagem traz um bloco cast do inventario do estilo.
type Elenco struct {
	SchemaVersion int                `json:"schema_version"`
	Estilo        string             `json:"estilo"`
	Personagens   []PersonagemElenco `json:"personagens"`
}

// PersonagemElenco liga um id de personagem ao bloco cast (skin, suit, hair, hat...).
type PersonagemElenco struct {
	ID   string            `json:"id"`
	Nome string            `json:"nome"`
	Cast map[string]string `json:"cast"`
}

// InventarioEstilo resume o estilo: partes que o renderizador desenha e ids aprovados por tipo.
type InventarioEstilo struct {
	Estilo           string              `json:"estilo"`
	PartesSuportadas []string            `json:"partes_suportadas"`
	Itens            map[string][]string `json:"itens"`
}

// Tem indica se o inventario do estilo guarda o item id do tipo informado.
func (i InventarioEstilo) Tem(tipo, id string) bool {
	return slices.Contains(i.Itens[tipo], id)
}

// Suporta indica se o renderizador do estilo desenha a parte informada.
func (i InventarioEstilo) Suporta(parte string) bool {
	return slices.Contains(i.PartesSuportadas, parte)
}

// tipoCampoCast mapeia o campo do cast ao tipo do inventario; forma exige partes_suportadas.
func tipoCampoCast(campo string) (tipo string, forma bool, ok bool) {
	switch campo {
	case "skin":
		return "pele", false, true
	case "suit", "shirt":
		return "traje", false, true
	case "tie":
		return "adorno", false, true
	case "hair":
		return "penteado", true, true
	case "hat", "glasses":
		return "adorno", true, true
	case "prop":
		return "aderecos", true, true
	}
	return "", false, false
}

// campoCor indica campo de cor do cast (hair_col, suit_col...), livre de inventario.
func campoCor(campo string) bool {
	return strings.HasSuffix(campo, "_col")
}

func fichasObrigatorias() []string {
	return []string{"turnaround", "expressoes_3x2", "poses_6x"}
}

type BoundingBox struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Sobrepoe indica se ha intersecao de area positiva entre duas caixas delimitadoras.
func (b BoundingBox) Sobrepoe(outro BoundingBox) bool {
	if b.W <= 0 || b.H <= 0 || outro.W <= 0 || outro.H <= 0 {
		return false
	}
	return b.X < outro.X+outro.W && b.X+b.W > outro.X &&
		b.Y < outro.Y+outro.H && b.Y+b.H > outro.Y
}

type FichaPersonagem struct {
	ID                string            `json:"id"`
	Nome              string            `json:"nome"`
	DescricaoImutavel string            `json:"descricao_imutavel"`
	Fichas            map[string]string `json:"fichas"`
}

// Validar verifica atributos obrigatorios e as 3 fichas do personagem.
func (f FichaPersonagem) Validar() error {
	if strings.TrimSpace(f.ID) == "" || strings.TrimSpace(f.Nome) == "" {
		return fmt.Errorf("%w: id e nome obrigatorios", ErrFichaInvalida)
	}
	if strings.TrimSpace(f.DescricaoImutavel) == "" {
		return fmt.Errorf("%w: descricao_imutavel obrigatoria para %q", ErrFichaInvalida, f.ID)
	}
	for _, chave := range fichasObrigatorias() {
		if strings.TrimSpace(f.Fichas[chave]) == "" {
			return fmt.Errorf("%w: ficha %q ausente em %q", ErrFichaInvalida, chave, f.ID)
		}
	}
	return nil
}

type DescritorSprite struct {
	ID           string      `json:"id"`
	PersonagemID string      `json:"personagem_id"`
	Tipo         string      `json:"tipo"`
	ArquivoPNG   string      `json:"arquivo_png"`
	Largura      int         `json:"largura"`
	Altura       int         `json:"altura"`
	AlfaValido   bool        `json:"alfa_valido"`
	DespillVerde bool        `json:"despill_verde"`
	BoundingBox  BoundingBox `json:"bounding_box"`
}

// Validar garante canal alfa, despill verde, dimensoes minimas e caixa valida.
func (d DescritorSprite) Validar() error {
	if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.PersonagemID) == "" {
		return fmt.Errorf("%w: id e personagem_id obrigatorios", ErrSpriteInvalido)
	}
	if !tipoSpriteValido(d.Tipo) || strings.TrimSpace(d.ArquivoPNG) == "" {
		return fmt.Errorf("%w: tipo ou arquivo_png invalido em %q", ErrSpriteInvalido, d.ID)
	}
	if !d.AlfaValido || !d.DespillVerde {
		return fmt.Errorf("%w: alfa_valido e despill_verde obrigatorios em %q", ErrSpriteInvalido, d.ID)
	}
	if d.Largura < DimensaoMinimaSprite || d.Altura < DimensaoMinimaSprite {
		return fmt.Errorf("%w: dimensao minima %dpx violada em %q", ErrSpriteInvalido, DimensaoMinimaSprite, d.ID)
	}
	if d.BoundingBox.X < 0 || d.BoundingBox.Y < 0 || d.BoundingBox.W <= 0 || d.BoundingBox.H <= 0 {
		return fmt.Errorf("%w: bounding_box invalida em %q", ErrSpriteInvalido, d.ID)
	}
	return nil
}

func tipoSpriteValido(tipo string) bool {
	return tipo == "pose" || tipo == "expressao" || tipo == "objeto"
}

type Cenario struct {
	ID          string   `json:"id"`
	Nome        string   `json:"nome"`
	Modo        string   `json:"modo"`
	Arquivo     string   `json:"arquivo"`
	HashEntrada string   `json:"hash_entrada"`
	Referencias []string `json:"referencias,omitempty"`
	FlowProjeto string   `json:"flow_projeto,omitempty"`
}

// ElementoMovel v2: personagem usa personagem_id, pose, humor, olhar e fala; sprite_id fica para estilos de sprite.
type ElementoMovel struct {
	ID           string  `json:"id"`
	Tipo         string  `json:"tipo"`
	SpriteID     string  `json:"sprite_id,omitempty"`
	Acao         string  `json:"acao,omitempty"`
	PersonagemID string  `json:"personagem_id,omitempty"`
	Pose         string  `json:"pose,omitempty"`
	Humor        string  `json:"humor,omitempty"`
	Olhar        string  `json:"olhar,omitempty"`
	Fala         bool    `json:"fala,omitempty"`
	X            float64 `json:"x,omitempty"`
	Y            float64 `json:"y,omitempty"`
	Escala       float64 `json:"escala,omitempty"`
	Entrada      string  `json:"entrada,omitempty"`
}

type Cena struct {
	IndiceLinha       int              `json:"indice_linha"`
	Texto             string           `json:"texto"`
	CenarioID         string           `json:"cenario_id"`
	CenarioID2        string           `json:"cenario_id_2,omitempty"`
	ParametrosCenario map[string]any   `json:"parametros_cenario,omitempty"`
	MetaforaVisual    string           `json:"metafora_visual"`
	Cartao            bool             `json:"cartao,omitempty"`
	TextoCartao       string           `json:"texto_cartao,omitempty"`
	TextoTela         string           `json:"texto_tela,omitempty"`
	CorTexto          string           `json:"cor_texto,omitempty"`
	TextoTelaX        float64          `json:"texto_tela_x,omitempty"`
	TextoTelaY        float64          `json:"texto_tela_y,omitempty"`
	TextoTela2        string           `json:"texto_tela_2,omitempty"`
	CorTexto2         string           `json:"cor_texto_2,omitempty"`
	Camera            string           `json:"camera,omitempty"`
	FX                []string         `json:"fx,omitempty"`
	Shake             []map[string]any `json:"shake,omitempty"`
	Grade             any              `json:"grade,omitempty"`
	TransicaoEntrada  string           `json:"tin,omitempty"`
	TransicaoSaida    string           `json:"tout,omitempty"`
	Camadas           []map[string]any `json:"camadas,omitempty"`
	Camadas2          []map[string]any `json:"camadas_2,omitempty"`
	ElementosMoveis   []ElementoMovel  `json:"elementos_moveis"`
}

type Storyboard struct {
	SchemaVersion int     `json:"schema_version"`
	ProjetoSlug   string  `json:"projeto_slug"`
	TotalLinhas   int     `json:"total_linhas"`
	AlvoCorteS    float64 `json:"alvo_corte_s,omitempty"`
	Musica        string  `json:"musica,omitempty"`
	Cenas         []Cena  `json:"cenas"`
}

type BuildNotes struct {
	CorpoMaterial     string `json:"corpo_material"`
	OlhosDetalhes     string `json:"olhos_detalhes"`
	AcessoriosEfeitos string `json:"acessorios_efeitos"`
	Regras            string `json:"regras"`
}

type FichaAtivo struct {
	ID                string     `json:"id"`
	Categoria         string     `json:"categoria"` // "personagem" | "objeto" | "motion"
	Nome              string     `json:"nome"`
	DescricaoImutavel string     `json:"descricao_imutavel"`
	Turnaround        []string   `json:"turnaround"`
	Expressoes        []string   `json:"expressoes"`
	Poses             []string   `json:"poses"`
	BuildNotes        BuildNotes `json:"build_notes"`
	FolhaMestrePNG    string     `json:"folha_mestre_png,omitempty"`
	PromptFlow        string     `json:"prompt_flow,omitempty"`
	HashEntrada       string     `json:"hash_entrada,omitempty"`
}

// Validar verifica as 4 faixas obrigatorias (turnaround, expressoes, poses e build_notes).
func (f FichaAtivo) Validar() error {
	if strings.TrimSpace(f.ID) == "" || strings.TrimSpace(f.Nome) == "" || strings.TrimSpace(f.DescricaoImutavel) == "" {
		return fmt.Errorf("%w: id, nome e descricao_imutavel obrigatorios", ErrFichaInvalida)
	}
	if f.Categoria != "personagem" && f.Categoria != "objeto" && f.Categoria != "motion" {
		return fmt.Errorf("%w: categoria invalida %q", ErrFichaInvalida, f.Categoria)
	}
	if len(f.Turnaround) < 3 || len(f.Expressoes) < 3 || len(f.Poses) < 3 {
		return fmt.Errorf("%w: turnaround, expressoes e poses exigem >= 3 itens em %q", ErrFichaInvalida, f.ID)
	}
	if strings.TrimSpace(f.BuildNotes.CorpoMaterial) == "" || strings.TrimSpace(f.BuildNotes.Regras) == "" {
		return fmt.Errorf("%w: build_notes (corpo_material e regras) obrigatorios em %q", ErrFichaInvalida, f.ID)
	}
	return nil
}

// ModoVisualValido verifica se o estilo_visual esta entre os modos suportados.
func ModoVisualValido(modo string) bool {
	return modo == "codigo" || modo == "flow" || modo == "hibrido"
}
