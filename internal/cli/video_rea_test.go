package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"legoloop/internal/cli"
)

func TestVideoReaCLI(t *testing.T) {
	cases := []struct {
		args     []string
		wantExit int
		contains string
	}{
		{[]string{"video-rea", "--help"}, 0, "Uso: legoloop video-rea"},
		{[]string{"video-rea", "ingest"}, 0, "[VideoREA] subcomando \"ingest\" executado"},
		{[]string{"video-rea", "forensics"}, 0, "[VideoREA] subcomando \"forensics\" executado"},
		{[]string{"video-rea", "curate"}, 0, "[VideoREA] subcomando \"curate\" executado"},
		{[]string{"video-rea", "gate"}, 0, "[VideoREA] subcomando \"gate\" executado"},
		{[]string{"video-rea", "govern"}, 0, "[VideoREA] subcomando \"govern\" executado"},
		{[]string{"video-rea", "book"}, 0, "[VideoREA] subcomando \"book\" executado"},
		{[]string{"video-rea", "run"}, 0, "[VideoREA] subcomando \"run\" executado"},
		{[]string{"video-rea", "invalido"}, 2, "subcomando desconhecido \"invalido\""},
	}

	for _, c := range cases {
		var stdout, stderr bytes.Buffer
		deps := cli.Deps{
			Stdout: &stdout,
			Stderr: &stderr,
		}
		exit := cli.Run(context.Background(), c.args, deps)
		if exit != c.wantExit {
			t.Errorf("Run(%v) exit = %d; esperava %d", c.args, exit, c.wantExit)
		}
		out := stdout.String() + stderr.String()
		if !strings.Contains(out, c.contains) {
			t.Errorf("Run(%v) saída = %q; esperava conter %q", c.args, out, c.contains)
		}
	}
}
