# PARTE IV — Prática: Receitas de Uso

## 14. Receita: ligar o 9Router

**Quando usar.** Você quer que o `agy` responda pelo 9Router da empresa.

1. Defina `NINEROUTER_KEY` no ambiente do usuário com a chave do 9Router.
2. Abra um terminal novo.
3. Rode `gw agy 9router`.
4. Confirme as três variáveis com `definida: True`.
5. Abra outro terminal novo e rode `agy -p "ok"`.

**Erro comum.** Rodar o passo 3 sem `NINEROUTER_KEY`. O `gw` falha com erro, porque não existe mais chave embutida.

## 15. Receita: voltar ao `agy` puro

**Quando usar.** Você quer o `agy` direto, com a conta Google logada e sem gateway.

1. Rode `gw agy off`. O comando apaga as três variáveis.
2. Feche o terminal e abra um novo.
3. Confira com `echo $env:AGY_LLM_GATEWAY_URL`. A saída deve vir vazia.

**Erro comum.** Testar no terminal que já estava aberto. Ele ainda tem as variáveis antigas.

## 16. Receita: usar outro provedor

**Quando usar.** Você tem a chave de um provedor com API compatível com OpenAI.

1. Guarde a chave em uma variável, por exemplo `MINHA_CHAVE`.
2. Defina `AGY_LLM_GATEWAY_URL` com a URL base do provedor, terminando em `/v1`.
3. Defina `AGY_LLM_GATEWAY_WIRE_PROTOCOL` como `openai`.
4. Defina `AGY_LLM_GATEWAY_API_KEY` com o conteúdo de `MINHA_CHAVE`.
5. Abra um terminal novo e teste.

Antes do passo 5, teste a URL e a chave com uma chamada à rota `/models`. Resposta HTTP 200 indica que a chave vale.

## 17. Receita: trocar a chave sem expor

**Quando usar.** A chave anterior vazou ou foi revogada.

1. Gere a chave nova no painel do gateway e revogue a antiga.
2. Guarde a chave nova em `MINHA_CHAVE_NOVA`.
3. Peça ao agente a troca (modelo da Parte II, Método C) ou edite a variável pela janela do sistema.
4. Abra um terminal novo e valide.

## 18. Receita: diagnosticar quando nada funciona

1. Confirme as três variáveis com `definida: True`.
2. Chame `/models` no gateway com a chave. Sem HTTP 200, o problema está na chave ou na URL.
3. Confira o protocolo: `openai` para Chat Completions e `genai` para Gemini.
4. Em máquina gerenciada, pergunte ao administrador sobre o `admin_settings.json`.
5. Em JetBrains ou Zed, confira `auth.type = gateway` e o campo `env`.
