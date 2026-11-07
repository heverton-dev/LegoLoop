package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"legoloop/internal/fechar"
)

const fecharUsage = `Uso: legoloop fechar <slug> [opções]

Executa a cadeia determinística de 9 etapas de canonização, telemetria e fechamento universal.

Opções:
  --dry-run        apenas simula a validação das 9 etapas sem gravar artefatos
  --skip-build     pula a etapa de recompilação do binário Go
  --tempo-min N    tempo total estimado em minutos para telemetria
  --tokens-in N    tokens de entrada estimados para telemetria
  --tokens-out N   tokens de saída estimados para telemetria
  --custo-usd N    custo financeiro estimado em USD
  --taxa-cambio N  taxa cambial BRL/USD (default: 5.80)
  --json           saída estruturada em JSON
`

func cmdFechar(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprintf(d.Stderr, "erro: informe o slug do processo.\n\n%s", fecharUsage)
		return 2
	}
	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		fmt.Fprint(d.Stdout, fecharUsage)
		return 0
	}
	slug := strings.TrimSpace(args[0])
	if slug == "" || strings.HasPrefix(slug, "-") {
		fmt.Fprintf(d.Stderr, "erro: informe o slug do processo.\n\n%s", fecharUsage)
		return 2
	}
	return parseFecharEExecutar(ctx, d, slug, args[1:])
}

func parseFecharEExecutar(ctx context.Context, d Deps, slug string, rawArgs []string) int {
	fs := newFlags("fechar", d)
	dryRun := fs.Bool("dry-run", false, "simula sem gravar")
	skipBuild := fs.Bool("skip-build", false, "pula build go")
	asJSON := fs.Bool("json", false, "saída em json")
	tempoMin := fs.Float64("tempo-min", 0, "tempo orçado em minutos")
	tokensIn := fs.Int("tokens-in", 0, "tokens de entrada orçados")
	tokensOut := fs.Int("tokens-out", 0, "tokens de saída orçados")
	custoUSD := fs.Float64("custo-usd", 0, "custo orçado em USD")
	taxaCambio := fs.Float64("taxa-cambio", 5.80, "taxa de câmbio BRL")
	if err := fs.Parse(rawArgs); err != nil {
		return 2
	}
	telem := &fechar.TelemetriaOrcamento{
		TempoEstimadoMin:   *tempoMin,
		TokensInEstimados:  *tokensIn,
		TokensOutEstimados: *tokensOut,
		CustoEstimadoUSD:   *custoUSD,
		TaxaCambioBRL:      *taxaCambio,
	}
	return runFechar(ctx, d, slug, *dryRun, *skipBuild, *asJSON, telem)
}

func runFechar(ctx context.Context, d Deps, slug string, dry, skipBuild, asJSON bool, t *fechar.TelemetriaOrcamento) int {
	req := fechar.Request{Slug: slug, DryRun: dry, SkipBuild: skipBuild, Telemetria: t}
	deps := fechar.Deps{Home: d.Home, SkipBuild: skipBuild}
	res, err := fechar.Executar(ctx, req, deps)
	if err != nil {
		fmt.Fprintf(d.Stderr, "erro no fechamento: %v\n", err)
		return 1
	}
	if asJSON {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Fprintln(d.Stdout, string(b))
		return 0
	}
	reportFechar(d, res)
	return 0
}

func reportFechar(d Deps, res *fechar.Result) {
	fmt.Fprintf(d.Stdout, "✓ Processo %s (%s) canonizado com sucesso!\n", res.Slug, res.CanonName)
	for _, e := range res.Etapas {
		fmt.Fprintf(d.Stdout, "  [%d/9] %s: %s (%s)\n", e.Etapa, e.Nome, e.Status, e.Artefato)
	}
	if res.ReportPath != "" {
		fmt.Fprintf(d.Stdout, "  📄 Relatório: %s\n", res.ReportPath)
	}
}
