package orchestrix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/dag"
	"legoloop/internal/loop"
	"legoloop/internal/runner"
)

// runSubloop executa o sub-loop em diretório isolado e estéril (R-LOOP-4).
func (o *Orchestrix) runSubloop(ctx context.Context, n dag.Node, workdir string) (runner.Result, error) {
	if o.LoopStore == nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: acervo de loops não configurado", n.ID)
	}
	lp, err := o.LoopStore.Get(n.Ref)
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: nó %q: sub-loop %q não encontrado: %w", n.ID, n.Ref, err)
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
	subReport, err := runner.Run(ctx, subDAG, runner.Options{
		Workdir:     subWorkdir,
		Parallelism: 1,
		Executor:    o,
	})
	if err != nil {
		return runner.Result{}, fmt.Errorf("orchestrix: sub-loop %q: %w", n.Ref, err)
	}
	if !subReport.OK {
		return buildFailedResult(n, subReport)
	}
	return o.finalizeSubloop(n, workdir, subWorkdir, subReport)
}

func parseSubDAG(lp *loop.Loop) (*dag.DAG, error) {
	raw, err := json.Marshal(struct {
		ID      string       `json:"id,omitempty"`
		Name    string       `json:"name"`
		Version string       `json:"version,omitempty"`
		Trigger *dag.Trigger `json:"trigger,omitempty"`
		Nodes   []dag.Node   `json:"nodes"`
	}{ID: lp.ID, Name: lp.Name, Version: lp.Version, Trigger: lp.Trigger, Nodes: lp.Nodes})
	if err != nil {
		return nil, err
	}
	return dag.Parse(raw)
}

func (o *Orchestrix) transferInputs(n dag.Node, lp *loop.Loop, srcDir, dstDir string) error {
	toCopy := make(map[string]bool)
	for _, art := range n.Inputs {
		if art != "" {
			toCopy[filepath.Clean(art)] = true
		}
	}
	if len(toCopy) == 0 {
		ports, err := loop.DerivePorts(lp)
		if err == nil {
			for _, c := range ports.Consumes {
				if c != "" && c != loop.AnyArtifact {
					toCopy[filepath.Clean(c)] = true
				}
			}
		}
	}
	for art := range toCopy {
		src := filepath.Join(srcDir, art)
		if data, err := os.ReadFile(src); err == nil {
			dst := filepath.Join(dstDir, art)
			if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(dst, data, 0600); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildFailedResult(n dag.Node, rep runner.Report) (runner.Result, error) {
	var msgs []string
	for _, nr := range rep.Nodes {
		if nr.Status != runner.StatusPassed {
			msgs = append(msgs, fmt.Sprintf("nó %q: %s", nr.ID, nr.Status))
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
