---
title: "Memória de Execução e Anticorpos no LegoLoop: Controle de Progresso e Gate Anti-Reincidência"
subtitle: "Como a Lei 11 registra o progresso de cada sessão e transforma erro conhecido em antídoto ou bloqueio"
author:
  - Heverton & Ecossistema LegoLoop
date: "16/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "GOVERNANÇA VIVA"
cover-title: "Memória e Anticorpos"
cover-accent: "Controle de Progresso na Lei 11"
cover-caption: "DEMANDA MEMORIA_CONCORRENTE_ANTICORPOS"
cover-stages:
  - Memória com Lock por Arquivo
  - Anticorpos no Schema anticorpos.v1
  - Gate Anti-Reincidência
  - Memória de Runtime por Nó
  - CLI de Sessões Ativas
  - Zeramento no Fechamento
header-mark: "LEGOLOOP MEMÓRIA-ANTICORPOS"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro documenta a memória de execução e o banco de anticorpos do LegoLoop, conforme a Lei 11 do AGENTS.md. A memória de execução dá controle de progresso por sessão com lock por arquivo. O banco de anticorpos persiste erros conhecidos no schema anticorpos.v1 e decide o fluxo do nó pelo gate anti-reincidência.
---


# Parte I — Visão Macro

## 1. O problema

O LegoLoop executa Processos com muitos nós e muitas sessões ao mesmo tempo.
Duas falhas custavam caro. A primeira: ninguém sabia em qual etapa cada sessão
estava. A segunda: o mesmo erro voltava a acontecer em Processos diferentes.

## 2. A resposta

A demanda MEMORIA_CONCORRENTE_ANTICORPOS ataca as duas falhas.

- A **memória de execução** grava `etapa_atual` e `proximo_passo` de cada
  sessão em `memoria_execucao.json`, com lock por arquivo.

- O **banco de anticorpos** grava cada erro conhecido com a assinatura SHA-256
  e o antídoto, no schema `anticorpos.v1`.

- O **gate anti-reincidência** decide, antes de executar o nó, entre o fluxo
  normal, o antídoto ou o bloqueio.

## 3. Os benefícios

| Benefício                     | Mecanismo                                          |
| :---------------------------- | :------------------------------------------------- |
| Controle de progresso         | `legoloop sessoes <slug> --json`                     |
| Concorrência segura           | Lock por arquivo com tomada de lock órfão          |
| Recuperação                   | Backup `.legoloop-bak` em arquivo corrompido         |
| Zero reincidência             | `brain.GateAntiReincidencia` antes do nó             |
| Prevenção                     | `brain.InjectPreventiveContext` no prompt            |
| Fechamento limpo              | Sessão zerada no Fechamento Canônico               |

## 4. A regra de ouro

Erro conhecido nunca vira improvisão. Ou o antídoto existe e é aplicado, ou o
gate bloqueia a execução e devolve a causa raiz. Essa é a Lei 11 em ação.


# Parte II — Visão Meso: os Componentes

## 1. Memória de execução (`internal/triagem`)

A memória guarda um envelope por Processo. O envelope tem `schema_version`,
`processo` e o mapa `sessoes`. Cada sessão tem `sessao_id`, `pid`,
`etapa_atual`, `proximo_passo`, `atualizado_em` e a telemetria de
`token_budget`.

Operação essencial:

| Função                              | Papel                                        |
| :---------------------------------- | :------------------------------------------- |
| `SalvarMemoria(dir, mem)`             | Grava a sessão com lock por arquivo          |
| `LerMemoria(dir, sessaoID)`           | Lê a sessão; sem ID, lê a legada             |
| `ListarSessoesAtivas(dir)`            | Lista só as sessões com PID vivo             |
| `ZerarSessao(dir, sessaoID)`          | Zera a sessão no fechamento                  |
| `LimparOrfaos(dir)`                   | Remove sessões de PID morto                  |

## 2. Anticorpos (`internal/brain`)

O banco tem dois escopos. `globais` vale para todo o ecossistema. `dominio`
vale para o domínio do Processo. Cada anticorpo tem `pacote`, `causa_raiz`,
`assinatura` SHA-256 e o `antidoto_aplicado`.

| Função                                    | Papel                                  |
| :---------------------------------------- | :------------------------------------- |
| `RecordComEscopo`                           | Grava anticorpo global ou de domínio   |
| `GateAntiReincidencia`                      | Decide: fluxo normal, antídoto ou bloqueio |
| `InjectPreventiveContext`                   | Injeta o antídoto no prompt do nó      |
| `NewAntibodiesDBPersistenteComGlobal`       | Abre o banco com anticorpos globais    |

## 3. Quarteto Lei 11 (`internal/fechar`)

O Fechamento Canônico zera a memória da sessão dona do Processo. O envelope
volta com `schema_version` igual a `memoria_execucao.v1` e sem sessões.

## 4. Memória em runtime (`internal/runner`)

A cada nó do DAG, o runner grava `etapa_atual` com o ID do nó e
`proximo_passo` com os sucessores do nó. Nó final aponta para `fim`.

## 5. CLI de sessões (`internal/cli`)

O comando `legoloop sessoes <slug> [--json] [--processos-dir DIR]` lista as
sessões ativas. A saída nunca expõe caminho absoluto.


# Parte III — Visão Micro: Schemas e Gates

## 1. Schema `memoria_execucao.v1`

