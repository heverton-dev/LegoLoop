// Package skills: custo.go mede quanto contexto as skills eager de um harness custam por turno.
package skills

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// bytesPorToken é a razão usada na medição; reproduziu o print do agy com 4% de erro.
const bytesPorToken = 4

var reFrontmatter = regexp.MustCompile(`(?s)^---\r?\n(.*?)\r?\n---`)

// CustoPasta é o resultado de uma pasta de skills lida sozinha.
type CustoPasta struct {
	Pasta     string `json:"pasta"`
	Skills    int    `json:"skills"`
	Invalidas int    `json:"invalidas"`
	Tokens    int    `json:"tokens"`
}

// CustoHarness é a união das pastas de um harness, deduplicada por nome de pasta.
type CustoHarness struct {
	Harness        string       `json:"harness"`
	Skills         int          `json:"skills"`
	Tokens         int          `json:"tokens"`
	Invalidas      int          `json:"invalidas"`
	Duplicadas     int          `json:"duplicadas"`
	Disabled       int          `json:"disabled"`
	TokensDisabled int          `json:"tokens_disabled"`
	Pastas         []CustoPasta `json:"pastas"`
}

type lidaSkill struct {
	nome   string
	tokens int
	valida bool
}

// CustoGrupo mede as pastas como o agy as lê: a primeira ocorrência de cada nome vale, as demais são duplicadas.
// Entradas *.disabled são contadas à parte e não entram em Skills nem em Tokens.
func CustoGrupo(harness string, pastas []string) CustoHarness {
	c := CustoHarness{Harness: harness}
	vistos := map[string]bool{}
	for _, p := range pastas {
		cp := CustoPasta{Pasta: p}
		for _, s := range lerPasta(p) {
			if !s.valida {
				cp.Invalidas++
				c.Invalidas++
				continue
			}
			cp.Skills++
			cp.Tokens += s.tokens
			somarUnica(&c, vistos, s)
		}
		c.Pastas = append(c.Pastas, cp)
	}
	return c
}

func somarUnica(c *CustoHarness, vistos map[string]bool, s lidaSkill) {
	if vistos[s.nome] {
		c.Duplicadas++
		return
	}
	vistos[s.nome] = true
	if strings.HasSuffix(s.nome, ".disabled") {
		// provado em 20/11 no /context do agy: entradas .disabled não pesam; ficam fora do total.
		c.Disabled++
		c.TokensDisabled += s.tokens
		return
	}
	c.Skills++
	c.Tokens += s.tokens
}

// lerPasta lê cada SKILL.md de primeiro nível. Pasta ausente devolve lista vazia.
func lerPasta(pasta string) []lidaSkill {
	ents, err := os.ReadDir(pasta)
	if err != nil {
		return nil
	}
	var out []lidaSkill
	for _, e := range ents {
		arq := filepath.Join(pasta, e.Name(), "SKILL.md")
		raw, err := os.ReadFile(arq)
		if err != nil {
			continue
		}
		out = append(out, medirSkill(e.Name(), arq, string(raw)))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].nome < out[j].nome })
	return out
}

// medirSkill conta frontmatter + caminho, o que o harness injeta por skill.
func medirSkill(nome, caminho, texto string) lidaSkill {
	m := reFrontmatter.FindStringSubmatch(strings.TrimPrefix(texto, string(rune(0xFEFF))))
	if m == nil {
		return lidaSkill{nome: nome}
	}
	return lidaSkill{nome: nome, valida: true, tokens: (len(m[1]) + len(caminho)) / bytesPorToken}
}
