---
name: legoloop-fronteira-ai-driven
description: Skill universal de fronteira-ai-driven
---
# legoloop-fronteira-ai-driven

## REGRAS NEGATIVAS (Proibido)

- Proibido executar comandos fora do escopo autorizado da fronteira AI-driven.
- Proibido ignorar falhas de gate ou suprimir erros em relatórios de execução.

## Grafo de Execução (DAG) & Nós Canônicos

| ID do Nó | Descrição / Saída |
| :--- | :--- |
| `T-01` | Prompt builder e orquestração estéril (`internal/orchestrix/prompt.go`) |
| `T-02` | Roteamento de modelos e fallbacks (`internal/router/`) |
| `T-03` | Execução isolada hermética (`internal/runner/isolated.go`) |
| `T-04` | Portões de fechamento canônico (`internal/fechar/`) |
| `T-05` | Relatório de homologação da fronteira (`docs/reports/2026-11-14_FRONTEIRA_AI_DRIVEN.md`) |
| `T-06` | Ponto de entrada CLI e comandos (`internal/cli/cli.go`) |

