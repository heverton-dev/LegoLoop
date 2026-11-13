# Catálogo de combos do 9Router (13/11/2026)

Fonte: lista real de `/api/combos` em roteador.signa-docs.com, lida ao gerar este arquivo. Coluna "função": chamada de ferramenta com `max_tokens` 800, 2 amostras (`2026-11-13_funcao_800_todos.json` e `2026-11-13_reteste20.json`).
Qualidade de resposta NÃO foi medida em nenhum combo. Tempos e HTTP: ver `2026-11-13_todos_combos.json`, `..._combos_restantes.json`, `..._piloto_t_agente.json`, `..._troca_4_combos.json`.

Avisos provados:
- Cloudflare: cota diária gratuita de 10.000 neurons esgotada (429 no corpo, `2026-11-13_reteste20.json`). Combos `m-cf-*` e `p-cloudflare-128k` ficam indisponíveis até renovar; membros `cf/` foram tirados dos combos de tarefa.
- Groq: HTTP 413 com pedido de 20k tokens ou mais (limite 8.000 tokens por minuto por organização). Só `t-micro-8k` usa Groq.
- Embedding: combos não funcionam em `/v1/embeddings` (HTTP 400). Use `gemini/gemini-embedding-2` direto.
- Hydrogen só funciona nas chaves DEV_06 e DEV_13; Argon e Carbon têm cota gratuita zero (ver mapa de contas).
- `code-pro`, `code-fast`, `code-free`, `task-micro` e `code-tests` são nomes usados pelo roteamento do código (`internal/plangen/route.go`); não apagar.

## Combos por tarefa (`t-*`) e legados mantidos (17)

| combo | membros | função |
|---|---|---|
| `code-fast` | `mistral/ministral-3b-latest`<br>`gemini/gemini-3.5-flash-lite`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`mistral/ministral-3b-2512` | sim |
| `code-free` | `openrouter/dots-studio/dots-3-note-preview:free`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | sim |
| `code-pro` | `gemini/gemini-3-flash-preview`<br>`gemini/gemini-3.5-flash`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`openrouter/dots-studio/dots-3-note-preview:free`<br>`mistral/ministral-3b-latest` | sim |
| `code-tests` | `mistral/codestral-latest`<br>`gemini/gemini-3.5-flash-lite` | não testado |
| `free-chat-L` | `mistral/ministral-14b-2512`<br>`mistral/ministral-8b-2512`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | não testado |
| `free-codigo-L` | `openrouter/cohere/north-mini-code:free`<br>`mistral/codestral-latest`<br>`mistral/mistral-code-latest` | não testado |
| `free-visao-L` | `mistral/ministral-14b-2512`<br>`mistral/ministral-8b-2512`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | não testado |
| `gemini-4-argon` | `gemini/gemini-4-argon`<br>`gemini/gemini-3.1-pro-preview`<br>`gemini/gemini-3.8-flash` | não testado |
| `gemini-4-hydrogen` | `gemini/gemini-4-hydrogen`<br>`gemini/gemini-3.8-flash`<br>`gemini/gemini-3.6-flash` | não testado |
| `megacombo-embedding` | `gemini/gemini-embedding-2` | não testado |
| `t-agente-128k` | `mistral/ministral-3b-latest`<br>`gemini/gemini-3.5-flash-lite`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`openrouter/dots-studio/dots-3-note-preview:free`<br>`mistral/ministral-3b-2512` | sim |
| `t-chat-128k` | `mistral/codestral-latest`<br>`gemini/gemini-3.5-flash-lite`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`openrouter/cohere/north-mini-code:free`<br>`mistral/ministral-3b-latest` | não |
| `t-codigo-128k` | `mistral/codestral-latest`<br>`openrouter/cohere/north-mini-code:free`<br>`mistral/codestral-2508`<br>`mistral/mistral-code-latest` | não |
| `t-micro-8k` | `groq/openai/gpt-oss-120b`<br>`groq/openai/gpt-oss-20b` | sim |
| `t-raciocinio-128k` | `gemini/gemini-3.5-flash-lite`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free`<br>`gemini/gemini-3-flash-preview` | sim |
| `t-visao-128k` | `mistral/ministral-3b-latest`<br>`nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`openrouter/dots-studio/dots-3-note-preview:free`<br>`mistral/ministral-3b-2512`<br>`nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | sim |
| `task-micro` | `mistral/ministral-3b-latest`<br>`gemini/gemini-3.5-flash-lite` | sim |

## Combos por provedor (`p-*`) (5)

| combo | membros | função |
|---|---|---|
| `p-cloudflare-128k` | `cf/@cf/moondream/moondream3.1-9B-A2B`<br>`cf/@cf/meta/llama-4-scout-17b-16e-instruct`<br>`cf/@cf/nvidia/nemotron-3-120b-a12b`<br>`cf/@cf/openai/gpt-oss-20b`<br>`cf/@cf/openai/gpt-oss-120b`<br>`cf/@cf/daily-co/phonellm-mini-alpha-2`<br>`cf/@cf/google/gemma-4-26b-a4b-it`<br>`cf/@cf/ibm-granite/granite-4.0-h-micro`<br>`cf/@cf/zai-org/glm-4.7-flash`<br>`cf/@cf/qwen/qwen3.8-27b`<br>`cf/@cf/mistralai/mistral-small-3.1-24b-instruct` | indisponível no teste |
| `p-gemini-128k` | `gemini/gemini-3.5-flash-lite`<br>`gemini/gemini-3-flash-preview`<br>`gemini/gemini-3.1-flash-lite-preview`<br>`gemini/gemini-3.1-flash-lite`<br>`gemini/gemini-flash-lite-latest`<br>`gemini/gemini-robotics-er-2-preview`<br>`gemini/gemini-3.5-flash` | sim |
| `p-mistral-128k` | `mistral/codestral-latest`<br>`mistral/ministral-3b-latest`<br>`mistral/codestral-2508`<br>`mistral/ministral-3b-2512`<br>`mistral/ministral-8b-2512`<br>`mistral/ministral-14b-2512`<br>`mistral/mistral-code-latest` | não |
| `p-nvidia-128k` | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning`<br>`nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b`<br>`nvidia/nvidia/nemotron-3-ultra-550b-a55b`<br>`nvidia/nvidia/nemotron-3.5-lightning-30b-a3b`<br>`nvidia/poolside/laguna-xs-2.1`<br>`nvidia/mistralai/mistral-large-4-1t-a52b`<br>`nvidia/nvidia/nemotron-3.5-bolt-e4b`<br>`nvidia/minimaxai/minimax-m3.1-flash` | sim |
| `p-openrouter-128k` | `openrouter/cohere/north-mini-code:free`<br>`openrouter/dots-studio/dots-3-note-preview:free`<br>`openrouter/recursive/flint-v1.0:free`<br>`openrouter/nvidia/nemotron-3.5-lightning:free`<br>`openrouter/liquid/lfm-3-1.7b:free`<br>`openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free`<br>`openrouter/nvidia/nemotron-3-ultra-550b-a55b:free`<br>`openrouter/openrouter/free`<br>`openrouter/nvidia/nemotron-3.5-super-vl:free` | sim |

