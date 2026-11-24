# 4. Processamento

Processar é transformar as entradas em material pronto para montar.
Esta parte mostra os scripts e a ordem em que rodam.
Os caminhos são de `ferramentas/` no repositório `proj--historias-palitos`.

## 4.1 Voz

`voz.py` lê as linhas do roteiro.
Gera uma narração por linha com edge-tts.
Grava os tempos de cada linha em `audio/timings.json`.
Comando: `python voz.py <projeto> --voz pt-BR-FranciscaNeural`.
Aviso: edge-tts envia o texto ao serviço online da Microsoft.
Fonte: `ferramentas/LEIAME.md`, linha 7; `processo/PROCESSO.md`, etapa 6.

## 4.2 Validação da timeline

`check_timeline.py` valida `timeline.json`.
Também mede o ritmo em cortes por minuto.
Fonte: `ferramentas/LEIAME.md`, linha 8.

## 4.3 Legendas

`srt.py` e `gerar_legendas.py` geram legendas a partir de `data/lines.json`.
No vídeo longo, a legenda sai em `.srt`.
No 9:16, a legenda é queimada no vídeo, em blocos de 3 palavras.
Fonte: `ferramentas/cortes.py` (docstring) e `ferramentas/LEIAME.md`, linha 17.

## 4.4 Trilha e volume

`mix.py` monta a trilha e os efeitos.
O volume final vai para -14 LUFS, com pico em -1,5 dB.
Comando: `python mix.py <projeto> --music dark`.
Fonte: `ferramentas/LEIAME.md`, linha 9; `processo/PROCESSO.md`, etapa 8.

## 4.5 Montagem da timeline

A timeline é um arquivo JSON com a lista de cenas (`shots`).
Cada cena tem duração, linhas de texto, cenário e personagem.
Scripts como `gerar_timeline_completa.py` montam esse arquivo.
Fonte: `curiosidades-malucas/gerar_timeline_completa.py` e `build_timeline.py` na raiz.

## 4.6 Recorte de cortes curtos

`cortes.py` lê uma especificação em `cortes/<video>.json`.
Cada corte é um intervalo de linhas do roteiro.
O script recorta cenas, voz e trilha, refaz a legenda e acrescenta um cartão final.
Não grava nada de novo.
Fonte: `ferramentas/cortes.py`.
