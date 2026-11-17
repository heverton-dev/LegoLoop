---
id: PLAN-MEGA-PLANO-PRODUCAO
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [SPEC-MEGA-PLANO-PRODUCAO]
covers: [AC-E01.1, AC-E01.2, AC-E01.3, AC-E01.4, AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.4, AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-E06.1, AC-E06.2, AC-E06.3, AC-E07.1, AC-E07.2, AC-E07.3, AC-E07.4, AC-E08.1, AC-E08.2, AC-E08.3, AC-E08.4, AC-E09.1, AC-E09.2, AC-E10.1, AC-E10.2, AC-E10.3, AC-NF.1, AC-NF.2]
---

# 🗺️ PLAN: Mega Plano LegoLoop para Produção

> **Status:** Em Planejamento
> **Data:** 2026-11-16
> **Versão:** 1.0.0
> **Autor:** Antigravity AI + Usuário
> **Iniciativa:** Mega Plano — LegoLoop Pronto para Produção (I0–I9)
> **Referência da Intenção:** [docs/intents/2026-11-16_INTENT_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md](../intents/2026-11-16_INTENT_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md)
> **Referência da Especificação:** [docs/specs/2026-11-16_SPEC_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md](../specs/2026-11-16_SPEC_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md)
> **Espelho Executável JSON:** [docs/routes/2026-11-16_ROUTE_MEGA_PLANO_PRODUCAO.json](../routes/2026-11-16_ROUTE_MEGA_PLANO_PRODUCAO.json)
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)

---

## 1. Estratégia de Execução & Governança de Loteamento

> [!CAUTION]
> **AUTORIZAÇÃO EXPLÍCITA DO USUÁRIO (2026-11-16 21:26, sessão `8e1ac50d`):** stop-gates humanos (`SG-*`) desta demanda são substituídos por gates de AUDITOR autônomo — evidência mecânica (roda os `verify` do nó) + julgamento, máx. 3 re-auditorias por nó, após isso `NAO_CONCLUIDO` com relatório e continuidade. Circuit breaker por épico: teto excedido → `NAO_CONCLUIDO` e segue. Checkpoint commits por nó e tag `pre-mega-plano` autorizados. **Push proibido.** Prazo: conclusão até 2026-11-17 06:00 (America/Sao_Paulo). Exceção humana única: SG-FINAL é reportado ao usuário ao amanhecer.

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `ephemeral` (AGENTS.md, regra 4 — as seis condições; SWARM nativo do harness como motor de execução paralela)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor; mutação prova cada teste novo)
* **Limite de Workers Paralelos:** `4` (SWARM: 1 worker por nano-ticket, gates por ticket)
* **Timeout Padrão por Ticket:** `300s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`
* **Ordem das Iniciativas:** I0→I1→I2/I3→I4→I5→I6→I7→I8→I9 (SPEC premissa 5; cada SG libera o épico seguinte)

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | I0 · Saneamento do `docs/` | 4 | `AC-E01.1`–`AC-E01.4` | `micro` |
| **E02** | I1 · Nascimento 5/5 + Templates | 4 | `AC-E02.1`–`AC-E02.4` | `micro` |
| **E03** | I2 · Memória Concorrente | 3 | `AC-E03.1`–`AC-E03.4` | `micro` |
| **E04** | I3 · Anticorpos Persistentes | 4 | `AC-E04.1`–`AC-E04.4` | `micro` |
| **E05** | I4 · VSA por Processo | 4 | `AC-E05.1`–`AC-E05.4` | `micro` |
| **E06** | I5 · Rename `loops/`→`processos/` | 3 | `AC-E06.1`–`AC-E06.3` | `micro` |
| **E07** | I6 · Tríplice Via de Invocação | 4 | `AC-E07.1`–`AC-E07.4` | `micro` |
| **E08** | I7 · Auditoria Universal + 6 Verbos | 4 | `AC-E08.1`–`AC-E08.4` | `micro` |
| **E09** | I8 · MCP Único via Router | 2 | `AC-E09.1`, `AC-E09.2` | `micro` |
| **E10** | I9 · Provas Ausentes & Higiene | 5 | `AC-E10.1`–`AC-E10.3`, `AC-NF.1`, `AC-NF.2` | `micro` |

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Saneamento do `docs/` | `AC-E01.1`–`AC-E01.4` | `TICK-01.1`–`TICK-01.4` |
| `R2` | Nascimento 5/5 + templates | `AC-E02.1`–`AC-E02.4` | `TICK-02.1`–`TICK-02.4` |
| `R3` | Memória concorrente | `AC-E03.1`–`AC-E03.4` | `TICK-03.1`–`TICK-03.3` |
| `R4` | Anticorpos persistentes | `AC-E04.1`–`AC-E04.4` | `TICK-04.1`–`TICK-04.4` |
| `R5` | VSA por processo | `AC-E05.1`–`AC-E05.4` | `TICK-05.1`–`TICK-05.4` |
| `R6` | Rename `loops/`→`processos/` | `AC-E06.1`–`AC-E06.3` | `TICK-06.1`–`TICK-06.3` |
| `R7` | Tríplice via | `AC-E07.1`–`AC-E07.4` | `TICK-07.1`–`TICK-07.4` |
| `R8` | Auditoria universal + 6 verbos | `AC-E08.1`–`AC-E08.4` | `TICK-08.1`–`TICK-08.4` |
| `R9` | MCP único via router | `AC-E09.1`, `AC-E09.2` | `TICK-09.1`, `TICK-09.2` |
| `R10` | Provas ausentes & higiene | `AC-E10.1`–`AC-E10.3`, `AC-NF.1`, `AC-NF.2` | `TICK-10.1`–`TICK-10.5` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

