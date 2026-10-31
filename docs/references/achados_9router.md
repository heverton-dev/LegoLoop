# Findings (verified 2026-10-03)

| Finding | Evidence | Rule |
|---|---|---|
| 9Router adds a fixed +2000 to `prompt_tokens`, `input_tokens`, `total_tokens` | Ollama log: 11 real vs 2011 reported; normalizer in the server bundle adds `2e3` | Real input = reported - 2000; dashboard cost is inflated |
| Caveman Ultra + Ponytail Ultra cost ~757 real input tokens per call | Ollama: 768 with savers, 12 without | Ponytail off; saver off on short calls |
| Kiro (`kr/*`) injects 3k-20k hidden input tokens | bench: 4k-19k real input for a 200-token prompt | Use Kiro where free credits matter more than tokens |
| Groq free tier: 413 above 8000 tokens/min, no fallback | Claude Code prompt alone exceeds it | Never in harness combos |
| Reasoning Gemini models (`gemini/*`, `ag/gemini-3.8-flash`) spend 5k-15k thinking tokens and stop at `max_tokens` without code | bench hard task | Use `ag/gemini-3.8-flash-low` |
| Default response is SSE | first test returned `data:` lines | Always `"stream": false` in scripts |
| Orca refuses Claude launch with Anthropic auth env when managed accounts are on | `This Claude launch defines explicit Anthropic auth environment variables` | Wrapper sets them in the child process |
| One invalid provider disables all custom providers in omp | `models.yml validation failed - custom providers disabled` | Local Ollama needs `auth: none` and `api` |
| `components sync` left files removed from the source in harness copies | verify flagged orphans after a clean sync | Fixed: sync now removes them inside the component folder |
| Reading Orca settings through the accessibility tree prints env values (the key) | key shown in tool output | Read Orca settings with `orca_9router.py --estado`, which masks the key |
| Gemini v1beta Schema: rejects empty `object` nodes, array `properties`, and extra JSON Schema fields | API Error 400: `Invalid value at 'tools[0].function_declarations[97].parameters.properties[0]'` | Recursive sanitizer in chunk 6922 converts empty child `object` to string and normalizes properties map |
| Context overflow on Kiro (`kr/*`): fails with 400 `CONTENT_LENGTH_EXCEEDS_THRESHOLD` above 100k tokens | API Error 400: `Input is too long` on large Claude Code sessions | Combo detects >100k tokens and auto-prioritizes 1M Gemini Flash models, keeping Kiro as last resort |

# Re-verificação na VPS (2026-10-31, `https://9router.vpsconexao.org`)

Método: mesma mensagem (curta e longa) enviada ao Ollama local `llama3.1:8b` (contagem real, `prompt_eval_count`) e ao 9Router em `cf/@cf/meta/llama-3.1-8b-instruct-fp8-fast` (mesmo modelo/tokenizador Llama 3).

| Rota | Curta | Longa | Δ longa−curta | Excesso fixo sobre o real |
|---|---|---|---|---|
| Real (Ollama) | 25 | 899 | 874 | — |
| `cf/…llama-3.1-8b` | 2430 | 3304 | **874** | **+2405** |
| `cf/…llama-3.3-70b` | 2430 | 3304 | 874 | +2405 |
| `code-fast` (Gemini) | 2411 | 3125 | 714 (outro tokenizador) | ≈ +2400 |
| `code-pro` (Claude Sonnet 4.5) | 6559 | 7531 | 972 (outro tokenizador) | ≈ +6550 |

| Conclusão | Regra no LegoLoop |
|---|---|
| Contagem por token do 9Router é exata (Δ idêntico) | — |
| Excesso fixo = 2000 do normalizador (falso) + ~405 reais (saver Caveman + template) | `Input = reportado − 2000`; os ~405 reais permanecem na conta; bruto em `tokens.reported_input` |
| `code-pro` tem ~4100 tokens reais extras (injeção do Kiro) | preferir `code-fast`/`code-free` em nós baratos |
| Fora do ar na medição: `cerebras/llama-3.3-70b` (503), `cf/…awq` (410), `code-free` (timeout 180 s) | combos/rotas instáveis ⇒ usar modo Sugerido com fallback |
