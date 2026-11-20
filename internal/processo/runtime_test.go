package processo_test

import (
	"strings"
	"testing"
)

// As projeções declaram o runtime do motor (legoloop), nunca um shell de um só sistema operacional.
func TestProjecoesDeclaramRuntimeDoMotor(t *testing.T) {
	pr := projetar(t, quatro())
	for nome, b := range map[string][]byte{"dag": pr.DAG, "studio": pr.Studio, "route": pr.Route} {
		texto := strings.ToLower(string(b))
		if strings.Contains(texto, "powershell") {
			t.Errorf("%s fixa powershell no runtime", nome)
		}
	}
}
