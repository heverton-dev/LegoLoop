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
	"io/fs"
	"legoloop/internal/caminhos"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/brain"
	"legoloop/internal/dag"
	"legoloop/internal/fabricar"
	"legoloop/internal/graph"
	"legoloop/internal/interview"
	"legoloop/internal/plangen"
	"legoloop/internal/processo"
	"legoloop/internal/provider"
	"legoloop/internal/reputation"
	"legoloop/internal/specgen"
	"legoloop/internal/studiodoc"
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

	matrixOnce sync.Once
	matrixReg  *reputation.RegistryMatrix

	canonMu    sync.Mutex
	fabricarMu sync.Mutex

	invMu   sync.Mutex
	inv     provider.Inventory
	invTime time.Time
}

// Handler monta o roteador com o middleware de segurança.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.registerStaticRoutes(mux)
	s.registerAPIRoutes(mux)
	s.registerProcessoRoutes(mux)
	return s.secure(mux)
}

func (s *Server) registerStaticRoutes(mux *http.ServeMux) {
	static, _ := fs.Sub(staticFS, "static")
	files := http.FileServerFS(static)
	for _, p := range []string{"GET /{$}", "GET /app.js", "GET /ui.js", "GET /icons.js", "GET /inventory.js", "GET /search.js", "GET /virtual.js", "GET /palette.js", "GET /slots.js", "GET /inspector.js", "GET /fabricar.js", "GET /processos.js", "GET /style.css", "GET /components/"} {
		mux.Handle(p, files)
	}
	mux.HandleFunc("GET /inventario", handlePage("inventario.html"))
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

	mux.HandleFunc("GET /api/loops", deprecada(s.handleProcessosLista))
	mux.HandleFunc("GET /api/loops/{id}", deprecada(s.handleProcessosObter))
	mux.HandleFunc("POST /api/loops", deprecada(s.handleProcessosGravar))
	mux.HandleFunc("POST /api/loops/{id}/duplicate", deprecada(s.handleProcessosDuplicar))
	mux.HandleFunc("POST /api/loops/{id}/fork", deprecada(s.handleProcessosDuplicar))
	mux.HandleFunc("DELETE /api/loops/{id}", deprecada(s.handleProcessosExcluir))
}

func deprecada(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Deprecation", "true")
		h(w, r)
	}
}

func (s *Server) handleDAGStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "idle",
		"engine": "orchestrix",
	})
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
		ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, MaxHeaderBytes: 64 << 10}
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
	data, err := d.Marshal()
	if err == nil {
		err = writeAtomic(s.ProjectDir, ArquivoDAGGravado, data)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao gravar "+ArquivoDAGGravado)
		return
	}

	s.writeOrchestratorSkill(d)
	s.registrarProcesso(d)

	res.Path = ArquivoDAGGravado
	res.DAG = s.maskDAG(res.DAG) // o arquivo gravado tem os caminhos reais; a resposta, não
	writeJSON(w, http.StatusOK, res)
}

// writeOrchestratorSkill grava a skill orquestradora (< 20 linhas) no Cofre canônico (R-RUN-3).
func (s *Server) writeOrchestratorSkill(d *dag.DAG) {
	vault := s.Vault
	if vault == "" {
		home := s.Home
		if home == "" {
			home, _ = os.UserHomeDir()
		}
		if home != "" {
			vault = filepath.Join(home, ".aidd", "canonical")
		}
	}
	if vault == "" || d == nil {
		return
	}
	skillDir := filepath.Join(vault, "skills", "loop-"+d.Name)
	_ = os.MkdirAll(skillDir, 0o755)
	skillContent := fmt.Sprintf(`---
name: loop-%s
description: Executa o Processo %s
---
# Processo %s

Execute este Processo autônomo via LegoLoop:
`+"```bash"+`
legoloop run
`+"```"+`
`, d.Name, d.Name, d.Name)
	_ = os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0o644)
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
