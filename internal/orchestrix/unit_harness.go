package orchestrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
)

const estadoHomologado = "homologado"

func (o *Orchestrix) callHarnessInDir(ctx context.Context, src provider.Source, estado, model, prompt, dir string, effort dag.EffortLevel) (string, runner.TokenUsage, error) {
	if dir == "" {
		return "", runner.TokenUsage{}, errors.New("callHarnessInDir: pasta de trabalho vazia")
	}
	if estado != estadoHomologado {
		return "", runner.TokenUsage{}, fmt.Errorf("harness %q recusado: estado %q, exige %q", src, estado, estadoHomologado)
	}
	spec, ok := harnessSpecs()[src]
	if !ok {
		return "", runner.TokenUsage{}, fmt.Errorf("executor desconhecido %q", src)
	}
	args, withFile, err := spec.args(harnessModel(model), prompt, effort)
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	if withFile {
		path := filepath.Join(dir, promptFile)
		if err := os.WriteFile(path, []byte(prompt), 0o600); err != nil {
			return "", runner.TokenUsage{}, err
		}
		defer os.Remove(path)
	}
	out, code, err := o.spawn(ctx, o.Inventory.CLIPath(string(src)), args, dir, nil)
	if err != nil {
		return "", runner.TokenUsage{}, err
	}
	text, in, outTok, perr := spec.parse(out)
	if perr == nil && code != 0 {
		perr = errors.New("resposta recebida, mas o processo falhou")
	}
	if perr != nil {
		return "", runner.TokenUsage{}, fmt.Errorf("%s saiu com código %d: %w", src, code, perr)
	}
	return text, runner.TokenUsage{Input: in, Output: outTok}, nil
}
