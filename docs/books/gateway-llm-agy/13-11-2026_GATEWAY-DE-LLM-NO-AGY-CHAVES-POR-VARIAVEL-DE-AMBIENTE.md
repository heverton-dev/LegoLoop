---
title: "Gateway de LLM no Antigravity CLI: Chaves por Variável de Ambiente"
subtitle: "Como apontar o agy para o 9Router ou para qualquer outro destino, sem trocar de conta Google"
author:
  - Heverton & Ecossistema LegoLoop
date: "13/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "OPERAÇÃO SEGURA DE AGENTES"
cover-title: "Gateway de LLM"
cover-accent: "Chaves por Variável de Ambiente no AGY"
cover-caption: "OS TRÊS MÉTODOS DE CONFIGURAÇÃO"
cover-stages:
  - Entender o Risco
  - Escolher o Destino
  - Configurar pelo Sistema
  - Configurar pelo Terminal
  - Pedir ao Agente
  - Validar e Desligar
header-mark: "GATEWAY AGY"
accent: "0f766e"
toc: true
toc-depth: 2
abstract: |
  Este livro ensina a apontar o Antigravity CLI (agy) para um gateway de LLM por variáveis de ambiente. O método troca o destino das chamadas de modelo e não toca na conta Google. Ele substitui o rodízio de contas, que causou o bloqueio de duas contas em 10/11/2026.

  A Parte I explica o problema e corrige três premissas comuns. A Parte II apresenta os três métodos de configuração: sistema operacional, linha de comando e pedido a um agente. A Parte III detalha as variáveis, os protocolos e a regra de precedência. A Parte IV traz receitas para os casos frequentes. A Parte V descreve as proibições que protegem as contas. A Parte VI reúne glossário, relatório de conformidade e fontes.

  Fonte primária: docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md, complementada pela documentação do Google Cloud sobre o gateway corporativo do Antigravity.
---

# Como ler este livro

## A quem este livro se destina

1. **Quem usa o `agy` no dia a dia.** Precisa ligar, trocar e desligar o gateway sem erro. As Partes II e IV são o núcleo.
2. **Quem lidera times que dividem contas.** Precisa entender por que o rodízio de contas é perigoso e o que o gateway muda. As Partes I e V respondem a isso.
3. **Quem mantém a plataforma LegoLoop.** Precisa saber quais comandos do `gw` foram bloqueados e por quê. A Parte V e o Apêndice B foram escritos para esse leitor.

## A estrutura da obra

| Parte            | Nível            | Conteúdo                                        |
| :--------------- | :--------------- | :---------------------------------------------- |
| I                | Visão geral      | O problema do rodízio e as três premissas corrigidas |
| II               | Método           | Os três métodos de configuração                 |
| III              | Detalhe          | Variáveis, protocolos e precedência             |
| IV               | Prática          | Receitas para os casos frequentes               |
| V                | Método           | Limites e proibições                            |
| VI               | Apêndices        | Glossário, conformidade e fontes                |


# PARTE I — O problema e as premissas

## 1. O que aconteceu

Em 10/11/2026 o `agy` mostrou a mensagem "service disabled ... violation of Terms of Service" na conta DEV_05.
A causa provável foi o `gw agy exec`. Ele rodava um rodízio de 13 contas Google e trocava de conta sozinho ao receber o erro 429.
Duas contas ficaram bloqueadas. Onze outras correm o mesmo risco enquanto o rodízio existir.

A empresa tem 13 desenvolvedores em uma máquina compartilhada por turnos. A cota de uma conta não basta, e a tentação de somar cotas é grande.
Somar cotas por rodízio de contas é o comportamento que o Google trata como violação.

## 2. A alternativa: trocar o destino, não a conta

O `agy` tem um recurso documentado pelo Google: o gateway de LLM.
Com ele, o `agy` envia as chamadas de modelo para uma URL que você escolhe, usando uma chave que o gateway exige.
A conta Google não participa. Nada troca de conta, e nenhum token Google é movido.


![Diagrama de Arquitetura](/img/diagrama_cb918c4c.png){width=85%}


## 3. Três premissas corrigidas

