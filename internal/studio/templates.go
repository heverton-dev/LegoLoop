package studio

import (
	_ "embed"
	"encoding/json"
	"net/http"
)

// Template representa uma esteira pré-fabricada pronta para uso no Canvas (R-UI-5).
type Template struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	Icon        string      `json:"icon"`
	Composition Composition `json:"composition"`
}

// builtinTemplatesJSON guarda os templates canônicos como dado (templates.json), fora do código Go.
//
//go:embed templates.json
var builtinTemplatesJSON []byte

// BuiltinTemplates retorna os templates canônicos do LegoLoop (≥ 3 templates validados).
func BuiltinTemplates() []Template {
	var out []Template
	if err := json.Unmarshal(builtinTemplatesJSON, &out); err != nil {
		panic("studio: templates.json inválido: " + err.Error())
	}
	return out
}

func (s *Server) handleTemplates(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"templates": BuiltinTemplates(),
	})
}
