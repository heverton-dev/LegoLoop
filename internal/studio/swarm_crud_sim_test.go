package studio

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
)

func TestSwarmUpdateRemoveESimulateAPI(t *testing.T) {
	proj := initSwarmGitRepo(t)
	if _, err := board.OpenSwarmBoard(proj); err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	s, _ := newServer(t)
	s.ProjectDir = proj
	h := s.Handler()
	procName := filepath.Base(proj)

	// 1. Add com peça e ticket
	recAdd := do(t, h, http.MethodPost, "/api/swarm/add", `{
		"process":"`+procName+`",
		"name":"w-crud",
		"role":"engenheiro-backend",
		"model":"t-codigo-128k",
		"piece":"aidd-tdd",
		"piece_kind":"skill",
		"ticket":"Construir API de usuários"
	}`)
	if recAdd.Code != http.StatusOK {
		t.Fatalf("add falhou: %d %s", recAdd.Code, recAdd.Body.String())
	}
	var w board.Worker
	decode(t, recAdd, &w)
	if w.Piece != "aidd-tdd" || w.Ticket != "Construir API de usuários" {
		t.Fatalf("peça/ticket ausentes na resposta: %+v", w)
	}

	// 2. Update altera papel, peça e ticket
	recUpd := do(t, h, http.MethodPost, "/api/swarm/update", `{
		"process":"`+procName+`",
		"name":"w-crud",
		"role":"engenheiro-tdd",
		"piece":"aidd-spec",
		"piece_kind":"skill",
		"ticket":"Validar especificação EARS"
	}`)
	if recUpd.Code != http.StatusOK {
		t.Fatalf("update falhou: %d %s", recUpd.Code, recUpd.Body.String())
	}
	decode(t, recUpd, &w)
	if w.Role != "engenheiro-tdd" || w.Piece != "aidd-spec" || w.Ticket != "Validar especificação EARS" {
		t.Fatalf("update não persistiu: %+v", w)
	}

	// 3. Remove exclui definitivamente do board
	recRem := do(t, h, http.MethodPost, "/api/swarm/remove", `{"process":"`+procName+`","name":"w-crud"}`)
	if recRem.Code != http.StatusOK {
		t.Fatalf("remove falhou: %d %s", recRem.Code, recRem.Body.String())
	}
	sb, _ := board.OpenSwarmBoard(proj)
	st, _ := sb.Status()
	if len(st.Workers) != 0 {
		t.Fatalf("esperava 0 workers após remove, obteve %+v", st.Workers)
	}

	// 4. Simulate cria processo de simulação com ordem síncrona/assíncrona
	recSim := do(t, h, http.MethodPost, "/api/swarm/simulate", `{
		"name":"meu-processo",
		"step_delay_ms": 10,
		"nodes":[
			{"id":"t1","label":"Manual","kind":"trigger_manual","role":"orquestrador","artifact":"Briefing"},
			{"id":"s1","label":"Skill A","kind":"skill","role":"arquiteto","artifact":"Spec"},
			{"id":"s2","label":"Skill B","kind":"skill","role":"planejador","artifact":"Plano"}
		],
		"edges":[
			{"from":"t1","to":"s1"},
			{"from":"t1","to":"s2"}
		]
	}`)
	if recSim.Code != http.StatusOK {
		t.Fatalf("simulate falhou: %d %s", recSim.Code, recSim.Body.String())
	}
	var simResp struct {
		OK      bool   `json:"ok"`
		Process string `json:"process"`
		URL     string `json:"url"`
		Steps   []struct {
			Step  int      `json:"step"`
			Mode  string   `json:"mode"`
			Nodes []string `json:"nodes"`
		} `json:"steps"`
	}
	decode(t, recSim, &simResp)
	if !simResp.OK || !strings.HasPrefix(simResp.Process, "sim-") || len(simResp.Steps) != 2 {
		t.Fatalf("resposta simulate inesperada: %+v", simResp)
	}
	if simResp.Steps[0].Mode != "sync" || simResp.Steps[1].Mode != "async" {
		t.Fatalf("esperava passo 1 sync e passo 2 async, obteve: %+v", simResp.Steps)
	}
	// Aguarda goroutine de simulação concluir antes do TempDir cleanup (condição binária, teto generoso).
	simDir := filepath.Join(proj, "scratch", simResp.Process, ".legoloop-swarm", "swarm.json")
	deadline := time.Now().Add(30 * time.Second)
	var finalData string
	for time.Now().Before(deadline) {
		if b, err := os.ReadFile(simDir); err == nil && strings.Contains(string(b), "Simulação concluída") {
			finalData = string(b)
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !strings.Contains(finalData, `"status": "done"`) {
		t.Fatalf("esperava workers finalizados com status done apos simulacao, obteve:\n%s", finalData)
	}
}

func TestSwarmMetaAndMultiPiecesInferTicketAPI(t *testing.T) {
	proj := initSwarmGitRepo(t)
	if _, err := board.OpenSwarmBoard(proj); err != nil {
		t.Fatalf("OpenSwarmBoard: %v", err)
	}
	s, _ := newServer(t)
	s.ProjectDir = proj
	h := s.Handler()
	procName := filepath.Base(proj)

	recMeta := do(t, h, http.MethodGet, "/api/swarm/meta", "")
	if recMeta.Code != http.StatusOK {
		t.Fatalf("GET /api/swarm/meta status=%d body=%s", recMeta.Code, recMeta.Body.String())
	}
	var meta struct {
		Roles      []map[string]any `json:"roles"`
		Models     []map[string]any `json:"models"`
		Workspaces []map[string]any `json:"workspaces"`
	}
	decode(t, recMeta, &meta)
	if len(meta.Roles) != 20 {
		t.Fatalf("esperava 20 papeis canonicos em /api/swarm/meta, obteve %d", len(meta.Roles))
	}
	if len(meta.Models) == 0 || len(meta.Workspaces) != 3 {
		t.Fatalf("esperava modelos e 3 workspaces em /api/swarm/meta, obteve %+v", meta)
	}

	recAdd := do(t, h, http.MethodPost, "/api/swarm/add", `{
		"process":"`+procName+`",
		"name":"w-infer",
		"role":"engenheiro-backend",
		"workspace_mode":"branch",
		"pieces":[{"kind":"skill","name":"aidd-tdd"},{"kind":"mcp","name":"context7"}]
	}`)
	if recAdd.Code != http.StatusOK {
		t.Fatalf("add com multi-pecas falhou: %d %s", recAdd.Code, recAdd.Body.String())
	}
	var w board.Worker
	decode(t, recAdd, &w)
	if len(w.Pieces) != 2 || w.Ticket == "" || !strings.Contains(w.Ticket, "aidd-tdd") {
		t.Fatalf("esperava multi-pecas e ticket inferido automaticamente, obteve: %+v", w)
	}
}
