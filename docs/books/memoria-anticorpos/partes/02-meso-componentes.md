# Parte II — Visão Meso: os Componentes

## 1. Memória de execução (`internal/triagem`)

A memória guarda um envelope por Processo. O envelope tem `schema_version`,
`processo` e o mapa `sessoes`. Cada sessão tem `sessao_id`, `pid`,
`etapa_atual`, `proximo_passo`, `atualizado_em` e a telemetria de
`token_budget`.

Operação essencial:

| Função | Papel |
| :--- | :--- |
| `SalvarMemoria(dir, mem)` | Grava a sessão com lock por arquivo |
| `LerMemoria(dir, sessaoID)` | Lê a sessão; sem ID, lê a legada |
| `ListarSessoesAtivas(dir)` | Lista só as sessões com PID vivo |
| `ZerarSessao(dir, sessaoID)` | Zera a sessão no fechamento |
| `LimparOrfaos(dir)` | Remove sessões de PID morto |

## 2. Anticorpos (`internal/brain`)

O banco tem dois escopos. `globais` vale para todo o ecossistema. `dominio`
vale para o domínio do Processo. Cada anticorpo tem `pacote`, `causa_raiz`,
`assinatura` SHA-256 e o `antidoto_aplicado`.

| Função | Papel |
| :--- | :--- |
| `RecordComEscopo` | Grava anticorpo global ou de domínio |
| `GateAntiReincidencia` | Decide: fluxo normal, antídoto ou bloqueio |
| `InjectPreventiveContext` | Injeta o antídoto no prompt do nó |
| `NewAntibodiesDBPersistenteComGlobal` | Abre o banco com anticorpos globais |

## 3. Quarteto Lei 11 (`internal/fechar`)

O Fechamento Canônico zera a memória da sessão dona do Processo. O envelope
volta com `schema_version` igual a `memoria_execucao.v1` e sem sessões.

## 4. Memória em runtime (`internal/runner`)

A cada nó do DAG, o runner grava `etapa_atual` com o ID do nó e
`proximo_passo` com os sucessores do nó. Nó final aponta para `fim`.

## 5. CLI de sessões (`internal/cli`)

O comando `legoloop sessoes <slug> [--json] [--processos-dir DIR]` lista as
sessões ativas. A saída nunca expõe caminho absoluto.
