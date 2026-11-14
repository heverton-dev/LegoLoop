// Comando legoloop: Composable Looping & Harness Fabric.
package main

import (
	"context"
	"os"
	"os/signal"

	"legoloop/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	var code int
	if len(os.Args) > 1 && os.Args[1] == "triagem" {
		code = runTriagem(ctx, os.Args[2:], os.Stdout, os.Stderr)
	} else if len(os.Args) > 1 && os.Args[1] == "migrar-loops" {
		code = runMigrarLoops(os.Args[2:], os.Stdout, os.Stderr)
	} else {
		code = cli.Run(ctx, os.Args[1:], cli.Deps{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr})
	}
	stop()
	os.Exit(code)
}
