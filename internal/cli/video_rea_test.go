package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/cli"
	"legoloop/internal/provider"
)

// fakeYtdlp substitui o yt-dlp: devolve metadados fixos e grava o áudio no destino pedido.
type fakeYtdlp struct{}

func (fakeYtdlp) Run(_ context.Context, _ string, args ...string) ([]byte, error) {
	for i, a := range args {
		if a == "-o" && i+1 < len(args) {
			if err := os.WriteFile(filepath.Join(filepath.Dir(args[i+1]), "audio.aac"), []byte("aac"), 0o644); err != nil {
				return nil, err
			}
		}
	}
	return []byte(`{"title":"","description":"","duration":42}`), nil
}

func TestVideoReaCLI_HelpAndInvalid(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, VideoRunner: fakeYtdlp{}}

	// 1. Help
	exit := cli.Run(context.Background(), []string{"video-rea", "--help"}, deps)
	if exit != 0 || !strings.Contains(stdout.String(), "Uso: legoloop video-rea") {
		t.Errorf("esperava exit 0 e usage, obteve %d", exit)
	}

	// 2. Subcomando desconhecido
	stdout.Reset()
	stderr.Reset()
	exit = cli.Run(context.Background(), []string{"video-rea", "desconhecido"}, deps)
	if exit != 2 || !strings.Contains(stderr.String(), "subcomando desconhecido") {
		t.Errorf("esperava exit 2 para desconhecido, obteve %d", exit)
	}
}

func TestVideoReaCLI_ForensicsTranscriptFile(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, VideoRunner: fakeYtdlp{}}

	videoJSON := filepath.Join(tempDir, "video.json")
	if exit := cli.Run(context.Background(), []string{
		"video-rea", "ingest",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--out", videoJSON,
	}, deps); exit != 0 {
		t.Fatalf("ingest falhou (exit %d): %s", exit, stderr.String())
	}

	transcriptPath := filepath.Join(tempDir, "transcript.txt")
	longTranscript := strings.Repeat("MiniMax H3 roda em PC fraco e gera vídeo grátis. ", 400)
	if err := os.WriteFile(transcriptPath, []byte(longTranscript), 0o644); err != nil {
		t.Fatalf("gravando transcrição: %v", err)
	}

	descPath := filepath.Join(tempDir, "desc.txt")
	if err := os.WriteFile(descPath, []byte("Adeus IAs pagas: Nova IA gera vídeos insanos em PC fraco (100% GRÁTIS)"), 0o644); err != nil {
		t.Fatalf("gravando descrição: %v", err)
	}

	dossierJSON := filepath.Join(tempDir, "dossier.json")
	stdout.Reset()
	stderr.Reset()
	exit := cli.Run(context.Background(), []string{
		"video-rea", "forensics",
		"--input", videoJSON,
		"--transcript-file", transcriptPath,
		"--desc-file", descPath,
		"--out", dossierJSON,
	}, deps)
	if exit != 0 {
		t.Fatalf("forensics com --transcript-file falhou (exit %d): %s", exit, stderr.String())
	}

	data, err := os.ReadFile(dossierJSON)
	if err != nil {
		t.Fatalf("lendo dossiê: %v", err)
	}
	var dossier struct {
		TranscribedText string `json:"transcribed_text"`
		Title           string `json:"title"`
	}
	if err := json.Unmarshal(data, &dossier); err != nil {
		t.Fatalf("parse dossiê: %v", err)
	}
	if dossier.TranscribedText != longTranscript {
		t.Errorf("transcrição do arquivo não foi preservada: obtido %d chars, esperado %d",
			len(dossier.TranscribedText), len(longTranscript))
	}
	if !strings.HasPrefix(dossier.Title, "Adeus IAs pagas") {
		t.Errorf("esperava título vindo do --desc-file, obteve %q", dossier.Title)
	}
}

