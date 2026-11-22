package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type manifestoScoopJSON struct {
	Version      string `json:"version"`
	Homepage     string `json:"homepage"`
	Architecture map[string]struct {
		URL  string `json:"url"`
		Hash string `json:"hash"`
		Bin  string `json:"bin"`
	} `json:"architecture"`
}

func TestManifestos(t *testing.T) {
	dir := t.TempDir()
	versao := "1.2.3"
	baseURL := "https://github.com/heverton-dev/LegoLoop/releases/download"
	sumsPath, esperados := criarPacotesESumsTeste(t, dir, versao)

	outDir := filepath.Join(dir, "manifestos")
	gerados, err := GerarManifestos(ConfigManifesto{
		Versao:   versao,
		BaseURL:  baseURL,
		SumsPath: sumsPath,
		OutDir:   outDir,
	})
	if err != nil {
		t.Fatalf("GerarManifestos falhou: %v", err)
	}

	validarWinget(t, gerados.Winget, versao, baseURL, esperados)
	validarScoop(t, gerados.Scoop, versao, baseURL, esperados)
	validarBrew(t, gerados.Brew, versao, baseURL, esperados)
}

func TestManifestosErros(t *testing.T) {
	dir := t.TempDir()
	sumsPath := filepath.Join(dir, "SHA256SUMS")
	if err := os.WriteFile(sumsPath, []byte("invalido\n"), 0644); err != nil {
		t.Fatalf("gravar sums: %v", err)
	}
	if _, err := GerarManifestos(ConfigManifesto{Versao: "", SumsPath: sumsPath, OutDir: dir}); err == nil {
		t.Fatal("esperava erro para versao vazia")
	}
	if _, err := GerarManifestos(ConfigManifesto{Versao: "1.0.0", SumsPath: sumsPath, OutDir: dir}); err == nil {
		t.Fatal("esperava erro para SHA256SUMS invalido")
	}
}

func criarPacotesESumsTeste(t *testing.T, dir, versao string) (string, map[string]string) {
	t.Helper()
	var arquivos []string
	esperados := make(map[string]string, len(AlvosPadrao))
	for _, alvo := range AlvosPadrao {
		bin := filepath.Join(dir, "bin_"+alvo.SO+"_"+alvo.Arch)
		conteudo := []byte("binario-" + alvo.SO + "-" + alvo.Arch + "-" + versao)
		if err := os.WriteFile(bin, conteudo, 0755); err != nil {
			t.Fatalf("gravar binario falso: %v", err)
		}
		pkg, err := Empacotar(alvo.SO, alvo.Arch, versao, bin, dir)
		if err != nil {
			t.Fatalf("Empacotar %s/%s: %v", alvo.SO, alvo.Arch, err)
		}
		hash, err := HashArquivo(pkg)
		if err != nil {
			t.Fatalf("HashArquivo %s: %v", pkg, err)
		}
		nome := filepath.Base(pkg)
		arquivos = append(arquivos, nome)
		esperados[nome] = hash
	}
	sumsPath, err := GerarSHA256SUMS(dir, arquivos)
	if err != nil {
		t.Fatalf("GerarSHA256SUMS: %v", err)
	}
	return sumsPath, esperados
}

func validarWinget(t *testing.T, caminho, versao, baseURL string, hashes map[string]string) {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler winget: %v", err)
	}
	texto := string(dados)
	if !strings.Contains(texto, "PackageVersion: "+versao) {
		t.Errorf("winget sem PackageVersion %s:\n%s", versao, texto)
	}
	for _, arch := range []string{"amd64", "arm64"} {
		nome := "legoloop_" + versao + "_windows_" + arch + ".zip"
		url := baseURL + "/v" + versao + "/" + nome
		if !strings.Contains(texto, "InstallerUrl: "+url) {
			t.Errorf("winget sem URL %s", url)
		}
		if !strings.Contains(texto, "InstallerSha256: "+hashes[nome]) {
			t.Errorf("winget sem SHA-256 exato de %s", nome)
		}
	}
}

func validarScoop(t *testing.T, caminho, versao, baseURL string, hashes map[string]string) {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler scoop: %v", err)
	}
	var m manifestoScoopJSON
	if err := json.Unmarshal(dados, &m); err != nil {
		t.Fatalf("scoop JSON invalido: %v", err)
	}
	if m.Version != versao {
		t.Errorf("scoop version = %q, esperado %q", m.Version, versao)
	}
	pares := map[string]string{"64bit": "amd64", "arm64": "arm64"}
	for chave, arch := range pares {
		nome := "legoloop_" + versao + "_windows_" + arch + ".zip"
		url := baseURL + "/v" + versao + "/" + nome
		item := m.Architecture[chave]
		if item.URL != url || item.Hash != hashes[nome] || item.Bin != "legoloop.exe" {
			t.Errorf("scoop %s divergente: %+v (esperado url=%s hash=%s)", chave, item, url, hashes[nome])
		}
	}
}

func validarBrew(t *testing.T, caminho, versao, baseURL string, hashes map[string]string) {
	t.Helper()
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("ler brew: %v", err)
	}
	texto := string(dados)
	if !strings.Contains(texto, "class Legoloop < Formula") || !strings.Contains(texto, `version "`+versao+`"`) {
		t.Errorf("brew sem cabecalho ou versao %s", versao)
	}
	alvos := []Alvo{
		{SO: "darwin", Arch: "arm64"},
		{SO: "darwin", Arch: "amd64"},
		{SO: "linux", Arch: "arm64"},
		{SO: "linux", Arch: "amd64"},
	}
	for _, alvo := range alvos {
		nome := "legoloop_" + versao + "_" + alvo.SO + "_" + alvo.Arch + ".tar.gz"
		url := baseURL + "/v" + versao + "/" + nome
		if !strings.Contains(texto, `url "`+url+`"`) {
			t.Errorf("brew sem URL %s", url)
		}
		if !strings.Contains(texto, `sha256 "`+hashes[nome]+`"`) {
			t.Errorf("brew sem SHA-256 exato de %s", nome)
		}
	}
}
