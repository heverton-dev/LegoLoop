package cli

import (
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/progresso"
)

// TestFabricarPlanoPreservaTicketAprovado cobre o A21: regerar o plano nao pode zerar ticket PASSED.
func TestFabricarPlanoPreservaTicketAprovado(t *testing.T) {
	raiz, spec, ticketsFile := prepararAmbienteFabricarPlano(t, ticketsDoFabricarPlano)
	args := []string{"plano", "--raiz", raiz, "--spec", spec, "--tickets", ticketsFile, "--data", "2026-11-21"}
	if code, saida := rodarFabricar(args...); code != 0 {
		t.Fatalf("primeira geracao (%d): %s", code, saida)
	}
	caminho := filepath.Join(raiz, "docs", "progress", "2026-11-21_plan-teste", "progresso.json")
	p, err := progresso.Carregar(caminho)
	if err != nil || len(p.Tickets) == 0 {
		t.Fatalf("carregar: %v", err)
	}
	p.Tickets[0].Estado, p.Tickets[0].Prova = progresso.EstadoAprovado, "go test ok"
	p.RecalcularResumo()
	if err := progresso.Salvar(caminho, progresso.CaminhoEspelhoMD(caminho), p); err != nil {
		t.Fatal(err)
	}
	if code, saida := rodarFabricar(args...); code != 0 {
		t.Fatalf("segunda geracao (%d): %s", code, saida)
	}
	depois, err := progresso.Carregar(caminho)
	if err != nil {
		t.Fatal(err)
	}
	if depois.Tickets[0].Estado != progresso.EstadoAprovado || depois.Tickets[0].Prova != "go test ok" || depois.Resumo.Aprovados != 1 {
		t.Fatalf("ticket aprovado foi zerado: %+v %+v", depois.Tickets[0], depois.Resumo)
	}
}
