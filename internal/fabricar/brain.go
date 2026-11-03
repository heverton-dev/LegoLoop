package fabricar

import (
	"errors"
	"fmt"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/provider"
)

// BrainTimeout cobre a geração da SPEC: prompt grande e modelo de raciocínio levam mais que os 45 s padrão
// do cliente, que servem às perguntas curtas da entrevista.
const BrainTimeout = 5 * time.Minute

// ErrBrain: o modelo escolhido não pode conduzir a fabricação (router fora do ar ou modelo de outro executor).
var ErrBrain = errors.New("fabricar: cérebro indisponível")

// BrainFor devolve o LLM que conduz a fabricação: o modelo escolhido pelo usuário (o "cérebro" da INTENT),
// servido pelo 9Router. Modelos de outro executor (ollama, claude CLI) não falam o protocolo do router e
// são recusados com uma explicação, em vez de trocados em silêncio por outro modelo. A chave só vai à base confiável.
func BrainFor(inv provider.Inventory, model string) (*interview.RouterLLMClient, error) {
	if inv.RouterURL == "" {
		return nil, fmt.Errorf("%w: a Tríade precisa de um LLM, suba o 9Router (veja legoloop doctor)", ErrBrain)
	}
	if src, _ := inv.Source(model); src != provider.SourceRouter {
		return nil, fmt.Errorf("%w: o modelo %q não é servido pelo 9Router (executor %q); escolha um modelo com executor 9router em legoloop doctor", ErrBrain, model, src)
	}
	c := interview.NewRouterClient(inv.RouterURL, provider.TrustedRouterKey(inv.RouterURL), model)
	c.HTTPClient.Timeout = BrainTimeout
	return c, nil
}
