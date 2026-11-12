# 🗺️ PLAN: Refatoração Completa da Lei #15 (ASD-STE100) e Orchestrix sobre Claude Code

> **Status:** Em execução. E00 a E03 concluídos. `SG-E03` aprovado em 2026-11-12. E04 a E06 pendentes. Ver o REPORT de execução.  
> **Data:** 2026-11-10  
> **Versão:** 1.1.0  
> **Autor:** Claude (a pedido do usuário)  
> **Referência da Auditoria:** `docs/2026-11-10_AUDITORIA_TAXONOMIA_ASD_STE100.md`  
> **Referência de Governança:** `AGENTS.md` itens 4, 6, 11 e 12  
> **Intenção:** [docs/intents/2026-11-10_INTENT_NORMALIZACAO_TAXONOMICA_LEI15.md](../intents/2026-11-10_INTENT_NORMALIZACAO_TAXONOMICA_LEI15.md)  
> **Especificação:** [docs/specs/2026-11-10_SPEC_NORMALIZACAO_TAXONOMICA_LEI15.md](../specs/2026-11-10_SPEC_NORMALIZACAO_TAXONOMICA_LEI15.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-10_ROUTE_NORMALIZACAO_TAXONOMICA_LEI15.json](../routes/2026-11-10_ROUTE_NORMALIZACAO_TAXONOMICA_LEI15.json)  
> **Registro de Execução:** [docs/reports/2026-11-12_REPORT_EXECUCAO_LEI15.md](../reports/2026-11-12_REPORT_EXECUCAO_LEI15.md)  
> **Pendências para o BUILD:** nenhuma de aprovação. O DAG espelho em `docs/dags/` nasce no `TICK-01.6`, depois do `TICK-01.2` (o parser precisa aceitar ponto no id).  

---

## 1. Decisões travadas (usuário, 2026-11-10)

| # | Decisão |
| :--- | :--- |
| D1 | Refatoração **completa e profunda**. Os contratos em disco entram na migração. |
| D2 | O nome da marca continua **LegoLoop**. O dicionário ganha uma exceção de marca. |
| D3 | Os loops existentes **não podem quebrar**: 15 em `loops/` e 6 no Cofre. O leitor dos nomes antigos permanece. |
| D4 | Ordem de roteamento (soft, o ticket sobrescreve com justificativa): 9Router, Antigravity, Mimocode, OpenCode, Claude Code, Ollama (ids no código: `9router`, `agy`, `mimo`, `opencode`, `claude`, `ollama`). O Claude Code é conta secundária por enquanto. |
| D5 | Modelo por natureza: `nano` Haiku. `micro` e `macro` Sonnet. Opus só com a marca `extreme` e aprovação no Stop-Gate. |
| D6 | Subagente efêmero permitido (regra da seção 3). Fork e workflow livre continuam proibidos. |
| D7 | `agy` em **raias fixas**: uma raia por Orca, conta da própria sessão, padrão 2 raias. |
| D8 | O 9Router já perdeu as contas `antigravity` e `kiro`. Outro agente cuida da VPS. Este PLAN não toca nela. |

## 2. Linha de base

