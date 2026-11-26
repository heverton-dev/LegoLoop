package a2ui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	subBufferSize    = 32
	actionBufferSize = 16
)

// SessionState define o estado do ciclo de vida da sessão.
type SessionState string

const (
	SessionRunning   SessionState = "running"
	SessionPaused    SessionState = "paused"
	SessionCancelled SessionState = "cancelled"
)

// SessionSummary resume o estado operacional da sessão para a gestão de processos.
type SessionSummary struct {
	ID          string       `json:"id"`
	Processo    string       `json:"processo"`
	State       SessionState `json:"state"`
	Surfaces    int          `json:"surfaces"`
	Subscribers int          `json:"subscribers"`
}

// Session guarda o estado das superfícies, assinantes SSE e canal de ações de uma execução.
type Session struct {
	ID       string
	Processo string
	State    SessionState

	mu          sync.RWMutex
	surfaces    map[string]Envelope
	order       []string
	subscribers map[chan Envelope]struct{}
	controlSubs map[chan ControlEvent]struct{}
	actions     chan Action
	timer       *SessionTimer
	status      SessionStatus
}

// Actions devolve o canal de leitura para receber ações do usuário.
func (s *Session) Actions() <-chan Action {
	return s.actions
}

// GetState devolve o estado atual da sessão sob lock de leitura.
func (s *Session) GetState() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.State == "" {
		return SessionRunning
	}
	return s.State
}

// SetState atualiza o estado operacional da sessão.
func (s *Session) SetState(st SessionState) {
	s.mu.Lock()
	s.State = st
	s.mu.Unlock()
}

