---
name: legoloop-video-render
author: LegoLoop
description: Use quando precisar executar o Processo legoloop-video-render (P3 do processo composto legoloop-video-animado): pre-voo, voz neural com conferencia Whisper, timeline com gate de ritmo, mix com loudness medido e render 9:16 com voz pelo motor portado. NAO use para roteiro (legoloop-video-roteiro), ativos (legoloop-video-ativos) nem distribuicao.
consumes: Artifact:VideoAssets
produces: Artifact:RenderedVideo
artifact: render.mp4
gate: exit_code == 0
suggested_model: sonnet
---

# legoloop-video-render — Voz, Timeline, Mix e Render com o Motor

> **Autor:** LegoLoop
> **Versao:** 1.0.0
> **Motor oficial:** `legoloop video-render` (implementacao em `internal/video/render`, `internal/video/motor` e `internal/cli/video_render_etapas.go`; motor Python em `motor-local/`, licenca MIT).
> **Caminho canonico:** o motor. O master 16:9 em ffmpeg (`video-render executar`) e fallback declarado.

---

## 1. Contexto e Contratos de I/O

**Quando usar:** depois dos Processos `legoloop-video-roteiro` e `legoloop-video-ativos`, quando a pasta do projeto tem `projeto.json`, `marca.json`, `script.json`, `cenas.json` e `elenco.json`.
**Quando NAO usar:** para criar roteiro, ativos ou publicar; para gerar imagem no Flow sem autorizacao do usuario.

**Pasta do projeto:** `projetos/video-animado/<nome-do-projeto>/`. Todos os artefatos do video ficam nela.

**Entradas:** `projeto.json` (inclui `estilo_visual`), `marca.json`, `script.json`, `cenas.json`, `elenco.json` e o inventario do estilo em `inventario/video/estilos/<estilo>/`.
**Saidas:** `audio/voice.wav`, `data/timings.json`, `timeline.json`, `audio/soundtrack.wav` e `render.mp4` (1080x1920, 30 fps, com voz).

---

## 2. Nos do Processo (ver o DAG do Processo)

Ordem obrigatoria. O `timeline.json` tem dois produtores (Go e motor); o motor grava por ultimo.

1. **n01_prevoo** — `legoloop video-render prevoo --projeto <DIR>`: contratos, provedor de voz, ffmpeg, faster-whisper, pastas de saida e motor local (python, Playwright, Chromium).
2. **n02_voz** — `legoloop video-render voz --projeto <DIR>`: sintetiza cada linha, confere com o Whisper (WER por linha) e grava `audio/voice.wav` e `data/timings.json`.
3. **n03_timeline** — `legoloop video-render timeline --projeto <DIR>`: compoe `timeline.json` e aplica o gate de ritmo.
4. **n04_mix** — `legoloop video-render mix --projeto <DIR>`: voz, SFX e trilha em `audio/soundtrack.wav`, com loudness medido.
5. **n05_motor** — `legoloop video-render motor --projeto <DIR> [--raiz DIR] [--motor-dir DIR]`: monta a timeline do motor, exige CHECK OK do `check_timeline.py` e renderiza `render.mp4` em 9:16 com voz.

Fallback: `legoloop video-render executar --projeto <DIR>` renderiza o master 16:9 em ffmpeg sobre imagens estaticas.

---

## 3. Regras Negativas (Guardrails)

- PROIBIDO abrir aba ou Chrome do Flow sem autorizacao explicita do usuario.
- PROIBIDO gerar imagem ou video no Flow sem autorizacao; em todo `flow_generate` de fundo, informar `model` Nano Banana e `type` image.
- PROIBIDO rodar o motor com estilo em `rascunho`; o pre-voo e o inventario bloqueiam.
- PROIBIDO pular o CHECK OK do `check_timeline.py`; o `motor` nao chama o render sem ele.
- PROIBIDO gravar artefato fora de `projetos/video-animado/<nome-do-projeto>/`.
- Nao recrie logica aqui: delegue sempre ao binario `legoloop`.
- Achado e corrigido na execucao e registrado em `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`, com anticorpo em `processos/legoloop-video-render/anticorpos.json` quando aplicavel.

## 4. Escopo declarado (o que este Processo ainda nao cobre)

- Cortes verticais 9:16 a partir de um master longo 16:9 (Etapa 8b da SPEC, `AC-E05.4` e `AC-E05.5`): adiado, registrado em `docs/backlog/BACKLOG_VIDEO_ANIMADO.md` (BL-03). O motor ja entrega o video em 9:16.

## 5. Comandos de Verificacao

```bash
legoloop video-render --help
legoloop validate processos/legoloop-video-render/pipeline.dag.json
go test -count=1 -cover ./internal/video/render/... ./internal/video/motor/...
legoloop fechar legoloop-video-render
```

## 6. Checklist de Conclusao

- [ ] `video-render prevoo` sem `FALTA`.
- [ ] `data/timings.json` com WER global dentro do limiar.
- [ ] `check_timeline.py` com CHECK OK.
- [ ] `render.mp4` em 1080x1920, 30 fps, com voz, duracao igual a da voz (1 quadro).
- [ ] Stop-Gate humano: o usuario assistiu e aprovou o video.
