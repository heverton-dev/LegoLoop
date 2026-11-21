---
id: PLAN-2026-11-21_ROUTE_INICIATIVA_4_I04_DO_MACROPROJE
status: DRAFT
date: 2026-11-21
schema_version: 2
upstream: [2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.5, AC-E02.6, AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E03.5, AC-E03.6, AC-E03.7, AC-E03.8, AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.4, AC-E04.5, AC-E04.6, AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-E05.5, AC-E05.6, AC-E05.7, AC-E06.1, AC-E06.2, AC-E06.3, AC-E06.4, AC-E06.5, AC-E06.6, AC-E06.7, AC-E06.8, AC-E07.1, AC-E07.2, AC-E07.3, AC-E07.4, AC-E07.5, AC-E07.6, AC-E07.7, AC-E07.8, AC-E07.9, AC-E07.10, AC-E08.1, AC-E08.2, AC-E08.3, AC-E08.4, AC-E08.5, AC-E09.1, AC-E09.2, AC-E09.3, AC-E09.4, AC-E09.5, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Espinha dorsal previsível, universal e determinística (I04)

> **Status:** Em Planejamento  
> **Data:** 2026-11-21  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md](../intents/2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md](../specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-21_ROUTE_INICIATIVA_4_I04_DO_MACROPROJE.json](../routes/2026-11-21_ROUTE_INICIATIVA_4_I04_DO_MACROPROJE.json)  
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
| **E01** | Linha de base executável | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Vocabulário e executor (G03, G04) | 6 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3`, `AC-E02.4`, `AC-E02.5`, `AC-E02.6` | `micro` |
| **E03** | Fachada de caminhos, nomes de Processo e código morto (G05, G06, G17) | 8 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-E03.4`, `AC-E03.5`, `AC-E03.6`, `AC-E03.7`, `AC-E03.8` | `micro` |
| **E04** | Descoberta de peças e Matriz (G01, G02) | 6 | `AC-E04.1`, `AC-E04.2`, `AC-E04.3`, `AC-E04.4`, `AC-E04.5`, `AC-E04.6` | `micro` |
| **E05** | Registro da Espinha e gate | 7 | `AC-E05.1`, `AC-E05.2`, `AC-E05.3`, `AC-E05.4`, `AC-E05.5`, `AC-E05.6`, `AC-E05.7`, `AC-NF.1` | `micro` |
| **E06** | Dívidas médias e baixas (G08 a G14) | 8 | `AC-E06.1`, `AC-E06.2`, `AC-E06.3`, `AC-E06.4`, `AC-E06.5`, `AC-E06.6`, `AC-E06.7`, `AC-E06.8` | `micro` |
| **E07** | CI e universalidade (G07, G13) | 8 | `AC-E07.6`, `AC-E07.7`, `AC-E07.8`, `AC-E07.9`, `AC-E07.10`, `AC-E07.4`, `AC-E07.1`, `AC-E07.2`, `AC-E07.3`, `AC-E07.5` | `micro` |
| **E08** | Mutação | 5 | `AC-E08.1`, `AC-E08.2`, `AC-E08.3`, `AC-E08.4`, `AC-E08.5` | `micro` |
| **E09** | Revisão das iniciativas congeladas e fechamento | 5 | `AC-E09.1`, `AC-E09.2`, `AC-E09.3`, `AC-E09.4`, `AC-E09.5`, `AC-NF.2` | `nano` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Linha de base executável

