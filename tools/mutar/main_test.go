package main

import (
	"bytes"
	"context"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type executorFalso struct {
	passou   bool
	saida    string
	err      error
	chamadas int
}

func (e *executorFalso) Rodar(ctx context.Context, dirPacote string) (bool, string, error) {
	e.chamadas++
	if ctx.Err() != nil {
		return false, "", ctx.Err()
	}
	return e.passou, e.saida, e.err
}

func TestParseArgs_SuportaFlagsAntesEDepoisDosPacotes(t *testing.T) {
	var errBuf bytes.Buffer
	cfg, err := parseArgs([]string{
		"./internal/dag",
		"./internal/runner",
		"--orcamento", "2m30s",
		"--justificados", "custom.json",
	}, &errBuf)
	if err != nil {
		t.Fatalf("parseArgs falhou: %v", err)
	}
	if len(cfg.Pacotes) != 2 || cfg.Pacotes[0] != "./internal/dag" || cfg.Pacotes[1] != "./internal/runner" {
		t.Fatalf("pacotes inesperados: %#v", cfg.Pacotes)
	}
	if cfg.Orcamento != 150*time.Second {
		t.Fatalf("orcamento esperado 2m30s, obtido %v", cfg.Orcamento)
	}
	if cfg.ArqJustificados != "custom.json" || !cfg.JustificadosExplicito {
		t.Fatalf("justificados inesperado: %q explicito=%v", cfg.ArqJustificados, cfg.JustificadosExplicito)
	}
}

func TestParseArgs_PadroesEErros(t *testing.T) {
	var errBuf bytes.Buffer
	cfg, err := parseArgs([]string{"./pkg/a"}, &errBuf)
	if err != nil {
		t.Fatalf("parseArgs padrao falhou: %v", err)
	}
	esperado := filepath.Join("espinha", "mutantes_justificados.json")
	if cfg.ArqJustificados != esperado || cfg.JustificadosExplicito {
		t.Fatalf("padrao esperado %q explicito=false, obtido %q explicito=%v", esperado, cfg.ArqJustificados, cfg.JustificadosExplicito)
	}

	if _, err := parseArgs([]string{}, &errBuf); err == nil {
		t.Fatal("esperava erro quando nenhum pacote e informado")
	}
	if _, err := parseArgs([]string{"./pkg/a", "--orcamento", "-1s"}, &errBuf); err == nil {
		t.Fatal("esperava erro para orcamento negativo")
	}
	if _, err := parseArgs([]string{"./pkg/a", "--orcamento", "invalido"}, &errBuf); err == nil {
		t.Fatal("esperava erro para orcamento invalido")
	}
}

func TestCarregarJustificados_ValidoEInvalido(t *testing.T) {
	dir := t.TempDir()
	valido := filepath.Join(dir, "justificados.json")
	conteudo := `[
  {
    "id": "M001",
    "pacote": "./internal/dag",
    "arquivo": "parse.go",
    "linha": 42,
    "operador": "eq_neq",
    "anticorpo": "ANT-001",
    "justificativa": "mutante equivalente no caminho morto"
  }
]`
	if err := os.WriteFile(valido, []byte(conteudo), 0o600); err != nil {
		t.Fatal(err)
	}
	lista, err := carregarJustificados(valido)
	if err != nil {
		t.Fatalf("carregarJustificados falhou: %v", err)
	}
	if len(lista) != 1 || lista[0].ID != "M001" || lista[0].Anticorpo != "ANT-001" {
		t.Fatalf("lista inesperada: %#v", lista)
	}

	invalidoJSON := filepath.Join(dir, "ruim.json")
	if err := os.WriteFile(invalidoJSON, []byte("{nao-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := carregarJustificados(invalidoJSON); err == nil {
		t.Fatal("esperava erro para JSON malformado")
	}

	semAnticorpo := filepath.Join(dir, "incompleto.json")
	if err := os.WriteFile(semAnticorpo, []byte(`[{"id":"M1","pacote":"p","arquivo":"a.go","linha":1,"operador":"op","justificativa":"ok"}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := carregarJustificados(semAnticorpo); err == nil {
		t.Fatal("esperava erro quando anticorpo esta vazio")
	}
}

func TestMutantes_ParseEAplicacaoOperador(t *testing.T) {
	dir := t.TempDir()
	fonte := "package exemplo\n\nfunc Soma(a, b int) int {\n\treturn a + b\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "soma.go"), []byte(fonte), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "soma_test.go"), []byte("package exemplo\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	op := Operador{
		Nome: "teste_noop",
		Gerar: func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error) {
			codigo, err := formatarAST(fset, arq)
			if err != nil {
				return nil, err
			}
			return []Mutante{{
				ID:       "M-1",
				Pacote:   pacote,
				Arquivo:  filepath.Base(caminho),
				Linha:    4,
				Operador: "teste_noop",
				Codigo:   codigo,
			}}, nil
		},
	}

	muts, err := Mutantes(dir, []Operador{op})
	if err != nil {
		t.Fatalf("Mutantes falhou: %v", err)
	}
	if len(muts) != 1 {
		t.Fatalf("esperava 1 mutante (ignorando _test.go), obtido %d", len(muts))
	}
	if !strings.Contains(string(muts[0].Codigo), "func Soma") {
		t.Fatalf("codigo formatado inesperado: %s", string(muts[0].Codigo))
	}
}

func TestRodar_ClassificaMortoSobreviventeEJustificado(t *testing.T) {
	m1 := Mutante{ID: "M1", Pacote: "./p", Arquivo: "a.go", Linha: 10, Operador: "eq"}
	m2 := Mutante{ID: "M2", Pacote: "./p", Arquivo: "a.go", Linha: 20, Operador: "add"}
	just := []Justificado{{
		ID:            "M1",
		Pacote:        "./p",
		Arquivo:       "a.go",
		Linha:         10,
		Operador:      "eq",
		Anticorpo:     "ANT-1",
		Justificativa: "equivalente",
	}}

	execSobrevive := &executorFalso{passou: true}
	rel := Rodar(context.Background(), []Mutante{m1, m2}, just, execSobrevive, 0)
	if rel.Total != 2 || rel.Justificados != 1 || rel.Sobreviventes != 1 || rel.Mortos != 0 {
		t.Fatalf("relatorio inesperado: %#v", rel)
	}
	if execSobrevive.chamadas != 1 {
		t.Fatalf("mutante justificado nao deveria executar testes; chamadas=%d", execSobrevive.chamadas)
	}

	execMata := &executorFalso{passou: false}
	rel2 := Rodar(context.Background(), []Mutante{m2}, nil, execMata, 5*time.Second)
	if rel2.Mortos != 1 || rel2.Sobreviventes != 0 {
		t.Fatalf("esperava 1 morto, obtido %#v", rel2)
	}
}

func TestExecutar_CodigosDeSaida(t *testing.T) {
	dir := t.TempDir()
	pkgDir := filepath.Join(dir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "a.go"), []byte("package pkg\nfunc F() int { return 1 }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	if code := executar([]string{}, &out, &errBuf, dir, &executorFalso{}, nil); code != 2 {
		t.Fatalf("sem pacotes esperava exit 2, obtido %d", code)
	}
	if code := executar([]string{"./pkg", "--justificados", "nao_existe.json"}, &out, &errBuf, dir, &executorFalso{}, nil); code != 2 {
		t.Fatalf("justificados inexistente esperava exit 2, obtido %d", code)
	}

	op := Operador{
		Nome: "op1",
		Gerar: func(fset *token.FileSet, arq *ast.File, pacote, caminho string) ([]Mutante, error) {
			return []Mutante{{ID: "M1", Pacote: pacote, Arquivo: filepath.Base(caminho), Linha: 2, Operador: "op1"}}, nil
		},
	}

	out.Reset()
	errBuf.Reset()
	if code := executar([]string{"./pkg"}, &out, &errBuf, dir, &executorFalso{passou: false}, []Operador{op}); code != 0 {
		t.Fatalf("todos mortos esperava exit 0, obtido %d (stderr=%s)", code, errBuf.String())
	}

	out.Reset()
	errBuf.Reset()
	if code := executar([]string{"./pkg"}, &out, &errBuf, dir, &executorFalso{passou: true}, []Operador{op}); code != 2 {
		t.Fatalf("sobrevivente esperava exit 2, obtido %d", code)
	}
}

func TestListarArquivosGo_RespeitaBuildTags(t *testing.T) {
	dir := t.TempDir()
	ativo := filepath.Join(dir, "comum.go")
	outroSO := filepath.Join(dir, "incompativel.go")
	if err := os.WriteFile(ativo, []byte("package p\nfunc F() int { return 1 }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outroSO, []byte("//go:build ignore\n\npackage p\nfunc G() int { return 2 }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	arqs, err := listarArquivosGo(dir)
	if err != nil {
		t.Fatalf("listarArquivosGo falhou: %v", err)
	}
	if len(arqs) != 1 || filepath.Base(arqs[0]) != "comum.go" {
		t.Fatalf("esperava apenas [comum.go], obtido %v", arqs)
	}
}
