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
