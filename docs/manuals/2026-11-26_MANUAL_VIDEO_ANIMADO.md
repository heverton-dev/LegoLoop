# Manual Operacional — `legoloop-video-animado`

Data: 2026-11-26. Versão: 1.0.0.

## 1. Finalidade

Fechar a SPEC v2 (Épicos E07 a E12) do vídeo animado:
- **E07 Inventário por estilo:** `NovoEstilo`, `GravarItem`, herança (`Herda`), validação humana por arquivo de aprovação.
- **E08 Contratos v2:** `modo_producao` + `estilo` no `projeto.json`; pré-voo bloqueia estilo `rascunho`.
- **E09 Peça 2:** elenco vetorial do inventário sem Flow; caminho de recorte preservado; entrega suspeita marcada.
- **E10 Motor portado:** `motor-local/` sob licença MIT; ponte `cenas.json`→`timeline.json`; render 9:16 determinístico com boca sincronizada.
- **E11 Renomeação:** `diorama`→`video`, alias depreciado, histórico preservado.
- **E12 Prova de reprodução:** cena 1 aprovada pelo usuário; First Brands reconstruído com ritmo no alvo.

## 2. Comandos CLI
```bash
legoloop video-inventario validar --estilo <id> --aprovacao <arq> --json
legoloop video-inventario verificar --estilo <id> --json
legoloop video-ativos prevoo --projeto <DIR> --json
legoloop video-ativos executar --projeto <DIR> --json
legoloop video-render timeline --projeto <DIR>
legoloop video-render executar --projeto <DIR>
```

## 3. Gate Mecânico
`go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .`, `legoloop espinha verificar --json`, cobertura `>= 95%` em `internal/video/inventario`, `internal/video/motor`, `internal/video/render` e `internal/video/roteiro`.

## 4. Stop-Gate Humano
Nenhum estilo vira `validado` sem arquivo de aprovação humana (`AC-E07.8`). Nenhuma cena entra como prova de reprodução sem aprovação visual do usuário gravada em arquivo (`AC-E12.1`, ver `aprovacao_cena1.json`).
