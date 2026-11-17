---
name: legoloop-triagem
description: Use when routing natural language requests or workspace contexts to the right LegoLoop process via deterministic triage.
---

# Triagem Universal LegoLoop

Antes de rotear qualquer pedido em linguagem natural no chat para um processo do LegoLoop, execute a triagem determinística:

```bash
legoloop triagem "<entrada>" --workspace "<DIR>" --json
```

## Regras de Operação

1. Leia o campo `processo`, `etapa_inicial`, `camada` e `precisa_confirmacao` do JSON retornado.
2. Se `precisa_confirmacao` for `true` ou `camada` for `"ambiguidade"`, apresente as opções ao usuário antes de iniciar qualquer processo.
3. Nunca adivinhe o processo ou pule a triagem quando houver múltiplos processos possíveis.

## REGRAS NEGATIVAS
- Nunca acione um LLM para entradas que possuem comando, slug ou estado de artefato inequívoco no disco.
- Nunca invente rotas ou processos que não constem no catálogo validado de `triagem.json`.

---

## DAG & Nós Canônicos

| ID do Nó | Nome do Nó | Descrição |
| :--- | :--- | :--- |
| `t1-motor-triagem` | Motor de Triagem | Classificação determinística em 4 camadas de precisão |
| `t2-fabricar-13-artefatos` | Fabricar 13 Artefatos | Geração dos manifestos canônicos de triagem |
| `t3-pontos-entrada` | Pontos de Entrada | Integração da triagem com CLI e comandos de chat |
| `t4-saneamento-processos` | Saneamento de Processos | Limpeza e alinhamento de processos legados |
| `t5-portoes-fechamento` | Portões e Fechamento | Validação mecânica e fechamento do ciclo de triagem |

