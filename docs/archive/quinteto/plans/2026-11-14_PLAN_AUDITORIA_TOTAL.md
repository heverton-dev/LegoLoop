# PLAN: Auditoria Total, Correção e Homologação LegoLoop

## Tickets de Correção (Paralelos por Domínio)
- **T-01 (`fix-core`):** Corrigir `TestInstalarSkills` (`internal/cli/instalar_skills_test.go`), `gofmt` (`internal/orchestrix/prompt_sanitize_test.go`), e `funclen` (`internal/processo/projecoes.go:64` e `internal/mcp/swarm_tools.go:34`).
- **T-02 (`fix-loops`):** Corrigir os 8 `pipeline.dag.json` inválidos em `processos/` (remover BOM de `matt-pocock` e `mega-pipeline`, corrigir campos desconhecidos em `addy-osmani`, `fechar`, `impeccable-design`, `lauren-tan`, `legoloop-triagem`, e gate de `fronteira-ai-driven`).
- **T-03 (`fix-studio`):** Corrigir bug do Canvas onde `Gatilho` (`trigger_*`) não aceita ligação de nenhuma peça (`internal/studio/static/app.js` / `internal/studio/pieces.go`) e adicionar teste unitário cobrindo `canConnect` com `trigger_*`.
- **T-04 (`fix-swarm-motores`):** Tratar erros ignorados em `internal/cli/swarm_ui.go:48-52` e `internal/runner/triggers.go:141` / `token_budget.go:91`, e adicionar teste unitário para CLI `swarm`.
- **T-05 (`verify-e2e`):** Rodar portões globais (`go test`, `go vet`, `gofmt`, `funclen`, `validate`), recompilar `bin/legoloop.exe`, validar Studio via Playwright (botões, modais, painéis, `Gatilho`) e validar CLI + MCP `mcp-router`.
