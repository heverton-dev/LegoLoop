package studio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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
	mux.HandleFunc("POST /api/sessao", s.handleA2UISessaoCriar)
	mux.HandleFunc("GET /api/sessao/{id}/stream", s.handleA2UISessaoStream)
	mux.HandleFunc("POST /api/sessao/{id}/action", s.handleA2UISessaoAction)
	mux.HandleFunc("GET /api/sessao/{id}/artefato", s.handleA2UISessaoArtefato)
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
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleA2UISessaoStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "SSE não suportado")
		return
	}
	ch, err := s.A2UIHub().Subscribe(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for env := range ch {
		data, mErr := json.Marshal(env)
		if mErr != nil {
			continue
		}
		_, _ = fmt.Fprintf(w, "event: a2ui\ndata: %s\n\n", data)
		flusher.Flush()
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"ok":true}`))
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
	_, _ = w.Write(data)
}
