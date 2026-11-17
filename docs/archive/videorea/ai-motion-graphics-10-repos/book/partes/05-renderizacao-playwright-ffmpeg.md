# Capítulo 5 — Renderização Playwright + FFmpeg

## Playwright (Captura de Frames)

- Abre HTML em Chromium headless
- Captura screenshots sequenciais (fps × duração)
- 90 frames para 3s @ 30fps
- `headless: false` proibido em CI

## FFmpeg (Conversão MP4)

```bash
ffmpeg -y -framerate 30 -i frame_%04d.png -c:v libx264 -pix_fmt yuv420p -s 1920x1080 -r 30 output.mp4
```

## Verificação

`ffprobe` confirma:
- Codec: h264
- Duração: ±0.1s do esperado
- Resolução: igual ao Request

## Limite

Renderização de 60s de vídeo ≤ 300s (5 min).
