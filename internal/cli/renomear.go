package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"legoloop/internal/jail"
	"legoloop/internal/renomear"
)

const renomearUsage = `Uso: legoloop renomear <antigo> <novo> [--dry-run] [--vault DIR] [--docs DIR] [--raiz DIR] [--yes] [--json]

Renomeia um Processo na pasta processos/, no DAG, em triagem.json, anticorpos.json,
vínculos do Cofre e referências em documentos, gerando backup com MANIFESTO.tsv
e gravando o histórico em espinha/renomeados.json.
`

type renomearFlags struct {
	dryRun bool
	vault  string
	docs   string
	raiz   string
	yes    bool
	asJSON bool
}

func cmdRenomear(_ context.Context, args []string, d Deps) int {
	antigo, novo, op, code := parseRenomearFlags(args, d)
	if code != 0 {
		return code
	}
	if err := validarNomesCLI(antigo, novo); err != nil {
		return failRenomear(d, err)
	}
	if !op.dryRun && !op.yes {
		if code := confirmarRenomear(d, antigo, novo); code != 0 {
			return code
		}
	}
	return executarRenomearCLI(d, antigo, novo, op)
}

func validarNomesCLI(antigo, novo string) error {
	if err := jail.ValidarNomeProcesso(antigo); err != nil {
		return err
	}
	return jail.ValidarNomeProcesso(novo)
}

func executarRenomearCLI(d Deps, antigo, novo string, op renomearFlags) int {
	home := d.Home
	if home == "" && op.vault != "" {
		home, _ = os.UserHomeDir()
	}
	rel, err := renomear.Executar(renomear.Opcoes{
		Antigo: antigo,
		Novo:   novo,
		Raiz:   op.raiz,
		Vault:  op.vault,
		Home:   home,
		Docs:   op.docs,
		DryRun: op.dryRun,
	})
	if err != nil {
		return failRenomear(d, err)
	}
	return imprimirRenomear(d, rel, op.asJSON)
}

func parseRenomearFlags(args []string, d Deps) (string, string, renomearFlags, int) {
	fs := newFlags("renomear", d)
	var op renomearFlags
	fs.BoolVar(&op.dryRun, "dry-run", false, "simula a renomeação sem alterar arquivos")
	fs.StringVar(&op.vault, "vault", "", "cofre canônico opcional")
	fs.StringVar(&op.docs, "docs", "", "diretório de documentos (padrão <raiz>/docs)")
	fs.StringVar(&op.raiz, "raiz", ".", "diretório raiz do projeto")
	fs.BoolVar(&op.yes, "yes", false, "confirma sem perguntar")
	fs.BoolVar(&op.asJSON, "json", false, "saída em JSON")
	pos, err := parseIntercalado(fs, args)
	if err != nil {
		return "", "", op, 2
	}
	if len(pos) != 2 {
		fmt.Fprintf(d.Stderr, "erro: informe <antigo> e <novo>.\n\n%s", renomearUsage)
		return "", "", op, 2
	}
	return pos[0], pos[1], op, 0
}

func confirmarRenomear(d Deps, antigo, novo string) int {
	if d.Stdin == nil {
		fmt.Fprintln(d.Stderr, "erro: a confirmação exige um terminal. Use --yes para confirmar.")
		return 2
	}
	ok, err := newPrompter(d.Stdin, d.Stdout).yesNo(fmt.Sprintf("Renomear %s para %s? [s/N]: ", antigo, novo))
	if err != nil {
		return fail(d, err)
	}
	if !ok {
		fmt.Fprintln(d.Stderr, "Renomeação cancelada.")
		return 1
	}
	return 0
}

func failRenomear(d Deps, err error) int {
	if errors.Is(err, jail.ErrEscapeRaiz) ||
		errors.Is(err, renomear.ErrNomeInvalido) ||
		errors.Is(err, renomear.ErrNaoEncontrado) ||
		errors.Is(err, renomear.ErrColisao) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

func imprimirRenomear(d Deps, rel renomear.Relatorio, asJSON bool) int {
	if asJSON {
		return printJSON(d, rel)
	}
	if rel.DryRun {
		fmt.Fprintf(d.Stdout, "[dry-run] %s -> %s (%d alterações previstas)\n", rel.Antigo, rel.Novo, len(rel.Alteracoes))
	} else {
		fmt.Fprintf(d.Stdout, "Renomeado: %s -> %s (%d alterações)\n", rel.Antigo, rel.Novo, len(rel.Alteracoes))
		fmt.Fprintf(d.Stdout, "Backup: %s\n", rel.Manifesto)
	}
	for _, a := range rel.Alteracoes {
		fmt.Fprintf(d.Stdout, "  - [%s] %s -> %s\n", a.Tipo, a.Origem, a.Destino)
	}
	return 0
}
