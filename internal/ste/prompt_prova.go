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

// promptProva gera o modo compacto: prosa curta em grupos, prova inline, Recomendação final.
func promptProva(pausa Pausa) string {
	return "# STE-PT modo prova\n" +
		"Pensamento: telegráfico, no máximo 5 linhas. Não carregue logs inteiros.\n" +
		"Escrita: PT-BR, voz ativa, uma ideia por frase, no máximo 20 palavras. Sem hedging. Um termo por entidade. Mantenha os artigos.\n" +
		"Formato: uma frase por linha, sem linha em branco entre as frases. Nunca junte frases em parágrafo. Sem cabeçalho, tabela, separador ou lista numerada. Vale também para mensagens intermediárias.\n" +
		"Fato que decide leva a prova entre parênteses: (arquivo:linha ou comando e resultado). Sem comando rodado nesta resposta, escreva \"Não verificado\". Citar documento não é prova.\n" +
		"Última linha: \"Recomendação:\" com uma ação, separada das frases por uma linha em branco.\n" +
		pausaProva(pausa) + "\n" +
		"Tetos: 15 linhas e 130 palavras. Detalhe longo só se o usuário escrever \"detalhe\"; ele vai para arquivo.\n" +
		"Resposta \"1\", \"sim\" ou \"seguir\": execute e reporte. Não reapresente o plano.\n" +
		"Exemplo:\n" +
		"A correção dos provedores concluiu com êxito (internal/cli/homologar.go:349).\n" +
		"Os arquivos de opencode e mimocode agora apontam para o host canônico.\n" +
		"As cópias de segurança foram criadas para cada arquivo modificado.\n" +
		"O próximo par para teste é claude|9router.\n" +
		"O teto estimado contempla 3 chamadas e até 6 minutos.\n" +
		"\n" +
		"Recomendação: autorizar a execução da homologação do par claude|9router.\n"
}
