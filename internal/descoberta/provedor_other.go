//go:build !windows && !linux && !darwin

package descoberta

import "runtime"

type provedorGenerico struct{}

// ProvedorPadrao devolve o provedor padrao para outros sistemas operacionais.
func ProvedorPadrao() Provedor {
	return provedorGenerico{}
}

func (provedorGenerico) SO() string                { return runtime.GOOS }
func (provedorGenerico) ListarUnidades() []Unidade { return nil }
