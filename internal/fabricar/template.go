package fabricar

import (
	"fmt"
	"strings"
)

// templateFileName é o nome do gabarito do processo em docs/templates.
func templateFileName(name string) string {
	return "TEMPLATE_" + upperSnake(strings.TrimPrefix(name, namePrefix)) + ".md"
}

// EmitTemplate é a Camada 5: o gabarito contratual das saídas do processo. Cada passo vira uma seção com o
// artefato esperado, o gate que o aprova e um placeholder. Os {{PLACEHOLDER}} são intencionais: é um template.
func EmitTemplate(ps ProcessSpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 📐 TEMPLATE: %s\n\n", ps.Name)
	fmt.Fprintf(&b, "> Gabarito contratual das saídas do processo `%s`.  \n", ps.Name)
	b.WriteString("> Preencha cada `{{PLACEHOLDER}}`. Gate de fechamento: `!has_pattern('<arquivo>', '\\{\\{')`.\n\n---\n\n")
	for i, s := range ps.Steps {
		fmt.Fprintf(&b, "## %d. `%s` (skill `%s`)\n\n* **Gate:** `%s`\n\n{{CONTEUDO_PASSO_%d}}\n\n", i+1, ps.artifactOf(i), refName(s.Piece), ps.gateOf(i), i+1)
	}
	return b.String()
}
