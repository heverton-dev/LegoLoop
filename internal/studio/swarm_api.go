package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"legoloop/internal/board"
)

func (s *Server) registerSwarmRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swarm", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(studioSwarmHTML))
	})
	mux.HandleFunc("GET /api/swarm/events", s.handleSwarmEvents)
	mux.HandleFunc("POST /api/swarm/add", s.handleSwarmAdd)
	mux.HandleFunc("POST /api/swarm/post", s.handleSwarmPost)
	mux.HandleFunc("POST /api/swarm/kill", s.handleSwarmKill)
}

func (s *Server) swarmRoot() string {
	if s.ProjectDir != "" {
		return filepath.Join(s.ProjectDir, "scratch")
	}
	return "scratch"
}

func studioBoardsPayload(root string) ([]byte, error) {
	boards := board.DiscoverBoards(root, "", "")
	names := make([]string, 0, len(boards))
	for k := range boards {
		names = append(names, k)
	}
	sort.Strings(names)
	procs := map[string]any{}
	for _, name := range names {
		sb, err := board.OpenSwarmBoard(boards[name])
		if err != nil {
			continue
		}
		st, _ := sb.Status()
		msgs, _ := sb.ReadMessages(0, nil)
		procs[name] = map[string]any{"workers": st.Workers, "messages": msgs}
	}
	return json.Marshal(map[string]any{"order": names, "processes": procs})
}

func (s *Server) handleSwarmAdd(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process       string              `json:"process"`
		Name          string              `json:"name"`
		Role          string              `json:"role"`
		Model         string              `json:"model"`
		WorkspaceMode board.WorkspaceMode `json:"workspace_mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	boards := board.DiscoverBoards(s.swarmRoot(), "", "")
	bdir, ok := boards[req.Process]
	if !ok {
		http.Error(w, "processo não encontrado", http.StatusNotFound)
		return
	}
	sb, err := board.OpenSwarmBoard(bdir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	worker, err := sb.AddWorker(board.AddWorkerRequest{
		Name:          req.Name,
		Role:          req.Role,
		Model:         req.Model,
		WorkspaceMode: req.WorkspaceMode,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(worker)
}

func (s *Server) handleSwarmEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE nao suportado", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	var last string
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()
	for {
		if data, err := studioBoardsPayload(s.swarmRoot()); err == nil && string(data) != last {
			last = string(data)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", last)
			flusher.Flush()
		}
		select {
		case <-r.Context().Done():
			return
		case <-tick.C:
		}
	}
}

func (s *Server) handleSwarmPost(w http.ResponseWriter, r *http.Request) {
	var req struct{ Proc, To, Body string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dir := board.DiscoverBoards(s.swarmRoot(), "", "")[req.Proc]
	if dir == "" {
		http.Error(w, "processo nao encontrado", http.StatusNotFound)
		return
	}
	sb, err := board.OpenSwarmBoard(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = sb.PostMessage("coordenador", req.To, req.Body)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSwarmKill(w http.ResponseWriter, r *http.Request) {
	var req struct{ Proc, Name string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dir := board.DiscoverBoards(s.swarmRoot(), "", "")[req.Proc]
	if dir == "" {
		http.Error(w, "processo nao encontrado", http.StatusNotFound)
		return
	}
	sb, err := board.OpenSwarmBoard(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = sb.KillWorker(req.Name, ".")
	w.WriteHeader(http.StatusOK)
}

const studioSwarmHTML = `<!doctype html>
<html lang="pt-BR"><head><meta charset="utf-8"><title>LegoLoop Studio — Swarm</title>
<link rel="stylesheet" href="/style.css"><link rel="stylesheet" href="/swarm.css">
</head><body>
<div class="app">
  <div class="top">
    <div class="brand"><span class="dot"></span><span>LegoLoop Swarm</span><span class="crumb" id="crumb"></span></div>
    <div class="hdr-right"><div class="meta" id="summary">Conectando SSE...</div><a class="btn link-clean" href="/">← Studio</a></div>
  </div>
  <div class="tabs" id="tabs"></div>
  <div class="grid">
    <div class="card"><div class="card-hd"><h3>Agentes</h3><span class="meta" id="counts"></span></div><div class="scroll" id="workers"></div></div>
    <div class="card"><div class="card-hd"><h3>Execução <span id="filter"></span></h3><span class="meta">Mais recente no topo</span></div><div id="ficha"></div><div class="scroll" id="msgs"></div><div class="pager" id="pager"></div></div>
  </div>
</div>
<script src="/swarm.js"></script></body></html>`
