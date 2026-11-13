package dissecar

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathJail(t *testing.T) {
	tmpDir := t.TempDir()

	// Caminho válido dentro da sandbox
	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)

	valid, err := ValidatePath(tmpDir, subDir)
	if err != nil {
		t.Fatalf("esperava sucesso no caminho valido, obteve erro: %v", err)
	}
	if valid == "" {
		t.Fatalf("esperava caminho limpo resolvido")
	}

	// Caminho inválido fora da sandbox (escape)
	_, err = ValidatePath(tmpDir, "../fora")
	if err != ErrPathEscape {
		t.Fatalf("esperava ErrPathEscape, obteve: %v", err)
	}
}

func TestExtractorsAndDAG(t *testing.T) {
	tmpDir := t.TempDir()

	// Criar arquivo de teste Go
	code := `package main
import "fmt"
type Dummy struct{}
func Main() { fmt.Println("hello") }`

	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(code), 0644)
	if err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	nodes, err := ExtractAST(tmpDir)
	if err != nil {
		t.Fatalf("ExtractAST falhou: %v", err)
	}
	if len(nodes) == 0 {
		t.Fatalf("esperava nos AST extraidos")
	}

	syscalls, err := ExtractSyscalls(tmpDir)
	if err != nil {
		t.Fatalf("ExtractSyscalls falhou: %v", err)
	}
	if len(syscalls) == 0 || syscalls[0] != "fmt" {
		t.Fatalf("esperava syscall/import fmt, obteve: %v", syscalls)
	}

	dag := BuildReverseDAG(nodes)
	if len(dag) == 0 {
		t.Fatalf("esperava DAG reverso sintetizado")
	}
}
