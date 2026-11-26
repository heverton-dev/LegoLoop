---
name: legoloop-diorama-ativos
author: LegoLoop
description: REDIRECIONADO em 2026-11-26 para legoloop-video-ativos. Use skills/legoloop-video-ativos/SKILL.md.
consumes: Artifact:ValidatedScript
produces: Artifact:DioramaAssets
artifact: cenas.json
gate: exit_code == 0
suggested_model: sonnet
---

# legoloop-diorama-ativos (redirecionado)

## Redirecionamento

Este processo foi renomeado para `legoloop-video-ativos` em 2026-11-26 (AC-E11.3).
Use `skills/legoloop-video-ativos/SKILL.md` e o comando `legoloop video-ativos`.

## Regras

- O comando `legoloop diorama-ativos` continua funcionando como alias depreciado (aviso em stderr, mesmo exit code).
- Nao recrie logica aqui: delegue sempre ao binario `legoloop`.
