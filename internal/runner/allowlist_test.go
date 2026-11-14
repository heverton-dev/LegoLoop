package runner

import (
	"testing"

	"legoloop/internal/dag"
)

func TestValidateToolAllowlist_BloqueiaNaoPermitida(t *testing.T) {
	n := dag.Node{
		ID: "n1",
		Tooling: map[string]any{
			"allowlist": []string{"read_file", "go_test"},
		},
	}
	if err := ValidateToolAllowlist(n, "read_file"); err != nil {
		t.Fatalf("ferramenta permitida falhou: %v", err)
	}
	if err := ValidateToolAllowlist(n, "rm_rf"); err == nil {
		t.Fatal("esperava erro para ferramenta fora da allowlist")
	}
}

func TestValidateToolAllowlist_PathJailEmContextIn(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	err := injectContextIn(src, dst, []string{"../fora.txt"})
	if err == nil {
		t.Fatal("esperava erro de Path Jail para ../fora.txt")
	}
}
