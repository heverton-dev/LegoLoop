package studio

import (
	"net/http"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
)

func (s *Server) handleA2UIFabricarPage(w http.ResponseWriter, r *http.Request) {
	if err := s.garantirFabricarSessao("sess_2"); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/sessao/sess_2", http.StatusSeeOther)
}

func (s *Server) garantirFabricarSessao(sessID string) error {
	hub := s.A2UIHub()
	if _, err := hub.GetSession(sessID); err == nil {
		return nil
	}
	hub.EnsureSession(sessID, "legoloop-fabricar")
	return s.publicarHierarquiaFabricar(sessID)
}

func (s *Server) publicarHierarquiaFabricar(sessID string) error {
	hub := s.A2UIHub()
	for _, env := range gerarEnvelopesFabricar(sessID) {
		if err := hub.Publish(env); err != nil {
			return err
		}
	}
	return nil
}

func gerarEnvelopesFabricar(sessID string) []a2ui.Envelope {
	return []a2ui.Envelope{
		envelopeFabricarRaiz(sessID),
		envelopeFabricarConcepcao(sessID),
		envelopeFabricarEspecificacao(sessID),
		envelopeFabricarStopGate(sessID),
	}
}

func envelopeFabricarRaiz(sessID string) a2ui.Envelope {
	return a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		FallbackTxt: "Processo Fabricar: orquestração das Etapas industriais.",
		Components: []a2ui.ComponentRef{
			{
				ID:   "fab-progresso-geral",
				Kind: "LegoProgressMeter",
				Props: map[string]any{
					"confianca": 40,
					"meta":      100,
					"label":     "Processo Fabricar (Etapa 1/5 ativa)",
				},
			},
			{
				ID:   "fab-resumo",
				Kind: "Panel",
				Props: map[string]any{
					"title": "Visão Geral do Processo Aninhado",
					"body":  "O Processo Fabricar governa Concepção, Especificação, Planejamento, Execução e Fechamento de forma contínua.",
				},
			},
		},
	}
}

func envelopeFabricarConcepcao(sessID string) a2ui.Envelope {
	return a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "proc:fabricar:concepcao",
		ParentID:    "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		FallbackTxt: "Subprocesso Concepção: levantamento de requisitos e espelho.",
		Components:  compsFabricarConcepcao(),
	}
}

func compsFabricarConcepcao() []a2ui.ComponentRef {
	return []a2ui.ComponentRef{
		{
			ID:   "concepcao-progresso",
			Kind: "LegoProgressMeter",
			Props: map[string]any{
				"confianca": 95,
				"meta":      95,
				"concluido": true,
				"label":     "Subprocesso: Concepção de Requisitos (Meta Atingida)",
			},
		},
		{
			ID:   "concepcao-espelho-1",
			Kind: "Panel",
			Props: map[string]any{
				"title": "Resultado Esperado",
				"body":  "Interface dinâmica A2UI com navegação em árvore para macroprocessos e subprocessos aninhados.",
			},
		},
		{
			ID:   "concepcao-espelho-2",
			Kind: "Panel",
			Props: map[string]any{
				"title": "Critérios de Sucesso",
				"body":  "Alternância imediata entre nós, sincronização em tempo real e foco automático em decisões humanas.",
			},
		},
	}
}

func envelopeFabricarEspecificacao(sessID string) a2ui.Envelope {
	return a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "proc:fabricar:especificacao",
		ParentID:    "proc:fabricar",
		ProtocolVer: "v0.9",
		MessageType: "createSurface",
		FallbackTxt: "Subprocesso Especificação: critérios EARS e contratos.",
		Components: []a2ui.ComponentRef{
			{
				ID:   "especificacao-progresso",
				Kind: "LegoProgressMeter",
				Props: map[string]any{
					"confianca": 10,
					"meta":      100,
					"label":     "Subprocesso: Especificação Técnica (Aguardando Concepção)",
				},
			},
			{
				ID:   "especificacao-ears",
				Kind: "Panel",
				Props: map[string]any{
					"title": "Critério de Aceite EARS",
					"body":  "Quando o operador selecionar um nó na árvore, a prancheta central deve renderizar exclusivamente as Peças daquela superfície.",
				},
			},
		},
	}
}

func envelopeFabricarStopGate(sessID string) a2ui.Envelope {
	return a2ui.Envelope{
		SessionID:         sessID,
		SurfaceID:         "proc:fabricar:concepcao:stop-gate-01",
		ParentID:          "proc:fabricar:concepcao",
		ProtocolVer:       "v0.9",
		MessageType:       "createSurface",
		RequiresAttention: true,
		FallbackTxt:       "Stop-Gate 01: Decisão humana necessária para avançar da Concepção para a Especificação.",
		Components: []a2ui.ComponentRef{
			{
				ID:   "gate-01-panel",
				Kind: "LegoStopGatePanel",
				Props: map[string]any{
					"title":   "Stop-Gate 01 • Aprovação da INTENT do Fabricar",
					"summary": "A concepção atingiu 95% de assertividade. Autoriza a transição para a Especificação do macroprocesso?",
				},
			},
		},
	}
}
