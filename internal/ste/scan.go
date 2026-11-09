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

func defHarness(id HarnessID, nome, base string) definicaoHarness {
	return definicaoHarness{
		id:       id,
		nome:     nome,
		baseRel:  base,
		regraRel: filepath.Join(base, "rules", "stept-rules.md"),
	}
}

func catalogoDefinicoes() []definicaoHarness {
	gemini := definicaoHarness{
		id:       HarnessGemini,
		nome:     "Google Gemini / Antigravity",
		baseRel:  ".gemini",
		regraRel: filepath.Join(".gemini", "config", "rules", "stept-rules.md"),
	}
	opencode := defHarness(HarnessOpenCode, "OpenCode", filepath.Join(".config", "opencode"))
	return []definicaoHarness{
		gemini,
		defHarness(HarnessClaude, "Claude Code", ".claude"),
		defHarness(HarnessOMP, "OMP CLI", ".omp"),
		opencode,
		defHarness(HarnessMimo, "MimoCode", ".mimocode"),
		defHarness(HarnessCursor, "Cursor IDE", ".cursor"),
		defHarness(HarnessWindsurf, "Windsurf IDE", ".windsurf"),
		defHarness(HarnessVSCode, "VS Code", ".vscode"),
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
