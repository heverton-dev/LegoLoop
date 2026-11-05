# MANUAL OPERACIONAL: TEMPLATES CANÔNICOS DO LEGOLOOP

> **Data:** 2026-11-05
> **Autoridade:** Ecossistema LegoLoop
> **Status:** Canônico
> **Escopo:** Como, onde, o que faz e por quê cada template de `docs/templates/` é utilizado no LegoLoop
> **Fontes de verdade:** `docs/templates/templates.go` · `docs/templates/templates_test.go` · `docs/templates/README.md` · `internal/assets/skills.go` · `internal/interview/intent_render.go` · `internal/specgen/specgen.go` · `internal/plangen/plan_md.go` · `internal/plangen/route.go`

---

## 1. Arquitetura Geral

`docs/templates/` é a **Single Source of Truth (SSOT)** do LegoLoop. Os arquivos `.md` e `.json` são embutidos no binário Go via `go:embed` em `docs/templates/templates.go`. Os motores (`interview`, `specgen`, `plangen`) importam `legoloop/docs/templates` e usam os templates como molde estrutural. Os testes em `templates_test.go` garantem que nenhum motor diverge do template.

O diretório contém **13 arquivos** organizados em duas camadas:

- **Cascata de Governança:** 4 templates sequenciais (INTENT → SPEC → PLAN → ROUTE)
- **Sextilha de Peças:** 6 templates atômicos (Skill, Script, Hook, MCP, Gate, Loop)
- **Infraestrutura:** `templates.go`, `templates_test.go`, `README.md`

---

## 2. Cascata de Governança (4 Estágios Sequenciais)

```text
┌──────────────┬──────────────┬──────────────┬────────────────────┐
│ 1. INTENT    │ 2. SPEC      │ 3. PLAN      │ 4. ROUTE           │
│ "O Porquê"   │ "O Quê"      │ "O Como"     │ "A Máquina"        │
└──────────────┴──────────────┴──────────────┴────────────────────┘
```

### 2.1. TEMPLATE_INTENT.md

- **Formato:** Markdown com placeholders `{{}}`
- **Motor Go:** `internal/interview/intent_render.go` → `RenderIntent()`
- **O que faz:** Materializa o documento de intenção a partir do espelho socrático da entrevista. Mesmas seções, mesma ordem do template. O que a entrevista não respondeu vira "Não informado", nunca placeholder.
- **Por quê:** Responde **"por quê"**. Sem intent aprovada (confiança ≥ 95%), nenhum artefato downstream nasce.
- **Gate de fechamento:** `!has_pattern('...', '\\{\\{')`
- **Saída em disco:** `docs/intents/YYYY-MM-DD_INTENT_<NOME>.md`

### 2.2. TEMPLATE_SPEC.md

- **Formato:** Markdown com placeholders `{{}}`
- **Motor Go:** `internal/specgen/specgen.go` → `Generator.Generate()`
- **O que faz:** Gera a especificação técnica formal. O prompt ao LLM embute o template inteiro. Cabeçalho carimbado pelo motor (status, versão, link da intent, alvo). Placeholder não preenchido é rejeitado.
- **Por quê:** Responde **"o quê"**. Combina Premissas Addy Osmani + Costuras Matt Pocock + Critérios Binários BDD `AC-XX`. Todo épico tem critério de caminho de falha; há critérios não funcionais `AC-NF`.
- **Gate de fechamento:** `has_binary_criteria() && has_acceptance_criteria()`
- **Saída em disco:** `docs/specs/YYYY-MM-DD_SPEC_<NOME>.md`

### 2.3. TEMPLATE_PLAN.md

