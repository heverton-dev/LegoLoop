# PLAN — Orchestrix executa `unit_task` (worktree + claude com ferramentas)

> **Origem:** achado do TICK-01.6 em 2026-11-10. `Orchestrix.Execute` só despacha `script`, `skill` e `subloop`. O DAG espelho da ROUTE tem 34 nós `unit_task` e `legoloop run` os recusa.
> **Decisão do usuário (2026-11-10):** modelo "worktree + claude com ferramentas".
> **Fora do PLAN da Lei #15:** esta trilha é paralela. Nenhum ticket dela renomeia termo da Lei #15.

## 1. Modelo de execução

1. A Orchestrix cria uma git worktree efêmera a partir de `HEAD`.
2. A Orchestrix roda `claude -p` dentro da worktree, com uma lista fechada de ferramentas.
3. A Orchestrix lista os arquivos alterados e compara com o `write_boundary` do nó.
4. A Orchestrix roda `command` (o primeiro argv de `verify`) dentro da worktree.
5. O gate do nó decide. Se o gate passar, a Orchestrix aplica o diff no `workdir` principal. Se falhar, descarta a worktree.
6. A Orchestrix grava `output.artifact` (JSON do resultado) no `workdir` principal.

## 2. Invariantes de segurança

- `unit_task` só roda com a flag explícita `--allow-unit-tasks` (mesmo padrão de `--allow-scripts`).
- Cofre real, `.git` e `~/.claude` ficam fora do alcance: a worktree não contém o Cofre.
- Caminho alterado fora de `output.artifacts` e do arquivo de teste do `test_seam` reprova o nó (`write_boundary: strict`).
- Link simbólico ou junction em caminho alterado reprova o nó.
- A Orchestrix nunca faz commit. `commit_on_node_pass` continua `false`.
- Sem sub-subagentes: o `claude` da worktree roda sem a ferramenta de agentes.

## 3. Tickets

| Ticket | Ação | Nature / Modelo | Dep. | Gate |
| :--- | :--- | :--- | :--- | :--- |
| `UT-01` | Funções puras: `buildUnitPrompt` (context_in, spec_ref, test_seam, gate, saídas permitidas) e `checkBoundary` (alterados contra permitidos; recusa caminho absoluto e `..`). | `micro` Sonnet | `[]` | `go test ./internal/orchestrix/...` verde, com teste de mutação do `checkBoundary` |
| `UT-02` | Gerenciador de worktree: criar, listar alterados (`git status --porcelain -z`), recusar link simbólico e junction em caminho alterado, aplicar diff no principal, descartar. Testes com repositório real em `t.TempDir()`. | `macro` Sonnet | `[]` | `go test ./internal/orchestrix/...` verde; worktree removida mesmo em erro; link em caminho alterado reprova |
| `UT-03` | `callClaudeTools`: `claude -p` na worktree com `--allowed-tools` fechado e sem agentes. Testes com o `claude` falso do `TestMain`. | `micro` Sonnet | `UT-01` | Argumentos exatos verificados no teste; ferramenta fora da lista reprova |
| `UT-04` | `runUnitTask` e a flag `--allow-unit-tasks`: junta UT-01 a UT-03, roda o `command`, avalia o gate, aplica ou descarta, grava o artefato. `Execute` despacha `unit_task`. | `macro` Sonnet | `UT-01`, `UT-02`, `UT-03` | `legoloop run` executa um DAG mínimo de 2 `unit_task` em repositório temporário; nó fora do `write_boundary` reprova |
| `UT-05` | Semear `context_in` existente do principal na worktree; excluir os semeados e inalterados de `Changed`; recusar `.git`, link e caminho fora do projeto. | `micro` Sonnet | `UT-04` | Command do nó 2 lê o arquivo do nó 1 sem commit; alterar arquivo semeado reprova |

#### 🛑 STOP-GATE HUMANO — `SG-UT`
* **Gate:** `prompt('A Orchestrix executa unit_task com segurança? Autoriza rodar o DAG espelho da Lei #15?')`
- [x] Usuário aprova depois de ler o diff de `internal/orchestrix/`.

## 4. Verificação por ticket

`go test ./...` verde e `go run ./tools/ticketgate -pkgs ./internal/orchestrix/...`.
