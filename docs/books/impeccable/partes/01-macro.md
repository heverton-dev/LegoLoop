# PARTE I — Macro: Por que Interfaces de IA Parecem Iguais

## 1. O problema do visual médio

Peça a um modelo de linguagem uma "landing page moderna" e o resultado costuma ser previsível: um título enorme sobre um gradiente roxo, um rótulo pequeno em caixa alta acima do título, três cartões do mesmo tamanho com ícone, título e parágrafo, um número grande com legenda pequena e sombras suaves em tudo. Peça de novo, para outro produto, e a página volta quase igual.

Isso não é falta de capacidade do modelo. É o efeito estatístico de como ele foi treinado. Um modelo aprende a produzir o que é mais provável dado o pedido. Para um pedido vago, o mais provável é a média de milhares de páginas parecidas. O resultado é uma interface que não tem erro óbvio, mas também não pertence a nenhum produto em particular.

O manual do Impeccable chama a saída desejada de *out-of-distribution*: uma entrega que sai da média, porque foi decidida e não apenas gerada. O método inteiro existe para trocar o padrão estatístico por decisões explícitas.

### Três sintomas que o leitor reconhece

* **A tela montada.** Os componentes estão todos lá, mas nenhum foi pensado para o conteúdo real. Os rótulos dizem "Lorem" ou "Título", os estados de erro não existem e a tela quebra em 375 pixels de largura.
* **O enfeite no lugar da hierarquia.** Gradiente no texto, brilho colorido sem deslocamento, vidro desfocado sem função. Cada um desses recursos tenta compensar a falta de uma hierarquia tipográfica clara.
* **A revisão infinita.** O agente tira uma captura de tela, ajusta um detalhe, tira outra, ajusta outro. Sem um critério de parada, o ciclo consome tokens e tempo sem convergir.

```{=typst}
#painel("A diferença entre construir e montar")[
  O arquivo de referência craft-floor.md, que acompanha o upstream do Impeccable, resume a distinção: as superfícies que ninguém desenhou (seleção de texto, cursor, barras de rolagem, anéis de foco, numerais em tabelas) continuam carregando o design. Tematizá-las a partir da paleta é, segundo o próprio documento, o sinal mais barato de que uma página foi construída e não montada, e o que os modelos mais deixam de fazer.
]
```

## 2. Os quatro pilares do método

O manual organiza o Impeccable em quatro pilares. Eles não são etapas; são regras que valem em qualquer etapa.

### Pilar 1: o briefing vence

A estética, a paleta, a época e as diretrizes definidas por quem pediu a interface têm prioridade sobre qualquer gosto do agente. O manual é explícito: nunca redirecione uma interface para gosto pessoal arbitrário.

Na prática, isso significa que se o briefing pede uma estética anos 90 com bordas grossas, o agente não "moderniza" a proposta por conta própria. O piso de qualidade, como veremos na Parte IV, trata vários recursos como padrões evitáveis e não como proibições: as palavras do briefing podem autorizar qualquer um deles.

### Pilar 2: refinamento não é redesenho

Toda tarefa de interface começa com uma pergunta: estamos refinando ou redesenhando?

| Abordagem       | O que preserva                    | O que muda                       | Quando usar                           |
| :-------------- | :-------------------------------- | :------------------------------- | :------------------------------------ |
| Refinamento     | Identidade, comportamento e texto | Alinhamento, ritmo, consistência | A tela funciona, mas está mal acabada |
| Redesenho       | Apenas o conteúdo e a função      | Todo o mundo visual              | A tela atual é anti-referência        |

No redesenho, a interface antiga vira evidência do que não fazer. O novo mundo visual precisa ser documentado, e não apenas aplicado, para que a próxima sessão o reconheça.

### Pilar 3: o piso de qualidade é inegociável

Antes de qualquer edição de UI, inclusive ajustes pontuais, a leitura do arquivo `craft-floor.md` é obrigatória. O manual cita alguns exemplos do que esse piso proíbe: gradientes sem contraste, bordas excessivas, tipografia desordenada e emojis na interface. A Parte IV dedica um capítulo inteiro ao conteúdo desse arquivo.

### Pilar 4: verificação delimitada

O quarto pilar ataca diretamente o sintoma da revisão infinita. A regra tem três partes:

1. Inspecione em lote: desktop e mobile na mesma rodada.
2. Corrija todos os apontamentos em um único bloco de edição.
3. Finalize com, no máximo, mais uma checagem.

```{=typst}
#frase-destaque[Uma rodada de captura, um bloco de correções, uma checagem final. Não existe quarta volta.]
```

O nome em inglês, *bounded pass*, descreve bem a ideia: a verificação tem uma fronteira conhecida antes de começar. Isso torna o custo previsível e força o agente a olhar para todos os problemas de uma vez, em vez de corrigi-los um por um à medida que aparecem.

## 3. O que o leitor leva desta parte

O Impeccable não é uma biblioteca de componentes nem um tema visual. É um processo com quatro regras permanentes: respeitar o briefing, decidir entre refinar e redesenhar, cumprir o piso de qualidade e limitar a verificação. As próximas partes mostram como essas regras viram estágios, gates e comandos concretos.
