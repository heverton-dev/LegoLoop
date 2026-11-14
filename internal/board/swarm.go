package board

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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
	WorkerStopped WorkerStatus = "stopped"
	WorkerKilled  WorkerStatus = "killed"
)

type Worker struct {
	Name          string        `json:"name"`
	Role          string        `json:"role"`
	Model         string        `json:"model,omitempty"`
	WorkspaceMode WorkspaceMode `json:"workspace_mode"`
	WorktreePath  string        `json:"worktree_path,omitempty"`
	BranchName    string        `json:"branch_name,omitempty"`
	Status        WorkerStatus  `json:"status"`
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
	Model         string        `json:"model,omitempty"`
	WorkspaceMode WorkspaceMode `json:"workspace_mode,omitempty"`
	RepoDir       string        `json:"repo_dir,omitempty"`
}

type swarmState struct {
	Workers  []Worker  `json:"workers"`
	Messages []Message `json:"messages"`
	NextSeq  int64     `json:"next_seq"`
}

type SwarmBoard struct {
	mu   sync.Mutex
	dir  string
	file string
}

// OpenSwarmBoard abre ou inicializa o quadro persistente do Swarm em dir.
func OpenSwarmBoard(dir string) (*SwarmBoard, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("board: diretorio obrigatorio")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	sb := &SwarmBoard{dir: dir, file: filepath.Join(dir, "swarm_board.json")}
	if _, err := os.Stat(sb.file); errors.Is(err, os.ErrNotExist) {
		if err := sb.saveLocked(swarmState{NextSeq: 1}); err != nil {
			return nil, err
		}
	}
	return sb, nil
}

func (sb *SwarmBoard) loadLocked() (swarmState, error) {
	data, err := os.ReadFile(sb.file)
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

func (sb *SwarmBoard) saveLocked(st swarmState) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := sb.file + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, sb.file)
}

func resolveWorkerName(req AddWorkerRequest, existing []Worker) string {
	if strings.TrimSpace(req.Name) != "" {
		return strings.TrimSpace(req.Name)
	}
	count := 1
	for _, w := range existing {
		if w.Role == req.Role {
			count++
		}
	}
	return fmt.Sprintf("%s-%d", req.Role, count)
}

func createWorktree(repoDir, wtDir, branch string) error {
	if strings.TrimSpace(repoDir) == "" {
		return errors.New("board: repo_dir obrigatorio para worktree")
	}
	if err := os.MkdirAll(filepath.Dir(wtDir), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("git", "worktree", "add", "-b", branch, wtDir)
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree add: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// AddWorker registra um worker validado pela Matriz de Papeis e cria worktree se necessario.
func (sb *SwarmBoard) AddWorker(req AddWorkerRequest) (Worker, error) {
	if err := taxonomy.ValidarPapel(req.Role); err != nil {
		return Worker{}, err
	}
	mode := req.WorkspaceMode
	if mode == "" {
		mode = WorkspaceInherit
	}
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return Worker{}, err
	}
	name := resolveWorkerName(req, st.Workers)
	now := time.Now().UTC()
	w := Worker{Name: name, Role: req.Role, Model: req.Model, WorkspaceMode: mode, Status: WorkerWaiting, CreatedAt: now, LastActiveAt: now}
	if mode == WorkspaceShare || mode == WorkspaceBranch {
		wtPath := filepath.Join(sb.dir, "worktrees", name)
		branch := fmt.Sprintf("swarm/%s-%d", name, now.UnixNano())
		if err := createWorktree(req.RepoDir, wtPath, branch); err != nil {
			return Worker{}, err
		}
		w.WorktreePath = wtPath
		w.BranchName = branch
	}
	st.Workers = append(st.Workers, w)
	if err := sb.saveLocked(st); err != nil {
		return Worker{}, err
	}
	return w, nil
}

// PostMessage publica uma mensagem sequencial no quadro.
func (sb *SwarmBoard) PostMessage(from, to, body string) (Message, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return Message{}, err
	}
	msg := Message{Seq: st.NextSeq, From: from, To: to, Body: body, CreatedAt: time.Now().UTC()}
	st.NextSeq++
	st.Messages = append(st.Messages, msg)
	for i := range st.Workers {
		if st.Workers[i].Name == from {
			st.Workers[i].LastActiveAt = msg.CreatedAt
		}
	}
	if err := sb.saveLocked(st); err != nil {
		return Message{}, err
	}
	return msg, nil
}

// ReadMessages retorna mensagens com Seq > since ou pertencentes a seqs.
func (sb *SwarmBoard) ReadMessages(since int64, seqs []int64) ([]Message, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return nil, err
	}
	seqSet := make(map[int64]bool, len(seqs))
	for _, s := range seqs {
		seqSet[s] = true
	}
	var out []Message
	for _, m := range st.Messages {
		if len(seqSet) > 0 && seqSet[m.Seq] {
			out = append(out, m)
		} else if len(seqSet) == 0 && m.Seq > since {
			out = append(out, m)
		}
	}
	return out, nil
}

func matchesTarget(m Message, target string) bool {
	if target == "" || m.To == "@all" || m.To == target {
		return true
	}
	return strings.Contains(m.Body, target) || strings.Contains(m.Body, "@all")
}

// WaitMessages aguarda passivamente mensagens direcionadas a target ou @all apos since.
func (sb *SwarmBoard) WaitMessages(ctx context.Context, target string, since int64, poll time.Duration) ([]Message, error) {
	if poll <= 0 {
		poll = 100 * time.Millisecond
	}
	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
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
		if len(matched) > 0 {
			return matched, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

// Status retorna o estado atual dos workers e contagem de mensagens.
func (sb *SwarmBoard) Status() (SwarmStatus, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return SwarmStatus{}, err
	}
	return SwarmStatus{Workers: st.Workers, MessageCount: len(st.Messages)}, nil
}

// TouchWorker atualiza status e timestamp de atividade de um worker.
func (sb *SwarmBoard) TouchWorker(name string, status WorkerStatus, at time.Time) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()
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
}

// CheckStall retorna workers em execucao sem atividade dentro de timeout.
func (sb *SwarmBoard) CheckStall(timeout time.Duration) ([]Worker, error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return nil, err
	}
	cutoff := time.Now().UTC().Add(-timeout)
	var stalled []Worker
	for _, w := range st.Workers {
		if w.Status == WorkerWorking && w.LastActiveAt.Before(cutoff) {
			stalled = append(stalled, w)
		}
	}
	return stalled, nil
}

// KillWorker encerra um worker e remove seu git worktree se existir.
func (sb *SwarmBoard) KillWorker(name, repoDir string) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	st, err := sb.loadLocked()
	if err != nil {
		return err
	}
	for i := range st.Workers {
		if st.Workers[i].Name != name {
			continue
		}
		if st.Workers[i].WorktreePath != "" && strings.TrimSpace(repoDir) != "" {
			cmd := exec.Command("git", "worktree", "remove", "--force", st.Workers[i].WorktreePath)
			cmd.Dir = repoDir
			_ = cmd.Run()
		}
		st.Workers[i].Status = WorkerKilled
		return sb.saveLocked(st)
	}
	return fmt.Errorf("board: worker %q nao encontrado", name)
}
