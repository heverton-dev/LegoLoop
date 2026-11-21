---
id: PLAN-2026-11-21_ROUTE_INICIATIVA_2_DO_MACROPROJETO_D
status: DRAFT
date: 2026-11-21
schema_version: 2
upstream: [2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.5, AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E04.1, AC-E04.2, AC-E04.3, AC-E05.1, AC-E05.2, AC-E05.3, AC-E06.1, AC-E06.2, AC-E07.1, AC-E07.2, AC-E07.3, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Processos de especificação e de planejamento (I02)

> **Status:** Em Planejamento  
> **Data:** 2026-11-21  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md](../intents/2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md)  
> **Referência da Especificação:** [docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md](../specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-21_ROUTE_INICIATIVA_2_DO_MACROPROJETO_D.json](../routes/2026-11-21_ROUTE_INICIATIVA_2_DO_MACROPROJETO_D.json)  
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
| **E01** | Pacote `internal/routedag` | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | `fabricar plano` | 5 | `AC-E02.1`, `AC-NF.1`, `AC-E02.2`, `AC-E02.3`, `AC-E02.4`, `AC-E02.5` | `micro` |
| **E03** | Processo de especificação | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-E03.4` | `micro` |
| **E04** | Processo de planejamento | 2 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3` | `micro` |
| **E05** | Triagem conhece os Processos novos | 3 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3` | `micro` |
| **E06** | Processo composto por `subprocesso` | 2 | `AC-E06.1`, `AC-E06.2` | `micro` |
| **E07** | Governança viva | 3 | `AC-E07.1`, `AC-E07.2`, `AC-E07.3`, `AC-NF.2` | `nano` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Pacote `internal/routedag`

#### Ticket `TICK-01.1` — Criar internal/routedag com Projetar, Converter e VerificarCadeia e migrar os testes do route2dag
* **Iniciativa:** `I02` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/routedag/routedag.go, internal/routedag/cadeia.go, internal/routedag/convert.go, internal/routedag/node.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/routedag/routedag_test.go`
* **Outputs:** Artifacts `internal/routedag/routedag.go, internal/routedag/cadeia.go, internal/routedag/convert.go, internal/routedag/node.go, internal/routedag/routedag_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/routedag/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/routedag/routedag.go') && file_exists('internal/routedag/cadeia.go') && file_exists('internal/routedag/convert.go') && file_exists('internal/routedag/node.go') && file_exists('internal/routedag/routedag_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Criar routedag.Executar e reduzir tools/route2dag/main.go a uma chamada
* **Iniciativa:** `I02` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/routedag/executar.go, tools/route2dag/main.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/routedag/executar_test.go`
* **Outputs:** Artifacts `internal/routedag/executar.go, internal/routedag/executar_test.go, tools/route2dag/main.go`
* **Verificação (argv, sem shell):** `go test ./internal/routedag/ ./tools/route2dag/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/routedag/executar.go') && file_exists('internal/routedag/executar_test.go') && file_exists('tools/route2dag/main.go')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.3` — Provar que o DAG da Forja sai idêntico depois da mudança
* **Iniciativa:** `I02` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/routedag/forja_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/routedag/forja_test.go`
* **Outputs:** Artifacts `internal/routedag/forja_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/routedag/ -run TestDagDaForjaIdentico`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/routedag/forja_test.go')`
* **Spec Ref:** `AC-E01.3`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — `fabricar plano`

#### Ticket `TICK-02.1` — Criar fabricar plano que grava PLAN, ROUTE e DAG em um comando
* **Iniciativa:** `I02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_plano.go, internal/cli/fabricar_entrevista.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_plano_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_plano.go, internal/cli/fabricar_plano_test.go, internal/cli/fabricar_entrevista.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarPlano`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_plano.go') && file_exists('internal/cli/fabricar_plano_test.go') && file_exists('internal/cli/fabricar_entrevista.go')`
* **Spec Ref:** `AC-E02.1, AC-NF.1`

#### Ticket `TICK-02.2` — Recusar sem INTENT aprovada com exit 2 e sem gravar nada
* **Iniciativa:** `I02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_plano.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_plano_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_plano.go, internal/cli/fabricar_plano_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarPlanoRecusa`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_plano.go') && file_exists('internal/cli/fabricar_plano_test.go')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.3` — Tornar o modo estrito o padrão e criar o flag --permissivo
* **Iniciativa:** `I02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_plano.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E02.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_plano_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_plano.go, internal/cli/fabricar_plano_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarPlanoEstrito`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_plano.go') && file_exists('internal/cli/fabricar_plano_test.go')`
* **Spec Ref:** `AC-E02.3`

