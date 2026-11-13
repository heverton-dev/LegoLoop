package dissecar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// ExtractAST analisa arquivos Go e extrai funções e tipos.
func ExtractAST(dir string) ([]string, error) {
	var nodes []string
	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Name != nil {
						nodes = append(nodes, "func:"+d.Name.Name)
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						if ts, ok := spec.(*ast.TypeSpec); ok {
							nodes = append(nodes, "type:"+ts.Name.Name)
						}
					}
				}
			}
		}
		return nil
	})

	return nodes, err
}

// ExtractSyscalls simula extração de chamadas de sistema ou imports de pacotes.
func ExtractSyscalls(dir string) ([]string, error) {
	var syscalls []string
	fset := token.NewFileSet()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}
			for _, imp := range file.Imports {
				if imp.Path != nil {
					syscalls = append(syscalls, stringsTrimQuotes(imp.Path.Value))
				}
			}
		}
		return nil
	})

	return uniqueStrings(syscalls), err
}

func stringsTrimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func uniqueStrings(input []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range input {
		if !keys[entry] {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// BuildReverseDAG sintetiza o DAG reverso a partir de nós.
func BuildReverseDAG(nodes []string) map[string][]string {
	dag := make(map[string][]string)
	for i, node := range nodes {
		if i > 0 {
			dag[node] = []string{nodes[i-1]}
		} else {
			dag[node] = []string{}
		}
	}
	return dag
}
