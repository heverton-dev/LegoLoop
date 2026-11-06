# Auditoria de harnesses e do 9Router (2026-11-06)

> Relatório gerado a partir de chamadas reais, registradas uma a uma em `docs/reports/data/2026-11-06_auditoria/`. Nada aqui vem de memória ou de documentação: cada afirmação cita o tipo de evidência.

## 1. Resumo

### 1.1 Harnesses
```
HARNESS       HEADLESS  MODELOS QUE FUNCIONAM               EFFORT
Claude Code   SIM       13 de 19 testados                   aceita low, high, xhigh, max; efeito visível só em alguns
Antigravity   SIM       4 de 7 (os 4 Gemini)                no id base: low, medium, high; xhigh e max recusados
MimoCode      SIM       10 de 31 testados                   aceito em 16 de 18 casos
OpenCode      SIM       20 de 68 (3 na sua configuração)    aceito em 32 de 32 casos
```

### 1.2 Combos do 9Router
```
COMBO                 VEREDITO                  MEMBROS OK  ATENDIDO POR
code-pro              FUNCIONA, COM DEFEITOS    3/5         claude-opus-4-6-thinking
code-fast             FUNCIONA, COM DEFEITOS    3/5         gemini-3.8-flash-n
code-free             FUNCIONA, COM DEFEITOS    2/3         agnes-3.0-flash
code-tests            FUNCIONA                  3/3         qwen3-coder-next
task-micro            FUNCIONA                  3/3         openai/gpt-oss-120b
megacombo-coding      FUNCIONA, COM DEFEITOS    2/4         gemma-4-26b-a4b-it
megacombo-context     FUNCIONA, COM DEFEITOS    3/4         gemini-3.6-flash
megacombo-embedding   NÃO CHAMÁVEL PELO NOME    1/1         -
megacombo-reasoning   FUNCIONA, COM DEFEITOS    1/3         gemini-3.6-flash
megacombo-vision      FUNCIONA, COM DEFEITOS    1/3         gemini-3.6-flash
```

### 1.3 O que mais importa

- **Headless:** os quatro harnesses rodam de uma tacada e terminam sozinhos (`claude -p`, `agy -p`, `mimo run`, `opencode run`). O `mimo` e o `opencode` travam se o stdin ficar aberto: rode com stdin fechado.
- **A sua configuração custa caro:** no `opencode`, só 3 dos 20 modelos que funcionam com configuração limpa seguem funcionando com a sua. **Observado e guardado** (`data/sondas/`): `opencode/beam-free` recusa com a sua configuração por excesso de ferramentas ("maximum number of items is 128") e funciona com configuração vazia. **Visto na tela, sem registro guardado:** `crouton-alpha-free` recusou por schema de ferramenta profundo demais. Os demais estouram o tempo, esbarram no filtro de conteúdo ou dão erro sem mensagem. **Inferido:** a sua configuração carrega 2 MCPs (`supabase` e `rea`); a causa exata de cada falha não foi isolada. No `mimo`, cada pergunta carrega ~110 mil tokens de entrada, mesmo com configuração vazia, então esse peso é do próprio `mimo`.
- **Combos:** todos respondem, mas vários têm membros mortos. `gemini/gemini-3.8-flash` está fora (403/429) e aparece em 5 combos; três membros `cloudflare-ai/*` não existem no Cloudflare; `ag/gemini-3.1-pro-low` e `ag/gemini-3.7-flash-low` dão 404. O fallback funciona, mas custa segundos a minutos.
- **`megacombo-embedding` não é chamável pelo nome** em `/v1/embeddings` (HTTP 400), embora o membro `gemini/gemini-embedding-2` responda.
- **`gw agy exec` troca de conta sozinho em erro 429.** Toda chamada `agy` (`agy.cmd` = `gw agy exec`) pode cair na DEV_04 (relogin) ou na DEV_05 (desativada por violação dos termos) depois de um erro de cota.

## 2. Método e limites

