package studio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

type revealRequest struct {
	Path string `json:"path"`
}

func isDisallowedLink(m os.FileMode) bool {
	return m&(os.ModeSymlink|os.ModeIrregular|os.ModeDevice|os.ModeNamedPipe|os.ModeSocket) != 0
}

func checkNoLinksRoot(root *os.Root, rel string) error {
	if rel == "." {
		return nil
	}
	cur := ""
	for part := range strings.SplitSeq(rel, string(filepath.Separator)) {
		if cur == "" {
			cur = part
		} else {
			cur = filepath.Join(cur, part)
		}
		info, err := root.Lstat(cur)
		if err != nil {
			return err
		}
		if isDisallowedLink(info.Mode()) {
			return fmt.Errorf("symlink ou junction não permitido: %q", filepath.ToSlash(cur))
		}
	}
	return nil
}

func resolveInRoot(rootDir, rel string, allowRoot bool) (string, bool, error) {
	if rootDir == "" {
		return "", false, errors.New("raiz não configurada")
	}
	rootInfo, err := os.Lstat(rootDir)
	if err != nil {
		return "", false, errors.New("caminho inexistente")
	}
	if isDisallowedLink(rootInfo.Mode()) {
		return "", false, errors.New("symlink ou junction não permitido na raiz")
	}
	cleanRel := filepath.Clean(filepath.FromSlash(strings.TrimSpace(rel)))
	if cleanRel == "." {
		if !allowRoot {
			return "", false, errors.New("path jail: operação direta na raiz bloqueada")
		}
		return filepath.Clean(rootDir), rootInfo.IsDir(), nil
	}
	if strings.Contains(cleanRel, ":") || !filepath.IsLocal(cleanRel) {
		return "", false, errors.New("path jail: caminho fora da raiz permitida")
	}
	return statInRoot(filepath.Clean(rootDir), cleanRel)
}

func statInRoot(rootDir, cleanRel string) (string, bool, error) {
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return "", false, errors.New("caminho inexistente")
	}
	defer root.Close()
	if err := checkNoLinksRoot(root, cleanRel); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, errors.New("caminho inexistente")
		}
		return "", false, err
	}
	st, err := root.Lstat(cleanRel)
	if err != nil {
		return "", false, errors.New("caminho inexistente")
	}
	return filepath.Join(rootDir, cleanRel), st.IsDir(), nil
}

func (s *Server) studioVaultDir() string {
	if s.Vault != "" {
		return s.Vault
	}
	if s.Home != "" {
		return filepath.Join(s.Home, ".aidd", "canonical")
	}
	if v, err := s.vault(); err == nil {
		return v
	}
	return ""
}

func (s *Server) revealRoots() []string {
	var out []string
	for _, dir := range []string{s.ProjectDir, s.studioVaultDir(), s.Home} {
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

// resolveSafePath valida e expande caminhos dentro do projeto, cofre ou home configurado, bloqueando links.
func (s *Server) resolveSafePath(raw string) (string, bool, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.ContainsRune(trimmed, 0) || strings.Contains(trimmed, "..") {
		return "", false, errors.New("caminho com path traversal rejeitado")
	}
	if trimmed == "~" || strings.HasPrefix(trimmed, "~/") || strings.HasPrefix(trimmed, `~\`) {
		rel := strings.TrimLeft(strings.TrimPrefix(trimmed, "~"), `/\`)
		return resolveInRoot(s.home(), cmpOrDot(rel), true)
	}
	if trimmed == "<projeto>" || strings.HasPrefix(trimmed, "<projeto>/") || strings.HasPrefix(trimmed, `<projeto>\`) {
		rel := strings.TrimLeft(strings.TrimPrefix(trimmed, "<projeto>"), `/\`)
		return resolveInRoot(s.ProjectDir, cmpOrDot(rel), true)
	}
	if !filepath.IsAbs(trimmed) && !strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "/") && !strings.HasPrefix(trimmed, `\`) {
		return resolveInRoot(s.ProjectDir, trimmed, false)
	}
	return s.resolveAbsInAllowedRoots(filepath.Clean(trimmed))
}

func cmpOrDot(rel string) string {
	if rel == "" {
		return "."
	}
	return rel
}

func (s *Server) resolveAbsInAllowedRoots(abs string) (string, bool, error) {
	for _, rootDir := range s.revealRoots() {
		rel, err := filepath.Rel(rootDir, abs)
		if err != nil {
			continue
		}
		if rel == "." || filepath.IsLocal(rel) {
			return resolveInRoot(rootDir, rel, true)
		}
	}
	return "", false, errors.New("path jail: caminho fora das raízes permitidas")
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
	launcher := s.RevealLauncher
	if launcher == nil {
		launcher = launchExplorer
	}
	if err := launcher(abs, isDir); err != nil {
		http.Error(w, "falha ao abrir explorador: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "path": s.displayPath(abs)})
}
