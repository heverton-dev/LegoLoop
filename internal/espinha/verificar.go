package espinha

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"legoloop/internal/dag"
)

// Violacao representa um desvio estrutural detectado na verificacao da espinha dorsal.
type Violacao struct {
	Regra  string `json:"regra"`
	Alvo   string `json:"alvo"`
	Motivo string `json:"motivo"`
}

// String formata a violacao no contrato canonico de saida linha a linha.
func (v Violacao) String() string {
	motivo := strings.ReplaceAll(strings.ReplaceAll(v.Motivo, "\r\n", "; "), "\n", "; ")
	return fmt.Sprintf("regra: %s; alvo: %s; motivo: %s", v.Regra, v.Alvo, motivo)
}

// Verificar executa todas as verificacoes da espinha dorsal sobre o diretorio raiz
// e devolve as violacoes em ordem deterministica.
func Verificar(raiz string) []Violacao {
	caminhoReg := filepath.Join(raiz, "espinha", "registro.json")
	dados, err := os.ReadFile(caminhoReg)
	if err != nil {
		return []Violacao{{
			Regra:  "registro-invalido",
			Alvo:   "espinha/registro.json",
			Motivo: err.Error(),
		}}
	}
	var reg Registro
	if err := json.Unmarshal(dados, &reg); err != nil {
		return []Violacao{{
			Regra:  "registro-invalido",
			Alvo:   "espinha/registro.json",
			Motivo: err.Error(),
		}}
	}
	var vs []Violacao
	if err := reg.Validar(); err != nil {
		vs = append(vs, Violacao{Regra: "registro-invalido", Alvo: "espinha/registro.json", Motivo: err.Error()})
	}
	vs = append(vs, verificarVerbos(raiz, reg)...)
	vs = append(vs, verificarLeis(raiz, reg)...)
	vs = append(vs, verificarProcessos(raiz, reg)...)
	vs = append(vs, verificarProcessosCompletos(raiz)...)
	vs = append(vs, verificarDAGs(raiz)...)
	vs = append(vs, verificarComandosExclusivosSO(raiz)...)
	vs = append(vs, verificarMotores(raiz, dados, reg)...)
	vs = append(vs, verificarFuncoesDuplicadas(raiz)...)
	vs = append(vs, verificarErrosIgnorados(raiz, reg.LimitesErroIgnorado)...)
	vs = append(vs, verificarGlobaisMutaveis(raiz, reg.GlobaisJustificadas)...)
	vs = append(vs, verificarInits(raiz, reg.InitsJustificados)...)
	vs = append(vs, verificarDebitoMarcado(raiz)...)
	ordenarViolacoes(vs)
	return vs
}

func ordenarViolacoes(vs []Violacao) {
	slices.SortFunc(vs, func(a, b Violacao) int {
		if c := cmp.Compare(a.Regra, b.Regra); c != 0 {
			return c
		}
		if c := cmp.Compare(a.Alvo, b.Alvo); c != 0 {
			return c
		}
		return cmp.Compare(a.Motivo, b.Motivo)
	})
}

func verificarVerbos(raiz string, reg Registro) []Violacao {
	caminho := filepath.Join(raiz, "internal", "cli", "cli.go")
	daCLI, err := extrairVerbosCLI(caminho)
	if err != nil {
		return []Violacao{{Regra: "verbo-sem-registro", Alvo: "internal/cli/cli.go", Motivo: err.Error()}}
	}
	setCLI := make(map[string]bool, len(daCLI))
	for _, v := range daCLI {
		setCLI[v] = true
	}
	setReg := make(map[string]bool, len(reg.Verbos))
	for _, v := range reg.Verbos {
		setReg[strings.TrimSpace(v.Nome)] = true
	}
	var vs []Violacao
	for _, nome := range daCLI {
		if !setReg[nome] {
			vs = append(vs, Violacao{
				Regra:  "verbo-sem-registro",
				Alvo:   nome,
				Motivo: "verbo declarado em internal/cli/cli.go ausente em espinha/registro.json",
			})
		}
	}
	for nome := range setReg {
		if nome != "" && !setCLI[nome] {
			vs = append(vs, Violacao{
				Regra:  "verbo-orfao",
				Alvo:   nome,
				Motivo: "verbo listado em espinha/registro.json ausente em internal/cli/cli.go",
			})
		}
	}
	return vs
}

