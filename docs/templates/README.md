# 📐 Diretório de Templates Canônicos do LegoLoop (SSOT)

Este diretório centraliza os modelos oficiais (**Single Source of Truth**) para os artefatos de governança, esteira e peças atômicas do ecossistema LegoLoop.

**O que mudou em relação à V1.** A V1 declarava "Padrão Ouro" e "100% universal" sem prova. A V2 troca promessa por verificação: cada template é aceito pelo componente do repositório que o lê (`dag.Parse`, `gate.Compile`, scanner de `assets`, `loop.Store`), os formatos de peça são os **nativos** de cada harness, e a seção 7 lista os comandos que reproduzem essas verificações. Os templates continuam agnósticos a harness no que é contrato (INTENT, SPEC, PLAN, ROUTE, DAG); onde a peça é, por natureza, de um harness (hook e MCP), o formato é o do harness.

---

## 1. A Cascata de Governança em 4 Estágios

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

1. **[`TEMPLATE_INTENT.md`](TEMPLATE_INTENT.md):** Alinhamento de intenção com a Entrevista Socrática (Addy Osmani, adaptada em PT-BR) em funil até confiança ≥ 95%. Responde "por quê"; não contém contratos.
2. **[`TEMPLATE_SPEC.md`](TEMPLATE_SPEC.md):** Premissas e Boundaries (Addy Osmani) + Costuras de Teste (Matt Pocock) + Critérios Binários BDD `AC-XX` (AIDD-Spec). Todo épico tem ao menos um critério de caminho de falha, e há critérios não funcionais `AC-NF`.
3. **[`TEMPLATE_PLAN.md`](TEMPLATE_PLAN.md):** Fatiamento vertical com 4 níveis de supervisão humana (`nano`, `micro`, `epic`, `autonomous`). É o espelho legível do ROUTE.
4. **[`TEMPLATE_ROUTE.json`](TEMPLATE_ROUTE.json):** Matriz isomórfica de nós para execução por agentes, com isolamento de duas camadas e contingência dupla (harness + modelo). É o espelho executável do PLAN.

> **Regra de espelho.** PLAN e ROUTE têm o mesmo conjunto de IDs, dependências, naturezas, limites (timeout, teto de tokens, rollback) e roteamento. Divergência entre os dois é defeito. Os exemplos dos dois templates coincidem campo a campo e isso é verificado (seção 7).

---

## 2. A Sextilha Canônica de Peças

Cada peça usa o formato que o repositório (ou o harness) realmente lê. Um template que nenhum componente reconhece não é template, é decoração.

| Tipo de Peça | Arquivo de Template | Formato real | Quem lê (no código) |
|:---|:---|:---|:---|
| **1. Skill** | [`TEMPLATE_PIECE_SKILL.md`](TEMPLATE_PIECE_SKILL.md) | `SKILL.md` com frontmatter `name`, `description`, `consumes`, `produces`, `artifact`, `gate`, `suggested_model` | `internal/assets/skills.go`, grafo `internal/graph` |
| **2. Script** | [`TEMPLATE_PIECE_SCRIPT.ps1`](TEMPLATE_PIECE_SCRIPT.ps1) | PowerShell 5.1+, UTF-8 **com BOM** | `internal/assets/scripts.go` (raízes `scripts/`, `tools/`, `.aidd/scripts`, `.claude/scripts`) |
| **3. Hook** | [`TEMPLATE_PIECE_HOOK.json`](TEMPLATE_PIECE_HOOK.json) | Trecho de `settings.json`: `hooks.<Evento>[{matcher, hooks:[{type, command, timeout}]}]` | `internal/assets/hooks.go` |
| **4. MCP** | [`TEMPLATE_PIECE_MCP.json`](TEMPLATE_PIECE_MCP.json) | Trecho de `.mcp.json`: `mcpServers.<nome>{command, args, env}` | `internal/assets/mcps.go` |
| **5. Gate** | [`TEMPLATE_PIECE_GATE.json`](TEMPLATE_PIECE_GATE.json) | Receitas cuja `expression` é DSL de `internal/gate`, para colar em `output.gate` ou `slots.gates` | `gate.Compile` (nenhum loader de arquivo: é material de autoria) |
| **6. Loop** | [`TEMPLATE_PIECE_LOOP.json`](TEMPLATE_PIECE_LOOP.json) | DAG v2 / `loop.Loop` (`nodes` com `depends_on`, `routing`, `output.{artifact,gate}`) | `dag.Parse`, `loop.Store`, `legoloop validate` e `legoloop run` |

