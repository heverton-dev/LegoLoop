package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/board"
	"legoloop/internal/caminhos"
	"legoloop/internal/taxonomy"
)

type simNodeInput struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Kind      string `json:"kind"`
	Role      string `json:"role"`
	Block     string `json:"block"`
	Piece     string `json:"piece"`
	Execution string `json:"execution"`
	Artifact  string `json:"artifact"`
	Gate      string `json:"gate"`
}

type simEdgeInput struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type swarmSimulateReq struct {
	Name        string         `json:"name"`
	Nodes       []simNodeInput `json:"nodes"`
	Edges       []simEdgeInput `json:"edges"`
	StepDelayMs int            `json:"step_delay_ms"`
}

type simStepSummary struct {
	Step  int      `json:"step"`
	Mode  string   `json:"mode"`
	Nodes []string `json:"nodes"`
}

var simSlugRe = regexp.MustCompile(`[^a-z0-9-]+`)

func normalizeSimSlug(name string) string {
	low := strings.ToLower(strings.TrimSpace(name))
	low = strings.ReplaceAll(low, " ", "-")
	clean := strings.Trim(simSlugRe.ReplaceAllString(low, "-"), "-")
	if clean == "" {
		clean = "processo"
	}
	return "sim-" + clean
}

func sanitizeWorkerID(raw string, used map[string]bool) string {
	clean := strings.ToLower(strings.TrimSpace(raw))
	clean = simSlugRe.ReplaceAllString(clean, "-")
	clean = strings.Trim(clean, "-")
	if clean == "" || !unicodeLetterOrDigitStart(clean) {
		clean = "no-1"
	}
	if len(clean) > 48 {
		clean = clean[:48]
	}
	base := clean
	for i := 2; used[clean]; i++ {
		clean = fmt.Sprintf("%s-%d", base, i)
	}
	used[clean] = true
	return clean
}

func unicodeLetterOrDigitStart(s string) bool {
	if s == "" {
		return false
	}
	c := s[0]
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

func resolveSimRole(n simNodeInput) string {
	role := strings.TrimSpace(n.Role)
	if role != "" && taxonomy.ValidarPapel(role) == nil {
		return role
	}
	kind := strings.ToLower(strings.TrimSpace(n.Kind))
	switch {
	case strings.HasPrefix(kind, "trigger_"), kind == string(assets.KindProcesso), kind == caminhos.TipoSubprocessoAntigo, kind == "subprocesso":
		return "orquestrador"
	case kind == "gate":
		return "auditor-qualidade"
	case kind == "script", kind == "hook":
		return "engenheiro-backend"
	default:
		return "arquiteto"
	}
}

func buildTopoLevels(nodes []simNodeInput, edges []simEdgeInput) ([][]simNodeInput, map[string]simNodeInput) {
	byID := make(map[string]simNodeInput, len(nodes))
	inDeg := make(map[string]int, len(nodes))
	adj := make(map[string][]string, len(nodes))
	for _, n := range nodes {
		byID[n.ID] = n
		inDeg[n.ID] = 0
	}
	for _, e := range edges {
		if _, okA := byID[e.From]; !okA {
			continue
		}
		if _, okB := byID[e.To]; !okB || e.From == e.To {
			continue
		}
		adj[e.From] = append(adj[e.From], e.To)
		inDeg[e.To]++
	}
	return collectLevels(nodes, byID, inDeg, adj), byID
}

func collectLevels(nodes []simNodeInput, byID map[string]simNodeInput, inDeg map[string]int, adj map[string][]string) [][]simNodeInput {
	var queue []string
	seen := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		if inDeg[n.ID] == 0 {
			queue = append(queue, n.ID)
			seen[n.ID] = true
		}
	}
	var levels [][]simNodeInput
	for len(queue) > 0 {
		var lvl []simNodeInput
		var next []string
		for _, id := range queue {
			lvl = append(lvl, byID[id])
			for _, dst := range adj[id] {
				inDeg[dst]--
				if inDeg[dst] <= 0 && !seen[dst] {
					seen[dst] = true
					next = append(next, dst)
				}
			}
		}
		levels = append(levels, lvl)
		queue = next
	}
	return appendRemainingNodes(levels, nodes, seen)
}

func appendRemainingNodes(levels [][]simNodeInput, nodes []simNodeInput, seen map[string]bool) [][]simNodeInput {
	var rest []simNodeInput
	for _, n := range nodes {
		if !seen[n.ID] {
			rest = append(rest, n)
		}
	}
	if len(rest) > 0 {
		levels = append(levels, rest)
	}
	return levels
}

