package studio

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func TestAndonCordPauseResumeCancel(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	sess := s.A2UIHub().EnsureSession("sess_teste_andon", "legoloop-fabricar")
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava estado running, obteve %s", sess.GetState())
	}

	// 1. Pausa via ação Andon Cord
	actionURL := "/api/sessao/" + sess.ID + "/action"
	pRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_pause"}`)
	if pRec.Code != http.StatusAccepted {
		t.Fatalf("POST pause code=%d body=%s", pRec.Code, pRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionPaused {
		t.Fatalf("esperava estado paused, obteve %s", sess.GetState())
	}

	// 2. Retomada via ação Andon Cord
	rRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_resume"}`)
	if rRec.Code != http.StatusAccepted {
		t.Fatalf("POST resume code=%d body=%s", rRec.Code, rRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava estado running, obteve %s", sess.GetState())
	}

	// 3. Cancelamento via ação Andon Cord
	cRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_cancel"}`)
	if cRec.Code != http.StatusAccepted {
		t.Fatalf("POST cancel code=%d body=%s", cRec.Code, cRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionCancelled {
		t.Fatalf("esperava estado cancelled, obteve %s", sess.GetState())
	}
}

func TestProcessosAtivosAPI(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	sess := s.A2UIHub().EnsureSession("sess_proc_api", "legoloop-fabricar")

	// GET /processos (página HTML)
	pageRec := do(t, h, http.MethodGet, "/processos", "")
	if pageRec.Code != http.StatusOK || !strings.Contains(pageRec.Body.String(), "Gestao de Processos") {
		t.Fatalf("esperava 200 na página /processos, obteve %d", pageRec.Code)
	}

	// GET /api/processos/ativos
	listRec := do(t, h, http.MethodGet, "/api/processos/ativos", "")
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/processos/ativos code=%d", listRec.Code)
	}
	var itens []ProcessoAtivoItem
	if err := json.Unmarshal(listRec.Body.Bytes(), &itens); err != nil {
		t.Fatalf("unmarshal itens: %v", err)
	}
	achou := false
	for _, it := range itens {
		if it.ID == sess.ID && it.Processo == "legoloop-fabricar" {
			achou = true
			break
		}
	}
	if !achou {
		t.Fatalf("sessão %s não listada nos processos ativos", sess.ID)
	}

	// POST /api/processos/{id}/andon
	andRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/andon", "{}")
	if andRec.Code != http.StatusOK {
		t.Fatalf("POST andon code=%d body=%s", andRec.Code, andRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionPaused {
		t.Fatalf("esperava paused após andon, obteve %s", sess.GetState())
	}

	// POST /api/processos/{id}/resume
	resRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/resume", "{}")
	if resRec.Code != http.StatusOK {
		t.Fatalf("POST resume code=%d body=%s", resRec.Code, resRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava running após resume, obteve %s", sess.GetState())
	}

	// POST /api/processos/{id}/cancel
	canRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/cancel", "{}")
	if canRec.Code != http.StatusOK {
		t.Fatalf("POST cancel code=%d body=%s", canRec.Code, canRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionCancelled {
		t.Fatalf("esperava cancelled após cancel, obteve %s", sess.GetState())
	}
}
