package studio

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type revealRequest struct {
	Path string `json:"path"`
}

// resolveSafePath valida e expande caminhos seguros (~ e <projeto>).
func (s *Server) resolveSafePath(raw string) (string, bool, error) {
	if strings.Contains(raw, "..") {
		return "", false, errors.New("caminho com path traversal rejeitado")
	}
	clean := raw
	if strings.HasPrefix(clean, "~/") || strings.HasPrefix(clean, `~\`) || clean == "~" {
		clean = filepath.Join(s.home(), clean[1:])
	} else if strings.HasPrefix(clean, "<projeto>/") || strings.HasPrefix(clean, "<projeto>\\") {
		clean = filepath.Join(s.ProjectDir, clean[len("<projeto>/"):])
	}
	abs, err := filepath.Abs(clean)
	if err != nil {
		return "", false, err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return "", false, errors.New("caminho inexistente")
	}
	return abs, st.IsDir(), nil
}

func launchExplorer(path string, isDir bool) error {
	if runtime.GOOS == "windows" {
		if isDir {
			return exec.Command("explorer.exe", path).Start()
		}
		return exec.Command("explorer.exe", "/select,"+path).Start()
	}
	if runtime.GOOS == "darwin" {
		return exec.Command("open", "-R", path).Start()
	}
	return exec.Command("xdg-open", filepath.Dir(path)).Start()
}

func (s *Server) handleAssetReveal(w http.ResponseWriter, r *http.Request) {
	var req revealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	abs, isDir, err := s.resolveSafePath(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := launchExplorer(abs, isDir); err != nil {
		http.Error(w, "falha ao abrir explorador: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "path": abs})
}
