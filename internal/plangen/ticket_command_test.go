package plangen_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/routedag"
)

func renderRoute(t *testing.T, raw string) []byte {
	t.Helper()
	var tickets []plangen.Ticket
	if err := json.Unmarshal([]byte(raw), &tickets); err != nil {
		t.Fatalf("tickets: %v", err)
	}
	base, spec := writeSpec(t)
	res, err := plangen.Render(base, spec, tickets, plangen.SupervisionEpic, fixedNow)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	data, err := os.ReadFile(res.RouteJSONPath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestTicketCommandNativeSwarm cobre a correção do A13: ticket native_swarm com command do executor
// gera ROUTE e DAG que o swarm consegue subir (papel e command presentes, command do nó vence o gate).
func TestTicketCommandNativeSwarm(t *testing.T) {
	route := renderRoute(t, `[
		{"id":"TICK-01.1","title":"Swarm","epic_id":"E01","nature":"micro","executor":"native_swarm","command":["claude","-p","faca o ticket"],"validation_cmd":"go test ./...","target_files":["a.go"]}
	]`)
	dagBytes, _, err := routedag.Converter(route)
	if err != nil {
		t.Fatalf("Converter: %v", err)
	}
	d, err := dag.Parse(dagBytes)
	if err != nil {
		t.Fatalf("dag.Parse: %v", err)
	}
	n := d.Node("TICK-01.1")
	want := []string{"claude", "-p", "faca o ticket"}
	if !reflect.DeepEqual(n.Command, want) {
		t.Errorf("command = %v; quer %v", n.Command, want)
	}
	levels, _, err := board.LevelsFromDAG(d, nil)
	if err != nil || len(levels) == 0 || levels[0][0].NodeID != "TICK-01.1" || !reflect.DeepEqual(levels[0][0].Command, want) {
		t.Fatalf("swarm não sobe o nó: levels=%v err=%v", levels, err)
	}
}

// TestTicketCommandAditivo prova que ticket sem command gera a mesma ROUTE de antes (sem a chave command).
func TestTicketCommandAditivo(t *testing.T) {
	route := renderRoute(t, `[
		{"id":"TICK-01.1","title":"Padrão","epic_id":"E01","nature":"micro","validation_cmd":"go test ./...","target_files":["a.go"]}
	]`)
	var rf struct {
		Nodes []map[string]any `json:"nodes"`
	}
	if err := json.Unmarshal(route, &rf); err != nil {
		t.Fatal(err)
	}
	for _, n := range rf.Nodes {
		if _, ok := n["command"]; ok {
			t.Errorf("nó %v ganhou command sem pedido", n["id"])
		}
		if _, ok := n["role"]; ok {
			t.Errorf("nó %v ganhou role sem pedido", n["id"])
		}
	}
}
