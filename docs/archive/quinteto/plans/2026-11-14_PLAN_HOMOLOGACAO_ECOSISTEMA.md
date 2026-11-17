---
id: PLAN-HOMOLOGACAO-ECOSISTEMA
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [SPEC-HOMOLOGACAO-ECOSISTEMA]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E02.1, AC-E02.2, AC-E03.1, AC-E03.2, AC-E04.1, AC-E04.2, AC-E05.1, AC-E05.2, AC-E06.1, AC-E06.2, AC-E07.1, AC-E07.2, AC-NF.1]
---

# 🗺️ PLAN: Correção Integral e Homologação do Ecossistema LegoLoop

> **Status:** Em Planejamento (Aguardando Stop-Gate Único Humano)  
> **Data:** 2026-11-14  
> **Versão:** 1.0.0  
> **Autor:** Antigravity Coordinator  
> **Iniciativas:** I-A (Arquitetura, Go, Segurança e Taxonomia) · I-B (Canais CLI, MCP Router, Motores N10 e Swarm) · I-C (Studio Web, Swarm UI e Fechamento)  
> **Referência da Intenção:** [docs/intents/2026-11-14_INTENT_HOMOLOGACAO_ECOSISTEMA.md](../intents/2026-11-14_INTENT_HOMOLOGACAO_ECOSISTEMA.md)  
> **Referência da Especificação:** [docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md](../specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-14_ROUTE_HOMOLOGACAO_ECOSISTEMA.json](../routes/2026-11-14_ROUTE_HOMOLOGACAO_ECOSISTEMA.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** Swarm Nativo (`workspace="share"` para autores por nível topológico e `workspace="inherit"` para verificadores cegos)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (workers dedicados no Swarm com contexto isolado por Ticket e escrita restrita às Peças atribuídas)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `4`
* **Timeout Padrão por Ticket:** `300s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Iniciativas, Épicos & Rastreabilidade

| Iniciativa | Épico ID | Nome do Épico | Tickets Micro / Nano | Critérios Cobertos | Peças Vinculadas |
|---|---|---|---|---|---|
| **I-A** | **E01** | Saneamento de Suíte, Schemas e Taxonomia | `TICK-01.1..01.3` (8 Nanos) | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `PECA-01`, `PECA-06` |
| **I-A** | **E02** | Segurança Path Jail, Concorrência e Mídia Real | `TICK-02.1..02.2` (5 Nanos) | `AC-E02.1`, `AC-E02.2`, `AC-NF.1` | `PECA-02`, `PECA-03`, `PECA-05` |
| **I-B** | **E03** | CLI, Harnesses, Skills, MCP Router e Fechar | `TICK-03.1..03.2` (5 Nanos) | `AC-E03.1`, `AC-E03.2` | `PECA-04`, `PECA-02` |
| **I-B** | **E04** | Integração dos 9 Motores N10 e Swarm Único | `TICK-04.1..04.2` (4 Nanos) | `AC-E04.1`, `AC-E04.2` | `PECA-02`, `PECA-04` |
| **I-C** | **E05** | Segurança Swarm UI e Jornada E2E Studio | `TICK-05.1..05.2` (3 Nanos) | `AC-E05.1`, `AC-E05.2`, `AC-NF.1` | `PECA-05` |
| **I-C** | **E06** | Design System, Responsividade, WCAG AA e ARIA | `TICK-06.1..06.2` (2 Nanos) | `AC-E06.1`, `AC-E06.2` | `PECA-05` |
| **I-C** | **E07** | Verificação Cega (Etapa 5) & Fechamento (Etapa 6) | `TICK-07.1..07.2` (3 Nanos) | `AC-E07.1`, `AC-E07.2` | `PECA-06`, `PECA-07` |

| R | Requisito | AC | Micro-Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Suíte Go, DAGs e remoção de segredos trackeados | `AC-E01.1` | `TICK-01.1` |
| `R2` | Paridade Schema v2, Parser Go, Gatilho e `ROUTE` | `AC-E01.2` | `TICK-01.2` |
| `R3` | Taxonomia `taxlinter` zerada e combos migrados | `AC-E01.3` | `TICK-01.3` |
| `R4` | Path Jail, CSRF/CORS, SSE sem race e `runner` | `AC-E02.1` | `TICK-02.1` |
| `R5` | Mídia real FFmpeg/Playwright em `motion` e `videorea` | `AC-E02.2` | `TICK-02.2` |
| `R6` | `mcp-router` operacional e executor de nó `mcp` | `AC-E03.1` | `TICK-03.1` |
| `R7` | Paridade CLI, harnesses/skills e `fechar` real | `AC-E03.2` | `TICK-03.2` |
| `R8` | Integração dos 9 motores N10 e resiliência | `AC-E04.1` | `TICK-04.1` |
| `R9` | Swarm Board único, file-lock e spawn/kill real | `AC-E04.2` | `TICK-04.2` |
| `R10` | Swarm UI unificado, zero XSS e gestão de workers | `AC-E05.1` | `TICK-05.1` |
| `R11` | Jornada E2E Studio, `#btn-run`, `ErrSemMatriz` e combos | `AC-E05.2` | `TICK-05.2` |
| `R12` | Layout responsivo, inspetor recuperável e tema | `AC-E06.1` | `TICK-06.1` |
| `R13` | Contraste WCAG AA, guildas legíveis e ARIA | `AC-E06.2` | `TICK-06.2` |
| `R14` | Verificação cega `>= 9.5/10`, NF e Fechamento Lei #11 | `AC-E07.1`, `AC-E07.2`, `AC-NF.1` | `TICK-07.1`, `TICK-07.2` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Iniciativa I-A · Épico 01 — Saneamento de Suíte, Schemas e Taxonomia

#### Micro-Ticket `TICK-01.1` — Desbloqueio de Suíte, DAGs e Remoção de Segredos do Git
* **Iniciativa:** `I-A` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-01`, `PECA-02`, `PECA-06` (`route2dag`, `isolated.go`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `docs/dags/2026-11-14_DAG_AUDITORIA_E2E_SWARM.json`, `internal/runner/isolated.go`, `RELATORIO_AUDITORIA_GEMINI4.md`, `AUDITORIA_FONTE.txt`, `legoloop.exe` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/runner/isolated_test.go:TestIsolatedExcludesCredenciaisAndSymlinks`
* **Outputs:** `docs/dags/2026-11-14_DAG_AUDITORIA_E2E_SWARM.json`, `internal/runner/isolated.go`, `internal/runner/isolated_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./tools/route2dag/... ./internal/runner/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/dags/2026-11-14_DAG_AUDITORIA_E2E_SWARM.json') && file_exists('internal/runner/isolated.go')`
* **Spec Ref:** `AC-E01.1`
* **Nano-Tickets:**
  * `NANO-01.1.1`: Corrigir `docs/dags/2026-11-14_DAG_AUDITORIA_E2E_SWARM.json` para contrato v2 válido (`output.artifact`, campos obrigatórios) restaurando `route2dag_test.go`.
  * `NANO-01.1.2`: Remover `RELATORIO_AUDITORIA_GEMINI4.md`, `AUDITORIA_FONTE.txt` e binário obsoleto `legoloop.exe` da raiz; blindar `internal/runner/isolated.go` bloqueando `.env*`, `credential_*.json`, erros silenciados em `copyPath` e symlinks.

#### Micro-Ticket `TICK-01.2` — Paridade 1:1 Schema v2 ↔ Parser Go, Peça Gatilho e Pares `ROUTE` ↔ `DAG`
* **Iniciativa:** `I-A` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-01` (`schemas/dag.v2.json`, `internal/dag`, `internal/graph`)
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `24000` tokens · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** `schemas/dag.v2.json`, `internal/dag/types.go`, `internal/dag/parse.go`, `internal/dag/schema_v2.go`, `internal/graph/ontology.go`, `internal/graph/index.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E01.2`
* **Alvo de Teste TDD:** `internal/dag/parity_v2_test.go:TestSchemaV2AndTriggerParity`
* **Outputs:** `schemas/dag.v2.json`, `internal/dag/types.go`, `internal/dag/parse.go`, `internal/dag/schema_v2.go`, `internal/dag/parity_v2_test.go`, `internal/graph/ontology.go`, `internal/graph/index.go`, `docs/dags/2026-11-10_DAG_DISSECAR.json`, `docs/dags/2026-11-13_DAG_STUDIO_PAPEIS_E_COMBOS.json`, `docs/dags/2026-11-14_DAG_DESIGN_UI.json`, `docs/dags/2026-11-14_DAG_SWARM.json`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/dag/... ./internal/graph/... ./tools/route2dag/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/parity_v2_test.go')`
* **Spec Ref:** `AC-E01.2`
* **Nano-Tickets:**
  * `NANO-01.2.1`: Alinhar `schemas/dag.v2.json` com `internal/dag` (`Role`, `Papel`, `Via`, `Combo`, `TokenBudget`, `TestSeam`, nó `mcp` permitido, `sensitiveDirs` com `.env*` e `credential_gemini_backup.json`).
  * `NANO-01.2.2`: Unificar contrato da Peça Gatilho (identificador único de fim de Processo, portas de entrada e saída em `ontology.go` e `index.go`).
  * `NANO-01.2.3`: Alinhar os 4 pares `ROUTE` ↔ `DAG` divergentes (`2026-11-10_DAG_DISSECAR`, `2026-11-13_DAG_STUDIO_PAPEIS_E_COMBOS`, `2026-11-14_DAG_DESIGN_UI`, `2026-11-14_DAG_SWARM`).

#### Micro-Ticket `TICK-01.3` — Taxonomia ASD-STE100 Zero Violações e Migração de Combos Honestos
* **Iniciativa:** `I-A` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-06` (`tools/taxlinter`, `docs/templates`, `internal/studiodoc`, `internal/router`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `20000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `cmd/legoloop/main.go`, `cmd/legoloop/migrar.go`, `cmd/legoloop/triagem.go`, `cmd/legoloop/triagem_test.go`, `docs/templates/`, `internal/studiodoc/store.go`, `internal/router/`, acervo de Processos, `skills/` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E01.3`
* **Alvo de Teste TDD:** `docs/templates/templates_test.go:TestTaxonomiaTemplatesZeroViolacoes`
* **Outputs:** `cmd/legoloop/main.go`, `docs/templates/templates_test.go`, `internal/studiodoc/store.go`, `internal/router/router.go`
* **Verificação (argv, sem shell):** `go run ./tools/taxlinter` nos alvos oficiais do AGENTS.md
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/templates/templates_test.go')`
* **Spec Ref:** `AC-E01.3`
* **Nano-Tickets:**
  * `NANO-01.3.1`: Zerar as 25 violações de `taxlinter` em `cmd/legoloop/` e mover `triagem` e migração de Processos para despacho consistente.
  * `NANO-01.3.2`: Zerar as 33 violações em `docs/templates/`, definir `limiteTaxonomia = 0` em `templates_test.go` e remover burla de string em `internal/studiodoc/store.go`.
  * `NANO-01.3.3`: Migrar as 171 ocorrências de combos legados no acervo de Processos, `skills/`, `docs/templates/` e `internal/router/` para os 6 combos certificados mantendo alias em `internal/combo/legado.go`.

---

### Iniciativa I-A · Épico 02 — Segurança Path Jail, Concorrência e Remoção de Stubs de Mídia

#### Micro-Ticket `TICK-02.1` — Blindagem Path Jail, CORS/CSRF, Data Race SSE e Ordem do Runner
* **Iniciativa:** `I-A` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-02`, `PECA-05` (`internal/studio`, `internal/runner`)
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `24000` tokens · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** `internal/studio/reveal.go`, `internal/studio/piece_actions.go`, `internal/studio/security.go`, `internal/studio/run_api.go`, `internal/studio/server.go`, `internal/studio/processos_api.go`, `internal/runner/runner.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/studio/security_hardening_test.go:TestPathJailDeleteOriginAndSSERace`
* **Outputs:** `internal/studio/reveal.go`, `internal/studio/piece_actions.go`, `internal/studio/security.go`, `internal/studio/run_api.go`, `internal/studio/server.go`, `internal/studio/security_hardening_test.go`, `internal/runner/runner.go`, `internal/runner/runner_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/... ./internal/runner/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/security_hardening_test.go')`
* **Spec Ref:** `AC-E02.1`
* **Nano-Tickets:**
  * `NANO-02.1.1`: Blindar `reveal.go`, `piece_actions.go` (mutex em `toggleFavorite`, Path Jail em `eliminatePiece`) e `security.go` (`Origin` em `POST`/`PUT`/`PATCH`/`DELETE`), removendo retorno de caminhos absolutos.
  * `NANO-02.1.2`: Corrigir data race em `run_api.go` (`sync.Mutex` no `newSSESender` e `ResponseController` para streams sem corte de `90s` em `server.go`), mascarar caminhos e trocar mock `"idle"` em `handleDAGStatus` por estado real.
  * `NANO-02.1.3`: Corrigir `runner.go` para executar `evalGate` na worktree isolada antes de `promoteArtifact` (com `os.Root`/`filepath.IsLocal`) e remover o stub `unsupported` retornando `nil`.

#### Micro-Ticket `TICK-02.2` — Substituição Integral de Stubs de Mídia (`motion` e `videorea`)
* **Iniciativa:** `I-A` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-03` (`internal/motion`, `internal/videorea`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/motion/encode.go`, `internal/motion/render.go`, `internal/videorea/cartoon_render.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E02.2`
* **Alvo de Teste TDD:** `internal/motion/motion_real_test.go:TestEncoderERendererSemStubPlaceholder`
* **Outputs:** `internal/motion/encode.go`, `internal/motion/render.go`, `internal/motion/motion_real_test.go`, `internal/videorea/cartoon_render.go`, `internal/videorea/cartoon_render_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/motion/... ./internal/videorea/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/motion_real_test.go')`
* **Spec Ref:** `AC-E02.2`
* **Nano-Tickets:**
  * `NANO-02.2.1`: Substituir `"mp4-placeholder"` e PNG 1x1 branco em `internal/motion/` por invocação real de FFmpeg/Playwright com `exec.LookPath` retornando erro explícito quando ausente.
  * `NANO-02.2.2`: Corrigir `internal/videorea/cartoon_render.go` para renderizar `script.Scenes` no FFmpeg sem fallback MP4 falso de 32 bytes.

---

### Iniciativa I-B · Épico 03 — CLI, Harnesses, Skills, MCP Router e `fechar` Real

#### Micro-Ticket `TICK-03.1` — MCP Router Operacional & Executor Real do Nó `mcp`
* **Iniciativa:** `I-B` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-04`, `PECA-02` (`internal/mcp`, `internal/cli/mcp_router.go`, `internal/orchestrix`)
* **Dependências:** `["TICK-01.2", "TICK-04.2"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `24000` tokens · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/mcp_router.go`, `internal/mcp/router.go`, `internal/mcp/swarm_tools.go`, `internal/orchestrix/orchestrix.go`, `internal/runner/runner.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E03.1`
* **Alvo de Teste TDD:** `internal/mcp/router_stdio_test.go:TestMCPRouterCofreSwarmEErro32700`
* **Outputs:** `internal/cli/mcp_router.go`, `internal/mcp/router.go`, `internal/mcp/swarm_tools.go`, `internal/mcp/router_stdio_test.go`, `internal/orchestrix/mcp_client.go`, `internal/orchestrix/mcp_client_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/mcp/... ./internal/orchestrix/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/orchestrix/mcp_client.go')`
* **Spec Ref:** `AC-E03.1`
* **Nano-Tickets:**
  * `NANO-03.1.1`: Conectar `cli/mcp_router.go` ao cofre/skills reais (`SetSkills`), incluir `swarm_*` no `tools/list`, responder JSON-RPC `-32700`, suportar linhas > 64 KiB via `bufio.Reader` e tornar `swarm_status` somente leitura.
  * `NANO-03.1.2`: Criar `internal/orchestrix/mcp_client.go` (cliente STDIO JSON-RPC `initialize` + `tools/call`) e conectar execução real do nó tipo `mcp` no `runner`.

#### Micro-Ticket `TICK-03.2` — Paridade CLI, Harnesses/Skills e `fechar` 100% Real
* **Iniciativa:** `I-B` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-04`, `PECA-07` (`internal/cli`, `internal/fechar`, `skills/embed.go`, `internal/processo`, `internal/assets`)
* **Dependências:** `["TICK-01.3", "TICK-02.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `360s` · Retries `1` · Entrada `28000` tokens · Teto de saída `900` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/`, `internal/fechar/`, `skills/embed.go`, `internal/processo/vinculo.go`, `internal/assets/scanner.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E03.2`
* **Alvo de Teste TDD:** `internal/fechar/fechar_real_test.go:TestFecharSemStubsETelemetriaReal`
* **Outputs:** `internal/cli/cli.go`, `internal/cli/remover.go`, `internal/cli/stept.go`, `internal/cli/sync_harnesses.go`, `internal/cli/fechar.go`, `internal/fechar/fechar.go`, `internal/fechar/types.go`, `internal/fechar/fechar_real_test.go`, `skills/embed.go`, `internal/processo/vinculo.go`, `internal/assets/scanner.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/cli/... ./internal/fechar/... ./internal/processo/... ./internal/assets/... ./skills/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fechar/fechar_real_test.go')`
* **Spec Ref:** `AC-E03.2`
* **Nano-Tickets:**
  * `NANO-03.2.1`: Corrigir `usage` vs handlers (`cli.go`), exit `0` em `--help`, `--dry-run`/`--force`/`--json` em `remover`, `--aspect` em `video-motion`, JSON puro em `canonize`, leitura pura em `stept status` e timeout de `cli_test.go:383`.
  * `NANO-03.2.2`: Incluir todas as skills nativas em `skills/embed.go`, adicionar `~/.gemini/config/skills` em `vinculo.go`, corrigir `.jsonc` e `mcp.rea` em `sync_harnesses.go` e ignorar `.trash`/`.vault_terceiros` em `assets/scanner.go`.
  * `NANO-03.2.3`: Remover stubs e tokens fictícios (`80000`/`65000`) de `internal/fechar/` e `cli/fechar.go`, ligando execução real de rebuild, governança, reindex e portões e corrigindo `duracao_ms` para milissegundos reais.

---

### Iniciativa I-B · Épico 04 — Integração dos 9 Motores N10 e Unificação Swarm Board/Runtime

#### Micro-Ticket `TICK-04.1` — Ligação Completa dos 9 Motores N10 em Produção e Resiliência Orchestrix
* **Iniciativa:** `I-B` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-02` (`internal/runner`, `internal/orchestrix`, `internal/brain`, `internal/combo`, `internal/gate`, `internal/dag`)
* **Dependências:** `["TICK-02.1", "TICK-03.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `360s` · Retries `1` · Entrada `28000` tokens · Teto de saída `900` tokens · Rollback `stash_changes`
* **Inputs:** `internal/runner/`, `internal/orchestrix/`, `internal/brain/`, `internal/combo/`, `internal/cli/commands.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E04.1`
* **Alvo de Teste TDD:** `internal/runner/n10_integration_test.go:TestN10EnginesConnectedInProduction`
* **Outputs:** `internal/runner/runner.go`, `internal/runner/n10_integration_test.go`, `internal/orchestrix/orchestrix.go`, `internal/orchestrix/backends.go`, `internal/orchestrix/unit_harness.go`, `internal/cli/commands.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/runner/... ./internal/orchestrix/... ./internal/brain/... ./internal/combo/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/n10_integration_test.go')`
* **Spec Ref:** `AC-E04.1`
* **Nano-Tickets:**
  * `NANO-04.1.1`: Integrar `ExecuteWithAndon`, `FormatSBAR`, `FlightRecord.SaveFlightRecord`, `ValidateRedBlueBoundary` e `MergeAST` no ciclo real de execução do `internal/runner`.
  * `NANO-04.1.2`: Integrar `AntibodiesDB`, `ExecuteSpeculative` (failover/hot-swap em `ModeAuto`/`ModeSuggested`, tratamento HTTP 413/429, captura `stderr`, timeout HTTP em `cli/commands.go`), `ValidarRequisicao`, `ValidateProvenance` e `ValidateSTE100`; propagar flags/Path Jail no executor de subprocesso (`internal/orchestrix`) e `envGateway` em `unit_harness.go`.

#### Micro-Ticket `TICK-04.2` — Contrato Único de Swarm Board, File-Lock e Spawn/Kill Real
* **Iniciativa:** `I-B` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-04` (`internal/board`)
* **Dependências:** `[]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `20000` tokens · Teto de saída `700` tokens · Rollback `stash_changes`
* **Inputs:** `internal/board/swarm.go`, `internal/board/swarm_discovery.go`, `internal/board/swarm_runtime.go`, `internal/board/swarm_runtime_test.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E04.2`
* **Alvo de Teste TDD:** `internal/board/swarm_unified_test.go:TestUnifiedPathFileLockExactMentionAndPIDKill`
* **Outputs:** `internal/board/swarm.go`, `internal/board/swarm_discovery.go`, `internal/board/swarm_runtime.go`, `internal/board/swarm_runtime_test.go`, `internal/board/swarm_unified_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/board/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/board/swarm_unified_test.go')`
* **Spec Ref:** `AC-E04.2`
* **Nano-Tickets:**
  * `NANO-04.2.1`: Padronizar escrita em `.legoloop-swarm/swarm.json` com leitura retrocompatível (`.swarm/swarm.json` e `swarm_board.json`), file-lock cross-processo, temp único (PID+ns), sanitização `req.Name`, distinção `share` vs `branch`, `@mention` exato e fim do leak em `swarm_runtime_test.go`.
  * `NANO-04.2.2`: Ligar `SpawnWorker`, `SetWorkerModel` e `KillWorker` (encerrando PID real no SO e removendo worktree) ao board unificado e CLI.

---

### Iniciativa I-C · Épico 05 — Segurança Swarm UI e Jornada E2E Studio Web

#### Micro-Ticket `TICK-05.1` — Swarm UI Unificado, Zero XSS, Preservação de Estado e Gestão de Workers
* **Iniciativa:** `I-C` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-05` (`internal/cli/swarm_ui.go`, `internal/studio/swarm_api.go`, `internal/studio/static/swarm.js`)
* **Dependências:** `["TICK-02.1", "TICK-04.2"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `24000` tokens · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/swarm_ui.go`, `internal/studio/swarm_api.go`, `internal/studio/static/swarm.js`, `internal/studio/static/swarm.html` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E05.1`
* **Alvo de Teste TDD:** `internal/studio/swarm_security_test.go:TestSwarmUIUnifiedNoXSSAndRealKill`
* **Outputs:** `internal/cli/swarm_ui.go`, `internal/studio/swarm_api.go`, `internal/studio/static/swarm.js`, `internal/studio/static/swarm.html`, `internal/studio/swarm_security_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/swarm_security_test.go')`
* **Spec Ref:** `AC-E05.1`
* **Nano-Tickets:**
  * `NANO-05.1.1`: Remover HTML duplicado de `cli/swarm_ui.go` delegando para `internal/studio` com `secure()`, erradicar `innerHTML` inseguro em `swarm.js` via `textContent`, preservar `#cmd-in` e `selectedAgent` por aba, reabrir abas fechadas, corrigir `s.ProjectDir` no `KillWorker` e incluir controles `+ Worker` e troca de modelo na UI.

#### Micro-Ticket `TICK-05.2` — Desbloqueio E2E do Studio (`#btn-run`, `ErrSemMatriz`, Gatilho, SVG e Combos)
* **Iniciativa:** `I-C` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-05` (`internal/studio/`)
* **Dependências:** `["TICK-01.2", "TICK-02.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `24000` tokens · Teto de saída `800` tokens · Rollback `stash_changes`
* **Inputs:** `internal/studio/static/app.js`, `internal/studio/static/inspector.js`, `internal/studio/static/lego-item-modal.js`, `internal/studio/fabricar_api.go`, `internal/studio/compose.go`, `internal/studio/templates.json` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E05.2`
* **Alvo de Teste TDD:** `internal/studio/studio_e2e_flow_test.go:TestFabricarErrSemMatrizAndComposeTrigger`
* **Outputs:** `internal/studio/static/app.js`, `internal/studio/static/inspector.js`, `internal/studio/static/lego-item-modal.js`, `internal/studio/fabricar_api.go`, `internal/studio/compose.go`, `internal/studio/templates.json`, `internal/studio/studio_e2e_flow_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/studio_e2e_flow_test.go')`
* **Spec Ref:** `AC-E05.2`
* **Nano-Tickets:**
  * `NANO-05.2.1`: Implementar `compileDAG` e jornada completa de `#btn-run` -> `#modal-run` SSE em `app.js`; mapear `fabricar.ErrSemMatriz` em `fabricar_api.go` (HTTP 422 com instrução clara).
  * `NANO-05.2.2`: Unificar identificador e portas do Gatilho (`app.js`, `inspector.js`, `compose.go`), alinhar `portPoint` SVG, atualizar `COMBO_HONESTO_SPECS` e detecção no modo Auto, corrigir `"Copiar JSON"` / `sonnet` hardcoded em `lego-item-modal.js` e acentos em `templates.json`.

---

### Iniciativa I-C · Épico 06 — Design System, Responsividade, Contraste WCAG AA e A11y

#### Micro-Ticket `TICK-06.1` — Layout Responsivo, Inspetor Recuperável, Modais e Tema
* **Iniciativa:** `I-C` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-05` (`internal/studio/static/{index.html,inventario.html,swarm.html,style.css,ui.js}`)
* **Dependências:** `["TICK-05.1", "TICK-05.2"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `20000` tokens · Teto de saída `700` tokens · Rollback `stash_changes`
* **Inputs:** `internal/studio/static/index.html`, `internal/studio/static/inventario.html`, `internal/studio/static/swarm.html`, `internal/studio/static/style.css`, `internal/studio/static/ui.js` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E06.1`
* **Alvo de Teste TDD:** `internal/studio/static_test.go:TestLayoutModaisNavegacaoETema`
* **Outputs:** `internal/studio/static/index.html`, `internal/studio/static/inventario.html`, `internal/studio/static/swarm.html`, `internal/studio/static/style.css`, `internal/studio/static/ui.js`, `internal/studio/static_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/static/style.css')`
* **Spec Ref:** `AC-E06.1`
* **Nano-Tickets:**
  * `NANO-06.1.1`: Corrigir layout/clipping em `1440`/`1024`/`390px` sem altura fixa de topbar, tornar `#btn-toggle-inspector` sempre visível/clicável fora de `#canvas-telemetry-bar`, corrigir `.modal-header` nos 4 modais, incluir link `Swarm` em `inventario.html` e criar seletor de tema Claro/Escuro persistido em armazenamento local do navegador.

#### Micro-Ticket `TICK-06.2` — Contraste WCAG AA, Guildas Legíveis e Acessibilidade ARIA
* **Iniciativa:** `I-C` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-05` (`internal/studio/static/{style.css,palette.js,swarm.html,swarm.js}`)
* **Dependências:** `["TICK-06.1"]`
* **Isolamento:** Workspace `share` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Modelo `t-agente-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `20000` tokens · Teto de saída `700` tokens · Rollback `stash_changes`
* **Inputs:** `internal/studio/static/style.css`, `internal/studio/static/palette.js`, `internal/studio/static/swarm.js`, `internal/studio/static_test.go` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E06.2`
* **Alvo de Teste TDD:** `internal/studio/static_test.go:TestWCAGTokensEGuildas`
* **Outputs:** `internal/studio/static/style.css`, `internal/studio/static/palette.js`, `internal/studio/static/swarm.js`, `internal/studio/static_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/studio/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/studio/static/palette.js')`
* **Spec Ref:** `AC-E06.2`
* **Nano-Tickets:**
  * `NANO-06.2.1`: Ajustar contraste WCAG AA (`>= 4.5:1`) em Claro e Escuro nas rotas `/`, `/inventario` e `/swarm`, aumentar `.guilda-badge` para `>= 11px` legível, definir tokens faltantes em `style.css`, atribuir guilda em `entryFromBlock` (`palette.js`) e corrigir semântica HTML/ARIA/teclado.

---

### Iniciativa I-C · Épico 07 — Verificação Cega Independente & Fechamento Canônico Universal

#### Micro-Ticket `TICK-07.1` — Verificação Cega Independente (Red-Team Gate - Etapa 5)
* **Iniciativa:** `I-C` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `script` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-06`, `PECA-07` (`playwright`, `chrome-devtools-mcp`, `rea`, `obscura`, `taxlinter`, `funclen`)
* **Dependências:** `["TICK-02.2", "TICK-03.2", "TICK-04.1", "TICK-06.2"]`
* **Isolamento:** Workspace `inherit` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-agente-128k` · Effort `high` · Fallback Modelo `t-codigo-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `true` · Allowlist `["playwright:browser_navigate", "chrome-devtools-mcp:list_pages", "rea:inspect_web_page"]`
* **Limites & Resiliência:** Timeout `600s` · Retries `1` · Entrada `32000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** `bin/legoloop.exe`, `cmd/`, `internal/`, `tools/`, `schemas/`, acervo de Processos, `skills/` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E07.1`
* **Alvo de Teste TDD:** `internal/e2e/e2e_test.go:TestE2EFullBinary`
* **Outputs:** `bin/legoloop.exe`, `scratch/gates/TICK-07.1.json`
* **Verificação (argv, sem shell):** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`
* **Spec Ref:** `AC-E07.1`
* **Nano-Tickets:**
  * `NANO-07.1.1`: Recompilar `bin/legoloop.exe` e executar `go test -count=1 ./...`, `-tags e2e`, `go vet ./...`, `gofmt -l .`, `funclen`, `taxlinter == 0` e `git diff --exit-code go.mod go.sum` por auditores cegos.
  * `NANO-07.1.2`: Validar todos os canais (CLI `validate`, `mcp-router` STDIO, `instalar-skills`, Studio + Swarm UI com Playwright, Chrome DevTools, REA e Obscura) comprovando notas `>= 9.5/10` nas 3 frentes; reabrir automaticamente no Swarm se alguma frente ficar `< 9.5/10`.

#### Micro-Ticket `TICK-07.2` — Fechamento Canônico Universal (Etapa 6 - Lei #11)
* **Iniciativa:** `I-C` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Peça Inventariada:** `PECA-07` (`legoloop fechar`, `legoloop legopress`, `legoloop instalar-skills`, `codebase-memory-mcp`)
* **Dependências:** `["TICK-07.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:** Modo `suggested` · Harness `9router` · Via `9router` · Modelo `t-agente-128k` · Effort `medium` · Fallback Modelo `t-codigo-128k`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `true` · Allowlist `["codebase-memory-mcp:index_repository"]`
* **Limites & Resiliência:** Timeout `600s` · Retries `1` · Entrada `32000` tokens · Teto de saída `1500` tokens · Rollback `stash_changes`
* **Inputs:** `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md`, `docs/dags/2026-11-14_DAG_HOMOLOGACAO_ECOSISTEMA.json` · Spec Anchor `docs/specs/2026-11-14_SPEC_HOMOLOGACAO_ECOSISTEMA.md#AC-E07.2`
* **Alvo de Teste TDD:** `internal/fechar/fechar_real_test.go:TestFecharHomologacaoIntegral`
* **Outputs:** `docs/manuals/2026-11-14_MANUAL_FUNCAO_HOMOLOGACAO_ECOSISTEMA_LEGOLOOP.md`, `docs/books/homologacao-ecosistema/livro.json`, `skills/legoloop-homologacao/SKILL.md`, `docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md`
* **Verificação (argv, sem shell):** `go install ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md') && file_exists('skills/legoloop-homologacao/SKILL.md')`
* **Spec Ref:** `AC-E07.2`
* **Nano-Tickets:**
  * `NANO-07.2.1`: Gerar Manual Operacional, compilar/auditar Livro LegoPress, criar Skill canônica `legoloop-homologacao`, sincronizar harnesses via `instalar-skills`, instalar `legoloop` no PATH, atualizar `AGENTS.md`/`MEMORY.md`/`README.md`, reindexar `codebase-memory-mcp` e emitir Relatório Final com comparativo Antes vs Depois.

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir os épicos, executar em ordem determinística (sem shell):
1. Suíte de regressão global: `go test -count=1 ./...` e `go test -count=1 -tags e2e ./internal/e2e/...`
2. Estática, formatação e taxonomia: `go vet ./...`, `gofmt -l .` (saída vazia), `go run ./tools/funclen internal cmd tools` e `go run ./tools/taxlinter` nos alvos oficiais do AGENTS.md
3. Recompilação e instalação: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`
4. Verificação de dependências intactas: `git diff --exit-code go.mod go.sum`
5. Emissão do relatório formal: `docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe') && file_exists('docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md')`

---

🛑 **STOP-GATE ÚNICO HUMANO:**  
*O planejamento hierarquizado (Iniciativas -> Épicos -> Micro-tickets -> Nano-tickets), a entrevista socrática e o Quinteto ISO completo estão aprovados para execução autônoma via Swarm nas Etapas 4, 5 e 6?*
