package studio

import (
	"net/http"

	"legoloop/internal/assets"
	"legoloop/internal/suggest"
)

var suggestKinds = map[string]bool{"skill": true, "script": true, "mcp": true, "subloop": true}

type suggestResponse struct {
	Gates     []suggest.Gate     `json:"gates"`
	Scripts   []assets.Item      `json:"scripts"`
	Hooks     suggestHooks       `json:"hooks"`
	Neighbors []suggest.Neighbor `json:"neighbors"`
}

type suggestHooks struct {
	Pre     []assets.Item `json:"pre"`
	Post    []assets.Item `json:"post"`
	Skipped int           `json:"skipped"`
}

// handleSuggest: GET /api/suggest?block=<id do catálogo> ou ?kind=<tipo>&artifact=<nome>.
// Tudo é calculado só do inventário (sem LLM) e sai sem caminhos absolutos nem comandos brutos.
func (s *Server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	q := r.URL.Query()
	ctx, code, msg := s.suggestContext(q.Get("block"), q.Get("kind"), q.Get("artifact"))
	if msg != "" {
		writeError(w, code, msg)
		return
	}
	res := suggest.Suggest(ctx, s.Index, s.snapshot(r.Context()))
	writeJSON(w, http.StatusOK, suggestResponse{
		Gates:     res.Gates,
		Scripts:   s.presentAll(res.Scripts),
		Hooks:     suggestHooks{Pre: s.presentAll(res.HooksPre), Post: s.presentAll(res.HooksPost), Skipped: res.SkippedHooks},
		Neighbors: res.Neighbors,
	})
}

func (s *Server) suggestContext(block, kind, artifact string) (suggest.Context, int, string) {
	switch {
	case block != "" && kind != "":
		return suggest.Context{}, http.StatusBadRequest, "informe block OU kind, não os dois"
	case block != "":
		b, ok := s.Index.Block(block)
		if !ok {
			return suggest.Context{}, http.StatusNotFound, "Peça não existe no catálogo"
		}
		return suggest.Context{Block: &b, Artifact: b.Artifact}, 0, ""
	case kind == "":
		return suggest.Context{}, http.StatusBadRequest, "informe block ou kind"
	case !suggestKinds[kind]:
		return suggest.Context{}, http.StatusBadRequest, "kind desconhecido"
	case artifact == "" || len(artifact) > maxArtifactLen:
		return suggest.Context{}, http.StatusBadRequest, "artifact é obrigatório (até 128 caracteres)"
	}
	return suggest.Context{Kind: kind, Artifact: artifact}, 0, ""
}

func (s *Server) presentAll(items []assets.Item) []assets.Item {
	out := make([]assets.Item, len(items))
	for i, it := range items {
		out[i] = s.present(it)
	}
	return out
}
