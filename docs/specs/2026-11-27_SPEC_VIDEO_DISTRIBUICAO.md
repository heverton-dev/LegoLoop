# SPEC — legoloop-video-distribuicao

- **Data:** 2026-11-27
- **Versao:** 1.0.0
- **Pacote:** `internal/video/distribuicao`; CLI `legoloop video-distribuicao`.
- **INTENT:** `docs/intents/2026-11-27_INTENT_VIDEO_DISTRIBUICAO.md`

## Criterios binarios de aceite

| ID | Criterio | Prova |
| :--- | :--- | :--- |
| `AC-E01.1` | Acervo ausente e vazio; gravacao atomica; id repetido, tipo, status, visibilidade, canais e validade do quente sao validados | `TestAcervo*, TestItemValidar` |
| `AC-E01.2` | Regra de originalidade: tema igual ou com 80% das palavras repetidas e recusado (ErrTemaRepetido) | `TestAdicionarAplicaOriginalidade, TestTemasSemelhantes` |
| `AC-E01.3` | Autorizacao de publicacao gravada e revogavel antes da publicacao | `TestAutorizarERevogar` |
| `AC-E02.1` | Quente fura a fila e sai no primeiro horario livre; passando de 48 h ou da validade, expira | `TestQuenteFuraAFila, TestQuenteForaDaJanelaEExpiraSemLugar` |
| `AC-E02.2` | Perene respeita cadencia (24 h) e teto diario; sem horario livre, pula com motivo | `TestPereneRespeitaCadenciaETeto, TestSemHorarioLivre*` |
| `AC-E02.3` | Item sem autorizacao nao e agendado; vencido expira; publicado e expirado nao recebem decisao | `TestSemAutorizacaoPula, TestValidadeVencidaExpira*` |
| `AC-E02.4` | Alerta de estoque baixo de perenes | `TestAlertaDeEstoqueBaixo` |
| `AC-E02.5` | Todo curto publico do YouTube gera o item de Reel do Instagram, sem duplicar | `TestGerarReels` |
| `AC-E03.1` | Parecer exige veredito ENVIAR ou NAO ENVIAR, minutagem mm:ss e, no NAO ENVIAR, ao menos um problema | `TestParecerValidar` |
| `AC-E03.2` | Parecer nunca e sobrescrito; so ENVIAR libera a vitrine | `TestGravarParecerELiberar` |
| `AC-E03.5` | Vitrine recusa video cuja versao nao foi liberada (registro review/vN_liberada.json com parecer ENVIAR) | `TestLiberarEExigirLiberada, TestGravarPublicacaoExigeLiberacao` |
| `AC-E03.3` | Folhas de contato (24 quadros por folha) e pacote do revisor com video, roteiro, regras e checklist | `TestGerarFolhas, TestGravarPacote` |
| `AC-E03.4` | Numero da proxima versao conta renders/vN.mp4 sem apagar nenhuma | `TestProximaVersao` |
| `AC-E04.1` | Titulos com nota: ate 60 caracteres, palavra principal no inicio, numero, sem caixa alta; acima de 100 zera | `TestNotaTitulo, TestPontuarTitulosOrdena` |
| `AC-E04.2` | Vitrine exige 5 titulos, 3 capas por formato, fontes, tags, resumo, rotulo de IA, capitulos mm:ss e arquivos no disco | `TestVitrineValidaOPacote` |
| `AC-E04.3` | publicacao.md grava titulo escolhido e descricao com capitulos, fontes e aviso de IA | `TestGravarPublicacaoEDescricao` |
| `AC-E04.4` | Capa 1280x720 com ate 4 palavras (3 no curto), fundo sem texto, aviso acima de 2 MB | `TestGerarCapa*, TestValidarTextoCapa` |
| `AC-E05.1` | Pedido de envio: titulo ate 100, visibilidade valida, publico exige confirmacao, rotulo de IA obrigatorio | `TestPedidoEnvioValidar` |
| `AC-E05.2` | Dry-run nao chama a rede; envio cria postagem privada com containsSyntheticMedia | `TestEnviarDryRun*, TestEnviarCriaPostagem*` |
| `AC-E05.3` | A chave nunca aparece em erro; sem POSTIZ_API_KEY ou POSTIZ_URL ha erro tipado | `TestEnviarOcultaAChave*, TestSegredoDoAmbienteENovoCliente` |
| `AC-E05.4` | Publicar do acervo exige autorizacao e visibilidade igual; falha nao altera o acervo | `TestPublicarDoAcervo` |
| `AC-E05.5` | Envio real grava a linha de auditoria em acervo/envios.json, sem a chave e com a resposta cortada em 120 caracteres | `TestRegistrarEnvioAcumulaESemSegredo` |
| `AC-E06.1` | Rascunho de resposta so para comentario sem resposta; so a aprovada sai | `TestRascunharEDecidir, TestRespostasGravarECarregar` |
| `AC-E06.2` | Aprendizado semanal com no minimo 3 medidas, validadas, em acervo/aprendizado.md | `TestGerarAprendizado, TestGravarAprendizado` |
| `AC-E07.1` | CLI video-distribuicao com 9 grupos de subcomando; publicar enviar e dry-run por padrao | `TestVideoDistribuicao*` |
| `AC-NF.1` | Cobertura >= 95% em internal/video/distribuicao; funclen e gofmt limpos | `go test -cover` |

## Comandos de verificacao

```bash
go test -count=1 -cover ./internal/video/distribuicao/... ./internal/cli/ -run VideoDistribuicao
go run ./tools/funclen internal
gofmt -l internal
legoloop fechar legoloop-video-distribuicao
```
