package interview_test

import (
	"context"
	"strings"
	"testing"

	"legoloop/internal/interview"
)

const confiancaAltaSemEspelho = "HIPÓTESE: entendi tudo.\nCONFIANÇA: 98%\nPERGUNTA: Qual o formato do erro?"

const espelhoCompleto = `HIPÓTESE: entendi tudo.
CONFIANÇA: 98%
ESPELHO CANÔNICO:
- Resultado Esperado: despachante no runner
- Quem se Beneficia: o operador
- Por que Agora: a I02 entrega a ROUTE
- Critério de Sucesso: testes verdes
- Restrições Inegociáveis: sem LLM
- Fora de Escopo: formato do DAG`

// Achado da I03: o motor concluía a sessão com confiança >= 95% mesmo sem Espelho. A sessão ficava sem saída:
// a INTENT não materializava (sem Espelho) e a entrevista não aceitava mais respostas (já finalizada).
func TestConfiancaAltaSemEspelhoNaoConclui(t *testing.T) {
	turn := interview.ParseTurnResult(confiancaAltaSemEspelho)
	if turn.Status == interview.StatusCompleted {
		t.Fatal("sem Espelho a sessão não pode ser concluída, mesmo com confiança 98%")
	}
}

func TestConfiancaAltaSemEspelhoPedeOEspelhoEConclui(t *testing.T) {
	llm := &sequenciaLLM{respostas: []string{confiancaAltaSemEspelho, espelhoCompleto}}
	eng := interview.NewEngine(llm)
	sess, turn, err := eng.Start(context.Background(), "ideia")
	if err != nil {
		t.Fatal(err)
	}
	if turn.Status != interview.StatusCompleted || turn.Mirror == nil {
		t.Fatalf("esperava COMPLETED com Espelho depois de pedi-lo; status=%s espelho=%v", turn.Status, turn.Mirror)
	}
	if len(llm.chamadas) != 2 {
		t.Fatalf("esperava 2 chamadas, obteve %d", len(llm.chamadas))
	}
	ultima := llm.chamadas[1][len(llm.chamadas[1])-1].Content
	if !strings.Contains(strings.ToUpper(ultima), "ESPELHO") {
		t.Fatalf("a correção deveria pedir o ESPELHO: %q", ultima)
	}
	if _, err := eng.MaterializeIntent(sess.ID, t.TempDir(), ""); err != nil {
		t.Fatalf("a INTENT deveria materializar: %v", err)
	}
}

// Se o cérebro insistir em não emitir o Espelho, a sessão continua aberta e aceita resposta.
func TestConfiancaAltaSemEspelhoInsistenteContinuaAberta(t *testing.T) {
	llm := &sequenciaLLM{respostas: []string{confiancaAltaSemEspelho}}
	eng := interview.NewEngine(llm)
	sess, turn, err := eng.Start(context.Background(), "ideia")
	if err != nil {
		t.Fatal(err)
	}
	if turn.Status == interview.StatusCompleted {
		t.Fatal("a sessão não pode concluir sem Espelho")
	}
	if strings.TrimSpace(turn.Question) == "" {
		t.Fatal("a pergunta não pode ficar vazia")
	}
	if _, err := eng.Step(context.Background(), sess.ID, "responda"); err != nil {
		t.Fatalf("a sessão deveria aceitar resposta, mas: %v", err)
	}
}
