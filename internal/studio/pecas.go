package studio

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"legoloop/internal/assets"
	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
)

// PecaProcesso é a peça invocável de um Processo no Studio (card arrastável).
type PecaProcesso struct {
	ID      string `json:"id"`
	Nome    string `json:"nome"`
	Kind    string `json:"kind"`
	Invocar string `json:"invocar"`
}

// pecasDeProcessos registra automaticamente uma peça invocável por Processo
// com contrato DAG no disco, sem registro manual.
func (s *Server) pecasDeProcessos() []PecaProcesso {
	out := []PecaProcesso{}
	if s.ProjectDir == "" {
		return out
	}
	dirProc := filepath.Join(s.ProjectDir, caminhos.DirProcessos)
	ents, err := os.ReadDir(dirProc)
	if err != nil {
		return out
	}
	for _, e := range ents {
		if !e.IsDir() || !isSafeSkillSlug(e.Name()) {
			continue
		}
		procDir := filepath.Join(dirProc, e.Name())
		if _, err := os.Stat(caminhos.CaminhoDAGLeitura(procDir)); err != nil {
			continue
		}
		out = append(out, PecaProcesso{
			ID: e.Name(), Nome: e.Name(), Kind: string(assets.KindProcesso),
			Invocar: "legoloop exec " + e.Name(),
		})
	}
	return out
}

func (s *Server) handlePecasListar(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"pecas": s.pecasDeProcessos()})
}

type pecaExecReq struct {
	AllowScripts bool `json:"allow_scripts"`
}

// handlePecasExec invoca o Processo pelo slug e executa o MESMO contrato DAG
// usado pela via CLI. Nós script exigem consentimento prévio (R-RUN-2).
func (s *Server) handlePecasExec(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("id")
	if !isSafeSkillSlug(slug) {
		writeError(w, http.StatusBadRequest, "id de peça inválido")
		return
	}
	d, ok := s.carregarDAGPeca(slug)
	if !ok {
		writeError(w, http.StatusNotFound, "processo não encontrado: "+slug)
		return
	}
	var req pecaExecReq
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	if !validateRunContext(w, d, req.AllowScripts) {
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "streaming sse não suportado"})
		return
	}
	startSSE(w, flusher)
	s.dispatchRunExecution(r, d, req.AllowScripts, newSSESender(w, flusher))
}

// carregarDAGPeca lê o contrato DAG do Processo sem escapar da pasta canônica.
func (s *Server) carregarDAGPeca(slug string) (*dag.DAG, bool) {
	procDir := filepath.Join(s.ProjectDir, caminhos.DirProcessos, slug)
	raw, err := os.ReadFile(caminhos.CaminhoDAGLeitura(procDir))
	if err != nil {
		return nil, false
	}
	d, err := dag.Parse(raw)
	if err != nil {
		return nil, false
	}
	return d, true
}
