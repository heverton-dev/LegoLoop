package forja

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/processo"
)

// DerivaItem descreve uma skill eager encontrada na pasta de um harness.
type DerivaItem struct {
	Harness string `json:"harness"`
	Skill   string `json:"skill"`
	Pasta   string `json:"pasta"`
}

// ResultadoDeriva resume a verificacao de deriva nos harnesses.
type ResultadoDeriva struct {
	Aprovado bool         `json:"aprovado"`
	Derivas  []DerivaItem `json:"derivas,omitempty"`
}

// DetectarDeriva inspeciona as pastas de skills dos harnesses no home e
// devolve qualquer skill eager (diferente de legoloop-skills).
func DetectarDeriva(home string) ([]DerivaItem, error) {
	if strings.TrimSpace(home) == "" {
		return nil, ErrSemHome
	}
	harnesses := processo.DetectarHarnesses(processo.OSFS(), home)
	var derivas []DerivaItem
	for _, h := range harnesses {
		itens, err := detectarDerivaHarness(h)
		if err != nil {
			return nil, err
		}
		derivas = append(derivas, itens...)
	}
	sort.Slice(derivas, func(i, j int) bool {
		if derivas[i].Harness != derivas[j].Harness {
			return derivas[i].Harness < derivas[j].Harness
		}
		return derivas[i].Skill < derivas[j].Skill
	})
	return derivas, nil
}

func detectarDerivaHarness(h processo.Harness) ([]DerivaItem, error) {
	info, err := os.Stat(h.Skills)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("forja: inspecionar pasta de skills %s: %w", h.Skills, err)
	}
	if !info.IsDir() {
		return nil, nil
	}
	entradas, err := os.ReadDir(h.Skills)
	if err != nil {
		return nil, fmt.Errorf("forja: listar skills %s: %w", h.Skills, err)
	}
	return filtrarSkillsEager(h, entradas), nil
}

func filtrarSkillsEager(h processo.Harness, entradas []os.DirEntry) []DerivaItem {
	var itens []DerivaItem
	for _, e := range entradas {
		nome := e.Name()
		if nome == processo.IndiceSkills || strings.HasPrefix(nome, ".") {
			continue
		}
		pasta := filepath.Join(h.Skills, nome)
		if ehDiretorioExistente(pasta) {
			itens = append(itens, DerivaItem{
				Harness: h.Nome,
				Skill:   nome,
				Pasta:   pasta,
			})
		}
	}
	return itens
}

func ehDiretorioExistente(caminho string) bool {
	info, err := os.Stat(caminho)
	return err == nil && info.IsDir()
}

// FormatarDerivas gera texto legivel listando harness, skill e pasta.
func FormatarDerivas(derivas []DerivaItem) string {
	if len(derivas) == 0 {
		return "deriva aprovada: nenhuma skill eager encontrada"
	}
	linhas := make([]string, 0, len(derivas))
	for _, d := range derivas {
		linhas = append(linhas, fmt.Sprintf(
			"harness=%s skill=%s pasta=%s",
			d.Harness, d.Skill, filepath.ToSlash(d.Pasta),
		))
	}
	return strings.Join(linhas, "\n")
}