#### Ticket `TICK-01.1` — Registrar a correção de G05 e as verificações novas no relatório de auditoria
* **Iniciativa:** `I04` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "Correção" docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Provar que legoloop run duas vezes gera relatório idêntico
* **Iniciativa:** `I04` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/run_determinismo_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/run_determinismo_test.go`
* **Outputs:** Artifacts `internal/cli/run_determinismo_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestRunScriptsDuasVezesRelatorioIdentico`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/run_determinismo_test.go')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.3` — Provar que legoloop exec duas vezes grava artefatos com o mesmo hash
* **Iniciativa:** `I04` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/exec_determinismo_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E01.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/exec_determinismo_test.go`
* **Outputs:** Artifacts `internal/cli/exec_determinismo_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestExecDuasVezesRelatorioIdentico`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/exec_determinismo_test.go')`
* **Spec Ref:** `AC-E01.3`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Vocabulário e executor (G03, G04)

#### Ticket `TICK-02.1` — Criar a tabela única unit_task e unit_no e usá-la na route2dag e no runner
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/tipos.go, tools/route2dag/convert.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/tipos_test.go`
* **Outputs:** Artifacts `internal/dag/tipos.go, internal/dag/tipos_test.go, tools/route2dag/convert.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestTabelaDeTipos`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/tipos.go') && file_exists('internal/dag/tipos_test.go') && file_exists('tools/route2dag/convert.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Remover o alias TypeSubloop e o literal unit_no do runner
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/types.go, internal/runner/runner.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/dag/types.go, internal/runner/runner.go`
* **Verificação (argv, sem shell):** `go build ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/types.go') && file_exists('internal/runner/runner.go')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.3` — Derivar o executor ausente pelo tipo no dag.Parse
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor.go, internal/dag/parse.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/executor.go, internal/dag/executor_test.go, internal/dag/parse.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestParseDerivaOExecutor`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/executor.go') && file_exists('internal/dag/executor_test.go') && file_exists('internal/dag/parse.go')`
* **Spec Ref:** `AC-E02.3`

#### Ticket `TICK-02.4` — Recusar executor incoerente com todos os nós listados
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.3","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/dag/executor.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_test.go`
* **Outputs:** Artifacts `internal/dag/executor.go, internal/dag/executor_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestExecutorIncoerenteListaTodos`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/dag/executor.go') && file_exists('internal/dag/executor_test.go')`
* **Spec Ref:** `AC-E02.4`

#### Ticket `TICK-02.5` — Corrigir os rótulos de executor dos DAGs do repositório
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.4","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos, docs/dags` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/dag/executor_repo_test.go`
* **Outputs:** Artifacts `processos, docs/dags, internal/dag/executor_repo_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/dag/ -run TestDAGsDoRepositorioValidam`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos') && file_exists('docs/dags') && file_exists('internal/dag/executor_repo_test.go')`
* **Spec Ref:** `AC-E02.5`

#### Ticket `TICK-02.6` — Fazer validate e run recusarem DAG incoerente com exit 2
* **Iniciativa:** `I04` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.5","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/commands.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E02.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/validate_executor_test.go`
* **Outputs:** Artifacts `internal/cli/validate_executor_test.go, internal/cli/commands.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestValidateERunRecusamIncoerenteComExit2`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/validate_executor_test.go') && file_exists('internal/cli/commands.go')`
* **Spec Ref:** `AC-E02.6`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Fachada de caminhos, nomes de Processo e código morto (G05, G06, G17)

#### Ticket `TICK-03.1` — Mover as constantes de paths para caminhos e renomear o jail para internal/jail
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.6","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/caminhos/nomes.go, internal/paths` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/caminhos/constantes_test.go`
* **Outputs:** Artifacts `internal/caminhos/nomes.go, internal/caminhos/constantes_test.go, internal/paths`
* **Verificação (argv, sem shell):** `go test ./internal/caminhos/ -run TestConstantesEFuncoesNoMesmoPacote`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/caminhos/nomes.go') && file_exists('internal/caminhos/constantes_test.go') && file_exists('internal/paths')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Migrar os chamadores e apagar as funções delegadas de paths
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/migrar/migrar.go, internal/refinar/refinar.go, internal/fabricar, internal/processo, internal/cli` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/migrar/migrar.go, internal/refinar/refinar.go, internal/fabricar, internal/processo, internal/cli`
* **Verificação (argv, sem shell):** `go build ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/migrar/migrar.go') && file_exists('internal/refinar/refinar.go') && file_exists('internal/fabricar') && file_exists('internal/processo') && file_exists('internal/cli')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `TICK-03.3` — Ligar o jail aos verbos que recebem nome de Processo
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/jail, internal/cli/commands.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/nome_hostil_test.go`
* **Outputs:** Artifacts `internal/jail, internal/cli/nome_hostil_test.go, internal/cli/commands.go`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestNomeHostilEhRecusado`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/jail') && file_exists('internal/cli/nome_hostil_test.go') && file_exists('internal/cli/commands.go')`
* **Spec Ref:** `AC-E03.3`

