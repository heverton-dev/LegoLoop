package board

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Consumidores de JSON (CLI, MCP, scripts) esperam lista vazia, nao null.
func TestSwarmJSONVazioNaoEhNull(t *testing.T) {
	sb, err := OpenSwarmBoard(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	st, err := sb.Status()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(st)
	if strings.Contains(string(b), `"workers":null`) {
		t.Fatalf("workers null no quadro vazio: %s", b)
	}
	stalled, err := sb.CheckStall(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	b, _ = json.Marshal(map[string]any{"stalled": stalled})
	if string(b) != `{"stalled":[]}` {
		t.Fatalf("stalled deveria ser lista vazia: %s", b)
	}
}
