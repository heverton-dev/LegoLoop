# PARTE I — Macro: O Fim da Geração Cega

## 1. O problema: código que chega antes da decisão

Um agente de IA consegue escrever centenas de linhas em segundos. Esse é o problema. Quando o pedido é "adicione um sistema de assinaturas", o agente não para para perguntar se uma assinatura pertence a um usuário ou a uma conta, se o cancelamento vale no fim do ciclo ou imediatamente, ou o que acontece com um pagamento recusado. Ele escolhe uma resposta plausível para cada pergunta, sem avisar, e escreve o código.

O resultado parece pronto. Os testes passam, porque foram escritos pelo mesmo agente, a partir das mesmas suposições. Semanas depois, alguém descobre que "conta" significa coisas diferentes em três módulos e que a regra de cancelamento contradiz o contrato comercial.

O método de Matt Pocock parte de um diagnóstico simples: o gargalo do desenvolvimento com agentes não é escrever código, é decidir o que escrever. Por isso o método gasta esforço deliberado antes da primeira linha de produção e mantém o humano no controle das decisões que são difíceis de desfazer.

### Três sintomas da geração cega

* **Premissas silenciosas.** O agente preencheu lacunas do pedido com suposições que ninguém aprovou.
* **Testes que concordam com o código por construção.** O teste recalcula o resultado do mesmo jeito que a implementação e, por isso, nunca discorda dela.
* **A sessão que apodrece.** Depois de horas de conversa, o agente começa a esquecer decisões do início, a misturar tentativas descartadas com a versão atual e a errar coisas que acertava antes.

## 2. Os cinco pilares do método

### Pilar 1: alinhar antes de construir

Em inglês, *grill before build*. Nenhuma linha de código de produção é gerada antes de a ideia passar por uma entrevista implacável que exponha premissas ocultas, casos de borda e ambiguidades de domínio. A entrevista não é uma formalidade: ela termina com decisões registradas em arquivo.

### Pilar 2: o desenvolvedor no controle

O agente opera em execuções curtas, de três a cinco passos, com gates de parada obrigatória. Ações críticas, como push, abertura de pull request, deploy e deleção, exigem autorização humana expressa.

```{=typst}
#painel("Por que execuções curtas")[
  Um agente que executa vinte passos sem parar acumula erros que se apoiam uns nos outros. Em três a cinco passos, o desvio é pequeno o bastante para ser notado e barato o bastante para ser desfeito.
]
```

### Pilar 3: módulos profundos

O pilar vem de John Ousterhout e do livro *A Philosophy of Software Design*. Um módulo profundo tem uma interface estreita e simples que esconde uma implementação rica. O contrário, o módulo raso, expõe quase tanta complexidade na interface quanto contém por dentro.

O método usa a palavra *seam*, costura, para a fronteira pública em que o comportamento é observado sem olhar para dentro. É nas costuras que os testes moram.

### Pilar 4: higiene severa de contexto

Modelos de linguagem degradam o raciocínio quando a janela de contexto cresce além de uma zona de clareza que o manual estima em cerca de 150 mil tokens. O método chama a faixa saudável de *Smart Zone* e organiza o trabalho para ficar dentro dela.

A regra prática é dupla: alinhamento, especificação e tickets acontecem em uma janela única e contínua, porque compartilham raciocínio; a implementação de cada ticket começa em contexto limpo, porque só precisa do ticket e do repositório.

### Pilar 5: retroalimentar o ambiente

Retrospectivas não mudam o código passado. Elas mudam as regras, os linters, as checagens automáticas e o `AGENTS.md`, para que a próxima sessão não repita o mesmo erro.

```{=typst}
#frase-destaque[Um erro corrigido no código volta. Um erro corrigido no ambiente não volta.]
```

## 3. O que o leitor leva desta parte

O método não é um conjunto de prompts melhores. É uma divisão de trabalho: o humano decide o que é difícil de reverter, o agente executa fatias pequenas e verificáveis, e cada falha vira uma regra automática. A tabela abaixo antecipa onde cada pilar se materializa no loop.

| Pilar                             | Onde aparece no loop                           |
| :-------------------------------- | :--------------------------------------------- |
| Alinhar antes de construir        | Sub-loop 1 e gate G1                           |
| Desenvolvedor no controle         | Gates G1 a G4 e ações críticas autorizadas     |
| Módulos profundos                 | `codebase-design` e testes em costuras           |
| Higiene de contexto               | Janela única até os tickets, `/clear` depois     |
| Retroalimentar o ambiente         | `retro` e o gate G4                              |
