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
