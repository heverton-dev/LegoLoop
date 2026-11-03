package studio

import (
	"encoding/json"
	"net/http"

	"legoloop/internal/interview"
	"legoloop/internal/specgen"
)

type specGenerateReq struct {
	IntentPath string `json:"intent_path"`
}

func (s *Server) getSpecGenerator() *specgen.Generator {
	if s.SpecGen != nil {
		return s.SpecGen
	}
	inv := s.inventory(nil)
	client := interview.NewRouterClient(inv.RouterURL, "", "code-pro")
	s.SpecGen = specgen.NewGenerator(client, s.Index)
	return s.SpecGen
}

func (s *Server) handleSpecGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	var req specGenerateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.IntentPath == "" {
		http.Error(w, "o campo 'intent_path' é obrigatório", http.StatusBadRequest)
		return
	}

	gen := s.getSpecGenerator()
	res, err := gen.Generate(r.Context(), s.ProjectDir, req.IntentPath)
	if err != nil {
		http.Error(w, "falha ao gerar especificação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
