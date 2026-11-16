# Relatório Consolidado Web-QA — legoloop-studio-v6

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7897
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
2. **[axe-core]** (serious) `/ (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
3. **[renderlint]** (critical) `http://127.0.0.1:7897/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @1280px light: 40 finding(s) in 6671ms
  / @768px light: 76 finding(s) in 6679ms
  / @375px dark: 83 finding(s) in 5488ms
  / @375px light: 99 finding(s) in 10025ms
  / @768px dark: 61 finding(s) in 4036ms
  /inventario @1280px light: 11 finding(s) in 1495ms
  / @1280px dark: 40 finding(s) in 5905ms
  /swarm @768px light: 5 finding(s) in 1267ms
  /inventario @768px light: 10 finding(s) in 5607ms
  /inventario @375px light: 4 finding(s) in 9567ms
  /swarm @375px light: 30 finding(s) in 7971ms
  /swarm @1280px light: 6 finding(s) in 6866ms
  /inventario @375px dark: 3 finding(s) in 5193ms
  /inventario @768px dark: 3 finding(s) in 5128ms
  /inventario @1280px dark: 3 finding(s) in 3281ms
  /swarm @1280px dark: 3 finding(s) in 765ms
  /swarm @375px dark: 13 finding(s) in 5102ms
  /swarm @768px dark: 10 finding(s) in 5135ms
4. **[renderlint]** (critical) `http://127.0.0.1:7897/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @768px light: 3 finding(s) in 6913ms
  /inventario @1280px light: 3 finding(s) in 7042ms
  /inventario @375px dark: 3 finding(s) in 5455ms
  /inventario @375px light: 4 finding(s) in 10376ms
  /inventario @1280px dark: 3 finding(s) in 3494ms
  /inventario @768px dark: 3 finding(s) in 5481ms
  / @375px light: 99 finding(s) in 8104ms
  / @1280px light: 40 finding(s) in 7398ms
  / @768px light: 77 finding(s) in 7892ms
  /swarm @375px light: 6 finding(s) in 7598ms
  /swarm @768px light: 11 finding(s) in 5756ms
  /swarm @1280px light: 6 finding(s) in 5722ms
  / @1280px dark: 44 finding(s) in 1028ms
  / @768px dark: 62 finding(s) in 3856ms
  / @375px dark: 84 finding(s) in 5739ms
  /swarm @768px dark: 6 finding(s) in 841ms
  /swarm @1280px dark: 5 finding(s) in 784ms
  /swarm @375px dark: 13 finding(s) in 5080ms
5. **[renderlint]** (critical) `http://127.0.0.1:7897/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @1280px light: 6 finding(s) in 6508ms
  /swarm @768px light: 11 finding(s) in 6659ms
  /swarm @375px dark: 13 finding(s) in 5171ms
  /swarm @375px light: 30 finding(s) in 8618ms
  /swarm @768px dark: 10 finding(s) in 5407ms
  /swarm @1280px dark: 6 finding(s) in 5458ms
  / @768px light: 76 finding(s) in 4529ms
  / @375px light: 99 finding(s) in 9890ms
  / @1280px light: 41 finding(s) in 6777ms
  /inventario @768px light: 3 finding(s) in 5615ms
  /inventario @375px light: 4 finding(s) in 9279ms
  / @1280px dark: 40 finding(s) in 1010ms
  /inventario @1280px light: 3 finding(s) in 6867ms
  / @768px dark: 62 finding(s) in 6264ms
  / @375px dark: 84 finding(s) in 6339ms
  /inventario @375px dark: 3 finding(s) in 5207ms
  /inventario @768px dark: 3 finding(s) in 5067ms
  /inventario @1280px dark: 3 finding(s) in 5087ms
6. **[unlighthouse]** (serious) `http://127.0.0.1:7897`: [Unlighthouse] ℹ Using system Chrome located at: C:\Program Files\Google\Chrome\Application\chrome.exe.
[Unlighthouse] ✔ Successfully connected to http://127.0.0.1:7897. (Status: 200).
[Unlighthouse] ℹ Sitemap appears to be missing, falling back to crawler mode.
[Unlighthouse] ✔ Unlighthouse has finished scanning http://127.0.0.1:7897: 3 routes in 56s.
[Unlighthouse] ℹ Running score budgets. 90
[Unlighthouse] ✔ Generated jsonSimple report ./.unlighthouse/ci-result.json


[Unlighthouse]  WARN  You seem to be missing a robots.txt.


[Unlighthouse]  ERROR  / has invalid score 0.77 for category performance.

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
