package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/orchestrix"
	"legoloop/internal/runner"
)

type runDAGRequest struct {
	DAG          json.RawMessage `json:"dag"`
	AllowScripts bool            `json:"allow_scripts"`
}

// handleRunDAG executa o pipeline via SSE, transmitindo progresso por nó em tempo real (R-RUN-1).
// Requer consentimento prévio para execução de scripts (R-RUN-2).
func (s *Server) handleRunDAG(w http.ResponseWriter, r *http.Request) {
	var req runDAGRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "corpo inválido: " + err.Error()})
		return
	}

	d, err := dag.Parse(req.DAG)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "dag inválido: " + err.Error()})
		return
	}

	// R-RUN-2: nós script exigem consentimento explícito
	for _, n := range d.Nodes {
		if n.Type == dag.TypeScript && !req.AllowScripts {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": fmt.Sprintf("execução de nós script (%s) requer consentimento prévio do usuário", n.ID),
			})
			return
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "streaming sse não suportado"})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	sendSSE := func(eventType string, payload map[string]any) {
		payload["type"] = eventType
		data, _ := json.Marshal(payload)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
		flusher.Flush()
	}

	sendSSE("dag_start", map[string]any{
		"dag":         d.Name,
		"nodes_total": len(d.Nodes),
		"time":        time.Now().Format(time.RFC3339),
	})

	exec := &orchestrix.Orchestrix{
		Inventory:    s.inventory(r.Context()),
		AllowScripts: req.AllowScripts,
		LoopStore:    s.Loops,
	}

	opts := runner.Options{
		Workdir:  s.ProjectDir,
		Executor: exec,
		OnNodeStart: func(nodeID string) {
			sendSSE("node_start", map[string]any{
				"id":   nodeID,
				"time": time.Now().Format(time.RFC3339),
			})
		},
		OnNodeFinish: func(nr runner.NodeReport) {
			sendSSE("node_finish", map[string]any{
				"id":          nr.ID,
				"status":      nr.Status,
				"tokens":      nr.Tokens,
				"model":       nr.Model,
				"duration_ms": nr.Duration.Milliseconds(),
				"notes":       nr.Notes,
				"error":       nr.Error,
			})
		},
	}

	rep, err := runner.Run(r.Context(), d, opts)
	if err != nil {
		sendSSE("dag_error", map[string]any{"error": err.Error()})
		return
	}

	sendSSE("dag_finish", map[string]any{
		"ok":     rep.OK,
		"tokens": rep.Tokens,
	})
}