```
DATA                06/11/2026 (manhã, horário de Brasília)
PROMPTS             ola: 'Responda apenas a palavra: OK'   conta: 'Quanto e 17 vezes 23?'   codigo: função Soma em Go   visao: PNG 16x16 vermelho
LIMITE              60 s; se estourar, uma repetição com 180 s (resposta entre 60 e 180 s = lento; sem resposta em 180 s = falha)
ASSINATURAS         claude e agy usaram as assinaturas, a pedido; opencode e mimo só modelos gratuitos (custo 0 no catálogo) com credencial configurada
CONFIGURAÇÃO        opencode e mimo medidos em duas condições: LIMPA (config vazia) e SUA (config real do usuário)
9ROUTER             por HTTP, no servidor remoto; a composição dos combos veio do banco local de 4/11 (instantâneo), validada contra o remoto
CHAVES              nenhuma chave, token ou e-mail entra nos dados (URLs e e-mails são ocultados na cópia)
```

**Limites que importam:**
- Prompts mínimos medem *se responde*, não qualidade. O efeito de effort foi medido só pelos tokens de raciocínio numa conta simples.
- O nível `medium` de effort não foi testado no `claude`, e o `xhigh` só nos modelos que o declaram.
- Um modelo que 'não respondeu no limite' pode ser lento, não morto; o tempo vem registrado.
- O banco local de combos tem 10 combos; o servidor remoto lista 8 (sem `task-micro` e `code-tests`).
- Os tempos dos membros do 9Router foram medidos com 4 chamadas em paralelo; membros do mesmo provedor se atrapalham e inflam o tempo (ex.: `ag/gemini-3.8-flash-low` aparece com 171 s, mas o combo que o contém respondeu em 63 s).
- A DEV_03 foi a conta do Antigravity em todos os testes finais; a cota dos modelos Claude e GPT dela estava esgotada, o que limita a prova deles.

## 3. Parte 1: harnesses

### 3.1 Claude Code 2.1.316 (assinatura claude.ai)

**Headless: SIM.** `claude -p "<prompt>" --model <m> --effort <n> --output-format json --no-session-persistence` (provado por chamada).

**Modelos** (rota oficial; o `modelo real` é o que o servidor reportou):
```
MODELO                      FUNCIONA  TEMPO    MODELO REAL                 OBSERVAÇÃO
claude-fable-5              NÃO       6.2 s    -                           exige créditos de uso (fora da assinatura)
claude-fable-5-1            NÃO       5.4 s    -                           exige créditos de uso (fora da assinatura)
claude-haiku-4-5-20251001   sim       5.9 s    claude-haiku-4-5-20251001   
claude-haiku-5-5            sim       6.6 s    claude-haiku-5-5            
claude-opus-4-6             sim       6.8 s    claude-opus-4-6             
claude-opus-4-7             sim       5.9 s    claude-opus-4-7             
claude-opus-4-8             sim       6.2 s    claude-opus-4-8             
claude-opus-5               sim       6.4 s    claude-opus-5               
claude-opus-5-5             sim       6.4 s    claude-opus-5-5             
claude-sonnet-4-6           sim       6.4 s    claude-sonnet-4-6           
claude-sonnet-5             sim       6.4 s    claude-sonnet-5             
claude-sonnet-5-5           sim       6.9 s    claude-sonnet-5-5           
code-fast                   NÃO       5.4 s    -                           modelo inexistente nesta rota
code-free                   NÃO       5.4 s    -                           modelo inexistente nesta rota
code-pro                    NÃO       5.4 s    -                           modelo inexistente nesta rota
fable                       NÃO       5.3 s    -                           exige créditos de uso (fora da assinatura)
haiku                       sim       7.2 s    claude-haiku-5-5            
opus                        sim       7.1 s    claude-opus-5-5             
sonnet                      sim       7.1 s    claude-sonnet-5-5           
```

- `fable-*` exigem crédito de uso, fora da assinatura.
- `code-fast`, `code-free` e `code-pro` **não existem na rota oficial**: só existem com `gw claude 9router --combo <nome>`.
- Cada chamada carrega ~37 mil tokens de contexto do seu perfil.

