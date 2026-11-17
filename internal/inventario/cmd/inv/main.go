// Comando inv opera o saneamento de docs/: gerar inventário, migrar arquivo morto e guardar a lei.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/inventario"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: inv gerar")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "gerar":
		gerar()
	default:
		fmt.Fprintln(os.Stderr, "comando desconhecido:", os.Args[1])
		os.Exit(2)
	}
}

// gerar escreve o inventário completo de docs/ em docs/reports/.
func gerar() {
	raiz, err := raizRepo()
	if err != nil {
		fmt.Fprintln(os.Stderr, "raiz do repositório:", err)
		os.Exit(2)
	}
	docs := filepath.Join(raiz, "docs")
	entradas, err := inventario.Varrer(docs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "varredura:", err)
		os.Exit(2)
	}
	saida := filepath.Join(docs, "reports", "2026-11-16_INVENTARIO_DOCS.md")
	if err := inventario.GerarRelatorio(entradas, saida); err != nil {
		fmt.Fprintln(os.Stderr, "relatório:", err)
		os.Exit(2)
	}
	fmt.Printf("inventário gravado em %s (%d arquivos)\n", saida, len(entradas))
}

// raizRepo sobe do diretório atual até o go.mod do repositório.
func raizRepo() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return "", fmt.Errorf("go.mod não encontrado")
		}
		dir = pai
	}
}
