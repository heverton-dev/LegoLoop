// Comando route2dag projeta uma ROUTE (route.v2) em um DAG executável (dag.v1).
//
// Regras da projeção:
//   - command do nó = primeiro argv de output.verify (sem verify, sem command);
//   - output.artifact = scratch/gates/<id>.json, único por nó;
//   - Stop-Gate vira nó com roteamento manual (human/operator);
//   - o gate de cada nó passa sem mudança;
//   - o DAG projetado precisa passar no parser (dag.Parse) antes de ser gravado.
//
// Uso: go run ./tools/route2dag <route.json> [saida.json]
// Sem saída explícita, grava docs/dags/<data>_DAG_<NOME>.json no diretório atual.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
)

// ErrArtefatoRepetido indica dois nós da projeção com o mesmo output.artifact.
var ErrArtefatoRepetido = errors.New("route2dag: dois nós gravam o mesmo artefato")

// Projetar converte a ROUTE e só devolve o DAG se ele tiver artefatos únicos e passar no parser.
func Projetar(raw []byte) ([]byte, string, error) {
	out, nome, err := Converter(raw)
	if err != nil {
		return nil, "", err
	}
	if out, err = copiarHarness(raw, out); err != nil {
		return nil, "", err
	}
	if err := artefatosUnicos(out); err != nil {
		return nil, "", err
	}
	if _, err := dag.Parse(out); err != nil {
		return nil, "", fmt.Errorf("route2dag: DAG projetado não valida: %w", err)
	}
	return out, nome, nil
}

func copiarHarness(raw, out []byte) ([]byte, error) {
	var r route
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("ROUTE inválida: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(out))
	dec.UseNumber()
	var d map[string]any
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}
	nodes, _ := d["nodes"].([]any)
	for i, n := range r.Nodes {
		nd, _ := nodes[i].(map[string]any)
		rt, _ := nd["routing"].(map[string]any)
		if rt != nil && n.Type != "stop_gate" {
			rt["harness"] = n.Routing.Harness
		}
	}
	res, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(res, '\n'), nil
}

// artefatosUnicos recusa dois nós com o mesmo output.artifact (comparação sem caixa, como no parser).
func artefatosUnicos(out []byte) error {
	var d struct {
		Nodes []struct {
			ID     string `json:"id"`
			Output struct {
				Artifact string `json:"artifact"`
			} `json:"output"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(out, &d); err != nil {
		return err
	}
	vistos := map[string]string{}
	for _, n := range d.Nodes {
		key := strings.ToLower(filepath.ToSlash(n.Output.Artifact))
		if outro, dup := vistos[key]; dup {
			return fmt.Errorf("%w: %q e %q gravam %q", ErrArtefatoRepetido, outro, n.ID, n.Output.Artifact)
		}
		vistos[key] = n.ID
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "uso: route2dag <route.json> [saida.json]")
		os.Exit(2)
	}
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	raw, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	out, nome, err := Projetar(raw)
	if err != nil {
		return err
	}
	dest := filepath.Join("docs", "dags", nome)
	if len(args) == 2 {
		dest = args[1]
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		return err
	}
	fmt.Println("DAG escrito em", dest)
	return nil
}
