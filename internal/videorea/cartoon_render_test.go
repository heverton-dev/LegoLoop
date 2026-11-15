package videorea

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type spyBinaryExecutor struct {
	lookPathFn func(file string) (string, error)
	runFn      func(ctx context.Context, cmd string, args ...string) error
	lookedUp   []string
	lastCmd    string
	lastArgs   []string
}

func (s *spyBinaryExecutor) LookPath(file string) (string, error) {
	s.lookedUp = append(s.lookedUp, file)
	if s.lookPathFn != nil {
		return s.lookPathFn(file)
	}
	return "/usr/bin/" + file, nil
}

func (s *spyBinaryExecutor) Run(ctx context.Context, cmd string, args ...string) error {
	s.lastCmd = cmd
	s.lastArgs = append([]string(nil), args...)
	if s.runFn != nil {
		return s.runFn(ctx, cmd, args...)
	}
	return nil
}

func TestRenderCartoonScenes_ValidacaoContratos(t *testing.T) {
	ctx := context.Background()
	spy := &spyBinaryExecutor{}

	if err := RenderCartoonScenes(ctx, nil, "audio.mp3", "out.mp4", spy); !errors.Is(err, ErrRenderFailed) {
		t.Fatalf("esperava ErrRenderFailed ao passar roteiro nulo, obteve: %v", err)
	}

	vazio := &CartoonScript{Title: "Vazio", Scenes: []CartoonScene{}}
	if err := RenderCartoonScenes(ctx, vazio, "audio.mp3", "out.mp4", spy); !errors.Is(err, ErrRenderFailed) {
		t.Fatalf("esperava ErrRenderFailed ao passar roteiro sem cenas, obteve: %v", err)
	}
}

func TestRenderCartoonScenes_ExecutaComandoMuxing(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "saida.mp4")

	spy := &spyBinaryExecutor{
		runFn: func(ctx context.Context, cmd string, args ...string) error {
			return os.WriteFile(outMP4, make([]byte, 256), 0644)
		},
	}

	script := &CartoonScript{
		Title: "Teste",
		Scenes: []CartoonScene{
			{Index: 1, Narration: "Cena inicial curta", Timestamp: "00:00", VisualCue: "Abertura verde", AccentColor: "#16a34a"},
			{Index: 2, Narration: "Segunda cena final", Timestamp: "00:04", VisualCue: "Climax azul", AccentColor: "#2563eb"},
		},
	}

	if err := RenderCartoonScenes(ctx, script, "audio.mp3", outMP4, spy); err != nil {
		t.Fatalf("erro inesperado ao renderizar cenas: %v", err)
	}

	if len(spy.lookedUp) != 1 || spy.lookedUp[0] != "ffmpeg" {
		t.Errorf("LookPath = %v, want [ffmpeg]", spy.lookedUp)
	}
	if !strings.Contains(strings.ToLower(spy.lastCmd), "ffmpeg") {
		t.Errorf("comando inesperado disparado: %s", spy.lastCmd)
	}
	joined := strings.Join(spy.lastArgs, " ")
	for _, want := range []string{
		"-filter_complex",
		"0x16a34a",
		"0x2563eb",
		"d=4",
		"concat=n=2",
		"1080x1920",
		"audio.mp3",
		outMP4,
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("argumentos ffmpeg nao usam cena %q: %s", want, joined)
		}
	}
}

func TestRenderCartoonScenes_ErroQuandoFFmpegAusenteNoPATH(t *testing.T) {
	ctx := context.Background()
	t.Setenv("PATH", t.TempDir())
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "ausente.mp4")

	script := &CartoonScript{
		Title:  "Teste",
		Scenes: []CartoonScene{{Index: 1, Narration: "Cena 1", Timestamp: "00:00"}},
	}

	err := RenderCartoonScenes(ctx, script, "audio.mp3", outMP4, nil)
	if !errors.Is(err, ErrRenderFailed) {
		t.Fatalf("esperava ErrRenderFailed quando ffmpeg ausente no PATH, obteve: %v", err)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "ffmpeg") {
		t.Errorf("mensagem deve citar ffmpeg: %v", err)
	}
	if _, statErr := os.Stat(outMP4); !os.IsNotExist(statErr) {
		t.Errorf("nao deve criar MP4 placeholder em %s", outMP4)
	}
}

