package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

// OperadoresPadrao retorna os 7 operadores de mutação da especificação.
func OperadoresPadrao() []Operador {
	return []Operador{
		operadorBinario("eq_neq", token.EQL, token.NEQ, nil),
		operadorBinario("lt_leq", token.LSS, token.LEQ, nil),
		operadorBinario("gt_geq", token.GTR, token.GEQ, nil),
		operadorBinario("and_or", token.LAND, token.LOR, nil),
		operadorBoolFlip(),
		operadorRemoverReturnErr(),
		operadorBinario("add_sub", token.ADD, token.SUB, filtrarSomaSemString),
	}
}

func montarID(pacote, arquivo string, linha int, operador string, indice int) string {
	return fmt.Sprintf("%s:%s:%d:%s:%d", pacote, arquivo, linha, operador, indice)
}

type filtroBinario func(arq *ast.File, expr *ast.BinaryExpr) bool

func operadorBinario(nome string, opA, opB token.Token, filtro filtroBinario) Operador {
	return Operador{
		Nome: nome,
		Gerar: func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error) {
			var muts []Mutante
			var errGen error
			indice := 0
			base := filepath.Base(caminho)
			ast.Inspect(arq, func(n ast.Node) bool {
				if errGen != nil {
					return false
				}
				bin, ok := n.(*ast.BinaryExpr)
				if !ok || (bin.Op != opA && bin.Op != opB) {
					return true
				}
				if filtro != nil && !filtro(arq, bin) {
					return true
				}
				orig := bin.Op
				novo := opB
				if orig == opB {
					novo = opA
				}
				bin.Op = novo
				codigo, err := formatarAST(fset, arq)
				bin.Op = orig
				if err != nil {
					errGen = err
					return false
				}
				indice++
				linha := fset.Position(bin.OpPos).Line
				muts = append(muts, Mutante{
					ID:       montarID(pacote, base, linha, nome, indice),
					Pacote:   pacote,
					Arquivo:  base,
					Linha:    linha,
					Operador: nome,
					Original: orig.String(),
					Mutado:   novo.String(),
					Codigo:   codigo,
				})
				return true
			})
			return muts, errGen
		},
	}
}

func operadorBoolFlip() Operador {
	const nome = "bool_flip"
	return Operador{
		Nome: nome,
		Gerar: func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error) {
			ignorados := coletarIdentsNaoValor(arq)
			var muts []Mutante
			var errGen error
			indice := 0
			base := filepath.Base(caminho)
			ast.Inspect(arq, func(n ast.Node) bool {
				if errGen != nil {
					return false
				}
				id, ok := n.(*ast.Ident)
				if !ok || ignorados[id] || (id.Name != "true" && id.Name != "false") {
					return true
				}
				orig := id.Name
				novo := "false"
				if orig == "false" {
					novo = "true"
				}
				id.Name = novo
				codigo, err := formatarAST(fset, arq)
				id.Name = orig
				if err != nil {
					errGen = err
					return false
				}
				indice++
				linha := fset.Position(id.NamePos).Line
				muts = append(muts, Mutante{
					ID:       montarID(pacote, base, linha, nome, indice),
					Pacote:   pacote,
					Arquivo:  base,
					Linha:    linha,
					Operador: nome,
					Original: orig,
					Mutado:   novo,
					Codigo:   codigo,
				})
				return true
			})
			return muts, errGen
		},
	}
}

func coletarIdentsNaoValor(arq *ast.File) map[*ast.Ident]bool {
	m := make(map[*ast.Ident]bool)
	ast.Inspect(arq, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			m[x.Sel] = true
		case *ast.KeyValueExpr:
			if id, ok := x.Key.(*ast.Ident); ok {
				m[id] = true
			}
		case *ast.Field:
			for _, id := range x.Names {
				m[id] = true
			}
		case *ast.ValueSpec:
			for _, id := range x.Names {
				m[id] = true
			}
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					m[id] = true
				}
			}
		case *ast.FuncDecl:
			m[x.Name] = true
		case *ast.TypeSpec:
			m[x.Name] = true
		case *ast.LabeledStmt:
			m[x.Label] = true
		case *ast.BranchStmt:
			if x.Label != nil {
				m[x.Label] = true
			}
		}
		return true
	})
	return m
}

func operadorRemoverReturnErr() Operador {
	const nome = "rm_return_err"
	return Operador{
		Nome: nome,
		Gerar: func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error) {
			ifSemElse := mapearIfSemElse(arq)
			var muts []Mutante
			var errGen error
			indice := 0
			base := filepath.Base(caminho)
			ast.Inspect(arq, func(n ast.Node) bool {
				if errGen != nil {
					return false
				}
				blk, ok := n.(*ast.BlockStmt)
				if !ok {
					return true
				}
				gerados, err := mutarBlocoReturnErr(fset, arq, blk, ifSemElse[blk], pacote, base, nome, &indice)
				if err != nil {
					errGen = err
					return false
				}
				muts = append(muts, gerados...)
				return true
			})
			return muts, errGen
		},
	}
}