- **Formato:** Markdown com placeholders `{{}}`
- **Motor Go:** `internal/plangen/plan_md.go` → `renderPlan()`
- **O que faz:** Materializa o plano de execução: fatiamento vertical, tickets `TICK-NN.M`, nível de supervisão, gates entre épicos. Espelha nó a nó o ROUTE.
- **Por quê:** Responde **"como"**. É o documento legível por humanos que o usuário aprova no Stop-Gate antes do BUILD.
- **Saída em disco:** `docs/plans/YYYY-MM-DD_PLAN_<NOME>.md`

### 2.4. TEMPLATE_ROUTE.json

- **Formato:** JSON com placeholders `{{}}`
- **Motor Go:** `internal/plangen/route.go` → `renderRoute()` + `templates.RouteKeys()`
- **O que faz:** Gera o grafo executável de nós (`nodes`) com isolamento, roteamento, timeout, token budget, gates DSL, rollback. `RouteKeys()` extrai o schema de chaves para validar ROUTEs gerados.
- **Por quê:** Responde **"a máquina"**. É o espelho executável do PLAN — `legoloop run` consome este JSON.
- **Regra:** PLAN e ROUTE divergem = defeito. Os testes verificam isomorfismo campo a campo.
- **Saída em disco:** `docs/routes/YYYY-MM-DD_ROUTE_<NOME>.json`

---

## 3. Sextilha de Peças (6 Tipos Atômicos)

### 3.1. TEMPLATE_PIECE_SKILL.md

- **Formato:** Markdown com frontmatter YAML
- **Quem lê:** `internal/assets/skills.go` → `parseSkill()` → `parseFrontmatter()`
- **Campos parseados do frontmatter:** `name`, `description`, `consumes`, `produces`, `artifact`, `gate`, `suggested_model` — propagados para `Item.Meta` (`map[string]string` em `internal/assets/types.go`)
- **O que faz:** Define uma skill como nó tipado no grafo semântico. O parser extrai os campos do frontmatter. O corpo (seções 1–4) é contexto para o LLM executor.
- **Por quê:** Peça de **transformação cognitiva** — o LLM segue o procedimento, o motor Go valida o gate. Tipagem forte (`Artifact:*`) permite composição de skills em loops.
- **Onde se instala:** `~/.claude/skills/<nome>/SKILL.md` (ou `.aidd/`, `.agents/`, etc.)

### 3.2. TEMPLATE_PIECE_SCRIPT.ps1

- **Formato:** PowerShell 5.1, UTF-8 **com BOM**
- **Quem lê:** `internal/assets/scripts.go` (raízes `scripts/`, `tools/`, `.aidd/scripts`, `.claude/scripts`)
- **O que faz:** Define um script determinístico: transformação atômica, idempotente, com confinamento à raiz do projeto, escrita via arquivo temporário, exit codes `0` (sucesso) / `1` (falha) / `2` (uso inválido). O exemplo normaliza texto (remove trailing whitespace, converte para LF).
- **Por quê:** Peça de **transformação mecânica** — sem LLM. Executado diretamente por `legoloop` ou como `verify[]` de um nó ROUTE.
- **BOM obrigatório:** PowerShell 5.1 sem BOM lê UTF-8 como ANSI.

### 3.3. TEMPLATE_PIECE_HOOK.json

- **Formato:** Fragmento de `settings.json` (formato Claude Code)
- **Quem lê:** `internal/assets/hooks.go`
- **O que faz:** Define um hook de harness: `matcher` (regex de nome de ferramenta), `command`, `timeout`. O exemplo intercepta `Write|Edit|MultiEdit` com `guard-secrets.ps1`.
- **Por quê:** Peça de **interceptação** — trava na borda do harness. Exit `0` permite, exit `2` bloqueia e devolve stderr ao modelo. O inventário redige o `command` (nunca é exposto publicamente).

### 3.4. TEMPLATE_PIECE_MCP.json

- **Formato:** Fragmento de `.mcp.json`
- **Quem lê:** `internal/assets/mcps.go`
- **O que faz:** Define um servidor MCP: `command`, `args`, `env` (segredos só por referência `${VAR}`, nunca literais).
- **Por quê:** Peça de **capacidade externa** — expande ferramentas do agente. O inventário nunca copia `env` nem `args` (proteção contra vazamento de segredos).

