---
name: legoloop-video-animado
author: LegoLoop
description: Use quando precisar executar o Macroprocesso legoloop-video-animado composto (P1 a P4): encadeia video-roteiro, video-ativos, video-render e video-distribuicao com gate final de prova e ritmo no alvo (12 a 21 cortes por minuto) em video real. O usuario invoca o macro e ele conduz todas as etapas.
consumes: Artifact:TopicOrHome
produces: Artifact:DistributedVideo
artifact: renders/cena1.mp4
gate: exit_code == 0
suggested_model: sonnet
---

# 🎬 legoloop-video-animado — Macroprocesso Composto de Producao de Video

> **Autor:** LegoLoop  
> **Versao:** 1.0.0  
> **Motor oficial:** Composicao deterministica de subprocessos (`legoloop-video-roteiro`, `legoloop-video-ativos`, `legoloop-video-render`, `legoloop-video-distribuicao`) e gate final de prova em video real.  
> **SPEC:** `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md`  

---

## 1. Contexto e Contratos de I/O

**Quando usar:** para produzir um video animado vertical completo do zero a publicacao, passando pelas 4 fases canonicas (roteiro, ativos, renderizacao e distribuicao). O usuario nao precisa fornecer tema antecipadamente: o processo faz as perguntas na Etapa 0 do P1 (nicho, fontes, formato, duracao, politica de gates e pautas).  
**Quando NAO usar:** para rodar apenas um processo isolado que ja possui skill propria (`video-roteiro`, `video-ativos`, `video-render` ou `video-distribuicao`).  

* **Entradas (Consumes):** `Artifact:TopicOrHome` (dialogo de alinhamento com o usuario na Etapa 0 do roteiro).
* **Saidas (Produces):** `Artifact:DistributedVideo` (`renders/cena1.mp4`, `render.mp4` e pacote de publicacao no acervo).

---

## 2. Nos do Processo Composto

1. **n01_video_roteiro** — P1: Subprocesso `legoloop-video-roteiro` (Etapas 0 a 2). Coleta configuracoes com o usuario, descobre pautas e gera `projeto.json`, `marca.json`, `script.json` e `plano.md`.
2. **n02_video_ativos** — P2: Subprocesso `legoloop-video-ativos` (Etapas 3 a 5). Seleciona estilo no inventario, gera elenco vetorial, cenarios e valida storyboard estrito gerando `cenas.json`.
3. **n03_video_render** — P3: Subprocesso `legoloop-video-render`. Executa voz neural com conferencia Whisper, timeline com gate de ritmo, mixagem sonora com ducking e render vertical 9:16 pelo motor portado (`motor-local`, licenca MIT).
4. **n04_video_distribuicao** — P4: Subprocesso `legoloop-video-distribuicao`. Monta pacote de revisao sem contexto, parecer de envio, vitrine, agenda e envio pelo Postiz (dry-run por padrao).
5. **n05_prova_reproducao** — Gate final e Stop-Gate de prova em video real: validacao da cena renderizada com ritmo de 12 a 21 cortes por minuto e aprovacao humana. Comando: `legoloop video-render executar --projeto .`.

---

## 3. Regras Negativas (Guardrails)

- PROIBIDO gerar imagem ou video no Flow sem autorizacao explicita do usuario.
- PROIBIDO aprovar um estilo (`video-inventario validar`) sem arquivo de aprovacao humana gravado (`AC-E07.8`).
- PROIBIDO renderizar com estilo `rascunho`; o pre-voo bloqueia (`AC-E08.1`).
- PROIBIDO mover `motor-local/` fora do `.gitignore` enquanto a licenca estiver pendente.
- Nao recrie logica aqui: delegue sempre aos comandos canonicos do binario `legoloop`.
- Achado e corrigido na execucao e registrado em `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`.
