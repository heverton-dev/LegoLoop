package studio

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
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
	if err != nil {
		return err
	}
	sess.SetState(a2ui.SessionPaused)
	s.persistirMemoriaAndon(sess.Processo, sessID, "pausado_andon_cord")
	return s.A2UIHub().Publish(a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "andon:status",
		ProtocolVer: "v0.9",
		MessageType: "updateDataModel",
		DataModel: map[string]any{
			"andon_state":  "paused",
			"timestamp_ms": time.Now().UnixMilli(),
		},
	})
}

func (s *Server) executarAndonResume(sessID string) error {
	sess, err := s.A2UIHub().GetSession(sessID)
	if err != nil {
		return err
	}
	sess.SetState(a2ui.SessionRunning)
	s.persistirMemoriaAndon(sess.Processo, sessID, "executando")
	return s.A2UIHub().Publish(a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "andon:status",
		ProtocolVer: "v0.9",
		MessageType: "updateDataModel",
		DataModel: map[string]any{
			"andon_state":  "running",
			"timestamp_ms": time.Now().UnixMilli(),
		},
	})
}

func (s *Server) executarAndonCancel(sessID string) error {
	sess, err := s.A2UIHub().GetSession(sessID)
	if err != nil {
		return err
	}
	sess.SetState(a2ui.SessionCancelled)
	s.limparMemoriaAndon(sess.Processo, sessID)
	return s.A2UIHub().Publish(a2ui.Envelope{
		SessionID:   sessID,
		SurfaceID:   "andon:status",
		ProtocolVer: "v0.9",
		MessageType: "updateDataModel",
		DataModel: map[string]any{
			"andon_state":  "cancelled",
			"timestamp_ms": time.Now().UnixMilli(),
		},
	})
}

func (s *Server) persistirMemoriaAndon(proc, sessID, etapa string) {
	if proc == "" {
		return
	}
	slug := strings.TrimPrefix(proc, "proc:")
	dir := filepath.Join(s.ProjectDir, "processos", slug)
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return
	}
	mem := triagem.MemoriaExecucao{
		Processo:     slug,
		SessaoID:     sessID,
		PID:          os.Getpid(),
		EtapaAtual:   etapa,
		ProximoPasso: "verificar_operador",
		AtualizadoEm: time.Now().UTC(),
	}
	if err := triagem.SalvarMemoria(dir, mem); err != nil {
		return
	}
}

func (s *Server) limparMemoriaAndon(proc, sessID string) {
	if proc == "" {
		return
	}
	slug := strings.TrimPrefix(proc, "proc:")
	dir := filepath.Join(s.ProjectDir, "processos", slug)
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return
	}
	if err := triagem.ZerarSessao(dir, sessID); err != nil {
		return
	}
}
