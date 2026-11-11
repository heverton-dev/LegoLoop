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

// handleRunDAG executa o Processo via SSE, transmitindo progresso por nó em tempo real (R-RUN-1).
// Requer consentimento prévio para execução de scripts (R-RUN-2).
func (s *Server) handleRunDAG(w http.ResponseWriter, r *http.Request) {
	req, d, ok := parseRunPayload(w, r)
	if !ok || !validateRunContext(w, d, req.AllowScripts) {
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "streaming sse não suportado"})
		return
	}
	startSSE(w, flusher)
	s.dispatchRunExecution(r, d, req.AllowScripts, newSSESender(w, flusher))
}

// parseRunPayload decodifica o corpo e o DAG; responde 400 e devolve ok=false se algum for inválido.
func parseRunPayload(w http.ResponseWriter, r *http.Request) (runDAGRequest, *dag.DAG, bool) {
	var req runDAGRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "corpo inválido: " + err.Error()})
		return req, nil, false
	}
	d, err := dag.Parse(req.DAG)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "dag inválido: " + err.Error()})
		return req, nil, false
	}
	return req, d, true
}

// validateRunContext aplica R-RUN-2: nós script exigem consentimento explícito.
func validateRunContext(w http.ResponseWriter, d *dag.DAG, allowScripts bool) bool {
	for _, n := range d.Nodes {
		if n.Type == dag.TypeScript && !allowScripts {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": fmt.Sprintf("execução de nós script (%s) requer consentimento prévio do usuário", n.ID),
			})
			return false
		}
	}
	return true
}

func startSSE(w http.ResponseWriter, flusher http.Flusher) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
}

type sseSender func(eventType string, payload map[string]any)

func newSSESender(w http.ResponseWriter, flusher http.Flusher) sseSender {
	return func(eventType string, payload map[string]any) {
		payload["type"] = eventType
		data, _ := json.Marshal(payload)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
		flusher.Flush()
	}
}

// dispatchRunExecution roda o DAG e transmite dag_start, node_* e dag_finish/dag_error.
func (s *Server) dispatchRunExecution(r *http.Request, d *dag.DAG, allowScripts bool, send sseSender) {
	send("dag_start", map[string]any{
		"dag":         d.Name,
		"nodes_total": len(d.Nodes),
		"time":        time.Now().Format(time.RFC3339),
	})
	exec := &orchestrix.Orchestrix{
		Inventory:    s.inventory(r.Context()),
		AllowScripts: allowScripts,
		LoopStore:    s.Processos,
	}
	start := time.Now()
	rep, err := runner.Run(r.Context(), d, s.runOptions(exec, send))
	elapsed := time.Since(start)
	if err != nil {
		send("dag_error", map[string]any{"error": err.Error(), "duration_ms": elapsed.Milliseconds()})
		return
	}
	// Estimativa média de custo base de tokens LLM (~$0.002/1k input, ~$0.008/1k output)
	costEst := (float64(rep.Tokens.Input)*0.002 + float64(rep.Tokens.Output)*0.008) / 1000.0
	send("dag_finish", map[string]any{
		"ok":          rep.OK,
		"tokens":      rep.Tokens,
		"duration_ms": elapsed.Milliseconds(),
		"cost_est":    costEst,
	})
}

func (s *Server) runOptions(exec *orchestrix.Orchestrix, send sseSender) runner.Options {
	return runner.Options{
		Workdir:  s.ProjectDir,
		Executor: exec,
		OnNodeStart: func(nodeID string) {
			send("node_start", map[string]any{"id": nodeID, "time": time.Now().Format(time.RFC3339)})
		},
		OnNodeFinish: func(nr runner.NodeReport) {
			send("node_finish", map[string]any{
				"id": nr.ID, "status": nr.Status, "tokens": nr.Tokens, "model": nr.Model,
				"duration_ms": nr.Duration.Milliseconds(), "notes": nr.Notes, "error": nr.Error,
			})
		},
	}
}
