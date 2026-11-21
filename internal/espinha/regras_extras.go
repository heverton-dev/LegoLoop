package espinha

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func verificarProcessosCompletos(raiz string) []Violacao {
	dirProc := filepath.Join(raiz, "processos")
	entradas, err := os.ReadDir(dirProc)
	if err != nil {
		return nil
	}
	var vs []Violacao
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(dirProc, e.Name())
		if !temArquivoDAG(dir) {
			vs = append(vs, Violacao{
				Regra:  "processo-incompleto",
				Alvo:   e.Name(),
				Motivo: "arquivo essencial ausente: *.dag.json",
			})
		}
		if _, err := os.Stat(filepath.Join(dir, "ANTICORPOS.md")); err != nil {
			vs = append(vs, Violacao{
				Regra:  "processo-incompleto",
				Alvo:   e.Name(),
				Motivo: "arquivo essencial ausente: ANTICORPOS.md",
			})
		}
	}
	return vs
}

func temArquivoDAG(dir string) bool {
	itens, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, it := range itens {
		if !it.IsDir() && strings.HasSuffix(it.Name(), ".dag.json") {
			return true
		}
	}
	return false
}

func verificarErrosIgnorados(raiz string, limites map[string]int) []Violacao {
	dirInternal := filepath.Join(raiz, "internal")
	if _, err := os.Stat(dirInternal); err != nil {
		return nil
	}
	contagem, err := contarErrosIgnoradosInternal(raiz, dirInternal)
	if err != nil {
		return []Violacao{{Regra: "erro-ignorado-acima-do-limite", Alvo: "internal", Motivo: err.Error()}}
	}
	var vs []Violacao
	for pacote, qtd := range contagem {
		lim := limites[pacote]
		if qtd > lim {
			vs = append(vs, Violacao{
				Regra:  "erro-ignorado-acima-do-limite",
				Alvo:   pacote,
				Motivo: fmt.Sprintf("%d erros ignorados excedem limite %d", qtd, lim),
			})
		}
	}
	return vs
}

func contarErrosIgnoradosInternal(raiz, dirInternal string) (map[string]int, error) {
	contagem := make(map[string]int)
	fset := token.NewFileSet()
	err := filepath.WalkDir(dirInternal, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil || d.IsDir() || !ehArquivoGoProducao(d.Name()) {
			return e
		}
		arq, err := parser.ParseFile(fset, caminho, nil, 0)
		if err != nil {
			return fmt.Errorf("analisar %s: %w", caminhoRelativo(raiz, caminho), err)
		}
		pacote := caminhoRelativo(raiz, filepath.Dir(caminho))
		if n := contarBlankAssignArquivo(arq); n > 0 {
			contagem[pacote] += n
		}
		return nil
	})
	return contagem, err
}

func contarBlankAssignArquivo(arq *ast.File) int {
	total := 0
	ast.Inspect(arq, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || as.Tok != token.ASSIGN || len(as.Lhs) == 0 {
			return true
		}
		if todosIdentificadoresBlank(as.Lhs) {
			total++
		}
		return true
	})
	return total
}

func todosIdentificadoresBlank(lhs []ast.Expr) bool {
	for _, expr := range lhs {
		id, ok := expr.(*ast.Ident)
		if !ok || id.Name != "_" {
			return false
		}
	}
	return true
}

func verificarGlobaisMutaveis(raiz string, justificadas []string) []Violacao {
	dirInternal := filepath.Join(raiz, "internal")
	if _, err := os.Stat(dirInternal); err != nil {
		return nil
	}
	permitidas := make(map[string]bool, len(justificadas))
	for _, j := range justificadas {
		permitidas[strings.TrimSpace(j)] = true
	}
	vs, err := coletarGlobaisMutaveis(raiz, dirInternal, permitidas)
	if err != nil {
		return []Violacao{{Regra: "global-mutavel-nao-justificada", Alvo: "internal", Motivo: err.Error()}}
	}
	return vs
}

func coletarGlobaisMutaveis(raiz, dirInternal string, permitidas map[string]bool) ([]Violacao, error) {
	fset := token.NewFileSet()
	var vs []Violacao
	err := filepath.WalkDir(dirInternal, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil || d.IsDir() || !ehArquivoGoProducao(d.Name()) {
			return e
		}
		arq, err := parser.ParseFile(fset, caminho, nil, 0)
		if err != nil {
			return fmt.Errorf("analisar %s: %w", caminhoRelativo(raiz, caminho), err)
		}
		pacote := caminhoRelativo(raiz, filepath.Dir(caminho))
		vs = append(vs, inspecionarGlobaisArquivo(pacote, arq, permitidas)...)
		return nil
	})
	return vs, err
}

func inspecionarGlobaisArquivo(pacote string, arq *ast.File, permitidas map[string]bool) []Violacao {
	var vs []Violacao
	for _, decl := range arq.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, sp := range gd.Specs {
			vspec, ok := sp.(*ast.ValueSpec)
			if !ok {
				continue
			}
			vs = append(vs, checarValueSpecGlobal(pacote, vspec, permitidas)...)
		}
	}
	return vs
}

