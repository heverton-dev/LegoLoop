# 🗺️ PLAN: Hierarquia de Iniciativas no LegoLoop Mega Pipeline

> **Status:** Concluído  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** LegoLoop Mega Pipeline & Antigravity  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_HIERARQUIA_INICIATIVA.md](../intents/2026-11-07_INTENT_HIERARQUIA_INICIATIVA.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md](../specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_HIERARQUIA_INICIATIVA.json](../routes/2026-11-07_ROUTE_HIERARQUIA_INICIATIVA.json)  
> **Template:** TEMPLATE_PLAN.md v2.0.0 (espelha TEMPLATE_ROUTE.json v2.0.0)  

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
| **E01** | Padronização de Templates & Schemas Canônicos (SSOT) | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Suporte a Iniciativas no Montador de Rotas (`internal/plangen`) | 2 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` | `micro` |
| **E03** | Governança da Skill & Recompilação Fresca Go | 2 | `AC-E03.1`, `AC-NF.1` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Padronização de Templates & Schemas Canônicos (SSOT)

#### Ticket `TICK-01.1` — Atualização dos Templates Canônicos com Nível Iniciativa Opcional
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `docs/templates/TEMPLATE_INTENT.md`, `docs/templates/TEMPLATE_SPEC.md`, `docs/templates/TEMPLATE_PLAN.md`, `docs/templates/TEMPLATE_ROUTE.json`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-E01.1`
* **Alvo de Teste TDD:** `docs/templates/templates_test.go:TestEmbeddedTemplatesMatchFiles`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/templates/TEMPLATE_*.md`, `docs/templates/TEMPLATE_ROUTE.json`
  * **Verificação:** `go test -v -run TestEmbeddedTemplatesMatchFiles ./docs/templates/...`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.1`
* **Escopo:** Documenta e insere o nível Iniciativa (`I01`) opcional nos templates canônicos mantendo a integridade estrutural.

#### Ticket `TICK-01.2` — Garantia de Estabilidade de Headings e Schemas de Templates
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `400` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `docs/templates/templates_test.go`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-E01.2`
* **Alvo de Teste TDD:** `docs/templates/templates_test.go:TestDocumentHeadingsAreNumberedAndStable`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/templates/templates_test.go`
  * **Verificação:** `go test -v ./docs/templates/...`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.2`
* **Escopo:** Valida que os 4 templates mantêm estabilidade de cabeçalhos numerados sem quebrar os gates mecânicos.

---

### Épico 02 — Suporte a Iniciativas no Montador de Rotas (`internal/plangen`)

#### Ticket `TICK-02.1` — Parsing e Modelagem de Iniciativas em `meta.go` e `assemble.go`
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/plangen/meta.go`, `internal/plangen/assemble.go`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/plangen/plangen_test.go:TestInitiativeParsingAndAssembly`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/plangen/meta.go`, `internal/plangen/assemble.go`, `internal/plangen/plangen_test.go`
  * **Verificação:** `go test -v ./internal/plangen/...`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`, `AC-E02.3`
* **Escopo:** Adiciona regexes, parsing de `Initiatives` e mapeamento robusto pai-filho (`Initiative -> Epics`) com fallback para planos sem iniciativa.

#### Ticket `TICK-02.2` — Emissão de Nós de Rota com `initiative_id` em `route.go`
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `500` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/plangen/route.go`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-E02.2`
* **Alvo de Teste TDD:** `internal/plangen/plangen_test.go:TestRouteJSONInitiativeEmission`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/plangen/route.go`
  * **Verificação:** `go test -v ./internal/plangen/...`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`
* **Escopo:** Propaga `initiative_id` opcional para os nós da rota JSON mantendo compatibilidade omitempty com rotas legadas.

---

### Épico 03 — Governança da Skill & Recompilação Fresca Go

#### Ticket `TICK-03.1` — Atualização da Matriz de 4 Níveis na Skill do Mega Pipeline
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `medium`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `C:\Users\trcnologia\.gemini\skills\aidd-legoloop-mega-pipeline\SKILL.md`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-E03.1`
* **Alvo de Teste TDD:** Inspeção de tabela de 4 níveis e guia de classificação
* **Outputs:**
  * **Artifacts Produzidos:** `skills/aidd-legoloop-mega-pipeline/SKILL.md`
  * **Verificação:** `Get-Content C:\Users\trcnologia\.gemini\skills\aidd-legoloop-mega-pipeline\SKILL.md | Select-String "Iniciativa"`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.1`
* **Escopo:** Formaliza a governança normativa com a hierarquia Iniciativa ➔ Épico ➔ Micro ➔ Nano na skill de agentes.

#### Ticket `TICK-03.2` — Recompilação Go, Checagem de Funclen e Invariantes
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PASSED`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `antigravity` · Modelo `flash` · Effort `low`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `400` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `cmd/legoloop/main.go`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_HIERARQUIA_INICIATIVA.md#AC-NF.1`
* **Alvo de Teste TDD:** Execução completa dos portões de qualidade
* **Outputs:**
  * **Artifacts Produzidos:** `bin/legoloop.exe`
  * **Verificação:** `go test -count=1 ./...; go vet ./...; go run ./tools/funclen internal; go build -o bin/legoloop.exe ./cmd/legoloop`
  * **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-NF.1`
* **Escopo:** Recompila o binário e certifica conformidade com todas as 15 invariantes do ecossistema.
