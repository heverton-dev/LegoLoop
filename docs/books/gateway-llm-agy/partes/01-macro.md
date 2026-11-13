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
