package orchestrix

import (
	"context"

	"legoloop/internal/dag"
	"legoloop/internal/provider"
	"legoloop/internal/runner"
)

// ChamarParaBateria roda uma chamada do harness numa pasta de bateria.
// Não confere a homologação: a bateria é quem homologa. A via 9router leva o gateway só ao filho.
func (o *Orchestrix) ChamarParaBateria(ctx context.Context, harness, via, model, prompt, dir string) (string, runner.TokenUsage, error) {
	if harness == "claude" {
		env, err := o.envGateway(dag.Node{ID: "bateria", Routing: dag.RoutingConfig{Harness: harness, Via: via}})
		if err != nil {
			return "", runner.TokenUsage{}, err
		}
		text, usage, err := o.callClaudeToolsEnv(ctx, model, prompt, dir, nil, unitLimitsFor(dag.Node{Nature: "micro"}), env)
		return text, usage, o.semChave(err)
	}
	text, usage, err := o.callHarnessInDir(ctx, provider.Source(harness), estadoHomologado, model, prompt, dir, dag.EffortLow)
	return text, usage, o.semChave(err)
}
