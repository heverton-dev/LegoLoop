---
id: PLAN-I7
status: DRAFT
date: 2026-11-13
schema_version: 2
upstream: [SPEC-I7]
covers: [AC-I7-01.1, AC-I7-01.2, AC-I7-01.3, AC-I7-02.1, AC-I7-02.2, AC-I7-03.1]
---

# 🗺️ PLAN: Combos Honestos do 9Router (Iniciativa I7)

> **Status:** Em Planejamento  
> **Data:** 2026-11-13  
> **Versão:** 1.0.0  
> **Autor:** i7-architect (Arquiteto I7)  
> **Iniciativa:** I7 — Combos Honestos do 9Router  
> **Referência da Intenção:** [docs/intents/2026-11-13_INTENT_COMBOS_HONESTOS_9ROUTER.md](../intents/2026-11-13_INTENT_COMBOS_HONESTOS_9ROUTER.md)  
> **Referência da Especificação:** [docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md](../specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-13_ROUTE_COMBOS_HONESTOS_9ROUTER.json](../routes/2026-11-13_ROUTE_COMBOS_HONESTOS_9ROUTER.json)  
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
| **E01** | Catálogo e Validação Go | 2 | `AC-I7-01.1`, `AC-I7-01.2`, `AC-I7-01.3` | `micro` |
| **E02** | Sincronização de Combos | 2 | `AC-I7-02.1`, `AC-I7-02.2` | `micro` |
| **E03** | Fechamento | 1 | `AC-I7-03.1` | `nano` |

### Rastreabilidade Requisito para Critério e Ticket

| R | Requisito | AC | Ticket |
| :--- | :--- | :--- | :--- |
| `R1` | Catálogo canônico dos 6 combos honestos em Go | `AC-I7-01.1` | `TICK-I7-01.1` |
| `R2` | Proteção contra erro 413 e barreira de tokens | `AC-I7-01.2`, `AC-I7-01.3` | `TICK-I7-01.2` |
| `R3` | Sincronização do manifesto e compatibilidade legada | `AC-I7-02.1`, `AC-I7-02.2`, `AC-I7-03.1` | `TICK-I7-02.1`, `TICK-I7-02.2`, `TICK-I7-03.1` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Catálogo e Validação Go

#### Ticket `TICK-I7-01.1` — Definição canônica dos 6 combos honestos em Go
* **Iniciativa:** `I7` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Modelo `sonnet` · Effort `high`
  * **Fallback:** Harness `claude` · Modelo `haiku`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` (backoff `1000ms`, fallback model `true`)
  * **Orçamento de Tokens:** Entrada `16000` tokens · Teto de saída `600` tokens
  * **Rollback:** `stash_changes`
* **Inputs:**
  * **Source Files:** `docs/plans/2026-11-13_RASCUNHO_COMBOS_9ROUTER.md`, `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md`
  * **Spec Anchor:** `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md#AC-I7-01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/combo/combo_test.go:TestCatalogoHonesto`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/combo/combo.go`, `internal/combo/combo_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/combo/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/combo/combo.go') && file_exists('internal/combo/combo_test.go')`
* **Spec Ref:** `AC-I7-01.1`
* **Escopo:** Implementar os tipos `MembroCombo` e `ComboHonesto` e o mapa de 6 combos certificados com seus tetos e propriedades de ferramentas.

#### Ticket `TICK-I7-01.2` — Validador de barreira mecânica contra HTTP 413
* **Iniciativa:** `I7` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-I7-01.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Modelo `sonnet` · Effort `high`
  * **Fallback:** Harness `claude` · Modelo `haiku`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/combo/combo.go`
  * **Spec Anchor:** `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md#AC-I7-01.2`
  * **Context Artifacts:** `["internal/combo/combo.go"]`
