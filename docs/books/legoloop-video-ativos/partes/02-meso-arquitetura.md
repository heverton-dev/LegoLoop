# Capítulo 2 — Arquitetura das Etapas 3, 4 e 5
Etapa 3 valida 3 fichas por personagem e `sprites/*.json` sem sobreposição; Etapa 4 gera SVGs 100% autorais ou via fila Flow resiliente com backoff e idempotência por hash; Etapa 5 compila `cenas.json` com 100% das linhas e `>= 4` elementos móveis.
