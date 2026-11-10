# 🗺️ PLAN: {{TEMA_DO_PLANO}}

> **Status:** {{Em Planejamento | Aprovado pelo Usuário | Em Execução | Concluído}}  
> **Data:** {{DATA}}  
> **Versão:** {{1.0.0}}  
> **Autor:** {{AUTOR}}  
> **Iniciativa:** {{NOME_OU_ID_DA_INICIATIVA | Opcional para demandas amplas}}  
> **Referência da Intenção:** [docs/intents/{{DATA}}_INTENT_{{NOME}}.md](../intents/{{DATA}}_INTENT_{{NOME}}.md)  
> **Referência da Especificação:** [docs/specs/{{DATA}}_SPEC_{{NOME}}.md](../specs/{{DATA}}_SPEC_{{NOME}}.md)  
> **Espelho Executável JSON:** [docs/routes/{{DATA}}_ROUTE_{{NOME}}.json](../routes/{{DATA}}_ROUTE_{{NOME}}.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

<!--
GUIA DE PREENCHIMENTO (apague este bloco ao concluir)
- Este PLAN e o ROUTE são espelhos 1:1: mesmo conjunto de IDs, dependências, naturezas, limites e roteamento.
  Qualquer divergência entre os dois é defeito. Valores do exemplo abaixo coincidem com o exemplo do ROUTE.
- Cada ticket é uma fatia vertical: entrega comportamento verificável de ponta a ponta, não uma camada.
- Cada ticket cita o critério AC da SPEC que ele fecha (Spec Ref). AC sem ticket = plano incompleto.
- Em demandas amplas, múltiplos Épicos podem ser agrupados opcionalmente sob uma Iniciativa (ex: I01 — Nome).
-->

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `{{nano | micro | epic | autonomous}}`
* **Permissão de Subagentes:** `false` ou `ephemeral` (AGENTS.md, regra 4: agente de harness proibido; subagente efêmero só com as seis condições)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | {{NOME_DO_EPICO_01}} | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | {{NOME_DO_EPICO_02}} | 1 | `AC-E02.1` | `macro` |

> Regra de cobertura: a união da coluna "Critérios de Aceite Cobertos" deve ser igual ao conjunto de `AC-*` da SPEC.

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — {{NOME_DO_EPICO_01}}

#### Ticket `TICK-01.1` — {{Descrição resumida da ação}}
* **Iniciativa:** `{{NOME_OU_ID_DA_INICIATIVA}}` | **Épico:** `E01` | **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green` | **Status:** `PENDING`
* **Dependências:** `[]`
* **Isolamento:**
  * **Workspace:** `local` (ou `git_worktree` / `tempdir`)
  * **Executor:** `orchestrix_hermetic` (subprocesso estéril)
  * **Subagentes Permitidos:** `false` (ou `ephemeral`, com a lista de ferramentas em Tooling > Allowlist)
  * **Write Boundary:** `strict`
  * **Extreme:** `false` (`true` libera Opus, só com aprovação humana)
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell` (ou `node`, `obscura`, etc.)
  * **MCPs Requeridos:** `false` · **Allowlist:** `[]` (formato `servidor:ferramenta`)
* **Limites & Resiliência:**
  * **Timeout:** `180s` · **Retries:** `1` (backoff `1000ms`, fallback model `true`)
  * **Orçamento de Tokens:** Entrada `16000` tokens · Teto de saída `600` tokens
  * **Rollback:** `stash_changes`
* **Inputs:**
  * **Source Files:** `internal/pkg/service.go`
  * **Spec Anchor:** `docs/specs/{{DATA}}_SPEC_{{NOME}}.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/pkg/service_test.go:TestFeatureCore`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/pkg/service.go`, `internal/pkg/service_test.go`
  * **Verificação (argv, sem shell):** `go test -v ./internal/pkg/...`
  * **Gate (DSL):** `exit_code == 0 && file_exists('internal/pkg/service.go') && file_exists('internal/pkg/service_test.go')`
* **Spec Ref:** `AC-E01.1`
* **Escopo (opcional):** {{Detalhe do que o teste do ticket precisa provar, com valores concretos. Use quando o gate sozinho não explica o ticket. Só existe no PLAN: o ROUTE não tem este campo.}}

#### Ticket `TICK-01.2` — {{Descrição resumida da ação}}
* **Iniciativa:** `{{NOME_OU_ID_DA_INICIATIVA}}` | **Épico:** `E01` | **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `low` · Fallback Harness `antigravity` · Modelo `flash`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `false` · Allowlist `[]`
* **Limites & Resiliência:** Timeout `180s` · Retries `1` · Entrada `8000` tokens · Teto de saída `300` tokens · Rollback `stash_changes`
* **Inputs:** `internal/pkg/service.go` · Spec Anchor `docs/specs/{{DATA}}_SPEC_{{NOME}}.md#AC-E01.2`
* **Outputs:** `internal/pkg/service.go`
* **Verificação (argv, sem shell):** `go vet ./internal/pkg/...`
* **Gate (DSL):** `exit_code == 0`
* **Spec Ref:** `AC-E01.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01 (`SG-E01`)
* **Gate (DSL):** `prompt('Evidências do Épico 01 aprovadas? Autoriza iniciar o Épico 02?')`
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

### Épico 02 — {{NOME_DO_EPICO_02}}

#### Ticket `TICK-02.1` — {{Descrição resumida da ação (subloop especialista)}}
* **Iniciativa:** `{{NOME_OU_ID_DA_INICIATIVA}}` | **Épico:** `E02` | **Nature:** `macro` | **Tipo:** `subloop` | **Modo de Execução:** `sequential` | **Status:** `PENDING`
* **Dependências:** `["SG-E01"]`
* **Ref do Subloop:** `exemplo-subloop` (id de um loop do acervo; nunca um caminho)
* **Isolamento:** Workspace `tempdir` · Executor `orchestrix_hermetic` · Subagentes `false` · Write Boundary `strict`
* **Roteamento:** Harness `9router` · Modelo `code-fast` · Effort `medium` · Fallback Harness `ollama` · Modelo `deepseek-coder:6.7b`
* **Tooling:** Runtime `powershell` · MCPs Requeridos `true` · Allowlist `["codebase-memory-mcp:search_code"]`
* **Limites & Resiliência:** Timeout `300s` · Retries `1` · Entrada `32000` tokens · Teto de saída `1200` tokens · Rollback `stash_changes`
* **Inputs:** Context Artifacts `["scratch/target_analysis.json"]` · Spec Anchor `docs/specs/{{DATA}}_SPEC_{{NOME}}.md#AC-E02.1`
* **Outputs:** `scratch/subloop_output.json`
* **Verificação (argv, sem shell):** nenhuma (o subloop se valida pelo próprio gate)
* **Gate (DSL):** `file_exists('scratch/subloop_output.json')`
* **Spec Ref:** `AC-E02.1`

---

## 4. Homologação Final & Fechamento da Demanda (`SG-FINAL`)

Ao concluir todos os épicos, executar nesta ordem (cada comando é um argv, sem shell):
1. Suíte de regressão global: `go test ./...`
2. Estática e formatação: `go vet ./...` e `gofmt -l .` (saída vazia). `golangci-lint run ./...` apenas se instalado e configurado.
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal: `docs/reports/{{DATA}}_REPORT_{{NOME}}.md`
5. Registro de histórico da sessão via `faz-session --id <ID> --format all`

* **Gate (DSL):** `exit_code == 0 && file_exists('bin/legoloop.exe')`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/{{DATA}}_ROUTE_{{NOME}}.json`, o BUILD pode ser iniciado.*
