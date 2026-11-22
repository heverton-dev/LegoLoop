package a2ui

const stopGateComponentKind = "LegoStopGateCard"

// SurfaceNode representa um nó na árvore hierárquica de superfícies de uma sessão.
type SurfaceNode struct {
	SurfaceID         string        `json:"surface_id"`
	ParentID          string        `json:"parent_surface_id,omitempty"`
	RequiresAttention bool          `json:"requires_attention,omitempty"`
	Children          []SurfaceNode `json:"children,omitempty"`
}

// Surface devolve uma cópia isolada do estado atual de uma superfície da sessão.
func (h *Hub) Surface(sessionID, surfaceID string) (Envelope, bool) {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return Envelope{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	env, ok := s.surfaces[surfaceID]
	if !ok {
		return Envelope{}, false
	}
	return cloneEnvelope(env), true
}

// Tree devolve a árvore hierárquica pai->filhos das superfícies ativas da sessão.
func (h *Hub) Tree(sessionID string) ([]SurfaceNode, error) {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.buildTreeLocked(), nil
}

// Breadcrumb devolve o caminho ordenado da raiz até surfaceID na sessão.
func (h *Hub) Breadcrumb(sessionID, surfaceID string) ([]string, error) {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.surfaces[surfaceID]; !ok {
		return nil, ErrSurfaceIDInvalido
	}
	var rev []string
	seen := make(map[string]bool)
	for cur := surfaceID; cur != "" && !seen[cur]; {
		seen[cur] = true
		rev = append(rev, cur)
		env, ok := s.surfaces[cur]
		if !ok {
			break
		}
		cur = env.ParentID
	}
	out := make([]string, len(rev))
	for i := range rev {
		out[i] = rev[len(rev)-1-i]
	}
	return out, nil
}

// PendingAttention devolve o primeiro surfaceID que aguarda atenção de Stop-Gate.
func (h *Hub) PendingAttention(sessionID string) string {
	s, err := h.GetSession(sessionID)
	if err != nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, id := range s.order {
		if env, ok := s.surfaces[id]; ok && env.RequiresAttention {
			return id
		}
	}
	return ""
}

func (s *Session) buildTreeLocked() []SurfaceNode {
	childrenByParent := make(map[string][]string)
	var roots []string
	for _, id := range s.order {
		env, ok := s.surfaces[id]
		if !ok {
			continue
		}
		if env.ParentID == "" || !s.hasSurfaceLocked(env.ParentID) {
			roots = append(roots, id)
			continue
		}
		childrenByParent[env.ParentID] = append(childrenByParent[env.ParentID], id)
	}
	visited := make(map[string]bool)
	out := make([]SurfaceNode, 0, len(roots))
	for _, rootID := range roots {
		out = append(out, s.buildNodeLocked(rootID, childrenByParent, visited))
	}
	return out
}

func (s *Session) buildNodeLocked(id string, children map[string][]string, visited map[string]bool) SurfaceNode {
	env := s.surfaces[id]
	node := SurfaceNode{
		SurfaceID:         env.SurfaceID,
		ParentID:          env.ParentID,
		RequiresAttention: env.RequiresAttention,
	}
	if visited[id] {
		return node
	}
	visited[id] = true
	for _, childID := range children[id] {
		node.Children = append(node.Children, s.buildNodeLocked(childID, children, visited))
	}
	return node
}

func (s *Session) hasSurfaceLocked(id string) bool {
	_, ok := s.surfaces[id]
	return ok
}

func (s *Session) clearAttention(surfaceID string) {
	if surfaceID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if env, ok := s.surfaces[surfaceID]; ok && env.RequiresAttention {
		env.RequiresAttention = false
		s.surfaces[surfaceID] = env
	}
}

func containsStopGate(comps []ComponentRef) bool {
	for _, c := range comps {
		if c.Kind == stopGateComponentKind {
			return true
		}
	}
	return false
}
