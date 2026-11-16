package board

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"legoloop/internal/taxonomy"
)

type WorkspaceMode string

const (
	WorkspaceInherit WorkspaceMode = "inherit"
	WorkspaceShare   WorkspaceMode = "share"
	WorkspaceBranch  WorkspaceMode = "branch"
)

type WorkerStatus string

const (
	WorkerWaiting WorkerStatus = "waiting"
	WorkerWorking WorkerStatus = "working"
	WorkerDone    WorkerStatus = "done"
	WorkerStopped WorkerStatus = "stopped"
	WorkerKilled  WorkerStatus = "killed"
)

// WorkerPiece associa uma peca do catalogo ao worker.
type WorkerPiece struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Worker struct {
	Name          string        `json:"name"`
	Role          string        `json:"role"`
	Parent        string        `json:"parent,omitempty"`
	Model         string        `json:"model,omitempty"`
	Piece         string        `json:"piece,omitempty"`
	PieceKind     string        `json:"piece_kind,omitempty"`
	Pieces        []WorkerPiece `json:"pieces,omitempty"`
	Ticket        string        `json:"ticket,omitempty"`
	Step          int           `json:"step,omitempty"`
	Mode          string        `json:"mode,omitempty"`
	WorkspaceMode WorkspaceMode `json:"workspace_mode"`
	WorktreePath  string        `json:"worktree_path,omitempty"`
	BranchName    string        `json:"branch_name,omitempty"`
	RepoDir       string        `json:"repo_dir,omitempty"`
	Status        WorkerStatus  `json:"status"`
	PID           int           `json:"pid,omitempty"`
	ExitCode      *int          `json:"exit_code,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	LastActiveAt  time.Time     `json:"last_active_at"`
}

type Message struct {
	Seq       int64     `json:"seq"`
	From      string    `json:"from"`
	To        string    `json:"to,omitempty"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type SwarmStatus struct {
	Workers      []Worker `json:"workers"`
	MessageCount int      `json:"message_count"`
}

type AddWorkerRequest struct {
	Name          string        `json:"name,omitempty"`
	Role          string        `json:"role"`
	Parent        string        `json:"parent,omitempty"`
	Model         string        `json:"model,omitempty"`
	Piece         string        `json:"piece,omitempty"`
	PieceKind     string        `json:"piece_kind,omitempty"`
	Pieces        []WorkerPiece `json:"pieces,omitempty"`
	Ticket        string        `json:"ticket,omitempty"`
	Step          int           `json:"step,omitempty"`
	Mode          string        `json:"mode,omitempty"`
	WorkspaceMode WorkspaceMode `json:"workspace_mode,omitempty"`
	RepoDir       string        `json:"repo_dir,omitempty"`
}

// UpdateWorkerRequest atualiza os atributos operacionais e associacoes de um worker.
type UpdateWorkerRequest struct {
	Name      string        `json:"name"`
	Role      string        `json:"role,omitempty"`
	Model     string        `json:"model,omitempty"`
	Piece     string        `json:"piece,omitempty"`
	PieceKind string        `json:"piece_kind,omitempty"`
	Pieces    []WorkerPiece `json:"pieces,omitempty"`
	Ticket    string        `json:"ticket,omitempty"`
}

type swarmState struct {
	Workers  []Worker  `json:"workers"`
	Messages []Message `json:"messages"`
	NextSeq  int64     `json:"next_seq"`
}

type SwarmBoard struct {
	mu      sync.Mutex
	projDir string
	dir     string
	file    string
	actMu   sync.Mutex
	active  map[int]chan struct{}
}

var (
	tmpWriteSeq         uint64
	validWorkerNameRe   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)
	exactMentionTokenRe = regexp.MustCompile(`@[a-zA-Z0-9_-]+`)
)

