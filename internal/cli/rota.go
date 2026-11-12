package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/processo"
	"legoloop/internal/roteamento"
)

const rotaTipoNo = "unit_task"

type RotaDeps struct {
	Inventario func(context.Context) ([]roteamento.ItemInv, error)
	Global     []byte
	Projeto    []byte
}

type rotaTicket struct {
	ID        string
	Titulo    string
	Natureza  string
	Arquivos  []string
	Depende   []string
	Validacao string
}

var (
	rotaCabecalho = regexp.MustCompile(`(?i)^###\s*\[(TICK(?:ET)?-[A-Za-z0-9_.-]+)\]\s*(.*)`)
	rotaNomePlano = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})_PLAN_(.+)$`)
)

func rotaCampo(linhas []string, nome string) string {
	re := regexp.MustCompile(`(?i)^\s*-\s*\*\*` + nome + `:\*\*\s*(.*)`)
	for _, l := range linhas {
		if m := re.FindStringSubmatch(l); len(m) > 1 {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}

func rotaLista(bruto string) []string {
	var out []string
	for _, p := range strings.Split(bruto, ",") {
		p = strings.Trim(strings.TrimSpace(p), "`")
		if p != "" && !strings.Contains(strings.ToLower(p), "nenhuma") {
			out = append(out, p)
		}
	}
	return out
}

func rotaTickets(plano string) []rotaTicket {
	var tickets []rotaTicket
	var atual *rotaTicket
	var linhas []string
	fechar := func() {
		if atual == nil {
			return
		}
		atual.Natureza = strings.ToLower(rotaCampo(linhas, "Natureza"))
		atual.Arquivos = rotaLista(rotaCampo(linhas, "Arquivos Afetados"))
		atual.Depende = rotaLista(rotaCampo(linhas, "Dependências"))
		atual.Validacao = rotaCampo(linhas, "Comando de Validação")
		tickets = append(tickets, *atual)
	}
	for _, l := range strings.Split(plano, "\n") {
		if m := rotaCabecalho.FindStringSubmatch(strings.TrimSpace(l)); len(m) > 2 {
			fechar()
			atual, linhas = &rotaTicket{ID: m[1], Titulo: strings.TrimSpace(m[2])}, nil
			continue
		}
		linhas = append(linhas, l)
	}
	fechar()
	return tickets
}

func rotaPacotes(arquivos []string) int {
	dirs := map[string]bool{}
	for _, a := range arquivos {
		dirs[path.Dir(filepath.ToSlash(a))] = true
	}
	return len(dirs)
}

func rotaNo(t rotaTicket, inv []roteamento.ItemInv, elegiveis map[string][]roteamento.CandidatoVia) (map[string]any, []roteamento.SubstituicaoVia, error) {
	res, err := roteamento.Classificar(roteamento.EntradaClasse{
		ArquivosSaida: len(t.Arquivos),
		Pacotes:       rotaPacotes(t.Arquivos),
		Destrutivo:    t.Natureza == roteamento.ClasseExtreme,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", t.ID, err)
	}
	r, err := roteamento.ResolverVia(t.ID, res.Classe, elegiveis, inv)
	if err != nil {
		return nil, nil, err
	}
	no := rotaNoBase(t, res)
	switch {
	case r.Usado.Harness != "":
		no["routing"] = map[string]any{"mode": "suggested", "harness": r.Usado.Harness, "via": string(r.Usado.Via), "model": r.Usado.Modelo, "effort": r.Usado.Esforco}
	case r.StopGate:
		no["type"] = "stop_gate"
		no["routing"] = map[string]any{"mode": "manual", "harness": "local", "via": string(roteamento.ViaOficial), "model": "human/operator", "effort": "low"}
	default:
		return nil, nil, fmt.Errorf("%s: nenhum candidato disponivel para a classe %s (politica: %s)", t.ID, res.Classe, r.Politica)
	}
	return no, r.Substituicoes, nil
}

func rotaNoBase(t rotaTicket, res roteamento.ResultadoClasse) map[string]any {
	verify := [][]string{}
	if f := strings.Fields(t.Validacao); len(f) > 0 {
		verify = append(verify, f)
	}
	no := map[string]any{
		"id":         t.ID,
		"label":      t.Titulo,
		"nature":     res.Classe,
		"class_rule": res.ClassRule,
		"type":       rotaTipoNo,
		"status":     "PENDING",
		"depends_on": append([]string{}, t.Depende...),
		"output":     map[string]any{"artifacts": append([]string{}, t.Arquivos...), "verify": verify, "gate": "exit_code == 0"},
	}
	if res.ClassOverrideReason != "" {
		no["class_override_reason"] = res.ClassOverrideReason
	}
	return no
}

func rotaNomeEData(arquivo string, agora time.Time) (string, string) {
	base := strings.TrimSuffix(filepath.Base(arquivo), filepath.Ext(arquivo))
	if m := rotaNomePlano.FindStringSubmatch(base); m != nil {
		return m[2], m[1]
	}
	return strings.ToUpper(base), agora.Format("2006-01-02")
}

func rotaDocumento(tickets []rotaTicket, inv []roteamento.ItemInv, elegiveis map[string][]roteamento.CandidatoVia, nome, data string) (map[string]any, []roteamento.SubstituicaoVia, error) {
	nos := make([]map[string]any, 0, len(tickets))
	var subs []roteamento.SubstituicaoVia
	for _, t := range tickets {
		no, s, err := rotaNo(t, inv, elegiveis)
		if err != nil {
			return nil, nil, err
		}
		nos = append(nos, no)
		subs = append(subs, s...)
	}
	return map[string]any{"name": nome, "version": "1.0.0", "date": data, "execution_strategy": map[string]any{}, "nodes": nos}, subs, nil
}

func rotaStore(inv []roteamento.ItemInv) map[string]roteamento.Homologacao {
	store := map[string]roteamento.Homologacao{}
	for _, it := range inv {
		store[it.Par.Chave()] = roteamento.Homologacao{Estado: it.Estado, PesoMB: it.PesoMB, Motivo: it.Motivo}
	}
	return store
}

func rotaAmbiente(ctx context.Context, d Deps, perfil string) (RotaDeps, error) {
	if d.Rota != nil {
		return *d.Rota, nil
	}
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	global, _ := os.ReadFile(filepath.Join(home, ".aidd", "roteamento.json"))
	if perfil == "" {
		perfil = filepath.Join(".aidd", "roteamento.json")
	}
	projeto, err := os.ReadFile(perfil)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return RotaDeps{}, err
	}
	store := rotaLerHomologacao(filepath.Join(home, ".aidd", "homologacao.json"))
	inventario := func(c context.Context) ([]roteamento.ItemInv, error) {
		return roteamento.LerInventarioUnificado(processo.OSFS(), home, d.discover(c), store)
	}
	return RotaDeps{Inventario: inventario, Global: global, Projeto: projeto}, nil
}