#### Ticket `TICK-02.4` — Gravar PLAN, ROUTE e DAG de forma atômica: falha em qualquer etapa não deixa arquivo
* **Iniciativa:** `I02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_plano.go, internal/cli/fabricar_plano_atomico.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E02.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_plano_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_plano.go, internal/cli/fabricar_plano_atomico.go, internal/cli/fabricar_plano_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarPlanoAtomico`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_plano.go') && file_exists('internal/cli/fabricar_plano_atomico.go') && file_exists('internal/cli/fabricar_plano_test.go')`
* **Spec Ref:** `AC-E02.4`

#### Ticket `TICK-02.5` — Remover legoloop plano e trocar toda referência por fabricar plano
* **Iniciativa:** `I02` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.4","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/plano.go, internal/cli/cli.go, AGENTS.md, README.md, docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E02.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/plano_test.go`
* **Outputs:** Artifacts `internal/cli/plano.go, internal/cli/plano_test.go, internal/cli/cli.go, AGENTS.md, README.md, docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md`
* **Verificação (argv, sem shell):** `go test ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/plano.go') && file_exists('internal/cli/plano_test.go') && file_exists('internal/cli/cli.go') && file_exists('AGENTS.md') && file_exists('README.md') && file_exists('docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md')`
* **Spec Ref:** `AC-E02.5`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Processo de especificação

#### Ticket `TICK-03.1` — Criar fabricar spec que grava a SPEC a partir de uma INTENT aprovada
* **Iniciativa:** `I02` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.5","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_spec.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_spec_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_spec.go, internal/cli/fabricar_spec_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarSpec`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_spec.go') && file_exists('internal/cli/fabricar_spec_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Recusar INTENT não aprovada com exit 2 sem chamar o LLM
* **Iniciativa:** `I02` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/fabricar_spec.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/fabricar_spec_test.go`
* **Outputs:** Artifacts `internal/cli/fabricar_spec.go, internal/cli/fabricar_spec_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestFabricarSpecRecusa`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fabricar_spec.go') && file_exists('internal/cli/fabricar_spec_test.go')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `TICK-03.3` — Criar o Processo legoloop-especificacao com contrato e nascimento 5 de 5
* **Iniciativa:** `I02` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-especificacao/pipeline.dag.json, processos/legoloop-especificacao/triagem.json, processos/legoloop-especificacao/memoria_execucao.json, processos/legoloop-especificacao/anticorpos.json, processos/legoloop-especificacao/ANTICORPOS.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E03.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-especificacao/pipeline.dag.json, processos/legoloop-especificacao/triagem.json, processos/legoloop-especificacao/memoria_execucao.json, processos/legoloop-especificacao/anticorpos.json, processos/legoloop-especificacao/ANTICORPOS.md`
* **Verificação (argv, sem shell):** `legoloop validate processos/legoloop-especificacao/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-especificacao/pipeline.dag.json') && file_exists('processos/legoloop-especificacao/triagem.json') && file_exists('processos/legoloop-especificacao/memoria_execucao.json') && file_exists('processos/legoloop-especificacao/anticorpos.json') && file_exists('processos/legoloop-especificacao/ANTICORPOS.md')`
* **Spec Ref:** `AC-E03.3, AC-E03.4`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Processo de planejamento

#### Ticket `TICK-04.1` — Criar o Processo legoloop-planejamento com contrato e nascimento 5 de 5
* **Iniciativa:** `I02` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-planejamento/pipeline.dag.json, processos/legoloop-planejamento/triagem.json, processos/legoloop-planejamento/memoria_execucao.json, processos/legoloop-planejamento/anticorpos.json, processos/legoloop-planejamento/ANTICORPOS.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-planejamento/pipeline.dag.json, processos/legoloop-planejamento/triagem.json, processos/legoloop-planejamento/memoria_execucao.json, processos/legoloop-planejamento/anticorpos.json, processos/legoloop-planejamento/ANTICORPOS.md`
* **Verificação (argv, sem shell):** `legoloop validate processos/legoloop-planejamento/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-planejamento/pipeline.dag.json') && file_exists('processos/legoloop-planejamento/triagem.json') && file_exists('processos/legoloop-planejamento/memoria_execucao.json') && file_exists('processos/legoloop-planejamento/anticorpos.json') && file_exists('processos/legoloop-planejamento/ANTICORPOS.md')`
* **Spec Ref:** `AC-E04.1, AC-E04.2`

