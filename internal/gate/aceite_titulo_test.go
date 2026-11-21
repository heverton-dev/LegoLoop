package gate

import (
	"os"
	"testing"
)

// TestAceiteReconheceOTituloCanonicoDoTemplate: o título "Critérios Binários de Aceite" é o do
// TEMPLATE_SPEC.md e de dez SPECs do projeto. O Gate has_acceptance_criteria precisa reconhecê-lo.
func TestAceiteReconheceOTituloCanonicoDoTemplate(t *testing.T) {
	for nome, arq := range map[string]string{
		"template":      "../../docs/templates/TEMPLATE_SPEC.md",
		"spec-da-forja": "../../docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP.md",
	} {
		b, err := os.ReadFile(arq)
		if err != nil {
			t.Fatalf("%s: %v", nome, err)
		}
		env := newEnv(t, map[string]string{"SPEC.md": string(b)}, 0)
		if !evalExpr(t, "has_acceptance_criteria('SPEC.md')", env) {
			t.Errorf("%s não passou em has_acceptance_criteria", nome)
		}
	}
	env := newEnv(t, map[string]string{"A.md": "## 5. Critérios Binários de Aceite por Nó\n- [ ] AC-1 ok\n", "B.md": "## Notas\n- [ ] algo\n"}, 0)
	if !evalExpr(t, "has_acceptance_criteria('A.md')", env) {
		t.Error("título com Binários e Aceite deve passar")
	}
	if evalExpr(t, "has_acceptance_criteria('B.md')", env) {
		t.Error("arquivo sem título de critérios não pode passar")
	}
}
