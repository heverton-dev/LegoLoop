package studio

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/interview"
)

func (s *Server) handleA2UIConcepcaoPage(w http.ResponseWriter, r *http.Request) {
	if err := s.garantirConcepcaoSessao("sess_1"); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/sessao/sess_1", http.StatusSeeOther)
}

func (s *Server) garantirConcepcaoSessao(sessID string) error {
	hub := s.A2UIHub()
	if _, err := hub.GetSession(sessID); err == nil {
		return nil
	}
	hub.EnsureSession(sessID, "legoloop-concepcao")
	initState := a2ui.ConcepcaoState{
		SessionID: sessID,
		SurfaceID: a2ui.DefaultConcepcaoSurfaceID,
		Meta:      95,
		Pergunta:  "Qual ideia ou demanda você deseja conceber no LegoLoop?",
		Palpite:   "Descreva o problema, quem se beneficia e o resultado esperado.",
	}
	if err := a2ui.PublishConcepcao(hub, initState); err != nil {
		return err
	}
	s.setConcepcaoInterviewID(sessID, "")
	return nil
}

func (s *Server) dispatchConcepcaoAction(ctx context.Context, act a2ui.Action) error {
	if _, ativa := s.getConcepcaoInterviewEntry(act.SessionID); !ativa {
		return nil
	}
	switch act.Name {
	case "submit_answer":
		return s.stepConcepcaoA2UI(ctx, act)
	case "stop_gate_decision":
		return s.stopGateConcepcaoA2UI(ctx, act)
	default:
		return nil
	}
}

func (s *Server) stepConcepcaoA2UI(ctx context.Context, act a2ui.Action) error {
	texto := extrairTextoResposta(act.Payload)
	if texto == "" {
		return nil
	}
	engine := s.getInterviewEngine(ctx)
	intvID := s.getConcepcaoInterviewID(act.SessionID)
	var turn *interview.TurnResult
	var err error
	if intvID == "" {
		var isess *interview.Session
		isess, turn, err = engine.Start(ctx, texto)
		if err != nil {
			return err
		}
		s.setConcepcaoInterviewID(act.SessionID, isess.ID)
		intvID = isess.ID
	} else {
		turn, err = engine.Step(ctx, intvID, texto)
		if err != nil {
			return err
		}
	}
	return a2ui.PublishConcepcao(s.A2UIHub(), turnToConcepcaoState(act.SessionID, turn, ""))
}

func (s *Server) stopGateConcepcaoA2UI(ctx context.Context, act a2ui.Action) error {
	dec, _ := act.Payload["decision"].(string)
	choice, _ := act.Payload["choice"].(string)
	d := strings.TrimSpace(dec)
	c := strings.TrimSpace(choice)

	if c == "0" || d == "cancelar" {
		return s.cancelarConcepcaoA2UI(act)
	}
	if c == "2" || d == "personalizar" {
		return s.personalizarConcepcaoA2UI(ctx, act)
	}
	if c == "1" || d == "aprovar" || d == "1" {
		return s.aprovarConcepcaoA2UI(ctx, act)
	}
	return nil
}

func (s *Server) cancelarConcepcaoA2UI(act a2ui.Action) error {
	s.setConcepcaoInterviewID(act.SessionID, "")
	st := a2ui.ConcepcaoState{
		SessionID:          act.SessionID,
		SurfaceID:          a2ui.DefaultConcepcaoSurfaceID,
		Confianca:          0,
		Meta:               95,
		Pergunta:           "Processo de concepção cancelado pelo usuário.",
		Palpite:            "Inicie uma nova demanda quando desejar.",
		AguardandoStopGate: false,
	}
	return a2ui.PublishConcepcao(s.A2UIHub(), st)
}

func (s *Server) personalizarConcepcaoA2UI(ctx context.Context, act a2ui.Action) error {
	intvID := s.getConcepcaoInterviewID(act.SessionID)
	if intvID == "" {
		return nil
	}
	note, _ := act.Payload["note"].(string)
	msg := strings.TrimSpace(note)
	if msg == "" {
		msg = "Desejo personalizar e ajustar os requisitos do espelho."
	}
	engine := s.getInterviewEngine(ctx)
	turn, err := engine.Step(ctx, intvID, msg)
	if err != nil {
		return err
	}
	return a2ui.PublishConcepcao(s.A2UIHub(), turnToConcepcaoState(act.SessionID, turn, ""))
}

func (s *Server) aprovarConcepcaoA2UI(ctx context.Context, act a2ui.Action) error {
	intvID := s.getConcepcaoInterviewID(act.SessionID)
	if intvID == "" {
		return nil
	}
	engine := s.getInterviewEngine(ctx)
	path, err := engine.MaterializeIntent(intvID, s.ProjectDir, "")
	if err != nil {
		return err
	}
	rel := path
	if r, rErr := filepath.Rel(s.ProjectDir, path); rErr == nil {
		rel = filepath.ToSlash(r)
	}
	st := a2ui.ConcepcaoState{
		SessionID:          act.SessionID,
		SurfaceID:          a2ui.DefaultConcepcaoSurfaceID,
		Confianca:          100,
		Meta:               95,
		IntentPath:         rel,
		AguardandoStopGate: false,
	}
	return a2ui.PublishConcepcao(s.A2UIHub(), st)
}

func extrairTextoResposta(payload map[string]any) string {
	for _, k := range []string{"answer", "text", "value"} {
		if v, ok := payload[k].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func turnToConcepcaoState(sessID string, turn *interview.TurnResult, intentPath string) a2ui.ConcepcaoState {
	st := a2ui.ConcepcaoState{
		SessionID:          sessID,
		SurfaceID:          a2ui.DefaultConcepcaoSurfaceID,
		Confianca:          turn.Confidence,
		Meta:               95,
		Pergunta:           turn.Question,
		Palpite:            turn.Guess,
		IntentPath:         intentPath,
		AguardandoStopGate: turn.Status == interview.StatusCompleted,
	}
	if turn.Status == interview.StatusCompleted {
		st.Pergunta = ""
		st.Palpite = ""
	}
	if turn.Mirror != nil {
		st.Espelho = map[string]string{
			"resultado_esperado":      turn.Mirror.Outcome,
			"quem_se_beneficia":       turn.Mirror.User,
			"por_que_agora":           turn.Mirror.WhyNow,
			"criterio_de_sucesso":     turn.Mirror.Success,
			"restricoes_inegociaveis": strings.Join(turn.Mirror.Constraints, "; "),
			"fora_de_escopo":          strings.Join(turn.Mirror.OutOfScope, "; "),
		}
	}
	return st
}

func (s *Server) getConcepcaoInterviewEntry(sessID string) (string, bool) {
	s.concepcaoMu.Lock()
	defer s.concepcaoMu.Unlock()
	if s.concepcaoMap == nil {
		return "", false
	}
	id, ok := s.concepcaoMap[sessID]
	return id, ok
}

func (s *Server) getConcepcaoInterviewID(sessID string) string {
	id, _ := s.getConcepcaoInterviewEntry(sessID)
	return id
}

func (s *Server) setConcepcaoInterviewID(sessID, intvID string) {
	s.concepcaoMu.Lock()
	defer s.concepcaoMu.Unlock()
	if s.concepcaoMap == nil {
		s.concepcaoMap = make(map[string]string)
	}
	s.concepcaoMap[sessID] = intvID
}
