---
id: EMENDA-INTENT-KRAFTEN-DEPLOY-20261126
status: APPROVED
date: 2026-11-26
schema_version: 2
emenda_de: [INTENT-intv_1795723484496401600]
upstream: [ANALISE_BDD_REDE_DOMINIO_VIAS]
---

# EMENDA da INTENT: `legoloop-kraften-deploy` (rede parametrizada, domínio mestre variável, 3 vias de entrada)

> **Status:** APPROVED (resposta do usuário "APROVADO", 2026-11-26).
> **Emenda de:** `docs/intents/2026-11-26_INTENT_TRANSFORMAR_O_PIPELINE_DE_DEPL.md` (APPROVED, Stop-Gate 01).
> **Base de fatos:** ponderações do usuário sobre rede, domínio mestre e vias de entrada + exercício BDD de reescrita de rede.

## 1. Decisões do usuário (2026-11-26)

| # | Decisão | Valor aprovado |
|---|---|---|
| D1 | Rede do `docker-stack.yml` | O nome da rede **não** é fixo (`network_conexao` era cópia do documento de referência). O `n01` valida a **estrutura** `networks` declarada. O nome da rede compartilhada do ecossistema é **parâmetro de configuração** (`.env` mestre do Launcher) |
| D2 | Comunicação entre apps | Todo serviço publicado anexa a **rede compartilhada do ecossistema** (mesma rede do Traefik e das demais apps). Redes privadas extras da app são preservadas. Rede externa declarada inexistente no Swarm reprova antes do build |
| D3 | Reescrita de rede no Launcher | **Aditiva e idempotente**: o Launcher acrescenta a anexação à rede compartilhada configurada em cada serviço do `docker-stack.yml`; nunca remove redes declaradas. Quarta mutação controlada (antes: `.env`, `Host(...)`, DNS) |
| D4 | Domínio mestre | Variável. O FQDN não deriva de padrão fixo (`server.vpsconexao.org`). Domínio base vem da configuração do Launcher (`.env` mestre) ou parâmetro explícito |
| D5 | Vias de entrada | As 3 vias do LegoLoop: **terminal** (CLI), **agente/chat** (linguagem natural ou slash command), **Studio** (botão Executar Processo). Parâmetros ausentes são **coletados interativamente, um a um**, em qualquer via |
| D6 | Respostas do chat | Protocolo ASD-STE100 (STE-PT, modo prova) a partir de 2026-11-26 |

## 2. O que muda na INTENT/SPEC aprovadas

### Entra no escopo (antes era fixo ou não existia)
- Nome da rede compartilhada como parâmetro de configuração (D1).
- Anexação aditiva à rede compartilhada no Launcher (D2, D3).
- Domínio mestre como parâmetro variável (D4).
- Coleta interativa de parâmetros + registro das 3 vias de entrada (D5).

### Sai do escopo ou muda de forma
- A exigência literal `network_conexao` no `n01` (AC-E01.3) vira: **bloco `networks` declarado no `docker-stack.yml`**.
- O gate `n04` ganha prova da **rede real anexada** no Swarm (`REDE=1`), além de `CF_OLD=0`, `CF_NEW=1`, `TLS=0`, `HTTP=200`.
- O `fqdnDe` de padrão fixo sai; entra o domínio base configurável.

### Restrições inegociáveis (mantidas)
- Zero LLM nos nós; `Dockerfile` da app nunca mutado; recibo só com gates verdes.

## 3. Handoff
Autoriza a atualização da SPEC (`AC-E01.3`, novo `AC-E02.3` de rede), do PLAN (tickets de coleta interativa e da mutação aditiva) e da implementação em `internal/kraften` + Launcher.
