---
schema: progresso.v1
iniciativa: "PLAN.md"
plano: "docs/plans/PLAN.md"
atualizado_em: "2026-11-26T20:55:29Z"
---

# Progresso Vivo: PLAN.md

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 7 | 0 | 0 | 0 | 0 | 7 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-KD-01.1` | E01 | Criar contratos estritos de DeployConfig e DeployReceipt com validacao sentinela ErrInvalidDeployInput | AC-E01.3 |
| `TICK-KD-01.2` | E01 | Implementar n01_preflight_zip: extracao anti-zip-slip e validacao dos 3 artefatos universais (Dockerfile, docker-stack.yml em network_conexao, HEALTHCHECK, ADMIN_USER/ADMIN_PASS) | AC-E01.1, AC-E01.2, AC-E01.3 |
| `TICK-KD-01.3` | E01 | Implementar n02_sync_metadata: cliente Cloudflare REST v4 idempotente (cfEnsureCNAME) e registro de metadados no Studio Launcher com CNAME criado antes de qualquer consulta DNS | AC-E01.4 |
| `TICK-KD-02.1` | E02 | Implementar n03_deploy_stream: envio do .zip para o Studio Launcher e acompanhamento do deploy-job ate status done | AC-E02.1 |
| `TICK-KD-02.2` | E02 | Implementar n04_verify_gates: gates binarios CF_OLD=0, CF_NEW=1, TLS=0, HTTP=200 com aborto sem recibo de sucesso em HTTP 500 | AC-E02.1, AC-E02.2 |
| `TICK-KD-02.3` | E02 | Implementar n05_receipt: gravacao do DEPLOY_RECEIPT.json integro na raiz do projeto | AC-NF.1 |
| `TICK-KD-02.4` | E02 | Implementar costura ExecutePipeline dos 5 nos e comando CLI legoloop kraften-deploy em cmd/legoloop | AC-E02.1, AC-NF.1 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| nenhum | - | - | - | - |