### Checklist de qualidade por tipo
* **Skill:** `name` igual ao nome do diretório; `description` ≤ 300 caracteres começando pelo gatilho de uso e dizendo quando NÃO usar; Target Files fechado; regras negativas; comandos de verificação com exit code; critério de parada. O campo `commands:` da V1 não era lido por nenhum componente e foi removido.
* **Script:** códigos de saída `0` sucesso, `1` falha de execução, `2` uso ou entrada inválida (convenção da CLI `legoloop`); alvo confinado à raiz; escrita atômica via temporário; `-DryRun` sem escrita; idempotente; transformação real (nada de placeholder que retorna sucesso sem fazer nada). Salvar com BOM.
* **Hook** (conforme a documentação de hooks do Claude Code; confirme na versão em uso): o script lê o evento em JSON no stdin; exit `0` permite, exit `2` bloqueia e devolve o stderr ao modelo; outros códigos não bloqueiam. `matcher` é regex de **nome de ferramenta**, não de arquivo; mantenha-o estreito. `timeout` em segundos. Nunca grave segredo literal no `command` (o inventário o redige e `assets.ContainsSecret` o acusa). Ajuste a sintaxe de `$CLAUDE_PROJECT_DIR` ao shell do seu harness.
* **MCP:** segredos só por referência `${VAR}`, nunca literais; o inventário nunca copia `env` nem `args`. No Windows, servidores via `npx` exigem `"command": "cmd"` com `"args": ["/c", "npx", ...]`. O protocolo é negociado em runtime: não fixe versão no config. O contrato de cada ferramenta vive no código do servidor e deve ser estrito, por exemplo:

```json
{
  "name": "analyze_target",
  "description": "Analisa um alvo e devolve o resultado estruturado, sem efeitos colaterais.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "target": { "type": "string", "description": "Caminho relativo do recurso alvo" },
      "mode": { "type": "string", "enum": ["analysis", "transformation", "validation"] }
    },
    "required": ["target", "mode"],
    "additionalProperties": false
  }
}
```
* **Gate:** só as funções que existem (`file_exists`, `word_count`, `line_count`, `has_pattern`, `has_binary_criteria`, `has_acceptance_criteria`, `schema_valid`, `interactive`, `prompt`) e o identificador `exit_code`. Comandos de terminal **não** entram na expressão: ficam em `verify` (argv, sem shell) e a expressão avalia o `exit_code`. Gate humano usa `prompt('pergunta')`. Caminhos literais, relativos ao projeto.
* **Loop:** passa em `legoloop validate`; `trigger` único; cada nó com `output.artifact` único e `output.gate` compilável; skills por `ref`, subloops por id de loop, nunca por caminho.

---

## 3. Régua Matricial de Maturação (Badges & Cores)

O motor `internal/reputation` dá a cada peça real uma nota de $0.00$ a $1.00$, com pesos fixos: Rigor $0.35$, Contratos $0.30$, Proveniência $0.20$, Histórico $0.15$.

* 🥇 **`ouro`** (`#F59E0B`): Score $\ge 0.85$.
* 🥈 **`prata`** (`#94A3B8`): Score $0.60 - 0.84$ (funcional, falta rigor de testes ou contratos formais).
* 🥉 **`bronze`** (`#F97316`): Score $0.35 - 0.59$ (peça anêmica, texto corrido sem verificação).
* ☣️ **`toxica`** (`#EF4444`): Score $< 0.35$ (prolixa, sem comandos, risco de alucinação).

**Como ler estes números (medido na V2).**
* **Teto sem histórico:** uma peça de autoria anônima, mesmo com Rigor e Contratos máximos, marca $0.82$ (prata). Chega a ouro ($0.90$) ao registrar execuções bem-sucedidas, ou com autoria reconhecida. Um template sozinho, portanto, não é "ouro": ele é o **caminho** até o ouro.
* **O scorer é textual.** Ele procura palavras-chave por regex no conteúdo. Hook e MCP, por serem fragmentos nativos sem texto, e os JSON em geral, ficam em bronze ou tóxica independentemente da qualidade real. O score de um template JSON não mede o template. Evoluir o scorer para avaliar estrutura é a recomendação da seção 8.

