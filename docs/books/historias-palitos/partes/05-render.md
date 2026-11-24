# 5. Render

Render é a etapa que transforma a timeline em arquivo de vídeo.
A timeline descreve o que aparece em cada instante.
O render desenha cada quadro e junta tudo num arquivo MP4.

## 5.1 O que entra

- `timeline.json`: lista de cenas, duração e texto.
- Imagens de cenário e personagens, em `flow/` e `sprites/`.
- Áudio de voz e trilha, em `audio/`.
- `marca.json`: fontes, cores e logo.
- Formato: 16:9, 9:16 ou 1:1.

## 5.2 O que o render faz

1. Abre o navegador Chromium, controlado pelo Playwright.
2. Cada cena vira um clipe curto. O log mostra as cenas em paralelo.
3. Os clipes ficam em `.render_cache/`, com `list.txt` e `silent.mp4`.
4. O FFmpeg junta os clipes e mistura o áudio.
5. O arquivo final vai para `renders/`.

Comando: `python engine/render.py <projeto>/timeline.json --formato 9:16 --out x.mp4`.
Fonte: `ferramentas/LEIAME.md`, linha 10.

## 5.3 Evidência do teste

Render do corte `02-invasao-muscular` (9:16).
Saída: 1080 x 1920 pixels, 31,7 s, em 1 min 20 s.
O log mostrou 5 cenas em paralelo e a linha "Done".
Arquivo: `renders/shorts/02-invasao-muscular.mp4`.
Verificado com `ffprobe` na própria execução do `cortes.py --render`.

## 5.4 Aviso sobre fontes

Em duas cenas do teste, o log mostrou "marca: falhou" para as fontes.
O render terminou mesmo assim.
A causa não foi verificada.
Hipótese a testar: falha ao carregar fontes do Google Fonts.

## 5.5 Pós-produção

`marca.py` adiciona vinheta de abertura, logo e vinheta de fechamento.
Comando: `python marca.py <projeto> x.mp4 --out final.mp4`.
Fonte: `ferramentas/LEIAME.md`, linha 11.
