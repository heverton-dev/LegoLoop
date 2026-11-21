# 📊 REPORT: Provas Ausentes & Débitos do Épico 10 (TICK-10.3 · AC-E10.3)

> **Data:** 2026-11-17
> **Épico:** E10 — I9 · Provas Ausentes & Higiene
> **Ticket:** TICK-10.3 (verificação de provas ausentes e registro de NAO_PROVADO com dono e prazo)
> **Executor:** SWARM `e10-higiene`

---

## 1. Resultado

As três provas ausentes permanecem **NAO_PROVADO** nesta máquina, com dono e prazo registrados abaixo.
Nenhuma delas foi executada aqui; nada foi arredondado para "pronto".

---

## 2. Veredicto por prova ausente

| Prova | Veredicto | Evidência mecânica | Dono | Prazo |
| :--- | :--- | :--- | :--- | :--- |
| `-race` com CGO | **NAO_PROVADO** | `go env CGO_ENABLED` → `0`; `gcc --version` → comando não encontrado (sem toolchain C) | Coordinator (`@coordinator`) | Antes do SG-FINAL, no CI Linux com CGO (o `-race` é obrigatório lá) |
| LLM real ponta a ponta | **NAO_PROVADO** | Suíte hermética usa LLM falso (AGENTS.md, regra 8: "testes e smoke usam LLM falso"); nenhum gateway real foi chamado | Usuário (homologação humana) | Homologação manual via gateway 9Router após SG-FINAL |
| Grafo semântico (`codebase-memory-mcp`) | **NAO_PROVADO** | MCP indisponível neste harness; Lei 14 (Graph-First) cumpriu-se por busca direta nos arquivos | Coordinator (`@coordinator`) | `codebase-memory-mcp cli index_repository` na governança viva do SG-FINAL |

---

## 3. O que ficou PROVADO no Épico 10 (mutações documentadas)

| # | Mutação | Teste que mordeu | Resultado |
| :--- | :--- | :--- | :--- |
| M1 | Mega SPEC volta a citar `ErrLockAtivo` (débito 2, Lei 15) | `internal/triagem/lei15_lock_test.go:TestLei15UmTermoParaLock` | FAIL → restaurado → PASS (Green-Red-Green) |
| M2 | Checagem de escape removida em `paths.DentroDaRaiz` (`rel == ".."` desativada) | `internal/paths/jail_test.go:TestEntradaHostil` | FAIL (`RenameHostil` criou arquivo fora da raiz) → restaurado → PASS |
| M3 | Lock de arquivo vira lock global (`adquirirLockMemoria` com arquivo único de temp) | `internal/triagem/carga_test.go:TestCinquentaWorkers` | FAIL (barreira dos 50 locks nunca fecha) → restaurado → PASS |

Provas anti-flake (TICK-10.2 + débito 3): `go test -count=3` verde em `internal/runner`, `internal/board`, `internal/mcp`, `internal/studio` após as correções de raiz (barreira de sobreposição determinística no `fakeExec`; espera ativa por condição binária; kill assíncrono do Windows aguardado em polling).

---

## 4. DÍVIDA REGISTRADA — exceções do taxlinter (débito 1, condição SG-E05)

O mapa `tools/taxlinter/exceptions.go` cobre contratos externos invioláveis (nomes de artefato e ferramentas de terceiros).
Contagem auditada por git:

| Marco | Entradas (arquivo) | Pares (arquivo→termo) |
| :--- | :---: | :---: |
| Baseline anterior ao `0e9bb5c` (E05) | 88 | 91 |
| Após `0e9bb5c` (migração do Quinteto) | 90 | 139 (+48) |
| Estado atual (após `886f783`, rename E06) | 94 | 143 (+4) |

