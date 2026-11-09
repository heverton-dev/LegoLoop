package videorea

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderCartoonScenes_ValidacaoContratos(t *testing.T) {
	ctx := context.Background()
	mock := &mockCommandExecutor{}

	// Roteiro nulo
	err := RenderCartoonScenes(ctx, nil, "audio.mp3", "out.mp4", mock)
	if err == nil {
		t.Fatalf("esperava erro ao passar roteiro nulo")
	}

	// Sem cenas
	vazio := &CartoonScript{Title: "Vazio", Scenes: []CartoonScene{}}
	err = RenderCartoonScenes(ctx, vazio, "audio.mp3", "out.mp4", mock)
	if err == nil {
		t.Fatalf("esperava erro ao passar roteiro sem cenas")
	}
}

func TestRenderCartoonScenes_ExecutaComandoMuxing(t *testing.T) {
	ctx := context.Background()
	mock := &mockCommandExecutor{}
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "saida.mp4")

	script := &CartoonScript{
		Title: "Teste",
		Scenes: []CartoonScene{
			{Index: 1, Narration: "Cena 1", Timestamp: "00:00"},
			{Index: 2, Narration: "Cena 2", Timestamp: "00:10"},
		},
	}

	err := RenderCartoonScenes(ctx, script, "audio.mp3", outMP4, mock)
	if err != nil {
		t.Fatalf("erro inesperado ao renderizar cenas: %v", err)
	}

	if mock.lastCmd != "ffmpeg" && mock.lastCmd != "node" && mock.lastCmd != "python" {
		t.Errorf("comando inesperado disparado: %s", mock.lastCmd)
	}
}

func TestGenerateCartoonFixture3s(t *testing.T) {
	tmpDir := t.TempDir()
	outMP4 := filepath.Join(tmpDir, "fixture_3s.mp4")

	err := GenerateCartoonFixture3s(outMP4)
	if err != nil {
		t.Fatalf("falha ao gerar fixture de 3s: %v", err)
	}

	data, err := os.ReadFile(outMP4)
	if err != nil || len(data) < 32 {
		t.Fatalf("arquivo de fixture não gerado ou muito pequeno: %d bytes", len(data))
	}

	// Verifica se contém assinatura MP4 (ftyp) nos primeiros bytes
	isMP4 := false
	for i := 0; i < 20; i++ {
		if i+4 <= len(data) && string(data[i:i+4]) == "ftyp" {
			isMP4 = true
			break
		}
	}
	if !isMP4 {
		t.Errorf("arquivo gerado não contém a assinatura ftyp de MP4")
	}
}