func rotaLerHomologacao(caminho string) map[string]roteamento.Homologacao {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil
	}
	var bruto map[string]struct {
		Estado string `json:"estado"`
		PesoMB int    `json:"peso_mb"`
		Motivo string `json:"motivo"`
	}
	if json.Unmarshal(dados, &bruto) != nil {
		return nil
	}
	store := map[string]roteamento.Homologacao{}
	for k, v := range bruto {
		store[k] = roteamento.Homologacao{Estado: v.Estado, PesoMB: v.PesoMB, Motivo: v.Motivo}
	}
	return store
}

func cmdRota(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("rota", d)
	plano := fs.String("plan", "", "arquivo PLAN")
	saida := fs.String("out", "", "arquivo ROUTE de saída")
	perfil := fs.String("perfil", "", "perfil de projeto (padrão .aidd/roteamento.json)")
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil || *plano == "" || *saida == "" {
		return 2
	}
	return executarRota(ctx, d, *plano, *saida, *perfil, *asJSON)
}

func executarRota(ctx context.Context, d Deps, plano, saida, perfil string, asJSON bool) int {
	texto, err := os.ReadFile(plano)
	if err != nil {
		return fail(d, err)
	}
	amb, err := rotaAmbiente(ctx, d, perfil)
	if err != nil {
		return fail(d, err)
	}
	inv, err := amb.Inventario(ctx)
	if err != nil {
		return fail(d, err)
	}
	elegiveis, err := roteamento.CarregarPerfil(amb.Global, amb.Projeto, rotaStore(inv))
	if err != nil {
		return fail(d, err)
	}
	nome, data := rotaNomeEData(plano, time.Now())
	doc, subs, err := rotaDocumento(rotaTickets(string(texto)), inv, elegiveis, nome, data)
	if err != nil {
		return fail(d, err)
	}
	return rotaGravar(d, doc, subs, saida, asJSON)
}

func rotaGravar(d Deps, doc map[string]any, subs []roteamento.SubstituicaoVia, saida string, asJSON bool) int {
	bytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fail(d, err)
	}
	if err := os.WriteFile(saida, append(bytes, '\n'), 0o644); err != nil {
		return fail(d, err)
	}
	for _, s := range subs {
		fmt.Fprintf(d.Stderr, "aviso: %s usou %s/%s no lugar de %s/%s: %s\n", s.No, s.Usado.Harness, s.Usado.Modelo, s.Pedido.Harness, s.Pedido.Modelo, s.Motivo)
	}
	if asJSON {
		return printJSON(d, doc)
	}
	fmt.Fprintf(d.Stdout, "ROUTE gravada em %s\n", saida)
	return 0
}
