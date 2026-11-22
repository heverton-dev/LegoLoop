package forja

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPathJail(t *testing.T) {
	home := t.TempDir()
	hostis := []string{
		"",
		"   ",
		".",
		"..",
		"../fora",
		"claude/../../fora",
		`..\fora`,
		"/etc/passwd",
		`C:\Windows\System32`,
		`\\servidor\compartilhamento`,
		filepath.Join(home, "claude"),
		"claude\x00invalido",
	}
	for _, entrada := range hostis {
		got, err := ResolverPastaHarness(home, entrada)
		if err == nil {
			t.Errorf("entrada hostil %q deveria falhar, mas retornou %q", entrada, got)
		}
	}

	if _, err := ResolverPastaHarness("", "claude"); err == nil {
		t.Error("home vazio deveria falhar")
	}
}

func TestPathJailValido(t *testing.T) {
	home := t.TempDir()
	validos := []string{"claude", ".claude", ".config/opencode"}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		t.Fatalf("filepath.Abs falhou: %v", err)
	}
	for _, nome := range validos {
		got, err := ResolverPastaHarness(home, nome)
		if err != nil {
			t.Fatalf("nome válido %q falhou: %v", nome, err)
		}
		rel, err := filepath.Rel(homeAbs, got)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			t.Errorf("caminho %q saiu do home %q (rel=%q)", got, homeAbs, rel)
		}
	}
}
