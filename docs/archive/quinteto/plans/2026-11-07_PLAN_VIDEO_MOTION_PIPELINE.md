# 🗺️ PLAN: Pipeline Canônico de Edição e Motion Vídeo Determinado (LegoLoop Video Motion)

> **Status:** Em Planejamento  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** Antigravity (LegoLoop Mega Pipeline)  
> **Iniciativa:** I01 — LegoLoop Creative Engines  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_VIDEO_MOTION_PIPELINE.md](../intents/2026-11-07_INTENT_VIDEO_MOTION_PIPELINE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md](../specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_VIDEO_MOTION_PIPELINE.json](../routes/2026-11-07_ROUTE_VIDEO_MOTION_PIPELINE.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Speech Sync & Audio Clock Engine | 4 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4`, `AC-E01.5` | `micro` |
| **E02** | Visual Motion & Compositor Engine (Remotion/Canvas) | 4 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3`, `AC-E02.4` | `micro` |
| **E03** | Deterministic Asset Pipeline & AI Studio Bridge | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Render, Multi-Aspect Exporter & Broadcast Bundle | 3 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` | `micro` |
| **E05** | Orquestrador LegoLoop, Livro Didático LegoPress & Meta-Processo | 3 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3`, `AC-E05.4` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Speech Sync & Audio Clock Engine

#### Ticket `TICK-01.1` — Detecção de Silêncio e Fatiamento de Pausas com 40ms Pré-Ataque
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/silence.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/videomotion/silence_test.go:TestDetectSilenceAndCut`
* **Outputs:** `internal/videomotion/silence.go`, `internal/videomotion/silence_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestDetectSilenceAndCut`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/silence.go') && file_exists('internal/videomotion/silence_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Parser e Validador Fonético Palavra por Palavra (`words.json`)
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/words.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E01.2`
* **Alvo de Teste TDD:** `internal/videomotion/words_test.go:TestParseAndValidateWords`
* **Outputs:** `internal/videomotion/words.go`, `internal/videomotion/words_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestParseAndValidateWords`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/words.go') && file_exists('internal/videomotion/words_test.go')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.3` — Detector de Ênfase Vocal por Análise de Energia RMS
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/energy.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E01.3`
* **Alvo de Teste TDD:** `internal/videomotion/energy_test.go:TestDetectPunchWords`
* **Outputs:** `internal/videomotion/energy.go`, `internal/videomotion/energy_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestDetectPunchWords`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/energy.go') && file_exists('internal/videomotion/energy_test.go')`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.4` — Normalizador de Loudness EBU R128 (-14 LUFS) e Tratamento de Erros
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/loudness.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E01.4`
* **Alvo de Teste TDD:** `internal/videomotion/loudness_test.go:TestMeasureAndNormalizeLUFS`
* **Outputs:** `internal/videomotion/loudness.go`, `internal/videomotion/loudness_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestMeasureAndNormalizeLUFS`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/loudness.go') && file_exists('internal/videomotion/loudness_test.go')`
* **Spec Ref:** `AC-E01.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.

---

### Épico 02 — Visual Motion & Compositor Engine (Remotion/Canvas)

#### Ticket `TICK-02.1` — Gerador Procedural de Cartões de Vidro 3D (Glass Cards)
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.4"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/motion_glass.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/videomotion/motion_glass_test.go:TestGlassCardComposition`
* **Outputs:** `internal/videomotion/motion_glass.go`, `internal/videomotion/motion_glass_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestGlassCardComposition`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/motion_glass.go') && file_exists('internal/videomotion/motion_glass_test.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Grid em Perspectiva e Kinetic Typography no Beat da Fala
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/kinetic.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E02.2`
* **Alvo de Teste TDD:** `internal/videomotion/kinetic_test.go:TestKineticWordFlipTiming`
* **Outputs:** `internal/videomotion/kinetic.go`, `internal/videomotion/kinetic_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestKineticWordFlipTiming`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/kinetic.go') && file_exists('internal/videomotion/kinetic_test.go')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.3` — Auto-Framing Facial e Safe Zones Verticais (9:16)
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/framing.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E02.3`
* **Alvo de Teste TDD:** `internal/videomotion/framing_test.go:TestAutoFramingCenterCrop`
* **Outputs:** `internal/videomotion/framing.go`, `internal/videomotion/framing_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestAutoFramingCenterCrop`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/framing.go') && file_exists('internal/videomotion/framing_test.go')`
* **Spec Ref:** `AC-E02.3`

#### Ticket `TICK-02.4` — Gerador Instantâneo de Folha de Contato QA (`qa_sheet.jpg`)
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/qasheet.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E02.4`
* **Alvo de Teste TDD:** `internal/videomotion/qasheet_test.go:TestGenerateQASheet`
* **Outputs:** `internal/videomotion/qasheet.go`, `internal/videomotion/qasheet_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestGenerateQASheet`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/qasheet.go') && file_exists('internal/videomotion/qasheet_test.go')`
* **Spec Ref:** `AC-E02.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.

---

### Épico 03 — Deterministic Asset Pipeline & AI Studio Bridge

#### Ticket `TICK-03.1` — Gerenciador de Cache Criptográfico SHA-256 de Prompts
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.4"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/cache.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E03.1`
* **Alvo de Teste TDD:** `internal/videomotion/cache_test.go:TestAssetCacheHitMiss`
* **Outputs:** `internal/videomotion/cache.go`, `internal/videomotion/cache_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestAssetCacheHitMiss`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/cache.go') && file_exists('internal/videomotion/cache_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Bridge Google AI Studio e Fallback Gracioso para Texturas
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/bridge.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E03.2`
* **Alvo de Teste TDD:** `internal/videomotion/bridge_test.go:TestBridgeGenerationAndFallback`
* **Outputs:** `internal/videomotion/bridge.go`, `internal/videomotion/bridge_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestBridgeGenerationAndFallback`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/bridge.go') && file_exists('internal/videomotion/bridge_test.go')`
* **Spec Ref:** `AC-E03.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o Épico 04?')`
- [ ] 100% dos tickets de E03 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.

---

### Épico 04 — Render, Multi-Aspect Exporter & Broadcast Bundle

#### Ticket `TICK-04.1` — Gerador do Pacote de Entrega Broadcast (`dist/<job_id>/`)
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/bundle.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E04.1`
* **Alvo de Teste TDD:** `internal/videomotion/bundle_test.go:TestGenerateBroadcastBundle`
* **Outputs:** `internal/videomotion/bundle.go`, `internal/videomotion/bundle_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestGenerateBroadcastBundle`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/bundle.go') && file_exists('internal/videomotion/bundle_test.go')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.2` — Renderizador Multi-Aspect Ratio (9:16, 16:9, 1:1)
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/render.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E04.2`
* **Alvo de Teste TDD:** `internal/videomotion/render_test.go:TestRenderMultiAspect`
* **Outputs:** `internal/videomotion/render.go`, `internal/videomotion/render_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestRenderMultiAspect`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/render.go') && file_exists('internal/videomotion/render_test.go')`
* **Spec Ref:** `AC-E04.2`

#### Ticket `TICK-04.3` — Exportador de Timeline FCPXML para DaVinci Resolve Studio
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/videomotion/fcpxml.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E04.3`
* **Alvo de Teste TDD:** `internal/videomotion/fcpxml_test.go:TestExportFCPXMLTimeline`
* **Outputs:** `internal/videomotion/fcpxml.go`, `internal/videomotion/fcpxml_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestExportFCPXMLTimeline`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/fcpxml.go') && file_exists('internal/videomotion/fcpxml_test.go')`
* **Spec Ref:** `AC-E04.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o Épico 05?')`
- [ ] 100% dos tickets de E04 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.

---

### Épico 05 — Orquestrador LegoLoop, Livro Didático LegoPress & Meta-Processo

#### Ticket `TICK-05.1` — Grafo DAG Canônico `pipeline.dag.json` e Comando CLI
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `loops/video-motion/pipeline.dag.json`, `internal/cli/video_motion.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E05.1`
* **Alvo de Teste TDD:** `internal/cli/video_motion_test.go:TestVideoMotionCLI`
* **Outputs:** `loops/video-motion/pipeline.dag.json`, `internal/cli/video_motion.go`, `internal/cli/video_motion_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/cli/... -run TestVideoMotionCLI`
* **Gate (DSL):** `exit_code == 0 && file_exists('loops/video-motion/pipeline.dag.json') && file_exists('internal/cli/video_motion.go')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Gerador do Livro Didático Canônico de Uso (LegoPress)
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `docs/manuals/livro-video-motion/livro.json`, `internal/videomotion/book.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E05.2`
* **Alvo de Teste TDD:** `internal/videomotion/book_test.go:TestGenerateDidacticBook`
* **Outputs:** `docs/manuals/livro-video-motion/livro.json`, `internal/videomotion/book.go`, `internal/videomotion/book_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/videomotion/... -run TestGenerateDidacticBook`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/manuals/livro-video-motion/livro.json') && file_exists('internal/videomotion/book.go')`
* **Spec Ref:** `AC-E05.2`

#### Ticket `TICK-05.3` — Manual do Meta-Processo de Engenharia Reversa de Vídeo
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/manuals/MANUAL_ENGENHARIA_REVERSA_VIDEO_PIPELINE.md` · Spec Anchor `docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md#AC-E05.3`
* **Outputs:** `docs/manuals/MANUAL_ENGENHARIA_REVERSA_VIDEO_PIPELINE.md`
* **Verificação (argv, sem shell):** `python -c "assert len(open('docs/manuals/MANUAL_ENGENHARIA_REVERSA_VIDEO_PIPELINE.md', encoding='utf-8').read()) > 500"`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/manuals/MANUAL_ENGENHARIA_REVERSA_VIDEO_PIPELINE.md')`
* **Spec Ref:** `AC-E05.3`

---

### 🛑 STOP-GATE FINAL DE HOMOLOGAÇÃO & ENTREGA (`SG-FINAL`)
* **Gate (DSL):** `prompt('Pipeline Video Motion completo e homologado? Autoriza entrega final?')`
- [ ] 100% dos tickets de E01 a E05 executados com sucesso (PASS nos gates mecânicos).
- [ ] Livro didático LegoPress gerado em `docs/manuals/livro-video-motion/`.
- [ ] Binário `bin/legoloop.exe` recompilado e atualizado.
- [ ] Suíte global de testes verde (`go test -count=1 ./...`).
- [ ] Manual metodológico de engenharia reversa completo.
