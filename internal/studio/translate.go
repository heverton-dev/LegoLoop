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

func (s *Server) validateSafeAssetPath(raw string) (string, error) {
	if strings.Contains(raw, "..") {
		return "", errors.New("caminho com path traversal rejeitado")
	}
	clean := raw
	if strings.HasPrefix(clean, "~/") || strings.HasPrefix(clean, `~\`) || clean == "~" {
		clean = filepath.Join(s.home(), clean[1:])
	} else if strings.HasPrefix(clean, "<projeto>/") || strings.HasPrefix(clean, "<projeto>\\") {
		clean = filepath.Join(s.ProjectDir, clean[len("<projeto>/"):])
	}
	abs, err := filepath.Abs(clean)
	if err != nil {
		return "", err
	}
	st, err := os.Stat(abs)
	if err != nil || st.IsDir() {
		return "", errors.New("arquivo de ativo inexistente ou inválido")
	}
	return abs, nil
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		http.Error(w, "erro ao ler ativo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	h := sha256.Sum256(content)
	hash := hex.EncodeToString(h[:])
	if rec, err := s.readCachedTranslation(hash); err == nil && rec != nil {
		writeJSON(w, http.StatusOK, rec)
		return
	}
	rec := buildDidacticExplanation(cleanPath, content)
	_ = s.saveCachedTranslation(&rec)
	writeJSON(w, http.StatusOK, rec)
}
