package interview_test

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/interview"
)

const espelhoSemPremonicao = `HIPÓTESE: o escopo está claro.
CONFIANÇA: 92%
ESPELHO CANÔNICO:
- Resultado Esperado: dois Processos canônicos
- Quem se Beneficia: o operador
- Por que Agora: a trava já existe
- Critério de Sucesso: quatro provas mecânicas
- Restrições Inegociáveis: TDD, sem LLM
- Fora de Escopo: alterar a INTENT da I01`

const premonicaoCompleta = `HIPÓTESE: o escopo está claro.
CONFIANÇA: 92%
PREMONIÇÃO:
1. Pergunta: Qual a ordem? | Resposta prevista: depois do TICK-01.2
2. Pergunta: Qual o formato? | Resposta prevista: Markdown
3. Pergunta: Qual a prova? | Resposta prevista: teste de ponta a ponta
ESPELHO CANÔNICO:
- Resultado Esperado: dois Processos canônicos`

// sequenciaLLM devolve as respostas em ordem e guarda as mensagens que recebeu.
type sequenciaLLM struct {
	respostas []string
	chamadas  [][]interview.Message
}

func (s *sequenciaLLM) Chat(_ context.Context, msgs []interview.Message) (string, error) {
	s.chamadas = append(s.chamadas, append([]interview.Message(nil), msgs...))
	i := len(s.chamadas) - 1
	if i >= len(s.respostas) {
		i = len(s.respostas) - 1
	}
	return s.respostas[i], nil
}

// Achado da I02: o cérebro emitia o Espelho com confiança abaixo de 95%, sem premonição e sem pergunta,
// e a sessão ficava parada. O motor deve pedir a premonição que falta, uma vez, e concluir.
func TestEspelhoSemPremonicaoPedeAPremonicaoEConclui(t *testing.T) {
	llm := &sequenciaLLM{respostas: []string{espelhoSemPremonicao, premonicaoCompleta}}
	_, turn, err := interview.NewEngine(llm).Start(context.Background(), "ideia")
	if err != nil {
		t.Fatal(err)
	}
	if turn.Status != interview.StatusCompleted {
		t.Fatalf("esperava COMPLETED depois de pedir a premonição, obteve %s", turn.Status)
	}
	if len(llm.chamadas) != 2 {
		t.Fatalf("esperava exatamente 2 chamadas ao cérebro, obteve %d", len(llm.chamadas))
	}
	ultima := llm.chamadas[1][len(llm.chamadas[1])-1].Content
	if !strings.Contains(strings.ToUpper(ultima), "PREMONIÇÃO") {
		t.Fatalf("a mensagem corretiva deveria pedir a PREMONIÇÃO, mas foi: %q", ultima)
	}
}

// Se o cérebro insistir em não emitir a premonição, o usuário recebe uma pergunta de confirmação do Espelho.
// Nunca uma pergunta vazia.
func TestEspelhoSemPremonicaoInsistenteGeraPerguntaDeConfirmacao(t *testing.T) {
	llm := &sequenciaLLM{respostas: []string{espelhoSemPremonicao}}
	_, turn, err := interview.NewEngine(llm).Start(context.Background(), "ideia")
	if err != nil {
		t.Fatal(err)
	}
	if turn.Status == interview.StatusCompleted {
		t.Fatal("sem premonição e com confiança 92% a sessão não pode concluir")
	}
	if strings.TrimSpace(turn.Question) == "" {
		t.Fatal("a pergunta não pode ficar vazia: o usuário precisa de um próximo passo")
	}
	if len(llm.chamadas) != 2 {
		t.Fatalf("a correção acontece uma única vez; chamadas: %d", len(llm.chamadas))
	}
}

// Um turno normal, com pergunta, não pode disparar a correção.
func TestTurnoComPerguntaNaoDisparaCorrecao(t *testing.T) {
	llm := &sequenciaLLM{respostas: []string{"HIPÓTESE: x\nCONFIANÇA: 60%\nPERGUNTA: Qual o prazo?\nPALPITE: uma semana"}}
	_, turn, err := interview.NewEngine(llm).Start(context.Background(), "ideia")
	if err != nil {
		t.Fatal(err)
	}
	if len(llm.chamadas) != 1 || turn.Question != "Qual o prazo?" {
		t.Fatalf("turno normal não deve ser alterado: chamadas=%d pergunta=%q", len(llm.chamadas), turn.Question)
	}
}
