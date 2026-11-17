# Diagnóstico: erro 413 do Groq nos combos do 9Router (13/11/2026)

## Sintoma
Harness com prompt de ~114k tokens falha com HTTP 413: `Request too large for model openai/gpt-oss-120b ... TPM: Limit 8000, Requested 114246`.
Aparece com `/model code-free` e com `/model Gemini 4 Hydrogen (Medium)`.

## Causa raiz (provada por teste ao vivo em https://roteador.signa-docs.com)
O erro não vem do Google AI Studio nem das 13 contas. Vem do plano gratuito do Groq.

1. O modelo `groq/openai/gpt-oss-120b` é o 1º membro de `code-free`, `task-micro`, `free-raciocinio-M` e `free-ferramentas-M`. É o 2º de `code-fast` e o 3º de `code-pro`.
2. O plano `on_demand` do Groq limita esse modelo a 8.000 tokens por minuto. Um único pedido acima disso é recusado com 413, mesmo que o contexto do modelo seja 131k.
3. Harnesses estilo Claude Code enviam 30k a 115k tokens por chamada (sistema, ferramentas, histórico). Todo pedido grande bate no 1º membro e é recusado.
4. O combo tenta os membros seguintes. Eles também falham (Cloudflare llama-3.1-70b tem janela curta; Gemini lite está com 429). O 9Router devolve ao cliente o erro do PRIMEIRO membro, o 413 do Groq, e esconde as falhas reais dos outros.

## Provas (sonda em scratchpad/sonda.py, pedido de ~98k tokens)
- `groq/openai/gpt-oss-120b` direto: HTTP 413, Limit 8000, Requested 97715, em 1,6 s.
- `code-free`: HTTP 413, mesmo texto e mesmo número (97715), em 4,9 s.
- `gemini/gemini-flash-lite-latest` direto: HTTP 503 com 429 de cota (último recurso do combo está esgotado).
- Pedido curto: `code-free`, `code-fast`, `code-pro` e `groq/gpt-oss-120b` respondem 200 em 1 a 7 s. Por isso o 9Router parece saudável em testes simples.
- Pedido de ~20k tokens em `code-free`: 200, atendido por gemini-3.5-flash-lite (Groq recusa, o failover funciona quando há membro vivo).

## Problema 2: combos Gemini 4 lentos
- `gemini-4-hydrogen` (definição: gemini-4-hydrogen, gemini-3.8-flash, gemini-3.6-flash): pedido curto levou 84 s e foi atendido pelo 3º membro (gemini-3.6-flash). Outra tentativa estourou 60 s e 120 s.
- `gemini/gemini-4-hydrogen` direto responde em 2,8 s com pedido curto. O atraso vem do combo passar por membros lentos ou em 429 antes de chegar a um vivo.
- `gemini-3.8-flash`: 429 de cota (503 no 9Router). `gemini-3.1-pro-preview`: 429. `gemini-4-argon`: 503 após 16,6 s.
- Não verificado: por que `gemini-4-hydrogen` combo não atende no 1º membro com pedido grande (teste de 114k estourou 120 s).

## Natureza do problema
Não é falha de API nem de chave. É erro de projeto dos combos: modelos com limite de TPM baixo (Groq free, 8k) ficaram na frente de rotas usadas por harnesses de contexto grande. O combo não filtra por tamanho do pedido e devolve o erro errado.
As 13 contas Google dão cota por conta, mas o limite de 8k TPM do Groq é da organização e não escala com elas.

## Correção sugerida (nenhuma aplicada; mexe no 9Router de produção)
1. Tirar `groq/openai/gpt-oss-120b` de `code-free`, `code-fast` e `code-pro`. Manter o Groq só em `task-micro` e `free-ferramentas-M` e em chamadas com menos de ~6k tokens.
2. Em `code-free`/`code-fast`/`code-pro`, pôr na frente só membros com janela e cota grandes: gemini-3-flash-preview, gemini-3.6-flash, mistral/codestral-latest (256k). Cloudflare llama-3.1/3.3-70b fica fora (janela curta).
3. `gemini-4-hydrogen`: pôr `gemini/gemini-3.6-flash` em 2º no lugar de `gemini-3.8-flash` (429 recorrente) e testar o 1º membro com pedido grande.
4. No harness: limitar o contexto enviado ao teto do combo e tratar 413 como "troque de combo", não como erro final.
5. Criar uma sonda de homologação por combo com 3 tamanhos (1k, 20k, 100k tokens) e registrar o tamanho máximo aceito. Hoje a sonda usa pedido curto e aprova combos que quebram com pedido real.
6. Pedir ao 9Router (ou aplicar no cliente) que, quando todos os membros falham, devolva o último erro ou a lista completa, não só o primeiro.
