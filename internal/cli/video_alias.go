package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/heverton-dev/LegoLoop/internal/video/render"
)

const usoVideoRender = `Uso: legoloop video-render <subcomando> [opções]

Subcomandos:
  timeline  --projeto DIR   compõe a timeline.json a partir dos contratos do projeto
  executar  --projeto DIR   renderiza o master 16:9 a partir da timeline salva
`

// avisoDepreciacao devolve a linha que o alias antigo imprime em stderr antes de executar o comando novo.
func avisoDepreciacao(antigo, novo string) string {
	return fmt.Sprintf("aviso: o comando %s está depreciado; use %s\n", antigo, novo)
}

// aliasDepreciado envolve o comando novo: imprime o aviso em stderr e devolve o mesmo exit code.
func aliasDepreciado(antigo, novo string, cmd func(context.Context, []string, Deps) int) func(context.Context, []string, Deps) int {
	return func(ctx context.Context, args []string, d Deps) int {
		fmt.Fprint(d.Stderr, avisoDepreciacao(antigo, novo))
		return cmd(ctx, args, d)
	}
}

func cmdVideoRender(ctx context.Context, args []string, d Deps) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Fprint(d.Stdout, usoVideoRender)
		return 0
	}
	switch args[0] {
	case "timeline":
		return cmdVideoRenderTimeline(ctx, args[1:], d)
	case "executar":
		return cmdVideoRenderExecutar(ctx, args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido %q\n%s", args[0], usoVideoRender)
		return 2
	}
}

// parseProjetoRender lê a flag obrigatória --projeto do subcomando de render.
func parseProjetoRender(sub string, args []string, stderr io.Writer) (string, bool) {
	fs := flag.NewFlagSet("video-render "+sub, flag.ContinueOnError)
	fs.SetOutput(stderr)
	projeto := fs.String("projeto", "", "diretório do projeto de vídeo")
	if err := fs.Parse(args); err != nil {
		return "", false
	}
	if *projeto == "" {
		fmt.Fprintf(stderr, "video-render %s: --projeto é obrigatório\n", sub)
		return "", false
	}
	return *projeto, true
}

func cmdVideoRenderTimeline(ctx context.Context, args []string, d Deps) int {
	projeto, ok := parseProjetoRender("timeline", args, d.Stderr)
	if !ok {
		return 2
	}
	tl, err := render.ExecutarTimeline(ctx, projeto)
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render timeline: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "timeline composta: %d batidas, %.2f s\n", len(tl.Batidas), tl.DuracaoS)
	return 0
}

func cmdVideoRenderExecutar(ctx context.Context, args []string, d Deps) int {
	projeto, ok := parseProjetoRender("executar", args, d.Stderr)
	if !ok {
		return 2
	}
	rel, err := render.ExecutarRender(ctx, projeto, render.RenderDeps{})
	if err != nil {
		fmt.Fprintf(d.Stderr, "video-render executar: %v\n", err)
		return 1
	}
	fmt.Fprintf(d.Stdout, "render concluído: encoder %s, %.1f LUFS, %d saídas\n", rel.Encoder, rel.LUFS, len(rel.Renders))
	return 0
}
