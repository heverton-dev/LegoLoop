package inventario

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPendentesSemDonoViramTickets prova o caminho de falha do AC-E01.4.
func TestPendentesSemDonoViramTickets(t *testing.T) {
	entradas := []Entrada{
		{Caminho: "plans/sem-padrao.md", Classe: ClassePendente, Motivo: "demanda sem padrão ISO"},
		{Caminho: "prompts/prompt.md", Classe: ClasseRever, Motivo: "diretório fora do previsto"},
		{Caminho: "intents/2026-11-08_INTENT_X.md", Classe: ClasseImplementado, Dono: "X"},
		{Caminho: "reports/historico.md", Classe: ClasseUltrapassado, Dono: "HIST"},
	}
	saida := filepath.Join(t.TempDir(), "INVENTARIO.md")
	if err := GerarRelatorio(entradas, saida); err != nil {
		t.Fatalf("GerarRelatorio: %v", err)
	}
	conteudo, err := os.ReadFile(saida)
	if err != nil {
		t.Fatalf("leitura do relatório: %v", err)
	}
	texto := string(conteudo)
	for _, pendente := range []string{"plans/sem-padrao.md", "prompts/prompt.md"} {
		if !strings.Contains(texto, pendente) {
			t.Errorf("pendente sem dono virou pendente silencioso: %s", pendente)
		}
	}
	if n := strings.Count(texto, "### TICK-INV-"); n != 2 {
		t.Errorf("seção TICKETS com %d tickets; esperado 2", n)
	}
	if n := strings.Count(texto, "has_binary_criteria"); n != 2 {
		t.Errorf("tickets sem critério binário: %d de 2", n)
	}
	for _, comDono := range []string{"INTENT_X.md", "reports/historico.md"} {
		if strings.Contains(texto, "TICK-INV-"+comDono) {
			t.Errorf("arquivo com dono virou ticket: %s", comDono)
		}
	}
	pendentes := PendentesSemDono(entradas)
	if len(pendentes) != 2 {
		t.Errorf("PendentesSemDono = %d entradas; esperado 2", len(pendentes))
	}
}