## Combos por modelo (`m-*`) (51)

| combo | membros | função |
|---|---|---|
| `m-cf-daily-co-phonellm-mini-alpha-2` | `cf/@cf/daily-co/phonellm-mini-alpha-2` | indisponível no teste |
| `m-cf-deepseek-ai-deepseek-r1-distill-qwen-32b` | `cf/@cf/deepseek-ai/deepseek-r1-distill-qwen-32b` | indisponível no teste |
| `m-cf-google-gemma-4-26b-a4b-it` | `cf/@cf/google/gemma-4-26b-a4b-it` | indisponível no teste |
| `m-cf-ibm-granite-granite-4-0-h-micro` | `cf/@cf/ibm-granite/granite-4.0-h-micro` | indisponível no teste |
| `m-cf-meta-llama-3-1-8b-instruct-fp8` | `cf/@cf/meta/llama-3.1-8b-instruct-fp8` | indisponível no teste |
| `m-cf-meta-llama-3-2-1b-instruct` | `cf/@cf/meta/llama-3.2-1b-instruct` | indisponível no teste |
| `m-cf-meta-llama-3-2-3b-instruct` | `cf/@cf/meta/llama-3.2-3b-instruct` | indisponível no teste |
| `m-cf-meta-llama-4-scout-17b-16e-instruct` | `cf/@cf/meta/llama-4-scout-17b-16e-instruct` | indisponível no teste |
| `m-cf-mistralai-mistral-small-3-1-24b-instruct` | `cf/@cf/mistralai/mistral-small-3.1-24b-instruct` | indisponível no teste |
| `m-cf-moondream-moondream3-1-9b-a2b` | `cf/@cf/moondream/moondream3.1-9B-A2B` | indisponível no teste |
| `m-cf-nvidia-nemotron-3-120b-a12b` | `cf/@cf/nvidia/nemotron-3-120b-a12b` | indisponível no teste |
| `m-cf-openai-gpt-oss-120b` | `cf/@cf/openai/gpt-oss-120b` | indisponível no teste |
| `m-cf-openai-gpt-oss-20b` | `cf/@cf/openai/gpt-oss-20b` | indisponível no teste |
| `m-cf-qwen-qwen2-5-coder-32b-instruct` | `cf/@cf/qwen/qwen2.5-coder-32b-instruct` | indisponível no teste |
| `m-cf-qwen-qwen3-30b-a3b-fp8` | `cf/@cf/qwen/qwen3-30b-a3b-fp8` | indisponível no teste |
| `m-cf-qwen-qwen3-8-27b` | `cf/@cf/qwen/qwen3.8-27b` | indisponível no teste |
| `m-cf-zai-org-glm-4-7-flash` | `cf/@cf/zai-org/glm-4.7-flash` | indisponível no teste |
| `m-gemini-4-hydrogen` | `gemini/gemini-4-hydrogen` | sim |
| `m-gemini-gemini-3-1-flash-lite` | `gemini/gemini-3.1-flash-lite` | sim |
| `m-gemini-gemini-3-1-flash-lite-preview` | `gemini/gemini-3.1-flash-lite-preview` | sim |
| `m-gemini-gemini-3-5-flash` | `gemini/gemini-3.5-flash` | sim |
| `m-gemini-gemini-3-5-flash-lite` | `gemini/gemini-3.5-flash-lite` | sim |
| `m-gemini-gemini-3-flash-preview` | `gemini/gemini-3-flash-preview` | sim |
| `m-gemini-gemini-flash-lite-latest` | `gemini/gemini-flash-lite-latest` | sim |
| `m-gemini-gemini-robotics-er-2-preview` | `gemini/gemini-robotics-er-2-preview` | sim |
| `m-mistral-codestral-2508` | `mistral/codestral-2508` | não |
| `m-mistral-codestral-latest` | `mistral/codestral-latest` | não |
| `m-mistral-ministral-14b-2512` | `mistral/ministral-14b-2512` | sim |
| `m-mistral-ministral-3b-2512` | `mistral/ministral-3b-2512` | sim |
| `m-mistral-ministral-3b-latest` | `mistral/ministral-3b-latest` | sim |
| `m-mistral-ministral-8b-2512` | `mistral/ministral-8b-2512` | sim |
| `m-mistral-mistral-code-latest` | `mistral/mistral-code-latest` | não |
| `m-mistral-voxtral-small-2507` | `mistral/voxtral-small-2507` | sim |
| `m-mistral-voxtral-small-latest` | `mistral/voxtral-small-latest` | sim |
| `m-nvidia-minimaxai-minimax-m3-1-flash` | `nvidia/minimaxai/minimax-m3.1-flash` | sim |
| `m-nvidia-mistralai-mistral-large-4-1t-a52b` | `nvidia/mistralai/mistral-large-4-1t-a52b` | sim |
| `m-nvidia-nemotron-3-5-bolt-e4b` | `nvidia/nvidia/nemotron-3.5-bolt-e4b` | não |
| `m-nvidia-nemotron-3-5-lightning-30b-a3b` | `nvidia/nvidia/nemotron-3.5-lightning-30b-a3b` | não |
| `m-nvidia-nemotron-3-5-super-vl-120b-a12b` | `nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b` | sim |
| `m-nvidia-nemotron-3-nano-omni-30b-a3b-reasoning` | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | sim |
| `m-nvidia-nemotron-3-ultra-550b-a55b` | `nvidia/nvidia/nemotron-3-ultra-550b-a55b` | sim |
| `m-nvidia-poolside-laguna-xs-2-1` | `nvidia/poolside/laguna-xs-2.1` | sim |
| `m-openrouter-cohere-north-mini-code-free` | `openrouter/cohere/north-mini-code:free` | sim |
| `m-openrouter-dots-studio-dots-3-note-preview-free` | `openrouter/dots-studio/dots-3-note-preview:free` | sim |
| `m-openrouter-liquid-lfm-3-1-7b-free` | `openrouter/liquid/lfm-3-1.7b:free` | sim |
| `m-openrouter-nvidia-nemotron-3-5-lightning-free` | `openrouter/nvidia/nemotron-3.5-lightning:free` | não |
| `m-openrouter-nvidia-nemotron-3-5-super-vl-free` | `openrouter/nvidia/nemotron-3.5-super-vl:free` | sim |
| `m-openrouter-nvidia-nemotron-3-nano-omni-30b-a3b-reasoning` | `openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | sim |
| `m-openrouter-nvidia-nemotron-3-ultra-550b-a55b-free` | `openrouter/nvidia/nemotron-3-ultra-550b-a55b:free` | sim |
| `m-openrouter-openrouter-free` | `openrouter/openrouter/free` | sim |
| `m-openrouter-recursive-flint-v1-0-free` | `openrouter/recursive/flint-v1.0:free` | sim |

Total de combos no 9Router: 73.
Backups (reversão): `scratch/backup_9router_combos_20261113_pre_*`.
