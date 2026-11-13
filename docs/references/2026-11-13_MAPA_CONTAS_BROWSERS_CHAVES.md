# Mapa de contas, navegadores e chaves (Google AI Studio e Groq)

Data: 13/11/2026. Fonte: leituras reais feitas nesta data. Nenhum valor de chave é gravado aqui, só os 4 últimos caracteres que o próprio console mostra.
Regra: só entra aqui o que foi lido ou testado. O que falta fica em "Não verificado".

## 1. Como o Claude controla o Chrome

O Claude usa a extensão "Claude in Chrome". A porta 9222 não funciona neste computador.
Prova: `http://127.0.0.1:9222/json/version` deu "Impossível conectar-se ao servidor remoto" e a opção `chrome://inspect/#remote-debugging` ficou em "starting…".

A extensão só controla o Chrome (perfil) onde ela está instalada e logada. Ela abre abas num grupo próprio e não vê as abas que você abriu.

Os nomes e IDs dos navegadores NÃO são estáveis:

| Momento | "Browser 1" | "Browser 2" |
|---|---|---|
| Primeira leitura | 827a2b98 (conta Heverton) | não existia |
| Após instalar na DEV_06 | 4a41bcc3 (conta da DEV_06) | 827a2b98 (Heverton) |
| Após trocar para a conta seguinte | ced9e670 | 827a2b98 (Heverton) |

Conclusão provada: "Browser 1" muda de ID a cada reconexão, então o nome não identifica a conta.
Como identificar sem perguntar: depois de selecionar o navegador, abra `https://aistudio.google.com/api-keys` e compare os 4 últimos caracteres das chaves com a seção 2.

Procedimento para trocar de conta:
1. `list_connected_browsers`.
2. Se houver mais de um navegador, escolher pelo ID mais recente (`connectedAt` maior) e conferir com a seção 2.
3. `select_browser` com o `deviceId`, depois `tabs_context_mcp` com `createIfEmpty`.
4. Ler `https://aistudio.google.com/api-keys` e comparar.

Observação: a leitura dessa página disparou o alerta do hook `leak-scan` sobre `CONTABO_SECRET_USER` nas 2 leituras com a lista de chaves da conta Heverton (nas contas seguintes não disparou). O valor não foi gravado. Rotação: decisão do usuário.

## 2. Contas Google e projetos (lidos em aistudio.google.com/api-keys)

| Conta no Chrome | Projeto | Chaves (4 finais) | Nível mostrado | Variável no .env-full |
|---|---|---|---|---|
| Heverton | gen-lang-client-0806139386 | OLlQ (9router), rufA (novela-frutas), j3vQ (omniroute-heverton) | Nível gratuito | nenhuma |
| Heverton (criada por "Zahara") | gen-lang-client-0286004490 | gFOA (Default Gemini Key) | Nível gratuito | nenhuma |
| DEV_06 | gen-lang-client-0880104697 | G0XA (Gemini API Key 2) | Teste sem custo financeiro, "Configurar o pré-pagamento" | DEV_06_GEMINI_API_KEY |
| DEV_06 (mesmo projeto) | gen-lang-client-0880104697 | 1j1w, ctfg, QN3w, eVXQ | Teste sem custo financeiro | nenhuma |
| DEV_08 (kr741963@gmail.com) | gen-lang-client-0392529727 | sJXQ, hVFA (9router), 7N7A (kr741963), JGow | Nível gratuito | nenhuma das 4 é a DEV_08_GEMINI_API_KEY |

Fatos provados:
- Chaves do mesmo projeto compartilham o projeto (3 chaves no 0806139386, 5 no 0880104697, 4 no 0392529727). A cota do Google é por projeto (doc oficial; não testado aqui).
- A conta aberta na 3ª leitura é a da DEV_08, segundo o usuário (lista `gw agy list`: DEV_08 = kr741963@gmail.com), e uma das chaves dela se chama "kr741963".
- Porém nenhuma das 4 chaves listadas é a `DEV_08_GEMINI_API_KEY` do `.env-full` (comparação pelos 4 finais, `mapa4.py`).
- A chave da DEV_08 é válida: deu 429 no Hydrogen, não 400 nem 403. Logo ela deve pertencer a um projeto que o AI Studio não lista (a página só mostra projetos importados). Não verificado.
- Só DEV_06 foi ligada a um projeto até agora.

