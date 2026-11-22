package studio

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
)

const dagPecaUmNo = `{"name":"p","version":"2.0.0","trigger":{"type":"manual"},"nodes":[
 {"id":"a","type":"script","command":["go","version"],"timeout_sec":60,
  "routing":{"mode":"manual","model":"none"},
  "output":{"artifact":"GO_VERSION_PECAS.txt","gate":"exit_code == 0"}}]}`

// TestRegistraPecaAutomatica (AC-E07.3): dado processo criado; a peça
// correspondente está registrada e invocável no Studio (arrastável na interface),
// sem registro manual.
func TestRegistraPecaAutomatica(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()
	procDir := filepath.Join(s.ProjectDir, caminhos.DirProcessos, "meu-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(caminhos.CaminhoDAG(procDir), []byte(dagPecaUmNo), 0o644); err != nil {
		t.Fatal(err)
	}

	// 1. A peça do processo está registrada sem registro manual.
	rec := do(t, h, "GET", "/api/pecas", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/pecas code=%d body=%s", rec.Code, rec.Body.String())
	}
	var lista struct {
		Pecas []struct {
			ID      string `json:"id"`
			Nome    string `json:"nome"`
			Invocar string `json:"invocar"`
		} `json:"pecas"`
	}
	decode(t, rec, &lista)
	var peca *struct {
		ID      string `json:"id"`
		Nome    string `json:"nome"`
		Invocar string `json:"invocar"`
	}
	for i := range lista.Pecas {
		if lista.Pecas[i].ID == "meu-proc" {
			peca = &lista.Pecas[i]
		}
	}
	if peca == nil {
		t.Fatalf("peça de meu-proc não registrada: %+v", lista.Pecas)
	}
	if !strings.Contains(peca.Invocar, "legoloop exec meu-proc") {
		t.Errorf("peça não delega a legoloop exec: %q", peca.Invocar)
	}

	// 2. A peça é invocável pela UI: exec roda o MESMO processo.
	rec2 := do(t, h, "POST", "/api/pecas/meu-proc/exec", `{"allow_scripts": true}`)
	if rec2.Code != http.StatusOK {
		t.Fatalf("POST /api/pecas/meu-proc/exec code=%d body=%s", rec2.Code, rec2.Body.String())
	}
	scanner := bufio.NewScanner(rec2.Body)
	var eventos []string
	for scanner.Scan() {
		if ln := scanner.Text(); strings.HasPrefix(ln, "data: ") {
			eventos = append(eventos, strings.TrimPrefix(ln, "data: "))
		}
	}
	if len(eventos) == 0 {
		t.Fatal("esperava eventos SSE da invocação")
	}
	var ultimo map[string]any
	if err := json.Unmarshal([]byte(eventos[len(eventos)-1]), &ultimo); err != nil {
		t.Fatalf("evento final inválido: %v", err)
	}
	if ultimo["type"] != "dag_finish" || ultimo["ok"] != true {
		t.Fatalf("esperava dag_finish ok, obteve %v", ultimo)
	}

	// 3. Mesmo contrato da via CLI: artefato do processo no workdir do Studio.
	arq := filepath.Join(s.ProjectDir, "GO_VERSION_PECAS.txt")
	if dados, err := os.ReadFile(arq); err != nil || !strings.Contains(string(dados), "go1.") {
		t.Errorf("artefato do processo ausente em %s: %v (%s)", arq, err, dados)
	}
}

// TestPecaExecPathJail: slug fora do padrão não escapa da pasta de Processos.
func TestPecaExecPathJail(t *testing.T) {
	s, _ := newServer(t)
	for _, slug := range []string{"a.b", "meu%20proc", "%2e%2e%2fpasswd"} {
		rec := do(t, s.Handler(), "POST", "/api/pecas/"+slug+"/exec", `{"allow_scripts": true}`)
		if rec.Code != http.StatusBadRequest && rec.Code != http.StatusNotFound {
			t.Errorf("slug %q deveria ser recusado, obteve %d (%s)", slug, rec.Code, rec.Body.String())
		}
	}
}
