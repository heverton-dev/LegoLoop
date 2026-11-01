# 🗺️ PLAN: Evolução do Motor Go com Harness Engineering

> **Status:** Proposta de Plano Único Hierárquico para Validação Humana (Stop-Gate 02)  
> **Data:** 2026-11-01  
> **Especificação Base:** [docs/specs/2026-11-01_SPEC_HARNESS_ENGINEERING_GO.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-01_SPEC_HARNESS_ENGINEERING_GO.md)  
> **Referência de Intenção:** [docs/intents/2026-11-01_INTENT_HARNESS_ENGINEERING_GO.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_HARNESS_ENGINEERING_GO.md)

---

## 🏛️ ÉPICO 01: Schema Enforcement (DAG, Builtin DSL & Runner Check)
> **Objetivo:** Garantir contratos estruturais rígidos para saídas de nós e pipelines, suportando declaração nativa no DAG e função builtin de validação na DSL de gates.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-01.1: Declaração de `output.schema` no DAG e Validação Sintática [CONCLUÍDO]
* **Tipo:** Micro-ticket (Modelagem de Dados e Parser)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/dag/types.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/types.go), [`internal/dag/parse.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/parse.go), [`internal/dag/dag_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/dag_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Adicionar campo `Schema string `json:"schema,omitempty"` no `HandoffContract`.
  2. Validar que caminhos de schema respeitem `filepath.IsLocal` e não escapem do workdir.
  3. Teste TDD unitário cobrindo DAGs com schemas válidos e inválidos.

### 🔹 Micro-ticket TICK-01.2: Implementação do Builtin `schema_valid` na DSL de Gates [CONCLUÍDO]
* **Tipo:** Micro-ticket (DSL de Expressões)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/gate/builtins.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/builtins.go), [`internal/gate/gate_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/gate_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Registrar `schema_valid` em `builtins()` com assinatura `(kStr, kStr) -> kBool`.
  2. Implementar validador de conformidade JSON (validação estrutural de tipos e propriedades obrigatórias contra o arquivo schema).
  3. Testes TDD: arquivos válidos retornam true; arquivos incompatíveis ou JSON quebrado retornam false.

### 🔹 Micro-ticket TICK-01.3: Enforcement Automático de Schema no Runner [CONCLUÍDO]
* **Tipo:** Micro-ticket (Ciclo de Vida do Runner)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/runner/runner.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner.go), [`internal/runner/runner_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. No término da execução do nó, se `Output.Schema` for declarado, validar o arquivo antes de promover ou liberar o nó como `passed`.
  2. Se a validação falhar, marcar `StatusGateFailed` com mensagem descritiva no `NodeReport.Error`.
  3. Testes TDD cobrindo nós inline e isolados.

### 🔸 Nano-ticket TICK-01.4: Refinamento de Erros e Caminhos Confinados [CONCLUÍDO]
* **Tipo:** Nano-ticket (Polimento)
* **Status:** ✅ Concluído
* **Arquivos Alvo:** [`internal/gate/builtins.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/builtins.go)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Mensagens claras em PT-BR para erros de leitura ou arquivo ausente no `schema_valid`.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 01
* [x] Validação dos testes unitários de Schema Enforcement.
* [ ] Aprovação do usuário para liberação do Épico 02.

---

