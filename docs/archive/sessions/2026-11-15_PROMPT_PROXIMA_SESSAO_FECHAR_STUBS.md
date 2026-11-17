Você trabalha no repositório C:\Users\trcnologia\Desktop\LegoLoop (Go, binário legoloop).
Leia AGENTS.md antes de qualquer ação. Responda em PT-BR, uma frase por linha. Sem subagentes. Sem commit e sem push sem autorização do operador.

Estado herdado (provado em 15/11/2026)
- video-rea: forensics ganhou Concepts, Processes e Repositories (internal/videorea/forensics_findings.go). Curate classifica em DEVE, PODE e NAO_PODE (internal/videorea/curator_classify.go).
- Os gates de n02 e n03 exigem "repositories" e "findings". has_pattern e word_count agora expandem ${SLUG} (internal/gate/builtins.go).
- Peças fabricadas e fechadas: legoloop-video-rea-transcrever e legoloop-video-rea-relatorio.
- Defeito aberto: fechar cria stub (SKILL de 3 linhas, DAG com nodes vazio) em internal/fechar/fechar.go:189-213 e imprime "canonizado com sucesso" com etapa em error.
- Defeito aberto: as peças do video-rea ainda não ligam ao processo (refinar e compor exigem processo.json no cofre).
- Falha alheia: internal/studio static_test.go:540 e :624 por edição concorrente de index.html.

Objetivo
1. Corrigir fechar: sem stub, sem sucesso falso (TDD Red-Green, mutação por teste).
2. Ligar transcrever ao n02 e relatorio ao n03 do video-rea.
Portão: go test -count=1 ./..., go vet ./..., go run ./tools/funclen internal, gofmt -l ., git diff --exit-code -- go.mod go.sum, go run ./tools/taxlinter internal cmd loops skills docs/templates.
Feche com legoloop fechar <slug> e leia a saída com --json.
