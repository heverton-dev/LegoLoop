package interview

import (
	"context"
	"fmt"
	"strings"
)

// msgPedirPremonicao é a correção enviada ao cérebro quando ele emite o Espelho sem a premonição exigida.
const msgPedirPremonicao = "Você emitiu o ESPELHO CANÔNICO com CONFIANÇA abaixo de 95% e sem a seção PREMONIÇÃO. " +
	"Emita agora a seção PREMONIÇÃO com exatamente 3 itens no formato `N. Pergunta: ... | Resposta prevista: ...`, " +
	"seguida do ESPELHO CANÔNICO. Se ainda houver lacuna real, faça UMA PERGUNTA em vez disso."

// msgPedirEspelho é a correção enviada quando o cérebro declara confiança de 95% ou mais sem emitir o Espelho.
const msgPedirEspelho = "Você declarou CONFIANÇA de 95% ou mais sem emitir o ESPELHO CANÔNICO. " +
	"Emita agora o ESPELHO CANÔNICO com os seis campos (Resultado Esperado, Quem se Beneficia, Por que Agora, " +
	"Critério de Sucesso, Restrições Inegociáveis, Fora de Escopo). Se ainda houver lacuna real, " +
	"reduza a CONFIANÇA e faça UMA PERGUNTA."

// perguntaDeConfirmacao é o próximo passo do usuário quando o cérebro insiste em não emitir a premonição.
const perguntaDeConfirmacao = "O Espelho acima está completo, mas a confiança ficou abaixo de 95% e o cérebro não emitiu a premonição. " +
	"Você reconhece o Espelho como está ou quer corrigir algum dos seis campos?"

// perguntaSemEspelho é o próximo passo quando o cérebro insiste em declarar confiança alta sem Espelho.
const perguntaSemEspelho = "O cérebro declarou confiança alta, mas não emitiu o Espelho Canônico. " +
	"Quer que ele o emita agora (responda \"emita o Espelho\") ou prefere responder a mais uma pergunta?"

// semEspelho diz se o turno declarou confiança de 95% ou mais sem Espelho e sem concluir.
func semEspelho(turn *TurnResult) bool {
	return turn.Status != StatusCompleted && turn.Mirror == nil && turn.Confidence >= 95
}

// emImpasse diz se o turno ficaria sem saída: o Espelho veio sem pergunta e sem conclusão,
// ou a confiança alta veio sem Espelho.
func emImpasse(turn *TurnResult) bool {
	if semEspelho(turn) {
		return true
	}
	return turn.Status != StatusCompleted && turn.Mirror != nil && strings.TrimSpace(turn.Question) == ""
}

// pedirPremonicao faz uma única correção ao cérebro. Se ele ainda não concluir nem perguntar,
// o motor devolve uma pergunta de confirmação, nunca uma pergunta vazia e nunca uma sessão sem saída.
func (e *Engine) pedirPremonicao(ctx context.Context, sess *Session, antes *TurnResult) (*TurnResult, error) {
	msg, fallback := msgPedirPremonicao, perguntaDeConfirmacao
	if semEspelho(antes) {
		msg, fallback = msgPedirEspelho, perguntaSemEspelho
	}
	sess.Messages = append(sess.Messages, Message{Role: "user", Content: msg})
	raw, err := e.llm.Chat(ctx, sess.Messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM: %w", err)
	}
	sess.Messages = append(sess.Messages, Message{Role: "assistant", Content: raw})
	turn := ParseTurnResult(raw)
	if turn.Status != StatusCompleted && (emImpasse(turn) || strings.TrimSpace(turn.Question) == "") {
		turn.Question = fallback
		if semEspelho(turn) {
			turn.Confidence = antes.Confidence
		}
	}
	return turn, nil
}
