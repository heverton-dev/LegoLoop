package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/progresso"
	"github.com/heverton-dev/LegoLoop/internal/routedag"
	"github.com/heverton-dev/LegoLoop/internal/upstream"
)

const fabricarPlanoUso = "uso: legoloop fabricar plano --spec SPEC --tickets tickets.json [--raiz DIR] [--nivel epic] [--data AAAA-MM-DD] [--permissivo] [--json]"

type opcoesFabricarPlano struct {
	raiz, spec, tickets, nivel, data string
	permissivo, asJSON               bool
}

type fabricarPlanoResult struct {
	PlanMarkdownPath string           `json:"plan_markdown_path"`
	RouteJSONPath    string           `json:"route_json_path"`
	DagJSONPath      string           `json:"dag_json_path"`
	Tickets          []plangen.Ticket `json:"tickets"`
	Warnings         []string         `json:"warnings,omitempty"`
}

func parseFabricarPlanoFlags(args []string, d Deps) (opcoesFabricarPlano, bool) {
	var o opcoesFabricarPlano
	var estrito bool
	fs := newFlags("fabricar plano", d)
	fs.StringVar(&o.raiz, "raiz", ".", "raiz do projeto (docs/plans, docs/routes e docs/dags ficam sob ela)")
	fs.StringVar(&o.spec, "spec", "", "arquivo SPEC aprovada")
	fs.StringVar(&o.tickets, "tickets", "", "arquivo JSON com a lista de Tickets")
	fs.StringVar(&o.nivel, "nivel", "epic", "supervisão: nano, micro, epic ou autonomous")
	fs.StringVar(&o.data, "data", "", "data dos documentos AAAA-MM-DD (padrão: hoje)")
	fs.BoolVar(&o.permissivo, "permissivo", false, "permite critérios da SPEC sem Ticket correspondente")
	fs.BoolVar(&estrito, "estrito", true, "modo estrito é o padrão (mantido para compatibilidade)")
	fs.BoolVar(&o.asJSON, "json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil || o.spec == "" || o.tickets == "" {
		fmt.Fprintln(d.Stderr, fabricarPlanoUso)
		return o, false
	}
	return o, true
}

// cmdFabricarPlano grava PLAN, ROUTE e DAG em um único comando a partir de SPEC e Tickets.
func cmdFabricarPlano(_ context.Context, args []string, d Deps) int {
	o, ok := parseFabricarPlanoFlags(args, d)
	if !ok {
		return 2
	}
	if err := upstream.VerificarSpec(o.raiz, o.spec); err != nil {
		fmt.Fprintf(d.Stderr, "erro: %v; execute o processo legoloop-concepcao para aprovar a INTENT\n", err)
		return 2
	}
	return executarFabricarPlano(o, d)
}

func executarFabricarPlano(o opcoesFabricarPlano, d Deps) int {
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
		limparArtefatosPlano(o.raiz, o.spec, nil, "")
		return fail(d, err)
	}
	if !o.permissivo && len(res.Warnings) > 0 {
		return abortarPlanoEstrito(d, o, res)
	}
	if err := gerarProgressoDoPlano(o.raiz, o.tickets, res.PlanMarkdownPath, quando); err != nil {
		limparArtefatosPlano(o.raiz, o.spec, res, "")
		return fail(d, err)
	}
	caminhoDAG, err := projetarEGravasDAG(o.raiz, res.RouteJSONPath)
	if err != nil {
		limparArtefatosPlano(o.raiz, o.spec, res, caminhoDAG)
		return fail(d, err)
	}
	return concluirFabricarPlano(o, d, res, caminhoDAG)
}

func abortarPlanoEstrito(d Deps, o opcoesFabricarPlano, res *plangen.PlanResult) int {
	limparArtefatosPlano(o.raiz, o.spec, res, "")
	for _, w := range res.Warnings {
		fmt.Fprintln(d.Stderr, "aviso:", w)
	}
	return fail(d, fmt.Errorf("%d aviso(s) no plano e modo estrito ativo (use --permissivo)", len(res.Warnings)))
}

func limparArtefatosPlano(raiz, specPath string, res *plangen.PlanResult, caminhoDAG string) {
	planPath, routePath := derivarCaminhosPlanoERota(raiz, specPath, res)
	removerArquivoSeExistir(planPath)
	removerArquivoSeExistir(routePath)
	limparProgressoResidual(raiz, planPath)
	limparDAGResidual(raiz, specPath, caminhoDAG)
}

