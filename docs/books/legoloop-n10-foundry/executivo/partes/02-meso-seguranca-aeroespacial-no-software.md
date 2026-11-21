# Parte II: Segurança Aeroespacial no Software

## Capítulo 2: Como Protocolos Industriais Protegem o Negócio

A confiabilidade do LegoLoop se apoia em metodologias comprovadas que operam onde falhas são intoleráveis.

### 1. O Princípio Toyota Jidoka (Andon Cord)
Na indústria tradicional, qualquer operador puxa uma corda (Andon) para paralisar a linha de montagem ao identificar um defeito.  
No LegoLoop, o sistema possui um limite estrito de **duas tentativas**. Se o agente de IA não resolve o problema na segunda iteração, o sistema puxa o Andon Cord digital, suspende a esteira e emite um relatório claro para intervenção humana. Isso impede custos descontrolados de API e código corrompido.

### 2. A Trava ICAO da Aviação Civil
Antes de qualquer alteração destrutiva (como remoção de arquivos ou reversão de commits), o operador recebe um **Cartão de Confirmação ICAO**.  
A operação só prossegue mediante autorização explícita com o comando `HEARBACK APPROVED`.
