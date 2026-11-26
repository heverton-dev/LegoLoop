package studio

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func TestA2UISessaoTimerAPI_ConfigurarPausarRetomarEstender(t *testing.T) {
	srv, _ := newServer(t)
	h := srv.Handler()
	sess := srv.A2UIHub().CreateSession("legoloop-concepcao")
	sess.Timer().Start(0)

	// Configurar para 10 minutos
	rec := do(t, h, http.MethodPost, "/api/sessao/"+sess.ID+"/timer", `{"acao":"configurar","timeout_ms":600000}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("configurar retornou %d: %s", rec.Code, rec.Body.String())
	}
	var snap a2ui.TimerSnapshot
	if err := json.Unmarshal(rec.Body.Bytes(), &snap); err != nil || snap.TimeoutMs != 600000 {
		t.Fatalf("snapshot inesperado apos configurar: %+v err=%v", snap, err)
	}

	// Pausar
	rec = do(t, h, http.MethodPost, "/api/sessao/"+sess.ID+"/timer", `{"acao":"pausar"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("pausar retornou %d", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &snap); err != nil || !snap.Paused {
		t.Fatalf("esperava paused=true, obteve %+v err=%v", snap, err)
	}

	// Retomar
	rec = do(t, h, http.MethodPost, "/api/sessao/"+sess.ID+"/timer", `{"acao":"retomar"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("retomar retornou %d", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &snap); err != nil || snap.Paused {
		t.Fatalf("esperava paused=false, obteve %+v err=%v", snap, err)
	}

	// GET timer
	rec = do(t, h, http.MethodGet, "/api/sessao/"+sess.ID+"/timer", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET timer retornou %d", rec.Code)
	}
}
