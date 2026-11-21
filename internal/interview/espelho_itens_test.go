package interview

import (
	"strings"
	"testing"
)

// Achado da I02: realce no meio do texto e itens quebrados nas vírgulas apareciam na INTENT.
func TestEspelhoNaoDeixaRealceNoMeioDoTexto(t *testing.T) {
	raw := "CONFIANÇA: 96%\nESPELHO CANÔNICO:\n- Resultado Esperado: x\n- Quem se Beneficia: **Operador**: executa o ciclo. **Outra sessão**: reutiliza.\n- Por que Agora: y\n- Critério de Sucesso: z\n- Restrições Inegociáveis: a\n- Fora de Escopo: b\n"
	m := ParseTurnResult(raw).Mirror
	if m == nil {
		t.Fatal("espelho ausente")
	}
	if strings.Contains(m.User, "**") {
		t.Fatalf("realce Markdown vazou para o campo: %q", m.User)
	}
}

func TestRestricoesEmLinhasNaoQuebramNasVirgulas(t *testing.T) {
	raw := "CONFIANÇA: 96%\nESPELHO CANÔNICO:\n- Resultado Esperado: x\n- Quem se Beneficia: u\n- Por que Agora: y\n- Critério de Sucesso: z\n- Restrições Inegociáveis:\n  - **Zero Stubs**: PLAN, ROUTE e DAG são gerados de verdade.\n  - **TDD**: o teste falha antes do código.\n- Fora de Escopo: b\n"
	m := ParseTurnResult(raw).Mirror
	if len(m.Constraints) != 2 {
		t.Fatalf("esperava 2 restrições (uma por linha), obteve %d: %q", len(m.Constraints), m.Constraints)
	}
	if !strings.Contains(m.Constraints[0], "PLAN, ROUTE e DAG") {
		t.Fatalf("a restrição foi cortada na vírgula: %q", m.Constraints[0])
	}
}

func TestListaEmUmaLinhaContinuaSeparadaPorVirgula(t *testing.T) {
	raw := "CONFIANÇA: 96%\nESPELHO CANÔNICO:\n- Resultado Esperado: x\n- Quem se Beneficia: u\n- Por que Agora: y\n- Critério de Sucesso: z\n- Restrições Inegociáveis: stdlib, TDD, sem LLM\n- Fora de Escopo: b\n"
	m := ParseTurnResult(raw).Mirror
	if len(m.Constraints) != 3 {
		t.Fatalf("esperava 3 itens em uma linha com vírgulas, obteve %d: %q", len(m.Constraints), m.Constraints)
	}
}
