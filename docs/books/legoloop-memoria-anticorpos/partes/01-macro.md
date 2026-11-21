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
