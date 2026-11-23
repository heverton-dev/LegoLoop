---
id: REPORT-2026-11-23_FECHAMENTO_I03
iniciativa: I03
status: APPROVED
fase: fechamento_aprovado
date: 2026-11-23
schema_version: 2
upstream:
  - 2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE
  - 2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE_REVISAO
  - 2026-11-21_PLAN_INICIATIVA_3_I03_DO_MACROPROJE_REVISAO
progresso: docs/progress/2026-11-21_plan-iniciativa-3-i03-do-macroproje-revisao/progresso.md
tickets: [TICK-I3-05.1, TICK-I3-05.2, TICK-I3-05.3, TICK-I3-05.4]
covers: [AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-NF.2]
interpretative: [1. Resumo Executivo, 6. Conclusão do Macroprojeto D e Stop-Gate Final]
---

# Relatório de Fechamento da I03 — Despachante de Executores no Runner

> **Status:** Fechamento pré-Stop-Gate final (`UNDER_REVIEW`)  
> **Data:** 2026-11-23  
> **Iniciativa:** I03 (Despachante de executores no runner)  
> **Referências:** [INTENT](../intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md) · [PLAN](../plans/2026-11-21_PLAN_INICIATIVA_3_I03_DO_MACROPROJE_REVISAO.md) · [Progresso](../progress/2026-11-21_plan-iniciativa-3-i03-do-macroproje-revisao/progresso.md)

---

## 1. Resumo Executivo

A Iniciativa I03 concluiu com êxito os 20 tickets de seu plano revisado (100% PASSED), entregando ao LegoLoop um despachante robusto de executores no motor de execução (`internal/runner`), plenamente integrado ao ecossistema do Swarm nativo (`internal/board`) e aderente a todas as leis operacionais e restrições de governança.

A entrega cumpre integralmente os 23 critérios de aceitação distribuídos em 5 épicos:
1. **Épico 01 (`internal/dag`):** Validação estrita de coerência entre tipo do nó e executor (`ErrExecutorIncoerente`), derivação automática de valor ausente no `dag.Parse`, confinamento exclusivo de subagentes efêmeros ao `native_swarm`, obrigatoriedade de `role`/`papel` e `command` para nós Swarm, e barreira mecânica no `legoloop validate` com código de saída 2.
2. **Épico 02 (`internal/plangen` e `internal/routedag`):** Incorporação do campo `Executor` ao `Ticket`, gravação atômica em `isolation.executor` na ROUTE, rejeição de tickets incoerentes antes da escrita em disco e projeção determinística no DAG com preservação do executor e papel padrão.
3. **Épico 03 (`internal/runner` e `internal/cli`):** Interface abstrata `SwarmRunner`, validação prévia de todos os nós antes do nível 0, despacho concorrente em níveis mistos (espera de ambos os lados), interrupção das camadas subsequentes em caso de falha no Swarm (`StatusSkipped`), e recusa no `legoloop run` com código 1 sem alocação indevida de worktree.
4. **Épico 04 (`internal/board` e CLI real):** Extração modular de `board.RunNodeOf`, manutenção da cobertura do pacote `internal/board` em 85,4% (acima do piso de 80%), abertura sob demanda do `SwarmBoard` apenas quando houver nós `native_swarm` e ausência de resíduos em execuções estritamente locais.
5. **Épico 05 (Paridade e Governança Viva):** Validação mecânica de todos os 36 processos e 26 DAGs históricos do repositório, teste de regressão golden assegurando saída idêntica para relatórios sem Swarm, sincronização documental viva (`AGENTS.md`, `README.md`, `MEMORY.md`), conformidade de tamanho e verificação completa do banco de anticorpos.

---

## 2. Rastreabilidade dos Tickets (20/20 PASSED)

