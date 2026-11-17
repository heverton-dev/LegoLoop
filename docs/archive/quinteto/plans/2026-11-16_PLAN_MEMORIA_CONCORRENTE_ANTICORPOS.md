---
id: PLAN-MEMORIA_CONCORRENTE_ANTICORPOS
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [SPEC-MEMORIA_CONCORRENTE_ANTICORPOS]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E01.4, AC-E01.5, AC-E01.6, AC-E01.7, AC-E01.8, AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.5, AC-E02.6, AC-E02.7, AC-E02.8, AC-E03.1, AC-E03.2, AC-E03.3]
---

# 🗺️ PLAN: Memória de Execução Concorrente + Anticorpos Persistidos

> **Status:** Em Planejamento (aguarda SG-M01)  
> **Data:** 2026-11-16  
> **Versão:** 1.0.0  
> **Autor:** planejador (enxame SWARM)  
> **Iniciativa:** MEMORIA_CONCORRENTE_ANTICORPOS  
> **Referência da Intenção:** [docs/intents/2026-11-16_INTENT_MEMORIA_EXECUCAO_ANTICORPOS.md](../intents/2026-11-16_INTENT_MEMORIA_EXECUCAO_ANTICORPOS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md](../specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-16_ROUTE_MEMORIA_CONCORRENTE_ANTICORPOS.json](../routes/2026-11-16_ROUTE_MEMORIA_CONCORRENTE_ANTICORPOS.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `ephemeral` (AGENTS.md, regra 4: 1 worker estéril por Ticket, seis condições cumpridas)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor; mutação prova o teste)
* **Limite de Workers Paralelos:** `2` (cadeias M1 ∥ M2 independentes)
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_processo`; por Nó `stash_changes`

**Gate de qualidade por Ticket (obrigatório, todos os itens):**
1. `go test -count=1 ./...`
2. `go vet ./...`
3. `go run ./tools/funclen internal` (nenhuma função > 35 linhas)
4. `gofmt -l .` com saída vazia
5. `git diff --exit-code -- go.mod go.sum`
6. `go run ./tools/taxlinter internal cmd loops skills docs/templates` (sem violação nova)

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Memória de Execução Concorrente | 4 | `AC-E01.1` … `AC-E01.8` | `micro` |
| **E02** | Anticorpos Persistidos | 4 | `AC-E02.1` … `AC-E02.8` | `micro` |
| **E03** | Fechamento Canônico Lei 11 | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `nano` |

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Memória de execução concorrente com lock e sessão | `AC-E01.1` … `AC-E01.8` | `M1.1`, `M1.2`, `M1.3`, `M1.4` |
| `R2` | Anticorpos persistidos com vacinação de reincidência | `AC-E02.1` … `AC-E02.8` | `M2.1`, `M2.2`, `M2.3`, `M2.4` |
| `R3` | Fechamento canônico Lei 11 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `M3.1`, `M3.2`, `M3.3` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### 🛑 STOP-GATE HUMANO `SG-M01` — Pós-Planejamento
* **Gate (DSL):** `prompt('Quinteto aprovado? Autoriza executar M1.1 e M2.1?')`
- [ ] SPEC, PLAN, ROUTE e DAG apresentados ao usuário.
- [ ] Autorização explícita concedida para iniciar os Tickets.

### Épico 01 — Memória de Execução Concorrente

#### Ticket `M1.1` — Modelo `MemoriaExecucao` ganha chave de sessão (ID de sessão + PID)
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-M01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` (backoff `1000ms`, fallback model `true`) · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/triagem/memoria.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/triagem/memoria_test.go:TestMemoriaSessao`
* **Outputs:** Artifacts `internal/triagem/memoria.go`, `internal/triagem/memoria_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/triagem/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/triagem/memoria.go') && file_exists('internal/triagem/memoria_test.go')`
* **Spec Ref:** `AC-E01.1`, `AC-E01.2`
* **Escopo:** Teste prova serialização de `sessao_id` e `pid` e rejeição `ErrSessaoInvalida` para ID vazio ou PID ≤ 0. Mutação remove o campo e o teste falha.

