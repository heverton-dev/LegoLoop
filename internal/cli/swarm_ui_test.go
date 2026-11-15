package cli

import (
	"context"
	"strings"
	"testing"
)

func TestCmdSwarmUI_Help(t *testing.T) {
	d, out, errBuf := deps(t)
	code := Run(context.Background(), []string{"swarm", "ui", "--help"}, d)
	if code != 0 {
		t.Fatalf("swarm ui --help retornou %d, want 0; stderr=%q", code, errBuf.String())
	}
	combined := out.String() + errBuf.String()
	for _, flagName := range []string{"-addr", "-dir", "-root"} {
		if !strings.Contains(combined, flagName) {
			t.Errorf("swarm ui --help deveria documentar %q:\n%s", flagName, combined)
		}
	}
}
