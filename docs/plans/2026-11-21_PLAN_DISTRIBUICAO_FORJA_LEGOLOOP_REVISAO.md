---
id: PLAN-2026-11-21_ROUTE_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO
status: DRAFT
date: 2026-11-22
schema_version: 2
upstream: [2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E03.1, AC-E03.2, AC-E03.3, AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.4, AC-E04.5, AC-E04.6, AC-E04.7, AC-NF.1, AC-NF.2, AC-NF.3, AC-E05.1, AC-E05.2, AC-E05.3, AC-E06.1, AC-E06.2, AC-E07.1, AC-E07.2, AC-E07.3, AC-E08.1, AC-E08.2, AC-E08.3, AC-E09.1, AC-E09.2, AC-E09.3, AC-E09.4, AC-E09.5, AC-E09.6, AC-E10.1, AC-E10.2, AC-E11.1, AC-E11.2]
---

# 🗺️ PLAN: Distribuição e Forja do LegoLoop (Revisão Emenda 2)

> **Status:** Em Planejamento  
> **Data:** 2026-11-22  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md](../intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md)  
> **Referência da Especificação:** [docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md](../specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-21_ROUTE_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.json](../routes/2026-11-21_ROUTE_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.json)  
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
| **E01** | Contrato e linha de base | 3 | `AC-E01.1`, `AC-E01.3`, `AC-E01.2` | `micro` |
| **E03** | Python do caminho de instalação para Go | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `micro` |
| **E04** | Motor da Forja | 11 | `AC-E04.6`, `AC-E04.1`, `AC-E04.2`, `AC-E04.3`, `AC-E04.4`, `AC-E04.5`, `AC-E04.6`, `AC-E04.7`, `AC-NF.1`, `AC-NF.2`, `AC-NF.3` | `micro` |
| **E05** | `legoloop instalar` | 3 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3` | `micro` |
| **E06** | Gate de orçamento de tokens | 2 | `AC-E06.1`, `AC-E06.2` | `micro` |
| **E07** | Canonizar o Processo Forja | 3 | `AC-E07.2`, `AC-E07.3`, `AC-E07.1` | `micro` |
| **E08** | Release reproduzível, assinado e com SBOM | 3 | `AC-E08.1`, `AC-E08.2`, `AC-E08.3` | `micro` |
| **E09** | Canais de distribuição | 6 | `AC-E09.1`, `AC-E09.2`, `AC-E09.3`, `AC-E09.6`, `AC-E09.4`, `AC-E09.5` | `micro` |
| **E10** | Detecção de deriva | 2 | `AC-E10.1`, `AC-E10.2` | `micro` |
| **E11** | Canal de atualização | 2 | `AC-E11.1`, `AC-E11.2` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Contrato e linha de base

#### Ticket `TICK-01.1` — Conferir a SPEC revisada aprovada contra os Gates de lacuna e de aceite
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md`
* **Verificação (argv, sem shell):** `go test ./internal/gate/ -run TestAceiteReconheceOTituloCanonicoDoTemplate`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Renomear o módulo para github.com/heverton-dev/LegoLoop (pré-requisitos I04: G01 Matriz por hash, G02 descoberta de peças, G05 fachada única de caminhos)
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `go.mod, cmd, internal, tools` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `go.mod, cmd, internal, tools`
* **Verificação (argv, sem shell):** `go build ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('go.mod') && file_exists('cmd') && file_exists('internal') && file_exists('tools')`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.3` — Gravar a linha de base de contexto por harness
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/reports/2026-11-21_BASELINE_FORJA.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/reports/2026-11-21_BASELINE_FORJA.json`
* **Verificação (argv, sem shell):** `legoloop skills-custo --json`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-21_BASELINE_FORJA.json')`
* **Spec Ref:** `AC-E01.2`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Python do caminho de instalação para Go

