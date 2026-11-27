---
name: legoloop-kraften-deploy
description: Use quando precisar executar o processo legoloop-kraften-deploy: deploy determinístico de aplicações do Studio Launcher em 5 nós zero LLM (preflight anti-zip-slip, sync de metadados e CNAME Cloudflare, deploy com stream do deploy-job, gates binários e recibo DEPLOY_RECEIPT.json). NÃO use para outros Processos.
consumes: Artifact:Zip
produces: Artifact:Receipt
artifact: DEPLOY_RECEIPT.json
gate: file_exists('DEPLOY_RECEIPT.json')
suggested_model: none
---

# 🚀 legoloop-kraften-deploy — Deploy Determinístico do Studio Launcher

> **Versão:** 1.0.0
> **Motor oficial:** `legoloop kraften-deploy` (implementação em `internal/kraften` e `internal/cli/kraften_deploy.go`, `internal/cli/kraften_etapas.go`).
> **Regra de ouro:** Processo 100% determinístico — zero LLM, zero agentes intermediários. Go nativo + Docker/Cloudflare REST v4.

---

## 1. Contexto & Contratos de I/O

**Quando usar:** para publicar uma aplicação qualquer (Go, Node, Python, Rust, PHP, estático) na VPS via Studio Launcher, com subdomínio e TLS automáticos.
**Quando NÃO usar:** para qualquer outro Processo.

* **Consumes:** pacote `.zip` com os 3 artefatos universais — `Dockerfile` (com `HEALTHCHECK`), `docker-stack.yml` com o campo `networks` declarado (o nome da rede é configuração do ambiente) e `.env.example` com `ADMIN_USER`/`ADMIN_PASS`.
* **Produces:** `DEPLOY_RECEIPT.json` auditável (prova dos gates `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`, `REDE=1`).

## 2. Procedimento (5 nós determinísticos em sequência)

1. **`n01_preflight_zip`** — `legoloop kraften-deploy preflight --zip PACOTE --dest DIR`: extração anti-zip-slip (todos os nomes validados antes de gravar qualquer byte) e validação dos 3 artefatos universais. Grava `preflight.json`.
2. **`n02_sync_metadata`** — `legoloop kraften-deploy sync --app ID --subdominio SUB --alvo HOST --launcher URL ...`: cria o CNAME Cloudflare de forma idempotente **antes** de qualquer consulta DNS e registra os metadados no Studio Launcher. Grava `sync.json`.
3. **`n03_deploy_stream`** — `legoloop kraften-deploy deploy --app ID --zip PACOTE --launcher URL ...`: envia o `.zip` (update-zip ou install-zip) e acompanha o `deploy-job` até `status == "done"`; `failed` ou timeout interrompem. Grava `deploy.json`.
4. **`n04_verify_gates`** — `legoloop kraften-deploy gates --health URL --host-novo FQDN --rede REDE --servico SERVICO ...`: prova matematicamente `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`, `REDE=1` (rede compartilhada anexada no Swarm, via `docker network inspect`); qualquer divergência reprova com exit 1. Grava `gates.json`.
5. **`n05_receipt`** — `legoloop kraften-deploy receipt --app ID --dest DIR`: consolida `deploy.json` + `gates.json` no `DEPLOY_RECEIPT.json`. Sem gates verdes não há recibo de sucesso.

**Modo completo (5 nós numa chamada):** `legoloop kraften-deploy --zip PACOTE --app ID --subdominio SUB --alvo HOST --dominio-base DOMINIO --launcher URL ...`

## 2b. Vias de entrada e coleta interativa

O Processo atende as 3 vias do LegoLoop: **terminal** (CLI), **agente/chat** (linguagem natural ou slash command) e **Studio** (botão Executar Processo). Parâmetros obrigatórios ausentes são **coletados interativamente, um a um**, em qualquer via — a coleta de parâmetros é a mesma (`completarKraften` no CLI; o agente pergunta conversando; o Studio envia o formulário). O domínio mestre (`--dominio-base`) e o nome da rede compartilhada são **configuração variável do ambiente** (`.env` mestre do Launcher), nunca valores fixos.

## 3. Invariantes & Regras Negativas

1. **NUNCA** edite o `Dockerfile` da aplicação (fronteira universal do `.zip`).
2. **NUNCA** insira agentes de IA/LLM nos nós.
3. **NUNCA** aceite pacotes sem os 3 artefatos universais, sem o campo `networks` declarado ou com caminhos Zip Slip (`..`).
3b. **NUNCA** fixe nomes de rede ou domínio mestre no código: ambos são parâmetros de configuração do ambiente.
4. **NUNCA** emita `DEPLOY_RECEIPT.json` com gates reprovados.
5. **NUNCA** duplique CNAMEs: `EnsureCNAME` consulta antes de criar.

## 4. Comandos de Verificação & Testes de Terminal

```bash
go test -count=1 ./internal/kraften/ ./internal/cli/ -run "TestKraften|TestPreflight|TestSync|TestDeploy|TestVerify|TestColetar|TestReceipt|TestExecutePipeline"
legoloop validate <caminho_do_DAG_do_processo>
```

## 5. Critério de Parada (Stopping Checklist)

- [ ] `n01_preflight_zip` aprovado com `preflight.json` gravado.
- [ ] `n02_sync_metadata` aprovado com `sync.json` e CNAME confirmado.
- [ ] `n03_deploy_stream` aprovado com `deploy.json` (`status == "done"`).
- [ ] `n04_verify_gates` aprovado com `gates.json` (provas binárias verdes, incluindo `REDE=1`).
- [ ] `n05_receipt` emitido: `DEPLOY_RECEIPT.json` íntegro.
- [ ] O DAG do processo passa em `legoloop validate` com exit code 0.

## 6. DAG & Nós Canônicos

| ID do Nó | Nome do Nó | Descrição |
| :--- | :--- | :--- |
| `n01_preflight_zip` | Preflight do pacote | Extração anti-zip-slip e validação dos 3 artefatos universais |
| `n02_sync_metadata` | Sincronização de metadados | CNAME idempotente antes de qualquer consulta DNS + registro no Launcher |
| `n03_deploy_stream` | Deploy com stream | Envio do `.zip` e acompanhamento do `deploy-job` até `done` |
| `n04_verify_gates` | Gates binários | `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`, `REDE=1` |
| `n05_receipt` | Recibo auditável | `DEPLOY_RECEIPT.json` com as provas consolidadas |