`go run ./tools/taxlinter internal cmd loops skills docs/templates` retornou exit 1 com **412 violações** em 2026-11-10. O número é um **piso**. Experimento de 2026-11-10: o linter lê só comentários e literais Go sem `/`, `\`, `-` ou `.`, chaves JSON simples e títulos de Markdown. Ele **não** lê identificadores Go (`type Loop`, `LoopID`), chaves compostas (`loop_end`, `subloop`), valores JSON, corpo de Markdown nem `.typ`. Por isso o gate "taxlinter exit 0" passaria sem renomear identificador. O `TICK-00.4` corrige o linter e mede o valor real.

| Termo proibido | Qtd | Canônico |
| :--- | ---: | :--- |
| loop | 198 | Processo |
| bloco | 57 | Peça |
| pipeline | 49 | Processo |
| card | 21 | Ticket |
| fase | 18 | Etapa |
| esteira | 17 | Processo |
| fluxo | 15 | Processo |
| demais (macro, issue, componente, capacidade, tarefa, módulo, task, storage, grade) | 37 | ver dicionário |

Superfície de contrato (medida): `pipeline.dag.json` em 39 literais Go e 195 arquivos de docs e skills. `loops/` em 90 literais Go. `internal/loop` com 22 imports. O `fabricar` e o `compor` gravam `*.loop.json`. O `fechar` exige `pipeline.dag.json` (`internal/fechar/fechar.go:201`). O schema tem `loop_end`, `subloop` e o campo `loop` (29 referências Go).

## 3. Estratégia de execução

* **Engine:** `orchestrix_hermetic`. Um ticket por contexto limpo.
* **Subagente efêmero (D6).** Permitido quando cumpre as seis condições:
  1. O contexto contém só os artefatos contratados do ticket.
  2. O ticket lista as ferramentas permitidas.
  3. Um subagente executa um ticket e termina ao concluir.
  4. A saída é um artefato em disco mais um Gate mecânico.
  5. O subagente não cria outros subagentes.
  6. O Write Boundary é `strict`.
  Mecanismo preferido: subprocesso estéril (`claude -p` com lista de ferramentas).
* **Roteamento (D4).** O Orchestrix percorre a ordem da D4. Em 429 ou falha, o ticket passa ao **próximo harness**. Nunca troca de conta.
* **Raias do `agy` (D7).** Cada Orca com `agy` é uma raia. O ticket fica fixo na raia até terminar. O número de raias é configurável (padrão 2; cai para 1 se os Orcas compartilharem login). O alarme "Terms of Service" desliga todas as raias na sessão. O Orchestrix confere antes do primeiro uso do dia que a conta logada não consta como bloqueada.
* **Modelo e esforço (D5).** `nano` Haiku `low`. `micro` Sonnet `medium`. `macro` Sonnet `high`. `extreme` Opus `high`, com Stop-Gate.
* **Paralelismo e assíncrono.** Tickets sem dependência mútua e sem arquivo em comum rodam em paralelo, com limite de workers por harness. Ticket longo roda em segundo plano e o Orchestrix consulta o Gate.
* **TDD estrito.** Red, Green, Refactor. Teste novo só vale se uma mutação o faz falhar (`AGENTS.md` item 8).
* **Teste dourado (D3).** O `TICK-00.2` congela o resultado de `legoloop validate` nos loops existentes. Todo ticket de E01 a E05 repete o teste.
* **Gate padrão de cada ticket:** uma peça única, `go run ./tools/ticketgate` (criada no `TICK-00.5`), com as opções `-pkgs`, `-tax`, `-gold`, `-build` e `-livros`. Ela roda `go test -count=1`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .` (falha se a saída não for vazia), `git diff --exit-code -- go.mod go.sum`, o `taxlinter` nas pastas pedidas e o teste dourado. O DAG aceita um único `command` por nó (`internal/dag/types.go:118`) e `gofmt -l` não falha por exit code. Por isso a bateria vira um único argv.
* **Parâmetros por natureza (espelhados no ROUTE):** `nano` 180 s, entrada 8.000 e saída 300 tokens, Haiku `low`, combo 9Router `code-fast`. `micro` 180 s, 16.000 e 600 tokens, Sonnet `medium`, combo `code-fast`. `macro` 300 s, 32.000 e 1.200 tokens, Sonnet `high`, combo `code-pro`. Retry 1, backoff 1.000 ms, modelo de reserva no retry. `stop_gate`: processo local, sem tokens.
* **Rollback:** `stash_changes` por ticket. Commit e push só quando o usuário pedir.
* **Dados reais:** nenhum ticket toca o Cofre real `~/.aidd` sem Stop-Gate. Testes usam `t.TempDir()`.

## 4. Mapa de épicos

| Épico | Nome | Tickets | Natureza dominante |
| :--- | :--- | ---: | :--- |
| **E00** | Fundação: medição, teste dourado, governança e peça `ticketgate` | 5 | `micro` |
| **E01** | Compatibilidade: leitor dual, schema v2, migrador e projeção ROUTE para DAG | 6 | `macro` |
| **E02** | Geradores, templates, prompts, CLI e skills | 5 | `micro` |
| **E03** | Código interno por pacote | 7 | `macro` |
| **E04** | Dados e diretórios: `loops/`, Cofre e documentação | 4 | `macro` |
| **E05** | Orchestrix sobre Claude Code: roteamento, raias, paralelismo | 6 | `macro` |
| **E06** | Fechamento | 1 | `micro` |