## 2b. Contas Google por tag (fonte: `gw agy list`, captura do usuário em 13/11/2026)

| Tag | E-mail |
|---|---|
| DEV_01 | khenperes001@gmail.com |
| DEV_02 | khenperes002@gmail.com |
| DEV_03 | khenperes003@gmail.com |
| DEV_04 | khenperes005@gmail.com |
| DEV_05 | peresg.claudia@gmail.com |
| DEV_06 | authenticparfum@gmail.com |
| DEV_07 | liriosagro@gmail.com |
| DEV_08 | kr741963@gmail.com |
| DEV_09 | dudanick741@gmail.com |
| DEV_10 | khenperes2015@gmail.com |
| DEV_11 | heverton.developer@gmail.com |
| DEV_12 | nickoleperes08012009@gmail.com |
| DEV_13 | hevertoneduardoperes@gmail.com |

Corroboração: a DEV_06 (authenticparfum) tem no AI Studio uma chave chamada "authentic", e a DEV_08 (kr741963) tem uma chamada "kr741963".
Esta lista cobre DEV_01 a DEV_13. As variáveis DEV_14 e DEV_15 existem no `.env-full` sem e-mail mapeado aqui.

Não verificado: projeto das chaves DEV_01 a DEV_05 e DEV_07 a DEV_15.

## 3. Resultado por chave nos modelos Gemini 4 (chamada direta ao Google)

Fonte: `docs/reports/data/2026-11-13_chaves3_prova.json` (1 amostra por célula).

| Chaves | Hydrogen | Argon e Carbon |
|---|---|---|
| DEV_06, DEV_13 | 200 | 429 com `limit=0` |
| DEV_04, 05, 08, 09, 10, 11, 12 | 429 com `limit=20` | 429 com `limit=0` |
| DEV_01, 02, 03, 07, 14, 15 | 403 | 403 |

Hipótese do usuário (Google Cloud ativo por chave libera cota dos modelos novos): ainda não provada.
Contra-indício: a DEV_06 responde 200 no Hydrogen e o console a mostra "sem faturamento".

## 4. Groq: organização por chave

Fonte: `docs/reports/data/2026-11-13_projetos_orgs_prova.json` (HTTP 413 com `org_id` no corpo, 13 chaves).

| Organização (id completo) | Chaves |
|---|---|
| org_01m64nvbceet7bnm3skvga5xky | DEV_01 |
| org_01m64semzje60b60cdj3022gbf | DEV_02, DEV_04 |
| org_01m641t4cjes0v5s38hdsj0401 | DEV_03 |
| org_01m40s5j4ye1n89qyhdec13xxv | DEV_05 |
| org_01m62efkfrejhs16wmyjnf4xtc | DEV_06 |
| org_01m64fac4reyar59fd7s14wkfs | DEV_07 |
| org_01m667cxetet6bbeaa9wqvm41d | DEV_08 |
| org_01m66762zne4xr6qw3s2a49dpt | DEV_09, DEV_11, DEV_12, DEV_13 |
| org_01m6663c50ey9tgwh94gq64qw8 | DEV_10 |

Total: 9 organizações para 13 chaves. A cota do Groq (8.000 tokens por minuto) é por organização.
Não verificado: quantas dessas 13 chaves o 9Router usa de fato.

## 5. Pendências

- Ler a conta dona de DEV_08 e das demais chaves (o 9Router tem 16 conexões Gemini; no `.env-full` há 15 variáveis).
- Comparar o rótulo "Teste sem custo financeiro" com "Nível gratuito" (significado não verificado).
- Provar ou refutar a hipótese do Google Cloud ativo, com uma chave que dê 429 no Hydrogen.
