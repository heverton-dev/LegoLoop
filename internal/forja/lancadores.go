package forja

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type definicaoLancador struct {
	nome     string
	argsWin  string
	argsUnix string
	envWin   string
	envUnix  string
}

func definicoesLancadores() []definicaoLancador {
	return []definicaoLancador{
		{
			nome:     "claude-minimo",
			argsWin:  "--allowedTools Bash,Read,Edit,Write,ToolSearch",
			argsUnix: "--allowedTools Bash,Read,Edit,Write,ToolSearch",
		},
		{
			nome:     "claude-gateway",
			argsWin:  "--setting-sources \"\" --allowedTools mcp__legoloop__call_mcp_tool",
			argsUnix: "--setting-sources \"\" --allowedTools mcp__legoloop__call_mcp_tool",
		},
		{
			nome:     "claude-completo",
			envWin:   "set ENABLE_CLAUDEAI_MCP_SERVERS=true",
			envUnix:  "ENABLE_CLAUDEAI_MCP_SERVERS=true",
			argsWin:  "--setting-sources project,local",
			argsUnix: "--setting-sources project,local",
		},
		{
			nome:    "claude-conectores",
			envWin:  "set ENABLE_CLAUDEAI_MCP_SERVERS=true",
			envUnix: "ENABLE_CLAUDEAI_MCP_SERVERS=true",
		},
	}
}

// GerarLancadores cria os lancadores do Claude adequados ao SO informado.
func GerarLancadores(dirBin, so string) ([]string, error) {
	dir := strings.TrimSpace(dirBin)
	if dir == "" {
		return nil, errors.New("forja: diretorio de binarios vazio")
	}
	alvo := strings.ToLower(strings.TrimSpace(so))
	if alvo != "windows" && alvo != "linux" && alvo != "darwin" {
		return nil, fmt.Errorf("forja: SO nao suportado: %s", so)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	defs := definicoesLancadores()
	criados := make([]string, 0, len(defs))
	for _, def := range defs {
		arq, err := gravarLancador(dir, alvo, def)
		if err != nil {
			return nil, err
		}
		criados = append(criados, arq)
	}
	return criados, nil
}

func gravarLancador(dir, so string, def definicaoLancador) (string, error) {
	if so == "windows" {
		caminho := filepath.Join(dir, def.nome+".cmd")
		conteudo := montarLancadorWindows(def)
		return caminho, os.WriteFile(caminho, []byte(conteudo), 0o644)
	}
	caminho := filepath.Join(dir, def.nome)
	conteudo := montarLancadorPosix(def)
	return caminho, os.WriteFile(caminho, []byte(conteudo), 0o755)
}

func montarLancadorWindows(def definicaoLancador) string {
	var linhas []string
	linhas = append(linhas, "@echo off", "setlocal")
	if def.envWin != "" {
		linhas = append(linhas, def.envWin)
	}
	cmd := "claude"
	if def.argsWin != "" {
		cmd += " " + def.argsWin
	}
	linhas = append(linhas, cmd+" %*", "endlocal")
	return strings.Join(linhas, "\r\n") + "\r\n"
}

func montarLancadorPosix(def definicaoLancador) string {
	var linhas []string
	linhas = append(linhas, "#!/usr/bin/env sh", "set -eu")
	prefixo := "exec"
	if def.envUnix != "" {
		prefixo += " " + def.envUnix
	}
	cmd := prefixo + " claude"
	if def.argsUnix != "" {
		cmd += " " + def.argsUnix
	}
	linhas = append(linhas, cmd+" \"$@\"")
	return strings.Join(linhas, "\n") + "\n"
}