**Effort** (`--effort low|medium|high|xhigh|max`; todos os níveis testados foram aceitos):
```
MODELO                      DECLARADOS (catálogo)             TESTADOS -> tokens de raciocínio
claude-haiku-4-5-20251001   nenhum                            low:95
claude-haiku-5-5            low,medium,high,xhigh,max         low:26  high:66  xhigh:64  max:153
claude-opus-4-6             low,medium,high,max               low:0  high:0  max:0
claude-opus-4-7             low,medium,high,xhigh,max         low:0  high:0  xhigh:0  max:0
claude-opus-4-8             low,medium,high,xhigh,max         low:0  high:0  xhigh:0  max:21
claude-opus-5               low,medium,high,xhigh,max         low:0  high:0  xhigh:0  max:0
claude-opus-5-5             low,medium,high,xhigh,max         low:15  high:0  xhigh:0  max:52
claude-sonnet-4-6           low,medium,high,max               low:0  high:0  max:0
claude-sonnet-5             low,medium,high,xhigh,max         low:0  high:0  xhigh:0  max:0
claude-sonnet-5-5           low,medium,high,xhigh,max         low:0  high:0  xhigh:0  max:0
```

Efeito visível nos tokens de raciocínio: `haiku-5-5` (26, 66, 64, 153) e `opus-5-5` (15, 0, 0, 52). Nos demais, 0 em qualquer nível: o raciocínio adaptativo decidiu não pensar numa conta fácil. Isso não prova que o effort não faz nada.

### 3.2 Antigravity (agy 1.5.2, assinatura Google AI Pro)

**Headless: SIM** (`agy -p`, JSON, termina sozinho). **Interativo: SIM** (terminal do Orca). Ambos provados com a conta DEV_03.

```
MODELO                      HEADLESS  TEMPO    OBSERVAÇÃO
gemini-3.8-flash            OK        14.1 s   
gemini-3.8-flash-high       OK        18.4 s   
gemini-3.8-flash-medium     OK        40.3 s   
gemini-3.8-flash-low        OK        13.9 s   
gpt-oss-120b-medium         FALHOU    8.9 s    cota esgotada (429)
claude-opus-4-6-thinking    FALHOU    10.6 s   cota esgotada (429)
claude-sonnet-4-6           FALHOU    11.7 s   cota esgotada (429)
```

- `claude-*` e `gpt-oss` via Antigravity: cota individual da conta esgotada (reinicia em ~1h09). É limite da conta, não do modelo.

**Effort:**
- Nos Gemini `-low/-medium/-high` o effort **vem embutido no nome**; `--effort` junto com eles é recusado localmente (`conflicts with --effort`).
- No id base `gemini-3.8-flash`, `--effort` aceita:
```
NÍVEL (--effort)  RESULTADO   TEMPO    RACIOCÍNIO (tokens) / MOTIVO
low               aceito      10.1 s   56 (resposta "391")
medium            aceito      11.8 s   74 (resposta "391")
high              aceito      12.9 s   61 (resposta "391")
xhigh             RECUSADO    3.0 s    gemini-3.8-flash has no "xhigh" effort (available: low, medium, high)
max               RECUSADO    3.1 s    gemini-3.8-flash has no "max" effort (available: low, medium, high)
```

**Interativo:**
```
MODELO (interativo)         EFFORT  RESULTADO  TEMPO
gemini-3.8-flash-low        -       OK (391)   12 s
gemini-3.8-flash-medium     -       OK (391)   12 s
gemini-3.8-flash-high       -       OK (391)   11 s
gemini-3.8-flash            high    OK (391)   11 s
```

**Atenção às contas:** `agy.cmd`, `agy-gw.cmd` e `antigravity.cmd` são todos `gw agy exec`, que faz **failover automático em erro 429** para a próxima conta. Na rodada inicial isso levou a DEV_04 (precisa de relogin: erro 'Verify your account') e a DEV_05 (desativada por violação dos termos). Os testes finais usaram o binário real, com a DEV_03 fixa.

