package cli

import (
	"context"
	"errors"
	"time"

	"legoloop/internal/orchestrix"
	"legoloop/internal/roteamento/homologacao"
)

// harnessReal liga a bateria ao motor: cada chamada passa pelo mesmo caminho dos nós de código.
type harnessReal struct {
	ox                  *orchestrix.Orchestrix
	harness, via, model string
}

// Executar roda uma chamada com o teto de duração. O motor não mede o pico de memória do processo
// filho nesta versão: PicoMB fica 0 e o peso vale o padrão de 750 MB.
func (h harnessReal) Executar(ctx context.Context, dir, prompt string, teto time.Duration) (homologacao.Saida, error) {
	ctx, cancelar := context.WithTimeout(ctx, teto)
	defer cancelar()
	modelo := h.model
	if prompt == homologacao.PromptB {
		modelo = "modelo-inexistente-da-bateria" // erro forçado de verdade: o C5 exige exit diferente de zero
	}
	inicio := time.Now()
	texto, uso, err := h.ox.ChamarParaBateria(ctx, h.harness, h.via, modelo, prompt, dir)
	saida := homologacao.Saida{Texto: texto, CustoUSD: uso.CostUSD, Tokens: uso.Input + uso.Output, Duracao: time.Since(inicio)}
	if err == nil {
		return saida, nil
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return saida, context.DeadlineExceeded
	}
	saida.ExitCode = 1
	return saida, nil
}

// adaptadorReal devolve a fábrica de harness da bateria sobre o motor real.
func adaptadorReal(ox *orchestrix.Orchestrix) func(harness, via, model string) homologacao.Harness {
	return func(harness, via, model string) homologacao.Harness {
		return harnessReal{ox: ox, harness: harness, via: via, model: model}
	}
}
