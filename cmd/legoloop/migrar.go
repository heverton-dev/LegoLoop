package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"legoloop/internal/migrar"
)

func runMigrarLoops(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("migrar-loops", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dryRun := fs.Bool("dry-run", false, "simula alterações sem gravar em disco")
	apply := fs.Bool("apply", false, "aplica alterações em loops/")
	dir := fs.String("dir", "loops", "diretório raiz dos processos")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	efetivo := *apply && !*dryRun
	rel, err := migrar.MigrarLoops(*dir, efetivo)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao migrar loops: %v\n", err)
		return 1
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(rel)
	return 0
}
