package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"legoloop/internal/orchestrix"
	"legoloop/internal/roteamento"
	"legoloop/internal/roteamento/homologacao"
)

const (
	chamadasPorPar = 3
	memoriaMinMB   = 750
)

type HomologarDeps struct {
	Harness      func(harness, via, model string) homologacao.Harness
	MemoriaLivre func() (uint64, error)
}

type opcoesHomologar struct {
	pares         string
	yes           bool
	naoInterativo bool
	asJSON        bool
	teto          time.Duration
}

type fabricaHarness func(harness, via, model string) homologacao.Harness

func cmdHomologar(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("homologar", d)
	o := opcoesHomologar{}
	fs.StringVar(&o.pares, "pares", "", "pares harness|via separados por vírgula (vazio = todos instalados)")
	fs.BoolVar(&o.yes, "yes", false, "confirma sem perguntar")
	fs.BoolVar(&o.naoInterativo, "non-interactive", false, "não pergunta: pares sem entrada ficam pendentes")
	fs.DurationVar(&o.teto, "teto", 120*time.Second, "duração máxima por chamada")
	fs.BoolVar(&o.asJSON, "json", false, "saída JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	return executarHomologar(ctx, d, o)
}

func homologarItens(inv []roteamento.ItemInv, pares string) []roteamento.ItemInv {
	pedidos := map[string]bool{}
	for _, p := range strings.Split(pares, ",") {
		if p = strings.TrimSpace(p); p != "" {
			pedidos[p] = true
		}
	}
	var out []roteamento.ItemInv
	for _, it := range inv {
		if it.Estado == roteamento.EstadoNaoInstalado {
			continue
		}
		if len(pedidos) == 0 || pedidos[it.Par.Chave()] {
			out = append(out, it)
		}
	}
	return out
}

func homologarMemoria(h *HomologarDeps) error {
	fn := orchestrix.MemoriaLivreSO
	if h != nil && h.MemoriaLivre != nil {
		fn = h.MemoriaLivre
	}
	livre, err := fn()
	if err != nil {
		return err
	}
	if livre < uint64(memoriaMinMB)<<20 {
		return fmt.Errorf("memória livre insuficiente: %d MB, mínimo %d MB", livre>>20, memoriaMinMB)
	}
	return nil
}

func homologarConfirmou(d Deps, o opcoesHomologar) bool {
	if o.yes {
		return true
	}
	if o.naoInterativo || d.Stdin == nil {
		return false
	}
	fmt.Fprint(d.Stdout, "Confirma? [s/N] ")
	linha, _ := bufio.NewReader(d.Stdin).ReadString('\n')
	r := strings.ToLower(strings.TrimSpace(linha))
	return r == "s" || r == "sim"
}

func homologarFabrica(ctx context.Context, d Deps) (fabricaHarness, error) {
	if d.Homologar != nil && d.Homologar.Harness != nil {
		return d.Homologar.Harness, nil
	}
	ox, err := newOrchestrix(ctx, d, "", nil)
	if err != nil {
		return nil, err
	}
	return adaptadorReal(ox), nil
}

func homologarHome(d Deps) string {
	if d.Home != "" {
		return d.Home
	}
	home, _ := os.UserHomeDir()
	return home
}

func executarHomologar(ctx context.Context, d Deps, o opcoesHomologar) int {
	amb, err := rotaAmbiente(ctx, d, "")
	if err != nil {
		return fail(d, err)
	}
	inv, err := amb.Inventario(ctx)
	if err != nil {
		return fail(d, err)
	}
	itens := homologarItens(inv, o.pares)
	chamadas := chamadasPorPar * len(itens)
	fmt.Fprintf(d.Stdout, "Teto total estimado: %d chamadas, até %s\n", chamadas, time.Duration(chamadas)*o.teto)
	chaves := make([]string, 0, len(itens))
	for _, it := range itens {
		chaves = append(chaves, it.Par.Chave())
	}
	fmt.Fprintln(d.Stdout, "Aviso: "+homologacao.AvisoGasto(chaves))
	if err := homologarMemoria(d.Homologar); err != nil {
		return fail(d, err)
	}
	caminho := homologacao.Caminho(homologarHome(d))
	store, err := homologacao.Ler(caminho)
	if err != nil {
		return fail(d, err)
	}
	if !homologarConfirmou(d, o) {
		return homologarPendentes(d, itens, store, caminho)
	}
	return homologarRodar(ctx, d, o, itens, store, caminho)
}

func homologarPendentes(d Deps, itens []roteamento.ItemInv, store homologacao.Store, caminho string) int {
	for _, it := range itens {
		chave := it.Par.Chave()
		if _, ok := store[chave]; !ok {
			store[chave] = homologacao.Entrada{Estado: homologacao.EstadoPendente, Criterios: map[string]bool{}}
		}
	}
	if err := store.Gravar(caminho); err != nil {
		return fail(d, err)
	}
	fmt.Fprintln(d.Stdout, "Sem confirmação: nenhuma chamada feita, pares sem entrada ficaram pendentes.")
	return 0
}

func homologarRodar(ctx context.Context, d Deps, o opcoesHomologar, itens []roteamento.ItemInv, store homologacao.Store, caminho string) int {
	fabrica, err := homologarFabrica(ctx, d)
	if err != nil {
		return fail(d, err)
	}
	linhas := []map[string]any{}
	for _, it := range itens {
		e, err := homologarPar(ctx, fabrica, it, o.teto)
		if err != nil {
			return fail(d, err)
		}
		chave := it.Par.Chave()
		store[chave] = e
		linhas = append(linhas, map[string]any{"par": chave, "estado": e.Estado, "peso_mb": e.PesoMB, "motivo": e.Motivo})
	}
	if err := store.Gravar(caminho); err != nil {
		return fail(d, err)
	}
	return homologarSaida(d, linhas, o.asJSON)
}

func homologarSaida(d Deps, linhas []map[string]any, asJSON bool) int {
	if asJSON {
		return printJSON(d, linhas)
	}
	w := table(d.Stdout)
	fmt.Fprintln(w, "PAR\tESTADO\tPESO_MB\tMOTIVO")
	for _, l := range linhas {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", l["par"], l["estado"], l["peso_mb"], l["motivo"])
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	return 0
}

func homologarPar(ctx context.Context, fabrica fabricaHarness, it roteamento.ItemInv, teto time.Duration) (homologacao.Entrada, error) {
	if len(it.Modelos) == 0 {
		return homologacao.Entrada{
			Estado: homologacao.EstadoReprovado, Data: time.Now().Format("2006-01-02"),
			Criterios: map[string]bool{}, Motivo: "sem modelo no inventário",
		}, nil
	}
	base, err := os.MkdirTemp("", "homologar-")
	if err != nil {
		return homologacao.Entrada{}, err
	}
	defer os.RemoveAll(base)
	r := homologacao.Avaliar(ctx, fabrica(it.Par.Harness, string(it.Par.Via), it.Modelos[0]), base, teto)
	return homologacao.Entrada{Estado: r.Estado, Data: r.Data, PesoMB: r.PesoMB, Criterios: r.Criterios, Motivo: r.Motivo}, nil
}