#### Ticket `TICK-03.1` — Criar ste.AplicarBloco idempotente para o bloco gerenciado
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/bloco.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/bloco_test.go`
* **Outputs:** Artifacts `internal/ste/bloco.go, internal/ste/bloco_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/ste/ -run TestAplicarBloco`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/bloco.go') && file_exists('internal/ste/bloco_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Provar a paridade com o script Python e remover o script
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/ste/testdata/gemini_golden.md, tools/steptsync/gemini_md.py` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/ste/bloco_test.go`
* **Outputs:** Artifacts `internal/ste/testdata/gemini_golden.md, internal/ste/bloco_test.go, tools/steptsync/gemini_md.py`
* **Verificação (argv, sem shell):** `go test ./internal/ste/ -run TestBlocoIgualAoGolden`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/ste/testdata/gemini_golden.md') && file_exists('internal/ste/bloco_test.go') && file_exists('tools/steptsync/gemini_md.py')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `TICK-03.3` — Provar que o caminho de instalação não usa Python e listar Python como opcional no doctor
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/doc.go, internal/cli/doctor.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E03.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/sempython_test.go`
* **Outputs:** Artifacts `internal/forja/doc.go, internal/forja/sempython_test.go, internal/cli/doctor.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestCaminhoSemPython`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/doc.go') && file_exists('internal/forja/sempython_test.go') && file_exists('internal/cli/doctor.go')`
* **Spec Ref:** `AC-E03.3`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Motor da Forja

#### Ticket `TICK-04.1` — Criar os contratos da Forja e a validação das opções
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/contrato.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/contrato_test.go`
* **Outputs:** Artifacts `internal/forja/contrato.go, internal/forja/contrato_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestContrato`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/contrato.go') && file_exists('internal/forja/contrato_test.go')`
* **Spec Ref:** `AC-E04.6`

#### Ticket `TICK-04.2` — Implementar a simulação --dry-run que não grava nada
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/forja.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/forja_test.go`
* **Outputs:** Artifacts `internal/forja/forja.go, internal/forja/forja_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestSimulacaoNaoGrava`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/forja.go') && file_exists('internal/forja/forja_test.go')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.3` — Implementar o aplicar com backup, MANIFESTO e idempotência
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/aplicar.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/aplicar_test.go`
* **Outputs:** Artifacts `internal/forja/aplicar.go, internal/forja/aplicar_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestAplicar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/aplicar.go') && file_exists('internal/forja/aplicar_test.go')`
* **Spec Ref:** `AC-E04.2`

#### Ticket `TICK-04.4` — Implementar o --desfazer que restaura byte a byte
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/desfazer.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/desfazer_test.go`
* **Outputs:** Artifacts `internal/forja/desfazer.go, internal/forja/desfazer_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestDesfazer`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/desfazer.go') && file_exists('internal/forja/desfazer_test.go')`
* **Spec Ref:** `AC-E04.3`

#### Ticket `TICK-04.5` — Gravar o bloco gerenciado que aponta para o AGENTS.md no GEMINI.md e no CLAUDE.md
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/ssot.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/ssot_test.go`
* **Outputs:** Artifacts `internal/forja/ssot.go, internal/forja/ssot_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestSSOT`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/ssot.go') && file_exists('internal/forja/ssot_test.go')`
* **Spec Ref:** `AC-E04.4`

#### Ticket `TICK-04.6` — Mesclar o settings.json do Claude preservando as chaves desconhecidas
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/settings_claude.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/settings_claude_test.go`
* **Outputs:** Artifacts `internal/forja/settings_claude.go, internal/forja/settings_claude_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestSettingsClaude`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/settings_claude.go') && file_exists('internal/forja/settings_claude_test.go')`
* **Spec Ref:** `AC-E04.5`

#### Ticket `TICK-04.7` — Implementar --sem e recusar perfil inválido
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/opcoes.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/opcoes_test.go`
* **Outputs:** Artifacts `internal/forja/opcoes.go, internal/forja/opcoes_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestSemEPerfil`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/opcoes.go') && file_exists('internal/forja/opcoes_test.go')`
* **Spec Ref:** `AC-E04.6`

#### Ticket `TICK-04.8` — Gerar os lançadores por SO sem PATH fixo
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/lancadores.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E04.7` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/lancadores_test.go`
* **Outputs:** Artifacts `internal/forja/lancadores.go, internal/forja/lancadores_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestLancadores`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/lancadores.go') && file_exists('internal/forja/lancadores_test.go')`
* **Spec Ref:** `AC-E04.7`