| Premissa                              | Veredito       | Por quê                                   |
| :------------------------------------ | :------------- | :---------------------------------------- |
| As variáveis são documentadas pelo Google | Confirmado     | Página "Connect Antigravity to an Enterprise Gateway" |
| Qualquer chave de qualquer provedor funciona | Parcial        | O destino precisa falar `genai` ou `openai`   |
| O uso é homologado e sem risco        | **Não provado**    | A página não cita termo de serviço nem gateway de terceiros |

O livro segue uma regra: **declare o que não foi provado**.
O que está provado é que o gateway não mexe na conta Google e que o rodízio causou o bloqueio.
O que não está provado é a posição do Google sobre gateways de terceiros. O leitor assume essa dúvida ao usar.

## 4. Gateway e login Google não convivem

A documentação do Google afirma que não se usa gateway e login com conta Google ao mesmo tempo.
Com o gateway ativo, as chamadas não consomem a cota de assinatura da conta.
Isso reforça a lógica da alternativa: o gateway substitui a conta, não a soma.


# PARTE II — Os três métodos de configuração

Os três métodos chegam ao mesmo resultado: as variáveis `AGY_LLM_GATEWAY_*` definidas no ambiente do usuário.


![Diagrama de Arquitetura](/img/diagrama_752f0502.png){width=85%}


## 5. Método A: pelo sistema operacional

No Windows:

1. Pressione `Win + R`.
2. Digite `rundll32 sysdm.cpl,EditEnvironmentVariables` e confirme.
3. Em **Variáveis de usuário**, clique em **Novo**.
4. Crie as três variáveis: `AGY_LLM_GATEWAY_URL`, `AGY_LLM_GATEWAY_WIRE_PROTOCOL` e `AGY_LLM_GATEWAY_API_KEY`.
5. Confirme as duas janelas e abra um terminal novo.

Para trocar a chave, edite o valor de `AGY_LLM_GATEWAY_API_KEY`. Para desligar, apague as três variáveis.

No macOS e no Linux, acrescente três linhas `export` ao `~/.zshrc` ou ao `~/.bashrc` e rode `source`.

## 6. Método B: pela linha de comando

Há dois modos.

**Modo de sessão.** As variáveis valem só no terminal atual. Serve para teste.

```powershell
$env:AGY_LLM_GATEWAY_URL = 'https://SEU-GATEWAY/v1'
$env:AGY_LLM_GATEWAY_WIRE_PROTOCOL = 'openai'
$env:AGY_LLM_GATEWAY_API_KEY = $env:MINHA_CHAVE
agy
```

**Modo persistente.** Grava no ambiente do usuário.

```powershell
[Environment]::SetEnvironmentVariable('AGY_LLM_GATEWAY_URL','https://SEU-GATEWAY/v1','User')
```

Nesta máquina existe o atalho do `gw`: `gw agy 9router` grava as três variáveis e `gw agy off` apaga as três.
O `gw` lê a chave da variável `NINEROUTER_KEY`.

## 7. Método C: por pedido a um agente

Você descreve o que quer, e o agente altera as variáveis.
A regra principal: **o valor da chave nunca vai no texto do pedido**. O pedido cita só o nome de uma variável que já guarda a chave.

Exemplo de pedido:

```text
Configure o agy para usar o gateway em https://SEU-GATEWAY/v1 com protocolo openai.
Leia a chave da variável MINHA_CHAVE e grave as três variáveis AGY_LLM_GATEWAY_*
no ambiente do usuário. Não imprima a chave. Confirme só com "definida: True/False".
```

Por que a regra existe: tudo que o agente imprime vira histórico da conversa.
Uma chave impressa fica exposta. Em 13/11/2026 a chave do 9Router vazou desse jeito, na saída de um comando.
O hook `leak-scan` do harness detecta o caso e avisa, mas só depois do vazamento.

## 8. Qual método escolher

| Situação                                         | Método                          |
| :----------------------------------------------- | :------------------------------ |
| Uso eventual, uma pessoa, sem automação          | A                               |
| Teste rápido ou script                           | B (modo de sessão)              |
| Máquina da equipe com `gw`                         | B (modo persistente)            |
| Várias máquinas ou configuração repetida         | C, com a chave em variável      |


