package studio

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/progresso"
)

func dagComProgresso(campo string) string {
	return `{"dag": {"name": "p", "version": "2.0.0", "trigger": {"type": "manual"}, "progresso": "` + campo + `",
	"nodes": [{"id": "TICK-01.1", "type": "script", "command": ["go", "version"],
	"routing": {"mode": "manual", "model": "none"}, "output": {"artifact": "res.txt", "gate": "true"}}]}, "allow_scripts": true}`
}

func progressoDoStudio(t *testing.T, projeto string) string {
	t.Helper()
	arq := filepath.Join(projeto, `docs`, `progress`, `x`, `progresso.json`)
	es := []progresso.EntradaTicket{{ID: `TICK-I99-01.1`, Title: `t`, EpicID: `E01`, InitiativeID: `I99`, ValidationCmd: `go version`}}
	p := progresso.NovoDeTickets(`I99`, `x`, es, time.Now().UTC().Format(time.RFC3339))
	if err := progresso.Salvar(arq, progresso.CaminhoEspelhoMD(arq), p); err != nil {
		t.Fatal(err)
	}
	return arq
}

func TestRunAPIAtualizaProgressoVivo(t *testing.T) {
	s, _ := newServer(t)
	arq := progressoDoStudio(t, s.ProjectDir)
	rec := do(t, s.Handler(), `POST`, `/api/run`, dagComProgresso(`docs/progress/x/progresso.json`))
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), `progresso_error`) {
		t.Fatalf(`status %d: %s`, rec.Code, rec.Body.String())
	}
	p, err := progresso.Carregar(arq)
	if err != nil {
		t.Fatal(err)
	}
	if it := p.Tickets[0]; it.Estado != progresso.EstadoAprovado || it.Responsavel != `studio` {
		t.Fatalf(`ticket não ficou vivo: %+v`, it)
	}
}

func TestRunAPIRecusaProgressoForaDoProjeto(t *testing.T) {
	s, _ := newServer(t)
	for _, campo := range []string{`../fora.json`, `docs/progress/nao-existe.json`, `docs/p.txt`} {
		rec := do(t, s.Handler(), `POST`, `/api/run`, dagComProgresso(campo))
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), `progresso vivo`) {
			t.Errorf(`campo %q: status %d: %s`, campo, rec.Code, rec.Body.String())
		}
	}
}
