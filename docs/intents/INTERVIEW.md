# INTERVIEW — legoloop-kraften-deploy

> **Sessão:** `intv_1795723484496401600` (cofre de sessões do motor `legoloop fabricar entrevista`)
> **Data:** 2026-11-26
> **Motor:** 9router/code-fast — entrevista estacionou em 10% de confiança; Espelho Canônico emitido e reconciliado manualmente.

## Resumo dos Turnos

**Turno 1 — Pergunta:** Os 5 nós especificados em `PROCESSO_DEPLOY_DETERMINISTICO_LEGOLOOP.md` utilizam exclusivamente comandos Go nativos e chamadas Docker/Cloudflare, sem dependência de agentes LLM internos?
**Resposta do usuário:** Confirmado — 100% determinístico zero-LLM; o `.md` é a especificação de referência, implementada sem alterações estruturais.

**Turno 2 — Pergunta (motor):** A estrutura de diretórios deve seguir `pkg/fabrica/` ou `cmd/legoloop/`?
**Reconciliação humana (pós-finalização do motor):**
1. `pkg/fabrica/` não existe — padrão real do repo: motor em `internal/kraften/`, comando em `cmd/legoloop` (como `internal/diorama/roteiro` + `legoloop diorama-roteiro`).
2. Fora de escopo corrigido: a esteira atende QUALQUER stack com `Dockerfile` e `docker-stack.yml` (Go, Node, Python, Rust, PHP, estático); Go nativo é a implementação da esteira, não restrição das apps.

## Espelho Canônico
Consolidado em `docs/intents/ESPELHO_CANONICO.md` (registro ISO: `2026-11-26_ESPELHO_CANONICO_KRAFTEN_DEPLOY.md`).
