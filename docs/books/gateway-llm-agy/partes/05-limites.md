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
