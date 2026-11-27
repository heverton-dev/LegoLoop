# 3. Passo a Passo

1. `legoloop video-inventario validar --estilo <id> --aprovacao <arq>` para sair de `rascunho` para `validado` (exige Stop-Gate humano gravado em arquivo).
2. `legoloop video-inventario verificar --estilo <id> --json` confere o hash de cada item pesado.
3. `legoloop video-ativos prevoo --projeto <dir> --json` confirma `modo_producao` e estilo validado antes de gerar qualquer ativo.
4. `legoloop video-ativos executar --projeto <dir> --json` monta `cenas.json` a partir do inventário (ou do Flow, para fundos ainda não aprovados).
5. `legoloop video-render timeline --projeto <dir>` chama a ponte do motor: lê `cenas.json` + elenco + timings + inventário e grava `timeline.json`, com `check_timeline.py` real aprovando (`CHECK OK`).
6. `legoloop video-render executar --projeto <dir>` chama o motor (Playwright headless + ffmpeg), mede o MP4 (1080x1920, 30fps, áudio, duração igual à voz, zero quadros em branco) e grava `render.mp4`.
7. Prova de reprodução: a cena 1 do `diorama-demo-pilula` foi renderizada, corrigida (estouro de data no calendário) e aprovada pelo usuário; o First Brands foi reconstruído do zero (10 cortes, voz local via edge-tts) e bateu o ritmo alvo (`cuts/min` 15,1; mediana 3,9s).
