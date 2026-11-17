---
name: legoloop-memoria-anticorpos
description: "Use quando o usuário pedir controle de progresso de sessão, memória de execução, anticorpos de erro, gate anti-reincidência ou o comando legoloop sessoes, e no fechamento da demanda MEMORIA_CONCORRENTE_ANTICORPOS."
consumes: Artifact:ProcessoDir
produces: Artifact:MemoriaExecucao
artifact: memoria_execucao.json
gate: exit_code == 0
suggested_model: sonnet
---

# Skill: legoloop-memoria-anticorpos

> **Nome canônico:** `legoloop-memoria-anticorpos`
> **Use quando:** o usuário pedir controle de progresso de sessão, memória de execução, anticorpos de erro, gate anti-reincidência ou o comando `legoloop sessoes`, e no fechamento da demanda MEMORIA_CONCORRENTE_ANTICORPOS.

## Propósito

Opera a memória de execução e o banco de anticorpos do LegoLoop, conforme a
Lei 11 (Governança Viva) do AGENTS.md. A memória de execução dá controle de
progresso por sessão. O banco de anticorpos evita reincidência de erro.

## Peças envolvidas

| Peça | Pacote | Contrato |
| :--- | :--- | :--- |
| Memória de execução | `internal/triagem` | `SalvarMemoria`, `LerMemoria`, `ListarSessoesAtivas`, `ZerarSessao`, `LimparOrfaos` |
| Anticorpos | `internal/brain` | `RecordComEscopo`, `GateAntiReincidencia`, `InjectPreventiveContext` |
| Quarteto Lei 11 | `internal/fechar` | Zera a sessão no Fechamento Canônico |
| Memória em runtime | `internal/runner` | Grava `etapa_atual` e `proximo_passo` a cada nó do DAG |
| CLI de sessões | `internal/cli` | `legoloop sessoes <slug> [--json] [--processos-dir DIR]` |

## Etapas

1. Confira as sessões ativas com `legoloop sessoes <slug> --json`.
2. A saída JSON traz `sessao_id`, `pid`, `etapa_atual` e `proximo_passo`.
3. A saída nunca expõe caminho absoluto (Lei 8).
4. Processo sem sessão viva devolve `[]` e exit 0.
5. Para erro conhecido, use `brain.GateAntiReincidencia` antes de executar o nó.
6. O gate devolve `AcaoNormal`, `AcaoAntidoto` ou `AcaoBloqueado`.
7. O bloqueio devolve `ErrGateAntiReincidencia` com a causa raiz na saída.

## Schemas

- `memoria_execucao.v1`: envelope com `schema_version`, `processo` e `sessoes`.
- `anticorpos.v1`: banco com `globais` e `dominio`; `assinatura` é SHA-256 do
  erro no formato `sha256:<hex>`.

## Recuperação de falha

- `ErrLockOcupado`: outro Processo detém o lock. Aguarde ou confira o `pid`.
- Lock órfão: PID morto; o runtime toma o lock sozinho.
- Arquivo corrompido: o runtime cria `<arquivo>.legoloop-bak` e segue.
- Anticorpo corrompido: backup criado; banco vazio funcional é liberado.

## Portões de qualidade

- `go test -count=1 ./...`
- `go vet ./...`
- `gofmt -l .` com saída vazia
- `go run ./tools/funclen internal` sem estouro novo
- `go run ./tools/taxlinter internal cmd skills docs/templates` sem violação nova
