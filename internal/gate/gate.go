package gate

import (
	"errors"
	"fmt"
	"strings"
)

// Prompter solicita confirmação interativa para gates humanos (Stop-Gates).
type Prompter interface {
	Prompt(message string) (bool, error)
}

// Env é o contexto de avaliação: arquivos confinados ao Root do workdir do nó,
// o código de saída do executor e o Prompter para Stop-Gates interativos.
type Env struct {
	Root     RootFS
	ExitCode int
	Prompter Prompter
}

// Expr é uma expressão de gate compilada e com tipos verificados.
type Expr struct {
	src  string
	root node
}

// Limites contra expressões abusivas (recursão profunda / payload gigante).
const (
	MaxExprBytes = 4 << 10
	MaxDepth     = 64
)

// Compile faz o parse e a verificação de tipos; a expressão final deve ser bool.
func Compile(src string) (*Expr, error) {
	if len(src) > MaxExprBytes {
		return nil, fmt.Errorf("gate: expressão excede %d bytes", MaxExprBytes)
	}
	src = strings.TrimSpace(src)
	if src == "" {
		return nil, errors.New("gate: expressão vazia")
	}
	toks, err := lex(src)
	if err != nil {
		return nil, err
	}
	p := &parser{toks: toks}
	root, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	if t := p.peek(); t.kind != tkEOF {
		return nil, fmt.Errorf("gate: token excedente %q na posição %d", t.text, t.pos)
	}
	if root.kind() != kBool {
		return nil, fmt.Errorf("gate: expressão deve ser bool, é %s", root.kind())
	}
	return &Expr{src: src, root: root}, nil
}

// Eval avalia a expressão no ambiente informado.
func (e *Expr) Eval(env Env) (bool, error) {
	if env.Root == nil {
		return false, errors.New("gate: Env.Root ausente")
	}
	v, err := e.root.eval(env)
	if err != nil {
		return false, err
	}
	return v.b, nil
}

// String devolve a expressão normalizada.
func (e *Expr) String() string { return e.src }

type litNode struct{ v value }

func (n litNode) kind() valKind                { return n.v.kind }
func (n litNode) eval(Env) (value, error)      { return n.v, nil }
func (exitCodeNode) kind() valKind             { return kNum }
func (exitCodeNode) eval(e Env) (value, error) { return value{kind: kNum, n: float64(e.ExitCode)}, nil }

type exitCodeNode struct{}

type notNode struct{ inner node }

func (notNode) kind() valKind { return kBool }
func (n notNode) eval(env Env) (value, error) {
	v, err := n.inner.eval(env)
	return value{kind: kBool, b: !v.b}, err
}

type logicNode struct {
	op   string
	l, r node
}

func (logicNode) kind() valKind { return kBool }

// eval usa curto-circuito: o lado direito só é avaliado se necessário.
func (n logicNode) eval(env Env) (value, error) {
	l, err := n.l.eval(env)
	if err != nil {
		return value{}, err
	}
	if (n.op == "&&" && !l.b) || (n.op == "||" && l.b) {
		return l, nil
	}
	return n.r.eval(env)
}

type cmpNode struct {
	op   string
	l, r node
}

func (cmpNode) kind() valKind { return kBool }
func (n cmpNode) eval(env Env) (value, error) {
	l, err := n.l.eval(env)
	if err != nil {
		return value{}, err
	}
	r, err := n.r.eval(env)
	if err != nil {
		return value{}, err
	}
	results := map[string]bool{
		"==": l.n == r.n, "!=": l.n != r.n,
		">": l.n > r.n, ">=": l.n >= r.n,
		"<": l.n < r.n, "<=": l.n <= r.n,
	}
	return value{kind: kBool, b: results[n.op]}, nil
}
