package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"legoloop/internal/processo"
	"legoloop/internal/refinar"
	"legoloop/internal/ste"
)

func init() {
	for id, msg := range map[string]string{
		"refinar.uso":          "Uso: legoloop refinar <nome> --etapa ID [--peca P] [--caso C] [--titulo T] [--vault DIR] [--docs DIR] [--yes] [--json]",
		"refinar.sem_terminal": "Não há terminal para confirmar. Use --yes para seguir sem perguntar.",
		"refinar.segue_yes":    "Segue sem perguntar, por causa do --yes.",
	} {
		ste.Default.Register(id, msg)
	}
}

// refinarOpts reúne as flags do comando.
type refinarOpts struct {
	yes, asJSON                      bool
	vault, docs, etapa, peca, titulo string
	caso                             listFlag
}

// cmdRefinar muda uma etapa de um processo do cofre. Mostra o diff e o cartão, pede a confirmação e grava.
func cmdRefinar(ctx context.Context, args []string, d Deps) int {
	var o refinarOpts
	fs := newFlags("refinar", d)
	fs.BoolVar(&o.yes, "yes", false, "segue sem perguntar")
	fs.BoolVar(&o.asJSON, "json", false, "imprime o resultado em JSON")
	fs.StringVar(&o.vault, "vault", "", "cofre canônico (padrão ~/.aidd/canonical)")
	fs.StringVar(&o.docs, "docs", "", "raiz do projeto que tem docs/ (padrão: diretório atual)")
	fs.StringVar(&o.etapa, "etapa", "", "id da etapa que muda (obrigatório)")
	fs.StringVar(&o.peca, "peca", "", "peça nova da etapa")
	fs.StringVar(&o.titulo, "titulo", "", "título novo da etapa")
	fs.Var(&o.caso, "caso", "caso de uso da etapa (repetível)")
	nome, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return 2
	}
	if nome == "" {
		nome = fs.Arg(0)
	}
	if nome == "" || o.etapa == "" {
		fmt.Fprintln(d.Stderr, ste.Msg("refinar.uso"))
		return 2
	}
	return runRefinar(ctx, nome, o, d)
}

func runRefinar(ctx context.Context, nome string, o refinarOpts, d Deps) int {
	roots, err := refinarRoots(o, d)
	if err != nil {
		return fail(d, err)
	}
	out := d.Stdout
	if o.asJSON {
		out = d.Stderr // mantém o stdout como JSON puro
	}
	confirm, err := confirmarRefinar(o, d, out)
	if err != nil {
		return fail(d, err)
	}
	req := refinar.Request{Nome: nome, Mudanca: refinar.Mudanca{Etapa: o.etapa, Peca: o.peca, Titulo: o.titulo, Casos: o.caso}}
	res, err := refinar.Run(ctx, req, refinar.Deps{Store: processo.Store{Roots: roots}, Confirm: confirm})
	if err != nil {
		return falhaRefinar(d, err)
	}
	return imprimirRefinar(d, res, o.asJSON)
}

func refinarRoots(o refinarOpts, d Deps) (processo.Roots, error) {
	vault, _, err := vaultListar(o.vault, d.Home)
	if err != nil {
		return processo.Roots{}, err
	}
	docs := o.docs
	if docs == "" {
		docs = "."
	}
	abs, err := filepath.Abs(docs)
	return processo.Roots{Vault: vault, Docs: abs}, err
}

// confirmarRefinar escolhe como o usuário confirma: --yes imprime o plano e segue; sem ele, o terminal pergunta.
func confirmarRefinar(o refinarOpts, d Deps, out io.Writer) (refinar.ConfirmFunc, error) {
	if o.yes {
		return func(pl refinar.Plano, _ func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
			fmt.Fprintf(out, "\n%s\n%s%s\n", pl.Diff(), pl.Cartao().Texto(), ste.Msg("refinar.segue_yes"))
			return true, nil
		}, nil
	}
	if d.Stdin == nil {
		return nil, errors.New(ste.Msg("refinar.sem_terminal"))
	}
	return newPrompter(d.Stdin, out).aprovarRefino, nil
}

// falhaRefinar traduz o erro: pedido inválido é uso incorreto (2); cancelar e falha de execução saem com 1.
func falhaRefinar(d Deps, err error) int {
	switch {
	case errors.Is(err, refinar.ErrNome), errors.Is(err, refinar.ErrMudanca):
		fmt.Fprintln(d.Stderr, "erro:", err)
		return 2
	case errors.Is(err, refinar.ErrCancelado):
		fmt.Fprintln(d.Stderr, ste.Msg("refinar.cancelado"))
		return 1
	}
	return fail(d, err)
}

func imprimirRefinar(d Deps, res *refinar.Resultado, asJSON bool) int {
	if asJSON {
		return printJSON(d, res)
	}
	fmt.Fprintf(d.Stdout, "\n%s\n\n", ste.Msg("refinar.concluido", res.Nome, res.De, res.Para))
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

// aprovarRefino mostra o diff e o cartão e conduz as três decisões. Enter vazio cancela. "Personalizar" pergunta
// a peça e os casos e volta ao cartão.
func (p *prompter) aprovarRefino(pl refinar.Plano, replan func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
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
			fmt.Fprintln(p.w, ste.Msg("refinar.entrada_invalida"))
		case dec == ste.Seguir:
			return true, nil
		case dec == ste.Personalizar:
			if pl, err = p.personalizarRefino(card, pl, replan); err != nil {
				return false, err
			}
			fmt.Fprintf(p.w, "\n%s", pl.Diff())
		default:
			return false, nil
		}
	}
	return false, errTooManyInvalid
}

// personalizarRefino pergunta o novo valor de cada parâmetro (Enter mantém) e devolve o plano recalculado. Se a
// mudança nova não vale, mostra o motivo e mantém o plano anterior.
func (p *prompter) personalizarRefino(card ste.Cartao, pl refinar.Plano, replan func(refinar.Mudanca) (refinar.Plano, error)) (refinar.Plano, error) {
	m := pl.Mudanca
	for _, par := range card.Parametros {
		fmt.Fprintln(p.w, par.Ajuda)
		in, err := p.line(fmt.Sprintf("%s [%s]: ", par.Nome, par.Valor))
		if err != nil {
			return pl, err
		}
		switch {
		case in == "":
		case par.Nome == "peca":
			m.Peca = in
		case par.Nome == "casos":
			m.Casos = parseNomes(in)
		}
	}
	novo, err := replan(m)
	if err != nil {
		fmt.Fprintln(p.w, "erro:", err)
		return pl, nil
	}
	return novo, nil
}
