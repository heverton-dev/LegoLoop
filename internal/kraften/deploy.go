package kraften

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DeployStream envia o .zip ao Studio Launcher e acompanha o deploy-job até
// status "done"; qualquer outro desfecho final interrompe o Processo.
func DeployStream(cfg DeployConfig, l *LauncherClient, acao string) (string, error) {
	if err := cfg.Validate(); err != nil {
		return "", err
	}
	if acao == "" {
		acao = "update-zip"
	}
	job, err := l.EnviarZip(cfg, acao)
	if err != nil {
		return "", err
	}
	return l.AguardarJob(job, 30*time.Second)
}

// EnviarZip faz POST multipart em /api/admin/apps/<id>/<acao> e devolve o job.
func (l *LauncherClient) EnviarZip(cfg DeployConfig, acao string) (string, error) {
	corpo, contentType, err := montarMultipart(cfg)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/api/admin/apps/%s/%s", l.base, cfg.AppID, acao)
	req, err := http.NewRequest(http.MethodPost, url, corpo)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(l.user, l.pass)
	req.Header.Set("Content-Type", contentType)
	resp, err := l.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("kraften: envio do zip recusado (%d)", resp.StatusCode)
	}
	var carga struct {
		Job string `json:"job"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&carga); err != nil {
		return "", err
	}
	if carga.Job == "" {
		return "", ErrInvalidDeployInput
	}
	return carga.Job, nil
}

// montarMultipart empacota o .zip como campo "file" do formulário.
func montarMultipart(cfg DeployConfig) (*strings.Reader, string, error) {
	arq, err := os.Open(cfg.ZipPath)
	if err != nil {
		return nil, "", ErrInvalidDeployInput
	}
	defer arq.Close()
	var buf strings.Builder
	w := multipart.NewWriter(&buf)
	campo, err := w.CreateFormFile("file", filepath.Base(cfg.ZipPath))
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(campo, arq); err != nil {
		return nil, "", err
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return strings.NewReader(buf.String()), w.FormDataContentType(), nil
}

// AguardarJob consulta GET /api/admin/jobs/<id> até done, failed ou timeout.
func (l *LauncherClient) AguardarJob(job string, limite time.Duration) (string, error) {
	intervalo := l.PollInterval
	if intervalo <= 0 {
		intervalo = 2 * time.Second
	}
	prazo := time.Now().Add(limite)
	for {
		status, err := l.consultarJob(job)
		if err != nil {
			return "", err
		}
		switch status {
		case "done":
			return status, nil
		case "failed":
			return status, fmt.Errorf("kraften: deploy-job %s terminou em failed", job)
		}
		if time.Now().After(prazo) {
			return status, fmt.Errorf("kraften: deploy-job %s excedeu o tempo limite (%s)", job, status)
		}
		time.Sleep(intervalo)
	}
}

func (l *LauncherClient) consultarJob(job string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/admin/jobs/%s", l.base, job), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(l.user, l.pass)
	resp, err := l.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var carga struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&carga); err != nil {
		return "", err
	}
	if carga.Status == "" {
		return "", ErrInvalidDeployInput
	}
	return carga.Status, nil
}
