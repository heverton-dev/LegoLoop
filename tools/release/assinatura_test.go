package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssinaturaDeclarada(t *testing.T) {
	pacotes := []ArtefatoAssinatura{
		{Arquivo: "legoloop_v1.0.0_linux_amd64.tar.gz", Tipo: "pacote", SHA256: "abc"},
		{Arquivo: "legoloop_v1.0.0_windows_amd64.zip", Tipo: "pacote", SHA256: "def"},
	}

	semCI := GerarRelatorioAssinatura(
		Config{Versao: "v1.0.0", Cosign: false, AssinarMSI: false},
		pacotes, "SHA256SUMS", "sbom.cdx.json", "github.com/heverton-dev/LegoLoop", "1.27.0",
	)
	checarStatusArtefato(t, semCI, "msi", StatusNaoAssinado)
	checarStatusArtefato(t, semCI, "pacote", StatusNaoAssinado)

	comCI := GerarRelatorioAssinatura(
		Config{Versao: "v1.0.0", Cosign: true, AssinarMSI: false},
		pacotes, "SHA256SUMS", "sbom.cdx.json", "github.com/heverton-dev/LegoLoop", "1.27.0",
	)
	checarStatusArtefato(t, comCI, "msi", StatusNaoAssinado)
	checarStatusArtefato(t, comCI, "pacote", StatusCosign)
	checarStatusArtefato(t, comCI, "checksums", StatusCosign)

	outDir := t.TempDir()
	caminho, err := GravarRelatorio(outDir, comCI)
	if err != nil {
		t.Fatalf("GravarRelatorio falhou: %v", err)
	}
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("erro lendo relatorio: %v", err)
	}
	var lido RelatorioRelease
	if err := json.Unmarshal(dados, &lido); err != nil {
		t.Fatalf("relatorio JSON invalido: %v", err)
	}
	if len(lido.Artefatos) != len(comCI.Artefatos) {
		t.Fatalf("quantidade de artefatos divergiu: %d != %d", len(lido.Artefatos), len(comCI.Artefatos))
	}
}

func checarStatusArtefato(t *testing.T, rel RelatorioRelease, tipo, esperado string) {
	t.Helper()
	achou := false
	for _, a := range rel.Artefatos {
		if a.Tipo == tipo {
			achou = true
			if a.Assinatura != esperado {
				t.Fatalf("tipo %s (%s): esperado %q, obtido %q", tipo, a.Arquivo, esperado, a.Assinatura)
			}
		}
	}
	if !achou {
		t.Fatalf("nenhum artefato do tipo %q encontrado no relatorio", tipo)
	}
}

func TestWorkflowReleaseConfigurado(t *testing.T) {
	raiz := raizRepoTeste(t)
	caminho := filepath.Join(raiz, ".github", "workflows", "release.yml")
	dados, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatalf("workflow release.yml ausente: %v", err)
	}
	texto := string(dados)
	for _, trecho := range []string{"v*", "go run ./tools/release --all", "cosign", "SHA256SUMS"} {
		if !strings.Contains(texto, trecho) {
			t.Fatalf("workflow release.yml nao contem %q", trecho)
		}
	}
}
