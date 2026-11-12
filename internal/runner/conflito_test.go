package runner

import (
	"context"
	"strings"
	"testing"
	"time"

	"legoloop/internal/dag"
)

func doisNos(t *testing.T, arquivoA, arquivoB string) *dag.DAG {
	t.Helper()
	mk := func(id, saida string) string {
		return `{"id":"` + id + `","type":"script","command":["x"],"context_out":["` + saida + `"],
		"routing":{"mode":"manual","model":"none"},
		"output":{"artifact":"` + id + `.md","gate":"file_exists('` + id + `.md') && exit_code == 0"}}`
	}
	d, err := dag.Parse([]byte(`{"name":"c","nodes":[` + strings.Join([]string{mk("a", arquivoA), mk("b", arquivoB)}, ",") + `]}`))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func maximoSimultaneo(t *testing.T, d *dag.DAG, opts Options) int32 {
	t.Helper()
	ex := &fakeExec{delay: 120 * time.Millisecond}
	opts.Workdir, opts.Executor, opts.Parallelism = t.TempDir(), ex, 4
	rep, err := Run(context.Background(), d, opts)
	if err != nil || !rep.OK {
		t.Fatalf("err = %v, rep = %+v", err, rep)
	}
	return ex.maxSeen.Load()
}

func TestDoisNosMesmoArquivoNaoJuntos(t *testing.T) {
	if got := maximoSimultaneo(t, doisNos(t, "src/x.go", "./src/x.go"), Options{}); got != 1 {
		t.Fatalf("nós com o mesmo arquivo rodaram juntos: máximo %d", got)
	}
	if got := maximoSimultaneo(t, doisNos(t, "src/x.go", "src/y.go"), Options{}); got != 2 {
		t.Fatalf("nós com arquivos distintos deveriam rodar juntos: máximo %d", got)
	}
}

func TestWorkersDaMemoriaLimitamOsNosJuntos(t *testing.T) {
	var pesosVistos [][]int
	limite := func(pesos []int) int {
		pesosVistos = append(pesosVistos, append([]int(nil), pesos...))
		return 1
	}
	if got := maximoSimultaneo(t, doisNos(t, "src/x.go", "src/y.go"), Options{Workers: limite, PesoMB: func(dag.Node) int { return 900 }}); got != 1 {
		t.Fatalf("a memória comporta 1 worker, mas rodaram %d juntos", got)
	}
	if len(pesosVistos) == 0 || pesosVistos[0][0] != 900 {
		t.Fatalf("o escalonador não consultou o limite com o peso do nó: %v", pesosVistos)
	}
}