### 3.3 MimoCode 0.1.16

**Headless: SIM** (`mimo run -m <provedor/modelo> --variant <nivel> --format json`, com stdin fechado). Modelos gratuitos com credencial: nvidia (37) e zenmux (6); `xiaomi-token-plan` fora (plano pago).

```
MODELO                                                LIMPA  TEMPO    SUA CONFIGURAÇÃO
nvidia/google/diffusiongemma-26b-a4b-it               OK     40.5 s   OK
nvidia/google/gemma-4-31b-it                          OK     58.5 s   OK
nvidia/meta/muse-glimmer-30b                          OK     28.6 s   OK
nvidia/minimaxai/minimax-m3.1-flash                   OK     19.0 s   OK
nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning  OK     20.2 s   OK
nvidia/nvidia/nemotron-3-ultra-550b-a55b              OK     27.1 s   OK
nvidia/nvidia/nemotron-3.5-lightning-30b-a3b          OK     23.2 s   OK
nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b         OK     25.6 s   OK
nvidia/poolside/laguna-xs-2.1                         OK     15.8 s   OK
nvidia/stepfun-ai/step-5-preview                      OK     22.0 s   FALHOU: erro sem mensagem
```

**Por que os outros falham (config limpa):**
```
 13  sem resposta no limite
  3  erro sem mensagem
  2  Bad Request: "auto" tool choice requires --enable-
  2  acesso negado (403)
  1  modelo removido (410)
```

**Effort** (`--variant`): aceito em 16 de 18 casos
```
MODELO                                                DECLARA                                       TESTADO -> raciocínio (tokens)
nvidia/google/diffusiongemma-26b-a4b-it               low,medium,high                               high:0  low:FALHOU
nvidia/google/gemma-4-31b-it                          low,medium,high                               low:138  high:136
nvidia/meta/muse-glimmer-30b                          low,medium,high                               low:0  high:0
nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning  low,medium,high                               low:49  high:46
nvidia/nvidia/nemotron-3-ultra-550b-a55b              low,medium,high                               low:50  high:68
nvidia/nvidia/nemotron-3.5-lightning-30b-a3b          low,medium,high                               low:105  high:99
nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b         low,medium,high                               low:70  high:107
nvidia/poolside/laguna-xs-2.1                         low,medium,high                               low:0  high:FALHOU
nvidia/stepfun-ai/step-5-preview                      low,medium,high                               high:0  low:0
```

- Peso fixo: ~110 mil tokens de entrada por chamada, mesmo com configuração vazia.
- Latência alta (15 a 60 s), dominada pelos endpoints NVIDIA.

### 3.4 OpenCode 1.19.1

**Headless: SIM** (`opencode run -m <provedor/modelo> --variant <nivel> --format json`, com stdin fechado). Provedores com credencial e modelos gratuitos: openrouter, nvidia, opencode, zenmux, groq, freellmapi, mimo, omniroute.

