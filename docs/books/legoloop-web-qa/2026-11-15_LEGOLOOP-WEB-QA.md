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


# Parte I — Visão Macro: Por que Auditar antes de Corrigir

## O problema das correções reativas em interfaces web

Aplicações web modernas falham silenciosamente em quatro dimensões que testes unitários de back-end não enxergam: controles visuais que não produzem efeito ao clique, elementos sobrepostos em telas menores, violações de contraste ou navegação por teclado, e links ou assets que retornam erro HTTP após um build. Quando agentes de código tentam corrigir um defeito visual diretamente a partir de uma queixa genérica, o resultado costuma ser regressão em outra rota ou viewport.

O processo `legoloop-web-qa` resolve esse ciclo separando rigorosamente **diagnóstico determinístico** de **mutação de código**. Primeiro, o sistema prova matematicamente o estado da aplicação através de oito ferramentas reais. Somente depois que o relatório consolidado existe em disco, o operador humano e o agente iniciam a entrevista socrática de afunilamento que origina o pacote de correção governado até `BUILD`.

## Três princípios arquiteturais do processo

1. **Zero Stubs no Preflight:** Se qualquer uma das oito ferramentas obrigatórias não estiver instalada e acessível, o processo aborta no nó `n01_preflight` com `exit_code == 1`. Nenhuma etapa simula resultado positivo ou substitui uma peça ausente por resposta vazia.
2. **Dupla Cobertura de Ambiente:** O runner `tools/webqa/webqa.mjs` aceita tanto uma URL já publicada (`--url`) quanto um comando e porta locais (`--cmd` + `--port`), subindo a aplicação, aguardando a porta responder e encerrando o processo filho ao final da execução.
3. **Relatório como Contrato de Entrada:** A saída em `docs/qa-web/YYYY-MM-DD_<alvo>/` não é um log descartável. Os arquivos `RELATORIO_QA_WEB.json` e `RELATORIO_QA_WEB.md` são insumos obrigatórios para o nó `n11_entrevista_correcao`.


# Parte II — O Arsenal de 8 Peças Complementares

## Por que exatamente estas oito ferramentas

Cada ferramenta do processo `legoloop-web-qa` cobre uma classe distinta de defeito sem duplicar esforço nem depender de serviços pagos fechados:

| Camada | Ferramenta | O que detecta | Modo de Execução |
| :--- | :--- | :--- | :--- |
| Integridade de Rede | `lychee` | Links, âncoras e assets com HTTP >= 400 ou timeout | Binário Rust CLI (`--format json`) |
| Render Rápido Local | `obscura` | Falhas de carregamento, bloqueio de rede local e DOM inicial | Binário Rust (`--allow-private-network`) |
| Runtime e Sessão | `playwright` | Erros de console, `pageerror`, falhas de rota e sessões com `storageState` | Chromium headless (`1440x900` e `390x844`) |
| UI Comportamental | `renderlint` | Botões mortos (`dead-control`), conteúdo encoberto (`obscured-content`) e `text-overflow` | CLI Node sobre navegador real |
| Consistência Visual | `aidd-linter-design` | Desvios de design tokens, espaçamento e contraste estrutural | Skill canônica de design tokens |
| Acessibilidade | `axe-core` | Violações WCAG (`critical` e `serious`), ARIA e foco por teclado | `@axe-core/cli` (`--exit`) |
| Saúde Multi-Rota | `unlighthouse` | Queda de performance, SEO e boas práticas abaixo de 90 em todas as rotas | `@unlighthouse/cli` (`--no-open`) |
| Regressão Visual | `lost-pixel` | Diferença pixel a pixel contra baseline aprovada | CLI open-source em `.lostpixel/baseline` |

## Separação entre `obscura` e `playwright`

