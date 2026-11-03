package studio

import (
	"encoding/json"
	"net/http"

	"legoloop/internal/interview"
	"legoloop/internal/plangen"
)

type planGenerateReq struct {
	SpecPath         string                   `json:"spec_path"`
	SupervisionLevel plangen.SupervisionLevel `json:"supervision_level"`
}

func (s *Server) getPlanGenerator() *plangen.Generator {
	if s.PlanGen != nil {
		return s.PlanGen
	}
	inv := s.inventory(nil)
	client := interview.NewRouterClient(inv.RouterURL, "", "code-pro")
	s.PlanGen = plangen.NewGenerator(client)
	return s.PlanGen
}

func (s *Server) handlePlanGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req planGenerateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.SpecPath == "" {
		http.Error(w, "o campo 'spec_path' é obrigatório", http.StatusBadRequest)
		return
	}

	gen := s.getPlanGenerator()
	res, err := gen.Generate(r.Context(), s.ProjectDir, req.SpecPath, req.SupervisionLevel)
	if err != nil {
		http.Error(w, "falha ao gerar plano e rota: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
