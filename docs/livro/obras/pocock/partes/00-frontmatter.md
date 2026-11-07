---
title: "Matt Pocock no LegoLoop: Rigor de Engenharia com Agentes"
subtitle: "Alinhar antes de construir, testar antes de implementar e corrigir o ambiente em vez do passado"
author:
  - Heverton & Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "ENGENHARIA DE SOFTWARE COM AGENTES"
cover-title: "Matt Pocock"
cover-accent: "Rigor de Engenharia com Agentes"
cover-caption: "OS SEIS NÓS DO LOOP matt-pocock"
cover-stages:
  - Alinhar Domínio
  - Especificar
  - Planejar Tickets
  - Construir com TDD
  - Revisar em Dois Eixos
  - Retroalimentar o Ambiente
header-mark: "MATT POCOCK"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro ensina a operar o método de Matt Pocock (upstream sob licença MIT) dentro do ecossistema LegoLoop. O método substitui a geração cega de código por um ciclo disciplinado: entrevistar a ideia até não restarem premissas ocultas, transformar o consenso em especificação, fatiar a especificação em tickets com dependências explícitas e implementar cada ticket com testes primeiro, em contexto limpo.

  A Parte I apresenta os cinco pilares do método. A Parte II descreve o loop mestre e seus quatro gates. A Parte III detalha os seis sub-loops e as habilidades de cada um. A Parte IV desce ao nível de implementação: roteamento de modelos, higiene de contexto e o arquivo loops/matt-pocock.loop.json. A Parte V traz receitas e a Parte VI reúne glossário, relatório de conformidade e referência das habilidades.

  Fonte primária: docs/manuals/2026-11-01_MANUAL_MATT_POCOCK_SKILLS_LOOP.md, complementada pelas habilidades upstream tdd, code-review e diagnosing-bugs.
---

# Como ler este livro

## A quem este livro se destina

1. **Quem lidera times que usam agentes.** Precisa entender por que "deixar a IA programar" produz código que ninguém consegue manter, e onde o humano precisa continuar decidindo. A Parte I responde a isso.
2. **Quem desenvolve com agentes no dia a dia.** Precisa saber qual habilidade chamar, quando limpar o contexto e o que cada gate exige. As Partes II, III e V são o núcleo para esse perfil.
3. **Quem mantém a plataforma LegoLoop.** Precisa saber onde o pipeline está declarado, quais artefatos cada nó exige e onde o manual e o código divergem. A Parte IV e o Apêndice B são escritos para esse leitor.

## A estrutura da obra

| Parte | Nível     | Conteúdo                                                 |
| :---- | :-------- | :------------------------------------------------------- |
| I     | Macro     | O problema da geração cega e os cinco pilares do método  |
| II    | Meso      | O loop mestre, da ideia à produção, com quatro gates     |
| III   | Meso      | Os seis sub-loops e suas habilidades                     |
| IV    | Micro     | Roteamento de modelos, higiene de contexto e integração  |
| V     | Prática   | Receitas para situações frequentes                       |
| VI    | Apêndices | Glossário, relatório de conformidade e referência rápida |

## As convenções da obra

Nomes de habilidades e caminhos de arquivo aparecem em fonte monoespaçada, como `grill-with-docs` ou `loops/matt-pocock.loop.json`. Exemplos ilustrativos são marcados como tal no texto. Toda afirmação sobre o LegoLoop é rastreável a um arquivo citado; quando manual e código discordam, o livro segue o código e registra a diferença no Apêndice B.
