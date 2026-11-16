package studio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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
	newStatus, err := s.toggleFavorite(kind + ":" + id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao atualizar favorito: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "id": id, "favorited": newStatus})
}

// toggleFavorite alterna a chave em favorites.json do Cofre canônico com exclusão mútua e tratamento real de erro.
func (s *Server) toggleFavorite(key string) (bool, error) {
	s.favMu.Lock()
	defer s.favMu.Unlock()
	vaultDir := s.studioVaultDir()
	if vaultDir == "" {
		return false, errors.New("cofre canônico não disponível")
	}
	if err := os.MkdirAll(vaultDir, 0o755); err != nil {
		return false, fmt.Errorf("falha ao criar diretório do cofre: %w", err)
	}
	favFile := filepath.Join(vaultDir, "favorites.json")
	favs, err := readFavoritesMap(favFile)
	if err != nil {
		return false, err
	}
	favs[key] = !favs[key]
	favBytes, err := json.MarshalIndent(favs, "", "  ")
	if err != nil {
		return false, fmt.Errorf("falha ao serializar favorites.json: %w", err)
	}
	if err := writeAtomic(vaultDir, "favorites.json", favBytes); err != nil {
		return false, fmt.Errorf("falha ao gravar favorites.json: %w", err)
	}
	return favs[key], nil
}

