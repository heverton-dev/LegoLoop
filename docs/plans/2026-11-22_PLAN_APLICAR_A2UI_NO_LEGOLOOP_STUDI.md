---
id: PLAN-2026-11-22_ROUTE_APLICAR_A2UI_NO_LEGOLOOP_STUDI
status: DRAFT
date: 2026-11-22
schema_version: 2
upstream: [2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E02.1, AC-E02.2, AC-E02.3, AC-E03.1, AC-E03.2, AC-E03.3, AC-E04.1, AC-E04.2, AC-E04.3, AC-E05.1, AC-E05.2, AC-E06.1, AC-E06.2, AC-NF.1, AC-NF.2, AC-NF.3]
---

# 🗺️ PLAN: Motor Único de Interface Declarativa A2UI nas 3 Vias do LegoLoop

> **Status:** Em Planejamento  
> **Data:** 2026-11-22  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-22_INTENT_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md](../intents/2026-11-22_INTENT_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md)  
> **Referência da Especificação:** [docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md](../specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-22_ROUTE_APLICAR_A2UI_NO_LEGOLOOP_STUDI.json](../routes/2026-11-22_ROUTE_APLICAR_A2UI_NO_LEGOLOOP_STUDI.json)  
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
| **E01** | Núcleo de Sessão Interativa & Validador A2UI (`internal/a2ui`) | 3 | `AC-E01.2`, `AC-E03.3`, `AC-E01.1`, `AC-E01.3`, `AC-NF.3`, `AC-E01.1`, `AC-E02.3`, `AC-NF.1` | `micro` |
| **E02** | Renderizador `@a2ui/lit`, Design System LegoLoop & Catálogo de Domínio (`internal/studio/static/`) | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Macroprocessos, Isolamento por `surfaceId` e Árvore/Breadcrumb | 1 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Integração nas 3 Vias (Terminal CLI, Agente/Slash via MCP e Studio Canvas) | 3 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` | `micro` |
| **E05** | Projeto Piloto: `legoloop-concepcao` nas 3 Vias | 1 | `AC-E05.1`, `AC-E05.2` | `micro` |
| **E06** | Épico E06 | 3 | `AC-E06.1`, `AC-E06.2`, `AC-NF.2`, `AC-NF.3` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Núcleo de Sessão Interativa & Validador A2UI (`internal/a2ui`)

#### Ticket `TICK-01.1` — Criar contratos, validador de catálogo confiável e bloqueio de payload hostil em internal/a2ui
* **Iniciativa:** `A2UI` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/envelope.go, internal/a2ui/catalog.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/envelope_test.go`
* **Outputs:** Artifacts `internal/a2ui/envelope.go, internal/a2ui/catalog.go, internal/a2ui/envelope_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/a2ui/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/a2ui/envelope.go') && file_exists('internal/a2ui/catalog.go') && file_exists('internal/a2ui/envelope_test.go')`
* **Spec Ref:** `AC-E01.2, AC-E03.3`

#### Ticket `TICK-01.2` — Implementar Hub de sessões interativas com SSE, snapshot de reconexão e entrega de Action sem vazamento
* **Iniciativa:** `A2UI` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/hub.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/hub_test.go`
* **Outputs:** Artifacts `internal/a2ui/hub.go, internal/a2ui/hub_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/a2ui/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/a2ui/hub.go') && file_exists('internal/a2ui/hub_test.go')`
* **Spec Ref:** `AC-E01.1, AC-E01.3, AC-NF.3`

#### Ticket `TICK-01.3` — Expor rotas HTTP /api/sessao, /stream, /action e /artefato protegidas por s.secure e internal/jail no Studio
* **Iniciativa:** `A2UI` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/studio/a2ui_api.go, internal/studio/server.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/studio/a2ui_api_test.go`
* **Outputs:** Artifacts `internal/studio/a2ui_api.go, internal/studio/a2ui_api_test.go, internal/studio/server.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 ./internal/studio/... -run "TestA2UI|TestSecurity"')`
* **Spec Ref:** `AC-E01.1, AC-E02.3, AC-NF.1`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Renderizador `@a2ui/lit`, Design System LegoLoop & Catálogo de Domínio (`internal/studio/static/`)

#### Ticket `TICK-02.1` — Empacotar bundle offline de @a2ui/lit e @a2ui/web_core em internal/studio/static/components via go:embed
* **Iniciativa:** `A2UI` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/studio/static/sessao.html, internal/studio/static/components/a2ui-bundle.js, internal/studio/server.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/studio/static_test.go`
* **Outputs:** Artifacts `internal/studio/static/sessao.html, internal/studio/static/components/a2ui-bundle.js, internal/studio/server.go, internal/studio/static_test.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 ./internal/studio/... -run "TestStatic|TestEmojis"')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Implementar catálogo de domínio LegoLoop sobre @a2ui/lit com Design System do Studio
* **Iniciativa:** `A2UI` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/studio/static/components/legoloop-catalog.js, internal/studio/static/sessao.css` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/studio/ui_logic_test.go`
* **Outputs:** Artifacts `internal/studio/static/components/legoloop-catalog.js, internal/studio/static/sessao.css, internal/studio/ui_logic_test.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 ./internal/studio/... -run "TestUI|TestStatic"')`
* **Spec Ref:** `AC-E02.2`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Macroprocessos, Isolamento por `surfaceId` e Árvore/Breadcrumb

