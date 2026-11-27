---
name: legoloop-video-animado
author: LegoLoop
description: Use quando precisar executar o Macroprocesso legoloop-video-animado (Epicos E07 a E12 - SPEC v2): inventario global por estilo (molde-de-pilula), contratos v2, Peca 2 vetorial sem Flow, motor portado com licenca MIT, renomeacao diorama->video com alias e prova de reproducao em video real. NAO use para um unico Processo isolado (video-roteiro, video-ativos tem skill propria).
consumes: Artifact:ValidatedScript
produces: Artifact:ReproducedVideo
artifact: renders/cena1.mp4
gate: exit_code == 0
suggested_model: sonnet
---

# 🎬 legoloop-video-animado — Inventario, Motor Portado e Prova de Reproducao

> **Autor:** LegoLoop
> **Versao:** 1.0.0
> **Motor oficial:** `legoloop video-inventario`, `legoloop video-ativos`, `legoloop video-render` (implementacao em `internal/video/*` e `motor-local/`, licenca MIT).
> **SPEC:** `docs/specs/2026-11-26_SPEC_VIDEO_ANIMADO_V2.md`

---

## 1. Contexto e Contratos de I/O

**Quando usar:** apos os Processos `legoloop-video-roteiro` e `legoloop-video-ativos` terem produzido `projeto.json` v2, `script.json` e `cenas.json`, para validar o estilo no inventario, ligar o render ao motor portado e provar a reproducao em video real (cena 1 e First Brands).
**Quando NAO usar:** para rodar so um dos 6 nos isoladamente fora de ordem, ou para gerar fundos no Flow sem autorizacao do usuario (ver regra do Flow em AGENTS.md).

---

## 2. Nos do Processo (ver o DAG do Processo)

1. **n01_inventario_por_estilo** — Inventario global por estilo: `NovoEstilo`, `GravarItem`, heranca e validacao humana (Epico E07). Comando: `legoloop video-inventario verificar --estilo molde-de-pilula --json`.
2. **n02_contratos_v2** — Contratos v2: pre-voo le `modo_producao` e `estilo` validado, storyboard confere elenco e inventario (Epico E08). Comando: `legoloop video-ativos prevoo --projeto . --json`.
3. **n03_peca2_vetorial** — Peca 2: elenco vetorial do inventario sem Flow, caminho de recorte preservado, entrega suspeita marcada (Epico E09). Comando: `legoloop video-ativos executar --projeto . --json`.
4. **n04_motor_portado** — Motor portado (`motor-local`, licenca MIT): ponte `cenas.json` -> `timeline.json`, render 9:16 deterministico, boca sincronizada (Epico E10). Comando: `legoloop video-render timeline --projeto projetos/video-diorama/first-brands`.
5. **n05_renomeacao_video** — Renomeacao `diorama`->`video` com alias depreciado, skills/processos/espinha migrados, historico preservado (Epico E11). Comando: `legoloop espinha verificar --json`.
6. **n06_prova_de_reproducao** — Prova de reproducao: cena 1 renderizada e aprovada pelo usuario (Stop-Gate humano), First Brands reproduzido com ritmo no alvo (Epico E12). Comando: `legoloop video-render executar --projeto projetos/video-diorama/diorama-demo-pilula`.

---

## 3. Regras Negativas (Guardrails)

- PROIBIDO gerar imagem ou video no Flow sem autorizacao explicita do usuario, mesmo imagem custando 0 credito.
- PROIBIDO aprovar um estilo (`video-inventario validar`) sem arquivo de aprovacao humana gravado (`AC-E07.8`).
- PROIBIDO renderizar com estilo `rascunho`; o pre-voo (`n02`) bloqueia (`AC-E08.1`).
- PROIBIDO mover `motor-local/` fora do `.gitignore` enquanto a licenca estiver `pendente` (`AC-E10.1`); so apos `liberada`/`reescrito`.
- Nao recrie logica aqui: delegue sempre ao binario `legoloop`.
- Achado e corrigido na execucao e registrado em `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`, com anticorpo em `processos/legoloop-video-animado/anticorpos.json` quando aplicavel.
