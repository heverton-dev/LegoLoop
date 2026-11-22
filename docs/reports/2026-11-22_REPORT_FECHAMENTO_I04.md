---
id: REPORT-2026-11-22_FECHAMENTO_I04
iniciativa: I04
status: UNDER_REVIEW
fase: fechamento_pre_stop_gate_final
date: 2026-11-22
schema_version: 2
upstream:
  - 2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE
  - 2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE
  - 2026-11-21_PLAN_INICIATIVA_4_I04_DO_MACROPROJE
insumo_auditoria: docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md
progresso: docs/progress/2026-11-21_plan-iniciativa-4-i04-do-macroproje/progresso.md
tickets: [TICK-I4-09.3, TICK-I4-09.4]
covers: [AC-E09.3, AC-E09.4]
interpretative: [1. Resumo Executivo, 6. Próximos Passos e Stop-Gate Final]
---

# Relatório de Fechamento da I04 — Espinha Dorsal Previsível, Universal e Determinística

> **Status:** Fechamento pré-Stop-Gate final (`UNDER_REVIEW`)  
> **Data:** 2026-11-22  
> **Iniciativa:** I04 (Espinha dorsal previsível, universal e determinística)  
> **Referências:** [Auditoria](2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md) · [PLAN](../plans/2026-11-21_PLAN_INICIATIVA_4_I04_DO_MACROPROJE.md) · [Progresso](../progress/2026-11-21_plan-iniciativa-4-i04-do-macroproje/progresso.md)

---

## 1. Resumo Executivo

A Iniciativa I04 fechou todas as 19 lacunas registradas na auditoria da espinha dorsal (G01 a G14 da seção 2 e G15 a G19 da seção 4A), estruturou o Registro da Espinha (`espinha/registro.json`) com o Gate `legoloop espinha verificar`, elevou a cobertura de `board`, `triagem` e `cli` acima de 80% (`espinha/cobertura.json`), configurou o workflow de CI multi-SO (`.github/workflows/ci.yml`) e concluiu a análise de mutação nos sete pacotes da espinha com 0 sobreviventes reais.

O congelamento das iniciativas I01, I02 e I03 foi preservado integralmente (`TICK-I4-09.3`), aguardando aprovação humana no Stop-Gate final da I04.

---

## 2. Desfecho das Lacunas da Auditoria (G01 a G19)

