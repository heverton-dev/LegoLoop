# Relatório Consolidado Web-QA — legoloop-studio-v7

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7897
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
2. **[axe-core]** (serious) `/ (mobile)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (1 elemento(s))
3. **[renderlint]** (critical) `http://127.0.0.1:7897/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 36 finding(s) in 7043ms
  / @1280px light: 39 finding(s) in 7235ms
  / @375px dark: 36 finding(s) in 5765ms
  / @375px light: 36 finding(s) in 10509ms
  / @768px dark: 37 finding(s) in 5579ms
  / @1280px dark: 40 finding(s) in 6700ms
  /inventario @1280px light: 7 finding(s) in 1546ms
  /inventario @768px light: 3 finding(s) in 6915ms
  /inventario @375px light: 4 finding(s) in 9216ms
  /inventario @375px dark: 11 finding(s) in 781ms
  /swarm @768px light: 1 finding(s) in 5603ms
  /swarm @375px light: 7 finding(s) in 8022ms
  /swarm @1280px light: 6 finding(s) in 5644ms
  /inventario @1280px dark: 3 finding(s) in 5272ms
  /inventario @768px dark: 3 finding(s) in 6058ms
  /swarm @375px dark: 7 finding(s) in 5054ms
  /swarm @768px dark: 6 finding(s) in 5078ms
  /swarm @1280px dark: 6 finding(s) in 5035ms
4. **[renderlint]** (critical) `http://127.0.0.1:7897/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @768px light: 3 finding(s) in 6814ms
  /inventario @1280px light: 3 finding(s) in 7002ms
  /inventario @375px dark: 3 finding(s) in 5434ms
  /inventario @375px light: 4 finding(s) in 10245ms
  /inventario @768px dark: 3 finding(s) in 5269ms
  /inventario @1280px dark: 3 finding(s) in 5335ms
  / @768px light: 36 finding(s) in 3046ms
  /swarm @768px light: 6 finding(s) in 5683ms
  / @1280px light: 40 finding(s) in 7013ms
  / @375px dark: 32 finding(s) in 979ms
  /swarm @375px light: 8 finding(s) in 7775ms
  / @375px light: 37 finding(s) in 11439ms
  /swarm @375px dark: 3 finding(s) in 732ms
  /swarm @1280px light: 6 finding(s) in 5610ms
  / @1280px dark: 41 finding(s) in 5622ms
  / @768px dark: 38 finding(s) in 5663ms
  /swarm @768px dark: 6 finding(s) in 5011ms
  /swarm @1280px dark: 6 finding(s) in 5078ms
5. **[renderlint]** (critical) `http://127.0.0.1:7897/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @1280px light: 6 finding(s) in 5857ms
  /swarm @768px light: 6 finding(s) in 5880ms
  /swarm @375px dark: 7 finding(s) in 5172ms
  /swarm @375px light: 8 finding(s) in 7736ms
  /swarm @768px dark: 6 finding(s) in 5257ms
  /swarm @1280px dark: 6 finding(s) in 5296ms
  / @768px light: 39 finding(s) in 4779ms
  / @1280px light: 27 finding(s) in 2752ms
  / @375px light: 36 finding(s) in 9692ms
  /inventario @1280px light: 3 finding(s) in 3944ms
  /inventario @768px light: 3 finding(s) in 5746ms
  / @375px dark: 24 finding(s) in 1078ms
  /inventario @375px dark: 9 finding(s) in 813ms
  /inventario @375px light: 4 finding(s) in 9597ms
  / @768px dark: 38 finding(s) in 3631ms
  /inventario @768px dark: 3 finding(s) in 5564ms
  / @1280px dark: 41 finding(s) in 6567ms
  /inventario @1280px dark: 3 finding(s) in 5107ms

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
