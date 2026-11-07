// Package skills embute os arquivos SKILL.md das funções nativas do LegoLoop.
package skills

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

// FS embute todos os SKILL.md de legoloop-*.
//
//go:embed legoloop-*/*.md
var FS embed.FS

// FuncoesNativas lista os nomes canônicos das 10 funções nativas do LegoLoop.
var FuncoesNativas = []string{
	"legoloop-fabricar",
	"legoloop-matriz",
	"legoloop-listar",
	"legoloop-refinar",
	"legoloop-remover",
	"legoloop-compor",
	"legoloop-legopress",
	"legoloop-piece-upgrader",
	"legoloop-mega-pipeline",
	"legoloop-loop-map",
}

// ObterSkill devolve o conteúdo do SKILL.md da função informada (ex: "legoloop-matriz").
func ObterSkill(nome string) ([]byte, error) {
	relPath := fmt.Sprintf("%s/SKILL.md", nome)
	data, err := fs.ReadFile(FS, relPath)
	if err != nil {
		return nil, fmt.Errorf("skill nativa não encontrada para %q: %w", nome, err)
	}
	return data, nil
}

// TodasAsSkills devolve um mapa nome -> conteúdo para todas as funções nativas.
func TodasAsSkills() (map[string][]byte, error) {
	out := make(map[string][]byte, len(FuncoesNativas))
	for _, fn := range FuncoesNativas {
		data, err := ObterSkill(fn)
		if err != nil {
			return nil, err
		}
		out[fn] = data
	}
	return out, nil
}

// CaminhoNoCofre devolve o caminho relativo padrão no cofre (ex: "skills/legoloop-matriz/SKILL.md").
func CaminhoNoCofre(nome string) string {
	return filepath.Join("skills", nome, "SKILL.md")
}
