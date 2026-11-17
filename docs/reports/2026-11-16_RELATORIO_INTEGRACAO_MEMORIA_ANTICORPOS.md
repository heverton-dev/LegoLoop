# Relatório de Integração — Memória Concorrente e Anticorpos (TICK-INT-01 + TICK-INT-02)

- **Data:** 2026-11-16
- **Worker:** `e03-integracao`
- **Demanda paralela:** MEMORIA_CONCORRENTE_ANTICORPOS (PLAN/SPEC/ROUTE/DAG de 2026-11-16)
- **Mega Plano:** Épico 03 (memória de execução), Épico 04 (anticorpos), Épico 11 (Lei 11)
- **Commits desta integração:** `98a132e`, `57fe915`, `4b7a512`, `9c0a705` + commit de fechamento

---

## 1. Divergências anotadas

| Tema | Brief / SPEC | Realidade no repositório | Decisão |
| :--- | :--- | :--- | :--- |
| Contagem de tickets | O brief cita 37 tickets paralelos | O PLAN paralelo tem 11 tickets (M1.1–M3.3) e a ROUTE tem 15 nós (11 tickets + 4 stop-gates) | Mapeei os 11 tickets, os 15 nós e os 27 ACs (19 paralelos + 8 mega); nenhuma linha ficou sem status |
| Nome de erro de lock | Mega SPEC cita `ErrLockAtivo` | O código usa `ErrLockOcupado` (`internal/triagem/memoria.go:33`) | Comportamento idêntico; divergência de nome registrada, sem renomear API |
| Pasta de Processos | Schemas citam `processos/<slug>/` | O acervo real vive em `processos/<slug>/` (26 Processos) | O conteúdo segue os schemas `anticorpos.v1`/`memoria_execucao.v1`; a pasta segue o dialeto do acervo |

## 2. TICK-INT-01 — Tabela A: tickets paralelos × ACs × status

| Ticket | Descrição | ACs paralelos | ACs mega | Status |
| :--- | :--- | :--- | :--- | :--- |
| M1.1 | Modelo de sessão serializa e valida | AC-E01.1, AC-E01.2 | — | COBERTO por 3153922 |
| M1.2 | Concorrência com lock por arquivo | AC-E01.3, AC-E01.4 | mega AC-E03.1 | COBERTO por 3153922 (ressalva de nome na divergência 2) |
| M1.3 | Fechar zera a sessão dona e limpa órfãos | AC-E01.5, AC-E01.6 | mega AC-E03.3 | COBERTO por 3153922; reforçado por TICK-INT-02 (`4b7a512`, `TestZeraMemoria` + mutação 14) |
| M1.4 | CLI/Studio expõe sessões ativas | AC-E01.7, AC-E01.8 | — | COBERTO por TICK-INT-02 (`4b7a512`, comando `legoloop sessoes`) |
| M2.1 | Anticorpos persistem entre execuções | AC-E02.1, AC-E02.2 | mega AC-E04.1, mega AC-E04.4 | COBERTO por TICK-INT-02 (`57fe915`, schema `anticorpos.v1` + backup) |
| M2.2 | Reincidência vacina com antídoto | AC-E02.3, AC-E02.4 | mega AC-E04.2 | COBERTO por 3153922 + TICK-INT-02 (`57fe915`, assinatura SHA-256) |
| M2.3 | Contexto preventivo injetado no nó | AC-E02.5, AC-E02.6 | — | COBERTO por 3153922 (`internal/orchestrix/orchestrix.go`) |
| M2.4 | Gate anti-reincidência em runtime | AC-E02.7, AC-E02.8 | mega AC-E04.3 | COBERTO por TICK-INT-02 (`57fe915`, `GateAntiReincidencia` + escopo global/domínio) |
| M3.1 | Manual operacional em `docs/manuals/` | AC-E03.1 (paralelo) | — | COBERTO por TICK-INT-02 (`4b7a512`) |
| M3.2 | Livro LegoPress + Skill + DAG validado | AC-E03.2 (paralelo) | — | COBERTO por TICK-INT-02 (`9c0a705`) |
| M3.3 | Build + install, governança viva, relatório | AC-E03.3 (paralelo) | — | COBERTO por TICK-INT-02 (commit de fechamento) |