Ordem: E00, E01, E02, E03, E04, E05, E06. E02 vem antes de E03 para os Processos novos nascerem corretos e o número de violações parar de crescer.

---

## 5. Tickets

### Épico 00 — Fundação

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-00.1` | Triar as 412 violações em `texto` ou `contrato`. A marca "LegoLoop" não é violação no linter atual (verificado em 2026-11-10). | `micro` Sonnet | `[]` | `file_exists('scratch/tax_triagem.json')` e soma das classes igual ao total |
| `TICK-00.2` | **Teste dourado:** congelar `legoloop validate` nos 15 loops de `loops/` e em cópias dos 6 do Cofre (cópia feita só com autorização explícita). Falha se qualquer loop deixar de validar. | `micro` Sonnet | `[]` | `go test ./internal/golden/...` verde |
| `TICK-00.3` | Atualizar `AGENTS.md` (regra 4: subagente efêmero; item 12: exceção de marca LegoLoop; ordem de roteamento e raias) e `TEMPLATE_PLAN.md` / `TEMPLATE_ROUTE.json` (`Subagentes: efêmero` com lista de ferramentas, campo `extreme`). Trocar o combo `code-nano`, que não existe no 9Router vivo, por `code-fast` em `plangen/route.go:221` e nos templates. | `micro` Sonnet | `[]` | `go test ./docs/templates/...` verde |
| `TICK-00.4` | Ampliar o `taxlinter`: identificadores Go (camelCase e snake_case), chaves JSON compostas, valores JSON, corpo de Markdown, `.typ` e literais com `/`, `-` ou `.`. Exceção para a marca "LegoLoop". Gerar o baseline real. | `macro` Sonnet | `TICK-00.1` | `go test ./tools/taxlinter/...` verde e `scratch/tax_baseline.txt` existe |
| `TICK-00.5` | Criar a peça `tools/ticketgate`: a bateria única de verificação (seção 3). `gofmt` sujo faz o comando falhar. | `macro` Sonnet | `[]` | `go test ./tools/ticketgate/...` verde e `file_exists('tools/ticketgate/main.go')` |

#### 🛑 STOP-GATE HUMANO — `SG-E00`
* **Gate:** `prompt('Triagem e baseline aprovados? Nomes canônicos novos aprovados (seção 6)?')`
- [x] Usuário aprova os nomes novos de arquivo, pasta, nó e campo (2026-11-10).
- [x] Usuário aprova o baseline real (2026-11-10).

---

### Épico 01 — Compatibilidade

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-01.1` | Centralizar os 90 literais `"loops"` e os 39 de `pipeline.dag.json` em uma constante e uma função de resolução únicas. Sem mudar comportamento. | `macro` Sonnet | `SG-E00` | Teste dourado verde e `grep` dos literais fora do pacote central retorna vazio |
| `TICK-01.2` | Leitor dual: o parser de DAG aceita o nome antigo e o novo, os tipos `subloop` e `loop_end` e os novos equivalentes, e o campo `loop` e o novo. Aceitar ponto no id de nó (`TICK-01.1`): o parser recusa hoje (`parse.go:24`) e `legoloop validate` falha no DAG existente `docs/dags/2026-11-09_DAG_QUINTETO_GOVERNANCA_ISO.json`. O `ref` de subprocesso mantém a regex estrita. | `macro` Sonnet | `TICK-01.1` | Teste de tabela com os dois dialetos, id com ponto aceito e teste dourado verde |
| `TICK-01.3` | `schemas/dag.v2.json` aceita os dois dialetos. O arquivo antigo `dag.v1.json` permanece. | `micro` Sonnet | `TICK-01.2` | Todos os loops de `loops/` validam contra v1 e v2 |
| `TICK-01.4` | `legoloop migrar`: converte um Processo do dialeto antigo para o novo. Idempotente, com backup e rollback. Usa só diretórios de teste. | `macro` Sonnet | `TICK-01.2` | Teste de ida e volta: migrar e validar equivale ao original |
| `TICK-01.5` | A paridade SHA-256 do `legoloop fechar` aceita os dois nomes de arquivo. | `micro` Sonnet | `TICK-01.1` | `go test ./internal/fechar/...` verde com os dois nomes |
| `TICK-01.6` | Criar `tools/route2dag`: projeta a ROUTE em DAG (`command` = primeiro argv de `verify`; `output.artifact` = `scratch/gates/<id>.json`, único por nó; Stop-Gate com roteamento manual). Gerar o DAG espelho desta ROUTE em `docs/dags/` e validar todos os DAGs de `docs/dags/`. | `macro` Sonnet | `TICK-01.2` | `go test ./tools/route2dag/...` verde e `file_exists('docs/dags/2026-11-10_DAG_NORMALIZACAO_TAXONOMICA_LEI15.json')` |