#### Ticket `TICK-03.4` — Remover o pacote vault sem importador em produção
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.3","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/vault` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/vault`
* **Verificação (argv, sem shell):** `go build ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/vault')`
* **Spec Ref:** `AC-E03.4`

#### Ticket `TICK-03.5` — Criar o verbo legoloop renomear com dry-run, backup, colisão e registro
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.4","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/renomear, internal/cli/renomear.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/renomear/renomear_test.go`
* **Outputs:** Artifacts `internal/renomear, internal/cli/renomear.go, internal/renomear/renomear_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/renomear/ -run TestRenomear`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/renomear') && file_exists('internal/cli/renomear.go') && file_exists('internal/renomear/renomear_test.go')`
* **Spec Ref:** `AC-E03.5`

#### Ticket `TICK-03.6` — Renomear os 9 Processos sem colisão com o verbo renomear
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.5","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos, espinha/renomeados.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos, espinha/renomeados.json`
* **Verificação (argv, sem shell):** `legoloop validate processos/legoloop-listar/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos') && file_exists('espinha/renomeados.json')`
* **Spec Ref:** `AC-E03.6`

#### Ticket `TICK-03.7` — Aposentar o Processo fechar e passar os anticorpos para legoloop-fechar
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.6","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/fechar, processos/legoloop-fechar, espinha/renomeados.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.7` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/fechar, processos/legoloop-fechar, espinha/renomeados.json`
* **Verificação (argv, sem shell):** `legoloop anticorpos status legoloop-fechar`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/fechar') && file_exists('processos/legoloop-fechar') && file_exists('espinha/renomeados.json')`
* **Spec Ref:** `AC-E03.7`

#### Ticket `TICK-03.8` — Provar que nenhum arquivo cita os nomes antigos dos Processos
* **Iniciativa:** `I04` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.7","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/nomes_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E03.8` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/nomes_test.go`
* **Outputs:** Artifacts `internal/espinha/nomes_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestSemReferenciaAoNomeAntigo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/nomes_test.go')`
* **Spec Ref:** `AC-E03.8`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Descoberta de peças e Matriz (G01, G02)

#### Ticket `TICK-04.1` — Criar a descoberta de unidades fixas por SO com provedor injetável
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.8","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/descoberta/unidades.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/descoberta/unidades_test.go`
* **Outputs:** Artifacts `internal/descoberta/unidades.go, internal/descoberta/unidades_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/descoberta/ -run TestUnidadesFixasPorSO`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/descoberta/unidades.go') && file_exists('internal/descoberta/unidades_test.go')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.2` — Fazer o TargetDrives restringir a busca e remover o aviso drives_ignorados
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fabricar/engine.go, internal/assets/discover.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/fabricar/drives_test.go`
* **Outputs:** Artifacts `internal/fabricar/engine.go, internal/fabricar/drives_test.go, internal/assets/discover.go`
* **Verificação (argv, sem shell):** `go test ./internal/fabricar/ -run TestTargetDrivesRestringeABusca`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fabricar/engine.go') && file_exists('internal/fabricar/drives_test.go') && file_exists('internal/assets/discover.go')`
* **Spec Ref:** `AC-E04.2`

