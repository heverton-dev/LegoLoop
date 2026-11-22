package cli

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/progresso"
)

func progressoSwarmDeTeste(t *testing.T, ids ...string) string {
	t.Helper()
	arq := filepath.Join(t.TempDir(), `progresso.json`)
	var es []progresso.EntradaTicket
	for _, id := range ids {
		es = append(es, progresso.EntradaTicket{ID: id, Title: id, EpicID: `E01`, InitiativeID: `I99`, ValidationCmd: `go version`})
	}
	p := progresso.NovoDeTickets(`I99`, `docs/plans/x.md`, es, time.Now().UTC().Format(time.RFC3339))
	if err := progresso.Salvar(arq, progresso.CaminhoEspelhoMD(arq), p); err != nil {
		t.Fatal(err)
	}
	return arq
}

func estadosDe(t *testing.T, arq string) map[string]string {
	t.Helper()
	p, err := progresso.Carregar(arq)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]string{}
	for _, it := range p.Tickets {
		m[it.ID] = it.Estado
	}
	return m
}

// TestSwarmRunMantemProgressoVivo: cada worker aprovado fecha o Ticket do nó e o espelho .md acompanha.
func TestSwarmRunMantemProgressoVivo(t *testing.T) {
	if testing.Short() {
		t.Skip(`E2E real sobe processos`)
	}
	repo := e2eGitRepo(t)
	dagPath := writeRunDAG(t, t.TempDir())
	arq := progressoSwarmDeTeste(t, `a`, `b`, `c`)
	code, joined, _ := runSwarmRun(t, repo, dagPath, ``, `-progresso`, arq)
	if code != 0 {
		t.Fatalf(`swarm run saiu %d: %s`, code, strings.ReplaceAll(joined, "\x00", ` | `))
	}
	for id, estado := range estadosDe(t, arq) {
		if estado != progresso.EstadoAprovado {
			t.Errorf(`ticket %s = %s, esperado PASSED`, id, estado)
		}
	}
	d, out, _ := deps(t)
	if c := run(t, d, `progresso`, `verificar`, `--arquivo`, arq); c != 0 {
		t.Fatalf(`verificar code=%d out=%s`, c, out.String())
	}
}

// TestSwarmRunFalhaMarcaProgresso: worker com código 7 reprova o Ticket; o nível seguinte não começa; nada fica IN_PROGRESS.
func TestSwarmRunFalhaMarcaProgresso(t *testing.T) {
	if testing.Short() {
		t.Skip(`E2E real sobe processos`)
	}
	repo := e2eGitRepo(t)
	dagPath := writeRunDAG(t, t.TempDir())
	arq := progressoSwarmDeTeste(t, `a`, `b`, `c`)
	if code, _, _ := runSwarmRun(t, repo, dagPath, `b`, `-progresso`, arq); code != 1 {
		t.Fatalf(`esperava código 1, veio %d`, code)
	}
	est := estadosDe(t, arq)
	if est[`b`] != progresso.EstadoReprovado || est[`c`] != progresso.EstadoPendente {
		t.Fatalf(`estados = %v`, est)
	}
	for id, estado := range est {
		if estado == progresso.EstadoEmAndamento {
			t.Errorf(`ticket %s ficou IN_PROGRESS (débito de estado)`, id)
		}
	}
}
