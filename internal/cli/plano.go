package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"legoloop/internal/plangen"
	"legoloop/internal/progresso"
	"legoloop/internal/upstream"
)

const planoUso = "uso: legoloop plano --spec SPEC --tickets tickets.json [--raiz DIR] [--nivel epic] [--data AAAA-MM-DD] [--estrito] [--json]"

// opcoesPlano reúne as opções do comando plano.
type opcoesPlano struct {
	raiz, spec, tickets, nivel, data string
	estrito, asJSON                  bool
}

func lerOpcoesPlano(args []string, d Deps) (opcoesPlano, bool) {
	var o opcoesPlano
	fs := newFlags("plano", d)
	fs.StringVar(&o.raiz, "raiz", ".", "raiz do projeto (docs/plans e docs/routes ficam sob ela)")
	fs.StringVar(&o.spec, "spec", "", "arquivo SPEC aprovada")
	fs.StringVar(&o.tickets, "tickets", "", "arquivo JSON com a lista de Tickets")
	fs.StringVar(&o.nivel, "nivel", "epic", "supervisão: nano, micro, epic ou autonomous")
	fs.StringVar(&o.data, "data", "", "data dos documentos AAAA-MM-DD (padrão: hoje)")
	fs.BoolVar(&o.estrito, "estrito", false, "sai com 1 se algum critério da SPEC ficar sem Ticket")
	fs.BoolVar(&o.asJSON, "json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil || o.spec == "" || o.tickets == "" {
		fmt.Fprintln(d.Stderr, planoUso)
		return o, false
	}
	return o, true
}

// cmdPlano renderiza o PLAN e a ROUTE a partir de uma lista de Tickets, sem LLM e de forma determinística.
// A SPEC precisa referenciar uma INTENT aprovada (trava upstream_ok). O passo seguinte é `legoloop rota`
// (harness e modelo) e `go run ./tools/route2dag` (DAG).
func cmdPlano(_ context.Context, args []string, d Deps) int {
	o, ok := lerOpcoesPlano(args, d)
	if !ok {
		return 2
	}
	if err := upstream.VerificarSpec(o.raiz, o.spec); err != nil {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return renderizarPlano(o, d)
}

func renderizarPlano(o opcoesPlano, d Deps) int {
	lista, err := lerTicketsDoPlano(o.tickets)
	if err != nil {
		return fail(d, err)
	}
	quando, err := dataDoPlano(o.data)
	if err != nil {
		return fail(d, err)
	}
	res, err := plangen.Render(o.raiz, o.spec, lista, plangen.SupervisionLevel(o.nivel), quando)
	if err != nil {
		return fail(d, err)
	}
	if err := gerarProgressoDoPlano(o.tickets, res.PlanMarkdownPath, quando); err != nil {
		return fail(d, err)
	}
	for _, w := range res.Warnings {
		fmt.Fprintln(d.Stderr, "aviso:", w)
	}
	if o.estrito && len(res.Warnings) > 0 {
		return fail(d, fmt.Errorf("%d aviso(s) no plano e --estrito ativo", len(res.Warnings)))
	}
	if o.asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "PLAN gravado em %s\nROUTE gravada em %s\n%d Tickets\n", res.PlanMarkdownPath, res.RouteJSONPath, len(res.Tickets))
	return 0
}

func lerTicketsDoPlano(caminho string) ([]plangen.Ticket, error) {
	b, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	var lista []plangen.Ticket
	if err := json.Unmarshal(b, &lista); err != nil {
		return nil, fmt.Errorf("tickets.json inválido: %w", err)
	}
	return lista, nil
}

func dataDoPlano(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("--data deve ser AAAA-MM-DD: %w", err)
	}
	return t, nil
}

func gerarProgressoDoPlano(ticketsArq, planMarkdownPath string, quando time.Time) error {
	entradas, err := progresso.CarregarEntradasTickets(ticketsArq)
	if err != nil {
		return err
	}
	progJSON := derivarCaminhoProgressoPlano(planMarkdownPath)
	p := progresso.NovoDeTickets("", planMarkdownPath, entradas, quando.UTC().Format(time.RFC3339))
	return progresso.Salvar(progJSON, progresso.CaminhoEspelhoMD(progJSON), p)
}
