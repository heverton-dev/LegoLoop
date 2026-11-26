---
id: PLAN-ROUTE
status: DRAFT
date: 2026-11-26
schema_version: 2
upstream: [SPEC]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E01.4, AC-E02.1, AC-E02.2, AC-NF.1]
---

# 🗺️ PLAN: Transformar pipeline de deploy determinístico do Studio Launcher (`legoloop-kraften-deploy`)

> **Status:** Em Planejamento  
> **Data:** 2026-11-26  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/INTENT.md](../intents/INTENT.md)  
> **Referência da Especificação:** [docs/specs/SPEC.md](../specs/SPEC.md)  
> **Espelho Executável JSON:** [docs/routes/ROUTE.json](../routes/ROUTE.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Peça, sem context bloat)
* **Nível de Supervisão:** `epic`
* **Permissão de Subagentes:** `false` ou `ephemeral` (AGENTS.md, regra 4: agente de harness proibido; subagente efêmero só com as seis condições)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Peça:** `180s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Nó 01 e 02: Preflight Zip Seguro e Sincronização de Metadados | 3 | `AC-E01.3`, `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4` | `micro` |
| **E02** | Nó 03, 04 e 05: Deploy Stream, Verificação de Gates e Recibo | 4 | `AC-E02.1`, `AC-E02.1`, `AC-E02.2`, `AC-NF.1`, `AC-E02.1`, `AC-NF.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Nó 01 e 02: Preflight Zip Seguro e Sincronização de Metadados

#### Ticket `TICK-01.1` — Criar contratos estritos de DeployConfig e DeployReceipt com validacao sentinela ErrInvalidDeployInput
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/schemas.go` · Spec Anchor `docs/specs/SPEC.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/schemas_test.go`
* **Outputs:** Artifacts `internal/kraften/schemas.go, internal/kraften/schemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestDeployConfig`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/schemas.go') && file_exists('internal/kraften/schemas_test.go')`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.2` — Implementar n01_preflight_zip: extracao anti-zip-slip e validacao dos 3 artefatos universais (Dockerfile, docker-stack.yml em network_conexao, HEALTHCHECK, ADMIN_USER/ADMIN_PASS)
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/preflight.go` · Spec Anchor `docs/specs/SPEC.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/preflight_test.go`
* **Outputs:** Artifacts `internal/kraften/preflight.go, internal/kraften/preflight_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestPreflightZip`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/preflight.go') && file_exists('internal/kraften/preflight_test.go')`
* **Spec Ref:** `AC-E01.1, AC-E01.2, AC-E01.3`

#### Ticket `TICK-01.3` — Implementar n02_sync_metadata: cliente Cloudflare REST v4 idempotente (cfEnsureCNAME) e registro de metadados no Studio Launcher com CNAME criado antes de qualquer consulta DNS
* **Iniciativa:** `nenhum` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/cloudflare.go, internal/kraften/sync.go` · Spec Anchor `docs/specs/SPEC.md#AC-E01.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/cloudflare_test.go`
* **Outputs:** Artifacts `internal/kraften/cloudflare.go, internal/kraften/cloudflare_test.go, internal/kraften/sync.go, internal/kraften/sync_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestSyncMetadata`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/cloudflare.go') && file_exists('internal/kraften/cloudflare_test.go') && file_exists('internal/kraften/sync.go') && file_exists('internal/kraften/sync_test.go')`
* **Spec Ref:** `AC-E01.4`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Nó 03, 04 e 05: Deploy Stream, Verificação de Gates e Recibo

#### Ticket `TICK-02.1` — Implementar n03_deploy_stream: envio do .zip para o Studio Launcher e acompanhamento do deploy-job ate status done
* **Iniciativa:** `nenhum` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/deploy.go` · Spec Anchor `docs/specs/SPEC.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/deploy_test.go`
* **Outputs:** Artifacts `internal/kraften/deploy.go, internal/kraften/deploy_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestDeployStream`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/deploy.go') && file_exists('internal/kraften/deploy_test.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Implementar n04_verify_gates: gates binarios CF_OLD=0, CF_NEW=1, TLS=0, HTTP=200 com aborto sem recibo de sucesso em HTTP 500
* **Iniciativa:** `nenhum` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/gates.go` · Spec Anchor `docs/specs/SPEC.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/gates_test.go`
* **Outputs:** Artifacts `internal/kraften/gates.go, internal/kraften/gates_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestVerifyGates`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/gates.go') && file_exists('internal/kraften/gates_test.go')`
* **Spec Ref:** `AC-E02.1, AC-E02.2`

#### Ticket `TICK-02.3` — Implementar n05_receipt: gravacao do DEPLOY_RECEIPT.json integro na raiz do projeto
* **Iniciativa:** `nenhum` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/receipt.go` · Spec Anchor `docs/specs/SPEC.md#AC-NF.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/receipt_test.go`
* **Outputs:** Artifacts `internal/kraften/receipt.go, internal/kraften/receipt_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestReceipt`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/receipt.go') && file_exists('internal/kraften/receipt_test.go')`
* **Spec Ref:** `AC-NF.1`

#### Ticket `TICK-02.4` — Implementar costura ExecutePipeline dos 5 nos e comando CLI legoloop kraften-deploy em cmd/legoloop
* **Iniciativa:** `nenhum` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/kraften/pipeline.go, internal/cli/kraften_deploy.go` · Spec Anchor `docs/specs/SPEC.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/kraften/pipeline_test.go`
* **Outputs:** Artifacts `internal/kraften/pipeline.go, internal/kraften/pipeline_test.go, internal/cli/kraften_deploy.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/kraften/ -run TestExecutePipeline`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/kraften/pipeline.go') && file_exists('internal/kraften/pipeline_test.go') && file_exists('internal/cli/kraften_deploy.go')`
* **Spec Ref:** `AC-E02.1, AC-NF.1`

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal em `docs/reports/`
5. Registro de histórico da sessão via `faz-session --id <ID> --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*Plano de execução fatiado com 7 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
