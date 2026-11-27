# Macro — Fundação

## O Problema

Publicar aplicações na VPS envolve mutação de domínio, credenciais, DNS e orquestração Docker. Executado manualmente (ou por agentes probabilísticos), o processo é lento, não auditável e propenso a erro.

## A Solução

Um Processo LegoLoop **100% determinístico** (`legoloop-kraften-deploy`): cada deploy é uma sequência fixa de 5 nós Go, sem nenhum LLM, com gates binários e recibo de prova.

## A Fronteira Universal do `.zip`

O Processo atende **qualquer stack** (Go, Node, Python, Rust, PHP, estático) desde que o pacote traga exatamente os 3 artefatos universais:

1. `Dockerfile` com `HEALTHCHECK` — gera a imagem e expõe a porta HTTP;
2. `docker-stack.yml` com bloco `networks` declarado e as labels do Traefik (o nome da rede compartilhada é configuração do ambiente, não regra do Processo);
3. `.env.example` declarando `ADMIN_USER` e `ADMIN_PASS`.

O `Dockerfile` da aplicação **nunca** é mutado: toda a mutação de domínio, credenciais e roteamento ocorre apenas sobre `.env`, a label `Host(...)` do `docker-stack.yml` e os registros DNS da Cloudflare.

## Restrições Inegociáveis

- Zero LLM e zero agentes intermediários nos nós.
- Extração anti-zip-slip: nomes validados antes de qualquer gravação.
- CNAME idempotente, criado antes de qualquer consulta DNS.
- Recibo só com gates verdes.
