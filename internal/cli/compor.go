package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"legoloop/internal/compor"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"compor.uso":          "Uso: legoloop compor <nome> <legoloop-x@versão> [--titulo T] [--depende ID] [--vault DIR] [--docs DIR] [--yes] [--json]",
		"compor.sem_terminal": "Não há terminal para confirmar. Use --yes para seguir sem perguntar.",
		"compor.segue_yes":    "Segue sem perguntar, por causa do --yes.",
	} {
		ste.Default.Register(id, msg)
	}
}

// comporOpts reúne as flags do comando.
type comporOpts struct {
	yes, asJSON         bool
	vault, docs, titulo string
	depende             listFlag
	nome, composto      string
}

// cmdCompor liga um processo do cofre como etapa de outro. Mostra o plano e o cartão, pede a confirmação e grava.
func cmdCompor(ctx context.Context, args []string, d Deps) int {
	var o comporOpts
	fs := newFlags("compor", d)
	fs.BoolVar(&o.yes, "yes", false, "segue sem perguntar")
	fs.BoolVar(&o.asJSON, "json", false, "imprime o resultado em JSON")
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto que tem docs/ (padrão: diretório atual)")
	fs.StringVar(&o.titulo, "titulo", "", "título da etapa nova")
	fs.Var(&o.depende, "depende", "etapa que vem antes da etapa nova (repetível; padrão: as etapas finais)")
	pos, rest := posicionais(args, 2)
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if pos = append(pos, fs.Args()...); len(pos) != 2 {
		fmt.Fprintln(d.Stderr, ste.Msg("compor.uso"))
		return 2
	}
	o.nome, o.composto = pos[0], pos[1]
	return runCompor(ctx, o, d)
}

// posicionais tira do começo de args até n argumentos que não são flag.
func posicionais(args []string, n int) ([]string, []string) {
	var pos []string
	for len(pos) < n && len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		pos, args = append(pos, args[0]), args[1:]
	}
	return pos, args
}

func runCompor(ctx context.Context, o comporOpts, d Deps) int {
	roots, err := refinarRoots(refinarOpts{vault: o.vault, docs: o.docs}, d)
	if err != nil {
		return fail(d, err)
	}
	out := d.Stdout
	if o.asJSON {
		out = d.Stderr // mantém o stdout como JSON puro
	}
	confirm, err := confirmarCompor(o, d, out)
	if err != nil {
		return fail(d, err)
	}
	req := compor.Request{Nome: o.nome, Ligacao: compor.Ligacao{Composto: o.composto, Titulo: o.titulo, Depende: o.depende}}
	res, err := compor.Run(ctx, req, compor.Deps{Store: processo.Store{Roots: roots}, Confirm: confirm})
	if err != nil {
		return falhaCompor(d, err)
	}
	return imprimirCompor(d, res, o.asJSON)
}

// confirmarCompor escolhe como o usuário confirma: --yes imprime o plano e segue; sem ele, o terminal pergunta.
func confirmarCompor(o comporOpts, d Deps, out io.Writer) (compor.ConfirmFunc, error) {
	if o.yes {
		return func(pl compor.Plano, _ func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
			fmt.Fprintf(out, "\n%s\n%s%s\n", pl.Diff(), pl.Cartao().Texto(), ste.Msg("compor.segue_yes"))
			return true, nil
		}, nil
	}
	if d.Stdin == nil {
		return nil, errors.New(ste.Msg("compor.sem_terminal"))
	}
	return newPrompter(d.Stdin, out).aprovarComposicao, nil
}

// falhaCompor traduz o erro: pedido inválido é uso incorreto (2); ciclo, processo inexistente, versão errada,
// cancelar e falha de execução saem com 1.
func falhaCompor(d Deps, err error) int {
	switch {
	case errors.Is(err, compor.ErrNome), errors.Is(err, compor.ErrPedido):
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	case errors.Is(err, compor.ErrCancelado):
		fmt.Fprintln(d.Stderr, ste.Msg("compor.cancelado"))
		return 1
	}
	return fail(d, err)
}

func imprimirCompor(d Deps, res *compor.Resultado, asJSON bool) int {
	if asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "\n%s\n\n", compor.MensagemConcluido(res))
	w := table(d.Stdout)
	fmt.Fprintln(w, "ARQUIVO\tESTADO\tCAMINHO\tBACKUP")
	for _, a := range res.Arquivos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Tipo, a.Estado, a.Caminho, ouTraco(a.Backup))
	}
	if err := w.Flush(); err != nil {
		return fail(d, err)
	}
	return 0
}

// aprovarComposicao mostra o plano e o cartão e conduz as três decisões. Enter vazio cancela. "Personalizar"
// pergunta o título e o depende e volta ao cartão.
func (p *prompter) aprovarComposicao(pl compor.Plano, replan func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
	fmt.Fprintf(p.w, "\n%s", pl.Diff())
	for invalidas := 0; invalidas < maxInvalid; {
		card := pl.Cartao()
		fmt.Fprintf(p.w, "\n%s", card.Render())
		in, err := p.line("> ")
		if err != nil {
			return false, err
		}
		dec, ok := ste.ParseDecisao(in)
		switch {
		case !ok:
			invalidas++
			fmt.Fprintln(p.w, ste.Msg("compor.entrada_invalida"))
		case dec == ste.Seguir:
			return true, nil
		case dec == ste.Personalizar:
			if pl, err = p.personalizarComposicao(card, pl, replan); err != nil {
				return false, err
			}
			fmt.Fprintf(p.w, "\n%s", pl.Diff())
		default:
			return false, nil
		}
	}
	return false, errTooManyInvalid
}

// personalizarComposicao pergunta o novo valor de cada parâmetro (Enter mantém) e devolve o plano recalculado.
// Se a ligação nova não vale, mostra o motivo e mantém o plano anterior.
func (p *prompter) personalizarComposicao(card ste.Cartao, pl compor.Plano, replan func(compor.Ligacao) (compor.Plano, error)) (compor.Plano, error) {
	l := pl.Ligacao
	for _, par := range card.Parametros {
		fmt.Fprintln(p.w, par.Ajuda)
		in, err := p.line(fmt.Sprintf("%s [%s]: ", par.Nome, par.Valor))
		if err != nil {
			return pl, err
		}
		switch {
		case in == "":
		case par.Nome == "titulo":
			l.Titulo = in
		case par.Nome == "depende":
			l.Depende = parseNomes(in)
		}
	}
	novo, err := replan(l)
	if err != nil {
		fmt.Fprintln(p.w, "erro:", err)
		return pl, nil
	}
	return novo, nil
}
