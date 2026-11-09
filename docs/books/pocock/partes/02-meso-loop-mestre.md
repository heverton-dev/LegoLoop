# PARTE II — Meso: O Loop Mestre de Rigor

## 4. Da ideia à produção em cinco estágios

O loop mestre atravessa cinco estágios. Quatro gates binários separam os estágios principais: G1 a G4. O segundo estágio, a prova de conceito, é opcional e funciona como um desvio que volta ao primeiro.

```{=typst}
#esteira((
  ("1. Alinhamento e Domínio", "grill-with-docs, domain-modeling"),
  ("3. Especificação e Grafo", "to-spec, to-tickets"),
  ("4. Construção TDD", "implement, tdd, codebase-design"),
  ("5. Fechamento e Entrega", "code-review, pr, retro"),
), gates: ("G1: decisões aprovadas", "G2: tickets e bloqueios aprovados", "G3: testes 100% verdes"))
```

Depois do estágio 5, o gate G4, review aprovado, marca a entrega como concluída.

### O desvio empírico: estágio 2

Às vezes uma pergunta do alinhamento não se resolve conversando. "Esse modelo de estado aguenta edição concorrente?" só se responde com código rodando. Nesse caso o fluxo sai do estágio 1, abre uma sessão isolada com `handoff`, constrói um protótipo descartável com `prototype` e volta com a resposta por outro `handoff`. O protótipo não entra na árvore principal.

```{=typst}
#painel("Protótipo não é rascunho de produção")[
  O protótipo responde uma pergunta e é jogado fora. Se o código do protótipo for aproveitado na implementação, ele entra sem spec, sem tickets e sem testes primeiro, que é exatamente o que o método existe para evitar.
]
```

## 5. Os gates e o que cada um exige

| Gate           | Entre                       | Condição para passar                    |
| :------------- | :-------------------------- | :-------------------------------------- |
| G1             | Alinhamento e Especificação | Ambiguidades de domínio sanadas e aprovadas |
| G2             | Especificação e Construção  | DAG de tickets aprovado, bloqueios resolvidos |
| G3             | Construção e Fechamento     | Testes 100% verdes, sem testes enfraquecidos |
| G4             | Fechamento e Entrega        | Review sem notas bloqueadoras, ambiente atualizado |

### G1: decisões aprovadas

O estágio de alinhamento termina quando o desenvolvedor aprova as decisões. Não basta o agente achar que entendeu; as decisões precisam estar escritas no `GLOSSARY.md`, em um documento de intenção ou em ADRs.

### G2: tickets e bloqueios aprovados

A especificação vira um grafo acíclico de tickets. O gate exige que cada aresta de bloqueio esteja explícita: o ticket 4 depende do 2, o 5 não depende de ninguém. Um ticket com dependência implícita quebra a execução em contexto limpo, porque a sessão que o implementa não sabe o que veio antes.

### G3: testes 100% verdes

O manual detalha o que conta como verde: nenhum teste enfraquecido e nenhuma supressão de compilador. Um teste que foi afrouxado para passar, ou um `@ts-ignore` colocado para calar um erro de tipo, reprovam o gate mesmo com a suíte verde.

### G4: review aprovado

O review em dois eixos não pode ter notas bloqueadoras, e a retrospectiva precisa ter atualizado o ambiente. Uma entrega sem retro passa no código, mas reprova no método.

## 6. Por que a ordem importa

Cada estágio produz o insumo exato do seguinte. A especificação depende de um vocabulário fixado; os tickets dependem da especificação; a implementação depende de tickets autocontidos; o review depende da especificação para o segundo eixo.

| Estágio               | Produz                                | Consumido por         |
| :-------------------- | :------------------------------------ | :-------------------- |
| Alinhamento           | `GLOSSARY.md`, intenção, ADRs           | Especificação         |
| Especificação         | Spec e DAG de tickets                 | Construção            |
| Construção            | Código, testes, relatório TDD         | Fechamento            |
| Fechamento            | Review, corpo do PR, `AGENTS.md`        | A próxima sessão      |
