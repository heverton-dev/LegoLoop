# Meso — Arquitetura

## Os 5 Nós Determinísticos

| Nó | Comando | Artefato | Gate |
| :--- | :--- | :--- | :--- |
| `n01_preflight_zip` | `legoloop kraften-deploy preflight` | `preflight.json` | 3 artefatos universais + anti-zip-slip |
| `n02_sync_metadata` | `legoloop kraften-deploy sync` | `sync.json` | CNAME criado antes de qualquer consulta DNS |
| `n03_deploy_stream` | `legoloop kraften-deploy deploy` | `deploy.json` | `deploy-job` com `status == "done"` |
| `n04_verify_gates` | `legoloop kraften-deploy gates` | `gates.json` | `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200` |
| `n05_receipt` | `legoloop kraften-deploy receipt` | `DEPLOY_RECEIPT.json` | recibo íntegro, só com gates verdes |

## Contratos de Dados

- **`DeployConfig`:** `app_id`, `zip_path`, `subdomain`, `target_host` — validação sentinela `ErrInvalidDeployInput` (vazio ou `..`).
- **`DeployReceipt`:** `app_id`, `status`, `timestamp`, `cf_status`, `http_code` — o artefato de prova.
- **`Gates`:** `cf_old`, `cf_new`, `tls_verdict`, `http_code` — as 4 provas binárias.

## Costuras de Teste

- **Primária:** `ExecutePipeline(ctx, cfg, deps)` em `internal/kraften/pipeline.go` — fronteiras injetáveis (HTTP, TLS, contagens DNS).
- **Secundárias:** `Preflight` (anti-zip-slip) e `EnsureCNAME` (idempotência).
- **Dublês apenas nas bordas:** clientes HTTP, relógio, filesystem temporário.

## Fluxo de Artefatos

```
pacote.zip → preflight.json → sync.json → deploy.json → gates.json → DEPLOY_RECEIPT.json
```

Cada nó consome o artefato do anterior e grava o seu: a cadeia inteira é auditável passo a passo.
