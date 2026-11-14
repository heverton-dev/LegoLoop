package studio

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"legoloop/internal/paths"
	"legoloop/internal/triagem"
)

type triagemRequest struct {
	Entrada   string `json:"entrada"`
	Workspace string `json:"workspace"`
}

func (s *Server) handleTriagem(w http.ResponseWriter, r *http.Request) {
	var req triagemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if req.Entrada == "" && req.Workspace == "" {
		http.Error(w, "entrada ou workspace obrigatório", http.StatusBadRequest)
		return
	}
	loopsDir := paths.DirAntigo
	if s.ProjectDir != "" {
		loopsDir = filepath.Join(s.ProjectDir, paths.DirAntigo)
	}
	catalogo, err := carregarCatalogoTriagemStudio(loopsDir)
	if err != nil {
		http.Error(w, "erro ao carregar catálogo de triagem", http.StatusInternalServerError)
		return
	}
	dec, err := decidirTriagemStudio(r, req, catalogo)
	if err != nil {
		http.Error(w, "erro ao executar triagem", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dec)
}

func decidirTriagemStudio(r *http.Request, req triagemRequest, catalogo []triagem.TriagemSpec) (triagem.Decisao, error) {
	dec, err := triagem.Decidir(r.Context(), triagem.Entrada{Texto: req.Entrada, Workspace: req.Workspace}, catalogo, nil)
	if errors.Is(err, triagem.ErrResolvedorAusente) {
		return triagem.Decisao{
			PrecisaConfirmacao: true,
			Camada:             "ambiguidade",
			Justificativa:      err.Error(),
		}, nil
	}
	return dec, err
}

func carregarCatalogoTriagemStudio(loopsDir string) ([]triagem.TriagemSpec, error) {
	entries, err := os.ReadDir(loopsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var catalogo []triagem.TriagemSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		spec, ok, err := lerTriagemSpecStudio(filepath.Join(loopsDir, e.Name(), "triagem.json"))
		if err != nil {
			return nil, err
		}
		if ok {
			catalogo = append(catalogo, spec)
		}
	}
	return catalogo, nil
}

func lerTriagemSpecStudio(caminho string) (triagem.TriagemSpec, bool, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			return triagem.TriagemSpec{}, false, nil
		}
		return triagem.TriagemSpec{}, false, err
	}
	var spec triagem.TriagemSpec
	if err := json.Unmarshal(dados, &spec); err != nil {
		return triagem.TriagemSpec{}, false, err
	}
	if err := spec.Validar(); err != nil {
		return triagem.TriagemSpec{}, false, err
	}
	return spec, true, nil
}