func resolveBoardPaths(dir string) (string, string, error) {
	trimmed := strings.TrimSpace(dir)
	if trimmed == "" {
		return "", "", errors.New("board: diretorio obrigatorio")
	}
	clean := filepath.Clean(trimmed)
	base := filepath.Base(clean)
	if base == ".legoloop-swarm" || base == ".swarm" {
		proj := filepath.Dir(clean)
		return proj, filepath.Join(proj, ".legoloop-swarm"), nil
	}
	return clean, filepath.Join(clean, ".legoloop-swarm"), nil
}

// OpenSwarmBoard abre ou inicializa o quadro persistente em .legoloop-swarm/swarm.json.
func OpenSwarmBoard(dir string) (*SwarmBoard, error) {
	projDir, boardDir, err := resolveBoardPaths(dir)
	if err != nil {
		return nil, err
	}
	sb := &SwarmBoard{
		projDir: projDir,
		dir:     boardDir,
		file:    filepath.Join(boardDir, "swarm.json"),
		active:  make(map[int]chan struct{}),
	}
	err = sb.withLock(func() error {
		if sb.findReadableFile() == "" {
			return sb.saveLocked(swarmState{NextSeq: 1})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sb, nil
}

func (sb *SwarmBoard) readCandidates() []string {
	return []string{
		sb.file,
		filepath.Join(sb.dir, "swarm_board.json"),
		filepath.Join(sb.projDir, ".swarm", "swarm.json"),
		filepath.Join(sb.projDir, ".swarm", "swarm_board.json"),
		filepath.Join(sb.projDir, "swarm_board.json"),
	}
}

func (sb *SwarmBoard) findReadableFile() string {
	for _, p := range sb.readCandidates() {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

func (sb *SwarmBoard) loadLocked() (swarmState, error) {
	src := sb.findReadableFile()
	if src == "" {
		return swarmState{}, os.ErrNotExist
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return swarmState{}, err
	}
	var st swarmState
	if err := json.Unmarshal(data, &st); err != nil {
		return swarmState{}, err
	}
	if st.NextSeq == 0 {
		st.NextSeq = 1
	}
	return st, nil
}

func atomicTempPath(target string) string {
	seq := atomic.AddUint64(&tmpWriteSeq, 1)
	return fmt.Sprintf("%s.tmp.%d.%d.%d", target, os.Getpid(), time.Now().UnixNano(), seq)
}

func (sb *SwarmBoard) saveLocked(st swarmState) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := atomicTempPath(sb.file)
	if err := writeSyncExclusive(tmp, data); err != nil {
		return err
	}
	if err := os.Rename(tmp, sb.file); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func writeSyncExclusive(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
}

func validateWorkerName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if strings.EqualFold(name, "all") {
		return "", fmt.Errorf("board: nome reservado %q", name)
	}
	if !validWorkerNameRe.MatchString(name) {
		return "", fmt.Errorf("board: nome de worker invalido %q", raw)
	}
	return name, nil
}

func workerExists(existing []Worker, name string) bool {
	for _, w := range existing {
		if w.Name == name {
			return true
		}
	}
	return false
}

func resolveWorkerName(req AddWorkerRequest, existing []Worker) (string, error) {
	if strings.TrimSpace(req.Name) != "" {
		name, err := validateWorkerName(req.Name)
		if err != nil {
			return "", err
		}
		if workerExists(existing, name) {
			return "", fmt.Errorf("board: worker %q ja existe", name)
		}
		return name, nil
	}
	count := 1
	for _, w := range existing {
		if w.Role == req.Role {
			count++
		}
	}
	for {
		cand := fmt.Sprintf("%s-%d", req.Role, count)
		if !workerExists(existing, cand) {
			return cand, nil
		}
		count++
	}
}

func normalizeWorkspaceMode(mode WorkspaceMode) (WorkspaceMode, error) {
	if mode == "" {
		return WorkspaceInherit, nil
	}
	switch mode {
	case WorkspaceInherit, WorkspaceShare, WorkspaceBranch:
		return mode, nil
	default:
		return "", fmt.Errorf("board: workspace_mode invalido %q", mode)
	}
}

func safeWorktreePath(boardDir, name string) (string, error) {
	root := filepath.Join(boardDir, "worktrees")
	target := filepath.Join(root, name)
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == "." || rel == ".." || filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("board: caminho de worktree invalido para %q", name)
	}
	return target, nil
}

func createWorktree(repoDir, wtDir, branch string, detach bool) error {
	if strings.TrimSpace(repoDir) == "" {
		return errors.New("board: repo_dir obrigatorio para worktree")
	}
	if err := os.MkdirAll(filepath.Dir(wtDir), 0o755); err != nil {
		return err
	}
	args := []string{"worktree", "add", "-b", branch, wtDir, "HEAD"}
	if detach {
		args = []string{"worktree", "add", "--detach", wtDir, "HEAD"}
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree add: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func provisionWorkspace(boardDir string, mode WorkspaceMode, name, repoDir string, now time.Time) (string, string, error) {
	if mode == WorkspaceInherit {
		return "", "", nil
	}
	wtPath, err := safeWorktreePath(boardDir, name)
	if err != nil {
		return "", "", err
	}
	if mode == WorkspaceShare {
		return wtPath, "", createWorktree(repoDir, wtPath, "", true)
	}
	branch := fmt.Sprintf("swarm/%s-%d", name, now.UnixNano())
	return wtPath, branch, createWorktree(repoDir, wtPath, branch, false)
}

// AddWorker registra um worker validado pela Matriz de Papeis e provisiona seu workspace.
func (sb *SwarmBoard) AddWorker(req AddWorkerRequest) (Worker, error) {
	if err := taxonomy.ValidarPapel(req.Role); err != nil {
		return Worker{}, err
	}
	mode, err := normalizeWorkspaceMode(req.WorkspaceMode)
	if err != nil {
		return Worker{}, err
	}
	var created Worker
	err = sb.withLock(func() error {
		return sb.addWorkerLocked(req, mode, &created)
	})
	return created, err
}

func (sb *SwarmBoard) addWorkerLocked(req AddWorkerRequest, mode WorkspaceMode, out *Worker) error {
	st, err := sb.loadLocked()
	if err != nil {
		return err
	}
	name, err := resolveWorkerName(req, st.Workers)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	wtPath, branch, err := provisionWorkspace(sb.dir, mode, name, req.RepoDir, now)
	if err != nil {
		return err
	}
	repoDir := strings.TrimSpace(req.RepoDir)
	if repoDir != "" {
		if abs, err := filepath.Abs(repoDir); err == nil {
			repoDir = abs
		}
	}
	pieces, firstPiece, firstKind := normalizeWorkerPieces(req.Pieces, req.Piece, req.PieceKind)
	w := Worker{
		Name: name, Role: req.Role, Parent: req.Parent, Model: req.Model,
		Piece: firstPiece, PieceKind: firstKind, Pieces: pieces,
		Ticket: strings.TrimSpace(req.Ticket), Step: req.Step, Mode: strings.TrimSpace(req.Mode),
		WorkspaceMode: mode, WorktreePath: wtPath, BranchName: branch,
		RepoDir: repoDir, Status: WorkerWaiting,
		CreatedAt: now, LastActiveAt: now,
	}
	st.Workers = append(st.Workers, w)
	if err := sb.saveLocked(st); err != nil {
		return err
	}
	*out = w
	return nil
}

// PostMessage publica uma mensagem sequencial no quadro.
func (sb *SwarmBoard) PostMessage(from, to, body string) (Message, error) {
	var msg Message
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		msg = Message{Seq: st.NextSeq, From: from, To: to, Body: body, CreatedAt: time.Now().UTC()}
		st.NextSeq++
		st.Messages = append(st.Messages, msg)
		applyMessageTargets(st.Workers, from, to, body, msg.CreatedAt)
		return sb.saveLocked(st)
	})
	return msg, err
}

// FinishSimulation publica a conclusao da simulacao e marca workers ativos como concluidos.
func (sb *SwarmBoard) FinishSimulation(body string) (Message, error) {
	var msg Message
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		msg = Message{Seq: st.NextSeq, From: "coordenador", To: "registro", Body: body, CreatedAt: now}
		st.NextSeq++
		st.Messages = append(st.Messages, msg)
		for i := range st.Workers {
			if st.Workers[i].Status != WorkerKilled {
				st.Workers[i].Status = WorkerDone
				st.Workers[i].LastActiveAt = now
			}
		}
		return sb.saveLocked(st)
	})
	return msg, err
}

func applyMessageTargets(workers []Worker, from, to, body string, at time.Time) {
	targets := collectMentionTargets(to, body)
	other := hasOtherTarget(targets, from)
	for i := range workers {
		if workers[i].Status == WorkerKilled {
			continue
		}
		if workers[i].Name == from {
			workers[i].LastActiveAt = at
			if other {
				workers[i].Status = WorkerWaiting
			}
		} else if targets["all"] || targets[workers[i].Name] {
			workers[i].Status = WorkerWorking
			workers[i].LastActiveAt = at
		}
	}
}

func collectMentionTargets(to, body string) map[string]bool {
	targets := make(map[string]bool)
	if t := strings.TrimPrefix(strings.TrimSpace(to), "@"); t != "" {
		targets[t] = true
	}
	for _, m := range extractMentions(body) {
		targets[m] = true
	}
	return targets
}

func hasOtherTarget(targets map[string]bool, from string) bool {
	for t := range targets {
		if t != "" && t != from {
			return true
		}
	}
	return false
}

// ReadMessages retorna mensagens com Seq > since ou pertencentes a seqs.
func (sb *SwarmBoard) ReadMessages(since int64, seqs []int64) ([]Message, error) {
	var out []Message
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		out = filterMessages(st.Messages, since, seqs)
		return nil
	})
	return out, err
}