#### Ticket `TICK-04.2` — Fazer o DAG do planejamento chamar fabricar plano em vez de duplicar a lógica
* **Iniciativa:** `I02` | **Épico:** `E04` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-planejamento/pipeline.dag.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-planejamento/pipeline.dag.json`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "legoloop fabricar plano" processos/legoloop-planejamento/pipeline.dag.json')`
* **Spec Ref:** `AC-E04.3`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Triagem conhece os Processos novos

#### Ticket `TICK-05.1` — Rotear gerar a SPEC para legoloop-especificacao sem LLM
* **Iniciativa:** `I02` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-especificacao/triagem.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/triagem/especificacao_test.go`
* **Outputs:** Artifacts `internal/triagem/especificacao_test.go, processos/legoloop-especificacao/triagem.json`
* **Verificação (argv, sem shell):** `go test ./internal/triagem/ -run TestTriagemRoteiaEspecificacao`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/triagem/especificacao_test.go') && file_exists('processos/legoloop-especificacao/triagem.json')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Rotear gerar o PLAN, a ROUTE e o DAG para legoloop-planejamento sem LLM
* **Iniciativa:** `I02` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-planejamento/triagem.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E05.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/triagem/planejamento_test.go`
* **Outputs:** Artifacts `internal/triagem/planejamento_test.go, processos/legoloop-planejamento/triagem.json`
* **Verificação (argv, sem shell):** `go test ./internal/triagem/ -run TestTriagemRoteiaPlanejamento`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/triagem/planejamento_test.go') && file_exists('processos/legoloop-planejamento/triagem.json')`
* **Spec Ref:** `AC-E05.2`

#### Ticket `TICK-05.3` — Provar que a triagem não roubou entradas dos Processos existentes
* **Iniciativa:** `I02` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/triagem/triagem_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E05.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/triagem/triagem_test.go`
* **Outputs:** Artifacts `internal/triagem/triagem_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/triagem/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/triagem/triagem_test.go')`
* **Spec Ref:** `AC-E05.3`

### 🛑 STOP-GATE HUMANO (`SG-E05`)
* **Gate (DSL):** `prompt('Evidências do Épico 05 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 06 — Processo composto por `subprocesso`

#### Ticket `TICK-06.1` — Criar legoloop-ciclo e compor concepção, especificação e planejamento por subprocesso
* **Iniciativa:** `I02` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.3","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-ciclo/pipeline.dag.json, processos/legoloop-ciclo/triagem.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E06.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-ciclo/pipeline.dag.json, processos/legoloop-ciclo/triagem.json`
* **Verificação (argv, sem shell):** `legoloop validate processos/legoloop-ciclo/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-ciclo/pipeline.dag.json') && file_exists('processos/legoloop-ciclo/triagem.json')`
* **Spec Ref:** `AC-E06.1`

#### Ticket `TICK-06.2` — Provar que o produces de cada subprocesso casa com o consumes do seguinte
* **Iniciativa:** `I02` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/ciclo_contrato_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E06.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/ciclo_contrato_test.go`
* **Outputs:** Artifacts `internal/cli/ciclo_contrato_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestCicloContratoEncadeia`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/ciclo_contrato_test.go')`
* **Spec Ref:** `AC-E06.2`

### 🛑 STOP-GATE HUMANO (`SG-E06`)
* **Gate (DSL):** `prompt('Evidências do Épico 06 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 07 — Governança viva

#### Ticket `TICK-07.1` — Atualizar AGENTS, README e MEMORY com os Processos novos
* **Iniciativa:** `I02` | **Épico:** `E07` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.2","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E07.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "fabricar plano" AGENTS.md README.md MEMORY.md')`
* **Spec Ref:** `AC-E07.1`

#### Ticket `TICK-07.2` — Registrar os anticorpos da I02 nos Processos novos
* **Iniciativa:** `I02` | **Épico:** `E07` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-especificacao/anticorpos.json, processos/legoloop-especificacao/ANTICORPOS.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E07.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-especificacao/anticorpos.json, processos/legoloop-especificacao/ANTICORPOS.md`
* **Verificação (argv, sem shell):** `legoloop anticorpos status legoloop-especificacao`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-especificacao/anticorpos.json') && file_exists('processos/legoloop-especificacao/ANTICORPOS.md')`
* **Spec Ref:** `AC-E07.2`

#### Ticket `TICK-07.3` — Fechar a Emenda 1, gerar o prompt de execução da I02 e rodar o portão completo
* **Iniciativa:** `I02` | **Épico:** `E07` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.2","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md, docs/prompts` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md#AC-E07.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md, docs/prompts`
* **Verificação (argv, sem shell):** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md') && file_exists('docs/prompts')`
* **Spec Ref:** `AC-E07.3, AC-NF.2`

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
*Plano de execução fatiado com 21 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
