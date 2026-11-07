---
title: "Verbos e Granularidade no LegoLoop"
subtitle: "Guia Didático BDD da Tríade de Triggers e do Encadeamento de Verbos"
author:
  - Heverton & Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "GUIA DIDÁTICO · BEHAVIOR-DRIVEN DEVELOPMENT"
tagline: |
  O mesmo comportamento, três vias de acionamento, contratos Dado-Quando-Então idênticos
toc: true
toc-depth: 2
abstract: |
  Este guia ensina a operar o LegoLoop por meio de seus verbos canônicos e dos cinco níveis de granularidade que organizam qualquer fluxo: do Mega Pipeline às Nano-Etapas.

  A obra percorre três camadas. No nível Macro, apresenta a filosofia de blocos componíveis governados por um grafo acíclico dirigido e a Tríade de Triggers (Terminal, Chat do Agente e LegoLoop Studio). No nível Meso, documenta as cinco trilhas canônicas de encadeamento de verbos. No nível Micro, especifica sete cenários BDD em que o mesmo comportamento é acionado pelas três vias, com resultado idêntico e auditável.
---

# Como ler este guia

## Público e pré-requisitos

O guia se destina a desenvolvedores, engenheiros de IA e arquitetos de software que operam o ecossistema LegoLoop e AIDD. Pressupõe familiaridade com linha de comando e com a noção de pipeline; nenhum conhecimento prévio do motor interno é exigido.

## Estrutura da obra

| Parte | Nível | Conteúdo |
| :--- | :--- | :--- |
| **I** | Macro | Filosofia Lego, níveis de granularidade e Tríade de Triggers |
| **II** | Meso | Matriz de encadeamento e as cinco trilhas canônicas |
| **III** | Micro | Sete cenários BDD multi-trigger |
| **IV** | Referência | Tabela de consulta rápida e síntese final |

## Convenção BDD

Cada cenário segue a forma **Dado-Quando-Então**. O bloco *Dado* fixa o estado inicial, o bloco *Quando* descreve a ação nas três vias de acionamento, e o bloco *Então* declara o resultado verificável. O comportamento contratado não depende da via escolhida.
