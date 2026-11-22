package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

const (
	sessaoBaseURLPadrao = "http://127.0.0.1:4242"
	sessaoTimeoutPadrao = 30 * time.Second
)

// SetSessaoHub conecta um barramento de sessoes A2UI e sua URL base ao Router.
func (r *Router) SetSessaoHub(hub *a2ui.Hub, baseURL string) {
	r.sessaoHub = hub
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = sessaoBaseURLPadrao
	}
	r.sessaoBaseURL = base
}

func (r *Router) handleSessaoTool(id any, nome string, args map[string]any) (JSONRPCResponse, bool) {
	switch nome {
	case "sessao_criar":
		return r.callSessaoCriar(id, args), true
	case "sessao_publicar":
		return r.callSessaoPublicar(id, args), true
	case "sessao_aguardar":
		return r.callSessaoAguardar(id, args), true
	default:
		return JSONRPCResponse{}, false
	}
}

func (r *Router) callSessaoCriar(id any, args map[string]any) JSONRPCResponse {
	if r.sessaoHub == nil {
		return r.toolCallError(id, "hub de sessao A2UI nao configurado")
	}
	processo, _ := args["processo"].(string)
	sessao := r.sessaoHub.CreateSession(strings.TrimSpace(processo))
	base := r.sessaoBaseURL
	if base == "" {
		base = sessaoBaseURLPadrao
	}
	return r.responderJSONSessao(id, map[string]any{
		"session_id": sessao.ID,
		"url":        fmt.Sprintf("%s/sessao/%s", base, sessao.ID),
		"stream_url": fmt.Sprintf("%s/api/sessao/%s/stream", base, sessao.ID),
		"action_url": fmt.Sprintf("%s/api/sessao/%s/action", base, sessao.ID),
	})
}

func (r *Router) callSessaoPublicar(id any, args map[string]any) JSONRPCResponse {
	if r.sessaoHub == nil {
		return r.toolCallError(id, "hub de sessao A2UI nao configurado")
	}
	env, err := decodificarEnvelopeSessao(args)
	if err != nil {
		return r.toolCallError(id, fmt.Sprintf("envelope A2UI invalido: %v", err))
	}
	if err := r.sessaoHub.Publish(env); err != nil {
		return r.toolCallError(id, err.Error())
	}
	return r.responderJSONSessao(id, map[string]any{
		"ok":         true,
		"session_id": env.SessionID,
		"surface_id": env.SurfaceID,
	})
}

func (r *Router) callSessaoAguardar(id any, args map[string]any) JSONRPCResponse {
	if r.sessaoHub == nil {
		return r.toolCallError(id, "hub de sessao A2UI nao configurado")
	}
	sessionID, _ := args["session_id"].(string)
	sessao, err := r.sessaoHub.GetSession(strings.TrimSpace(sessionID))
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	timer := time.NewTimer(duracaoTimeoutSessao(args))
	defer timer.Stop()
	select {
	case act := <-sessao.Actions():
		return r.responderJSONSessao(id, act)
	case <-timer.C:
		return r.responderJSONSessao(id, map[string]any{
			"fallback_textual": true,
			"session_id":       sessao.ID,
			"motivo":           "timeout",
		})
	}
}

func decodificarEnvelopeSessao(args map[string]any) (a2ui.Envelope, error) {
	origem := args
	if aninhado, ok := args["envelope"].(map[string]any); ok {
		origem = aninhado
	}
	bruto, err := json.Marshal(origem)
	if err != nil {
		return a2ui.Envelope{}, err
	}
	var env a2ui.Envelope
	if err := json.Unmarshal(bruto, &env); err != nil {
		return a2ui.Envelope{}, err
	}
	return env, nil
}

func duracaoTimeoutSessao(args map[string]any) time.Duration {
	if v, ok := args["timeout_ms"].(float64); ok && v > 0 {
		return time.Duration(v) * time.Millisecond
	}
	if v, ok := args["timeout_ms"].(int); ok && v > 0 {
		return time.Duration(v) * time.Millisecond
	}
	return sessaoTimeoutPadrao
}

func (r *Router) responderJSONSessao(id any, valor any) JSONRPCResponse {
	data, err := json.Marshal(valor)
	if err != nil {
		return r.toolCallError(id, fmt.Sprintf("falha ao serializar resposta: %v", err))
	}
	return r.toolCallText(id, string(data))
}

func sessaoToolDefs() []map[string]any {
	return []map[string]any{
		defSessaoCriar(),
		defSessaoPublicar(),
		defSessaoAguardar(),
	}
}

func defSessaoCriar() map[string]any {
	return map[string]any{
		"name":        "sessao_criar",
		"description": "Cria uma sessao interativa A2UI e devolve session_id e url.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"processo": map[string]any{"type": "string", "description": "Nome do Processo associado a sessao"},
			},
		},
	}
}

func defSessaoPublicar() map[string]any {
	return map[string]any{
		"name":        "sessao_publicar",
		"description": "Valida e publica um envelope A2UI na sessao interativa.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id":       map[string]any{"type": "string", "description": "ID da sessao A2UI"},
				"surface_id":       map[string]any{"type": "string", "description": "ID hierarquico da superficie"},
				"protocol_version": map[string]any{"type": "string", "description": "Versao do protocolo A2UI"},
				"message_type":     map[string]any{"type": "string", "description": "Tipo de mensagem A2UI"},
				"components":       map[string]any{"type": "array", "description": "Lista de items do catalogo"},
				"data_model":       map[string]any{"type": "object", "description": "Estado de dados da superficie"},
				"fallback_text":    map[string]any{"type": "string", "description": "Texto de fallback"},
			},
			"required": []string{"session_id", "surface_id"},
		},
	}
}

func defSessaoAguardar() map[string]any {
	return map[string]any{
		"name":        "sessao_aguardar",
		"description": "Aguarda uma Action do usuario na sessao A2UI ou devolve fallback_textual no timeout.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "ID da sessao A2UI"},
				"timeout_ms": map[string]any{"type": "integer", "description": "Tempo limite em milissegundos"},
			},
			"required": []string{"session_id"},
		},
	}
}
