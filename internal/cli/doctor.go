package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/provider"
)

func cmdDoctor(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("doctor", d)
	asJSON := fs.Bool("json", false, "saída JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	inv := d.discover(ctx)
	if *asJSON {
		return printJSON(d, map[string]any{"inventory": inv, "models": inv.Models(), "python": "opcional"})
	}
	w := table(d.Stdout)
	fmt.Fprintln(w, "MODELO\tEXECUTOR")
	for _, m := range inv.Models() {
		fmt.Fprintf(w, "%s\t%s\n", m.ID, m.Source)
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	router := "offline"
	if inv.RouterURL != "" {
		router = "online em " + inv.RouterURL
	}
	fmt.Fprintf(d.Stdout, "\n9Router: %s\n", router)
	fmt.Fprintf(d.Stdout, "Harnesses no PATH: %s\n", joinOrNone(cliNames(inv.CLIs)))
	fmt.Fprintf(d.Stdout, "Chaves de API detectadas: %s\n", joinOrNone(inv.APIKeys))
	fmt.Fprintf(d.Stdout, "Python: opcional\n")
	for _, e := range inv.Errors {
		fmt.Fprintf(d.Stdout, "aviso: %s\n", e)
	}
	return 0
}

func cliNames(cs []provider.CLI) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Name)
	}
	return out
}

func joinOrNone(xs []string) string {
	if len(xs) == 0 {
		return "nenhum"
	}
	return strings.Join(xs, ", ")
}
