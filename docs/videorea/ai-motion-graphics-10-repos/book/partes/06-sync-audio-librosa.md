# Capítulo 6 — Sync de Áudio com Librosa

## Por Que Sync Importa

Todo visual change que acontece no beat da música torna o vídeo "premium". É o diferencial entre amador e profissional.

## Como Funciona

1. `scripts/beat_detect.py` carrega áudio via Librosa
2. `librosa.beat.beat_track()` detecta beats
3. Retorna timestamps (segundos) para sincronizar animações

## Fallback Gracioso

- Áudio ausente → sem sync, sem erro
- Arquivo corrompido → sem sync, sem erro
- Pipeline continua normalmente

## Exemplo (120 BPM)

```json
[0.0, 0.5, 1.0, 1.5, 2.0, 2.5]
```