### 3.5. TEMPLATE_PIECE_GATE.json

- **Formato:** JSON com receitas de gate
- **Quem lê:** `gate.Compile()` (material de autoria, sem loader de arquivo dedicado)
- **O que faz:** Catálogo de 5 receitas reutilizáveis:
  - `mechanical` — exit_code == 0
  - `artifact-quality` — file_exists + word_count > 50 + !TODO/FIXME/TBD
  - `spec-contract` — has_binary_criteria + has_acceptance_criteria
  - `schema-contract` — schema_valid
  - `human-stop-gate` — prompt('pergunta')
- **Por quê:** Peça de **validação** — as `expression`s compilam no DSL de `internal/gate`. Cola-se em `output.gate` de um nó ROUTE ou Loop.

### 3.6. TEMPLATE_PIECE_LOOP.json

- **Formato:** DAG v2 / `loop.Loop`
- **Quem lê:** `dag.Parse()` + `loop.Store` + `legoloop validate` + `legoloop run`
- **O que faz:** Define um loop modular: nós com `depends_on`, `routing`, `output.artifact`, `output.gate`. O exemplo encadeia `aidd-spec` → `aidd-tickets` com stop-gate humano antes do BUILD.
- **Por quê:** Peça de **composição** — orquestra skills e scripts num grafo acíclico dirigido.
- **Sem placeholders:** IDs concretos (`exemplo-loop-spec-plano`) porque `loop.ValidateID` rejeita `{{`.

---

## 4. Infraestrutura de Suporte

### 4.1. templates.go

Pacote Go que embute os 4 templates de governança via `go:embed`. Funções exportadas:


| Função              | Retorno                                                                                              |
| :--- | :--- |
| `Intent()`          | Conteúdo de `TEMPLATE_INTENT.md`                                                                     |
| `Spec()`            | Conteúdo de `TEMPLATE_SPEC.md`                                                                       |
| `Plan()`            | Conteúdo de `TEMPLATE_PLAN.md`                                                                       |
| `Route()`           | Conteúdo de `TEMPLATE_ROUTE.json` (clone)                                                            |
| `Headings(md)`      | Lista ordenada de títulos `## N. ...` de um documento                                                |
| `StripGuidance(md)` | Remove comentários HTML `<!-- ... -->` de orientação ao autor                                        |
| `RouteKeys()`       | Extrai o `RouteShape` (chaves de todos os níveis) do template ROUTE para validação de ROUTEs gerados |
| `Version`           | Constante `"2.0.0"`                                                                                  |


Importado por 10 pacotes: `interview`, `specgen`, `plangen`, `fabricar`, `processo`, `studio`, e os respectivos testes.

### 4.2. templates_test.go

17+ verificações de conformidade executadas por `go test ./docs/templates/`:

1. Todos os JSON são JSON válido
2. Conteúdo embutido (`go:embed`) é idêntico ao arquivo em disco
3. Seções numeradas são estáveis (INTENT: 6, SPEC: 6, PLAN: 4)
4. `StripGuidance` remove comentários HTML
5. Todos os gates (ROUTE + PLAN + Gate recipes + Loop) compilam em `gate.Compile`
6. `TEMPLATE_PIECE_LOOP.json` passa em `dag.Parse` e é carregado por `loop.Store`
7. O inventário (`assets.Scanner`) reconhece Skill, Hook, MCP, Script e Loop
8. Hooks não contêm segredos literais (`assets.ContainsSecret`)
9. PLAN espelha ROUTE nó a nó (ID, natureza, timeout, token budget, modelo, spec_ref)
10. Fragmentos nativos (Hook, Loop) não contêm placeholders `{{`
11. Script tem BOM UTF-8
12. `RouteKeys()` retorna conjuntos não vazios e ordenados com chaves obrigatórias

