# 🗺️ PLAN: Fábrica Universal de Processos Canônicos (legoloop-fabricar)

> **Status:** Proposta para Revisão (UNDER_REVIEW)  
> **Data:** 2026-11-03  
> **Versão:** 1.0.0  
> **Autor:** Antigravity AI + Everton  
> **Referência da Intenção:** [docs/intents/2026-11-03_INTENT_LEGOLOOP_CANONIZER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-03_INTENT_LEGOLOOP_CANONIZER.md)  
> **Referência da Especificação:** [docs/specs/2026-11-03_SPEC_LEGOLOOP_CANONIZER_V2.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-03_SPEC_LEGOLOOP_CANONIZER_V2.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Comando Dedicado da CLI:** `legoloop fabricar <nome>` (Elimina 100% de colisão com o `canonize --root` legado).
* **Engine:** `orchestrix_hermetic` (Execução isolada por bloco sem context bloat).
* **Permissão de Subagentes:** `false` (Subprocessos estéreis sem clonagem de system prompt).
* **TDD Estrito:** `true` (Ciclo Red ➔ Green ➔ Refactor antes de qualquer código de produção).
* **Limite de Workers Paralelos:** `1` (Execução sequencial por dependência estrita).
* **Timeout Padrão por Bloco:** `180s`.
* **Política de Rollback:** `abort_pipeline` / `stash_changes`.

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|:---|:---|:---:|:---|:---:|
| **E01** | Núcleo Go da Fábrica (`internal/canonizer`) | 3 | `AC-E01.1`, `AC-E01.2`, `AC-E01.3`, `AC-E01.4`, `AC-E01.5` | `micro` |
| **E02** | Subcomando CLI `legoloop fabricar` (`internal/cli`) | 2 | `AC-E02.1`, `AC-E02.2`, `AC-E02.3`, `AC-E02.4` | `micro` |
| **E03** | Interface Web & Canvas no Studio (`internal/studio`) | 2 | `AC-E03.1`, `AC-E03.2`, `AC-E03.3` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Núcleo Go da Fábrica (`internal/canonizer`)

#### Ticket `FAB-E01.1` — Sanitização de Nome com Prefixo `legoloop-` e Limite de 64 chars
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.5 Sonnet` / `Qwen-Coder` · Effort `medium`
* **Tooling:** Runtime `powershell` · Allowlist: `[]` (somente compilador Go)
* **Inputs:** `internal/canonizer/name.go`, `internal/canonizer/name_test.go`
* **Alvo de Teste TDD:** `internal/canonizer/name_test.go:TestNormalizeProcessName`
* **Outputs:** `internal/canonizer/name.go`, `internal/canonizer/name_test.go`
* **Mechanical Gate:** `go test -v -run TestNormalizeProcessName ./internal/canonizer/...`
* **Spec Ref:** `AC-E01.1`

#### Ticket `FAB-E01.2` — Scanner Multi-Drive com Timeout e Detecção Agnóstica a SO
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Dependências:** `["FAB-E01.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.5 Sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/canonizer/scanner.go`, `internal/canonizer/scanner_test.go`
* **Alvo de Teste TDD:** `internal/canonizer/scanner_test.go:TestDiscoverDrives`
* **Outputs:** `internal/canonizer/scanner.go`, `internal/canonizer/scanner_test.go`
* **Mechanical Gate:** `go test -v -run TestDiscoverDrives ./internal/canonizer/...`
* **Spec Ref:** `AC-E01.2`, `AC-E01.4`

#### Ticket `FAB-E01.3` — Orquestrador Fino e Geração Tripartida (SKILL.md, .loop.json, ROUTE.json)
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Dependências:** `["FAB-E01.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.7 Sonnet` / `Pro` · Effort `high`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/canonizer/engine.go`, `internal/canonizer/engine_test.go`
* **Alvo de Teste TDD:** `internal/canonizer/engine_test.go:TestCanonizePipeline`
* **Outputs:** `internal/canonizer/engine.go`, `internal/canonizer/engine_test.go`
* **Mechanical Gate:** `go test -v -run TestCanonizePipeline ./internal/canonizer/...`
* **Spec Ref:** `AC-E01.3`, `AC-E01.5`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 01
- [ ] 100% dos testes de `internal/canonizer` passando com exit code 0.
- [ ] Sem quebra na suíte de testes existente (`go test ./...`).
- [ ] Aval humano para avançar para a CLI (Épico 02).

---

### Épico 02 — Subcomando CLI `legoloop fabricar` (`internal/cli`)

#### Ticket `FAB-E02.1` — Implementação do Comando `legoloop fabricar`
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Dependências:** `["FAB-E01.3"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.5 Sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/cli/fabricar.go`, `internal/cli/fabricar_test.go`
* **Alvo de Teste TDD:** `internal/cli/fabricar_test.go:TestRunFabricar`
* **Outputs:** `internal/cli/fabricar.go`, `internal/cli/fabricar_test.go`
* **Mechanical Gate:** `go test -v -run TestRunFabricar ./internal/cli/...`
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`

#### Ticket `FAB-E02.2` — Registro no Dispatcher Central e Skill para o Chat
* **Nature:** `nano` | **Tipo:** `unit_task` | **Modo de Execução:** `sequential`
* **Dependências:** `["FAB-E02.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Gemini 2.0 Flash-Lite` · Effort `low`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/cli/cli.go`, `skills/legoloop-canonizer/SKILL.md`
* **Outputs:** `internal/cli/cli.go`, `skills/legoloop-canonizer/SKILL.md`
* **Mechanical Gate:** `go test -v ./internal/cli/...`
* **Spec Ref:** `AC-E02.3`, `AC-E02.4`

---

### 🛑 STOP-GATE HUMANO — ÉPICO 02
- [ ] Subcomando `legoloop fabricar --help` e execução com `--dry-run` funcionando perfeitamente no terminal.
- [ ] Aval humano para avançar para o Studio Web (Épico 03).

---

### Épico 03 — Interface Web & Canvas no Studio (`internal/studio`)

#### Ticket `FAB-E03.1` — Endpoint HTTP `POST /api/canonizer/run`
* **Nature:** `micro` | **Tipo:** `unit_task` | **Modo de Execução:** `tdd_red_green`
* **Dependências:** `["FAB-E02.2"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.5 Sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/studio/canonizer_api.go`, `internal/studio/canonizer_api_test.go`
* **Alvo de Teste TDD:** `internal/studio/canonizer_api_test.go:TestHandleCanonizerRun`
* **Outputs:** `internal/studio/canonizer_api.go`, `internal/studio/canonizer_api_test.go`
* **Mechanical Gate:** `go test -v -run TestHandleCanonizerRun ./internal/studio/...`
* **Spec Ref:** `AC-E03.3`

#### Ticket `FAB-E03.2` — Modal "+ Novo Processo" e Plotagem Automática no Canvas
* **Nature:** `micro` | **Tipo:** `ui_task` | **Modo de Execução:** `sequential`
* **Dependências:** `["FAB-E03.1"]`
* **Isolamento:** Workspace `local` · Executor `orchestrix_hermetic` · Subagentes `false`
* **Roteamento Cognitivo:** Modelo `Claude 3.5 Sonnet` · Effort `medium`
* **Tooling:** Runtime `powershell` · Allowlist: `[]`
* **Inputs:** `internal/studio/static/app.js`, `internal/studio/static/index.html`, `internal/studio/static/style.css`
* **Outputs:** `internal/studio/static/app.js`, `internal/studio/static/index.html`, `internal/studio/static/style.css`
* **Mechanical Gate:** `go test -v -run "TestZeroEmojisInFrontendAndTemplates|TestNoNativeDialogs" ./internal/studio/...`
* **Spec Ref:** `AC-E03.1`, `AC-E03.2`

---

## 4. Homologação Final & Fechamento da Demanda

Ao concluir os 3 épicos:
1. Execução de suíte de regressão completa: `go test ./...`
2. Recompilação obrigatória do binário: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`
3. Smoke test real no terminal: `legoloop fabricar "teste-smoke" --yes --dry-run`
4. Emissão do relatório formal: `docs/reports/2026-11-03_REPORT_LEGOLOOP_FABRICAR.md`
5. Registro de histórico da sessão via `faz-session --id 1da38997-8926-43dc-a9c7-4c885fc12303 --format all`
