# Nó e Peça no LegoLoop

## Nó é sinônimo de peça?

Não. O glossário canônico separa os dois termos (AGENTS.md:71-80).
Nó é o "vértice de execução unitária" do DAG. Peça é a "unidade reutilizável".
O nó é o lugar do grafo onde algo roda. A peça é o que o nó materializa.
O campo `peca` do nó guarda o nome do Processo que a peça materializa (internal/dag/types.go:165).
O comentário do código diz "Node é uma peça Lego do Processo" (types.go:135). Esse comentário mistura os dois termos, mas o glossário e a struct os mantêm separados.
"Etapa" é outro termo: no glossário ela é o passo interativo da CLI. Para o DAG, o termo correto é nó.

## Os seis tipos de nó

Todo nó tem os campos obrigatórios `id`, `type`, `routing` e `output` (com `artifact` e `gate`), conforme schemas/dag.v2.json. O molde do nó inteiro é docs/templates/TEMPLATE_DAG.json. Cada tipo acrescenta a sua regra de validação (internal/dag/parse.go:150-176) e o seu executor (internal/orchestrix/orchestrix.go:105-120).

**skill.** O objetivo é despachar o prompt de uma skill para um modelo, em execução isolada. Exige o campo `ref` com o nome da skill. O template da peça é docs/templates/TEMPLATE_PIECE_SKILL.md, um SKILL.md com frontmatter (name, description, consumes, produces, artifact, gate, requires).

**script.** O objetivo é executar um comando local determinístico, sem shell. Exige `command` como lista de argumentos (argv), validada contra injeção. O template da peça é docs/templates/TEMPLATE_PIECE_SCRIPT.sh e TEMPLATE_PIECE_SCRIPT.ps1, com uso, códigos de saída 0, 1 e 2, e os campos de suporte a SO e harness.

**subprocesso.** O objetivo é encaixar outro Processo inteiro como um único nó. Exige `ref` com o nome de um Processo válido. Roda isolado e estéril, em diretório temporário (internal/orchestrix/subloop.go:17-25). O template da peça é docs/templates/TEMPLATE_PIECE_PROCESSO.json, um DAG v2 completo. O campo antigo `loop_end` virou `processo_fim`, e o leitor converte.

**stop_gate.** O objetivo é parar o Processo e exigir decisão humana. Não exige `ref` nem `command`. O executor não faz nada, e o runner avalia o prompt do gate depois (orchestrix.go:116-118). O template é a receita `human-stop-gate` em docs/templates/TEMPLATE_PIECE_GATE.json, cuja expressão é `prompt('{{PERGUNTA_AO_HUMANO}}')`. Os gates binários automáticos (`exit_code == 0`, `file_exists`, `schema_valid`) são outras receitas do mesmo arquivo e vão em `output.gate` ou `slots.gates`.

**unit_no.** O objetivo é representar uma unidade cognitiva de engenharia (types.go:31-32). O `command` é opcional. Por segurança, só roda com consentimento explícito: sem a flag `--allow-unit-tasks`, o nó falha (internal/orchestrix/unit_run.go:73-75). Não existe template de peça próprio para este tipo: `grep -rln "unit_no|unit_task" docs/templates` devolve só TEMPLATE_INTENT.md, TEMPLATE_PLAN.md, TEMPLATE_ROUTE.json e TEMPLATE_ROUTE_LEGOPRESS.json, e a pasta não tem arquivo TEMPLATE_PIECE_UNIT. O molde do nó está em TEMPLATE_ROUTE.json:56-100.

**mcp.** O objetivo é acionar uma ferramenta de um servidor MCP. Exige `ref` sem espaço, barra invertida, dois-pontos nem `..` (parse.go:174-181). O `command` é opcional. O template da peça é docs/templates/TEMPLATE_PIECE_MCP.json, que registra o servidor e a lista de ferramentas (`tools`, com name e description) para busca sob demanda.

## Verificação interna dos executores (2026-11-20)

Método: li o código de cada executor e rodei os testes. `go test -count=1 ./internal/orchestrix ./internal/runner ./internal/gate ./internal/dag` passou nos quatro pacotes (orchestrix 17,3 s, runner 15,0 s, gate 2,7 s, dag 3,3 s). Rodei também, com `-v`, doze testes nomeados abaixo, e todos passaram.

### Ciclo comum a todo nó (internal/runner)

O runner despacha o nó por `ExecuteWithAndon` (runner/retry.go:15). O nó recebe no máximo 2 tentativas externas (`MaxNodeRetries = 2`, retry.go:11). Só vale sucesso quando `err == nil` e `ExitCode == 0` (retry.go:24). Na segunda falha, o runner puxa o Andon Cord e gera um relatório SBAR (retry.go:29-31).

