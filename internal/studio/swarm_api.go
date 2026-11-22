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

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/taxonomy"
)

func (s *Server) registerSwarmRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swarm", handlePage("swarm.html"))
	mux.HandleFunc("GET /api/swarm/meta", s.handleSwarmMeta)
	mux.HandleFunc("GET /api/swarm/events", s.handleSwarmEvents)
	mux.HandleFunc("GET /api/events", s.handleSwarmEvents)
	mux.HandleFunc("POST /api/swarm/add", s.handleSwarmAdd)
	mux.HandleFunc("POST /api/swarm/update", s.handleSwarmUpdate)
	mux.HandleFunc("POST /api/swarm/remove", s.handleSwarmRemove)
	mux.HandleFunc("POST /api/swarm/simulate", s.handleSwarmSimulate)
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
	Piece         string              `json:"piece"`
	PieceKind     string              `json:"piece_kind"`
	Pieces        []board.WorkerPiece `json:"pieces"`
	Ticket        string              `json:"ticket"`
	WorkspaceMode board.WorkspaceMode `json:"workspace_mode"`
	Command       []string            `json:"command"`
}

func (s *Server) handleSwarmMeta(w http.ResponseWriter, r *http.Request) {
	models := s.inventory(r.Context()).Models()
	if models == nil {
		models = nil
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"roles":      taxonomy.ListarPapeis(),
		"guilds":     taxonomy.ListarGuildas(),
		"models":     models,
		"workspaces": swarmWorkspaceOptions(),
	})
}

func swarmWorkspaceOptions() []map[string]string {
	return []map[string]string{
		{"id": string(board.WorkspaceInherit), "label": "Herdado (mesmo workspace)"},
		{"id": string(board.WorkspaceShare), "label": "Compartilhado (worktree proprio)"},
		{"id": string(board.WorkspaceBranch), "label": "Nova Branch (clone isolado)"},
	}
}

func inferSwarmTicket(role string, pieces []board.WorkerPiece, single string, mode board.WorkspaceMode) string {
	var names []string
	for _, p := range pieces {
		if strings.TrimSpace(p.Name) != "" {
			names = append(names, strings.TrimSpace(p.Name))
		}
	}
	if len(names) == 0 && strings.TrimSpace(single) != "" {
		names = append(names, strings.TrimSpace(single))
	}
	wsLabel := string(mode)
	if wsLabel == "" {
		wsLabel = "inherit"
	}
	if len(names) > 0 {
		return fmt.Sprintf("Executar %s como %s (%s)", strings.Join(names, ", "), firstNonEmpty(role, "worker"), wsLabel)
	}
	if strings.TrimSpace(role) != "" {
		return fmt.Sprintf("Executar papel %s em workspace %s", role, wsLabel)
	}
	return ""
}

func (s *Server) handleSwarmAdd(w http.ResponseWriter, r *http.Request) {
	var req swarmAddReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sb, err := s.openSwarmBoardForProc(firstNonEmpty(req.Process, req.Proc))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	ticket := strings.TrimSpace(req.Ticket)
	if ticket == "" {
		ticket = inferSwarmTicket(req.Role, req.Pieces, req.Piece, req.WorkspaceMode)
	}
	worker, err := sb.AddWorker(board.AddWorkerRequest{
		Name: req.Name, Role: req.Role, Parent: req.Parent, Model: req.Model,
		Piece: req.Piece, PieceKind: req.PieceKind, Pieces: req.Pieces,
		Ticket: ticket, WorkspaceMode: req.WorkspaceMode, RepoDir: s.swarmRepoDir(),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Command) > 0 {
		if err := spawnAddedWorker(sb, &worker, req.Command); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, worker)
}

func spawnAddedWorker(sb *board.SwarmBoard, w *board.Worker, cmd []string) error {
	pid, err := sb.SpawnWorker(context.Background(), w.Name, cmd)
	if err != nil {
		return err
	}
	w.PID = pid
	w.Status = board.WorkerWorking
	return nil
}

func (s *Server) handleSwarmUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process   string              `json:"process"`
		Proc      string              `json:"proc"`
		Name      string              `json:"name"`
		Role      string              `json:"role"`
		Model     string              `json:"model"`
		Piece     string              `json:"piece"`
		PieceKind string              `json:"piece_kind"`
		Pieces    []board.WorkerPiece `json:"pieces"`
		Ticket    string              `json:"ticket"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sb, err := s.openSwarmBoardForProc(firstNonEmpty(req.Process, req.Proc))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	ticket := strings.TrimSpace(req.Ticket)
	if ticket == "" && (len(req.Pieces) > 0 || strings.TrimSpace(req.Piece) != "") {
		ticket = inferSwarmTicket(req.Role, req.Pieces, req.Piece, "")
	}
	worker, err := sb.UpdateWorker(board.UpdateWorkerRequest{
		Name: req.Name, Role: req.Role, Model: req.Model,
		Piece: req.Piece, PieceKind: req.PieceKind, Pieces: req.Pieces, Ticket: ticket,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (s *Server) handleSwarmRemove(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Process string `json:"process"`
		Proc    string `json:"proc"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sb, err := s.openSwarmBoardForProc(firstNonEmpty(req.Process, req.Proc))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := sb.RemoveWorker(req.Name, s.swarmRepoDir()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed", "name": req.Name})
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
