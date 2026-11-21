# PARTE I — Macro: Medir, Restringir, Provar

## 1. Dois problemas com a mesma raiz

O manual apresenta o método de Addy Osmani como a síntese de mais de duas décadas de engenharia web no Google Chrome e no ecossistema JavaScript. O método combate dois problemas que, à primeira vista, parecem distintos.

O primeiro é o **achismo de desempenho**. Alguém diz que a página está lenta, outra pessoa troca uma biblioteca, uma terceira adiciona cache, e ninguém sabe se a página ficou mais rápida, porque ninguém mediu antes. Às vezes ficou mais lenta.

O segundo é a **proliferação de código assistido por IA**. Agentes produzem muito código, rápido. Sem barreiras, esse código cresce em complexidade, ignora limites de segurança e acumula comportamentos que ninguém especificou.

A raiz comum é a ausência de um número ou de uma regra verificável antes de agir. Os dois problemas se resolvem do mesmo jeito: medir, fixar limites por escrito e provar cada mudança.

```{=typst}
#frase-destaque[Sem número de partida, otimizar é adivinhar.]
```

## 2. Os seis pilares

### Pilar 1: medir antes de otimizar

Toda otimização exige telemetria prévia: LCP, INP, CLS, TTFB, perfil de memória, tamanho de bundle. E exige medição posterior com prova de ganho. Uma otimização sem as duas medições não é aceita, mesmo que pareça óbvia.

### Pilar 2: desenvolvimento orientado a restrições

Limites de arquitetura, tetos de complexidade ciclomática e contratos de segurança são barreiras ativas, registradas em `CONSTRAINTS.md`. Nenhum agente pode violá-las ou silenciá-las. A diferença entre uma restrição e uma recomendação é que a restrição bloqueia a entrega.

### Pilar 3: ancoragem na documentação oficial

Agentes de IA às vezes inventam sintaxes ou comportamentos de bibliotecas. O método exige que decisões sobre APIs e padrões se baseiem na documentação primária oficial das bibliotecas e frameworks em uso, e não na memória do modelo.

### Pilar 4: fatiamento fino e prova por testes

Um bug exige um teste vermelho que reproduza a falha antes da correção. O manual chama esse padrão de *Prove-It*. Funcionalidades avançam em fatias atômicas e testadas, sem stubs e sem quebrar a suíte.

### Pilar 5: runtime real e auditoria em cinco eixos

Testes unitários e mocks não substituem a inspeção do navegador de verdade: o DOM, o CSS renderizado, os deslocamentos de layout, os erros de console e o comportamento em viewport real. A revisão final avalia cinco eixos e é acompanhada de contestação adversarial, o *Doubt-Driven Development*.

### Pilar 6: entrega observável e segura

O ciclo termina com rastreabilidade: logs estruturados, métricas, traces, um plano de migração do tipo expandir e contrair, e checklists de entrega rigorosos.

```{=typst}
#painel("O padrão Prove-It")[
  Um relato de bug não é corrigido até existir um teste que falha pelo motivo exato do bug. A correção mínima faz esse teste passar, e ele fica na suíte para sempre. Sem o vermelho primeiro, não há como provar que a correção resolveu o problema relatado e não outro.
]
```

## 3. O que o leitor leva desta parte

| Pilar                             | Pergunta que impõe antes de agir               |
| :-------------------------------- | :--------------------------------------------- |
| Medir antes de otimizar           | Qual é o número hoje?                          |
| Restrições                        | Qual limite esta mudança não pode passar?      |
| Documentação oficial              | Onde a documentação diz que funciona assim?    |
| Fatiamento e Prove-It             | Qual teste falha agora e passará depois?       |
| Runtime real e cinco eixos        | O navegador real confirma? Quem contestou?     |
| Entrega observável                | Como saberemos que quebrou em produção?        |
