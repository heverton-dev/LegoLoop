# 🗺️ PLAN: legoloop-web-qa

> **Status:** Alinhado à SPEC v1.1.0  
> **Template:** TEMPLATE_PLAN.md v2.2.0  

---

## Épico 01 — Processo Canônico de Auditoria Web-QA e Encadeamento até BUILD

| Ticket | Nó DAG | Camada | Descrição | Gate Mecânico |
| :--- | :--- | :---: | :--- | :--- |
| `TICK-01.1` | `n01_preflight` | 1 | Validar presença das 8 ferramentas obrigatórias via `webqa.mjs preflight` | `exit_code == 0` |
| `TICK-01.2` | `n02_boot_e_descoberta` | 2 | Subir app local opcional e descobrir rotas/viewports (`1440x900`, `390x844`) | `file_exists('scratch/webqa/routes_manifest.json')` |
| `TICK-01.3` | `n03_links_lychee` | 3 | Auditar links e assets com `lychee` | `exit_code == 0` |
| `TICK-01.4` | `n04_obscura_e_playwright` | 3 | Validar render, rotas e console com `obscura` + `playwright` | `exit_code == 0` |
| `TICK-01.5` | `n05_renderlint_e_design` | 3 | Auditar botões mortos, layout e design tokens | `exit_code == 0` |
| `TICK-01.6` | `n06_axe_acessibilidade` | 3 | Auditar WCAG (`critical`/`serious` = 0) com `axe-core` | `exit_code == 0` |
| `TICK-01.7` | `n07_unlighthouse_rotas` | 3 | Auditar todas as rotas com `unlighthouse` (nota mínima >= 90) | `exit_code == 0` |
| `TICK-01.8` | `n08_lost_pixel_visual` | 4 | Executar regressão visual / snapshot inicial com `lost-pixel` | `file_exists('scratch/webqa/lost_pixel.json')` |
| `TICK-01.9` | `n09_sg_baseline_visual` | 5 | Stop-Gate humano de aprovação da baseline visual | `prompt(...)` |
| `TICK-01.10` | `n10_relatorio_consolidado` | 6 | Consolidar `RELATORIO_QA_WEB.{json,md}` em `docs/qa-web/` | `file_exists('docs/qa-web/RELATORIO_QA_WEB.md')` |
| `TICK-01.11` | `n11_entrevista_correcao` | 7 | Conduzir entrevista socrática de afunilamento sobre correções | `file_exists('scratch/webqa/INTERVIEW_CORRECAO.md')` |
| `TICK-01.12` | `n12_governanca_build` | 8 | Gerar pacote `INTERVIEW..SWARM` para ciclo `BUILD` | `file_exists('docs/qa-web/PACOTE_BUILD_SWARM.json')` |