```
MODELO                                                LIMPA  TEMPO    SUA CONFIGURAÇÃO
nvidia/google/diffusiongemma-26b-a4b-it               OK     18.0 s   FALHOU: sem resposta no limite
nvidia/google/gemma-4-31b-it                          OK     48.1 s   FALHOU: input rejected by content safety screening
nvidia/meta/muse-glimmer-30b                          OK     18.7 s   FALHOU: {"message":"Internal server error","type":"interna
nvidia/minimaxai/minimax-m3.1-flash                   OK     6.7 s    FALHOU: sem resposta no limite
nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning  OK     6.9 s    FALHOU: sem resposta no limite
nvidia/nvidia/nemotron-3-ultra-550b-a55b              OK     9.2 s    OK
nvidia/nvidia/nemotron-3.5-lightning-30b-a3b          OK     8.2 s    FALHOU: sem resposta no limite
nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b         OK     10.3 s   OK
nvidia/poolside/laguna-xs-2.1                         OK     8.5 s    FALHOU: Bad Request: Invalid request
nvidia/stepfun-ai/step-5-preview                      OK     10.7 s   OK
opencode/beam-free                                    OK     10.2 s   FALHOU: erro sem mensagem
opencode/big-pickle                                   OK     8.4 s    FALHOU: sem resposta no limite
opencode/crouton-alpha-free                           OK     9.4 s    FALHOU: erro sem mensagem
opencode/deepseek-v4.1-flash-free                     OK     8.7 s    FALHOU: erro sem mensagem
opencode/flint-40b-free                               OK     8.6 s    FALHOU: sem resposta no limite
opencode/glm-5.3-flash-free                           OK     11.3 s   FALHOU: sem resposta no limite
opencode/mimo-v2.6-flash-free                         OK     42.6 s   FALHOU: sem resposta no limite
opencode/nemotron-3.5-lightning-free                  OK     12.6 s   FALHOU: sem resposta no limite
opencode/nemotron-3.5-super-vl-free                   OK     10.2 s   FALHOU: sem resposta no limite
opencode/pulse-alpha-free                             OK     25.7 s   FALHOU: sem resposta no limite
```

**Por que os outros falham (config limpa):**
```
 29  modelo não encontrado (404)
  6  Unexpected server error. Check server logs for det
  4  sem resposta no limite
  3  Invalid API Key
  3  acesso negado (403)
  2  Bad Request: "auto" tool choice requires --enable-
  1  modelo removido (410)
```

- 29 dos 68 candidatos gratuitos (segundo o catálogo) **não existem mais no provedor**.
- Na **sua** configuração só 3 de 20 funcionam. Observado e guardado em `beam-free`: "maximum number of items is 128" (ferramentas). `crouton-alpha-free` mostrou "JSON schema exceeds the maximum nesting depth" na verificação do executor, sem registro guardado. Nos outros 15 as causas são tempo esgotado, filtro de conteúdo, erro do servidor ou erro sem mensagem; a ligação com os MCPs `supabase` e `rea` é inferência, não prova.

**Effort** (`--variant`): aceito em 32 de 32 casos
```
MODELO                                                DECLARA                                       TESTADO -> raciocínio (tokens)
nvidia/google/diffusiongemma-26b-a4b-it               none,minimal,low,medium,high,xhigh,max        none:0  max:0
nvidia/google/gemma-4-31b-it                          low,medium,high                               low:63  high:95
nvidia/meta/muse-glimmer-30b                          none,minimal,low,medium,high,max              none:0  max:0
nvidia/minimaxai/minimax-m3.1-flash                   low,medium,high,xhigh,max                     max:0  low:0
nvidia/nvidia/nemotron-3-nano-omni-30b-a3b-reasoning  low,medium,high                               low:34  high:233
nvidia/nvidia/nemotron-3-ultra-550b-a55b              none,medium,high                              none:0  high:26
nvidia/nvidia/nemotron-3.5-lightning-30b-a3b          low,medium,high                               low:74  high:86
nvidia/nvidia/nemotron-3.5-super-vl-120b-a12b         low,medium,high                               low:95  high:95
nvidia/stepfun-ai/step-5-preview                      none,minimal,low,medium,high,xhigh,max        max:0  none:0
opencode/beam-free                                    low,medium,high,xhigh,max                     low:70  max:64
opencode/crouton-alpha-free                           minimal,low,medium,high,xhigh,max             minimal:0  max:0
opencode/deepseek-v4.1-flash-free                     low,high,max                                  max:8  low:0
opencode/flint-40b-free                               low,medium,high,xhigh,max                     low:0  max:0
opencode/glm-5.3-flash-free                           low,high,max                                  max:23  low:62
opencode/nemotron-3.5-super-vl-free                   none,minimal,low,medium,high,xhigh            xhigh:114  none:0
opencode/pulse-alpha-free                             low,high,max                                  max:31  low:30
```

