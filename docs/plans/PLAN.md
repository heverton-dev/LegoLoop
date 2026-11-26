---
id: PLAN-2026-11-25_ROUTE_LEGOLOOP_DIORAMA_RENDER
status: APPROVED
date: 2026-11-25
schema_version: 2
upstream: [2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER]
covers: [AC-E01.1, AC-E01.2, AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.5, AC-E02.7, AC-E02.6, AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E03.5, AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.4, AC-E04.5, AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-E05.5, AC-E06.1, AC-E06.2, AC-NF.1, AC-NF.2, AC-NF.3]
---

# 🗺️ PLAN: Processo Canonizado `legoloop-diorama-render` (Ciclo 3 — Etapas 6, 7, 8 e 8b)

> **Status:** APPROVED (Stop-Gate 03 aprovado em 2026-11-25; Emendas 1 e 3 autorizadas pelo usuário)  
> **Data:** 2026-11-25  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-25_INTENT_LEGOLOOP_DIORAMA_RENDER.md](../intents/2026-11-25_INTENT_LEGOLOOP_DIORAMA_RENDER.md)  
> **Referência da Especificação:** [docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md](../specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-25_ROUTE_LEGOLOOP_DIORAMA_RENDER.json](../routes/2026-11-25_ROUTE_LEGOLOOP_DIORAMA_RENDER.json)  
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
| **E02** | Épico E02 | 3 | `AC-E02.1`, `AC-E02.5`, `AC-E02.2`, `AC-E02.3`, `AC-E02.4`, `AC-E02.7`, `AC-E02.6` | `micro` |
| **E03** | Épico E03 | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-E03.4`, `AC-E03.5` | `micro` |
| **E04** | Épico E04 | 2 | `AC-E04.4`, `AC-E04.1`, `AC-E04.2`, `AC-E04.3`, `AC-E04.5` | `micro` |
| **E05** | Épico E05 | 3 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3`, `AC-E05.4`, `AC-E05.5` | `micro` |
| **E06** | Épico E06 | 2 | `AC-E06.2`, `AC-NF.1`, `AC-NF.2`, `AC-NF.3`, `AC-E06.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Épico E01

#### Ticket `TICK-01.1` — Criar contratos estritos de ConfigRender, Timings, Timeline, Batida, EventoSFX e RelatorioRender
* **Iniciativa:** `D03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/schemas.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/schemas_test.go`
* **Outputs:** Artifacts `internal/diorama/render/schemas.go, internal/diorama/render/schemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestSchemas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/schemas.go') && file_exists('internal/diorama/render/schemas_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Implementar pre-voo do Ciclo 3: cobertura de cenas, sprites, provedor registrado, ferramentas e render.json
* **Iniciativa:** `D03` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/prevoo.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/prevoo_test.go`
* **Outputs:** Artifacts `internal/diorama/render/prevoo.go, internal/diorama/render/prevoo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestPreVoo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/prevoo.go') && file_exists('internal/diorama/render/prevoo_test.go')`
* **Spec Ref:** `AC-E01.1, AC-E01.2, AC-NF.2`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Épico E02

#### Ticket `TICK-02.1` — Implementar Etapa 6 (voz): interface Sintetizador, registro de provedores, adaptador edge-tts e gate de creditos
* **Iniciativa:** `D03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/voz.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/voz_test.go`
* **Outputs:** Artifacts `internal/diorama/render/voz.go, internal/diorama/render/voz_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestVoz`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/voz.go') && file_exists('internal/diorama/render/voz_test.go')`
* **Spec Ref:** `AC-E02.1, AC-E02.5`

#### Ticket `TICK-02.2` — Implementar Etapa 6 (Whisper): Transcritor, WER por linha (15%) e global (5%), timings.json e regravacao seletiva com parada em audio identico
* **Iniciativa:** `D03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/whisper.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/whisper_test.go`
* **Outputs:** Artifacts `internal/diorama/render/whisper.go, internal/diorama/render/whisper_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestWhisper`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/whisper.go') && file_exists('internal/diorama/render/whisper_test.go')`
* **Spec Ref:** `AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.7`

#### Ticket `TICK-02.3` — Implementar mapeamento deterministico de fonemas PT-BR para visemas
* **Iniciativa:** `D03` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/visemas.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E02.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/visemas_test.go`
* **Outputs:** Artifacts `internal/diorama/render/visemas.go, internal/diorama/render/visemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestVisemas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/visemas.go') && file_exists('internal/diorama/render/visemas_test.go')`
* **Spec Ref:** `AC-E02.6`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Épico E03

#### Ticket `TICK-03.1` — Implementar Etapa 8: leitura estrita, camadas 2.5D por profundidade, camera e troca de pose/expressao
* **Iniciativa:** `D03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","TICK-02.3","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/timeline.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/timeline_test.go`
* **Outputs:** Artifacts `internal/diorama/render/timeline.go, internal/diorama/render/timeline_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestTimeline`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/timeline.go') && file_exists('internal/diorama/render/timeline_test.go')`
* **Spec Ref:** `AC-E03.1, AC-E03.2, AC-E03.3`

#### Ticket `TICK-03.2` — Implementar segmentacao em batidas e gate check_timeline (12-21 cortes/min, mediana 2,5-4 s)
* **Iniciativa:** `D03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/ritmo.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E03.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/ritmo_test.go`
* **Outputs:** Artifacts `internal/diorama/render/ritmo.go, internal/diorama/render/ritmo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestRitmo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/ritmo.go') && file_exists('internal/diorama/render/ritmo_test.go')`
* **Spec Ref:** `AC-E03.4`

