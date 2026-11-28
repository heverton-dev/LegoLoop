# Arquitetura

- `internal/video/distribuicao`: acervo, agenda, revisao, vitrine, capa, publicacao, comentarios e aprendizado.
- `internal/video/distribuicao/scripts/capa.py`: capa 1280x720 com Pillow; o texto entra pelo codigo, nunca pelo modelo de imagem.
- `internal/cli/video_distribuicao.go`: os 9 grupos de subcomando.

O acervo (`acervo/acervo.json`) guarda cada item com tipo quente ou perene, validade, canais, status e autorizacao. O agendador e deterministico: dado o mesmo acervo, a mesma configuracao e o mesmo instante, o plano e o mesmo.
