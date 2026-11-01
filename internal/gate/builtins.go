package gate

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// RootFS é o subconjunto de *os.Root usado pelos builtins: todo acesso fica
// confinado ao diretório do nó (bloqueia "..", caminhos absolutos e symlinks que escapam).
type RootFS interface {
	Stat(name string) (fs.FileInfo, error)
	Open(name string) (*os.File, error)
}

// Limite de leitura por arquivo de gate, evita estouro de memória em artefatos gigantes.
const maxGateFileBytes = 16 << 20

var (
	checklistRe  = regexp.MustCompile(`(?m)^\s*[-*]\s+\[[ xX]\]\s+\S`)
	acceptanceRe = regexp.MustCompile(`(?im)^#{1,6}\s+.*(crit[ée]rios?\s+de\s+aceita|acceptance\s+criteria)`)
)

type builtin struct {
	args []valKind
	ret  valKind
	fn   func(env Env, args []value, re *regexp.Regexp) (value, error)
}

func builtins() map[string]builtin {
	return map[string]builtin{
		"file_exists":             {[]valKind{kStr}, kBool, fileExists},
		"word_count":              {[]valKind{kStr}, kNum, countWith(bytes.Fields)},
		"line_count":              {[]valKind{kStr}, kNum, countWith(splitLines)},
		"has_pattern":             {[]valKind{kStr, kStr}, kBool, hasPattern},
		"has_binary_criteria":     {[]valKind{kStr}, kBool, matchFile(checklistRe)},
		"has_acceptance_criteria": {[]valKind{kStr}, kBool, hasAcceptance},
		"schema_valid":            {[]valKind{kStr, kStr}, kBool, schemaValid},
		"interactive":             {[]valKind{kStr}, kBool, interactivePrompt},
		"prompt":                  {[]valKind{kStr}, kBool, interactivePrompt},
	}
}

type callNode struct {
	name string
	b    builtin
	args []node
	re   *regexp.Regexp // pré-compilada quando o padrão é literal
}

func newCall(name string, args []node) (node, error) {
	b, ok := builtins()[name]
	if !ok {
		return nil, fmt.Errorf("gate: função desconhecida %q", name)
	}
	if len(args) != len(b.args) {
		return nil, fmt.Errorf("gate: %s espera %d argumento(s), recebeu %d", name, len(b.args), len(args))
	}
	for i, a := range args {
		if a.kind() != b.args[i] {
			return nil, fmt.Errorf("gate: %s argumento %d deve ser %s", name, i+1, b.args[i])
		}
	}
	return compileCallLiterals(callNode{name: name, b: b, args: args})
}

// compileCallLiterals valida caminhos e regex literais já na compilação.
func compileCallLiterals(c callNode) (node, error) {
	if c.name == "interactive" || c.name == "prompt" {
		return c, nil
	}
	if lit, ok := c.args[0].(litNode); ok && !filepath.IsLocal(lit.v.s) {
		return nil, fmt.Errorf("gate: caminho %q fora do workdir", lit.v.s)
	}
	if c.name == "schema_valid" {
		if lit, ok := c.args[1].(litNode); ok && !filepath.IsLocal(lit.v.s) {
			return nil, fmt.Errorf("gate: caminho %q fora do workdir", lit.v.s)
		}
	}
	if c.name == "has_pattern" {
		lit, ok := c.args[1].(litNode)
		if !ok {
			return nil, errors.New("gate: has_pattern exige regex literal")
		}
		re, err := regexp.Compile(lit.v.s)
		if err != nil {
			return nil, fmt.Errorf("gate: regex inválida: %w", err)
		}
		c.re = re
	}
	return c, nil
}

func (c callNode) kind() valKind { return c.b.ret }

func (c callNode) eval(env Env) (value, error) {
	vals := make([]value, len(c.args))
	for i, a := range c.args {
		v, err := a.eval(env)
		if err != nil {
			return value{}, err
		}
		vals[i] = v
	}
	return c.b.fn(env, vals, c.re)
}

func fileExists(env Env, args []value, _ *regexp.Regexp) (value, error) {
	info, err := env.Root.Stat(filepath.FromSlash(args[0].s))
	if errors.Is(err, fs.ErrNotExist) {
		return value{kind: kBool}, nil
	}
	if err != nil {
		return value{}, fmt.Errorf("gate: file_exists(%q): %w", args[0].s, err)
	}
	return value{kind: kBool, b: info.Mode().IsRegular()}, nil
}

// readGateFile abre uma única vez, exige arquivo regular e lê no máximo o limite
// (sem janela entre Stat e leitura; FIFO/diretório recusados).
func readGateFile(env Env, name string) ([]byte, error) {
	f, err := env.Root.Open(filepath.FromSlash(name))
	if err != nil {
		return nil, fmt.Errorf("gate: %q: %w", name, err)
	}
	defer f.Close()
	if info, err := f.Stat(); err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("gate: %q não é arquivo regular", name)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxGateFileBytes+1))
	if err == nil && len(data) > maxGateFileBytes {
		err = fmt.Errorf("gate: %q excede %d bytes", name, maxGateFileBytes)
	}
	return data, err
}

func countWith(split func([]byte) [][]byte) func(Env, []value, *regexp.Regexp) (value, error) {
	return func(env Env, args []value, _ *regexp.Regexp) (value, error) {
		data, err := readGateFile(env, args[0].s)
		if err != nil {
			return value{}, err
		}
		return value{kind: kNum, n: float64(len(split(data)))}, nil
	}
}

func splitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	return bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n"))
}

func hasPattern(env Env, args []value, re *regexp.Regexp) (value, error) {
	return matchFile(re)(env, args, nil)
}

func matchFile(re *regexp.Regexp) func(Env, []value, *regexp.Regexp) (value, error) {
	return func(env Env, args []value, _ *regexp.Regexp) (value, error) {
		data, err := readGateFile(env, args[0].s)
		if err != nil {
			return value{}, err
		}
		return value{kind: kBool, b: re.Match(data)}, nil
	}
}

// hasAcceptance exige seção "Critérios de Aceitação" e ao menos um item binário.
func hasAcceptance(env Env, args []value, _ *regexp.Regexp) (value, error) {
	data, err := readGateFile(env, args[0].s)
	if err != nil {
		return value{}, err
	}
	ok := acceptanceRe.Match(data) && checklistRe.Match(data)
	return value{kind: kBool, b: ok}, nil
}

func schemaValid(env Env, args []value, _ *regexp.Regexp) (value, error) {
	schemaData, err := readGateFile(env, args[0].s)
	if err != nil {
		return value{kind: kBool, b: false}, nil
	}
	targetData, err := readGateFile(env, args[1].s)
	if err != nil {
		return value{kind: kBool, b: false}, nil
	}
	return value{kind: kBool, b: ValidateJSONSchema(schemaData, targetData)}, nil
}

func interactivePrompt(env Env, args []value, _ *regexp.Regexp) (value, error) {
	if env.Prompter == nil {
		return value{kind: kBool, b: false}, errors.New("gate: interactive exige prompter configurado")
	}
	ok, err := env.Prompter.Prompt(args[0].s)
	if err != nil {
		return value{kind: kBool, b: false}, fmt.Errorf("gate: prompt falhou: %w", err)
	}
	return value{kind: kBool, b: ok}, nil
}


