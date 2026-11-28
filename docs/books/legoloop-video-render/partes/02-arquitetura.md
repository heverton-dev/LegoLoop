# Arquitetura

- `internal/video/render`: pre-voo, voz, whisper, timeline, ritmo, mix, SFX, legendas e render ffmpeg.
- `internal/video/motor`: ponte `cenas.json` para `timeline.json` do motor, pre-voo do motor, render de referencia e determinismo.
- `internal/cli/video_render_etapas.go`: os subcomandos `prevoo`, `voz`, `mix` e `motor`.
- `motor-local/`: `puppet.js`, `sets.js`, `render.py`, `run_render.py` e `check_timeline.py`.

O `timeline.json` tem dois produtores. O Go o grava no passo `timeline`. O motor o regrava no passo `motor`. Por isso a ordem e fixa.
