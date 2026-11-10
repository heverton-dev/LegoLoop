// Comando route2dag projeta uma ROUTE (route.v2) em um DAG executável (dag.v1).
//
// Regras da projeção:
//   - command do nó = primeiro argv de output.verify (sem verify, sem command);
//   - output.artifact = scratch/gates/<id>.json, único por nó;
//   - Stop-Gate vira nó com roteamento manual (human/operator);
//   - o gate de cada nó passa sem mudança.
//
// Uso: go run ./tools/route2dag <route.json> [saida.json]
// Sem saída explícita, grava docs/dags/<data>_DAG_<NOME>.json no diretório atual.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "uso: route2dag <route.json> [saida.json]")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	out, nome, err := Converter(raw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	dest := filepath.Join("docs", "dags", nome)
	if len(os.Args) == 3 {
		dest = os.Args[2]
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	fmt.Println("DAG escrito em", dest)
}