func readFavoritesMap(favFile string) (map[string]bool, error) {
	favs := make(map[string]bool)
	data, err := os.ReadFile(favFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return favs, nil
		}
		return nil, fmt.Errorf("falha ao ler favorites.json: %w", err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return favs, nil
	}
	if err := json.Unmarshal(data, &favs); err != nil {
		return nil, fmt.Errorf("favorites.json corrompido: %w", err)
	}
	return favs, nil
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

// handlePieceUpgrade executa a elevação determinística para Padrão Ouro (ou entrevista quando solicitado).
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
	var req struct {
		Mode string `json:"mode"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	piecePath := s.findPiecePath(kind, id)
	if strings.TrimSpace(req.Mode) != "interview" && piecePath != "" {
		s.upgradeDeterministic(w, r, kind, id, piecePath)
		return
	}
	s.upgradeWithInterview(w, r, kind, id, piecePath)
}

func (s *Server) upgradeDeterministic(w http.ResponseWriter, r *http.Request, kind, id, piecePath string) {
	rootDir, rel, ok := matchSubpathInRoots(piecePath, s.upgradeRoots())
	if !ok {
		writeError(w, http.StatusForbidden, "path jail: elevação fora das raízes permitidas")
		return
	}
	if _, _, err := resolveInRoot(rootDir, rel, false); err != nil {
		writeError(w, http.StatusForbidden, "path jail: "+err.Error())
		return
	}
	rep, err := reputation.UpgradePieceFile(piecePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao elevar peça: "+err.Error())
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(r.Context())
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "UPGRADE_COMPLETED",
		"mode":           "deterministic",
		"kind":           kind,
		"target_piece":   id,
		"changed":        rep.Changed,
		"before_rating":  rep.Before,
		"current_rating": rep.After,
		"added_sections": rep.AddedSections,
	})
}

func (s *Server) upgradeWithInterview(w http.ResponseWriter, r *http.Request, kind, id, piecePath string) {
	rating := reputation.PieceRating{Ref: id, Name: id, TotalScore: 0.50, Tier: reputation.TierBronze}
	if piecePath != "" {
		if scored, err := reputation.ScoreFile(piecePath); err == nil {
			rating = scored
		}
	}
	engine := s.getInterviewEngine(r.Context())
	initialIdea := fmt.Sprintf("Elevar a peça %s/%s (Score Atual: %.2f - Tier: %s) para o Padrão Ouro (Score >= 0.85).",
		kind, id, rating.TotalScore, rating.Tier)
	sess, turn, err := engine.Start(r.Context(), initialIdea)
	if err != nil {
		http.Error(w, "falha ao iniciar entrevista de elevação: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "UPGRADE_INITIATED",
		"target_piece":   id,
		"current_rating": rating,
		"session":        sess,
		"first_turn":     turn,
		"processo_alvo":  "docs/intents/INTENT_UPGRADE_" + strings.ToUpper(id) + ".md",
	})
}

func (s *Server) upgradeRoots() []string {
	roots := s.deleteRoots()
	if s.Home != "" {
		cleanHome := filepath.Clean(s.Home)
		if !slices.Contains(roots, cleanHome) {
			roots = append(roots, cleanHome)
		}
	}
	return roots
}

// handlePieceDelete executa a eliminação profunda (modo local ou full).
func (s *Server) handlePieceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
		return
	}
	kind, id := r.PathValue("kind"), r.PathValue("id")
	mode := r.URL.Query().Get("mode")
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
	writeJSON(w, http.StatusOK, map[string]any{"status": "DELETED", "mode": mode, "eliminated": id})
}

// eliminatePiece aplica o modo de eliminação com Path Jail estrito; devolve status HTTP != 0 e mensagem se falhar.
func (s *Server) eliminatePiece(mode, targetPath, id string) (int, string) {
	switch mode {
	case "full":
		return s.eliminateFull(targetPath, id)
	case "local":
		return s.eliminateLocal(targetPath)
	default:
		return http.StatusBadRequest, "modo de eliminação inválido (use 'local' ou 'full')"
	}
}

func (s *Server) eliminateFull(targetPath, id string) (int, string) {
	target := targetPath
	if !strings.EqualFold(filepath.Base(targetPath), id) {
		target = filepath.Dir(targetPath)
	}
	rootDir, rel, ok := matchSubpathInRoots(target, s.deleteRoots())
	if !ok {
		return http.StatusForbidden, "path jail: exclusão fora do projeto ou cofre bloqueada"
	}
	if _, _, err := resolveInRoot(rootDir, rel, false); err != nil {
		return http.StatusForbidden, "path jail: " + err.Error()
	}
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return http.StatusInternalServerError, "falha ao abrir raiz de exclusão: " + err.Error()
	}
	defer root.Close()
	if err := root.RemoveAll(rel); err != nil {
		return http.StatusInternalServerError, "falha ao eliminar arquivo físico: " + err.Error()
	}
	return 0, ""
}

func (s *Server) eliminateLocal(targetPath string) (int, string) {
	rootDir, rel, ok := matchSubpathInRoots(targetPath, []string{s.ProjectDir})
	if !ok {
		return http.StatusForbidden, "path jail: exclusão local permitida apenas dentro do projeto"
	}
	if _, _, err := resolveInRoot(rootDir, rel, false); err != nil {
		return http.StatusForbidden, "path jail: " + err.Error()
	}
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return http.StatusInternalServerError, "falha ao abrir raiz do projeto: " + err.Error()
	}
	defer root.Close()
	if err := root.Remove(rel); err != nil {
		return http.StatusInternalServerError, "falha ao remover arquivo local: " + err.Error()
	}
	return 0, ""
}

func (s *Server) deleteRoots() []string {
	var out []string
	for _, dir := range []string{s.ProjectDir, s.studioVaultDir()} {
		if dir == "" {
			continue
		}
		clean := filepath.Clean(dir)
		if !slices.Contains(out, clean) {
			out = append(out, clean)
		}
	}
	return out
}

func matchSubpathInRoots(target string, roots []string) (string, string, bool) {
	if strings.TrimSpace(target) == "" || strings.Contains(target, "..") {
		return "", "", false
	}
	cleanTarget := filepath.Clean(target)
	for _, base := range roots {
		if base == "" {
			continue
		}
		cleanBase := filepath.Clean(base)
		rel, err := filepath.Rel(cleanBase, cleanTarget)
		if err == nil && rel != "." && filepath.IsLocal(rel) {
			return cleanBase, rel, true
		}
	}
	return "", "", false
}
