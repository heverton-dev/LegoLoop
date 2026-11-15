package orchestrix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/runner"
	"legoloop/internal/studiodoc"
)

// runSubloop executa o sub-loop em diretório isolado e estéril (R-LOOP-4).
func (o *Orchestrix) runSubloop(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	lp, err := o.loadSubProcess(n)
	if err != nil {
		return runner.Result{}, err
	}
	subWorkdir, err := os.MkdirTemp("", "legoloop-subloop-*")
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: criar workspace estéril para sub-loop: %w", err)
	}
	defer os.RemoveAll(subWorkdir)

	if err := o.transferInputs(n, lp, workdir, subWorkdir); err != nil {
		return runner.Result{}, err
	}
	subDAG, err := parseSubDAG(lp)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: %w", n.ID, err)
	}
	subReport, err := runner.Run(ctx, subDAG, o.subRunnerOptions(subWorkdir))
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: sub-loop %q: %w", n.Ref, err)
	}
	if !subReport.OK {
		return buildFailedResult(n, subReport)
	}
	return o.finalizeSubloop(n, workdir, subWorkdir, subReport)
}

func (o *Orchestrix) loadSubProcess(n dag.Node) (*studiodoc.Processo, error) {
	if o.LoopStore == nil {
		return nil, fmt.Errorf("orchestrix: nó %q: acervo de loops não configurado", n.ID)
	}
	lp, err := o.LoopStore.Get(n.Ref)
	if err != nil {
		return nil, fmt.Errorf("orchestrix: nó %q: sub-loop %q não encontrado: %w", n.ID, n.Ref, err)
	}
	return lp, nil
}

func (o *Orchestrix) subRunnerOptions(subWorkdir string) runner.Options {
	return runner.Options{
		Workdir:        subWorkdir,
		Parallelism:    1,
		DefaultTimeout: o.DefaultTimeout,
		Executor:       o,
		AllowScripts:   o.AllowScripts,
		AutoApprove:    o.AutoApprove,
		NonInteractive: o.NonInteractive,
		Vars:           o.Vars,
	}
}

func parseSubDAG(lp *studiodoc.Processo) (*dag.DAG, error) {
	raw, err := json.Marshal(struct {
		Schema             string         `json:"$schema,omitempty"`
		ID                 string         `json:"id,omitempty"`
		Name               string         `json:"name"`
		Description        string         `json:"description,omitempty"`
		Version            string         `json:"version,omitempty"`
		SchemaVersion      int            `json:"schema_version,omitempty"`
		Date               string         `json:"date,omitempty"`
		UpstreamSignatures map[string]any `json:"upstream_signatures,omitempty"`
		ExecutionStrategy  map[string]any `json:"execution_strategy,omitempty"`
		Trigger            *dag.Trigger   `json:"trigger,omitempty"`
		Nodes              []dag.Node     `json:"nodes"`
	}{
		Schema:             lp.Schema,
		ID:                 lp.ID,
		Name:               lp.Name,
		Description:        lp.Description,
		Version:            lp.Version,
		SchemaVersion:      lp.SchemaVersion,
		Date:               lp.Date,
		UpstreamSignatures: lp.UpstreamSignatures,
		ExecutionStrategy:  lp.ExecutionStrategy,
		Trigger:            lp.Trigger,
		Nodes:              lp.Nodes,
	})
	if err != nil {
		return nil, err
	}
	return dag.Parse(raw)
}

func (o *Orchestrix) transferInputs(n dag.Node, lp *studiodoc.Processo, srcDir, dstDir string) error {
	toCopy, err := collectSubloopInputs(n, lp)
	if err != nil {
		return err
	}
	if len(toCopy) == 0 {
		return nil
	}
	srcRoot, err := os.OpenRoot(srcDir)
	if err != nil {
		return err
	}
	defer srcRoot.Close()
	dstRoot, err := os.OpenRoot(dstDir)
	if err != nil {
		return err
	}
	defer dstRoot.Close()
	for _, rel := range toCopy {
		if err := copySubloopInput(srcRoot, dstRoot, rel); err != nil {
			return err
		}
	}
	return nil
}

