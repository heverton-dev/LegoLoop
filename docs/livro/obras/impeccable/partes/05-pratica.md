# PARTE V — Prática: Receitas de Uso

As receitas a seguir combinam sub-loops para situações frequentes. São roteiros ilustrativos montados a partir do manual; os nomes de tela e os números servem de exemplo e não descrevem um projeto real.

## 18. Receita: auditar uma tela existente

**Quando usar.** Alguém diz que "o dashboard está estranho", mas ninguém sabe dizer o quê.

1. Classifique o modo. Um dashboard é quase sempre Operate.
2. Rode `critique` e `audit` em sequência, sem corrigir nada.
3. Consolide os achados em uma lista com severidade alta, média e baixa.
4. Apresente a lista e pare. Este é o gate P1.

**Saída esperada.** Uma lista como esta, ilustrativa:

| Severidade | Achado                                     | Origem     |
| :--------- | :----------------------------------------- | :--------- |
| Alta       | Ação principal abaixo da dobra em 375px    | `critique` |
| Alta       | Contraste de 3,1:1 no texto secundário     | `audit`    |
| Média      | Quatro tamanhos de título sem escala clara | `critique` |
| Baixa      | Scrollbar com o estilo padrão do navegador | `audit`    |

**Erro comum.** Corrigir o contraste enquanto ainda está avaliando. A correção isolada parece produtiva, mas impede de ver que o problema de hierarquia talvez peça uma reorganização maior.

## 19. Receita: polir antes de uma entrega

**Quando usar.** A tela está aprovada e funcionando; falta acabamento.

1. Leia `craft-floor.md`.
2. Rode `harden` para cobrir estados de erro, vazio e carregamento.
3. Rode `clarify` sobre todos os rótulos e mensagens.
4. Faça uma captura em lote, mobile e desktop.
5. Rode `polish` e corrija tudo em um único bloco.
6. Faça no máximo uma checagem final.

```{=typst}
#painel("O limite da receita")[
  Se a checagem final encontrar um problema estrutural, como a hierarquia que não funciona no mobile, a receita de polimento não é o lugar para resolvê-lo. Volte ao sub-loop Melhorar com layout e typeset, e só depois retome o polimento.
]
```

## 20. Receita: redesenhar uma landing page

**Quando usar.** A página atual é anti-referência: ninguém quer preservar nada dela além do conteúdo.

1. Rode `init` se não existir `PRODUCT.md`.
2. Use a página antiga apenas como evidência: liste o que ela faz mal.
3. Classifique o modo como Persuade.
4. Rode `shape` para planejar o fluxo sem código e peça aprovação da direção (P2).
5. Execute `layout`, `typeset` e `colorize`, nesta ordem.
6. Considere `overdrive`, que só é permitido aqui por causa do modo.
7. Verifique em lote (P4) e rode `document` para registrar o novo mundo visual.

## 21. Receita: decidir entre duas variantes

**Quando usar.** A direção está aprovada, mas o bloco de preços tem duas propostas e a equipe não decide.

1. Rode `generate` sobre o bloco para gerar as variações lado a lado no navegador.
2. Use `live` para ajustar um detalhe da variante preferida.
3. Registre a escolha e o motivo.
4. Volte ao estágio 4 com a variante escolhida.

## 22. Receita: corrigir um formulário em que usuários erram

**Quando usar.** O formulário de cadastro tem alta taxa de abandono e o suporte recebe perguntas sobre as mesmas mensagens de erro.

1. Rode `critique` focado no formulário para localizar onde o usuário hesita.
2. Rode `clarify` sobre rótulos, textos de ajuda e mensagens de validação.
3. Rode `harden` para cobrir os limites: nome muito longo, campo vazio, colagem de texto com espaços.
4. Faça a captura em lote, mobile e desktop, com os estados de erro visíveis.

O critério do `clarify` vem do piso de qualidade: cada erro precisa nomear o problema e a forma de recuperação. A tabela abaixo mostra o tipo de troca esperada; as mensagens são ilustrativas.

| Antes             | Depois                                              |
| :---------------- | :-------------------------------------------------- |
| Erro inválido     | O CPF precisa ter 11 números. Confira os dígitos.   |
| Campo obrigatório | Informe um e-mail para receber a confirmação.       |
| Senha fraca       | Use 8 caracteres ou mais, com pelo menos um número. |
| Enviar            | Criar conta                                         |

## 23. Receita: tornar uma tela mais calma

**Quando usar.** Usuários de um painel interno reclamam de cansaço visual.

1. Confirme o modo Operate.
2. Rode `distill` para remover elementos que não ajudam a tarefa.
3. Rode `quieter` para reduzir a estimulação.
4. Revise com `audit` para garantir que o contraste mínimo continua respeitado depois da redução.

O passo 4 existe porque reduzir estímulo visual costuma significar clarear cores, e cores claras demais derrubam o contraste abaixo de 4,5:1.
