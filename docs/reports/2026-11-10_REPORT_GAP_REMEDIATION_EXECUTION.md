# REPORT: Execução da Remediação de Gaps (TICK-GAP-01 a TICK-GAP-06)

Data: 2026-11-10

## 1. Resultado por ticket

| Ticket | Estado | Prova |
| :--- | :---: | :--- |
| TICK-GAP-01 | Feito | `TestDownloadMediaAndMetadata`, `TestDownloadMediaAndMetadataFailures`; mutação (remover `--extract-audio`) falha 4 testes |
| TICK-GAP-02 | Feito | `TestDossierValidateStrict`, `TestVideoReaCLI_ForensicsRejeitaDossieSemAchados`; mutação (`ValidateStrict` devolve nil) falha 4 testes; mutação no CLI falha o teste novo |
| TICK-GAP-03 | Feito | `TestBundlePathJailSecurity`; mutação (remover a checagem de separador) falha 4 testes |
| TICK-GAP-04 | Feito | `go run ./tools/funclen internal` sai com 0 (as 22 funções refatoradas) |
| TICK-GAP-05 | Feito no manual | Linha 66 do manual usa `impeccable-finish-reviewer` |
| TICK-GAP-06 | Feito | Seção `## Fallback Gracioso` documenta os dois modos |

## 2. Mudanças relevantes

- `internal/videorea/crawler_download.go` (novo): `Runner`, `ExecRunner`, `DownloadMediaAndMetadata`, `DownloadWithRunner`. Acrescentei `--no-simulate` aos flags pedidos. Sem ele o `yt-dlp` só imprime o JSON e não baixa o áudio, pois `--dump-json` implica simulação.
- `internal/cli/video_rea.go`: `runIngest` chama o download real e sai com 1 em falha. URL `file://` e `local://` não passam pelo `yt-dlp`. `Deps.VideoRunner` injeta o executor nos testes. `executeForensics` chama `ValidateStrict` e sai com 1.
- `internal/videorea/bundle.go`: `SanitizeBundlePath` e `NewTargetBundle` devolvem erro; `bundle_factory.go` propaga o erro.
- `loops/legoloop-video-rea/pipeline.dag.json`, o `.loop.json` e o espelho em `docs/dags/`: todos os caminhos `scratch/videorea/` viraram `docs/videorea/${SLUG}/scratch/`.
- Mudança extra e necessária: o runner não expandia `${VAR}` nos gates. Adicionei `gate.Env.Vars`, `runner.Options.Vars` e expansão em `file_exists`. Teste: `TestFileExistsExpandsVars` (mutação falha 4 testes). O `legoloop run` precisa receber `--var SLUG=<slug>`.
- `internal/motion`: as interfaces `frameEncoder`, `frameCapturer` e `beatDetector` estavam definidas só em arquivos `_test.go`, e `go build ./...` e `go vet ./...` falhavam. Movi as três para os arquivos de produção.

## 3. Desvios do ticket (declarados)

1. GAP-04, `BuiltinTemplates`: `//go:embed` não alcança `loops/` (fora do pacote) e o formato de `Composition` difere do `.loop.json`. Os templates agora vivem em `internal/studio/templates.json` (gerado a partir do conteúdo anterior; a saída é idêntica byte a byte). A função tem 8 linhas.
2. GAP-05: `grep -r impeccabe docs/` ainda acha 5 arquivos. Todos citam o typo como objeto (tabela de erros do livro Impeccable e relatórios de auditoria). Corrigir esses textos apagaria o registro da correção. O manual está limpo.
3. A segunda checagem de prefixo em `SanitizeBundlePath` é defesa em profundidade: a mutação dela não falha nenhum teste, porque as checagens anteriores já rejeitam todo caso.
4. O `yt-dlp` real não foi executado neste ambiente. O teste usa executor falso. O caminho real (`ExecRunner`) não tem prova de execução.

## 4. Portão de qualidade

| Comando | Resultado |
| :--- | :--- |
| `go test -count=1 ./internal/videorea/...` | ok |
| `go test -count=1 ./internal/studio/...` | ok |
| `go test -count=1 ./internal/canon/...` | ok |
| `go vet ./...` | limpo |
| `go run ./tools/funclen internal` | exit 0 |
| `gofmt -l .` | vazio |
| `go build -o bin/legoloop.exe ./cmd/legoloop` | ok |
| `git diff -- go.mod go.sum` | vazio (só stdlib) |
| `go test -count=1 ./...` | 38 pacotes ok, nenhuma falha |

## 5. Pendências

- A primeira execução da suíte completa estourou o tempo no `internal/cli`: meu teste novo não injetava `Discover` e sondava a rede real. Corrigi com `noModels`; o subconjunto `VideoRea` caiu de 63 s para 1,2 s.
- `TestValidateText` (`internal/taxonomy`) falhou uma vez numa execução anterior, enquanto outro agente editava `taxonomy.go` e `taxlinter/main.go`. Na execução final passou.
- O `yt-dlp` real não foi executado (ver seção 3, item 4).
- Sem commit e sem push: você não pediu.