---

### Épico 02 — Geradores, templates, prompts, CLI e skills

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-02.1` | `fabricar` e `compor` gravam os nomes canônicos novos (`camada4.go:115`, `engine.go:160`, `compor.go:475`) [CONCLUÍDO]. | `micro` Sonnet | `TICK-01.2` | Artefato gerado em `t.TempDir()` passa no `taxlinter` e em `legoloop validate` |
| `TICK-02.2` | Normalizar `docs/templates/*` e o `templates_test.go`. | `micro` Sonnet | `SG-E00` | `taxlinter docs/templates` exit 0 e `go test ./docs/templates/...` verde |
| `TICK-02.3` | Normalizar os prompts do `fabricar` e da entrevista (10 arquivos). Teste com LLM falso: a saída gerada não contém termo proibido. | `macro` Sonnet | `TICK-02.1`, `TICK-02.2` | Teste do LLM falso verde e mutação do prompt faz o teste falhar |
| `TICK-02.4` | Normalizar o texto de ajuda e as mensagens da CLI (`cli.go:53`, `commands.go:344`, `commands.go:375` e demais). | `nano` Haiku | `TICK-01.2` | `taxlinter internal/cli` exit 0 nos textos de ajuda |
| `TICK-02.5` | Normalizar as skills `skills/*/SKILL.md` e reinstalar com `legoloop instalar-skills`. | `micro` Sonnet | `TICK-02.4` | `taxlinter skills` exit 0 e instalação verde nos 5 harnesses |

#### 🛑 STOP-GATE HUMANO — `SG-E02`
* **Depende de:** todos os tickets de E01 e E02. O `TICK-01.4` e o `TICK-03.2` tocam `internal/dag`, então E01 termina antes de E03.
* **Gate:** `prompt('Novos Processos nascem sem termo proibido? Autoriza a refatoração interna?')`

---

### Épico 03 — Código interno por pacote

| Ticket | Pasta(s) | Violações (piso) | Nature / Modelo | Dep. |
| :--- | :--- | ---: | :--- | :--- |
| `TICK-03.1` | `internal/studio` (rotas `/api/loops` ganham alias de transição) | 58 | `macro` Sonnet | `TICK-03.3` |
| `TICK-03.2` | `internal/dag`, `internal/graph` | 57 | `macro` Sonnet | `TICK-03.3` |
| `TICK-03.3` | `internal/loop` vira `internal/studiodoc`. Roda sozinho primeiro: toca os 22 importadores. `internal/processo` recebe só ajuste de texto. | 67 | `macro` Sonnet | `SG-E02` |
| `TICK-03.4` | `internal/cli`, `internal/fabricar`, `internal/refinar` | 42 | `micro` Sonnet | `TICK-03.3` |
| `TICK-03.5` | `internal/assets`, `internal/board`, `internal/orchestrix`, `internal/suggest` | 54 | `micro` Sonnet | `TICK-03.3` |
| `TICK-03.6` | `internal/runner`, `internal/videorea` e demais pastas de `internal/` e `cmd/` | a medir | `micro` Sonnet | `TICK-03.3` |
| `TICK-03.7` | Verificação global: `taxlinter internal cmd` com exit 0 | n/a | `nano` Haiku | `TICK-03.1`, `TICK-03.2`, `TICK-03.3`, `TICK-03.4`, `TICK-03.5`, `TICK-03.6` |

* **Gate de cada ticket:** `taxlinter <pasta>` exit 0, gate padrão da seção 3 e teste dourado verde.
* O `TICK-03.3` depende de `SG-E02` e roda sozinho primeiro. Os demais dependem dele e rodam em paralelo (pastas disjuntas).

#### 🛑 STOP-GATE HUMANO — `SG-E03`
* **Depende de:** `TICK-03.7`.
* **Gate:** `prompt('Código interno normalizado e loops intactos? Autoriza mover dados e diretórios?')`

---

### Épico 04 — Dados e diretórios

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-04.1` | Renomear `loops/` para `processos/` (`git mv`) e migrar os 15 itens (9 pastas com `pipeline.dag.json` e 6 `.loop.json`) para os nomes novos com `legoloop migrar`. O teste dourado continua lendo cópias dos nomes antigos. | `macro` Sonnet | `SG-E03` | Teste dourado verde, os 15 validam sob os nomes novos e `file_exists('processos/legoloop-fechar/processo.dag.json')` |
| `TICK-04.2` | Migrar os 6 loops do Cofre real com `legoloop migrar`. **Ação destrutiva: exige backup e autorização explícita.** | `micro` Sonnet | `SG-E04-COFRE` | `legoloop validate` verde nos 6 e backup existe |
| `TICK-04.3` | Atualizar as referências em docs, manuais e livros (195 arquivos). Recompilar os livros pelo LegoPress. | `macro` Sonnet | `TICK-04.1` | `go run ./tools/ticketgate -livros docs/books,docs/videorea` exit 0 (executa `legoloop legopress auditar <pdf>` em cada um dos 27 PDFs) |
| `TICK-04.4` | Marcar o leitor dos nomes antigos como depreciado com aviso na CLI. **Não remover** (D3). | `nano` Haiku | `TICK-04.2` | Aviso aparece e os loops antigos continuam válidos |