| ID | Sev. | Origem | Desfecho | Pacote / Arquivo | Teste / Gate de Fechamento |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **G01** | Crítica | Seção 2 / 4A | Fechada: `fabricar` detecta Matriz defasada e atualiza em modo sem cérebro; Matriz em dia não dispara atualização. | `internal/fabricar/engine.go` | `internal/fabricar/atualiza_matriz_test.go` (`TestFabricarAtualizaMatrizDefasada`, `TestFabricarNaoAtualizaMatrizEmDia`, `TestMatrizStatusZeroDepoisDoFabricar`) |
| **G02** | Crítica | Seção 2 / 4A | Fechada: descoberta de unidades fixas por SO com provedor injetável; `TargetDrives` restringe a busca e o aviso `drives_ignorados` foi removido. | `internal/descoberta/unidades.go`, `internal/descoberta/provedor_*.go`, `internal/fabricar/engine.go`, `internal/assets/discover.go` | `internal/descoberta/unidades_test.go`, `internal/descoberta/leitura_test.go`, `internal/fabricar/drives_test.go` (`TestTargetDrivesRestringeABusca`) |
| **G03** | Alta | Seção 2 | Fechada: `dag.Parse` deriva `isolation.executor` ausente pelo tipo do nó, recusa combinações incoerentes listando todos os nós e `validate`/`run` saem com exit 2; DAGs do repositório corrigidos. | `internal/dag/executor.go`, `internal/dag/parse.go`, `internal/cli/commands.go`, `processos/*/pipeline.dag.json`, `docs/dags/` | `internal/dag/executor_test.go`, `internal/dag/executor_repo_test.go`, `internal/cli/validate_executor_test.go` (`TestValidateERunRecusamIncoerenteComExit2`) |
| **G04** | Alta | Seção 2 | Fechada: tabela única `unit_task` e `unit_no` aplicada em `route2dag` e `runner`; alias duplicado `TypeSubloop` e literal solto removidos. | `internal/dag/tipos.go`, `internal/dag/types.go`, `tools/route2dag/convert.go`, `internal/runner/runner.go` | `internal/dag/tipos_test.go` (`TestTabelaDeTipos`), `go build ./...` |
| **G05** | Alta | Seção 2 / 4A | Fechada: constantes consolidadas em `internal/caminhos/nomes.go`, fachada `internal/paths` removida e jail extraído para `internal/jail` ligado aos verbos da CLI. | `internal/caminhos/nomes.go`, `internal/jail/jail.go`, `internal/cli/commands.go` (`internal/paths` removido) | `internal/caminhos/constantes_test.go`, `internal/cli/nome_hostil_test.go` (`TestNomeHostilEhRecusado`), `legoloop espinha verificar` |
| **G06** | Média | Seção 2 | Fechada: pacote morto `internal/vault` removido por não possuir importador em produção. | `internal/vault` (removido) | `go build ./...`, `go test -count=1 ./...` |
| **G07** | Média | Seção 2 | Fechada: resolução de fontes, `python`/`python3` e shell adaptada por SO; literais fixos de Windows removidos do código comum. | `internal/assets/argv.go`, `internal/videorea/cartoon_render.go` | `internal/assets/argv_test.go`, `go test ./internal/assets/ ./internal/videorea/ -run TestPorSO` |
| **G08** | Média | Seção 2 | Fechada: verbo descontinuado `matrix` removido da CLI e do catálogo; todos os verbos registrados em `espinha/registro.json` e verificados via `cli.Verbos`. | `internal/cli/cli.go`, `espinha/registro.json` | `internal/cli` (`TestVerboMatrixRemovido`), `internal/espinha/verbos_test.go` (`TestTodoVerboEstaNoRegistro`) |
| **G09** | Média | Seção 2 | Fechada: `ANTICORPOS.md` criado nos 3 Processos faltantes (`legoloop-agent-desktop`, `legoloop-matriz`, `legoloop-listar`); 9 Processos sem prefixo renomeados via `legoloop renomear` e registrados em `espinha/renomeados.json`. | `internal/renomear/renomear.go`, `internal/cli/renomear.go`, `espinha/renomeados.json`, `processos/*/ANTICORPOS.md` | `internal/renomear/renomear_test.go`, `internal/espinha/nomes_test.go` (`TestSemReferenciaAoNomeAntigo`), `legoloop espinha verificar` |
| **G10** | Média | Seção 2 | Fechada: criado Gate `tools/cobertura` com piso de 80% em `espinha/cobertura.json`; medições obtidas: `board` 85,5%, `cli` 82,5%, `triagem` 96,6%. | `tools/cobertura/main.go`, `espinha/cobertura.json`, `internal/board`, `internal/triagem`, `internal/cli` | `go run ./tools/cobertura` (exit 0) |
| **G11** | Média | Seção 2 / 4A | Fechada: criada gravação atômica `escrita.GravarAtomico` e aplicada nos pontos críticos identificados em `fechar`, `brain` e `cli`. | `internal/escrita/atomico.go`, `internal/fechar`, `internal/brain`, `internal/cli` | `internal/escrita/atomico_test.go`, `go test ./internal/escrita/... ./internal/fechar/... ./internal/brain/...` |
| **G12** | Baixa | Seção 2 | Fechada: falha ao gravar telemetria agora é relatada e o limite de erros ignorados foi fixado no Registro da Espinha. | `internal/runner`, `espinha/registro.json` | `internal/runner/telemetria_test.go` (`TestFalhaDeTelemetriaEhRelatada`), `legoloop espinha verificar` |
| **G13** | Baixa | Seção 2 | Fechada (definição do workflow): criado CI multi-SO (`ubuntu-latest`, `windows-latest`, `macos-latest`) com `go test`, `go vet`, `funclen`, `gofmt`, `taxlinter`, `espinha verificar`, `cobertura`, `go.mod`/`go.sum` limpos e `-race` em Linux/macOS. | `.github/workflows/ci.yml` | `.github/workflows/ci.yml` (ver ressalva de execução pós-push na seção 5) |
| **G14** | Baixa | Seção 2 / 4A | Fechada (sem achado na seção 4A + trava mecânica): nenhum pacote central da espinha tem estado global mutável; `init()` e tabelas somente leitura justificados no Registro. | `espinha/registro.json`, `internal/espinha/globais_test.go` | `internal/espinha/globais_test.go` (`TestGlobaisDaEspinhaEstaoJustificados`), `legoloop espinha verificar` |
| **G15** | Alta | Seção 4A | Fechada: jail recusa caminhos hostis (`..\`, `C:\`, `\\server\`) em qualquer SO; lock da memória da triagem exclusivo; guarda de deleção normaliza separadores; 6 testes dependentes de ferramentas externas tornados herméticos. | `internal/jail/jail.go`, `internal/triagem/memoria.go`, `internal/inventario/guard.go`, `internal/cli`, `internal/legopress` | `internal/jail/jail_test.go` (`TestEntradaHostil`), `internal/triagem/memoria_test.go` (`TestMemoriaLockOcupado`), `internal/inventario/guard_test.go` (`TestRecusaDelecaoTemplate`), `go test ./internal/cli/ ./internal/legopress/ -count=1` |
| **G16** | Alta | Seção 4A | Fechada: chamadas diretas a `powershell` e `python` nos DAGs dos Processos substituídas por iniciador portável, bloqueadas por teste de regressão. | `processos/*/pipeline.dag.json`, `internal/espinha/so_test.go` | `internal/espinha/so_test.go` (`TestProcessosSemComandoExclusivoDoSO`) |
| **G17** | Média | Seção 4A | Fechada: Processo duplicado `processos/fechar` aposentado e removido; anticorpos consolidados em `processos/legoloop-fechar` e mapeamento gravado em `espinha/renomeados.json`. | `processos/legoloop-fechar/ANTICORPOS.md`, `espinha/renomeados.json` (`processos/fechar` removido) | `legoloop anticorpos status legoloop-fechar`, `internal/espinha/nomes_test.go` |
| **G18** | Baixa | Seção 4A | Fechada: relato SBAR do Andon diferencia recusa pré-execução (ex.: sem `--allow-scripts`) de execução com exit code 0. | `internal/runner/sbar.go` | `internal/runner/sbar_test.go` (`TestSBARDeRecusaPreExecucao`) |
| **G19** | Baixa | Seção 4A | Fechada: arquivo temporário `scratch/i6_inventario_legados.json` removido do rastreamento Git (`git ls-files scratch` vazio). | `scratch/i6_inventario_legados.json` (removido) | `git ls-files scratch` (saída vazia) |

### Infraestrutura Estrutural Entregue Além das Lacunas

* **Registro da Espinha e Gate (`E05`):** `espinha/registro.json`, `internal/espinha/verificar.go`, `internal/espinha/duplicadas.go`, `internal/espinha/regras_test.go`, `internal/espinha/mutantes_test.go` e integração no `fechar` (`internal/fechar/fechar.go`). Verificado nesta rodada: `go run ./cmd/legoloop espinha verificar` retornou exit 0.
* **Determinismo de Execução (`E01`):** `internal/cli/run_determinismo_test.go` e `internal/cli/exec_determinismo_test.go` provam relatório e hashes idênticos em execuções repetidas.
* **Ferramenta de Mutação (`E08`):** `tools/mutar` construída apenas com biblioteca padrão (`go.mod` e `go.sum` intactos), com 7 operadores, timeout por mutante, `GOMEMLIMIT`, encerramento de árvore de processos e suporte a build tags.

---

## 3. Mutação dos Sete Pacotes da Espinha (`TICK-I4-08.4`)

Todos os sete pacotes da espinha dorsal foram submetidos ao `tools/mutar` individualmente com orçamento controlado e `MUTAR_WORKERS=1`, atingindo **0 sobreviventes reais** e **0 timeouts** nas rodadas individuais. Os mutantes genuinamente equivalentes e guardas defensivas de I/O inatingíveis foram consolidados em `espinha/mutantes_justificados.json` (40 entradas catalogadas com anticorpos `ANT-MUT-*`).

| Pacote | Arquivo de Testes de Mutação | Sobreviventes Reais | Entradas em `espinha/mutantes_justificados.json` | Observação / Prova |
| :--- | :--- | :---: | :---: | :--- |
| `internal/upstream` | `internal/upstream/upstream_test.go` | 0 | 1 (`ANT-MUT-UPSTREAM-001`) | Mutante equivalente em `upstream.go:42` (`lt_leq`). |
| `internal/interview` | `internal/interview/mutacao_test.go` | 0 | 6 (`ANT-MUT-INTERVIEW-001` a `006`) | 3 equivalentes (`engine.go`, `intent_texto.go`) e 3 guardas defensivas de I/O (`client.go`, `sessions.go`). |
| `internal/triagem` | `internal/triagem/mutacao_test.go` | 0 | 7 (`ANT-MUT-TRIAGEM-001` a `007`) | Equivalentes e retry defensivo em `lock_windows.go`, `memoria.go`, `pidvivo_windows.go`. |
| `internal/plangen` | `internal/plangen/mutacao_test.go` | 0 | 7 (`ANT-MUT-PLANGEN-001` a `007`) | Equivalentes em `concurrency.go`, `meta.go`, `plangen.go`, `verify.go`. |
| `internal/dag` | `internal/dag/mutacao_test.go` | 0 | 1 (`ANT-MUT-DAG-001`) | `total=245 mortos=244 sobreviventes=0 justificados=1 timeouts=0`. |
| `internal/gate` | `internal/gate/mutacao_test.go` | 0 | 2 (`ANT-MUT-GATE-001` a `002`) | `total=259 mortos=257 sobreviventes=0 justificados=2 timeouts=0`. |
| `internal/runner` | `internal/runner/mutacao_test.go` | 0 | 16 (`ANT-MUT-RUNNER-001` a `016`) | `total=379 mortos=362 sobreviventes=0 justificados=17 timeouts=0` (`runner.go:497:and_or` cobre os 2 operadores `||` da linha). |
| **Total** | **7 pacotes** | **0** | **40 entradas** | **100% dos mutantes matáveis eliminados nas rodadas individuais.** |

> Nota sobre `internal/espinha`: as próprias regras do verificador da espinha têm prova dedicada de 100% de morte de mutantes em `internal/espinha/mutantes_test.go` (`TICK-I4-08.5`, aprovado).

---

## 4. Conferência do Congelamento das Iniciativas I01, I02 e I03 (`TICK-I4-09.3`)

Conforme a Emenda 2 e o critério `AC-E09.3`, nenhum aviso de congelamento foi removido antes do Stop-Gate final humano. Conferência direta realizada em 2026-11-22:

| Prompt | Ocorrências de `EXECUÇÃO CONGELADA` | Situação |
| :--- | :---: | :--- |
| `docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md` | **1** | Congelamento mantido |
| `docs/prompts/2026-11-21_PROMPT_EXECUCAO_I02.md` | **1** | Congelamento mantido |
| `docs/prompts/2026-11-21_PROMPT_EXECUCAO_I03.md` | **1** | Congelamento mantido |

---

## 5. Não verificado

Em cumprimento à regra de verdade operacional e ao padrão STE-PT, os itens abaixo **não foram verificados** nesta execução e permanecem explicitamente declarados:

1. **CI verde em Linux e macOS após o push final da I04:** `.github/workflows/ci.yml` está criado com matriz `ubuntu-latest`, `windows-latest` e `macos-latest`, mas a última execução remota registrada antes do fechamento final concluiu em falha (`progresso.md`, `TICK-I4-07.8`); a validação verde nos três SOs depende do push final da I04.
2. **Execução única combinada dos 7 pacotes juntos no `tools/mutar`:** os sete pacotes da espinha foram provados individualmente com 0 sobreviventes reais, mas não houve rodada única agregada invocando todos os sete pacotes em um único comando `go run ./tools/mutar ./internal/dag ./internal/runner ./internal/gate ./internal/plangen ./internal/triagem ./internal/upstream ./internal/interview`.
3. **Execução nativa em macOS:** não há máquina macOS local; o juiz será o job `macos-latest` do CI após o push.
4. **Unidades de rede e mídias removíveis na descoberta:** os testes de `internal/descoberta` cobrem unidades fixas por SO com provedor injetável, sem teste em volumes SMB/NFS ou discos removíveis reais.
5. **Desempenho e consumo de memória sob DAG de grande porte:** não foi executado benchmark de estresse com DAGs de centenas de nós.
6. **Reconciliação final de `progresso.md` para `TICK-I4-01.1`, `TICK-I4-07.7` e Épico 09:** no snapshot de `progresso.md` de `2026-11-22T14:48:15Z`, `TICK-I4-01.1` e `TICK-I4-07.7` apareciam como `FAILED` por exigirem confirmação manual (`grep -c`), e os tickets de `E09` seguiam `PENDING` aguardando o fechamento consolidado.

---

## 6. Próximos Passos e Stop-Gate Final

1. Concluir os tickets paralelos do Épico 09 (`TICK-I4-09.1`, `TICK-I4-09.2`, `TICK-I4-09.5`), rodar o portão global completo (`go test -count=1 ./...`, `go vet ./...`, `funclen`, `gofmt`, `taxlinter`, `espinha verificar`, `cobertura`) e reconciliar `progresso.md`.
2. Apresentar este relatório e as evidências consolidadas ao usuário no **Stop-Gate Final da I04**.
3. Somente após aprovação explícita do usuário: remover `EXECUÇÃO CONGELADA` dos prompts de I01, I02 e I03, efetuar commit/push e confirmar o CI verde nos três SOs.