---

## 4. Sub-Loop de Elevação de Peças (Upgrade)

Qualquer peça em nível Bronze ou Prata pode ser elevada acionando a esteira canônica:
`Diagnóstico de Gaps` ➔ `Entrevista Socrática` ➔ `SPEC` ➔ `PLAN` ➔ `ROUTE` ➔ `BUILD TDD` ➔ `Peça elevada`.

O diagnóstico é `legoloop upgrade <arquivo> --diagnose`; a validação pós-reescrita é `--apply`.

---

## 5. Convenção de Placeholders

* Todo valor a preencher usa `{{NOME_EM_MAIUSCULAS}}` (ou `{{nome-em-kebab}}` quando vira identificador). Não há outra sintaxe de placeholder.
* Orientação ao autor fica em comentário HTML `<!-- ... -->` nos Markdown e é apagada ao concluir. Os JSON não têm comentários: a orientação está nesta seção e na seção 2.
* **Gate de fechamento de qualquer documento preenchido:** `!has_pattern('<arquivo>', '\{\{')`.
* Os JSON de hook e de loop não têm placeholders: são exemplos válidos como estão e devem ser editados. `TEMPLATE_PIECE_LOOP.json` usa ids concretos (`exemplo-loop-spec-plano`) porque `loop.ValidateID` não aceita `{{`.

---

## 6. Dicionário de Campos e Vocabulário

### Tipos de nó nos três vocabulários
| INTENT (bloco) | ROUTE (`type`) | DAG / Loop (`type`) |
|:---|:---|:---|
| `skill` | `unit_task` com skill no roteamento | `skill` (`ref: skills/...`) |
| `script` | `unit_task` com `verify` | `script` (`command` argv) |
| `subloop` | `subloop` (`ref` = id do loop) | `subloop` (`ref` = id do loop) |
| `unit_task`, `ui_task` | `unit_task` | `skill` ou `script` |
| `stop_gate` | `stop_gate` | nó cujo `output.gate` usa `prompt(...)` |

### Enumerações do ROUTE
* `nature`: `nano` | `micro` | `macro`
* `status`: `PENDING` | `RUNNING` | `PASSED` | `FAILED` | `SKIPPED` | `BLOCKED`
* `execution_mode`: `tdd_red_green` | `sequential` | `parallel`
* `isolation.workspace`: `local` | `git_worktree` | `tempdir`; `isolation.executor`: `orchestrix_hermetic` | `local_process`
* `isolation.subagents_allowed`: **sempre `false`** (AGENTS.md, regra 4). O campo existe para tornar a proibição explícita e auditável.
* `rollback_strategy` (por nó) e `rollback_on_failure` (global): `abort_pipeline` | `stash_changes` | `skip_optional`
* `routing.mode`: `suggested` | `manual` | `auto`; `routing.effort`: `low` | `medium` | `high` | `null`
* `defects[]`: `{ "id", "severity": "blocker|major|minor", "summary", "resolved" }`, preenchido durante a execução
* `output.verify[]`: lista de argv (cada comando é um array, **sem shell**); `output.gate`: DSL de gate avaliado com o `exit_code` da verificação
* `spec_ref`: id do `AC-*` que o nó fecha; `inputs.spec_anchor`: `arquivo#AC-*`
* `Escopo` (opcional): campo só do PLAN, com o detalhe do que o teste do ticket prova; não existe no ROUTE

### Convenção de exit code (scripts e CLI)
`0` sucesso · `1` falha de execução · `2` uso ou entrada inválida.

---

## 7. Verificação Mecânica dos Templates

```powershell
# Loop é um DAG válido (esperado: exit 0)
legoloop validate docs/templates/TEMPLATE_PIECE_LOOP.json

# Script sem erro de sintaxe (esperado: 0 erros)
$e=$null; [System.Management.Automation.Language.Parser]::ParseFile("$PWD\docs\templates\TEMPLATE_PIECE_SCRIPT.ps1",[ref]$null,[ref]$e) | Out-Null; $e.Count

# Todos os JSON são JSON válido
Get-ChildItem docs/templates/*.json | ForEach-Object { $null = Get-Content $_ -Raw | ConvertFrom-Json }
```

