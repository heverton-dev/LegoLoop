# Manual Operacional — `legoloop-kraften-deploy`

- **Data:** 2026-11-26
- **Versão:** 1.0.0
- **Motor CLI:** `legoloop kraften-deploy`
- **Pacote Go:** `internal/kraften` + `internal/cli/kraften_deploy.go` + `internal/cli/kraften_etapas.go`
- **DAG Canônico:** `processos/legoloop-kraften-deploy/pipeline.dag.json`
- **Skill Universal:** `skills/legoloop-kraften-deploy/SKILL.md`

---

## 1. Propósito

O Processo `legoloop-kraften-deploy` publica aplicações na VPS via Studio Launcher de forma 100% determinística (zero LLM): valida o pacote, sincroniza subdomínio e DNS na Cloudflare, executa o deploy no Docker Swarm, prova os gates binários e emite o recibo auditável `DEPLOY_RECEIPT.json`. Atende qualquer stack com `Dockerfile` e `docker-stack.yml` (Go, Node, Python, Rust, PHP, estático).

## 2. Comandos e Subcomandos

| Subcomando | Entrada | Saída e Gate |
| :--- | :--- | :--- |
| `legoloop kraften-deploy preflight --zip PACOTE [--dest DIR]` | `.zip` da app | Extrai sem Zip Slip, valida os 3 artefatos universais; grava `preflight.json` (`exit 0` se válido) |
| `legoloop kraften-deploy sync --app ID --subdominio SUB --alvo HOST --launcher URL --user U --pass P --cloudflare URL --zone ID --token T [--dest DIR]` | metadados + credenciais | Cria o CNAME idempotente antes de qualquer consulta DNS e registra no Launcher; grava `sync.json` |
| `legoloop kraften-deploy deploy --app ID --zip PACOTE --launcher URL --user U --pass P [--acao update-zip\|install-zip] [--dest DIR]` | `.zip` validado | Envio multipart e acompanhamento do `deploy-job` até `status == "done"`; grava `deploy.json` |
| `legoloop kraften-deploy gates --health URL --host-novo FQDN --rede REDE --servico SERVICO [--host-antigo FQDN] --cloudflare URL --zone ID --token T [--dest DIR]` | endpoints + DNS + Swarm | Prova `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`, `REDE=1`; grava `gates.json` (`exit 0` só com as 5 provas verdes) |
| `legoloop kraften-deploy receipt --app ID [--dest DIR]` | `deploy.json` + `gates.json` | Emite `DEPLOY_RECEIPT.json` consolidado |
| `legoloop kraften-deploy --zip ... --app ... --subdominio ... --alvo ...` | tudo acima | Executa os 5 nós numa chamada (costura `ExecutePipeline`) |

**Modo DAG:** `legoloop run processos/legoloop-kraften-deploy/pipeline.dag.json -var ZIP=... -var APP=... -var SUBDOMINIO=... -var ALVO=... -var LAUNCHER=... ...` — os 5 nós rodam em sequência topológica com artefatos intermediários (`preflight.json` → `sync.json` → `deploy.json` → `gates.json` → `DEPLOY_RECEIPT.json`).

## 3. Regras Inegociáveis

1. O `Dockerfile` da aplicação **nunca** é mutado pelo Processo.
2. Zero LLM e zero agentes intermediários nos nós.
3. Pacote sem `Dockerfile` (com `HEALTHCHECK`), sem `docker-stack.yml` com bloco `networks` declarado ou sem `.env.example` com `ADMIN_USER`/`ADMIN_PASS` é rejeitado. O nome da rede é configuração do ambiente, não regra do Processo.
4. Caminhos com `..` (Zip Slip) são recusados **antes** de qualquer gravação.
5. `DEPLOY_RECEIPT.json` só nasce com gates verdes (`CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`, `REDE=1`).
6. O CNAME novo é criado antes de qualquer consulta DNS (anti cache-negativo) e nunca duplicado.

## 4. Prova Mecânica

```bash
go test -count=1 ./internal/kraften/ ./internal/cli/ -run "TestKraften|TestPreflight|TestSync|TestDeploy|TestVerify|TestColetar|TestReceipt|TestExecutePipeline"
legoloop validate processos/legoloop-kraften-deploy/pipeline.dag.json
```
