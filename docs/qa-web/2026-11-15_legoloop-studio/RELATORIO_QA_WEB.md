# Relatório Consolidado Web-QA — legoloop-studio

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7892
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[axe-core]** (critical) `/ (desktop)`: aria-required-children: Certain ARIA roles must contain particular children (1 elemento(s))
2. **[axe-core]** (serious) `/ (desktop)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
3. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (6 elemento(s))
4. **[axe-core]** (serious) `/inventario (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
5. **[axe-core]** (serious) `/swarm (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
6. **[axe-core]** (critical) `/ (mobile)`: aria-required-children: Certain ARIA roles must contain particular children (1 elemento(s))
7. **[axe-core]** (serious) `/ (mobile)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
8. **[axe-core]** (serious) `/ (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
9. **[axe-core]** (serious) `/inventario (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
10. **[axe-core]** (serious) `/swarm (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
11. **[renderlint]** (critical) `http://127.0.0.1:7892`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 94 finding(s) in 7692ms
  / @1280px light: 60 finding(s) in 7754ms
  / @375px dark: 54 finding(s) in 4210ms
  / @375px light: 63 finding(s) in 10228ms
  / @768px dark: 78 finding(s) in 4975ms
  / @1280px dark: 59 finding(s) in 5316ms
  /inventario @768px light: 60 finding(s) in 4957ms
  /inventario @375px light: 33 finding(s) in 7392ms
  /inventario @1280px light: 73 finding(s) in 5718ms
  data stress skipped on http://127.0.0.1:7892/swarm: page.reload: net::ERR_CONNECTION_REFUSED
  /swarm @768px light: 13 finding(s) in 3345ms
  /inventario @375px dark: 0 finding(s) in 237ms
  data stress skipped on http://127.0.0.1:7892/swarm: page.reload: net::ERR_CONNECTION_REFUSED
  /swarm @1280px light: 13 finding(s) in 3403ms
  /inventario @768px dark: 0 finding(s) in 163ms
  /inventario @1280px dark: 0 finding(s) in 118ms
  /swarm @375px dark: 0 finding(s) in 157ms
  /swarm @768px dark: 0 finding(s) in 129ms
  /swarm @1280px dark: 0 finding(s) in 97ms
  /swarm @375px light: 19 finding(s) in 6325ms

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
