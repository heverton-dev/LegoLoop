---
id: PLAN-2026-11-24_ROUTE_LEGOLOOP_DIORAMA_ATIVOS
status: DRAFT
date: 2026-11-24
schema_version: 2
upstream: [2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS]
covers: [AC-E01.1, AC-E01.2, AC-E02.1, AC-E02.2, AC-E03.1, AC-E03.2, AC-E03.3, AC-E04.1, AC-E04.2, AC-E04.3, AC-E05.1, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Processo Canonizado `legoloop-diorama-ativos` (Ciclo 2 — Etapas 3, 4 e 5)

> **Status:** Em Planejamento  
> **Data:** 2026-11-24  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-24_INTENT_LEGOLOOP_DIORAMA_ATIVOS.md](../intents/2026-11-24_INTENT_LEGOLOOP_DIORAMA_ATIVOS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md](../specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-24_ROUTE_LEGOLOOP_DIORAMA_ATIVOS.json](../routes/2026-11-24_ROUTE_LEGOLOOP_DIORAMA_ATIVOS.json)  
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
| **E01** | Épico E01 | 2 | `AC-E01.1`, `AC-E01.1`, `AC-E01.2`, `AC-NF.2` | `micro` |
| **E02** | Épico E02 | 1 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Épico E03 | 1 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `micro` |
| **E04** | Épico E04 | 2 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3`, `AC-NF.1`, `AC-NF.2` | `micro` |
| **E05** | Épico E05 | 1 | `AC-E05.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Épico E01

#### Ticket `TICK-01.1` — Criar contratos estritos de Elenco, Sprites, Cenarios e Storyboard
* **Iniciativa:** `D02` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/schemas.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/schemas_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/schemas.go, internal/diorama/ativos/schemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestSchemas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/schemas.go') && file_exists('internal/diorama/ativos/schemas_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Implementar pre-voo do Ciclo 2 com Path Jail, contratos do Ciclo 1 e trava flow_projeto
* **Iniciativa:** `D02` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/prevoo.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/prevoo_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/prevoo.go, internal/diorama/ativos/prevoo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestPreVoo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/prevoo.go') && file_exists('internal/diorama/ativos/prevoo_test.go')`
* **Spec Ref:** `AC-E01.1, AC-E01.2, AC-NF.2`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Épico E02

#### Ticket `TICK-02.1` — Implementar Etapa 3: 3 fichas por personagem, recorte rembg+despill e sprites/*.json sem sobreposicao
* **Iniciativa:** `D02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/elenco.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/elenco_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/elenco.go, internal/diorama/ativos/elenco_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestElenco`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/elenco.go') && file_exists('internal/diorama/ativos/elenco_test.go')`
* **Spec Ref:** `AC-E02.1, AC-E02.2`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Épico E03

#### Ticket `TICK-03.1` — Implementar Etapa 4: cenarios SVG autorais (codigo), fila Flow resiliente com backoff/--ref/trava e idempotencia por hash
* **Iniciativa:** `D02` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/cenarios.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/cenarios_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/cenarios.go, internal/diorama/ativos/cenarios_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestCenarios`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/cenarios.go') && file_exists('internal/diorama/ativos/cenarios_test.go')`
* **Spec Ref:** `AC-E03.1, AC-E03.2, AC-E03.3`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Épico E04

#### Ticket `TICK-04.1` — Implementar Etapa 5: compilador e validador de cenas.json (100% linhas, metafora visual e >= 4 elementos moveis)
* **Iniciativa:** `D02` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","TICK-03.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/storyboard.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/storyboard_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/storyboard.go, internal/diorama/ativos/storyboard_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/ -run TestStoryboard`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/storyboard.go') && file_exists('internal/diorama/ativos/storyboard_test.go')`
* **Spec Ref:** `AC-E04.1, AC-E04.2`

#### Ticket `TICK-04.2` — Orquestrar Processo completo com politica_gates (humano/mecanico/hibrido) e cobertura >= 95%
* **Iniciativa:** `D02` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2","TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/ativos/processo.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/ativos/processo_test.go`
* **Outputs:** Artifacts `internal/diorama/ativos/processo.go, internal/diorama/ativos/processo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -cover ./internal/diorama/ativos/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/ativos/processo.go') && file_exists('internal/diorama/ativos/processo_test.go')`
* **Spec Ref:** `AC-E04.3, AC-NF.1, AC-NF.2`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Épico E05

#### Ticket `TICK-05.1` — Integrar CLI legoloop diorama-ativos, DAG v2, Skill e fechamento 5 de 5 via legoloop fechar
* **Iniciativa:** `D02` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/diorama_ativos.go, internal/cli/cli.go, processos/legoloop-diorama-ativos/pipeline.dag.json, skills/legoloop-diorama-ativos/SKILL.md` · Spec Anchor `docs/specs/2026-11-24_SPEC_LEGOLOOP_DIORAMA_ATIVOS.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/diorama_ativos_test.go`
* **Outputs:** Artifacts `internal/cli/diorama_ativos.go, internal/cli/diorama_ativos_test.go, internal/cli/cli.go, processos/legoloop-diorama-ativos/pipeline.dag.json, skills/legoloop-diorama-ativos/SKILL.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/ativos/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/diorama_ativos.go') && file_exists('internal/cli/diorama_ativos_test.go') && file_exists('internal/cli/cli.go') && file_exists('processos/legoloop-diorama-ativos/pipeline.dag.json') && file_exists('skills/legoloop-diorama-ativos/SKILL.md')`
* **Spec Ref:** `AC-E05.1`

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