Antes do executor, `preExecNode` aplica a fronteira Red/Blue, roda os hooks `pre` dos slots e aplica o orçamento de tokens (runner.go:329-340). Hook ou script de slot só roda com `--allow-scripts` (runner.go:289-291). Cada comando de slot tem 30 s e ambiente reduzido a PATH e SYSTEMROOT (runner/slots.go:20-25).

Depois do executor, `postExecAndPromote` roda os scripts e hooks `post`, valida o schema do artefato, avalia o gate do nó, os gates de slot e, com a etiqueta `ste100`, o validador STE-100 (runner.go:362-379, 453-468). Em modo `isolated`, o runner promove o artefato ao workdir principal só se o nó passou (runner.go:380-388).

O contexto do nó usa `timeout_sec`, ou o padrão do runner (runner.go:496-505).

### skill

O executor é `runSkill` (orchestrix.go:123). Com harness `direto`, ele chama o 9Router por HTTP sem harness e exige via `router` (orchestrix.go:291-314). Nos demais casos, o fluxo é: resolver rotas, montar o prompt, validar o combo, despachar, gravar o artefato.

As rotas seguem o modo (orchestrix.go:346-377). No modo `auto`, tenta o combo e depois um fallback do inventário. No modo `manual`, usa só o modelo travado. No modo `suggested`, usa o modelo e depois o fallback. Modelo fora do inventário não roda.

O prompt é estéril (prompt.go:24-40). Ele junta o núcleo inviolável (prompt.go:19), o SKILL.md lido do cofre, as regras por etiqueta, o tooling, o `context_in` e as entradas. O `context_in` passa por sanitização contra injeção de prompt (prompt.go:76-93). Cada arquivo lido tem teto de 256 KiB (prompt.go:16). Entrada `var.X` vira variável, `file('P')` vira artefato e o resto é literal (prompt.go:144-157).

Com duas ou mais rotas em modo `auto` ou `suggested`, o executor corre as duas primeiras de forma especulativa (orchestrix.go:172-217). Falha de rota alimenta um anticorpo e um contexto negativo no prompt seguinte (orchestrix.go:242-252).

O despacho por fonte cobre 9Router por HTTP, Ollama por HTTP, `claude -p` em diretório temporário vazio sem ferramentas, e harnesses (orchestrix.go:324-335; backends.go:59-198). O HTTP tenta de novo uma vez em 429 (backends.go:114-122). O artefato é gravado por arquivo temporário e rename, e recusa diretório-pai que seja link (orchestrix.go:383-426).

### script

O executor é `runScript` (backends.go:243). Sem `--allow-scripts`, o nó falha (backends.go:244-246). O executável é resolvido pelo PATH, e o nome `legoloop` aponta para o binário em execução (backends.go:263-268). Os argumentos expandem variáveis do `Vars` e do ambiente (backends.go:270-284). O processo roda sem shell, com ambiente estéril e saída limitada (backends.go:301-363). Se o script não gravou o artefato, o stdout vira o artefato (backends.go:257-259). O código de saída vai para o resultado e o runner o julga. Prova: `TestScriptRequiresConsent` passou.

### subprocesso

O executor é `runSubloop` (subloop.go:17). Ele carrega o Processo do acervo pelo `ref` (subloop.go:45-54). Cria um diretório temporário estéril, removido ao fim (subloop.go:22-26). Copia só as entradas declaradas, com path jail que recusa caminho fora do workspace, `.git`, `.env*`, `credential_*` e symlink (subloop.go:101-202). Roda o sub-DAG com paralelismo 1, herdando `AllowScripts`, `AutoApprove`, `NonInteractive` e `Vars` (subloop.go:56-66). Se o sub-DAG falha, devolve o resultado com as mensagens dos nós (subloop.go:204-221). Se passa, lê o artefato final (`sub_artifact`, ou `artifact`) e o grava no pai (subloop.go:225-267). Prova: `TestRunSubloopSuccess`, `TestRunSubloopNotFound` e `TestRunSubloopInputTransfer` passaram.

### stop_gate

O executor não faz nada: devolve o modelo `stop_gate/human` (orchestrix.go:116-118). A barreira vive no runner, que avalia `output.gate`, normalmente `prompt('pergunta')` (gate/builtins.go:212-220). Com `--yes`, o prompt aprova sozinho. Com `--non-interactive`, o prompt reprova com erro. Nos demais casos, usa o prompter do terminal (runner.go:518-536; cli/commands.go:213-218, 286-294). Sem prompter configurado, o gate falha (gate/builtins.go:212-214). Prova: `TestStopGateNodeFollowsHumanPrompt` passou.