O motor `obscura` executa verificação ultrarrápida em V8 nativo (~30MB de RAM) com permissão explícita de rede local (`--allow-private-network`), validando prontidão da página antes de acionar sessões completas de Chromium no `playwright`. O `playwright`, por sua vez, garante fidelidade total de eventos JavaScript, captura de `console.error`, injeção opcional de `storageState` para fluxos autenticados e teste nos dois viewports padrão.


# Parte III — Anatomia do DAG, Viewports e Limiares

## Topologia em 8 Camadas

O grafo `loops/legoloop-web-qa/pipeline.dag.json` (`schema_version: 2`) organiza 12 nós em 8 camadas topológicas, garantindo isolamento estrito (`orchestrix_hermetic`) e `write_boundary: strict`:

1. **Camada 1 (`n01_preflight`):** Verifica presença real das 8 peças no PATH, `node_modules` local e cofre de skills.
2. **Camada 2 (`n02_boot_e_descoberta`):** Sobe o servidor local se `--cmd` e `--port` forem fornecidos e descobre rotas via crawl HTML + `/sitemap.xml` + lista explícita `--routes`.
3. **Camada 3 (Execução Paralela):** Dispara simultaneamente cinco auditores independentes: `n03_links_lychee`, `n04_obscura_e_playwright`, `n05_renderlint_e_design`, `n06_axe_acessibilidade` e `n07_unlighthouse_rotas`.
4. **Camada 4 (`n08_lost_pixel_visual`):** Compara capturas nos viewports `desktop (1440x900)` e `mobile (390x844)` após `n04_obscura_e_playwright`.
5. **Camada 5 (`n09_sg_baseline_visual`):** Stop-Gate humano obrigatório sempre que `.lostpixel/baseline` ainda não existir.
6. **Camada 6 (`n10_relatorio_consolidado`):** Consolida métricas e achados em `docs/qa-web/YYYY-MM-DD_<alvo>/`.
7. **Camada 7 (`n11_entrevista_correcao`):** Conduz entrevista socrática de afunilamento lendo o relatório gerado.
8. **Camada 8 (`n12_governanca_build`):** Materializa o pacote `INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM` para o ciclo `BUILD`.

## Matriz de Limiares e Códigos de Saída

A função `evaluateThresholds` em `tools/webqa/webqa.mjs` aplica critérios binários sem margem subjetiva:

- **`exit_code == 0` (`PASSED`):** Zero links quebrados, zero erros de console/rota, zero violações `critical`/`serious` no `axe-core`, zero violações no `renderlint`/`aidd-linter-design`/`obscura`, nota mínima >= 90 no `unlighthouse`, baseline existente e `diffCount == 0` no `lost-pixel`.
- **`exit_code == 2` (`NEEDS_BASELINE_STOP_GATE`):** Todas as métricas funcionais passaram, mas não há baseline prévia do `lost-pixel`. O processo para e solicita confirmação humana no nó `n09_sg_baseline_visual`.
- **`exit_code == 1` (`FAILED`):** Qualquer violação de limiar acima ou ferramenta ausente no preflight.


# Parte IV — Prática Operacional e Cadeia até BUILD

## Executando o Runner em Projetos Locais ou Publicados

Antes de rodar qualquer auditoria, valide a integridade do grafo e do ambiente:

```bash
legoloop validate loops/legoloop-web-qa/pipeline.dag.json
node --test tools/webqa/webqa.test.mjs
node tools/webqa/webqa.mjs preflight --json
```

Para auditar uma aplicação local subindo o servidor automaticamente e preservando sessão autenticada:

```bash
node tools/webqa/webqa.mjs run \
  --cmd "npm run dev" \
  --port 3000 \
  --alvo meu-portal \
  --routes "/,/dashboard,/config" \
  --storage-state auth.json
```

## Estrutura do Relatório em `docs/qa-web/`

A função `writeConsolidatedReport` aplica Path Jail estrito (via `sanitizeTargetSlug`) e grava dois arquivos complementares dentro de `docs/qa-web/YYYY-MM-DD_<alvo>/`:

