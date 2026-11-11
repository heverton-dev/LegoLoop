package brain

import (
	"errors"
	"fmt"
	"strings"

	"legoloop/internal/interview"
	"legoloop/internal/orchestrix"
	"legoloop/internal/provider"
	"legoloop/internal/ste"
)

var (
	// ErrSemCerebro: nenhum executor de modelo está disponível nesta máquina.
	ErrSemCerebro = errors.New("brain: nenhum cérebro disponível")
	// ErrModeloDesconhecido: o modelo pedido não está no inventário.
	ErrModeloDesconhecido = errors.New("brain: modelo fora do inventário")
)

func init() {
	ste.Default.Register("brain.sem_cerebro",
		"Nenhum modelo está disponível. Suba o 9Router, instale o Ollama ou instale um harness com CLI. Depois rode legoloop doctor.")
	ste.Default.Register("brain.escolha_auto",
		"Cérebro escolhido: %s. O LegoLoop usa o primeiro executor disponível: 9Router, Ollama, depois CLI de harness. Use --model para escolher outro.")
	ste.Default.Register("brain.escolha_pedido",
		"Cérebro escolhido: %s, por pedido seu.")
	ste.Default.Register("brain.modelo_desconhecido",
		"O modelo %q não está disponível. Modelos disponíveis: %s.")
	ste.Default.Register("brain.lista_resumida", "%s e mais %d. Veja todos com legoloop doctor")
}

// Pick escolhe o cérebro: o modelo pedido, se houver, senão o primeiro executor disponível na ordem
// 9Router, Ollama, CLI de harness (a ordem de Inventory.Models). Devolve também a explicação da escolha.
func Pick(inv provider.Inventory, requested string, opts ...Option) (Brain, string, error) {
	models := inv.Models()
	if len(models) == 0 {
		return nil, "", fmt.Errorf("%w: %s", ErrSemCerebro, ste.Msg("brain.sem_cerebro"))
	}
	cfg := newConfig(opts)
	if requested != "" {
		src, ok := inv.Source(requested)
		if !ok {
			return nil, "", fmt.Errorf("%w: %s", ErrModeloDesconhecido, ste.Msg("brain.modelo_desconhecido", requested, ids(models)))
		}
		b := build(inv, cfg, src, requested)
		return b, ste.Msg("brain.escolha_pedido", b.Name()), nil
	}
	m := preferredByEffort(models, cfg.effort)
	b := build(inv, cfg, m.Source, m.ID)
	return b, ste.Msg("brain.escolha_auto", b.Name()), nil
}

// preferredByEffort seleciona o modelo de acordo com o esforço cognitivo pedido.
func preferredByEffort(models []provider.Model, effort Effort) provider.Model {
	switch effort {
	case EffortLow:
		for _, m := range models {
			id := strings.ToLower(m.ID)
			if strings.Contains(id, "fast") || strings.Contains(id, "flash") || strings.Contains(id, "haiku") || strings.Contains(id, "mini") {
				return m
			}
		}
	case EffortHigh:
		for _, m := range models {
			id := strings.ToLower(m.ID)
			if strings.Contains(id, "pro") || strings.Contains(id, "opus") || strings.Contains(id, "reasoning") || strings.Contains(id, "coder") {
				return m
			}
		}
	}
	return preferred(models)
}

// preferred é o primeiro modelo da lista; no claude CLI prefere o sonnet (bom custo e qualidade).
func preferred(models []provider.Model) provider.Model {
	first := models[0]
	if first.Source != provider.SourceClaudeCLI {
		return first
	}
	for _, m := range models {
		if m.Source == provider.SourceClaudeCLI && m.ID == "sonnet" {
			return m
		}
	}
	return first
}

// maxListados limita quantos modelos a mensagem de erro mostra: uma lista de centenas ninguém lê.
const maxListados = 15

// ids lista os modelos para o erro. Passando do limite, mostra os primeiros, quantos faltam e onde ver todos.
func ids(models []provider.Model) string {
	n := min(len(models), maxListados)
	out := make([]string, n)
	for i := range out {
		out[i] = models[i].ID
	}
	lista := strings.Join(out, ", ")
	if len(models) > maxListados {
		return ste.Msg("brain.lista_resumida", lista, len(models)-maxListados)
	}
	return lista
}

func build(inv provider.Inventory, cfg config, src provider.Source, model string) Brain {
	name := string(src) + "/" + model
	switch src {
	case provider.SourceRouter:
		c := interview.NewRouterClient(inv.RouterURL, provider.TrustedRouterKey(inv.RouterURL), model)
		c.HTTPClient.Timeout = cfg.timeout
		return routerBrain{RouterLLMClient: c, name: name, effort: cfg.effort}
	case provider.SourceClaudeCLI:
		name = "claude/" + model
	case provider.SourceOllama:
		name = "ollama/" + model
	default:
		name = model // harness: "agy:default" já carrega o nome do executor
	}
	o := &orchestrix.Orchestrix{Inventory: inv, HTTP: cfg.http, OllamaURL: cfg.ollamaURL}
	return orchBrain{o: o, src: src, model: model, name: name, limit: cfg.timeout, effort: cfg.effort}
}