func TestVideoReaCLI_PipelineReal(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, VideoRunner: fakeYtdlp{}}

	// 1. Ingest
	videoJSON := filepath.Join(tempDir, "video.json")
	exit := cli.Run(context.Background(), []string{
		"video-rea", "ingest",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--out", videoJSON,
	}, deps)
	if exit != 0 {
		t.Fatalf("ingest falhou (exit %d): %s", exit, stderr.String())
	}
	if _, err := os.Stat(videoJSON); os.IsNotExist(err) {
		t.Fatalf("arquivo de saída do ingest não foi criado: %s", videoJSON)
	}

	// 2. Forensics
	stdout.Reset()
	stderr.Reset()
	dossierJSON := filepath.Join(tempDir, "dossier.json")
	exit = cli.Run(context.Background(), []string{
		"video-rea", "forensics",
		"--input", videoJSON,
		"--transcript", "Tutorial de como editar vídeos com Claude Code e FFmpeg. Prompt: corte silêncios +40ms",
		"--out", dossierJSON,
	}, deps)
	if exit != 0 {
		t.Fatalf("forensics falhou (exit %d): %s", exit, stderr.String())
	}
	if _, err := os.Stat(dossierJSON); os.IsNotExist(err) {
		t.Fatalf("arquivo de saída do forensics não foi criado: %s", dossierJSON)
	}

	// 3. Curate
	stdout.Reset()
	stderr.Reset()
	curationJSON := filepath.Join(tempDir, "curation.json")
	exit = cli.Run(context.Background(), []string{
		"video-rea", "curate",
		"--dossier", dossierJSON,
		"--out", curationJSON,
	}, deps)
	if exit != 0 {
		t.Fatalf("curate falhou (exit %d): %s", exit, stderr.String())
	}
	if _, err := os.Stat(curationJSON); os.IsNotExist(err) {
		t.Fatalf("arquivo de saída do curate não foi criado: %s", curationJSON)
	}

	// 4. Gate
	stdout.Reset()
	stderr.Reset()
	exit = cli.Run(context.Background(), []string{
		"video-rea", "gate",
		"--report", curationJSON,
		"--approve",
	}, deps)
	if exit != 0 {
		t.Fatalf("gate falhou (exit %d): %s", exit, stderr.String())
	}

	// 5. Govern
	stdout.Reset()
	stderr.Reset()
	bundleDir := filepath.Join(tempDir, "bundle")
	exit = cli.Run(context.Background(), []string{
		"video-rea", "govern",
		"--approval", curationJSON,
		"--bundle", bundleDir,
		"--repo", tempDir,
	}, deps)
	if exit != 0 {
		t.Fatalf("govern falhou (exit %d): %s", exit, stderr.String())
	}

	// Ler curation para saber o suggested_name gerado
	curationData, err := os.ReadFile(curationJSON)
	if err != nil {
		t.Fatalf("erro lendo curation: %v", err)
	}
	var report struct {
		SuggestedName string `json:"suggested_name"`
	}
	_ = json.Unmarshal(curationData, &report)

	// 6. Book
	stdout.Reset()
	stderr.Reset()
	exit = cli.Run(context.Background(), []string{
		"video-rea", "book",
		"--bundle", filepath.Join(bundleDir, report.SuggestedName),
	}, deps)
	if exit != 0 {
		t.Fatalf("book falhou (exit %d): %s", exit, stderr.String())
	}

	// 7. Run (Full Pipeline)
	stdout.Reset()
	stderr.Reset()
	runBundleDir := filepath.Join(tempDir, "run_bundle")
	exit = cli.Run(context.Background(), []string{
		"video-rea", "run",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--bundle", runBundleDir,
		"--repo", tempDir,
	}, deps)
	if exit != 0 {
		t.Fatalf("run full falhou (exit %d): %s", exit, stderr.String())
	}
}