#### Ticket `TICK-04.9` — Garantir que segredos não entram no relatório nem na saída
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/relatorio.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-NF.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/relatorio_test.go`
* **Outputs:** Artifacts `internal/forja/relatorio.go, internal/forja/relatorio_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestSegredos`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/relatorio.go') && file_exists('internal/forja/relatorio_test.go')`
* **Spec Ref:** `AC-NF.1`

#### Ticket `TICK-04.10` — Aplicar path jail na resolução da pasta do harness
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/jail.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-NF.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/jail_test.go`
* **Outputs:** Artifacts `internal/forja/jail.go, internal/forja/jail_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestPathJail`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/jail.go') && file_exists('internal/forja/jail_test.go')`
* **Spec Ref:** `AC-NF.2`

#### Ticket `TICK-04.11` — Medir o desempenho do dry-run com 24 harnesses
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/desempenho_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-NF.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/desempenho_test.go`
* **Outputs:** Artifacts `internal/forja/desempenho_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestDesempenhoDryRun`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/desempenho_test.go')`
* **Spec Ref:** `AC-NF.3`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — `legoloop instalar`

#### Ticket `TICK-05.1` — Criar legoloop instalar encadeando a Forja e entregando o relatório
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/instalar/instalar.go, internal/cli/instalar.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/instalar/instalar_test.go`
* **Outputs:** Artifacts `internal/instalar/instalar.go, internal/instalar/instalar_test.go, internal/cli/instalar.go`
* **Verificação (argv, sem shell):** `go test ./internal/instalar/ -run TestInstalarEntregaRelatorio`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/instalar/instalar.go') && file_exists('internal/instalar/instalar_test.go') && file_exists('internal/cli/instalar.go')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Criar o atalho do grafo semântico no PATH, com dry-run
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/instalar/path.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E05.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/instalar/path_test.go`
* **Outputs:** Artifacts `internal/instalar/path.go, internal/instalar/path_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/instalar/ -run TestAtalhoDoGrafo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/instalar/path.go') && file_exists('internal/instalar/path_test.go')`
* **Spec Ref:** `AC-E05.2`

#### Ticket `TICK-05.3` — Marcar o bootstrap como substituído pelo instalar
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/bootstrap.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E05.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/bootstrap_test.go`
* **Outputs:** Artifacts `internal/cli/bootstrap.go, internal/cli/bootstrap_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestBootstrap`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/bootstrap.go') && file_exists('internal/cli/bootstrap_test.go')`
* **Spec Ref:** `AC-E05.3`

### 🛑 STOP-GATE HUMANO (`SG-E05`)
* **Gate (DSL):** `prompt('Evidências do Épico 05 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 06 — Gate de orçamento de tokens

#### Ticket `TICK-06.1` — Implementar forjar --verificar --orcamento com reprovação por teto estourado
* **Iniciativa:** `I01` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/orcamento.go, docs/governanca/orcamento_tokens.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E06.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/orcamento_test.go`
* **Outputs:** Artifacts `internal/forja/orcamento.go, internal/forja/orcamento_test.go, docs/governanca/orcamento_tokens.json`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestOrcamentoEstourado`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/orcamento.go') && file_exists('internal/forja/orcamento_test.go') && file_exists('docs/governanca/orcamento_tokens.json')`
* **Spec Ref:** `AC-E06.1`

#### Ticket `TICK-06.2` — Aprovar quando o custo respeita o teto
* **Iniciativa:** `I01` | **Épico:** `E06` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/orcamento_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E06.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/orcamento_test.go`
* **Outputs:** Artifacts `internal/forja/orcamento_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestOrcamentoRespeitado`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/orcamento_test.go')`
* **Spec Ref:** `AC-E06.2`

### 🛑 STOP-GATE HUMANO (`SG-E06`)
* **Gate (DSL):** `prompt('Evidências do Épico 06 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 07 — Canonizar o Processo Forja

#### Ticket `TICK-07.1` — Criar o Processo legoloop-forjar: DAG, skill e slash command
* **Iniciativa:** `I01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-forjar, skills/legoloop-forjar/SKILL.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E07.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-forjar, skills/legoloop-forjar/SKILL.md`
* **Verificação (argv, sem shell):** `legoloop validate processos/legoloop-forjar/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-forjar') && file_exists('skills/legoloop-forjar/SKILL.md')`
* **Spec Ref:** `AC-E07.2`

