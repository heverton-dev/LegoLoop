# Parte V — Micro: Acabamento e Sonorização

## 5.1 Os Três Efeitos de Acabamento de Matt Loui

No tutorial original, o criador gasta 5 minutos dentro do CapCut aplicando três efeitos manuais que transformam uma animação digital comum em uma peça cinematográfica:

1. **Aberração Cromática nas Bordas (Chromatic Edges):** Separação sutil dos canais RGB nas extremidades da tela, simulando a refração óptica de uma lente real de cinema.
2. **Blur Radial nos Cantos (Radial Blur 2):** Desfocamento sutil de vinheta circular que puxa o olhar do espectador para o centro dinâmico do vídeo.
3. **Sonorização na Frenagem (Sound Design):** Inserção sincronizada de efeitos sonoros de transição:
   * *Whoosh:* No início do movimento de aceleração;
   * *UI Pop / Click:* No momento exato em que a mola física atinge o ponto de frenagem;
   * *Line Sound / Swish:* Ao traçar linhas ou bordas.

---

## 5.2 A Substituição Soberana via FFmpeg Determinístico

O LegoLoop extirpa a necessidade de aplicar esses efeitos manualmente no CapCut Pro. Um script Python orquestra filtros de FFmpeg e mesclagem de áudio com precisão de frame:

```bash
ffmpeg -i overlay.mov -i fundo.mp4 -i whoosh.wav -i pop.wav   -filter_complex "[0:v]chromatic_aberration=0.002,radial_blur=radius=0.1[v_fx]; [1:v][v_fx]overlay=0:0[out_v]; [2:a]adelay=400|400[a1]; [3:a]adelay=750|750[a2]; [a1][a2]amix=inputs=2[out_a]"   -map "[out_v]" -map "[out_a]" -c:v libx264 -crf 17 dist/final_master.mp4
```
