package studio

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

type ProcessoAtivoItem struct {
	ID          string            `json:"id"`
	Processo    string            `json:"processo"`
	State       a2ui.SessionState `json:"state"`
	Surfaces    int               `json:"surfaces"`
	Subscribers int               `json:"subscribers"`
	URL         string            `json:"url"`
	TimestampMs int64             `json:"timestamp_ms"`
}

func (s *Server) handleProcessosAtivos(w http.ResponseWriter, r *http.Request) {
	sessions := s.A2UIHub().ListSessions()
	itens := make([]ProcessoAtivoItem, 0, len(sessions))
	now := time.Now().UnixMilli()
	for _, sess := range sessions {
		itens = append(itens, ProcessoAtivoItem{
			ID:          sess.ID,
			Processo:    sess.Processo,
			State:       sess.State,
			Surfaces:    sess.Surfaces,
			Subscribers: sess.Subscribers,
			URL:         "/sessao/" + sess.ID,
			TimestampMs: now,
		})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(itens); err != nil {
		return
	}
}

func (s *Server) handleProcessoAndonPause(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonPause(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "paused"})
}

func (s *Server) handleProcessoAndonResume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonResume(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "running"})
}

func (s *Server) handleProcessoAndonCancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonCancel(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "cancelled"})
}
