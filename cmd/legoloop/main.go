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
	code := cli.Run(ctx, os.Args[1:], cli.Deps{Stdout: os.Stdout, Stderr: os.Stderr})
	stop()
	os.Exit(code)
}
