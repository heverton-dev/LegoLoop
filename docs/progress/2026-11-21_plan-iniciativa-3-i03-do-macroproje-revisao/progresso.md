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
| 20 | 20 | 0 | 0 | 0 | 0 |

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
| nenhum | - | - | - |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| `TICK-I3-01.1` | E01 | 0 | go test ./internal/dag/ -run "TestExecutorCoerenteEDerivado/TestValidarExecutoresListaTodos" (exit 0) | 2026-11-23T14:21:00Z |
| `TICK-I3-01.2` | E01 | 0 | go test ./internal/dag/ -run TestSubagentesEfemerosSoNoSwarm (exit 0) | 2026-11-23T14:30:00Z |
| `TICK-I3-01.3` | E01 | 0 | go test ./internal/dag/ -run TestNativeSwarmExigePapelECommand (exit 0) | 2026-11-23T14:40:00Z |
| `TICK-I3-01.4` | E01 | 0 | go test ./internal/dag/... ./internal/cli/ -run TestValidateRecusaExecutor (exit 0) | 2026-11-23T14:50:00Z |
| `TICK-I3-02.1` | E02 | 0 | go test ./internal/plangen/ -run TestTicketExecutor (exit 0) | 2026-11-23T15:00:00Z |
| `TICK-I3-02.2` | E02 | 0 | go test ./internal/plangen/ -run TestTicketExecutorIncoerente (exit 0) | 2026-11-23T15:30:00Z |
| `TICK-I3-02.3` | E02 | 0 | go test ./internal/routedag/... (exit 0) | 2026-11-23T16:00:00Z |
| `TICK-I3-03.1` | E03 | 0 | go test ./internal/runner/ -run TestRunRecusaAntesDoNivelZero (exit 0) | 2026-11-23T16:40:00Z |
| `TICK-I3-03.2` | E03 | 0 | go test ./internal/runner/ -run TestDespachoPorNivel (exit 0) | 2026-11-23T17:20:00Z |
| `TICK-I3-03.3` | E03 | 0 | go test ./internal/runner/ -run TestNivelMistoEsperaOsDoisLados (exit 0) | 2026-11-23T18:00:00Z |
| `TICK-I3-03.4` | E03 | 0 | go test ./internal/runner/ -run TestFalhaNoSwarmMarcaSkipped (exit 0) | 2026-11-23T18:00:00Z |
| `TICK-I3-03.5` | E03 | 0 | go test ./internal/runner/ -run TestSemNativeSwarmNaoUsaSwarm (exit 0) | 2026-11-23T18:00:00Z |
| `TICK-I3-03.6` | E03 | 0 | go test ./internal/cli/ -run TestRunRecusaExecutorSemWorktree (exit 0) | 2026-11-23T18:25:00Z |
| `TICK-I3-04.1` | E04 | 0 | go test -cover ./internal/board/... (exit 0, coverage: 85.4%) | 2026-11-23T20:53:00Z |
| `TICK-I3-04.2` | E04 | 0 | go test ./internal/cli/ -run TestRunNativeSwarmBinario (exit 0) | 2026-11-23T21:18:00Z |
| `TICK-I3-04.3` | E04 | 0 | go test ./internal/cli/ -run TestRunSemNativeSwarmNaoCriaDiretorio (exit 0) | 2026-11-23T21:18:00Z |
| `TICK-I3-05.1` | E05 | 0 | go test ./internal/dag/ -run TestDAGsDoRepositorioValidam (exit 0) | 2026-11-23T21:40:00Z |
| `TICK-I3-05.2` | E05 | 0 | go test ./internal/runner/ -run TestParidadeDoRelatorio (exit 0) | 2026-11-23T21:40:00Z |
| `TICK-I3-05.3` | E05 | 0 | AGENTS.md (11212 bytes <= 24000); native_swarm em AGENTS.md, README.md, MEMORY.md | 2026-11-23T21:40:00Z |
| `TICK-I3-05.4` | E05 | 0 | Emenda 2 atualizada; anticorpos verificados; portão completo verde | 2026-11-23T21:40:00Z |
