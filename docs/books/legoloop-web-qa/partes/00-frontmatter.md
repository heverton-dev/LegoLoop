---
title: "Web-QA Determinístico no LegoLoop: Da Auditoria de Interface ao Ciclo BUILD"
subtitle: "Como auditar aplicações web com oito ferramentas reais sem mocks e converter achados em correções governadas"
author:
  - Heverton & Ecossistema LegoLoop
date: "15/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "ENGENHARIA DE QUALIDADE WEB"
cover-title: "Web-QA Determinístico"
cover-accent: "8 Ferramentas Reais até o Ciclo BUILD"
cover-caption: "CADEIA CANÔNICA DO LEGOLOOP-WEB-QA"
cover-stages:
  - Preflight Estrito
  - Descoberta e Viewports
  - Auditoria Paralela
  - Regressão Visual
  - Relatório Consolidado
  - Entrevista até BUILD
header-mark: "LEGOLOOP WEB-QA"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro documenta a arquitetura, os contratos e a operação do processo canônico `legoloop-web-qa` no ecossistema LegoLoop. O processo une oito ferramentas reais (`playwright`, `obscura`, `aidd-linter-design`, `renderlint`, `axe-core`, `unlighthouse`, `lychee` e `lost-pixel`) sob um DAG determinístico com zero tolerância a stubs ou simulações.

  A Parte I apresenta o problema de interfaces quebradas e o modelo mental macro de auditoria antes da mutação de código. A Parte II detalha as oito peças do arsenal e suas responsabilidades complementares. A Parte III disseca as oito camadas do DAG `pipeline.dag.json`, os viewports (`1440x900` e `390x844`) e os limiares numéricos de reprovação. A Parte IV demonstra a operação prática via `tools/webqa/webqa.mjs` e o encadeamento obrigatório pós-relatório (`INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM -> BUILD`). A Parte V consolida guardrails negativos, matriz de falhas FMEA e apêndices de rastreabilidade.

  Fonte primária: `docs/manuals/2026-11-15_MANUAL_WEB_QA.md`, `docs/specs/2026-11-15_SPEC_WEB_QA.md` e `loops/legoloop-web-qa/pipeline.dag.json`.
---

# Como ler este livro

## A quem este livro se destina

1. **Quem opera ou desenvolve aplicações web com agentes.** Precisa detectar botões mortos, sobreposições, erros de console, violações WCAG e quedas de Lighthouse antes de publicar. As Partes II, III e IV são o núcleo prático.
2. **Quem governa processos no LegoLoop.** Precisa garantir que nenhum relatório de QA dispare edição cega de código sem entrevista socrática e cadeia completa até `BUILD`. As Partes I, III e V estabelecem esses contratos.
3. **Quem mantém o runner determinístico `tools/webqa/webqa.mjs`.** Precisa compreender os contratos de preflight, Path Jail, descoberta de rotas e suíte TDD em `tools/webqa/webqa.test.mjs`.

## Estrutura da obra

- **Parte I — Visão Macro:** Por que auditorias superficiais falham e como o `web-qa` separa verificação de correção.
- **Parte II — O Arsenal de 8 Peças:** Papel técnico de cada ferramenta e ausência de sobreposição redundante.
- **Parte III — Anatomia do DAG e Limiares:** As 8 camadas topológicas, paralelismo, viewports e Stop-Gate visual.
- **Parte IV — Prática e Cadeia até BUILD:** Execução do runner, geração de `docs/qa-web/` e condução da entrevista socrática de correção.
- **Parte V — Limites, FMEA e Apêndices:** Regras negativas invioláveis, matriz de riscos e rastreabilidade dos artefatos.