```json
{
  "schema_version": "memoria_execucao.v1",
  "processo": "meu-processo",
  "sessoes": {
    "abc123": {
      "processo": "meu-processo",
      "sessao_id": "abc123",
      "pid": 4321,
      "etapa_atual": "M1.2",
      "proximo_passo": "M1.3",
      "atualizado_em": "2026-11-16T22:00:00Z"
    }
  }
}
```

Regras: chave do mapa é o `sessao_id`; `pid` positivo é obrigatório; escrita é
atômica; leitura corrompida gera backup `.legoloop-bak` e sentinela.

## 2. Schema `anticorpos.v1`

```json
{
  "schema_version": "anticorpos.v1",
  "globais": [],
  "dominio": {
    "meu-processo": [
      {
        "pacote": "internal/triagem",
        "causa_raiz": "lock sem PID",
        "assinatura": "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
        "antidoto_aplicado": "limparLockOrfao"
      }
    ]
  }
}
```

Regras: `assinatura` é SHA-256 do erro no formato `sha256:<hex>`; carga legado
é migrada sem perda; banco corrompido vira backup e o sistema segue vazio.

## 3. Gate anti-reincidência

```go
acao, antidoto, err := db.GateAntiReincidencia(pkg, errSig)
```

| Ação                | Quando                           | Resultado                   |
| :------------------ | :------------------------------- | :-------------------------- |
| `AcaoFluxoNormal`     | Erro desconhecido                | Executa o nó                |
| `AcaoAntidoto`        | Erro conhecido com antídoto      | Aplica o antídoto           |
| `AcaoBloqueado`       | Erro conhecido sem antídoto novo | `ErrGateAntiReincidencia`     |

## 4. Lock por arquivo

- O lock guarda o PID do detentor.
- PID morto vira lock órfão e o runtime toma o lock.
- `ErrLockOcupado` significa detentor vivo: aguarde.


# Parte IV — Prática: Operação do Dia a Dia

## 1. Consultar o progresso

```powershell
legoloop sessoes meu-processo --json
```

Saída esperada com duas sessões vivas:

```json
[
  {
    "sessao_id": "abc123",
    "pid": 4321,
    "etapa_atual": "M1.2",
    "proximo_passo": "M1.3"
  },
  {
    "sessao_id": "def456",
    "pid": 4322,
    "etapa_atual": "M2.1",
    "proximo_passo": "M2.2"
  }
]
```

Processo sem sessão viva devolve `[]` e exit 0.

## 2. Antes de executar um nó

1. Calcule a assinatura do erro possível.
2. Chame `brain.GateAntiReincidencia`.
3. `AcaoFluxoNormal`: execute o nó.
4. `AcaoAntidoto`: aplique o antídoto e execute.
5. `AcaoBloqueado`: pare e reporte a causa raiz.

## 3. Depois de um erro novo

1. Registre com `brain.RecordComEscopo`.
2. Escolha o escopo: global para erro de ecossistema, domínio para erro local.
3. Gravar o antídoto transforma bloqueio futuro em `AcaoAntidoto`.

## 4. No Fechamento Canônico

1. `fechar.ZerarSessao` roda no quarteto Lei 11.
2. O envelope volta sem sessões e com `schema_version` correto.
3. Confira com `legoloop sessoes <slug> --json`: deve devolver `[]`.

## 5. Roteiro de recuperação

| Sintoma                            | Ação                                          |
| :--------------------------------- | :-------------------------------------------- |
| `ErrLockOcupado`                     | Aguarde o detentor ou confira o `pid` do lock   |
| Lock órfão                         | Nada a fazer: o runtime toma o lock           |
| `ErrMemoriaCorrompida`               | Confira o backup `.legoloop-bak`                |
| `ErrAnticorposCorrompidos`           | Confira o backup e regrave os anticorpos      |
| `ErrGateAntiReincidencia`            | Corrija a causa raiz e registre o antídoto    |


# Parte V — Limites e Apêndices

## 1. Limites declarados

- O lock por arquivo cobre o mesmo host. Duas máquinas não compartilham lock.
- O gate anti-reincidência depende da assinatura do erro. Erro com texto variável gera assinaturas diferentes.
- A telemetria de `token_budget` só aparece para nós com orçamento declarado.
- O zeramento da sessão cobre a sessão dona do Fechamento. Sessões alheias seguem vivas.

## 2. Apêndice A — Divergências anotadas na integração

| Tema                    | Registro                                                 |
| :---------------------- | :------------------------------------------------------- |
| Contagem de tickets     | O brief cita 37 tickets; o plano paralelo tem 11 tickets e a ROUTE tem 15 nós |
| Nome de erro            | Termo único `ErrLockOcupado` (Lei 15); divergência `ErrLockAtivo` resolvida em 2026-11-17 |
| Pasta de Processos      | Os schemas citam `processos/<slug>`; o acervo real vive em `processos/<slug>` |

## 3. Apêndice B — Portões de qualidade

- `go test -count=1 ./...`
- `go vet ./...`
- `gofmt -l .` com saída vazia
- `go run ./tools/funclen internal` sem estouro novo
- `go run ./tools/taxlinter internal cmd loops skills docs/templates` sem violação nova
- `git diff --exit-code -- go.mod go.sum`

## 4. Apêndice C — Artefatos da demanda

| Artefato                | Caminho                                                             |
| :---------------------- | :------------------------------------------------------------------ |
| Manual operacional      | `docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md`     |
| Skill canônica          | `skills/legoloop-memoria-anticorpos/SKILL.md`                         |
| Livro didático          | `docs/books/memoria-anticorpos/`                                      |
| Relatório de integração | `docs/reports/2026-11-16_RELATORIO_INTEGRACAO_MEMORIA_ANTICORPOS.md`      |