func TestRenderCartoonScenes_ErroQuandoExecutorFalhaRemoveParcial(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "parcial.mp4")
	_ = os.WriteFile(outMP4, []byte("lixo"), 0644)

	spy := &spyBinaryExecutor{
		runFn: func(ctx context.Context, cmd string, args ...string) error {
			return errors.New("exit status 1")
		},
	}
	script := &CartoonScript{
		Title:  "Falha",
		Scenes: []CartoonScene{{Index: 1, Narration: "Cena unica", Timestamp: "00:00"}},
	}

	err := RenderCartoonScenes(ctx, script, "audio.mp3", outMP4, spy)
	if !errors.Is(err, ErrRenderFailed) {
		t.Fatalf("esperava ErrRenderFailed, obteve: %v", err)
	}
	if _, statErr := os.Stat(outMP4); !os.IsNotExist(statErr) {
		t.Errorf("arquivo parcial deveria ser removido apos falha")
	}
}

func TestGenerateCartoonFixture3s_ErroQuandoFFmpegAusente(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	outMP4 := filepath.Join(t.TempDir(), "sem_ffmpeg.mp4")

	err := GenerateCartoonFixture3s(outMP4)
	if !errors.Is(err, ErrRenderFailed) {
		t.Fatalf("esperava ErrRenderFailed quando ffmpeg ausente no PATH, obteve: %v", err)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "ffmpeg") {
		t.Errorf("mensagem de erro deveria mencionar ffmpeg: %v", err)
	}
	if _, statErr := os.Stat(outMP4); !os.IsNotExist(statErr) {
		t.Errorf("nao deve gravar MP4 falso de 32 bytes em %s", outMP4)
	}
}

func TestGenerateCartoonFixture3s(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg nao disponivel neste host")
	}
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "fixture_3s.mp4")

	if err := GenerateCartoonFixture3s(outMP4); err != nil {
		t.Fatalf("falha ao gerar fixture de 3s: %v", err)
	}

	data, err := os.ReadFile(outMP4)
	if err != nil {
		t.Fatalf("erro ao ler fixture gerada: %v", err)
	}
	if len(data) <= 512 {
		t.Fatalf("fixture pequena demais (%d bytes) — parece placeholder", len(data))
	}
	if !bytes.Contains(data[:32], []byte("ftyp")) {
		t.Errorf("arquivo gerado nao contem assinatura ftyp de MP4")
	}
}

func TestRenderCartoonScenes_RealFFmpegComCenas(t *testing.T) {
	bin, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg nao disponivel neste host")
	}
	tmpDir := t.TempDir()
	audioPath := filepath.Join(tmpDir, "voz.wav")
	cmd := exec.Command(bin, "-y", "-f", "lavfi", "-i", "sine=frequency=440:sample_rate=44100:duration=2", audioPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha preparando audio de teste: %v (%s)", err, string(out))
	}

	outMP4 := filepath.Join(tmpDir, "cartoon_cenas.mp4")
	script := &CartoonScript{
		Title: "Varginha",
		Scenes: []CartoonScene{
			{Index: 1, Timestamp: "00:00", Narration: "Cena um", VisualCue: "Cena um", AccentColor: "#1a1a24"},
			{Index: 2, Timestamp: "00:01", Narration: "Cena dois", VisualCue: "Cena dois", AccentColor: "#16a34a"},
		},
	}

	if err := RenderCartoonScenes(context.Background(), script, audioPath, outMP4, nil); err != nil {
		t.Fatalf("RenderCartoonScenes real falhou: %v", err)
	}
	data, err := os.ReadFile(outMP4)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) <= 1024 || !bytes.Contains(data[:32], []byte("ftyp")) {
		t.Fatalf("video MP4 multi-cena invalido: len=%d", len(data))
	}
}
