package progresso

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

const agoraTeste = `2026-11-22T12:00:00Z`

func arquivoDeTeste(t *testing.T, entradas []EntradaTicket) string {
	t.Helper()
	caminho := filepath.Join(t.TempDir(), `docs`, `progress`, `p`, `progresso.json`)
	p := NovoDeTickets(`I99`, `docs/plans/x.md`, entradas, agoraTeste)
	if err := Salvar(caminho, CaminhoEspelhoMD(caminho), p); err != nil {
		t.Fatal(err)
	}
	return caminho
}

func entradasDeTeste(n int, comando string) []EntradaTicket {
	var es []EntradaTicket
	for i := 1; i <= n; i++ {
		es = append(es, EntradaTicket{ID: fmt.Sprintf(`TICK-I99-01.%d`, i), Title: `t`, EpicID: `E01`, InitiativeID: `I99`, ValidationCmd: comando})
	}
	return es
}

func execFixo(code int) Executor {
	return func(context.Context, string, string) (int, string, error) { return code, `saida`, nil }
}

func TestAtualizarConcorrenteNaoPerdeEscrita(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(12, `true`))
	var wg sync.WaitGroup
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			err := Atualizar(arq, func(p *Progresso) error {
				return Iniciar(p, id, `w`, `n`, agoraTeste)
			})
			if err != nil {
				t.Error(err)
			}
		}(fmt.Sprintf(`TICK-I99-01.%d`, i))
	}
	wg.Wait()
	p, err := Carregar(arq)
	if err != nil {
		t.Fatal(err)
	}
	if p.Resumo.EmAndamento != 12 {
		t.Fatalf(`em andamento = %d, esperado 12`, p.Resumo.EmAndamento)
	}
}

func TestAtualizarSemGravarQuandoFnFalha(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(1, `true`))
	err := Atualizar(arq, func(p *Progresso) error {
		p.Tickets[0].Estado = EstadoAprovado
		return fmt.Errorf(`falhou`)
	})
	if err == nil {
		t.Fatal(`esperava erro`)
	}
	p, _ := Carregar(arq)
	if p.Tickets[0].Estado != EstadoPendente {
		t.Fatalf(`estado gravado apesar do erro: %s`, p.Tickets[0].Estado)
	}
}

func TestTravaObsoletaEhRemovida(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(1, `true`))
	trava := arq + `.lock`
	if err := os.WriteFile(trava, []byte(`1`), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Minute)
	if err := os.Chtimes(trava, old, old); err != nil {
		t.Fatal(err)
	}
	if err := Atualizar(arq, func(*Progresso) error { return nil }); err != nil {
		t.Fatal(err)
	}
}

func TestVerificarCoerenteNaoAchaNada(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(2, `true`))
	p, _ := Carregar(arq)
	md, _ := os.ReadFile(CaminhoEspelhoMD(arq))
	texto := string(md)
	if a := Verificar(p, OpcoesVerificar{Agora: time.Now(), IdadeMaxima: time.Hour, Espelho: &texto}); len(a) != 0 {
		t.Fatalf(`achados inesperados: %+v`, a)
	}
}

func regrasDe(achados []Achado) map[string]bool {
	m := map[string]bool{}
	for _, a := range achados {
		m[a.Regra] = true
	}
	return m
}

func TestVerificarDetectaDivergencias(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(3, `true`))
	p, _ := Carregar(arq)
	p.Tickets[0].Estado = EstadoAprovado
	p.Tickets[0].ExitCode = 1
	p.Tickets[1].Estado = EstadoEmAndamento
	p.Tickets[1].AtualizadoEm = `2026-11-20T00:00:00Z`
	p.Tickets[2].Comando = ``
	antigo := `outro texto`
	achados := Verificar(p, OpcoesVerificar{Agora: time.Date(2026, 11, 22, 12, 0, 0, 0, time.UTC), IdadeMaxima: 24 * time.Hour, Espelho: &antigo})
	regras := regrasDe(achados)
	for _, esperada := range []string{`resumo-divergente`, `aprovado-sem-exit-zero`, `orfao`, `sem-comando`, `espelho-divergente`} {
		if !regras[esperada] {
			t.Errorf(`regra %s ausente em %+v`, esperada, achados)
		}
	}
}

func TestVerificarIgnoraQuebraDeLinhaDoWindows(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(1, `true`))
	p, _ := Carregar(arq)
	md, _ := os.ReadFile(CaminhoEspelhoMD(arq))
	crlf := normalizarQuebras(string(md))
	crlf = replaceTodos(crlf, "\n", "\r\n")
	if a := Verificar(p, OpcoesVerificar{Espelho: &crlf}); len(a) != 0 {
		t.Fatalf(`CRLF não devia divergir: %+v`, a)
	}
}

func replaceTodos(s, de, para string) string {
	out := ``
	for _, r := range s {
		if string(r) == de {
			out += para
		} else {
			out += string(r)
		}
	}
	return out
}

func TestSondarApontaGateVerdeSemAprovacao(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(2, `true`))
	p, _ := Carregar(arq)
	p.Tickets[1].Estado = EstadoAprovado
	achados := Sondar(context.Background(), p, `.`, execFixo(0))
	if len(achados) != 1 || achados[0].Ticket != `TICK-I99-01.1` || achados[0].Regra != `gate-verde-sem-aprovacao` {
		t.Fatalf(`achados = %+v`, achados)
	}
	if a := Sondar(context.Background(), p, `.`, execFixo(1)); len(a) != 0 {
		t.Fatalf(`gate vermelho não devia gerar achado: %+v`, a)
	}
}

