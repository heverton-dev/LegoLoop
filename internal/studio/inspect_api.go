package studio

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"legoloop/internal/assets"
	"legoloop/internal/caminhos"
	"legoloop/internal/graph"
)

type assetInspectResponse struct {
	OK             bool              `json:"ok"`
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Kind           string            `json:"kind"`
	Description    string            `json:"description"`
	WhatItIs       string            `json:"what_it_is"`
	WhatItDoes     string            `json:"what_it_does"`
	WhatItDelivers string            `json:"what_it_delivers"`
	Origin         string            `json:"origin,omitempty"`
	Scope          string            `json:"scope,omitempty"`
	Path           string            `json:"path,omitempty"`
	Artifact       string            `json:"artifact,omitempty"`
	Gate           string            `json:"gate,omitempty"`
	Consumes       []string          `json:"consumes,omitempty"`
	Produces       string            `json:"produces,omitempty"`
	SuggestedModel string            `json:"suggested_model,omitempty"`
	Cost           int               `json:"cost,omitempty"`
	SHA256         string            `json:"sha256,omitempty"`
	Tier           string            `json:"tier,omitempty"`
	Score          string            `json:"score,omitempty"`
	Guilda         string            `json:"guilda,omitempty"`
	Papel          string            `json:"papel,omitempty"`
	Meta           map[string]string `json:"meta,omitempty"`
	Sections       []string          `json:"sections,omitempty"`
	LineCount      int               `json:"line_count,omitempty"`
	SizeBytes      int               `json:"size_bytes,omitempty"`
	SourceContent  string            `json:"source_content,omitempty"`
}

func (s *Server) handleAssetInspect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := strings.TrimSpace(q.Get("id"))
	name := strings.TrimSpace(q.Get("name"))
	kind := strings.TrimSpace(q.Get("kind"))
	rawPath := strings.TrimSpace(q.Get("path"))

	if blk, ok := s.findBlockForInspect(id, name); ok {
		resp := buildBlockInspect(blk)
		s.enrichBlockWithInventory(r, &resp, blk)
		writeJSON(w, http.StatusOK, resp)
		return
	}
	if s.Assets == nil {
		writeError(w, http.StatusNotFound, "inventário indisponível")
		return
	}
	it, ok := s.findItemForInspect(r, id, name, kind, rawPath)
	if !ok {
		writeError(w, http.StatusNotFound, "peça não encontrada no inventário")
		return
	}
	resp := s.buildItemInspect(it, rawPath)
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) findBlockForInspect(id, name string) (graph.Block, bool) {
	if s.Index != nil {
		if b, ok := matchBlockInSlice(s.Index.Blocks(), id, name); ok {
			return b, true
		}
	}
	if b, ok := matchBlockInSlice(graph.TriggerBlocks(), id, name); ok {
		return b, true
	}
	return matchBlockInSlice(graph.ControlBlocks(), id, name)
}

func matchBlockInSlice(list []graph.Block, id, name string) (graph.Block, bool) {
	for _, b := range list {
		if id != "" && (strings.EqualFold(b.ID, id) || strings.EqualFold(b.Kind, id) || strings.EqualFold(b.Ref, id)) {
			return b, true
		}
		if name != "" && strings.EqualFold(b.Label, name) {
			return b, true
		}
	}
	return graph.Block{}, false
}

func buildBlockInspect(b graph.Block) assetInspectResponse {
	cons := make([]string, len(b.Consumes))
	for i, c := range b.Consumes {
		cons[i] = string(c)
	}
	what := fmtBlockWhat(b)
	return assetInspectResponse{
		OK:             true,
		ID:             b.ID,
		Name:           b.Label,
		Kind:           firstNonEmpty(b.Kind, "skill"),
		Description:    what,
		WhatItIs:       what,
		WhatItDoes:     fmtBlockDoes(b),
		WhatItDelivers: "Entrega o artefato " + string(b.Produces),
		Origin:         "catalogo-nativo",
		Scope:          "core",
		Path:           b.Ref,
		Artifact:       b.Artifact,
		Gate:           b.Gate,
		Consumes:       cons,
		Produces:       string(b.Produces),
		SuggestedModel: b.SuggestedModel,
		Cost:           b.Cost,
	}
}

func fmtBlockWhat(b graph.Block) string {
	if strings.HasPrefix(b.Kind, "trigger_") {
		return "Gatilho nativo do motor LegoLoop (" + b.Label + ")."
	}
	return "Peça tipada do catálogo LegoLoop (" + b.Label + ")."
}

func fmtBlockDoes(b graph.Block) string {
	if len(b.Consumes) > 0 {
		parts := make([]string, len(b.Consumes))
		for i, c := range b.Consumes {
			parts[i] = string(c)
		}
		return "Consome " + strings.Join(parts, ", ") + " e produz " + string(b.Produces) + "."
	}
	return "Dispara o processo produzindo " + string(b.Produces) + "."
}

func (s *Server) findItemForInspect(r *http.Request, id, name, kind, rawPath string) (assets.Item, bool) {
	snap := s.snapshot(r.Context())
	for _, it := range snap.Items {
		if matchInspectItem(it, id, name, kind, s.displayPath(it.Path), rawPath) {
			return it, true
		}
	}
	for _, it := range s.processoItems() {
		if matchInspectItem(it, id, name, kind, it.Path, rawPath) {
			return it, true
		}
	}
	return assets.Item{}, false
}

