package processo

import (
	"os"
	"path/filepath"
	"strings"
)

// IndiceSkills é a única skill que o modo lazy deixa nas pastas dos harnesses.
// Ela manda buscar as demais por skill_search e ler por call_skill.
const IndiceSkills = "legoloop-skills"

// Politica decide quais skills vão para as pastas dos harnesses.
// Prova de 20/11: o agy lê o frontmatter de toda skill achada e cobra até ~20,5k tokens por turno.
type Politica struct {
	Eager    bool            // true = comportamento antigo: toda skill vai para toda pasta
	Excecoes map[string]bool // skills que ficam sempre visíveis mesmo em modo lazy
}

// Vincula diz se a skill vai para a pasta dos harnesses.
func (p Politica) Vincula(nome string) bool {
	return p.Eager || nome == IndiceSkills || p.Excecoes[nome]
}

// LerExcecoes lê ~/.legoloop/skills-eager.txt: um nome por linha; linha com # é comentário.
func LerExcecoes(home string) map[string]bool {
	out := map[string]bool{}
	if home == "" {
		return out
	}
	raw, err := os.ReadFile(filepath.Join(home, ".legoloop", "skills-eager.txt"))
	if err != nil {
		return out
	}
	for _, ln := range strings.Split(string(raw), "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" && !strings.HasPrefix(ln, "#") {
			out[ln] = true
		}
	}
	return out
}