# PARTE III — Variáveis, protocolos e precedência

## 9. As seis variáveis

| Variável                      | Obrigatória     | Função                               |
| :---------------------------- | :-------------: | :----------------------------------- |
| `AGY_LLM_GATEWAY_URL`              | Sim             | URL HTTPS do gateway                 |
| `AGY_LLM_GATEWAY_API_KEY`           | Sim             | Chave que o gateway exige            |
| `AGY_LLM_GATEWAY_WIRE_PROTOCOL`     | Não             | `genai` (padrão) ou `openai`             |
| `AGY_LLM_GATEWAY_CA_CERT`           | Não             | Caminho do pacote de certificados raiz |
| `AGY_LLM_GATEWAY_MODELS`           | Não             | Modelos permitidos, separados por vírgula |
| `AGY_LLM_GATEWAY_HEADERS`          | Não             | Cabeçalhos extras no formato `CHAVE: VALOR` |

A chave segue como token no cabeçalho `Authorization` de cada chamada.

## 10. Os dois protocolos

| Protocolo      | Quando usar                           | Observação                   |
| :------------- | :------------------------------------ | :--------------------------- |
| `genai`          | Gateway que fala a API do Gemini      | Padrão do `agy`                |
| `openai`         | Gateway ou provedor com API Chat Completions | Não funciona no servidor ACP |

O servidor ACP é o que o JetBrains e o Zed usam. Nesses editores, só `genai` funciona.

## 11. Qual chave serve

A chave é a do gateway escolhido. Não há "chave Google" obrigatória.

| Destino                        | Protocolo      | URL                                    | Chave             |
| :----------------------------- | :------------: | :------------------------------------- | :---------------- |
| 9Router da empresa             | `openai`         | `https://roteador.signa-docs.com/v1`     | chave do 9Router  |
| Provedor compatível com OpenAI | `openai`         | URL base do provedor                   | chave do provedor |
| Gateway que fala Gemini        | `genai`          | URL do gateway                         | chave do gateway  |
| Roteador próprio               | `openai`         | URL do roteador                        | chave do roteador |

Um provedor que não fala `genai` nem `openai` não funciona direto. Coloque um roteador na frente dele.

## 12. Precedência

1. O `admin_settings.json` do administrador vence a variável de ambiente. Em máquina gerenciada, a variável pode não valer.
2. No JetBrains e no Zed, defina `"auth": { "type": "gateway" }` em `~/.gemini/antigravity-acp/settings.json` e ponha as variáveis no campo `env` do agente.
3. Programas abertos pelo menu Iniciar ou pelo Dock não leem o perfil do shell.
4. O `agy` na web não aceita gateway.

## 13. Terminal novo, sempre

Variável de usuário só chega aos processos que nascem depois dela.
Todo terminal aberto antes da mudança mantém o ambiente antigo.
Esse é o motivo mais comum de "mudei e não funcionou".


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


# PARTE V — Limites e proibições

## 19. O que o gateway não resolve

O gateway troca o destino das chamadas. Ele não dá permissão do Google para nada.
Duas dúvidas continuam abertas:

- A posição do Google sobre gateways de terceiros. A página oficial não a descreve.
- A política de uso do provedor que você escolher. Cada provedor tem as suas regras.

## 20. As proibições do `gw`

Depois do bloqueio de 10/11/2026, o `gw` recusa os comandos que trocam de conta ou consultam várias contas.

| Comando                        | Estado           | Motivo                           |
| :----------------------------- | :--------------: | :------------------------------- |
| `gw agy next`                    | Bloqueado        | troca de conta por código        |
| `gw agy failover`                | Bloqueado        | troca de conta por código        |
| `gw agy devNN`                   | Bloqueado        | seleciona outra conta            |
| `gw agy sync` e `wincred`          | Bloqueado        | injeta tokens de outra conta     |
| `gw agy probe --all`             | Bloqueado        | usa o token de 13 contas         |
| `gw agy refresh all`             | Bloqueado        | renova tokens de 13 contas       |
| `gw failover agy`                | Bloqueado        | troca de conta por código        |
| `gw agy probe` (conta ativa)     | Liberado         | uma conta, a logada              |
| `gw agy list`                    | Liberado         | local, sem rede                  |