* **Alvo de Teste TDD:** `internal/combo/validar_test.go:TestValidarRequisicaoCombo`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/combo/validar.go`, `internal/combo/validar_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/combo/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/combo/validar.go')`
* **Spec Ref:** `AC-I7-01.2`
* **Escopo:** Criar a barreira de validação que recusa chamadas acima de 8.000 tokens no `t-micro-8k` retornando `ErrTetoTokensExcedido` e rejeita combos sem suporte a ferramentas quando exigido.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-I7-E01`)
* **Gate (DSL):** `prompt('Catálogo dos 6 combos e barreira mecânica de tokens validados com sucesso? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários TDD de barreira de tokens apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o Épico 02.

---

### Épico 02 — Sincronização de Combos

#### Ticket `TICK-I7-02.1` — Gerador de manifesto JSON para o 9Router
* **Iniciativa:** `I7` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-I7-E01"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Modelo `sonnet` · Effort `medium`
  * **Fallback:** Harness `claude` · Modelo `haiku`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/combo/combo.go`
  * **Spec Anchor:** `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md#AC-I7-02.1`
  * **Context Artifacts:** `["internal/combo/combo.go"]`
* **Alvo de Teste TDD:** `internal/combo/manifesto_test.go:TestGerarManifesto9Router`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/combo/manifesto.go`, `internal/combo/manifesto_test.go`, `scratch/combos_9router_manifest.json`
  * **Verificação (argv, sem shell):** `go test -v -run TestGerarManifesto9Router ./internal/combo/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('scratch/combos_9router_manifest.json')`
* **Spec Ref:** `AC-I7-02.1`
* **Escopo:** Gerar o arquivo JSON de sincronização contendo a definição formal dos 6 combos no padrão aceito pelo gateway 9Router.

#### Ticket `TICK-I7-02.2` — Mapeamento e resolução de retrocompatibilidade
* **Iniciativa:** `I7` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-I7-02.1"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Modelo `sonnet` · Effort `medium`
  * **Fallback:** Harness `claude` · Modelo `haiku`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/combo/combo.go`, `internal/brain/brain.go`
  * **Spec Anchor:** `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md#AC-I7-02.2`
  * **Context Artifacts:** `["internal/combo/combo.go"]`
* **Alvo de Teste TDD:** `internal/combo/legado_test.go:TestMapeamentoLegado`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/combo/legado.go`, `internal/combo/legado_test.go`
  * **Verificação (argv, sem shell):** `go test -v -run TestMapeamentoLegado ./internal/combo/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/combo/legado.go')`
* **Spec Ref:** `AC-I7-02.2`
* **Escopo:** Garantir resolução transparente de nomes legados (`code-pro`, `code-fast`, `code-free`) para os novos combos honestos sem modelos Groq ou Cloudflare de janela curta como primeiro recurso.

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-I7-E02`)
* **Gate (DSL):** `prompt('Manifesto exportado e compatibilidade de legados testada? Autoriza iniciar o Fechamento?')`
- [ ] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Arquivo `scratch/combos_9router_manifest.json` gerado e inspecionado.
- [ ] Autorização explícita do usuário concedida para iniciar o Épico 03.

---

### Épico 03 — Fechamento

#### Ticket `TICK-I7-03.1` — Validação mecânica global e auditoria de qualidade
* **Iniciativa:** `I7` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `script` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-I7-E02"]`
* **Isolamento:**
  * **Workspace:** `local`
  * **Executor:** `orchestrix_hermetic`
  * **Subagentes Permitidos:** `false`
  * **Write Boundary:** `strict`
  * **Extreme:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `claude` · Modelo `haiku` · Effort `low`
  * **Fallback:** Harness `claude` · Modelo `sonnet`
* **Tooling:**
  * **Runtime:** `powershell`
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]`
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/combo/...`
  * **Spec Anchor:** `docs/specs/2026-11-13_SPEC_COMBOS_HONESTOS_9ROUTER.md#AC-I7-03.1`
  * **Context Artifacts:** `["scratch/combos_9router_manifest.json"]`
* **Alvo de Teste TDD:** Suíte global Go
* **Outputs:**
  * **Artifacts Produzidos:** `scratch/i7_validacao_final.json`
  * **Verificação (argv, sem shell):** `go test -count=1 ./...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('scratch/i7_validacao_final.json')`
* **Spec Ref:** `AC-I7-03.1`
* **Escopo:** Executar a suíte de testes completa, auditoria estática `go vet`, checagem de funclen e validação de governança.

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test -count=1 ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia).
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Validação do DAG: `go run ./cmd/legoloop validate docs/dags/2026-11-13_DAG_COMBOS_HONESTOS_9ROUTER.json`
5. Emissão do relatório formal: `docs/reports/2026-11-13_REPORT_COMBOS_HONESTOS_9ROUTER.md`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/2026-11-13_ROUTE_COMBOS_HONESTOS_9ROUTER.json`, o BUILD pode ser iniciado.*
