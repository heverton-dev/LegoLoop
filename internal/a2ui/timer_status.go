package a2ui

import (
	"context"
	"sync"
	"time"
)

const defaultSessionTimeout = 5 * time.Minute

// TimerSnapshot representa o estado observavel do cronometro regressivo da sessao.
type TimerSnapshot struct {
	TimeoutMs   int64 `json:"timeout_ms"`
	RemainingMs int64 `json:"remaining_ms"`
	Paused      bool  `json:"paused"`
	Active      bool  `json:"active"`
	Expired     bool  `json:"expired"`
}

// EtapaStatus representa uma etapa da arvore viva de execucao exibida na A2UI.
type EtapaStatus struct {
	ID      string `json:"id"`
	Titulo  string `json:"titulo"`
	Status  string `json:"status"`
	Detalhe string `json:"detalhe,omitempty"`
}

// SessionStatus descreve a atividade atual do motor e as etapas de execucao.
type SessionStatus struct {
	Etapa              string        `json:"etapa,omitempty"`
	Mensagem           string        `json:"mensagem,omitempty"`
	Etapas             []EtapaStatus `json:"etapas,omitempty"`
	AguardandoResposta bool          `json:"aguardando_resposta,omitempty"`
	AtualizadoEm       string        `json:"atualizado_em,omitempty"`
}

// ControlEvent transporta atualizacoes de timer e status via SSE.
type ControlEvent struct {
	Kind   string         `json:"kind"`
	Timer  *TimerSnapshot `json:"timer,omitempty"`
	Status *SessionStatus `json:"status,omitempty"`
}

// SessionTimer controla contagem regressiva, pausa, retomada e extensao por etapa.
type SessionTimer struct {
	mu        sync.Mutex
	timeout   time.Duration
	remaining time.Duration
	paused    bool
	active    bool
	expired   bool
	lastTick  time.Time
	notifyCh  chan struct{}
}

// NewSessionTimer cria um temporizador com timeout padrao.
func NewSessionTimer(timeout time.Duration) *SessionTimer {
	if timeout <= 0 {
		timeout = defaultSessionTimeout
	}
	return &SessionTimer{
		timeout:   timeout,
		remaining: timeout,
		notifyCh:  make(chan struct{}, 1),
	}
}

func (t *SessionTimer) signalLocked() {
	select {
	case t.notifyCh <- struct{}{}:
	default:
	}
}

func (t *SessionTimer) advanceLocked(now time.Time) {
	if !t.active || t.paused || t.expired {
		t.lastTick = now
		return
	}
	if !t.lastTick.IsZero() {
		elapsed := now.Sub(t.lastTick)
		if elapsed > 0 {
			t.remaining -= elapsed
			if t.remaining <= 0 {
				t.remaining = 0
				t.expired = true
				t.active = false
			}
		}
	}
	t.lastTick = now
}

// SetTimeout altera o tempo limite padrao por etapa da sessao.
func (t *SessionTimer) SetTimeout(d time.Duration) TimerSnapshot {
	if d <= 0 {
		return t.Snapshot()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.advanceLocked(time.Now())
	t.timeout = d
	if !t.active || t.remaining > d || !t.expired {
		t.remaining = d
		t.expired = false
	}
	t.signalLocked()
	return t.snapshotLocked()
}

// Start inicia ou reinicia a contagem regressiva para uma nova espera.
func (t *SessionTimer) Start(override time.Duration) TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	dur := t.timeout
	if override > 0 {
		dur = override
	}
	t.remaining = dur
	t.active = true
	t.paused = false
	t.expired = false
	t.lastTick = time.Now()
	t.signalLocked()
	return t.snapshotLocked()
}

// Pause congela a contagem regressiva mantendo a sessao ativa.
func (t *SessionTimer) Pause() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.advanceLocked(time.Now())
	if t.active && !t.expired {
		t.paused = true
	}
	t.signalLocked()
	return t.snapshotLocked()
}

// Resume retoma a contagem regressiva do ponto congelado.
func (t *SessionTimer) Resume() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.active && t.paused && !t.expired {
		t.paused = false
		t.lastTick = time.Now()
	}
	t.signalLocked()
	return t.snapshotLocked()
}

// Extend soma tempo extra ao contador restante.
func (t *SessionTimer) Extend(extra time.Duration) TimerSnapshot {
	if extra <= 0 {
		return t.Snapshot()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.advanceLocked(time.Now())
	t.remaining += extra
	t.expired = false
	t.signalLocked()
	return t.snapshotLocked()
}

// Stop encerra a contagem ativa quando a resposta do usuario chega.
func (t *SessionTimer) Stop() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.active = false
	t.paused = false
	t.expired = false
	t.remaining = t.timeout
	t.signalLocked()
	return t.snapshotLocked()
}

// Snapshot devolve uma fotografia consistente do estado atual do timer.
func (t *SessionTimer) Snapshot() TimerSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.advanceLocked(time.Now())
	return t.snapshotLocked()
}

