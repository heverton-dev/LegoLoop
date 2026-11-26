package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/heverton-dev/LegoLoop/internal/ste"
)

// executarSteptLint aplica a lei da prova ao texto da entrada padrao: exit 1 quando ha evasao.
// Serve de gancho de parada para qualquer harness: o texto da resposta entra por stdin.
func executarSteptLint(d Deps, asJSON, hook bool) int {
	if d.Stdin == nil {
		fmt.Fprintln(d.Stderr, "stept lint: informe o texto pela entrada padrão")
		return 2
	}
	if hook {
		return executarGanchoLint(d)
	}
	b, err := io.ReadAll(d.Stdin)
	if err != nil {
		fmt.Fprintf(d.Stderr, "stept lint: %v\n", err)
		return 1
	}
	vs := ste.LintEvasao(string(b))
	if asJSON {
		if vs == nil {
			vs = []ste.Violation{}
		}
		saida, _ := json.Marshal(vs)
		fmt.Fprintln(d.Stdout, string(saida))
	} else {
		for _, v := range vs {
			fmt.Fprintf(d.Stderr, "%s: %s\n", v.Regra, v.Msg)
		}
	}
	if len(vs) > 0 {
		return 1
	}
	return 0
}