### mcp

O executor é `runMCP` (mcp_client.go:248). O comando vem de `command` ou, na falta dele, é `legoloop mcp-router` (mcp_client.go:273-286). Com `ref` no formato `servidor/ferramenta` e sem `command`, a chamada vai por `call_mcp_tool` ao roteador. Com `command`, chama só a ferramenta (mcp_client.go:288-302). As entradas expandem variáveis (mcp_client.go:289-293).

O cliente fala JSON-RPC 2.0 por STDIO: `initialize`, `notifications/initialized`, `tools/call` (mcp_client.go:142-188). O timeout padrão é 60 s, ou o `timeout_sec` do nó (mcp_client.go:20, 254). Resposta de erro JSON-RPC ou `isError` vira falha (mcp_client.go:219-235). O stderr do servidor entra na mensagem de erro (mcp_client.go:237-246). O artefato só é gravado se o servidor não o criou antes (mcp_client.go:304-314). Prova: `TestMCPClientSTDIOAndRunnerIntegration` passou.

### unit_no

O executor é `runUnitTask` (unit_run.go:69). A ordem é esta:

1. Recusa sem `--allow-unit-tasks` (unit_run.go:73-75).
2. Prepara o plano: resolve o modelo, compila o gate, exige que o workdir seja a raiz do repositório git e monta o prompt (unit_run.go:51-66, 385-397).
3. Cria uma git worktree destacada de HEAD em diretório temporário (worktree.go:39-54) e a descarta sempre, inclusive em erro (worktree.go:57-67).
4. Semeia na worktree os arquivos de `context_in` que existem no repositório principal e ainda não foram commitados (unit_seed.go:20-49).
5. Roda o `command` do nó uma vez antes do modelo (a fase Red do TDD). Se houver comando de diagnóstico em `verification.diagnostic_command`, roda esse no lugar (unit_run.go:255-275; unit_retention.go:31-66). A saída entra no prompt (unit_run.go:277-305).
6. Chama o modelo na worktree. Com o claude CLI, as ferramentas são só Read, Edit, Write, Glob e Grep, sem Bash e sem Agent ou Task, com esforço e orçamento por natureza: `micro` e `nano` usam esforço `low` e US$ 0,80, as demais usam `medium` e US$ 2,50 (unit_claude.go:15-59). Outro harness só roda se estiver homologado (unit_run.go:346-354; unit_harness.go:22-24).
7. Lista os arquivos alterados, descontando os semeados e inalterados. Reprova se algum caminho atravessa link ou fica fora do write_boundary (unit_seed.go:89-106; unit.go:15-65).
8. Aplica a higiene em arquivos `.go`: remove comentários narrativos com termo proibido e roda gofmt (hygiene.go:21-119).
9. Roda o `command` e avalia o gate com a worktree como raiz (unit_run.go:201-238).
10. Se o gate reprova, repete até 2 vezes (`maxUnitRetries = 2`), com a falha anexada ao prompt, o que dá 3 chamadas ao modelo no máximo (unit_run.go:117-138, 172-184).
11. Se as tentativas esgotam, mas o número de violações caiu e o `go build` passa, retém o resultado como `parcial` (unit_run.go:141-158). Senão, salva o patch em `scratch/gates/<id>.failed.patch` e falha (unit_run.go:156, 407-421).
12. Se o gate passa, copia para o repositório principal só os arquivos alterados, com validação prévia de links e de arquivo regular (worktree.go:177-249). Grava o relatório JSON do nó como artefato (unit_run.go:97-101).

O write_boundary permitido é `context_out` mais `test_seam.file` (unit.go:17-23). Em caso de falha do modelo, quando a natureza é `nano` ou `micro`, o despacho tenta os demais candidatos elegíveis. As outras naturezas falham sem trocar de harness (despacho.go:27-42, 62-87).

Achado 1: o runner repete o nó até 2 vezes por fora (retry.go:19), e o unit_no repete até 3 chamadas por dentro. O teto real de chamadas ao modelo por nó é, portanto, 6, quando o erro vem do modelo ou do gate e não de uma recusa de segurança.

Achado 2: o campo `isolation.write_boundary` é exigido pelo schema v2 (dag/schema_v2.go:62), mas nenhum código de `internal/orchestrix` ou `internal/runner` lê o valor. O executor sempre aplica o limite estrito. O campo `retry_policy.max_retries` do template de ROUTE também só aparece no gerador (plangen/route.go:43), e o executor usa a constante 2. O template de peça do unit_no deve documentar isso para não prometer um controle que não existe.

