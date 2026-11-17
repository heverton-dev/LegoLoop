# SPEC: Auditoria E2E CLI, Studio e Portões (Swarm)

- **Data:** 2026-11-14
- **Slug:** `auditoria-e2e-swarm`
- **Critérios de Aceitação (EARS):**
  1. Quando `bin/legoloop.exe` for compilado, o comando `go build -o bin/legoloop.exe ./cmd/legoloop` deve retornar exit code 0.
  2. Quando os comandos CLI (`listar`, `validate`, `matriz ver`, `stept status`, `mcp-router`) rodarem em ambiente hermético, o comportamento e os códigos de saída devem ser validados sem tocar `~/.aidd` real.
  3. Quando o LegoLoop Studio rodar em porta dedicada, o Playwright deve validar Canvas, Processos, Peças e modais com zero erros de console, CSP ou HTTP >= 400.
  4. Todos os portões (`go test -count=1 ./...`, `go vet ./...`, `gofmt -l .`, `funclen`, `taxlinter`) devem passar com exit code 0 e saída limpa.
