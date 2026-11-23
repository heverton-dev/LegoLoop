package studio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

func (s *Server) dispatchAndonAction(ctx context.Context, act a2ui.Action) error {
	switch act.Name {
	case "andon_cord_pause":
		return s.executarAndonPause(act.SessionID)
	case "andon_cord_resume":
		return s.executarAndonResume(act.SessionID)
	case "andon_cord_cancel":
		return s.executarAndonCancel(act.SessionID)
	default:
		return nil
	}
}

func (s *Server) executarAndonPause(sessID string) error {
	sess, err := s.A2UIHub().GetSession(sessID)
	if err == nil {
		sess.SetState(a2ui.SessionPaused)
		if errMem := s.persistirMemoriaAndon(sess.Processo, sessID, "pausado_andon_cord"); errMem != nil {
			return errMem
		}
		return s.publicarAndonState(sessID, "paused")
	}
	proc, errDisco := s.buscarProcessoPorSessao(sessID)
	if errDisco != nil {
		return fmt.Errorf("sessão %s não encontrada: %w", sessID, err)
	}
	return s.persistirMemoriaAndon(proc, sessID, "pausado_andon_cord")
}

func (s *Server) executarAndonResume(sessID string) error {
	sess, err := s.A2UIHub().GetSession(sessID)
	if err == nil {
		sess.SetState(a2ui.SessionRunning)
		if errMem := s.persistirMemoriaAndon(sess.Processo, sessID, "executando"); errMem != nil {
			return errMem
		}
		return s.publicarAndonState(sessID, "running")
	}
	proc, errDisco := s.buscarProcessoPorSessao(sessID)
	if errDisco != nil {
		return fmt.Errorf("sessão %s não encontrada: %w", sessID, err)
	}
	return s.persistirMemoriaAndon(proc, sessID, "executando")
}

func (s *Server) executarAndonCancel(sessID string) error {
	sess, err := s.A2UIHub().GetSession(sessID)
	if err == nil {
		sess.SetState(a2ui.SessionCancelled)
		if errMem := s.limparMemoriaAndon(sess.Processo, sessID); errMem != nil {
			return errMem
		}
		return s.publicarAndonState(sessID, "cancelled")
	}
	proc, errDisco := s.buscarProcessoPorSessao(sessID)
	if errDisco != nil {
		return fmt.Errorf("sessão %s não encontrada: %w", sessID, err)
	}
	return s.limparMemoriaAndon(proc, sessID)
}

func (s *Server) publicarAndonState(sessID, state string) error {
	return s.A2UIHub().Publish(a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "andon:status",
		ProtocolVer: "v0.9",
		MessageType: "updateDataModel",
		DataModel: map[string]any{
			"andon_state":  state,
			"timestamp_ms": time.Now().UnixMilli(),
		},
	})
}

func (s *Server) buscarProcessoPorSessao(sessID string) (string, error) {
	if sessID == "" {
		return "", fmt.Errorf("sessão vazia")
	}
	dirProc := caminhos.DirProcessosEm(s.ProjectDir)
	entries, err := os.ReadDir(dirProc)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		slug := entry.Name()
		arq := filepath.Join(dirProc, slug, triagem.ArquivoMemoria)
		mems := extrairMemoriasDeArquivo(arq, slug)
		for _, m := range mems {
			if m.SessaoID == sessID {
				return slug, nil
			}
		}
	}
	return "", fmt.Errorf("sessão %q não encontrada em disco", sessID)
}

func (s *Server) persistirMemoriaAndon(proc, sessID, etapa string) error {
	if proc == "" {
		return nil
	}
	slug := strings.TrimPrefix(proc, "proc:")
	dir := filepath.Join(caminhos.DirProcessosEm(s.ProjectDir), slug)
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return nil
	}
	mem, err := triagem.LerMemoria(dir, sessID)
	if err != nil {
		mem = triagem.MemoriaExecucao{
			Processo:     slug,
			SessaoID:     sessID,
			PID:          os.Getpid(),
			ProximoPasso: "verificar_operador",
		}
	}
	mem.Processo = slug
	mem.SessaoID = sessID
	mem.EtapaAtual = etapa
	if mem.PID <= 0 {
		mem.PID = os.Getpid()
	}
	mem.AtualizadoEm = time.Now().UTC()
	return triagem.SalvarMemoria(dir, mem)
}

func (s *Server) limparMemoriaAndon(proc, sessID string) error {
	if proc == "" {
		return nil
	}
	slug := strings.TrimPrefix(proc, "proc:")
	dir := filepath.Join(caminhos.DirProcessosEm(s.ProjectDir), slug)
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return nil
	}
	return triagem.ZerarSessao(dir, sessID)
}
