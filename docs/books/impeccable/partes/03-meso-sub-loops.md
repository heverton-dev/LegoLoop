# PARTE III — Meso: Os Seis Sub-Loops

## 7. Visão geral da taxonomia

Os cerca de 25 comandos do Impeccable estão agrupados em seis sub-loops reutilizáveis. Cada sub-loop tem um objetivo único, o que permite chamá-lo isoladamente quando a tarefa não exige o fluxo completo.

| Sub-loop        | Nome                     | Objetivo em uma frase                    |
| :-------------- | :----------------------- | :--------------------------------------- |
| 1               | Avaliar                  | Radiografar o estado atual sem alterar nada |
| 2               | Construir e Direcionar   | Fixar contratos de produto antes de codificar |
| 3               | Melhorar                 | Elevar a qualidade pelos fundamentos visuais |
| 4               | Refinar                  | Preparar a interface para produção       |
| 5               | Corrigir                 | Resolver fricções específicas            |
| 6               | Iterar                   | Explorar alternativas visuais ao vivo    |

## 8. Sub-loop 1: Avaliar

O objetivo é diagnosticar sem fazer alterações prematuras. Um agente que corrige enquanto avalia perde a visão do conjunto e tende a atacar o primeiro problema que encontra, não o mais grave.

* `critique` faz a revisão heurística: usabilidade, hierarquia de informação e carga cognitiva. Responde à pergunta "o usuário entende o que fazer?".
* `audit` faz a auditoria técnica: acessibilidade segundo a WCAG, performance de renderização e responsividade em vários viewports. Responde à pergunta "a tela funciona para todos?".

Os dois comandos são complementares. Uma tela pode passar no `audit` com contraste correto e foco visível e ainda assim falhar no `critique` por esconder a ação principal.

## 9. Sub-loop 2: Construir e Direcionar

Este sub-loop estabelece as fundações conceituais e os contratos duráveis de produto antes de qualquer código.

* `init` captura o contexto durável da aplicação em `PRODUCT.md`: para quem é, o que resolve, em que tom fala.
* `shape` planeja fluxos de tela e interações sem escrever código.
* `document` gera o `DESIGN.md` a partir da verdade do código.
* `extract` garimpa tokens e componentes consolidados para o design system.

```{=typst}
#painel("Dois documentos, dois papéis")[
  PRODUCT.md responde "o quê e para quem". DESIGN.md responde "como está construído". O primeiro é escrito antes do código, a partir de decisões. O segundo é escrito depois, a partir do código entregue. Inverter essa ordem produz um DESIGN.md cheio de intenções que o código não cumpre.
]
```

## 10. Sub-loop 3: Melhorar

O objetivo é elevar a qualidade sensorial e tátil da interface pelos fundamentos visuais. São seis comandos, e a ordem em que aparecem é a ordem recomendada.

| Comando        | Atua sobre                           | Pergunta que responde            |
| :------------- | :----------------------------------- | :------------------------------- |
| `layout`         | Respiros, grids, alinhamento, ritmo  | O olho encontra o caminho?       |
| `typeset`        | Proporção, peso, escala tipográfica  | A hierarquia é óbvia sem cor?    |
| `colorize`       | Paleta com significado semântico     | Cada cor quer dizer algo?        |
| `animate`        | Micro-interações funcionais          | O movimento explica uma mudança? |
| `delight`        | Detalhes memoráveis e estados táteis | Há um momento que se lembra?     |
| `overdrive`      | Expressividade estética ampliada     | Só em Persuade ou Experience     |

O comando `colorize` evita dois extremos que o manual nomeia: o monocromatismo insosso e o carnaval de cores. A cor só entra quando carrega significado: estado, categoria, ação primária.

O comando `overdrive` merece atenção especial. Ele é reservado exclusivamente aos modos Persuade e Experience. Usá-lo em uma tela Operate é um erro de classificação, não uma escolha estética.

## 11. Sub-loop 4: Refinar

Este sub-loop blinda a interface para ambientes de missão crítica.

* `polish` é o passe de micropolimento e consistência antes da entrega.
* `harden` trata estados de erro, carregamento, internacionalização e limites de dados: o nome com 80 caracteres, a lista vazia, a lista com dez mil itens, a tradução alemã que dobra o tamanho do botão.
* `onboard` cria os fluxos de primeira execução e os estados vazios.
* `distill` remove complexidade visual desnecessária e elementos ruidosos.
* `bolder` e `quieter` calibram o tom visual: o primeiro aumenta contraste, o segundo reduz estimulação excessiva.

```{=typst}
#painel("bolder e quieter não são opostos de gosto")[
  Os dois comandos corrigem um desvio em relação ao modo da superfície. Uma tela Persuade sem energia pede bolder. Uma tela Operate que grita pede quieter. A decisão vem do modo, não da preferência de quem revisa.
]
```

## 12. Sub-loop 5: Corrigir

O objetivo é resolver fricções ergonômicas e falhas operacionais pontuais.

* `clarify` refatora a UX copy: microtexto, rótulos de botões e mensagens de validação. O piso de qualidade dá o critério: controles nomeiam sua ação; erros nomeiam o problema e a forma de recuperação.
* `adapt` adequa a tela a extremos: mobile compacto, ultra-wide e layouts nativos.
* `optimize` diagnostica e corrige gargalos de renderização na árvore de componentes.

Um botão que diz "Enviar" está correto, mas um botão que diz "Publicar artigo" está melhor, porque nomeia a consequência. Uma mensagem "Erro inválido" falha nos dois critérios: não diz qual é o problema nem como sair dele.

## 13. Sub-loop 6: Iterar

O último sub-loop serve para explorar alternativas visuais em tempo real.

* `live` permite inspecionar e selecionar elementos diretamente no navegador para gerar variantes ao vivo.
* `generate` cria variações simultâneas de um elemento no navegador para uma decisão imediata.

Iterar é útil quando a direção está aprovada, mas um elemento específico ainda não convence. Em vez de três rodadas de descrição verbal, o agente mostra três versões lado a lado.

## 14. Como escolher o sub-loop certo

| Situação                               | Sub-loop           | Comando de entrada      |
| :------------------------------------- | :----------------- | :---------------------- |
| Não sei o que está errado              | Avaliar            | `critique`                |
| Tela nova, sem código ainda            | Construir          | `shape`                   |
| Funciona, mas parece amadora           | Melhorar           | `layout`                  |
| Vai para produção amanhã               | Refinar            | `harden`                  |
| Usuários erram um formulário           | Corrigir           | `clarify`                 |
| Duas opções e nenhuma decisão          | Iterar             | `generate`                |
