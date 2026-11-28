package cli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/kraften"
)

const usoKraftenDeploy = `uso: legoloop kraften-deploy --zip PACOTE.zip --app ID --subdominio SUB --alvo HOST [opções]
opções:
  --launcher URL     base da API do Studio Launcher (ex.: https://laucher.vpsconexao.org)
  --user USR         usuário administrativo do Launcher
  --pass PASS        senha administrativa do Launcher
  --cloudflare URL   base da API Cloudflare v4 (padrão: https://api.cloudflare.com/client/v4)
  --zone ID          ID da zona DNS na Cloudflare
  --token TOKEN      token de API da Cloudflare
  --health URL       rota de healthcheck da app publicada
  --dest DIR         diretório raiz de trabalho (padrão: .)
  --json             emite o recibo em JSON
`

// flagsKraften agrupa as opções do comando kraften-deploy.
type flagsKraften struct {
	zip, app, sub, alvo, launcher, user, pass, cf, zone, token, health, dest, hostAntigo string
	hostNovo, acao, baseDomain                                                           string
	rede, servico                                                                        string
	emJSON                                                                               bool
}

func parseKraften(args []string) (flagsKraften, bool) {
	fs := flag.NewFlagSet("kraften-deploy", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o flagsKraften
	fs.StringVar(&o.zip, "zip", "", "caminho do pacote .zip")
	fs.StringVar(&o.app, "app", "", "ID da aplicação")
	fs.StringVar(&o.sub, "subdominio", "", "subdomínio desejado")
	fs.StringVar(&o.alvo, "alvo", "", "host alvo do CNAME")
	fs.StringVar(&o.launcher, "launcher", "", "base da API do Studio Launcher")
	fs.StringVar(&o.user, "user", "", "usuário do Launcher")
	fs.StringVar(&o.pass, "pass", "", "senha do Launcher")
	fs.StringVar(&o.cf, "cloudflare", "https://api.cloudflare.com/client/v4", "base da API Cloudflare")
	fs.StringVar(&o.zone, "zone", "", "ID da zona DNS")
	fs.StringVar(&o.token, "token", "", "token da Cloudflare")
	fs.StringVar(&o.health, "health", "", "URL do healthcheck")
	fs.StringVar(&o.dest, "dest", ".", "diretório raiz de trabalho")
	fs.StringVar(&o.hostAntigo, "host-antigo", "", "FQDN do CNAME antigo (vazio = primeiro deploy)")
	fs.StringVar(&o.hostNovo, "host-novo", "", "FQDN do CNAME novo")
	fs.StringVar(&o.acao, "acao", "", "ação do envio do zip (update-zip ou install-zip)")
	fs.StringVar(&o.baseDomain, "dominio-base", "", "domínio mestre do ecossistema (configuração variável)")
	fs.StringVar(&o.rede, "rede", "", "nome da rede compartilhada do ecossistema")
	fs.StringVar(&o.servico, "servico", "", "nome do serviço no Docker Swarm")
	fs.BoolVar(&o.emJSON, "json", false, "emite o recibo em JSON")
	if err := fs.Parse(args); err != nil {
		return o, false
	}
	return o, true
}

// cmdKraftenDeploy executa o Processo determinístico de deploy de 5 nós sem LLM.
func cmdKraftenDeploy(ctx context.Context, args []string, d Deps) int {
	if len(args) > 0 && ehEtapaKraften(args[0]) {
		if code, ok := etapaKraften(args[0], args[1:], d); ok {
			return code
		}
	}
	o, ok := parseKraften(args)
	if !ok {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	exigidos := append(exigidosDeEtapa("preflight"), exigidosDeEtapa("sync")...)
	exigidos = append(exigidos, exigidosDeEtapa("deploy")...)
	exigidos = append(exigidos, exigidosDeEtapa("gates")...)
	if err := completarKraften(&o, exigidos, d); err != nil {
		fmt.Fprintf(d.Stderr, "kraften-deploy: %v\n", err)
		return 1
	}
	if semValor(&o, exigidos) {
		fmt.Fprint(d.Stderr, usoKraftenDeploy)
		return 2
	}
	cfg := configDeKraften(&o, o.zip)
	recibo, err := kraften.ExecutePipeline(cfg, depsDeKraften(&o))
	if err != nil {
		fmt.Fprintf(d.Stderr, "kraften-deploy: %v\n", err)
		return 1
	}
	return emitirReciboKraften(recibo, o.emJSON, d)
}

// depsDeKraften monta as fronteiras injetáveis do modo completo.
func depsDeKraften(o *flagsKraften) kraften.PipelineDeps {
	clienteCF := kraften.NovoCloudflareClient(o.cf, o.zone, o.token, http.DefaultClient)
	return kraften.PipelineDeps{
		Launcher:    kraften.NovoLauncherClient(o.launcher, o.user, o.pass, http.DefaultClient),
		Cloudflare:  clienteCF,
		Dest:        o.dest,
		HealthURL:   o.health,
		HostAntigo:  o.hostAntigo,
		Contagens:   contagemCNAME(clienteCF),
		TLSVerify:   verificarTLS,
		HTTPCheck:   checarHTTP,
		RedeAnexada: func() int { return redeAnexadaKraften(o.rede, o.servico) },
	}
}

// contagemCNAME mede quantos registros CNAME existem para o nome informado.
// Nome vazio indica ausência de host anterior (primeiro deploy) e conta zero.
func contagemCNAME(cf *kraften.CloudflareClient) func(string) int {
	return func(nome string) int {
		if nome == "" {
			return 0
		}
		n, err := cf.ContarCNAME(nome)
		if err != nil {
			return -1
		}
		return n
	}
}

// verificarTLS devolve 0 quando o handshake TLS é válido.
func verificarTLS(host string) int {
	conexao, err := tls.DialWithDialer(&net.Dialer{Timeout: 15 * time.Second}, "tcp", host+":443", &tls.Config{ServerName: host})
	if err != nil {
		return 1
	}
	defer conexao.Close()
	return 0
}

// checarHTTP devolve o código de status da rota informada.
func checarHTTP(url string) (int, error) {
	cliente := &http.Client{Timeout: 15 * time.Second}
	resp, err := cliente.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// redeAnexadaNoSwarm prova a anexação: consulta `docker network inspect` e
// conta os serviços anexados à rede configurada. Devolve 1 quando o serviço
// informado está anexado, 0 quando não está e -1 em erro de consulta.
func redeAnexadaNoSwarm(rede, servico string) int {
	saida, err := exec.Command("docker", "network", "inspect", rede,
		"--format", "{{range .Services}}{{.Name}} {{end}}").Output()
	if err != nil {
		return -1
	}
	for _, nome := range strings.Fields(string(saida)) {
		if nome == servico {
			return 1
		}
	}
	return 0
}

// agoraISO devolve o instante atual em formato RFC3339 UTC.
func agoraISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func emitirReciboKraften(r kraften.DeployReceipt, emJSON bool, d Deps) int {
	if emJSON {
		corpo, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			fmt.Fprintf(d.Stderr, "kraften-deploy: %v\n", err)
			return 1
		}
		fmt.Fprintln(d.Stdout, string(corpo))
		return 0
	}
	fmt.Fprintf(d.Stdout, "deploy %s: %s (CF=%d HTTP=%d) em %s\n", r.AppID, r.Status, r.CFStatus, r.HTTPCode, r.Timestamp)
	return 0
}