Efeito visível: `nemotron-3-nano-omni` (34 low, 233 high), `nemotron-3-ultra` (0 none, 26 high), `gemma-4-31b` (63, 95). Vários aceitam o nível e não mudam o raciocínio.

## 4. Parte 2: 9Router

### 4.1 Provedores que realmente funcionam (membros dos combos)
```
PROVEDOR        MEMBROS OK  FALHARAM  CAUSAS DAS FALHAS
ag              3/5         2         gemini-3.1-pro-low: modelo não encontrado (404); gemini-3.7-flash-low: modelo não encontrado (404)
agnes           1/1         0         
cloudflare-ai   0/3         3         qwen3.8-27b: modelo inexistente no provedor; gemma-4-26b-a4b-it: modelo inexistente no provedor; llama-3.2-11b-vision-instruct: modelo inexistente no provedor
gemini          6/8         2         gemini-3.1-pro-preview: cota esgotada (429); gemini-3.8-flash: cota esgotada (429)
groq            1/1         0         
kr              3/3         0         
```

### 4.2 Combos: composição, modelos que respondem e veredito

### code-pro  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 28.6 s   | atendido por: claude-opus-4-6-thinking
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     ag/gemini-3.1-pro-low                         NÃO       23.1 s   modelo não encontrado (404)
2     ag/claude-opus-4-6-thinking                   sim       3.1 s    
3     gemini/gemini-3.1-pro-preview                 NÃO       1.0 s    cota esgotada (429)
4     ag/claude-sonnet-4-6                          sim       2.1 s    
5     kr/claude-sonnet-4.5                          sim       1.9 s    
```

Capacidade: código: OK; conta: OK.
Defeitos: membros indisponíveis agora: ag/gemini-3.1-pro-low (modelo não encontrado (404)), gemini/gemini-3.1-pro-preview (cota esgotada (429)) | o 1º membro vivo é o nº 2; o anterior falha e o combo gasta 28.6 s para responder.

### code-fast  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 62.7 s   | atendido por: gemini-3.8-flash-n
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     ag/gemini-3.8-flash-low                       sim       171.4 s  atendido por gemini-3.8-flash-n
2     gemini/gemini-3.8-flash                       NÃO       1.3 s    cota esgotada (429)
3     ag/gemini-3.7-flash-low                       NÃO       23.6 s   modelo não encontrado (404)
4     ag/claude-sonnet-4-6                          sim       2.1 s    
5     kr/qwen3-coder-next                           sim       3.0 s    
```

Capacidade: código: OK; conta: OK.
Defeitos: membros indisponíveis agora: gemini/gemini-3.8-flash (cota esgotada (429)), ag/gemini-3.7-flash-low (modelo não encontrado (404)).

### code-free  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 135.7 s   | atendido por: agnes-3.0-flash
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemini-3.8-flash                       NÃO       1.3 s    cota esgotada (429)
2     agnes/agnes-3.0-flash                         sim       15.1 s   
3     kr/qwen3-coder-next                           sim       3.0 s    
```

Capacidade: código: OK; conta: OK.
Defeitos: membros indisponíveis agora: gemini/gemini-3.8-flash (cota esgotada (429)) | o 1º membro vivo é o nº 2; o anterior falha e o combo gasta 135.7 s para responder.

### code-tests  ->  FUNCIONA

```
no servidor remoto: NÃO (só no banco local de 4/11)   | respondeu: sim em 1.3 s   | atendido por: qwen3-coder-next
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     kr/qwen3-coder-next                           sim       3.0 s    
2     kr/claude-haiku-4.5                           sim       1.0 s    
3     ag/gemini-3.8-flash-low                       sim       171.4 s  atendido por gemini-3.8-flash-n
```

Capacidade: código: OK; conta: OK.
Sem defeitos.

### task-micro  ->  FUNCIONA

```
no servidor remoto: NÃO (só no banco local de 4/11)   | respondeu: sim em 1.5 s   | atendido por: openai/gpt-oss-120b
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     groq/openai/gpt-oss-120b                      sim       1.5 s    
2     agnes/agnes-3.0-flash                         sim       15.1 s   
3     kr/qwen3-coder-next                           sim       3.0 s    
```

Capacidade: código: OK; conta: OK.
Sem defeitos.

### megacombo-coding  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 3.4 s   | atendido por: gemma-4-26b-a4b-it
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemma-4-26b-a4b-it                     sim       6.7 s    
2     gemini/gemini-3.8-flash                       NÃO       1.3 s    cota esgotada (429)
3     cloudflare-ai/qwen3.8-27b                     NÃO       0.5 s    modelo inexistente no provedor
4     gemini/gemma-4-31b-it                         sim       28.7 s   
```

