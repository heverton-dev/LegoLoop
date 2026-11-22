package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/caminhos"
	"legoloop/internal/progresso"
)

const dagComTicket = `{"name":"p","progresso":"docs/progress/x/progresso.json","nodes":[
 {"id":"TICK-01.1","type":"script","command":["go","version"],"timeout_sec":60,
  "routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"GO_VERSION.txt","gate":"exit_code == 0 && has_pattern('GO_VERSION.txt', 'go1\\.')"}}]}`

func progressoDeTeste(t *testing.T, work, comando string) string {
	t.Helper()
	arq := filepath.Join(work, `docs`, `progress`, `x`, `progresso.json`)
	es := []progresso.EntradaTicket{{ID: `TICK-I99-01.1`, Title: `t`, EpicID: `E01`, InitiativeID: `I99`, ValidationCmd: comando}}
	p := progresso.NovoDeTickets(`I99`, `docs/plans/x.md`, es, time.Now().UTC().Format(time.RFC3339))
	if err := progresso.Salvar(arq, progresso.CaminhoEspelhoMD(arq), p); err != nil {
		t.Fatal(err)
	}
	return arq
}

func ticketDe(t *testing.T, arq string) progresso.ItemTicket {
	t.Helper()
	p, err := progresso.Carregar(arq)
	if err != nil {
		t.Fatal(err)
	}
	return p.Tickets[0]
}

func TestRunAtualizaProgressoVivoPeloCampoDoDAG(t *testing.T) {
	work := t.TempDir()
	arq := progressoDeTeste(t, work, `go version`)
	p := caminhos.CaminhoDAG(work)
	os.WriteFile(p, []byte(dagComTicket), 0o644)
	d, out, errb := deps(t)
	if code := run(t, d, `run`, p, `--workdir`, work, `--allow-scripts`); code != 0 {
		t.Fatalf(`run code=%d out=%s err=%s`, code, out.String(), errb.String())
	}
	it := ticketDe(t, arq)
	if it.Estado != progresso.EstadoAprovado || it.ExitCode != 0 || it.Responsavel != `runner` {
		t.Fatalf(`ticket não ficou vivo: %+v`, it)
	}
	d2, out2, _ := deps(t)
	if code := run(t, d2, `progresso`, `verificar`, `--arquivo`, arq); code != 0 {
		t.Fatalf(`verificar code=%d out=%s`, code, out2.String())
	}
}

func TestRunMarcaTicketReprovadoQuandoGateFalha(t *testing.T) {
	work := t.TempDir()
	arq := progressoDeTeste(t, work, `go version`)
	p := caminhos.CaminhoDAG(work)
	os.WriteFile(p, []byte(strings.Replace(dagComTicket, `go1\\.`, `nunca-casa`, 1)), 0o644)
	d, _, _ := deps(t)
	if code := run(t, d, `run`, p, `--workdir`, work, `--allow-scripts`); code != 1 {
		t.Fatalf(`gate reprovado devia sair 1, code=%d`, code)
	}
	if it := ticketDe(t, arq); it.Estado != progresso.EstadoReprovado {
		t.Fatalf(`ticket devia estar FAILED: %+v`, it)
	}
}

func TestRunFalhaSeOArquivoDeProgressoNaoExiste(t *testing.T) {
	work := t.TempDir()
	p := caminhos.CaminhoDAG(work)
	os.WriteFile(p, []byte(dagComTicket), 0o644)
	d, _, errb := deps(t)
	if code := run(t, d, `run`, p, `--workdir`, work, `--allow-scripts`); code == 0 || !strings.Contains(errb.String(), `progresso vivo`) {
		t.Fatalf(`devia recusar progresso ausente, code=%d err=%s`, code, errb.String())
	}
}

func TestCaminhoProgressoVivo(t *testing.T) {
	abs := filepath.Join(t.TempDir(), `a.json`)
	casos := []struct{ base, flag, dag, want string }{
		{`w`, ``, ``, ``},
		{`w`, ``, `docs/p.json`, filepath.Join(`w`, `docs`, `p.json`)},
		{`w`, `f.json`, `docs/p.json`, filepath.Join(`w`, `f.json`)},
		{`w`, abs, ``, abs},
	}
	for _, c := range casos {
		got, err := caminhoProgressoVivo(c.base, c.flag, c.dag)
		if err != nil || got != c.want {
			t.Errorf(`caminhoProgressoVivo(%q,%q,%q) = %q, %v; esperado %q`, c.base, c.flag, c.dag, got, err, c.want)
		}
	}
}

// Path Jail: o campo `progresso` de um DAG de terceiros não escapa do projeto.
func TestCaminhoProgressoVivoRecusaCampoForaDoProjeto(t *testing.T) {
	abs := filepath.Join(t.TempDir(), `a.json`)
	for _, campo := range []string{`../fora.json`, `docs/../../fora.json`, abs, `docs/p.txt`} {
		if got, err := caminhoProgressoVivo(`w`, ``, campo); err == nil {
			t.Errorf(`campo %q devia ser recusado, veio %q`, campo, got)
		}
	}
}

func TestProgressoVerificarReprovaDivergencia(t *testing.T) {
	work := t.TempDir()
	arq := progressoDeTeste(t, work, `go version`)
	if err := os.WriteFile(progresso.CaminhoEspelhoMD(arq), []byte(`editado à mão`), 0o644); err != nil {
		t.Fatal(err)
	}
	d, out, _ := deps(t)
	if code := run(t, d, `progresso`, `verificar`, `--arquivo`, arq); code != 1 || !strings.Contains(out.String(), `espelho-divergente`) {
		t.Fatalf(`code=%d out=%s`, code, out.String())
	}
}

func TestProgressoVerificarReconciliaOrfao(t *testing.T) {
	work := t.TempDir()
	arq := progressoDeTeste(t, work, `go version`)
	old := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	if err := progresso.Atualizar(arq, func(p *progresso.Progresso) error {
		return progresso.Iniciar(p, `TICK-I99-01.1`, `sessao-caiu`, `n`, old)
	}); err != nil {
		t.Fatal(err)
	}
	d, out, _ := deps(t)
	if code := run(t, d, `progresso`, `verificar`, `--arquivo`, arq); code != 1 || !strings.Contains(out.String(), `orfao`) {
		t.Fatalf(`órfão devia reprovar, code=%d out=%s`, code, out.String())
	}
	d2, out2, _ := deps(t)
	if code := run(t, d2, `progresso`, `verificar`, `--arquivo`, arq, `--reconciliar`, `--raiz`, work); code != 0 {
		t.Fatalf(`reconciliar code=%d out=%s`, code, out2.String())
	}
	if it := ticketDe(t, arq); it.Estado != progresso.EstadoAprovado {
		t.Fatalf(`gate verde devia aprovar o órfão: %+v`, it)
	}
}
