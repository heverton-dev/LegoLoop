# Parte IV — Receitas Multi-Aspect Ratio e Broadcast Bundle

## Matriz de Formatos e Áreas de Segurança

Cada rede social possui restrições visuais específicas (sobreposição de botões de like, comentários e legenda de rodapé). O motor ajusta o layout com base na matriz abaixo:

| Formato         | Resolução       | Destino Principal         | Margem de Segurança     |
| :-------------- | :-------------- | :------------------------ | :---------------------- |
| **9:16**            | 1080x1920       | TikTok, Reels, Shorts     | 120px topo, 280px base  |
| **16:9**            | 1920x1080       | YouTube Desktop, Cursos   | 60px uniforme           |
| **1:1**             | 1080x1080       | Instagram Feed, LinkedIn  | 80px topo e base        |

## Anatomia do Pacote de Entrega (`dist/<job_id>/`)

Ao concluir a renderização, a pasta de entrega é montada com todos os ativos requeridos para distribuição profissional:

```
dist/job_20261107_01/
├── video_final_9x16.mp4      # Master de transmissão codificado em H.264 / AAC
├── legendas.srt              # Legenda universal para players externos
├── legendas.vtt              # Legenda web com estilização CSS
├── timeline.fcpxml           # XML v1.9 importável no DaVinci Resolve Studio
├── qa_sheet.jpg              # Prova de contato dos 12 frames-chave
└── POST_COPY.md              # Texto de postagem com hashtags e timing de publicação
```
