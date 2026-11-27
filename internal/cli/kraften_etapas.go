package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/heverton-dev/LegoLoop/internal/kraften"
)

// fronteiraTLSKraften mede o handshake TLS (injetável em testes).
var fronteiraTLSKraften = verificarTLS

// redeAnexadaKraften prova a anexação da rede compartilhada no Swarm.
// Devolve 1 quando o serviço tem a rede configurada anexada (injetável em testes).
var redeAnexadaKraften = redeAnexadaNoSwarm

// ehEtapaKraften informa se o subcomando pertence aos 5 nós determinísticos.
func ehEtapaKraften(sub string) bool {
	switch sub {
	case "preflight", "sync", "deploy", "gates", "receipt":
		return true
	}
	return false
}

// etapaKraften despacha o subcomando de etapa do Processo kraften-deploy.
func etapaKraften(sub string, args []string, d Deps) (int, bool) {
	switch sub {
	case "preflight":
		return etapaPreflightKraften(args, d), true
	case "sync":
		return etapaSyncKraften(args, d), true
	case "deploy":
		return etapaDeployKraften(args, d), true
	case "gates":
		return etapaGatesKraften(args, d), true
	case "receipt":
		return etapaReceiptKraften(args, d), true
	}
	return 0, false
}

// etapaPreflightKraften executa n01_preflight_zip e grava preflight.json.
func etapaPreflightKraften(args []string, d Deps) int {
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := completarKraften(&o, exigidosDeEtapa("preflight"), d); err != nil {
		fmt.Fprintf(d.Stderr, "preflight: %v\n", err)
		return 1
	}
	if semValor(&o, exigidosDeEtapa("preflight")) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := kraften.Preflight(o.zip, o.dest); err != nil {
		fmt.Fprintf(d.Stderr, "preflight: %v\n", err)
		return 1
	}
	return gravarEtapaKraften(o.dest, "preflight.json", map[string]any{"zip": o.zip, "ok": true}, d)
}

// etapaSyncKraften executa n02_sync_metadata e grava sync.json.
func etapaSyncKraften(args []string, d Deps) int {
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := completarKraften(&o, exigidosDeEtapa("sync"), d); err != nil {
		fmt.Fprintf(d.Stderr, "sync: %v\n", err)
		return 1
	}
	if semValor(&o, exigidosDeEtapa("sync")) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	cfg := configDeKraften(&o, "sync")
	l := kraften.NovoLauncherClient(o.launcher, o.user, o.pass, http.DefaultClient)
	cf := kraften.NovoCloudflareClient(o.cf, o.zone, o.token, http.DefaultClient)
	if err := kraften.SyncMetadata(cfg, l, cf); err != nil {
		fmt.Fprintf(d.Stderr, "sync: %v\n", err)
		return 1
	}
	fqdn := kraften.FqdnDe(o.sub, o.baseDomain)
	return gravarEtapaKraften(o.dest, "sync.json", map[string]any{"fqdn": fqdn}, d)
}

// etapaDeployKraften executa n03_deploy_stream e grava deploy.json.
func etapaDeployKraften(args []string, d Deps) int {
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := completarKraften(&o, exigidosDeEtapa("deploy"), d); err != nil {
		fmt.Fprintf(d.Stderr, "deploy: %v\n", err)
		return 1
	}
	if semValor(&o, exigidosDeEtapa("deploy")) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	cfg := configDeKraften(&o, o.zip)
	l := kraften.NovoLauncherClient(o.launcher, o.user, o.pass, http.DefaultClient)
	status, err := kraften.DeployStream(cfg, l, o.acao)
	if err != nil {
		fmt.Fprintf(d.Stderr, "deploy: %v\n", err)
		return 1
	}
	return gravarEtapaKraften(o.dest, "deploy.json", map[string]any{"status": status}, d)
}

// etapaGatesKraften executa n04_verify_gates e grava gates.json.
func etapaGatesKraften(args []string, d Deps) int {
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := completarKraften(&o, exigidosDeEtapa("gates"), d); err != nil {
		fmt.Fprintf(d.Stderr, "gates: %v\n", err)
		return 1
	}
	if semValor(&o, exigidosDeEtapa("gates")) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	cf := kraften.NovoCloudflareClient(o.cf, o.zone, o.token, http.DefaultClient)
	gates, err := kraften.ColetarGates(contagemCNAME(cf), fronteiraTLSKraften, checarHTTP,
		func() int { return redeAnexadaKraften(o.rede, o.servico) }, o.health, o.hostAntigo, o.hostNovo)
	if err != nil {
		fmt.Fprintf(d.Stderr, "gates: %v\n", err)
		gravarEtapaKraften(o.dest, "gates.json", gates, d)
		return 1
	}
	return gravarEtapaKraften(o.dest, "gates.json", gates, d)
}

// etapaReceiptKraften executa n05_receipt e grava DEPLOY_RECEIPT.json.
func etapaReceiptKraften(args []string, d Deps) int {
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	if err := completarKraften(&o, exigidosDeEtapa("receipt"), d); err != nil {
		fmt.Fprintf(d.Stderr, "receipt: %v\n", err)
		return 1
	}
	if semValor(&o, exigidosDeEtapa("receipt")) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	rcp, err := montarReceipt(o.dest, o.app)
	if err != nil {
		fmt.Fprintf(d.Stderr, "receipt: %v\n", err)
		return 1
	}
	if _, err := kraften.GravarReceipt(rcp, o.dest); err != nil {
		fmt.Fprintf(d.Stderr, "receipt: %v\n", err)
		return 1
	}
	return 0
}

// montarReceipt consolida deploy.json e gates.json no recibo do Processo.
func montarReceipt(dest, app string) (kraften.DeployReceipt, error) {
	var deploy struct {
		Status string `json:"status"`
	}
	var gates kraften.Gates
	if err := lerJSONEtapa(dest, "deploy.json", &deploy); err != nil {
		return kraften.DeployReceipt{}, err
	}
	if err := lerJSONEtapa(dest, "gates.json", &gates); err != nil {
		return kraften.DeployReceipt{}, err
	}
	return kraften.DeployReceipt{
		AppID:     app,
		Status:    deploy.Status,
		Timestamp: agoraISO(),
		CFStatus:  gates.CFNew,
		HTTPCode:  gates.HTTPCode,
	}, nil
}

func lerJSONEtapa(dest, nome string, alvo any) error {
	corpo, err := os.ReadFile(filepath.Join(dest, nome))
	if err != nil {
		return err
	}
	return json.Unmarshal(corpo, alvo)
}

func gravarEtapaKraften(dest, nome string, carga any, d Deps) int {
	corpo, err := json.MarshalIndent(carga, "", "  ")
	if err != nil {
		fmt.Fprintf(d.Stderr, "%s: %v\n", nome, err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(dest, nome), corpo, 0o644); err != nil {
		fmt.Fprintf(d.Stderr, "%s: %v\n", nome, err)
		return 1
	}
	return 0
}