## 3. TICK-INT-01 — Tabela B: 19 ACs paralelos × status

| AC | Descrição | Status | Prova |
| :--- | :--- | :--- | :--- |
| AC-E01.1 | Modelo de sessão serializa e valida | COBERTO por 3153922 | `internal/triagem/memoria.go` (`EstadoMemoria`, `Validate`) |
| AC-E01.2 | Sessão sem identidade é rejeitada | COBERTO por 3153922 | `internal/triagem/memoria.go:30` (`ErrSessaoInvalida`) |
| AC-E01.3 | Duas sessões simultâneas gravam sem corrupção | COBERTO por 3153922 | `internal/triagem/memoria.go:104` (`adquirirLockMemoria`) |
| AC-E01.4 | Lock retido não corrompe estado | COBERTO por 3153922 | `internal/triagem/memoria.go` (`ErrLockOcupado`) |
| AC-E01.5 | Fechar zera só a sessão dona e limpa órfãos | COBERTO por 3153922 | `internal/fechar/quarteto.go:54` + `triagem.LimparOrfaos` |
| AC-E01.6 | Fechar com sessão inexistente falha claro | COBERTO por 3153922 | `internal/fechar/quarteto_test.go` (`TestFecharSessaoDona`) |
| AC-E01.7 | CLI expõe sessões ativas | COBERTO por TICK-INT-02 (`4b7a512`) | `internal/cli/sessoes.go` (`cmdSessoes`), `internal/cli/sessoes_test.go` (`TestCmdSessoesAtivasListaDuas`) |
| AC-E01.8 | Processo sem sessões devolve lista vazia | COBERTO por TICK-INT-02 (`4b7a512`) | `internal/cli/sessoes_test.go` (`TestCmdSessoesVazio`: `[]`, exit 0) |
| AC-E02.1 | Anticorpos persistem entre execuções | COBERTO por 3153922 + `57fe915` | `internal/brain/antibodies.go` (banco em disco `anticorpos.v1`) |
| AC-E02.2 | Arquivo corrompido falha com erro sentinela | COBERTO por TICK-INT-02 (`57fe915`) | `internal/brain/antibodies.go` (`ErrAnticorposCorrompidos` + backup `.legoloop-bak`) |
| AC-E02.3 | Reincidência vacina com antídoto | COBERTO por `57fe915` | `internal/brain/antibodies.go` (`AcaoAntidoto`, `AntidotoAplicado`) |
| AC-E02.4 | Registro sem assinatura é rejeitado | COBERTO por `57fe915` | `internal/brain/antibodies.go` (`assinatura` SHA-256 `sha256:<hex>`), mutação 8 |
| AC-E02.5 | Contexto preventivo injeta antes do nó | COBERTO por 3153922 | `internal/orchestrix/orchestrix.go` (`InjectPreventiveContext`) |
| AC-E02.6 | Sem anticorpos a execução segue normal | COBERTO por 3153922 | `internal/brain/antibodies_test.go` (injeção preventiva) |
| AC-E02.7 | Gate bloqueia erro conhecido em runtime | COBERTO por TICK-INT-02 (`57fe915`) | `internal/brain/antibodies.go` (`GateAntiReincidencia`, `ErrGateAntiReincidencia`), mutação 7 |
| AC-E02.8 | Erro desconhecido não é bloqueado | COBERTO por `57fe915` | `internal/brain/antibodies.go` (`AcaoNormal`) |
| AC-E03.1 (paralelo) | Manual operacional emitido | COBERTO por TICK-INT-02 (`4b7a512`) | `docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md` |
| AC-E03.2 (paralelo) | Livro, Skill e DAG validados | COBERTO por TICK-INT-02 (`9c0a705`) | `docs/books/memoria-anticorpos/`, `skills/legoloop-memoria-anticorpos/SKILL.md`, `legoloop validate` (15 nós válidos) |
| AC-E03.3 (paralelo) | Build, governança viva e relatório | COBERTO por TICK-INT-02 (commit de fechamento) | Seção 6 deste relatório |

