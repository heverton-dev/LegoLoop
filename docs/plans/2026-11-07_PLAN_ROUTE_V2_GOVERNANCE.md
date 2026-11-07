# 🗺️ PLAN: GOVERNANCA_ROUTE_V2_CONCORRENCIA

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** Antigravity & Engenharia LegoLoop  
> **Iniciativa:** I01 — Governança e Resiliência da Tríade Executável  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_ROUTE_V2_GOVERNANCE.md](../intents/2026-11-07_INTENT_ROUTE_V2_GOVERNANCE.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md](../specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_ROUTE_V2_GOVERNANCE.json](../routes/2026-11-07_ROUTE_ROUTE_V2_GOVERNANCE.json)  
> **Template:** TEMPLATE_PLAN.md v2.0.0  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (nesta fase de bootstrap do motor, agentes de harness estão desligados; execução estéril em Go puro)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Templates Canônicos e Schemas da ROUTE v2.1 | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Cálculo Criptográfico Merkle e Validação de Integridade | 2 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3` | `micro` |
| **E03** | Governança de Concorrência e Teto de Subagentes Limpos | 1 | `AC-E03.1`, `AC-E03.2`, `AC-NF.1` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Templates Canônicos e Schemas da ROUTE v2.1

#### Ticket `TICK-01.1` — Atualização do TEMPLATE_ROUTE.json para v2.1.0
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modo `suggested` · Primário `9router:code-pro` (high) · Fallback `antigravity:flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:**
  * **Source Files:** `docs/templates/TEMPLATE_ROUTE.json`
  * **Spec Anchor:** `docs/specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md#AC-E01.1`
* **Outputs:**
  * **Artifacts Produzidos:** `docs/templates/TEMPLATE_ROUTE.json`
  * **Verificação (argv, sem shell):** `go test -v -count=1 ./docs/templates/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('docs/templates/TEMPLATE_ROUTE.json')`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — Sincronização de templates.go e testes de conformidade
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Primário `9router:code-nano` (low) · Fallback `antigravity:flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/templates/templates.go`, `docs/templates/templates_test.go`
* **Outputs:** `docs/templates/templates.go`, `docs/templates/templates_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 ./docs/templates/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] Templates v2.1.0 e testes de templates passando com exit code 0.
- [ ] Evidências apresentadas ao usuário no chat.
- [ ] Autorização explícita do usuário concedida.

---

### Épico 02 — Cálculo Criptográfico Merkle e Validação de Integridade

#### Ticket `TICK-02.1` — Implementação do Hashing Normalizado e Structs Upstream em plangen
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Primário `9router:code-pro` (high) · Fallback `antigravity:flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/plangen/route.go`, `internal/plangen/route_test.go`
* **Alvo de Teste TDD:** `internal/plangen/route_test.go:TestUpstreamSignatures`
* **Outputs:** `internal/plangen/route.go`, `internal/plangen/route_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 ./internal/plangen/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.1`, `AC-E02.3`

#### Ticket `TICK-02.2` — Implementação do Validador de Integridade com Bloqueio de Boot
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Primário `9router:code-pro` (high) · Fallback `antigravity:flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/plangen/verify.go`, `internal/plangen/verify_test.go`
* **Alvo de Teste TDD:** `internal/plangen/verify_test.go:TestVerifyRouteIntegrity`
* **Outputs:** `internal/plangen/verify.go`, `internal/plangen/verify_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 ./internal/plangen/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E02.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('Evidências do Épico 02 aprovadas? Autoriza iniciar o Épico 03?')`
- [ ] Testes de cálculo de hash e bloqueio de boot por adulteração passando com sucesso.
- [ ] Autorização explícita do usuário concedida.

---

### Épico 03 — Governança de Concorrência e Teto de Subagentes Limpos

#### Ticket `TICK-03.1` — ConcurrencyStrategy com Teto de Subagentes e Degradação Serial
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Modo `suggested` · Primário `9router:code-pro` (high) · Fallback `antigravity:flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites:** Timeout `180s` · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/plangen/route.go`, `internal/plangen/route_test.go`
* **Alvo de Teste TDD:** `internal/plangen/route_test.go:TestConcurrencyStrategy`
* **Outputs:** `internal/plangen/route.go`, `internal/plangen/route_test.go`
* **Verificação (argv, sem shell):** `go test -v -count=1 ./internal/plangen/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.1`, `AC-E03.2`, `AC-NF.1`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 03 (`SG-E03`)
* **Gate (DSL):** `prompt('Evidências do Épico 03 aprovadas? Autoriza prosseguir para homologação e fechamento?')`
- [ ] ConcurrencyStrategy validada e testada.
- [ ] Autorização explícita do usuário concedida.

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem:
1. Suíte de regressão global: `go test -v -count=1 ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída estritamente vazia).
3. Linter de funções enxutas: `go run ./tools/funclen internal` (todas as funções $\le 35$ linhas).
4. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
5. Emissão do relatório formal: `docs/reports/2026-11-07_REPORT_ROUTE_V2_GOVERNANCE.md`
6. Registro de histórico da sessão via `faz-session --id 6353f637-cb0a-41a1-92ee-8c1400e75faa --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/2026-11-07_ROUTE_ROUTE_V2_GOVERNANCE.json`, o BUILD pode ser iniciado.*