func extrairVerbosCLI(caminho string) ([]string, error) {
	fset := token.NewFileSet()
	arq, err := parser.ParseFile(fset, caminho, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("analisar %s: %w", caminho, err)
	}
	var verbos []string
	for _, decl := range arq.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "mapaComandos" || fn.Body == nil {
			continue
		}
		verbos = extrairChavesMapa(fn.Body)
		break
	}
	slices.Sort(verbos)
	return verbos, nil
}

func extrairChavesMapa(body *ast.BlockStmt) []string {
	var chaves []string
	ast.Inspect(body, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if _, isMap := lit.Type.(*ast.MapType); !isMap {
			return true
		}
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			bl, ok := kv.Key.(*ast.BasicLit)
			if !ok || bl.Kind != token.STRING {
				continue
			}
			if nome, err := strconv.Unquote(bl.Value); err == nil {
				chaves = append(chaves, nome)
			}
		}
		return false
	})
	return chaves
}

func verificarLeis(raiz string, reg Registro) []Violacao {
	caminho := filepath.Join(raiz, "AGENTS.md")
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return []Violacao{{Regra: "lei-sem-verificador", Alvo: "AGENTS.md", Motivo: err.Error()}}
	}
	leis := ExtrairLeisAgents(string(dados))
	if err := ConferirLeis(leis, reg); err == nil {
		return nil
	}
	if len(leis) == 0 {
		return []Violacao{{Regra: "lei-sem-verificador", Alvo: "AGENTS.md", Motivo: "nenhuma lei encontrada em AGENTS.md"}}
	}
	indice, err := indexarRequisitos(reg.Requisitos)
	if err != nil {
		return []Violacao{{Regra: "lei-sem-verificador", Alvo: "espinha/registro.json", Motivo: err.Error()}}
	}
	var vs []Violacao
	for _, num := range leis {
		chave := fmt.Sprintf("Lei %d", num)
		req, ok := indice[chave]
		if !ok {
			vs = append(vs, Violacao{Regra: "lei-sem-verificador", Alvo: chave, Motivo: "ausente em espinha/registro.json"})
			continue
		}
		if err := validarRequisitoLei(chave, req); err != nil {
			vs = append(vs, Violacao{Regra: "lei-sem-verificador", Alvo: chave, Motivo: err.Error()})
		}
	}
	return vs
}

func verificarProcessos(raiz string, reg Registro) []Violacao {
	dirProc := filepath.Join(raiz, "processos")
	semPrefixo, err := VerificarPrefixoProcessos(dirProc)
	if err != nil {
		return []Violacao{{Regra: "processo-sem-prefixo", Alvo: "processos", Motivo: err.Error()}}
	}
	var vs []Violacao
	for _, nome := range semPrefixo {
		vs = append(vs, Violacao{
			Regra:  "processo-sem-prefixo",
			Alvo:   nome,
			Motivo: "diretorio em processos/ sem prefixo legoloop-",
		})
	}
	emDisco, err := listarDirsProcessos(dirProc)
	if err != nil {
		return append(vs, Violacao{Regra: "processo-duplicado", Alvo: "processos", Motivo: err.Error()})
	}
	return append(vs, compararProcessosRegistro(emDisco, reg.Processos)...)
}

func listarDirsProcessos(dirProc string) (map[string]bool, error) {
	entradas, err := os.ReadDir(dirProc)
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool)
	for _, e := range entradas {
		if e.IsDir() {
			set[e.Name()] = true
		}
	}
	return set, nil
}

func compararProcessosRegistro(emDisco map[string]bool, processos []Processo) []Violacao {
	noRegistro := make(map[string]bool, len(processos))
	var vs []Violacao
	for _, p := range processos {
		slug := strings.TrimSpace(p.Slug)
		if noRegistro[slug] {
			vs = append(vs, Violacao{Regra: "processo-duplicado", Alvo: slug, Motivo: "slug repetido em espinha/registro.json"})
			continue
		}
		noRegistro[slug] = true
		if slug != "" && !emDisco[slug] {
			vs = append(vs, Violacao{
				Regra:  "processo-duplicado",
				Alvo:   slug,
				Motivo: "processo listado em espinha/registro.json ausente em processos/",
			})
		}
	}
	for slug := range emDisco {
		if !noRegistro[slug] {
			vs = append(vs, Violacao{
				Regra:  "processo-duplicado",
				Alvo:   slug,
				Motivo: "processo em processos/ ausente em espinha/registro.json",
			})
		}
	}
	return vs
}

type dagBruto struct {
	Name  string    `json:"name"`
	Nodes []noBruto `json:"nodes"`
}

