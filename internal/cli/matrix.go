package cli

import (
	"fmt"
	"io"
	"strings"

	"legoloop/internal/ste"
)

// RunMatrix é o comando `matrix`, descontinuado: vale como alias de `matriz ver`. Mantém as flags antigas
// (--path, --kind, --tier, --min, --json) e avisa no stderr. A saída agora é a da `matriz ver`.
func RunMatrix(stdout, stderr io.Writer, args []string) int {
	d := Deps{Stdout: stdout, Stderr: stderr}
	fmt.Fprintln(stderr, "aviso:", ste.Msg("matriz.matrix_descontinuado"))
	return cmdMatrizVer(traduzirMatrix(args), d)
}

// traduzirMatrix troca as flags do `matrix` pelas da `matriz ver`: --kind vira --tipo e --path vira --registro.
func traduzirMatrix(args []string) []string {
	renomear := map[string]string{"kind": "tipo", "path": "registro"}
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = a
		base, val, tem := strings.Cut(strings.TrimLeft(a, "-"), "=")
		novo, ok := renomear[base]
		if !ok || !strings.HasPrefix(a, "-") {
			continue
		}
		out[i] = "--" + novo
		if tem {
			out[i] += "=" + val
		}
	}
	return out
}
