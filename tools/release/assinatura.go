package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// StatusNaoAssinado indica artefato sem assinatura criptografica.
	StatusNaoAssinado = "nao_assinado"
	// StatusCosign indica assinatura via cosign OIDC no CI.
	StatusCosign = "cosign"
	// StatusAuthenticode indica assinatura Authenticode do instalador MSI.
	StatusAuthenticode = "authenticode"
)

// ArtefatoAssinatura declara o estado de assinatura de um artefato do release.
type ArtefatoAssinatura struct {
	Arquivo    string `json:"arquivo"`
	Tipo       string `json:"tipo"`
	Assinatura string `json:"assinatura"`
	SHA256     string `json:"sha256,omitempty"`
}

// RelatorioRelease consolida os artefatos, checksums, SBOM e status de assinatura.
type RelatorioRelease struct {
	Versao     string               `json:"versao"`
	Modulo     string               `json:"modulo"`
	GoVersao   string               `json:"go_versao"`
	SHA256SUMS string               `json:"sha256sums"`
	SBOM       string               `json:"sbom"`
	Artefatos  []ArtefatoAssinatura `json:"artefatos"`
}

// GerarRelatorioAssinatura monta o relatorio declarando o status de cada artefato.
func GerarRelatorioAssinatura(
	cfg Config,
	pacotes []ArtefatoAssinatura,
	sumsFile, sbomFile, modulo, goVer string,
) RelatorioRelease {
	statusPacote := StatusNaoAssinado
	if cfg.Cosign {
		statusPacote = StatusCosign
	}
	statusMSI := StatusNaoAssinado
	if cfg.AssinarMSI {
		statusMSI = StatusAuthenticode
	}
	lista := make([]ArtefatoAssinatura, 0, len(pacotes)+2)
	for _, p := range pacotes {
		p.Assinatura = statusPacote
		lista = append(lista, p)
	}
	lista = append(lista,
		ArtefatoAssinatura{
			Arquivo:    fmt.Sprintf("legoloop_%s_windows_amd64.msi", cfg.Versao),
			Tipo:       "msi",
			Assinatura: statusMSI,
		},
		ArtefatoAssinatura{
			Arquivo:    sumsFile,
			Tipo:       "checksums",
			Assinatura: statusPacote,
		},
	)
	return RelatorioRelease{
		Versao:     cfg.Versao,
		Modulo:     modulo,
		GoVersao:   goVer,
		SHA256SUMS: sumsFile,
		SBOM:       sbomFile,
		Artefatos:  lista,
	}
}

// GravarRelatorio persiste o relatorio de release em release-relatorio.json.
func GravarRelatorio(outDir string, rel RelatorioRelease) (string, error) {
	dados, err := json.MarshalIndent(rel, "", "  ")
	if err != nil {
		return "", err
	}
	caminho := filepath.Join(outDir, "release-relatorio.json")
	if err := os.WriteFile(caminho, append(dados, '\n'), 0644); err != nil {
		return "", err
	}
	return caminho, nil
}
