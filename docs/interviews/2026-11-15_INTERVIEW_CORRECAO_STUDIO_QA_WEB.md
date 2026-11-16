---
id: INTERVIEW-intv_1794763194418751600
status: APPROVED
date: 2026-11-15
schema_version: 2
upstream: []
covers: [Q1, Q2]
---

# INTERVIEW: Correção dos Achados QA Web do LegoLoop Studio v2

> **Protocolo:** uma pergunta por vez, em prosa. Entrevista conduzida via `legoloop web-qa corrigir --relatorio docs/qa-web/2026-11-15_legoloop-studio-v2/RELATORIO_QA_WEB.json --model gemini-4-hydrogen` (`sessao: intv_1794763194418751600`).
> **Regra de prova:** nenhuma pergunta nasce de inferência. Cada fato traz prova verificável em arquivo e linha ou comando e saída.
> **Confiança final:** 88% (aprovada pelo usuário para consolidação em INTENT -> SPEC -> PLAN -> ROUTE -> DAG -> SWARM -> BUILD)
> **Iniciativa:** CORRECAO_STUDIO_QA_WEB
> **Template:** TEMPLATE_INTERVIEW.md v2.3.0

## 1. Respostas

| ID | Pergunta | Fatos | Decisão do usuário |
| :--- | :--- | :--- | :--- |
| Q1 | Demanda cobre todos 4 eixos (crash/recusa servidor 7893, `axe-core`, layout `renderlint` e 328 tokens CSS), ou `ERR_CONNECTION_REFUSED` foi instabilidade externa de ambiente e devemos focar só código frontend/acessibilidade/tokens? | F1, F2, F3, F4 | Aceitar palpite: cobrir todos os 4 eixos em `internal/studio/` (concorrência em `assets_api.go:246`, ARIA crítica, layout `renderlint` e extração dos 328 tokens CSS para `tokens.css`). |
| Q2 | Critério de sucesso exige zero achados nas 4 ferramentas QA mantendo paleta visual idêntica light/dark em `tokens.css`, e algum arquivo HTML/JS/Go além de `assets_api.go` e `static/` está proibido mexer? | F1, F3, F4 | Aceitar palpite: zero achados críticos/sérios no QA Web mantendo aparência light/dark via `tokens.css`, com permissão para ajustar templates HTML/JS/Go em `internal/studio/`. |

## 2. Confiança por rodada

| Rodada | Perguntas feitas | Confiança | Motivo da variação |
| :--- | :--- | :--- | :--- |
| 1 | Q1 | 65% | Escopo dos 4 eixos definido com base em `RELATORIO_QA_WEB.json` e stack trace da porta 7894. |
| 2 | Q2 | 88% | Critérios binários de aceite e fronteira em `internal/studio/` confirmados pelo usuário. |

## 3. Fatos medidos durante a entrevista

| ID | Afirmação | Prova (arquivo:linha ou comando) | Resultado |
| :--- | :--- | :--- | :--- |
| F1 | Relatório QA Web registrou 338 achados (`playwright` rotas=5/console=1, `axe` critical=1, `renderlint` violations=3, `aidd-linter-design` violations=328). | `docs/qa-web/2026-11-15_legoloop-studio-v2/RELATORIO_QA_WEB.json:24-47` | VERIFICADO |
| F2 | Studio encerrou sob requisições paralelas com `fatal error: concurrent map writes` em `(*Server).present` ao gravar `it.Meta["tier"]`. | `internal/studio/assets_api.go:246` e `task-592` (`STUDIO_EXITED=True`) | VERIFICADO |
| F3 | Rotas `/` possuem `role="tablist"` (`#palette-cats`) e `role="listbox"` (`#palette-list`) causando `aria-required-children`. | `internal/studio/static/index.html:95-99` e `internal/studio/static/palette.js:10` | VERIFICADO |
| F4 | Não existe `tokens.css` em `internal/studio/static/`; variáveis CSS e cores literais estão misturadas em `style.css` e `swarm.css`. | `internal/studio/static/style.css:18-37`, `internal/studio/static/swarm.css:62-66` | VERIFICADO |

## 4. Pendências para a SPEC

* Especificar a clonagem segura de `it.Meta` com `maps.Clone` em `internal/studio/assets_api.go:240-255` e teste unitário concorrente.
* Especificar correção de `aria-required-children` em `index.html` e `palette.js`, além das falhas de layout do `renderlint` em `/`, `/inventario` e `/swarm`.
* Especificar extração completa do `:root` e `[data-theme="dark"]` para `internal/studio/static/tokens.css` zerando violações no `aidd-linter-design`.
