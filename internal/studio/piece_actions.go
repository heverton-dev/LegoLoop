package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"legoloop/internal/reputation"
)

// handlePieceFavorite alterna o status de favorito/SSOT de uma peça.
func (s *Server) handlePieceFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind := r.PathValue("kind")
	id := r.PathValue("id")
	if kind == "" || id == "" {
		http.Error(w, "kind e id são obrigatórios", http.StatusBadRequest)
		return
	}

	// Persiste a preferência no cofre canônico
	vaultDir := s.Vault
	if vaultDir == "" {
		home := s.Home
		if home == "" {
			home, _ = os.UserHomeDir()
		}
		vaultDir = filepath.Join(home, ".aidd", "canonical")
	}
	_ = os.MkdirAll(vaultDir, 0755)

	favFile := filepath.Join(vaultDir, "favorites.json")
	favs := make(map[string]bool)

	if data, err := os.ReadFile(favFile); err == nil {
		_ = json.Unmarshal(data, &favs)
	}

	key := kind + ":" + id
	newStatus := !favs[key]
	favs[key] = newStatus

	favBytes, _ := json.MarshalIndent(favs, "", "  ")
	_ = os.WriteFile(favFile, favBytes, 0644)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"kind":      kind,
		"id":        id,
		"favorited": newStatus,
	})
}

// handlePieceUpgrade dispara o início do sub-loop de elevação para Padrão Ouro.
func (s *Server) handlePieceUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind := r.PathValue("kind")
	id := r.PathValue("id")
	if kind == "" || id == "" {
		http.Error(w, "kind e id são obrigatórios", http.StatusBadRequest)
		return
	}

	// 1. Localizar arquivo físico da peça
	piecePath := ""
	if s.Assets != nil {
		for _, item := range s.Assets.Snapshot().Items {
			if strings.EqualFold(string(item.Kind), kind) && (strings.EqualFold(item.ID, id) || strings.EqualFold(item.Name, id)) {
				piecePath = item.Path
				break
			}
		}
	}

	// Diagnóstico matricial
	rating := reputation.PieceRating{
		Ref:        id,
		Name:       id,
		TotalScore: 0.50,
		Tier:       reputation.TierBronze,
	}

	if piecePath != "" {
		if r, err := reputation.ScoreFile(piecePath); err == nil {
			rating = r
		}
	}

	// 2. Iniciar entrevista socrática de elevação com o Addy Osmani
	engine := s.getInterviewEngine(r.Context())
	initialIdea := fmt.Sprintf("Elevar a peça %s/%s (Score Atual: %.2f - Tier: %s) para o Padrão Ouro (Score >= 0.85).",
		kind, id, rating.TotalScore, rating.Tier)

	sess, turn, err := engine.Start(r.Context(), initialIdea)
	if err != nil {
		http.Error(w, "falha ao iniciar entrevista de elevação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":          "UPGRADE_INITIATED",
		"target_piece":    id,
		"current_rating":  rating,
		"session":         sess,
		"first_turn":      turn,
		"pipeline_target": "docs/intents/INTENT_UPGRADE_" + strings.ToUpper(id) + ".md",
	})
}

// handlePieceDelete executa a eliminação profunda (modo local ou full).
func (s *Server) handlePieceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind := r.PathValue("kind")
	id := r.PathValue("id")
	mode := r.URL.Query().Get("mode") // "local" ou "full"
	if mode == "" {
		mode = "local"
	}

	targetPath := ""
	if s.Assets != nil {
		for _, item := range s.Assets.Snapshot().Items {
			if strings.EqualFold(string(item.Kind), kind) && (strings.EqualFold(item.ID, id) || strings.EqualFold(item.Name, id)) {
				targetPath = item.Path
				break
			}
		}
	}

	if targetPath == "" {
		http.Error(w, "peça não encontrada no inventário", http.StatusNotFound)
		return
	}

	switch mode {
	case "full":
		// Eliminação completa do disco (expurgo profundo)
		targetDir := targetPath
		if !strings.EqualFold(filepath.Base(targetPath), id) {
			targetDir = filepath.Dir(targetPath)
		}
		if err := os.RemoveAll(targetDir); err != nil {
			http.Error(w, "falha ao eliminar arquivo físico: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case "local":
		// Eliminação a nível de projeto (remove do índice local)
		if strings.HasPrefix(targetPath, s.ProjectDir) {
			_ = os.Remove(targetPath)
		}
	default:
		http.Error(w, "modo de eliminação inválido (use 'local' ou 'full')", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "DELETED",
		"mode":       mode,
		"eliminated": id,
	})
}
