package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"legoloop/internal/fabricar"
	"legoloop/internal/interview"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("fabricar.confianca", "Confiança: %d%% (meta: %d%%).")
	ste.Default.Register("fabricar.confirmar_rascunho", "Gerar a SPEC mesmo assim?")
}

// metaConfianca é a confiança a partir da qual a entrevista pode terminar.
const metaConfianca = fabricar.ConfidenceTarget

// maxInvalid é quantas entradas inválidas seguidas o terminal tolera antes de abortar.
const maxInvalid = 3

var (
	errTooManyInvalid = fmt.Errorf("entrada inválida %d vezes seguidas", maxInvalid)
	errInputClosed    = errors.New("entrada encerrada (EOF)")
)

// prompter conduz o diálogo no terminal. Usa um único bufio.Reader: um leitor por etapa perderia as linhas
// que o anterior já deixou no buffer.
type prompter struct {
	r *bufio.Reader
	w io.Writer
}

func newPrompter(in io.Reader, out io.Writer) *prompter {
	return &prompter{r: bufio.NewReader(in), w: out}
}

// line imprime o prompt e lê uma linha. EOF sem texto aborta; EOF com texto parcial ainda vale como resposta.
func (p *prompter) line(prompt string) (string, error) {
	fmt.Fprint(p.w, prompt)
	s, err := p.r.ReadString('\n')
	if err != nil && strings.TrimSpace(s) == "" {
		return "", errInputClosed
	}
	return strings.TrimSpace(s), nil
}

// setupInteractive liga o terminal ao motor: drives, entrevista e confirmação final usam o mesmo prompter.
func setupInteractive(ctx context.Context, req *fabricar.Request, fd *fabricar.Deps, d Deps) error {
	if d.Stdin == nil {
		return errors.New("o modo interativo exige um terminal (stdin); use --yes para o modo autônomo")
	}
	p := newPrompter(d.Stdin, d.Stdout)
	fd.Ask, fd.Aprovar, fd.ConfirmDraft = p.ask, p.aprovar, p.confirmDraft
	fd.Context = interview.Context{Goal: goalOf(*req), Docs: interview.LoadProjectContext(fd.Roots.Docs)}
	fd.OnTurn = func(c int) { fmt.Fprintln(d.Stdout, ste.Msg("fabricar.confianca", c, metaConfianca)) }
	return nil
}

// ask apresenta a pergunta da entrevista com opções numeradas: aceitar o palpite, responder em texto livre ou cancelar.
func (p *prompter) ask(question, guess string) (string, error) {
	accept, free := 1, 2
	if strings.TrimSpace(guess) == "" {
		accept, free = 0, 1 // sem palpite não há o que aceitar
	}
	for attempt := 0; attempt < maxInvalid; attempt++ {
		p.showQuestion(question, guess, accept, free)
		in, err := p.line("Escolha: ")
		if err != nil {
			return "", err
		}
		switch n, convErr := strconv.Atoi(in); {
		case convErr != nil:
		case n == 0:
			return "", fabricar.ErrAborted
		case n == accept && accept != 0:
			return guess, nil
		case n == free:
			if text, err := p.freeText(); err != nil || text != "" {
				return text, err
			}
		}
		fmt.Fprintln(p.w, "Entrada inválida.")
	}
	return "", errTooManyInvalid
}

func (p *prompter) showQuestion(question, guess string, accept, free int) {
	fmt.Fprintf(p.w, "\n%s\n", question)
	if accept != 0 {
		fmt.Fprintf(p.w, "  [%d] Aceitar o palpite: %s\n", accept, guess)
	}
	fmt.Fprintf(p.w, "  [%d] Responder com minhas palavras\n  [0] Cancelar\n", free)
}

func (p *prompter) freeText() (string, error) { return p.line("Sua resposta: ") }

// confirmDraft mostra a confiança medida e as lacunas e pergunta se a SPEC deve ser gerada assim mesmo.
func (p *prompter) confirmDraft(c fabricar.Clarification) (bool, error) {
	for _, w := range c.Warnings() {
		fmt.Fprintln(p.w, w)
	}
	return p.yesNo(ste.Msg("fabricar.confirmar_rascunho") + " [s/N]: ")
}

// yesNo pede [s/N]. Enter vale N; resposta fora do padrão reapresenta, até maxInvalid vezes.
func (p *prompter) yesNo(prompt string) (bool, error) {
	for attempt := 0; attempt < maxInvalid; attempt++ {
		in, err := p.line(prompt)
		if err != nil {
			return false, err
		}
		switch strings.ToLower(in) {
		case "s", "sim", "y", "yes":
			return true, nil
		case "", "n", "nao", "não", "no":
			return false, nil
		}
		fmt.Fprintln(p.w, "Entrada inválida: responda s ou n.")
	}
	return false, errTooManyInvalid
}

// goalOf é o objetivo que a entrevista recebe: --goal ou, sem ele, o próprio nome do processo.
func goalOf(req fabricar.Request) string {
	if g := strings.TrimSpace(req.Goal); g != "" {
		return g
	}
	return strings.TrimSpace(req.RawName)
}
