---
name: legoloop-diorama-roteiro
author: LegoLoop
description: REDIRECIONADO em 2026-11-26 para legoloop-video-roteiro. Use skills/legoloop-video-roteiro/SKILL.md.
consumes: Artifact:Niche
produces: Artifact:ValidatedScript
artifact: script.json
gate: exit_code == 0
suggested_model: sonnet
---

# legoloop-diorama-roteiro (redirecionado)

## Redirecionamento

Este processo foi renomeado para `legoloop-video-roteiro` em 2026-11-26 (AC-E11.3).
Use `skills/legoloop-video-roteiro/SKILL.md` e o comando `legoloop video-roteiro`.

## Regras

- O comando `legoloop diorama-roteiro` continua funcionando como alias depreciado (aviso em stderr, mesmo exit code).
- Nao recrie logica aqui: delegue sempre ao binario `legoloop`.