func (s *Server) handleSwarmSimulate(w http.ResponseWriter, r *http.Request) {
	var req swarmSimulateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Nodes) == 0 {
		writeError(w, http.StatusBadRequest, "adicione ao menos uma peça antes de simular")
		return
	}
	proc := normalizeSimSlug(req.Name)
	sb, err := s.resetSimulationBoard(proc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	levels, _ := buildTopoLevels(req.Nodes, req.Edges)
	steps, workerMap := registerSimulationWorkers(sb, levels)
	delay := resolveSimDelay(req.StepDelayMs)
	go runSimulationSequence(sb, req.Name, levels, workerMap, delay)
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"process": proc,
		"url":     "/swarm?proc=" + proc,
		"steps":   steps,
	})
}

func (s *Server) resetSimulationBoard(proc string) (*board.SwarmBoard, error) {
	dir := filepath.Join(s.swarmRoot(), proc)
	_ = os.RemoveAll(filepath.Join(dir, ".legoloop-swarm"))
	return board.OpenSwarmBoard(dir)
}

func registerSimulationWorkers(sb *board.SwarmBoard, levels [][]simNodeInput) ([]simStepSummary, map[string]string) {
	used := map[string]bool{"coordenador": true}
	workerMap := make(map[string]string)
	steps := make([]simStepSummary, 0, len(levels))
	for i, lvl := range levels {
		mode := "sync"
		if len(lvl) > 1 {
			mode = "async"
		}
		var nodeIDs []string
		for _, n := range lvl {
			wName := sanitizeWorkerID(n.ID, used)
			workerMap[n.ID] = wName
			nodeIDs = append(nodeIDs, n.ID)
			piece := firstNonEmpty(n.Piece, n.Block, n.Label, n.ID)
			ticketDesc := fmt.Sprintf("Passo %d (%s): %s -> artefato %s", i+1, modeLabelPT(mode), firstNonEmpty(n.Label, n.ID), firstNonEmpty(n.Artifact, "saida"))
			_, _ = sb.AddWorker(board.AddWorkerRequest{
				Name: wName, Role: resolveSimRole(n), WorkspaceMode: board.WorkspaceInherit,
				Model: "simulador-dag", Piece: piece, PieceKind: firstNonEmpty(n.Kind, "skill"),
				Ticket: ticketDesc, Step: i + 1, Mode: mode,
			})
		}
		steps = append(steps, simStepSummary{Step: i + 1, Mode: mode, Nodes: nodeIDs})
	}
	return steps, workerMap
}

func resolveSimDelay(ms int) time.Duration {
	if ms >= 5 && ms <= 5000 {
		return time.Duration(ms) * time.Millisecond
	}
	return 320 * time.Millisecond
}

func modeLabelPT(mode string) string {
	if mode == "async" {
		return "paralelo / assíncrono"
	}
	return "sequencial / síncrono"
}

func runSimulationSequence(sb *board.SwarmBoard, procName string, levels [][]simNodeInput, workerMap map[string]string, delay time.Duration) {
	_, _ = sb.PostMessage("coordenador", "registro", fmt.Sprintf("Simulação do processo %q iniciada (%d níveis topológicos).", firstNonEmpty(procName, "dag"), len(levels)))
	for i, lvl := range levels {
		mode := "sync"
		if len(lvl) > 1 {
			mode = "async"
		}
		simulateStepLevel(sb, i+1, mode, lvl, workerMap, delay)
	}
	_, _ = sb.FinishSimulation(fmt.Sprintf("Simulação concluída com sucesso em %d passos.", len(levels)))
}

func simulateStepLevel(sb *board.SwarmBoard, stepNum int, mode string, lvl []simNodeInput, workerMap map[string]string, delay time.Duration) {
	for _, n := range lvl {
		wName := workerMap[n.ID]
		_ = sb.TouchWorker(wName, board.WorkerWorking, time.Now().UTC())
		_, _ = sb.PostMessage("coordenador", "@"+wName, fmt.Sprintf("[Passo %d · %s] Executar %q (%s).", stepNum, modeLabelPT(mode), firstNonEmpty(n.Label, n.ID), firstNonEmpty(n.Kind, "peca")))
	}
	time.Sleep(delay)
	for _, n := range lvl {
		wName := workerMap[n.ID]
		art := firstNonEmpty(n.Artifact, "concluido")
		_, _ = sb.PostMessage(wName, "@coordenador", fmt.Sprintf("[OK Passo %d] Artefato entregue: %s (gate: %s).", stepNum, art, firstNonEmpty(n.Gate, "valido")))
		_ = sb.TouchWorker(wName, board.WorkerWaiting, time.Now().UTC())
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
