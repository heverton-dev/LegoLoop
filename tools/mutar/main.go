// Comando mutar executa testes de mutação em pacotes Go usando apenas biblioteca padrão.
//
// Uso: go run ./tools/mutar <pacote>... [--orcamento DUR] [--justificados ARQ]
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const caminhoJustificadosPadrao = "espinha/mutantes_justificados.json"

// Config guarda os parâmetros da execução via CLI.
type Config struct {
	Pacotes               []string
	Orcamento             time.Duration
	ArqJustificados       string
	JustificadosExplicito bool
}

// Operador define uma transformação AST geradora de mutantes.
type Operador struct {
	Nome  string
	Gerar func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error)
}

// Mutante representa uma mutação individual sobre um arquivo Go.
type Mutante struct {
	ID       string
	Pacote   string
	Arquivo  string
	Linha    int
	Operador string
	Original string
	Mutado   string
	Codigo   []byte
}

// Justificado registra um mutante sobrevivente justificado por anticorpo.
type Justificado struct {
	ID            string `json:"id"`
	Pacote        string `json:"pacote"`
	Arquivo       string `json:"arquivo"`
	Linha         int    `json:"linha"`
	Operador      string `json:"operador"`
	Anticorpo     string `json:"anticorpo"`
	Justificativa string `json:"justificativa"`
}

// StatusMutante classifica o desfecho de um mutante após os testes.
type StatusMutante string

const (
	StatusMorto        StatusMutante = "morto"
	StatusSobrevivente StatusMutante = "sobrevivente"
	StatusJustificado  StatusMutante = "justificado"
	StatusTimeout      StatusMutante = "timeout"
)

// Resultado consolida o desfecho da avaliação de um mutante.
type Resultado struct {
	Mutante       Mutante
	Status        StatusMutante
	Saida         string
	Justificativa *Justificado
}

// Relatorio sumariza a rodada de mutação.
type Relatorio struct {
	Total         int
	Mortos        int
	Sobreviventes int
	Justificados  int
	Timeouts      int
	Resultados    []Resultado
}

func main() {
	os.Exit(executarCLI(os.Args[1:], os.Stdout, os.Stderr))
}

func parseArgs(args []string, stderr io.Writer) (Config, error) {
	flags, pacotes := separarFlagsEPosicionais(args)
	fs := flag.NewFlagSet("mutar", flag.ContinueOnError)
	fs.SetOutput(stderr)
	orcStr := fs.String("orcamento", "", "orçamento máximo de tempo (ex: 5m, 30s)")
	arqJust := fs.String("justificados", filepath.FromSlash(caminhoJustificadosPadrao), "arquivo JSON de mutantes justificados")
	if err := fs.Parse(flags); err != nil {
		return Config{}, err
	}
	pacotes = append(pacotes, fs.Args()...)
	if len(pacotes) == 0 {
		return Config{}, errors.New("informe ao menos um pacote")
	}
	orc, err := parseOrcamento(*orcStr)
	if err != nil {
		return Config{}, err
	}
	explicito := flagInformada(fs, "justificados")
	return Config{
		Pacotes:               pacotes,
		Orcamento:             orc,
		ArqJustificados:       *arqJust,
		JustificadosExplicito: explicito,
	}, nil
}

func separarFlagsEPosicionais(args []string) (flags, posicionais []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			posicionais = append(posicionais, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") && arg != "-" {
			flags = append(flags, arg)
			if !strings.Contains(arg, "=") && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
			continue
		}
		posicionais = append(posicionais, arg)
	}
	return flags, posicionais
}

func parseOrcamento(bruto string) (time.Duration, error) {
	if strings.TrimSpace(bruto) == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(bruto)
	if err != nil {
		return 0, fmt.Errorf("orçamento inválido %q: %w", bruto, err)
	}
	if d < 0 {
		return 0, fmt.Errorf("orçamento não pode ser negativo: %v", d)
	}
	return d, nil
}

func flagInformada(fs *flag.FlagSet, nome string) bool {
	achou := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == nome {
			achou = true
		}
	})
	return achou
}

func resolverJustificados(dir string, cfg Config) ([]Justificado, error) {
	caminho := cfg.ArqJustificados
	if !filepath.IsAbs(caminho) {
		caminho = filepath.Join(dir, caminho)
	}
	if _, err := os.Stat(caminho); errors.Is(err, os.ErrNotExist) && !cfg.JustificadosExplicito {
		return nil, nil
	}
	return carregarJustificados(caminho)
}

func carregarJustificados(caminho string) ([]Justificado, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	var lista []Justificado
	if err := json.Unmarshal(dados, &lista); err != nil {
		return nil, fmt.Errorf("JSON inválido em %s: %w", caminho, err)
	}
	for i, j := range lista {
		if err := validarJustificado(j); err != nil {
			return nil, fmt.Errorf("entrada %d em %s: %w", i, caminho, err)
		}
	}
	return lista, nil
}

func validarJustificado(j Justificado) error {
	if strings.TrimSpace(j.ID) == "" {
		return errors.New("campo id obrigatório")
	}
	if strings.TrimSpace(j.Pacote) == "" || strings.TrimSpace(j.Arquivo) == "" {
		return errors.New("campos pacote e arquivo obrigatórios")
	}
	if j.Linha <= 0 || strings.TrimSpace(j.Operador) == "" {
		return errors.New("campos linha (>0) e operador obrigatórios")
	}
	if strings.TrimSpace(j.Anticorpo) == "" || strings.TrimSpace(j.Justificativa) == "" {
		return errors.New("campos anticorpo e justificativa obrigatórios")
	}
	return nil
}

func coletarMutantes(dir string, pacotes []string, operadores []Operador) ([]Mutante, error) {
	var todos []Mutante
	for _, pkg := range pacotes {
		caminhoPkg := pkg
		if !filepath.IsAbs(caminhoPkg) {
			caminhoPkg = filepath.Join(dir, pkg)
		}
		muts, err := Mutantes(caminhoPkg, operadores)
		if err != nil {
			return nil, err
		}
		todos = append(todos, muts...)
	}
	return todos, nil
}

// Mutantes gera os mutantes de um diretório de pacote sem executar testes.
func Mutantes(pacote string, operadores []Operador) ([]Mutante, error) {
	arquivos, err := listarArquivosGo(pacote)
	if err != nil {
		return nil, err
	}
	var resultado []Mutante
	for _, arqPath := range arquivos {
		fset := token.NewFileSet()
		arq, err := parser.ParseFile(fset, arqPath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", arqPath, err)
		}
		for _, op := range operadores {
			if op.Gerar == nil {
				continue
			}
			gerados, err := op.Gerar(fset, arq, pacote, arqPath)
			if err != nil {
				return nil, fmt.Errorf("operador %s em %s: %w", op.Nome, arqPath, err)
			}
			resultado = append(resultado, gerados...)
		}
	}
	return resultado, nil
}

func listarArquivosGo(dirPacote string) ([]string, error) {
	entradas, err := os.ReadDir(dirPacote)
	if err != nil {
		return nil, err
	}
	var arquivos []string
	for _, e := range entradas {
		nome := e.Name()
		if e.IsDir() || !strings.HasSuffix(nome, ".go") || strings.HasSuffix(nome, "_test.go") {
			continue
		}
		arquivos = append(arquivos, filepath.Join(dirPacote, nome))
	}
	sort.Strings(arquivos)
	return arquivos, nil
}

func formatarAST(fset *token.FileSet, arq *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, arq); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
