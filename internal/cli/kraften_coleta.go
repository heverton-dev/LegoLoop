package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/kraften"
)

// perguntar pede um parâmetro ao usuário (fronteira injetável em testes).
// As 3 vias de entrada (terminal, agente/chat, Studio) alimentam este fluxo.
var perguntar = perguntarStdin

// perguntarStdin lê uma linha do stdin com rótulo no stderr.
func perguntarStdin(rotulo string, d Deps) (string, error) {
	fmt.Fprintf(d.Stderr, "%s: ", rotulo)
	leitor := bufio.NewReader(d.Stdin)
	linha, err := leitor.ReadString('\n')
	if err != nil && linha == "" {
		return "", err
	}
	return strings.TrimSpace(linha), nil
}

// completarKraften pergunta, um a um, os parâmetros obrigatórios ausentes.
func completarKraften(o *flagsKraften, exigidos []string, d Deps) error {
	alvos := map[string]*string{
		"--zip": &o.zip, "--app": &o.app, "--subdominio": &o.sub, "--alvo": &o.alvo,
		"--dominio-base": &o.baseDomain, "--launcher": &o.launcher,
		"--health": &o.health, "--host-novo": &o.hostNovo,
		"--rede": &o.rede, "--servico": &o.servico,
	}
	for _, nome := range exigidos {
		campo := alvos[nome]
		if campo == nil || *campo != "" {
			continue
		}
		valor, err := perguntar(nome, d)
		if err != nil {
			return err
		}
		*campo = valor
	}
	return nil
}

// exigidosDeEtapa lista os parâmetros obrigatórios de cada subcomando.
func exigidosDeEtapa(sub string) []string {
	switch sub {
	case "preflight":
		return []string{"--zip"}
	case "sync":
		return []string{"--app", "--subdominio", "--alvo", "--dominio-base", "--launcher"}
	case "deploy":
		return []string{"--app", "--zip", "--launcher"}
	case "gates":
		return []string{"--health", "--host-novo", "--rede", "--servico"}
	case "receipt":
		return []string{"--app"}
	}
	return nil
}

// semValor informa se o parâmetro obrigatório continua vazio após a coleta.
func semValor(o *flagsKraften, exigidos []string) bool {
	alvos := map[string]string{
		"--zip": o.zip, "--app": o.app, "--subdominio": o.sub, "--alvo": o.alvo,
		"--dominio-base": o.baseDomain, "--launcher": o.launcher,
		"--health": o.health, "--host-novo": o.hostNovo,
		"--rede": o.rede, "--servico": o.servico,
	}
	for _, nome := range exigidos {
		if alvos[nome] == "" {
			return true
		}
	}
	return false
}

// configDeKraften monta o DeployConfig com domínio base configurável.
func configDeKraften(o *flagsKraften, zipPath string) kraften.DeployConfig {
	return kraften.DeployConfig{
		AppID: o.app, ZipPath: zipPath, Subdomain: o.sub,
		TargetHost: o.alvo, BaseDomain: o.baseDomain,
	}
}