> Todos os tickets: Iniciativa `Mega-Plano-Producao` · Workspace `local` (auditorias em `tempdir`) · Executor `orchestrix_hermetic` · Subagentes `ephemeral` · Write Boundary `strict` · Rollback `stash_changes` · Fallback Harness `antigravity` · Modelo `flash`. Roteamento padrão: Harness `9router` · Modelo `code-pro` · Effort `high` (nanos: `code-fast` · `low`).

### Épico 01 — I0 · Saneamento do `docs/`

#### Ticket `TICK-01.1` — Inventário classificado de todo `docs/` (SWARM por diretório)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]` · **Tooling:** Runtime `powershell` · MCPs `false` · Allowlist `[]`
* **Limites:** Timeout `600s` · Retries `1` · Entrada `32000` · Saída `1200`
* **Inputs:** `docs/**` · Spec Anchor `#AC-E01.1`
* **Alvo de Teste TDD:** `internal/inventario/inventario_test.go:TestClassificaTodosArquivos`
* **Outputs:** `docs/reports/2026-11-16_INVENTARIO_DOCS.md`, `internal/inventario/inventario.go`
* **Verificação:** `go test -v ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-16_INVENTARIO_DOCS.md')`
* **Spec Ref:** `AC-E01.1` · **Escopo:** 3598 arquivos classificados em {IMPLEMENTADO, ULTRAPASSADO, SUPERADO, PENDENTE, REVER} com dono — zero sem classificação.

#### Ticket `TICK-01.2` — Migração para `docs/archive/` (videorea, sessions, reports históricos)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `16000` · Saída `600`
* **Inputs:** `docs/videorea/**`, `docs/sessions/**`, `docs/reports/**` · Spec Anchor `#AC-E01.2`
* **Alvo de Teste TDD:** `internal/inventario/migrar_test.go:TestMigraParaArchive`
* **Outputs:** `docs/archive/**`, `internal/inventario/migrar.go`
* **Verificação:** `go test -v ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0 && dir_exists('docs/archive')`
* **Spec Ref:** `AC-E01.2`

