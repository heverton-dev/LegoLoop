# Manual Operacional — `legoloop-video-render`

- **Data:** 2026-11-27
- **Versao:** 1.0.0
- **Motor CLI:** `legoloop video-render`
- **Pacotes Go:** `internal/video/render`, `internal/video/motor` e `internal/cli/video_render_etapas.go`
- **Motor Python:** `motor-local/` (licenca MIT, `motor-local/LICENCA.json`)
- **DAG Canonico:** `processos/legoloop-video-render/pipeline.dag.json`
- **Skill Universal:** `skills/legoloop-video-render/SKILL.md`

---

## 1. Proposito

O Processo `legoloop-video-render` e o P3 do processo composto `legoloop-video-animado`. Ele transforma `script.json`, `cenas.json` e `elenco.json` em `render.mp4`: voz neural conferida pelo Whisper, timeline com gate de ritmo, mix com loudness medido e render 9:16 com voz pelo motor portado.

A pasta do projeto e `projetos/video-animado/<nome-do-projeto>/`. Todos os artefatos do video ficam nela.

## 2. Comandos e Subcomandos

| Subcomando | Entrada | Saida e Gate |
| :--- | :--- | :--- |
| `legoloop video-render prevoo --projeto DIR` | contratos, ambiente local | Bloqueia em contrato, provedor, ferramenta, pasta ou motor ausente (`exit 0` se tudo OK) |
| `legoloop video-render voz --projeto DIR` | `script.json`, `marca.json` | `audio/voice.wav`, `data/timings.json`; reprova linha com WER acima do limiar |
| `legoloop video-render timeline --projeto DIR` | `cenas.json`, `data/timings.json` | `timeline.json` com gate de ritmo (12 a 21 cortes por minuto, mediana de 2,5 a 4 s) |
| `legoloop video-render mix --projeto DIR` | `audio/voice.wav`, `timeline.json` | `audio/soundtrack.wav` com LUFS e true peak medidos |
| `legoloop video-render motor --projeto DIR [--raiz DIR] [--motor-dir DIR]` | contratos, inventario, `audio/soundtrack.wav` | `timeline.json` do motor com CHECK OK e `render.mp4` em 1080x1920 |
| `legoloop video-render executar --projeto DIR` | `timeline.json` | Fallback: master 16:9 em ffmpeg sobre imagens estaticas |

## 3. Regras Inegociaveis

1. **Ordem fixa:** prevoo, voz, timeline, mix, motor. O `timeline.json` tem dois produtores; o motor grava por ultimo.
2. **Sem CHECK OK nao ha render:** o `motor` so chama o render depois do `check_timeline.py` aprovar.
3. **Formato 9:16:** o comando `motor` passa `--formato 9:16`; sem isso o motor renderiza em 16:9 e o gate de 1080x1920 reprova.
4. **Estilo validado:** o estilo do projeto (`estilo_visual`) precisa estar `validado` no inventario, com fundos aprovados.
5. **Flow so com autorizacao:** imagem custa 0 credito, video custa; sempre informar `model` e `type` no `flow_generate`.
6. **Stop-Gate humano:** o usuario assiste e aprova o video final.

## 4. Escopo declarado

Cortes verticais a partir de um master longo 16:9 (Etapa 8b da SPEC) estao adiados (BL-03 em `docs/backlog/BACKLOG_VIDEO_ANIMADO.md`). O motor ja entrega o video em 9:16.

## 5. Solucao de Problemas

| Sintoma | Causa | Acao |
| :--- | :--- | :--- |
| `render.mp4` em 1920x1080 | `--formato 9:16` ausente | usar `video-render motor`, nao chamar o `run_render.py` direto |
| `fundo ausente do inventario` | cenario fora do estilo ou nao aprovado | aprovar o item no inventario do estilo |
| `check_timeline reprovou` | timeline fora do contrato | corrigir `cenas.json` ou o texto das linhas |
| `ferramenta ausente` no prevoo | ffmpeg, faster-whisper, Playwright ou Chromium | instalar o item citado |
