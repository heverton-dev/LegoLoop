package orchestrix

import (
	"context"
	"errors"
	"fmt"

	"legoloop/internal/dag"
	"legoloop/internal/roteamento"
	"legoloop/internal/runner"
)

// harnessDireto é a chamada HTTP ao 9Router sem harness (SPEC 4.1).
const harnessDireto = "direto"

// chamada é o resultado de um despacho: o texto, o consumo, as notas e o modelo realmente usado.
// O modelo fica vazio quando o candidato declarado respondeu.
type chamada struct {
	text   string
	usage  runner.TokenUsage
	notas  []string
	modelo string
	// harness e via ficam vazios quando o candidato declarado respondeu.
	harness, via string
}

// despachar roda o nó no candidato declarado. Se ele falhar, a política da classe decide:
// nano e micro seguem para o próximo candidato elegível; as demais classes falham sem trocar de harness.
func (o *Orchestrix) despachar(ctx context.Context, n dag.Node, plan unitPlan, prompt, dir string) (chamada, error) {
	text, usage, err := o.chamarNo(ctx, n, plan.model, prompt, dir)
	if err == nil {
		return chamada{text: text, usage: usage}, nil
	}
	if ctx.Err() != nil {
		return chamada{usage: usage}, err
	}
	if roteamento.PoliticaDeFalha(n.Nature) != roteamento.PoliticaDegrada {
		return chamada{usage: usage}, fmt.Errorf("harness %s/%s indisponível: divergência entre o DAG e o inventário; a classe %q não troca de harness: %w",
			harnessDe(n), n.Routing.Via, n.Nature, err)
	}
	return o.degradar(ctx, n, prompt, dir, usage, err)
}

func harnessDe(n dag.Node) string {
	if n.Routing.Harness == "" {
		return "claude"
	}
	return n.Routing.Harness
}

func pedidoDe(n dag.Node) roteamento.CandidatoVia {
	return roteamento.CandidatoVia{Harness: harnessDe(n), Via: roteamento.Via(n.Routing.Via), Modelo: n.Routing.Model}
}

// noComCandidato devolve o nó com o roteamento trocado pelo candidato.
func noComCandidato(n dag.Node, c roteamento.CandidatoVia) dag.Node {
	n.Routing.Harness, n.Routing.Via = c.Harness, string(c.Via)
	n.Routing.Model, n.Routing.Effort = c.Modelo, dag.EffortLevel(c.Esforco)
	return n
}

// degradar tenta os demais candidatos elegíveis da classe, na ordem do perfil.
func (o *Orchestrix) degradar(ctx context.Context, n dag.Node, prompt, dir string, gasto runner.TokenUsage, falha error) (chamada, error) {
	pedido := pedidoDe(n)
	erros := []error{falha}
	for _, c := range o.Elegiveis[n.Nature] {
		if c == pedido {
			continue
		}
		alt := noComCandidato(n, c)
		model, err := o.unitModelFor(alt)
		if err == nil {
			var text string
			var usage runner.TokenUsage
			text, usage, err = o.chamarNo(ctx, alt, model, prompt, dir)
			gasto = addUsage(gasto, usage)
			if err == nil {
				return o.registrarSubstituicao(n, pedido, c, falha, text, gasto), nil
			}
		}
		erros = append(erros, fmt.Errorf("candidato %s/%s: %w", c.Harness, c.Modelo, err))
		if ctx.Err() != nil {
			break
		}
	}
	return chamada{usage: gasto}, fmt.Errorf("nenhum candidato da classe %q respondeu: %w", n.Nature, errors.Join(erros...))
}

// registrarSubstituicao grava a Substituicao nas notas do nó e avisa na tela.
func (o *Orchestrix) registrarSubstituicao(n dag.Node, pedido, usado roteamento.CandidatoVia, falha error, text string, gasto runner.TokenUsage) chamada {
	s := roteamento.SubstituicaoVia{No: n.ID, Classe: n.Nature, Pedido: pedido, Usado: usado,
		Motivo: "candidato declarado indisponível em execução: " + o.semChaveTexto(falha.Error())}
	nota := fmt.Sprintf("substituicao: %s usou %s/%s/%s no lugar de %s/%s/%s: %s", s.No,
		usado.Harness, usado.Via, usado.Modelo, pedido.Harness, pedido.Via, pedido.Modelo, s.Motivo)
	if o.Aviso != nil {
		o.Aviso("aviso: " + nota)
	}
	return chamada{text: text, usage: gasto, notas: []string{nota}, modelo: usado.Modelo, harness: usado.Harness, via: string(usado.Via)}
}