#### Ticket `TICK-04.3` — Fazer o fabricar atualizar a Matriz defasada em modo sem cérebro
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fabricar/engine.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/fabricar/atualiza_matriz_test.go`
* **Outputs:** Artifacts `internal/fabricar/engine.go, internal/fabricar/atualiza_matriz_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/fabricar/ -run TestFabricarAtualizaMatrizDefasada`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fabricar/engine.go') && file_exists('internal/fabricar/atualiza_matriz_test.go')`
* **Spec Ref:** `AC-E04.3`

#### Ticket `TICK-04.4` — Provar que Matriz em dia não dispara atualização
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.3","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fabricar/atualiza_matriz_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/fabricar/atualiza_matriz_test.go`
* **Outputs:** Artifacts `internal/fabricar/atualiza_matriz_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/fabricar/ -run TestFabricarNaoAtualizaMatrizEmDia`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fabricar/atualiza_matriz_test.go')`
* **Spec Ref:** `AC-E04.4`

#### Ticket `TICK-04.5` — Provar que a descoberta é determinística e somente leitura, com relatório
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.4","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/descoberta/leitura_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/descoberta/leitura_test.go`
* **Outputs:** Artifacts `internal/descoberta/leitura_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/descoberta/ -run TestDescobertaSomenteLeitura`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/descoberta/leitura_test.go')`
* **Spec Ref:** `AC-E04.5`

#### Ticket `TICK-04.6` — Provar que matriz status fica sem diferença depois do fabricar
* **Iniciativa:** `I04` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.5","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fabricar/atualiza_matriz_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E04.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/fabricar/atualiza_matriz_test.go`
* **Outputs:** Artifacts `internal/fabricar/atualiza_matriz_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/fabricar/ -run TestMatrizStatusZeroDepoisDoFabricar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fabricar/atualiza_matriz_test.go')`
* **Spec Ref:** `AC-E04.6`

### 🛑 STOP-GATE HUMANO (`SG-E04`)
* **Gate (DSL):** `prompt('Evidências do Épico 04 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 05 — Registro da Espinha e gate

#### Ticket `TICK-05.1` — Criar o Registro da Espinha e o schema espinha.v1
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.6","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/registro.go, espinha/registro.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/registro_test.go`
* **Outputs:** Artifacts `internal/espinha/registro.go, internal/espinha/registro_test.go, espinha/registro.json`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestRegistroValida`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/registro.go') && file_exists('internal/espinha/registro_test.go') && file_exists('espinha/registro.json')`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Exportar cli.Verbos e provar que todo verbo está no Registro
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/cli.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/verbos_test.go`
* **Outputs:** Artifacts `internal/cli/cli.go, internal/espinha/verbos_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestTodoVerboEstaNoRegistro`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/cli.go') && file_exists('internal/espinha/verbos_test.go')`
* **Spec Ref:** `AC-E05.2`

#### Ticket `TICK-05.3` — Ligar cada lei do AGENTS.md a um verificador ou justificativa
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `espinha/registro.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/leis_test.go`
* **Outputs:** Artifacts `internal/espinha/leis_test.go, espinha/registro.json`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestTodaLeiTemVerificador`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/leis_test.go') && file_exists('espinha/registro.json')`
* **Spec Ref:** `AC-E05.3`

#### Ticket `TICK-05.4` — Criar legoloop espinha verificar com a conferência de motores e Processos
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.3","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/verificar.go, internal/cli/espinha.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/espinha/verificar.go, internal/cli/espinha.go`
* **Verificação (argv, sem shell):** `legoloop espinha verificar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/verificar.go') && file_exists('internal/cli/espinha.go')`
* **Spec Ref:** `AC-E05.4`

