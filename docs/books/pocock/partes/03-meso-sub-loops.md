# PARTE III — Meso: Os Seis Sub-Loops

## 7. Visão geral

O manual decompõe o ecossistema de Matt Pocock em seis sub-loops, disponíveis como blocos conectáveis no Studio Canvas e em pipelines DAG.

| Sub-loop       | Nome                           | Gate de saída                            |
| :------------- | :----------------------------- | :--------------------------------------- |
| 1              | Concepção, Grilling e Domínio  | Ambiguidades sanadas e aprovadas         |
| 2              | Exploração Empírica e Pesquisa | Resposta empírica devolvida via `handoff`  |
| 3              | Especificação e Grafo          | DAG aprovado, bloqueios resolvidos       |
| 4              | Construção Cirúrgica e TDD     | Testes verdes, sem testes enfraquecidos  |
| 5              | Diagnóstico e Arquitetura      | Falha reproduzida, corrigida e protegida |
| 6              | Auditoria, Entrega e Ambiente  | Review aprovado e ambiente atualizado    |

## 8. Sub-loop 1: Concepção, Grilling e Domínio

O objetivo é fixar um vocabulário unívoco, desafiar premissas e registrar decisões arquiteturais antes de qualquer código.

* `grill-with-docs` é a entrevista socrática focada no repositório. Deixa rastro durável em `GLOSSARY.md` e em ADRs.
* `grill-me` é a variante sem estado, para alinhamento rápido sem pasta de projeto.
* `grilling` é o primitivo da entrevista: rodadas numeradas e separação clara entre o que o agente sabe e o que o usuário decide.
* `domain-modeling` audita termos sobrecarregados, como "usuário" e "conta", e documenta decisões difíceis de reverter.

**Saídas:** `GLOSSARY.md`, `docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md` e ADRs em `docs/adr/`.

```{=typst}
#painel("Fatos do agente e decisões do usuário")[
  A separação que o grilling impõe é a parte mais valiosa da entrevista. O agente pode afirmar fatos que leu no repositório: o campo tal existe, a tabela tal tem tal índice. Ele não pode afirmar decisões: se o cancelamento é imediato é uma escolha de produto, e só o usuário a faz.
]
```

## 9. Sub-loop 2: Exploração Empírica e Pesquisa

O objetivo é responder perguntas que só código executável ou fontes primárias resolvem, sem poluir a árvore principal.

* `prototype` cria um programa descartável na branch `prototype/<nome>`, para testar um modelo de estado ou a sensação de uma interface.
* `handoff` passa o contexto de forma estéril para uma nova pasta, harness ou sessão.
* `research` roda uma pesquisa profunda em segundo plano sobre documentação oficial e fontes primárias.

**Saídas:** relatório técnico com citações e protótipo referenciado na branch isolada.

## 10. Sub-loop 3: Especificação e Grafo

O objetivo é sintetizar o consenso em uma especificação formal e fatiá-la em um grafo acíclico de tarefas independentes.

* `to-spec` consolida a discussão em uma spec executável e determinística, sem entrevistar de novo.
* `to-tickets` quebra a spec em tickets verticais, os *tracer bullets*, cada um com suas arestas de bloqueio (`blocked-by`).
* `wayfinder` traça um mapa de tickets de decisão para iniciativas nebulosas ou projetos grandes começando do zero.

**Saídas:** `docs/specs/YYYY-MM-DD_SPEC_<TEMA>.md` e tickets em `.scratch/<feature>/issues/NN-slug.md` ou `docs/plans/`.

Um ticket vertical atravessa todas as camadas necessárias para entregar um comportamento observável, da interface ao banco. O contrário, o ticket horizontal ("criar todas as tabelas"), não entrega nada testável sozinho.

## 11. Sub-loop 4: Construção Cirúrgica e TDD

O objetivo é implementar orientado a testes, com contexto isolado por ticket.

* `implement` consome um ticket por vez, a partir de uma janela limpa, atacando nós desbloqueados.
* `implement-spec` distribui subagentes em paralelo sobre os tickets prontos e unifica em uma branch de integração.
* `tdd` conduz o ciclo vermelho-verde por fatia de comportamento, sem mocks de produção e sem stubs vazios.
* `codebase-design` fornece o vocabulário de módulos profundos, interfaces estreitas e costuras.

**Saídas:** código de produção, testes e `docs/reports/TDD_REPORT.md`.

### O que a habilidade tdd exige

O texto upstream da habilidade `tdd` é mais preciso do que o resumo do manual em três pontos.

1. **Testes só em costuras combinadas.** Antes de escrever qualquer teste, o agente lista as costuras que serão testadas e confirma com o usuário. Nenhum teste é escrito em costura não confirmada.
2. **Três anti-padrões nomeados.** O teste acoplado à implementação, que quebra em refatorações sem mudança de comportamento; o teste tautológico, que recalcula o esperado do mesmo jeito que o código; e o fatiamento horizontal, que escreve todos os testes antes de toda a implementação.
3. **Refatorar não faz parte do loop.** O ciclo é vermelho e verde. A refatoração pertence ao estágio de review, conduzido pela habilidade `code-review`.

```{=typst}
#painel("O teste tautológico", cor: cor-erro)[
  expect(soma(a, b)).toBe(a + b) passa por construção: o esperado é calculado do mesmo jeito que o código calcula. O valor esperado precisa vir de uma fonte independente: um literal sabidamente correto, um exemplo resolvido à mão ou a especificação.
]
```

## 12. Sub-loop 5: Diagnóstico, Arquitetura e Triagem

O objetivo é resolver falhas complexas e manter a saúde arquitetural.

* `diagnosing-bugs` proíbe teorizar antes de obter um loop de feedback reproduzível: um comando que fique vermelho por causa deste bug. A correção vem acompanhada de teste de regressão.
* `improve-codebase-architecture` procura módulos rasos e oportunidades de aprofundamento.
* `triage` qualifica e enriquece issues externas antes de despachá-las para agentes.

O texto upstream de `diagnosing-bugs` chama a construção do loop de feedback de "a habilidade em si" e lista formas de obtê-lo, em ordem: teste falhando, script HTTP contra o servidor local, chamada de CLI comparando a saída com um snapshot, script de navegador headless e reprodução de um trace capturado. A mesma habilidade exige ocultar todo segredo como `<REDACTED>` antes de mostrar comandos e saídas.

## 13. Sub-loop 6: Auditoria, Entrega e Ambiente

O objetivo é revisar a conformidade, redigir o PR e blindar o ambiente.

* `code-review` audita em dois eixos ortogonais: *Standards*, as regras documentadas do repositório, e *Spec*, a fidelidade à tarefa de origem.
* `pr` estrutura o pull request com a menor evidência visual necessária, antes e depois, e classifica o risco como porta de mão única ou de mão dupla.
* `retro` converte falhas mecânicas em checagens automáticas e divergências em regras de linter ou do `AGENTS.md`.
* `wizard` gera scripts interativos para etapas que exigem presença humana, como tokens e senhas.
* `wait-what` re-explica um conceito sem jargão quando o usuário se perde.

**Saídas:** `docs/reports/CODE_REVIEW.md`, corpo do PR e atualizações no `AGENTS.md`.

O upstream de `code-review` acrescenta um detalhe de mecânica: os dois eixos rodam em subagentes paralelos, para que um não contamine o contexto do outro, sobre o diff de três pontos entre um ponto fixo informado pelo usuário e o `HEAD`.
