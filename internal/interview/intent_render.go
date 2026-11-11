package interview

import (
	"fmt"
	"strings"

	"legoloop/docs/templates"
)

// IntentView reúne tudo que o TEMPLATE_INTENT precisa para ser preenchido a partir da entrevista.
type IntentView struct {
	Theme      string
	Date       string
	SessionID  string
	Confidence int
	Mirror     *Mirror
	SpecRel    string // caminho do próximo documento, relativo à raiz do projeto
	// Natureza do processo (vocabulário). Vazia ou "software" mantém as Peças e critérios de engenharia de
	// software; qualquer outra natureza usa Peças e critérios sem comandos de Go.
	Natureza string
}

const notInformed = "Não informado na entrevista"

// RenderIntent materializa o TEMPLATE_INTENT (mesmas seções, mesma ordem) com os dados do espelho socrático.
// Nada é inventado: o que a entrevista não respondeu é declarado como "a definir", nunca como placeholder.
func RenderIntent(v IntentView) string {
	var b strings.Builder
	writeIntentHeader(&b, v)
	writeIntentContext(&b, v)
	writeIntentScope(&b, v)
	writeIntentBlocks(&b, v.Natureza)
	writeIntentGovernance(&b, v)
	writeIntentCriteria(&b, v)
	writeIntentHandoff(&b, v)
	return b.String()
}

// intentStatus: a aprovação é sempre humana. Com confiança >= 95% a INTENT segue para revisão.
func intentStatus(confidence int) string {
	if confidence >= 95 {
		return "UNDER_REVIEW"
	}
	return "DRAFT"
}

func orDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

func bullets(items []string, indent, empty string) string {
	var b strings.Builder
	for _, it := range items {
		if it = strings.TrimSpace(it); it != "" {
			fmt.Fprintf(&b, "%s* %s\n", indent, it)
		}
	}
	if b.Len() == 0 {
		fmt.Fprintf(&b, "%s* %s\n", indent, empty)
	}
	return b.String()
}

func cell(s string) string { return strings.ReplaceAll(s, "|", `\|`) }

func writeIntentHeader(b *strings.Builder, v IntentView) {
	fmt.Fprintf(b, "# 🎯 INTENT: %s\n\n", v.Theme)
	fmt.Fprintf(b, "> **Status:** %s  \n", intentStatus(v.Confidence))
	fmt.Fprintf(b, "> **Data:** %s  \n", v.Date)
	b.WriteString("> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  \n")
	b.WriteString("> **Alvo:** A definir na SPEC  \n")
	fmt.Fprintf(b, "> **Sessão:** %s  \n", v.SessionID)
	fmt.Fprintf(b, "> **Confiança da Entrevista:** %d%% (mínimo de 95%% para sair de DRAFT)  \n", v.Confidence)
	fmt.Fprintf(b, "> **Template:** TEMPLATE_INTENT.md v%s  \n\n---\n\n", templates.Version)
}

func writeIntentContext(b *strings.Builder, v IntentView) {
	m := v.Mirror
	b.WriteString("## 1. Contexto e Motivação\n")
	fmt.Fprintf(b, "* **O Problema Real:** %s\n", v.Theme)
	fmt.Fprintf(b, "* **Quem Sofre:** %s\n", orDefault(m.User, notInformed))
	fmt.Fprintf(b, "* **O Porquê:** %s\n", orDefault(m.WhyNow, notInformed))
	b.WriteString("* **Estado Atual vs. Estado Desejado:**\n")
	fmt.Fprintf(b, "  * **Hoje:** Conforme a ideia inicial da entrevista: %s\n", v.Theme)
	fmt.Fprintf(b, "  * **Amanhã:** %s\n\n---\n\n", orDefault(m.Outcome, notInformed))
}

func writeIntentScope(b *strings.Builder, v IntentView) {
	m := v.Mirror
	b.WriteString("## 2. Escopo Inegociável\n\n### O que ESTÁ no Escopo (Objetivos Primários):\n")
	fmt.Fprintf(b, "* [ ] Objetivo 1: %s\n", orDefault(m.Outcome, notInformed))
	fmt.Fprintf(b, "* [ ] Objetivo 2: %s\n\n", orDefault(m.Success, notInformed))
	b.WriteString("### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):\n")
	b.WriteString("* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.\n")
	b.WriteString("* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.\n")
	b.WriteString("* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.\n")
	for _, o := range m.OutOfScope {
		if o = strings.TrimSpace(o); o != "" {
			fmt.Fprintf(b, "* ⛔ %s\n", o)
		}
	}
	b.WriteString("\n---\n\n")
}

