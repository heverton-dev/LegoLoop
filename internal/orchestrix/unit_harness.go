package orchestrix

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/provider"
	"github.com/heverton-dev/LegoLoop/internal/runner"
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
	extra, gwErr := o.harnessGatewayEnv(src, model)
	if gwErr != nil {
		return "", runner.TokenUsage{}, o.semChave(gwErr)
	}
	out, stderrOut, code, err := o.spawnEnvCaptured(ctx, o.Inventory.CLIPath(string(src)), args, dir, nil, extra)
	if err != nil {
		return "", runner.TokenUsage{}, o.semChave(appendStderr(err, stderrOut))
	}
	return parseHarnessExecResult(string(src), spec, out, stderrOut, code)
}

func (o *Orchestrix) harnessGatewayEnv(src provider.Source, model string) ([]string, error) {
	wantsRouter := strings.HasPrefix(model, "9router/") || o.RouterKey != "" || o.Inventory.RouterURL != ""
	if !wantsRouter {
		return nil, nil
	}
	env, err := o.envGateway(dag.Node{ID: string(src), Routing: dag.RoutingConfig{Via: viaRouter}})
	if err != nil {
		return nil, err
	}
	if src != provider.SourceClaudeCLI {
		env = append(env, "OPENAI_BASE_URL="+o.Inventory.RouterURL, "OPENAI_API_KEY="+o.routerCredential())
	}
	return env, nil
}

func parseHarnessExecResult(src string, spec harnessSpec, out, stderrOut []byte, code int) (string, runner.TokenUsage, error) {
	text, in, outTok, perr := spec.parse(out)
	if perr == nil && code != 0 {
		perr = errors.New("resposta recebida, mas o processo falhou")
	}
	if perr != nil {
		return "", runner.TokenUsage{}, appendStderr(fmt.Errorf("%s saiu com código %d: %w", src, code, perr), stderrOut)
	}
	return text, runner.TokenUsage{Input: in, Output: outTok}, nil
}