func filterMessages(msgs []Message, since int64, seqs []int64) []Message {
	seqSet := make(map[int64]bool, len(seqs))
	for _, s := range seqs {
		seqSet[s] = true
	}
	var out []Message
	for _, m := range msgs {
		if len(seqSet) > 0 && seqSet[m.Seq] {
			out = append(out, m)
		} else if len(seqSet) == 0 && m.Seq > since {
			out = append(out, m)
		}
	}
	return out
}

func isIdentChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '-'
}

func extractMentions(body string) []string {
	idxs := exactMentionTokenRe.FindAllStringIndex(body, -1)
	var out []string
	for _, pair := range idxs {
		start, end := pair[0], pair[1]
		if start > 0 && isIdentChar(body[start-1]) {
			continue
		}
		if end < len(body) && isIdentChar(body[end]) {
			continue
		}
		token := body[start+1 : end]
		if token == "all" || validWorkerNameRe.MatchString(token) {
			out = append(out, token)
		}
	}
	return out
}

func matchesTarget(m Message, target string) bool {
	want := strings.TrimPrefix(strings.TrimSpace(target), "@")
	if want == "" {
		return true
	}
	to := strings.TrimPrefix(strings.TrimSpace(m.To), "@")
	if to == "all" || to == want {
		return true
	}
	for _, mention := range extractMentions(m.Body) {
		if mention == "all" || mention == want {
			return true
		}
	}
	return false
}