#### Ticket `TICK-05.5` — Detectar função duplicada entre pacotes
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.4","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/duplicadas.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/duplicadas_test.go`
* **Outputs:** Artifacts `internal/espinha/duplicadas.go, internal/espinha/duplicadas_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestFuncaoDuplicadaEntrePacotes`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/duplicadas.go') && file_exists('internal/espinha/duplicadas_test.go')`
* **Spec Ref:** `AC-E05.5`

#### Ticket `TICK-05.6` — Provar as doze regras com fixture com defeito e exit 2
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.5","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/regras_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/regras_test.go`
* **Outputs:** Artifacts `internal/espinha/regras_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestCadaRegraSaiComExitDois`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/regras_test.go')`
* **Spec Ref:** `AC-E05.6`

#### Ticket `TICK-05.7` — Rodar o gate no fechar e provar que a espinha não usa LLM nem relógio
* **Iniciativa:** `I04` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.6","SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fechar/fechar.go, internal/espinha` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E05.7` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/fechar/fechar.go, internal/espinha`
* **Verificação (argv, sem shell):** `go test ./internal/fechar/ -run TestFecharRodaEspinhaVerificar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fechar/fechar.go') && file_exists('internal/espinha')`
* **Spec Ref:** `AC-E05.7, AC-NF.1`

### 🛑 STOP-GATE HUMANO (`SG-E05`)
* **Gate (DSL):** `prompt('Evidências do Épico 05 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 06 — Dívidas médias e baixas (G08 a G14)

#### Ticket `TICK-06.1` — Completar o nascimento 5 de 5 dos três Processos sem ANTICORPOS.md
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.7","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos/legoloop-agent-desktop/ANTICORPOS.md, processos/legoloop-matriz/ANTICORPOS.md, processos/legoloop-listar/ANTICORPOS.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `processos/legoloop-agent-desktop/ANTICORPOS.md, processos/legoloop-matriz/ANTICORPOS.md, processos/legoloop-listar/ANTICORPOS.md`
* **Verificação (argv, sem shell):** `legoloop espinha verificar`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/legoloop-agent-desktop/ANTICORPOS.md') && file_exists('processos/legoloop-matriz/ANTICORPOS.md') && file_exists('processos/legoloop-listar/ANTICORPOS.md')`
* **Spec Ref:** `AC-E06.1`

#### Ticket `TICK-06.2` — Remover o verbo matrix descontinuado
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/cli.go, internal/cli/matrix.go, README.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/cli/cli.go, internal/cli/matrix.go, README.md`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ -run TestVerboMatrixRemovido`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/cli.go') && file_exists('internal/cli/matrix.go') && file_exists('README.md')`
* **Spec Ref:** `AC-E06.2`

#### Ticket `TICK-06.3` — Criar tools/cobertura e subir board, triagem e cli a 80 por cento
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.2","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/cobertura, espinha/cobertura.json, internal/board, internal/triagem, internal/cli` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `tools/cobertura, espinha/cobertura.json, internal/board, internal/triagem, internal/cli`
* **Verificação (argv, sem shell):** `go run ./tools/cobertura`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/cobertura') && file_exists('espinha/cobertura.json') && file_exists('internal/board') && file_exists('internal/triagem') && file_exists('internal/cli')`
* **Spec Ref:** `AC-E06.3`

#### Ticket `TICK-06.4` — Criar escrita.GravarAtomico e usá-la nas gravações de estado crítico
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.3","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/escrita, internal/fechar, internal/brain, internal/cli` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/escrita, internal/fechar, internal/brain, internal/cli`
* **Verificação (argv, sem shell):** `go test ./internal/escrita/... ./internal/fechar/... ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/escrita') && file_exists('internal/fechar') && file_exists('internal/brain') && file_exists('internal/cli')`
* **Spec Ref:** `AC-E06.4`

#### Ticket `TICK-06.5` — Relatar a falha de telemetria e fixar o limite de erros ignorados
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.4","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/telemetria.go, espinha/registro.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/telemetria_test.go`
* **Outputs:** Artifacts `internal/runner/telemetria.go, internal/runner/telemetria_test.go, espinha/registro.json`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestFalhaDeTelemetriaEhRelatada`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/telemetria.go') && file_exists('internal/runner/telemetria_test.go') && file_exists('espinha/registro.json')`
* **Spec Ref:** `AC-E06.5`

