package studio

import (
	"context"
	"encoding/json"
	"net/http"

	"legoloop/internal/canon"
)

type canonPieceRequest struct {
	SHA256 string `json:"sha256"`
}

func filterPlanBySHA(full canon.Plan, hash string) canon.Plan {
	res := canon.Plan{Vault: full.Vault}
	for _, g := range full.Groups {
		if g.SHA256 == hash {
			res.Groups = append(res.Groups, g)
			res.SavedBytes += g.SavedBytes
		}
	}
	return res
}

func (s *Server) handleCanonizePiece(w http.ResponseWriter, r *http.Request) {
	var req canonPieceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SHA256 == "" {
		writeError(w, http.StatusBadRequest, "sha256 obrigatório")
		return
	}
	if !s.canonMu.TryLock() {
		writeError(w, http.StatusConflict, "já existe uma canonização em andamento")
		return
	}
	defer s.canonMu.Unlock()
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), canonTimeout)
	defer cancel()

	plan, _, err := s.buildPlan(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao gerar plano: "+err.Error())
		return
	}
	filtered := filterPlanBySHA(plan, req.SHA256)
	if len(filtered.Groups) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "linked": 0, "message": "Sem duplicatas pendentes"})
		return
	}
	rep := canon.Apply(ctx, filtered, canon.DefaultLinker())
	s.Assets.Refresh(ctx)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "linked": len(rep.Linked), "saved_bytes": rep.SavedBytes})
}