// WaitMessages aguarda passivamente mensagens direcionadas a target ou @all apos since.
func (sb *SwarmBoard) WaitMessages(ctx context.Context, target string, since int64, poll time.Duration) ([]Message, error) {
	if poll <= 0 {
		poll = 100 * time.Millisecond
	}
	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
		matched, err := sb.pollMatchedMessages(target, since)
		if err != nil || len(matched) > 0 {
			return matched, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (sb *SwarmBoard) pollMatchedMessages(target string, since int64) ([]Message, error) {
	msgs, err := sb.ReadMessages(since, nil)
	if err != nil {
		return nil, err
	}
	var matched []Message
	for _, m := range msgs {
		if matchesTarget(m, target) {
			matched = append(matched, m)
		}
	}
	return matched, nil
}

// Status retorna o estado atual dos workers e contagem de mensagens.
func (sb *SwarmBoard) Status() (SwarmStatus, error) {
	var out SwarmStatus
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		if reconcileDeadPIDs(st.Workers) {
			if err := sb.saveLocked(st); err != nil {
				return err
			}
		}
		workers := st.Workers
		if workers == nil {
			workers = []Worker{}
		}
		out = SwarmStatus{Workers: workers, MessageCount: len(st.Messages)}
		return nil
	})
	return out, err
}

