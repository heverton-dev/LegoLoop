package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"legoloop/internal/fechar"
)

const fecharUsage = `Uso: legoloop fechar <slug> [opções]

Executa a cadeia determinística de 8 etapas de canonização e fechamento universal.

Opções:
  --dry-run     apenas simula a validação das 8 etapas sem gravar artefatos
  --skip-build  pula a etapa de recompilação do binário Go
  --json        saída estruturada em JSON
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
	fs := newFlags("fechar", d)
	dryRun := fs.Bool("dry-run", false, "simula sem gravar")
	skipBuild := fs.Bool("skip-build", false, "pula build go")
	asJSON := fs.Bool("json", false, "saída em json")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	slug := strings.TrimSpace(args[0])
	if slug == "" || strings.HasPrefix(slug, "-") {
		fmt.Fprintf(d.Stderr, "erro: informe o slug do processo.\n\n%s", fecharUsage)
		return 2
	}
	return runFechar(ctx, d, slug, *dryRun, *skipBuild, *asJSON)
}

func runFechar(ctx context.Context, d Deps, slug string, dry, skipBuild, asJSON bool) int {
	req := fechar.Request{Slug: slug, DryRun: dry, SkipBuild: skipBuild}
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
		fmt.Fprintf(d.Stdout, "  [%d/8] %s: %s (%s)\n", e.Etapa, e.Nome, e.Status, e.Artefato)
	}
}
