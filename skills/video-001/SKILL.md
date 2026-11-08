---
name: video-001
description: Use quando precisar editar um vídeo gravado a partir de roteiro e vídeo de referência com custo zero e 100% local - transcrição com whisper-cli, alinhamento do roteiro às palavras, plano de cenas aprovado por humano, animações renderizadas com HyperFrames (HTML para WebM com alfa) e montagem final 1080p60 com FFmpeg. Baseado no método de Chiara Costa (Claude Code + HyperFrames).
tags: [audiovisual, video-rea, hyperframes, whisper, ffmpeg, claude-code, chiara-costa]
---

# 🎬 Skill Universal: `video-001`

Pipeline canônico de edição por roteiro, derivado do vídeo "Automatizei a edição dos meus vídeos com Claude Code" (Chiara Costa), sem APIs pagas.

## 1. Ativação
* **CLI:** `legoloop run loops/video-001/pipeline.dag.json --var raw_video="caminho/do/video.mp4"`
* **Chat:** `/video-001`, "editar vídeo por roteiro com HyperFrames".

## 2. Entradas e saídas
* **Consome:** vídeo bruto, `roteiro.txt` (uma fala por linha) e, opcionalmente, vídeo de referência de estilo.
* **Produz:** `segmentos.json`, `plano.md`, `cenas/cenas.webm`, `final.mp4` (1920x1080, 60 fps) e `verificacao.json`, tudo em `docs/videorea/video-001/scratch/`.

## 3. Topologia (8 nós lineares)
`n01_extrair_audio` → `n02_transcrever` → `n03_alinhar_roteiro` → `n04_plano_cenas` → 🛑 `n05_gate_plano` → `n06_renderizar_cenas` → `n07_montar` → `n08_verificar`

## 4. Protocolo do nó `n04_plano_cenas`
1. Ler `segmentos.json` e o roteiro.
2. Escrever `plano.md`: uma seção por cena com `id`, trecho do roteiro, estilo e duração. Nunca calcular tempos; usar os de `segmentos.json`.
3. Se houver vídeo de referência, descrever o estilo observado em texto.
4. Parar. O `n05_gate_plano` exige aprovação humana antes de qualquer render.

## 5. Guardrails negativos
* **NUNCA** renderizar sem plano aprovado no `n05`.
* **NUNCA** chamar API paga ou de nuvem (HeyGen, ElevenLabs, Higgsfield). O áudio é o do operador.
* **NUNCA** deixar `<script src="http...">` no projeto HyperFrames; o GSAP é local.
* **NUNCA** gravar fora de `docs/videorea/video-001/scratch/`.
* **NUNCA** declarar o vídeo pronto sem o `verificar_video.py` retornar exit 0.

## 6. Pré-requisitos do host
`ffmpeg` com `libvpx-vp9`, `node`, `python`, `hyperframes` 0.9.x (`bun install -g hyperframes`, com `hyperframes telemetry disable`) e `whisper-cli` com modelo `ggml-small.bin`.

## 7. Limitações conhecidas
* O `n08_verificar` usa `--duracao 3` fixo (valor do Golden Fixture); ajuste para vídeos reais.
* O `n06` gera um único WebM de cenas, não um por cena.
* Os caminhos do `whisper-cli` e do modelo estão fixos no `n02` para este host.
