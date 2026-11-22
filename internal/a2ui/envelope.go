package a2ui

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrSessaoNaoEncontrada = errors.New("a2ui: sessão não encontrada")
	ErrSessaoInvalida      = errors.New("a2ui: session_id inválido")
	ErrSurfaceIDInvalido   = errors.New("a2ui: surface_id inválido")
	ErrMensagemInvalida    = errors.New("a2ui: message_type inválido")
	ErrItemProibido        = errors.New("a2ui: item fora do catálogo confiável")
	ErrPayloadHostil       = errors.New("a2ui: conteúdo executável ou HTML bruto proibido")
	ErrAcaoInvalida        = errors.New("a2ui: ação de usuário inválida")
)

var (
	surfaceIDPattern = regexp.MustCompile(`^[a-zA-Z0-9._:/-]{1,128}$`)
	hostilPattern    = regexp.MustCompile(`(?i)(<\s*script|<\s*iframe|<\s*svg|<\s*img|javascript\s*:|data\s*:\s*text/html|on[a-z]+\s*=)`)
)

func mensagemValida(t string) bool {
	switch t {
	case "createSurface", "updateComponents", "updateDataModel", "deleteSurface":
		return true
	default:
		return false
	}
}

// Envelope encapsula uma mensagem do protocolo A2UI vinculada a uma superfície isolada.
type Envelope struct {
	SessionID         string         `json:"session_id"`
	SurfaceID         string         `json:"surface_id"`
	ParentID          string         `json:"parent_surface_id,omitempty"`
	ProtocolVer       string         `json:"protocol_version"`
	MessageType       string         `json:"message_type"`
	Components        []ComponentRef `json:"components,omitempty"`
	DataModel         map[string]any `json:"data_model,omitempty"`
	FallbackTxt       string         `json:"fallback_text,omitempty"`
	RequiresAttention bool           `json:"requires_attention,omitempty"`
}

// ComponentRef referencia um item declarado no catálogo confiável.
type ComponentRef struct {
	ID    string         `json:"id"`
	Kind  string         `json:"kind"`
	Props map[string]any `json:"props,omitempty"`
}

// Action representa uma resposta do usuário enviada pela UI, stdin ou MCP.
type Action struct {
	SessionID string         `json:"session_id"`
	SurfaceID string         `json:"surface_id"`
	Origin    string         `json:"origin"`
	Name      string         `json:"name"`
	Payload   map[string]any `json:"payload,omitempty"`
}

// Validate valida todos os campos do envelope antes de qualquer emissão SSE.
func (e Envelope) Validate(cat Catalog) error {
	if strings.TrimSpace(e.SessionID) == "" {
		return ErrSessaoInvalida
	}
	if !validSurfaceID(e.SurfaceID) || (e.ParentID != "" && !validSurfaceID(e.ParentID)) {
		return ErrSurfaceIDInvalido
	}
	if !mensagemValida(e.MessageType) {
		return ErrMensagemInvalida
	}
	for _, c := range e.Components {
		if !cat.Allows(c.Kind) {
			return ErrItemProibido
		}
		if contemHostil(c.Props) {
			return ErrPayloadHostil
		}
	}
	if contemHostil(e.DataModel) {
		return ErrPayloadHostil
	}
	return nil
}

// Validate checa os campos mínimos de uma ação recebida do usuário.
func (a Action) Validate() error {
	if strings.TrimSpace(a.SessionID) == "" || strings.TrimSpace(a.Name) == "" {
		return ErrAcaoInvalida
	}
	if a.SurfaceID != "" && !validSurfaceID(a.SurfaceID) {
		return ErrSurfaceIDInvalido
	}
	return nil
}

func validSurfaceID(id string) bool {
	if !surfaceIDPattern.MatchString(id) {
		return false
	}
	return !strings.Contains(id, "..")
}

func contemHostil(v any) bool {
	switch x := v.(type) {
	case string:
		return hostilPattern.MatchString(x)
	case map[string]any:
		for _, val := range x {
			if contemHostil(val) {
				return true
			}
		}
	case []any:
		for _, val := range x {
			if contemHostil(val) {
				return true
			}
		}
	}
	return false
}
