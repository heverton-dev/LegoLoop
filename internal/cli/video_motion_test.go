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
