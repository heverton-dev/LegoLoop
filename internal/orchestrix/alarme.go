package orchestrix

import (
	"sync"

	"legoloop/internal/roteamento/homologacao"
)

type RegistroReprovados struct {
	mu     sync.Mutex
	motivo map[string]string
}

func chavePar(harness, via string) string { return harness + "|" + via }

func (r *RegistroReprovados) Marcar(harness, via, motivo string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.motivo == nil {
		r.motivo = map[string]string{}
	}
	r.motivo[chavePar(harness, via)] = motivo
}

func (r *RegistroReprovados) Motivo(harness, via string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.motivo[chavePar(harness, via)]
	return m, ok
}

func (o *Orchestrix) estadoEfetivo(harness, via string) string {
	if o.Reprovados != nil {
		if _, ok := o.Reprovados.Motivo(harness, via); ok {
			return "reprovado"
		}
	}
	if o.EstadoHarness == nil {
		return ""
	}
	return o.EstadoHarness(harness, via)
}

func motivoBanimento(texto string, err error) string {
	t := homologacao.TrechoBanimento(texto)
	if t == "" && err != nil {
		t = homologacao.TrechoBanimento(err.Error())
	}
	if t == "" {
		return ""
	}
	return "banimento: " + t
}

func (o *Orchestrix) vigiarBanimento(harness, via, texto string, err error) {
	motivo := motivoBanimento(texto, err)
	if motivo == "" {
		return
	}
	if o.Reprovados != nil {
		o.Reprovados.Marcar(harness, via, motivo)
	}
	if o.AoReprovar != nil {
		o.AoReprovar(harness, via, motivo)
	}
}
