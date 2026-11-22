package progresso

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Achado descreve uma divergência entre o arquivo de progresso e a realidade.
type Achado struct {
	Regra  string `json:"regra"`
	Ticket string `json:"ticket,omitempty"`
	Motivo string `json:"motivo"`
}

// OpcoesVerificar parametriza Verificar.
type OpcoesVerificar struct {
	Agora       time.Time
	IdadeMaxima time.Duration // Ticket IN_PROGRESS mais antigo que isto é órfão; 0 desliga a checagem
	Espelho     *string       // conteúdo atual do .md; nil não confere o espelho
}

// Verificar confere a coerência do progresso sem executar nenhum comando.
func Verificar(p Progresso, opt OpcoesVerificar) []Achado {
	achados := achadosEstrutura(p)
	achados = append(achados, achadosTickets(p)...)
	achados = append(achados, achadosEspelho(p, opt.Espelho)...)
	for _, id := range Orfaos(p, opt.Agora, opt.IdadeMaxima) {
		achados = append(achados, Achado{Regra: `orfao`, Ticket: id, Motivo: `IN_PROGRESS sem registro recente`})
	}
	return achados
}

func achadosEstrutura(p Progresso) []Achado {
	var out []Achado
	if p.Schema != SchemaVersao {
		out = append(out, Achado{Regra: `schema-invalido`, Motivo: fmt.Sprintf(`schema %q, esperado %q`, p.Schema, SchemaVersao)})
	}
	calc := p
	calc.RecalcularResumo()
	if calc.Resumo != p.Resumo {
		out = append(out, Achado{Regra: `resumo-divergente`, Motivo: fmt.Sprintf(`resumo %+v, recalculado %+v`, p.Resumo, calc.Resumo)})
	}
	return out
}

func achadosTickets(p Progresso) []Achado {
	var out []Achado
	vistos := make(map[string]bool, len(p.Tickets))
	for _, it := range p.Tickets {
		if vistos[it.ID] {
			out = append(out, Achado{Regra: `ticket-duplicado`, Ticket: it.ID, Motivo: `ID repetido`})
		}
		vistos[it.ID] = true
		if it.Estado == EstadoAprovado && it.ExitCode != 0 {
			out = append(out, Achado{Regra: `aprovado-sem-exit-zero`, Ticket: it.ID, Motivo: fmt.Sprintf(`PASSED com exit_code %d`, it.ExitCode)})
		}
		if strings.TrimSpace(it.Comando) == `` {
			out = append(out, Achado{Regra: `sem-comando`, Ticket: it.ID, Motivo: `sem comando de verificação`})
		}
	}
	return out
}

func achadosEspelho(p Progresso, atual *string) []Achado {
	if atual == nil {
		return nil
	}
	p.RecalcularResumo()
	if normalizarQuebras(*atual) != normalizarQuebras(RenderizarMarkdown(p)) {
		return []Achado{{Regra: `espelho-divergente`, Motivo: `o .md difere do que o JSON gera`}}
	}
	return nil
}

func normalizarQuebras(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// Orfaos lista os Tickets IN_PROGRESS sem registro dentro da idade máxima (sessão que caiu).
func Orfaos(p Progresso, agora time.Time, idade time.Duration) []string {
	if idade <= 0 {
		return nil
	}
	var ids []string
	for _, it := range p.Tickets {
		if it.Estado != EstadoEmAndamento {
			continue
		}
		t, err := time.Parse(time.RFC3339, it.AtualizadoEm)
		if err != nil || agora.Sub(t) > idade {
			ids = append(ids, it.ID)
		}
	}
	return ids
}

// Sondar roda o Gate dos Tickets não aprovados sem alterar o estado e devolve os que já passam.
func Sondar(ctx context.Context, p Progresso, dir string, execFn Executor) []Achado {
	if execFn == nil {
		execFn = ExecutorPadrao
	}
	var out []Achado
	for _, it := range p.Tickets {
		if it.Estado == EstadoAprovado || strings.TrimSpace(it.Comando) == `` {
			continue
		}
		if code, _, err := execFn(ctx, dir, it.Comando); err == nil && code == 0 {
			out = append(out, Achado{Regra: `gate-verde-sem-aprovacao`, Ticket: it.ID, Motivo: `o Gate passa, mas o estado é ` + it.Estado})
		}
	}
	return out
}

// ReconciliarOrfaos reexecuta o Gate de cada órfão: verde vira PASSED, vermelho volta a PENDING.
func ReconciliarOrfaos(ctx context.Context, p *Progresso, dir, agora string, ids []string, execFn Executor) {
	for _, id := range ids {
		if passou, _ := ValidarTicket(ctx, p, id, dir, agora, execFn); !passou {
			reabrir(p, id, agora)
		}
	}
}

func reabrir(p *Progresso, id, agora string) {
	idx := buscarIndice(p.Tickets, id)
	if idx < 0 {
		return
	}
	it := &p.Tickets[idx]
	it.Estado = EstadoPendente
	it.Responsavel = ``
	it.Prova = resumirSaida(`órfão reaberto: ` + it.Prova)
	it.AtualizadoEm = agora
	p.AtualizadoEm = agora
	p.RecalcularResumo()
}
