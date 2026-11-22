package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPacotesSumsESBOM(t *testing.T) {
	raiz := raizRepoTeste(t)
	outDir := t.TempDir()
	cfg := Config{
		Raiz:   raiz,
		OutDir: outDir,
		Versao: "v1.0.0-test",
		Todos:  true,
	}

	rel, err := Montar(cfg)
	if err != nil {
		t.Fatalf("Montar falhou: %v", err)
	}
	verificarSeisPacotes(t, outDir, cfg.Versao)
	verificarSumsValidoEAdulterado(t, outDir, rel.SHA256SUMS)
	verificarSBOMConteudo(t, outDir, rel.SBOM)
}

func verificarSeisPacotes(t *testing.T, outDir, versao string) {
	t.Helper()
	esperados := []string{
		"legoloop_" + versao + "_windows_amd64.zip",
		"legoloop_" + versao + "_windows_arm64.zip",
		"legoloop_" + versao + "_linux_amd64.tar.gz",
		"legoloop_" + versao + "_linux_arm64.tar.gz",
		"legoloop_" + versao + "_darwin_amd64.tar.gz",
		"legoloop_" + versao + "_darwin_arm64.tar.gz",
	}
	for _, nome := range esperados {
		info, err := os.Stat(filepath.Join(outDir, nome))
		if err != nil {
			t.Fatalf("pacote ausente %s: %v", nome, err)
		}
		if info.Size() == 0 {
			t.Fatalf("pacote vazio: %s", nome)
		}
	}
}

func verificarSumsValidoEAdulterado(t *testing.T, outDir, sumsNome string) {
	t.Helper()
	sumsPath := filepath.Join(outDir, sumsNome)
	if err := VerificarSHA256SUMS(outDir, sumsPath); err != nil {
		t.Fatalf("SHA256SUMS invalido: %v", err)
	}
	dados, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatalf("erro lendo SHA256SUMS: %v", err)
	}
	adulterado := filepath.Join(t.TempDir(), "SHA256SUMS.bad")
	linhaRuim := strings.Replace(string(dados), string(dados[:4]), "0000", 1)
	if err := os.WriteFile(adulterado, []byte(linhaRuim), 0644); err != nil {
		t.Fatalf("erro gravando sums adulterado: %v", err)
	}
	if err := VerificarSHA256SUMS(outDir, adulterado); err == nil {
		t.Fatal("esperado erro ao verificar SHA256SUMS adulterado")
	}
}

func verificarSBOMConteudo(t *testing.T, outDir, sbomNome string) {
	t.Helper()
	dados, err := os.ReadFile(filepath.Join(outDir, sbomNome))
	if err != nil {
		t.Fatalf("erro lendo SBOM: %v", err)
	}
	var doc SBOM
	if err := json.Unmarshal(dados, &doc); err != nil {
		t.Fatalf("SBOM JSON invalido: %v", err)
	}
	if doc.Metadata.Component.Name != "github.com/heverton-dev/LegoLoop" {
		t.Fatalf("modulo incorreto no SBOM: %q", doc.Metadata.Component.Name)
	}
	if doc.Metadata.GoVersion == "" {
		t.Fatal("versao do Go ausente no SBOM")
	}
}

func TestEmpacotarReproduzivel(t *testing.T) {
	tmp := t.TempDir()
	binFake := filepath.Join(tmp, "legoloop")
	if err := os.WriteFile(binFake, []byte("binario-teste-deterministico"), 0755); err != nil {
		t.Fatal(err)
	}
	out1, out2 := t.TempDir(), t.TempDir()
	for _, so := range []string{"windows", "linux"} {
		p1, err := Empacotar(so, "amd64", "v1.0.0", binFake, out1)
		if err != nil {
			t.Fatalf("Empacotar %s 1 falhou: %v", so, err)
		}
		p2, err := Empacotar(so, "amd64", "v1.0.0", binFake, out2)
		if err != nil {
			t.Fatalf("Empacotar %s 2 falhou: %v", so, err)
		}
		h1, _ := HashArquivo(p1)
		h2, _ := HashArquivo(p2)
		if h1 != h2 {
			t.Fatalf("pacote %s nao reproduzivel: %s != %s", so, h1, h2)
		}
	}
}
