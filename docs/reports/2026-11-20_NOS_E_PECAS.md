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

## Ação 2: proposta do template da peça unit_no (para aprovação, nada criado ainda)

### Diretrizes lidas

As diretrizes vêm de docs/templates/README.md e dos testes de docs/templates (templates_test.go e sexteto_test.go). Elas impõem o seguinte:

1. Toda peça declara `os_support`, `execution_modes` e `harness_support`, usa `load: lazy` e `requires.class` (README seção 1; `TestPecasAgnosticas`).
2. Skill não fixa modelo e MCP não fixa `autoApprove`. Para o unit_no, vale o mesmo espírito: o template declara classe e natureza, e o modelo é exemplo editável.
3. Os JSON não têm comentários. A orientação vive no README (seção 2) e num documento de apoio.
4. Os JSON de exemplo de Processo e de nascimento não têm placeholders `{{...}}` e são válidos como estão (`TestTemplatesNascimento`, `TestNativeFragmentsHaveNoPlaceholdersAndScriptHasBOM`).
5. O texto não pode ter sinônimo proibido da taxonomia: limite zero em `*.md`, `*.json`, `*.typ` e `*.go` do diretório (`TestTaxonomiaTemplatesZeroViolacoes`).
6. Toda expressão de gate precisa compilar (`TestEveryGateExpressionCompiles`).
7. O nó com `subagents_allowed` só aceita `false` ou `"ephemeral"` com `write_boundary: strict` (`TestContratoSubagenteEfemeroNosTemplates`).

### Achados novos nos templates existentes (comando rodado: `go run ./cmd/legoloop validate <arquivo>`)

`TEMPLATE_PIECE_PROCESSO.json` não passa em `legoloop validate`: `dag: json inválido: json: unknown field "execution_modes"`. O parser é estrito, e os campos de peça agnóstica (`execution_modes`, `os_support`, `harness_support`, `load`, `requires`) não existem no contrato do DAG. O README (seção 2) afirma que o Processo "passa em `legoloop validate`", e isso é falso para este arquivo. O arquivo que passa é `TEMPLATE_PIECE_LOOP.json`, e é ele que o teste `TestProcessoTemplateIsValidDAGAndDocument` valida.

`TEMPLATE_DAG.json` também não passa: `campo obrigatório do schema_version 2 ausente: nó "n01_ingestao": routing.via`.

Os testes atuais não detectam nenhum dos dois. Não corrigi, para não ampliar o escopo sem sua decisão.

### Como o unit_no funciona (resumo da verificação interna)

O unit_no é o nó que edita código com um modelo. Ele trabalha numa git worktree descartável criada de HEAD, roda o `command` antes (Red) e depois (gate) da edição, e só copia os arquivos para o repositório se o gate passar. Exige `--allow-unit-tasks`, workdir na raiz do repositório git e modelo do claude CLI (aliases `code-fast` e `code-pro` viram `sonnet`, `flash` vira `haiku`). Usa só as ferramentas Read, Edit, Write, Glob e Grep. Limita a escrita a `context_out` mais `test_seam.file`. Repete até `retry_policy.max_retries` vezes (0 a 2) e, se esgotar, falha de forma definitiva sem repetir o nó, ou retém resultado `parcial` quando as violações diminuíram e o build passa.

### Forma proposta do template

Arquivo: `docs/templates/TEMPLATE_PIECE_UNIT.json`, um DAG v2 de um único nó `unit_no`, válido no parser real (`dag.Parse` e `legoloop validate`).

Topo: `$schema`, `schema_version: 2`, `id` e `name` concretos (`exemplo-peca-unit-no`), `description`, `version`, `date`, `upstream_signatures`, `execution_strategy` e `trigger: manual`.

O nó (todos os campos que o schema v2 exige para unit_no, mais os que o executor lê):
- `type: unit_no`, `nature: micro`, `name`, `spec_ref`, `status`.
- `routing`: modo `suggested`, harness `claude` e modelo `sonnet` como preferência editável, `fallback` para um modelo de outra origem, esforço `medium` e `nature: micro`. Com a correção de universalidade abaixo, o nó não quebra sem o `claude`: ele cai no próximo candidato homologado da classe.
- `command`: argv de verificação sem shell (por exemplo `go test ./...`), que define o `exit_code` do gate.
- `context_in` (arquivos que o modelo lê), `context_out` (os únicos que ele pode escrever) e `test_seam` com `file` e `signature` (o teste-alvo do TDD).
- `isolation`: `write_boundary: strict`, `subagents_allowed: false`, `workspace: local`, `executor: orchestrix_hermetic`.
- `retry_policy: {max_retries: 2}` (campo corrigido na Ação 1).
- `limits`, `tooling`, `rollback_strategy`, `verification` (com `diagnostic_command` opcional), `timeout_sec`.
- `output`: `artifact` (relatório JSON do nó) e `gate` como `exit_code == 0 && file_exists('...')`.

Para cumprir a diretriz 1 sem quebrar o parser, a decisão final é uma terceira opção, [C], que substitui as opções [A] e [B] da primeira versão: os campos agnósticos (`execution_modes`, `os_support`, `harness_support`, `load`, `requires.class`) ficam dentro de `execution_strategy`, o mapa livre do contrato. Opinião técnica: [C] é a melhor porque mantém a declaração no próprio arquivo da peça (a diretriz e o teste `TestPecasAgnosticas` continuam valendo), passa no parser real sem alterar o contrato, e preserva os campos no Studio, cujo `Processo` já copia `execution_strategy`. A opção de ampliar o contrato com campos novos repetiria o defeito do Achado 2 (campos que nenhum código lê). O template do unit_no segue o mesmo padrão, e o documento de apoio `TEMPLATE_PIECE_UNIT.md` traz só a mecânica.

