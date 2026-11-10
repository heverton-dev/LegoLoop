package taxonomy_test

import (
	"strings"
	"testing"

	"legoloop/internal/taxonomy"
)

func TestCheckForbidden(t *testing.T) {
	tests := []struct {
		input         string
		wantCanonical string
		wantForbidden bool
	}{
		{"esteira", "Processo", true},
		{"ESTEIRA", "Processo", true},
		{"loop", "Processo", true},
		{"pipeline", "Processo", true},
		{"task", "Nó", true},
		{"componente", "Peça", true},
		{"Processo", "", false},
		{"palavra_qualquer", "", false},
	}

	for _, tt := range tests {
		gotCanonical, gotForbidden := taxonomy.CheckForbidden(tt.input)
		if gotForbidden != tt.wantForbidden {
			t.Errorf("CheckForbidden(%q) gotForbidden = %v, want %v", tt.input, gotForbidden, tt.wantForbidden)
		}
		if gotCanonical != tt.wantCanonical {
			t.Errorf("CheckForbidden(%q) gotCanonical = %q, want %q", tt.input, gotCanonical, tt.wantCanonical)
		}
	}
}

func TestValidateText(t *testing.T) {
	sample := `
# Documento de Teste
Este sistema usa uma esteira para executar cada task.
Também possui um componente legado.
O Processo é canônico e o Nó está correto.
`
	violations, err := taxonomy.ValidateText(sample)
	if err != nil {
		t.Fatalf("ValidateText falhou: %v", err)
	}

	if len(violations) < 3 {
		t.Fatalf("Esperava pelo menos 3 violações, obteve %d", len(violations))
	}

	foundEsteira := false
	for _, v := range violations {
		if strings.EqualFold(v.ForbiddenTerm, "esteira") && v.ApprovedTerm == "Processo" {
			foundEsteira = true
		}
	}
	if !foundEsteira {
		t.Errorf("Violação 'esteira' -> 'Processo' não detectada")
	}
}
