package plangen_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/plangen"
)

type routeFileMini struct {
	Nodes []struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Isolation struct {
			Executor string `json:"executor"`
		} `json:"isolation"`
	} `json:"nodes"`
}

// TestTicketExecutor cobre o critério AC-E02.1:
// Tickets em JSON com um ticket native_swarm e os demais sem o campo geram ROUTE
// com native_swarm nesse nó e orchestrix_hermetic nos outros nós de unidade.
func TestTicketExecutor(t *testing.T) {
	rawTickets := `[
		{"id":"TICK-01.1","title":"Node 1 Swarm","epic_id":"E01","nature":"micro","executor":"native_swarm","validation_cmd":"go test ./...","target_files":["a.go"]},
		{"id":"TICK-01.2","title":"Node 2 Default","epic_id":"E01","nature":"micro","validation_cmd":"go test ./...","target_files":["b.go"]}
	]`
	var tickets []plangen.Ticket
	if err := json.Unmarshal([]byte(rawTickets), &tickets); err != nil {
		t.Fatalf("falha ao deserializar tickets: %v", err)
	}

	base, spec := writeSpec(t)
	res, err := plangen.Render(base, spec, tickets, plangen.SupervisionEpic, fixedNow)
	if err != nil {
		t.Fatalf("Render falhou: %v", err)
	}

	data, err := os.ReadFile(res.RouteJSONPath)
	if err != nil {
		t.Fatalf("falha ao ler route.json: %v", err)
	}

	var rf routeFileMini
	if err := json.Unmarshal(data, &rf); err != nil {
		t.Fatalf("falha ao interpretar route.json: %v", err)
	}

	execMap := make(map[string]string)
	for _, n := range rf.Nodes {
		execMap[n.ID] = n.Isolation.Executor
	}

	if got := execMap["TICK-01.1"]; got != "native_swarm" {
		t.Errorf("TICK-01.1 executor = %q; quer native_swarm", got)
	}
	if got := execMap["TICK-01.2"]; got != "orchestrix_hermetic" {
		t.Errorf("TICK-01.2 executor = %q; quer orchestrix_hermetic", got)
	}
}

// TestTicketExecutorIncoerente cobre o critério AC-E02.2:
// Ticket com executor incoerente ou desconhecido é recusado com dag.ErrExecutorIncoerente
// e nenhum arquivo de plan/route é deixado em disco.
func TestTicketExecutorIncoerente(t *testing.T) {
	for _, badExec := range []string{"xyz", "local_process"} {
		t.Run(badExec, func(t *testing.T) {
			rawTickets := fmt.Sprintf(`[
				{"id":"TICK-01.1","title":"Invalid Node","epic_id":"E01","nature":"micro","executor":%q,"validation_cmd":"go test ./...","target_files":["a.go"]}
			]`, badExec)
			var tickets []plangen.Ticket
			if err := json.Unmarshal([]byte(rawTickets), &tickets); err != nil {
				t.Fatalf("falha ao deserializar tickets: %v", err)
			}

			base, spec := writeSpec(t)
			res, err := plangen.Render(base, spec, tickets, plangen.SupervisionEpic, fixedNow)
			if err == nil {
				t.Fatalf("Render deveria ter falhado para executor %q, mas obteve res=%+v", badExec, res)
			}
			if !errors.Is(err, dag.ErrExecutorIncoerente) {
				t.Fatalf("esperava erro embrulhando dag.ErrExecutorIncoerente, obteve: %v", err)
			}

			for _, dir := range []string{filepath.Join(base, "docs", "plans"), filepath.Join(base, "docs", "routes")} {
				matches, _ := filepath.Glob(filepath.Join(dir, "*"))
				if len(matches) > 0 {
					t.Errorf("arquivo(s) gravado(s) em %s: %v", dir, matches)
				}
			}
		})
	}
}
