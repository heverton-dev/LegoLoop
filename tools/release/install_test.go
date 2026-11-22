package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const versaoTesteInstall = "v0.1.0-test"

func prepararServidorRelease(t *testing.T, corromperSha bool) (string, func()) {
	t.Helper()
	raiz := raizRepoTeste(t)
	serveDir := t.TempDir()
	binDir := t.TempDir()
	binPath := filepath.Join(binDir, NomeBinario(runtime.GOOS))
	if _, err := CompilarBinario(raiz, runtime.GOOS, runtime.GOARCH, versaoTesteInstall, binPath); err != nil {
		t.Fatalf("erro ao compilar binario de teste: %v", err)
	}
	pkgPath, err := Empacotar(runtime.GOOS, runtime.GOARCH, versaoTesteInstall, binPath, serveDir)
	if err != nil {
		t.Fatalf("erro ao empacotar binario de teste: %v", err)
	}
	pkgNome := filepath.Base(pkgPath)
	if corromperSha {
		linha := strings.Repeat("0", 64) + "  " + pkgNome + "\n"
		if err := os.WriteFile(filepath.Join(serveDir, "SHA256SUMS"), []byte(linha), 0644); err != nil {
			t.Fatalf("erro ao gravar SHA256SUMS adulterado: %v", err)
		}
	} else if _, err := GerarSHA256SUMS(serveDir, []string{pkgNome}); err != nil {
		t.Fatalf("erro ao gerar SHA256SUMS: %v", err)
	}
	srv := httptest.NewServer(http.FileServer(http.Dir(serveDir)))
	return srv.URL, srv.Close
}

func executarScriptPS1(raiz, baseURL, installDir, homeDir string, dryRun bool) ([]byte, error) {
	script := filepath.Join(raiz, "install.ps1")
	args := []string{
		"-NoProfile", "-ExecutionPolicy", "Bypass",
		"-File", script,
		"-BaseUrl", baseURL,
		"-InstallDir", installDir,
		"-Versao", versaoTesteInstall,
		"-HomeDir", homeDir,
	}
	if dryRun {
		args = append(args, "-DryRun")
	}
	return exec.Command("powershell", args...).CombinedOutput()
}

func executarScriptSH(shPath, raiz, baseURL, installDir, homeDir string, dryRun bool) ([]byte, error) {
	script := filepath.ToSlash(filepath.Join(raiz, "install.sh"))
	args := []string{
		script,
		"--base-url", baseURL,
		"--install-dir", filepath.ToSlash(installDir),
		"--versao", versaoTesteInstall,
		"--home", filepath.ToSlash(homeDir),
	}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return exec.Command(shPath, args...).CombinedOutput()
}

func verificarDiretorioSemBinario(t *testing.T, dir string) {
	t.Helper()
	bin := filepath.Join(dir, NomeBinario(runtime.GOOS))
	if _, err := os.Stat(bin); !os.IsNotExist(err) {
		t.Fatalf("binario nao deveria existir em %s", bin)
	}
}

func verificarRelatorioForja(t *testing.T, homeDir string) {
	t.Helper()
	padrao := filepath.Join(homeDir, ".legoloop", "forja", "*", "RELATORIO.json")
	encontrados, err := filepath.Glob(padrao)
	if err != nil {
		t.Fatalf("erro ao procurar relatorio: %v", err)
	}
	if len(encontrados) == 0 {
		t.Fatalf("nenhum RELATORIO.json encontrado em %s", padrao)
	}
}

func verificarFalhaSha(t *testing.T, out []byte, err error, instDir string) {
	t.Helper()
	if err == nil {
		t.Fatalf("esperado erro com SHA-256 errado; saida:\n%s", string(out))
	}
	if !strings.Contains(strings.ToLower(string(out)), "sha-256") {
		t.Fatalf("erro inesperado (sem mencao a SHA-256): %v\nsaida:\n%s", err, string(out))
	}
	verificarDiretorioSemBinario(t, instDir)
}

func TestScriptAbortaComShaErrado(t *testing.T) {
	raiz := raizRepoTeste(t)
	baseURL, fechar := prepararServidorRelease(t, true)
	defer fechar()

	if runtime.GOOS == "windows" {
		t.Run("install.ps1", func(t *testing.T) {
			instDir := filepath.Join(t.TempDir(), "bin")
			out, err := executarScriptPS1(raiz, baseURL, instDir, t.TempDir(), false)
			verificarFalhaSha(t, out, err, instDir)
		})
	}
	if shPath, err := exec.LookPath("sh"); err == nil {
		t.Run("install.sh", func(t *testing.T) {
			instDir := filepath.Join(t.TempDir(), "bin")
			out, err := executarScriptSH(shPath, raiz, baseURL, instDir, t.TempDir(), false)
			verificarFalhaSha(t, out, err, instDir)
		})
	}
}

func TestScriptInstalaEDisparaForja(t *testing.T) {
	raiz := raizRepoTeste(t)
	baseURL, fechar := prepararServidorRelease(t, false)
	defer fechar()

	if runtime.GOOS == "windows" {
		t.Run("install.ps1", func(t *testing.T) {
			validarProcessoCompletoScript(t, func(instDir, homeDir string, dry bool) ([]byte, error) {
				return executarScriptPS1(raiz, baseURL, instDir, homeDir, dry)
			})
		})
	}
	if shPath, err := exec.LookPath("sh"); err == nil {
		t.Run("install.sh", func(t *testing.T) {
			validarProcessoCompletoScript(t, func(instDir, homeDir string, dry bool) ([]byte, error) {
				return executarScriptSH(shPath, raiz, baseURL, instDir, homeDir, dry)
			})
		})
	}
}

func validarProcessoCompletoScript(t *testing.T, rodar func(instDir, homeDir string, dry bool) ([]byte, error)) {
	t.Helper()
	instDry := filepath.Join(t.TempDir(), "bin-dry")
	homeDry := t.TempDir()
	if out, err := rodar(instDry, homeDry, true); err != nil {
		t.Fatalf("dry-run falhou: %v\nsaida:\n%s", err, string(out))
	}
	verificarDiretorioSemBinario(t, instDry)
	if entradas, err := os.ReadDir(homeDry); err != nil || len(entradas) != 0 {
		t.Fatalf("dry-run nao deveria gravar em HOME: entradas=%d err=%v", len(entradas), err)
	}

	instReal := filepath.Join(t.TempDir(), "bin-real")
	homeReal := t.TempDir()
	if out, err := rodar(instReal, homeReal, false); err != nil {
		t.Fatalf("instalacao real falhou: %v\nsaida:\n%s", err, string(out))
	}
	bin := filepath.Join(instReal, NomeBinario(runtime.GOOS))
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("binario nao instalado em %s: %v", bin, err)
	}
	verificarRelatorioForja(t, homeReal)
}
