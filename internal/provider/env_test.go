package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	tmp := t.TempDir()
	envPath := filepath.Join(tmp, ".env")
	content := `# Comentário
NINEROUTER_URL="https://9router.vpsconexao.org"
NINEROUTER_KEY=sk-test-12345
IGNORAR_VAZIA=
`
	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	vals := parseEnvFile(envPath)
	if vals["NINEROUTER_URL"] != "https://9router.vpsconexao.org" {
		t.Errorf("url incorreta: %s", vals["NINEROUTER_URL"])
	}
	if vals["NINEROUTER_KEY"] != "sk-test-12345" {
		t.Errorf("key incorreta: %s", vals["NINEROUTER_KEY"])
	}
}

func TestLoadRouterEnvFallback(t *testing.T) {
	tmp := t.TempDir()
	d9 := filepath.Join(tmp, "Desktop", "9router")
	if err := os.MkdirAll(d9, 0755); err != nil {
		t.Fatal(err)
	}
	envPath := filepath.Join(d9, ".env")
	if err := os.WriteFile(envPath, []byte("NINEROUTER_URL=https://9router.vpsconexao.org\nNINEROUTER_KEY=chave-vps\n"), 0600); err != nil {
		t.Fatal(err)
	}

	u, k := findRouterCreds(tmp, func(k string) string { return "" })
	if u != "https://9router.vpsconexao.org" {
		t.Errorf("esperava url vps, obteve %q", u)
	}
	if k != "chave-vps" {
		t.Errorf("esperava chave-vps, obteve %q", k)
	}
}