func TestVideoReaCLI_CurateNomeCanônico(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, VideoRunner: fakeYtdlp{}}

	videoJSON := filepath.Join(tempDir, "video.json")
	if exit := cli.Run(context.Background(), []string{
		"video-rea", "ingest",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--out", videoJSON,
	}, deps); exit != 0 {
		t.Fatalf("ingest falhou: %s", stderr.String())
	}

	dossierJSON := filepath.Join(tempDir, "dossier.json")
	if exit := cli.Run(context.Background(), []string{
		"video-rea", "forensics",
		"--input", videoJSON,
		"--transcript", "Minimx H3 rodando no Confi para gerar video gratis.",
		"--out", dossierJSON,
	}, deps); exit != 0 {
		t.Fatalf("forensics falhou: %s", stderr.String())
	}

	curationJSON := filepath.Join(tempDir, "curation.json")
	stdout.Reset()
	stderr.Reset()
	if exit := cli.Run(context.Background(), []string{
		"video-rea", "curate",
		"--dossier", dossierJSON,
		"--name", "legoloop-minimax-h3-local",
		"--out", curationJSON,
	}, deps); exit != 0 {
		t.Fatalf("curate falhou (exit %d): %s", exit, stderr.String())
	}

	data, err := os.ReadFile(curationJSON)
	if err != nil {
		t.Fatalf("lendo curadoria: %v", err)
	}
	var report struct {
		SuggestedName string `json:"suggested_name"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("parse curadoria: %v", err)
	}
	if report.SuggestedName != "legoloop-minimax-h3-local" {
		t.Errorf("esperava nome canônico sobrescrito, obteve %q", report.SuggestedName)
	}
}

func TestVideoReaCLI_SuporteAModelEEffort(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{
		Stdout:      &stdout,
		Stderr:      &stderr,
		VideoRunner: fakeYtdlp{},
		Discover: func(context.Context) provider.Inventory {
			return provider.Inventory{
				RouterURL:    "http://127.0.0.1:20128",
				RouterModels: []string{"fast-brain", "deep-brain"},
			}
		},
	}

	videoJSON := filepath.Join(tempDir, "video.json")
	if exit := cli.Run(context.Background(), []string{
		"video-rea", "ingest",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--out", videoJSON,
	}, deps); exit != 0 {
		t.Fatalf("ingest falhou: %s", stderr.String())
	}

	dossierJSON := filepath.Join(tempDir, "dossier.json")
	// Teste 1: Forensics com --model e --effort
	exit := cli.Run(context.Background(), []string{
		"video-rea", "forensics",
		"--input", videoJSON,
		"--transcript", "Testando extração com cérebro alocado usando MiniMax H3 e Claude Code.",
		"--model", "fast-brain",
		"--effort", "high",
		"--out", dossierJSON,
	}, deps)
	if exit != 0 {
		t.Fatalf("forensics com --model e --effort falhou (exit %d): %s", exit, stderr.String())
	}

	// Teste 2: Run completo com --model e --effort
	stdout.Reset()
	stderr.Reset()
	exit = cli.Run(context.Background(), []string{
		"video-rea", "run",
		"--url", "https://youtu.be/W7E7Ln0z6aY",
		"--transcript", "Pipeline completo com modelo roteado usando MiniMax H3 e Claude Code.",
		"--model", "deep-brain",
		"--effort", "medium",
		"--bundle", filepath.Join(tempDir, "bundle"),
		"--repo", tempDir,
	}, deps)
	if exit != 0 {
		t.Fatalf("run com --model e --effort falhou (exit %d): %s", exit, stderr.String())
	}

	// Teste 3: Modelo inexistente deve falhar com exit code 1
	stdout.Reset()
	stderr.Reset()
	exit = cli.Run(context.Background(), []string{
		"video-rea", "forensics",
		"--input", videoJSON,
		"--transcript", "Teste com modelo inexistente.",
		"--model", "modelo-fantasma",
		"--out", filepath.Join(tempDir, "fantasma.json"),
	}, deps)
	if exit != 1 || !strings.Contains(stderr.String(), "erro ao selecionar cérebro") {
		t.Errorf("esperava exit 1 com erro de cérebro, obteve %d: %s", exit, stderr.String())
	}
}

func TestVideoReaCLI_ForensicsRejeitaDossieSemAchados(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr, VideoRunner: fakeYtdlp{}}
	videoJSON := filepath.Join(tempDir, "video.json")
	if exit := cli.Run(context.Background(), []string{"video-rea", "ingest", "--url", "https://youtu.be/abc", "--out", videoJSON}, deps); exit != 0 {
		t.Fatalf("ingest falhou: %s", stderr.String())
	}
	out := filepath.Join(tempDir, "dossier.json")
	exit := cli.Run(context.Background(), []string{"video-rea", "forensics", "--input", videoJSON,
		"--transcript", "Bom dia, hoje o tempo está agradável.", "--out", out}, deps)
	if exit != 1 || !strings.Contains(stderr.String(), "sem achados forenses") {
		t.Fatalf("dossiê vazio deve sair com 1: exit=%d stderr=%s", exit, stderr.String())
	}
	if _, err := os.Stat(out); err == nil {
		t.Error("dossiê inválido não pode ser gravado")
	}
}
