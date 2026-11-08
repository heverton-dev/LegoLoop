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
)

func TestVideoReaCLI_HelpAndInvalid(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

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

func TestVideoReaCLI_PipelineReal(t *testing.T) {
	tempDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}

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
