package espinha

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// padraoDebito casa o marcador que adia uma correção: `TODO:`, `FIXME(autor):`, `LACUNA:` etc.
// É sensível a maiúsculas e exige `:` para não confundir com a palavra "todo" do português.
var padraoDebito = regexp.MustCompile(`^(TODO|FIXME|HACK|XXX|LACUNA|DÉBITO|DEBITO)(\([^)]*\))?:`)

// verificarDebitoMarcado aplica a Lei 16 (lacuna zero): todo achado se corrige na execução,
// então nenhum comentário em Go pode adiar correção. A regra vale em qualquer harness, porque roda no `go test`.
func verificarDebitoMarcado(raiz string) []Violacao {
	var vs []Violacao
	for _, sub := range []string{`internal`, `tools`, `cmd`} {
		vs = append(vs, debitosDoDiretorio(raiz, filepath.Join(raiz, sub))...)
	}
	return vs
}

func debitosDoDiretorio(raiz, dir string) []Violacao {
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	var vs []Violacao
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() && d.Name() == `testdata` {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), `.go`) {
			return nil
		}
		arq, err := parser.ParseFile(fset, caminho, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf(`analisar %s: %w`, caminhoRelativo(raiz, caminho), err)
		}
		vs = append(vs, debitosDoArquivo(fset, raiz, caminho, arq.Comments)...)
		return nil
	})
	if err != nil {
		vs = append(vs, Violacao{Regra: `debito-marcado`, Alvo: caminhoRelativo(raiz, dir), Motivo: err.Error()})
	}
	return vs
}

func debitosDoArquivo(fset *token.FileSet, raiz, caminho string, grupos []*ast.CommentGroup) []Violacao {
	var vs []Violacao
	for _, grupo := range grupos {
		for _, c := range grupo.List {
			texto := strings.TrimSpace(strings.Trim(c.Text, `/*`))
			if !padraoDebito.MatchString(texto) {
				continue
			}
			vs = append(vs, Violacao{
				Regra:  `debito-marcado`,
				Alvo:   fmt.Sprintf(`%s:%d`, caminhoRelativo(raiz, caminho), fset.Position(c.Pos()).Line),
				Motivo: `comentário adia correção (Lei 16): corrija na execução, com teste`,
			})
		}
	}
	return vs
}
