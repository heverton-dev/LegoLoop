package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
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

	sondaProvider    bool
	corrigirProvider bool
}

type ProviderDeps struct {
	HTTP     *http.Client
	Arquivos []string
	Esquema  string
	Canonico string
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
	fs.BoolVar(&o.sondaProvider, "provider", false, "sonda o host de cada harness (resolve, versão, combos)")
	fs.BoolVar(&o.corrigirProvider, "corrigir-provider", false, "troca o host que não resolve ou está antigo pelo canônico")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if o.sondaProvider || o.corrigirProvider {
		return executarProvider(ctx, d, o)
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

type sondador struct {
	ctx     context.Context
	cli     *http.Client
	esquema string
	canon   string
	cache   map[string]homologacao.Sonda
}

func (s *sondador) sondar(host string) homologacao.Sonda {
	if r, ok := s.cache[host]; ok {
		return r
	}
	chave := ""
	if host == s.canon {
		chave = provider.TrustedRouterKey(s.esquema + "://" + host)
	}
	r := homologacao.SondarHost(s.ctx, s.cli, s.esquema, host, chave)
	if host == s.canon {
		r.Comparar(r.Versao)
	} else {
		r.Comparar(s.sondar(s.canon).Versao)
	}
	s.cache[host] = r
	return r
}

func providerArquivos(d Deps) []string {
	if d.Provider != nil && d.Provider.Arquivos != nil {
		return d.Provider.Arquivos
	}
	h := homologarHome(d)
	return []string{
		filepath.Join(h, ".config", "opencode", "opencode.json"),
		filepath.Join(h, ".config", "mimocode", "mimocode.jsonc"),
		filepath.Join(h, ".omp", "agent", "models.yml"),
	}
}

func novoSondador(ctx context.Context, d Deps) *sondador {
	s := &sondador{
		ctx: ctx, cli: &http.Client{Timeout: 10 * time.Second},
		esquema: "https", canon: homologacao.HostCanonico, cache: map[string]homologacao.Sonda{},
	}
	if p := d.Provider; p != nil {
		if p.HTTP != nil {
			s.cli = p.HTTP
		}
		if p.Esquema != "" {
			s.esquema = p.Esquema
		}
		if p.Canonico != "" {
			s.canon = p.Canonico
		}
	}
	return s
}

func executarProvider(ctx context.Context, d Deps, o opcoesHomologar) int {
	s := novoSondador(ctx, d)
	var ends []homologacao.Endereco
	for _, arq := range providerArquivos(d) {
		e, err := homologacao.LerEnderecos(arq)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(d.Stderr, "aviso: %v\n", err)
			}
			continue
		}
		ends = append(ends, e...)
	}
	if o.sondaProvider {
		providerImprimir(d, s, ends)
	}
	if !o.corrigirProvider {
		return 0
	}
	return providerCorrigir(d, o, s, ends)
}

func providerImprimir(d Deps, s *sondador, ends []homologacao.Endereco) {
	w := table(d.Stdout)
	fmt.Fprintln(w, "ARQUIVO\tHOST\tRESOLVE\tVERSAO\tDISPONIVEL\tCOMBOS")
	for _, e := range ends {
		r := s.sondar(e.Host)
		fmt.Fprintf(w, "%s\t%s\t%t\t%s\t%s\t%s\n", e.Arquivo, e.Host, r.Resolve, r.Versao, r.VersaoDisponivel, strings.Join(r.Combos, ","))
	}
	_ = w.Flush()
}

func providerPlano(s *sondador, ends []homologacao.Endereco) []homologacao.Endereco {
	vistos := map[string]bool{}
	var plano []homologacao.Endereco
	for _, e := range ends {
		r := s.sondar(e.Host)
		chave := e.Arquivo + "|" + e.Host
		if e.Host != s.canon && (!r.Resolve || r.Antigo) && !vistos[chave] {
			vistos[chave] = true
			plano = append(plano, e)
		}
	}
	return plano
}

func providerCorrigir(d Deps, o opcoesHomologar, s *sondador, ends []homologacao.Endereco) int {
	plano := providerPlano(s, ends)
	if len(plano) == 0 {
		fmt.Fprintln(d.Stdout, "Nenhum endereço precisa de correção.")
		return 0
	}
	for _, e := range plano {
		fmt.Fprintf(d.Stdout, "Plano: %s: %s -> %s\n", e.Arquivo, e.Host, s.canon)
	}
	if !homologarConfirmou(d, o) {
		fmt.Fprintln(d.Stdout, "Sem confirmação: nenhum arquivo foi alterado.")
		return 0
	}
	for _, e := range plano {
		m, err := homologacao.Corrigir(e.Arquivo, e.Host, s.canon, true)
		if err != nil {
			return fail(d, err)
		}
		fmt.Fprintf(d.Stdout, "Corrigido: %s conferido=%t backup=%s\n", m.Arquivo, m.Conferida, m.Backup)
	}
	return 0
}
