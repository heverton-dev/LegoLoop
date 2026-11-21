# Capítulo 6 — Sync de Áudio com Librosa

## Por Que Sync Importa

Todo visual change que acontece no beat da música torna o vídeo "premium". É o diferencial entre amador e profissional.

## Como Funciona

1. `scripts/beat_detect.py` carrega áudio via Librosa
2. `librosa.beat.beat_track()` detecta beats
3. Retorna timestamps (segundos) para sincronizar animações

## Fallback Gracioso

O contrato de áudio tem dois modos. O pedido declara o modo.

### Modo Broadcast (Shorts, Reels, Standalone)

- O áudio integrado a -14 LUFS é obrigatório.
- Áudio ausente quebra o gate. O Processo para com erro.
- Arquivo corrompido quebra o gate.

### Modo Motion Overlay (camadas transparentes para DaVinci Resolve)

- A trilha de áudio é opcional.
- Áudio ausente produz um MP4 com canal mudo. O gate não quebra.
- Arquivo corrompido produz um MP4 com canal mudo, sem sync e sem erro.

Os dois modos resolvem a contradição entre o contrato BDD de -14 LUFS (`harness-engineering`) e o fallback silencioso (`ai-motion-graphics`).

## Exemplo (120 BPM)

```json
[0.0, 0.5, 1.0, 1.5, 2.0, 2.5]
```
