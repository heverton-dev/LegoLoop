package espinha

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var tiposPredeclarados = map[string]bool{
	"any": true, "bool": true, "byte": true, "comparable": true,
	"complex64": true, "complex128": true, "error": true,
	"float32": true, "float64": true, "int": true, "int8": true,
	"int16": true, "int32": true, "int64": true, "rune": true,
	"string": true, "uint": true, "uint8": true, "uint16": true,
	"uint32": true, "uint64": true, "uintptr": true,
}

func verificarFuncoesDuplicadas(raiz string) []Violacao {
	dirInternal := filepath.Join(raiz, "internal")
	if _, err := os.Stat(dirInternal); err != nil {
		return nil
	}
	ocorrencias, err := coletarFuncoesInternal(raiz, dirInternal)
	if err != nil {
		return []Violacao{{
			Regra:  "funcao-duplicada-entre-pacotes",
			Alvo:   "internal",
			Motivo: err.Error(),
		}}
	}
	return construirViolacoesDuplicadas(ocorrencias)
}

func coletarFuncoesInternal(raiz, dirInternal string) (map[string]map[string]bool, error) {
	ocorrencias := make(map[string]map[string]bool)
	fset := token.NewFileSet()
	err := filepath.WalkDir(dirInternal, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() || !ehArquivoGoProducao(d.Name()) {
			return nil
		}
		return registrarFuncoesArquivo(fset, raiz, caminho, ocorrencias)
	})
	return ocorrencias, err
}

func ehArquivoGoProducao(nome string) bool {
	return strings.HasSuffix(nome, ".go") && !strings.HasSuffix(nome, "_test.go")
}

func registrarFuncoesArquivo(fset *token.FileSet, raiz, caminho string, ocorrencias map[string]map[string]bool) error {
	arq, err := parser.ParseFile(fset, caminho, nil, 0)
	if err != nil {
		return fmt.Errorf("analisar %s: %w", caminhoRelativo(raiz, caminho), err)
	}
	pacote := caminhoRelativo(raiz, filepath.Dir(caminho))
	for _, decl := range arq.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !ehFuncaoVerificavel(fn) {
			continue
		}
		chave := fn.Name.Name + assinaturaCanonica(fset, pacote, fn.Type)
		if ocorrencias[chave] == nil {
			ocorrencias[chave] = make(map[string]bool)
		}
		ocorrencias[chave][pacote] = true
	}
	return nil
}

func ehFuncaoVerificavel(fn *ast.FuncDecl) bool {
	if fn.Recv != nil || fn.Name == nil {
		return false
	}
	nome := fn.Name.Name
	if nome == "init" || nome == "main" || !ast.IsExported(nome) {
		return false
	}
	return true
}

func assinaturaCanonica(fset *token.FileSet, pacote string, ft *ast.FuncType) string {
	tparams := extrairNomesTypeParams(ft.TypeParams)
	var b strings.Builder
	if len(tparams) > 0 {
		b.WriteByte('[')
		b.WriteString(formatarListaTipos(fset, pacote, ft.TypeParams, tparams))
		b.WriteByte(']')
	}
	b.WriteByte('(')
	b.WriteString(formatarListaTipos(fset, pacote, ft.Params, tparams))
	b.WriteString(")->(")
	b.WriteString(formatarListaTipos(fset, pacote, ft.Results, tparams))
	b.WriteByte(')')
	return b.String()
}

func extrairNomesTypeParams(fl *ast.FieldList) map[string]bool {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}
	res := make(map[string]bool)
	for _, f := range fl.List {
		for _, n := range f.Names {
			res[n.Name] = true
		}
	}
	return res
}

func formatarListaTipos(fset *token.FileSet, pacote string, fl *ast.FieldList, tparams map[string]bool) string {
	if fl == nil {
		return ""
	}
	var tipos []string
	for _, f := range fl.List {
		t := formatarExprTipo(fset, pacote, f.Type, tparams)
		n := max(len(f.Names), 1)
		for range n {
			tipos = append(tipos, t)
		}
	}
	return strings.Join(tipos, ",")
}

func formatarExprTipo(fset *token.FileSet, pacote string, expr ast.Expr, tparams map[string]bool) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if tiposPredeclarados[t.Name] || tparams[t.Name] {
			return t.Name
		}
		return pacote + "." + t.Name
	case *ast.StarExpr:
		return "*" + formatarExprTipo(fset, pacote, t.X, tparams)
	case *ast.ArrayType:
		return formatarArrayTipo(fset, pacote, t, tparams)
	case *ast.MapType:
		return "map[" + formatarExprTipo(fset, pacote, t.Key, tparams) + "]" + formatarExprTipo(fset, pacote, t.Value, tparams)
	case *ast.Ellipsis:
		return "..." + formatarExprTipo(fset, pacote, t.Elt, tparams)
	case *ast.FuncType:
		return "func" + assinaturaCanonica(fset, pacote, t)
	default:
		return formatarNoBruto(fset, expr)
	}
}

func formatarArrayTipo(fset *token.FileSet, pacote string, t *ast.ArrayType, tparams map[string]bool) string {
	if t.Len == nil {
		return "[]" + formatarExprTipo(fset, pacote, t.Elt, tparams)
	}
	return "[" + formatarNoBruto(fset, t.Len) + "]" + formatarExprTipo(fset, pacote, t.Elt, tparams)
}

func formatarNoBruto(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return strings.ReplaceAll(buf.String(), " ", "")
}

func construirViolacoesDuplicadas(ocorrencias map[string]map[string]bool) []Violacao {
	var vs []Violacao
	for chave, pacsMap := range ocorrencias {
		if len(pacsMap) < 2 {
			continue
		}
		pacs := make([]string, 0, len(pacsMap))
		for p := range pacsMap {
			pacs = append(pacs, p)
		}
		slices.Sort(pacs)
		vs = append(vs, Violacao{
			Regra:  "funcao-duplicada-entre-pacotes",
			Alvo:   chave,
			Motivo: "funcao duplicada em pacotes distintos: " + strings.Join(pacs, ", "),
		})
	}
	return vs
}
