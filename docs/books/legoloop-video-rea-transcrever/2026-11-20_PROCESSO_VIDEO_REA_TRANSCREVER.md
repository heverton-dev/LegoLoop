# Frontmatter

Manual Didático: Manual Didático: legoloop-video-rea-transcrever
Processo legoloop-video-rea-transcrever no LegoLoop.
Versão 1.0.0 — Data 2026-11-20.


# Macro Fundação

Visão macro do processo legoloop-video-rea-transcrever.
Converter o áudio aac em wav 16 kHz mono com ffmpeg e transcrever com C:\Users\trcnologia\tools\whisper-cpp\Release\whisper-cli.exe e o modelo C:\Users\trcnologia\.cache\whisper-models\ggml-small.bin. Verificar que os dois caminhos existem antes de usar. Plano B: read_file multimodal do áudio. Saída: transcricao.txt.
Operação orientada aos 7 Protocolos Industriais e 5 Prisões Mecânicas.


# Meso Arquitetura

Arquitetura técnica e grafo DAG de legoloop-video-rea-transcrever.
Contratos de entrada (context_in), saída (context_out) e gates mecânicos.


# Micro Passo a Passo

Passo a passo determinístico de execução do processo legoloop-video-rea-transcrever.
Comandos CLI: `legoloop run processos/legoloop-video-rea-transcrever/pipeline.dag.json`.
Validação e testes contínuos.


