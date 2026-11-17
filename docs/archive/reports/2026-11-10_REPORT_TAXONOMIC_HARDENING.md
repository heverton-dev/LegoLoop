# REPORT: Endurecimento Taxonômico ASD-STE100

Data: 2026-11-10

## 1. Resultado por ticket

| Ticket | Estado | Observação |
| :--- | :---: | :--- |
| 1. Dicionário e `internal/taxonomy` | Existia antes desta sessão | Não alterado por mim. Outro agente editou `taxonomy.go` durante a sessão. |
| 2. `tools/taxlinter` | Existia antes desta sessão | Não alterado por mim. Outro agente editou `main.go` e `linter_test.go` durante a sessão. |
| 3. Studio, API REST e Cofre | Feito | Ver seção 2. |
| 4. Governança e templates | Feito, com ressalva | Ver seção 3. |

## 2. Ticket 3

- `internal/studio/server.go`: `/api/processos` (listar e gerenciar) e `/api/dag` (status e leitura) são as rotas canônicas. `/api/loops` continua registrada como rota legada deprecada.
- `internal/studio/static/app.js` e `index.html`: rótulos "esteira" viraram "DAG" ou "Processo". O `app.js` chama `/api/processos/`.
- Corrigi `TestFabricarInspectorHasAProcessPanel`, que exigia a string `/api/loops/` no `app.js`; agora exige `/api/processos/`.
- `internal/vault/vault.go`: comentários dizem "Cofre de Segredos". A CLI já usava "cofre canônico" para `~/.aidd/canonical/`.
- Não existe `static/processos.js` com rótulos a trocar além dos já tratados no `app.js` e no `index.html`.

## 3. Ticket 4

- `AGENTS.md`: Lei #15 (item 12) com a tabela canônica resumida e o portão `go run ./tools/taxlinter ...`.
- `docs/templates/TEMPLATE_DAG.json`: ids `node_01_*` viraram `n01_*`. O linter retorna zero violações em `docs/templates`.
- `TEMPLATE_ROUTE.json` mantém `"type": "unit_task"`. Esse valor pertence ao enum de `schemas/dag.v1.json` e de `internal/dag/types.go`, usado em 12 pacotes. Renomear exige um ticket de contrato próprio.
- Os ids `TICK-*` da ROUTE ficaram: o domínio `plan_unit` define "Ticket" como termo aprovado.

## 4. Ressalva do portão do linter

O comando `go run ./tools/taxlinter internal cmd loops skills docs/templates` retorna exit 1 (469 violações na última execução, contagem instável porque o dicionário e o linter mudaram durante a sessão). Quase todas são o termo "loop" em identificadores e comentários Go. Só `docs/templates` está limpo. A Lei #15 registra o portão e determina um ticket de normalização para cada alvo com violação herdada.

## 5. O que não foi provado

- Nenhuma mutação foi aplicada nos Tickets 1 e 2 por mim.
- A UI do Studio não foi aberta no navegador nesta sessão.
