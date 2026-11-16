---
id: PLAN-FABRICAR-AGNOSTICO
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [SPEC-FABRICAR-AGNOSTICO]
covers: [AC-E01.1, AC-E01.2, AC-E02.1, AC-E02.2, AC-E03.1, AC-E03.2, AC-E03.3, AC-E04.1, AC-E04.2, AC-E05.1, AC-E05.2, AC-NF.1]
---

# 🗺️ PLAN: Fabricar Agnóstico — Contrato de I/O, Molde por Natureza e Verbos Derivados

> **Status:** Em Planejamento
> **Data:** 2026-11-16
> **Versão:** 1.0.0
> **Autor:** audit-verbos
> **Iniciativa:** I01 — Agnosticidade Universal dos Verbos de Processo
> **Referência da Intenção:** [docs/intents/2026-11-16_INTENT_FABRICAR_AGNOSTICO.md](../intents/2026-11-16_INTENT_FABRICAR_AGNOSTICO.md)
> **Referência da Especificação:** [docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md](../specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md)
> **Espelho Executável JSON:** [docs/routes/2026-11-16_ROUTE_FABRICAR_AGNOSTICO.json](../routes/2026-11-16_ROUTE_FABRICAR_AGNOSTICO.json)
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por Ticket, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Ticket:** `180s`
* **Política de Rollback:** global `abort_processo`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Contrato de I/O da Peça | 2 | `AC-E01.1`, `AC-E01.2`, `AC-NF.1` | `micro` |
| **E02** | Molde Agnóstico por Natureza | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Cadeia Sexteto e Verbos Derivados | 3 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `micro` |
| **E04** | Gate Anti-Reimplementação | 1 | `AC-E04.1`, `AC-E04.2` | `micro` |
| **E05** | Varredura dos Outros Verbos | 2 | `AC-E05.1`, `AC-E05.2` | `micro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" é igual ao conjunto de `AC-*` da SPEC.

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Contrato de I/O declarado por peça | `AC-E01.1`, `AC-E01.2`, `AC-NF.1` | `TICK-01.1`, `TICK-01.2` |
| `R2` | Molde agnóstico por natureza de nó | `AC-E02.1`, `AC-E02.2` | `TICK-02.1`, `TICK-02.2` |
| `R3` | Cadeia Sexteto e verbos derivados | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `TICK-03.1`, `TICK-03.2`, `TICK-03.3` |
| `R4` | Gate anti-reimplementação | `AC-E04.1`, `AC-E04.2` | `TICK-04.1` |
| `R5` | Varredura dos outros 6 verbos | `AC-E05.1`, `AC-E05.2` | `TICK-05.1`, `TICK-05.2` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Contrato de I/O da Peça

#### Ticket `TICK-01.1` — Contrato de dados, validação e path jail
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `git_worktree`
  * **Executor:** `orchestrix_hermetic` (subprocesso estéril)
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` (backoff `1000ms`, fallback model `true`)
  * **Orçamento de Tokens:** Entrada `16000` tokens · Teto de saída `600` tokens
  * **Rollback:** `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/processo/processo.go`, `internal/processo/gates.go`
  * **Spec Anchor:** `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E01.2`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/processo/contrato_test.go:TestContratoValidate`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/processo/contrato_peca.go`, `internal/processo/contrato_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/processo/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/processo/contrato_peca.go') && file_exists('internal/processo/contrato_test.go')`
* **Spec Ref:** `AC-E01.2`, `AC-NF.1`
* **Escopo (opcional):** Tipos `Contrato`, `Entrada`, `Saida`, `Validate`, `ContratoDe` (carregamento do `contrato.json` do cofre). Path jail: `caminho_padrao` reprovado se `..`, absoluto ou `!filepath.IsLocal`. Erro sentinela `ErrContratoInvalido`. Mutação obrigatória: invalidar `Validate` deve fazer o teste falhar.

