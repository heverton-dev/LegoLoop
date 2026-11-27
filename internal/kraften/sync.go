package kraften

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LauncherClient fala com a API administrativa do Studio Launcher.
type LauncherClient struct {
	base         string
	user         string
	pass         string
	hc           *http.Client
	PollInterval time.Duration
}

// NovoLauncherClient cria o cliente da API do Launcher com autenticação básica.
func NovoLauncherClient(base, user, pass string, hc *http.Client) *LauncherClient {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &LauncherClient{base: strings.TrimRight(base, "/"), user: user, pass: pass, hc: hc}
}

// RegistrarApp faz upsert dos metadados da app em POST /api/admin/apps.
func (l *LauncherClient) RegistrarApp(cfg DeployConfig) error {
	carga := map[string]any{
		"app_id":      cfg.AppID,
		"subdomain":   cfg.Subdomain,
		"target_host": cfg.TargetHost,
	}
	buf, err := json.Marshal(carga)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, l.base+"/api/admin/apps", strings.NewReader(string(buf)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(l.user, l.pass)
	req.Header.Set("Content-Type", "application/json")
	resp, err := l.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("kraften: registro de metadados recusado (%d)", resp.StatusCode)
	}
	return nil
}

// SyncMetadata garante o CNAME antes de qualquer consulta DNS e em seguida
// registra os metadados da app no Studio Launcher.
func SyncMetadata(cfg DeployConfig, l *LauncherClient, cf *CloudflareClient) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if err := cf.EnsureCNAME(FqdnDe(cfg.Subdomain, cfg.BaseDomain), cfg.TargetHost); err != nil {
		return err
	}
	return l.RegistrarApp(cfg)
}

// FqdnDe monta o FQDN sobre o domínio base configurado (variável por ambiente).
// O subdomínio já completo (com ponto) passa direto.
func FqdnDe(subdomain, baseDomain string) string {
	if strings.Contains(subdomain, ".") {
		return subdomain
	}
	return subdomain + "." + baseDomain
}
