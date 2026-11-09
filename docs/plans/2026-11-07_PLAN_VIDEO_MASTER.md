# 🗺️ PLAN: Pipeline Master de Produção Audiovisual Soberana

> **Status:** Concluído  
> **Data:** 2026-11-07  
> **Versão:** 1.1.0  
> **Autor:** LegoLoop REA Engine  
> **Iniciativa:** Audiovisual Soberano 2026  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_VIDEO_MASTER.md](../intents/2026-11-07_INTENT_VIDEO_MASTER.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md](../specs/2026-11-07_SPEC_VIDEO_MASTER.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_VIDEO_MASTER.json](../routes/2026-11-07_ROUTE_VIDEO_MASTER.json)  
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
| **E01** | Síntese Forense & Matriz de Soberania | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Grafo DAG Acíclico & Golden Fixture | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Livro LegoPress & Skill Universal | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Síntese Forense & Matriz de Soberania

#### Ticket `TICK-01.1` — Consolidação do Dossiê Forense dos 43 Vídeos
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json`
  * **Verificação (argv, sem shell):** `Test-Path docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json`
  * **Gate (DSL):** `exit_code == 0 && file_exists('docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Emissão da Matriz de Curadoria e Stop-Gate
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `stop_gate` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E01.2`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/videorea/2026-11-07_CURATION_VIDEO_MASTER.json`
  * **Verificação (argv, sem shell):** `Test-Path docs/videorea/2026-11-07_CURATION_VIDEO_MASTER.json`
  * **Gate (DSL):** `exit_code == 0 && file_exists('docs/videorea/2026-11-07_CURATION_VIDEO_MASTER.json')`
* **Spec Ref:** `AC-E01.2`

---

### Épico 02 — Grafo DAG Acíclico & Golden Fixture

#### Ticket `TICK-02.1` — Construção e Validação do Grafo DAG Mestre
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E02.1`
* **Outputs:**
  * **Artifacts Produzidos:** `loops/legoloop-video-master/pipeline.dag.json`, `loops/legoloop-video-master/legoloop-video-master.loop.json`
  * **Verificação (argv, sem shell):** `legoloop validate loops/legoloop-video-master/pipeline.dag.json`
  * **Gate (DSL):** `exit_code == 0 && file_exists('loops/legoloop-video-master/pipeline.dag.json')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Geração Determinística da Golden Fixture de 3s
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E02` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videorea/factory_bridge.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E02.2`
* **Outputs:**
  * **Artifacts Produzidos:** `testdata/fixtures/video_master_fixture_3s.mp4`
  * **Verificação (argv, sem shell):** `Test-Path testdata/fixtures/video_master_fixture_3s.mp4`
  * **Gate (DSL):** `exit_code == 0 && file_exists('testdata/fixtures/video_master_fixture_3s.mp4')`
* **Spec Ref:** `AC-E02.2`

---

### Épico 03 — Livro LegoPress & Skill Universal

#### Ticket `TICK-03.1` — Compilação e Estruturação do Livro LegoPress
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `docs/videorea/2026-11-07_DOSSIER_VIDEO_MASTER.json` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E03.1`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/books/video-master/livro.json`, `docs/books/video-master/2026-11-07_LIVRO_LEGOLOOP_VIDEO_MASTER.pdf`
  * **Verificação (argv, sem shell):** `legoloop legopress auditar docs/books/video-master/2026-11-07_LIVRO_LEGOLOOP_VIDEO_MASTER.pdf`
  * **Gate (DSL):** `exit_code == 0 && file_exists('docs/books/video-master/2026-11-07_LIVRO_LEGOLOOP_VIDEO_MASTER.pdf')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Registro da Skill Universal nos Harnesses
* **Iniciativa:** `Audiovisual Soberano 2026` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `DONE`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `skills/legoloop-video-master/SKILL.md` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md#AC-E03.2`
* **Outputs:**
  * **Artifacts Produzidos:** `skills/legoloop-video-master/SKILL.md`
  * **Verificação (argv, sem shell):** `Test-Path skills/legoloop-video-master/SKILL.md`
  * **Gate (DSL):** `exit_code == 0 && file_exists('skills/legoloop-video-master/SKILL.md')`
* **Spec Ref:** `AC-E03.2`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

* **Critério de Parada Mecânico:**
  * Todos os 6 tickets marcados como `DONE`.
  * `legoloop validate loops/legoloop-video-master/pipeline.dag.json` com retorno 0.
  * `legoloop legopress auditar docs/books/video-master/2026-11-07_LIVRO_LEGOLOOP_VIDEO_MASTER.pdf` com retorno 0 (geometria 100% válida).
  * Suíte global `go test -count=1 ./...` passando 100% verde sem warnings.
  * Relatório final registrado em `docs/reports/2026-11-07_REPORT_VIDEO_MASTER.md` com balanço orçado vs realizado.
