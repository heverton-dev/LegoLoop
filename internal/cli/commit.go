package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/heverton-dev/LegoLoop/internal/commit"
)

func cmdCommit(ctx context.Context, args []string, d Deps) int {
	fs := flag.NewFlagSet("commit", flag.ContinueOnError)
	fs.SetOutput(d.Stderr)
	var opt commit.Opcoes
	fs.StringVar(&opt.Mensagem, "m", "", "mensagem de commit (opcional)")
	fs.BoolVar(&opt.DryRun, "dry-run", false, "simula o commit sem alterar o git")
	fs.BoolVar(&opt.NoPush, "no-push", false, "realiza apenas o commit local sem push")
	fs.BoolVar(&opt.Verbose, "verbose", false, "exibe saída detalhada")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if opt.Mensagem == "" && fs.NArg() > 0 {
		opt.Mensagem = fs.Arg(0)
	}

	if err := commit.Executar(ctx, d.Stdout, d.Stderr, opt); err != nil {
		fmt.Fprintf(d.Stderr, "erro: %v\n", err)
		return 1
	}
	return 0
}