#### Ticket `TICK-01.2` — Derivação de I/O e gate por contrato no motor
* **Iniciativa:** `I01` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `git_worktree` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fabricar/natureza.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/fabricar/natureza_test.go:TestArtefatoDeContrato`
* **Outputs:** `internal/fabricar/natureza.go`, `internal/fabricar/natureza_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fabricar/natureza_test.go')`
* **Spec Ref:** `AC-E01.1`
* **Escopo (opcional):** `etapaDe`/`ArtefatoDe` passam a preferir o contrato da peça (`Produz.Caminho` = `caminho_padrao`, `Gate` = `gate_sugerido`) e a cair no molde atual quando ausente. Peça sem contrato não muda comportamento (regressão).

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Contrato de I/O validado? Autoriza iniciar o Épico 02 (molde agnóstico)?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Molde Agnóstico por Natureza

#### Ticket `TICK-02.1` — Nó por natureza (peça | subprocesso | processo)
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `git_worktree` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/processo/processo.go`, `internal/fabricar/natureza.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/fabricar/natureza_test.go:TestNoPorNatureza`
* **Outputs:** `internal/processo/processo.go`, `internal/fabricar/natureza.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/processo/... ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`
* **Escopo (opcional):** `Etapa` ganha campo `Natureza` (peca|subprocesso|processo) derivado do nó conforme `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md:60`: `Compoe` ≠ `""` ⇒ `processo` (composição); `Contrato.Tipo` == `processo` ⇒ `subprocesso` (peça que embrulha processo); senão `peca` (executável). Artefato derivado do `Sai` do contrato por natureza; molde atual como default.

#### Ticket `TICK-02.2` — Vocabulário aberto por acréscimo
* **Iniciativa:** `I01` | **Épico:** `E02` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/matriz/vocab.go`, `internal/fabricar/natureza.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E02.2`
* **Alvo de Teste TDD:** `internal/fabricar/natureza_test.go:TestVocabularioPorAcrescimo`
* **Outputs:** `internal/matriz/vocab.go`, `internal/fabricar/natureza.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/matriz/... ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`
* **Escopo (opcional):** termos do cofre somam-se aos da semente `vocab.default.json`; termo novo não é descartado por `confere`; semente nunca é reescrita.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Molde agnóstico por natureza validado? Autoriza iniciar o Épico 03 (Sexteto e verbos derivados)?')`
- [ ] 100% dos tickets de E02 executados com sucesso.
- [ ] Prova de etapa com artefato/gate do contrato apresentada.
- [ ] Autorização explícita concedida.

---

### Épico 03 — Cadeia Sexteto e Verbos Derivados

#### Ticket `TICK-03.1` — Sexteto completo para qualquer natureza
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `git_worktree` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fabricar/materialize.go`, `internal/fabricar/template.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E03.1`
* **Alvo de Teste TDD:** `internal/fabricar/materialize_test.go:TestSextetoCompleto`
* **Outputs:** `internal/fabricar/materialize.go`, `internal/fabricar/template.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.1`
* **Escopo (opcional):** o plano de gravação lista os 6 artefatos do Sexteto (INTERVIEW, INTENT, SPEC, PLAN, ROUTE/DAG, TEMPLATE/REPORT) com caminhos ISO `YYYY-MM-DD_*`, para qualquer natureza.

#### Ticket `TICK-03.2` — refinar e compor universais
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `git_worktree` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/refinar/adicionar.go`, `internal/compor/compor.go`, `internal/cli/refinar.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E03.2`
* **Alvo de Teste TDD:** `internal/refinar/refinar_test.go:TestRefinarSai`
* **Outputs:** `internal/refinar/`, `internal/compor/`, `internal/cli/refinar.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/refinar/... ./internal/compor/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.2`
* **Escopo (opcional):** `refinar` aceita `--sai`/`--entra`; `compor` aceita nó de qualquer natureza; sem mudança pedida, `ErrMudanca` preservado.

