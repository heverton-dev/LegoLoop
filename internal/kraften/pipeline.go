package kraften

import (
	"time"
)

// PipelineDeps reúne as fronteiras injetáveis do Processo determinístico.
type PipelineDeps struct {
	Launcher     *LauncherClient
	Cloudflare   *CloudflareClient
	Dest         string
	HealthURL    string
	HostAntigo   string
	Contagens    func(nome string) int
	TLSVerify    func(host string) int
	HTTPCheck    func(url string) (int, error)
	RedeAnexada  func() int
	PollInterval time.Duration
}

// ExecutePipeline encadeia os 5 nós determinísticos e emite o recibo de prova.
func ExecutePipeline(cfg DeployConfig, deps PipelineDeps) (DeployReceipt, error) {
	if err := cfg.Validate(); err != nil {
		return DeployReceipt{}, err
	}
	if err := Preflight(cfg.ZipPath, deps.Dest); err != nil {
		return DeployReceipt{}, err
	}
	if err := SyncMetadata(cfg, deps.Launcher, deps.Cloudflare); err != nil {
		return DeployReceipt{}, err
	}
	if err := etapaDeploy(cfg, deps); err != nil {
		return DeployReceipt{}, err
	}
	gates, err := ColetarGates(deps.Contagens, deps.TLSVerify, deps.HTTPCheck, deps.RedeAnexada, deps.HealthURL, deps.HostAntigo, FqdnDe(cfg.Subdomain, cfg.BaseDomain))
	if err != nil {
		return DeployReceipt{}, err
	}
	return emitirRecibo(cfg, gates, deps.Dest)
}

// etapaDeploy executa o nó n03 com o intervalo de polling configurado.
func etapaDeploy(cfg DeployConfig, deps PipelineDeps) error {
	if deps.PollInterval > 0 {
		deps.Launcher.PollInterval = deps.PollInterval
	}
	_, err := DeployStream(cfg, deps.Launcher, "")
	return err
}

// emitirRecibo grava o DEPLOY_RECEIPT.json apenas em caso de gates verdes.
func emitirRecibo(cfg DeployConfig, g Gates, dest string) (DeployReceipt, error) {
	rcp := DeployReceipt{
		AppID:     cfg.AppID,
		Status:    "done",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		CFStatus:  g.CFNew,
		HTTPCode:  g.HTTPCode,
	}
	if _, err := GravarReceipt(rcp, dest); err != nil {
		return DeployReceipt{}, err
	}
	return rcp, nil
}
