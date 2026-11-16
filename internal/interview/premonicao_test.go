package interview

import "testing"

func TestPremonicaoExigeTresPerguntasAbaixoDeNoventaECinco(t *testing.T) {
	rawSemPremonicao := `HIPÓTESE: CLI em Go
CONFIANÇA: 80%
ESPELHO CANÔNICO:
- Resultado Esperado: CLI funcional
- Quem se Beneficia: Devs
- Por que Agora: Automação
- Critério de Sucesso: Testes verdes
- Restrições Inegociáveis: Stdlib
- Fora de Escopo: UI`

	turn := ParseTurnResult(rawSemPremonicao)
	if turn.Status == StatusCompleted {
		t.Fatalf("esperava StatusInProgress sem premonição e confiança < 95%%, obteve %s", turn.Status)
	}
}

func TestPremonicaoComDuasPerguntasNaoConclui(t *testing.T) {
	rawDuas := `HIPÓTESE: CLI em Go
CONFIANÇA: 85%
PREMONIÇÃO:
1. Pergunta: Qual linguagem? | Resposta prevista: Go
2. Pergunta: Qual banco? | Resposta prevista: Nenhum
ESPELHO CANÔNICO:
- Resultado Esperado: CLI funcional`

	turn := ParseTurnResult(rawDuas)
	if turn.Status == StatusCompleted {
		t.Fatalf("esperava StatusInProgress com apenas 2 perguntas na premonição, obteve %s", turn.Status)
	}
}

func TestPremonicaoComTresPerguntasConcluiAbaixoDeNoventaECinco(t *testing.T) {
	rawTres := `HIPÓTESE: CLI em Go
CONFIANÇA: 85%
PREMONIÇÃO:
1. Pergunta: Qual linguagem? | Resposta prevista: Go
2. Pergunta: Qual banco? | Resposta prevista: Nenhum
3. Pergunta: Qual entrega? | Resposta prevista: Binário único
ESPELHO CANÔNICO:
- Resultado Esperado: CLI funcional`

	turn := ParseTurnResult(rawTres)
	if turn.Status != StatusCompleted {
		t.Fatalf("esperava StatusCompleted com 3 perguntas na premonição, obteve %s", turn.Status)
	}
}

func TestConfiancaNoventaECincoConcluiSemPremonicao(t *testing.T) {
	raw := `HIPÓTESE: CLI em Go
CONFIANÇA: 95%
ESPELHO CANÔNICO:
- Resultado Esperado: CLI funcional`

	turn := ParseTurnResult(raw)
	if turn.Status != StatusCompleted {
		t.Fatalf("esperava StatusCompleted com confiança >= 95%%, obteve %s", turn.Status)
	}
}
