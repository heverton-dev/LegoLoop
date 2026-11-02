# 📐 Diretório de Templates Canônicos do LegoLoop

Este diretório centraliza os modelos oficiais (**Single Source of Truth**) para todos os artefatos de governança e orquestração do ecossistema LegoLoop.

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                       A CASCATA CANÔNICA EM 4 ESTÁGIOS                     │
├──────────────┬──────────────┬──────────────┬────────────────────────────────┤
│ 1. INTENT    │ 2. SPEC      │ 3. PLAN      │ 4. ROUTE                       │
│ "O Porquê"   │ "O Quê"      │ "O Como"     │ "A Máquina"                    │
│ TEMPLATE_    │ TEMPLATE_    │ TEMPLATE_    │ TEMPLATE_                      │
│ INTENT.md    │ SPEC.md      │ PLAN.md      │ ROUTE.json                     │
└──────────────┴──────────────┴──────────────┴────────────────────────────────┘
```

---

## 1. Contrato do `TEMPLATE_INTENT.md`

O arquivo [`TEMPLATE_INTENT.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/templates/TEMPLATE_INTENT.md) padroniza a captura da dor do usuário e a intenção de negócio:
* **Escopo Inegociável:** Separação estrita entre o que *ESTÁ* no escopo e o que *NÃO ESTÁ* (evitando escopo rastejante / scope creep).
* **Mapeamento Preliminar de Blocos:** Identifica de largada se a demanda precisará de skills especializadas, subloops de terceiros ou apenas scripts de compilação.
* **Governança & Restrições Inegociáveis:** Define o `Trigger` operacional (manual, cron, file_change, loop_end), o nível de risco/blast radius, a política de privacidade (nuvem vs local Ollama) e o contrato de evidência factual.
* **Gate de Transição Obrigatório:** Nenhum documento avança para `SPEC.md` sem aprovação humana expressa do *Intent*.

---

## 2. Contrato do `TEMPLATE_ROUTE.json`

O arquivo [`TEMPLATE_ROUTE.json`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/templates/TEMPLATE_ROUTE.json) define a **Matriz Determinística Canônica** para rotas de execução autônoma no LegoLoop.

