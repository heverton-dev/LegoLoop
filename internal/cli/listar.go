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
	"legoloop/internal/ste"
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
	imprimirLacunas(d, res.Processos)
	return 0
}

// imprimirLacunas mostra cada etapa em lacuna e, quando há peça candidata, os comandos prontos para copiar.
func imprimirLacunas(d Deps, its []listar.Item) {
	for _, it := range its {
		for _, etapa := range it.Lacunas {
			fmt.Fprintf(d.Stdout, "\nAVISO: %s\n", ste.Msg("listar.lacuna", etapa, it.Nome))
			n := 0
			for _, s := range it.Sugestoes {
				if s.Etapa == etapa {
					n++
					imprimirSugestao(d, s)
				}
			}
			if n == 0 {
				fmt.Fprintln(d.Stdout, ste.Msg("listar.lacuna_sem_peca", etapa, etapa))
			}
		}
	}
}

func imprimirSugestao(d Deps, s listar.Sugestao) {
	fmt.Fprintln(d.Stdout, ste.Msg("listar.lacuna_sugestao", s.Etapa, s.Caso))
	for _, c := range s.Comandos[:min(len(s.Comandos), listar.MaxComandosTexto)] {
		fmt.Fprintf(d.Stdout, "  %s\n", c)
	}
	switch mais := len(s.Comandos) - listar.MaxComandosTexto; {
	case mais == 1:
		fmt.Fprintln(d.Stdout, ste.Msg("listar.lacuna_mais_uma", s.Etapa))
	case mais > 1:
		fmt.Fprintln(d.Stdout, ste.Msg("listar.lacuna_mais", mais, s.Etapa))
	}
}

func ouTraco(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
