# PARTE V — Prática: Receitas de Uso

As receitas combinam sub-loops para situações frequentes. Os nomes de funcionalidades e os números são ilustrativos.

## 17. Receita: uma funcionalidade nova, do zero

**Quando usar.** Há uma ideia, por exemplo "permitir que clientes pausem a assinatura", e nenhuma linha de código.

1. Rode `grill-with-docs`. Responda as rodadas até que "pausar" tenha um significado único: quanto tempo, quem paga, o que acontece com o ciclo.
2. Confirme que as decisões estão em `GLOSSARY.md` e, se alguma for difícil de reverter, em uma ADR. Este é o G1.
3. Na mesma sessão, rode `to-spec` e depois `to-tickets`.
4. Revise o DAG: cada ticket precisa dizer de quem depende. Este é o G2.
5. Limpe o contexto. Para cada ticket desbloqueado, abra uma sessão limpa e rode `implement` com `tdd`.
6. Com a suíte verde e sem testes enfraquecidos (G3), rode `code-review`, `pr` e `retro`.

**Erro comum.** Pular o passo 5 e implementar os tickets na mesma sessão do alinhamento. A sessão já carrega toda a conversa, e o agente passa a implementar a partir dela, e não do ticket.

## 18. Receita: um bug difícil

**Quando usar.** "O total do carrinho às vezes sai errado." Ninguém sabe reproduzir.

1. Rode `diagnosing-bugs` e não aceite hipóteses antes do passo 2.
2. Construa um loop de feedback: um comando que fique vermelho por causa deste bug. Tente, nesta ordem, teste, script HTTP, CLI com snapshot, navegador headless e reprodução de trace.
3. Com o vermelho em mãos, investigue.
4. Corrija com o menor patch possível e veja o mesmo comando ficar verde.
5. Transforme o comando em teste de regressão permanente.
6. Rode `retro`: se o bug veio de uma classe de erro detectável, converta-a em checagem automática.

```{=typst}
#painel("Sem loop, sem diagnóstico")[
  Se nenhuma das formas de construir o loop funcionar, o passo seguinte é dizer isso ao usuário e pedir mais informação, não teorizar. Teorizar sem um sinal vermelho produz correções que parecem plausíveis e não se pode provar que funcionaram.
]
```

## 19. Receita: uma dúvida que só código responde

**Quando usar.** No meio do alinhamento surge a pergunta "uma fila em memória aguenta o pico de pedidos?".

1. Interrompa o grilling e rode `handoff` para uma sessão isolada.
2. Na sessão nova, rode `prototype` na branch `prototype/fila-pico`.
3. Meça, anote o resultado e descarte o código.
4. Volte com `handoff`, levando só a resposta e a evidência.
5. Retome o grilling com a decisão informada.

## 20. Receita: revisar o trabalho de outro agente

**Quando usar.** Um agente entregou uma branch dizendo que está pronta.

1. Fixe o ponto de comparação, por exemplo `main`, e confirme que ele existe e que o diff não está vazio.
2. Localize a spec de origem: referência de issue nos commits, caminho passado como argumento ou arquivo em `docs/specs/`.
3. Rode `code-review`. Os dois eixos, Standards e Spec, rodam em subagentes separados.
4. Trate cada nota bloqueadora antes do PR. Notas não bloqueadoras viram tickets.
5. Rode `retro` sobre o que o review encontrou que uma checagem automática poderia ter pego.

O passo 1 existe porque um ponto de comparação inválido deve falhar logo, e não dentro de dois subagentes paralelos.
