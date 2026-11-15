package studio

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// TranslationRecord armazena a tradução didática sem tocar no arquivo original.
type TranslationRecord struct {
	SHA256          string `json:"sha256"`
	SourcePath      string `json:"source_path"`
	WhatItIs        string `json:"what_it_is"`
	WhatItDoes      string `json:"what_it_does"`
	WhatItDelivers  string `json:"what_it_delivers"`
	FullContentPTBR string `json:"full_content_ptbr"`
}

func (s *Server) translationCacheDir() (string, error) {
	home := s.Home
	if home == "" && s.Vault != "" {
		dir := filepath.Join(s.Vault, ".cache", "translations")
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", err
		}
		return dir, nil
	}
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	dir := filepath.Join(home, ".aidd", "cache", "translations")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

func (s *Server) findRevealRootRel(abs string) (string, string, bool) {
	cleanAbs := filepath.Clean(abs)
	for _, rootDir := range s.revealRoots() {
		rel, err := filepath.Rel(rootDir, cleanAbs)
		if err != nil {
			continue
		}
		if rel != "." && filepath.IsLocal(rel) {
			return rootDir, rel, true
		}
	}
	return "", "", false
}

func (s *Server) validateSafeAssetPath(raw string) (string, error) {
	abs, isDir, err := s.resolveSafePath(raw)
	if err != nil {
		return "", err
	}
	if isDir {
		return "", errors.New("arquivo de ativo inexistente ou inválido")
	}
	if _, _, ok := s.findRevealRootRel(abs); !ok {
		return "", errors.New("path jail: caminho fora das raízes permitidas")
	}
	return abs, nil
}

func (s *Server) readAssetInRevealRoots(abs string) ([]byte, error) {
	rootDir, rel, ok := s.findRevealRootRel(abs)
	if !ok {
		return nil, errors.New("path jail: caminho fora das raízes permitidas")
	}
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return nil, errors.New("raiz do ativo inacessível")
	}
	defer root.Close()
	if err := checkNoLinksRoot(root, rel); err != nil {
		return nil, err
	}
	return root.ReadFile(rel)
}

func (s *Server) maskAssetText(text string) string {
	out := s.maskPaths(text)
	if s.Vault != "" {
		v := filepath.Clean(s.Vault)
		fwd := filepath.ToSlash(v)
		out = strings.NewReplacer(
			v, "<cofre>",
			fwd, "<cofre>",
			strings.ReplaceAll(fwd, "/", `\`), "<cofre>",
			strings.ReplaceAll(v, `\`, `\\`), "<cofre>",
		).Replace(out)
	}
	return out
}

func (s *Server) maskAssetSourcePath(abs string) string {
	masked := s.maskAssetText(abs)
	if filepath.IsAbs(masked) {
		return s.displayPath(abs)
	}
	return masked
}

func (s *Server) readCachedTranslation(hash string) (*TranslationRecord, error) {
	dir, err := s.translationCacheDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, hash+".json"))
	if err != nil {
		return nil, err
	}
	var rec TranslationRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *Server) saveCachedTranslation(rec *TranslationRecord) error {
	dir, err := s.translationCacheDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, rec.SHA256+".json"), data, 0600)
}

func buildDidacticExplanation(path string, content []byte) TranslationRecord {
	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return TranslationRecord{
		SHA256:          hash,
		SourcePath:      path,
		WhatItIs:        fmt.Sprintf("Peça do ecossistema: %s.", name),
		WhatItDoes:      "Executa rotinas especializadas de automação, engenharia de software e testes.",
		WhatItDelivers:  "Entrega artefatos validados, relatórios estruturados e código de produção em conformidade.",
		FullContentPTBR: string(content),
	}
}

func (s *Server) handleAssetTranslate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Path == "" {
		http.Error(w, "JSON inválido ou path ausente", http.StatusBadRequest)
		return
	}
	cleanPath, err := s.validateSafeAssetPath(req.Path)
	if err != nil {
		http.Error(w, s.maskAssetText(err.Error()), http.StatusBadRequest)
		return
	}
	content, err := s.readAssetInRevealRoots(cleanPath)
	if err != nil {
		http.Error(w, "erro ao ler ativo: "+s.maskAssetText(err.Error()), http.StatusInternalServerError)
		return
	}
	maskedSource := s.maskAssetSourcePath(cleanPath)
	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])
	if rec, err := s.readCachedTranslation(hash); err == nil && rec != nil {
		rec.SourcePath = maskedSource
		writeJSON(w, http.StatusOK, rec)
		return
	}
	rec := buildDidacticExplanation(maskedSource, content)
	_ = s.saveCachedTranslation(&rec)
	writeJSON(w, http.StatusOK, rec)
}
