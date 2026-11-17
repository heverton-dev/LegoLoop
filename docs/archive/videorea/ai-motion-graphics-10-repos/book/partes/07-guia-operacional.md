# Capítulo 7 — Guia Operacional

## Comando Rápido

```bash
legoloop run ai-motion-graphics
```

## Entrada Típica

> "Crie um vídeo explicativo de 30 segundos sobre como um painel solar funciona, 1920x1080"

## Fluxo do Operador

1. Pipeline sugere repos → **confirme** ou ajuste
2. Geração HTML/WebGL → aguarde
3. Renderização Playwright + FFmpeg → aguarde
4. MP4 entregue em `output/`

## Solução de Problemas

| Problema | Solução |
|:---|:---|
| Playwright não inicia | Verificar `npx playwright install chromium` |
| FFmpeg não encontrado | `winget install ffmpeg` ou adicionar ao PATH |
| Librosa não detecta beats | Verificar `pip install librosa` |
| MP4 sem áudio | Fornecer `AudioPath` no Request |

## Custo

**R$ 0,00/mês** — todos os 10 repos são open-source.
