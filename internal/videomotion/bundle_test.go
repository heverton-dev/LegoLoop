package videomotion_test

import (
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/videomotion"
)

func TestGenerateBroadcastBundle(t *testing.T) {
	tempDir := t.TempDir()

	words := []videomotion.WordItem{
		{Index: 0, Word: "Este", StartMs: 200, EndMs: 600},
		{Index: 1, Word: "é", StartMs: 650, EndMs: 900},
		{Index: 2, Word: "o", StartMs: 950, EndMs: 1100},
		{Index: 3, Word: "LegoLoop", StartMs: 1150, EndMs: 2200, IsPunchWord: true},
	}

	cfg := videomotion.BundleConfig{
		JobID:       "job_test_01",
		OutputDir:   tempDir,
		Title:       "LegoLoop Motion Video",
		CTA:         "Comente 'LEGO' para receber o processo completo!",
		Words:       words,
		AspectRatio: videomotion.AspectVertical,
	}

	bundle, err := videomotion.GenerateBroadcastBundle(cfg)
	if err != nil {
		t.Fatalf("GenerateBroadcastBundle falhou: %v", err)
	}

	// 1. Validar que as pastas estruturais foram criadas
	expectedDirs := []string{
		filepath.Join(tempDir, "job_test_01", "video"),
		filepath.Join(tempDir, "job_test_01", "captions"),
		filepath.Join(tempDir, "job_test_01", "copy"),
		filepath.Join(tempDir, "job_test_01", "studio"),
	}
	for _, d := range expectedDirs {
		if fi, err := os.Stat(d); os.IsNotExist(err) || !fi.IsDir() {
			t.Errorf("diretório esperado não existe: %s", d)
		}
	}

	// 2. Validar que subtitles.srt e subtitles.vtt foram gerados
	srtPath := filepath.Join(tempDir, "job_test_01", "captions", "subtitles.srt")
	srtContent, err := os.ReadFile(srtPath)
	if err != nil {
		t.Fatalf("subtitles.srt não foi gerado: %v", err)
	}
	if !containsStr(string(srtContent), "00:00:00,200 --> 00:00:02,200") {
		t.Errorf("timestamp SRT incorreto:\n%s", string(srtContent))
	}

	vttPath := filepath.Join(tempDir, "job_test_01", "captions", "subtitles.vtt")
	vttContent, err := os.ReadFile(vttPath)
	if err != nil {
		t.Fatalf("subtitles.vtt não foi gerado: %v", err)
	}
	if !containsStr(string(vttContent), "WEBVTT") {
		t.Errorf("cabeçalho WEBVTT ausente no arquivo VTT")
	}

	// 3. Validar que POST_COPY.md foi gerado com CTAs e tags
	copyPath := filepath.Join(tempDir, "job_test_01", "copy", "POST_COPY.md")
	copyContent, err := os.ReadFile(copyPath)
	if err != nil {
		t.Fatalf("POST_COPY.md não foi gerado: %v", err)
	}
	if !containsStr(string(copyContent), "Comente 'LEGO'") {
		t.Errorf("CTA ausente no POST_COPY.md")
	}

	if bundle.BundleDir != filepath.Join(tempDir, "job_test_01") {
		t.Errorf("BundleDir esperado %s, obteve %s", filepath.Join(tempDir, "job_test_01"), bundle.BundleDir)
	}
}

func TestGenerateBroadcastBundle_Validation(t *testing.T) {
	_, err := videomotion.GenerateBroadcastBundle(videomotion.BundleConfig{})
	if err == nil {
		t.Fatal("esperava erro para configuração vazia")
	}
}
