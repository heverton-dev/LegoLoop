package fabricar

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/processo"
)

// limpo deixa o texto livre numa linha só e sem crase: título não abre seção no documento.
func limpo(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "`", "'")
}

// interviewDe renderiza o registro da entrevista já feita (o espelho canônico), para que o Sexteto
// tenha a INTERVIEW mesmo quando o pedido chega autônomo. Caminho ISO em docs/interviews.
func interviewDe(name string, clar Clarification, now time.Time) namedFile {
	date := now.Format("2006-01-02")
	slug := strings.TrimPrefix(name, namePrefix)
	file := fmt.Sprintf("%s_INTERVIEW_%s.md", date, upperSnake(slug))
	var b strings.Builder
	fmt.Fprintf(&b, "# 🎤 INTERVIEW: %s\n\n> Registro da entrevista que deu origem ao processo `%s`.\n\n", clar.Theme, name)
	fmt.Fprintf(&b, "* **Tema:** %s\n* **Confiança:** %d%%\n* **Turnos:** %d\n* **Alcançou a meta:** %t\n\n",
		clar.Theme, clar.Confidence, clar.Turns, clar.Reached)
	m := clar.Mirror
	fmt.Fprintf(&b, "## Espelho Canônico\n\n* **Resultado Esperado:** %s\n* **Quem se Beneficia:** %s\n* **Por que Agora:** %s\n* **Critério de Sucesso:** %s\n",
		limpo(m.Outcome), limpo(m.User), limpo(m.WhyNow), limpo(m.Success))
	for _, c := range m.Constraints {
		fmt.Fprintf(&b, "* **Restrição:** %s\n", limpo(c))
	}
	for _, o := range m.OutOfScope {
		fmt.Fprintf(&b, "* **Fora de Escopo:** %s\n", limpo(o))
	}
	return namedFile{Name: file, Data: []byte(b.String())}
}

// reportDe renderiza o relatório do processo em docs/reports com caminho ISO: cada etapa com
// natureza, artefato e gate, para qualquer natureza de processo.
func reportDe(name string, p processo.Processo, now time.Time) namedFile {
	date := now.Format("2006-01-02")
	slug := strings.TrimPrefix(name, namePrefix)
	file := fmt.Sprintf("%s_REPORT_%s.md", date, upperSnake(slug))
	var b strings.Builder
	fmt.Fprintf(&b, "# 📊 REPORT: %s\n\n> Relatório do processo `%s` (natureza `%s`).\n\n", limpo(p.Objetivo), name, p.Natureza)
	b.WriteString("| Etapa | Natureza | Artefato | Gate |\n| :--- | :--- | :--- | :--- |\n")
	for _, e := range p.Etapas {
		fmt.Fprintf(&b, "| `%s` | %s | `%s` | `%s` |\n", e.ID, e.Natureza, e.Produz.Caminho, e.Gate)
	}
	return namedFile{Name: file, Data: []byte(b.String())}
}

// contratosDe carrega o manifesto contrato.json de cada peça usada na Fusão: o contrato mora ao
// lado do SKILL.md da peça. Peça sem contrato entra sem manifesto (nil); manifesto quebrado erra.
func contratosDe(m Merged) (map[string]*processo.Contrato, error) {
	out := map[string]*processo.Contrato{}
	for _, pc := range m.Used {
		caminho := filepath.Join(filepath.Dir(pc.Path), "contrato.json")
		if _, err := os.Stat(caminho); err != nil {
			continue
		}
		c, err := processo.ContratoDe(caminho)
		if err != nil {
			return nil, err
		}
		out[pc.Name] = &c
	}
	return out, nil
}
