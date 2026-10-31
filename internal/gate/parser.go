package gate

import (
	"fmt"
	"strconv"
)

type valKind int

const (
	kBool valKind = iota
	kNum
	kStr
)

func (k valKind) String() string {
	return [...]string{"bool", "número", "string"}[k]
}

type value struct {
	kind valKind
	b    bool
	n    float64
	s    string
}

// node é um nó tipado da AST; o tipo é verificado na compilação.
type node interface {
	kind() valKind
	eval(env Env) (value, error)
}

type parser struct {
	toks  []token
	pos   int
	depth int
}

// enter limita o aninhamento ("(" e "!") para a recursão nunca estourar a pilha.
func (p *parser) enter() error {
	p.depth++
	if p.depth > MaxDepth {
		return fmt.Errorf("gate: aninhamento excede %d níveis", MaxDepth)
	}
	return nil
}

func (p *parser) peek() token { return p.toks[p.pos] }

func (p *parser) next() token {
	t := p.toks[p.pos]
	if t.kind != tkEOF {
		p.pos++
	}
	return t
}

func (p *parser) acceptOp(ops ...string) (string, bool) {
	t := p.peek()
	if t.kind != tkOp {
		return "", false
	}
	for _, op := range ops {
		if t.text == op {
			p.pos++
			return op, true
		}
	}
	return "", false
}

func (p *parser) parseOr() (node, error) {
	return p.parseLogical("||", p.parseAnd)
}

func (p *parser) parseAnd() (node, error) {
	return p.parseLogical("&&", p.parseUnary)
}

func (p *parser) parseLogical(op string, sub func() (node, error)) (node, error) {
	left, err := sub()
	if err != nil {
		return nil, err
	}
	for {
		if _, ok := p.acceptOp(op); !ok {
			return left, nil
		}
		right, err := sub()
		if err != nil {
			return nil, err
		}
		if left.kind() != kBool || right.kind() != kBool {
			return nil, fmt.Errorf("gate: %q exige operandos bool", op)
		}
		left = logicNode{op: op, l: left, r: right}
	}
}

func (p *parser) parseUnary() (node, error) {
	if _, ok := p.acceptOp("!"); ok {
		if err := p.enter(); err != nil {
			return nil, err
		}
		defer func() { p.depth-- }()
		inner, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if inner.kind() != kBool {
			return nil, fmt.Errorf("gate: '!' exige bool, recebeu %s", inner.kind())
		}
		return notNode{inner}, nil
	}
	return p.parseCmp()
}

func (p *parser) parseCmp() (node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	op, ok := p.acceptOp("==", "!=", ">", ">=", "<", "<=")
	if !ok {
		return left, nil
	}
	right, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	if left.kind() != kNum || right.kind() != kNum {
		return nil, fmt.Errorf("gate: comparação %q exige números", op)
	}
	return cmpNode{op: op, l: left, r: right}, nil
}

func (p *parser) parseTerm() (node, error) {
	t := p.next()
	switch t.kind {
	case tkNumber:
		n, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return nil, fmt.Errorf("gate: número inválido %q", t.text)
		}
		return litNode{value{kind: kNum, n: n}}, nil
	case tkString:
		return litNode{value{kind: kStr, s: t.text}}, nil
	case tkIdent:
		return p.parseIdent(t)
	case tkLParen:
		if err := p.enter(); err != nil {
			return nil, err
		}
		defer func() { p.depth-- }()
		inner, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.next().kind != tkRParen {
			return nil, fmt.Errorf("gate: ')' esperado")
		}
		return inner, nil
	}
	return nil, fmt.Errorf("gate: termo inesperado %q na posição %d", t.text, t.pos)
}

func (p *parser) parseIdent(t token) (node, error) {
	switch t.text {
	case "true", "false":
		return litNode{value{kind: kBool, b: t.text == "true"}}, nil
	case "exit_code":
		return exitCodeNode{}, nil
	}
	if p.peek().kind != tkLParen {
		return nil, fmt.Errorf("gate: identificador desconhecido %q", t.text)
	}
	p.next()
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}
	return newCall(t.text, args)
}

func (p *parser) parseArgs() ([]node, error) {
	var args []node
	if p.peek().kind == tkRParen {
		p.next()
		return args, nil
	}
	for {
		arg, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		switch p.next().kind {
		case tkRParen:
			return args, nil
		case tkComma:
			continue
		default:
			return nil, fmt.Errorf("gate: ',' ou ')' esperado nos argumentos")
		}
	}
}
