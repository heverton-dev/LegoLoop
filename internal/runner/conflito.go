package runner

import (
	"context"
	"path"
	"path/filepath"
	"sync"

	"legoloop/internal/dag"
)

// pesoPadraoMB é o peso assumido de um nó quando o escalonador não recebe a medição.
const pesoPadraoMB = 750

type ativo struct {
	arquivos []string
	pesoMB   int
}

// agenda decide quais nós prontos começam agora. Um nó espera se divide arquivo do context_out
// com um nó ativo, se o teto de paralelismo está cheio ou se a memória não comporta mais um worker.
type agenda struct {
	mu      sync.Mutex
	ativos  map[string]ativo
	mudou   chan struct{}
	teto    int
	workers func(pesosMB []int) int
}

func novaAgenda(teto int, workers func([]int) int) *agenda {
	return &agenda{ativos: map[string]ativo{}, mudou: make(chan struct{}), teto: teto, workers: workers}
}

func arquivosDe(n dag.Node) []string {
	out := make([]string, 0, len(n.ContextOut))
	for _, f := range n.ContextOut {
		out = append(out, path.Clean(filepath.ToSlash(f)))
	}
	return out
}

func compartilham(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}

// cabe informa se o nó pode começar agora. Com a agenda vazia, o nó sempre cabe (mínimo de 1 worker).
func (a *agenda) cabe(arquivos []string, pesoMB int) bool {
	if len(a.ativos) == 0 {
		return true
	}
	if len(a.ativos) >= a.teto {
		return false
	}
	pesos := []int{pesoMB}
	for _, at := range a.ativos {
		if compartilham(arquivos, at.arquivos) {
			return false
		}
		pesos = append(pesos, at.pesoMB)
	}
	return a.workers == nil || len(a.ativos)+1 <= a.workers(pesos)
}

// entrar bloqueia até o nó caber ou até o contexto acabar.
func (a *agenda) entrar(ctx context.Context, id string, arquivos []string, pesoMB int) error {
	for {
		a.mu.Lock()
		if a.cabe(arquivos, pesoMB) {
			a.ativos[id] = ativo{arquivos: arquivos, pesoMB: pesoMB}
			a.mu.Unlock()
			return nil
		}
		espera := a.mudou
		a.mu.Unlock()
		select {
		case <-espera:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *agenda) sair(id string) {
	a.mu.Lock()
	delete(a.ativos, id)
	close(a.mudou)
	a.mudou = make(chan struct{})
	a.mu.Unlock()
}

func (r *run) pesoMB(n dag.Node) int {
	if r.opts.PesoMB != nil {
		return r.opts.PesoMB(n)
	}
	return pesoPadraoMB
}
