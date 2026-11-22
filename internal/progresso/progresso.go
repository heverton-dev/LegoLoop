package progresso

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
)

// SchemaVersao identifica a versão canônica do arquivo de progresso vivo.
const SchemaVersao = "progresso.v1"

// Estados canônicos de um Ticket no arquivo de progresso vivo.
const (
	EstadoPendente    = "PENDING"
	EstadoEmAndamento = "IN_PROGRESS"
	EstadoAprovado    = "PASSED"
	EstadoReprovado   = "FAILED"
	EstadoBloqueado   = "BLOCKED"
)

// EntradaTicket espelha o formato de entrada de *_TICKETS_*.json.
type EntradaTicket struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	EpicID        string   `json:"epic_id"`
	InitiativeID  string   `json:"initiative_id"`
	TargetFiles   []string `json:"target_files"`
	SpecCriteria  []string `json:"spec_criteria"`
	ValidationCmd string   `json:"validation_cmd"`
}

// Resumo consolida a contagem de Tickets por estado.
type Resumo struct {
	Total       int `json:"total"`
	Aprovados   int `json:"aprovados"`
	EmAndamento int `json:"em_andamento"`
	Reprovados  int `json:"reprovados"`
	Bloqueados  int `json:"bloqueados"`
	Pendentes   int `json:"pendentes"`
}

// ItemTicket guarda o estado vivo e a prova do Gate de um Ticket.
type ItemTicket struct {
	ID           string   `json:"id"`
	Epico        string   `json:"epico"`
	Titulo       string   `json:"titulo"`
	Estado       string   `json:"estado"`
	Responsavel  string   `json:"responsavel"`
	Comando      string   `json:"comando"`
	ExitCode     int      `json:"exit_code"`
	Prova        string   `json:"prova"`
	Alvos        []string `json:"alvos,omitempty"`
	Criterios    []string `json:"criterios,omitempty"`
	AtualizadoEm string   `json:"atualizado_em"`
}

// Progresso representa o arquivo vivo progresso.v1 de uma iniciativa.
type Progresso struct {
	Schema       string       `json:"schema"`
	Iniciativa   string       `json:"iniciativa"`
	Plano        string       `json:"plano"`
	AtualizadoEm string       `json:"atualizado_em"`
	Resumo       Resumo       `json:"resumo"`
	Tickets      []ItemTicket `json:"tickets"`
}

// NovoDeTickets cria um Progresso inicial com todos os Tickets em PENDING.
func NovoDeTickets(iniciativa, plano string, entradas []EntradaTicket, agora string) Progresso {
	if iniciativa == "" && len(entradas) > 0 {
		iniciativa = entradas[0].InitiativeID
	}
	itens := make([]ItemTicket, 0, len(entradas))
	for _, e := range entradas {
		itens = append(itens, ItemTicket{
			ID:           strings.TrimSpace(e.ID),
			Epico:        strings.TrimSpace(e.EpicID),
			Titulo:       strings.TrimSpace(e.Title),
			Estado:       EstadoPendente,
			Comando:      strings.TrimSpace(e.ValidationCmd),
			ExitCode:     -1,
			Alvos:        e.TargetFiles,
			Criterios:    e.SpecCriteria,
			AtualizadoEm: agora,
		})
	}
	p := Progresso{
		Schema:       SchemaVersao,
		Iniciativa:   strings.TrimSpace(iniciativa),
		Plano:        filepath.ToSlash(strings.TrimSpace(plano)),
		AtualizadoEm: agora,
		Tickets:      itens,
	}
	p.RecalcularResumo()
	return p
}

// RecalcularResumo atualiza os contadores a partir da lista atual de Tickets.
func (p *Progresso) RecalcularResumo() {
	r := Resumo{Total: len(p.Tickets)}
	for _, it := range p.Tickets {
		switch it.Estado {
		case EstadoAprovado:
			r.Aprovados++
		case EstadoEmAndamento:
			r.EmAndamento++
		case EstadoReprovado:
			r.Reprovados++
		case EstadoBloqueado:
			r.Bloqueados++
		default:
			r.Pendentes++
		}
	}
	p.Resumo = r
}

// Validar verifica o contrato progresso.v1 e unicidade de IDs.
func (p Progresso) Validar() error {
	if p.Schema != SchemaVersao {
		return fmt.Errorf("schema invalido: %q (esperado %q)", p.Schema, SchemaVersao)
	}
	if strings.TrimSpace(p.Iniciativa) == "" {
		return errors.New("iniciativa obrigatoria")
	}
	vistos := make(map[string]bool, len(p.Tickets))
	for _, it := range p.Tickets {
		if err := validarItem(it, vistos); err != nil {
			return err
		}
	}
	return nil
}

func validarItem(it ItemTicket, vistos map[string]bool) error {
	if strings.TrimSpace(it.ID) == "" {
		return errors.New("ticket sem id")
	}
	if vistos[it.ID] {
		return fmt.Errorf("ticket duplicado: %s", it.ID)
	}
	vistos[it.ID] = true
	if !estadoValido(it.Estado) {
		return fmt.Errorf("estado invalido em %s: %q", it.ID, it.Estado)
	}
	return nil
}

func estadoValido(e string) bool {
	switch e {
	case EstadoPendente, EstadoEmAndamento, EstadoAprovado, EstadoReprovado, EstadoBloqueado:
		return true
	default:
		return false
	}
}

// CarregarEntradasTickets lê um arquivo *_TICKETS_*.json.
func CarregarEntradasTickets(caminho string) ([]EntradaTicket, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, fmt.Errorf("ler tickets %s: %w", caminho, err)
	}
	var lista []EntradaTicket
	if err := json.Unmarshal(dados, &lista); err != nil {
		return nil, fmt.Errorf("decodificar tickets %s: %w", caminho, err)
	}
	return lista, nil
}

// Carregar lê e valida um arquivo progresso.v1 em disco.
func Carregar(caminho string) (Progresso, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return Progresso{}, fmt.Errorf("ler progresso %s: %w", caminho, err)
	}
	var p Progresso
	if err := json.Unmarshal(dados, &p); err != nil {
		return Progresso{}, fmt.Errorf("decodificar progresso %s: %w", caminho, err)
	}
	p.RecalcularResumo()
	return p, nil
}

// Salvar grava atomicamente o JSON canônico e, quando informado, o espelho Markdown.
func Salvar(caminhoJSON, caminhoMD string, p Progresso) error {
	p.RecalcularResumo()
	if err := p.Validar(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(caminhoJSON), 0o755); err != nil {
		return fmt.Errorf("criar pasta de progresso: %w", err)
	}
	dados, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("serializar progresso: %w", err)
	}
	dados = append(dados, '\n')
	if err := escrita.GravarAtomico(caminhoJSON, dados); err != nil {
		return err
	}
	if strings.TrimSpace(caminhoMD) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(caminhoMD), 0o755); err != nil {
		return fmt.Errorf("criar pasta do espelho Markdown: %w", err)
	}
	return escrita.GravarAtomico(caminhoMD, []byte(RenderizarMarkdown(p)))
}

// CaminhoEspelhoMD deriva o caminho do espelho Markdown a partir do JSON.
func CaminhoEspelhoMD(caminhoJSON string) string {
	ext := filepath.Ext(caminhoJSON)
	if strings.EqualFold(ext, ".json") {
		return strings.TrimSuffix(caminhoJSON, ext) + ".md"
	}
	return caminhoJSON + ".md"
}
