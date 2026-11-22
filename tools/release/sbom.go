package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// SBOM representa um documento CycloneDX enxuto gerado sem dependencias externas.
type SBOM struct {
	BOMFormat   string       `json:"bomFormat"`
	SpecVersion string       `json:"specVersion"`
	Version     int          `json:"version"`
	Metadata    SBOMMetadata `json:"metadata"`
	Components  []SBOMComp   `json:"components"`
}

// SBOMMetadata descreve o modulo principal e a versao do compilador Go.
type SBOMMetadata struct {
	GoVersion string   `json:"go_version"`
	Component SBOMComp `json:"component"`
}

// SBOMComp descreve uma Peça de software listada no SBOM.
type SBOMComp struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
	PURL    string `json:"purl,omitempty"`
}

// LerModuloEGo extrai o nome do modulo e a versao do Go de go.mod.
func LerModuloEGo(raiz string) (string, string, error) {
	dados, err := os.ReadFile(filepath.Join(raiz, "go.mod"))
	if err != nil {
		return "", "", err
	}
	var modulo, goVer string
	for _, linha := range strings.Split(string(dados), "\n") {
		campos := strings.Fields(strings.TrimSpace(linha))
		if len(campos) < 2 {
			continue
		}
		if campos[0] == "module" && modulo == "" {
			modulo = campos[1]
		}
		if campos[0] == "go" && goVer == "" {
			goVer = campos[1]
		}
	}
	if modulo == "" || goVer == "" {
		return "", "", errors.New("go.mod sem diretiva module ou go")
	}
	return modulo, goVer, nil
}

// GerarSBOM cria o arquivo sbom.cdx.json em outDir.
func GerarSBOM(raiz, versao, outDir string) (string, error) {
	modulo, goVer, err := LerModuloEGo(raiz)
	if err != nil {
		return "", err
	}
	doc := SBOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.5",
		Version:     1,
		Metadata: SBOMMetadata{
			GoVersion: goVer,
			Component: SBOMComp{
				Type:    "application",
				Name:    modulo,
				Version: versao,
				PURL:    "pkg:golang/" + modulo + "@" + versao,
			},
		},
		Components: []SBOMComp{
			{Type: "runtime", Name: "go", Version: goVer},
		},
	}
	dados, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	caminho := filepath.Join(outDir, "sbom.cdx.json")
	if err := os.WriteFile(caminho, append(dados, '\n'), 0644); err != nil {
		return "", err
	}
	return caminho, nil
}