func writeIntentBlocks(b *strings.Builder, natureza string) {
	b.WriteString("## 3. Mapeamento Preliminar de Blocos e Subloops\n")
	b.WriteString("Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.\n\n")
	b.WriteString("| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |\n| :--- | :--- | :--- | :--- |\n")
	b.WriteString("| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |\n")
	if !ehSoftware(natureza) {
		writePecasGenericas(b)
		return
	}
	b.WriteString("| `PECA-02` | `unit_task` | Código de produção | Implementação da lógica via TDD Red-Green. |\n")
	b.WriteString("| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |\n")
	b.WriteString("| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |\n\n---\n\n")
}

// ehSoftware: natureza vazia (INTENT antigo ou entrevista sem natureza) e "software" usam o texto de engenharia de software.
func ehSoftware(natureza string) bool { return natureza == "" || natureza == "software" }

// writePecasGenericas completa a tabela de Peças de um processo que não é de software: sem comandos de Go.
func writePecasGenericas(b *strings.Builder) {
	b.WriteString("| `PECA-02` | `unit_task` | Execução das etapas | Cada etapa grava o seu artefato e passa por um gate de conteúdo. |\n")
	b.WriteString("| `PECA-03` | `script` | `legoloop validate` | Validação mecânica do grafo do processo. |\n")
	b.WriteString("| `PECA-04` | `stop_gate` | Aprovação humana | Pausa antes de qualquer gravação definitiva. |\n\n---\n\n")
}

func writeIntentGovernance(b *strings.Builder, v IntentView) {
	m := v.Mirror
	b.WriteString("## 4. Governança Operacional e Restrições Inegociáveis\n")
	b.WriteString("* **Gatilho Operacional de Disparo (Trigger):** `manual`  \n  *(Define como este Processo acordará no dia a dia quando estiver construído.)*\n")
	b.WriteString("* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01\n")
	b.WriteString("* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01\n")
	fmt.Fprintf(b, "* **Contrato de Evidência Factual (Definition of Done Visual/Real):** %s\n", orDefault(m.Success, notInformed))
	b.WriteString("* **Restrições Inegociáveis:**\n")
	b.WriteString(bullets(m.Constraints, "  ", "Nenhuma restrição registrada na entrevista"))
	b.WriteString("* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.\n")
	fmt.Fprintf(b, "* **Riscos e Perguntas em Aberto:** %s\n\n---\n\n", openRisks(v.Confidence))
}

func openRisks(confidence int) string {
	if confidence >= 100 {
		return "Nenhum registrado."
	}
	return fmt.Sprintf("Confiança de %d%%: as lacunas restantes devem ser tratadas no Stop-Gate 01.", confidence)
}

func writeIntentCriteria(b *strings.Builder, v IntentView) {
	b.WriteString("## 5. Critérios de Sucesso de Alto Nível\n")
	b.WriteString("Cada critério precisa ter uma forma mecânica de prova.\n\n")
	b.WriteString("| # | Critério | Como se prova | Comando / Evidência |\n| :--- | :--- | :--- | :--- |\n")
	if ehSoftware(v.Natureza) {
		b.WriteString("| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |\n")
		b.WriteString("| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |\n")
	} else {
		b.WriteString("| 1 | Cada etapa grava o artefato prometido, e o gate de conteúdo dela passa. | Gate da etapa avaliado sobre o artefato | `legoloop run` com exit code 0 |\n")
		b.WriteString("| 2 | O grafo do processo é válido antes de executar. | Validação mecânica do grafo | `legoloop validate` com exit code 0 |\n")
	}
	fmt.Fprintf(b, "| 3 | %s | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |\n\n---\n\n",
		cell(orDefault(v.Mirror.Success, notInformed)))
}

func writeIntentHandoff(b *strings.Builder, v IntentView) {
	b.WriteString("## 6. Handoff de Governança\n")
	fmt.Fprintf(b, "* **Próximo Documento Mandatório:** [`%s`](../specs/%s)\n", v.SpecRel, specBase(v.SpecRel))
	b.WriteString("* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.\n\n---\n\n")
	b.WriteString("🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  \n")
	fmt.Fprintf(b, "*A intenção consolidada com %d%% de certeza factual atende à sua expectativa?*  \n", v.Confidence)
	b.WriteString("*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*\n")
}

func specBase(rel string) string {
	if i := strings.LastIndex(rel, "/"); i >= 0 {
		return rel[i+1:]
	}
	return rel
}

// SpecRelFor deriva o caminho da SPEC a partir do nome da INTENT (mesma convenção de datas e nomes).
func SpecRelFor(fileName, theme, date string) string {
	if strings.Contains(fileName, "_INTENT_") {
		return "docs/specs/" + strings.Replace(fileName, "_INTENT_", "_SPEC_", 1)
	}
	return fmt.Sprintf("docs/specs/%s_SPEC_%s.md", date, sanitizeFileName(theme))
}