func limparProgressoResidual(raiz, planPath string) {
	if planPath == "" {
		return
	}
	progJSON := derivarCaminhoProgressoPlano(raiz, planPath)
	removerArquivoSeExistir(progJSON)
	removerArquivoSeExistir(progresso.CaminhoEspelhoMD(progJSON))
	removerArquivoSeExistir(filepath.Dir(progJSON))
}

func limparDAGResidual(raiz, specPath, caminhoDAG string) {
	if caminhoDAG != "" {
		removerArquivoSeExistir(caminhoDAG)
		return
	}
	if specPath != "" {
		base := filepath.Base(specPath)
		dagName := strings.Replace(strings.Replace(base, "SPEC", "DAG", 1), ".md", ".json", 1)
		removerArquivoSeExistir(filepath.Join(raiz, "docs", "dags", dagName))
	}
}

func derivarCaminhosPlanoERota(raiz, specPath string, res *plangen.PlanResult) (string, string) {
	if res != nil {
		return res.PlanMarkdownPath, res.RouteJSONPath
	}
	if specPath == "" {
		return "", ""
	}
	base := filepath.Base(specPath)
	planName := strings.Replace(base, "SPEC", "PLAN", 1)
	routeName := strings.Replace(strings.Replace(base, "SPEC", "ROUTE", 1), ".md", ".json", 1)
	return filepath.Join(raiz, "docs", "plans", planName), filepath.Join(raiz, "docs", "routes", routeName)
}

func removerArquivoSeExistir(caminho string) {
	if caminho == "" {
		return
	}
	if err := os.Remove(caminho); err != nil && !os.IsNotExist(err) {
		return
	}
}

func projetarEGravasDAG(raiz, caminhoRoute string) (string, error) {
	rawRoute, err := os.ReadFile(caminhoRoute)
	if err != nil {
		return "", err
	}
	if err := routedag.VerificarCadeia(rawRoute, raiz); err != nil {
		return "", err
	}
	dagBytes, nomeDAG, err := routedag.Converter(rawRoute)
	if err != nil {
		return "", err
	}
	if _, err := dag.Parse(dagBytes); err != nil {
		return "", fmt.Errorf("DAG inválido: %w", err)
	}
	caminhoDAG := filepath.Join(raiz, "docs", "dags", nomeDAG)
	if err := os.MkdirAll(filepath.Dir(caminhoDAG), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(caminhoDAG, dagBytes, 0o644); err != nil {
		return "", err
	}
	return caminhoDAG, nil
}

func concluirFabricarPlano(o opcoesFabricarPlano, d Deps, res *plangen.PlanResult, caminhoDAG string) int {
	for _, w := range res.Warnings {
		fmt.Fprintln(d.Stderr, "aviso:", w)
	}
	if o.asJSON {
		return printJSON(d, fabricarPlanoResult{
			PlanMarkdownPath: res.PlanMarkdownPath,
			RouteJSONPath:    res.RouteJSONPath,
			DagJSONPath:      caminhoDAG,
			Tickets:          res.Tickets,
			Warnings:         res.Warnings,
		})
	}
	fmt.Fprintf(d.Stdout, "PLAN gravado em %s\nROUTE gravada em %s\nDAG gravado em %s\n%d Tickets\n",
		res.PlanMarkdownPath, res.RouteJSONPath, caminhoDAG, len(res.Tickets))
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

func gerarProgressoDoPlano(raiz, ticketsArq, planMarkdownPath string, quando time.Time) error {
	entradas, err := progresso.CarregarEntradasTickets(ticketsArq)
	if err != nil {
		return err
	}
	iniciativa := ""
	if len(entradas) > 0 {
		iniciativa = entradas[0].InitiativeID
	}
	if iniciativa == "" {
		iniciativa = filepath.Base(planMarkdownPath)
	}
	progJSON := derivarCaminhoProgressoPlano(raiz, planMarkdownPath)
	p := progresso.NovoDeTickets(iniciativa, planMarkdownPath, entradas, quando.UTC().Format(time.RFC3339))
	return progresso.Salvar(progJSON, progresso.CaminhoEspelhoMD(progJSON), p)
}
