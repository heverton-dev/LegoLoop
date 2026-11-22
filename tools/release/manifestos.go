package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const urlBaseReleasePadrao = "https://github.com/heverton-dev/LegoLoop/releases/download"
const homepageLegoLoop = "https://github.com/heverton-dev/LegoLoop"

// ConfigManifesto configura a geracao dos manifestos de distribuicao.
type ConfigManifesto struct {
	Versao   string
	BaseURL  string
	SumsPath string
	OutDir   string
}

// ManifestosGerados guarda os caminhos dos arquivos gerados.
type ManifestosGerados struct {
	Winget string
	Scoop  string
	Brew   string
}

// GerarManifestos gera manifestos de winget, scoop e brew a partir do SHA256SUMS.
func GerarManifestos(cfg ConfigManifesto) (ManifestosGerados, error) {
	if strings.TrimSpace(cfg.Versao) == "" {
		return ManifestosGerados{}, errors.New("versao vazia")
	}
	if strings.TrimSpace(cfg.OutDir) == "" {
		return ManifestosGerados{}, errors.New("diretorio de saida vazio")
	}
	if strings.TrimSpace(cfg.BaseURL) == "" {
		cfg.BaseURL = urlBaseReleasePadrao
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	hashes, err := LerSHA256SUMS(cfg.SumsPath)
	if err != nil {
		return ManifestosGerados{}, err
	}
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return ManifestosGerados{}, err
	}
	return gravarTodosManifestos(cfg, hashes)
}

func gravarTodosManifestos(cfg ConfigManifesto, hashes map[string]string) (ManifestosGerados, error) {
	wingetPath, err := gerarManifestoWinget(cfg, hashes)
	if err != nil {
		return ManifestosGerados{}, err
	}
	scoopPath, err := gerarManifestoScoop(cfg, hashes)
	if err != nil {
		return ManifestosGerados{}, err
	}
	brewPath, err := gerarManifestoBrew(cfg, hashes)
	if err != nil {
		return ManifestosGerados{}, err
	}
	return ManifestosGerados{Winget: wingetPath, Scoop: scoopPath, Brew: brewPath}, nil
}

// LerSHA256SUMS le o arquivo SHA256SUMS e retorna mapa arquivo -> sha256.
func LerSHA256SUMS(sumsPath string) (map[string]string, error) {
	dados, err := os.ReadFile(sumsPath)
	if err != nil {
		return nil, err
	}
	hashes := make(map[string]string)
	for _, linha := range strings.Split(strings.TrimSpace(string(dados)), "\n") {
		if strings.TrimSpace(linha) == "" {
			continue
		}
		partes := strings.SplitN(linha, "  ", 2)
		if len(partes) != 2 || len(partes[0]) != 64 {
			return nil, fmt.Errorf("linha invalida em SHA256SUMS: %q", linha)
		}
		hashes[strings.TrimSpace(partes[1])] = partes[0]
	}
	if len(hashes) == 0 {
		return nil, errors.New("SHA256SUMS vazio")
	}
	return hashes, nil
}

func urlEHashPacote(cfg ConfigManifesto, hashes map[string]string, so, arch string) (string, string, error) {
	ext := ".tar.gz"
	if so == "windows" {
		ext = ".zip"
	}
	nome := fmt.Sprintf("legoloop_%s_%s_%s%s", cfg.Versao, so, arch, ext)
	hash, ok := hashes[nome]
	if !ok {
		return "", "", fmt.Errorf("checksum ausente para %s", nome)
	}
	tag := "v" + strings.TrimPrefix(cfg.Versao, "v")
	return fmt.Sprintf("%s/%s/%s", cfg.BaseURL, tag, nome), hash, nil
}

const templateWinget = `PackageIdentifier: heverton-dev.LegoLoop
PackageVersion: %s
PackageName: LegoLoop
Publisher: heverton-dev
License: MIT
ShortDescription: Plataforma de processos deterministas LegoLoop
PackageUrl: %s
ManifestType: singleton
ManifestVersion: 1.6.0
Installers:
  - Architecture: x64
    InstallerType: zip
    NestedInstallerType: portable
    NestedInstallerFiles:
      - RelativeFilePath: legoloop.exe
        PortableCommandAlias: legoloop
    Scope: user
    InstallerUrl: %s
    InstallerSha256: %s
  - Architecture: arm64
    InstallerType: zip
    NestedInstallerType: portable
    NestedInstallerFiles:
      - RelativeFilePath: legoloop.exe
        PortableCommandAlias: legoloop
    Scope: user
    InstallerUrl: %s
    InstallerSha256: %s
`

