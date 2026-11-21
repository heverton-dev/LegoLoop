---
title: "Processo Canônico de Gestão da Matriz de Peças"
subtitle: "Governança matemática, catalogação atômica, tiers de reputação e vocabulário controlado no LegoLoop"
author:
  - Heverton & Ecossistema LegoLoop
date: "20/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "CATÁLOGO E REPUTAÇÃO DE PEÇAS"
cover-title: "Matriz de Peças"
cover-accent: "Governança e Tiers de Reputação"
cover-caption: "ARQUITETURA CANÔNICA DO SUBSISTEMA MATRIZ"
cover-stages:
  - Ingestão Incremental
  - Cálculo de Hashes
  - Auditoria de Qualidade
  - Tiers de Reputação
  - Vocabulário ASD-STE100
  - Consulta Multidimensional
header-mark: "LEGOLOOP MATRIZ"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro documenta a arquitetura, os contratos e a operação do processo canônico `legoloop-matriz` no ecossistema LegoLoop. O subsistema governa o inventário central de peças (`registry_matrix.json`), assegurando que cada componente possua proveniência, integridade mecânica por SHA-256 e reputação auditável.

  A Parte I apresenta a fundação conceitual da Matriz, o modelo matemático de pontuação de peças e a hierarquia dos tiers de reputação. A Parte II disseca a arquitetura técnica, o vocabulário controlado ASD-STE100 e as dimensões de perfil (etapa, natureza e caso de uso). A Parte III detalha a operação da interface de linha de comando (`legoloop matriz`), as garantias transacionais de gravação atômica e os procedimentos de recuperação.

  Fonte primária: `docs/manuals/2026-11-06_MANUAL_FUNCAO_MATRIZ_LEGOLOOP.md`, `skills/legoloop-matriz/SKILL.md` e `internal/matriz/`.
---

# Como ler este livro

## A quem este livro se destina

1. **Operadores de Agentes e Desenvolvedores.** Profissionais que precisam consultar peças aptas, inspecionar tiers de reputação e garantir reutilização sem duplicidade. A Parte I e a Parte III são leituras prioritárias.
2. **Arquitetos de Governança.** Responsáveis pela integridade do cofre canônico, aprovação de vocabulário e certificação de novos componentes. A Parte II detalha os contratos invariantes.
3. **Engenheiros de Plataforma LegoLoop.** Mantenedores do núcleo em Go que estendem os algoritmos de score, filtros de busca e integração com os harnesses. Todas as partes constituem especificação formal.

## Estrutura da obra

| Parte          | Título                    | Conteúdo Principal                       |
| :------------- | :------------------------ | :--------------------------------------- |
| I              | Fundação e Modelo Mental  | Problema da defasagem, pontuação e tiers de reputação |
| II             | Arquitetura e Vocabulário | Subsistemas, metadados, ASD-STE100 e perfis de uso |
| III            | Operação e Invariantes    | Comandos da CLI, gravação atômica e recuperação |
