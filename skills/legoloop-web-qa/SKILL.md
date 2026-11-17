---
name: legoloop-web-qa
description: "Use quando precisar executar o processo legoloop-web-qa: auditoria determinística de web apps com playwright, obscura, aidd-linter-design, renderlint, axe-core, unlighthouse, lychee e lost-pixel, gerando relatório consolidado em docs/qa-web/ como insumo obrigatório para entrevista socrática de correção e geração dos artefatos INTERVIEW, INTENT, SPEC, PLAN, ROUTE, DAG e SWARM para ciclo BUILD."
consumes: Artifact:TargetURL
produces: Artifact:QAReport
artifact: docs/qa-web/RELATORIO_QA_WEB.md
gate: exit_code == 0
suggested_model: sonnet
---

# 🧩 legoloop-web-qa — Auditoria Determinística de Web Apps até BUILD

> **Versão:** 1.1.0  
> **DAG Canônico:** `docs/dags/2026-11-15_DAG_WEB_QA.json`  
> **Runner Determinístico:** `tools/webqa/webqa.mjs`

---

## 1. Contexto & Arsenal das 8 Peças Obrigatórias

Esta skill governa a execução do processo `web-qa`. Nenhuma ferramenta pode ser simulada ou ignorada. Se faltar qualquer peça no preflight, a execução bloqueia com `exit_code != 0`.

1. **Preflight & Orquestração:** `node tools/webqa/webqa.mjs preflight --json`
2. **Descoberta & Viewports:** crawl + `sitemap.xml` + `--routes`, com arquivo de sessão opcional do Playwright em `desktop (1440x900)` e `mobile (390x844)`.
3. **Camada Paralela de Auditoria:**
   - `lychee`: integridade de links e assets (falha se `brokenLinks > 0`).
   - `obscura` + `playwright`: boot/fetch local (`--allow-private-network`), erros de rota HTTP >= 400 e erros de console.
   - `renderlint` + `aidd-linter-design`: botões mortos (`dead-control`), sobreposição (`obscured-content`), corte de texto e tokens.
   - `axe-core`: acessibilidade WCAG (falha se `critical > 0` ou `serious > 0`).
   - `unlighthouse`: varredura multi-rota Lighthouse (falha se nota mínima < 90).
4. **Regressão Visual & Stop-Gate:**
   - `lost-pixel`: se não houver baseline em `.lostpixel/baseline`, cria snapshot inicial e exige aprovação no Stop-Gate `n09_sg_baseline_visual`.
5. **Relatório & Encadeamento Obrigatório até BUILD:**
   - Grava `RELATORIO_QA_WEB.json` e `RELATORIO_QA_WEB.md` em `docs/qa-web/YYYY-MM-DD_<alvo>/`.
   - Inicia obrigatoriamente a **entrevista socrática de afunilamento** (`aidd-grill`) lendo o relatório, seguida da geração da cadeia completa `INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM` para o ciclo de `BUILD` das correções.

---

## 2. Procedimento Determinístico

```bash
# 1. Validar o grafo DAG v2
legoloop validate docs/dags/2026-11-15_DAG_WEB_QA.json

# 2. Verificar todas as 8 ferramentas (bloqueia se faltar alguma)
legoloop web-qa preflight --json

# 3. Executar contra URL publicada, comando local ou diretório de projeto (--project)
legoloop web-qa run --url http://127.0.0.1:3000 --alvo meu-app
legoloop web-qa run --cmd "npm run dev" --port 3000 --alvo meu-app --routes "/,/login"
legoloop web-qa run --project . --alvo legoloop-studio
legoloop web-qa run --project ../meu-webapp --port 4173 --alvo meu-webapp

# 4. Iniciar entrevista socrática de correção diretamente a partir do relatório
legoloop web-qa corrigir --relatorio docs/qa-web/YYYY-MM-DD_meu-app/RELATORIO_QA_WEB.json
```

---

## 3. Regras Negativas (Zero Tolerância)

1. **NUNCA** crie mocks ou pule ferramenta ausente no preflight (`Zero Stubs`).
2. **NUNCA** edite o código-fonte do app alvo diretamente após o relatório sem concluir a entrevista socrática e gerar `INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM`.
3. **NUNCA** promova baseline visual do `lost-pixel` sem aprovação humana explícita no Stop-Gate `n09_sg_baseline_visual`.
4. **NUNCA** grave relatórios fora de `docs/qa-web/YYYY-MM-DD_<alvo>/` (Path Jail).

---

## 4. Critério de Parada (Stopping Checklist)

- [ ] `legoloop validate docs/dags/2026-11-15_DAG_WEB_QA.json` retornou exit code 0.
- [ ] `node --test tools/webqa/webqa.test.mjs` retornou exit code 0.
- [ ] Relatório consolidado salvo em `docs/qa-web/YYYY-MM-DD_<alvo>/RELATORIO_QA_WEB.{json,md}`.
- [ ] Entrevista socrática iniciada sobre os achados antes do ciclo `BUILD`.

---

## 5. DAG & Nós Canônicos

| ID do Nó | Nome do Nó | Descrição |
| :--- | :--- | :--- |
| `n01_preflight` | Preflight estrito das 8 ferramentas obrigatórias | Validação e verificação binária de todo o ferramental de QA |
| `n02_boot_e_descoberta` | Boot opcional de app local e descoberta de rotas e viewports | Inicialização do servidor local e mapeamento de rotas |
| `n03_links_lychee` | Auditoria determinística de links e assets com lychee | Verificação de integridade de links internos, externos e assets |
| `n04_obscura_e_playwright` | Inspeção de rotas, console e navegação com obscura e playwright | Checagem de status HTTP e captura de erros de runtime/console |
| `n05_renderlint_e_design` | Auditoria de botões mortos, layout e tokens com renderlint e aidd-linter-design | Análise de layout, controles mortos e conformidade de design |
| `n06_axe_acessibilidade` | Auditoria de acessibilidade WCAG (critical e serious) com axe-core | Testes automatizados de conformidade com padrões de acessibilidade |
| `n07_unlighthouse_rotas` | Varredura Lighthouse multi-rota (nota mínima 90) com unlighthouse | Auditoria de performance, SEO e melhores práticas web |
| `n08_lost_pixel_visual` | Regressão visual open-source e snapshot inicial com lost-pixel | Captura de snapshots e comparação de regressão visual |
| `n10_relatorio_consolidado` | Consolidação de relatório JSON + Markdown em docs/qa-web/ | Unificação dos diagnósticos em relatório JSON e Markdown |
| `n11_entrevista_correcao` | Entrevista socrática de afunilamento sobre correções do relatório | Afunilamento investigativo das falhas encontradas |
| `n12_governanca_build` | Geração da cadeia INTERVIEW -> INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM para BUILD | Materialização do ciclo formal para remediação e build |

