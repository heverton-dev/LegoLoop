// Comando funclen é o gate "funções pequenas" (AGENTS/PROMPT: < 35 linhas): lista funções
// de produção (sem _test.go) cujo corpo excede o limite e sai com código 1 se houver alguma.
//
// Uso: go run ./tools/funclen [dir] [limite]
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir, limit := "internal", 35
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil {
			limit = n
		}
	}
	found, err := scan(dir, limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(2)
	}
	for _, f := range found {
		fmt.Println(f)
	}
	if len(found) > 0 {
		os.Exit(1)
	}
}

func scan(dir string, limit int) ([]string, error) {
	var found []string
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return err
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return err
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			if n := fset.Position(fn.Body.Rbrace).Line - fset.Position(fn.Body.Lbrace).Line - 1; n > limit {
				found = append(found, fmt.Sprintf("%s:%d %s %d linhas", p, fset.Position(fn.Pos()).Line, fn.Name.Name, n))
			}
		}
		return nil
	})
	return found, err
}
