---
title: "Addy Osmani no LegoLoop: Performance e Engenharia Medidas"
subtitle: "Medir antes de otimizar, restringir antes de construir e provar antes de entregar"
author:
  - Heverton & Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "WEB PERFORMANCE E QUALIDADE DE ENGENHARIA"
cover-title: "Addy Osmani"
cover-accent: "Performance e Engenharia Medidas"
cover-caption: "OS SETE NÓS DO LOOP addy-osmani"
cover-stages:
  - Alinhar Contexto
  - Especificar Restrições
  - Planejar Tickets
  - Construir com TDD
  - Runtime e Web Vitals
  - Review e Hardening
  - Observabilidade e Ship
header-mark: "ADDY OSMANI"
accent: "2563eb"
toc: true
toc-depth: 2
abstract: |
  Este livro ensina a operar o método de Addy Osmani (upstream sob licença MIT) dentro do ecossistema LegoLoop. O método reúne 25 habilidades de engenharia que combatem dois problemas: o achismo de desempenho, em que se otimiza sem medir, e a proliferação de código assistido por IA sem controle de qualidade.

  A Parte I apresenta os seis pilares do método. A Parte II descreve o loop mestre de sete estágios e seus seis gates. A Parte III detalha os sub-loops e as habilidades de cada um. A Parte IV desce ao nível de implementação: limiares de Core Web Vitals, os cinco eixos de revisão, os playbooks da skill e o arquivo loops/addy-osmani.loop.json. A Parte V traz receitas e a Parte VI reúne glossário, relatório de conformidade e referência das habilidades.

  Fonte primária: docs/manuals/2026-11-02_MANUAL_ADDY_OSMANI_LOOP.md, complementada pelos arquivos upstream da skill aidd-osmani.
---

# Como ler este livro

## A quem este livro se destina

1. **Quem responde pela qualidade de um produto web.** Precisa entender por que "o site está lento" não se resolve com intuição e como orçamentos de performance viram regra. A Parte I responde a isso.
2. **Quem desenvolve com agentes.** Precisa saber qual habilidade chamar, o que cada gate exige e como provar uma correção. As Partes II, III e V são o núcleo para esse perfil.
3. **Quem mantém a plataforma LegoLoop.** Precisa saber onde o pipeline está declarado, quais artefatos cada nó exige e onde o manual e o código divergem. A Parte IV e o Apêndice B são escritos para esse leitor.

## A estrutura da obra

| Parte | Nível     | Conteúdo                                                     |
| :---- | :-------- | :----------------------------------------------------------- |
| I     | Macro     | Achismo de desempenho, código sem controle e os seis pilares |
| II    | Meso      | O loop mestre de sete estágios e seis gates                  |
| III   | Meso      | Os sub-loops e suas habilidades                              |
| IV    | Micro     | Limiares, eixos de revisão, playbooks e integração           |
| V     | Prática   | Receitas para situações frequentes                           |
| VI    | Apêndices | Glossário, relatório de conformidade e referência rápida     |

## As convenções da obra

Nomes de habilidades e caminhos de arquivo aparecem em fonte monoespaçada, como `performance-optimization` ou `loops/addy-osmani.loop.json`. Exemplos ilustrativos são marcados como tal no texto. Toda afirmação sobre o LegoLoop é rastreável a um arquivo citado; quando manual e código discordam, o livro segue o código e registra a diferença no Apêndice B.