#### Ticket `TICK-06.6` — Provar que os globais e init da espinha estão justificados
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.5","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `espinha/registro.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/globais_test.go`
* **Outputs:** Artifacts `internal/espinha/globais_test.go, espinha/registro.json`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestGlobaisDaEspinhaEstaoJustificados`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/globais_test.go') && file_exists('espinha/registro.json')`
* **Spec Ref:** `AC-E06.6`

#### Ticket `TICK-06.7` — Corrigir o relato do Andon quando o nó é recusado antes de executar
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.6","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/runner/sbar.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.7` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/runner/sbar_test.go`
* **Outputs:** Artifacts `internal/runner/sbar.go, internal/runner/sbar_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/runner/ -run TestSBARDeRecusaPreExecucao`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/runner/sbar.go') && file_exists('internal/runner/sbar_test.go')`
* **Spec Ref:** `AC-E06.7`

#### Ticket `TICK-06.8` — Remover o arquivo temporário rastreado em scratch
* **Iniciativa:** `I04` | **Épico:** `E06` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.7","SG-E05"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `scratch/i6_inventario_legados.json` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E06.8` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `scratch/i6_inventario_legados.json`
* **Verificação (argv, sem shell):** `git ls-files scratch`
* **Gate (DSL):** `exit_code == 0 && file_exists('scratch/i6_inventario_legados.json')`
* **Spec Ref:** `AC-E06.8`

### 🛑 STOP-GATE HUMANO (`SG-E06`)
* **Gate (DSL):** `prompt('Evidências do Épico 06 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 07 — CI e universalidade (G07, G13)

#### Ticket `TICK-07.1` — Fazer o jail recusar caminho hostil em qualquer SO
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.8","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/jail/jail.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.6` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/jail/jail_test.go`
* **Outputs:** Artifacts `internal/jail/jail.go, internal/jail/jail_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/jail/ -run TestEntradaHostil`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/jail/jail.go') && file_exists('internal/jail/jail_test.go')`
* **Spec Ref:** `AC-E07.6`

#### Ticket `TICK-07.2` — Tornar o lock da memória da triagem exclusivo em Linux
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/triagem/memoria.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.7` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/triagem/memoria_test.go`
* **Outputs:** Artifacts `internal/triagem/memoria.go, internal/triagem/memoria_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/triagem/ -run TestMemoriaLockOcupado`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/triagem/memoria.go') && file_exists('internal/triagem/memoria_test.go')`
* **Spec Ref:** `AC-E07.7`

#### Ticket `TICK-07.3` — Normalizar os separadores na guarda de deleção
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.2","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/inventario/guard.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.8` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/inventario/guard_test.go`
* **Outputs:** Artifacts `internal/inventario/guard.go, internal/inventario/guard_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/inventario/ -run TestRecusaDelecaoTemplate`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/inventario/guard.go') && file_exists('internal/inventario/guard_test.go')`
* **Spec Ref:** `AC-E07.8`

#### Ticket `TICK-07.4` — Tornar herméticos os seis testes que dependem de ferramenta instalada
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.3","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli, internal/legopress` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.9` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `internal/cli, internal/legopress`
* **Verificação (argv, sem shell):** `go test ./internal/cli/ ./internal/legopress/ -count=1`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli') && file_exists('internal/legopress')`
* **Spec Ref:** `AC-E07.9`

#### Ticket `TICK-07.5` — Trocar powershell e python direto por iniciador portável nos Processos
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.4","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `processos` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.10` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/so_test.go`
* **Outputs:** Artifacts `processos, internal/espinha/so_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestProcessosSemComandoExclusivoDoSO`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos') && file_exists('internal/espinha/so_test.go')`
* **Spec Ref:** `AC-E07.10`