Verificações permanentes, executadas por `go test ./docs/templates/` (usam os pacotes reais do repositório):
* os 17 gates de ROUTE, PLAN, receitas e loop compilam em `gate.Compile`;
* `TEMPLATE_PIECE_LOOP.json` passa em `dag.Parse` (que rejeita campos desconhecidos) e é carregado por `loop.Store`;
* o scanner de `assets` descobre a skill, o hook (`PreToolUse:Write|Edit|MultiEdit`), o MCP (transporte `stdio`), o script e o loop;
* todo nó do ROUTE tem o ticket correspondente no PLAN, com mesma natureza, timeout, teto de tokens, modelo e `spec_ref`;
* o script responde `0`, `0`, `0`, `2`, `2`, `2`, `1` nos sete cenários (dry-run, execução, idempotência, alvo inexistente, fora da raiz, raiz-irmã, arquivo somente leitura) e não deixa `.legoloop-tmp`.

---

## 8. Limitações Conhecidas e Decisões Pendentes

1. **O ROUTE v2 ainda não tem schema publicado**, então os motores e os testes validam o conjunto de chaves contra o próprio `TEMPLATE_ROUTE.json` (`templates.RouteKeys()`), não contra um JSON Schema.
2. **Não existe o schema publicado** `https://legoloop.dev/schemas/route.v2.json` (nem `gate-recipe.v1.json`); as URLs são identificadores, como já é `dag.v1.json`. Publicar os schemas permitiria impor os templates com o gate `schema_valid(...)`.
3. **Três dialetos de ROUTE convivem** em `docs/routes/`: o legado (`id`, `name`, `target_files`, `gate`, 6 arquivos), o gerado pelo `plangen` e o da matriz. Migrar os legados não faz parte desta entrega.
4. **Scorer textual** (seção 3): recomenda-se avaliar JSON e `.ps1` estruturalmente, e permitir que o histórico de execuções suba o score de peças que não são Markdown.
5. **Formato de hook e MCP é o do Claude Code.** Outros harnesses (Gemini, Cursor, OpenCode) têm formatos próprios; o `internal/assets/mcps.go` já mapeia vários para MCP, mas há template de hook para um só.
6. **Versão e histórico.** Estes templates (`template_version` 2.0.0) substituíram a V1 em `docs/templates/`. A V1 permanece recuperável pelo histórico do git (`git log -- docs/templates`).

---

## 9. Motores que Seguem Estes Templates

Os templates são embutidos no binário (docs/templates/templates.go, go:embed), então o arquivo desta pasta é a única fonte de verdade. Cada motor tem teste de conformidade que **falha** se ele divergir do template.

| Template | Motor | O que garante o teste |
|:---|:---|:---|
| TEMPLATE_INTENT.md | interview.MaterializeIntent | Mesmas seções na mesma ordem, sem placeholder, status UNDER_REVIEW (a aprovação é sempre humana; abaixo de 95% de confiança, DRAFT) |
| TEMPLATE_SPEC.md | specgen.Generate | O prompt embute o template; o cabeçalho é carimbado pelo motor (status, versão, link da intenção, alvo, template); placeholder sem preencher é rejeitado |
| TEMPLATE_PLAN.md | plangen.Generate | Mesmas seções; espelha o ROUTE nó a nó; épicos com o nome que a SPEC lhes dá |
| TEMPLATE_ROUTE.json | plangen.Generate | Mesmo conjunto de chaves em todos os níveis; todo output.gate compila em gate.Compile; subagentes sempre alse; dependências existentes e sem ciclo; SG-FINAL ao fim |
| TEMPLATE_PIECE_* | legoloop upgrade | A recomendação aponta o template do tipo da peça |

Regras de geração do plangen:
* IDs canônicos TICK-NN.M; os aliases que o modelo escreva (TICKET-E01.1, TICK-E01.1) são resolvidos para o mesmo ticket.
* Nível de supervisão: 
ano gate após cada ticket; micro após cada ticket micro; epic entre épicos; utonomous nenhum. Entre épicos há gate nos três primeiros níveis, salvo se o último ticket do épico já tem o seu. O último épico fecha no SG-FINAL.
* Quem dependia de um ticket com gate passa a depender do gate, e cada épico espera a fronteira do anterior.
* O comando de validação vira argv (separado por &&). Se depender de shell (pipe, redirecionamento, aspas, glob) ou de executável relativo, o nó **não** inventa um comando: cai em gate humano prompt(...).
* Critério da SPEC sem ticket e dependência que não existe viram warnings no resultado.