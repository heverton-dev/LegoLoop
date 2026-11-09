package videorea

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockCommandExecutor struct {
	lastCmd  string
	lastArgs []string
	err      error
}

func (m *mockCommandExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	m.lastCmd = cmd
	m.lastArgs = args
	return m.err
}

func TestSynthesizeEdgeTTS_TextoVazio(t *testing.T) {
	ctx := context.Background()
	mock := &mockCommandExecutor{}

	err := SynthesizeEdgeTTS(ctx, "   ", "pt-BR-AntonioNeural", "scratch/audio.mp3", mock)
	if !errors.Is(err, ErrTTSFailed) {
		t.Fatalf("esperava ErrTTSFailed ao passar texto vazio, obteve: %v", err)
	}
}

func TestSynthesizeEdgeTTS_ComandoExecutadoCorretamente(t *testing.T) {
	ctx := context.Background()
	mock := &mockCommandExecutor{}
	tmpDir := t.TempDir()
	outAudio := filepath.Join(tmpDir, "audio.mp3")

	texto := "Janeiro de 1996. [reflexivo] Três moças avistam algo estranho."
	err := SynthesizeEdgeTTS(ctx, texto, "pt-BR-AntonioNeural", outAudio, mock)
	if err != nil {
		t.Fatalf("erro inesperado na síntese Edge-TTS: %v", err)
	}

	if mock.lastCmd != "python" && mock.lastCmd != "edge-tts" {
		t.Errorf("comando inesperado: %s", mock.lastCmd)
	}
	encontrouVoz := false
	for _, a := range mock.lastArgs {
		if a == "pt-BR-AntonioNeural" {
			encontrouVoz = true
			break
		}
	}
	if !encontrouVoz {
		t.Errorf("esperava argumento de voz 'pt-BR-AntonioNeural' na chamada")
	}
}

func TestSynthesizeEdgeTTS_FalhaLimpaArquivo(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	corrompido := filepath.Join(tmpDir, "corrompido.mp3")
	_ = os.WriteFile(corrompido, []byte("lixo"), 0644)

	mock := &mockCommandExecutor{err: errors.New("falha de conexao")}
	err := SynthesizeEdgeTTS(ctx, "teste de erro", "pt-BR-AntonioNeural", corrompido, mock)
	if !errors.Is(err, ErrTTSFailed) {
		t.Fatalf("esperava ErrTTSFailed, obteve: %v", err)
	}
	if _, statErr := os.Stat(corrompido); !os.IsNotExist(statErr) {
		t.Errorf("arquivo corrompido deveria ter sido removido em caso de falha")
	}
}
