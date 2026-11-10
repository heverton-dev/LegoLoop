package studio

import (
	"encoding/json"
	"errors"
	"legoloop/internal/caminhos"
	"net/http"

	"legoloop/internal/loop"
)

func (s *Server) loopStore() (*loop.Store, error) {
	if s.Loops != nil {
		return s.Loops, nil
	}
	v, err := s.vault()
	if err != nil {
		return nil, err
	}
	dir := caminhos.DirLoopsEm(v)
	store, err := loop.NewStore(dir)
	if err != nil {
		return nil, err
	}
	s.Loops = store
	return store, nil
}

func (s *Server) handleLoopsList(w http.ResponseWriter, r *http.Request) {
	store, err := s.loopStore()
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
		list = []loop.Summary{}
	}
	seen := map[string]bool{}
	for _, l := range list {
		seen[l.ID] = true
	}
	if s.ProjectDir != "" {
		if pStore, err := loop.NewStore(caminhos.DirLoopsEm(s.ProjectDir)); err == nil {
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

func (s *Server) handleLoopsGet(w http.ResponseWriter, r *http.Request) {
	store, err := s.loopStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	lp, err := store.Get(id)
	if err != nil && s.ProjectDir != "" {
		if pStore, pErr := loop.NewStore(caminhos.DirLoopsEm(s.ProjectDir)); pErr == nil {
			if pLoop, pGetErr := pStore.Get(id); pGetErr == nil {
				lp = pLoop
				err = nil
			}
		}
	}
	if err != nil {
		if errors.Is(err, loop.ErrNotFound) || errors.Is(err, loop.ErrInvalidID) {
			http.Error(w, "loop não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, lp)
}

func (s *Server) handleLoopsSave(w http.ResponseWriter, r *http.Request) {
	store, err := s.loopStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lp loop.Loop
	if err := json.NewDecoder(r.Body).Decode(&lp); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if lp.ID == "" {
		lp.ID = loop.ToKebabCase(lp.Name)
	} else {
		lp.ID = loop.ToKebabCase(lp.ID)
	}
	if err := store.Save(&lp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.Assets != nil {
		s.Assets.Refresh(r.Context())
	}
	writeJSON(w, http.StatusOK, lp)
}

func (s *Server) handleLoopsDuplicate(w http.ResponseWriter, r *http.Request) {
	store, err := s.loopStore()
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
	targetID := loop.ToKebabCase(req.ID)
	if targetID == "" {
		targetID = loop.ToKebabCase(req.Name)
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

func (s *Server) handleLoopsDelete(w http.ResponseWriter, r *http.Request) {
	store, err := s.loopStore()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := r.PathValue("id")
	if err := store.Remove(id); err != nil {
		if errors.Is(err, loop.ErrNotFound) {
			http.Error(w, "loop não encontrado", http.StatusNotFound)
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