#### 🛑 STOP-GATE HUMANO — `SG-E04-COFRE` (antes do `TICK-04.2`)
* **Depende de:** `TICK-04.1` e `TICK-01.4`.
* **Gate:** `prompt('Backup do Cofre confirmado? Autoriza migrar os 6 loops reais?')`

---

### Épico 05 — Orchestrix sobre Claude Code

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-05.1` | Executor de subagente efêmero: subprocesso `claude -p` com lista de ferramentas, contexto só dos Inputs do ticket, Write Boundary `strict`, timeout e saída em disco. | `macro` Sonnet | `SG-E03` | Teste com subprocesso falso: ferramenta fora da lista é recusada e escrita fora do limite falha |
| `TICK-05.2` | Roteador: ordem da D4, sobrescrita por ticket com justificativa, failover só para o próximo harness, nunca troca de conta. | `macro` Sonnet | `TICK-05.1` | Teste de tabela da ordem e teste de 429 passando ao próximo harness |
| `TICK-05.3` | Raias do `agy`: número configurável (padrão 2), ticket fixo na raia, conferência de conta bloqueada, desligamento global ao ver "Terms of Service". | `macro` Sonnet | `TICK-05.2` | Teste: 429 não move o ticket de raia e o alarme desliga todas as raias |
| `TICK-05.4` | Modelo e esforço por natureza (D5). Campo `extreme` no schema da ROUTE com Stop-Gate obrigatório. | `micro` Sonnet | `TICK-05.3` | Teste: `extreme` sem Stop-Gate falha na validação |
| `TICK-05.5` | Paralelismo e assíncrono: limite de workers por harness, tickets longos em segundo plano com consulta de Gate, detecção de conflito de arquivo. | `macro` Sonnet | `TICK-05.3` | Teste: dois tickets com o mesmo arquivo não rodam juntos |
| `TICK-05.6` | Telemetria de tokens, tempo e custo por ticket e por harness no REPORT. | `micro` Sonnet | `TICK-05.4`, `TICK-05.5` | `docs/reports/` recebe o balanço no formato do `TEMPLATE_REPORT.md` |

* **Smoke com sistema real:** depois do `TICK-05.5`, um Processo de duas etapas roda de ponta a ponta com o binário real e LLM real (`AGENTS.md` item 8). O resultado entra no REPORT.

---

### Épico 06 — Fechamento (`SG-FINAL`)

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `TICK-06.1` | Cadeia de fechamento do `AGENTS.md` item 11: manual, livro LegoPress, skill, loop, build e install, governança viva, reindexação, portões e REPORT. Incluir os cinco alvos no gate global do `taxlinter` e remover a ressalva "alvos já normalizados". | `micro` Sonnet | `TICK-04.3`, `TICK-04.4`, `TICK-05.6` | `taxlinter internal cmd processos skills docs/templates` exit 0 (`processos` é o nome de `loops/` após o `TICK-04.1`) e `file_exists('bin/legoloop.exe')` |

---

## 6. Nomes canônicos propostos (o `SG-E00` aprova ou troca)

| Atual | Proposto | Observação |
| :--- | :--- | :--- |
| `pipeline.dag.json` | `processo.dag.json` | O leitor aceita os dois. |
| `<nome>.loop.json` | `<nome>.studio.json` (aprovado em 2026-11-10) | O Cofre tem **três** arquivos por Processo: `.loop.json` (documento do Studio), `.dag.json` (contrato da CLI, gravado por `fabricar`, `compor`, `refinar`) e `.processo.json` (registro v2). Os nomes `.dag.json` e `.processo.json` **já estão em uso** e **permanecem**. Verificado em 2026-11-10: os 6 loops do Cofre real têm só `.loop.json` (v1). |
| `loops/` | `processos/` | `git mv` no `TICK-04.1`. |
| Nó `subloop` | `subprocesso` | Schema v2. |
| Gatilho `loop_end` | `processo_fim` | Schema v2. |
| Campo `loop` | `processo` | Schema v2. |
| Rota `/api/loops` | `/api/processos` | Alias de transição no Studio. |
| Pacote `internal/loop` (tipo `Loop`) | `internal/studiodoc`, tipo `Documento` (aprovado em 2026-11-10) | Não funde com `internal/processo`: o `loop` guarda documentos DAG (1.102 linhas, 22 importadores). O `processo` é o núcleo de domínio (3.052 linhas, 53 importadores). |

## 7. Riscos

| Risco | Mitigação |
| :--- | :--- |
| Renomear um contrato quebra Processos gravados | Leitor dual, teste dourado e Stop-Gate antes do Cofre real |
| O número 412 é piso | `TICK-00.4` mede o real antes de planejar E03 |
| Conflito de arquivo entre tickets paralelos | `TICK-05.5` detecta. Até lá, E03 roda em pastas disjuntas |
| Banimento do `agy` | Raias fixas, nenhuma troca de conta, alarme global (D7) |
| O Claude Code é conta secundária | Fica na posição 5 da ordem. Só entra por sobrescrita justificada ou por falha dos anteriores |

---

## 8. Evidências verificadas em 2026-11-10

| Fato | Como verifiquei |
| :--- | :--- |
| `loops/` tem 15 itens: 9 pastas com `pipeline.dag.json` e 6 `.loop.json`. Os 15 passam em `legoloop validate`. | `ls` e `legoloop validate` em cada item |
| O Cofre real tem 6 `.loop.json` (v1), sem `.dag.json` nem `.processo.json`. | `ls ~/.aidd/canonical/loops` |
| `.dag.json` e `.processo.json` já estão em uso. `.studio.json`, `processo.dag.json`, `processos/`, `internal/{golden,paths,migrar,studiodoc}` e `schemas/dag.v2.json` estão livres. | `grep` e `ls` |
| O linter não lê identificadores, chaves compostas, valores JSON, corpo de Markdown nem `.typ`. A marca "LegoLoop" não é violação. | Seis arquivos de teste em diretório temporário, saída do `taxlinter` |
| `file_exists` só aceita arquivo regular (`info.Mode().IsRegular()`). | `internal/gate/builtins.go:135` |
| O DAG tem um único `command` por nó e o `stop_gate` não exige `command`. | `internal/dag/types.go:118`, `parse.go:153` |
| O parser de DAG recusa ponto no id de nó (`^[A-Za-z0-9_-]{1,64}$`). O DAG `docs/dags/2026-11-09_DAG_QUINTETO_GOVERNANCA_ISO.json` **já existente falha** em `legoloop validate` por esse motivo. | `legoloop validate` e `internal/dag/parse.go:24` |
| O DAG recusa dois nós que gravam o mesmo `output.artifact` e exige `artifact` relativo e não vazio, inclusive no Stop-Gate. | `legoloop validate` (mensagens de erro reais) |
| A projeção desta ROUTE em DAG, com ids sem ponto, valida: 39 nós, sem ciclo. | Script em diretório temporário e `legoloop validate` |
| O `schemas/route.v2.json` recusava `null` em `spec_ref`, `test_target` e `epic_id` dos nós, e o `TEMPLATE_ROUTE.json` oficial usa `null` (7 erros no template, 23 na ROUTE). Corrigi o schema para aceitar `null`. Depois: 0 erros no template e na ROUTE. | `jsonschema` Draft 7 antes e depois, e `go test ./internal/gate/...` |
| `legoloop legopress auditar` recebe um PDF. Há 27 PDFs em `docs/books` e `docs/videorea`. | `legoloop legopress` e `git ls-files` |
| O comando `legoloop migrar` não existe. `instalar-skills` existe. | `legoloop help` |
| `claude` tem as opções `-p`, `--allowedTools`, `--tools`, `--model`, `--effort`, `--add-dir`, `--permission-mode` e `--no-session-persistence`. | `claude --help` |
| O 9Router vivo expõe os combos `code-fast`, `code-free`, `code-preview`, `code-pro` e `megacombo-*`. **Não** expõe `code-nano`. | `legoloop doctor --json` |
| Os ids de harness no código são `agy`, `mimo`, `opencode`, `claude`, `codex`, `gemini`, `cursor`, `omp` e `kiro`. O 9Router e o Ollama são provedores. | `internal/assets/system.go:19`, `internal/provider/inventory.go` |
| `fallback` com modelo `flash` é convenção do repositório (DAGs existentes e `brain/pick.go:60`). Não é id de modelo do `agy`. | `docs/dags/`, `internal/brain/pick.go` |

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano acima cobre 100% das decisões D1 a D8 e protege os loops existentes com Gates mecânicos?*  
*Aprovado este documento, gero a SPEC com critérios `AC-*` e o `docs/routes/2026-11-10_ROUTE_NORMALIZACAO_TAXONOMICA_LEI15.json`. O BUILD só começa depois disso.*


---

## 9. Registro de Execução (2026-11-12)

| Épico | Estado | Commit de fechamento |
| :--- | :--- | :--- |
| E00, E01, E02 | Implementados em 2026-11-10. `SG-E00` e `SG-E02` aprovados pelo usuário em 2026-11-12 (retroativo). | ver REPORT |
| E03 | Concluído. `SG-E03` aprovado em 2026-11-12. | `64b1858` |
| E04 | Não iniciado. `SG-E04-COFRE` bloqueia `TICK-04.2` e `TICK-04.4`. | n/a |
| E05 | Não iniciado. | n/a |
| E06 | Não iniciado. `SG-FINAL` é do usuário. | n/a |

**Divergências aceitas pelo usuário durante a execução:**

1. O E03 rodou antes da aprovação formal de `SG-E00` e `SG-E02`. O usuário aprovou os dois gates de forma retroativa em 2026-11-12.
2. Os nós do E03 foram divididos em fatias de até 12 arquivos (`ffadc27`). O DAG e a ROUTE já refletem a divisão.
3. O roteamento da seção 3 (D4, D5) não foi aplicado como escrito. O DAG declara `routing.mode: suggested`. O `unit_task` traduz `code-fast` e `code-pro` para Sonnet e `flash` para Haiku. A decisão D7 (raias do `agy`) está obsoleta: o `agy` foi banido em 2026-11-10.
4. O commit por nó foi autorizado pelo usuário. Nenhum push foi feito.
5. O motor ganhou mudanças não previstas (retentativa, limite de custo, `write_boundary` por diretório, higiene antes do gate, retenção de edições parciais). A lista com commits está no REPORT, seção 4.

**Tabela do Épico 03 (seção 5).** Os tickets `03.1`, `03.4`, `03.5`, `03.6` e `03.7` têm fatias `a`, `b`, `c` e até `h`. O livro-razão nó a nó está no REPORT, seção 3.
