# Relatório Consolidado Web-QA — legoloop-studio-v5

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7897
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[playwright]** (critical) `/ (desktop)`: Failed to load resource: the server responded with a status of 429 (Too Many Requests) | Failed to load resource: the server responded with a status of 429 (Too Many Requests) | Failed to load resource: the server responded with a status of 429 (Too Many Requests)
2. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
3. **[axe-core]** (serious) `/swarm (desktop)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
4. **[axe-core]** (serious) `/ (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
5. **[axe-core]** (serious) `/swarm (mobile)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
6. **[renderlint]** (critical) `http://127.0.0.1:7897/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 91 finding(s) in 6752ms
  / @1280px light: 43 finding(s) in 6759ms
  / @1280px dark: 45 finding(s) in 1103ms
  / @768px dark: 53 finding(s) in 1200ms
  / @375px dark: 90 finding(s) in 5786ms
  /inventario @375px light: 9 finding(s) in 1323ms
  /inventario @768px light: 7 finding(s) in 1436ms
  / @375px light: 106 finding(s) in 10224ms
  /inventario @1280px light: 3 finding(s) in 6021ms
  /swarm @768px light: 11 finding(s) in 6027ms
  /swarm @1280px light: 11 finding(s) in 5929ms
  /swarm @375px light: 11 finding(s) in 7945ms
  /inventario @375px dark: 3 finding(s) in 5211ms
  /inventario @768px dark: 3 finding(s) in 5220ms
  /inventario @1280px dark: 3 finding(s) in 5213ms
  /swarm @375px dark: 8 finding(s) in 5130ms
  /swarm @768px dark: 9 finding(s) in 5106ms
  /swarm @1280px dark: 9 finding(s) in 5098ms
7. **[renderlint]** (critical) `http://127.0.0.1:7897/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @375px dark: 3 finding(s) in 5202ms
  /inventario @768px light: 3 finding(s) in 6052ms
  /inventario @1280px light: 3 finding(s) in 6200ms
  /inventario @375px light: 5 finding(s) in 7559ms
  / @768px light: 93 finding(s) in 3001ms
  /inventario @1280px dark: 3 finding(s) in 5315ms
  /inventario @768px dark: 3 finding(s) in 5352ms
  / @375px light: 106 finding(s) in 8388ms
  /swarm @768px light: 11 finding(s) in 6002ms
  /swarm @375px light: 11 finding(s) in 7951ms
  /swarm @1280px light: 11 finding(s) in 5764ms
  / @375px dark: 91 finding(s) in 5678ms
  / @768px dark: 67 finding(s) in 5961ms
  /swarm @375px dark: 8 finding(s) in 5105ms
  /swarm @768px dark: 9 finding(s) in 5104ms
  / @1280px dark: 44 finding(s) in 10281ms
  / @1280px light: 44 finding(s) in 21167ms
  /swarm @1280px dark: 9 finding(s) in 5153ms
8. **[renderlint]** (critical) `http://127.0.0.1:7897/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @1280px light: 11 finding(s) in 6876ms
  /swarm @768px light: 11 finding(s) in 6947ms
  /swarm @375px dark: 8 finding(s) in 5144ms
  /swarm @375px light: 11 finding(s) in 8894ms
  /swarm @768px dark: 9 finding(s) in 5220ms
  /swarm @1280px dark: 9 finding(s) in 5215ms
  / @375px light: 106 finding(s) in 10201ms
  / @1280px light: 44 finding(s) in 6938ms
  /inventario @375px light: 4 finding(s) in 9362ms
  /inventario @768px light: 3 finding(s) in 6207ms
  /inventario @1280px light: 3 finding(s) in 5998ms
  / @375px dark: 91 finding(s) in 6692ms
  / @1280px dark: 44 finding(s) in 5878ms
  / @768px dark: 67 finding(s) in 6575ms
  / @768px light: 92 finding(s) in 23484ms
  /inventario @375px dark: 3 finding(s) in 5225ms
  /inventario @1280px dark: 3 finding(s) in 5474ms
  /inventario @768px dark: 3 finding(s) in 6169ms

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
