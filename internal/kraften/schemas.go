// Package kraften implementa o Processo determinístico de deploy do Studio
// Launcher (zero LLM): preflight de pacote, sincronização de metadados e DNS,
// deploy com stream, gates binários e recibo JSON auditável.
package kraften

import (
	"errors"
	"strings"
)

// ErrInvalidDeployInput sinaliza entrada de deploy inválida ou insegura.
var ErrInvalidDeployInput = errors.New("kraften: entrada de deploy inválida ou insegura")

// DeployConfig descreve o alvo de um deploy determinístico.
type DeployConfig struct {
	AppID      string `json:"app_id"`
	ZipPath    string `json:"zip_path"`
	Subdomain  string `json:"subdomain"`
	TargetHost string `json:"target_host"`
}

// DeployReceipt é o artefato de prova emitido ao fim do Processo.
type DeployReceipt struct {
	AppID     string `json:"app_id"`
	Status    string `json:"status"` // "done" ou "failed"
	Timestamp string `json:"timestamp"`
	CFStatus  int    `json:"cf_status"`
	HTTPCode  int    `json:"http_code"`
}

// Validate aplica as regras estritas de entrada do Processo.
func (c DeployConfig) Validate() error {
	if strings.TrimSpace(c.AppID) == "" || strings.TrimSpace(c.ZipPath) == "" {
		return ErrInvalidDeployInput
	}
	if strings.Contains(c.ZipPath, "..") {
		return ErrInvalidDeployInput
	}
	return nil
}