### 4.3. README.md

Documentação humana da cascata, da sextilha, do scorer de reputação (`internal/reputation`), das convenções de placeholders e dos comandos de verificação mecânica.

---

## 5. Diagrama de Dependências

```text
TEMPLATE_INTENT.md ──► interview.RenderIntent()
         │
         ▼
TEMPLATE_SPEC.md ───► specgen.Generate()
         │
         ▼
TEMPLATE_PLAN.md ──┐
                   ├► plangen.renderPlan() + plangen.renderRoute()
TEMPLATE_ROUTE.json┘

TEMPLATE_PIECE_SKILL.md ──► assets.parseSkill() ──► Item.Meta ──► graph.Index
TEMPLATE_PIECE_SCRIPT.ps1 ─► assets.scanScripts()
TEMPLATE_PIECE_HOOK.json ──► assets.scanHooks()
TEMPLATE_PIECE_MCP.json ───► assets.scanMCPs()
TEMPLATE_PIECE_GATE.json ──► gate.Compile() (autoria)
TEMPLATE_PIECE_LOOP.json ──► dag.Parse() + loop.Store ──► legoloop run
```

Cada seta é coberta por teste em `templates_test.go`. Template que nenhum componente lê não é template — é decoração (regra do README, seção 2).

---

## 6. Regras Transversais

1. **Placeholders:** todo valor a preencher usa `{{NOME}}`. Gate de fechamento: `!has_pattern('<arquivo>', '\\{\\{')`. Fragmentos nativos (Hook, Loop) são exemplos válidos como estão, sem placeholders.
2. **Espelho PLAN ↔ ROUTE:** mesmo conjunto de IDs, dependências, naturezas, limites e roteamento. Divergência é defeito. Verificado por `TestPlanMirrorsRoute`.
3. **Embute no binário:** o arquivo em `docs/templates/` é a única fonte de verdade. `templates.go` usa `go:embed`; os motores importam o pacote, nunca copiam strings.
4. **Scorer de reputação:** `internal/reputation` dá nota 0.00–1.00 a cada peça real (Rigor 0.35, Contratos 0.30, Proveniência 0.20, Histórico 0.15). Ouro ≥ 0.85, prata 0.60–0.84, bronze 0.35–0.59, tóxica < 0.35.
5. **Subagentes proibidos:** `isolation.subagents_allowed` é sempre `false` nos templates ROUTE e Loop (AGENTS.md, regra 4).
6. **Exit codes:** `0` sucesso, `1` falha de execução, `2` uso/entrada inválida. Convenção da CLI `legoloop` e de todos os scripts.

&nbsp;

&nbsp;

```txt
Dear Contabo Trust & Compliance Team,

Regarding our suspended server:
- IP Address: 167.86.69.79
- Instance: vmi3278564
- Reference: Ticket#2026110110005763

Incident Analysis & Root Cause:
An unauthorized Alpine container exploited an exposed Docker Daemon API configuration (port 2375) to perform outbound network scanning.

Immediate Remediations & Prevention Plan:
- Docker Hardening: We are removing all external TCP bindings (port 2375) from the Docker daemon configuration, restricting dockerd exclusively to local Unix domain sockets (unix:///var/run/docker.sock).
- Outbound Filtering: We are applying strict iptables drop rules (OUTPUT and DOCKER-USER chains) for destination port 2375 to block all outbound scanning traffic permanently.
- Malware Cleanup: We are terminating and removing all rogue Alpine containers, verifying cron tasks/processes, and rotating all administrative credentials.

Currently, all instance controls (Power, Reboot, VNC Console) are disabled/locked in our Customer Panel and network traffic is suspended.
Please lift the administrative suspension and enable our VNC/network access so we can execute these remediations immediately.

Best regards,
Renan Borges da Silva
Infrastructure & IT Team - Conexão
infra.ti@conexao.com.br
```
