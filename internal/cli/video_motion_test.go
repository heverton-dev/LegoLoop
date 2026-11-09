package cli_test

import (
	"bytes"
	"context"
	"testing"

	"legoloop/internal/cli"
)

func TestVideoMotionCLI_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	code := cli.Run(context.Background(), []string{"video-motion", "--help"}, deps)
	if code != 0 {
		t.Fatalf("esperava exit code 0 para --help, obteve %d", code)
	}

	out := stdout.String()
	if !bytes.Contains(stdout.Bytes(), []byte("legoloop video-motion")) {
		t.Errorf("ajuda de video-motion ausente no stdout:\n%s", out)
	}
}

func TestVideoMotionCLI_UnknownSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	code := cli.Run(context.Background(), []string{"video-motion", "inexistente"}, deps)
	if code != 2 {
		t.Fatalf("esperava exit code 2 para subcomando desconhecido, obteve %d", code)
	}
}

func TestVideoMotionCLI_RejeitaStubsSemFlagsObrigatorias(t *testing.T) {
	subcmds := []string{"silence", "punch", "assets", "preview", "bundle", "run"}
	for _, sc := range subcmds {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
		code := cli.Run(context.Background(), []string{"video-motion", sc}, deps)
		if code == 0 {
			t.Errorf("subcomando %q retornou exit code 0 sem parâmetros obrigatórios (stub detectado)", sc)
		}
	}
}

func TestVideoMotionCLI_SuporteAModelEEffort(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := cli.Deps{Stdout: &stdout, Stderr: &stderr}
	code := cli.Run(context.Background(), []string{"video-motion", "run", "video.mp4", "--model", "modelo-inexistente-xyz"}, deps)
	if code == 0 {
		t.Fatalf("esperava falha ao requisitar modelo inexistente com --model")
	}
	if !bytes.Contains(stderr.Bytes(), []byte("cérebro")) && !bytes.Contains(stderr.Bytes(), []byte("modelo")) {
		t.Errorf("esperava erro de modelo no stderr, obteve:\n%s", stderr.String())
	}
}
