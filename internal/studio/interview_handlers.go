package studio

import (
	"context"
	"encoding/json"
	"net/http"

	"legoloop/internal/interview"
	"legoloop/internal/provider"
)

type interviewStartReq struct {
	Idea string `json:"idea"`
}

type interviewMsgReq struct {
	SessionID string `json:"session_id"`
	Answer    string `json:"answer"`
}

type interviewFinalizeReq struct {
	SessionID string `json:"session_id"`
	FileName  string `json:"file_name,omitempty"`
}

func (s *Server) getInterviewEngine(ctx context.Context) *interview.Engine {
	if s.Interview != nil {
		return s.Interview
	}
	inv := s.inventory(ctx)
	routerURL := inv.RouterURL
	client := interview.NewRouterClient(routerURL, provider.TrustedRouterKey(routerURL), "code-pro")
	s.Interview = interview.NewEngine(client)
	return s.Interview
}

func (s *Server) handleInterviewStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req interviewStartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Idea == "" {
		http.Error(w, "o campo 'idea' é obrigatório", http.StatusBadRequest)
		return
	}

	engine := s.getInterviewEngine(r.Context())
	sess, turn, err := engine.Start(r.Context(), req.Idea)
	if err != nil {
		http.Error(w, "falha ao iniciar entrevista: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"session": sess,
		"turn":    turn,
	})
}

func (s *Server) handleInterviewMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req interviewMsgReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.SessionID == "" || req.Answer == "" {
		http.Error(w, "campos 'session_id' e 'answer' são obrigatórios", http.StatusBadRequest)
		return
	}

	engine := s.getInterviewEngine(r.Context())
	turn, err := engine.Step(r.Context(), req.SessionID, req.Answer)
	if err != nil {
		http.Error(w, "falha ao avançar turno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"turn": turn,
	})
}

func (s *Server) handleInterviewFinalize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req interviewFinalizeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.SessionID == "" {
		http.Error(w, "o campo 'session_id' é obrigatório", http.StatusBadRequest)
		return
	}

	engine := s.getInterviewEngine(r.Context())
	path, err := engine.MaterializeIntent(req.SessionID, s.ProjectDir, req.FileName)
	if err != nil {
		http.Error(w, "falha ao materializar intenção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"artifact_path": path,
		"status":        "MATERIALIZED",
	})
}