#### Ticket `TICK-03.3` — Implementar legendas karaoke palavra a palavra e legenda.srt
* **Iniciativa:** `D03` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/legendas.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E03.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/legendas_test.go`
* **Outputs:** Artifacts `internal/diorama/render/legendas.go, internal/diorama/render/legendas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestLegendas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/legendas.go') && file_exists('internal/diorama/render/legendas_test.go')`
* **Spec Ref:** `AC-E03.5`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Épico E04

#### Ticket `TICK-04.1` — Corrigir teto de True Peak para -1,5 dBTP em videomotion.ValidateLoudnessConformity (achado F2)
* **Iniciativa:** `D03` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/videomotion/loudness.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E04.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/videomotion/loudness_test.go`
* **Outputs:** Artifacts `internal/videomotion/loudness.go, internal/videomotion/loudness_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/videomotion/ -run TestValidateLoudnessConformity`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/videomotion/loudness.go') && file_exists('internal/videomotion/loudness_test.go')`
* **Spec Ref:** `AC-E04.4`

#### Ticket `TICK-04.2` — Implementar Etapa 7: eventos de SFX, BGM com ducking, normalizacao em duas passadas com margem de pico e verificacao independente
* **Iniciativa:** `D03` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/sfx.go, internal/diorama/render/ativos_sonoros.go, internal/diorama/render/mix.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/sfx_test.go`
* **Outputs:** Artifacts `internal/diorama/render/sfx.go, internal/diorama/render/sfx_test.go, internal/diorama/render/ativos_sonoros.go, internal/diorama/render/ativos_sonoros_test.go, internal/diorama/render/mix.go, internal/diorama/render/mix_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestMix`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/sfx.go') && file_exists('internal/diorama/render/sfx_test.go') && file_exists('internal/diorama/render/ativos_sonoros.go') && file_exists('internal/diorama/render/ativos_sonoros_test.go') && file_exists('internal/diorama/render/mix.go') && file_exists('internal/diorama/render/mix_test.go')`
* **Spec Ref:** `AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.5`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Épico E05

#### Ticket `TICK-05.1` — Implementar sonda e escolha de encoder NVENC/QSV/AMF com fallback libx264
* **Iniciativa:** `D03` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/encoder.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/encoder_test.go`
* **Outputs:** Artifacts `internal/diorama/render/encoder.go, internal/diorama/render/encoder_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestEncoder`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/encoder.go') && file_exists('internal/diorama/render/encoder_test.go')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Implementar render master 16:9 em 720p/1080p/4K, versoes numeradas e gate de integridade
* **Iniciativa:** `D03` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","TICK-04.2","TICK-03.3","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/render.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E05.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/render_test.go`
* **Outputs:** Artifacts `internal/diorama/render/render.go, internal/diorama/render/render_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestRender`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/render.go') && file_exists('internal/diorama/render/render_test.go')`
* **Spec Ref:** `AC-E05.2, AC-E05.3`

#### Ticket `TICK-05.3` — Implementar cortes 9:16 por re-render da timeline com area segura e legenda fora do narrador
* **Iniciativa:** `D03` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2","TICK-03.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/vertical.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E05.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/vertical_test.go`
* **Outputs:** Artifacts `internal/diorama/render/vertical.go, internal/diorama/render/vertical_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/ -run TestVertical`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/vertical.go') && file_exists('internal/diorama/render/vertical_test.go')`
* **Spec Ref:** `AC-E05.4, AC-E05.5`

### 🛑 STOP-GATE HUMANO (`SG-E05`)
* **Gate (DSL):** `prompt('Evidências do Épico 05 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 06 — Épico E06

#### Ticket `TICK-06.1` — Orquestrar Processo completo, E2E hermetico no diorama-demo, determinismo e cobertura >= 95%
* **Iniciativa:** `D03` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2","TICK-03.2","TICK-03.3","TICK-04.2","TICK-05.3","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/render/processo.go` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E06.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/render/processo_test.go`
* **Outputs:** Artifacts `internal/diorama/render/processo.go, internal/diorama/render/processo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -cover ./internal/diorama/render/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/render/processo.go') && file_exists('internal/diorama/render/processo_test.go')`
* **Spec Ref:** `AC-E06.2, AC-NF.1, AC-NF.2, AC-NF.3`

#### Ticket `TICK-06.2` — Integrar CLI legoloop diorama-render, DAG v2, Skill e fechamento 5 de 5 via legoloop fechar
* **Iniciativa:** `D03` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/diorama_render.go, internal/cli/cli.go, processos/legoloop-diorama-render/processo.dag.json, skills/legoloop-diorama-render/SKILL.md` · Spec Anchor `docs/specs/2026-11-25_SPEC_LEGOLOOP_DIORAMA_RENDER.md#AC-E06.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/diorama_render_test.go`
* **Outputs:** Artifacts `internal/cli/diorama_render.go, internal/cli/diorama_render_test.go, internal/cli/cli.go, processos/legoloop-diorama-render/processo.dag.json, skills/legoloop-diorama-render/SKILL.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/render/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/diorama_render.go') && file_exists('internal/cli/diorama_render_test.go') && file_exists('internal/cli/cli.go') && file_exists('processos/legoloop-diorama-render/processo.dag.json') && file_exists('skills/legoloop-diorama-render/SKILL.md')`
* **Spec Ref:** `AC-E06.1`

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
*Plano de execução fatiado com 15 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