## 4. TICK-INT-01 — Tabela C: 8 ACs do mega SPEC × status

| AC | Descrição | Status | Prova |
| :--- | :--- | :--- | :--- |
| mega AC-E03.1 | Lock por arquivo (segunda sessão recebe erro, arquivo íntegro) | COBERTO por 3153922 | `internal/triagem/memoria.go:33` (`ErrLockOcupado`; divergência 2) |
| mega AC-E03.2 | Lock órfão tomado e `atualizado` avança | COBERTO por TICK-INT-02 (`98a132e`) | `internal/triagem/memoria.go` (`limparLockOrfao`), mutações 1 |
| mega AC-E03.3 | Fechar zera `memoria_execucao.json` (schema válido) | COBERTO por TICK-INT-02 (`4b7a512`) | `internal/fechar/quarteto_test.go` (`TestZeraMemoria`), mutação 14 |
| mega AC-E03.4 | Corrompido: sentinela + backup `.legoloop-bak` + sem escrita cega | COBERTO por TICK-INT-02 (`98a132e`) | `internal/triagem/memoria.go` (`ErrMemoriaCorrompida`, `backupCorrompido`), mutações 2–3 |
| mega AC-E04.1 | Persistência em `anticorpos.json` (schema `anticorpos.v1`) | COBERTO por TICK-INT-02 (`57fe915`) | `internal/brain/antibodies.go` (`bancoJSON`); divergência 3 (pasta `processos/<slug>`) |
| mega AC-E04.2 | Reincidência incrementa e antídoto é injetado antes do nó | COBERTO por 3153922 + `57fe915` | `Find` + `InjectPreventiveContext` em `orchestrix.go` |
| mega AC-E04.3 | Escopo global (infra/LLM) vs domínio (regra local) | COBERTO por TICK-INT-02 (`57fe915`) | `RecordComEscopo`, `NewAntibodiesDBPersistenteComGlobal`, mutação 5 |
| mega AC-E04.4 | Corrompido: sentinela + backup; processo segue sem vacina | COBERTO por TICK-INT-02 (`57fe915`) | `ErrAnticorposCorrompidos` + banco vazio funcional, mutação 6 |

## 5. TICK-INT-01 — Tabela D: 15 nós da ROUTE paralela × status

| Nó | Tipo | Status |
| :--- | :--- | :--- |
| SG-M01 | Stop-Gate humano | PENDENTE (Stop-Gate humano; coberto pelo SG-EINT desta entrega) |
| M1.1, M1.2, M1.3 | Tickets | COBERTO por 3153922 (Tabela A) |
| M1.4 | Ticket | COBERTO por TICK-INT-02 (`4b7a512`) |
| M2.1, M2.2, M2.3 | Tickets | COBERTO por 3153922 + `57fe915` (Tabela A) |
| M2.4 | Ticket | COBERTO por TICK-INT-02 (`57fe915`) |
| SG-M02, SG-M03 | Stop-Gates humanos | PENDENTE (Stop-Gates humanos; cobertos pelo SG-EINT desta entrega) |
| M3.1, M3.2, M3.3 | Tickets Lei 11 | COBERTO por TICK-INT-02 (`4b7a512`, `9c0a705`, commit de fechamento) |
| SG-FINAL-M | Stop-Gate humano | PENDENTE (Stop-Gate humano; este relatório é a peça de decisão) |

## 6. TICK-INT-02 — Gaps implementados (TDD estrito)

