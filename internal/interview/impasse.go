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

// perguntaDeConfirmacao é o próximo passo do usuário quando o cérebro insiste em não emitir a premonição.
const perguntaDeConfirmacao = "O Espelho acima está completo, mas a confiança ficou abaixo de 95% e o cérebro não emitiu a premonição. " +
	"Você reconhece o Espelho como está ou quer corrigir algum dos seis campos?"

// emImpasse diz se o turno trouxe o Espelho, não concluiu e não deixou pergunta: a sessão ficaria parada.
func emImpasse(turn *TurnResult) bool {
	return turn.Status != StatusCompleted && turn.Mirror != nil && strings.TrimSpace(turn.Question) == ""
}

// pedirPremonicao faz uma única correção. Se o cérebro ainda não concluir nem perguntar,
// o motor devolve uma pergunta de confirmação do Espelho, nunca uma pergunta vazia.
func (e *Engine) pedirPremonicao(ctx context.Context, sess *Session) (*TurnResult, error) {
	sess.Messages = append(sess.Messages, Message{Role: "user", Content: msgPedirPremonicao})
	raw, err := e.llm.Chat(ctx, sess.Messages)
	if err != nil {
		return nil, fmt.Errorf("falha ao invocar LLM: %w", err)
	}
	sess.Messages = append(sess.Messages, Message{Role: "assistant", Content: raw})
	turn := ParseTurnResult(raw)
	if emImpasse(turn) || (turn.Status != StatusCompleted && strings.TrimSpace(turn.Question) == "") {
		turn.Question = perguntaDeConfirmacao
	}
	return turn, nil
}
