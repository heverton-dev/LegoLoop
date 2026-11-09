package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCartoonShortsCLI_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdout: &stdout, Stderr: &stderr}

	exit := Run(context.Background(), []string{"cartoon-shorts", "--help"}, deps)
	if exit != 0 {
		t.Fatalf("esperava exit code 0 para --help, obteve: %d", exit)
	}
	if !strings.Contains(stdout.String(), "cartoon-shorts") {
		t.Errorf("saída de help não contém 'cartoon-shorts': %s", stdout.String())
	}
}

func TestCartoonShortsCLI_SemArgumentos(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdout: &stdout, Stderr: &stderr}

	exit := Run(context.Background(), []string{"cartoon-shorts"}, deps)
	if exit != 2 {
		t.Fatalf("esperava exit code 2 (erro de uso) sem argumentos, obteve: %d", exit)
	}
}

func TestCartoonShortsCLI_TemaValidoComDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	deps := Deps{Stdout: &stdout, Stderr: &stderr}

	exit := Run(context.Background(), []string{"cartoon-shorts", "--theme", "ET de Varginha", "--dry-run"}, deps)
	if exit != 0 {
		t.Fatalf("esperava exit code 0 em --dry-run, obteve: %d (stderr: %s)", exit, stderr.String())
	}
	if !strings.Contains(stdout.String(), "ET de Varginha") {
		t.Errorf("saída esperada não contém o tema informado: %s", stdout.String())
	}
}
