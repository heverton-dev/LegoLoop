package instalar

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func criarBinarioFalso(t *testing.T, dir, nome string) string {
	t.Helper()
	caminho := filepath.Join(dir, nome)
	if err := os.WriteFile(caminho, []byte("bin"), 0o755); err != nil {
		t.Fatalf("falha ao criar binário falso: %v", err)
	}
	return caminho
}

func TestAtalhoDoGrafo(t *testing.T) {
	t.Run("windows cria cmd", func(t *testing.T) {
		dirBin := filepath.Join(t.TempDir(), "bin")
		fonte := criarBinarioFalso(t, t.TempDir(), "codebase-memory-mcp.exe")
		destino, err := CriarAtalhoGrafo(dirBin, fonte, "windows", true)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !strings.HasSuffix(destino, "codebase-memory-mcp.cmd") {
			t.Fatalf("destino=%q sem sufixo .cmd", destino)
		}
		dados, err := os.ReadFile(destino)
		if err != nil {
			t.Fatalf("atalho .cmd não gravado: %v", err)
		}
		if !strings.Contains(string(dados), fonte) || !strings.Contains(string(dados), "%*") {
			t.Fatalf("conteúdo inesperado do .cmd: %q", string(dados))
		}
	})

	t.Run("posix cria sh executavel", func(t *testing.T) {
		dirBin := filepath.Join(t.TempDir(), "bin")
		fonte := criarBinarioFalso(t, t.TempDir(), "codebase-memory-mcp")
		destino, err := CriarAtalhoGrafo(dirBin, fonte, "linux", true)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if filepath.Base(destino) != "codebase-memory-mcp" {
			t.Fatalf("nome inesperado: %q", destino)
		}
		info, err := os.Stat(destino)
		if err != nil {
			t.Fatalf("atalho posix não gravado: %v", err)
		}
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0o755 {
			t.Fatalf("permissão=%o, esperava 0755", info.Mode().Perm())
		}
		dados, err := os.ReadFile(destino)
		if err != nil {
			t.Fatalf("falha ao ler atalho posix: %v", err)
		}
		if !strings.HasPrefix(string(dados), "#!/usr/bin/env sh\n") || !strings.Contains(string(dados), fonte) {
			t.Fatalf("conteúdo inesperado do script posix: %q", string(dados))
		}
	})
}

func TestAtalhoDoGrafoDryRunEAusente(t *testing.T) {
	dirBin := filepath.Join(t.TempDir(), "bin")
	fonte := criarBinarioFalso(t, t.TempDir(), "codebase-memory-mcp")
	destino, err := CriarAtalhoGrafo(dirBin, fonte, "darwin", false)
	if err != nil {
		t.Fatalf("dry-run falhou: %v", err)
	}
	if destino == "" {
		t.Fatal("esperava caminho planejado no dry-run")
	}
	if _, statErr := os.Stat(dirBin); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run não deveria criar %s", dirBin)
	}
	vazio, err := CriarAtalhoGrafo(dirBin, filepath.Join(t.TempDir(), "inexistente"), "linux", true)
	if err != nil {
		t.Fatalf("binário ausente não deve falhar: %v", err)
	}
	if vazio != "" {
		t.Fatalf("esperava destino vazio para binário ausente, obteve %q", vazio)
	}
}