#### Ticket `TICK-07.2` — Provar que os quatro gatilhos geram o mesmo RELATORIO.json
* **Iniciativa:** `I01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/gatilhos_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E07.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/gatilhos_test.go`
* **Outputs:** Artifacts `internal/forja/gatilhos_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestQuatroGatilhosMesmoRelatorio`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/gatilhos_test.go')`
* **Spec Ref:** `AC-E07.3`

#### Ticket `TICK-07.3` — Fechar o Processo com legoloop fechar
* **Iniciativa:** `I01` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.2","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/manuals, docs/books/legoloop-forjar, docs/reports` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E07.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/manuals, docs/books/legoloop-forjar, docs/reports`
* **Verificação (argv, sem shell):** `legoloop fechar legoloop-forjar --json`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/manuals') && file_exists('docs/books/legoloop-forjar') && file_exists('docs/reports')`
* **Spec Ref:** `AC-E07.1`

### 🛑 STOP-GATE HUMANO (`SG-E07`)
* **Gate (DSL):** `prompt('Evidências do Épico 07 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 08 — Release reproduzível, assinado e com SBOM

#### Ticket `TICK-08.1` — Criar tools/release com build reproduzível
* **Iniciativa:** `I01` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/release/release.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E08.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/release_test.go`
* **Outputs:** Artifacts `tools/release/release.go, tools/release/release_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestBuildReproduzivel`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/release/release.go') && file_exists('tools/release/release_test.go')`
* **Spec Ref:** `AC-E08.1`

#### Ticket `TICK-08.2` — Gerar pacotes por SO e arquitetura, SHA256SUMS e SBOM
* **Iniciativa:** `I01` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/release/pacotes.go, tools/release/sbom.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E08.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/pacotes_test.go`
* **Outputs:** Artifacts `tools/release/pacotes.go, tools/release/sbom.go, tools/release/pacotes_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestPacotesSumsESBOM`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/release/pacotes.go') && file_exists('tools/release/sbom.go') && file_exists('tools/release/pacotes_test.go')`
* **Spec Ref:** `AC-E08.2`

#### Ticket `TICK-08.3` — Declarar a assinatura no relatório do release e assinar com cosign no CI
* **Iniciativa:** `I01` | **Épico:** `E08` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.2","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/release/assinatura.go, .github/workflows/release.yml` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E08.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/assinatura_test.go`
* **Outputs:** Artifacts `tools/release/assinatura.go, tools/release/assinatura_test.go, .github/workflows/release.yml`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestAssinaturaDeclarada`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/release/assinatura.go') && file_exists('tools/release/assinatura_test.go') && file_exists('.github/workflows/release.yml')`
* **Spec Ref:** `AC-E08.3`

### 🛑 STOP-GATE HUMANO (`SG-E08`)
* **Gate (DSL):** `prompt('Evidências do Épico 08 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 09 — Canais de distribuição

#### Ticket `TICK-09.1` — Criar install.sh e install.ps1 que conferem o SHA-256 e abortam se divergir
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `install.sh, install.ps1` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/install_test.go`
* **Outputs:** Artifacts `install.sh, install.ps1, tools/release/install_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestScriptAbortaComShaErrado`
* **Gate (DSL):** `exit_code == 0 && file_exists('install.sh') && file_exists('install.ps1') && file_exists('tools/release/install_test.go')`
* **Spec Ref:** `AC-E09.1`

#### Ticket `TICK-09.2` — Fazer os scripts instalarem o binário e dispararem a Forja
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.1","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `install.sh, install.ps1` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/install_test.go`
* **Outputs:** Artifacts `install.sh, install.ps1, tools/release/install_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestScriptInstalaEDisparaForja`
* **Gate (DSL):** `exit_code == 0 && file_exists('install.sh') && file_exists('install.ps1') && file_exists('tools/release/install_test.go')`
* **Spec Ref:** `AC-E09.2`

#### Ticket `TICK-09.3` — Gerar os manifestos de winget, scoop e brew
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.1","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/release/manifestos.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/manifestos_test.go`
* **Outputs:** Artifacts `tools/release/manifestos.go, tools/release/manifestos_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestManifestos`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/release/manifestos.go') && file_exists('tools/release/manifestos_test.go')`
* **Spec Ref:** `AC-E09.3`