#### Ticket `TICK-03.1` — Implementar hierarquia de macroprocessos isolados por surfaceId com árvore, breadcrumb e foco automático em Stop-Gate
* **Iniciativa:** `A2UI` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/hierarchy.go, internal/runner/subprocesso.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/hierarchy_test.go`
* **Outputs:** Artifacts `internal/a2ui/hierarchy.go, internal/a2ui/hierarchy_test.go, internal/runner/subprocesso.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/a2ui/... ./internal/runner/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/a2ui/hierarchy.go') && file_exists('internal/a2ui/hierarchy_test.go') && file_exists('internal/runner/subprocesso.go')`
* **Spec Ref:** `AC-E03.1, AC-E03.2`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Integração nas 3 Vias (Terminal CLI, Agente/Slash via MCP e Studio Canvas)

#### Ticket `TICK-04.1` — Implementar AwaitAction (corrida stdin vs navegador) e integração da Via 1 (Terminal CLI)
* **Iniciativa:** `A2UI` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/await.go, internal/cli/exec.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/await_test.go`
* **Outputs:** Artifacts `internal/a2ui/await.go, internal/a2ui/await_test.go, internal/cli/exec.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/a2ui/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/a2ui/await.go') && file_exists('internal/a2ui/await_test.go') && file_exists('internal/cli/exec.go')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.2` — Adicionar ferramentas MCP sessao_criar, sessao_publicar e sessao_aguardar no legoloop mcp-router (Via 2)
* **Iniciativa:** `A2UI` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/mcp/sessao_tools.go, internal/mcp/router.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E04.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/mcp/sessao_tools_test.go`
* **Outputs:** Artifacts `internal/mcp/sessao_tools.go, internal/mcp/sessao_tools_test.go, internal/mcp/router.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/mcp/sessao_tools.go') && file_exists('internal/mcp/sessao_tools_test.go') && file_exists('internal/mcp/router.go')`
* **Spec Ref:** `AC-E04.2`

#### Ticket `TICK-04.3` — Conectar botão Executar do Studio Canvas à rota de sessão interativa A2UI preservando POST /api/run (Via 3)
* **Iniciativa:** `A2UI` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/studio/static/app.js, internal/studio/run_api.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/studio/run_api_test.go`
* **Outputs:** Artifacts `internal/studio/static/app.js, internal/studio/run_api.go, internal/studio/run_api_test.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 ./internal/studio/... -run "TestRun|TestStatic"')`
* **Spec Ref:** `AC-E04.3`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Projeto Piloto: `legoloop-concepcao` nas 3 Vias

#### Ticket `TICK-05.1` — Implementar adaptador A2UI do piloto legoloop-concepcao (confiança, perguntas, Espelho, Stop-Gate 01 e INTENT.md)
* **Iniciativa:** `A2UI` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","TICK-04.1","TICK-04.2","TICK-04.3","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/concepcao.go, internal/cli/fabricar_entrevista.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/concepcao_test.go`
* **Outputs:** Artifacts `internal/a2ui/concepcao.go, internal/a2ui/concepcao_test.go, internal/cli/fabricar_entrevista.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/a2ui/... ./internal/interview/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/a2ui/concepcao.go') && file_exists('internal/a2ui/concepcao_test.go') && file_exists('internal/cli/fabricar_entrevista.go')`
* **Spec Ref:** `AC-E05.1, AC-E05.2`

### 🛑 STOP-GATE HUMANO (`SG-E05`)
* **Gate (DSL):** `prompt('Evidências do Épico 05 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 06 — Épico E06

#### Ticket `TICK-06.1` — Fase 2: Implementar superfície LegoSwarmBoard alimentada por swarm_api.go com paridade funcional e segurança
* **Iniciativa:** `A2UI` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/a2ui/swarm_surface.go, internal/studio/static/components/legoloop-catalog.js, internal/studio/swarm_api.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E06.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/a2ui/swarm_surface_test.go`
* **Outputs:** Artifacts `internal/a2ui/swarm_surface.go, internal/a2ui/swarm_surface_test.go, internal/studio/static/components/legoloop-catalog.js, internal/studio/swarm_api.go`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 ./internal/studio/... ./internal/a2ui/... -run "TestSwarm|TestA2UI"')`
* **Spec Ref:** `AC-E06.1`

#### Ticket `TICK-06.2` — Fase 2: Remover UI web antiga do Swarm após paridade provada, mantendo motor/backend e redirecionamento
* **Iniciativa:** `A2UI` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/studio/server.go, internal/studio/swarm_api.go` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-E06.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/studio/swarm_crud_sim_test.go`
* **Outputs:** Artifacts `internal/studio/server.go, internal/studio/swarm_api.go, internal/studio/swarm_crud_sim_test.go, internal/studio/swarm_security_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/... -run TestSwarm`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/server.go') && file_exists('internal/studio/swarm_api.go') && file_exists('internal/studio/swarm_crud_sim_test.go') && file_exists('internal/studio/swarm_security_test.go')`
* **Spec Ref:** `AC-E06.2`

#### Ticket `TICK-06.3` — Portão final global: race, vet, funclen, gofmt, taxlinter, espinha verificar e recompilação fresca
* **Iniciativa:** `A2UI` | **Épico:** `E06` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.2","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `AGENTS.md, README.md, MEMORY.md` · Spec Anchor `docs/specs/2026-11-22_SPEC_APLICAR_A2UI_NO_LEGOLOOP_STUDI.md#AC-NF.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `AGENTS.md, README.md, MEMORY.md`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: go test -count=1 -race ./... && go vet ./... && go run ./tools/funclen internal && go run ./tools/taxlinter internal cmd loops skills docs/templates && go build -o bin/legoloop.exe ./cmd/legoloop && ./bin/legoloop.exe espinha verificar')`
* **Spec Ref:** `AC-NF.2, AC-NF.3`

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
*Plano de execução fatiado com 13 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
