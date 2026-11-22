---
schema: progresso.v1
iniciativa: "I03"
plano: "docs/plans/2026-11-21_PLAN_INICIATIVA_3_I03_DO_MACROPROJE_REVISAO.md"
atualizado_em: "2026-11-21T00:00:00Z"
---

# Progresso Vivo: I03

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 20 | 0 | 0 | 0 | 0 | 20 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-I3-01.1` | E01 | Consolidar testes de coerência executor-tipo, derivação de ausente no dag.Parse e listagem completa via ErrExecutorIncoerente | AC-E01.1, AC-E01.2 |
| `TICK-I3-01.2` | E01 | Recusar subagents_allowed ephemeral fora do native_swarm com ErrExecutorIncoerente | AC-E01.3 |
| `TICK-I3-01.3` | E01 | Exigir papel e command e recusar stop_gate no native_swarm com ErrExecutorIncoerente | AC-E01.4 |
| `TICK-I3-01.4` | E01 | Provar que dag.Parse e legoloop validate recusam executor incoerente ou desconhecido | AC-E01.5 |
| `TICK-I3-02.1` | E02 | Criar o campo executor do Ticket coerente com o tipo e gravá-lo na ROUTE | AC-E02.1 |
| `TICK-I3-02.2` | E02 | Recusar Ticket com executor incoerente ou desconhecido sem gravar PLAN nem ROUTE | AC-E02.2 |
| `TICK-I3-02.3` | E02 | Provar que a projeção do DAG preserva native_swarm e recusa executor incoerente | AC-E02.3 |
| `TICK-I3-03.1` | E03 | Criar Options.Swarm e a validação de executores antes do nível 0 em runner.Run | AC-E03.1 |
| `TICK-I3-03.2` | E03 | Despachar por nível: native_swarm ao SwarmRunner e os demais ao Executor | AC-E03.2, AC-NF.1 |
| `TICK-I3-03.3` | E03 | Rodar o nível misto em paralelo e terminar só quando os dois lados terminam | AC-E03.3 |
| `TICK-I3-03.4` | E03 | Converter o resultado do Swarm em NodeReport e marcar os níveis seguintes como skipped | AC-E03.4 |
| `TICK-I3-03.5` | E03 | Garantir que DAG sem native_swarm não toca no Swarm | AC-E03.5 |
| `TICK-I3-03.6` | E03 | Provar que legoloop run recusa executor inválido com exit 1 e sem criar worktree | AC-E03.6 |
| `TICK-I3-04.1` | E04 | Extrair board.RunNodeOf de LevelsFromDAG preservando cobertura do board >= 80% | AC-E04.1 |
| `TICK-I3-04.2` | E04 | Ligar legoloop run ao SwarmBoard real quando houver nó native_swarm | AC-E04.2 |
| `TICK-I3-04.3` | E04 | Provar que sem native_swarm o diretório do Swarm não é criado | AC-E04.3 |
| `TICK-I3-05.1` | E05 | Provar que todos os DAGs do repositório continuam válidos | AC-E05.1 |
| `TICK-I3-05.2` | E05 | Provar com golden que o relatório de DAG sem native_swarm não mudou | AC-E05.2 |
| `TICK-I3-05.3` | E05 | Atualizar AGENTS, README e MEMORY com o despachante e o native_swarm | AC-E05.3 |
| `TICK-I3-05.4` | E05 | Registrar os anticorpos, atualizar emendas, prompt da I03 e rodar o portão completo | AC-E05.4, AC-NF.2 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| nenhum | - | - | - | - |
