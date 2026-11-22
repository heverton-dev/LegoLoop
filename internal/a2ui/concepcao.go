package a2ui

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// DefaultConcepcaoSurfaceID é o identificador padrão da superfície do piloto de concepção.
	DefaultConcepcaoSurfaceID = "proc:concepcao"
	defaultConcepcaoMeta      = 95
	concepcaoProtocolVersion  = "v0.9"
)

// ConcepcaoState representa o estado de um turno do Processo legoloop-concepcao.
type ConcepcaoState struct {
	SessionID          string            `json:"session_id"`
	SurfaceID          string            `json:"surface_id"`
	Confianca          int               `json:"confianca"`
	Meta               int               `json:"meta"`
	Turno              int               `json:"turno"`
	Pergunta           string            `json:"pergunta,omitempty"`
	Palpite            string            `json:"palpite,omitempty"`
	Lacunas            []string          `json:"lacunas,omitempty"`
	Espelho            map[string]string `json:"espelho,omitempty"`
	IntentPath         string            `json:"intent_path,omitempty"`
	AguardandoStopGate bool              `json:"aguardando_stop_gate,omitempty"`
}

// BuildConcepcaoEnvelopes traduz o estado da concepção em envelopes A2UI validáveis.
func BuildConcepcaoEnvelopes(st ConcepcaoState) []Envelope {
	surfaceID := strings.TrimSpace(st.SurfaceID)
	if surfaceID == "" {
		surfaceID = DefaultConcepcaoSurfaceID
	}
	meta := st.Meta
	if meta <= 0 {
		meta = defaultConcepcaoMeta
	}
	fallback := buildConcepcaoFallback(st, meta)
	comps := buildConcepcaoComponents(st, meta)
	data := buildConcepcaoDataModel(st, meta)

	return []Envelope{
		{
			SessionID:         st.SessionID,
			SurfaceID:         surfaceID,
			ProtocolVer:       concepcaoProtocolVersion,
			MessageType:       "createSurface",
			Components:        comps,
			FallbackTxt:       fallback,
			RequiresAttention: st.AguardandoStopGate,
		},
		{
			SessionID:         st.SessionID,
			SurfaceID:         surfaceID,
			ProtocolVer:       concepcaoProtocolVersion,
			MessageType:       "updateDataModel",
			DataModel:         data,
			FallbackTxt:       fallback,
			RequiresAttention: st.AguardandoStopGate,
		},
	}
}

// PublishConcepcao valida e publica o estado da concepção no Hub sem exigir assinantes SSE.
func PublishConcepcao(hub *Hub, st ConcepcaoState) error {
	if hub == nil {
		return ErrSessaoNaoEncontrada
	}
	for _, env := range BuildConcepcaoEnvelopes(st) {
		if err := hub.Publish(env); err != nil {
			return err
		}
	}
	return nil
}

func buildConcepcaoComponents(st ConcepcaoState, meta int) []ComponentRef {
	comps := []ComponentRef{
		{
			ID:   "concepcao-progresso",
			Kind: "LegoProgressMeter",
			Props: map[string]any{
				"confianca": st.Confianca,
				"meta":      meta,
				"turno":     st.Turno,
				"lacunas":   append([]string(nil), st.Lacunas...),
			},
		},
	}
	comps = append(comps, buildPerguntaComponents(st)...)
	comps = append(comps, buildEspelhoCards(st.Espelho)...)
	if st.AguardandoStopGate {
		comps = append(comps, buildStopGateComponent())
	}
	if strings.TrimSpace(st.IntentPath) != "" {
		comps = append(comps, ComponentRef{
			ID:   "concepcao-intent",
			Kind: "LegoArtifactViewer",
			Props: map[string]any{
				"title": "INTENT.md",
				"path":  st.IntentPath,
			},
		})
	}
	return comps
}

func buildPerguntaComponents(st ConcepcaoState) []ComponentRef {
	if strings.TrimSpace(st.Pergunta) == "" && strings.TrimSpace(st.Palpite) == "" {
		return nil
	}
	return []ComponentRef{
		{
			ID:   "concepcao-resposta",
			Kind: "TextField",
			Props: map[string]any{
				"label":       st.Pergunta,
				"value":       st.Palpite,
				"placeholder": "Edite o palpite ou escreva sua resposta",
			},
		},
		{
			ID:   "concepcao-enviar",
			Kind: "Button",
			Props: map[string]any{
				"label":  "Enviar resposta",
				"action": "submit_answer",
			},
		},
	}
}

func buildEspelhoCards(espelho map[string]string) []ComponentRef {
	if len(espelho) == 0 {
		return nil
	}
	chaves := make([]string, 0, len(espelho))
	for k := range espelho {
		chaves = append(chaves, k)
	}
	sort.Strings(chaves)

	cards := make([]ComponentRef, 0, len(chaves))
	for _, k := range chaves {
		cards = append(cards, ComponentRef{
			ID:   "concepcao-espelho-" + k,
			Kind: "Card",
			Props: map[string]any{
				"title": k,
				"body":  espelho[k],
			},
		})
	}
	return cards
}

func buildStopGateComponent() ComponentRef {
	return ComponentRef{
		ID:   "concepcao-stop-gate-01",
		Kind: "LegoStopGateCard",
		Props: map[string]any{
			"title":   "Stop-Gate 01 — Aprovação da INTENT",
			"action":  "stop_gate_decision",
			"options": []string{"[1] seguir", "[2] personalizar", "[0] cancelar"},
		},
	}
}

func buildConcepcaoDataModel(st ConcepcaoState, meta int) map[string]any {
	espelho := make(map[string]any, len(st.Espelho))
	for k, v := range st.Espelho {
		espelho[k] = v
	}
	return map[string]any{
		"confianca":            st.Confianca,
		"meta":                 meta,
		"turno":                st.Turno,
		"pergunta":             st.Pergunta,
		"palpite":              st.Palpite,
		"lacunas":              append([]string(nil), st.Lacunas...),
		"espelho":              espelho,
		"intent_path":          st.IntentPath,
		"aguardando_stop_gate": st.AguardandoStopGate,
	}
}

func buildConcepcaoFallback(st ConcepcaoState, meta int) string {
	var linhas []string
	linhas = append(linhas, fmt.Sprintf("Concepção (turno %d): confiança %d%% / meta %d%%.", st.Turno, st.Confianca, meta))
	if len(st.Lacunas) > 0 {
		linhas = append(linhas, "Lacunas: "+strings.Join(st.Lacunas, ", ")+".")
	}
	if strings.TrimSpace(st.Pergunta) != "" {
		linhas = append(linhas, "Pergunta: "+st.Pergunta)
	}
	if strings.TrimSpace(st.Palpite) != "" {
		linhas = append(linhas, "Palpite: "+st.Palpite)
	}
	if len(st.Espelho) > 0 {
		linhas = append(linhas, fmt.Sprintf("Espelho Canônico com %d itens preenchidos.", len(st.Espelho)))
	}
	if strings.TrimSpace(st.IntentPath) != "" {
		linhas = append(linhas, "INTENT: "+st.IntentPath)
	}
	if st.AguardandoStopGate {
		linhas = append(linhas, "Stop-Gate 01: [1] seguir | [2] personalizar | [0] cancelar")
	}
	return strings.Join(linhas, "\n")
}