type noBruto struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func verificarDAGs(raiz string) []Violacao {
	arquivos, err := coletarDAGsProcessos(filepath.Join(raiz, "processos"))
	if err != nil {
		return []Violacao{{Regra: "executor-incoerente", Alvo: "processos", Motivo: err.Error()}}
	}
	nomesVistos := make(map[string]string, len(arquivos))
	var vs []Violacao
	for _, arq := range arquivos {
		rel := caminhoRelativo(raiz, arq)
		vs = append(vs, verificarUmDAG(arq, rel, nomesVistos)...)
	}
	return vs
}

func coletarDAGsProcessos(dirProc string) ([]string, error) {
	var lista []string
	err := filepath.WalkDir(dirProc, func(caminho string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".dag.json") {
			lista = append(lista, caminho)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(lista)
	return lista, nil
}

func verificarUmDAG(caminho, rel string, nomesVistos map[string]string) []Violacao {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return []Violacao{{Regra: "executor-incoerente", Alvo: rel, Motivo: err.Error()}}
	}
	var vs []Violacao
	var bruto dagBruto
	if json.Unmarshal(dados, &bruto) == nil {
		vs = append(vs, checarNomeENos(rel, bruto, nomesVistos)...)
	}
	if temRegra(vs, "tipo-de-no-fora-da-tabela") {
		return vs
	}
	d, err := dag.Parse(dados)
	if err != nil {
		if errors.Is(err, dag.ErrExecutorIncoerente) {
			return append(vs, Violacao{Regra: "executor-incoerente", Alvo: rel, Motivo: err.Error()})
		}
		return append(vs, Violacao{Regra: "dag-invalido", Alvo: rel, Motivo: err.Error()})
	}
	if err := d.ValidarExecutores(); err != nil {
		vs = append(vs, Violacao{Regra: "executor-incoerente", Alvo: rel, Motivo: err.Error()})
	}
	return vs
}

func checarNomeENos(rel string, bruto dagBruto, nomesVistos map[string]string) []Violacao {
	var vs []Violacao
	nome := strings.TrimSpace(bruto.Name)
	if nome != "" {
		if anterior, dup := nomesVistos[nome]; dup {
			vs = append(vs, Violacao{
				Regra:  "processo-duplicado",
				Alvo:   rel,
				Motivo: fmt.Sprintf("nome %q ja usado em %s", nome, anterior),
			})
		} else {
			nomesVistos[nome] = rel
		}
	}
	for _, n := range bruto.Nodes {
		if _, err := dag.DoDAG(n.Type); err != nil {
			vs = append(vs, Violacao{
				Regra:  "tipo-de-no-fora-da-tabela",
				Alvo:   rel + ":" + n.ID,
				Motivo: err.Error(),
			})
		}
	}
	return vs
}

func temRegra(vs []Violacao, regra string) bool {
	for _, v := range vs {
		if v.Regra == regra {
			return true
		}
	}
	return false
}

func caminhoRelativo(raiz, caminho string) string {
	rel, err := filepath.Rel(raiz, caminho)
	if err != nil {
		return filepath.ToSlash(caminho)
	}
	return filepath.ToSlash(rel)
}

func verificarMotores(raiz string, dadosJSON []byte, reg Registro) []Violacao {
	var bruto struct {
		Motores []map[string]any `json:"motores"`
	}
	var vs []Violacao
	if json.Unmarshal(dadosJSON, &bruto) == nil {
		for _, m := range bruto.Motores {
			pacote, _ := m["pacote"].(string)
			val, existe := m["deterministico"]
			if _, ok := val.(bool); !existe || !ok {
				vs = append(vs, Violacao{
					Regra:  "motor-sem-deterministico",
					Alvo:   pacote,
					Motivo: "chave booleana deterministico ausente em espinha/registro.json",
				})
			}
		}
	}
	for _, m := range reg.Motores {
		if !motorExisteEmInternal(raiz, m.Pacote) {
			vs = append(vs, Violacao{
				Regra:  "motor-inexistente",
				Alvo:   m.Pacote,
				Motivo: "pacote de motor ausente sob internal/",
			})
		}
	}
	return vs
}

func motorExisteEmInternal(raiz, pacote string) bool {
	if !strings.HasPrefix(pacote, "internal/") {
		return false
	}
	dir := filepath.Join(raiz, filepath.FromSlash(pacote))
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entradas {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			return true
		}
	}
	return false
}
