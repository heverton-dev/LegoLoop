package cli

import (
	"context"
	"fmt"

	"github.com/heverton-dev/LegoLoop/internal/espinha"
)

const espinhaUso = `uso: legoloop espinha <verificar> [--raiz DIR] [--json]

Comandos:
  verificar  Executa a verificacao deterministica da espinha dorsal sobre a raiz do repositorio.
`

func cmdEspinha(_ context.Context, args []string, d Deps) int {
	if len(args) == 0 {
		fmt.Fprint(d.Stderr, espinhaUso)
		return 2
	}
	switch args[0] {
	case "verificar":
		return cmdEspinhaVerificar(args[1:], d)
	default:
		fmt.Fprintf(d.Stderr, "subcomando desconhecido: %s\n\n%s", args[0], espinhaUso)
		return 2
	}
}

func cmdEspinhaVerificar(args []string, d Deps) int {
	fs := newFlags("espinha verificar", d)
	raiz := fs.String("raiz", ".", "diretorio raiz do repositorio")
	asJSON := fs.Bool("json", false, "emite as violacoes em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	violacoes := espinha.Verificar(*raiz)
	if *asJSON {
		if violacoes == nil {
			violacoes = []espinha.Violacao{}
		}
		if code := printJSON(d, violacoes); code != 0 {
			return code
		}
	} else {
		for _, v := range violacoes {
			fmt.Fprintln(d.Stdout, v.String())
		}
	}
	if len(violacoes) > 0 {
		return 2
	}
	return 0
}