func matchInspectItem(it assets.Item, id, name, kind, maskedPath, rawPath string) bool {
	if kind != "" && !strings.EqualFold(string(it.Kind), kind) && kind != caminhos.TipoSubprocessoAntigo {
		return false
	}
	if id != "" && (strings.EqualFold(it.ID, id) || strings.EqualFold(it.Name, id) || strings.EqualFold(it.Meta[caminhos.ChaveMetaProcessoID], id)) {
		return true
	}
	if name != "" && strings.EqualFold(it.Name, name) {
		return true
	}
	if rawPath != "" && (it.Path == rawPath || maskedPath == rawPath) {
		return true
	}
	return false
}

func (s *Server) buildItemInspect(it assets.Item, rawPath string) assetInspectResponse {
	presented := s.present(it)
	resp := assetInspectResponse{
		OK:          true,
		ID:          presented.ID,
		Name:        presented.Name,
		Kind:        string(presented.Kind),
		Description: presented.Description,
		Origin:      presented.Origin,
		Scope:       presented.Scope,
		Path:        presented.Path,
		Artifact:    presented.Meta["artifact"],
		Gate:        presented.Meta["gate"],
		Tier:        presented.Meta["tier"],
		Score:       presented.Meta["score"],
		Guilda:      presented.Meta["guilda"],
		Papel:       firstNonEmpty(presented.Meta["papel_nome"], presented.Meta["papel_primario"]),
		Meta:        presented.Meta,
	}
	s.enrichInspectFromDisk(&resp, it.Path, rawPath)
	return resp
}

func (s *Server) enrichInspectFromDisk(resp *assetInspectResponse, itemPath, queryPath string) {
	target := firstNonEmpty(itemPath, queryPath)
	if target == "" {
		return
	}
	cleanPath, err := s.validateSafeAssetPath(target)
	if err != nil {
		return
	}
	content, err := s.readAssetInRevealRoots(cleanPath)
	if err != nil {
		return
	}
	sum := sha256.Sum256(content)
	resp.SHA256 = hex.EncodeToString(sum[:])
	resp.SizeBytes = len(content)
	text := string(content)
	resp.LineCount = strings.Count(text, "\n") + 1
	desc, body, secs := parseInspectDocument(text)
	if desc != "" && resp.Description == "" {
		resp.Description = desc
	}
	resp.Sections = secs
	resp.WhatItIs = firstNonEmpty(desc, resp.Description, "Peça física: "+resp.Name)
	resp.WhatItDoes = extractInspectSummary(body, resp.Kind, resp.Name)
	resp.WhatItDelivers = firstNonEmpty(resp.Artifact, "Artefato estruturado conforme contrato da peça")
	if len(text) > 16384 {
		text = text[:16384] + "\n... (conteúdo truncado em 16 KiB)"
	}
	resp.SourceContent = s.maskAssetText(text)
}

func (s *Server) enrichBlockWithInventory(r *http.Request, resp *assetInspectResponse, blk graph.Block) {
	if s.Assets == nil {
		return
	}
	it, ok := s.findItemForInspect(r, blk.ID, blk.Label, "", blk.Ref)
	if !ok {
		it, ok = s.findItemForInspect(r, blk.Ref, blk.Label, "", "")
	}
	if !ok {
		return
	}
	inv := s.buildItemInspect(it, blk.Ref)
	if inv.Description != "" {
		resp.Description = inv.Description
		resp.WhatItIs = firstNonEmpty(inv.WhatItIs, inv.Description)
	}
	if inv.WhatItDoes != "" {
		resp.WhatItDoes = inv.WhatItDoes
	}
	resp.SHA256 = inv.SHA256
	resp.SizeBytes = inv.SizeBytes
	resp.LineCount = inv.LineCount
	resp.Sections = inv.Sections
	resp.SourceContent = inv.SourceContent
	resp.Tier = inv.Tier
	resp.Score = inv.Score
	resp.Guilda = inv.Guilda
	resp.Papel = inv.Papel
	if inv.Path != "" {
		resp.Path = inv.Path
	}
}

func parseInspectDocument(text string) (string, string, []string) {
	desc := ""
	body := text
	if strings.HasPrefix(text, "---\n") || strings.HasPrefix(text, "---\r\n") {
		if fm, rest, ok := splitFrontmatter(text); ok {
			desc = extractYAMLField(fm, "description")
			body = rest
		}
	}
	var secs []string
	for line := range strings.SplitSeq(body, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "#") {
			h := strings.TrimSpace(strings.TrimLeft(t, "#"))
			if h != "" && len(secs) < 10 {
				secs = append(secs, h)
			}
		}
	}
	return desc, body, secs
}

func splitFrontmatter(text string) (string, string, bool) {
	norm := strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(norm, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return "", text, false
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return strings.Join(lines[1:i], "\n"), strings.Join(lines[i+1:], "\n"), true
		}
	}
	return "", text, false
}

func extractYAMLField(fm, key string) string {
	prefix := key + ":"
	for line := range strings.SplitSeq(fm, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, prefix) {
			val := strings.TrimSpace(strings.TrimPrefix(t, prefix))
			return strings.Trim(val, `"'`)
		}
	}
	return ""
}

func extractInspectSummary(body, kind, name string) string {
	for line := range strings.SplitSeq(body, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "```") || strings.HasPrefix(t, "---") {
			continue
		}
		if len(t) > 180 {
			return t[:180] + "..."
		}
		return t
	}
	return "Execução de " + kind + " (" + name + ") segundo o conteúdo fonte."
}
