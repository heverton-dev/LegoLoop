# 🐝 Plano de Swarm para BUILD: Correção Studio QA Web

> **Iniciativa:** CORRECAO_STUDIO_QA_WEB  
> **DAG Origem:** [2026-11-15_DAG_CORRECAO_STUDIO_QA_WEB.json](../dags/2026-11-15_DAG_CORRECAO_STUDIO_QA_WEB.json)  
> **Modelo Swarm:** Nativo (workers isolados por fronteira estrita de arquivos, sem colisão)

## Topologia de Execução

| Lote | Worker Swarm | Ticket | Fronteira Exclusiva de Escrita | Critério de Saída |
| :--- | :--- | :--- | :--- | :--- |
| Lote 1 (Paralelo) | `studio-concurrency` | `TICK-01.1` | `internal/studio/assets_api.go`, `internal/studio/assets_api_test.go` | `go test -count=1 ./internal/studio/...` passa com teste de 50 goroutines concorrentes |
| Lote 1 (Paralelo) | `studio-a11y-layout` | `TICK-02.1` | `internal/studio/static/*.html`, `palette.js`, `app.js`, `components/*.js` | Zero erros críticos/sérios no `axe-core` e `renderlint` |
| Lote 2 (Sequencial pós-L1) | `studio-tokens` | `TICK-03.1` | `internal/studio/static/tokens.css`, `style.css`, `swarm.css` + link CSS nos HTMLs | `aidd-linter-design` retorna `0` violações preservando light/dark |
| Lote 3 (Join Barrier) | `studio-verifier` | `TICK-04.1` | Leitura + compilação `bin/legoloop.exe` + execução real | `./bin/legoloop.exe web-qa run --project . --port 7895 --alvo legoloop-studio-v3` retorna `PASSED` (`exit 0`) |
