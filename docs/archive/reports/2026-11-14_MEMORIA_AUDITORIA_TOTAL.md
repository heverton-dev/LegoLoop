# Memória Dedicada: Auditoria Total LegoLoop (2026-11-14)

## 1. Contexto e Escopo
Auditoria profunda executada via Swarm nativo em `C:\Users\trcnologia\Desktop\LegoLoop` cobrindo Core Go, Motores, Processos/Peças, Studio e Swarm.

## 2. Achados Consolidados por Domínio

### 2.1 Core Go (`audit-core`, seq #7)
- `go test -count=1 ./...`: FAIL em `internal/cli` (`TestInstalarSkills` esperava 15 funções nativas, obteve 16) e `internal/gate` (`TestSchemaV2AceitaDialetos` falhou em `processos/matt-pocock/pipeline.dag.json` e `processos/mega-pipeline/pipeline.dag.json`).
- `go vet ./...`: PASS (0 warnings).
- `gofmt -l .`: FAIL em `internal/orchestrix/prompt_sanitize_test.go`.
- `go run ./tools/funclen internal`: FAIL com 2 funções >35 linhas: `internal/mcp/swarm_tools.go:34` (`handleSwarmTool`, 59 linhas) e `internal/processo/projecoes.go:64` (`Projetar`, 38 linhas).
- `go run ./tools/taxlinter internal cmd loops skills docs/templates`: FAIL com 883 violações ASD-STE100.

### 2.2 Motores (`audit-motores`, seq #9)
- Todos os 13 pacotes (`runner`, `orchestrix`, `wave`, `fabricar`, `fechar`, `compor`, `refinar`, `remover`, `videomotion`, `videorea`, `legopress`, `redc`, `mcp`) passaram em `go test -count=1` e `go vet`.
- Erros de I/O ignorados (`_ = ...`) em `internal/fechar/fechar.go`, `internal/orchestrix/unit_run.go:419-420`, `internal/videorea/bundle_factory.go:228`, `internal/runner/triggers.go:141` e `internal/runner/token_budget.go:91`.
- 126 violações `taxlinter` nos motores.

### 2.3 Processos e Peças (`audit-processos`, seq #8)
- `docs/templates/...`: PASS.
- `processos/` (19 processos): 11 válidos e 8 falhando em `legoloop validate`:
  - `addy-osmani`, `fechar`, `impeccable-design`, `lauren-tan`: `unknown field ""`.
  - `fronteira-ai-driven`: `nó "T-01": gate: caractere inesperado '.' na posição 8`.
  - `legoloop-triagem`: `unknown field "description"`.
  - `matt-pocock`, `mega-pipeline`: UTF-8 BOM (`\ufeff`).
- 457 violações `taxlinter` em `processos/`, `skills/`, `docs/templates/`.

### 2.4 LegoLoop Studio (`audit-studio`, seq #12)
- `go test -count=1 ./internal/studio/...`: PASS.
- Bug do Canvas (`Gatilho` não aceita ligação de nenhuma peça): `internal/studio/static/app.js:526-532` (`canConnect`) exige `a.piece || b.piece` ou `state.matrix[a.block].includes(b.block)`, mas `internal/studio/pieces.go` não possui entradas `trigger_*` na matriz.
- Colapso de painéis (`#btn-toggle-sidebar`, `#btn-toggle-inspector`), grade 4x2 e CSP estrita (`internal/studio/security.go:17-18,25`) presentes e testados em Go; pendente validação E2E real via Playwright.

### 2.5 Swarm (`audit-swarm`, seq #10)
- `internal/board/...` e `internal/mcp/...`: PASS.
- `internal/cli/... -run "Test.*Swarm"`: sem testes unitários.
- Paridade CLI vs MCP em `internal/mcp/swarm_tools.go` (falta `Parent`, `seqs`, `timeout` configurável e `dir` absoluto) e erros silenciosos em `internal/cli/swarm_ui.go:48-52`.
