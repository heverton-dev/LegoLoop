package studio

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

type ProcessoAtivoItem struct {
	ID           string            `json:"id"`
	Processo     string            `json:"processo"`
	EtapaAtual   string            `json:"etapa_atual,omitempty"`
	ProximoPasso string            `json:"proximo_passo,omitempty"`
	State        a2ui.SessionState `json:"state"`
	Surfaces     int               `json:"surfaces"`
	Subscribers  int               `json:"subscribers"`
	URL          string            `json:"url"`
	TimestampMs  int64             `json:"timestamp_ms"`
}

func (s *Server) handleProcessosAtivos(w http.ResponseWriter, r *http.Request) {
	vistas := make(map[string]struct{})
	itens := s.coletarSessoesHub(vistas)
	itens = append(itens, s.coletarSessoesDisco(vistas)...)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(itens); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) coletarSessoesHub(vistas map[string]struct{}) []ProcessoAtivoItem {
	sessions := s.A2UIHub().ListSessions()
	itens := make([]ProcessoAtivoItem, 0, len(sessions))
	now := time.Now().UnixMilli()
	for _, sess := range sessions {
		vistas[sess.ID] = struct{}{}
		itens = append(itens, s.montarItemHub(sess, now))
	}
	return itens
}

func (s *Server) montarItemHub(sess a2ui.SessionSummary, now int64) ProcessoAtivoItem {
	slug := strings.TrimPrefix(sess.Processo, "proc:")
	dir := filepath.Join(caminhos.DirProcessosEm(s.ProjectDir), slug)
	etapa, prox := "", ""
	if mem, err := triagem.LerMemoria(dir, sess.ID); err == nil {
		etapa = mem.EtapaAtual
		prox = mem.ProximoPasso
	}
	return ProcessoAtivoItem{
		ID:           sess.ID,
		Processo:     sess.Processo,
		EtapaAtual:   etapa,
		ProximoPasso: prox,
		State:        sess.State,
		Surfaces:     sess.Surfaces,
		Subscribers:  sess.Subscribers,
		URL:          "/sessao/" + sess.ID,
		TimestampMs:  now,
	}
}

func (s *Server) coletarSessoesDisco(vistas map[string]struct{}) []ProcessoAtivoItem {
	dirProc := caminhos.DirProcessosEm(s.ProjectDir)
	entries, err := os.ReadDir(dirProc)
	if err != nil {
		return nil
	}
	var itens []ProcessoAtivoItem
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		slug := entry.Name()
		arq := filepath.Join(dirProc, slug, triagem.ArquivoMemoria)
		mems := extrairMemoriasDeArquivo(arq, slug)
		for _, m := range mems {
			if _, ok := vistas[m.SessaoID]; ok {
				continue
			}
			vistas[m.SessaoID] = struct{}{}
			itens = append(itens, itemDeMemoriaDisco(m))
		}
	}
	return itens
}

func extrairMemoriasDeArquivo(caminho, slug string) []triagem.MemoriaExecucao {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		return nil
	}
	var env triagem.EstadoMemoria
	if err := json.Unmarshal(dados, &env); err == nil && len(env.Sessoes) > 0 {
		itens := make([]triagem.MemoriaExecucao, 0, len(env.Sessoes))
		for id, m := range env.Sessoes {
			if id == "legado" && m.EtapaAtual == "" && m.PID <= 0 {
				continue
			}
			if m.SessaoID == "" {
				m.SessaoID = id
			}
			if m.Processo == "" {
				m.Processo = slug
			}
			itens = append(itens, m)
		}
		return itens
	}
	return extrairMemoriaSimples(dados, slug)
}

func extrairMemoriaSimples(dados []byte, slug string) []triagem.MemoriaExecucao {
	var m triagem.MemoriaExecucao
	if err := json.Unmarshal(dados, &m); err != nil || (m.SessaoID == "" && m.EtapaAtual == "") {
		return nil
	}
	if m.Processo == "" {
		m.Processo = slug
	}
	if m.SessaoID == "" {
		m.SessaoID = slug
	}
	return []triagem.MemoriaExecucao{m}
}

func itemDeMemoriaDisco(m triagem.MemoriaExecucao) ProcessoAtivoItem {
	st := a2ui.SessionRunning
	etapaLower := strings.ToLower(m.EtapaAtual)
	if strings.Contains(etapaLower, "pausad") || strings.Contains(etapaLower, "paused") {
		st = a2ui.SessionPaused
	} else if strings.Contains(etapaLower, "cancel") {
		st = a2ui.SessionCancelled
	}
	ts := m.AtualizadoEm.UnixMilli()
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}
	return ProcessoAtivoItem{
		ID:           m.SessaoID,
		Processo:     m.Processo,
		EtapaAtual:   m.EtapaAtual,
		ProximoPasso: m.ProximoPasso,
		State:        st,
		Surfaces:     0,
		Subscribers:  0,
		URL:          "/sessao/" + m.SessaoID,
		TimestampMs:  ts,
	}
}

func (s *Server) handleProcessoAndonPause(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonPause(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "paused"})
}

func (s *Server) handleProcessoAndonResume(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonResume(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "running"})
}

func (s *Server) handleProcessoAndonCancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.executarAndonCancel(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "state": "cancelled"})
}
