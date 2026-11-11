package studio

import (
	"encoding/json"
	"errors"
	"legoloop/internal/caminhos"
	"net/http"

	"legoloop/internal/studiodoc"
)

func (s *Server) loopStore() (*studiodoc.Store, error) {
	return s.acervoStore()
}

func (s *Server) acervoStore() (*studiodoc.Store, error) {
	if s.Processos != nil {
		return s.Processos, nil
	}
	v, err := s.vault()
	if err != nil {
		return nil, err
	}
	store, err := studiodoc.NewStore(caminhos.DirLoopsEm(v))
	if err != nil {
		return nil, err
	}
	s.Processos = store
	return store, nil
}

func (s *Server) handleProcessosLista(w http.ResponseWriter, r *http.Request) {
	store, err := s.acervoStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	list, err := store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []studiodoc.Summary{}
	}
	seen := map[string]bool{}
	for _, l := range list {
		seen[l.ID] = true
	}
	if s.ProjectDir != "" {
		if pStore, err := studiodoc.NewStore(caminhos.DirLoopsEm(s.ProjectDir)); err == nil {
			if pList, err := pStore.List(); err == nil {
				for _, l := range pList {
					if !seen[l.ID] {
						seen[l.ID] = true
						list = append(list, l)
					}
				}
			}
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{caminhos.ChaveAPILoops: list})
}

func (s *Server) handleProcessosObter(w http.ResponseWriter, r *http.Request) {
	store, err := s.acervoStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	doc, err := store.Get(id)
	if err != nil && s.ProjectDir != "" {
		if pStore, pErr := studiodoc.NewStore(caminhos.DirLoopsEm(s.ProjectDir)); pErr == nil {
			if pDoc, pGetErr := pStore.Get(id); pGetErr == nil {
				doc = pDoc
				err = nil
			}
		}
	}
	if err != nil {
		if errors.Is(err, studiodoc.ErrNotFound) || errors.Is(err, studiodoc.ErrInvalidID) {
			http.Error(w, "Processo não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) handleProcessosGravar(w http.ResponseWriter, r *http.Request) {
	store, err := s.acervoStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var doc studiodoc.Processo
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if doc.ID == "" {
		doc.ID = studiodoc.ToKebabCase(doc.Name)
	} else {
		doc.ID = studiodoc.ToKebabCase(doc.ID)
	}
	if err := store.Save(&doc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(r.Context())
	}
	writeJSON(w, http.StatusOK, doc)
}

func (s *Server) handleProcessosDuplicar(w http.ResponseWriter, r *http.Request) {
	store, err := s.acervoStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	var req struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	targetID := studiodoc.ToKebabCase(req.ID)
	if targetID == "" {
		targetID = studiodoc.ToKebabCase(req.Name)
	}
	dup, err := store.Duplicate(id, targetID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(r.Context())
	}
	writeJSON(w, http.StatusOK, dup)
}

func (s *Server) handleProcessosExcluir(w http.ResponseWriter, r *http.Request) {
	store, err := s.acervoStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	if err := store.Remove(id); err != nil {
		if errors.Is(err, studiodoc.ErrNotFound) {
			http.Error(w, "Processo não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(r.Context())
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
