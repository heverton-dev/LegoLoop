package orchestrix

import (
	"context"
	"strings"
	"testing"
)

func TestAlarmeMarcaParReprovadoEmExecucao(t *testing.T) {
	repo := newRepoT(t)
	o, no := nodeOutroHarness("homologado")
	o.Reprovados = &RegistroReprovados{}
	var avisos []string
	o.AoReprovar = func(harness, via, motivo string) {
		avisos = append(avisos, harness+"|"+via+"|"+motivo)
	}
	primeiro := no("n1")
	primeiro.SpecRef = "Terms of Service"
	_, err := o.Execute(context.Background(), primeiro, repo)
	if err == nil || !strings.Contains(err.Error(), "banimento") {
		t.Fatalf("o nó deveria falhar por banimento: %v", err)
	}
	if len(avisos) != 1 || !strings.HasPrefix(avisos[0], "omp|oficial|banimento:") {
		t.Fatalf("AoReprovar = %v", avisos)
	}
	if _, ok := o.Reprovados.Motivo("omp", "oficial"); !ok {
		t.Fatal("o par deveria constar como reprovado")
	}
	_, err = o.Execute(context.Background(), no("n2"), repo)
	if err == nil || !strings.Contains(err.Error(), "recusado") {
		t.Fatalf("o segundo nó deveria ser recusado: %v", err)
	}
	if len(avisos) != 1 {
		t.Fatalf("a recusa não deve reprovar de novo: %v", avisos)
	}
}

func TestEstadoEfetivoSemRegistro(t *testing.T) {
	o := &Orchestrix{}
	if got := o.estadoEfetivo("omp", "oficial"); got != "" {
		t.Fatalf("estado = %q", got)
	}
	o.EstadoHarness = func(string, string) string { return "homologado" }
	if got := o.estadoEfetivo("omp", "oficial"); got != "homologado" {
		t.Fatalf("estado = %q", got)
	}
	o.Reprovados = &RegistroReprovados{}
	o.Reprovados.Marcar("omp", "oficial", "banimento: x")
	if got := o.estadoEfetivo("omp", "oficial"); got != "reprovado" {
		t.Fatalf("estado = %q", got)
	}
	if got := o.estadoEfetivo("omp", "9router"); got != "homologado" {
		t.Fatalf("outro par = %q", got)
	}
}
