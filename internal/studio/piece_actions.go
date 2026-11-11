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
	kind, id := r.PathValue("kind"), r.PathValue("id")
	if kind == "" || id == "" {
		http.Error(w, "kind e id são obrigatórios", http.StatusBadRequest)
		return
	}
	newStatus := s.toggleFavorite(kind + ":" + id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"kind": kind, "id": id, "favorited": newStatus})
}

// toggleFavorite alterna a chave em favorites.json do Cofre canônico e devolve o novo estado.
func (s *Server) toggleFavorite(key string) bool {
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
	favs[key] = !favs[key]
	favBytes, _ := json.MarshalIndent(favs, "", "  ")
	_ = os.WriteFile(favFile, favBytes, 0644)
	return favs[key]
}

// findPiecePath localiza o arquivo físico da peça no inventário (vazio se ausente).
func (s *Server) findPiecePath(kind, id string) string {
	if s.Assets == nil {
		return ""
	}
	for _, item := range s.Assets.Snapshot().Items {
		if strings.EqualFold(string(item.Kind), kind) && (strings.EqualFold(item.ID, id) || strings.EqualFold(item.Name, id)) {
			return item.Path
		}
	}
	return ""
}

// handlePieceUpgrade dispara o início do subprocesso de elevação para Padrão Ouro.
func (s *Server) handlePieceUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind, id := r.PathValue("kind"), r.PathValue("id")
	if kind == "" || id == "" {
		http.Error(w, "kind e id são obrigatórios", http.StatusBadRequest)
		return
	}
	rating := reputation.PieceRating{Ref: id, Name: id, TotalScore: 0.50, Tier: reputation.TierBronze}
	if piecePath := s.findPiecePath(kind, id); piecePath != "" {
		if scored, err := reputation.ScoreFile(piecePath); err == nil {
			rating = scored
		}
	}
	// Inicia a entrevista socrática de elevação com o Addy Osmani
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
		"status":         "UPGRADE_INITIATED",
		"target_piece":   id,
		"current_rating": rating,
		"session":        sess,
		"first_turn":     turn,
		"processo_alvo":  "docs/intents/INTENT_UPGRADE_" + strings.ToUpper(id) + ".md",
	})
}

// handlePieceDelete executa a eliminação profunda (modo local ou full).
func (s *Server) handlePieceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind, id := r.PathValue("kind"), r.PathValue("id")
	mode := r.URL.Query().Get("mode") // "local" ou "full"
	if mode == "" {
		mode = "local"
	}
	targetPath := s.findPiecePath(kind, id)
	if targetPath == "" {
		http.Error(w, "peça não encontrada no inventário", http.StatusNotFound)
		return
	}
	if status, msg := s.eliminatePiece(mode, targetPath, id); status != 0 {
		http.Error(w, msg, status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "DELETED", "mode": mode, "eliminated": id})
}

// eliminatePiece aplica o modo de eliminação; devolve status HTTP != 0 e mensagem se falhar.
func (s *Server) eliminatePiece(mode, targetPath, id string) (int, string) {
	switch mode {
	case "full":
		// Eliminação completa do disco (expurgo profundo)
		targetDir := targetPath
		if !strings.EqualFold(filepath.Base(targetPath), id) {
			targetDir = filepath.Dir(targetPath)
		}
		if err := os.RemoveAll(targetDir); err != nil {
			return http.StatusInternalServerError, "falha ao eliminar arquivo físico: " + err.Error()
		}
	case "local":
		// Eliminação a nível de projeto (remove do índice local)
		if strings.HasPrefix(targetPath, s.ProjectDir) {
			_ = os.Remove(targetPath)
		}
	default:
		return http.StatusBadRequest, "modo de eliminação inválido (use 'local' ou 'full')"
	}
	return 0, ""
}
