# PARTE IV — Fechamento Universal e Apêndices

## 1. Cadeia Canônica de Fechamento (`TICK-07.2`)

Para impedir entregas parciais ou divergentes entre código e documentação, a Etapa 6 consolida nove vínculos obrigatórios:

1. Manual operacional em `docs/manuals/2026-11-14_MANUAL_FUNCAO_HOMOLOGACAO_ECOSISTEMA_LEGOLOOP.md`.
2. Livro técnico compilado em Markdown e PDF (`170x240mm`) e auditado por `legoloop legopress auditar`.
3. Skill `legoloop-homologacao` distribuída para o cofre e harnesses via `legoloop instalar-skills`.
4. Quinteto ISO aprovado (`APPROVED`) com paridade SHA-256 e DAG validado por `legoloop validate`.
5. Binário recompilado em `bin/legoloop.exe` e instalado no PATH com `go install ./cmd/legoloop`.
6. Documentação viva atualizada (`AGENTS.md`, `MEMORY.md` e `README.md`).
7. Grafo semântico reindexado em `codebase-memory-mcp`.
8. Suíte completa executada após a última edição (`go test`, `-tags e2e`, `go vet`, `funclen`, `taxlinter`, `gofmt -l .`, `go.mod` limpo).
9. Relatório final com fatos observados em `docs/reports/2026-11-14_REPORT_HOMOLOGACAO_ECOSISTEMA.md`.

## 2. Apêndice A — O que Permanece Declarado como Não Verificado Localmente

Seguindo o protocolo de honestidade técnica do LegoLoop:

- **Detector de corrida com `-race`:** não executado no Windows por ausência de toolchain CGO ativa; permanece exigido no CI Linux.
- **Observações não bloqueantes da Rodada 2:** texto de ajuda resumido em `printSwarmUsage` (`internal/cli/swarm.go`) frente às flags detalhadas nos subcomandos, e rótulo padrão de fallback em `lego-item-modal.js` quando `suggested_model` está ausente.

## 3. Apêndice B — Comandos Rápidos de Auditoria Contínua

| Objetivo                          | Comando                                        |
| :-------------------------------- | :--------------------------------------------- |
| Recompilar binário local e PATH   | `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop` |
| Rodar testes unitários e E2E      | `go test -count=1 ./...` e `go test -count=1 -tags e2e ./internal/e2e/...` |
| Auditar tamanho de funções        | `go run ./tools/funclen internal cmd tools`      |
| Auditar taxonomia STE             | `go run ./tools/taxlinter internal cmd skills docs/templates` |
| Auditar formatação e go.mod       | `gofmt -l .` e `git diff --exit-code -- go.mod go.sum` |