O `gw agy exec` não faz mais failover automático em erro 429 (`maxFailoverRetries = 0`).
O wrapper `agy.cmd` chama o `agy.exe` direto, e o `gw harness sync-wrappers` não gera wrapper do `agy`.

## 21. Por que `probe --all` não é só uma consulta

A leitura comum diz que `probe` só consulta. O código mostra outra coisa.
Para cada conta, o `probe` renova o token e chama a API interna do Google que informa a cota, com o identificador do Antigravity.
Com `--all`, o processo repete isso para 13 contas, na mesma máquina, em sequência.
O padrão se parece com o rodízio. A relação com os bloqueios é hipótese, sem prova. O risco justifica o bloqueio.

## 22. Segredos

1. Nunca grave chave em código, em log, em commit ou em resposta de chat.
2. Chave que apareceu em saída de comando está comprometida. Revogue e troque.
3. Verifique variáveis só com `definida: True/False`.
4. Mantenha o hook `leak-scan` ativo no harness.


# PARTE VI — Apêndices

## Apêndice A — Glossário

| Termo                        | Significado                                         |
| :--------------------------- | :-------------------------------------------------- |
| `agy`                          | Antigravity CLI, o harness de linha de comando do Google |
| Gateway                      | Endereço HTTPS que recebe as chamadas de modelo no lugar da conta Google |
| Protocolo (`WIRE_PROTOCOL`)     | Formato das chamadas: `genai` ou `openai`               |
| Combo                        | Conjunto nomeado de modelos no 9Router              |
| Rodízio                      | Troca de conta Google por código para somar cotas. Proibido |
| Terminal novo                | Terminal aberto depois da mudança das variáveis     |
| `leak-scan`                    | Hook que bloqueia saída com segredo                 |

## Apêndice B — Relatório de conformidade

| Fato                                   | Prova                                       |
| :------------------------------------- | :------------------------------------------ |
| As seis variáveis `AGY_LLM_GATEWAY_*` são documentadas | Google Cloud, "Connect Antigravity to an Enterprise Gateway" |
| `gw agy 9router` grava as três variáveis | `pkg/driver/driver.go`, `SetRoute`              |
| `gw agy off` apaga as três variáveis     | `pkg/driver/driver.go`, ramo `RouteOff`         |
| Gateway do 9Router respondeu           | `/v1/models` com HTTP 200 em 13/11/2026       |
| O `agy` respondeu pelo 9Router           | confirmação do usuário em terminal novo     |
| Comandos de rodízio estão bloqueados   | execução de `next`, `dev03`, `probe --all` e `sync` com código de saída 2 |
| O Google aprova gateway de terceiros   | **Não verificado**                              |
| `probe --all` causou os bloqueios        | **Não verificado** (hipótese)                   |
| `gw agy list` não usa rede               | leitura de `pkg/driver/list.go`, sem chamadas HTTP |

## Apêndice C — Referência rápida

| Quero                             | Faço                                           |
| :-------------------------------- | :--------------------------------------------- |
| Ligar o 9Router                   | `gw agy 9router`, terminal novo                  |
| Voltar ao `agy` puro                | `gw agy off`, terminal novo                      |
| Testar sem persistir              | variáveis `$env:` na sessão                      |
| Conferir sem mostrar a chave      | `definida: True/False`                           |
| Trocar a chave                    | janela do sistema ou pedido ao agente com variável |

## Apêndice D — Fontes

- Google Cloud, "Connect Antigravity to an Enterprise Gateway": https://docs.cloud.google.com/gemini/enterprise/docs/ai-developer-tools/connect-enterprise-gateway
- Google Cloud, "Configure an enterprise gateway for Google Antigravity": https://docs.cloud.google.com/gemini/enterprise/docs/ai-developer-tools/deploy-enterprise-gateway
- Portkey, "Google Antigravity" (guia de terceiros): https://portkey.ai/docs/integrations/libraries/antigravity
- Manual operacional: `docs/manuals/2026-11-13_MANUAL_GATEWAY_LLM_AGY_LEGOLOOP.md`


