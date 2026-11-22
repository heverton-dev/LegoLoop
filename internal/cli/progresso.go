package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/progresso"
)

const progressoUso = `uso: legoloop progresso <gerar|iniciar|validar|status|verificar> [opções]
  gerar   --tickets T.json --out ARQ.json [--plano P] [--iniciativa I]
  iniciar --arquivo ARQ.json --ticket ID --responsavel R [--nota N]
  validar --arquivo ARQ.json --ticket ID [--raiz DIR]
  status  --arquivo ARQ.json [--json]
  verificar --arquivo ARQ.json [--max-idade 24h] [--sondar] [--reconciliar] [--raiz DIR] [--json]`

func cmdProgresso(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	switch args[0] {
	case "gerar":
		return cmdProgressoGerar(args[1:], d)
	case "iniciar":
		return cmdProgressoIniciar(args[1:], d)
	case "validar":
		return cmdProgressoValidar(ctx, args[1:], d)
	case "status":
		return cmdProgressoStatus(args[1:], d)
	case "verificar":
		return cmdProgressoVerificar(ctx, args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s\n", args[0], progressoUso)
		return 2
	}
}

func cmdProgressoGerar(args []string, d Deps) int {
	fs := newFlags("progresso gerar", d)
	ticketsArq := fs.String("tickets", "", "arquivo *_TICKETS_*.json")
	outArq := fs.String("out", "", "arquivo JSON de progresso")
	plano := fs.String("plano", "", "caminho do PLAN correspondente")
	iniciativa := fs.String("iniciativa", "", "código da iniciativa")
	if err := fs.Parse(args); err != nil || *ticketsArq == "" || *outArq == "" {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	entradas, err := progresso.CarregarEntradasTickets(*ticketsArq)
	if err != nil {
		return fail(d, err)
	}
	p := progresso.NovoDeTickets(*iniciativa, *plano, entradas, time.Now().UTC().Format(time.RFC3339))
	mdArq := progresso.CaminhoEspelhoMD(*outArq)
	if err := progresso.Salvar(*outArq, mdArq, p); err != nil {
		return fail(d, err)
	}
	fmt.Fprintf(d.Stdout, "Progresso gravado em %s e %s (%d Tickets)\n", *outArq, mdArq, p.Resumo.Total)
	return 0
}

func cmdProgressoIniciar(args []string, d Deps) int {
	fs := newFlags("progresso iniciar", d)
	arq := fs.String("arquivo", "", "arquivo JSON de progresso")
	ticket := fs.String("ticket", "", "ID do Ticket")
	resp := fs.String("responsavel", "", "nome do agente ou worker")
	nota := fs.String("nota", "", "ponto exato de retomada")
	if err := fs.Parse(args); err != nil || *arq == "" || *ticket == "" || *resp == "" {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	p, err := progresso.Carregar(*arq)
	if err != nil {
		return fail(d, err)
	}
	if err := progresso.Iniciar(&p, *ticket, *resp, *nota, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return fail(d, err)
	}
	if err := progresso.Salvar(*arq, progresso.CaminhoEspelhoMD(*arq), p); err != nil {
		return fail(d, err)
	}
	fmt.Fprintf(d.Stdout, "%s marcado como IN_PROGRESS (%s)\n", *ticket, *resp)
	return 0
}

func cmdProgressoValidar(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("progresso validar", d)
	arq := fs.String("arquivo", "", "arquivo JSON de progresso")
	ticket := fs.String("ticket", "", "ID do Ticket")
	raiz := fs.String("raiz", ".", "diretório de execução do Gate")
	if err := fs.Parse(args); err != nil || *arq == "" || *ticket == "" {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	p, err := progresso.Carregar(*arq)
	if err != nil {
		return fail(d, err)
	}
	passou, errExec := progresso.ValidarTicket(ctx, &p, *ticket, *raiz, time.Now().UTC().Format(time.RFC3339), nil)
	if err := progresso.Salvar(*arq, progresso.CaminhoEspelhoMD(*arq), p); err != nil {
		return fail(d, err)
	}
	if errExec != nil || !passou {
		fmt.Fprintf(d.Stderr, "%s reprovado no Gate\n", *ticket)
		return 2
	}
	fmt.Fprintf(d.Stdout, "%s promovido para PASSED\n", *ticket)
	return 0
}

func cmdProgressoStatus(args []string, d Deps) int {
	fs := newFlags("progresso status", d)
	arq := fs.String("arquivo", "", "arquivo JSON de progresso")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil || *arq == "" {
		fmt.Fprintln(d.Stderr, progressoUso)
		return 2
	}
	p, err := progresso.Carregar(*arq)
	if err != nil {
		return fail(d, err)
	}
	if *asJSON {
		return printJSON(d, p)
	}
	fmt.Fprintf(d.Stdout, "%s: %d/%d PASSED | %d IN_PROGRESS | %d FAILED | %d PENDING\n",
		p.Iniciativa, p.Resumo.Aprovados, p.Resumo.Total, p.Resumo.EmAndamento, p.Resumo.Reprovados, p.Resumo.Pendentes)
	return 0
}

// derivarCaminhoProgressoPlano devolve docs/progress/<data>_plan-<nome>/progresso.json.
// A data e o nome vêm do arquivo do PLAN, então o caminho não muda entre execuções.
func derivarCaminhoProgressoPlano(raiz, planPath string) string {
	base := strings.TrimSuffix(filepath.Base(planPath), filepath.Ext(planPath))
	data, nome := "", base
	if i := strings.Index(base, "_PLAN_"); i >= 0 {
		data, nome = base[:i], base[i+len("_PLAN_"):]
	}
	nome = strings.ToLower(strings.ReplaceAll(nome, "_", "-"))
	pasta := "plan-" + nome
	if data != "" {
		pasta = data + "_" + pasta
	}
	return filepath.Join(raiz, "docs", "progress", pasta, "progresso.json")
}
