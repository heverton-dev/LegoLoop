# Parte IV — Prática Operacional e Cadeia até BUILD

## Executando o Runner em Projetos Locais ou Publicados

Antes de rodar qualquer auditoria, valide a integridade do grafo e do ambiente:

```bash
legoloop validate processos/legoloop-web-qa/pipeline.dag.json
node --test tools/webqa/webqa.test.mjs
node tools/webqa/webqa.mjs preflight --json
```

Para auditar uma aplicação local subindo o servidor automaticamente e preservando sessão autenticada:

```bash
node tools/webqa/webqa.mjs run \
  --cmd "npm run dev" \
  --port 3000 \
  --alvo meu-portal \
  --routes "/,/dashboard,/config" \
  --storage-state auth.json
```

## Estrutura do Relatório em `docs/qa-web/`

A função `writeConsolidatedReport` aplica Path Jail estrito (via `sanitizeTargetSlug`) e grava dois arquivos complementares dentro de `docs/qa-web/YYYY-MM-DD_<alvo>/`:

- **`RELATORIO_QA_WEB.json`:** Payload estruturado com `targetUrl`, `routes`, `viewports`, `status`, `failures`, `metrics` e `findings`.
- **`RELATORIO_QA_WEB.md`:** Sumário legível para humanos e agentes contendo todos os achados numerados por ferramenta, severidade e rota.

## O Fluxo Obrigatório Pós-Relatório até `BUILD`

Diferente de scripts avulsos de CI, o `legoloop-web-qa` proíbe correções improvisadas após um relatório vermelho. Quando `RELATORIO_QA_WEB.md` registra falhas, a Camada 7 (`n11_entrevista_correcao`) lê o relatório como insumo e conduz uma entrevista socrática de afunilamento (uma pergunta por vez até confiança >= 95%) para decidir prioridade, escopo e riscos de cada correção.

Concluída a entrevista, a Camada 8 (`n12_governanca_build`) emite a cadeia completa de governança:

1. `INTERVIEW` — decisões validadas para cada achado do relatório.
2. `INTENT` e `SPEC` — critérios binários de correção por componente ou rota afetada.
3. `PLAN`, `ROUTE` e `DAG` — fatiamento topológico das correções sem conflito de escrita.
4. `SWARM` e ciclo `BUILD` — execução isolada em worktree com reexecução do próprio `legoloop-web-qa` como gate final de aceitação.
