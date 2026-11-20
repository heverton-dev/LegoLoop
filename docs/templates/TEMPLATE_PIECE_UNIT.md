# Peça de nó unit_no (template v2.4.0)

O arquivo `TEMPLATE_PIECE_UNIT.json` é um DAG v2 de um único nó `unit_no`. Ele é um exemplo válido: passa em `legoloop validate` e em `dag.Parse`. Edite os valores para o seu Processo. Este documento explica a mecânica do nó e cada campo que o executor lê.

## O que o nó faz

O `unit_no` é o nó que edita código com um modelo. A Orchestrix executa o nó nesta ordem:

1. Recusa o nó sem o consentimento do usuário: rode com `legoloop run <dag> --allow-unit-no`. O nome antigo da flag continua aceito como alias legado.
2. Exige que o diretório de trabalho seja a raiz de um repositório git.
3. Cria uma git worktree descartável a partir de HEAD, fora do repositório.
4. Copia para a worktree os arquivos de `context_in` que existem no repositório e ainda não estão no commit.
5. Roda o `command` uma vez antes da edição (etapa Red do TDD) e anexa a falha ao prompt.
6. Chama o modelo na worktree, com as ferramentas Read, Edit, Write, Glob e Grep. O modelo não tem Bash e não cria subagentes.
7. Aceita só as escritas de `context_out` mais o arquivo de `test_seam.file`. Qualquer outro caminho alterado reprova o nó.
8. Para arquivos `.go`, remove comentários narrativos com termo proibido e roda gofmt.
9. Roda o `command` de novo e avalia o gate, com o `exit_code` do `command` como entrada.
10. Se o gate reprova, repete a chamada ao modelo até `retry_policy.max_retries` vezes, com a falha anexada ao prompt.
11. Se o gate passa, copia para o repositório só os arquivos alterados e grava o relatório JSON em `output.artifact`.

Se as tentativas esgotam, o nó falha de forma definitiva e o runner não o repete. Se o número de violações caiu e o `go build` passa, o nó termina como `parcial`. Nos dois casos, a worktree some e o repositório fica intacto, e o patch das tentativas é salvo em `scratch/gates/<id>.failed.patch`.

## Campos que o executor lê

* `command`: lista de argumentos sem shell. O executável é um nome no PATH ou um caminho absoluto. Use o comando de verificação do seu projeto. O exemplo usa `go test`.
* `context_in`: arquivos e pastas que o modelo lê. Arquivo ausente vira "ainda não existe" no prompt. Pasta vira uma indicação para o modelo explorar.
* `context_out`: os únicos arquivos que o modelo pode escrever. Pasta com barra final vale para tudo dentro dela.
* `test_seam`: `file` é o teste-alvo do TDD e também uma escrita permitida; `signature` é o nome do teste.
* `isolation.write_boundary`: aceita só `strict`. O validador recusa outro valor num nó `unit_no`.
* `retry_policy.max_retries`: de 0 a 2. Sem o campo, vale 2.
* `nature`: `nano` ou `micro` limita o esforço do modelo a `low` e o gasto a US$ 0,80 por chamada; as demais naturezas usam `medium` e US$ 2,50. A natureza também define a política de falha do roteamento (veja abaixo).
* `verification.diagnostic_command`: opcional. Lista de argumentos que substitui o `command` na etapa Red, quando o diagnóstico completo é mais útil ao modelo.
* `output.artifact`: o relatório JSON do nó. Ele não pode estar entre as escritas permitidas.
* `output.gate`: expressão que compila no DSL de gates. Use `exit_code == 0` mais as condições do seu artefato.

## Roteamento universal

O `routing` do template é uma preferência editável, não uma exigência. O exemplo declara `claude` e `sonnet` com fallback `haiku`, mas o nó não quebra numa máquina sem o `claude`:

* Nós `nano` e `micro` degradam para o próximo candidato homologado da classe, na ordem do perfil de roteamento, e a substituição aparece nas notas do nó e num aviso.
* As demais classes não trocam de harness e falham com uma mensagem que cita o harness declarado.
* O harness `direto` não edita arquivos e não serve para o `unit_no`.
* Para regenerar o roteamento de uma máquina, rode `legoloop rota`.

Por isso o template usa `nature: micro`. Mude a natureza para uma classe maior só quando o nó exigir um harness específico.

## Peça agnóstica

Como o parser do DAG recusa campo desconhecido no topo do documento, a declaração agnóstica (`execution_modes`, `os_support`, `harness_support`, `load` e `requires.class`) fica dentro de `execution_strategy`.

## Checklist antes de usar

* O `command` roda em qualquer sistema operacional que o projeto suporta e não fixa um shell.
* `context_out` lista só o que o modelo deve escrever, e nada do que o gate lê como prova.
* O gate falha de verdade quando o código está errado (rode o `command` antes da edição e veja o exit code diferente de zero).
* O nó passa em `legoloop validate`.
