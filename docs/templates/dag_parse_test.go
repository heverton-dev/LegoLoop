package templates_test

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
)

var placeholderRe = regexp.MustCompile(`\{\{[^}]+\}\}`)

// runtimeShellRe acha "runtime" seguido de powershell na mesma linha (JSON ou Markdown).
var runtimeShellRe = regexp.MustCompile(`(?i)runtime[^
]{0,12}powershell`)

// concretos troca cada {{PLACEHOLDER}} por um valor concreto, único por ocorrência, para o parser real avaliar a forma do documento.
func concretos(data []byte) []byte {
	n := 0
	return placeholderRe.ReplaceAllFunc(data, func(_ []byte) []byte {
		n++
		return []byte("exemplo" + string(rune('a'+n%26)) + string(rune('a'+n/26)))
	})
}

// Todo template em forma de DAG passa no parser real (dag.Parse), o mesmo de legoloop validate.
func TestTemplatesEmFormaDeDAGPassamNoParser(t *testing.T) {
	for _, nome := range []string{"TEMPLATE_PIECE_PROCESSO.json", arquivoProcesso, "TEMPLATE_DAG.json", "TEMPLATE_DAG_VIDEOREA.json", "TEMPLATE_DAG_LEGOPRESS.json", "TEMPLATE_PIECE_UNIT.json"} {
		t.Run(nome, func(t *testing.T) {
			if _, err := dag.Parse(concretos(read(t, nome))); err != nil {
				t.Fatalf("%s não passa em dag.Parse: %v", nome, err)
			}
		})
	}
}

// Nenhum template em forma de DAG fixa um shell de um só sistema operacional.
func TestTemplatesEmFormaDeDAGSemShellDeUmSO(t *testing.T) {
	for _, nome := range []string{"TEMPLATE_PIECE_PROCESSO.json", arquivoProcesso, "TEMPLATE_DAG.json", "TEMPLATE_DAG_VIDEOREA.json", "TEMPLATE_DAG_LEGOPRESS.json", "TEMPLATE_PIECE_UNIT.json"} {
		texto := strings.ToLower(string(read(t, nome)))
		for _, shell := range []string{`"powershell"`, `"pwsh"`, `"cmd"`, `"bash"`, `"sh"`} {
			if strings.Contains(texto, shell) {
				t.Errorf("%s fixa o shell %s: o template deve rodar em qualquer SO", nome, shell)
			}
		}
	}
}

// Nenhum template JSON ou Markdown fixa um shell como runtime: tooling.runtime declara o motor (legoloop).
func TestTemplatesSemRuntimeDeUmSO(t *testing.T) {
	for _, padrao := range []string{"*.json", "*.md"} {
		arquivos, err := filepath.Glob(padrao)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range arquivos {
			if runtimeShellRe.Match(read(t, a)) {
				t.Errorf("%s fixa powershell no runtime: use o runtime legoloop", a)
			}
		}
	}
}
