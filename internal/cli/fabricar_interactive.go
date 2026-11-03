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
)

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
	chosen, err := p.pickDrives(driveRoots(ctx, d))
	if err != nil {
		return err
	}
	req.TargetDrives = chosen
	fd.Ask, fd.Confirm = p.ask, p.confirm
	return nil
}

func driveRoots(ctx context.Context, d Deps) []string {
	if d.Fabricar != nil && d.Fabricar.DriveRoots != nil {
		return d.Fabricar.DriveRoots(ctx)
	}
	roots, _ := fabricar.SystemRoots(ctx)
	return roots
}

// pickDrives lista os drives detectados e deixa o usuário desmarcar alguns. Enter mantém todos.
func (p *prompter) pickDrives(drives []string) ([]string, error) {
	if len(drives) == 0 {
		fmt.Fprintln(p.w, "Nenhum drive detectado; a varredura usará as raízes padrão do sistema.")
		return nil, nil
	}
	for attempt := 0; attempt < maxInvalid; attempt++ {
		fmt.Fprintln(p.w, "Drives detectados:")
		for i, dr := range drives {
			fmt.Fprintf(p.w, "  [%d] %s\n", i+1, dr)
		}
		in, err := p.line("Desmarcar quais? (números separados por vírgula; Enter mantém todos): ")
		if err != nil {
			return nil, err
		}
		if kept, ok := keepDrives(drives, in); ok {
			return kept, nil
		}
		fmt.Fprintln(p.w, "Entrada inválida: use números da lista e deixe ao menos um drive.")
	}
	return nil, errTooManyInvalid
}

// keepDrives devolve os drives que sobram. ok é falso se algum número for inválido ou se nenhum drive sobrar.
func keepDrives(drives []string, input string) ([]string, bool) {
	drop := map[int]bool{}
	for _, f := range strings.FieldsFunc(input, func(r rune) bool { return r == ',' || r == ' ' }) {
		n, err := strconv.Atoi(f)
		if err != nil || n < 1 || n > len(drives) {
			return nil, false
		}
		drop[n-1] = true
	}
	var kept []string
	for i, dr := range drives {
		if !drop[i] {
			kept = append(kept, dr)
		}
	}
	return kept, len(kept) > 0
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

// confirm mostra o resumo dos caminhos e pede [s/N]. Enter vale N; resposta fora do padrão reapresenta.
func (p *prompter) confirm(paths []string) (bool, error) {
	fmt.Fprintln(p.w, "\nResumo: serão gravados")
	for _, path := range paths {
		fmt.Fprintf(p.w, "  %s\n", path)
	}
	for attempt := 0; attempt < maxInvalid; attempt++ {
		in, err := p.line("Gravar? [s/N]: ")
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
