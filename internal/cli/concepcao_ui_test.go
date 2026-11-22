package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestEntrevistaUISobeAutomaticamenteSemStudioExterno(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	var out, errOut bytes.Buffer
	deps := Deps{Stdout: &out, Stderr: &errOut, Home: t.TempDir()}
	code := cmdEntrevista(ctx, []string{"ui", "--docs", t.TempDir()}, deps)
	if code != 0 {
		t.Fatalf("cmdEntrevista ui code=%d stderr=%s", code, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, "http://127.0.0.1:") || !strings.Contains(got, "/concepcao") {
		t.Fatalf("esperava URL clicável /concepcao no stdout, obteve: %q", got)
	}
}
