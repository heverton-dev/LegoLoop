# 5. Operação e Limites
## Executar
`legoloop run loops/video-001/pipeline.dag.json --var raw_video="caminho/do/video.mp4"`
## Pré-requisitos
`ffmpeg` com `libvpx-vp9`, `node`, `python`, `hyperframes` 0.9.x com telemetria desativada (`hyperframes telemetry disable`) e `whisper-cli` com `ggml-small.bin`.
## Verificação do Golden Fixture
O fixture `fixture_3s.mp4` tem 1920x1080, 60 fps, 3,0 s e áudio. O comando `python -I scripts/verificar_video.py --video <arquivo> --duracao 3` retorna exit 0 quando conforme.
## Limites conhecidos
| Limite                                        | Efeito                             |
| :-------------------------------------------- | :--------------------------------- |
| `n08` com duração fixa de 3 s                   | Ajustar para vídeos reais          |
| `n06` gera um único WebM                        | Não há um arquivo por cena         |
| Caminhos do Whisper fixos no `n02`              | Específicos deste host             |
| DAG não executado de ponta a ponta pelo `legoloop run` | Fixture montado nó a nó            |
| Validação do `whisper-cli`                      | Apenas sobre 40 s de áudio real    |
