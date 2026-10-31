package assets

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"sync"
	"time"
)

// Watcher mantém o inventário atualizado: re-varre a cada intervalo (barato graças ao cache do
// Scanner) e só publica uma nova Version quando o conteúdo realmente muda. Usa varredura
// periódica, não notificações do SO, para manter a regra "somente biblioteca padrão".
type Watcher struct {
	scanner  *Scanner
	interval time.Duration
	done     chan struct{}

	refreshMu sync.Mutex // serializa Refresh (varredura + publicação)
	mu        sync.RWMutex
	snap      Snapshot
	sig       string
}

// NewWatcher cria o vigia; nada roda até Start ou Refresh.
func NewWatcher(s *Scanner, interval time.Duration) *Watcher {
	return &Watcher{scanner: s, interval: interval, done: make(chan struct{})}
}

// Start varre agora e a cada intervalo, até o contexto terminar. Não bloqueia.
func (w *Watcher) Start(ctx context.Context) {
	go func() {
		defer close(w.done)
		w.Refresh(ctx)
		t := time.NewTicker(w.interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				w.Refresh(ctx)
			}
		}
	}()
}

// Roots devolve os diretórios de ativos existentes (base da canonização).
func (w *Watcher) Roots() []string { return w.scanner.Roots() }

// Done fecha quando o vigia iniciado por Start termina.
func (w *Watcher) Done() <-chan struct{} { return w.done }

// Snapshot devolve o último inventário publicado (Version 0 = ainda não houve varredura).
func (w *Watcher) Snapshot() Snapshot {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.snap
}

// Refresh varre agora e publica se algo mudou. Varredura interrompida por cancelamento
// é descartada: um inventário parcial nunca substitui o anterior.
func (w *Watcher) Refresh(ctx context.Context) Snapshot {
	w.refreshMu.Lock()
	defer w.refreshMu.Unlock()
	snap := w.scanner.Scan(ctx)
	if ctx.Err() != nil {
		return w.Snapshot()
	}
	sig := signature(snap.Items)
	w.mu.Lock()
	defer w.mu.Unlock()
	if sig != w.sig {
		snap.Version = w.snap.Version + 1
		w.snap, w.sig = snap, sig
	}
	return w.snap
}

// signature resume o conteúdo do inventário; independe de ScannedAt.
func signature(items []Item) string {
	h := sha256.New()
	var b strings.Builder
	for _, it := range items {
		b.Reset()
		b.WriteString(string(it.Kind) + "\x00" + it.Name + "\x00" + it.Description + "\x00" + it.Path + "\x00" + it.Scope + "\x00" + it.Origin)
		keys := make([]string, 0, len(it.Meta))
		for k := range it.Meta {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			b.WriteString("\x00" + k + "=" + it.Meta[k])
		}
		b.WriteString("\n")
		h.Write([]byte(b.String()))
	}
	return hex.EncodeToString(h.Sum(nil))
}