// reconcileDeadPIDs marca como stopped os workers cujo processo morreu sem que o dono do spawn
// (outro processo do CLI) pudesse registrar a saida. Retorna true se alterou algo.
func reconcileDeadPIDs(workers []Worker) bool {
	changed := false
	now := time.Now().UTC()
	for i := range workers {
		w := &workers[i]
		if w.PID <= 0 || w.Status == WorkerKilled || isOSProcessAlive(w.PID) {
			continue
		}
		w.Status = WorkerStopped
		w.PID = 0
		w.LastActiveAt = now
		changed = true
	}
	return changed
}

// TouchWorker atualiza status e timestamp de atividade de um worker.
func (sb *SwarmBoard) TouchWorker(name string, status WorkerStatus, at time.Time) error {
	return sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for i := range st.Workers {
			if st.Workers[i].Name == name {
				st.Workers[i].Status = status
				st.Workers[i].LastActiveAt = at.UTC()
				return sb.saveLocked(st)
			}
		}
		return fmt.Errorf("board: worker %q nao encontrado", name)
	})
}

// CheckStall retorna workers em execucao sem atividade dentro de timeout.
func (sb *SwarmBoard) CheckStall(timeout time.Duration) ([]Worker, error) {
	stalled := []Worker{}
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		cutoff := time.Now().UTC().Add(-timeout)
		for _, w := range st.Workers {
			if w.Status == WorkerWorking && w.LastActiveAt.Before(cutoff) {
				stalled = append(stalled, w)
			}
		}
		return nil
	})
	return stalled, err
}

// KillWorker encerra o PID real do worker no SO, remove seu worktree e marca killed.
func (sb *SwarmBoard) KillWorker(name, repoDir string) error {
	w, err := sb.lookupWorker(name)
	if err != nil {
		return err
	}
	if w.PID > 0 {
		_ = killOSProcess(w.PID)
		sb.waitActiveExit(w.PID, 3*time.Second)
	}
	effectiveRepo := strings.TrimSpace(repoDir)
	if w.RepoDir != "" && (effectiveRepo == "" || effectiveRepo == ".") {
		effectiveRepo = w.RepoDir
	}
	if w.WorktreePath != "" {
		removeWorktree(effectiveRepo, w.WorktreePath)
	}
	return sb.markWorkerKilled(name)
}

func (sb *SwarmBoard) lookupWorker(name string) (Worker, error) {
	var found Worker
	err := sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for _, w := range st.Workers {
			if w.Name == name {
				found = w
				return nil
			}
		}
		return fmt.Errorf("board: worker %q nao encontrado", name)
	})
	return found, err
}

func (sb *SwarmBoard) markWorkerKilled(name string) error {
	return sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for i := range st.Workers {
			if st.Workers[i].Name == name {
				st.Workers[i].Status = WorkerKilled
				st.Workers[i].PID = 0
				st.Workers[i].LastActiveAt = time.Now().UTC()
				return sb.saveLocked(st)
			}
		}
		return fmt.Errorf("board: worker %q nao encontrado", name)
	})
}

// UpdateWorker atualiza papel, modelo, peça associada e ticket de um worker existente.
func (sb *SwarmBoard) UpdateWorker(req UpdateWorkerRequest) (Worker, error) {
	if req.Role != "" {
		if err := taxonomy.ValidarPapel(req.Role); err != nil {
			return Worker{}, err
		}
	}
	var updated Worker
	err := sb.withLock(func() error {
		return sb.updateWorkerLocked(req, &updated)
	})
	return updated, err
}

