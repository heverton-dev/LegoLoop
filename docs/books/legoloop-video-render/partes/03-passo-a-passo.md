# Passo a Passo

1. `legoloop video-render prevoo --projeto <DIR>` confere tudo antes de gastar tempo.
2. `legoloop video-render voz --projeto <DIR>` grava `audio/voice.wav` e `data/timings.json`.
3. `legoloop video-render timeline --projeto <DIR>` grava `timeline.json` e aplica o gate de ritmo.
4. `legoloop video-render mix --projeto <DIR>` grava `audio/soundtrack.wav`.
5. `legoloop video-render motor --projeto <DIR>` grava `render.mp4`.

Se um passo falhar, corrija a causa e repita so aquele passo. O usuario assiste ao `render.mp4` e aprova no Stop-Gate humano.
