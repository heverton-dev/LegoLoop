# Relatório: `swarm run` sobre o DAG real da peça `legoloop-swarm`

Data: 2026-11-16. Janela da execução: 17:11:34Z a 17:11:38Z (4 s).

## 1. O que foi executado

Comando (binário compilado do repositório, sem harness externo):

```bash
legoloop swarm run processos/legoloop-swarm/pipeline.dag.json \
  -dir <repo-externo> -repo <repo-externo> -timeout 150s \
  -- bash scripts/swarm_peca_etapa.sh
```

- Cada nó `TICK-01.N` do DAG virou um worker real (processo do sistema operacional) em worktree git próprio (`workspace: share`).
- O corpo do worker (`scripts/swarm_peca_etapa.sh`) roda a etapa do nó com comandos `legoloop swarm` reais, num quadro de demonstração separado (`demo`).
- `SG-FINAL` é `stop_gate` humano: o `swarm run` o pulou (`skipped`).

## 2. Resultado por nível

Todos os 5 níveis saíram com código 0 e status `stopped`. O nível seguinte só começou depois do fim do anterior.

| Nível | Worker | Início (UTC) | Fim (UTC) | Código |
| :--- | :--- | :--- | :--- | :--- |
| 0 | TICK-01-1 | 17:11:35.011 | 17:11:35.418 | 0 |
| 1 | TICK-01-2 | 17:11:35.418 | 17:11:36.248 | 0 |
| 2 | TICK-01-3 | 17:11:36.248 | 17:11:36.874 | 0 |
| 3 | TICK-01-4 | 17:11:36.874 | 17:11:37.484 | 0 |
| 4 | TICK-01-5 | 17:11:37.484 | 17:11:38.421 | 0 |

`swarm run` saiu com código 1 em cada falha simulada nos testes e com código 0 aqui.

## 3. Gates

Os 5 gates do DAG foram avaliados pelo motor `internal/gate` (`Compile` e `Eval`) sobre os artefatos gerados:

| Nó | Artefato | Gate |
| :--- | :--- | :--- |
| TICK-01.1 | `out/01_niveis.md` | verdadeiro |
| TICK-01.2 | `out/02_workers.md` | verdadeiro |
| TICK-01.3 | `out/03_spawn.md` | verdadeiro |
| TICK-01.4 | `out/04_coordenacao.md` | verdadeiro |
| TICK-01.5 | `out/05_encerramento.md` | verdadeiro |

Controle negativo (feito antes, sobre os gates do `processo.json`, que têm o mesmo texto): num `out/` vazio os 5 gates falham.

## 4. Evidências nos artefatos

- `01_niveis.md`: `legoloop validate` saiu com código 0 e listou os níveis 0 a 5.
- `02_workers.md`: `w1` (engenheiro-backend) e `w2` (engenheiro-tdd), cada um com worktree próprio.
- `03_spawn.md`: `w1` pid 32508 e `w2` pid 20412, ambos `working`, "paralelo: 2 pids distintos vivos".
- `04_coordenacao.md`: mensagem `@w1` encerrou só o `w1` (`stopped concluido`); `w2` seguiu `working`; `monitor` listou 0 parados.
- `05_encerramento.md`: `w1` e `w2` `killed`, PID 0, worktree removido, `git worktree list` com 1 entrada.

## 5. Limpeza

- Quadro externo: os 5 workers `killed`, PID 0, `git worktree list` com 1 entrada.
- Quadro de demonstração: `w1` e `w2` `killed`, PID 0, `git worktree list` com 1 entrada.
- Processos: os 4 `legoloop.exe` ativos depois da execução são `legoloop mcp-router` da sessão, nenhum worker.

## 6. Limites desta prova

1. O DAG da peça é serial: o paralelismo real aparece dentro da etapa e03 (2 workers `w1` e `w2` no quadro de demonstração), não entre níveis do DAG externo. O paralelismo entre nós do mesmo nível está provado em `internal/cli/swarm_run_e2e_test.go` (barreira) e em `internal/mcp/swarm_run_test.go`.
2. Os workers executam um script determinístico, não agentes de LLM. A prova cobre o orquestrador (níveis, worktrees, códigos de saída, limpeza), não a qualidade de um agente.
3. O `SG-FINAL` (stop-gate humano) não foi avaliado.
4. O quadro externo usou um repositório git descartável, não o repositório do LegoLoop.

## 7. Como repetir

1. Compile: `go build -o legoloop.exe ./cmd/legoloop`.
2. Crie um repositório git descartável com um commit.
3. Exporte `BIN`, `PROJ`, `OUT` e `DEMO` e rode o comando da seção 1 com o `scripts/swarm_peca_etapa.sh` do ensaio.
4. Avalie os gates do DAG com `internal/gate` sobre `OUT`.
