package roteiro

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrSchemaInvalido       = errors.New("diorama-roteiro: schema invalido")
	ErrConsistenciaConfig   = errors.New("diorama-roteiro: inconsistencia entre projeto.json e marca.json")
	ErrPreVooBloqueado      = errors.New("diorama-roteiro: pre-voo bloqueado")
	ErrTemaRepetido         = errors.New("diorama-roteiro: tema ja registrado no acervo")
	ErrGateRoteiroReprovado = errors.New("diorama-roteiro: gate final de roteiro reprovado")
)

type PoliticaGates struct {
	Modo         string   `json:"modo"`
	GatesHumanos []string `json:"gates_humanos,omitempty"`
}

type Projeto struct {
	SchemaVersion      int           `json:"schema_version"`
	Slug               string        `json:"slug"`
	Canal              string        `json:"canal"`
	Nicho              string        `json:"nicho"`
	PerfisFonte        []string      `json:"perfis_fonte"`
	FontesExtras       []string      `json:"fontes_extras,omitempty"`
	FontesExcluidas    []string      `json:"fontes_excluidas,omitempty"`
	EstruturaNarrativa string        `json:"estrutura_narrativa"`
	Formatos           []string      `json:"formatos"`
	DuracaoAlvoS       int           `json:"duracao_alvo_s"`
	Voz                string        `json:"voz"`
	EstiloVisual       string        `json:"estilo_visual"`
	Animacao           string        `json:"animacao"`
	PoliticaGates      PoliticaGates `json:"politica_gates"`
	VideoReferencia    string        `json:"video_referencia,omitempty"`
	FlowProjeto        string        `json:"flow_projeto,omitempty"`
}

type VozMarca struct {
	Motor      string `json:"motor"`
	Voz        string `json:"voz"`
	Velocidade string `json:"velocidade"`
}

type Marca struct {
	SchemaVersion int      `json:"schema_version"`
	Nome          string   `json:"nome"`
	Formatos      []string `json:"formatos"`
	Voz           VozMarca `json:"voz"`
}

type Fonte struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Nivel    string `json:"nivel"`
	Entidade string `json:"entidade"`
	Data     string `json:"data"`
}

type LinhaRoteiro struct {
	Text       string   `json:"text"`
	Natureza   string   `json:"natureza"`
	Fontes     []string `json:"fontes,omitempty"`
	Atribuicao string   `json:"atribuicao,omitempty"`
	Escopo     string   `json:"escopo,omitempty"`
	Preprint   bool     `json:"preprint,omitempty"`
	Amostra    string   `json:"amostra,omitempty"`
}

type Script struct {
	SchemaVersion int            `json:"schema_version"`
	Title         string         `json:"title"`
	Idioma        string         `json:"idioma"`
	Estrutura     string         `json:"estrutura"`
	TipoPauta     string         `json:"tipo_pauta"`
	Sources       []Fonte        `json:"sources"`
	Lines         []LinhaRoteiro `json:"lines"`
}

// ValidarConfiguracao valida Projeto, Marca, perfis no catalogo e consistencia cruzada.
func ValidarConfiguracao(p Projeto, m Marca, cat CatalogoFontes) error {
	if err := validarProjeto(p, cat); err != nil {
		return err
	}
	if err := validarMarca(m); err != nil {
		return err
	}
	return validarConsistencia(p, m)
}

func validarProjeto(p Projeto, cat CatalogoFontes) error {
	if strings.TrimSpace(p.Nicho) == "" {
		return fmt.Errorf("%w: nicho obrigatorio", ErrSchemaInvalido)
	}
	if len(p.PerfisFonte) == 0 {
		return fmt.Errorf("%w: perfis_fonte obrigatorio", ErrSchemaInvalido)
	}
	for _, perfil := range p.PerfisFonte {
		if !cat.PerfilValido(perfil) {
			return fmt.Errorf("%w: perfil %q ausente no catalogo", ErrSchemaInvalido, perfil)
		}
	}
	if len(p.Formatos) == 0 || p.DuracaoAlvoS <= 0 {
		return fmt.Errorf("%w: formatos e duracao_alvo_s obrigatorios", ErrSchemaInvalido)
	}
	return validarPoliticaGates(p.PoliticaGates)
}

func validarPoliticaGates(pg PoliticaGates) error {
	switch pg.Modo {
	case "humano", "mecanico":
		return nil
	case "hibrido":
		if len(pg.GatesHumanos) == 0 {
			return fmt.Errorf("%w: modo hibrido exige gates_humanos", ErrSchemaInvalido)
		}
		return nil
	default:
		return fmt.Errorf("%w: politica_gates.modo invalido %q", ErrSchemaInvalido, pg.Modo)
	}
}

func validarMarca(m Marca) error {
	if strings.TrimSpace(m.Nome) == "" {
		return fmt.Errorf("%w: marca.nome obrigatorio", ErrSchemaInvalido)
	}
	if len(m.Formatos) == 0 {
		return fmt.Errorf("%w: marca.formatos obrigatorio", ErrSchemaInvalido)
	}
	return nil
}

func validarConsistencia(p Projeto, m Marca) error {
	for _, fmtProj := range p.Formatos {
		if !slices.Contains(m.Formatos, fmtProj) {
			return fmt.Errorf("%w: formato %q ausente em marca.json", ErrConsistenciaConfig, fmtProj)
		}
	}
	if strings.TrimSpace(p.Voz) != "" && strings.TrimSpace(m.Voz.Voz) != "" && p.Voz != m.Voz.Voz {
		return fmt.Errorf("%w: voz %q diverge de marca.json %q", ErrConsistenciaConfig, p.Voz, m.Voz.Voz)
	}
	return nil
}

// ExigeAprovacaoHumana informa se o gateId requer parada humana conforme a politica configurada.
func (pg PoliticaGates) ExigeAprovacaoHumana(gateID string) bool {
	if pg.Modo == "humano" {
		return true
	}
	if pg.Modo == "hibrido" {
		return slices.Contains(pg.GatesHumanos, gateID)
	}
	return false
}
