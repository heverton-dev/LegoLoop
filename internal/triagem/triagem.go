package triagem

import (
	"errors"
	"fmt"
	"strings"
)

const SchemaVersionV1 = "triagem.v1"

// RegraDeterministica define padrões de roteamento sem custo de LLM.
type RegraDeterministica struct {
	PadraoSlug      string   `json:"padrao_slug"`
	ComandosCLI     []string `json:"comandos_cli"`
	ArtefatosAtivos []string `json:"artefatos_ativos"`
	Prioridade      int      `json:"prioridade"`
}

// CriterioAmbiguidade fornece contexto resumido para a camada probabilística.
type CriterioAmbiguidade struct {
	Descricao string   `json:"descricao"`
	Etapas    []string `json:"etapas"`
}

// RegraEtapa define padrões determinísticos para uma etapa específica.
type RegraEtapa struct {
	ID            string   `json:"id"`
	PalavrasChave []string `json:"palavras_chave,omitempty"`
	SlashCommand  string   `json:"slash_command,omitempty"`
	FlagCLI       string   `json:"flag_cli,omitempty"`
}

// TriagemSpec representa o contrato triagem.v1 de um processo.
type TriagemSpec struct {
	SchemaVersion         string                `json:"schema_version"`
	Processo              string                `json:"processo"`
	RegrasDeterministicas []RegraDeterministica `json:"regras_deterministicas"`
	CriterioAmbiguidade   CriterioAmbiguidade   `json:"criterio_ambiguidade"`
	Etapas                []RegraEtapa          `json:"etapas,omitempty"`
	PontosEntrada         []string              `json:"pontos_entrada"`
}

// Validar verifica a conformidade estrutural de TriagemSpec com triagem.v1.
func (s TriagemSpec) Validar() error {
	if s.SchemaVersion != SchemaVersionV1 {
		return fmt.Errorf("schema_version inválido: %q", s.SchemaVersion)
	}
	if strings.TrimSpace(s.Processo) == "" {
		return errors.New("processo obrigatório")
	}
	if len(s.RegrasDeterministicas) == 0 {
		return errors.New("regras_deterministicas obrigatórias")
	}
	if len(s.PontosEntrada) == 0 {
		return errors.New("pontos_entrada obrigatórios")
	}
	return nil
}

// GerarTriagemPadrao cria uma Peça de Triagem triagem.v1 padrão para um processo.
func GerarTriagemPadrao(processo string, etapas []string) TriagemSpec {
	slug := strings.TrimSpace(processo)
	etapasCopia := append([]string(nil), etapas...)
	regrasEtapas := make([]RegraEtapa, 0, len(etapasCopia))
	for _, et := range etapasCopia {
		limpa := limparNomeEtapa(et)
		frase := strings.ReplaceAll(limpa, "-", " ")
		palavras := []string{frase}
		if limpa != frase {
			palavras = append(palavras, limpa)
		}
		regrasEtapas = append(regrasEtapas, RegraEtapa{
			ID:            et,
			PalavrasChave: palavras,
			SlashCommand:  "/" + slug + "-" + limpa,
			FlagCLI:       "--etapa " + et,
		})
	}
	return TriagemSpec{
		SchemaVersion: SchemaVersionV1,
		Processo:      slug,
		RegrasDeterministicas: []RegraDeterministica{
			{
				PadraoSlug:  slug,
				ComandosCLI: []string{"legoloop " + slug, "/legoloop-" + slug},
				Prioridade:  100,
			},
		},
		CriterioAmbiguidade: CriterioAmbiguidade{
			Descricao: fmt.Sprintf("Classificação para o processo %s", slug),
			Etapas:    etapasCopia,
		},
		Etapas:        regrasEtapas,
		PontosEntrada: []string{"cli", "chat", "studio"},
	}
}

func limparNomeEtapa(et string) string {
	s := strings.ReplaceAll(strings.TrimSpace(et), "_", "-")
	partes := strings.Split(s, "-")
	if len(partes) > 1 && temDigito(partes[0]) {
		return strings.Join(partes[1:], "-")
	}
	return s
}

func temDigito(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}
