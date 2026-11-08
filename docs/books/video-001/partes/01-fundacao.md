# 1. Fundação e Princípios
## O problema
Editar vídeos gravados e criar animações é lento. O fluxo original usa Claude Code com a skill HyperFrames, narração com voz clonada e cenas realistas geradas por serviço externo.
## O que este pipeline decide
| Decisão               | Escolha                                                    |
| :-------------------- | :--------------------------------------------------------- |
| Transcrição           | `whisper-cli` local com modelo `ggml-small`                    |
| Animações             | HyperFrames local, HTML renderizado para WebM com alfa     |
| Montagem              | FFmpeg, saída 1920x1080 a 60 fps                           |
| Voz                   | Áudio do operador; sem clonagem por API                    |
| Revisão humana        | Plano de cenas aprovado antes de qualquer render           |
## Fora do escopo
Serviços de nuvem pagos, voz neural por API externa, cenas realistas de gerador de vídeo e interface web própria.
