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
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/graph"
	"legoloop/internal/interview"
	"legoloop/internal/loop"
	"legoloop/internal/plangen"
	"legoloop/internal/provider"
	"legoloop/internal/specgen"
)

//go:embed static
var staticFS embed.FS

// PipelineFile é o contrato gravado pelo botão "Materializar Esteira".
const PipelineFile = "pipeline.dag.json"

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
	// Loops é o acervo de sub-loops; nil inicializa em ~/.aidd/canonical/loops.
	Loops *loop.Store
	// CanonRoots são raízes extras da canonização, definidas pelo operador. O cliente nunca envia caminhos.
	CanonRoots []string

	// Interview é o motor socrático do Addy Osmani; se nil, inicializa sob demanda.
	Interview *interview.Engine
	// SpecGen é o gerador híbrido de especificações técnicas; se nil, inicializa sob demanda.
	SpecGen *specgen.Generator
	// PlanGen é o gerador híbrido de planos e rotas; se nil, inicializa sob demanda.
	PlanGen *plangen.Generator

	canonMu sync.Mutex

	invMu   sync.Mutex
	inv     provider.Inventory
	invTime time.Time
}

// Handler monta o roteador com o middleware de segurança.
func (s *Server) Handler() http.Handler {
	static, _ := fs.Sub(staticFS, "static")
	mux := http.NewServeMux()
	files := http.FileServerFS(static)
	for _, p := range []string{"GET /{$}", "GET /app.js", "GET /ui.js", "GET /icons.js", "GET /inventory.js", "GET /search.js", "GET /virtual.js", "GET /palette.js", "GET /slots.js", "GET /inspector.js", "GET /style.css", "GET /components/"} {
		mux.Handle(p, files)
	}
	mux.HandleFunc("GET /inventario", handlePage("inventario.html"))
	mux.HandleFunc("GET /api/suggest", s.handleSuggest)
	mux.HandleFunc("GET /api/assets", s.handleAssets)
	mux.HandleFunc("GET /api/assets/{kind}", s.handleAssetsKind)
	mux.HandleFunc("POST /api/assets/rescan", s.handleAssetsRescan)
	mux.HandleFunc("GET /api/canonize/status", s.handleCanonizeStatus)
	mux.HandleFunc("POST /api/canonize/plan", s.handleCanonizePlan)
	mux.HandleFunc("POST /api/canonize/apply", s.handleCanonizeApply)
	mux.HandleFunc("GET /api/blocks", s.handleBlocks)
	mux.HandleFunc("GET /api/inventory", s.handleInventory)
	mux.HandleFunc("GET /api/loops", s.handleLoopsList)
	mux.HandleFunc("GET /api/loops/{id}", s.handleLoopsGet)
	mux.HandleFunc("POST /api/loops", s.handleLoopsSave)
	mux.HandleFunc("POST /api/loops/{id}/duplicate", s.handleLoopsDuplicate)
	mux.HandleFunc("POST /api/loops/{id}/fork", s.handleLoopsDuplicate)
	mux.HandleFunc("DELETE /api/loops/{id}", s.handleLoopsDelete)
	mux.HandleFunc("POST /api/assets/translate", s.handleAssetTranslate)
	mux.HandleFunc("GET /api/assets/compatible", s.handleAssetCompatible)
	mux.HandleFunc("POST /api/assets/reveal", s.handleAssetReveal)
	mux.HandleFunc("POST /api/assets/import", s.handleAssetImport)
	mux.HandleFunc("POST /api/canonize/piece", s.handleCanonizePiece)
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
	return s.secure(mux)
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
		err = writeAtomic(s.ProjectDir, PipelineFile, data)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha ao gravar "+PipelineFile)
		return
	}

	// R-RUN-3: Grava a skill orquestradora (< 20 linhas) no cofre canônico
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
	if vault != "" && d != nil {
		skillDir := filepath.Join(vault, "skills", "loop-"+d.Name)
		_ = os.MkdirAll(skillDir, 0o755)
		skillPath := filepath.Join(skillDir, "SKILL.md")
		skillContent := fmt.Sprintf(`---
name: loop-%s
description: Executa o loop LegoLoop %s
---
# Loop %s

Execute esta esteira autônoma via LegoLoop:
`+"```bash"+`
legoloop run
`+"```"+`
`, d.Name, d.Name, d.Name)
		_ = os.WriteFile(skillPath, []byte(skillContent), 0o644)
	}

	// R-RUN-3: Registra o loop no acervo se disponível
	if s.Loops != nil && d != nil {
		loopObj := &loop.Loop{
			ID:          d.Name,
			Name:        d.Name,
			Description: d.Description,
			Version:     d.Version,
			Trigger:     d.Trigger,
			Nodes:       d.Nodes,
		}
		_ = s.Loops.Save(loopObj)
	}

	res.Path = PipelineFile
	res.DAG = s.maskDAG(res.DAG) // o arquivo gravado tem os caminhos reais; a resposta, não
	writeJSON(w, http.StatusOK, res)
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
