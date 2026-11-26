package kraften

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// criarZip grava um .zip temporário com as entradas informadas.
func criarZip(t *testing.T, entradas map[string]string) string {
	t.Helper()
	caminho := filepath.Join(t.TempDir(), "pacote.zip")
	arq, err := os.Create(caminho)
	if err != nil {
		t.Fatalf("criar zip: %v", err)
	}
	defer arq.Close()
	w := zip.NewWriter(arq)
	for nome, corpo := range entradas {
		f, err := w.Create(nome)
		if err != nil {
			t.Fatalf("entrada %s: %v", nome, err)
		}
		if _, err := f.Write([]byte(corpo)); err != nil {
			t.Fatalf("gravar %s: %v", nome, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("fechar zip: %v", err)
	}
	return caminho
}

const (
	dockerfileSaudavel = "FROM alpine:3.20\nHEALTHCHECK CMD wget -q -O- http://localhost:8080/health || exit 1\nCMD [\"./app\"]\n"
	stackSaudavel      = "version: \"3.8\"\nservices:\n  web:\n    image: app:latest\n    networks:\n      - network_conexao\n"
	envSaudavel        = "ADMIN_USER=admin\nADMIN_PASS=segredo\nAPP_PORT=8080\n"
)

func entradasSaudaveis() map[string]string {
	return map[string]string{
		"app/Dockerfile":       dockerfileSaudavel,
		"app/docker-stack.yml": stackSaudavel,
		"app/.env.example":     envSaudavel,
	}
}

func TestPreflightZipExtraiSemEscapar(t *testing.T) {
	zipPath := criarZip(t, entradasSaudaveis())
	dest := t.TempDir()
	if err := Preflight(zipPath, dest); err != nil {
		t.Fatalf("pacote válido rejeitado: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "app", "Dockerfile")); err != nil {
		t.Fatalf("Dockerfile não extraído dentro da raiz: %v", err)
	}
}

func TestPreflightZipRejeitaZipSlip(t *testing.T) {
	zipPath := criarZip(t, map[string]string{
		"app/Dockerfile":       dockerfileSaudavel,
		"app/docker-stack.yml": stackSaudavel,
		"app/.env.example":     envSaudavel,
		"../evil.txt":          "fora da raiz",
	})
	dest := t.TempDir()
	if err := Preflight(zipPath, dest); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("esperava ErrInvalidDeployInput para zip slip, obteve %v", err)
	}
	raiz := filepath.Dir(dest)
	if _, err := os.Stat(filepath.Join(raiz, "evil.txt")); err == nil {
		t.Fatal("arquivo escapou do diretório raiz")
	}
}

func TestPreflightZipRejeitaSemDockerfile(t *testing.T) {
	entradas := entradasSaudaveis()
	delete(entradas, "app/Dockerfile")
	if err := Preflight(criarZip(t, entradas), t.TempDir()); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("esperava ErrInvalidDeployInput sem Dockerfile, obteve %v", err)
	}
}

func TestPreflightZipRejeitaSemHealthcheck(t *testing.T) {
	entradas := entradasSaudaveis()
	entradas["app/Dockerfile"] = "FROM alpine:3.20\nCMD [\"./app\"]\n"
	if err := Preflight(criarZip(t, entradas), t.TempDir()); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("esperava ErrInvalidDeployInput sem HEALTHCHECK, obteve %v", err)
	}
}

func TestPreflightZipRejeitaStackSemRede(t *testing.T) {
	entradas := entradasSaudaveis()
	entradas["app/docker-stack.yml"] = "services:\n  web:\n    networks:\n      - rede_errada\n"
	if err := Preflight(criarZip(t, entradas), t.TempDir()); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("esperava ErrInvalidDeployInput sem network_conexao, obteve %v", err)
	}
}

func TestPreflightZipRejeitaSemVariaveisPadrao(t *testing.T) {
	for nome, corpo := range map[string]string{
		"app/.env.example": "APP_PORT=8080\n",
		"app/.env":         "ADMIN_USER=admin\nADMIN_PASS=segredo\n",
	} {
		entradas := entradasSaudaveis()
		delete(entradas, "app/.env.example")
		if nome == "app/.env" {
			entradas["app/.env"] = corpo
		} else {
			entradas[nome] = corpo
		}
		err := Preflight(criarZip(t, entradas), t.TempDir())
		if !errors.Is(err, ErrInvalidDeployInput) {
			t.Fatalf("esperava ErrInvalidDeployInput para %s, obteve %v", nome, err)
		}
	}
}

func TestPreflightZipEncontraSubpastaDoDockerfile(t *testing.T) {
	entradas := entradasSaudaveis()
	renomeado := map[string]string{}
	for nome, corpo := range entradas {
		renomeado["release-1.2.3/"+nome] = corpo
	}
	dest := t.TempDir()
	if err := Preflight(criarZip(t, renomeado), dest); err != nil {
		t.Fatalf("pacote com subpasta rejeitado: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "release-1.2.3", "app", "Dockerfile")); err != nil {
		t.Fatalf("estrutura original não preservada: %v", err)
	}
}

func TestPreflightZipCaminhoBackslash(t *testing.T) {
	zipPath := criarZip(t, map[string]string{
		"app\\Dockerfile":       dockerfileSaudavel,
		"app\\docker-stack.yml": stackSaudavel,
		"app\\.env.example":     envSaudavel,
	})
	if err := Preflight(zipPath, t.TempDir()); err != nil {
		t.Fatalf("pacote com separadores Windows rejeitado: %v", err)
	}
	if !strings.Contains(dockerfileSaudavel, "HEALTHCHECK") {
		t.Fatal("sanidade do fixture")
	}
}
