# Prompt de Ativação Hermética — Ciclos 2 a 5 (`legoloop-video-diorama`)

Cole em uma **nova sessão limpa** (ou inicie com `/swarm`) para economizar tokens e paralelizar os Ciclos 2, 3 e 4 com interface A2UI:

```text
Leia apenas:
1. docs/progress/2026-11-24_macro_video_diorama/estado.json
2. docs/intents/2026-11-24_MAPA_MACRO_VIDEO_DIORAMA.md

Regras desta execução:
- O Ciclo 1 (legoloop-diorama-roteiro) já está DONE com 95,7% de cobertura.
- Ofereça primeiro a abertura da interface interativa A2UI (via MCP sessao_criar/sessao_publicar/sessao_aguardar ou CLI legoloop fabricar entrevista ui).
- Execute os Ciclos 2 (legoloop-diorama-ativos), 3 (legoloop-diorama-render) e 4 (legoloop-diorama-distribuicao) em paralelo via /swarm com workspace="share" (um worktree por pacote internal/diorama/<modulo>), seguindo TDD >= 95% cobertura, funclen <= 35 linhas e taxlinter zero violações.
- Faça o merge sequencial apenas nos arquivos listados em arquivos_merge_sequencial de estado.json.
- Ao concluir cada Ciclo, execute legoloop fechar <slug> e atualize estado.json (status: DONE).
- Por fim, componha legoloop-video-diorama (Ciclo 5) via legoloop compor.
```
