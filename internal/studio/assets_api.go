package studio

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/canon"
)

// maxListItems limita cada resposta de lista; o total real vai em "total".
const maxListItems = 20000

// canonTimeout é o teto de uma canonização aplicada (ela termina mesmo se a aba fechar).
const canonTimeout = 10 * time.Minute

// maxHookDisplay é o tamanho, em runas, do comando de hook exibido (após mascarar caminhos).
const maxHookDisplay = 200

// Abas do inventário que não vêm do scanner de disco.
const (
	kindModel = "model"
	kindLoop  = "loop"
)

func errNoInventory(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "inventário indisponível neste servidor")
}

// displayPath nunca devolve caminho absoluto: troca o projeto por <projeto> e o HOME por ~.
func (s *Server) displayPath(p string) string {
	if p == "" {
		return ""
	}
	for _, base := range [][2]string{{s.ProjectDir, "<projeto>"}, {s.home(), "~"}} {
		if base[0] == "" {
			continue
		}
		if rel, err := filepath.Rel(base[0], p); err == nil && (rel == "." || filepath.IsLocal(rel)) {
			return filepath.ToSlash(filepath.Join(base[1], rel))
		}
	}
	return filepath.Base(p)
}

// maskPaths troca, dentro de um texto livre (ex.: comando de hook), o projeto por <projeto> e o
// HOME por ~, em qualquer variante de barra. O projeto vem primeiro: costuma ficar dentro do HOME.
func (s *Server) maskPaths(text string) string {
	var pairs []string
	for _, base := range [][2]string{{s.ProjectDir, "<projeto>"}, {s.home(), "~"}} {
		if base[0] == "" {
			continue
		}
		fwd := filepath.ToSlash(base[0])
		for _, v := range []string{base[0], fwd, strings.ReplaceAll(fwd, "/", `\`), strings.ReplaceAll(base[0], `\`, `\\`)} {
			pairs = append(pairs, v, base[1])
		}
	}
	return strings.NewReplacer(pairs...).Replace(text)
}

func truncateRunes(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}

func (s *Server) home() string {
	if s.Home != "" {
		return s.Home
	}
	h, _ := os.UserHomeDir()
	return h
}

// snapshot devolve o inventário publicado; na primeira chamada faz a varredura inicial.
func (s *Server) snapshot(ctx context.Context) assets.Snapshot {
	if snap := s.Assets.Snapshot(); snap.Version != 0 {
		return snap
	}
	return s.Assets.Refresh(ctx)
}

func (s *Server) counts(ctx context.Context, snap assets.Snapshot) map[string]int {
	loopCount := 0
	if store, err := s.loopStore(); err == nil {
		if list, err := store.List(); err == nil {
			loopCount = len(list)
		}
	}
	out := map[string]int{kindModel: len(s.inventory(ctx).Models()), kindLoop: loopCount}
	for k, n := range snap.Counts {
		out[string(k)] = n
	}
	return out
}

func (s *Server) summary(ctx context.Context, snap assets.Snapshot) map[string]any {
	return map[string]any{"version": snap.Version, "scanned_at": snap.ScannedAt, "counts": s.counts(ctx, snap)}
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	writeJSON(w, http.StatusOK, s.summary(r.Context(), s.snapshot(r.Context())))
}

func (s *Server) handleAssetsRescan(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	var req struct{}
	if !decodeStrict(w, r, &req) {
		return
	}
	writeJSON(w, http.StatusOK, s.summary(r.Context(), s.Assets.Refresh(r.Context())))
}

func (s *Server) handleAssetsKind(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	kind := r.PathValue("kind")
	items, ok := s.itemsOf(r.Context(), kind)
	if !ok {
		writeError(w, http.StatusNotFound, "categoria desconhecida")
		return
	}
	total := len(items)
	if total > maxListItems {
		items = items[:maxListItems]
	}
	writeJSON(w, http.StatusOK, map[string]any{"kind": kind, "items": items, "total": total, "truncated": total > len(items)})
}

// itemsOf devolve os itens da aba já com caminhos de exibição; lista vazia, nunca null.
func (s *Server) itemsOf(ctx context.Context, kind string) ([]assets.Item, bool) {
	switch kind {
	case kindModel:
		out := []assets.Item{}
		for _, m := range s.inventory(ctx).Models() {
			out = append(out, assets.Item{Kind: kindModel, Name: m.ID, Origin: string(m.Source), Scope: "system"})
		}
		return out, true
	case kindLoop:
		return s.loopItems(), true
	}
	if !slices.Contains(assets.Kinds(), assets.Kind(kind)) {
		return nil, false
	}
	out := []assets.Item{}
	for _, it := range s.snapshot(ctx).Of(assets.Kind(kind)) {
		out = append(out, s.present(it))
	}
	return out, true
}

func (s *Server) loopItems() []assets.Item {
	out := []assets.Item{}
	if store, err := s.loopStore(); err == nil {
		list, _ := store.List()
		for _, l := range list {
			out = append(out, assets.Item{
				ID:          l.ID,
				Kind:        assets.KindLoop,
				Name:        l.Name,
				Description: l.Description,
				Origin:      "acervo",
				Scope:       "global",
				Meta: map[string]string{
					"nodes":    strconv.Itoa(l.NodesCount),
					"artifact": l.FinalArtifact,
					"trigger":  l.TriggerType,
				},
			})
		}
	}
	return out
}

// present devolve o item como o navegador pode vê-lo: sem caminho absoluto e com o comando de
// hook mascarado e curto. O comando bruto (Command) nunca é serializado.
func (s *Server) present(it assets.Item) assets.Item {
	it.Path = s.displayPath(it.Path)
	it.Description = s.maskPaths(it.Description)
	if it.Kind == assets.KindHook {
		it.Description = truncateRunes(it.Description, maxHookDisplay)
	}
	return it
}

// resolver devolve o inventário para resolver referências de peças e slots.
func (s *Server) resolver(ctx context.Context) Resolver {
	if s.Assets == nil {
		return noInventory{}
	}
	return s.snapshot(ctx)
}

// maskDAG troca caminhos locais por ~ e <projeto> no DAG devolvido ao navegador. Só a resposta é
// mascarada: o pipeline.dag.json gravado em disco guarda os caminhos reais, que o runner precisa.
func (s *Server) maskDAG(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}
	return json.RawMessage(s.maskPaths(string(raw)))
}

// ---------- Canonização ----------

type canonGroup struct {
	SHA256    string   `json:"sha256"`
	Canonical string   `json:"canonical"`
	Links     []string `json:"links"`
	Saved     int64    `json:"saved_bytes"`
}

type canonPlanResponse struct {
	Vault           string       `json:"vault"`
	FilesScanned    int          `json:"files_scanned"`
	Groups          []canonGroup `json:"groups"`
	SavedBytes      int64        `json:"saved_bytes"`
	SkippedGitRoots []string     `json:"skipped_git_roots"`
}

// canonRoots junta as raízes do inventário com as do operador, sem repetir.
func (s *Server) canonRoots() []string {
	var out []string
	for _, r := range append(s.Assets.Roots(), s.CanonRoots...) {
		if !slices.Contains(out, r) {
			out = append(out, r)
		}
	}
	return out
}

func (s *Server) vault() (string, error) {
	if s.Vault != "" {
		return s.Vault, nil
	}
	return canon.DefaultVault()
}

func (s *Server) buildPlan(ctx context.Context) (canon.Plan, canonPlanResponse, error) {
	vault, err := s.vault()
	if err != nil {
		return canon.Plan{}, canonPlanResponse{}, err
	}
	roots := s.canonRoots()
	files, err := canon.HashTree(ctx, canon.ScanOptions{Roots: roots, Exclude: []string{vault}, Workers: canon.DefaultWorkers})
	if err != nil {
		return canon.Plan{}, canonPlanResponse{}, err
	}
	plan := canon.BuildPlan(files, vault)
	resp := canonPlanResponse{Vault: s.displayPath(vault), FilesScanned: len(files), SavedBytes: plan.SavedBytes,
		Groups: []canonGroup{}, SkippedGitRoots: []string{}}
	for _, g := range plan.Groups {
		cg := canonGroup{SHA256: g.SHA256, Canonical: s.displayPath(g.Canonical.Path), Links: []string{}, Saved: g.SavedBytes}
		for _, l := range g.Links {
			cg.Links = append(cg.Links, s.displayPath(l))
		}
		resp.Groups = append(resp.Groups, cg)
	}
	for _, root := range roots {
		if _, inRepo := canon.GitRoot(root, s.home()); inRepo {
			resp.SkippedGitRoots = append(resp.SkippedGitRoots, s.displayPath(root))
		}
	}
	return plan, resp, nil
}

func (s *Server) handleCanonizePlan(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	var req struct{}
	if !decodeStrict(w, r, &req) {
		return
	}
	_, resp, err := s.buildPlan(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao varrer os arquivos para o plano")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

type canonFailure struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

type canonApplyResponse struct {
	Linked     int            `json:"linked"`
	Methods    map[string]int `json:"methods"`
	Failed     []canonFailure `json:"failed"`
	SavedBytes int64          `json:"saved_bytes"`
}

func (s *Server) handleCanonizeApply(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	var req struct {
		Apply bool `json:"apply"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	if !req.Apply {
		writeError(w, http.StatusBadRequest, "confirmação ausente: envie apply:true")
		return
	}
	if !s.canonMu.TryLock() {
		writeError(w, http.StatusConflict, "já existe uma canonização em andamento")
		return
	}
	defer s.canonMu.Unlock()
	// A operação mexe em arquivos do usuário: não é interrompida pelo fechamento da aba.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), canonTimeout)
	defer cancel()
	resp, err := s.applyPlan(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao varrer os arquivos para o plano")
		return
	}
	s.Assets.Refresh(ctx)
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) applyPlan(ctx context.Context) (canonApplyResponse, error) {
	plan, _, err := s.buildPlan(ctx)
	if err != nil {
		return canonApplyResponse{}, err
	}
	rep := canon.Apply(ctx, plan, canon.DefaultLinker())
	resp := canonApplyResponse{Linked: len(rep.Linked), Methods: map[string]int{}, Failed: []canonFailure{}, SavedBytes: rep.SavedBytes}
	for _, l := range rep.Linked {
		resp.Methods[string(l.Method)]++
	}
	for _, f := range rep.Failed {
		resp.Failed = append(resp.Failed, canonFailure{Path: s.displayPath(f.Path), Error: f.Error})
	}
	return resp, nil
}

// ---------- Páginas ----------

// handleAssetCompatible busca peças ontologicamente compatíveis com o item informado.
func (s *Server) handleAssetCompatible(w http.ResponseWriter, r *http.Request) {
	if s.Assets == nil {
		errNoInventory(w)
		return
	}
	id := r.URL.Query().Get("id")
	snap := s.Assets.Snapshot()
	var target assets.Item
	found := false
	for _, it := range snap.Items {
		if it.ID == id || it.Path == id || it.Name == id {
			target = it
			found = true
			break
		}
	}
	if !found {
		writeJSON(w, http.StatusOK, map[string]any{"compatible": []assets.CompatibleItem{}})
		return
	}
	comp := assets.FindCompatible(target, snap.Items)
	writeJSON(w, http.StatusOK, map[string]any{"compatible": comp})
}

// handlePage serve um HTML embutido (rota limpa, sem extensão).
func handlePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		data, err := fs.ReadFile(staticFS, "static/"+name)
		if err != nil {
			writeError(w, http.StatusNotFound, "página não encontrada")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	}
}

// decodeStrict decodifica o corpo JSON recusando campos desconhecidos e dados extras.
func decodeStrict(w http.ResponseWriter, r *http.Request, dst any) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	var tooBig *http.MaxBytesError
	switch {
	case errors.As(err, &tooBig):
		writeError(w, http.StatusRequestEntityTooLarge, "corpo excede 1 MiB")
		return false
	case err != nil:
		writeError(w, http.StatusBadRequest, "corpo inválido: "+err.Error())
		return false
	case !atEOF(dec):
		writeError(w, http.StatusBadRequest, "corpo inválido: dados extras após o objeto JSON")
		return false
	}
	return true
}
