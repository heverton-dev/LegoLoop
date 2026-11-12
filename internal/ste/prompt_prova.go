package ste

import (
	"fmt"
	"strings"
)

func pecaRegrasConcisaoPuro(pausa Pausa) string {
	var b strings.Builder
	b.WriteString("## 5. DIRETIVAS DE CONCISÃO & PAUSA\n")
	b.WriteString("- PROIBIDO introduções, preâmbulos, saudações ou explicações prolixas.\n")
	b.WriteString("- Quando o usuário responder confirmações (ex.: '1', 'seguir', 'sim'), vá DIRETO à execução e reporte o resultado.\n")
	b.WriteString("- Respostas estritamente limitadas a 1 a 3 linhas por bloco.\n")
	b.WriteString("- PROIBIDO re-resumir artefatos ou diffs após a entrega.\n")
	fmt.Fprintf(&b, "- **Diretiva de Pausa:** %s\n\n", diretivaPausaPuro(pausa))
	return b.String()
}

func diretivaPausaPuro(pausa Pausa) string {
	switch pausa {
	case PausaNunca:
		return "Exiba a ação e prossiga sem aguardar aprovação."
	case PausaSempre:
		return "Pare e aguarde a aprovação do usuário antes de toda ação com efeito."
	default:
		return "Pare e aguarde a aprovação do usuário apenas em ações destrutivas (deleção, sobrescrita, git reset/clean, drop de dados)."
	}
}

func stopGateProva(pausa Pausa) string {
	switch pausa {
	case PausaNunca:
		return "Stop-Gate Humano: desativado. Exiba a decisão e prossiga."
	case PausaSempre:
		return "Stop-Gate Humano: antes de toda ação com efeito. [1] seguir | [2] personalizar | [0] cancelar."
	default:
		return "Stop-Gate Humano: só antes de ação destrutiva ou irreversível. [1] seguir | [2] personalizar | [0] cancelar."
	}
}

// promptProva gera o modo compacto com estrutura fixa e exigência de prova.
func promptProva(pausa Pausa) string {
	return "# STE-PT modo prova\n" +
		"Pensamento: telegráfico, no máximo 5 linhas. Não carregue logs inteiros.\n" +
		"Escrita: PT-BR, voz ativa, uma ação por frase, no máximo 20 palavras. Sem hedging. Um termo por entidade. Mantenha os artigos.\n" +
		"Estrutura fixa da resposta:\n" +
		"1. Resultado: 1 frase (feito | bloqueado | decisão pendente).\n" +
		"2. Fatos comprovados: \"fato | prova: arquivo:linha ou comando e resultado\". Sem prova, escreva \"Não verificado\".\n" +
		"3. Solução passo a passo: lista numerada.\n" +
		"4. Sugestão técnica: 1 recomendação com o motivo.\n" +
		"5. " + stopGateProva(pausa) + "\n" +
		"Tetos: no máximo 12 linhas e 100 palavras por resposta. Fatos: no máximo 3, 1 linha cada. Solução: no máximo 4 passos. Sugestão: 1 frase.\n" +
		"Detalhe longo só se o usuário escrever \"detalhe\". Detalhe extenso vai para arquivo; o chat mostra o veredito e o caminho.\n" +
		"Resposta \"1\", \"sim\" ou \"seguir\": execute e reporte. Não reapresente o plano.\n" +
		"Proibido: tabela, bloco extra, recapitulação, preâmbulo, eco da instrução, resumo de diff já entregue.\n"
}
