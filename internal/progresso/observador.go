package progresso

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"
)

var errSemTicket = errors.New(`nó sem Ticket no progresso`)

var reIniciativaNoID = regexp.MustCompile(`^([A-Za-z]+)-[A-Za-z]+\d+-(.+)$`)

// IDCanonico iguala o ID do nó do DAG ao ID do Ticket: TICK-I4-01.1 e TICK-01.1 viram TICK-01.1.
func IDCanonico(id string) string {
	id = strings.TrimSpace(id)
	if m := reIniciativaNoID.FindStringSubmatch(id); m != nil {
		return m[1] + `-` + m[2]
	}
	return id
}

// ResolverTicket acha o Ticket do progresso que corresponde ao nó do DAG.
func ResolverTicket(p Progresso, noID string) (string, bool) {
	if idx := buscarIndice(p.Tickets, noID); idx >= 0 {
		return p.Tickets[idx].ID, true
	}
	alvo := IDCanonico(noID)
	for _, it := range p.Tickets {
		if IDCanonico(it.ID) == alvo {
			return it.ID, true
		}
	}
	return ``, false
}

// Observador mantém o progresso vivo durante uma execução: o runner e o swarm avisam o início e o fim de cada nó.
type Observador struct {
	ctx         context.Context
	arquivo     string
	dir         string
	responsavel string
	exec        Executor
	agora       func() time.Time
	mu          sync.Mutex
	erros       []error
}

// NovoObservador cria o observador de um arquivo de progresso; dir é a raiz onde roda o comando de cada Ticket.
func NovoObservador(ctx context.Context, arquivo, dir, responsavel string) *Observador {
	return &Observador{ctx: ctx, arquivo: arquivo, dir: dir, responsavel: responsavel, agora: time.Now}
}

// Erros devolve as falhas de gravação ocorridas até agora.
func (o *Observador) Erros() []error {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]error(nil), o.erros...)
}

// NoInicio marca o Ticket do nó como IN_PROGRESS.
func (o *Observador) NoInicio(noID string) {
	o.aplicar(noID, func(p *Progresso, id string) error {
		return Iniciar(p, id, o.responsavel, `nó iniciado: `+noID, o.agoraStr())
	})
}

// NoFim fecha o Ticket do nó: nó reprovado vira FAILED; nó aprovado roda o comando de verificação do Ticket.
func (o *Observador) NoFim(noID string, passou bool, exitCode int, detalhe string) {
	if !passou {
		o.aplicar(noID, func(p *Progresso, id string) error {
			marcarFim(p, id, EstadoReprovado, exitCode, detalhe, o.agoraStr())
			return nil
		})
		return
	}
	p, err := Carregar(o.arquivo)
	if err != nil {
		o.registrar(err)
		return
	}
	id, ok := ResolverTicket(p, noID)
	if !ok {
		return
	}
	comando := strings.TrimSpace(p.Tickets[buscarIndice(p.Tickets, id)].Comando)
	if comando == `` {
		o.aplicar(noID, func(p *Progresso, id string) error {
			marcarFim(p, id, EstadoAprovado, 0, `Gate do nó aprovado`, o.agoraStr())
			return nil
		})
		return
	}
	code, saida := o.executar(comando)
	o.aplicar(noID, func(p *Progresso, id string) error {
		it := &p.Tickets[buscarIndice(p.Tickets, id)]
		aplicarResultadoValidacao(p, it, code, saida, o.agoraStr())
		return nil
	})
}

// executar roda o comando do Ticket fora da trava, para um comando lento não bloquear os outros nós.
func (o *Observador) executar(comando string) (int, string) {
	execFn := o.exec
	if execFn == nil {
		execFn = ExecutorPadrao
	}
	code, saida, err := execFn(o.ctx, o.dir, comando)
	if err != nil {
		return 1, err.Error()
	}
	return code, saida
}

func (o *Observador) registrar(err error) {
	o.mu.Lock()
	o.erros = append(o.erros, err)
	o.mu.Unlock()
}

func marcarFim(p *Progresso, id, estado string, exitCode int, detalhe, agora string) {
	it := &p.Tickets[buscarIndice(p.Tickets, id)]
	it.Estado = estado
	it.ExitCode = exitCode
	it.Prova = resumirSaida(detalhe)
	it.AtualizadoEm = agora
	p.AtualizadoEm = agora
	p.RecalcularResumo()
}

func (o *Observador) aplicar(noID string, fn func(*Progresso, string) error) {
	err := Atualizar(o.arquivo, func(p *Progresso) error {
		id, ok := ResolverTicket(*p, noID)
		if !ok {
			return errSemTicket
		}
		return fn(p, id)
	})
	if err != nil && !errors.Is(err, errSemTicket) {
		o.registrar(err)
	}
}

func (o *Observador) agoraStr() string {
	return o.agora().UTC().Format(time.RFC3339)
}