| Ticket | Épico | Título | Alvo Principal | Critérios | Prova Mecânica |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `TICK-I3-01.1` | E01 | Coerência executor-tipo e listagem completa | `internal/dag/executor.go` | `AC-E01.1`, `AC-E01.2` | `go test ./internal/dag/ -run "TestExecutorCoerenteEDerivado\|TestValidarExecutoresListaTodos"` (exit 0) |
| `TICK-I3-01.2` | E01 | Subagentes efêmeros só no `native_swarm` | `internal/dag/executor.go` | `AC-E01.3` | `go test ./internal/dag/ -run TestSubagentesEfemerosSoNoSwarm` (exit 0) |
| `TICK-I3-01.3` | E01 | Papel, command e recusa de `stop_gate` no Swarm | `internal/dag/executor.go` | `AC-E01.4` | `go test ./internal/dag/ -run TestNativeSwarmExigePapelECommand` (exit 0) |
| `TICK-I3-01.4` | E01 | Recusa em `dag.Parse` e `legoloop validate` | `internal/cli/validate_executor_test.go` | `AC-E01.5` | `go test ./internal/dag/... ./internal/cli/ -run TestValidateRecusaExecutor` (exit 0) |
| `TICK-I3-02.1` | E02 | Campo `Executor` no Ticket e gravação na ROUTE | `internal/plangen/plangen.go` | `AC-E02.1` | `go test ./internal/plangen/ -run TestTicketExecutor` (exit 0) |
| `TICK-I3-02.2` | E02 | Recusa de Ticket incoerente sem tocar disco | `internal/plangen/plangen.go` | `AC-E02.2` | `go test ./internal/plangen/ -run TestTicketExecutorIncoerente` (exit 0) |
| `TICK-I3-02.3` | E02 | Projeção do DAG preserva Swarm e atribui papel | `internal/routedag/executor.go` | `AC-E02.3` | `go test ./internal/routedag/...` (exit 0) |
| `TICK-I3-03.1` | E03 | Interface `SwarmRunner` e checagem pré-nível zero | `internal/runner/runner.go` | `AC-E03.1` | `go test ./internal/runner/ -run TestRunRecusaAntesDoNivelZero` (exit 0) |
| `TICK-I3-03.2` | E03 | Despacho por nível conforme executor | `internal/runner/runner.go` | `AC-E03.2` | `go test ./internal/runner/ -run TestDespachoPorNivel` (exit 0) |
| `TICK-I3-03.3` | E03 | Nível misto concorrente com barreira de sincronização | `internal/runner/runner.go` | `AC-E03.3` | `go test ./internal/runner/ -run TestNivelMistoEsperaOsDoisLados` (exit 0) |
| `TICK-I3-03.4` | E03 | Falha no Swarm interrompe camadas com `skipped` | `internal/runner/runner.go` | `AC-E03.4` | `go test ./internal/runner/ -run TestFalhaNoSwarmMarcaSkipped` (exit 0) |
| `TICK-I3-03.5` | E03 | DAG sem Swarm não aciona infraestrutura externa | `internal/runner/runner.go` | `AC-E03.5` | `go test ./internal/runner/ -run TestSemNativeSwarmNaoUsaSwarm` (exit 0) |
| `TICK-I3-03.6` | E03 | `legoloop run` recusa incoerente sem criar worktree | `internal/cli/run_executor_test.go` | `AC-E03.6` | `go test ./internal/cli/ -run TestRunRecusaExecutorSemWorktree` (exit 0) |
| `TICK-I3-04.1` | E04 | Extração de `board.RunNodeOf` com cobertura >= 80% | `internal/board/swarm_dag.go` | `AC-E04.1` | `go test -cover ./internal/board/...` (exit 0, cob: 85.4%) |
| `TICK-I3-04.2` | E04 | `runVivo` abre SwarmBoard quando há `native_swarm` | `internal/cli/progresso_vivo.go` | `AC-E04.2` | `go test ./internal/cli/ -run TestRunNativeSwarmBinario` (exit 0) |
| `TICK-I3-04.3` | E04 | Ausência de pasta `.legoloop-swarm/` em DAG comum | `internal/cli/run_swarm_test.go` | `AC-E04.3` | `go test ./internal/cli/ -run TestRunSemNativeSwarmNaoCriaDiretorio` (exit 0) |
| `TICK-I3-05.1` | E05 | Validação contínua de todos os DAGs do repo | `internal/dag/executor_repo_test.go` | `AC-E05.1` | `go test ./internal/dag/ -run TestDAGsDoRepositorioValidam` (exit 0) |
| `TICK-I3-05.2` | E05 | Teste golden de paridade de relatório do runner | `internal/runner/paridade_test.go` | `AC-E05.2` | `go test ./internal/runner/ -run TestParidadeDoRelatorio` (exit 0) |
| `TICK-I3-05.3` | E05 | Sincronização viva de AGENTS, README e MEMORY | Documentação viva | `AC-E05.3` | `AGENTS.md` (11212 bytes <= 24000); citações validadas |
| `TICK-I3-05.4` | E05 | Verificação de anticorpos, emenda e portão geral | Governança e gates | `AC-E05.4`, `AC-NF.2` | `legoloop espinha verificar` e `progresso verificar` (exit 0) |

---

## 3. Conclusão do Macroprojeto D

Com a entrega da Iniciativa I03, as quatro iniciativas planejadas no Macroprojeto D estão plenamente concluídas e integradas:
- **I01 (Distribuição e Forja):** 38/38 tickets PASSED. Instalação universal, forja industrial em duas camadas e execução sob rigorosas restrições herméticas.
- **I04 (Espinha Dorsal):** 56/56 tickets PASSED. Descoberta de unidades multi-SO, path jail, matriz de componentes, pisos mecânicos de cobertura e mutação sem sobreviventes.
- **I02 (Processos de Ciclo):** 22/22 tickets PASSED. Geração de SPEC a partir de INTENT com checagem upstream, geração atômica de PLAN/ROUTE/DAG e macroprocesso `legoloop-ciclo`.
- **I03 (Despachante de Executores):** 20/20 tickets PASSED. Despacho dinâmico por nível para `local_process`, `orchestrix_hermetic` e `native_swarm`, com paridade e integridade garantidas.
