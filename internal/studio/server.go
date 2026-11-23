// Package studio serve o Looping Studio (canvas N8N) e sua API local.
// Só escuta em loopback, valida Host (anti DNS-rebinding), exige JSON + mesma origem
// em POST (anti-CSRF), limita corpo e taxa, e nunca expõe caminhos ou segredos locais.
package studio

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/assets"
	"github.com/heverton-dev/LegoLoop/internal/brain"
	"github.com/heverton-dev/LegoLoop/internal/dag"
	"github.com/heverton-dev/LegoLoop/internal/fabricar"
	"github.com/heverton-dev/LegoLoop/internal/graph"
	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/plangen"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/provider"
	"github.com/heverton-dev/LegoLoop/internal/reputation"
	"github.com/heverton-dev/LegoLoop/internal/runner"
	"github.com/heverton-dev/LegoLoop/internal/specgen"
	"github.com/heverton-dev/LegoLoop/internal/studiodoc"
)

//go:embed static
var staticFS embed.FS

// ArquivoDAGGravado é o contrato gravado pelo botão "Materializar Processo".
const ArquivoDAGGravado = caminhos.ArquivoDAG

// inventoryTTL evita sondar a rede a cada clique no canvas.
const inventoryTTL = 10 * time.Second

// Server agrega dependências do Studio; sem estado global.
type Server struct {
	Index         *graph.Index
	ProjectDir    string
	Discover      func(context.Context) provider.Inventory
	Host          string
	RatePerSecond float64
	Burst         int

	// Inventário de ativos (skills, MCPs, hooks...). Nil desliga as rotas /api/assets e /api/canonize.
	Assets *assets.Watcher
	// Home é só para exibição (~) e para o limite de busca de repositórios git na canonização.
	Home string
	// Vault é o cofre canônico; vazio usa ~/.aidd/canonical.
	Vault string
	// Processos é o acervo de Processos; nil inicializa no Cofre canônico.
	Processos *studiodoc.Store
	// CanonRoots são raízes extras da canonização, definidas pelo operador. O cliente nunca envia caminhos.
	CanonRoots []string

	// Interview é o motor socrático do Addy Osmani; se nil, inicializa sob demanda.
	Interview *interview.Engine
	// SpecGen é o gerador híbrido de especificações técnicas; se nil, inicializa sob demanda.
	SpecGen *specgen.Generator
	// PlanGen é o gerador híbrido de planos e rotas; se nil, inicializa sob demanda.
	PlanGen *plangen.Generator
	// Fabricar cria processos novos (POST /api/fabricar/run); nil usa a fábrica real com o 9Router do inventário.
	Fabricar FabricarFunc
	// BrainOpts ajusta o cérebro da fábrica real (URL do Ollama, prazo, cliente HTTP); usado em testes.
	BrainOpts []brain.Option
	// FabricarMatriz substitui a matriz do cofre na fábrica (nil = registry_matrix.json do cofre); usado em testes.
	FabricarMatriz fabricar.MatrizFunc
	// StoreOverride força o store do processo diretamente (usado em testes).
	StoreOverride *processo.Store
	// RevealLauncher substitui a abertura do explorador do SO (usado em testes).
	RevealLauncher func(path string, isDir bool) error

	// A2UI é o barramento de sessões declarativas; nil inicializa sob demanda.
	A2UI         *a2ui.Hub
	a2uiOnce     sync.Once
	concepcaoMu  sync.Mutex
	concepcaoMap map[string]string

	matrixOnce sync.Once
	matrixReg  *reputation.RegistryMatrix

	canonMu    sync.Mutex
	fabricarMu sync.Mutex
	favMu      sync.Mutex

	dagMu    sync.Mutex
	dagState dagRunState

	invMu   sync.Mutex
	inv     provider.Inventory
	invTime time.Time
}

type dagRunState struct {
	Status     string             `json:"status"`
	Engine     string             `json:"engine"`
	DAG        string             `json:"dag,omitempty"`
	NodesTotal int                `json:"nodes_total,omitempty"`
	NodesDone  int                `json:"nodes_done,omitempty"`
	ActiveNode string             `json:"active_node,omitempty"`
	DurationMs int64              `json:"duration_ms,omitempty"`
	Tokens     *runner.TokenUsage `json:"tokens,omitempty"`
	Error      string             `json:"error,omitempty"`
	ActiveRuns int                `json:"-"`
}

// Handler monta o roteador com o middleware de segurança.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.registerStaticRoutes(mux)
	s.registerAPIRoutes(mux)
	s.registerProcessoRoutes(mux)
	s.registerSwarmRoutes(mux)
	s.registerA2UIRoutes(mux)
	return s.secure(mux)
}

const faviconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><rect width="32" height="32" rx="7" fill="#0f172a"/><rect x="7" y="11" width="18" height="12" rx="2" fill="#38bdf8"/><rect x="10" y="7" width="4" height="4" rx="1" fill="#38bdf8"/><rect x="18" y="7" width="4" height="4" rx="1" fill="#38bdf8"/></svg>`

func handleFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(faviconSVG))
}

func (s *Server) registerStaticRoutes(mux *http.ServeMux) {
	static, _ := fs.Sub(staticFS, "static")
	files := http.FileServerFS(static)
	for _, p := range []string{"GET /{$}", "GET /app.js", "GET /ui.js", "GET /icons.js", "GET /inventory.js", "GET /search.js", "GET /virtual.js", "GET /palette.js", "GET /slots.js", "GET /inspector.js", "GET /fabricar.js", "GET /processos.js", "GET /processos-vivos.js", "GET /style.css", "GET /components/"} {
		mux.Handle(p, files)
	}
	mux.HandleFunc("GET /favicon.ico", handleFavicon)
	mux.HandleFunc("GET /processos", handlePage("processos.html"))
	mux.HandleFunc("GET /inventario", handlePage("inventario.html"))
	mux.HandleFunc("GET /sessao/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.garantirSessaoLaboratorio(r.PathValue("id"))
		handlePage("sessao.html")(w, r)
	})
}

func (s *Server) registerAPIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/suggest", s.handleSuggest)
	mux.HandleFunc("GET /api/assets", s.handleAssets)
	mux.HandleFunc("GET /api/assets/{kind}", s.handleAssetsKind)
	mux.HandleFunc("POST /api/assets/rescan", s.handleAssetsRescan)
	mux.HandleFunc("GET /api/canonize/status", s.handleCanonizeStatus)
	mux.HandleFunc("POST /api/canonize/plan", s.handleCanonizePlan)
	mux.HandleFunc("POST /api/canonize/apply", s.handleCanonizeApply)
	mux.HandleFunc("GET /api/blocks", s.handleBlocks)
	mux.HandleFunc("GET /api/inventory", s.handleInventory)
	mux.HandleFunc("POST /api/assets/translate", s.handleAssetTranslate)
	mux.HandleFunc("GET /api/assets/inspect", s.handleAssetInspect)
	mux.HandleFunc("GET /api/assets/compatible", s.handleAssetCompatible)
	mux.HandleFunc("POST /api/assets/reveal", s.handleAssetReveal)
	mux.HandleFunc("POST /api/assets/import", s.handleAssetImport)
	mux.HandleFunc("POST /api/canonize/piece", s.handleCanonizePiece)
	mux.HandleFunc("POST /api/fabricar/run", s.handleFabricarRun)
	mux.HandleFunc("POST /api/assets/{kind}/{id}/favorite", s.handlePieceFavorite)
	mux.HandleFunc("POST /api/assets/{kind}/{id}/upgrade", s.handlePieceUpgrade)
	mux.HandleFunc("DELETE /api/assets/{kind}/{id}", s.handlePieceDelete)
	mux.HandleFunc("GET /api/route", s.handleRoute)
	mux.HandleFunc("POST /api/validate", s.handleValidate)
	mux.HandleFunc("POST /api/materialize", s.handleMaterialize)
	mux.HandleFunc("POST /api/run", s.handleRunDAG)
	mux.HandleFunc("GET /api/templates", s.handleTemplates)
	mux.HandleFunc("GET /api/mcps/catalog", s.handleMCPCatalog)
	mux.HandleFunc("POST /api/mcps/install", s.handleMCPInstall)
	mux.HandleFunc("POST /api/interview/start", s.handleInterviewStart)
	mux.HandleFunc("POST /api/interview/message", s.handleInterviewMessage)
	mux.HandleFunc("POST /api/interview/finalize", s.handleInterviewFinalize)
	mux.HandleFunc("POST /api/spec/generate", s.handleSpecGenerate)
	mux.HandleFunc("POST /api/plan/generate", s.handlePlanGenerate)
	mux.HandleFunc("POST /api/triagem", s.handleTriagem)
}

func (s *Server) registerProcessoRoutes(mux *http.ServeMux) {
	// Canonical Processos & DAG APIs (ASD-STE100)
	mux.HandleFunc("GET /api/processos", s.handleProcessosListar)
	mux.HandleFunc("GET /api/processos/{id}", s.handleProcessosObter)
	mux.HandleFunc("POST /api/processos", s.handleProcessosGravar)
	mux.HandleFunc("POST /api/processos/{id}/duplicate", s.handleProcessosDuplicar)
	mux.HandleFunc("POST /api/processos/{id}/fork", s.handleProcessosDuplicar)
	mux.HandleFunc("DELETE /api/processos/{id}", s.handleProcessosExcluir)
	mux.HandleFunc("POST /api/processos/refinar", s.handleProcessosRefinar)
	mux.HandleFunc("POST /api/processos/remover", s.handleProcessosRemover)
	mux.HandleFunc("POST /api/processos/compor", s.handleProcessosCompor)
	mux.HandleFunc("GET /api/dag", s.handleDAGStatus)
	mux.HandleFunc("GET /api/dag/{id}", s.handleDAGGet)
	mux.HandleFunc("GET /api/processos/ativos", s.handleProcessosAtivos)
	mux.HandleFunc("POST /api/processos/{id}/andon", s.handleProcessoAndonPause)
	mux.HandleFunc("POST /api/processos/{id}/resume", s.handleProcessoAndonResume)
	mux.HandleFunc("POST /api/processos/{id}/cancel", s.handleProcessoAndonCancel)
	mux.HandleFunc("GET /api/pecas", s.handlePecasListar)
	mux.HandleFunc("POST /api/pecas/{id}/exec", s.handlePecasExec)
}

func deprecada(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Deprecation", "true")
		h(w, r)
	}
}

func (s *Server) handleDAGStatus(w http.ResponseWriter, r *http.Request) {
	st := s.snapshotDAGState()
	writeJSON(w, http.StatusOK, st)
}

func (s *Server) snapshotDAGState() dagRunState {
	s.dagMu.Lock()
	defer s.dagMu.Unlock()
	st := s.dagState
	if st.Status == "" {
		st.Status = "idle"
	}
	st.Engine = "orchestrix"
	if st.Error != "" {
		st.Error = s.maskPaths(st.Error)
	}
	return st
}

func (s *Server) handleDAGGet(w http.ResponseWriter, r *http.Request) {
	s.handleProcessosObter(w, r)
}

// ListenAndServe escuta apenas em loopback e encerra graciosamente quando ctx termina.
// Se ready != nil, recebe o endereço efetivo (útil com porta 0).
func (s *Server) ListenAndServe(ctx context.Context, addr string, ready chan<- string) error {
	if err := requireLoopback(addr); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.Host = ln.Addr().String()
	srv := &http.Server{Handler: s.Handler(), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second, WriteTimeout: 90 * time.Second, MaxHeaderBytes: 64 << 10}
	if ready != nil {
		ready <- s.Host
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(ln) }()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func requireLoopback(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host == "localhost" {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return errors.New("studio: só é permitido escutar em loopback (127.0.0.1, ::1, localhost)")
}

func (s *Server) handleBlocks(w http.ResponseWriter, _ *http.Request) {
	types := []string{}
	for _, t := range graph.ArtifactTypes() {
		types = append(types, string(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"blocks": s.Index.Blocks(), "types": types, "matrix": s.Index.Matrix(),
	})
}

func (s *Server) inventory(ctx context.Context) provider.Inventory {
	s.invMu.Lock()
	defer s.invMu.Unlock()
	if s.invTime.IsZero() || time.Since(s.invTime) > inventoryTTL {
		s.inv, s.invTime = s.Discover(ctx), time.Now()
	}
	return s.inv
}

// handleInventory expõe só o necessário ao canvas: modelos verdes e nomes de harness.
func (s *Server) handleInventory(w http.ResponseWriter, r *http.Request) {
	inv := s.inventory(r.Context())
	clis := []string{}
	for _, c := range inv.CLIs {
		clis = append(clis, c.Name)
	}
	models := inv.Models()
	if models == nil {
		models = []provider.Model{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"models": models, "clis": clis, "api_keys": inv.APIKeys, "router_online": inv.RouterURL != "",
	})
}

func (s *Server) handleRoute(w http.ResponseWriter, r *http.Request) {
	from := graph.ArtifactType(r.URL.Query().Get("from"))
	to := graph.ArtifactType(r.URL.Query().Get("to"))
	path, cost, err := s.Index.Route(from, to)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	ids := []string{}
	for _, b := range path {
		ids = append(ids, b.ID)
	}
	writeJSON(w, http.StatusOK, map[string]any{"path": ids, "cost": cost})
}

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	c, ok := decodeComposition(w, r)
	if !ok {
		return
	}
	res, _ := compose(s.Index, s.inventory(r.Context()), s.resolver(r.Context()), c)
	res.DAG = s.maskDAG(res.DAG)
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleMaterialize(w http.ResponseWriter, r *http.Request) {
	c, ok := decodeComposition(w, r)
	if !ok {
		return
	}
	res, d := compose(s.Index, s.inventory(r.Context()), s.resolver(r.Context()), c)
	if !res.OK {
		writeJSON(w, http.StatusUnprocessableEntity, res)
		return
	}
	if !isSafeSkillSlug(d.Name) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("%v: %q", errInvalidSkillSlug, d.Name))
		return
	}
	data, err := d.Marshal()
	if err == nil {
		err = writeAtomic(s.ProjectDir, ArquivoDAGGravado, data)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao gravar "+ArquivoDAGGravado)
		return
	}
	if err := s.writeOrchestratorSkill(d); err != nil {
		writeError(w, http.StatusInternalServerError, s.maskPaths(err.Error()))
		return
	}
	s.registrarProcesso(d)
	res.Path = ArquivoDAGGravado
	res.DAG = s.maskDAG(res.DAG) // o arquivo gravado tem os caminhos reais; a resposta, não
	writeJSON(w, http.StatusOK, res)
}

var errInvalidSkillSlug = errors.New("nome de processo inválido para skill")

func isSafeSkillSlug(name string) bool {
	if name == "" {
		return false
	}
	for i := range len(name) {
		c := name[i]
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func (s *Server) resolveSkillVaultDir() (string, error) {
	if s.Vault != "" {
		return filepath.Clean(s.Vault), nil
	}
	home := s.Home
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("falha ao obter diretório home: %w", err)
		}
	}
	return filepath.Join(filepath.Clean(home), ".aidd", "canonical"), nil
}

func ensureDirInRoot(root *os.Root, rel string) error {
	clean := filepath.Clean(rel)
	cur := ""
	for part := range strings.SplitSeq(clean, string(filepath.Separator)) {
		if cur == "" {
			cur = part
		} else {
			cur = filepath.Join(cur, part)
		}
		info, err := root.Lstat(cur)
		if err == nil {
			if isDisallowedLink(info.Mode()) {
				return fmt.Errorf("symlink ou junction não permitido: %q", filepath.ToSlash(cur))
			}
			if !info.IsDir() {
				return fmt.Errorf("caminho não é diretório: %q", filepath.ToSlash(cur))
			}
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := root.Mkdir(cur, 0o755); err != nil {
			return err
		}
	}
	return checkNoLinksRoot(root, clean)
}

// writeOrchestratorSkill grava a skill orquestradora (< 20 linhas) no Cofre canônico (R-RUN-3).
func (s *Server) writeOrchestratorSkill(d *dag.DAG) error {
	if d == nil {
		return errors.New("dag nulo")
	}
	if !isSafeSkillSlug(d.Name) {
		return fmt.Errorf("%w: %q", errInvalidSkillSlug, d.Name)
	}
	vault, err := s.resolveSkillVaultDir()
	if err != nil {
		return err
	}
	if info, err := os.Lstat(vault); err == nil && isDisallowedLink(info.Mode()) {
		return errors.New("symlink ou junction não permitido no cofre")
	}
	if err := os.MkdirAll(vault, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(vault)
	if err != nil {
		return err
	}
	defer root.Close()
	skillDir := filepath.Join("skills", "processo-"+d.Name)
	if err := ensureDirInRoot(root, skillDir); err != nil {
		return err
	}
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if info, err := root.Lstat(skillFile); err == nil && isDisallowedLink(info.Mode()) {
		return fmt.Errorf("symlink ou junction não permitido: %q", filepath.ToSlash(skillFile))
	}
	skillContent := fmt.Sprintf("---\nname: processo-%s\ndescription: Executa o Processo %s\n---\n# Processo %s\n\nExecute este Processo autônomo via LegoLoop:\n```bash\nlegoloop run\n```\n", d.Name, d.Name, d.Name)
	return root.WriteFile(skillFile, []byte(skillContent), 0o644)
}

// registrarProcesso registra o Processo no acervo, se disponível (R-RUN-3).
func (s *Server) registrarProcesso(d *dag.DAG) {
	if s.Processos == nil || d == nil {
		return
	}
	_ = s.Processos.Save(&studiodoc.Processo{
		ID:          d.Name,
		Name:        d.Name,
		Description: d.Description,
		Version:     d.Version,
		Trigger:     d.Trigger,
		Nodes:       d.Nodes,
	})
}

// writeAtomic grava via arquivo temporário + rename, confinado ao diretório do projeto.
func writeAtomic(dir, name string, data []byte) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	tmp := name + ".tmp"
	if err := root.WriteFile(tmp, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return root.Rename(tmp, name)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	// Cabeçalhos já enviados: um erro aqui só significa cliente desconectado.
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
