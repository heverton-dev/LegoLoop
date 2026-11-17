# Rascunho de combos do 9Router (13/11/2026) - NADA APLICADO

Fonte única: `docs/reports/data/2026-11-13_homologacao_membros.json` e `..._homologacao_extra.json` (chamadas reais via roteador.signa-docs.com, 1 amostra por teste).
Critério de entrada: HTTP 200 em pedido curto, ~20k e ~100k tokens. Tempo = segundos medidos em 1 amostra. Qualidade de resposta NÃO foi medida.
Reteste (`..._reteste.json`, 3 amostras por tamanho): um membro retestado entra se tiver HTTP 200 em pelo menos 2 de 3 amostras em cada tamanho; 503/429 em função ou imagem contam como 'não provado'.
Critério estrito: todos os testes executados do modelo (inclusive função e imagem) devem dar HTTP 200. Falhas 503 do Mistral e do OpenRouter vinham com [429] no corpo, ou seja, limite momentâneo: precisam de reteste antes de condenar o modelo.
Cuidado: `ministral-3b-latest` e `ministral-3b-2512` são o mesmo provedor e provavelmente o mesmo modelo; não contam como diversidade de cota (não verificado se são o mesmo modelo).
Fora por prova: Groq (HTTP 413 em 20k e 100k, 3 de 3 modelos), Cohere (HTTP 405, 11 de 11), gemini-3.8-flash e Argon/Carbon (sem cota gratuita ou 403).

## Combos propostos

### `t-agente-128k` - agentes com ferramentas (chamada de função provada)
| # | membro | provedor | 100k (s) | função | imagem |
|---|---|---|---|---|---|
| 1 | `mistral/ministral-3b-latest` | mistral | 2.4 | sim | sim |
| 2 | `gemini/gemini-3.5-flash-lite` | gemini | 2.8 | sim | não provada |
| 3 | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | nvidia | 3.1 | sim | sim |
| 4 | `openrouter/dots-studio/dots-3-note-preview:free` | openrouter | 5.2 | sim | sim |
| 5 | `cf/@cf/meta/llama-4-scout-17b-16e-instruct` | cloudflare | 5.9 | sim | não provada |
| 6 | `mistral/ministral-3b-2512` | mistral | 7.1 | sim | sim |

### `t-chat-128k` - chat geral, contexto grande
| # | membro | provedor | 100k (s) | função | imagem |
|---|---|---|---|---|---|
| 1 | `mistral/codestral-latest` | mistral | 2.1 | NÃO | não provada |
| 2 | `gemini/gemini-3.5-flash-lite` | gemini | 2.8 | sim | não provada |
| 3 | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | nvidia | 3.1 | sim | sim |
| 4 | `openrouter/cohere/north-mini-code:free` | openrouter | 4.9 | NÃO | não provada |
| 5 | `cf/@cf/meta/llama-4-scout-17b-16e-instruct` | cloudflare | 5.9 | sim | não provada |
| 6 | `mistral/ministral-3b-latest` | mistral | 2.4 | sim | sim |
Aviso: membros com função 'NÃO' respondem sem `tool_calls`; não usar este combo com agentes que exigem ferramentas (use `t-agente-128k`).

