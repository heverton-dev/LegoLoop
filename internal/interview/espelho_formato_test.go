package interview_test

import (
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/interview"
)

// espelhoFechaNegritoDepoisDosDoisPontos reproduz o formato que o claude emite:
// o realce fecha depois dos dois-pontos (**Campo:** texto).
const espelhoFechaNegritoDepoisDosDoisPontos = "CONFIANÇA: 96%\n\nESPELHO CANÔNICO:\n" +
	"- **Resultado Esperado:** instalação universal do LegoLoop\n" +
	"- **Quem se Beneficia:** usuários de qualquer canal\n" +
	"- **Por que Agora:** não existe instalação universal\n" +
	"- **Critério de Sucesso:** o instalar roda sem erro\n" +
	"- **Restrições Inegociáveis:** sem segredo, sem sobrescrever\n" +
	"- **Fora de Escopo:** nada\n"

// espelhoComListas reproduz o campo cujo valor é uma lista em várias linhas.
const espelhoComListas = "CONFIANÇA: 96%\n\nESPELHO CANÔNICO:\n" +
	"- **Resultado Esperado:** instalação universal\n" +
	"- **Quem se Beneficia:** todos os usuários\n" +
	"- **Por que Agora:** não existe instalador\n" +
	"- **Critério de Sucesso:** em máquina limpa:\n" +
	"  - o instalar roda sem erro\n" +
	"  - o segundo dry-run não acha pendência\n" +
	"- **Restrições Inegociáveis:**\n" +
	"  - nunca sobrescrever sem backup\n" +
	"  - nunca imprimir segredos\n" +
	"- **Fora de Escopo:** nada\n"

func TestEspelhoNegritoDepoisDosDoisPontosNaoVazaAsteriscos(t *testing.T) {
	m := interview.ParseTurnResult(espelhoFechaNegritoDepoisDosDoisPontos).Mirror
	if m == nil {
		t.Fatal("espelho não extraído")
	}
	if m.Outcome != "instalação universal do LegoLoop" {
		t.Errorf("Outcome = %q", m.Outcome)
	}
	for nome, v := range map[string]string{"User": m.User, "WhyNow": m.WhyNow, "Success": m.Success} {
		if strings.Contains(v, "*") {
			t.Errorf("%s vazou realce: %q", nome, v)
		}
	}
	if len(m.Constraints) != 2 || m.Constraints[0] != "sem segredo" {
		t.Errorf("Constraints = %q", m.Constraints)
	}
}

func TestEspelhoComListaEmVariasLinhasNaoTruncaNoPrimeiro(t *testing.T) {
	m := interview.ParseTurnResult(espelhoComListas).Mirror
	if m == nil {
		t.Fatal("espelho não extraído")
	}
	if !strings.Contains(m.Success, "segundo dry-run") {
		t.Errorf("Success truncado: %q", m.Success)
	}
	if len(m.Constraints) != 2 || m.Constraints[1] != "nunca imprimir segredos" {
		t.Errorf("Constraints = %q", m.Constraints)
	}
	if len(m.OutOfScope) != 1 || m.OutOfScope[0] != "nada" {
		t.Errorf("OutOfScope = %q", m.OutOfScope)
	}
}

func TestMarcadoresComNegritoDepoisDosDoisPontos(t *testing.T) {
	raw := "**HIPÓTESE:** a instalação é uma só\n**CONFIANÇA:** 96%\n**PERGUNTA:** qual canal sai primeiro?\n**PALPITE:** o script de terminal\n"
	turn := interview.ParseTurnResult(raw)
	if !turn.ConfidenceReported || turn.Confidence != 96 {
		t.Errorf("confiança = %d (informada=%v), quer 96", turn.Confidence, turn.ConfidenceReported)
	}
	if turn.Hypothesis != "a instalação é uma só" || turn.Question != "qual canal sai primeiro?" || turn.Guess != "o script de terminal" {
		t.Errorf("marcadores = %q | %q | %q", turn.Hypothesis, turn.Question, turn.Guess)
	}
}

func TestEspelhoNaoAbsorveTextoDepoisDeLinhaEmBranco(t *testing.T) {
	raw := "CONFIANÇA: 96%\n\nESPELHO CANÔNICO:\n" +
		"- **Resultado Esperado:** x\n- **Quem se Beneficia:** y\n- **Por que Agora:** z\n" +
		"- **Critério de Sucesso:** w\n- **Restrições Inegociáveis:** a\n- **Fora de Escopo:** nada\n\n" +
		"Pontos em aberto para a especificação, que não são escopo.\n"
	m := interview.ParseTurnResult(raw).Mirror
	if m == nil || len(m.OutOfScope) != 1 || m.OutOfScope[0] != "nada" {
		t.Errorf("OutOfScope absorveu texto solto: %+v", m)
	}
}
