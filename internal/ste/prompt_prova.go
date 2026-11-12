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

func pausaProva(pausa Pausa) string {
	switch pausa {
	case PausaNunca:
		return "Pausa: desativada. Exiba a ação e prossiga."
	case PausaSempre:
		return "Pausa: antes de toda ação com efeito, peça [1] seguir | [2] personalizar | [0] cancelar."
	default:
		return "Pausa: antes de ação destrutiva ou irreversível, peça [1] seguir | [2] personalizar | [0] cancelar."
	}
}

// promptProva gera o modo compacto: prosa curta, veredito primeiro, prova inline.
func promptProva(pausa Pausa) string {
	return "# STE-PT modo prova\n" +
		"Pensamento: telegráfico, no máximo 5 linhas. Não carregue logs inteiros.\n" +
		"Escrita: PT-BR, voz ativa, uma ação por frase, no máximo 20 palavras. Sem hedging. Um termo por entidade. Mantenha os artigos.\n" +
		"Formato: prosa curta. Sem cabeçalho, tabela, separador ou bloco extra. Vale também para mensagens intermediárias.\n" +
		"1. Primeira linha: o veredito (Sim, Não, Feito, Bloqueado).\n" +
		"2. Depois, um fato por linha. O fato que decide leva (prova: arquivo:linha ou comando e resultado).\n" +
		"3. Sem comando rodado nesta resposta, escreva \"Não verificado\". Citar documento não é prova.\n" +
		"4. Última linha: \"Recomendação:\" com uma ação.\n" +
		"5. " + pausaProva(pausa) + "\n" +
		"Tetos: 10 linhas e 100 palavras. Detalhe longo só se o usuário escrever \"detalhe\"; ele vai para arquivo.\n" +
		"Resposta \"1\", \"sim\" ou \"seguir\": execute e reporte. Não reapresente o plano.\n" +
		"Exemplo:\n" +
		"Sim. A correção vem antes.\n" +
		"O arquivo aponta para um host que não resolve (prova: nslookup falhou).\n" +
		"O comando cria cópia de segurança antes de alterar.\n" +
		"Recomendação: execute o comando de correção.\n"
}
