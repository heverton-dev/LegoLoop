package progresso

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNovoDeTicketsEValidar(t *testing.T) {
	entradas := []EntradaTicket{
		{ID: "TICK-01.1", EpicID: "E01", Title: "Primeiro", ValidationCmd: "go test ./...", TargetFiles: []string{"a.go"}, SpecCriteria: []string{"AC-1"}},
		{ID: "TICK-01.2", EpicID: "E01", Title: "Segundo", ValidationCmd: "go vet ./...", TargetFiles: []string{"b.go"}, SpecCriteria: []string{"AC-2"}},
	}
	p := NovoDeTickets("I04", "docs/plans/PLAN.md", entradas, "2026-11-22T11:00:00Z")
	if err := p.Validar(); err != nil {
		t.Fatalf("progresso valido recusado: %v", err)
	}
	if p.Resumo.Total != 2 || p.Resumo.Pendentes != 2 || p.Resumo.Aprovados != 0 {
		t.Fatalf("resumo incorreto: %+v", p.Resumo)
	}
	invalido := p
	invalido.Schema = "outro"
	if err := invalido.Validar(); err == nil {
		t.Fatal("esperado erro para schema invalido")
	}
}

func TestIniciarEValidarPromoveApenasComExitZero(t *testing.T) {
	p := NovoDeTickets("I04", "PLAN.md", []EntradaTicket{
		{ID: "TICK-01.1", EpicID: "E01", Title: "Teste", ValidationCmd: "go test ./a"},
		{ID: "TICK-01.2", EpicID: "E01", Title: "Falha", ValidationCmd: "go test ./b"},
	}, "2026-11-22T11:00:00Z")

	if err := Iniciar(&p, "TICK-01.1", "worker-1", "iniciando TDD", "2026-11-22T11:01:00Z"); err != nil {
		t.Fatalf("Iniciar falhou: %v", err)
	}
	if p.Tickets[0].Estado != EstadoEmAndamento || p.Resumo.EmAndamento != 1 {
		t.Fatalf("estado esperado IN_PROGRESS, obtido %+v", p.Tickets[0])
	}

	okExec := func(_ context.Context, _, _ string) (int, string, error) {
		return 0, "ok legoloop/a 0.1s", nil
	}
	passou, err := ValidarTicket(context.Background(), &p, "TICK-01.1", ".", "2026-11-22T11:02:00Z", okExec)
	if err != nil || !passou {
		t.Fatalf("esperado PASSED, obtido passou=%v err=%v", passou, err)
	}
	if p.Tickets[0].Estado != EstadoAprovado || p.Tickets[0].ExitCode != 0 || p.Resumo.Aprovados != 1 {
		t.Fatalf("ticket nao promovido para PASSED: %+v", p.Tickets[0])
	}

	failExec := func(_ context.Context, _, _ string) (int, string, error) {
		return 2, "FAIL legoloop/b", nil
	}
	passou, err = ValidarTicket(context.Background(), &p, "TICK-01.2", ".", "2026-11-22T11:03:00Z", failExec)
	if err != nil || passou {
		t.Fatalf("esperado FAILED, obtido passou=%v err=%v", passou, err)
	}
	if p.Tickets[1].Estado != EstadoReprovado || p.Tickets[1].ExitCode != 2 || p.Resumo.Reprovados != 1 {
		t.Fatalf("ticket deveria ficar FAILED: %+v", p.Tickets[1])
	}
}

func TestSalvarCarregarERenderizarMarkdown(t *testing.T) {
	dir := t.TempDir()
	arqJSON := filepath.Join(dir, "PROGRESSO.json")
	arqMD := filepath.Join(dir, "PROGRESSO.md")
	p := NovoDeTickets("I04", "PLAN.md", []EntradaTicket{
		{ID: "TICK-01.1", EpicID: "E01", Title: "Em curso", ValidationCmd: "go test ./..."},
	}, "2026-11-22T11:00:00Z")
	_ = Iniciar(&p, "TICK-01.1", "worker-a", "editando teste", "2026-11-22T11:05:00Z")

	if err := Salvar(arqJSON, arqMD, p); err != nil {
		t.Fatalf("Salvar falhou: %v", err)
	}
	carregado, err := Carregar(arqJSON)
	if err != nil {
		t.Fatalf("Carregar falhou: %v", err)
	}
	if carregado.Tickets[0].Responsavel != "worker-a" {
		t.Fatalf("responsavel perdido: %+v", carregado.Tickets[0])
	}
	mdBytes, err := os.ReadFile(arqMD)
	if err != nil {
		t.Fatalf("ler MD falhou: %v", err)
	}
	md := string(mdBytes)
	for _, trecho := range []string{"# Progresso Vivo: I04", "TICK-01.1", "worker-a", "editando teste"} {
		if !strings.Contains(md, trecho) {
			t.Fatalf("Markdown sem trecho %q:\n%s", trecho, md)
		}
	}
}

func TestTemplateProgressoJSONValida(t *testing.T) {
	caminho := filepath.Join("..", "..", "docs", "templates", "TEMPLATE_PROGRESSO.json")
	p, err := Carregar(caminho)
	if err != nil {
		t.Fatalf("template JSON invalido: %v", err)
	}
	if err := p.Validar(); err != nil {
		t.Fatalf("template nao passa em Validar: %v", err)
	}
}
