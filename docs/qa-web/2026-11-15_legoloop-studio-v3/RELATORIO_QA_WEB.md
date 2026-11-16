# Relatório Consolidado Web-QA — legoloop-studio-v3

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7878
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[playwright]** (critical) `/ (desktop)`: Failed to load resource: net::ERR_CONNECTION_REFUSED
2. **[axe-core]** (critical) `/ (desktop)`: aria-required-children: Certain ARIA roles must contain particular children (1 elemento(s))
3. **[axe-core]** (serious) `/ (desktop)`: label-content-name-mismatch: Elements must have their visible text as part of their accessible name (2 elemento(s))
4. **[playwright]** (critical) `/inventario (desktop)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/inventario
Call log:
[2m  - navigating to "http://127.0.0.1:7878/inventario", waiting until "networkidle"[22m

5. **[playwright]** (critical) `/swarm (desktop)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/swarm
Call log:
[2m  - navigating to "http://127.0.0.1:7878/swarm", waiting until "networkidle"[22m

6. **[playwright]** (critical) `/ (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/
Call log:
[2m  - navigating to "http://127.0.0.1:7878/", waiting until "networkidle"[22m

7. **[playwright]** (critical) `/inventario (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/inventario
Call log:
[2m  - navigating to "http://127.0.0.1:7878/inventario", waiting until "networkidle"[22m

8. **[playwright]** (critical) `/swarm (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/swarm
Call log:
[2m  - navigating to "http://127.0.0.1:7878/swarm", waiting until "networkidle"[22m

9. **[renderlint]** (critical) `http://127.0.0.1:7878/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 92 finding(s) in 16222ms
  / @1280px light: 44 finding(s) in 16312ms
  / @375px dark: 91 finding(s) in 6955ms
  / @375px light: 107 finding(s) in 19726ms
  / @768px dark: 67 finding(s) in 7189ms
  / @1280px dark: 50 finding(s) in 7120ms
  /inventario @768px light: 3 finding(s) in 7109ms
  /inventario @375px light: 4 finding(s) in 10746ms
  /inventario @1280px light: 3 finding(s) in 5041ms
  /swarm @1280px light: 0 finding(s) in 117ms
  /swarm @768px light: 5 finding(s) in 1724ms
  /inventario @375px dark: 0 finding(s) in 139ms
  /inventario @768px dark: 0 finding(s) in 177ms
  /inventario @1280px dark: 0 finding(s) in 144ms
  /swarm @375px dark: 0 finding(s) in 171ms
  /swarm @768px dark: 0 finding(s) in 141ms
  /swarm @1280px dark: 0 finding(s) in 145ms
  loading-state check skipped on http://127.0.0.1:7878/swarm: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7878/swarm
  /swarm @375px light: 10 finding(s) in 6213ms
10. **[renderlint]** (critical) `http://127.0.0.1:7878/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @768px light: 3 finding(s) in 7631ms
  /inventario @1280px light: 3 finding(s) in 7829ms
  /inventario @375px dark: 29 finding(s) in 4636ms
  /inventario @375px light: 4 finding(s) in 10917ms
  /inventario @768px dark: 3 finding(s) in 5812ms
  /inventario @1280px dark: 3 finding(s) in 5816ms
  / @768px light: 92 finding(s) in 6670ms
  / @375px light: 107 finding(s) in 11979ms
  / @1280px light: 44 finding(s) in 9008ms
  /swarm @375px light: 11 finding(s) in 8856ms
  /swarm @768px light: 11 finding(s) in 6521ms
  /swarm @1280px light: 1 finding(s) in 6714ms
  / @768px dark: 67 finding(s) in 5968ms
  / @1280px dark: 29 finding(s) in 4372ms
  / @375px dark: 91 finding(s) in 6029ms
  /swarm @768px dark: 0 finding(s) in 131ms
  /swarm @1280px dark: 0 finding(s) in 145ms
  data stress skipped on http://127.0.0.1:7878/swarm: page.reload: net::ERR_CONNECTION_REFUSED
  /swarm @375px dark: 8 finding(s) in 3104ms
11. **[renderlint]** (critical) `http://127.0.0.1:7878/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @1280px light: 11 finding(s) in 12215ms
  /swarm @768px light: 11 finding(s) in 12234ms
  /swarm @375px dark: 8 finding(s) in 5608ms
  /swarm @375px light: 11 finding(s) in 14260ms
  / @768px light: 77 finding(s) in 3424ms
  /swarm @1280px dark: 9 finding(s) in 5863ms
  /swarm @768px dark: 9 finding(s) in 6025ms
  /inventario @768px light: 3 finding(s) in 6605ms
  / @375px light: 107 finding(s) in 12050ms
  / @1280px light: 44 finding(s) in 8885ms
  / @768px dark: 2 finding(s) in 1051ms
  / @1280px dark: 0 finding(s) in 103ms
  /inventario @375px light: 4 finding(s) in 9691ms
  /inventario @375px dark: 0 finding(s) in 118ms
  /inventario @768px dark: 0 finding(s) in 132ms
  /inventario @1280px dark: 0 finding(s) in 147ms
  data stress skipped on http://127.0.0.1:7878/: page.reload: net::ERR_CONNECTION_REFUSED
  / @375px dark: 51 finding(s) in 2034ms
  data stress skipped on http://127.0.0.1:7878/inventario: page.reload: net::ERR_CONNECTION_REFUSED
  /inventario @1280px light: 3 finding(s) in 3942ms

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
