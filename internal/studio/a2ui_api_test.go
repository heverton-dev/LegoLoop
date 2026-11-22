package studio

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func TestA2UISessaoCriarStreamEAction(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	rec := do(t, h, http.MethodPost, "/api/sessao", `{"processo":"legoloop-concepcao"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/sessao code=%d body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		SessionID string `json:"session_id"`
		URL       string `json:"url"`
		StreamURL string `json:"stream_url"`
		ActionURL string `json:"action_url"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal resposta: %v", err)
	}
	if created.SessionID == "" || !strings.Contains(created.URL, created.SessionID) {
		t.Fatalf("resposta incompleta: %+v", created)
	}

	env := a2ui.Envelope{
		SessionID:   created.SessionID,
		SurfaceID:   "proc:concepcao",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		Components:  []a2ui.ComponentRef{{ID: "c1", Kind: "Card", Props: map[string]any{"title": "Oi"}}},
	}
	if err := s.A2UIHub().Publish(env); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	ts := httptest.NewServer(h)
	defer ts.Close()
	s.Host = strings.TrimPrefix(ts.URL, "http://")

	resp, err := http.Get(ts.URL + created.StreamURL)
	if err != nil {
		t.Fatalf("GET stream: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("stream headers inesperados: %d %s", resp.StatusCode, resp.Header.Get("Content-Type"))
	}

	sc := bufio.NewScanner(resp.Body)
	achou := false
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "data: ") && strings.Contains(line, `"surface_id":"proc:concepcao"`) {
			achou = true
			break
		}
	}
	if !achou {
		t.Fatal("esperava evento SSE com surface_id proc:concepcao")
	}

	s.Host = host
	sess, err := s.A2UIHub().GetSession(created.SessionID)
	if err != nil {
		t.Fatalf("GetSession: %v", err)
	}
	actRec := do(t, h, http.MethodPost, created.ActionURL, `{"surface_id":"proc:concepcao","origin":"ui","name":"submit_answer","payload":{"answer":"sim"}}`)
	if actRec.Code != http.StatusAccepted {
		t.Fatalf("POST action code=%d body=%s", actRec.Code, actRec.Body.String())
	}
	select {
	case got := <-sess.Actions():
		if got.Name != "submit_answer" || got.Payload["answer"] != "sim" {
			t.Fatalf("action inesperada: %+v", got)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout esperando action na sessão")
	}
}

func TestA2UIArtefatoRespeitaJail(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()
	sess := s.A2UIHub().CreateSession("piloto")

	arq := filepath.Join(s.ProjectDir, "docs", "INTENT.md")
	if err := os.MkdirAll(filepath.Dir(arq), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(arq, []byte("# INTENT OK\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	okRec := do(t, h, http.MethodGet, "/api/sessao/"+sess.ID+"/artefato?path=docs/INTENT.md", "")
	if okRec.Code != http.StatusOK || !strings.Contains(okRec.Body.String(), "# INTENT OK") {
		t.Fatalf("esperava 200 com conteúdo do artefato, obteve %d %q", okRec.Code, okRec.Body.String())
	}

	escapeRec := do(t, h, http.MethodGet, "/api/sessao/"+sess.ID+"/artefato?path=../../etc/passwd", "")
	if escapeRec.Code != http.StatusForbidden {
		t.Fatalf("esperava 403 para path traversal, obteve %d %q", escapeRec.Code, escapeRec.Body.String())
	}
}

func TestA2UISecurityBloqueiaHostEOrigin(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	recHost := do(t, h, http.MethodPost, "/api/sessao", `{"processo":"x"}`, func(r *http.Request) {
		r.Host = "evil.example:7878"
	})
	if recHost.Code != http.StatusMisdirectedRequest {
		t.Fatalf("esperava 421 para host externo, obteve %d", recHost.Code)
	}

	recOrigin := do(t, h, http.MethodPost, "/api/sessao", `{"processo":"x"}`, func(r *http.Request) {
		r.Header.Set("Origin", "https://evil.example")
	})
	if recOrigin.Code != http.StatusForbidden {
		t.Fatalf("esperava 403 para Origin cruzado, obteve %d", recOrigin.Code)
	}
}