func (sb *SwarmBoard) updateWorkerLocked(req UpdateWorkerRequest, out *Worker) error {
	st, err := sb.loadLocked()
	if err != nil {
		return err
	}
	for i := range st.Workers {
		if st.Workers[i].Name == req.Name {
			applyWorkerPatch(&st.Workers[i], req)
			st.Workers[i].LastActiveAt = time.Now().UTC()
			*out = st.Workers[i]
			return sb.saveLocked(st)
		}
	}
	return fmt.Errorf("board: worker %q nao encontrado", req.Name)
}

func applyWorkerPatch(w *Worker, req UpdateWorkerRequest) {
	if strings.TrimSpace(req.Role) != "" {
		w.Role = strings.TrimSpace(req.Role)
	}
	if strings.TrimSpace(req.Model) != "" {
		w.Model = strings.TrimSpace(req.Model)
	}
	pieces, firstPiece, firstKind := normalizeWorkerPieces(req.Pieces, req.Piece, req.PieceKind)
	w.Piece = firstPiece
	w.PieceKind = firstKind
	w.Pieces = pieces
	w.Ticket = strings.TrimSpace(req.Ticket)
}

func normalizeWorkerPieces(in []WorkerPiece, single, kind string) ([]WorkerPiece, string, string) {
	var out []WorkerPiece
	seen := make(map[string]bool)
	for _, p := range in {
		name := strings.TrimSpace(p.Name)
		k := strings.TrimSpace(p.Kind)
		if name == "" {
			continue
		}
		key := strings.ToLower(k + ":" + name)
		if !seen[key] {
			seen[key] = true
			out = append(out, WorkerPiece{Kind: k, Name: name})
		}
	}
	if len(out) == 0 && strings.TrimSpace(single) != "" {
		out = append(out, WorkerPiece{Kind: strings.TrimSpace(kind), Name: strings.TrimSpace(single)})
	}
	if len(out) == 0 {
		return nil, "", ""
	}
	return out, out[0].Name, out[0].Kind
}

// RemoveWorker encerra PID/worktree do worker e o remove definitivamente do quadro.
func (sb *SwarmBoard) RemoveWorker(name, repoDir string) error {
	w, err := sb.lookupWorker(name)
	if err != nil {
		return err
	}
	sb.cleanupWorkerResources(w, repoDir)
	return sb.deleteWorkerRecord(name)
}

func (sb *SwarmBoard) cleanupWorkerResources(w Worker, repoDir string) {
	if w.PID > 0 {
		_ = killOSProcess(w.PID)
		sb.waitActiveExit(w.PID, 3*time.Second)
	}
	effectiveRepo := strings.TrimSpace(repoDir)
	if w.RepoDir != "" && (effectiveRepo == "" || effectiveRepo == ".") {
		effectiveRepo = w.RepoDir
	}
	if w.WorktreePath != "" {
		removeWorktree(effectiveRepo, w.WorktreePath)
	}
}

func (sb *SwarmBoard) deleteWorkerRecord(name string) error {
	return sb.withLock(func() error {
		st, err := sb.loadLocked()
		if err != nil {
			return err
		}
		for i := range st.Workers {
			if st.Workers[i].Name == name {
				st.Workers = append(st.Workers[:i], st.Workers[i+1:]...)
				return sb.saveLocked(st)
			}
		}
		return fmt.Errorf("board: worker %q nao encontrado", name)
	})
}

func removeWorktree(repoDir, wtPath string) {
	if strings.TrimSpace(repoDir) != "" {
		cmd := exec.Command("git", "worktree", "remove", "--force", wtPath)
		cmd.Dir = repoDir
		_ = cmd.Run()
		prune := exec.Command("git", "worktree", "prune")
		prune.Dir = repoDir
		_ = prune.Run()
	}
	_ = os.RemoveAll(wtPath)
}

// LiveView devolve uma copia dos workers com o PID morto marcado como stopped, sem gravar nada. Serve a leituras
// somente leitura (MCP swarm_status) que precisam mostrar o mesmo estado que Status() grava.
func LiveView(workers []Worker) []Worker {
	out := append([]Worker{}, workers...)
	reconcileDeadPIDs(out)
	return out
}
