# 4. Arquitetura do DAG
O DAG tem oito nós em cadeia linear.
| Nó                     | Tipo           | Função                                | Saída                  |
| :--------------------- | :------------- | :------------------------------------ | :--------------------- |
| `n01_extrair_audio`        | script         | FFmpeg extrai áudio 16 kHz mono       | `audio.wav`              |
| `n02_transcrever`         | script         | `whisper-cli` com `-ojf`                  | `audio.json`             |
| `n03_alinhar_roteiro`      | script         | Casa as falas do roteiro às palavras  | `segmentos.json`         |
| `n04_plano_cenas`          | skill          | Escreve o plano cena por cena         | `plano.md`               |
| `n05_gate_plano`           | stop_gate       | Aprovação humana do plano             | `plano_aprovado.json`     |
| `n06_renderizar_cenas`     | script         | HyperFrames renderiza WebM com alfa a 60 fps | `cenas.webm`             |
| `n07_montar`              | script         | FFmpeg sobrepõe as cenas ao vídeo base | `final.mp4`              |
| `n08_verificar`           | script         | `verificar_video.py` confere o resultado | `verificacao.json`       |
## Contrato de `segmentos.json`
Lista de objetos com `id`, `texto`, `inicio_ms` e `fim_ms`. O fim é sempre maior que o início e não há sobreposição entre segmentos.
## Alinhamento
A função `alinhar` ignora acentos, caixa e pontuação, procura as duas primeiras palavras de cada fala numa janela de 80 palavras e falha com erro se a fala não existir na transcrição.