- **Condição (a) — listagem:** a relação completa das exceções vigentes está em [seção 4.1](#41-relação-completa-das-exceções-vigentes).
- **Condição (b) — não crescimento:** `48` (pares migrados do Quinteto em `0e9bb5c`) `+ 91` (baseline anterior) `= 139` homologados. Os `+4` pares atuais (`886f783`, Épico 06) são os MESMOS termos sob os caminhos renomeados para `processos/` no Épico 06 — renome herdado homologado no SG-E06, sem vocabulário novo. **Nenhum termo novo entrou no mapa após o homologado.**

### 4.1 Relação completa das exceções vigentes

*(gerada mecanicamente de `tools/taxlinter/exceptions.go`)*

Total: 94 entradas.

- `internal/cli/rota.go` — task
- `internal/roteamento/classe.go` — macro
- `internal/orchestrix/unit_claude.go` — task
- `internal/orchestrix/orchestrix.go` — task, loop
- `internal/orchestrix/subloop.go` — loop
- `internal/orchestrix/subloop_test.go` — loop
- `internal/orchestrix/subloop_artefato_test.go` — loop
- `internal/orchestrix/unit.go` — task
- `internal/orchestrix/prompt.go` — tarefa
- `internal/orchestrix/unit_run.go` — task
- `internal/orchestrix/unit_run_test.go` — task
- `internal/orchestrix/unit_claude_test.go` — macro
- `internal/orchestrix/unit_test.go` — task
- `internal/legopress/auditor_test.go` — pipeline
- `internal/golden/testdata/processos/legoloop-video-motion.dag.json` — pipeline
- `internal/golden/testdata/processos/legoloop-mega-pipeline.loop.json` — pipeline
- `internal/interview/intent_render.go` — task
- `docs/templates/TEMPLATE_DAG.json` — pipeline
- `skills/legoloop-loop-map/SKILL.md` — loop
- `skills/legoloop-mega-pipeline/SKILL.md` — pipeline
- `internal/caminhos/nomes.go` — pipeline
- `internal/caminhos/nomes_test.go` — loop
- `internal/cli/commands.go` — task, loop
- `internal/studio/swarm_api.go` — card
- `internal/studio/ui_logic_test.go` — loop
- `internal/studio/static_test.go` — loop
- `internal/studio/loops_api.go` — loop
- `internal/studio/loops_api_test.go` — loop
- `internal/studio/processos_test.go` — loop
- `internal/studio/processos_rest_test.go` — loop
- `internal/studio/pieces.go` — loop
- `internal/studio/pieces_test.go` — loop
- `internal/studio/materialize_test.go` — loop
- `internal/studio/fabricar_plot_test.go` — loop
- `internal/studio/triagem_api_test.go` — loop
- `internal/studio/run_api.go` — loop
- `internal/studio/server_test.go` — tarefa
- `internal/assets/loops.go` — loop
- `internal/assets/skills.go` — pipeline
- `internal/assets/types.go` — loop
- `internal/cli/cartao_test.go` — fluxo
- `internal/cli/compor_binary_test.go` — loop
- `internal/cli/compor_real_binary_test.go` — tarefa
- `internal/cli/fabricar_report_texto_test.go` — loop
- `internal/cli/fabricar_test.go` — pipeline
- `internal/cli/listar_test.go` — loop
- `internal/cli/processo_real_binary_test.go` — pipeline
- `internal/cli/remover_test.go` — loop
- `internal/cli/upgrade.go` — loop
- `internal/cli/upgrade_hint_test.go` — fluxo
- `internal/compor/compor.go` — loop
- `internal/compor/compor_test.go` — loop
- `internal/dag/dialetos_test.go` — loop
- `internal/dag/parse.go` — task
- `internal/dag/parse_dual_test.go` — loop
- `internal/dag/types.go` — loop, task
- `internal/dag/v2_test.go` — loop
- `internal/fabricar/camada4.go` — loop
- `internal/fabricar/curadoria_fronteira.go` — pipeline
- `internal/fabricar/pieces.go` — loop
- `internal/fabricar/prompts_test.go` — tarefa
- `internal/fechar/fechar.go` — loop
- `internal/fechar/fechar_test.go` — loop
- `internal/matriz/registro.go` — loop
- `internal/matriz/vocab.default.json` — pipeline
- `internal/migrar/migrar.go` — pipeline
- `internal/migrar/migrar_test.go` — pipeline
- `internal/orchestrix/hygiene_test.go` — loop
- `internal/paths/paths.go` — loop
- `internal/processo/cofre_test.go` — loop
- `internal/processo/origem_test.go` — loop
- `internal/processo/persistencia.go` — loop
- `internal/processo/projecoes.go` — loop
- `internal/processo/projecoes_test.go` — loop
- `internal/processo/store_test.go` — loop
- `internal/refinar/refinar.go` — loop
- `internal/refinar/refinar_test.go` — loop
- `internal/remover/remover.go` — loop
- `internal/runner/n10_test.go` — pipeline
- `internal/runner/telemetry.go` — pipeline
- `internal/runner/triggers.go` — loop
- `internal/ste/glossario.json` — componente
- `internal/ste/prompt_prova.go` — bloco
- `internal/ste/ste_test.go` — esteira
- `internal/studio/assets_api.go` — loop
- `internal/studio/compose.go` — loop
- `internal/studio/dissecar_handler.go` — pipeline
- `internal/studio/dissecar_handler_test.go` — pipeline
- `internal/caminhos/caminhos_test.go` — loop
- `internal/cli/legado.go` — loop
- `internal/cli/legado_test.go` — loop, pipeline
- `internal/fechar/quarteto_test.go` — loop
- `docs/intents/2026-11-15_INTENT_VIDEO_REA_TRANSCREVER.md` — task
- `docs/routes/2026-11-15_ROUTE_VIDEO_REA_TRANSCREVER.json` — loop, pipeline

---

## 5. Débito 4 — contagem do plano paralelo ("37 tickets")

Confirmado: o relatório de integração já registra a divergência e nada pendura.
`docs/reports/2026-11-16_RELATORIO_INTEGRACAO_MEMORIA_ANTICORPOS.md:15`:

> | Contagem de tickets | O brief cita 37 tickets paralelos | O PLAN paralelo tem 11 tickets (M1.1–M3.3) e a ROUTE tem 15 nós (11 tickets + 4 stop-gates) | Mapeei os 11 tickets, os 15 nós e os 27 ACs (19 paralelos + 8 mega); nenhuma linha ficou sem status |

A Tabela A do mesmo relatório declara status por ticket (M1.1–M3.3, todos COBERTOS com commit).
Não há ticket pendurado nem sobreposição a executar duas vezes.

---

## 6. NAO_PROVADO (declaração honesta)

1. **`-race`**: não executado nesta máquina (sem CGO, sem gcc). Obrigatório no CI Linux.
2. **LLM real ponta a ponta**: não executado (hermeticidade da suíte).
3. **Grafo semântico**: não usado (MCP ausente no harness).
4. Os números de desempenho do AC-NF.2 (50 workers) provam não-serialização e integridade; não provam latência absoluta sob carga de produção.
