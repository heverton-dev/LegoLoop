package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"legoloop/internal/taxonomy"
)

var (
	marcaRe  = regexp.MustCompile(`(?i)legoloop`)
	naoLetra = regexp.MustCompile(`[^\p{L}\n]+`)
	palavra  = regexp.MustCompile(`[\p{L}]+`)
)

// termosProibidos devolve "linha:termo" para cada termo proibido do dicionário canônico.
// A marca LegoLoop não conta. Pontuação, hífen e barra separam palavras, como no taxlinter.
func termosProibidos(texto string) []string {
	limpo := naoLetra.ReplaceAllString(marcaRe.ReplaceAllString(texto, " "), " ")
	var achados []string
	vs, err := taxonomy.ValidateText(limpo)
	if err != nil {
		return []string{"erro no dicionário: " + err.Error()}
	}
	for _, v := range vs {
		achados = append(achados, fmt.Sprintf("%d:%s", v.Line, v.ForbiddenTerm))
	}
	for i, linha := range strings.Split(limpo, "\n") {
		for _, w := range palavra.FindAllString(linha, -1) {
			if achou := pluralProibido(w); achou != "" {
				achados = append(achados, fmt.Sprintf("%d:%s", i+1, achou))
			}
		}
	}
	return achados
}

func pluralProibido(w string) string {
	lw := strings.ToLower(w)
	if len(lw) > 2 && strings.HasSuffix(lw, "s") {
		if _, ok := taxonomy.CheckForbidden(lw[:len(lw)-1]); ok {
			return lw
		}
	}
	return ""
}

// semContratoDoFrontmatter esvazia o name e os itens de commands do frontmatter e mantém as linhas.
func semContratoDoFrontmatter(texto string) string {
	linhas := strings.Split(texto, "\n")
	if len(linhas) == 0 || strings.TrimSpace(linhas[0]) != "---" {
		return texto
	}
	emCommands := false
	for i := 1; i < len(linhas); i++ {
		l := linhas[i]
		if strings.TrimSpace(l) == "---" {
			break
		}
		switch {
		case strings.HasPrefix(l, "name:"):
			linhas[i], emCommands = "", false
		case strings.HasPrefix(l, "commands:"):
			linhas[i], emCommands = "", true
		case emCommands && strings.HasPrefix(l, " "):
			linhas[i] = ""
		default:
			emCommands = false
		}
	}
	return strings.Join(linhas, "\n")
}

func TestDetectorDeTermoProibido(t *testing.T) {
	for _, ruim := range []string{"dois lo" + "ops"} {
		if len(termosProibidos(ruim)) == 0 {
			t.Errorf("%q deveria ter termo proibido", ruim)
		}
	}
	for _, bom := range []string{"LegoLoop", "legoloop-fabricar", "cada Etapa do Processo", "ênfase vocal"} {
		if got := termosProibidos(bom); len(got) != 0 {
			t.Errorf("%q não deveria ter termo proibido: %v", bom, got)
		}
	}
}

func TestSkillsSemTermoProibido(t *testing.T) {
	raiz := filepath.Join("..", "..", "skills")
	arquivos, err := filepath.Glob(filepath.Join(raiz, "*", "SKILL.md"))
	if err != nil || len(arquivos) == 0 {
		t.Fatalf("nenhum SKILL.md em %s: %v", raiz, err)
	}
	arquivos = append(arquivos, filepath.Join(raiz, "embed.go"))
	for _, arq := range arquivos {
		data, err := os.ReadFile(arq)
		if err != nil {
			t.Fatal(err)
		}
		if achados := termosProibidos(semContratoDoFrontmatter(string(data))); len(achados) > 0 {
			t.Errorf("%s tem %d termo(s) proibido(s) (linha:termo): %v", arq, len(achados), achados)
		}
	}
}

// O name e os commands do frontmatter são contrato (identificam a skill): não entram na checagem.
func TestSemContratoDoFrontmatter(t *testing.T) {
	const contrato = "lo" + "op-map"
	entrada := "---\nname: legoloop-" + contrato + "\ndescription: texto ruim\ncommands:\n  - \"/" + contrato + "\"\n---\nCorpo\n"
	got := semContratoDoFrontmatter(entrada)
	if strings.Contains(got, contrato) {
		t.Errorf("contrato do frontmatter deveria sair: %q", got)
	}
	if !strings.Contains(got, "description: texto ruim") || !strings.Contains(got, "Corpo") {
		t.Errorf("descrição e corpo precisam ficar: %q", got)
	}
	if strings.Count(got, "\n") != strings.Count(entrada, "\n") {
		t.Errorf("o número de linhas precisa se manter")
	}
}
