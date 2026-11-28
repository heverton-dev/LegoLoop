# Manual Operacional — `legoloop-video-animado`

Data: 2026-11-27. Versão: 1.0.0.

## 1. Finalidade

Operar o Macroprocesso composto `legoloop-video-animado` integrando de ponta a ponta:
- **P1 `legoloop-video-roteiro`:** Entrevista na Etapa 0, descoberta de pautas e redação de roteiro (`projeto.json`, `marca.json`, `script.json`).
- **P2 `legoloop-video-ativos`:** Inventário por estilo, elenco vetorial sem dependência de Flow e storyboard estrito (`cenas.json`).
- **P3 `legoloop-video-render`:** Voz neural com Whisper, timeline com ritmo controlado, mix com loudness medido e render 9:16 pelo motor local com licença MIT.
- **P4 `legoloop-video-distribuicao`:** Pacote de revisão sem contexto, vitrine, agenda e publicação Postiz.
- **Gate final de prova:** Verificação da renderização em vídeo real com ritmo no alvo (12 a 21 cortes por minuto) e aprovação humana.

## 2. Comandos CLI dos Subprocessos e Macro
```bash
legoloop video-roteiro executar --projeto <DIR>
legoloop video-ativos executar --projeto <DIR> --json
legoloop video-render prevoo --projeto <DIR>
legoloop video-render voz --projeto <DIR>
legoloop video-render mix --projeto <DIR>
legoloop video-render motor --projeto <DIR>
legoloop video-distribuicao prevoo --projeto <DIR>
legoloop video-distribuicao publicar --projeto <DIR>
legoloop video-render executar --projeto <DIR>
```

## 3. Gate Mecânico
`go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .`, `legoloop espinha verificar --json`, cobertura `>= 95%` em `internal/video/inventario`, `internal/video/motor`, `internal/video/render`, `internal/video/roteiro` e `internal/video/distribuicao`.

## 4. Stop-Gate Humano
Nenhum estilo vira validado sem aprovação humana gravada em arquivo. Nenhuma publicação externa ocorre sem aprovação e flag explícita (`--confirmo-publico`). Nenhuma cena entra como prova sem verificação visual e ritmo no alvo de 12 a 21 cortes por minuto.