func mapearIfSemElse(arq *ast.File) map[*ast.BlockStmt]bool {
	m := make(map[*ast.BlockStmt]bool)
	ast.Inspect(arq, func(n ast.Node) bool {
		if ifStmt, ok := n.(*ast.IfStmt); ok && ifStmt.Else == nil && ifStmt.Body != nil {
			m[ifStmt.Body] = true
		}
		return true
	})
	return m
}

func mutarBlocoReturnErr(
	fset *token.FileSet,
	arq *ast.File,
	blk *ast.BlockStmt,
	emIfSemElse bool,
	pacote, base, nome string,
	indice *int,
) ([]Mutante, error) {
	var muts []Mutante
	origList := blk.List
	for i, stmt := range origList {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		posErr := indiceIdentErr(ret.Results)
		if posErr < 0 {
			continue
		}
		codigo, mutDesc, err := aplicarMutacaoReturnErr(fset, arq, blk, origList, i, ret, posErr, emIfSemElse)
		if err != nil {
			return nil, err
		}
		*indice++
		linha := fset.Position(ret.Return).Line
		muts = append(muts, Mutante{
			ID:       montarID(pacote, base, linha, nome, *indice),
			Pacote:   pacote,
			Arquivo:  base,
			Linha:    linha,
			Operador: nome,
			Original: "return err",
			Mutado:   mutDesc,
			Codigo:   codigo,
		})
	}
	return muts, nil
}

func aplicarMutacaoReturnErr(
	fset *token.FileSet,
	arq *ast.File,
	blk *ast.BlockStmt,
	origList []ast.Stmt,
	idx int,
	ret *ast.ReturnStmt,
	posErr int,
	emIfSemElse bool,
) ([]byte, string, error) {
	descarte := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("_")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{ast.NewIdent("err")},
	}
	if emIfSemElse {
		nova := make([]ast.Stmt, len(origList))
		copy(nova, origList)
		nova[idx] = descarte
		blk.List = nova
		codigo, err := formatarAST(fset, arq)
		blk.List = origList
		return codigo, "_ = err", err
	}
	origExpr := ret.Results[posErr]
	ret.Results[posErr] = ast.NewIdent("nil")
	nova := make([]ast.Stmt, 0, len(origList)+1)
	nova = append(nova, origList[:idx]...)
	nova = append(nova, descarte, ret)
	nova = append(nova, origList[idx+1:]...)
	blk.List = nova
	codigo, err := formatarAST(fset, arq)
	blk.List = origList
	ret.Results[posErr] = origExpr
	return codigo, "return nil", err
}

func indiceIdentErr(results []ast.Expr) int {
	for i := len(results) - 1; i >= 0; i-- {
		if id, ok := results[i].(*ast.Ident); ok && id.Name == "err" {
			return i
		}
	}
	return -1
}

func filtrarSomaSemString(arq *ast.File, bin *ast.BinaryExpr) bool {
	if bin.Op != token.ADD {
		return true
	}
	varsString := coletarVarsString(arq)
	return !pareceString(bin.X, varsString) && !pareceString(bin.Y, varsString)
}

func coletarVarsString(arq *ast.File) map[string]bool {
	vars := make(map[string]bool)
	ast.Inspect(arq, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Field:
			if ehTipoString(x.Type) {
				for _, id := range x.Names {
					vars[id.Name] = true
				}
			}
		case *ast.ValueSpec:
			if ehTipoString(x.Type) {
				for _, id := range x.Names {
					vars[id.Name] = true
				}
			} else {
				for i, id := range x.Names {
					if i < len(x.Values) && pareceString(x.Values[i], vars) {
						vars[id.Name] = true
					}
				}
			}
		case *ast.AssignStmt:
			if x.Tok == token.DEFINE || x.Tok == token.ASSIGN {
				for i, lhs := range x.Lhs {
					id, ok := lhs.(*ast.Ident)
					if ok && i < len(x.Rhs) && pareceString(x.Rhs[i], vars) {
						vars[id.Name] = true
					}
				}
			}
		}
		return true
	})
	return vars
}

func ehTipoString(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == "string"
}

func pareceString(expr ast.Expr, vars map[string]bool) bool {
	switch x := expr.(type) {
	case *ast.BasicLit:
		return x.Kind == token.STRING
	case *ast.Ident:
		return vars[x.Name]
	case *ast.ParenExpr:
		return pareceString(x.X, vars)
	case *ast.BinaryExpr:
		return x.Op == token.ADD && (pareceString(x.X, vars) || pareceString(x.Y, vars))
	case *ast.CallExpr:
		return chamadaRetornaString(x)
	}
	return false
}

func chamadaRetornaString(call *ast.CallExpr) bool {
	if id, ok := call.Fun.(*ast.Ident); ok {
		return id.Name == "string"
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	switch sel.Sel.Name {
	case "String", "Sprintf", "Sprint", "Join", "TrimSpace", "Itoa", "FormatInt", "ToLower", "ToUpper", "Base", "Dir", "Clean":
		return true
	}
	return false
}