### Regras Invioláveis da Matriz de Nós:
1. **Determinismo Estrutural:** Todos os nós no array `nodes` possuem rigorosamente as mesmas chaves (matriz isomórfica).
2. **Sem Omissão de Campos:** Parâmetros que não se aplicam a um tipo de nó recebem `null` ou `[]` (nunca são omitidos).
3. **Isolamento em Duas Dimensões:** Cada nó define seu próprio `workspace` (onde os arquivos ficam) e seu `executor` (como o processo roda), com herança do `default_isolation` global.
4. **Fallback Duplo (BDD - Resiliência Total):** O fallback suporta contingência tanto de **Modelo** quanto de **Harness** (`{"harness": "antigravity", "model": "flash"}`).
5. **Isolamento de Ferramentas sob Demanda:** Se `mcps_required` for `false`, `mcps_allowed` deve ser `[]`, garantindo 0 queima de tokens de ferramentas.
6. **Gates Determinísticos:** Todo nó possui um critério de validação binária testável via terminal ou expressão lógica.
7. **Recompilação Obrigatória:** Toda rota que altera código em Go deve ter um nó de fechamento (`stop_gate`) com a recompilação fresca do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`.

---

## 2. As Duas Camadas de Isolamento (`node.isolation`)

```text
┌─────────────────────────────────────────────────────────────┐
│                 CAMADA 1: WORKSPACE SANDBOX                 │
│                 "Onde os arquivos existem?"                 │
├─────────────────┬───────────────────────────────────────────┤
│ local           │ Árvore do projeto (inline). Padrão seguro.│
│ git_worktree    │ Branch Git temporária (para alto risco).  │
│ tempdir         │ Pasta temporária na memória (descartável).│
├─────────────────┴───────────────────────────────────────────┤
│                 CAMADA 2: EXECUTOR DISPATCH                 │
│                 "Como o processo é acionado?"               │
├─────────────────┬───────────────────────────────────────────┤
│ orchestrix_     │ Subprocesso CLI estéril com contexto zero │
│ hermetic        │ (2k tokens). Padrão de máxima economia.   │
│ local_process   │ Comando local de terminal (go test/build).│
│ subagent        │ Thread interna rodando contida na redoma. │
└─────────────────┴───────────────────────────────────────────┘
```

---

## 3. Dicionário de Campos da Matriz Canônica

| Campo | Tipo | Valores Válidos / Formato | Descrição |
| :--- | :--- | :--- | :--- |
| `id` | `string` | `TICK-XX.X` ou `SG-XX` | Identificador único do nó no DAG. |
| `label` | `string` | Texto descritivo conciso | Nome legível da ação executada pelo nó. |
| `nature` | `string` | `nano`, `micro`, `macro`, `spike` | Granularidade de esforço cognitivo do ticket. |
| `type` | `string` | `unit_task`, `unit_fix`, `stop_gate`, `skill`, `subloop`, `audit` | Propósito funcional do nó no pipeline. |
| `ref` | `string\|null` | Ex: `skills/aidd-spec` ou caminho de rota | Caminho da Skill ou do Subloop invocado. |
| `status` | `string` | `PENDING`, `IN_PROGRESS`, `COMPLETED`, `FAILED` | Estado factual da execução do nó. |
| `execution_mode`| `string` | `tdd_red_green`, `sequential`, `parallel` | Modalidade de concorrência e ciclo de teste. |
| `isolation.workspace`| `string` | `local`, `git_worktree`, `tempdir` | Local físico onde os arquivos são modificados. |
| `isolation.executor` | `string` | `orchestrix_hermetic`, `local_process`, `subagent` | Mecanismo de acionamento do processo/IA. |
| `isolation.subagents_allowed`| `boolean` | `true`, `false` | Se a IA pode fatiar tarefas internamente na redoma. |
| `timeout_sec` | `integer` | Ex: `120`, `180`, `300` | Teto de execução em segundos para abortar loops. |
| `token_budget` | `object` | `{"max_output_tokens": 600}` | Teto severo de resposta (Lei #3 do AGENTS.md). |
| `depends_on` | `array[string]` | Ex: `["TICK-01.1"]` | Nós predecessores obrigatórios para liberação deste nó. |
| `routing.mode` | `string` | `suggested`, `manual`, `auto` | Modo de despacho do 9Router. |
| `routing.harness` | `string` | `9router`, `antigravity`, `claude-code`, `ollama`, `local` | O harness/ambiente do SO encarregado da inferência. |
| `routing.model` | `string` | `code-pro`, `code-fast`, `sonnet`, `flash`, `script/local` | Modelo cognitivo primário. |
| `routing.effort`| `string\|null` | `low`, `medium`, `high` ou `null` | Nível de raciocínio profundo solicitado. |
| `routing.fallback`| `object\|null` | `{"harness": "antigravity", "model": "flash"}` | **Fallback Duplo:** Harness e Modelo de contingência. |
| `tooling.runtime`| `string\|null` | `powershell`, `bash`, `node` ou `null` | Runtime local caso o nó execute scripts. |
| `tooling.mcps_required` | `boolean` | `true`, `false` | Se o nó necessita de MCPs ativos. |
| `tooling.mcps_allowed` | `array[string]` | Ex: `["codebase-memory-mcp:search_code"]` | Lista cirúrgica de ferramentas MCP autorizadas. |
| `inputs.source_files` | `array[string]` | Caminhos relativos | Arquivos exatos que o nó tem permissão de inspecionar. |
| `inputs.spec_anchor` | `string\|null` | Ex: `docs/specs/SPEC.md#secao` | Seção exata da especificação que ancora a regra. |
| `inputs.context_artifacts` | `array[string]` | Caminhos de relatórios / scratch | Dados transitórios produzidos por nós anteriores. |
| `test_target` | `string\|null` | Ex: `caminho_test.go:TestNome` | Teste unitário exato que valida o nó. |
| `defects` | `array[string]` | Ex: `["BUG-01"]` | Lista de IDs de defeitos solucionados pelo ticket. |
| `rollback_strategy` | `string` | `abort_pipeline`, `stash_changes`, `retry_3x` | Ação de contingência caso o gate falhe. |
| `output.artifacts` | `array[string]` | Caminhos gerados/alterados | Arquivos de saída entregues pelo nó. |
| `output.gate` | `string` | Expressão booleana / comando de terminal | Condição binária obrigatória de aprovação (Zero Stubs). |