func collectSubloopInputs(n dag.Node, lp *studiodoc.Processo) ([]string, error) {
	seen := make(map[string]bool)
	var out []string
	add := func(raw string) error {
		rel, err := sanitizeSubloopRel(raw)
		if err != nil {
			return err
		}
		if rel != "" && !seen[rel] {
			seen[rel] = true
			out = append(out, rel)
		}
		return nil
	}
	for _, art := range n.Inputs {
		if err := add(art); err != nil {
			return nil, err
		}
	}
	if len(out) == 0 && lp != nil {
		if ports, err := studiodoc.DerivePorts(lp); err == nil {
			for _, c := range ports.Consumes {
				if c == "" || c == studiodoc.AnyArtifact || (strings.Contains(c, ":") && !filepath.IsAbs(c)) {
					continue
				}
				if err := add(c); err != nil {
					return nil, err
				}
			}
		}
	}
	return out, nil
}

func sanitizeSubloopRel(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}
	rel := filepath.Clean(filepath.FromSlash(trimmed))
	if !filepath.IsLocal(rel) || rel == "." {
		return "", fmt.Errorf("path jail: entrada %q fora do workspace", raw)
	}
	if isSensitiveSubloopRel(rel) {
		return "", fmt.Errorf("path jail: entrada sensível %q proibida", raw)
	}
	return rel, nil
}

func isSensitiveSubloopRel(rel string) bool {
	for _, seg := range strings.Split(filepath.ToSlash(rel), "/") {
		low := strings.ToLower(seg)
		if low == ".git" || strings.HasPrefix(low, ".env") || strings.HasPrefix(low, "credential_") {
			return true
		}
	}
	return false
}

func copySubloopInput(srcRoot, dstRoot *os.Root, rel string) error {
	info, err := srcRoot.Lstat(rel)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("path jail: symlink %q proibido", rel)
	}
	data, err := srcRoot.ReadFile(rel)
	if err != nil {
		return err
	}
	return writeArtifact(dstRoot, rel, data)
}

func buildFailedResult(n dag.Node, rep runner.Report) (runner.Result, error) {
	var msgs []string
	for _, nr := range rep.Nodes {
		if nr.Status != runner.StatusPassed {
			if nr.Error != "" {
				msgs = append(msgs, fmt.Sprintf("nó %q: %s (%s)", nr.ID, nr.Status, nr.Error))
			} else {
				msgs = append(msgs, fmt.Sprintf("nó %q: %s", nr.ID, nr.Status))
			}
		}
	}
	return runner.Result{
		ExitCode: 1,
		Tokens:   rep.Tokens,
		Model:    "subloop/" + n.Ref,
		Notes:    msgs,
	}, fmt.Errorf("orchestrix: sub-loop %q falhou (%s)", n.Ref, strings.Join(msgs, "; "))
}

// subArtifact é o caminho do artefato final dentro do sub-loop: sub_artifact, ou o próprio artifact quando
// o DAG não separa os dois caminhos (contrato antigo).
func subArtifact(n dag.Node) string {
	if n.Output.SubArtifact != "" {
		return n.Output.SubArtifact
	}
	return n.Output.Artifact
}

// readSubArtifact lê o artefato final dentro do diretório do sub-loop, sem sair dele.
func readSubArtifact(subWorkdir, rel string) ([]byte, error) {
	root, err := os.OpenRoot(subWorkdir)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return root.ReadFile(filepath.FromSlash(rel))
}

// finalizeSubloop lê o artefato final no sub-loop (subArtifact) e grava no pai em n.Output.Artifact.
func (o *Orchestrix) finalizeSubloop(n dag.Node, parentDir, subWorkdir string, rep runner.Report) (runner.Result, error) {
	sub := subArtifact(n)
	data, err := readSubArtifact(subWorkdir, sub)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: artefato final %q não produzido pelo sub-loop: %w", sub, err)
	}
	root, err := os.OpenRoot(parentDir)
	if err != nil {
		return runner.Result{}, err
	}
	defer root.Close()
	if err := writeArtifact(root, n.Output.Artifact, data); err != nil {
		return runner.Result{}, err
	}
	var notes []string
	for _, nr := range rep.Nodes {
		notes = append(notes, nr.Notes...)
	}
	return runner.Result{
		ExitCode: 0,
		Tokens:   rep.Tokens,
		Model:    "subloop/" + n.Ref,
		Notes:    notes,
	}, nil
}