// Hub gerencia o ciclo de vida das sessões A2UI em memória.
type Hub struct {
	catalog  Catalog
	seq      atomic.Uint64
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewHub cria um barramento de sessões validando envelopes contra cat.
func NewHub(cat Catalog) *Hub {
	return &Hub{
		catalog:  cat,
		sessions: make(map[string]*Session),
	}
}

// ListSessions devolve a lista de resumos de todas as sessões registradas.
func (h *Hub) ListSessions() []SessionSummary {
	h.mu.RLock()
	defer h.mu.RUnlock()
	res := make([]SessionSummary, 0, len(h.sessions))
	for _, s := range h.sessions {
		s.mu.RLock()
		st := s.State
		if st == "" {
			st = SessionRunning
		}
		res = append(res, SessionSummary{
			ID:          s.ID,
			Processo:    s.Processo,
			State:       st,
			Surfaces:    len(s.surfaces),
			Subscribers: len(s.subscribers),
		})
		s.mu.RUnlock()
	}
	return res
}

// CreateSession registra uma nova sessão interativa e devolve seu ponteiro.
func (h *Hub) CreateSession(processo string) *Session {
	id := fmt.Sprintf("sess_%d", h.seq.Add(1))
	s := &Session{
		ID:          id,
		Processo:    processo,
		State:       SessionRunning,
		surfaces:    make(map[string]Envelope),
		subscribers: make(map[chan Envelope]struct{}),
		controlSubs: make(map[chan ControlEvent]struct{}),
		actions:     make(chan Action, actionBufferSize),
		timer:       NewSessionTimer(defaultSessionTimeout),
	}
	h.mu.Lock()
	h.sessions[id] = s
	h.mu.Unlock()
	return s
}

// EnsureSession devolve a sessão com o ID especificado ou a cria se ainda não existir.
func (h *Hub) EnsureSession(id, processo string) *Session {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s, ok := h.sessions[id]; ok {
		return s
	}
	s := &Session{
		ID:          id,
		Processo:    processo,
		State:       SessionRunning,
		surfaces:    make(map[string]Envelope),
		subscribers: make(map[chan Envelope]struct{}),
		controlSubs: make(map[chan ControlEvent]struct{}),
		actions:     make(chan Action, actionBufferSize),
		timer:       NewSessionTimer(defaultSessionTimeout),
	}
	h.sessions[id] = s
	if strings.HasPrefix(id, "sess_") {
		var n uint64
		if _, err := fmt.Sscanf(id, "sess_%d", &n); err == nil {
			for {
				cur := h.seq.Load()
				if n <= cur || h.seq.CompareAndSwap(cur, n) {
					break
				}
			}
		}
	}
	return s
}

// GetSession retorna a sessão pelo ID ou ErrSessaoNaoEncontrada.
func (h *Hub) GetSession(id string) (*Session, error) {
	h.mu.RLock()
	s, ok := h.sessions[id]
	h.mu.RUnlock()
	if !ok {
		return nil, ErrSessaoNaoEncontrada
	}
	return s, nil
}

// Publish valida o envelope, atualiza o snapshot da superfície e transmite aos assinantes.
func (h *Hub) Publish(env Envelope) error {
	if err := env.Validate(h.catalog); err != nil {
		return err
	}
	s, err := h.GetSession(env.SessionID)
	if err != nil {
		return err
	}
	if containsStopGate(env.Components) {
		env.RequiresAttention = true
	}
	s.applyAndBroadcast(env)
	if !strings.HasPrefix(env.SurfaceID, "andon:") {
		s.mu.RLock()
		estavaAguardando := s.status.AguardandoResposta
		s.mu.RUnlock()
		if estavaAguardando {
			return h.PublishStatus(env.SessionID, SessionStatus{Etapa: "interacao", AguardandoResposta: false})
		}
	}
	return nil
}

// Subscribe anexa um assinante SSE à sessão, enviando o snapshot atual antes de eventos novos.
func (h *Hub) Subscribe(ctx context.Context, sessionID string) (<-chan Envelope, error) {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	ch := make(chan Envelope, subBufferSize)
	s.attachSubscriber(ch)
	go func() {
		<-ctx.Done()
		s.detachSubscriber(ch)
	}()
	return ch, nil
}

// SubmitAction valida e entrega a ação do usuário no canal da sessão.
func (h *Hub) SubmitAction(act Action) error {
	if err := act.Validate(); err != nil {
		return err
	}
	s, err := h.GetSession(act.SessionID)
	if err != nil {
		return err
	}
	select {
	case s.actions <- act:
		s.clearAttention(act.SurfaceID)
		if !strings.HasPrefix(act.Name, "andon_cord_") {
			return h.PublishStatus(act.SessionID, SessionStatus{
				Etapa:              "processando",
				Mensagem:           "Resposta enviada (" + act.Name + "). Processo em execução, aguarde a próxima etapa...",
				AguardandoResposta: true,
			})
		}
		return nil
	default:
		return ErrAcaoInvalida
	}
}

// SubscriberCount devolve o número de assinantes ativos na sessão.
func (h *Hub) SubscriberCount(sessionID string) int {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

func (s *Session) applyAndBroadcast(env Envelope) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mergeSurfaceLocked(env)
	for ch := range s.subscribers {
		select {
		case ch <- env:
		default:
		}
	}
}

func (s *Session) mergeSurfaceLocked(env Envelope) {
	if env.MessageType == "deleteSurface" {
		delete(s.surfaces, env.SurfaceID)
		return
	}
	cur, existe := s.surfaces[env.SurfaceID]
	if !existe {
		s.order = append(s.order, env.SurfaceID)
		s.surfaces[env.SurfaceID] = cloneEnvelope(env)
		return
	}
	if len(env.Components) > 0 {
		cur.Components = append([]ComponentRef(nil), env.Components...)
	}
	cur.DataModel = mergeMaps(cur.DataModel, env.DataModel)
	if env.FallbackTxt != "" {
		cur.FallbackTxt = env.FallbackTxt
	}
	if env.ParentID != "" {
		cur.ParentID = env.ParentID
	}
	if env.RequiresAttention {
		cur.RequiresAttention = true
	}
	cur.MessageType = env.MessageType
	s.surfaces[env.SurfaceID] = cur
}

func (s *Session) attachSubscriber(ch chan Envelope) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range s.order {
		if env, ok := s.surfaces[id]; ok {
			select {
			case ch <- cloneEnvelope(env):
			default:
			}
		}
	}
	s.subscribers[ch] = struct{}{}
}

func (s *Session) detachSubscriber(ch chan Envelope) {
	s.mu.Lock()
	if _, ok := s.subscribers[ch]; ok {
		delete(s.subscribers, ch)
		close(ch)
	}
	s.mu.Unlock()
}

func cloneEnvelope(env Envelope) Envelope {
	cp := env
	if len(env.Components) > 0 {
		cp.Components = append([]ComponentRef(nil), env.Components...)
	}
	if len(env.DataModel) > 0 {
		cp.DataModel = mergeMaps(nil, env.DataModel)
	}
	return cp
}

func mergeMaps(dst, src map[string]any) map[string]any {
	if len(dst) == 0 && len(src) == 0 {
		return nil
	}
	out := make(map[string]any, len(dst)+len(src))
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range src {
		out[k] = v
	}
	return out
}