Capacidade: código: OK; conta: OK.
Defeitos: membros que não existem no provedor: cloudflare-ai/qwen3.8-27b | membros indisponíveis agora: gemini/gemini-3.8-flash (cota esgotada (429)).

### megacombo-context  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 1.6 s   | atendido por: gemini-3.6-flash
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemini-3.6-flash                       sim       1.5 s    
2     gemini/gemini-3.5-flash                       sim       1.9 s    atendido por gemini-3.6-flash
3     gemini/gemini-flash-lite-latest               sim       0.9 s    atendido por gemini-3.5-flash-lite
4     gemini/gemini-3.8-flash                       NÃO       1.3 s    cota esgotada (429)
```

Capacidade: código: OK; conta: OK.
Defeitos: membros indisponíveis agora: gemini/gemini-3.8-flash (cota esgotada (429)).

### megacombo-embedding  ->  NÃO CHAMÁVEL PELO NOME

```
no servidor remoto: sim   | respondeu: não em 0.2 s   | atendido por: -
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemini-embedding-2                     sim       0.8 s    
```

Capacidade: embedding: ver defeito.
Defeitos: o combo NÃO é chamável pelo nome em /v1/embeddings (HTTP 400 Invalid model format), embora o membro responda.

### megacombo-reasoning  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 9.8 s   | atendido por: gemini-3.6-flash
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemini-3.1-pro-preview                 NÃO       1.0 s    cota esgotada (429)
2     gemini/gemini-3.6-flash                       sim       1.5 s    
3     cloudflare-ai/gemma-4-26b-a4b-it              NÃO       0.6 s    modelo inexistente no provedor
```

Capacidade: código: OK; conta: OK.
Defeitos: membros que não existem no provedor: cloudflare-ai/gemma-4-26b-a4b-it | membros indisponíveis agora: gemini/gemini-3.1-pro-preview (cota esgotada (429)) | o 1º membro vivo é o nº 2; o anterior falha e o combo gasta 9.8 s para responder.

### megacombo-vision  ->  FUNCIONA, COM DEFEITOS

```
no servidor remoto: sim   | respondeu: sim em 1.9 s   | atendido por: gemini-3.6-flash
ORDEM MEMBRO                                        RESPONDE  TEMPO    OBSERVAÇÃO
1     gemini/gemini-3.6-flash                       sim       1.5 s    
2     gemini/gemini-3.8-flash                       NÃO       1.3 s    cota esgotada (429)
3     cloudflare-ai/llama-3.2-11b-vision-instruct   NÃO       0.5 s    modelo inexistente no provedor
```

Capacidade: visão: OK (imagem 16x16 vermelha reconhecida).
Defeitos: membros que não existem no provedor: cloudflare-ai/llama-3.2-11b-vision-instruct | membros indisponíveis agora: gemini/gemini-3.8-flash (cota esgotada (429)).

### 4.3 Os combos foram construídos corretamente?

Critérios: todo membro existe e responde; nenhum membro morto; ordem de fallback faz sentido; a capacidade combina com a finalidade.

