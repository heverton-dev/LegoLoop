# 🗺️ PLAN: Matriz Canônica de Papéis e Casamento Determinístico de Peças

> **Status:** Aprovado pelo Usuário  
> **Data:** 2026-11-10  
> **Versão:** 1.0.0  
> **Autor:** Antigravity  
> **Iniciativa:** Fábrica Determinística LegoLoop  
> **Referência da Intenção:** [docs/intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md](../intents/2026-11-10_INTENT_MATRIZ_PAPEIS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md](../specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-10_ROUTE_MATRIZ_PAPEIS.json](../routes/2026-11-10_ROUTE_MATRIZ_PAPEIS.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false`
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Catálogo de Papéis e Validação Estrita (Taxonomia) | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Contrato da Peça na Matriz e Filtros na CLI | 2 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` | `micro` |
| **E03** | Casamento Determinístico no Gerador (Fabricar) | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Catálogo de Papéis e Validação Estrita (Taxonomia)

#### Ticket `TICK-01.1` — Implementar Catálogo dos 18 Papéis e Validador em Go
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `antigravity` · Modelo `flash` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1`
  * **Orçamento de Tokens:** Entrada `16000` tokens · Teto de saída `600` tokens
  * **Rollback:** `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/taxonomy/roles.go`
  * **Spec Anchor:** `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/taxonomy/roles_test.go:TestValidarPapel`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/taxonomy/roles.go`, `internal/taxonomy/roles_test.go`
  * **Verificação (argv, sem shell):** `go test -count=1 ./internal/taxonomy/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/taxonomy/roles.go') && file_exists('internal/taxonomy/roles_test.go')`
* **Spec Ref:** `AC-E01.1`, `AC-E01.2`

#### Ticket `TICK-01.2` — Sincronizar Dicionário JSON e Schemas Oficiais
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json`, `schemas/dag.v1.json` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-E01.1`
* **Outputs:** `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json`, `schemas/dag.v1.json`
* **Verificação (argv, sem shell):** `go run ./tools/taxlinter internal cmd loops skills docs/templates`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.1`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).

---

### Épico 02 — Contrato da Peça na Matriz e Filtros na CLI

#### Ticket `TICK-02.1` — Atualizar Contrato da Peça e Validação na Matriz
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `high`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/matriz/matriz.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-E02.1`
* **Outputs:** `internal/matriz/matriz.go`, `internal/matriz/matriz_test.go`
* **Alvo de Teste TDD:** `internal/matriz/matriz_test.go:TestPecaValidateRoles`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/matriz/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Adicionar Filtros --papel e --guilda na CLI da Matriz
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `high`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `cmd/legoloop/matriz.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-E02.2`
* **Outputs:** `cmd/legoloop/matriz.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./cmd/legoloop/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`, `AC-E02.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).

---

### Épico 03 — Casamento Determinístico no Gerador (Fabricar)

#### Ticket `TICK-03.1` — Implementar Algoritmo de Casamento em 2 Níveis no Fabricar
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `high`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/generator/generator.go` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md#AC-E03.1`
* **Outputs:** `internal/generator/generator.go`, `internal/generator/generator_test.go`
* **Alvo de Teste TDD:** `internal/generator/generator_test.go:TestCasamentoPapeis`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/generator/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.1`, `AC-E03.2`

#### Ticket `TICK-03.2` — Verificação Integrada e Recompilação Fresca
* **Iniciativa:** `Fábrica Determinística LegoLoop` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `bin/legoloop.exe` · Spec Anchor `docs/specs/2026-11-10_SPEC_MATRIZ_PAPEIS.md`
* **Outputs:** `bin/legoloop.exe`
* **Verificação (argv, sem shell):** `go build -o bin/legoloop.exe ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`
* **Spec Ref:** `AC-E03.1`, `AC-E03.2`

---

### 🛑 STOP-GATE HUMANO FINAL (`SG-FINAL`)
* **Gate (DSL):** `prompt('Todos os Épicos foram concluídos e validados? Autoriza o fechamento da demanda?')`
- [ ] 100% dos tickets executados com sucesso.
- [ ] Binário `bin/legoloop.exe` atualizado e testado.

