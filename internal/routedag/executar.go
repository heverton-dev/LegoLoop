package routedag

import (
	"fmt"
	"io"
	"os"
)

// Executar processa os argumentos de linha de comando para conversão route2dag.
func Executar(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "uso: route2dag <arquivo_route.json>")
		return 2
	}
	raw, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "erro: %v\n", err)
		return 1
	}
	dados, _, err := Converter(raw)
	if err != nil {
		fmt.Fprintf(stderr, "erro: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, string(dados))
	return 0
}