func TestReconciliarOrfaos(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(2, `true`))
	p, _ := Carregar(arq)
	for i := range p.Tickets {
		p.Tickets[i].Estado = EstadoEmAndamento
		p.Tickets[i].Responsavel = `sessao-caiu`
		p.Tickets[i].AtualizadoEm = `2026-11-20T00:00:00Z`
	}
	agora := time.Date(2026, 11, 22, 12, 0, 0, 0, time.UTC)
	orfaos := Orfaos(p, agora, 24*time.Hour)
	if len(orfaos) != 2 {
		t.Fatalf(`órfãos = %v`, orfaos)
	}
	exec := func(_ context.Context, _ string, _ string) (int, string, error) { return 0, `ok`, nil }
	ReconciliarOrfaos(context.Background(), &p, `.`, agoraTeste, orfaos[:1], exec)
	ReconciliarOrfaos(context.Background(), &p, `.`, agoraTeste, orfaos[1:], execFixo(1))
	if p.Tickets[0].Estado != EstadoAprovado {
		t.Errorf(`gate verde devia aprovar: %s`, p.Tickets[0].Estado)
	}
	if p.Tickets[1].Estado != EstadoPendente || p.Tickets[1].Responsavel != `` {
		t.Errorf(`gate vermelho devia reabrir: %+v`, p.Tickets[1])
	}
}

func TestIDCanonicoEResolverTicket(t *testing.T) {
	casos := map[string]string{`TICK-I4-01.1`: `TICK-01.1`, `TICK-01.1`: `TICK-01.1`, `e01`: `e01`}
	for in, want := range casos {
		if got := IDCanonico(in); got != want {
			t.Errorf(`IDCanonico(%q) = %q, esperado %q`, in, got, want)
		}
	}
	p := NovoDeTickets(`I99`, `x`, entradasDeTeste(1, `true`), agoraTeste)
	if id, ok := ResolverTicket(p, `TICK-01.1`); !ok || id != `TICK-I99-01.1` {
		t.Fatalf(`resolver = %q %v`, id, ok)
	}
	if _, ok := ResolverTicket(p, `TICK-09.9`); ok {
		t.Fatal(`nó desconhecido não devia resolver`)
	}
}

func TestObservadorMantemProgressoVivo(t *testing.T) {
	es := entradasDeTeste(3, `verifica`)
	es[2].ValidationCmd = ``
	arq := arquivoDeTeste(t, es)
	o := NovoObservador(context.Background(), arq, `.`, `runner`)
	o.exec = execFixo(0)
	o.NoInicio(`TICK-01.1`)
	p, _ := Carregar(arq)
	if p.Tickets[0].Estado != EstadoEmAndamento || p.Tickets[0].Responsavel != `runner` {
		t.Fatalf(`início não registrado: %+v`, p.Tickets[0])
	}
	o.NoFim(`TICK-01.1`, true, 0, ``)
	o.NoFim(`TICK-01.2`, false, 2, `gate reprovou`)
	o.NoFim(`TICK-01.3`, true, 0, ``)
	o.NoFim(`no-sem-ticket`, true, 0, ``)
	p, _ = Carregar(arq)
	estados := []string{p.Tickets[0].Estado, p.Tickets[1].Estado, p.Tickets[2].Estado}
	if estados[0] != EstadoAprovado || estados[1] != EstadoReprovado || estados[2] != EstadoAprovado {
		t.Fatalf(`estados = %v`, estados)
	}
	if len(o.Erros()) != 0 {
		t.Fatalf(`erros inesperados: %v`, o.Erros())
	}
	md, _ := os.ReadFile(CaminhoEspelhoMD(arq))
	texto := string(md)
	regras := regrasDe(Verificar(p, OpcoesVerificar{Espelho: &texto}))
	if regras[`espelho-divergente`] || regras[`resumo-divergente`] {
		t.Fatalf(`o espelho e o resumo devem acompanhar o JSON: %v`, regras)
	}
}

func TestObservadorNoFimReprovaQuandoComandoFalha(t *testing.T) {
	arq := arquivoDeTeste(t, entradasDeTeste(1, `verifica`))
	o := NovoObservador(context.Background(), arq, `.`, `swarm`)
	o.exec = execFixo(3)
	o.NoFim(`TICK-01.1`, true, 0, ``)
	p, _ := Carregar(arq)
	if p.Tickets[0].Estado != EstadoReprovado || p.Tickets[0].ExitCode != 3 {
		t.Fatalf(`ticket = %+v`, p.Tickets[0])
	}
}

func TestObservadorRegistraErroDeArquivoAusente(t *testing.T) {
	o := NovoObservador(context.Background(), filepath.Join(t.TempDir(), `nao-existe.json`), `.`, `runner`)
	o.NoInicio(`TICK-01.1`)
	if len(o.Erros()) != 1 {
		t.Fatalf(`erros = %v`, o.Erros())
	}
}

// TestProgressoDoRepositorioEstaCoerente é o portão: todo progresso.json em docs/progress precisa ser coerente.
func TestProgressoDoRepositorioEstaCoerente(t *testing.T) {
	arquivos, err := filepath.Glob(filepath.Join(`..`, `..`, `docs`, `progress`, `*`, `progresso.json`))
	if err != nil {
		t.Fatal(err)
	}
	for _, arq := range arquivos {
		p, err := Carregar(arq)
		if err != nil {
			t.Errorf(`%s: %v`, arq, err)
			continue
		}
		md, errMD := os.ReadFile(CaminhoEspelhoMD(arq))
		if errMD != nil {
			t.Errorf(`%s: espelho .md ausente`, arq)
			continue
		}
		texto := string(md)
		for _, a := range Verificar(p, OpcoesVerificar{Espelho: &texto}) {
			t.Errorf(`%s: %s %s %s`, arq, a.Regra, a.Ticket, a.Motivo)
		}
	}
}