Provas: `TestUnitTaskRefusedWithoutFlag`, `TestUnitTaskPassAppliesAndWritesArtifact`, `TestUnitTaskOutsideBoundaryFails`, `TestUnitTaskGateFailureDiscards`, `TestUnitTaskRequiresRepoRoot` e `TestUnitTaskStopsAfterMaxRetries` passaram.

### Correções dos achados 1 e 2 (2026-11-20)

Os achados acima descrevem o estado antes da correção. A correção seguiu TDD: escrevi os testes, vi que não compilavam (Red), implementei e rodei (Green).

Achado 1, corrigido. A reprovação do gate depois das retentativas agora é uma falha definitiva (`runner.Permanent`, internal/runner/permanent.go). O `ExecuteWithAndon` não repete o nó nesse caso e devolve `andon_cord: falha definitiva no nó ..., sem nova tentativa` (internal/runner/retry.go:24-28). O unit_no marca assim o erro final em `handleExhaustedRetries` (unit_run.go:156-157). Falha comum continua com 2 tentativas externas. O teto de chamadas ao modelo por nó caiu de 6 para 3. Provas: `TestAndonDoesNotRetryPermanentFailure`, `TestAndonStillRetriesOrdinaryFailure`, `TestPermanentKeepsCause` e `TestUnitTaskGateExhaustionIsNotRetriedByRunner`.

Achado 2, corrigido pela opção [A]. O campo `retry_policy.max_retries` virou parte do contrato do nó (dag.Node, internal/dag/types.go), com valor de 0 a 2 e válido só em nó unit_no (internal/dag/parse.go, `validateUnitPolicy` e `validateNodeCommon`). O executor o lê em `unitRetries` (unit_run.go), limitado ao teto `dag.MaxUnitRetries = 2`. Sem o campo, vale o teto 2. O `isolation.write_boundary` de um unit_no só aceita `"strict"`; o stop_gate continua aceitando `"none"`. O schema schemas/dag.v2.json ganhou a propriedade `retry_policy` (0 a 2). Os campos `backoff_ms` e `fallback_model_on_retry` do ROUTE continuam fora do DAG e o parser estrito os recusa. Provas: `TestUnitTaskHonorsRetryPolicyMaxRetries` (0, 1 e 2 retentativas dão 1, 2 e 3 chamadas), `TestUnitRetriesCapsAtMaxUnitRetries`, `TestUnitNoWriteBoundaryOnlyStrict`, `TestStopGateKeepsWriteBoundaryNone` e `TestRetryPolicyRules`.

Verificação: `go build ./...` passou. Os pacotes orchestrix, runner, gate, dag, golden, docs/templates, plangen, fabricar, studiodoc e studio passaram. Dois testes falhavam, e os dois já falhavam no HEAD limpo (git worktree temporário, removido). Nenhum vinha da minha mudança. Corrigi os dois, e `go test ./...` passa inteiro (55 pacotes ok, zero FAIL, código de saída 0).

`TestInstalarSkills` (internal/cli): o teste fixava 38 funções nativas, e a skill `legoloop-diagram-design` entrou em 2026-11-18 (commit 624bc9f) sem atualizar a contagem. A lista vem das pastas embutidas (skills/embed.go), então o teste ficou desatualizado. Atualizei para 39 e incluí `legoloop-diagram-design` na checagem de skills embutidas.

`TestProcessosCincoPorCinco` (internal/processo) tinha quatro causas, todas em processos gerados pelo diagram-design:
1. `processos/listar/pipeline.dag.json` usava a função de gate `file_size`, que não existe no catálogo (internal/gate/catalog.go). Troquei por `word_count(...) > 200`, que existe.
2. O nó `e05` do mesmo DAG tinha `type: "gate"`, que não é tipo de nó. Ele roda `verify_geometry.py` por `command`, então o tipo certo é `script`, e o gate fica em `output.gate`.
3. `processos/listar/anticorpos.json` usava um formato próprio (chave `anticorpos`). Converti para `anticorpos.v1`: `schema_version`, `processo`, `globais` e `dominio`, com `assinatura` (SHA-256 da descrição), `erro`, `antidoto` e `recorrencias`. Mantive os campos originais (`id`, `no`, `tipo`, `severidade`, `verificacao`) como extras.
4. Faltavam `triagem.json` e `memoria_execucao.json` em `listar`, e `memoria_execucao.json` em `legoloop-agent-desktop`. Gerei a triagem pelo gerador oficial `triagem.GerarTriagemPadrao` e criei as memórias vazias no formato `memoria_execucao.v1`.