#### Ticket `TICK-07.6` — Resolver fonte, python e shell por SO
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.5","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/assets/argv.go, internal/videorea/cartoon_render.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/assets/argv_test.go`
* **Outputs:** Artifacts `internal/assets/argv.go, internal/videorea/cartoon_render.go, internal/assets/argv_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/assets/ ./internal/videorea/ -run TestPorSO`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/assets/argv.go') && file_exists('internal/videorea/cartoon_render.go') && file_exists('internal/assets/argv_test.go')`
* **Spec Ref:** `AC-E07.4`

#### Ticket `TICK-07.7` — Criar o workflow de CI nos três SOs com os gates e go test -race
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.6","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `.github/workflows/ci.yml` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `.github/workflows/ci.yml`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -cE "windows-latest|ubuntu-latest|macos-latest" .github/workflows/ci.yml')`
* **Spec Ref:** `AC-E07.1, AC-E07.2, AC-E07.3`

#### Ticket `TICK-07.8` — Conferir o CI verde nos três SOs
* **Iniciativa:** `I04` | **Épico:** `E07` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.7","SG-E06"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `.github/workflows/ci.yml` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E07.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `.github/workflows/ci.yml`
* **Verificação (argv, sem shell):** `gh run list --workflow ci.yml --limit 1 --json conclusion`
* **Gate (DSL):** `exit_code == 0 && file_exists('.github/workflows/ci.yml')`
* **Spec Ref:** `AC-E07.5`

### 🛑 STOP-GATE HUMANO (`SG-E07`)
* **Gate (DSL):** `prompt('Evidências do Épico 07 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 08 — Mutação

#### Ticket `TICK-08.1` — Criar o tools/mutar só com biblioteca padrão
* **Iniciativa:** `I04` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.8","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/mutar/main.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E08.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/mutar/main_test.go`
* **Outputs:** Artifacts `tools/mutar/main.go, tools/mutar/main_test.go`
* **Verificação (argv, sem shell):** `git diff --exit-code go.mod`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/mutar/main.go') && file_exists('tools/mutar/main_test.go')`
* **Spec Ref:** `AC-E08.1`

#### Ticket `TICK-08.2` — Implementar os sete operadores de mutação
* **Iniciativa:** `I04` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/mutar/operadores.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E08.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/mutar/operadores_test.go`
* **Outputs:** Artifacts `tools/mutar/operadores.go, tools/mutar/operadores_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/mutar/ -run TestOperadores`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/mutar/operadores.go') && file_exists('tools/mutar/operadores_test.go')`
* **Spec Ref:** `AC-E08.2`

#### Ticket `TICK-08.3` — Reprovar com exit 2 quando um mutante sobrevive
* **Iniciativa:** `I04` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.2","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `tools/mutar/rodar.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E08.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `tools/mutar/rodar_test.go`
* **Outputs:** Artifacts `tools/mutar/rodar.go, tools/mutar/rodar_test.go`
* **Verificação (argv, sem shell):** `go test ./tools/mutar/ -run TestSobreviventeReprova`
* **Gate (DSL):** `exit_code == 0 && file_exists('tools/mutar/rodar.go') && file_exists('tools/mutar/rodar_test.go')`
* **Spec Ref:** `AC-E08.3`

#### Ticket `TICK-08.4` — Matar ou justificar todo mutante dos sete pacotes da espinha
* **Iniciativa:** `I04` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.3","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `espinha/mutantes_justificados.json, internal/dag, internal/runner, internal/gate, internal/plangen, internal/triagem, internal/upstream, internal/interview` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E08.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `espinha/mutantes_justificados.json, internal/dag, internal/runner, internal/gate, internal/plangen, internal/triagem, internal/upstream, internal/interview`
* **Verificação (argv, sem shell):** `go run ./tools/mutar ./internal/dag ./internal/runner ./internal/gate ./internal/plangen ./internal/triagem ./internal/upstream ./internal/interview`
* **Gate (DSL):** `exit_code == 0 && file_exists('espinha/mutantes_justificados.json') && file_exists('internal/dag') && file_exists('internal/runner') && file_exists('internal/gate') && file_exists('internal/plangen') && file_exists('internal/triagem') && file_exists('internal/upstream') && file_exists('internal/interview')`
* **Spec Ref:** `AC-E08.4`

#### Ticket `TICK-08.5` — Provar que os mutantes das regras da espinha morrem em 100 por cento
* **Iniciativa:** `I04` | **Épico:** `E08` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.4","SG-E07"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/espinha/mutantes_test.go` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E08.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/espinha/mutantes_test.go`
* **Outputs:** Artifacts `internal/espinha/mutantes_test.go`
* **Verificação (argv, sem shell):** `go test ./internal/espinha/ -run TestMutantesDasRegras`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/espinha/mutantes_test.go')`
* **Spec Ref:** `AC-E08.5`

