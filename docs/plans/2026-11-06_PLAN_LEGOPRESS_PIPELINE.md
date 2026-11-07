# 🗺️ PLAN: LegoPress™ — Pipeline Industrial e Template Canônico de Livros Técnicos

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-06  
> **Versão:** 1.0.0  
> **Autor:** Heverton & Ecossistema LegoLoop  
> **Referência da Intenção:** [docs/intents/2026-11-06_INTENT_LEGOPRESS_PIPELINE.md](../intents/2026-11-06_INTENT_LEGOPRESS_PIPELINE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md](../specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-06_ROUTE_LEGOPRESS_PIPELINE.json](../routes/2026-11-06_ROUTE_LEGOPRESS_PIPELINE.json)  
> **Template:** TEMPLATE_PLAN.md v2.0.0 (espelha TEMPLATE_ROUTE.json v2.0.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução cirúrgica isolada por bloco, sem context bloat)
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
| **E01** | Fundação de Tipos, Validação e Particionador | 2 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Motor de Compilação & Template Typst 170x240mm | 1 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Motor de Auditoria Mecânica & Verificação Geométrica | 1 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Motores de Adequação Didática & Sincronização por Diff | 1 | `AC-E04.1`, `AC-E04.2` | `micro` |
| **E05** | Integração na CLI & Template Canônico de Rota | 1 | `AC-E05.1`, `AC-E05.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Fundação de Tipos, Validação e Particionador

#### Ticket `TICK-01.1` — Implementação dos Contratos e Validador de Manifesto
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `[]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E01.1`
  * `context_artifacts`: `[]`
* **Test Target:** `internal/legopress/types_test.go:TestManifestValidation`
* **Output:**
  * `artifacts`: `["internal/legopress/types.go", "internal/legopress/types_test.go"]`
  * `gate`: `exit_code == 0 && go test -v -run TestManifestValidation ./internal/legopress/...`

#### Ticket `TICK-01.2` — Motor de Particionamento Macro/Meso/Micro
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-fast` · Effort `medium`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `["internal/legopress/types.go"]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E01.3`
  * `context_artifacts`: `[]`
* **Test Target:** `internal/legopress/partitioner_test.go:TestPartitionHierarchy`
* **Output:**
  * `artifacts`: `["internal/legopress/partitioner.go", "internal/legopress/partitioner_test.go"]`
  * `gate`: `exit_code == 0 && go test -v -run TestPartitionHierarchy ./internal/legopress/...`

---

### Épico 02 — Motor de Compilação & Template Typst 170x240mm

#### Ticket `TICK-02.1` — Compilador Pandoc + Typst com Verificação de Dependências
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-fast` · Effort `medium`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `["internal/legopress/types.go"]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E02.1`
  * `context_artifacts`: `["docs/livro/ativos/livro.typst"]`
* **Test Target:** `internal/legopress/compiler_test.go:TestCompilerBuildArgsAndDoctor`
* **Output:**
  * `artifacts`: `["internal/legopress/compiler.go", "internal/legopress/compiler_test.go"]`
  * `gate`: `exit_code == 0 && go test -v -run TestCompilerBuildArgsAndDoctor ./internal/legopress/...`

---

### Épico 03 — Motor de Auditoria Mecânica & Verificação Geométrica

#### Ticket `TICK-03.1` — Auditor de Geometria de PDF e Validador de Layout
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `["internal/legopress/types.go"]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E03.1`
  * `context_artifacts`: `["docs/livro/06-11-2026_LEGOLOOP-A-REVOLU-O-DA-IA-MODULAR-EM-PE-AS-SIMPLES.pdf"]`
* **Test Target:** `internal/legopress/auditor_test.go:TestAuditorGeometry`
* **Output:**
  * `artifacts`: `["internal/legopress/auditor.go", "internal/legopress/auditor_test.go"]`
  * `gate`: `exit_code == 0 && go test -v -run TestAuditorGeometry ./internal/legopress/...`

---

### Épico 04 — Motores de Adequação Didática & Sincronização por Diff

#### Ticket `TICK-04.1` — Transmutador Pedagógico e Sincronizador de Diffs do Git
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-fast` · Effort `medium`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `["internal/legopress/types.go"]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E04.1`
  * `context_artifacts`: `[]`
* **Test Target:** `internal/legopress/transmuter_test.go:TestTransmuterAndSync`
* **Output:**
  * `artifacts`: `["internal/legopress/transmuter.go", "internal/legopress/sync.go", "internal/legopress/transmuter_test.go"]`
  * `gate`: `exit_code == 0 && go test -v -run TestTransmuterAndSync ./internal/legopress/...`

---

### Épico 05 — Integração na CLI & Template Canônico de Rota

#### Ticket `TICK-05.1` — Exposição CLI `legoloop legopress` e Template de Rota
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1", "TICK-04.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens
* **Inputs:**
  * `source_files`: `["cmd/legoloop/main.go", "internal/legopress/types.go"]`
  * `spec_anchor`: `docs/specs/2026-11-06_SPEC_LEGOPRESS_PIPELINE.md#AC-E05.1`
  * `context_artifacts`: `[]`
* **Test Target:** `cmd/legoloop/legopress_test.go:TestLegoPressCLIIntegration`
* **Output:**
  * `artifacts`: `["cmd/legoloop/legopress.go", "cmd/legoloop/legopress_test.go", "docs/templates/TEMPLATE_ROUTE_LEGOPRESS.json"]`
  * `gate`: `exit_code == 0 && go build -o bin/legoloop.exe ./cmd/legoloop && go test -v -run TestLegoPressCLIIntegration ./cmd/legoloop/...`

---

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**  
*O plano estruturado de 6 tickets verticais em 5 épicos fecha com rigor todos os critérios de aceitação da SPEC?*  
*A aprovação humana explícita autoriza o início do Loop 4 (construção TDD Épico por Épico).*
