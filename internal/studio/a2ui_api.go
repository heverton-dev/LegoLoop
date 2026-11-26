package studio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

type sessaoCriarReq struct {
	Processo string `json:"processo"`
}

type sessaoCriarResp struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
	StreamURL string `json:"stream_url"`
	ActionURL string `json:"action_url"`
}

// A2UIHub devolve o barramento de sessões A2UI do Studio, inicializando sob demanda.
func (s *Server) A2UIHub() *a2ui.Hub {
	s.a2uiOnce.Do(func() {
		if s.A2UI == nil {
			s.A2UI = a2ui.NewHub(a2ui.DefaultCatalog())
		}
	})
	return s.A2UI
}

func (s *Server) registerA2UIRoutes(mux *http.ServeMux) {
	s.garantirSessaoLaboratorio("sess_1")
	s.garantirSessaoLaboratorio("sess_2")
	mux.HandleFunc("GET /concepcao", s.handleA2UIConcepcaoPage)
	mux.HandleFunc("GET /fabricar", s.handleA2UIFabricarPage)
	mux.HandleFunc("POST /api/sessao", s.handleA2UISessaoCriar)
	mux.HandleFunc("GET /api/sessao/{id}/stream", s.handleA2UISessaoStream)
	mux.HandleFunc("POST /api/sessao/{id}/action", s.handleA2UISessaoAction)
	mux.HandleFunc("GET /api/sessao/{id}/timer", s.handleA2UISessaoTimerGet)
	mux.HandleFunc("POST /api/sessao/{id}/timer", s.handleA2UISessaoTimerPost)
	mux.HandleFunc("GET /api/sessao/{id}/artefato", s.handleA2UISessaoArtefato)
}

func (s *Server) garantirSessaoLaboratorio(id string) {
	switch id {
	case "sess_1":
		if err := s.garantirConcepcaoSessao(id); err != nil {
			return
		}
	case "sess_2":
		if err := s.garantirFabricarSessao(id); err != nil {
			return
		}
	}
}

func (s *Server) handleA2UISessaoCriar(w http.ResponseWriter, r *http.Request) {
	var req sessaoCriarReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	sess := s.A2UIHub().CreateSession(req.Processo)
	resp := sessaoCriarResp{
		SessionID: sess.ID,
		URL:       "/sessao/" + sess.ID,
		StreamURL: "/api/sessao/" + sess.ID + "/stream",
		ActionURL: "/api/sessao/" + sess.ID + "/action",
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}

func (s *Server) handleA2UISessaoStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "SSE não suportado")
		return
	}
	id := r.PathValue("id")
	s.garantirSessaoLaboratorio(id)
	ch, err := s.A2UIHub().Subscribe(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	ctrlCh, err := s.A2UIHub().SubscribeControl(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	bombarStreamA2UI(w, flusher, ch, ctrlCh)
}

func bombarStreamA2UI(w http.ResponseWriter, flusher http.Flusher, ch <-chan a2ui.Envelope, ctrlCh <-chan a2ui.ControlEvent) {
	for ch != nil || ctrlCh != nil {
		select {
		case env, ok := <-ch:
			if !ok {
				ch = nil
				continue
			}
			if !emitirEventoSSE(w, flusher, "a2ui", env) {
				return
			}
		case ev, ok := <-ctrlCh:
			if !ok {
				ctrlCh = nil
				continue
			}
			if !emitirEventoControleSSE(w, flusher, ev) {
				return
			}
		}
	}
}

func emitirEventoControleSSE(w http.ResponseWriter, flusher http.Flusher, ev a2ui.ControlEvent) bool {
	if ev.Kind == "timer" && ev.Timer != nil {
		return emitirEventoSSE(w, flusher, "timer", ev.Timer)
	}
	if ev.Kind == "status" && ev.Status != nil {
		return emitirEventoSSE(w, flusher, "status", ev.Status)
	}
	return true
}

func emitirEventoSSE(w http.ResponseWriter, flusher http.Flusher, nome string, payload any) bool {
	data, err := json.Marshal(payload)
	if err != nil {
		return true
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", nome, data); err != nil {
		return false
	}
	flusher.Flush()
	return true
}

type sessaoTimerReq struct {
	Acao      string `json:"acao"`
	TimeoutMs int64  `json:"timeout_ms,omitempty"`
}

func (s *Server) handleA2UISessaoTimerGet(w http.ResponseWriter, r *http.Request) {
	sess, err := s.A2UIHub().GetSession(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(sess.Timer().Snapshot()); err != nil {
		return
	}
}

func (s *Server) handleA2UISessaoTimerPost(w http.ResponseWriter, r *http.Request) {
	sess, err := s.A2UIHub().GetSession(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	var req sessaoTimerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	snap, ok := aplicarAcaoTimerSessao(sess, req)
	if !ok {
		writeError(w, http.StatusBadRequest, "ação de timer inválida")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(snap); err != nil {
		return
	}
}

func aplicarAcaoTimerSessao(sess *a2ui.Session, req sessaoTimerReq) (a2ui.TimerSnapshot, bool) {
	switch req.Acao {
	case "pausar":
		return sess.PauseTimer(), true
	case "retomar":
		return sess.ResumeTimer(), true
	case "configurar":
		if req.TimeoutMs <= 0 {
			return a2ui.TimerSnapshot{}, false
		}
		return sess.ConfigureTimer(time.Duration(req.TimeoutMs) * time.Millisecond), true
	case "estender":
		ms := req.TimeoutMs
		if ms <= 0 {
			ms = 300000
		}
		return sess.ExtendTimer(time.Duration(ms) * time.Millisecond), true
	default:
		return a2ui.TimerSnapshot{}, false
	}
}

func (s *Server) handleA2UISessaoAction(w http.ResponseWriter, r *http.Request) {
	var act a2ui.Action
	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	act.SessionID = r.PathValue("id")
	if act.Origin == "" {
		act.Origin = "ui"
	}
	if err := s.A2UIHub().SubmitAction(act); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, a2ui.ErrSessaoNaoEncontrada) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	if err := s.dispatchAndonAction(r.Context(), act); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := s.dispatchConcepcaoAction(r.Context(), act); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
		return
	}
}

func (s *Server) handleA2UISessaoArtefato(w http.ResponseWriter, r *http.Request) {
	if _, err := s.A2UIHub().GetSession(r.PathValue("id")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	alvo, err := jail.DentroDaRaiz(s.ProjectDir, r.URL.Query().Get("path"))
	if err != nil {
		writeError(w, http.StatusForbidden, "caminho fora da raiz permitida")
		return
	}
	data, err := os.ReadFile(alvo)
	if err != nil {
		writeError(w, http.StatusNotFound, "artefato não encontrado")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		return
	}
}
