# 3. Etapas

Cada etapa lê entradas, faz um trabalho e gera saídas.
Ao fim, há um portão.
Sem aprovação, o agente não avança.
Fonte: `processo/PROCESSO.md`, seção "Portões de aprovação".

| Etapa          | Nome              | Entrada             | Trabalho                               | Saída                      | Portão         |
| :------------- | :---------------- | :------------------ | :------------------------------------- | :------------------------- | :------------- |
| 0              | Configuração      | respostas da pessoa | pergunta uma a uma                     | `projeto.json`               | sim            |
| 1              | Pauta             | nicho e fontes      | busca tendências e demanda             | 10 ideias, 3 marcadas      | escolha        |
| 2              | Roteiro           | pauta aprovada      | pesquisa em fontes A ou B e escreve em batidas | `roteiro.json`, batidas      | sim            |
| 3              | Elenco            | roteiro             | fichas do personagem no Flow           | personagem salvo           | sim            |
| 4              | Estilo e cenários | elenco              | prompt mestre, cenários, objetos       | prompts e quadros de teste | sim            |
| 5              | Storyboard        | cenários e elenco   | um quadro por cena                     | `storyboard.json`            | sim            |
| 6              | Voz               | roteiro             | narração por linha com tempos          | `audio/timings.json`         | checagem       |
| 7              | Animação          | storyboard          | motor em código ou vídeo do Flow       | cenas                      | cena teste     |
| 8              | Montagem          | cenas, voz, marca   | timeline com marca e legendas          | `timeline.json`              | não            |
| 8b             | Cortes curtos     | timeline            | recorta trechos em 9:16                | `renders/shorts/`            | sim            |
| 9              | Revisão           | render              | folhas de contato e revisor            | "enviar" ou "não enviar"   | sim            |
| 10             | Vitrine           | vídeo aprovado      | títulos, capas, descrição              | pacote                     | sim            |
| 11             | Publicação        | acervo              | agenda no Postiz                       | postagem                   | autorizada     |
| 12             | Comentários       | comentários         | rascunhos de resposta                  | respostas                  | sim            |
| 13             | Aprendizado       | números             | recomendações semanais                 | `acervo/aprendizado.md`      | não            |

## 3.1 Etapa 0: a régua da duração

A duração define quantos cortes e cenas o vídeo precisa.
A regra dos 3 segundos vale para o vídeo longo: um corte ou movimento a cada 3 s.
Em 9:16 e 1:1, a regra vira 2 a 2,5 s.
Fonte: `processo/PROCESSO.md`, tabela de ritmo.

## 3.2 Etapa 2: regras do roteiro

Abra com o momento mais dramático.
Ligue as batidas com "mas" ou "portanto", nunca com "e então".
Use linhas de 8 a 25 palavras.
Separe fato de alegação com "segundo os promotores".
Fonte: `processo/PROCESSO.md`, etapa 2.

## 3.3 Etapa 3: o elenco vem primeiro

Cada personagem tem três fichas no Flow: referência, expressões e poses.
A descrição é igual em todos os prompts do mesmo personagem.
O personagem aprovado vai para o recurso "Personagens" do Flow.
Fonte: `processo/PROCESSO.md`, etapa 3.

## 3.4 Etapa 7: três formas de animar

A: motor em código, grátis e local.
B: vídeo do Flow por cena, que consome créditos.
C: híbrido, com vídeo do Flow só nos momentos-chave.
Antes de tudo, faça uma cena teste de 10 s.
Fonte: `processo/PROCESSO.md`, etapa 7.
