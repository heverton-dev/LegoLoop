# 🗺️ PLAN: AI Motion Graphics — Pipeline Soberano com 10 GitHub Repos

> **Status:** Em Planejamento  
> **Data:** 2026-11-09  
> **Versão:** 1.0.0  
> **Autor:** MiMoCode  
> **Iniciativa:** I-AIMG — AI Motion Graphics  
> **Referência da Intenção:** [docs/intents/2026-11-09_INTENT_AI_MOTION_GRAPHICS.md](../intents/2026-11-09_INTENT_AI_MOTION_GRAPHICS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md](../specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-09_ROUTE_AI_MOTION_GRAPHICS.json](../routes/2026-11-09_ROUTE_AI_MOTION_GRAPHICS.json)  
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
| **E01** | Motor de Seleção de Repos | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Motor de Geração HTML/WebGL | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Motor de Renderização Playwright + FFmpeg | 3 | `AC-E03.1`, `AC-E03.2`, `AC-NF.1` | `micro` |
| **E04** | Sync de Áudio Librosa | 2 | `AC-E04.1`, `AC-E04.2` | `micro` |
| **E05** | Skill Canônica & Livro Didático | 2 | `AC-E05.1`, `AC-E05.2` | `macro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Motor de Seleção de Repos

#### Ticket `TICK-01.1` — Tipos, contrato e validação de Request
* **Iniciativa:** `I-AIMG` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/types.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E01.3`
* **Alvo de Teste TDD:** `internal/motion/types_test.go:TestRequestValidate`
* **Outputs:** `internal/motion/types.go`, `internal/motion/types_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/types.go') && file_exists('internal/motion/types_test.go')`
* **Spec Ref:** `AC-E01.3`
* **Escopo:** Provar que `Request.Validate()` rejeita prompt vazio, duration <= 0 e scenario inválido com `ErrInvalidRequest`. Aceita cenários `explainer`, `product_launch`, `pixel_art`.

#### Ticket `TICK-01.2` — Motor de sugestão de repos por cenário
* **Iniciativa:** `I-AIMG` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/suggest.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/motion/suggest_test.go:TestSuggestRepos`
* **Outputs:** `internal/motion/suggest.go`, `internal/motion/suggest_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/suggest.go') && file_exists('internal/motion/suggest_test.go')`
* **Spec Ref:** `AC-E01.1`, `AC-E01.2`
* **Escopo:** `suggestRepos()` retorna `threejs` + `react-postprocessing` para `explainer`, `pixelfx` primeiro para `pixel_art`, `video-shotcraft` + `librosa` para `product_launch`. Confidence >= 0.7.

#### Ticket `TICK-01.3` — Verificação estática do motor de seleção
* **Iniciativa:** `I-AIMG` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/suggest.go`, `internal/motion/types.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E01.3`
* **Outputs:** `internal/motion/suggest.go`, `internal/motion/types.go`
* **Verificação (argv, sem shell):** `go vet ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Motor de Geração HTML/WebGL

#### Ticket `TICK-02.1` — Gerador HTML/WebGL com Three.js
* **Iniciativa:** `I-AIMG` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/build.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/motion/build_test.go:TestBuildHTML`
* **Outputs:** `internal/motion/build.go`, `internal/motion/build_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/build.go') && file_exists('internal/motion/build_test.go')`
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`
* **Escopo:** `buildHTML()` gera HTML com `<canvas>` e import Three.js para repos `threejs`. Retorna `ErrNoRepoSelected` para lista vazia. HTML < 50KB.

#### Ticket `TICK-02.2` — Verificação estática do gerador HTML
* **Iniciativa:** `I-AIMG` | **Épico:** `E02` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/build.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E02.2`
* **Outputs:** `internal/motion/build.go`
* **Verificação (argv, sem shell):** `go vet ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] 100% dos tickets de E02 executados com sucesso.
- [ ] Evidências de testes verdes apresentadas ao usuário.
- [ ] Autorização explícita concedida.

---

### Épico 03 — Motor de Renderização Playwright + FFmpeg

#### Ticket `TICK-03.1` — Captura de frames via Playwright
* **Iniciativa:** `I-AIMG` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/render.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E03.1`
* **Alvo de Teste TDD:** `internal/motion/render_test.go:TestCaptureFrames`
* **Outputs:** `internal/motion/render.go`, `internal/motion/render_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/render.go') && file_exists('internal/motion/render_test.go')`
* **Spec Ref:** `AC-E03.1`
* **Escopo:** `captureFrames()` abre Playwright headless, carrega HTML, captura N frames (fps × duration). 90 frames para 3s @ 30fps. Cada frame com dimensões do Request.