#### Ticket `TICK-01.3` — Guard de proteção de lei (templates/vocab indeletáveis)
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.2"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Inputs:** `internal/inventario/migrar.go` · Spec Anchor `#AC-E01.3`
* **Alvo de Teste TDD:** `internal/inventario/guard_test.go:TestRecusaDelecaoTemplate`
* **Outputs:** `internal/inventario/guard.go`
* **Verificação:** `go test -v ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.3`

#### Ticket `TICK-01.4` — Pendentes sem dono viram tickets do PLAN
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.3"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Inputs:** `docs/reports/2026-11-16_INVENTARIO_DOCS.md` · Spec Anchor `#AC-E01.4`
* **Outputs:** inventário com seção PENDENTE→tickets
* **Verificação:** `go vet ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Inventário e arquivo morto do docs/ aprovados? Autoriza iniciar E02 (Nascimento 5/5)?')`
- [ ] 100% dos tickets de E01 com PASS nos gates mecânicos.
- [ ] Inventário completo apresentado (3598 arquivos classificados).
- [ ] Autorização explícita do usuário para iniciar o Épico 02.

---

### Épico 02 — I1 · Nascimento 5/5 + Templates

#### Ticket `TICK-02.1` — Templates canônicos: TRIAGEM, MEMORIA_EXECUCAO, ANTICORPOS
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `16000` · Saída `600`
* **Inputs:** `docs/templates/templates_test.go` · Spec Anchor `#AC-E02.2`
* **Alvo de Teste TDD:** `docs/templates/templates_test.go:TestTemplatesNascimento`
* **Outputs:** `docs/templates/TEMPLATE_TRIAGEM.json`, `docs/templates/TEMPLATE_MEMORIA_EXECUCAO.json`, `docs/templates/TEMPLATE_ANTICORPOS.json`
* **Verificação:** `go test -v ./docs/templates/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/templates/TEMPLATE_TRIAGEM.json') && file_exists('docs/templates/TEMPLATE_MEMORIA_EXECUCAO.json') && file_exists('docs/templates/TEMPLATE_ANTICORPOS.json')`
* **Spec Ref:** `AC-E02.2`

