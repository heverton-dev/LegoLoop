# 🗺️ PLAN: [Tema do Plano de Execução]

> **Status:** [Em Planejamento | Aprovado pelo Usuário | Em Execução | Concluído]  
> **Data:** YYYY-MM-DD  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Usuário  
> **Referência da Intenção:** [docs/intents/YYYY-MM-DD_INTENT_NOME.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/YYYY-MM-DD_INTENT_NOME.md)  
> **Referência da Especificação:** [docs/specs/YYYY-MM-DD_SPEC_NOME.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/YYYY-MM-DD_SPEC_NOME.md)  
> **Espelho Executável JSON:** [docs/routes/YYYY-MM-DD_ROUTE_NOME.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/YYYY-MM-DD_ROUTE_NOME.json)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (Execução isolada por bloco sem context bloat)
* **Permissão de Subagentes:** `false` (Agentes de harness proibidos; subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** `abort_pipeline` / `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | [Nome do Épico 01] | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | [Nome do Épico 02] | 1 | `AC-E02.1` | `nano` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — [Nome do Épico 01]

#### Ticket `TICK-01.1` — [Descrição Resumida da Ação]
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Isolamento:**
  * **Workspace:** `local` (ou `git_worktree` / `tempdir`)
  * **Executor:** `orchestrix_hermetic` (subprocesso estéril)
  * **Subagentes Permitidos:** `false`
* **Roteamento Cognitivo:**
  * **Modo:** `suggested`
  * **Primário:** Harness `9router` · Modelo `code-pro` · Effort `high`
  * **Fallback:** Harness `antigravity` · Modelo `flash`
* **Tooling:**
  * **Runtime:** `powershell` (ou `node`, `obscura`, etc.)
  * **MCPs Requeridos:** `false` (ou lista de `servidor:ferramenta`)
* **Inputs:**
  * **Source Files:** `internal/pkg/service.go`
  * **Spec Anchor:** `docs/specs/YYYY-MM-DD_SPEC_NOME.md#AC-E01.1`
  * **Context Artifacts:** `[]`
* **Alvo de Teste TDD:** `internal/pkg/service_test.go:TestFeatureCore`
* **Outputs:**
  * **Artifacts Produzidos:** `internal/pkg/service.go`, `internal/pkg/service_test.go`
  * **Mechanical Gate:** `exit_code == 0 && go test -v ./internal/pkg/...`
* **Spec Ref:** `AC-E01.1`

#### Ticket `TICK-01.2` — [Descrição Resumida da Ação]
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential`
* **Dependências:** `["TICK-01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento:** Harness `9router` · Modelo `code-nano` · Effort `low`
* **Tooling:** Runtime `powershell`
* **Inputs:** `internal/pkg/service.go`
* **Outputs:** `internal/pkg/service.go`
* **Mechanical Gate:** `exit_code == 0 && golangci-lint run internal/pkg/...`
* **Spec Ref:** `AC-E01.2`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01
- [ ] 100% dos tickets de E01 executados com sucesso (PASS nos gates mecânicos).
- [ ] Evidências de testes unitários verdes apresentadas ao usuário.
- [ ] Autorização explícita do usuário concedida para iniciar o próximo Épico.

---

## 4. Homologação Final & Fechamento da Demanda

Ao concluir todos os épicos:
1. Execução de suíte de regressão global: `go test ./...`
2. Validação de cobertura e linters: `golangci-lint run ./...`
3. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`
4. Emissão do relatório formal: `docs/reports/YYYY-MM-DD_REPORT_NOME.md`
5. Registro de histórico da sessão via `faz-session --id <ID> --format all`

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima mapeia 100% dos requisitos da SPEC e possui isolamento e gates mecânicos claros?*  
*Aprovado este documento e gerado o `docs/routes/YYYY-MM-DD_ROUTE_NOME.json`, o BUILD pode ser iniciado.*
