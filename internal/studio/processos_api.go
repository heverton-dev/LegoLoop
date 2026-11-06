package studio

import (
	"encoding/json"
	"errors"
	"net/http"

	"legoloop/internal/compor"
	"legoloop/internal/listar"
	"legoloop/internal/processo"
	"legoloop/internal/refinar"
	"legoloop/internal/remover"
)

func (s *Server) processoStore() (processo.Store, error) {
	if s.StoreOverride != nil {
		return *s.StoreOverride, nil
	}
	v, err := s.vault()
	if err != nil {
		return processo.Store{}, err
	}
	docs := s.ProjectDir
	if docs == "" {
		docs = "."
	}
	return processo.Store{Roots: processo.Roots{Vault: v, Docs: docs}}, nil
}

func (s *Server) handleProcessosListar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	st, err := s.processoStore()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := listar.Run(r.Context(), listar.Deps{Store: st, Home: s.Home})
	if err != nil && !errors.Is(err, listar.ErrSemCofre) {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

type procRefinarReq struct {
	Name   string   `json:"name"`
	Etapa  string   `json:"etapa"`
	Peca   string   `json:"peca"`
	Caso   []string `json:"caso"`
	Titulo string   `json:"titulo"`
	Apply  bool     `json:"apply"`
}

func (s *Server) handleProcessosRefinar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req procRefinarReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	st, err := s.processoStore()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refReq := refinar.Request{
		Nome:    req.Name,
		Mudanca: refinar.Mudanca{Etapa: req.Etapa, Peca: req.Peca, Titulo: req.Titulo, Casos: req.Caso},
	}
	if !req.Apply {
		s.execRefinarPreview(w, refReq, st)
		return
	}
	s.execRefinarApply(w, r, refReq, st)
}

func (s *Server) execRefinarPreview(w http.ResponseWriter, req refinar.Request, st processo.Store) {
	pl, err := refinar.Planejar(req, refinar.Deps{Store: st})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := map[string]any{
		"nome":     pl.Nome,
		"de":       pl.De,
		"para":     pl.Para,
		"diff":     pl.Diff(),
		"cartao":   pl.Cartao(),
		"mudanca":  pl.Mudanca,
		"avisos":   pl.Avisos,
		"arquivos": pl.Arquivos,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) execRefinarApply(w http.ResponseWriter, r *http.Request, req refinar.Request, st processo.Store) {
	confirm := func(pl refinar.Plano, _ func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		return true, nil
	}
	res, err := refinar.Run(r.Context(), req, refinar.Deps{Store: st, Confirm: confirm})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "resultado": res})
}

type procRemoverReq struct {
	Name  string `json:"name"`
	Apply bool   `json:"apply"`
}

func (s *Server) handleProcessosRemover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req procRemoverReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	st, err := s.processoStore()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pl, err := remover.Planejar(r.Context(), remover.Deps{Store: st, Home: s.Home}, req.Name)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, remover.ErrEmUso) {
			code = http.StatusConflict
		}
		writeError(w, code, err.Error())
		return
	}
	if !req.Apply {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pl)
		return
	}
	if err := remover.Executar(r.Context(), pl); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "nome": pl.Nome})
}

type procComporReq struct {
	Name     string   `json:"name"`
	Composto string   `json:"composto"`
	Titulo   string   `json:"titulo"`
	Depende  []string `json:"depende"`
	Apply    bool     `json:"apply"`
}

func (s *Server) handleProcessosCompor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		return
	}
	var req procComporReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	st, err := s.processoStore()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	compReq := compor.Request{
		Nome:    req.Name,
		Ligacao: compor.Ligacao{Composto: req.Composto, Titulo: req.Titulo, Depende: req.Depende},
	}
	if !req.Apply {
		s.execComporPreview(w, compReq, st)
		return
	}
	s.execComporApply(w, r, compReq, st)
}

func (s *Server) execComporPreview(w http.ResponseWriter, req compor.Request, st processo.Store) {
	pl, err := compor.Planejar(req, compor.Deps{Store: st})
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, processo.ErrCiclo) {
			code = http.StatusConflict
		}
		writeError(w, code, err.Error())
		return
	}
	resp := map[string]any{
		"nome":     pl.Nome,
		"de":       pl.De,
		"para":     pl.Para,
		"diff":     pl.Diff(),
		"cartao":   pl.Cartao(),
		"arquivos": pl.Arquivos,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) execComporApply(w http.ResponseWriter, r *http.Request, req compor.Request, st processo.Store) {
	confirm := func(pl compor.Plano, _ func(compor.Ligacao) (compor.Plano, error)) (bool, error) {
		return true, nil
	}
	res, err := compor.Run(r.Context(), req, compor.Deps{Store: st, Confirm: confirm})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "resultado": res})
}
