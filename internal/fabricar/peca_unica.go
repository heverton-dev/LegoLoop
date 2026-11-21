package fabricar

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"legoloop/internal/caminhos"
	"legoloop/internal/processo"
)

// TiposPeca são os tipos que o Fabricar cria como peça única (nó solto).
var TiposPeca = []string{"script", "skill", "gate", "hook", "mcp", "processo"}

// TipoPecaValido diz se o tipo é um dos aceitos em --peca.
func TipoPecaValido(tipo string) bool { return slices.Contains(TiposPeca, tipo) }

// ferramentasDoWorker são as ferramentas que o worker efêmero da peça pode usar (Lei #4, condição 2).
var ferramentasDoWorker = map[string][]string{
	"script": {"Read", "Write", "Edit", "Bash"}, "gate": {"Read", "Write", "Edit", "Bash"},
	"hook": {"Read", "Write", "Edit", "Bash"}, "mcp": {"Read", "Write", "Edit", "Bash"},
	"skill": {"Read", "Write", "Edit"}, "processo": {"Read", "Write", "Edit"},
}

// SwarmPlan é o contrato declarativo do worker efêmero que lapida a peça. Nada roda ao gravá-lo:
// o operador ativa o worker com legoloop swarm add.
type SwarmPlan struct {
	Peca          string   `json:"peca"`
	Tipo          string   `json:"tipo"`
	Workers       int      `json:"workers"`
	Ticket        string   `json:"ticket"`
	Ferramentas   []string `json:"ferramentas"`
	WriteBoundary string   `json:"write_boundary"`
	Subagentes    string   `json:"subagentes"`
}

func swarmPlanDe(nome, tipo string) ([]byte, error) {
	sp := SwarmPlan{Peca: nome, Tipo: tipo, Workers: 1, Ticket: "T1-" + strings.TrimPrefix(nome, namePrefix),
		Ferramentas: ferramentasDoWorker[tipo], WriteBoundary: "strict", Subagentes: "ephemeral"}
	return json.MarshalIndent(sp, "", "  ")
}

// ProcessoPeca monta o Processo de uma peça única: uma etapa do tipo da peça, ligada à própria peça.
func ProcessoPeca(nome, objetivo, tipo string, now time.Time) (processo.Processo, error) {
	art := ArtefatoDe(tipo, naturezaQualquer, 1)
	g, err := processo.GateParaArtefato(art)
	if err != nil {
		return processo.Processo{}, err
	}
	e := processo.Etapa{ID: "e01", Titulo: tituloDe(tipo, 1), Tipo: tipo, Peca: nome, Produz: art, Gate: g}
	return processo.Processo{Nome: nome, Versao: "1.0.0", Natureza: naturezaQualquer, Objetivo: objetivo,
		CriadoEm: now.UTC(), Etapas: []processo.Etapa{e}}, nil
}

// ajustarPeca tira a TRIAGEM e põe o swarm.plan.json: a triagem roteia pedidos entre vários nós e uma peça só tem um.
func ajustarPeca(files []PlannedFile, roots Roots, p processo.Processo, peca string) ([]PlannedFile, error) {
	if peca == "" {
		return files, nil
	}
	sp, err := swarmPlanDe(p.Nome, peca)
	if err != nil {
		return nil, err
	}
	out := slices.DeleteFunc(files, func(f PlannedFile) bool { return f.Kind == "triagem" })
	path := filepath.Join(roots.Vault, caminhos.DirProcessos, p.Nome, "swarm.plan.json")
	return append(out, PlannedFile{Kind: "swarm_plan", Path: path, Data: sp}), nil
}

func errTipoPeca(tipo string) error {
	return fmt.Errorf("%w: tipo de peça %q (use %s)", ErrInvalidRequest, tipo, strings.Join(TiposPeca, ", "))
}
