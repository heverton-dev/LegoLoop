package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"legoloop/internal/caminhos"
	"legoloop/internal/migrar"
)

func runMigrarProcessos(subcmd string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet(subcmd, flag.ContinueOnError)
	fs.SetOutput(stderr)
	dryRun := fs.Bool("dry-run", false, "simula alterações sem gravar em disco")
	apply := fs.Bool("apply", false, "aplica alterações no diretório de processos")
	dir := fs.String("dir", caminhos.DirProcessos, "diretório raiz dos processos")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	efetivo := *apply && !*dryRun
	rel, err := migrar.MigrarProcessos(*dir, efetivo)
	if err != nil {
		fmt.Fprintf(stderr, "erro ao migrar processos: %v\n", err)
		return 1
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(rel)
	return 0
}
