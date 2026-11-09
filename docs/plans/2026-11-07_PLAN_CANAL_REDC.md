# 🗺️ PLAN: Mega-Pipeline REDC (Engenharia Reversa Profunda de Canais e Clonagem Metodológica/Estrutural)

> **Status:** Em Planejamento  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** LegoLoop  
> **Iniciativa:** I04 — Mega-Pipelines Audiovisuais e Automação de Canais  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_CANAL_REDC.md](../intents/2026-11-07_INTENT_CANAL_REDC.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_CANAL_REDC.md](../specs/2026-11-07_SPEC_CANAL_REDC.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_CANAL_REDC.json](../routes/2026-11-07_ROUTE_CANAL_REDC.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Ingestão, Catalogação e Clusterização dos 100 Vídeos | 2 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Motor Deep REDC Top 15 (DNA, Tokens e Ganchos) | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Gerador de Pipelines de Métodos e Brand Swap | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Integração CLI, Grafo DAG e Skill Canônica | 3 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Ingestão, Catalogação e Clusterização dos 100 Vídeos

#### Ticket `TICK-01.1` — Modelagem de dados, parser de métricas e ranking de viralidade
* **Iniciativa:** `I04` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` (backoff `1000ms`, fallback model `true`) · Orçamento Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/redc/catalog.go`, `internal/redc/types.go`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/redc/catalog_test.go:TestCatalogVideos`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/redc/types.go`, `internal/redc/catalog.go`, `internal/redc/catalog_test.go`
  * **Verificação (argv, sem shell):** `go test -v -count=1 ./internal/redc/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/catalog.go') && file_exists('internal/redc/catalog_test.go')`
* **Spec Ref:** `AC-E01.1`
* **Escopo:** Validar ordenação por tração (`views * 0.7 + engajamento * 0.3`) e rejeição de entradas corrompidas com `ErrDadosCorrompidos`.

#### Ticket `TICK-01.2` — Clusterização de métodos e desduplicação
* **Iniciativa:** `I04` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/redc/cluster.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E01.2`
* **Alvo de Teste TDD:** `internal/redc/cluster_test.go:TestClusterMethods`
* **Outputs:** `internal/redc/cluster.go`, `internal/redc/cluster_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestClusterMethods ./internal/redc/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/cluster.go')`
* **Spec Ref:** `AC-E01.2`
* **Escopo:** Agrupar métodos de 100 vídeos em clusters temáticos sem gerar processos duplicados para variações do mesmo título.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [x] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [x] Testes unitários de catalogação e clusterização verdes.

---

### Épico 02 — Motor Deep REDC Top 15 (DNA, Tokens e Ganchos)

#### Ticket `TICK-02.1` — Síntese do Dossiê DNA (tom, ritmo, design tokens, motion tokens e ganchos)
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/redc/dna.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/redc/dna_test.go:TestExtractDNA`
* **Outputs:** `internal/redc/dna.go`, `internal/redc/dna_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestExtractDNA ./internal/redc/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/dna.go')`
* **Spec Ref:** `AC-E02.1`
* **Escopo:** Sintetizar tom de voz, ritmo de cortes, paleta e matriz de ganchos dos 15 vídeos virais de pico.

#### Ticket `TICK-02.2` — Persistência determinística da pasta `repc/`
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/redc/writer.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E02.2`
* **Alvo de Teste TDD:** `internal/redc/writer_test.go:TestWriteREPC`
* **Outputs:** `internal/redc/writer.go`, `internal/redc/writer_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestWriteREPC ./internal/redc/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/writer.go')`
* **Spec Ref:** `AC-E02.2`
* **Escopo:** Gravar `01_dossie_dna.json`, `02_design_tokens.json`, `03_motion_tokens.json` e `04_matriz_ganchos.md` com integridade e sem caracteres corrompidos.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [x] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [x] Artefatos gerados em pasta temporária verificados e aprovados.

---

### Épico 03 — Gerador de Pipelines de Métodos e Brand Swap

#### Ticket `TICK-03.1` — Gerador de micro-pipelines canônicos dos métodos em `pipelines/`
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/redc/pipeline_gen.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E03.1`
* **Alvo de Teste TDD:** `internal/redc/pipeline_gen_test.go:TestGenerateMethodPipelines`
* **Outputs:** `internal/redc/pipeline_gen.go`, `internal/redc/pipeline_gen_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestGenerateMethodPipelines ./internal/redc/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/pipeline_gen.go')`
* **Spec Ref:** `AC-E03.1`
* **Escopo:** Construir manifestos `.loop.json` e DAGs dos métodos ensinados nos vídeos.

#### Ticket `TICK-03.2` — Gerador estrutural de vídeo com injeção de marca (Brand Swap)
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/redc/brand_swap.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E03.2`
* **Alvo de Teste TDD:** `internal/redc/brand_swap_test.go:TestGenerateStructuralVideo`
* **Outputs:** `internal/redc/brand_swap.go`, `internal/redc/brand_swap_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestGenerateStructuralVideo ./internal/redc/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/redc/brand_swap.go')`
* **Spec Ref:** `AC-E03.2`
* **Escopo:** Substituir tokens visuais da referência pela marca do usuário sem quebrar o timing de retenção.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o Épico 04?')`
- [x] 100% dos tickets de E03 executados com sucesso (PASS nos gates mecânicos).
- [x] Brand Swap testado e verificado com tema autônomo.

---

### Épico 04 — Integração CLI, Grafo DAG e Skill Canônica

#### Ticket `TICK-04.1` — Comandos CLI `legoloop canal redc` e `legoloop canal gerar`
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/cli/canal.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E04.1`
* **Alvo de Teste TDD:** `internal/cli/canal_test.go:TestCanalCLI`
* **Outputs:** `internal/cli/canal.go`, `internal/cli/canal_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 -run TestCanalCLI ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/canal.go')`
* **Spec Ref:** `AC-E04.1`
* **Escopo:** Plugar comandos na CLI do LegoLoop com suporte a `--json`, `--amostra` e `--marca`.

#### Ticket `TICK-04.2` — Grafo DAG acíclico `loops/legoloop-canal-redc/pipeline.dag.json`
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-04.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `loops/legoloop-canal-redc/pipeline.dag.json` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E04.3`
* **Outputs:** `loops/legoloop-canal-redc/pipeline.dag.json`, `loops/legoloop-canal-redc/legoloop-canal-redc.loop.json`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop validate loops/legoloop-canal-redc/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('loops/legoloop-canal-redc/pipeline.dag.json')`
* **Spec Ref:** `AC-E04.3`
* **Escopo:** Montar o grafo acíclico canônico com camadas de isolamento e validação estrita.

#### Ticket `TICK-04.3` — Skill Canônica `skills/legoloop-canal-redc/SKILL.md`
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `COMPLETED`
* **Dependências:** `["TICK-04.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `skills/legoloop-canal-redc/SKILL.md` · Spec Anchor `docs/specs/2026-11-07_SPEC_CANAL_REDC.md#AC-E04.1`
* **Outputs:** `skills/legoloop-canal-redc/SKILL.md`
* **Verificação (argv, sem shell):** `go test -v -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('skills/legoloop-canal-redc/SKILL.md')`
* **Spec Ref:** `AC-E04.1`
* **Escopo:** Criar a skill canônica universal e registrar nos harnesses.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Mega-Pipeline REDC 100% implementado, testado e validado? Autoriza a canonização final?')`
- [x] 100% dos tickets executados com sucesso.
- [x] `go test ./...` 100% verde e recompilação do binário fresca com `go build -o bin/legoloop.exe ./cmd/legoloop`.