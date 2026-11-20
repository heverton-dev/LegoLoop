package plangen_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/plangen"
)

// O ROUTE e o PLAN gerados declaram o runtime do motor (legoloop), nunca um shell de um só sistema operacional.
func TestRouteGeradoDeclaraRuntimeDoMotor(t *testing.T) {
	saida := "### [TICKET-E01.1] Criar tipos\n- **Épico:** E01\n- **Natureza:** micro\n- **Modo:** tdd_red_green\n- **Dependências:** nenhuma\n- **Arquivos Afetados:** internal/service/types.go\n- **Critérios Atendidos:** AC-E01.1\n- **Comando de Validação:** go test ./internal/service/...\n"
	dir := t.TempDir()
	spec := filepath.Join(dir, "2026-11-02_SPEC_SERVICO.md")
	if err := os.WriteFile(spec, []byte("# SPEC Servico"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := plangen.NewGenerator(&mockPlanLLM{output: saida}).Generate(context.Background(), dir, spec, plangen.SupervisionMicro)
	if err != nil {
		t.Fatal(err)
	}
	for _, caminho := range []string{res.RouteJSONPath, res.PlanMarkdownPath} {
		b, err := os.ReadFile(caminho)
		if err != nil {
			t.Fatal(err)
		}
		texto := strings.ToLower(string(b))
		if strings.Contains(texto, "powershell") {
			t.Errorf("%s fixa powershell no runtime", filepath.Base(caminho))
		}
		if !strings.Contains(texto, "legoloop") {
			t.Errorf("%s não declara o runtime legoloop", filepath.Base(caminho))
		}
	}
}
