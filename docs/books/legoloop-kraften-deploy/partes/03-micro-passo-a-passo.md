# Micro — Passo a Passo

## 1. Preflight (`n01_preflight_zip`)

```bash
legoloop kraften-deploy preflight --zip pacote.zip --dest .
```

- Abre o `.zip` e valida **todos** os nomes de entrada antes de gravar qualquer byte (anti-zip-slip: recusa `..`, absolutos e `:`).
- Localiza a subpasta do `Dockerfile` e confere: `HEALTHCHECK` no `Dockerfile`, bloco `networks` declarado no `docker-stack.yml`, `ADMIN_USER`/`ADMIN_PASS` no `.env.example`.
- Qualquer violação → `ErrInvalidDeployInput` e exit 1.

## 2. Sync (`n02_sync_metadata`)

```bash
legoloop kraften-deploy sync --app app01 --subdominio sub01 --alvo server.vpsconexao.org \
  --launcher https://laucher.vpsconexao.org --user U --pass P \
  --cloudflare https://api.cloudflare.com/client/v4 --zone ZONE --token TOKEN
```

- Calcula o FQDN (`sub01.vpsconexao.org`) e consulta `GET /zones/{zone}/dns_records?name=...`.
- Se não existir, cria `POST .../dns_records` (`CNAME`, `proxied: false`, `ttl: 1`); se existir, não duplica.
- Só então registra os metadados em `POST /api/admin/apps` — o CNAME nasce antes de qualquer consulta DNS do navegador.

## 3. Deploy (`n03_deploy_stream`)

```bash
legoloop kraften-deploy deploy --app app01 --zip pacote.zip --launcher URL --user U --pass P [--acao install-zip]
```

- Envia o `.zip` como multipart em `POST /api/admin/apps/<id>/<acao>` (padrão `update-zip`).
- Consulta `GET /api/admin/jobs/<id>` até `status == "done"`; `failed` interrompe com erro; timeout de 30s aborta.

## 4. Gates (`n04_verify_gates`)

```bash
legoloop kraften-deploy gates --health https://sub01.vpsconexao.org/health \
  --host-novo sub01.vpsconexao.org --host-antigo sub01-velho.vpsconexao.org ...
```

Prova binária (exit 0 só com as quatro):

| Prova | Esperado | Significado |
| :--- | :--- | :--- |
| `CF_OLD` | `0` | o CNAME antigo zerou |
| `CF_NEW` | `1` | o CNAME novo está ativo |
| `TLS` | `0` | handshake TLS válido (Let's Encrypt) |
| `HTTP` | `200` | rota `/health` respondendo |
| `REDE` | `1` | rede compartilhada anexada ao serviço no Swarm |

## 5. Recibo (`n05_receipt`)

```bash
legoloop kraften-deploy receipt --app app01 --dest .
```

- Consolida `deploy.json` + `gates.json` no `DEPLOY_RECEIPT.json`.
- Sem gates verdes o recibo de sucesso não existe — a falha do nó 4 interrompe a cadeia.

## Modo Completo e Modo DAG

- **Uma chamada:** `legoloop kraften-deploy --zip ... --app ... --subdominio ... --alvo ...` (costura `ExecutePipeline`).
- **DAG:** `legoloop run processos/legoloop-kraften-deploy/pipeline.dag.json -var ZIP=... -var APP=...` — os 5 nós em sequência topológica.
