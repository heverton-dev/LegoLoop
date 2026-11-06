package ste

import (
	"os"
	"path/filepath"
)

// HarnessInfo consolida metadados de detecção de um harness no sistema de arquivos.
type HarnessInfo struct {
	ID            HarnessID `json:"id"`
	Nome          string    `json:"nome"`
	Instalado     bool      `json:"instalado"`
	PastaBase     string    `json:"pasta_base"`
	CaminhoRegras string    `json:"caminho_regras"`
	Ativo         bool      `json:"ativo"`
}

type definicaoHarness struct {
	id       HarnessID
	nome     string
	baseRel  string
	regraRel string
}

func catalogoDefinicoes() []definicaoHarness {
	return []definicaoHarness{
		{
			id:       HarnessGemini,
			nome:     "Google Gemini / Antigravity",
			baseRel:  ".gemini",
			regraRel: filepath.Join(".gemini", "config", "rules", "stept-rules.md"),
		},
		{
			id:       HarnessClaude,
			nome:     "Claude Code",
			baseRel:  ".claude",
			regraRel: filepath.Join(".claude", "rules", "stept-rules.md"),
		},
		{
			id:       HarnessOMP,
			nome:     "OMP CLI",
			baseRel:  ".omp",
			regraRel: filepath.Join(".omp", "rules", "stept-rules.md"),
		},
		{
			id:       HarnessOpenCode,
			nome:     "OpenCode",
			baseRel:  filepath.Join(".config", "opencode"),
			regraRel: filepath.Join(".config", "opencode", "rules", "stept-rules.md"),
		},
		{
			id:       HarnessMimo,
			nome:     "MimoCode",
			baseRel:  ".mimocode",
			regraRel: filepath.Join(".mimocode", "rules", "stept-rules.md"),
		},
	}
}

// HomeDirPadrao obtém o diretório do usuário ou ponto de montagem atual.
func HomeDirPadrao() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func avaliarHarness(def definicaoHarness, homeDir string) HarnessInfo {
	pastaBase := filepath.Join(homeDir, def.baseRel)
	caminhoRegra := filepath.Join(homeDir, def.regraRel)

	infoBase, errBase := os.Stat(pastaBase)
	instalado := errBase == nil && infoBase.IsDir()

	infoRegra, errRegra := os.Stat(caminhoRegra)
	ativo := errRegra == nil && !infoRegra.IsDir()

	return HarnessInfo{
		ID:            def.id,
		Nome:          def.nome,
		Instalado:     instalado,
		PastaBase:     pastaBase,
		CaminhoRegras: caminhoRegra,
		Ativo:         ativo,
	}
}

// ScanHarnesses detecta a presença física dos harnesses sem criar diretórios fantasmas.
func ScanHarnesses(homeDir string) []HarnessInfo {
	if homeDir == "" {
		homeDir = HomeDirPadrao()
	}
	defs := catalogoDefinicoes()
	res := make([]HarnessInfo, 0, len(defs))
	for _, def := range defs {
		res = append(res, avaliarHarness(def, homeDir))
	}
	return res
}
