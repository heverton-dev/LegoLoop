# SPEC: Auditoria Total, Correção e Homologação LegoLoop

## Critérios de Aceitação (EARS)
1. **Quando** `go test -count=1 ./...` for executado, **então** 100% dos pacotes devem passar (incluindo `TestInstalarSkills` e `TestSchemaV2AceitaDialetos`).
2. **Quando** `go vet ./...`, `gofmt -l .` e `go run ./tools/funclen internal` forem executados, **então** todos devem retornar saída limpa e exit code 0.
3. **Quando** `go run ./cmd/legoloop validate` for executado nos 19 processos em `loops/`, **então** todos os 19 `pipeline.dag.json` devem validar sem erro.
4. **Quando** um bloco `Gatilho` (`trigger_*`) for conectado a qualquer peça ou bloco no Canvas do Studio, **então** `canConnect` deve permitir a ligação e renderizar a aresta sem erros de console.
5. **Quando** o Studio for inspecionado via Playwright, **então** botões, modais, grade 4x2, `#btn-toggle-sidebar`, `#btn-toggle-inspector` e conexão do `Gatilho` devem operar com zero erros de console.
6. **Quando** `bin/legoloop.exe` for recompilado, **então** a criação/execução de processo nas 3 vertentes (CLI, MCP `mcp-router` e Studio) deve ser comprovada.