- **`RELATORIO_QA_WEB.json`:** Payload estruturado com `targetUrl`, `routes`, `viewports`, `status`, `failures`, `metrics` e `findings`.
- **`RELATORIO_QA_WEB.md`:** Sumário legível para humanos e agentes contendo todos os achados numerados por ferramenta, severidade e rota.

## O Fluxo Obrigatório Pós-Relatório até `BUILD`

Diferente de scripts avulsos de CI, o `legoloop-web-qa` proíbe correções improvisadas após um relatório vermelho. Quando `RELATORIO_QA_WEB.md` registra falhas, a Camada 7 (`n11_entrevista_correcao`) lê o relatório como insumo e conduz uma entrevista socrática de afunilamento (uma pergunta por vez até confiança >= 95%) para decidir prioridade, escopo e riscos de cada correção.

Concluída a entrevista, a Camada 8 (`n12_governanca_build`) emite a cadeia completa de governança:

1. `INTERVIEW` — decisões validadas para cada achado do relatório.
2. `INTENT` e `SPEC` — critérios binários de correção por componente ou rota afetada.
3. `PLAN`, `ROUTE` e `DAG` — fatiamento topológico das correções sem conflito de escrita.
4. `SWARM` e ciclo `BUILD` — execução isolada em worktree com reexecução do próprio `legoloop-web-qa` como gate final de aceitação.


# Parte V — Guardrails Negativos, FMEA e Apêndices

## Regras Negativas Invioláveis

1. **Proibição de Fallback Falso:** Se `renderlint`, `lychee`, `axe`, `unlighthouse` ou `lost-pixel` não existirem no ambiente, o runner interrompe no `n01_preflight`. É proibido retornar `exit_code == 0` simulando ferramentas ausentes.
2. **Proibição de Edição Cega:** Nenhum agente pode alterar componentes React/HTML diretamente após ler `RELATORIO_QA_WEB.md` sem concluir `n11_entrevista_correcao` e `n12_governanca_build`.
3. **Path Jail em `docs/qa-web/`:** Slugs contendo `..`, `/` ou `\\` são rejeitados antes de criar diretórios em disco.
4. **Proteção da Baseline Visual:** O snapshot inicial do `lost-pixel` nunca se autopromove a baseline oficial sem confirmação humana em `n09_sg_baseline_visual`.

## Matriz de Riscos (FMEA Simplificada)

| Nó | Modo de Falha | Efeito | Contenção Mecânica |
| :--- | :--- | :--- | :--- |
| `n01_preflight` | Ferramenta ausente no PATH | Auditoria incompleta | Abort imediato com comando exato de instalação |
| `n02_boot_e_descoberta` | Porta ocupada ou servidor não responde | Travamento de pipeline | Timeout de 30s no healthcheck e `kill()` no `finally` |
| `n08_lost_pixel_visual` | Mudança legítima de layout confundida com bug | Falso positivo contínuo | Stop-Gate humano `n09_sg_baseline_visual` |
| `n10_relatorio_consolidado` | Slug malicioso no `--alvo` | Escrita fora do repositório | `sanitizeTargetSlug` + verificação `startsWith(qaRoot)` |

## Apêndice: Rastreabilidade dos Artefatos

- **Skill Canônica:** `skills/legoloop-web-qa/SKILL.md`
- **DAG Executivo (v2):** `loops/legoloop-web-qa/pipeline.dag.json` e `docs/dags/2026-11-15_DAG_WEB_QA.json`
- **Rota v2:** `docs/routes/2026-11-15_ROUTE_WEB_QA.json`
- **SPEC e PLAN:** `docs/specs/2026-11-15_SPEC_WEB_QA.md` e `docs/plans/2026-11-15_PLAN_WEB_QA.md`
- **Fichas Técnicas:** `docs/fichas/legoloop-web-qa/n01.md` a `n12.md`
- **Runner e Suíte TDD:** `tools/webqa/webqa.mjs` e `tools/webqa/webqa.test.mjs`
