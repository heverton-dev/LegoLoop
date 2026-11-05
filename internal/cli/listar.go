package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/canon"
	"legoloop/internal/listar"
	"legoloop/internal/processo"
)

// cmdListar mostra os processos do cofre. Sem --vault, um cofre padrão que ainda não existe conta como vazio.
func cmdListar(ctx context.Context, args []string, d Deps) int {
	fs := newFlags("listar", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	asJSON := fs.Bool("json", false, "saída em JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(d.Stderr, "erro: argumento inesperado %q (uso: legoloop listar [--vault DIR] [--json])\n", fs.Arg(0))
		return 2
	}
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	v, padrao, err := vaultListar(*vault, d.Home)
	if err != nil {
		return fail(d, err)
	}
	if _, err := os.Stat(v); padrao && os.IsNotExist(err) {
		return imprimirListar(d, listar.Resultado{Vault: v, Processos: []listar.Item{}}, *asJSON)
	}
	res, err := listar.Run(ctx, listar.Deps{Store: processo.Store{Roots: processo.Roots{Vault: v}}, Home: home})
	if err != nil {
		return fail(d, err)
	}
	return imprimirListar(d, res, *asJSON)
}

// vaultListar resolve o cofre: --vault, senão HOME/.aidd/canonical. padrao diz se veio do padrão.
func vaultListar(flagVault, home string) (string, bool, error) {
	v := flagVault
	switch {
	case v != "":
	case home != "":
		v = filepath.Join(home, ".aidd", "canonical")
	default:
		dv, err := canon.DefaultVault()
		if err != nil {
			return "", true, err
		}
		v = dv
	}
	abs, err := filepath.Abs(v)
	return abs, flagVault == "", err
}

func imprimirListar(d Deps, res listar.Resultado, asJSON bool) int {
	for _, it := range res.Processos {
		if it.Erro != "" {
			fmt.Fprintf(d.Stderr, "aviso: %s não pôde ser lido: %s\n", it.Nome, it.Erro)
		}
	}
	if asJSON {
		if len(res.Processos) == 0 {
			fmt.Fprintln(d.Stderr, listar.MensagemVazio(res.Vault))
		}
		return printJSON(d, res)
	}
	if len(res.Processos) == 0 {
		fmt.Fprintln(d.Stdout, listar.MensagemVazio(res.Vault))
		return 0
	}
	fmt.Fprintf(d.Stdout, "\n=== PROCESSOS NO COFRE (%d) ===\n\n", len(res.Processos))
	w := table(d.Stdout)
	fmt.Fprintln(w, "NOME\tVERSÃO\tNATUREZA\tHARNESSES\tUSADO POR")
	for _, it := range res.Processos {
		natureza := it.Natureza
		if it.Erro != "" {
			natureza = "inválido"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", it.Nome, ouTraco(it.Versao), ouTraco(natureza),
			ouTraco(strings.Join(it.Harnesses, ", ")), ouTraco(strings.Join(it.UsadoPor, ", ")))
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	return 0
}

func ouTraco(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
