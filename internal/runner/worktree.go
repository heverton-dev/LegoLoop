package runner

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// MergeAST realiza a fusão determinística de duas árvores sintáticas Go (Motor II).
// Anexa declarações do arquivo incoming ao arquivo base preservando pacotes.
func MergeAST(baseSrc, incomingSrc []byte) ([]byte, error) {
	fset := token.NewFileSet()
	baseAST, err := parser.ParseFile(fset, "base.go", baseSrc, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("worktree merge: erro ao analisar base: %w", err)
	}

	incomingAST, err := parser.ParseFile(fset, "incoming.go", incomingSrc, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("worktree merge: erro ao analisar incoming: %w", err)
	}

	mergeDecls(baseAST, incomingAST)

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, baseAST); err != nil {
		return nil, fmt.Errorf("worktree merge: erro ao formatar código mesclado: %w", err)
	}

	return buf.Bytes(), nil
}

func mergeDecls(base, incoming *ast.File) {
	existing := make(map[string]bool)
	for _, decl := range base.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			existing[fn.Name.Name] = true
		}
	}

	for _, decl := range incoming.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if !existing[fn.Name.Name] {
				base.Decls = append(base.Decls, decl)
			}
		}
	}
}
