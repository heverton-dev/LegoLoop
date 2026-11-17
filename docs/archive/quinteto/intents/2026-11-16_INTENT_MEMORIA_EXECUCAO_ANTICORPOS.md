# INTENT — Memória de Execução Concorrente + Anticorpos Persistidos

- **Data:** 2026-11-16
- **Origem:** Sessão `8e1ac50d-4232-4ce2-883a-c4bd3bbc222c` (enxame Fabricar Agnóstico) — lacunas encontradas em auditoria final do verbo `fabricar`.
- **Status:** REGISTRADA — aguarda SPEC/PLAN/ROUTE/DAG em sessão estéril dedicada.

## Problema (2 lacunas)

### Lacuna 1 — Gestão concorrente da memória de execução

O arquivo `processos/<slug>/memoria_execucao.json` é o checkpoint vivo do processo (etapa atual, próximo passo, zerado no fechar). Hoje é **arquivo único por processo**, sem chave de sessão e sem lock (`internal/triagem/memoria.go:11-22`). Se 5 agentes executam a mesma etapa simultaneamente, todos gravam o mesmo arquivo: corrida e corrupção de estado.

**Desejo:** gestão de sessões concorrentes com lock e estado por sessão; zera no fechar sem inflar disco.

### Lacuna 2 — Anticorpos só com casca

O Motor III de anticorpos existe (`internal/brain/antibodies.go:10-67`) com tripla erro→causa→diff e injeção preventiva, mas `AntibodiesDB` é mapa **em memória**: sem persistência, sem arquivo por processo, sem vacinação automática de reincidência.

**Desejo (padrão Tesla/SpaceX):** arquivo de anticorpos persistido por processo; a cada execução, erro com reincidência é eliminado por antídoto registrado — quanto mais o processo roda, mais blindado ele fica.

## Fora de escopo

- Não altera os 6 verbos além do necessário.
- Não toca no cofre fora de verbos oficiais.

## Referências

- Prova Lacuna 1: `internal/triagem/memoria.go:11-22`; zera no fechar: `internal/fechar/quarteto.go:50`.
- Prova Lacuna 2: `internal/brain/antibodies.go:22` (mapa em memória); injeção preventiva: `antibodies.go:63-67`.
- Relatório da sessão: `docs/reports/2026-11-16_REPORT_FABRICAR_AGNOSTICO.md`.

---

## PROMPT DE ABERTURA DA PRÓXIMA SESSÃO (copiar e colar)

> Leia a INTENT `docs/intents/2026-11-16_INTENT_MEMORIA_EXECUCAO_ANTICORPOS.md` e execute a demanda MEMÓRIA-CONCORRENTE+ANTICORPOS do fim ao fim, sem me interromper exceto nos stop-gates.
> Regras obrigatórias: PT-BR em tudo; TDD estrito (RED antes do GREEN, mutação prova o teste, zero stubs); funções ≤35 linhas; só stdlib; hermético (`t.TempDir()`, nunca tocar `~/.aidd` fora de verbos oficiais); nunca editar `.dag.json`/`.processo.json`/`.loop.json` à mão — só verbos oficiais; nunca `taskkill /IM legoloop.exe`; commit e push SOMENTE se eu pedir; declare sempre o que não foi provado.
> Sequência executiva:
> 1. Fechar o Quinteto: gerar SPEC (`docs/specs/2026-11-16_SPEC_MEMORIA_CONCORRENTE_ANTICORPOS.md`) com critérios EARS binários, PLAN com tickets, ROUTE e DAG espelho. Me mostrar em SG-M01.
> 2. Executar os tickets M1.1→M3.3 (tabela abaixo) em enxame SWARM com workers estéreis (1 worker por ticket, contexto só do ticket, saída em disco, gate mecânico por ticket).
> 3. Ao fim de M1 e M2, apresentar SG-M02 e SG-M03 com provas arquivo:linha.
> 4. Fechamento canônico (Lei 11) em SG-FINAL-M: manual em docs/manuals/, livro LegoPress, skill canônica, DAG validado, `go build` + `go install`, AGENTS.md/MEMORY.md/README, reindexação do grafo, relatório em docs/reports/ com telemetria de tempo, tokens e custo USD/BRL orçado vs realizado.
> Gate de qualidade por ticket: `go test -count=1 ./...`, `go vet ./...`, funclen ≤35, `gofmt -l .` vazio, `git diff --exit-code -- go.mod go.sum`, taxlinter sem violação nova.

## DECOMPOSIÇÃO COMPLETA DE TRABALHO (12 unidades + 4 stop-gates)

| ID | Épico | Unidade de trabalho | Prova de pronto |
| :--- | :--- | :--- | :--- |
| M0.1 | Planejamento | SPEC EARS + PLAN + ROUTE + DAG do Quinteto | SG-M01 humano |
| M1.1 | Memória concorrente | Modelo `MemoriaExecucao` ganha chave de sessão (ID de sessão + PID) | teste RED→GREEN + mutação |
| M1.2 | Memória concorrente | Lock de arquivo com PID (padrão `internal/board`) em `SalvarMemoria`/`LerMemoria` | teste de corrida (2 sessões simultâneas sem corrupção) + mutação |
| M1.3 | Memória concorrente | `runner` grava por sessão; `fechar` zera só a sessão dona e limpa arquivo órfão | teste ciclo completo + mutação |
| M1.4 | Memória concorrente | CLI/Studio expõe sessões ativas do processo | contrato de CLI testado |
| M2.1 | Anticorpos | `AntibodiesDB` persiste em `processos/<slug>/anticorpos.json` (carrega ao iniciar, grava ao registrar) | teste de persistência entre execuções + mutação |
| M2.2 | Anticorpos | Vacinação de reincidência: `Record` deduplica por `errSig`; reincidência aplica antídoto automaticamente | teste: mesmo erro 2× ⇒ antídoto aplicado, erro não repete + mutação |
| M2.3 | Anticorpos | Injeção preventiva no início da execução (`InjectPreventiveContext` vira passo do processo) | teste: contexto injetado antes do nó 1 + mutação |
| M2.4 | Anticorpos | Gate anti-reincidência em runtime (erro conhecido ⇒ antídoto ou bloqueio com saída clara) | teste de gate + mutação |
| M3.1 | Fechamento | Manual operacional em `docs/manuals/` | arquivo + checklist Lei 11 |
| M3.2 | Fechamento | Livro LegoPress + skill canônica + DAG validado (`legoloop validate`) | validação exit 0 |
| M3.3 | Fechamento | Build + install, governança viva, reindex do grafo, relatório com telemetria | `go install` EXIT=0; relatório em `docs/reports/` |

## ORÇAMENTO E ALOCAÇÃO DE EXECUÇÃO

- **Execução:** enxame SWARM, 1 worker estéril por unidade de trabalho (até 4 em paralelo nas cadeias independentes M1∥M2); coordenador não executa trabalho técnico.
- **Estimativa:** 12 unidades + 4 stop-gates; ~4-6h de sessão; orçamento de tokens e custo USD/BRL deve ser registrado orçado vs realizado no relatório final (TEMPLATE_REPORT v2.2.0).
- **Paradas humanas obrigatórias:** SG-M01 (pós-planejamento), SG-M02 (pós-M1), SG-M03 (pós-M2), SG-FINAL-M (pós-fechamento). Nenhum stop-gate é auto-aprovado.
