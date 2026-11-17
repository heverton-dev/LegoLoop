# 🗺️ PLAN: Motor Canônico Universal de Fechamento e Canonização (legoloop fechar)

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** Ecossistema LegoLoop  
> **Iniciativa:** Canonização Universal do Ciclo de Vida do LegoLoop  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_FECHAR_CANONICO.md](../intents/2026-11-07_INTENT_FECHAR_CANONICO.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md](../specs/2026-11-07_SPEC_FECHAR_CANONICO.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_FECHAR_CANONICO.json](../routes/2026-11-07_ROUTE_FECHAR_CANONICO.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

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
| **E01** | Motor Central `internal/fechar` em Go | 2 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3` | `micro` |
| **E02** | CLI `legoloop fechar` e Skill Universal | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Grafo DAG Acíclico e Acoplamento nos Motores | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Motor Central `internal/fechar` em Go

#### Ticket `TICK-01.1` — Modelagem de Tipos e Validação do Slug
* **Iniciativa:** `Canonização Universal` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E01.1`
* **Alvo de Teste TDD:** `internal/fechar/fechar_test.go:TestValidateAndNormalizeSlug`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/fechar/types.go`, `internal/fechar/fechar_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/fechar -run TestValidateAndNormalizeSlug`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/fechar/types.go')`
* **Spec Ref:** `AC-E01.1`, `AC-E01.2`

#### Ticket `TICK-01.2` — Orquestrador Hermético das 8 Etapas
* **Iniciativa:** `Canonização Universal` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `240s` · Retries `1` · Entrada `16000` tokens · Teto de saída `600` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fechar/types.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E01.3`
* **Alvo de Teste TDD:** `internal/fechar/fechar_test.go:TestExecutarOitoEtapas`
* **Outputs:** `internal/fechar/fechar.go`, `internal/fechar/fechar_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/fechar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/fechar/fechar.go')`
* **Spec Ref:** `AC-E01.3`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Motor internal/fechar testado com 100% verde? Autoriza iniciar CLI e Skill?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — Subcomando CLI e Skill Universal

#### Ticket `TICK-02.1` — Implementação do Comando CLI `legoloop fechar`
* **Iniciativa:** `Canonização Universal` | **Épico:** `E02` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-pro` · Effort `high` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `12000` tokens · Teto de saída `400` tokens · Rollback `stash_changes`
* **Inputs:** `internal/fechar/fechar.go` · Spec Anchor `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E02.1`
* **Alvo de Teste TDD:** `internal/cli/fechar_test.go:TestCmdFechar`
* **Outputs:** `internal/cli/fechar.go`, `internal/cli/fechar_test.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/cli -run TestCmdFechar`
* **Gate (DSL):** `exit_code == 0 && file_exists('internal/cli/fechar.go')`
* **Spec Ref:** `AC-E02.1`

#### Ticket `TICK-02.2` — Criação da Skill Universal `skills/legoloop-fechar/SKILL.md`
* **Iniciativa:** `Canonização Universal` | **Épico:** `E02` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `120s` · Retries `1` · Entrada `6000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Spec Anchor `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E02.2`
* **Outputs:** `skills/legoloop-fechar/SKILL.md`, `skills/embed.go`
* **Verificação (argv, sem shell):** `go test -v ./internal/cli -run TestInstalarSkills`
* **Gate (DSL):** `exit_code == 0 && file_exists('skills/legoloop-fechar/SKILL.md')`
* **Spec Ref:** `AC-E02.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02 (`SG-E02`)
* **Gate (DSL):** `prompt('CLI e Skill legoloop-fechar criados e testados? Autoriza DAG e Acoplamento?')`
- [ ] 100% dos tickets de E02 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 03 — Grafo DAG Acíclico e Acoplamento nos Motores

#### Ticket `TICK-03.1` — Grafo DAG `loops/fechar/pipeline.dag.json`
* **Iniciativa:** `Canonização Universal` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E02"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `120s` · Retries `1` · Entrada `6000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** Spec Anchor `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E03.1`
* **Outputs:** `loops/fechar/pipeline.dag.json`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop validate loops/fechar/pipeline.dag.json`
* **Gate (DSL):** `exit_code == 0 && file_exists('loops/fechar/pipeline.dag.json')`
* **Spec Ref:** `AC-E03.1`

#### Ticket `TICK-03.2` — Acoplamento no `mega-pipeline.loop.json`
* **Iniciativa:** `Canonização Universal` | **Épico:** `E03` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `120s` · Retries `1` · Entrada `6000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `loops/mega-pipeline.loop.json` · Spec Anchor `docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md#AC-E03.2`
* **Outputs:** `loops/mega-pipeline.loop.json`
* **Verificação (argv, sem shell):** `go run ./cmd/legoloop validate loops/mega-pipeline.loop.json`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E03.2`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test -count=1 ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia)
3. Limite de tamanho de função: `go run ./tools/funclen internal` (<= 35 linhas)
4. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`
5. Distribuição universal de skills: `legoloop instalar-skills`
6. Reindexação do grafo semântico via `codebase-memory-mcp cli index_repository`
7. Emissão do relatório formal: `docs/reports/2026-11-07_REPORT_FECHAR_CANONICO.md`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/2026-11-07_ROUTE_FECHAR_CANONICO.json`, o BUILD pode ser iniciado.*
