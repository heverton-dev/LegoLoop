package a2ui

import (
	"fmt"
	"strings"
)

const (
	// DefaultSwarmSurfaceID é o identificador padrão da superfície do painel Swarm.
	DefaultSwarmSurfaceID = "proc:swarm"
	swarmProtocolVersion  = "v0.9"
	swarmPostActionName   = "swarm_post"
)

// SwarmWorker representa um agente listado no quadro do Swarm.
type SwarmWorker struct {
	Name  string `json:"name"`
	Role  string `json:"role,omitempty"`
	State string `json:"state"`
	Model string `json:"model,omitempty"`
}

// SwarmMessage representa uma mensagem publicada no quadro do Swarm.
type SwarmMessage struct {
	Seq   int    `json:"seq"`
	Agent string `json:"agent"`
	Text  string `json:"text"`
	Time  string `json:"time,omitempty"`
}

// SwarmState concentra o estado atual do Swarm para renderização em A2UI.
type SwarmState struct {
	SessionID string         `json:"session_id"`
	SurfaceID string         `json:"surface_id"`
	Title     string         `json:"title,omitempty"`
	Workers   []SwarmWorker  `json:"workers,omitempty"`
	Messages  []SwarmMessage `json:"messages,omitempty"`
	Progress  int            `json:"progress"`
}

// BuildSwarmEnvelopes gera os envelopes A2UI (createSurface e updateDataModel) do Swarm.
func BuildSwarmEnvelopes(st SwarmState) []Envelope {
	surfaceID := strings.TrimSpace(st.SurfaceID)
	if surfaceID == "" {
		surfaceID = DefaultSwarmSurfaceID
	}
	progress := clampProgress(st.Progress)
	fallback := buildSwarmFallback(st, progress)
	comps := buildSwarmComponents(st, progress)
	data := buildSwarmDataModel(st, progress)

	return []Envelope{
		{
			SessionID:   st.SessionID,
			SurfaceID:   surfaceID,
			ProtocolVer: swarmProtocolVersion,
			MessageType: "createSurface",
			Components:  comps,
			FallbackTxt: fallback,
		},
		{
			SessionID:   st.SessionID,
			SurfaceID:   surfaceID,
			ProtocolVer: swarmProtocolVersion,
			MessageType: "updateDataModel",
			DataModel:   data,
			FallbackTxt: fallback,
		},
	}
}

// PublishSwarm valida e publica os envelopes do Swarm no Hub.
func PublishSwarm(hub *Hub, st SwarmState) error {
	if hub == nil {
		return ErrSessaoNaoEncontrada
	}
	for _, env := range BuildSwarmEnvelopes(st) {
		if err := hub.Publish(env); err != nil {
			return err
		}
	}
	return nil
}

// ParseSwarmPostAction valida uma ação swarm_post e devolve o texto sanitizado.
func ParseSwarmPostAction(act Action) (string, error) {
	if act.Name != swarmPostActionName || act.Payload == nil {
		return "", ErrAcaoInvalida
	}
	raw, ok := act.Payload["text"].(string)
	if !ok {
		return "", ErrAcaoInvalida
	}
	text := strings.TrimSpace(raw)
	if text == "" {
		return "", ErrAcaoInvalida
	}
	if contemHostil(act.Payload) {
		return "", ErrPayloadHostil
	}
	return text, nil
}

func buildSwarmComponents(st SwarmState, progress int) []ComponentRef {
	title := strings.TrimSpace(st.Title)
	if title == "" {
		title = "Swarm"
	}
	return []ComponentRef{
		{
			ID:   "swarm-progresso",
			Kind: "LegoProgressMeter",
			Props: map[string]any{
				"title":    title,
				"progress": progress,
				"meta":     100,
			},
		},
		{
			ID:   "swarm-board",
			Kind: "LegoSwarmBoard",
			Props: map[string]any{
				"title":    title,
				"action":   swarmPostActionName,
				"workers":  swarmWorkersToAny(st.Workers),
				"messages": swarmMessagesToAny(st.Messages),
				"progress": progress,
			},
		},
	}
}

func buildSwarmDataModel(st SwarmState, progress int) map[string]any {
	return map[string]any{
		"title":    st.Title,
		"progress": progress,
		"workers":  swarmWorkersToAny(st.Workers),
		"messages": swarmMessagesToAny(st.Messages),
	}
}

func buildSwarmFallback(st SwarmState, progress int) string {
	title := strings.TrimSpace(st.Title)
	if title == "" {
		title = "Swarm"
	}
	linhas := []string{
		fmt.Sprintf("%s: progresso %d%% (%d workers, %d mensagens).", title, progress, len(st.Workers), len(st.Messages)),
	}
	for _, w := range st.Workers {
		linhas = append(linhas, fmt.Sprintf("- worker %s [%s]", w.Name, w.State))
	}
	if len(st.Messages) > 0 {
		ult := st.Messages[len(st.Messages)-1]
		linhas = append(linhas, fmt.Sprintf("Última mensagem (#%d %s): %s", ult.Seq, ult.Agent, ult.Text))
	}
	return strings.Join(linhas, "\n")
}

func swarmWorkersToAny(workers []SwarmWorker) []any {
	out := make([]any, 0, len(workers))
	for _, w := range workers {
		out = append(out, map[string]any{
			"name":  w.Name,
			"role":  w.Role,
			"state": w.State,
			"model": w.Model,
		})
	}
	return out
}

func swarmMessagesToAny(messages []SwarmMessage) []any {
	out := make([]any, 0, len(messages))
	for _, m := range messages {
		out = append(out, map[string]any{
			"seq":   m.Seq,
			"agent": m.Agent,
			"text":  m.Text,
			"time":  m.Time,
		})
	}
	return out
}

func clampProgress(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