func (t *SessionTimer) snapshotLocked() TimerSnapshot {
	rem := t.remaining
	if rem < 0 {
		rem = 0
	}
	return TimerSnapshot{
		TimeoutMs:   t.timeout.Milliseconds(),
		RemainingMs: rem.Milliseconds(),
		Paused:      t.paused,
		Active:      t.active,
		Expired:     t.expired,
	}
}

// Timer devolve o temporizador da sessao, criando sob demanda se necessario.
func (s *Session) Timer() *SessionTimer {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer == nil {
		s.timer = NewSessionTimer(defaultSessionTimeout)
	}
	return s.timer
}

// ConfigureTimer define o timeout por etapa e notifica assinantes SSE.
func (s *Session) ConfigureTimer(d time.Duration) TimerSnapshot {
	snap := s.Timer().SetTimeout(d)
	s.broadcastControl(ControlEvent{Kind: "timer", Timer: &snap})
	return snap
}

// PauseTimer pausa o cronometro regressivo e notifica assinantes SSE.
func (s *Session) PauseTimer() TimerSnapshot {
	snap := s.Timer().Pause()
	s.broadcastControl(ControlEvent{Kind: "timer", Timer: &snap})
	return snap
}

// ResumeTimer retoma o cronometro regressivo e notifica assinantes SSE.
func (s *Session) ResumeTimer() TimerSnapshot {
	snap := s.Timer().Resume()
	s.broadcastControl(ControlEvent{Kind: "timer", Timer: &snap})
	return snap
}

// ExtendTimer adiciona tempo extra ao cronometro e notifica assinantes SSE.
func (s *Session) ExtendTimer(extra time.Duration) TimerSnapshot {
	snap := s.Timer().Extend(extra)
	s.broadcastControl(ControlEvent{Kind: "timer", Timer: &snap})
	return snap
}

// WaitAction aguarda uma Action do usuario respeitando pausa, retomada e extensao.
func (s *Session) WaitAction(ctx context.Context, override time.Duration) (Action, bool) {
	t := s.Timer()
	snap := t.Start(override)
	s.broadcastControl(ControlEvent{Kind: "timer", Timer: &snap})

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case act := <-s.actions:
			stopped := t.Stop()
			s.broadcastControl(ControlEvent{Kind: "timer", Timer: &stopped})
			return act, true
		case <-ctx.Done():
			stopped := t.Stop()
			s.broadcastControl(ControlEvent{Kind: "timer", Timer: &stopped})
			return Action{}, false
		case <-ticker.C:
			cur := t.Snapshot()
			if cur.Expired {
				s.broadcastControl(ControlEvent{Kind: "timer", Timer: &cur})
				return Action{}, false
			}
		case <-t.notifyCh:
			cur := t.Snapshot()
			if cur.Expired {
				s.broadcastControl(ControlEvent{Kind: "timer", Timer: &cur})
				return Action{}, false
			}
		}
	}
}

// PublishStatus atualiza o status vivo e a arvore de etapas da sessao.
func (h *Hub) PublishStatus(sessionID string, st SessionStatus) error {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return err
	}
	if st.AtualizadoEm == "" {
		st.AtualizadoEm = time.Now().UTC().Format(time.RFC3339)
	}
	s.mu.Lock()
	s.status = st
	s.mu.Unlock()
	cp := st
	s.broadcastControl(ControlEvent{Kind: "status", Status: &cp})
	return nil
}

// SubscribeControl anexa um assinante SSE de eventos de timer e status vivo.
func (h *Hub) SubscribeControl(ctx context.Context, sessionID string) (<-chan ControlEvent, error) {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	ch := make(chan ControlEvent, subBufferSize)
	s.attachControl(ch)
	go func() {
		<-ctx.Done()
		s.detachControl(ch)
	}()
	return ch, nil
}

func (s *Session) attachControl(ch chan ControlEvent) {
	snap := s.Timer().Snapshot()
	s.mu.Lock()
	if s.controlSubs == nil {
		s.controlSubs = make(map[chan ControlEvent]struct{})
	}
	st := s.status
	s.controlSubs[ch] = struct{}{}
	s.mu.Unlock()

	select {
	case ch <- ControlEvent{Kind: "timer", Timer: &snap}:
	default:
	}
	if st.Mensagem != "" || len(st.Etapas) > 0 {
		select {
		case ch <- ControlEvent{Kind: "status", Status: &st}:
		default:
		}
	}
}

func (s *Session) detachControl(ch chan ControlEvent) {
	s.mu.Lock()
	if _, ok := s.controlSubs[ch]; ok {
		delete(s.controlSubs, ch)
		close(ch)
	}
	s.mu.Unlock()
}

func (s *Session) broadcastControl(ev ControlEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.controlSubs {
		select {
		case ch <- ev:
		default:
		}
	}
}
