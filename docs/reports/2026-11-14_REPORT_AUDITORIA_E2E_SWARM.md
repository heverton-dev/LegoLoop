# REPORT: Auditoria E2E CLI, Studio e Portões (Swarm)

- **Data:** 2026-11-14
- **Commit base:** `eec49b5`
- **Execução:** Swarm paralelo com `workspace:"share"` (`cli-prover`, `studio-prover`, `arch-auditor`).

## 1. Prova Real CLI & MCP STDIO (Hermético)
- `go build -o bin/legoloop.exe ./cmd/legoloop`: exit 0.
- `legoloop listar`: exit 0.
- `legoloop validate` (em `pipeline.dag.json` real): exit 0 (`7 nós válidos, sem ciclos`).
- `legoloop matriz ver` (cofre vazio): exit 1 esperado (`erro: a matriz não existe...`).
- `legoloop stept status`: exit 0 (`MODO=prova`, `PAUSA=destrutivo`).
- `legoloop mcp-router` (STDIO `initialize` + `tools/list`): exit 0, retornou `mcp_search` e `call_mcp_tool`.

## 2. Prova Real LegoLoop Studio + Playwright
- Servidor iniciado em `http://127.0.0.1:18089/` (PID 6948, encerrado no teardown).
- Fluxos validados: Home/Canvas, Processos (`modal-processos`), Peças (`btn-toggle-palette`), modal `+ Novo Processo`.
- Erros: `console=0`, `pageerror=0`, `CSP=0`, `HTTP>=400=0`.
- Evidências: `docs/reports/img/2026-11-14_studio_home.png`, `docs/reports/img/2026-11-14_studio_processos.png`.

## 3. Portões de Qualidade
- `go test -count=1 ./...`: exit 0 (51 pacotes ok).
- `go vet ./...`: exit 0.
- `gofmt -l .`: corrigidos 4 arquivos (`internal/assets/types.go`, `internal/processo/projecoes.go`, `internal/studio/compose.go`, `internal/studio/loops_api_test.go`), saída final vazia.
- `go run ./tools/funclen internal`: exit 0.
- `go run ./tools/taxlinter internal loops skills`: exit 0.