func checarValueSpecGlobal(pacote string, vspec *ast.ValueSpec, permitidas map[string]bool) []Violacao {
	var vs []Violacao
	tipoMutavel := ehTipoMutavelAST(vspec.Type)
	for i, nome := range vspec.Names {
		if nome.Name == "_" {
			continue
		}
		mut := tipoMutavel
		if !mut && i < len(vspec.Values) {
			mut = ehValorMutavelAST(vspec.Values[i])
		}
		alvo := pacote + "." + nome.Name
		if mut && !permitidas[alvo] {
			vs = append(vs, Violacao{
				Regra:  "global-mutavel-nao-justificada",
				Alvo:   alvo,
				Motivo: "variavel global mutavel ausente em globais_justificadas de espinha/registro.json",
			})
		}
	}
	return vs
}

func ehTipoMutavelAST(expr ast.Expr) bool {
	if expr == nil {
		return false
	}
	switch t := expr.(type) {
	case *ast.MapType, *ast.StarExpr:
		return true
	case *ast.ArrayType:
		return t.Len == nil
	default:
		return false
	}
}

func ehValorMutavelAST(expr ast.Expr) bool {
	switch v := expr.(type) {
	case *ast.CompositeLit:
		return ehTipoMutavelAST(v.Type)
	case *ast.UnaryExpr:
		return v.Op == token.AND
	case *ast.CallExpr:
		return ehChamadaAlocacaoMutavel(v)
	default:
		return false
	}
}

func ehChamadaAlocacaoMutavel(call *ast.CallExpr) bool {
	id, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if id.Name == "new" {
		return true
	}
	if id.Name == "make" && len(call.Args) > 0 {
		return ehTipoMutavelAST(call.Args[0])
	}
	return false
}

func verificarInits(raiz string, justificadas map[string]string) []Violacao {
	dirInternal := filepath.Join(raiz, "internal")
	if _, err := os.Stat(dirInternal); err != nil {
		return nil
	}
	vs, err := coletarInitsNaoJustificados(raiz, dirInternal, justificadas)
	if err != nil {
		return []Violacao{{Regra: "init-nao-justificado", Alvo: "internal", Motivo: err.Error()}}
	}
	return vs
}

func coletarInitsNaoJustificados(raiz, dirInternal string, justificadas map[string]string) ([]Violacao, error) {
	fset := token.NewFileSet()
	var vs []Violacao
	err := filepath.WalkDir(dirInternal, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil || d.IsDir() || !ehArquivoGoProducao(d.Name()) {
			return e
		}
		arq, err := parser.ParseFile(fset, caminho, nil, 0)
		if err != nil {
			return fmt.Errorf("analisar %s: %w", caminhoRelativo(raiz, caminho), err)
		}
		if !arquivoTemInit(arq) {
			return nil
		}
		alvo := caminhoRelativo(raiz, caminho)
		if strings.TrimSpace(justificadas[alvo]) == "" {
			vs = append(vs, Violacao{
				Regra:  "init-nao-justificado",
				Alvo:   alvo,
				Motivo: "funcao init() ausente ou sem justificativa em inits_justificados de espinha/registro.json",
			})
		}
		return nil
	})
	return vs, err
}

func arquivoTemInit(arq *ast.File) bool {
	for _, decl := range arq.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && fn.Name.Name == "init" {
			return true
		}
	}
	return false
}

type dagComandosBruto struct {
	Nodes []noComandoBruto `json:"nodes"`
}

type noComandoBruto struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Command []string `json:"command"`
}

func verificarComandosExclusivosSO(raiz string) []Violacao {
	arquivos, err := coletarDAGsProcessos(filepath.Join(raiz, "processos"))
	if err != nil {
		return nil
	}
	var vs []Violacao
	for _, arq := range arquivos {
		rel := caminhoRelativo(raiz, arq)
		dados, err := os.ReadFile(arq)
		if err != nil {
			continue
		}
		var bruto dagComandosBruto
		if json.Unmarshal(dados, &bruto) != nil {
			continue
		}
		vs = append(vs, checarComandosNos(rel, bruto.Nodes)...)
	}
	return vs
}

func checarComandosNos(rel string, nos []noComandoBruto) []Violacao {
	var vs []Violacao
	for _, n := range nos {
		if n.Type != "script" || len(n.Command) == 0 {
			continue
		}
		if ehComandoExclusivoSO(n.Command[0]) {
			vs = append(vs, Violacao{
				Regra:  "comando-exclusivo-de-so",
				Alvo:   rel + ":" + n.ID,
				Motivo: fmt.Sprintf("comando exclusivo de SO %q em no script", n.Command[0]),
			})
		}
	}
	return vs
}

func ehComandoExclusivoSO(cmd string) bool {
	base := strings.ToLower(filepath.Base(strings.TrimSpace(cmd)))
	base = strings.TrimSuffix(base, ".exe")
	switch base {
	case "powershell", "pwsh", "cmd", "python", "python3":
		return true
	default:
		return false
	}
}