```
COMBO                 VEREDITO                  MEMBROS OK  ATENDIDO POR
code-pro              FUNCIONA, COM DEFEITOS    3/5         claude-opus-4-6-thinking
code-fast             FUNCIONA, COM DEFEITOS    3/5         gemini-3.8-flash-n
code-free             FUNCIONA, COM DEFEITOS    2/3         agnes-3.0-flash
code-tests            FUNCIONA                  3/3         qwen3-coder-next
task-micro            FUNCIONA                  3/3         openai/gpt-oss-120b
megacombo-coding      FUNCIONA, COM DEFEITOS    2/4         gemma-4-26b-a4b-it
megacombo-context     FUNCIONA, COM DEFEITOS    3/4         gemini-3.6-flash
megacombo-embedding   NÃO CHAMÁVEL PELO NOME    1/1         -
megacombo-reasoning   FUNCIONA, COM DEFEITOS    1/3         gemini-3.6-flash
megacombo-vision      FUNCIONA, COM DEFEITOS    1/3         gemini-3.6-flash
```

**Defeitos de construção (não são falhas passageiras):**
- `cloudflare-ai/qwen3.8-27b`, `cloudflare-ai/gemma-4-26b-a4b-it` e `cloudflare-ai/llama-3.2-11b-vision-instruct` não existem no Cloudflare: o id está errado. Estão em `megacombo-coding`, `megacombo-reasoning` e `megacombo-vision`.
- `ag/gemini-3.1-pro-low` e `ag/gemini-3.7-flash-low` dão 404 no Antigravity (modelo retirado). Estão em `code-pro` e `code-fast`.
- `megacombo-embedding` não é chamável pelo nome em `/v1/embeddings`.

**Indisponibilidade passageira (cota ou acesso, no momento do teste):** `gemini/gemini-3.8-flash` e `gemini/gemini-3.1-pro-preview` (429 e 403).

### 4.4 Achado de protocolo
O 9Router responde `Content-Type: text/event-stream`, mas o corpo é um JSON único seguido de `data: [DONE]` colado; os combos respondem em fluxo `data:`. Clientes que decidem o parser pelo tipo de conteúdo podem falhar. Modelos de raciocínio devolvem `content` vazio se `max_tokens` for baixo (use 300 ou mais).

## 5. Riscos e achados transversais

- **`gw agy exec` rotaciona contas sozinho em 429** e pode cair em conta ruim. O `gw agy probe` mostra 'Disponível' para a DEV_05, que está desativada: ele só mede cota.
- **Sessão do Antigravity e testes disputam a mesma credencial** (Gerenciador de Credenciais do Windows, `gemini:antigravity`).
- **Memória:** `mimo` (~750 MB por processo) e `agy` (dois em paralelo derrubam a máquina) exigem concorrência baixa.
- O primeiro teste de visão do 9Router usou um PNG inválido e foi refeito; os dados antigos não entram aqui.
- **Estado que a auditoria deixou:** o `agy` ficou fixado na conta **DEV_03** (antes estava na DEV_04, que precisa de relogin). As worktrees descartáveis e os diretórios temporários criados pela auditoria são removidos ao final.

## 6. Próximo passo (fora desta auditoria)

A matriz `docs/reports/2026-11-06_MATRIZ_ROTEAMENTO.json` serve de insumo para corrigir o gerador da ROUTE (`scratch/gen-plan/gen.mjs`, que fixa `claude` + `sonnet`). Isso fica para depois da aprovação deste relatório.

## 7. Arquivos de dados

```
docs/reports/data/2026-11-06_auditoria/
  9router.jsonl  9router_visao.jsonl  combos.json  9router_modelos_remotos.json
  claude.jsonl  opencode.jsonl  mimo.jsonl  agy_headless.jsonl  agy_interativo.jsonl
  agy_rodada_conta_bloqueada.jsonl  agy_rodada_rotacao_gw.jsonl   (rodadas descartadas, guardadas como evidência)
  catalogo_opencode.json  catalogo_mimo.json  cand_claude.json
  scripts/   (r9aud, r9vision, hrun, hrun2, catalogo, cand_sub, combos_json, agy_inter, gerar)
```