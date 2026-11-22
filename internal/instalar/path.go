package instalar

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/escrita"
)

const nomeBinarioGrafo = "codebase-memory-mcp"

// CriarAtalhoGrafo cria o atalho do binário do grafo na pasta de binários do usuário
// (.cmd no Windows e script sh 0o755 em Linux/macOS) quando o binário fonte existe.
// Com aplicar=false apenas devolve o destino planejado sem gravar nada em disco.
func CriarAtalhoGrafo(dirBin, binarioFonte, so string, aplicar bool) (string, error) {
	fonte := strings.TrimSpace(binarioFonte)
	if fonte == "" || !arquivoExiste(fonte) {
		return "", nil
	}
	dir := strings.TrimSpace(dirBin)
	if dir == "" {
		return "", errors.New("instalar: diretorio de binarios vazio")
	}
	alvo, err := normalizarSO(so)
	if err != nil {
		return "", err
	}
	destino, conteudo, perm := montarAtalhoGrafo(dir, fonte, alvo)
	if !aplicar {
		return destino, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := escrita.GravarAtomicoComPerm(destino, []byte(conteudo), perm); err != nil {
		return "", err
	}
	return destino, nil
}

func normalizarSO(so string) (string, error) {
	alvo := strings.ToLower(strings.TrimSpace(so))
	if alvo == "" {
		alvo = runtime.GOOS
	}
	if alvo != "windows" && alvo != "linux" && alvo != "darwin" {
		return "", fmt.Errorf("instalar: SO nao suportado: %s", so)
	}
	return alvo, nil
}

func montarAtalhoGrafo(dir, fonte, so string) (string, string, os.FileMode) {
	if so == "windows" {
		destino := filepath.Join(dir, nomeBinarioGrafo+".cmd")
		conteudo := fmt.Sprintf("@echo off\r\nsetlocal\r\n\"%s\" %%*\r\nendlocal\r\n", fonte)
		return destino, conteudo, 0o644
	}
	destino := filepath.Join(dir, nomeBinarioGrafo)
	conteudo := fmt.Sprintf("#!/usr/bin/env sh\nset -eu\nexec \"%s\" \"$@\"\n", fonte)
	return destino, conteudo, 0o755
}

func resolverDirBinarios(op Opcao) string {
	if strings.TrimSpace(op.DirBinarios) != "" {
		return op.DirBinarios
	}
	return filepath.Join(op.Home, ".local", "bin")
}

func resolverBinarioGrafo(op Opcao) string {
	if strings.TrimSpace(op.BinarioGrafo) != "" {
		return op.BinarioGrafo
	}
	nome := nomeBinarioGrafo
	if strings.EqualFold(op.SO, "windows") {
		nome += ".exe"
	}
	return filepath.Join(op.Home, "tools", nomeBinarioGrafo, nome)
}

func arquivoExiste(caminho string) bool {
	info, err := os.Stat(caminho)
	return err == nil && !info.IsDir()
}