### `t-raciocinio-128k` - raciocínio (capacidade provada; qualidade NÃO provada)
| # | membro | provedor | 100k (s) | função | imagem |
|---|---|---|---|---|---|
| 1 | `gemini/gemini-3.5-flash-lite` | gemini | 2.8 | sim | não provada |
| 2 | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | nvidia | 3.1 | sim | sim |
| 3 | `openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | openrouter | 6.3 | sim | sim |
| 4 | `cf/@cf/openai/gpt-oss-20b` | cloudflare | 9.6 | NÃO | não provada |
| 5 | `gemini/gemini-3-flash-preview` | gemini | 4.1 | sim | não provada |
| 6 | `cf/@cf/openai/gpt-oss-120b` | cloudflare | 9.8 | sim | não provada |

### `t-visao-128k` - imagem (cor de PNG sintético respondida corretamente)
| # | membro | provedor | 100k (s) | função | imagem |
|---|---|---|---|---|---|
| 1 | `mistral/ministral-3b-latest` | mistral | 2.4 | sim | sim |
| 2 | `nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning` | nvidia | 3.1 | sim | sim |
| 3 | `openrouter/dots-studio/dots-3-note-preview:free` | openrouter | 5.2 | sim | sim |
| 4 | `mistral/ministral-3b-2512` | mistral | 7.1 | sim | sim |
| 5 | `nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b` | nvidia | 3.5 | sim | sim |
| 6 | `openrouter/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning:free` | openrouter | 6.3 | sim | sim |

### `t-codigo-128k` - código sem ferramentas (qualidade NÃO provada)
| # | membro | provedor | 100k (s) | função | imagem |
|---|---|---|---|---|---|
| 1 | `mistral/codestral-latest` | mistral | 2.1 | NÃO | não provada |
| 2 | `openrouter/cohere/north-mini-code:free` | openrouter | 4.9 | NÃO | não provada |
| 3 | `mistral/codestral-2508` | mistral | 5.9 | NÃO | não provada |
| 4 | `mistral/mistral-code-latest` | mistral | n/d (só reteste) | não provada | não provada |
Aviso: membros com função 'NÃO' respondem sem `tool_calls`; não usar este combo com agentes que exigem ferramentas (use `t-agente-128k`).

### `t-micro-8k` - só pedidos pequenos (Groq)
| # | membro | curto | 20k |
|---|---|---|---|
| 1 | `groq/openai/gpt-oss-120b` | HTTP 200 em 1.6s | HTTP 413 (TPM 8.000) |
| 2 | `groq/openai/gpt-oss-20b` | HTTP 200 em 2.0s | HTTP 413 (TPM 8.000) |
Regra: nunca receber pedido de harness. Teto provado: falha já em ~20k tokens.

### `m-gemini-4-hydrogen` - ponto focal em modelo
| # | membro | prova |
|---|---|---|
| 1 | `gemini/gemini-4-hydrogen` | 100k: HTTP 200 em 15,9 s (passo 2); mas no log da VPS (15:00-15:18 UTC) só 6 de 57 chamadas ao combo terminaram no Hydrogen, as demais caíram em 429 de todas as 16 contas; reteste via 9Router: 200 em 2 de 9 chamadas (curto 200/200/503, 20k e 100k 503 em 6 de 6); direto no Google, só DEV_06 e DEV_13 deram 200 e as outras 7 chaves vivas deram 429 com limit=20 |
Sem membros de reserva de outra família: se esgotar, falha rápida e o cliente escolhe outro combo. Esforço (low/medium/high) segue como parâmetro `reasoning_effort` (provado: 235 vs 445 tokens de raciocínio).

### `t-embedding` - vetores
| # | membro | prova |
|---|---|---|
| 1 | `gemini/gemini-embedding-2` | `/v1/embeddings`: HTTP 200 em 3 de 3 chamadas (reteste) |
Só 1 membro provado: sem reserva. Outros membros de embedding: não testados.

### Reprovados no reteste (3 amostras por tamanho)
| membro | resultado (curto, 20k, 100k) |
|---|---|
| `cloudflare-ai/qwen3.8-27b` | curto=400/400/400, 20k=400/400/400, 100k=400/400/400 |
| `cf/@cf/meta/llama-3.3-70b-instruct-fp8-fast` | curto=200/200/200, 20k=200/200/200, 100k=413/413/413 |
| `cloudflare-ai/gemma-4-26b-a4b-it` | curto=400/400/400, 20k=400/400/400, 100k=400/400/400 |
| `mistral/mistral-large-4` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
| `cloudflare-ai/llama-3.2-11b-vision-instruct` | curto=400/400/400, 20k=400/400/400, 100k=400/400/400 |
| `cf/@cf/meta/llama-3.1-70b-instruct-fp8-fast` | curto=200/200/200, 20k=200/200/200, 100k=413/413/413 |
| `gemini/gemma-4-26b-a4b-it` | curto=200/200/200, 20k=200/200/200, 100k=503/503/503 |
| `mistral/mistral-large-4-0` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
| `gemini/gemini-3.8-flash` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
| `gemini/gemma-4-31b-it` | curto=200/0/200, 20k=0/200/200, 100k=503/503/503 |
| `gemini/gemini-3.1-pro-preview` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
| `gemini/gemini-3.6-flash` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
| `gemini/gemini-4-hydrogen` | curto=200/200/503, 20k=503/503/503, 100k=503/503/503 |
| `gemini/gemini-4-argon` | curto=503/503/503, 20k=503/503/503, 100k=503/503/503 |
Leitura: 413 = limite de tamanho do provedor (prova de capacidade); 503 em 9 de 9 = indisponível durante o teste (pode ser momentâneo); 400 nos `cloudflare-ai/*` = prefixo ou rota inválidos (causa não verificada).

## Combos atuais no 9Router e destino sugerido
| combo atual | membros | membros reprovados pelas provas | destino sugerido |
|---|---|---|---|
| `code-pro` | 4 | cf/@cf/meta/llama-3.3-70b-instruct-fp8-fast (reteste 3x: curto=200/200/200, 20k=200/200/200, 100k=413/413/413); groq/openai/gpt-oss-120b (413) | substituir |
| `code-fast` | 4 | groq/openai/gpt-oss-120b (413) | substituir |
| `code-free` | 3 | groq/openai/gpt-oss-120b (413); cf/@cf/meta/llama-3.1-70b-instruct-fp8-fast (reteste 3x: curto=200/200/200, 20k=200/200/200, 100k=413/413/413) | substituir |
| `task-micro` | 3 | groq/openai/gpt-oss-120b (413) | substituir |
| `code-tests` | 3 | - | manter |
| `megacombo-coding` | 4 | gemini/gemma-4-26b-a4b-it (reteste 3x: curto=200/200/200, 20k=200/200/200, 100k=503/503/503); gemini/gemini-3.8-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); cloudflare-ai/qwen3.8-27b (reteste 3x: curto=400/400/400, 20k=400/400/400, 100k=400/400/400); gemini/gemma-4-31b-it (reteste 3x: curto=200/0/200, 20k=0/200/200, 100k=503/503/503) | substituir |
| `megacombo-reasoning` | 3 | gemini/gemini-3.1-pro-preview (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.6-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); cloudflare-ai/gemma-4-26b-a4b-it (reteste 3x: curto=400/400/400, 20k=400/400/400, 100k=400/400/400) | substituir |
| `megacombo-context` | 4 | gemini/gemini-3.6-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.8-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
| `megacombo-embedding` | 1 | - | manter |
| `megacombo-vision` | 3 | gemini/gemini-3.6-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.8-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); cloudflare-ai/llama-3.2-11b-vision-instruct (reteste 3x: curto=400/400/400, 20k=400/400/400, 100k=400/400/400) | substituir |
| `free-chat-M` | 5 | groq/qwen/qwen3.8-27b (413) | substituir |
| `free-chat-L` | 5 | - | manter |
| `free-chat-XL` | 5 | mistral/mistral-large-4 (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
| `free-codigo-L` | 3 | - | manter |
| `free-raciocinio-M` | 5 | groq/openai/gpt-oss-120b (413); groq/openai/gpt-oss-20b (413) | substituir |
| `free-raciocinio-XL` | 4 | mistral/mistral-large-4 (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); mistral/mistral-large-4-0 (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
| `free-ferramentas-L` | 5 | gemini/gemini-3.6-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); mistral/mistral-large-4 (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
| `free-visao-L` | 3 | - | manter |
| `free-ferramentas-M` | 2 | groq/openai/gpt-oss-120b (413) | substituir |
| `gemini-4-hydrogen` | 3 | gemini/gemini-4-hydrogen (reteste 3x: curto=200/200/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.8-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.6-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
| `gemini-4-argon` | 3 | gemini/gemini-4-argon (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.1-pro-preview (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503); gemini/gemini-3.8-flash (reteste 3x: curto=503/503/503, 20k=503/503/503, 100k=503/503/503) | substituir |