#### Ticket `M1.2` — Lock de arquivo com PID (padrão `internal/board`) em `SalvarMemoria`/`LerMemoria`
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["M1.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/triagem/memoria.go`, `internal/board/swarm_lock.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E01.3`
* **Alvo de Teste TDD:** `internal/triagem/memoria_test.go:TestMemoriaLockConcorrente`
* **Outputs:** Artifacts `internal/triagem/memoria.go`, `internal/triagem/memoria_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 -race ./internal/triagem/...` (sem CGO: `go test -count=1 ./internal/triagem/...`)
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.3`, `AC-E01.4`
* **Escopo:** Teste de corrida com 2 sessões simultâneas no mesmo diretório (`t.TempDir()`): JSON íntegro nas duas leituras. Lock retido devolve `ErrLockOcupado` sem escrita parcial. Mutação remove o lock e o teste de corrida falha.

#### Ticket `M1.3` — `runner` grava por sessão; `fechar` zera só a sessão dona e limpa arquivo órfão
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["M1.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/fechar/quarteto.go`, `internal/triagem/memoria.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E01.5`
* **Alvo de Teste TDD:** `internal/fechar/quarteto_test.go:TestFecharSessaoDona`
* **Outputs:** Artifacts `internal/fechar/quarteto.go`, `internal/fechar/quarteto_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/fechar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.5`, `AC-E01.6`
* **Escopo:** Ciclo completo em `t.TempDir()`: 2 sessões ativas + 1 órfã (PID morto). `fechar` zera a dona, preserva a ativa, limpa a órfã. Sessão inexistente devolve erro citando o ID sem alterar disco.

#### Ticket `M1.4` — CLI/Studio expõe sessões ativas do processo
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["M1.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `12000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `cmd/legoloop/main.go`, `internal/triagem/memoria.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E01.7`
* **Alvo de Teste TDD:** `cmd/legoloop/main_test.go:TestSessoesAtivas`
* **Outputs:** Artifacts `cmd/legoloop/main.go`, `cmd/legoloop/main_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./cmd/legoloop/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.7`, `AC-E01.8`
* **Escopo:** Contrato de CLI testado: 2 sessões ⇒ JSON com `sessao_id`, `pid`, `etapa_atual`; 0 sessões ⇒ `[]` e exit 0. Saída sem caminho absoluto (Lei 8).

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-M02`)
* **Gate (DSL):** `prompt('Evidências de M1 aprovadas? Autoriza iniciar M2?')`
- [ ] 100% dos Tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes (incluindo corrida) apresentadas com arquivo:linha.
- [ ] Autorização explícita do usuário concedida.

### Épico 02 — Anticorpos Persistidos

#### Ticket `M2.1` — `AntibodiesDB` persiste em `processos/<slug>/anticorpos.json`
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-M01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/brain/antibodies.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestAnticorposPersistencia`
* **Outputs:** Artifacts `internal/brain/antibodies.go`, `internal/brain/antibodies_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`
* **Escopo:** Carrega ao iniciar e grava ao registrar, em `t.TempDir()`. JSON corrompido devolve `ErrAnticorposCorrompidos` sem panic e sem sobrescrever. Mutação remove a gravação e o teste de persistência falha.

#### Ticket `M2.2` — Vacinação de reincidência: `Record` deduplica por `errSig`; reincidência aplica antídoto
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["M2.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/brain/antibodies.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E02.3`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestVacinaReincidencia`
* **Outputs:** Artifacts `internal/brain/antibodies.go`, `internal/brain/antibodies_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.3`, `AC-E02.4`
* **Escopo:** Mesmo erro 2× ⇒ `Recorrencias == 2`, antídoto aplicado e erro não repete no teste. `errSig` vazio ⇒ `ErrAnticorpoInvalido` e banco inalterado. Mutação remove a deduplicação e o teste falha.

#### Ticket `M2.3` — Injeção preventiva no início da execução (`InjectPreventiveContext` vira passo do processo)
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["M2.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/brain/antibodies.go`, `internal/triagem` (execução) · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E02.5`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestInjecaoPreventiva`
* **Outputs:** Artifacts `internal/brain/antibodies.go`, `internal/brain/antibodies_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/brain/... ./internal/triagem/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.5`, `AC-E02.6`
* **Escopo:** Teste observa o contexto do Nó 1 e confirma a advertência preventiva presente antes dele. Sem anticorpos ⇒ contexto vazio e execução segue. Mutação remove o passo de injeção e o teste falha.

#### Ticket `M2.4` — Gate anti-reincidência em runtime (erro conhecido ⇒ antídoto ou bloqueio com saída clara)
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["M2.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `internal/brain/antibodies.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E02.7`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestGateAntiReincidencia`
* **Outputs:** Artifacts `internal/brain/antibodies.go`, `internal/brain/antibodies_test.go`
* **Verificação (argv, sem shell):** `go test -count=1 ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.7`, `AC-E02.8`
* **Escopo:** Erro conhecido ⇒ antídoto aplicado; sem antídoto aplicável ⇒ bloqueio citando o `errSig`. Erro desconhecido ⇒ fluxo normal. Mutação desativa o gate e o teste falha.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-M03`)
* **Gate (DSL):** `prompt('Evidências de M2 aprovadas? Autoriza iniciar o Fechamento M3?')`
- [ ] 100% dos Tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de vacinação e gate anti-reincidência apresentadas com arquivo:linha.
- [ ] Autorização explícita do usuário concedida.

### Épico 03 — Fechamento Canônico Lei 11

#### Ticket `M3.1` — Manual operacional em `docs/manuals/`
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-M02", "SG-M03"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `12000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Context Artifacts `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E03.1`
* **Outputs:** Artifacts `docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md`
* **Verificação (argv, sem shell):** `powershell -NoProfile -Command "Test-Path docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md"`
* **Gate (DSL):** `file_exists('docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md') && !has_pattern('docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md', '\\{\\{')`
* **Spec Ref:** `AC-E03.1`
* **Escopo:** Manual com checklist Lei 11; sem placeholders.

#### Ticket `M3.2` — Livro LegoPress + skill canônica + DAG validado (`legoloop validate`)
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `subprocesso` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["M3.1"]`
* **Ref do Processo filho:** `legoloop-legopress` (id de Processo do acervo; nunca um caminho)
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Context Artifacts `docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E03.2`
* **Outputs:** Artifacts `docs/books/memoria-anticorpos/livro.json`, `docs/books/memoria-anticorpos/livro.pdf`, `skills/legoloop-memoria-anticorpos/SKILL.md`
* **Verificação (argv, sem shell):** `legoloop validate`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/books/memoria-anticorpos/livro.json') && file_exists('skills/legoloop-memoria-anticorpos/SKILL.md')`
* **Spec Ref:** `AC-E03.2`

#### Ticket `M3.3` — Build + install, governança viva, reindex do grafo, relatório com telemetria
* **Iniciativa:** `MEMORIA_CONCORRENTE_ANTICORPOS` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["M3.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Source Files `cmd/legoloop/main.go`, `AGENTS.md`, `MEMORY.md`, `README.md` · Spec Anchor `docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md#AC-E03.3`
* **Outputs:** Artifacts `bin/legoloop.exe`, `docs/reports/2026-11-16_REPORT_MEMORIA_CONCORRENTE_ANTICORPOS.md`
* **Verificação (argv, sem shell):** `go install ./cmd/legoloop`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-16_REPORT_MEMORIA_CONCORRENTE_ANTICORPOS.md')`
* **Spec Ref:** `AC-E03.3`
* **Escopo:** Relatório conforme TEMPLATE_REPORT v2.2.0 com telemetria de tempo, tokens e custo USD/BRL orçado vs realizado.

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL-M`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test -count=1 ./...`
2. Estática e formatação: `go vet ./...`, `gofmt -l .` (saída vazia), `go run ./tools/funclen internal`, `go run ./tools/taxlinter internal cmd loops skills docs/templates`
3. Recompilação e instalação frescas: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`
4. Reindexação do grafo semântico: `codebase-memory-mcp cli index_repository`
5. Emissão do relatório formal: `docs/reports/2026-11-16_REPORT_MEMORIA_CONCORRENTE_ANTICORPOS.md`
6. Registro de histórico da sessão via `faz-session --id <ID> --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA (`SG-M01`):**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e a ROUTE, o BUILD dos Tickets M1.1→M3.3 pode ser iniciado.*
