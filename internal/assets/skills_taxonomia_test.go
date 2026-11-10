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

// pluralProibido pega o plural simples ("loops") de um termo proibido.
func pluralProibido(w string) string {
	lw := strings.ToLower(w)
	if len(lw) > 2 && strings.HasSuffix(lw, "s") {
		if _, ok := taxonomy.CheckForbidden(lw[:len(lw)-1]); ok {
			return lw
		}
	}
	return ""
}

func TestDetectorDeTermoProibido(t *testing.T) {
	for _, ruim := range []string{"um loop", "dois Pipelines", "legoloop-loop-map", "cada fase", "a unit_task"} {
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
		if achados := termosProibidos(string(data)); len(achados) > 0 {
			t.Errorf("%s tem %d termo(s) proibido(s) (linha:termo): %v", arq, len(achados), achados)
		}
	}
}
