package studio

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"legoloop/internal/board"
)

func (s *Server) registerSwarmRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swarm", handlePage("swarm.html"))
	mux.HandleFunc("GET /api/swarm/events", s.handleSwarmEvents)
	mux.HandleFunc("GET /api/events", s.handleSwarmEvents)
	mux.HandleFunc("POST /api/swarm/add", s.handleSwarmAdd)
	mux.HandleFunc("POST /api/swarm/spawn", s.handleSwarmSpawn)
	mux.HandleFunc("POST /api/swarm/model", s.handleSwarmModel)
	mux.HandleFunc("POST /api/swarm/post", s.handleSwarmPost)
	mux.HandleFunc("POST /api/post", s.handleSwarmPost)
	mux.HandleFunc("POST /api/swarm/kill", s.handleSwarmKill)
	mux.HandleFunc("POST /api/kill", s.handleSwarmKill)
}

func (s *Server) swarmRoot() string {
	if s.ProjectDir != "" {
		return filepath.Join(s.ProjectDir, "scratch")
	}
	return "scratch"
}

func (s *Server) discoverSwarmBoards() map[string]string {
	out := make(map[string]string)
	if s.ProjectDir != "" {
		proj := filepath.Clean(s.ProjectDir)
		base := filepath.Base(proj)
		if base == ".legoloop-swarm" || base == ".swarm" {
			proj = filepath.Dir(proj)
		}
		maps.Copy(out, board.DiscoverBoards(proj, "", ""))
		maps.Copy(out, board.DiscoverBoards(filepath.Join(proj, "scratch"), "", ""))
		procName := filepath.Base(proj)
		if procName == "" || procName == "." || procName == string(filepath.Separator) {
			procName = "swarm"
		}
		if _, exists := out[procName]; !exists {
			out[procName] = filepath.Join(proj, ".legoloop-swarm")
		}
		return out
	}
	maps.Copy(out, board.DiscoverBoards("scratch", "", ""))
	if len(out) == 0 {
		out["swarm"] = ".legoloop-swarm"
	}
	return out
}

func (s *Server) openSwarmBoardForProc(proc string) (*board.SwarmBoard, error) {
	boards := s.discoverSwarmBoards()
	dir := boards[strings.TrimSpace(proc)]
	if dir == "" && len(boards) == 1 {
		for _, v := range boards {
			dir = v
		}
	}
	if dir == "" {
		return nil, fmt.Errorf("processo nao encontrado")
	}
	return board.OpenSwarmBoard(dir)
}

func (s *Server) studioBoardsPayload() ([]byte, error) {
	boards := s.discoverSwarmBoards()
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

type swarmAddReq struct {
	Process       string              `json:"process"`
	Proc          string              `json:"proc"`
	Name          string              `json:"name"`
	Role          string              `json:"role"`
	Parent        string              `json:"parent"`
	Model         string              `json:"model"`
	WorkspaceMode board.WorkspaceMode `json:"workspace_mode"`
	Command       []string            `json:"command"`
}

func (s *Server) handleSwarmAdd(w http.ResponseWriter, r *http.Request) {
	var req swarmAddReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	proc := req.Process
	if proc == "" {
		proc = req.Proc
	}
	sb, err := s.openSwarmBoardForProc(proc)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	worker, err := sb.AddWorker(board.AddWorkerRequest{
		Name: req.Name, Role: req.Role, Parent: req.Parent, Model: req.Model,
		WorkspaceMode: req.WorkspaceMode, RepoDir: s.swarmRepoDir(),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Command) > 0 {
		pid, err := sb.SpawnWorker(context.Background(), worker.Name, req.Command)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		worker.PID = pid
		worker.Status = board.WorkerWorking
	}
	writeJSON(w, http.StatusOK, worker)
}

func (s *Server) handleSwarmSpawn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process string   `json:"process"`
		Proc    string   `json:"proc"`
		Name    string   `json:"name"`
		Command []string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	proc := req.Process
	if proc == "" {
		proc = req.Proc
	}
	sb, err := s.openSwarmBoardForProc(proc)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	cmdArgs := req.Command
	if len(cmdArgs) == 0 {
		exe, err := os.Executable()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		cmdArgs = []string{exe, "swarm", "status"}
	}
	pid, err := sb.SpawnWorker(context.Background(), req.Name, cmdArgs)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"name": req.Name, "pid": pid, "status": board.WorkerWorking})
}

func (s *Server) handleSwarmModel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process string `json:"process"`
		Proc    string `json:"proc"`
		Name    string `json:"name"`
		Model   string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	proc := req.Process
	if proc == "" {
		proc = req.Proc
	}
	sb, err := s.openSwarmBoardForProc(proc)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	worker, err := sb.SetWorkerModel(req.Name, req.Model)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (s *Server) handleSwarmEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE nao suportado", http.StatusInternalServerError)
		return
	}
	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	var last string
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()
	for {
		if data, err := s.studioBoardsPayload(); err == nil && string(data) != last {
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
	var req struct {
		Process string `json:"process"`
		Proc    string `json:"proc"`
		To      string `json:"to"`
		Body    string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	proc := req.Proc
	if proc == "" {
		proc = req.Process
	}
	sb, err := s.openSwarmBoardForProc(proc)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	msg, err := sb.PostMessage("coordenador", req.To, req.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, msg)
}

func (s *Server) swarmRepoDir() string {
	if s.ProjectDir != "" {
		return s.ProjectDir
	}
	return "."
}

func (s *Server) handleSwarmKill(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process string `json:"process"`
		Proc    string `json:"proc"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	proc := req.Proc
	if proc == "" {
		proc = req.Process
	}
	sb, err := s.openSwarmBoardForProc(proc)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := sb.KillWorker(req.Name, s.swarmRepoDir()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "killed", "name": req.Name})
}