#### Ticket `TICK-03.3` — fechar parametrizável por natureza
* **Iniciativa:** `I01` | **Épico:** `E03` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fechar/fechar.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E03.3`
* **Alvo de Teste TDD:** `internal/fechar/fechar_test.go:TestFecharPorNatureza`
* **Outputs:** `internal/fechar/fechar.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/fechar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.3`
* **Escopo (opcional):** etapas Go-cêntricas ficam condicionais por natureza; cadeia restante inalterada; `--dry-run` prova a omissão.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Sexteto e verbos derivados validados? Autoriza iniciar o Épico 04 (gate anti-reimplementação)?')`
- [ ] 100% dos tickets de E03 executados com sucesso.
- [ ] Dry-run do `fabricar` com Sexteto completo apresentado.
- [ ] Autorização explícita concedida.

---

### Épico 04 — Gate Anti-Reimplementação

#### Ticket `TICK-04.1` — Gate de reimplementação por compatibilidade de I/O
* **Iniciativa:** `I01` | **Épico:** `E04` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E03"]`
* **Isolamento:** Workspace `git_worktree` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/processo/contrato_peca.go`, `internal/canon` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E04.1`
* **Alvo de Teste TDD:** `internal/processo/reimplementacao_test.go:TestGateReprovaDuplicata`
* **Outputs:** `internal/processo/reimplementacao.go`, `internal/processo/reimplementacao_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/processo/...`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/processo/reimplementacao.go')`
* **Spec Ref:** `AC-E04.1`, `AC-E04.2`
* **Escopo (opcional):** `ErrReimplementacao` quando `entra`/`sai` de peça nova são estruturalmente iguais aos de peça do cofre (compatibilidade estrutural, não hash); peça com I/O distinto passa.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 04 (`SG-E04`)
* **Gate (DSL):** `prompt('Gate anti-reimplementação validado? Autoriza iniciar o Épico 05 (varredura dos verbos)?')`
- [ ] Ticket de E04 executado com sucesso.
- [ ] Prova de reprovação de duplicata apresentada.
- [ ] Autorização explícita concedida.

---

### Épico 05 — Varredura dos Outros Verbos

#### Ticket `TICK-05.1` — matriz: aptidão por compatibilidade de I/O
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `micro` | **Tipo:** `unit_no` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E04"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-codigo-128k` · Effort `medium` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/matriz/rank.go`, `internal/fabricar/merger_fit.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E05.1`
* **Alvo de Teste TDD:** `internal/matriz/rank_test.go:TestAptidaoPorIO`
* **Outputs:** `internal/matriz/rank.go`, `internal/fabricar/merger_fit.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/matriz/... ./internal/fabricar/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.1`
* **Escopo (opcional):** `Aptidao` considera compatibilidade `Sai`→`Consome` quando há contrato; sem contrato, comportamento atual preservado.

#### Ticket `TICK-05.2` — listar e remover com contrato
* **Iniciativa:** `I01` | **Épico:** `E05` | **Nature:** `nano` | **Tipo:** `unit_no` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-05.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `t-chat-128k` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/listar/listar.go`, `internal/cli/remover.go` · Spec Anchor `docs/specs/2026-11-16_SPEC_FABRICAR_AGNOSTICO.md#AC-E05.2`
* **Alvo de Teste TDD:** `internal/listar/listar_test.go:TestListarPecaComContrato`
* **Outputs:** `internal/listar/listar.go`, `internal/cli/remover.go`, testes
* **Verificação (argv, sem shell):** `go test -v ./internal/listar/... ./internal/cli/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E05.2`
* **Escopo (opcional):** `listar` expõe resumo do contrato; `remover` aceita peça com contrato mantendo rollback `.legoloop-bak` e recusa de nome fora do padrão.

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.
3. Taxonomia: `go run ./tools/taxlinter internal cmd loops skills docs/templates`
4. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
5. Emissão do relatório formal: `docs/reports/2026-11-16_REPORT_FABRICAR_AGNOSTICO.md`
6. Registro de histórico da sessão via `faz-session --id ab938316-779d-4aa8-9045-fb0d43336880 --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*
*Aprovado este documento e gerado o `docs/routes/2026-11-16_ROUTE_FABRICAR_AGNOSTICO.json`, o BUILD pode ser iniciado.*