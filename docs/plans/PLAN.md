---
id: PLAN-2026-11-24_ROUTE_CRIAR_O_PROCESSO_CANONIZADO_LE
status: DRAFT
date: 2026-11-24
schema_version: 2
upstream: [2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E02.1, AC-E02.2, AC-E03.1, AC-E03.2, AC-E03.3, AC-E04.1, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Processo Canonizado `legoloop-diorama-roteiro` (Ciclo 1 do Macroprocesso)

> **Status:** Em Planejamento  
> **Data:** 2026-11-24  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/2026-11-24_INTENT_CRIAR_O_PROCESSO_CANONIZADO_LE.md](../intents/2026-11-24_INTENT_CRIAR_O_PROCESSO_CANONIZADO_LE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md](../specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-24_ROUTE_CRIAR_O_PROCESSO_CANONIZADO_LE.json](../routes/2026-11-24_ROUTE_CRIAR_O_PROCESSO_CANONIZADO_LE.json)  
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
| **E01** | Épico E01 | 2 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | Épico E02 | 2 | `AC-E02.1`, `AC-E02.2`, `AC-NF.1` | `micro` |
| **E03** | Épico E03 | 2 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3`, `AC-NF.2` | `micro` |
| **E04** | Épico E04 | 1 | `AC-E04.1` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Épico E01

#### Ticket `TICK-01.1` — Criar contratos e validacao estrita de Projeto, Marca e consistencia cruzada
* **Iniciativa:** `D01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/schemas.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E01.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/schemas_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/schemas.go, internal/diorama/roteiro/schemas_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/ -run TestValidarConfiguracao`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/schemas.go') && file_exists('internal/diorama/roteiro/schemas_test.go')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Implementar pre-voo modular com 5 bloqueantes FALTA e itens de P2 e P3 como AVISO
* **Iniciativa:** `D01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/prevoo.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E01.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/prevoo_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/prevoo.go, internal/diorama/roteiro/prevoo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/ -run TestPreVoo`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/prevoo.go') && file_exists('internal/diorama/roteiro/prevoo_test.go')`
* **Spec Ref:** `AC-E01.2, AC-E01.3`

### 🛑 STOP-GATE HUMANO (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 02 — Épico E02

#### Ticket `TICK-02.1` — Criar catalogo declarativo por nicho e classificacao fail-closed por dominio A, B, preprint e C
* **Iniciativa:** `D01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/catalogo.go, internal/diorama/roteiro/catalogo_padrao.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E02.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/catalogo_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/catalogo.go, internal/diorama/roteiro/catalogo_padrao.go, internal/diorama/roteiro/catalogo_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/ -run TestClassificarFonte`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/catalogo.go') && file_exists('internal/diorama/roteiro/catalogo_padrao.go') && file_exists('internal/diorama/roteiro/catalogo_test.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Implementar coletor HTTP gratuito por perfis_fonte, filtro em acervo/temas.json e politica_gates
* **Iniciativa:** `D01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/coletor.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E02.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/coletor_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/coletor.go, internal/diorama/roteiro/coletor_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/ -run TestColetarPautas`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/coletor.go') && file_exists('internal/diorama/roteiro/coletor_test.go')`
* **Spec Ref:** `AC-E02.2, AC-NF.1`

### 🛑 STOP-GATE HUMANO (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 03 — Épico E03

#### Ticket `TICK-03.1` — Implementar validador epistemico universal para as 4 naturezas e bloqueio de nivel C
* **Iniciativa:** `D01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/validador.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E03.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/validador_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/validador.go, internal/diorama/roteiro/validador_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/ -run TestValidarRoteiro`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/validador.go') && file_exists('internal/diorama/roteiro/validador_test.go')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Implementar pipeline com limite de 2 correcoes automaticas e gravacao pos-Gate em acervo/temas.json
* **Iniciativa:** `D01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2","TICK-02.2","TICK-03.1","SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/diorama/roteiro/acervo.go, internal/diorama/roteiro/pipeline.go` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E03.2` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/diorama/roteiro/pipeline_test.go`
* **Outputs:** Artifacts `internal/diorama/roteiro/acervo.go, internal/diorama/roteiro/pipeline.go, internal/diorama/roteiro/pipeline_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/diorama/roteiro/acervo.go') && file_exists('internal/diorama/roteiro/pipeline.go') && file_exists('internal/diorama/roteiro/pipeline_test.go')`
* **Spec Ref:** `AC-E03.2, AC-E03.3, AC-NF.2`

### 🛑 STOP-GATE HUMANO (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza iniciar o próximo épico?')`
- [ ] Tickets anteriores executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para prosseguir.

---

### Épico 04 — Épico E04

#### Ticket `TICK-04.1` — Expor comando CLI legoloop diorama-roteiro, DAG v2, Skill e registro 5 de 5 na Espinha
* **Iniciativa:** `D01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2","SG-E03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `legoloop` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/cli/diorama_roteiro.go, cmd/legoloop/main.go, processos/legoloop-diorama-roteiro/pipeline.dag.json, skills/legoloop-diorama-roteiro/SKILL.md` · Spec Anchor `docs/specs/2026-11-24_SPEC_CRIAR_O_PROCESSO_CANONIZADO_LE.md#AC-E04.1` · Context Artifacts `[]`
* **Alvo de Teste TDD:** `internal/cli/diorama_roteiro_test.go`
* **Outputs:** Artifacts `internal/cli/diorama_roteiro.go, internal/cli/diorama_roteiro_test.go, cmd/legoloop/main.go, processos/legoloop-diorama-roteiro/pipeline.dag.json, skills/legoloop-diorama-roteiro/SKILL.md`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/diorama/roteiro/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/diorama_roteiro.go') && file_exists('internal/cli/diorama_roteiro_test.go') && file_exists('cmd/legoloop/main.go') && file_exists('processos/legoloop-diorama-roteiro/pipeline.dag.json') && file_exists('skills/legoloop-diorama-roteiro/SKILL.md')`
* **Spec Ref:** `AC-E04.1`

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
*Plano de execução fatiado com 7 tickets. Deseja aprovar o plano e liberar a execução autônoma do BUILD?*  
*Aprovado este documento e gerado o ROUTE correspondente, o BUILD pode ser iniciado.*
