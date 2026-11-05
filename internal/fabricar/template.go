package fabricar

import (
	"fmt"
	"strings"

	"legoloop/internal/processo"
)

// templateFileName é o nome do gabarito do processo em docs/templates.
func templateFileName(name string) string {
	return "TEMPLATE_" + upperSnake(strings.TrimPrefix(name, namePrefix)) + ".md"
}

// EmitTemplate é a Camada 5: o gabarito contratual das saídas do processo. Cada etapa vira uma seção com o
// artefato que ela produz, quem a executa e o gate que a aprova (os mesmos do Processo, do DAG e do PLAN).
// Os {{PLACEHOLDER}} são intencionais: é um template.
func EmitTemplate(p processo.Processo) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 📐 TEMPLATE: %s\n\n", p.Nome)
	fmt.Fprintf(&b, "> Gabarito contratual das saídas do processo `%s`.  \n", p.Nome)
	b.WriteString("> Preencha cada `{{PLACEHOLDER}}`. Gate de fechamento: `!has_pattern('<arquivo>', '\\{\\{')`.\n\n---\n\n")
	for i, e := range p.Etapas {
		fmt.Fprintf(&b, "## %d. Etapa `%s`: `%s` (%s)\n\n* **Gate:** `%s`\n\n{{CONTEUDO_ETAPA_%d}}\n\n",
			i+1, e.ID, e.Produz.Caminho, quemExecuta(e), e.Gate, i+1)
	}
	return b.String()
}

// quemExecuta diz o que executa a etapa: a peça ou o processo composto.
func quemExecuta(e processo.Etapa) string {
	if e.Compoe != "" {
		return "processo `" + e.Compoe + "`"
	}
	return "skill `" + e.Peca + "`"
}