### 🛑 STOP-GATE HUMANO (`SG-E08`)
* **Gate (DSL):** `prompt('Evidências do Épico 08 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 09 — Revisão das iniciativas congeladas e fechamento

#### Ticket `TICK-09.1` — Revisar a SPEC e o PLAN da I01, da I02 e da I03 pela Emenda 2
* **Iniciativa:** `I04` | **Épico:** `E09` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.5","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/specs, docs/plans, docs/routes, docs/dags` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E09.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/specs, docs/plans, docs/routes, docs/dags`
* **Verificação (argv, sem shell):** `ls docs/specs docs/plans`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/specs') && file_exists('docs/plans') && file_exists('docs/routes') && file_exists('docs/dags')`
* **Spec Ref:** `AC-E09.1`

#### Ticket `TICK-09.2` — Tirar o Épico de CI do plano revisado da I01 e listar G01, G02 e G05 como pré-requisito
* **Iniciativa:** `I04` | **Épico:** `E09` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.1","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/plans` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E09.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/plans`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "AC-E02" docs/plans/2026-11-21_PLAN_DISTRIBUICAO_FORJA_LEGOLOOP_REVISAO.md')`
* **Spec Ref:** `AC-E09.2`

#### Ticket `TICK-09.3` — Manter o congelamento até o Stop-Gate final e removê-lo só com a aprovação do usuário
* **Iniciativa:** `I04` | **Épico:** `E09` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.2","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/prompts` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E09.3` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/prompts`
* **Verificação (argv, sem shell):** nenhuma (confirmação manual pelo gate)
* **Gate (DSL):** `prompt('Execute manualmente e confirme exit code 0: grep -c "EXECUÇÃO CONGELADA" docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md')`
* **Spec Ref:** `AC-E09.3`

#### Ticket `TICK-09.4` — Gravar o relatório de fechamento da I04 com o desfecho das lacunas
* **Iniciativa:** `I04` | **Épico:** `E09` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.3","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `docs/reports` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E09.4` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `docs/reports`
* **Verificação (argv, sem shell):** `ls docs/reports`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports')`
* **Spec Ref:** `AC-E09.4`

#### Ticket `TICK-09.5` — Atualizar os documentos vivos e rodar o portão completo
* **Iniciativa:** `I04` | **Épico:** `E09` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.4","SG-E08"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md` · Spec Anchor `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md#AC-E09.5` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `nenhum`
* **Outputs:** Artifacts `AGENTS.md, docs/governanca/AGENTS_COMPLETO.md, README.md, MEMORY.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('AGENTS.md') && file_exists('docs/governanca/AGENTS_COMPLETO.md') && file_exists('README.md') && file_exists('MEMORY.md')`
* **Spec Ref:** `AC-E09.5, AC-NF.2`

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
*Plano de execução fatiado com 56 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