#### Ticket `TICK-09.4` — Criar a definição WiX do MSI por usuário, sem administrador
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.3","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/release/wix/legoloop.wxs` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/release/msi_test.go`
* **Outputs:** Artifacts `tools/release/wix/legoloop.wxs, tools/release/msi_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/release/ -run TestMSIPorUsuario`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/release/wix/legoloop.wxs') && file_exists('tools/release/msi_test.go')`
* **Spec Ref:** `AC-E09.6`

#### Ticket `TICK-09.5` — Provar git clone e go install remoto
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.2","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/manuals` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/manuals`
* **Verificação (argv, sem shell):** `go install ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/manuals')`
* **Spec Ref:** `AC-E09.4`

#### Ticket `TICK-09.6` — Expor o instalar como skill, slash command e ferramenta do MCP
* **Iniciativa:** `I01` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.2","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `skills/legoloop-instalar/SKILL.md, internal/mcp/instalar_tool.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E09.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/mcp/instalar_tool_test.go`
* **Outputs:** Artifacts `skills/legoloop-instalar/SKILL.md, internal/mcp/instalar_tool.go, internal/mcp/instalar_tool_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/mcp/ -run TestInstalarViaMCP`
* **Gate (DSL):** `exit_code == 0 && file_exists('skills/legoloop-instalar/SKILL.md') && file_exists('internal/mcp/instalar_tool.go') && file_exists('internal/mcp/instalar_tool_test.go')`
* **Spec Ref:** `AC-E09.5`

### 🛑 STOP-GATE HUMANO (`SG-E09`)
* **Gate (DSL):** `prompt('Evidências do Épico 09 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 10 — Detecção de deriva

#### Ticket `TICK-10.1` — Implementar forjar --verificar com detecção de deriva
* **Iniciativa:** `I01` | **Épico:** `E10` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/deriva.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E10.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/deriva_test.go`
* **Outputs:** Artifacts `internal/forja/deriva.go, internal/forja/deriva_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestDerivaReprova`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/deriva.go') && file_exists('internal/forja/deriva_test.go')`
* **Spec Ref:** `AC-E10.1`

#### Ticket `TICK-10.2` — Implementar --corrigir que devolve a skill ao Cofre com backup
* **Iniciativa:** `I01` | **Épico:** `E10` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-10.1","SG-E09"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/forja/corrigir.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E10.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/forja/corrigir_test.go`
* **Outputs:** Artifacts `internal/forja/corrigir.go, internal/forja/corrigir_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/forja/ -run TestCorrigirDeriva`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/forja/corrigir.go') && file_exists('internal/forja/corrigir_test.go')`
* **Spec Ref:** `AC-E10.2`

### 🛑 STOP-GATE HUMANO (`SG-E10`)
* **Gate (DSL):** `prompt('Evidências do Épico 10 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 11 — Canal de atualização

#### Ticket `TICK-11.1` — Criar legoloop atualizar que aborta com SHA-256 divergente
* **Iniciativa:** `I01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/atualizar/atualizar.go, internal/cli/atualizar.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E11.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/atualizar/atualizar_test.go`
* **Outputs:** Artifacts `internal/atualizar/atualizar.go, internal/atualizar/atualizar_test.go, internal/cli/atualizar.go`
* **Verificação (argv, sem shell):** `go test ./internal/atualizar/ -run TestShaDivergenteMantemBinario`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/atualizar/atualizar.go') && file_exists('internal/atualizar/atualizar_test.go') && file_exists('internal/cli/atualizar.go')`
* **Spec Ref:** `AC-E11.1`

#### Ticket `TICK-11.2` — Guardar o binário anterior e revalidar a Forja depois de atualizar
* **Iniciativa:** `I01` | **Épico:** `E11` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-11.1","SG-E10"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/atualizar/prev.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md#AC-E11.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/atualizar/prev_test.go`
* **Outputs:** Artifacts `internal/atualizar/prev.go, internal/atualizar/prev_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/atualizar/ -run TestGuardaPrevERevalida`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/atualizar/prev.go') && file_exists('internal/atualizar/prev_test.go')`
* **Spec Ref:** `AC-E11.2`

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
*Plano de execução fatiado com 38 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
