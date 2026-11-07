---
title: "Impeccable no LegoLoop: Design de Interfaces sem Clichês"
subtitle: "O loop de cinco estágios, os seis sub-loops e o piso de qualidade que separa uma tela construída de uma tela montada"
author:
  - Heverton & Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "ENGENHARIA VISUAL DE INTERFACES"
cover-title: "Impeccable"
cover-accent: "Design de Interfaces sem Clichês"
cover-caption: "OS SEIS NÓS DO LOOP impeccable-design"
cover-stages:
  - Avaliar
  - Construir
  - Melhorar
  - Refinar
  - Corrigir
  - Iterar e Documentar
header-mark: "IMPECCABLE"
accent: "7c3aed"
toc: true
toc-depth: 2
abstract: |
  Este livro ensina a operar o método Impeccable, criado por Paul Bakaus (upstream v4.7.0), dentro do ecossistema LegoLoop. O Impeccable é um conjunto de cerca de 25 comandos de design de interface organizados em seis sub-loops: Avaliar, Construir, Melhorar, Refinar, Corrigir e Iterar.

  A obra segue a progressão Macro, Meso e Micro. A Parte I explica por que interfaces geradas por IA convergem para o mesmo visual genérico e quais pilares o método usa para evitar isso. As Partes II e III descrevem o loop mestre de cinco estágios, seus quatro gates de parada e cada comando dos sub-loops. A Parte IV desce ao nível de implementação: o arquivo loops/impeccable-design.loop.json, o template do Studio Canvas e os agentes especializados. A Parte V traz receitas de uso e a Parte VI reúne glossário, referência rápida e um relatório de conformidade entre o manual e o código.

  Fonte primária: docs/manuals/2026-11-01_MANUAL_IMPECCABLE_DESIGN_LOOP.md, complementada pelo arquivo de referência upstream craft-floor.md.
---

# Como ler este livro

## A quem este livro se destina

O livro atende três leitores diferentes, e cada um pode começar por um ponto distinto.

1. **Quem decide sobre produto.** Precisa entender por que uma interface "bonita" gerada por IA ainda parece amadora e como um processo com gates reduz retrabalho. A Parte I responde a isso.
2. **Quem opera o LegoLoop.** Precisa saber qual comando chamar, em que ordem e o que esperar como saída. As Partes II, III e V são o núcleo para esse perfil.
3. **Quem mantém a plataforma.** Precisa saber onde o pipeline está declarado, como o Studio o instancia e onde o manual e o código divergem. A Parte IV e o Apêndice B são escritos para esse leitor.

## A estrutura da obra

| Parte | Nível     | Conteúdo                                                           |
| :---- | :-------- | :----------------------------------------------------------------- |
| I     | Macro     | O problema do visual genérico e os quatro pilares do método        |
| II    | Meso      | Os quatro modos de superfície e o loop mestre de cinco estágios    |
| III   | Meso      | Os seis sub-loops e cada um dos seus comandos                      |
| IV    | Micro     | Piso de qualidade, agentes especializados e integração no LegoLoop |
| V     | Prática   | Receitas para situações frequentes                                 |
| VI    | Apêndices | Glossário, referência rápida e relatório de conformidade           |

## As convenções da obra

Comandos e caminhos de arquivo aparecem em fonte monoespaçada, como `polish` ou `loops/impeccable-design.loop.json`. Exemplos ilustrativos são marcados como tal no texto. Toda afirmação sobre o LegoLoop é rastreável a um arquivo citado; quando manual e código discordam, o livro segue o código e registra a diferença no Apêndice B.
