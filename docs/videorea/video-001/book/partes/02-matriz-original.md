# 2. Matriz Original e Curadoria
## Stack declarada no vídeo
| Ferramenta               | Papel                      | Custo declarado                  | Substituição local       |
| :----------------------- | :------------------------- | :------------------------------- | :----------------------- |
| Claude Code com Opus 5.5 | Orquestra a edição         | R$ 100,00 por mês (plano)        | Gateway local 9Router    |
| Higgsfield Plus          | Cenas realistas            | Não informado no vídeo           | Opcional, fora do núcleo |
| HeyGen                   | Voz clonada                | Não informado no vídeo           | Áudio do operador        |
| ElevenLabs               | Voz por API (alternativa)  | Não informado no vídeo           | Áudio do operador        |
| HyperFrames              | Animações em código        | Gratuito                         | Mantido                  |
| Remotion                 | Alternativa ao HyperFrames | Licença comercial conforme porte | Dispensado               |
| FFmpeg, Node.js, Python  | Base técnica               | Gratuito                         | Mantidos                 |
## Consumo observado
Uma edição consumiu cerca de 10% do limite do plano (de 35% para 45%). O vídeo citado somou mais de 76 mil visualizações orgânicas.
## Diagnóstico do host
`ffmpeg`, `node`, `python` e `whisper-cli` já estavam instalados. O `hyperframes` estava com a instalação quebrada (link para um clone apagado) e foi reinstalado com `bun install -g hyperframes`.