### Testes que acompanham (TDD, escritos antes do arquivo)

1. O template passa em `dag.Parse`, e o único nó é `unit_no`.
2. O template não tem `{{`.
3. O gate compila, e `retry_policy.max_retries` está entre 0 e 2.
4. `isolation.write_boundary` é `strict` e `subagents_allowed` é `false`.
5. Todo caminho de `context_out` e o `test_seam.file` aparecem como escrita permitida (espelha `allowedWrites`).
6. O README cita o template na tabela de peças, e o documento de apoio existe.
7. A taxonomia dá zero violações nos dois arquivos novos.


## Resolução dos achados novos nos templates e da universalidade (2026-11-20)

Regra do usuário: todo achado se corrige no ato, e decisão técnica sem consenso sobe ao humano antes da correção.

Corrigido, com TDD (testes Red antes, Green depois):

1. `TEMPLATE_PIECE_PROCESSO.json` agora passa em `legoloop validate` e em `dag.Parse`. Reescrevi como DAG v2 completo, com os campos agnósticos em `execution_strategy`, sem o campo `routing.class` (que não existe no contrato), com o stop_gate no roteamento canônico de `legoloop rota` (`manual`, `local`, `oficial`, `human/operator`) e com artefato próprio para o stop_gate (o validador exige artefato único por nó).
2. `TEMPLATE_DAG.json` agora passa no parser: faltava `routing.via` em todos os nós. Os dois scripts dependiam de PowerShell (`Test-Path`, `Copy-Item`), o que quebra em Linux e macOS; troquei por `legoloop version` e `git rev-parse HEAD`, que rodam em qualquer sistema.
3. A afirmação do README de que o Processo "passa em `legoloop validate`" agora é verdadeira e imposta por teste: `TestTemplatesEmFormaDeDAGPassamNoParser` (cinco templates em forma de DAG) e `TestTemplatesEmFormaDeDAGSemShellDeUmSO` (nenhum `command` fixa `powershell`, `pwsh`, `cmd`, `bash` ou `sh`). O README também documenta onde fica a declaração agnóstica do Processo e a regra de roteamento como preferência.
4. Universalidade do executor: um nó `unit_no` com `claude` declarado quebrava numa máquina sem o `claude`, porque `prepareUnit` falhava antes do despacho (erro `nenhum modelo de [sonnet ] está no inventário`) e a degradação de `nano` e `micro` nunca rodava. Agora, nas classes que degradam, o plano marca o candidato como indisponível e o despacho vai direto ao próximo candidato homologado da classe, registrando a substituição nas notas e no aviso (internal/orchestrix/unit_run.go e despacho.go). Nas demais classes, o nó falha sem trocar de harness, com mensagem que cita o harness declarado. Provas: `TestNoDeCodigoSemClaudeDegradaNanoMicro` e `TestNoDeCodigoSemClaudeFalhaForaDeNanoMicro`.

Verificação: `go test -count=1 ./...` passa em 55 pacotes, com código de saída 0.

Decisão humana ([1]): `tooling.runtime` passa a declarar o motor, `"legoloop"` (o Orchestrix executa sem shell), no lugar de `"powershell"`. Aplicado e provado:

- Geradores: internal/plangen/route.go (2 pontos) e internal/processo/projecoes.go.
- Templates: TEMPLATE_ROUTE.json, TEMPLATE_ROUTE_LEGOPRESS.json e TEMPLATE_PLAN.md.
- Processos vivos: 33 arquivos `.json` em `processos/`.
- Fixtures de teste: dag_test.go, subloop_parse_test.go e store_test.go.
- Provas (Red antes, Green depois): `TestRouteGeradoDeclaraRuntimeDoMotor` (plangen), `TestProjecoesDeclaramRuntimeDoMotor` (processo) e `TestTemplatesSemRuntimeDeUmSO` (templates).
- Não alterei os registros datados `docs/dags/` (28 arquivos) e `docs/routes/` (13 arquivos), que são histórico imutável por convenção de data, nem o texto da auditoria de 2026-11-10 que já descrevia o defeito.

Achado extra, corrigido: na rodada completa, `TestUnifiedPathFileLockExactMentionAndPIDKill` (internal/board) falhou uma vez com `TempDir RemoveAll cleanup ... being used by another process`. A pasta `004` é a do quadro `sbNat`, usada pelo worker que sai sozinho. O Windows libera os handles de um processo recém-terminado com atraso, e a remoção única do `t.TempDir` falha sob carga. O pacote passava isolado, então a falha é de corrida de limpeza e não de lógica. Troquei a pasta por `tempDirRobusto` (swarm_test.go), que remove com retentativa por até 10 s. Prova: a suíte inteira passou (55 pacotes, código 0) e o pacote passou 3 vezes seguidas com `-count=3`. A falha original não reproduzia isolada, então a prova é por ausência de nova falha, não por reprodução.

Verificação final: `go test -count=1 ./...` passou em 55 pacotes, com código de saída 0.
