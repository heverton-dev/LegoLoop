---
title: "Gateway de LLM no Antigravity CLI: Chaves por Variável de Ambiente"
subtitle: "Como apontar o agy para o 9Router ou para qualquer outro destino, sem trocar de conta Google"
author:
  - Heverton & Ecossistema LegoLoop
date: "13/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "OPERAÇÃO SEGURA DE AGENTES"
cover-title: "Gateway de LLM"
cover-accent: "Chaves por Variável de Ambiente no AGY"
cover-caption: "OS TRÊS MÉTODOS DE CONFIGURAÇÃO"
cover-stages:
  - Entender o Risco
  - Escolher o Destino
  - Configurar pelo Sistema
  - Configurar pelo Terminal
  - Pedir ao Agente
  - Validar e Desligar
header-mark: "GATEWAY AGY"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro ensina a apontar o Antigravity CLI (agy) para um gateway de LLM por variáveis de ambiente. O método troca o destino das chamadas de modelo e não toca na conta Google. Ele substitui o rodízio de contas, que causou o bloqueio de duas contas em 10/11/2026.

  A Parte I explica o problema e corrige três premissas comuns. A Parte II apresenta os três métodos de configuração: sistema operacional, linha de comando e pedido a um agente. A Parte III detalha as variáveis, os protocolos e a regra de precedência. A Parte IV traz receitas para os casos frequentes. A Parte V descreve as proibições que protegem as contas. A Parte VI reúne glossário, relatório de conformidade e fontes.

  Fonte primária: docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md, complementada pela documentação do Google Cloud sobre o gateway corporativo do Antigravity.
---

# Como ler este livro

## A quem este livro se destina

1. **Quem usa o `agy` no dia a dia.** Precisa ligar, trocar e desligar o gateway sem erro. As Partes II e IV são o núcleo.
2. **Quem lidera times que dividem contas.** Precisa entender por que o rodízio de contas é perigoso e o que o gateway muda. As Partes I e V respondem a isso.
3. **Quem mantém a plataforma LegoLoop.** Precisa saber quais comandos do `gw` foram bloqueados e por quê. A Parte V e o Apêndice B foram escritos para esse leitor.

## A estrutura da obra

| Parte            | Nível            | Conteúdo                                        |
| :--------------- | :--------------- | :---------------------------------------------- |
| I                | Visão geral      | O problema do rodízio e as três premissas corrigidas |
| II               | Método           | Os três métodos de configuração                 |
| III              | Detalhe          | Variáveis, protocolos e precedência             |
| IV               | Prática          | Receitas para os casos frequentes               |
| V                | Método           | Limites e proibições                            |
| VI               | Apêndices        | Glossário, conformidade e fontes                |
