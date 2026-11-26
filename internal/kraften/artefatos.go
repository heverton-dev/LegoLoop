package kraften

import (
	"os"
	"path/filepath"
	"strings"
)

// validarArtefatos confere os 3 artefatos universais do pacote na raiz do app.
func validarArtefatos(raiz string) error {
	if err := validarDockerfile(raiz); err != nil {
		return err
	}
	if err := validarStack(raiz); err != nil {
		return err
	}
	return validarVarsPadrao(raiz)
}

// validarDockerfile exige Dockerfile com HEALTHCHECK.
func validarDockerfile(raiz string) error {
	corpo, err := os.ReadFile(filepath.Join(raiz, "Dockerfile"))
	if err != nil {
		return ErrInvalidDeployInput
	}
	if !strings.Contains(string(corpo), "HEALTHCHECK") {
		return ErrInvalidDeployInput
	}
	return nil
}

// validarStack exige docker-stack.yml conectado à rede network_conexao.
func validarStack(raiz string) error {
	corpo, err := os.ReadFile(filepath.Join(raiz, "docker-stack.yml"))
	if err != nil {
		return ErrInvalidDeployInput
	}
	if !strings.Contains(string(corpo), "network_conexao") {
		return ErrInvalidDeployInput
	}
	return nil
}

// validarVarsPadrao exige .env.example declarando ADMIN_USER e ADMIN_PASS.
func validarVarsPadrao(raiz string) error {
	corpo, err := os.ReadFile(filepath.Join(raiz, ".env.example"))
	if err != nil {
		return ErrInvalidDeployInput
	}
	texto := string(corpo)
	for _, chave := range []string{"ADMIN_USER", "ADMIN_PASS"} {
		if !strings.Contains(texto, chave) {
			return ErrInvalidDeployInput
		}
	}
	return nil
}
