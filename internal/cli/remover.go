package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"legoloop/internal/processo"
	"legoloop/internal/remover"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("remover.uso", "Uso: legoloop remover <nome> [--vault DIR] [--yes].")
	ste.Default.Register("remover.pergunta", "Remover o processo? [s/N]: ")
	ste.Default.Register("remover.sem_terminal", "A confirmação exige um terminal. Use --yes para confirmar sem perguntar.")
}

// cmdRemover apaga um processo do cofre e os vínculos dele. Mostra o cartão e pede confirmação; --yes confirma.
// Um processo em uso falha com exit 1 e nada muda. Sem confirmação, nada muda.
func cmdRemover(ctx context.Context, args []string, d Deps) int {
	nome, rest := splitPositional(args)
	fs := newFlags("remover", d)
	vault := fs.String("vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	ws := fs.Bool("workspace", false, "apaga também os artefatos locais do workspace")
	yes := fs.Bool("yes", false, "confirma sem perguntar")
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if nome == "" || fs.NArg() > 0 {
		fmt.Fprintln(d.Stderr, "erro:", ste.Msg("remover.uso"))
		return 2
	}
	home := d.Home
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	v, _, err := vaultListar(*vault, d.Home)
	if err != nil {
		return fail(d, err)
	}
	var wsDir string
	if *ws {
		wsDir, _ = os.Getwd()
	}
	rd := remover.Deps{Store: processo.Store{Roots: processo.Roots{Vault: v}}, Home: home, Workspace: wsDir}
	plano, err := remover.Planejar(ctx, rd, nome)
	if err != nil {
		return fail(d, err)
	}
	return confirmarERemover(ctx, d, plano, *yes)
}

// confirmarERemover mostra o cartão, pede a confirmação e só então apaga.
func confirmarERemover(ctx context.Context, d Deps, plano remover.Plano, yes bool) int {
	fmt.Fprintf(d.Stdout, "\n%s", plano.Cartao.Texto())
	if !yes {
		if d.Stdin == nil {
			return fail(d, errors.New(ste.Msg("remover.sem_terminal")))
		}
		ok, err := newPrompter(d.Stdin, d.Stdout).yesNo(ste.Msg("remover.pergunta"))
		if err != nil {
			return fail(d, err)
		}
		if !ok {
			fmt.Fprintln(d.Stderr, remover.MensagemCancelado())
			return 1
		}
	}
	if err := remover.Executar(ctx, plano); err != nil {
		return fail(d, err)
	}
	fmt.Fprintln(d.Stdout, remover.MensagemFeito(plano))
	return 0
}
