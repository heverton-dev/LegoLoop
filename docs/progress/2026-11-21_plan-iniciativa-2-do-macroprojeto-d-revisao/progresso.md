---
schema: progresso.v1
iniciativa: "I02"
plano: "docs/plans/2026-11-21_PLAN_INICIATIVA_2_DO_MACROPROJETO_D_REVISAO.md"
atualizado_em: "2026-11-22T22:04:10Z"
---

# Progresso Vivo: I02

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 22 | 4 | 1 | 0 | 0 | 17 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| `TICK-I2-02.1` | E02 | i02-e02-fabricar-plano | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-I2-02.2` | E02 | Recusar sem INTENT aprovada com exit 2 e sem gravar nada | AC-E02.2 |
| `TICK-I2-02.3` | E02 | Tornar o modo estrito o padrão e criar o flag --permissivo | AC-E02.3 |
| `TICK-I2-02.4` | E02 | Gravar PLAN, ROUTE e DAG de forma atômica: falha em qualquer etapa não deixa arquivo | AC-E02.4 |
| `TICK-I2-02.5` | E02 | Remover legoloop plano e trocar toda referência por fabricar plano | AC-E02.5 |
| `TICK-I2-03.1` | E03 | Criar fabricar spec que grava a SPEC a partir de uma INTENT aprovada | AC-E03.1 |
| `TICK-I2-03.2` | E03 | Recusar INTENT não aprovada com exit 2 sem chamar o LLM | AC-E03.2 |
| `TICK-I2-03.3` | E03 | Criar o Processo legoloop-especificacao com contrato, nascimento 5 de 5 e Registro da Espinha | AC-E03.3, AC-E03.4 |
| `TICK-I2-04.1` | E04 | Criar o Processo legoloop-planejamento com contrato, nascimento 5 de 5 e Registro da Espinha | AC-E04.1, AC-E04.2 |
| `TICK-I2-04.2` | E04 | Fazer o DAG do planejamento chamar fabricar plano em vez de duplicar a lógica | AC-E04.3 |
| `TICK-I2-05.1` | E05 | Rotear gerar a SPEC para legoloop-especificacao sem LLM | AC-E05.1 |
| `TICK-I2-05.2` | E05 | Rotear gerar o PLAN, a ROUTE e o DAG para legoloop-planejamento sem LLM | AC-E05.2 |
| `TICK-I2-05.3` | E05 | Provar que a triagem não roubou entradas dos Processos existentes | AC-E05.3 |
| `TICK-I2-06.1` | E06 | Criar legoloop-ciclo 5 de 5, registrar na Espinha e compor concepção, especificação e planejamento | AC-E06.1 |
| `TICK-I2-06.2` | E06 | Provar que o produces de cada subprocesso casa com o consumes do seguinte | AC-E06.2 |
| `TICK-I2-07.1` | E07 | Atualizar AGENTS, README e MEMORY com os Processos novos | AC-E07.1 |
| `TICK-I2-07.2` | E07 | Registrar os anticorpos da I02 nos Processos novos | AC-E07.2 |
| `TICK-I2-07.3` | E07 | Atualizar emendas, prompt de execução da I02 e rodar o portão completo | AC-E07.3, AC-NF.2 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| `TICK-I2-01.1` | E01 | 0 | ok github.com/heverton-dev/LegoLoop/internal/routedag 0.072s (18 testes) | 2026-11-22T20:43:25Z |
| `TICK-I2-01.2` | E01 | 0 | ok internal/routedag e tools/route2dag e main.go 11 linhas | 2026-11-22T21:00:00Z |
| `TICK-I2-01.3` | E01 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/routedag	0.025s / | 2026-11-23T00:23:58Z |
| `TICK-I2-01.4` | E01 | 0 | ok  	github.com/heverton-dev/LegoLoop/internal/routedag	0.014s / | 2026-11-23T00:43:05Z |