## 🏛️ ÉPICO 02: Effort Clamping Cognitivo (Tipagem, Roteamento & Telemetria)
> **Objetivo:** Permitir controle explícito do esforço cognitivo (`effort: "low" | "medium" | "high"`) nos nós, propagando aos executores de harness com comportamento no-op tolerante e telemetria nos Notes.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-02.1: Tipagem de `EffortLevel` e Validação no DAG
* **Tipo:** Micro-ticket (Modelagem de Tipos)
* **Arquivos Alvo:** [`internal/dag/types.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/types.go), [`internal/dag/parse.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/parse.go), [`internal/dag/dag_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/dag/dag_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Definir `EffortLevel` (`low`, `medium`, `high`) e campo `Effort` em `RoutingConfig`.
  2. Validador em `dag.Parse` que aceita vazio (omissão) ou valores canônicos, rejeitando strings desconhecidas.
  3. Testes unitários com casos válidos e de erro.

### 🔹 Micro-ticket TICK-02.2: Injeção de Effort nos Harnesses da Orchestrix
* **Tipo:** Micro-ticket (Integração de Executores)
* **Arquivos Alvo:** [`internal/orchestrix/harness.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/orchestrix/harness.go), [`internal/orchestrix/orchestrix.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/orchestrix/orchestrix.go), [`internal/orchestrix/harness_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/orchestrix/harness_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Suporte a repasse de effort para `agy` (ex: flag `--effort`), `omp`, `opencode` e `9Router`.
  2. Fallback gracioso: quando o harness/backend não aceita reasoning budget, registrar nota explicativa em `Result.Notes`.
  3. Testes TDD cobrindo despacho com `effort: "low"` e fallback.

### 🔸 Nano-ticket TICK-02.3: Telemetria de Clamping no `NodeReport`
* **Tipo:** Nano-ticket (Telemetria)
* **Arquivos Alvo:** [`internal/runner/runner.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner.go)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Assegurar que os `Notes` gerados pela Orchestrix são repassados integralmente ao `NodeReport`.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 02
* [ ] Validação dos testes unitários de Effort Clamping.
* [ ] Aprovação do usuário para liberação do Épico 03.

---

## 🏛️ ÉPICO 03: Interactive Gates (Auditoria Interativa & Prompter CLI)
> **Objetivo:** Permitir que gates peçam intervenção ou aprovação humana explícita antes do pipeline prosseguir, integrando interface Prompter no motor runner e flags no CLI.  
> **Modelo Recomendado:** `Claude 3.7 Sonnet` ou `Gemini 2.0 Pro` (`--effort high`)

### 🔹 Micro-ticket TICK-03.1: Interface `Prompter` e Builtin `interactive` na DSL
* **Tipo:** Micro-ticket (Extensão de DSL)
* **Arquivos Alvo:** [`internal/gate/gate.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/gate.go), [`internal/gate/builtins.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/builtins.go), [`internal/gate/gate_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/gate/gate_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Definir `Prompter interface { Prompt(message string) (bool, error) }`.
  2. Adicionar `Prompter` ao `Env` de avaliação.
  3. Builtin `interactive(msg)` e alias `prompt(msg)` chamando `env.Prompter`. Se `env.Prompter == nil`, retornar erro claro ou recusa segura.
  4. Testes TDD mockando o Prompter (retornando true, false ou erro).

### 🔹 Micro-ticket TICK-03.2: Gerenciamento de Prompter e AutoApprove no Runner
* **Tipo:** Micro-ticket (Ciclo de Execução)
* **Arquivos Alvo:** [`internal/runner/runner.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner.go), [`internal/runner/runner_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Configurar `Prompter` e flag `AutoApprove` (`--yes`) em `runner.Options`.
  2. Se `AutoApprove == true`, implementar prompter automático que aprova e anota em `Notes`.
  3. Se o gate for rejeitado pelo usuário, marcar `StatusGateFailed` com `"rejeitado pelo operador"`.
  4. Testes TDD com aprovação, rejeição e auto-approve.

### 🔹 Micro-ticket TICK-03.3: Suporte a Terminal Interativo e Flags no CLI
* **Tipo:** Micro-ticket (Interface CLI)
* **Arquivos Alvo:** [`internal/cli/commands.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/cli/commands.go), [`internal/cli/cli_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/cli/cli_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis & Testes TDD:**
  1. Flags `-y, --yes` e `--non-interactive` adicionadas a `legoloop run`.
  2. Terminal Prompter que imprime a mensagem com destaque e lê `[s/n]` ou `[a/r]` do `Stdin`.
  3. Em ambiente `--non-interactive` sem `--yes`, falhar com aviso de bloqueio.
  4. Testes TDD simulando `Stdin` e flags.

### 🔸 Nano-ticket TICK-03.4: Resiliência contra Context Cancellation
* **Tipo:** Nano-ticket (Robustez)
* **Arquivos Alvo:** [`internal/runner/runner.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/runner.go)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. Se o contexto for cancelado (Ctrl+C ou timeout) durante a espera do prompter, encerrar imediatamente sem travar goroutine órfã.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 03
* [ ] Validação dos testes unitários de Interactive Gates.
* [ ] Aprovação do usuário para liberação da Homologação.

---

## 🏛️ ÉPICO 04: Homologação E2E, Recompilação & Fechamento
> **Objetivo:** Testar a suíte completa de ponta a ponta, recompilar o binário canônico e emitir o relatório de conformidade.  
> **Modelo Recomendado:** `Claude 3.5 Sonnet` ou `Gemini 2.0 Flash` (`--effort medium`)

### 🔹 Micro-ticket TICK-04.1: Teste Integrado E2E do Harness Engineering
* **Tipo:** Micro-ticket (Validação E2E)
* **Arquivos Alvo:** [`internal/runner/harness_engineering_e2e_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/runner/harness_engineering_e2e_test.go)
* **Rota:** `Claude 3.5 Sonnet` / `Gemini 2.0 Flash` (`--effort medium`)
* **Entregáveis:**
  1. Pipeline completo executando com nó contendo `effort: "low"`, `output.schema` validando JSON e gate com `interactive(...)`.

### 🔸 Nano-ticket TICK-04.2: Recompilação Obrigatória e Relatório Final
* **Tipo:** Nano-ticket (Entrega & Auditoria)
* **Arquivos Alvo:** [`bin/legoloop.exe`](file:///C:/Users/trcnologia/Desktop/LegoLoop/bin/legoloop.exe), [`docs/reports/2026-11-01_REPORT_HARNESS_ENGINEERING_GO.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-01_REPORT_HARNESS_ENGINEERING_GO.md)
* **Rota:** `Gemini 2.0 Flash-Lite` (`--effort low`)
* **Entregáveis:**
  1. `go build -o bin/legoloop.exe ./cmd/legoloop`.
  2. Relatório formal de homologação com telemetria e validação dos critérios G1 a G6.

### 🛑 STOP-GATE DE AUDITORIA HUMANA FINAL
* [ ] Suíte 100% verde (`go test ./...`).
* [ ] Binário atualizado no disco.
* [ ] Relatório aprovado pelo usuário.
