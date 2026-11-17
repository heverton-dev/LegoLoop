# REPORT E06 — Rename `loops/` → `processos/` (TICK-E06)

**Data:** 2026-11-17 · **Worker:** e06-rename · **Repositório:** LegoLoop

## Resultado

Escrita canônica migrou para `processos/` com leitura dual preservando `loops/` legado.

## Fatos comprovados

| Fato | Prova |
| :--- | :--- |
| Escrita sempre em `processos/` | `internal/processo/persistencia.go` `CaminhoNoCofre` grava `vault/processos/` |
| Leitura dual | `caminhos.DirsLeitura` (união `processos/` + `loops/`) em `internal/processo/cofre.go` |
| AC-E06.2 `fechar` grava em `processos/<slug>/` | `internal/fechar/quarteto_test.go` `TestGravaEmProcessos` GREEN |
| Referências sem órfão `loops/` | `internal/caminhos/caminhos_test.go` `TestReferenciasProcessos` GREEN |
| Compat de leitura do legado | `TestGravaEmProcessos`/`TestLoopsLegadoExitDois` cobrem `loops/` |
| Hook de caminho legado exit 2 | `internal/cli/legado.go` `TestLoopsLegadoExitDois` GREEN |
| Mutações provadas (Green-Red-Green) | 3/3: `mutacao_tmp.txt` FAIL, revert `resolverDirProcesso` FAIL, hook `caminhoLegado` FAIL |
| Gates mecânicos | `go test ./...` 0 · `go vet` 0 · `gofmt -l` vazio · `funclen` baseline 3 · `taxlinter` baseline 9 arquivos · `go.mod/go.sum` intactos |
| Build fresco | `go build -o bin/legoloop.exe` + `go install` exit 0 |
| `legoloop validate` 26/26 | `bin\legoloop.exe validate processos/<slug>/pipeline.dag.json` TOTAL=26 FALHAS=0 |

## Solução passo a passo

1. Produção: escrita canônica em `processos/` (fabricar, scaffold, peca_unica, textos, migrar, triagem, refinar).
2. Leitura dual: `CaminhoNoCofreLeitura`, `DirAcervoLeitura`, `DirsLeitura`, detecção v1 `.loop.json` nos dois dialetos.
3. Fixtures de teste migradas para `DirProcessos`; legado mantido onde testa compat.
4. Paridade de assinaturas: `route_sha256` em 3 DAGs espelho atualizados.
5. TDD: 3 testes novos + 3 mutações provadas e restauradas.

## Sugestão técnica

Próximo épico pode normalizar a dívida baseline de `taxlinter` (9 arquivos) e `funclen` (3 estouros), hoje fora do escopo E06.

## Stop-Gate Humano

SG-E06: [1] seguir | [2] personalizar | [0] cancelar

## NÃO_PROVADO / Pendências

- `docs/templates/templates_test.go:31,185` ainda usa `paths.DirAntigo`/`DirAntigoEm` — intocável por regra; migração pendente.
- Dívida baseline `taxlinter` (9 arquivos com termos proibidos) e `funclen` (3 estouros: `internal/board/swarm_run.go` `runLevel`, `internal/cli/swarm_run.go` `cmdSwarmRun`, `internal/gate/cites.go` `CiteVerdicts`) — pertencem a outros épicos.
- `-race` não rodado (exige CGO/CI Linux).
- Fluxo com LLM real não exercitado (testes herméticos usam LLM falso).
