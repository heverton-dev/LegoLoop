# Relatório Consolidado Web-QA — webapp

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7878
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
2. **[axe-core]** (serious) `/swarm (desktop)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
3. **[axe-core]** (serious) `/ (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
4. **[axe-core]** (serious) `/swarm (mobile)`: color-contrast: Elements must meet minimum color contrast ratio thresholds (2 elemento(s))
5. **[renderlint]** (critical) `http://127.0.0.1:7878/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 92 finding(s) in 7811ms
  / @1280px light: 44 finding(s) in 7888ms
  / @375px dark: 91 finding(s) in 5858ms
  / @375px light: 107 finding(s) in 11014ms
  / @768px dark: 67 finding(s) in 3923ms
  /inventario @1280px light: 7 finding(s) in 1532ms
  / @1280px dark: 44 finding(s) in 5728ms
  /inventario @768px light: 3 finding(s) in 4001ms
  /inventario @375px light: 4 finding(s) in 9701ms
  /swarm @375px light: 11 finding(s) in 5956ms
  /inventario @375px dark: 5 finding(s) in 990ms
  /swarm @768px light: 11 finding(s) in 6002ms
  /swarm @1280px light: 11 finding(s) in 5920ms
  /swarm @375px dark: 8 finding(s) in 5058ms
  /inventario @1280px dark: 3 finding(s) in 5238ms
  /inventario @768px dark: 3 finding(s) in 5626ms
  /swarm @768px dark: 1 finding(s) in 5558ms
  /swarm @1280px dark: 9 finding(s) in 5029ms
6. **[renderlint]** (critical) `http://127.0.0.1:7878/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @768px light: 3 finding(s) in 4679ms
  /inventario @1280px light: 3 finding(s) in 4831ms
  /inventario @375px dark: 3 finding(s) in 2879ms
  /inventario @375px light: 3 finding(s) in 6656ms
  /inventario @768px dark: 3 finding(s) in 3028ms
  / @768px light: 92 finding(s) in 4952ms
  /swarm @375px light: 5 finding(s) in 1106ms
  / @375px light: 107 finding(s) in 7948ms
  /inventario @1280px dark: 3 finding(s) in 9397ms
  / @1280px light: 44 finding(s) in 8092ms
  / @768px dark: 49 finding(s) in 1311ms
  /swarm @768px light: 11 finding(s) in 5903ms
  / @375px dark: 91 finding(s) in 5867ms
  /swarm @1280px light: 11 finding(s) in 6080ms
  /swarm @768px dark: 3 finding(s) in 887ms
  / @1280px dark: 44 finding(s) in 3963ms
  /swarm @375px dark: 8 finding(s) in 3170ms
  /swarm @1280px dark: 9 finding(s) in 5181ms
7. **[renderlint]** (critical) `http://127.0.0.1:7878/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @375px dark: 8 finding(s) in 10433ms
  /swarm @1280px light: 11 finding(s) in 13159ms
  /swarm @768px light: 11 finding(s) in 13172ms
  / @375px light: 78 finding(s) in 1739ms
  /swarm @375px light: 11 finding(s) in 15402ms
  /swarm @768px dark: 9 finding(s) in 5350ms
  /swarm @1280px dark: 9 finding(s) in 5348ms
  / @768px light: 94 finding(s) in 3677ms
  /inventario @375px light: 9 finding(s) in 1465ms
  /inventario @1280px light: 5 finding(s) in 1816ms
  / @1280px light: 44 finding(s) in 5187ms
  / @1280px dark: 37 finding(s) in 1027ms
  /inventario @768px light: 11 finding(s) in 3599ms
  /inventario @768px dark: 3 finding(s) in 3292ms
  /inventario @375px dark: 4 finding(s) in 3834ms
  / @375px dark: 91 finding(s) in 5841ms
  / @768px dark: 66 finding(s) in 6992ms
  /inventario @1280px dark: 3 finding(s) in 5052ms

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
