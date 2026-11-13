---
name: legoloop-gateway-agy
description: "Use quando o usuário pedir para apontar, trocar a chave, validar ou desligar o gateway de LLM do Antigravity CLI (agy) por variáveis de ambiente AGY_LLM_GATEWAY_*, com 9Router ou outro provedor, sem rodízio de contas Google."
---

# LegoLoop Gateway AGY — Chaves por Variável de Ambiente

Esta skill altera o destino das chamadas de modelo do `agy`. Ela nunca troca a conta Google.
Manual: `docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md`. Livro: `docs/books/gateway-llm-agy/livro.json`.

## 1. Regras inegociáveis

1. Nunca peça, leia em voz alta, imprima ou grave o valor da chave. Use só o nome de uma variável que já guarda a chave.
2. Confirme variáveis só com `definida: True/False`.
3. Nunca rode `gw agy next`, `failover`, `devNN`, `sync`, `wincred`, `probe --all` ou `refresh all`. O `gw` bloqueia esses comandos.
4. Nunca gere wrapper do `agy` nem rode `gw harness sync-wrappers` para o `agy`.
5. Declare "Não verificado" para a posição do Google sobre gateways de terceiros.
6. Antes de editar arquivo de configuração, faça cópia de segurança.

## 2. Variáveis

| Variável | Obrigatória | Função |
| :--- | :---: | :--- |
| `AGY_LLM_GATEWAY_URL` | Sim | URL HTTPS do gateway |
| `AGY_LLM_GATEWAY_API_KEY` | Sim | Chave que o gateway exige |
| `AGY_LLM_GATEWAY_WIRE_PROTOCOL` | Não | `genai` (padrão) ou `openai` |
| `AGY_LLM_GATEWAY_CA_CERT` | Não | Certificado raiz |
| `AGY_LLM_GATEWAY_MODELS` | Não | Modelos permitidos |
| `AGY_LLM_GATEWAY_HEADERS` | Não | Cabeçalhos extras `CHAVE: VALOR` |

## 3. Processo

1. **Ligar.** Receba a URL, o protocolo e o nome da variável com a chave. Grave `URL`, `WIRE_PROTOCOL` e `API_KEY` no ambiente do usuário (PowerShell: `[Environment]::SetEnvironmentVariable(nome, valor, 'User')`). Nesta máquina, `gw agy 9router` faz o mesmo para o 9Router e lê `NINEROUTER_KEY`.
2. **Trocar a chave.** Leia a variável da chave nova e regrave só `AGY_LLM_GATEWAY_API_KEY`.
3. **Desligar.** Apague as três variáveis. Nesta máquina, `gw agy off`.
4. **Validar.** Confirme `definida: True/False`. Chame `<URL>/models` com a chave e espere HTTP 200, sem imprimir a chave.
5. **Avisar.** Diga ao usuário para abrir um terminal novo. Terminal aberto antes mantém o ambiente antigo.

## 4. Limites

- O `admin_settings.json` do administrador vence a variável de ambiente.
- O protocolo `openai` não funciona no JetBrains nem no Zed (servidor ACP). Use `genai` neles.
- Chave que apareceu em saída de comando está comprometida. Mande o usuário revogar e trocar.