| # | Gap | Commit | Mutações provadas (falha → restaurada) |
| :--- | :--- | :--- | :--- |
| 1 | Schemas `memoria_execucao.v1` como fonte da persistência; lock órfão tomado; backup `.legoloop-bak`; escrita cega bloqueada | `98a132e` | (1) sem `limparLockOrfao` → FAIL; (2) sem `SchemaVersion` → FAIL; (3) sem escrita do backup → FAIL |
| 2 | Persistência real de anticorpos via `anticorpos.v1` (template `TEMPLATE_ANTICORPOS.json`); escopo global/domínio; gate anti-reincidência; carga legado sem perda | `57fe915` | (4) sem fundir globais → FAIL; (5) escopo desviado → FAIL; (6) sem backup → FAIL; (7) antídoto vira bloqueio → FAIL; (8) assinatura sem SHA-256 → FAIL |
| 3 | CLI `legoloop sessoes <slug> [--json] [--processos-dir DIR]` (AC-E01.7/E01.8) | `4b7a512` | (9) `ListarSessoesAtivas` sem filtro de PID vivo → FAIL; (10) `cmdSessoes` sem impressão JSON → FAIL; (11) resolução de pasta nunca acha → FAIL |
| 4 | Integração runtime de `memoria_execucao.json`: `etapa_atual` e `proximo_passo` a cada nó do DAG | `4b7a512` | (12) `preExecNode` sem `registrarEtapaMemoria` → FAIL; (13) `proximosDe` sem `DependsOn` → FAIL |
| 5 | `TestZeraMemoria` (mega AC-E03.3) no Fechamento Canônico Lei 11 | `4b7a512` | (14) `quarteto.go` sem `zerarMemoriaSessao` → FAIL em `TestZeraMemoria` e `TestFecharSessaoDona` |
| 6 | Épico 03 paralelo (Lei 11): manual, livro LegoPress, Skill canônica, DAG validado, governança viva | `4b7a512`, `9c0a705`, commit de fechamento | Sem mutação (artefatos de documentação); auditoria geométrica do PDF: 16 páginas, 100% válida |

Todas as 14 mutações foram restauradas; nenhuma alteração de mutação ficou no disco.

## 7. Gates de qualidade (saída colada)

Saída literal de cada portão, rodada em 2026-11-16 no estado final da árvore:

1. **`go test -count=1 ./...`** — `TEST_EXIT=0` (cauda literal: `ok legoloop/tools/route2dag 3.637s`, `ok legoloop/tools/taxlinter 3.881s`, `ok legoloop/tools/ticketgate 5.072s`; nenhum FAIL em nenhum pacote).
2. **`go vet ./...`** — `VET_EXIT=0`, saída vazia.
3. **`gofmt -l .`** — `GOFMT_LIST=[]` (lista vazia, nenhum arquivo fora do formato).
4. **`go run ./tools/funclen internal`** — `FUNCLEN_EXIT=1` apenas pelos 3 estouros **pré-existentes** (baseline herdado): `internal\board\swarm_run.go:90 runLevel 45 linhas`, `internal\cli\swarm_run.go:13 cmdSwarmRun 42 linhas`, `internal\gate\cites.go:129 CiteVerdicts 47 linhas`. **Zero estouro novo.**
5. **`go run ./tools/taxlinter internal cmd loops skills docs/templates`** — `TAX_EXIT=1`, `VIOLACOES=31`, **todas pré-existentes** (arquivos herdados: `swarm_run_e2e_test.go`, `sync_yaml.go`, golden, `processos/legoloop-swarm`, etc.). **Zero violação nova** (correção final: `AcaoFluxoNormal` → `AcaoNormal` também no SKILL.md, linha 39).
6. **`git diff --exit-code -- go.mod go.sum`** — `GOMOD_EXIT=0`, saída vazia (somente stdlib).

## 8. NAO_PROVADO

- `-race` (exige CGO neste ambiente Windows).
- Fluxo completo com LLM real (testes usam LLM falso hermeticamente).
- Reindexação do grafo semântico (`codebase-memory-mcp` ausente no PATH).
- Studio aberto continua no binário antigo; reiniciar após `go install`.
- `go install ./cmd/legoloop` pode falhar se o Studio estiver com o binário em uso; nesse caso, reiniciar o Studio resolve.
- Stop-gates humanos SG-M01/SG-M02/SG-M03/SG-FINAL-M do plano paralelo seguem abertos por natureza (Lei 8): a decisão humana é o SG-EINT.

## 9. Stop-Gate SG-EINT

`[1] seguir | [2] personalizar | [0] cancelar`
