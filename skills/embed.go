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

// naoNativas lista as skills legoloop-* embutidas que não são funções nativas.
var naoNativas = map[string]bool{
	"legoloop-canal-redc":     true,
	"legoloop-cartoon-shorts": true,
	"legoloop-n10-foundry":    true,
	"legoloop-video-edition":  true,
}

// FuncoesNativas lista, em ordem alfabética, os nomes das 14 funções nativas do LegoLoop.
// A lista vem das pastas embutidas em FS, menos as de naoNativas.
var FuncoesNativas = listarNativas()

func listarNativas() []string {
	entradas, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil
	}
	var nomes []string
	for _, e := range entradas {
		if e.IsDir() && !naoNativas[e.Name()] {
			nomes = append(nomes, e.Name())
		}
	}
	return nomes
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
