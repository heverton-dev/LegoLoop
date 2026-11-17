# 🎯 INTENT: Evolução do Motor Go do LegoLoop com Harness Engineering

> **Status:** Proposta para Alinhamento Socrático e Aprovação Humana (Stop-Gate 01)  
> **Data:** 2026-11-01  
> **Autor:** Antigravity AI + Usuário  
> **Alvo:** Motor Go do LegoLoop (`internal/dag/`, `internal/runner/`, `internal/gate/`, `internal/orchestrix/`, `internal/cli/`)

---

## 1. Contexto e Motivação

O LegoLoop orquestra pipelines de IA estéreis, operando nós de skill, script e subloops sob contratos de handoff rígidos. Contudo, para blindar a execução autônoma contra alucinações de saída, desperdício de cotas de raciocínio e transições cegas entre fases críticas, é necessário integrar ao motor Go as 3 propostas canônicas de **Harness Engineering**:

1. **Schema Enforcement:** Impor validação estrutural estrita aos artefatos e aos próprios DAGs, garantindo que saídas geradas obedeçam a schemas definidos (ex: JSON Schema, validações estruturais de campos) antes da promoção de artefatos ou liberação de dependências downstream.
2. **Effort Clamping:** Adicionar ao modelo de roteamento (`RoutingConfig` / `Node`) o controle explícito de esforço cognitivo (`effort: "low" | "medium" | "high"` ou teto de reasoning tokens), propagando essas restrições diretamente aos executores de harness (`agy`, `omp`, `opencode`, `mimo`, `ollama` e `9Router`) para impedir que micro/nano-tickets gastem cotas desnecessárias de deep thinking.
3. **Interactive Gates:** Expandir o sistema de gates (`internal/gate` e `internal/runner`), hoje puramente baseados em expressões booleanas estáticas de arquivos/código de saída, para permitir **Stop-Gates Interativos**. Esses gates suspendem a execução no terminal (ou via API/Studio), solicitando confirmação explícita (`[A]provar / [R]ejeitar`) do operador humano antes de avançar para a próxima camada ou nó crítico.

---

## 2. Objetivos Principais

1. **Schema Enforcement Contratual:**
   * Permitir que o nó defina um schema de saída no `HandoffContract` (ex: `output.schema` apontando para um arquivo de schema ou tipo aceito) e/ou validação do pipeline contra o schema JSON formal do LegoLoop.
   * Validador embutido no runner que rejeita artefatos com `StatusGateFailed` se não cumprirem o schema contratado.
2. **Effort Clamping Nativo:**
   * Adicionar o campo `Effort` (`low`, `medium`, `high`) no `RoutingConfig` de `internal/dag/types.go`.
   * Atualizar os argumentos de invocação dos harnesses em `internal/orchestrix/harness.go` para injetar a flag/parâmetro correspondente de esforço ou tokens de reasoning (ex: flag `--effort` no agy/omp/opencode onde aplicável, ou via parâmetros HTTP no 9router).
3. **Interactive Gates (Auditoria Humana no Runner):**
   * Criar suporte à expressão de gate `interactive()` ou configuração `interactive: true` / `prompt: "..."` no contrato de gate.
   * Quando o runner encontrar um gate interativo, se estiver em ambiente TTY/CLI interativo, pausar e aguardar confirmação do usuário (com suporte a `--yes` / `--non-interactive` para CI/CD automatizado).
   * Emitir eventos de ciclo de vida (`OnGatePrompt` / `OnGateApprove` / `OnGateReject`) para integração com o Studio via SSE.
4. **TDD Estrito & Recompilação:**
   * Cobrir 100% das novas capacidades com testes unitários em Go (Red-Green-Refactor).
   * Recompilar o binário canônico `bin/legoloop.exe` ao final da homologação.

---

## 3. Escopo

### O que entra:
* **`internal/dag/`:**
  * Atualização de `types.go`, `parse.go` e `v2.go` para suportar `effort` no `RoutingConfig`, `schema` no `HandoffContract` e tipagem de gates interativos.
  * Validações semânticas de DAG para rejeitar valores inválidos de esforço e caminhos inexistentes de schema.
* **`internal/gate/`:**
  * Suporte a função `schema_valid(schema_path, file_path)` e/ou função `interactive(prompt)`.
* **`internal/runner/`:**
  * Integração da validação de schema pós-execução e antes do promote.
  * Interrupção / prompt para gates interativos em `runner.Options` com handlers configuráveis (`InteractivePrompt`).
* **`internal/orchestrix/`:**
  * Mapeamento do `effort` para argumentos dos backends (`agy`, `omp`, `opencode`, `mimo`, `ollama`, `9Router`).
* **`internal/cli/`:**
  * Adição de flags de controle como `--non-interactive` / `--yes` para automações e comandos CLI.

### O que NÃO entra nesta etapa:
* Modificação das interfaces visuais estáticas legadas do Studio que não consumam a nova API.
* Modificação de pipelines de terceiros que não solicitem as novas funcionalidades (retrocompatibilidade 100% garantida).

---

## 4. Questões para Validação Socrática (Stop-Gate 01)

Para travarmos esta Intenção e avançarmos para a `SPEC` formal:

1. **Sintaxe do Schema Enforcement no DAG:**
   * Preferência: Você prefere declarar o schema diretamente no handoff (`output: { artifact: "out.json", schema: "schemas/out.schema.json", gate: "..." }`) ou como uma função embutida no gate (`gate: "schema_valid('schemas/out.json', 'out.json') && exit_code == 0"`)?
   *(Recomendação: Ambos os suportes — `output.schema` como contrato de primeira classe e `schema_valid(...)` como função do gate DSL).*
2. **Níveis e Comportamento de Effort Clamping:**
   * Confirmamos os três níveis canônicos: `"low"`, `"medium"`, `"high"` com default `"medium"` para skills quando omitido? Como o harness deve se comportar se o backend local (ex: script puro ou ollama sem reasoning) não suportar a flag `--effort`? (Recomendação: No-op tolerante com warning/telemetria nos `Notes` do `NodeReport`).
3. **Mecanismo de Interactive Gate no CLI:**
   * Em execução via CLI interativa (`legoloop run pipeline.dag.json`), o motor deve pausar a goroutine do nó e ler da entrada padrão `Stdin` (`[A]provar / [R]ejeitar`), respeitando um timeout se configurado?
