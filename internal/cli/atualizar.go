package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/atualizar"
)

type opcoesCmdAtualizar struct {
	canal   string
	baseURL string
	binario string
	home    string
	asJSON  bool
}

func cmdAtualizar(ctx context.Context, args []string, d Deps) int {
	op, code := parseFlagsAtualizar(args, d)
	if code != 0 {
		return code
	}
	home, err := resolverHomeInstalar(d, op.home)
	if err != nil {
		return fail(d, err)
	}
	binario, err := resolverBinarioAtualizar(op.binario)
	if err != nil {
		return fail(d, err)
	}
	res, err := atualizar.Executar(ctx, atualizar.Opcao{
		Canal:   op.canal,
		BaseURL: op.baseURL,
		Binario: binario,
		Home:    home,
	})
	if err != nil {
		return tratarErroAtualizar(d, err)
	}
	return exibirResultadoAtualizar(d, res, op.asJSON)
}

func parseFlagsAtualizar(args []string, d Deps) (opcoesCmdAtualizar, int) {
	var op opcoesCmdAtualizar
	fs := newFlags("atualizar", d)
	fs.StringVar(&op.canal, "canal", atualizar.CanalEstavel, "canal de atualizacao: estavel ou beta")
	fs.StringVar(&op.baseURL, "base-url", atualizar.BaseURLPadrao, "URL base de download do pacote e SHA256SUMS")
	fs.StringVar(&op.binario, "binario", "", "caminho do binario alvo (padrao: executavel atual)")
	fs.StringVar(&op.home, "home", "", "diretorio HOME alvo para revalidacao da Forja")
	fs.BoolVar(&op.asJSON, "json", false, "emite o resultado em JSON")
	if err := fs.Parse(args); err != nil {
		return op, 2
	}
	return op, 0
}

func resolverBinarioAtualizar(flagBinario string) (string, error) {
	if strings.TrimSpace(flagBinario) != "" {
		return flagBinario, nil
	}
	return os.Executable()
}

func tratarErroAtualizar(d Deps, err error) int {
	if errors.Is(err, atualizar.ErrCanalInvalido) ||
		errors.Is(err, atualizar.ErrSemBaseURL) ||
		errors.Is(err, atualizar.ErrSemBinario) ||
		errors.Is(err, atualizar.ErrSemHome) {
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	}
	return fail(d, err)
}

func exibirResultadoAtualizar(d Deps, res atualizar.Resultado, asJSON bool) int {
	if asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "Atualizar (%s): binario=%s prev=%s sha256=%s revalidado=%t\n",
		res.Canal, res.Binario, res.BinarioPrev, res.SHA256Novo, res.Revalidado)
	return 0
}