#### Ticket `TICK-02.2` — `fabricar` instancia scaffold 5/5 validado por schema
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Inputs:** `internal/fabricar/*`, `docs/templates/*` · Spec Anchor `#AC-E02.1`
* **Alvo de Teste TDD:** `internal/fabricar/scaffold_test.go:TestNasceCincoPorCinco`
* **Outputs:** `internal/fabricar/scaffold.go` + testes
* **Verificação:** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.3` — Adequação dos processos existentes a 5/5
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `8000` · Saída `300`
* **Inputs:** `loops/**` (5 processos) · Spec Anchor `#AC-E02.3`
* **Outputs:** processos com 5/5 ou exceção documentada
* **Verificação:** `legoloop validate` em cada processo
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.3`

#### Ticket `TICK-02.4` — Slug inválido: exit 2, zero diretório criado
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.2"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/cli/fabricar_test.go:TestNomeInvalidoExitDois`
* **Outputs:** `internal/cli/fabricar.go` (validação de slug)
* **Verificação:** `go test -v ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Nascimento 5/5 e templates aprovados? Autoriza iniciar E03/E04 (memória e anticorpos)?')`
- [ ] 100% dos tickets de E02 com PASS nos gates mecânicos.
- [ ] Processo exemplar criado com 5/5 apresentado.
- [ ] Autorização explícita do usuário.

---

### Épico 03 — I2 · Memória de Execução Concorrente

#### Ticket `TICK-03.1` — Biblioteca de lock por arquivo (Session+PID) em `internal/triagem`
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Inputs:** `internal/triagem/memoria.go` · Spec Anchor `#AC-E03.1` (cobre `#AC-E03.2`)
* **Alvo de Teste TDD:** `internal/triagem/memoria_test.go:TestLockConcorrente` + `:TestLockOrfao`
* **Outputs:** `internal/triagem/lock.go` + testes
* **Verificação:** `go test -race -v ./internal/triagem/...` (ou `-count=1` sem CGO)
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.1`, `AC-E03.2` · **Escopo:** 2 sessões gravam; segunda com PID vivo recebe `ErrLockAtivo`; PID morto → lock tomado; integridade byte a byte.

#### Ticket `TICK-03.2` — Zerar memória no `fechar`
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/fechar/quarteto_test.go:TestZeraMemoria`
* **Outputs:** `internal/fechar/quarteto.go`
* **Verificação:** `go test -v ./internal/fechar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.3`

#### Ticket `TICK-03.3` — Memória corrompida: sentinela + backup `.legoloop-bak`
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/triagem/memoria_test.go:TestCorrompidoBackup`
* **Outputs:** `internal/triagem/memoria.go`
* **Verificação:** `go test -v ./internal/triagem/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Memória concorrente segura aprovada? Autoriza iniciar E05 (VSA)?')`
- [ ] 100% dos tickets de E03 com PASS.
- [ ] Prova de concorrência (2 sessões) apresentada.
- [ ] Autorização explícita do usuário.

---

### Épico 04 — I3 · Motor de Anticorpos Persistente

#### Ticket `TICK-04.1` — Persistência em `anticorpos.json` + schema `anticorpos.v1`
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Inputs:** `internal/brain/antibodies.go:22` · Spec Anchor `#AC-E04.1`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestPersisteVacina`
* **Outputs:** `schemas/anticorpos.v1.json`, `internal/brain/antibodies.go` (persistência)
* **Verificação:** `go test -v ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('schemas/anticorpos.v1.json')`
* **Spec Ref:** `AC-E04.1`

#### Ticket `TICK-04.2` — Vacinação de reincidência (SHA-256 + injeção preventiva)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestVacinaReincidencia`
* **Outputs:** `internal/brain/antibodies.go` (Find+InjectPreventiveContext)
* **Verificação:** `go test -v ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E04.2` · **Escopo:** mesmo erro (assinatura) reocorre → `recorrencias`++ e antídoto injetado ANTES da execução.

#### Ticket `TICK-04.3` — Escopo global (infra/LLM) vs domínio (local)
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.2"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestEscopoGlobalDominio`
* **Outputs:** `internal/brain/antibodies.go`
* **Verificação:** `go test -v ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E04.3`

#### Ticket `TICK-04.4` — `anticorpos.json` corrompido: backup, não trava
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-04.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/brain/antibodies_test.go:TestCorrompidoNaoTrava`
* **Outputs:** `internal/brain/antibodies.go`
* **Verificação:** `go test -v ./internal/brain/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E04.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Anticorpos persistentes e vacinação aprovados? Autoriza iniciar E05 (VSA)?')`
- [ ] 100% dos tickets de E04 com PASS.
- [ ] Prova de vacinação de reincidência apresentada.
- [ ] Autorização explícita do usuário.

---

### Épico 05 — I4 · VSA por Processo (Fatia Vertical Autocontida)

#### Ticket `TICK-05.1` — Scaffold do slice (manifestos + `docs/` do processo)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03", "SG-E04"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `9000`
* **Inputs:** `internal/fabricar/scaffold.go` · Spec Anchor `#AC-E05.1`
* **Alvo de Teste TDD:** `internal/fabricar/scaffold_test.go:TestSliceAutocontido`
* **Outputs:** `internal/fabricar/scaffold.go` (fatia vertical)
* **Verificação:** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.1`

#### Ticket `TICK-05.2` — Migração do Quinteto (286 docs) para os slices
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1", "SG-E01"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `32000` · Saída `1200`
* **Inputs:** `docs/intents|specs|plans|routes|dags/**` · Spec Anchor `#AC-E05.2`
* **Alvo de Teste TDD:** `internal/inventario/migrar_test.go:TestMigraQuinteto`
* **Outputs:** slices com Quinteto + `docs/archive/` para órfãos
* **Verificação:** `go test -v ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.2` · **Escopo:** maior blast radius — cada documento ganha dono (slice) ou vai ao archive com registro.

#### Ticket `TICK-05.3` — Abstração por SO (peças `.ps1`/`.sh` resolvidas por SO)
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/fabricar/scaffold_test.go:TestResolvePorSO`
* **Outputs:** `internal/fabricar/scaffold.go`
* **Verificação:** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.3`

#### Ticket `TICK-05.4` — Órfãos do Quinteto → `docs/archive/` com registro
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.2"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `8000` · Saída `300`
* **Spec Anchor:** `#AC-E05.4`
* **Outputs:** inventário com seção ÓRFÃOS
* **Verificação:** `go vet ./internal/inventario/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 05 (`SG-E05`)
* **Gate (DSL):** `prompt('Arquitetura VSA consolidada? Autoriza iniciar E06 (rename loops/→processos/)?')`
- [ ] 100% dos tickets de E05 com PASS.
- [ ] Slice exemplar apresentado; docs/ raiz enxuto comprovado.
- [ ] Autorização explícita do usuário (rename mexe em Lei — mudança deliberada).

---

### Épico 06 — I5 · Renomeação `loops/` → `processos/`

#### Ticket `TICK-06.1` — Rename + migração de todas as referências vivas
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E05"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `32000` · Saída `1200`
* **Inputs:** `loops/**`, código Go, schemas, skills, docs · Spec Anchor `#AC-E06.1`
* **Alvo de Teste TDD:** `internal/caminhos/caminhos_test.go:TestReferenciasProcessos`
* **Outputs:** `processos/**` + referências migradas
* **Verificação:** `go test -v ./...`
* **Gate (DSL):** `exit_code == 0 && dir_exists('processos')`
* **Spec Ref:** `AC-E06.1`

#### Ticket `TICK-06.2` — Binário e schemas apontam `processos/`
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/fechar/quarteto_test.go:TestGravaEmProcessos`
* **Outputs:** `internal/fechar/quarteto.go`, `internal/caminhos/*`
* **Verificação:** `legoloop validate` + `go test -v ./internal/fechar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E06.2`

#### Ticket `TICK-06.3` — Comando legado `loops/`: exit 2 com guia de migração
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-06.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/cli/legado_test.go:TestLoopsLegadoExitDois`
* **Outputs:** `internal/cli/legado.go`
* **Verificação:** `go test -v ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E06.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 06 (`SG-E06`)
* **Gate (DSL):** `prompt('Rename concluído sem quebra? Autoriza iniciar E07 (tríplice via)?')`
- [ ] 100% dos tickets de E06 com PASS.
- [ ] `git grep "loops/"` sem referências vivas.
- [ ] Autorização explícita do usuário.

---

### Épico 07 — I6 · Tríplice Via de Invocação

#### Ticket `TICK-07.1` — Via CLI: `legoloop exec <slug>`
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E06"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Alvo de Teste TDD:** `internal/cli/exec_test.go:TestExecProcesso`
* **Outputs:** `internal/cli/exec.go`
* **Verificação:** `go test -v ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E07.1`

#### Ticket `TICK-07.2` — Via harness: SKILL + slash + `instalar-skills`
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Alvo de Teste TDD:** `internal/skills/instalar_test.go:TestDistribuiTodosHarness`
* **Outputs:** `internal/skills/instalar.go`
* **Verificação:** `go test -v ./internal/skills/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E07.2`

#### Ticket `TICK-07.3` — Via Studio: peça registrada automaticamente
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.1"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `24000` · Saída `900`
* **Alvo de Teste TDD:** `internal/studio/pecas_test.go:TestRegistraPecaAutomatica`
* **Outputs:** `internal/studio/pecas.go`
* **Verificação:** `go test -v ./internal/studio/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E07.3`

#### Ticket `TICK-07.4` — SKILL inválido: exit ≠ 0 aponta arquivo, harness íntegro
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-07.2"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/skills/instalar_test.go:TestSkillInvalidoExit`
* **Outputs:** `internal/skills/instalar.go`
* **Verificação:** `go test -v ./internal/skills/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E07.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 07 (`SG-E07`)
* **Gate (DSL):** `prompt('Tríplice via funcionando? Autoriza iniciar E08 (auditoria universal)?')`
- [ ] 100% dos tickets de E07 com PASS.
- [ ] Processo invocado nas 3 vias apresentado.
- [ ] Autorização explícita do usuário.

---

### Épico 08 — I7 · Auditoria Universal Agnóstica + 6 Verbos

#### Ticket `TICK-08.1` — Processo `auditoria-atomica` (agnóstico de alvo)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E07"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `32000` · Saída `1200`
* **Alvo de Teste TDD:** `internal/fabricar/auditoria_test.go:TestAlvoArbitrario`
* **Outputs:** `processos/auditoria-atomica/` (5/5) + `skills/auditoria-atomica/SKILL.md`
* **Verificação:** `legoloop validate` + `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('processos/auditoria-atomica/pipeline.dag.json')`
* **Spec Ref:** `AC-E08.1` · **Escopo:** alvo = verbo|processo|peça|função|documento|binário|repo externo. Saída: disseção+lacunas+tickets+refatoração+veredito.

#### Ticket `TICK-08.2` — Auditoria profunda dos 6 verbos (matriz, listar, refinar, remover, compor, fechar)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1"]` · **Limites:** Timeout `1800s` · Retries `1` · Entrada `32000` · Saída `1200`
* **Inputs:** `internal/{matriz,listar,refinar,remover,compor,fechar}/**` · Spec Anchor `#AC-E08.2`
* **Outputs:** `docs/reports/2026-11-17_REPORT_AUDITORIA_6_VERBOS.md` + tickets de correção executados
* **Verificação:** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0 && file_exists('docs/reports/2026-11-17_REPORT_AUDITORIA_6_VERBOS.md')`
* **Spec Ref:** `AC-E08.2` · **Escopo:** mesma profundidade da auditoria do `fabricar` — `compor` e `fechar` intocados hoje.

#### Ticket `TICK-08.3` — Mutações provadas (Green-Red-Green) em todo teste novo
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.2"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `8000` · Saída `300`
* **Outputs:** relatório com matriz de mutações restauradas
* **Verificação:** `go test -count=1 ./...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E08.3`

#### Ticket `TICK-08.4` — Alvo inexistente: exit 2, zero artefato órfão
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-08.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/fabricar/auditoria_test.go:TestAlvoInexistenteExitDois`
* **Outputs:** `internal/fabricar/auditoria.go`
* **Verificação:** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E08.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 08 (`SG-E08`)
* **Gate (DSL):** `prompt('Auditoria universal e 6 verbos aprovados? Autoriza iniciar E09 (MCP)?')`
- [ ] 100% dos tickets de E08 com PASS.
- [ ] Relatórios dos 6 verbos apresentados.
- [ ] Autorização explícita do usuário.

---

### Épico 09 — I8 · MCP Único via `mcp-router`

#### Ticket `TICK-09.1` — Router como interface única de tools dos processos
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E08"]` · **Limites:** Timeout `300s` · Retries `1` · Entrada `16000` · Saída `600`
* **Alvo de Teste TDD:** `internal/mcp/router_test.go:TestToolsViaRouter`
* **Outputs:** `internal/mcp/router.go`
* **Verificação:** `go test -v ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E09.1`

#### Ticket `TICK-09.2` — Conflito de nome de tool: registro recusado
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-09.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/mcp/router_test.go:TestConflitoToolExit`
* **Outputs:** `internal/mcp/router.go`
* **Verificação:** `go test -v ./internal/mcp/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E09.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 09 (`SG-E09`)
* **Gate (DSL):** `prompt('MCP único via router aprovado? Autoriza iniciar E10 (provas e higiene)?')`
- [ ] 100% dos tickets de E09 com PASS.
- [ ] Autorização explícita do usuário.

---

### Épico 10 — I9 · Provas Ausentes & Higiene

#### Ticket `TICK-10.1` — Débitos zerados: taxlinter 14 + funclen 3
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E09"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `24000` · Saída `900`
* **Outputs:** normalização `internal/refinar/curado_v2*`, `loops/legoloop-swarm`, `swarm_run.go`, `cites.go`
* **Verificação:** `go run ./tools/taxlinter internal cmd loops skills docs/templates` + `go run ./tools/funclen internal`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E10.1`

#### Ticket `TICK-10.2` — 2 testes flaky estáveis (3 execuções consecutivas)
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-10.1"]` · **Limites:** Timeout `600s` · Retries `3` · Entrada `16000` · Saída `600`
* **Alvo de Teste TDD:** `internal/board/*_test.go:TestUnifiedPathFileLockExactMentionAndPIDKill` + `internal/mcp/*_test.go:TestSwarmSpawnToolRealProcess`
* **Outputs:** correção dos 2 testes
* **Verificação:** `go test -count=1 ./internal/board/... ./internal/mcp/...` (3×)
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E10.2`

#### Ticket `TICK-10.3` — Provas ausentes: `-race`, LLM real, grafo semântico
* **Nature:** `micro` | **Tipo:** `unit_no` | **Modo:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-10.2"]` · **Limites:** Timeout `900s` · Retries `1` · Entrada `16000` · Saída `600`
* **Outputs:** `docs/reports/2026-11-17_REPORT_PROVAS_AUSENTES.md` (prova ou NAO_PROVADO com dono+prazo)
* **Verificação:** `go test -race -count=1 ./...` (com CGO) / registro documental
* **Gate (DSL):** `file_exists('docs/reports/2026-11-17_REPORT_PROVAS_AUSENTES.md')`
* **Spec Ref:** `AC-E10.3`

#### Ticket `TICK-10.4` — Path Jail hostil: `..`, JSON malicioso, PID falso
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-10.1"]` · **Limites:** Timeout `180s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/paths/jail_test.go:TestEntradaHostil`
* **Outputs:** `internal/paths/jail.go`
* **Verificação:** `go test -v ./internal/paths/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-NF.1`

#### Ticket `TICK-10.5` — Carga: 50 workers, sem lock global, memória íntegra
* **Nature:** `nano` | **Tipo:** `unit_no` | **Modo:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-10.4"]` · **Limites:** Timeout `600s` · Retries `1` · Entrada `8000` · Saída `300`
* **Alvo de Teste TDD:** `internal/triagem/carga_test.go:TestCinquentaWorkers`
* **Outputs:** `internal/triagem/carga_test.go`
* **Verificação:** `go test -v -count=1 ./internal/triagem/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-NF.2`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test -count=1 ./...` (3 execuções)
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia); `go run ./tools/taxlinter internal cmd loops skills docs/templates`; `go run ./tools/funclen internal`
3. Recompilação obrigatória: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`
4. Relatório formal: `docs/reports/2026-11-17_REPORT_MEGA_PLANO_PRODUCAO.md` (TEMPLATE_REPORT v2.2.0 — telemetria tempo/tokens/custo USD/BRL)
5. Registro de sessão: `faz-session --id 8e1ac50d-4232-4ce2-883a-c4bd3bbc222c --format all`
6. Governança viva: `AGENTS.md`, `MEMORY.md`, `README.md` + reindexação do grafo (`codebase-memory-mcp cli index_repository`)

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*
*Aprovado este documento e gerado o `docs/routes/2026-11-16_ROUTE_MEGA_PLANO_PRODUCAO.json`, o BUILD pode ser iniciado.*
