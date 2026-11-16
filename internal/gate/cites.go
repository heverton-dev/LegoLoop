package gate

import (
	"bytes"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// maxCitacoes limita o trabalho por artefato: acima disso o gate reprova em vez de varrer sem fim.
const maxCitacoes = 500

// bt e o acento grave, que nao cabe dentro de string bruta.
const bt = "`"

var (
	citePathRe = regexp.MustCompile(`(?:[A-Za-z0-9_.-]+/)*[A-Za-z0-9_.-]+\.(?:go|md|json|sh|py|js|ts|tsx|css|html|ya?ml|toml|mod|sql)`)
	// citeColonRe: o numero vem logo apos o caminho, como "arquivo.go:12" ou "arquivo.go:12-15".
	citeColonRe = regexp.MustCompile(`^[` + bt + `')\]]*:(\d+)(?:\s*[-–]\s*(\d+))?`)
	// citeLooseRe: "(linha 12)", "linhas 12-15" ou "line 12" ate 80 caracteres depois do caminho, antes do proximo caminho.
	citeLooseRe = regexp.MustCompile(`(?i)^.{0,80}?(?:linhas?|lines?)\s*(\d+)(?:\s*[-–]\s*(\d+))?`)
)

type citation struct {
	file       string
	from, to   int
	artifactLn int
}

// extractCitations acha, linha a linha do texto, cada "caminho + número de linha". O número pertence ao caminho
// mais próximo à sua esquerda na mesma linha. Caminho sem número não é citação.
func extractCitations(text string) []citation {
	var out []citation
	for ln, line := range strings.Split(text, "\n") {
		idx := citePathRe.FindAllStringIndex(line, -1)
		for i, loc := range idx {
			rest := line[loc[1]:]
			if i+1 < len(idx) {
				rest = line[loc[1]:idx[i+1][0]]
			}
			m := citeColonRe.FindStringSubmatch(rest)
			if m == nil {
				m = citeLooseRe.FindStringSubmatch(rest)
			}
			if m == nil {
				continue
			}
			from, _ := strconv.Atoi(m[1])
			to := from
			if m[2] != "" {
				to, _ = strconv.Atoi(m[2])
			}
			out = append(out, citation{file: line[loc[0]:loc[1]], from: from, to: to, artifactLn: ln + 1})
		}
	}
	return out
}

// citesExist e o builtin cites_exist(artefato, base): verdadeiro só se o artefato traz ao menos uma citação
// "arquivo:linha" e todas existem: o arquivo e um arquivo regular sob base e o intervalo citado cabe nele.
// Citação fora de base ou com ".." conta como inexistente. Prova a existência, nao o sentido da linha citada.
func citesExist(env Env, args []value, _ *regexp.Regexp) (value, error) {
	data, err := readGateFile(env, args[0].s)
	if err != nil {
		return value{}, err
	}
	cites := extractCitations(string(data))
	if len(cites) == 0 || len(cites) > maxCitacoes {
		return value{kind: kBool}, nil
	}
	base := path.Clean(filepath.ToSlash(expandGateVars(env, args[1].s)))
	cache := map[string]int{}
	for _, c := range cites {
		target := path.Join(base, c.file)
		if !filepath.IsLocal(filepath.FromSlash(target)) || c.from < 1 || c.to < c.from {
			return value{kind: kBool}, nil
		}
		total, ok := cache[target]
		if !ok {
			content, err := readGateFile(Env{Root: env.Root}, target)
			if err != nil {
				total = -1 // inexistente, diretorio ou grande demais: citacao invalida
			} else {
				total = bytes.Count(content, []byte("\n"))
				if len(content) > 0 && content[len(content)-1] != '\n' {
					total++
				}
			}
			cache[target] = total
		}
		if total < c.to {
			return value{kind: kBool}, nil
		}
	}
	return value{kind: kBool, b: true}, nil
}
