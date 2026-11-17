# 🗺️ PLAN: Peças Híbridas, Triggers e Governança de Slots Canônicos

> **Status:** Proposta para Revisão (Stop-Gate 02)  
> **Data:** 2026-11-02  
> **Especificação Base:** [docs/specs/2026-11-02_SPEC_PECAS_HIBRIDAS_E_GOVERNANCA_SLOTS.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-02_SPEC_PECAS_HIBRIDAS_E_GOVERNANCA_SLOTS.md)  

---

## ÉPICO 01 — Taxonomia & Expansão de Tipos no DAG

### Node E01.1 — Expansão de NodeType e Configuração Dual-Mode de Gate

| Campo | Valor |
| :--- | :--- |
| **id** | `E01.1` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `[]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/dag/types.go` |
| **output** | `internal/dag/types.go` atualizado com novos `NodeType` e `GateConfig` |
| **gate** | `go test ./internal/dag/... -run TestTypes` |
| **spec_ref** | `AC-E01.1` |

### Node E01.2 — Validação Semântica de Novos Nós de Grafo

| Campo | Valor |
| :--- | :--- |
| **id** | `E01.2` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E01.1"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/dag/validate.go`, `internal/dag/validate_test.go` |
| **output** | `internal/dag/validate.go` com regras para gate, hook, mcp e binary |
| **gate** | `go test ./internal/dag -run TestValidateNodeTypes` |
| **spec_ref** | `AC-E01.2`, `AC-E01.3`, `AC-E01.4` |

### 🛑 STOP-GATE — ÉPICO 01
* [ ] Validação de novos tipos de nós implementada e testada
* [ ] Suíte `./internal/dag` 100% verde
* [ ] Usuário aprovou avanço para o Épico 02

---

## ÉPICO 02 — Padronização Universal de Triggers

### Node E02.1 — Normalização e Validação do Trigger Raiz

| Campo | Valor |
| :--- | :--- |
| **id** | `E02.1` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E01.2"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/dag/validate.go`, `internal/dag/types.go` |
| **output** | Validação estrita de triggers com fallback seguro para manual |
| **gate** | `go test ./internal/dag -run TestTriggerValidation` |
| **spec_ref** | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` |

### 🛑 STOP-GATE — ÉPICO 02
* [ ] Triggers validados como ponto de partida obrigatório de qualquer loop
* [ ] Compatibilidade retroativa com DAGs v1 preservada
* [ ] Usuário aprovou avanço para o Épico 03

---

## ÉPICO 03 — Governança de Slots & Mecanismo de Override

### Node E03.1 — Resolução Determinística de Slots e Overrides

| Campo | Valor |
| :--- | :--- |
| **id** | `E03.1` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E02.1"]` |
| **model** | `9router/code-pro` · effort `high` · fallback `sonnet` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/dag/slots.go`, `internal/dag/slots_test.go` |
| **output** | Função `ResolveEffectiveSlots` com testes unitários exaustivos |
| **gate** | `go test ./internal/dag -run TestSlotsOverride` |
| **spec_ref** | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-E03.4` |

### 🛑 STOP-GATE — ÉPICO 03
* [ ] `slots_override` implementado e isolado de peças canônicas originais
* [ ] Teste unitário de zero-override e delta-override 100% verde
* [ ] Usuário aprovou avanço para o Épico 04

---

## ÉPICO 04 — Adaptação Visual no Looping Studio

### Node E04.1 — Renderização de Triggers e Nós Híbridos no Canvas

| Campo | Valor |
| :--- | :--- |
| **id** | `E04.1` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E03.1"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/studio/static/app.js`, `internal/studio/static/ui.js` |
| **output** | Nós de trigger `⚡`, badges de `Stop-Gate` e saídas duplas no canvas |
| **gate** | `go test ./internal/studio/...` |
| **spec_ref** | `AC-E04.1`, `AC-E04.2` |

### Node E04.2 — Inspetor de Slots & Overrides no Studio

| Campo | Valor |
| :--- | :--- |
| **id** | `E04.2` |
| **nature** | `micro` |
| **execution_mode** | `sequential` |
| **depends_on** | `["E04.1"]` |
| **model** | `9router/code-fast` · effort `medium` · fallback `flash` |
| **tooling** | `edit_file`, `run_command` |
| **input** | `internal/studio/static/inspector.js` |
| **output** | Interface de configuração de overrides por nó |
| **gate** | `go test ./internal/studio/...` |
| **spec_ref** | `AC-E04.3` |

### 🛑 STOP-GATE — ÉPICO 04
* [ ] Visualização completa no Studio de nós de trigger, gates e slots overrides
* [ ] Suite `go test ./...` 100% verde
* [ ] Recompilação de `bin/legoloop.exe` realizada com sucesso
* [ ] Usuário aprovou entrega final