#### Ticket `TICK-03.2` — Conversão frames → MP4 via FFmpeg
* **Iniciativa:** `I-AIMG` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/encode.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E03.2`
* **Alvo de Teste TDD:** `internal/motion/encode_test.go:TestFramesToMP4`
* **Outputs:** `internal/motion/encode.go`, `internal/motion/encode_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/encode.go') && file_exists('internal/motion/encode_test.go')`
* **Spec Ref:** `AC-E03.2`
* **Escopo:** `framesToMP4()` invoca FFmpeg com `h264`, aceita frames + áudio opcional. `ffprobe` confirma duração ±0.1s e resolução correta.

#### Ticket `TICK-03.3` — Teste de integração com limite de tempo
* **Iniciativa:** `I-AIMG` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/render.go`, `internal/motion/encode.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-NF.1`
* **Outputs:** `internal/motion/render_test.go`, `internal/motion/encode_test.go`
* **Verificação (argv, sem shell):** `go test -v -timeout 300s ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-NF.1`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o Épico 04?')`
- [ ] 100% dos tickets de E03 executados com sucesso.
- [ ] Evidências de renderização apresentadas ao usuário.
- [ ] Autorização explícita concedida.

---

### Épico 04 — Sync de Áudio Librosa

#### Ticket `TICK-04.1` — Detecção de beats via Librosa
* **Iniciativa:** `I-AIMG` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03"]`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/sync.go`, `scripts/beat_detect.py` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E04.1`
* **Alvo de Teste TDD:** `internal/motion/sync_test.go:TestDetectBeats`
* **Outputs:** `internal/motion/sync.go`, `internal/motion/sync_test.go`, `scripts/beat_detect.py`
* **Verificação (argv, sem shell):** `go test -v ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/motion/sync.go') && file_exists('scripts/beat_detect.py')`
* **Spec Ref:** `AC-E04.1`, `AC-E04.2`
* **Escopo:** `detectBeats()` invoca `scripts/beat_detect.py` (Librosa). Retorna timestamps de beats ±50ms para 120 BPM. Fallback gracioso para áudio inexistente.

#### Ticket `TICK-04.2` — Verificação estática do sync de áudio
* **Iniciativa:** `I-AIMG` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/motion/sync.go` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E04.2`
* **Outputs:** `internal/motion/sync.go`
* **Verificação (argv, sem shell):** `go vet ./internal/motion/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E04.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o Épico 05?')`
- [ ] 100% dos tickets de E04 executados com sucesso.
- [ ] Evidências de sync apresentadas ao usuário.
- [ ] Autorização explícita concedida.

---

### Épico 05 — Skill Canônica & Livro Didático

#### Ticket `TICK-05.1` — Skill canônica universal
* **Iniciativa:** `I-AIMG` | **Épico:** `E05` | **Nature:** `macro` | **Tipo:** `subloop` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E04"]`
* **Ref do Subloop:** `legoloop-instalar-skills`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `medium` · Fallback Harness `ollama` · Modelo `deepseek-coder:6.7b`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `32000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Context Artifacts `["docs/videorea/ai-motion-graphics-10-repos/dossier.json"]` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E05.1`
* **Outputs:** `skills/ai-motion-graphics/SKILL.md`
* **Verificação (argv, sem shell):** `legoloop instalar-skills --dry-run --json`
* **Gate (DSL):** `exit_code == 0 && file_exists('skills/ai-motion-graphics/SKILL.md')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Livro Didático LegoPress
* **Iniciativa:** `I-AIMG` | **Épico:** `E05` | **Nature:** `macro` | **Tipo:** `subloop` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1"]`
* **Ref do Subloop:** `legoloop-legopress`
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `medium` · Fallback Harness `ollama` · Modelo `deepseek-coder:6.7b`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `32000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Context Artifacts `["docs/videorea/ai-motion-graphics-10-repos/curation.json"]` · Spec Anchor `docs/specs/2026-11-09_SPEC_AI_MOTION_GRAPHICS.md#AC-E05.2`
* **Outputs:** `docs/books/ai-motion-graphics/`
* **Verificação (argv, sem shell):** `legoloop legopress auditar`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/books/ai-motion-graphics/livro.json')`
* **Spec Ref:** `AC-E05.2`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test -count=1 ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia)
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal: `docs/reports/2026-11-09_REPORT_AI_MOTION_GRAPHICS.md`
5. Registro de histórico da sessão via `faz-session --id ses_ffe5e3d81734fffeURpuQcduSB --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA (`SG-ROUTE`):**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/2026-11-09_ROUTE_AI_MOTION_GRAPHICS.json`, o BUILD pode ser iniciado.*
