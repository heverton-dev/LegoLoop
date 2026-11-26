package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

const (
	sessaoBaseURLPadrao = "http://127.0.0.1:4242"
	sessaoTimeoutPadrao = 5 * time.Minute
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

// SetSessaoAbridor configura a funcao usada por sessao_criar para abrir a URL no navegador.
func (r *Router) SetSessaoAbridor(fn func(string) error) {
	r.sessaoAbridor = fn
}

func (r *Router) handleSessaoTool(id any, nome string, args map[string]any) (JSONRPCResponse, bool) {
	switch nome {
	case "sessao_criar":
		return r.callSessaoCriar(id, args), true
	case "sessao_publicar":
		return r.callSessaoPublicar(id, args), true
	case "sessao_aguardar":
		return r.callSessaoAguardar(id, args), true
	case "sessao_status":
		return r.callSessaoStatus(id, args), true
	default:
		return JSONRPCResponse{}, false
	}
}

func (r *Router) callSessaoCriar(id any, args map[string]any) JSONRPCResponse {
	if r.sessaoHub == nil {
		return r.toolCallError(id, "hub de sessao A2UI nao configurado")
	}
	processo, _ := args["processo"].(string)
	if processo == "" {
		if s, ok := args["slug"].(string); ok && s != "" {
			processo = s
		} else if m, ok := args["modulo"].(string); ok && m != "" {
			processo = m
		}
	}
	sessao := r.sessaoHub.CreateSession(strings.TrimSpace(processo))
	if override := duracaoTimeoutOpcional(args); override > 0 {
		sessao.ConfigureTimer(override)
	}
	base := r.sessaoBaseURL
	if base == "" {
		base = sessaoBaseURLPadrao
	}
	sessURL := fmt.Sprintf("%s/sessao/%s", base, sessao.ID)
	aberto := r.abrirSessaoSeHabilitado(sessURL, args)
	return r.responderJSONSessao(id, map[string]any{
		"session_id":       sessao.ID,
		"url":              sessURL,
		"stream_url":       fmt.Sprintf("%s/api/sessao/%s/stream", base, sessao.ID),
		"action_url":       fmt.Sprintf("%s/api/sessao/%s/action", base, sessao.ID),
		"navegador_aberto": aberto,
		"timeout_ms":       sessao.Timer().Snapshot().TimeoutMs,
	})
}

func (r *Router) abrirSessaoSeHabilitado(url string, args map[string]any) bool {
	if r.sessaoAbridor == nil {
		return false
	}
	if v, ok := args["abrir_navegador"].(bool); ok && !v {
		return false
	}
	return r.sessaoAbridor(url) == nil
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
	if act, ok := sessao.WaitAction(context.Background(), duracaoTimeoutOpcional(args)); ok {
		return r.responderJSONSessao(id, act)
	}
	return r.responderJSONSessao(id, map[string]any{
		"fallback_textual": true,
		"session_id":       sessao.ID,
		"motivo":           "timeout",
	})
}

func (r *Router) callSessaoStatus(id any, args map[string]any) JSONRPCResponse {
	if r.sessaoHub == nil {
		return r.toolCallError(id, "hub de sessao A2UI nao configurado")
	}
	sessionID, _ := args["session_id"].(string)
	bruto, err := json.Marshal(args)
	if err != nil {
		return r.toolCallError(id, err.Error())
	}
	var st a2ui.SessionStatus
	if err := json.Unmarshal(bruto, &st); err != nil {
		return r.toolCallError(id, err.Error())
	}
	if err := r.sessaoHub.PublishStatus(strings.TrimSpace(sessionID), st); err != nil {
		return r.toolCallError(id, err.Error())
	}
	return r.responderJSONSessao(id, map[string]any{"ok": true, "session_id": sessionID, "etapas": len(st.Etapas)})
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

func duracaoTimeoutOpcional(args map[string]any) time.Duration {
	if v, ok := args["timeout_ms"].(float64); ok && v > 0 {
		return time.Duration(v) * time.Millisecond
	}
	if v, ok := args["timeout_ms"].(int); ok && v > 0 {
		return time.Duration(v) * time.Millisecond
	}
	return 0
}

func duracaoTimeoutSessao(args map[string]any) time.Duration {
	if d := duracaoTimeoutOpcional(args); d > 0 {
		return d
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
		defSessaoStatus(),
	}
}

func defSessaoCriar() map[string]any {
	return map[string]any{
		"name":        "sessao_criar",
		"description": "Cria uma sessao interativa A2UI, abre automaticamente no navegador (padrao) e devolve session_id e url. Aceita timeout_ms por etapa. Em seguida, chame sessao_publicar e sessao_aguardar no mesmo turno.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"processo":        map[string]any{"type": "string", "description": "Nome do Processo associado a sessao (ou slug/modulo)"},
				"slug":            map[string]any{"type": "string", "description": "Alias para processo"},
				"modulo":          map[string]any{"type": "string", "description": "Alias para processo"},
				"abrir_navegador": map[string]any{"type": "boolean", "description": "Abre a URL no navegador automaticamente (padrao true)"},
				"timeout_ms":      map[string]any{"type": "integer", "description": "Tempo limite padrao por etapa em ms (padrao 300000)"},
			},
		},
	}
}

func defSessaoPublicar() map[string]any {
	return map[string]any{
		"name":        "sessao_publicar",
		"description": "Valida e publica um envelope A2UI na sessao interativa. Quando esperar resposta do usuario, chame sessao_aguardar imediatamente no mesmo turno sem encerrar o turno pedindo aviso no chat.",
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
		"description": "Aguarda uma Action do usuario na sessao A2UI respeitando pausa/retomada do cronometro sem exigir mensagem no chat; devolve fallback_textual apenas no timeout.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "ID da sessao A2UI"},
				"timeout_ms": map[string]any{"type": "integer", "description": "Tempo limite em ms (opcional; usa o configurado na sessao)"},
			},
			"required": []string{"session_id"},
		},
	}
}

func defSessaoStatus() map[string]any {
	return map[string]any{
		"name":        "sessao_status",
		"description": "Publica status ao vivo e arvore de etapas na sessao A2UI sem bloquear (use durante processamento e apos Stop-Gates para o usuario acompanhar no navegador).",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "ID da sessao A2UI"},
				"etapa":      map[string]any{"type": "string", "description": "Etapa atual (ex: concepcao, execucao, validacao)"},
				"mensagem":   map[string]any{"type": "string", "description": "Atividade em andamento exibida ao usuario"},
				"etapas":     map[string]any{"type": "array", "description": "Lista de etapas {id, titulo, status: pending|running|done|error, detalhe}"},
			},
			"required": []string{"session_id"},
		},
	}
}
