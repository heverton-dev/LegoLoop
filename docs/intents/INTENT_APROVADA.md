# INTENT APROVADA — legoloop-kraften-deploy

> **Data:** 2026-11-26
> **Stop-Gate 01:** APROVADO pelo usuário no chat (sessão `ses_ffe5de6a05040ffelVrsrUnDoh`).
> **INTENT:** `docs/intents/2026-11-26_INTENT_TRANSFORMAR_O_PIPELINE_DE_DEPL.md` (alias `docs/intents/INTENT.md`, `status: APPROVED`).
> **INTERVIEW:** `intv_1795723484496401600` (cofre de sessões).
> **Espelho Canônico:** `docs/intents/2026-11-26_ESPELHO_CANONICO_KRAFTEN_DEPLOY.md`.

## Decisão
O usuário aprova formalmente a INTENT de transformar o pipeline de deploy determinístico do Studio Launcher no Processo LegoLoop `legoloop-kraften-deploy` (5 nós determinísticos, zero LLM), com as divergências reconciliadas:
1. Esteira implementada em Go nativo (`internal/kraften/` + `cmd/legoloop`); atende qualquer stack com `Dockerfile`/`docker-stack.yml`.
2. Fora de escopo: mutação do `Dockerfile` da app, agentes LLM intermediários, stacks sem os 3 artefatos universais, rede fora de `network_conexao`.

## Handoff
Autoriza a geração da SPEC (`docs/specs/`, processo `legoloop-especificacao`, Stop-Gate 02).