func gerarManifestoWinget(cfg ConfigManifesto, hashes map[string]string) (string, error) {
	urlAmd, hashAmd, err := urlEHashPacote(cfg, hashes, "windows", "amd64")
	if err != nil {
		return "", err
	}
	urlArm, hashArm, err := urlEHashPacote(cfg, hashes, "windows", "arm64")
	if err != nil {
		return "", err
	}
	conteudo := fmt.Sprintf(templateWinget, cfg.Versao, homepageLegoLoop, urlAmd, hashAmd, urlArm, hashArm)
	dest := filepath.Join(cfg.OutDir, "heverton-dev.LegoLoop.yaml")
	return dest, os.WriteFile(dest, []byte(conteudo), 0644)
}

func gerarManifestoScoop(cfg ConfigManifesto, hashes map[string]string) (string, error) {
	urlAmd, hashAmd, err := urlEHashPacote(cfg, hashes, "windows", "amd64")
	if err != nil {
		return "", err
	}
	urlArm, hashArm, err := urlEHashPacote(cfg, hashes, "windows", "arm64")
	if err != nil {
		return "", err
	}
	doc := map[string]any{
		"version":     cfg.Versao,
		"description": "Plataforma de processos deterministas LegoLoop",
		"homepage":    homepageLegoLoop,
		"license":     "MIT",
		"architecture": map[string]any{
			"64bit": map[string]string{"url": urlAmd, "hash": hashAmd, "bin": "legoloop.exe"},
			"arm64": map[string]string{"url": urlArm, "hash": hashArm, "bin": "legoloop.exe"},
		},
	}
	dados, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	dest := filepath.Join(cfg.OutDir, "legoloop.json")
	return dest, os.WriteFile(dest, append(dados, '\n'), 0644)
}

func gerarManifestoBrew(cfg ConfigManifesto, hashes map[string]string) (string, error) {
	urlMacArm, hashMacArm, err := urlEHashPacote(cfg, hashes, "darwin", "arm64")
	if err != nil {
		return "", err
	}
	urlMacAmd, hashMacAmd, err := urlEHashPacote(cfg, hashes, "darwin", "amd64")
	if err != nil {
		return "", err
	}
	urlLinArm, hashLinArm, err := urlEHashPacote(cfg, hashes, "linux", "arm64")
	if err != nil {
		return "", err
	}
	urlLinAmd, hashLinAmd, err := urlEHashPacote(cfg, hashes, "linux", "amd64")
	if err != nil {
		return "", err
	}
	conteudo := montarFormulaBrew(cfg.Versao, urlMacArm, hashMacArm, urlMacAmd, hashMacAmd, urlLinArm, hashLinArm, urlLinAmd, hashLinAmd)
	dest := filepath.Join(cfg.OutDir, "legoloop.rb")
	return dest, os.WriteFile(dest, []byte(conteudo), 0644)
}

const templateBrew = `class Legoloop < Formula
  desc "Plataforma de processos deterministas LegoLoop"
  homepage "%s"
  version "%s"
  license "MIT"

  on_macos do
    on_arm do
      url "%s"
      sha256 "%s"
    end
    on_intel do
      url "%s"
      sha256 "%s"
    end
  end

  on_linux do
    on_arm do
      url "%s"
      sha256 "%s"
    end
    on_intel do
      url "%s"
      sha256 "%s"
    end
  end

  def install
    bin.install "legoloop"
  end

  test do
    system "#{bin}/legoloop", "--version"
  end
end
`

func montarFormulaBrew(versao, uMacArm, hMacArm, uMacAmd, hMacAmd, uLinArm, hLinArm, uLinAmd, hLinAmd string) string {
	return fmt.Sprintf(templateBrew, homepageLegoLoop, versao, uMacArm, hMacArm, uMacAmd, hMacAmd, uLinArm, hLinArm, uLinAmd, hLinAmd)
}
