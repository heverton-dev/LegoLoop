# PARTE II — Meso: O Loop Mestre de Engenharia

## 4. Sete estágios, seis gates

O loop mestre leva uma ideia à produção em sete estágios. Entre cada par de estágios há um gate binário, de G1 a G6.

```{=typst}
#esteira((
  ("1. Alinhamento e Contexto", "using-agent-skills, interview-me, idea-refine, context-engineering"),
  ("2. Restrições e Especificação", "constraint-driven, spec-driven, api-and-interface-design"),
  ("3. Fatiamento e Planejamento", "planning-and-task-breakdown"),
  ("4. Construção TDD e UI", "source-driven, test-driven, frontend-ui, incremental"),
  ("5. Runtime Real e Web Performance", "browser-testing-with-devtools, performance-optimization"),
  ("6. Auditoria, Hardening e Contestação", "security, doubt-driven, simplification, review"),
  ("7. Migração, Observabilidade e Ship", "deprecation, observability, ci-cd, git, shipping"),
), gates: ("G1: alvo e intenção aprovados", "G2: spec e orçamentos fixados", "G3: DAG de tickets e testes definidos", "G4: suíte TDD 100% verde", "G5: CWV e runtime auditados", "G6: zero bloqueios em review"))
```

### Estágio 1: Alinhamento e Contexto

`using-agent-skills` roteia a intenção para as habilidades certas. `interview-me` e `idea-refine` conduzem a entrevista socrática. `context-engineering` delimita o escopo e registra a baseline: os números e o estado dos testes antes de qualquer mudança.

**G1.** Alvo e intenção aprovados pelo usuário.

### Estágio 2: Restrições e Especificação

`constraint-driven-development` produz o `CONSTRAINTS.md`. `spec-driven-development` fixa os critérios de aceite. `api-and-interface-design` define contratos e tipagem.

**G2.** Spec e orçamentos de performance fixados. A partir daqui, "rápido" deixa de ser uma opinião e passa a ser um número, por exemplo LCP de até 2,5 segundos.

### Estágio 3: Fatiamento e Planejamento

`planning-and-task-breakdown` decompõe a spec em um grafo de tickets verticais.

**G3.** O DAG de tickets está aprovado e cada ticket tem o teste que vai provar sua conclusão.

### Estágio 4: Construção TDD, Diagnóstico e UI

`source-driven-development` busca a documentação primária. `test-driven-development` conduz o ciclo vermelho-verde. `frontend-ui-engineering` e `incremental-implementation` constroem a fatia. Se algo quebra no caminho, `debugging-and-error-recovery` entra como desvio para encontrar a causa raiz.

**G4.** Suíte TDD 100% verde.

### Estágio 5: Runtime Real e Web Performance

`browser-testing-with-devtools` inspeciona o comportamento no navegador real. `performance-optimization` mede os Core Web Vitals e faz o profiling.

**G5.** CWV e runtime auditados, com números antes e depois.

### Estágio 6: Auditoria Multi-Eixo, Hardening e Contestação

`security-and-hardening` cobre o OWASP Top 10. `doubt-driven-development` contesta os achados e as premissas. `code-simplification` e `code-review-and-quality` fecham a revisão.

**G6.** Zero bloqueios em review.

### Estágio 7: Migração, Observabilidade e Ship

`deprecation-and-migration` planeja a migração sem downtime. `observability-and-instrumentation` garante logs e métricas. `ci-cd-and-automation`, `git-workflow-and-versioning` e `shipping-and-launch` levam à entrega.

## 5. O que os gates protegem

| Gate                 | Sem ele, o risco é                                          |
| :------------------- | :---------------------------------------------------------- |
| G1                   | Resolver com perfeição o problema errado                    |
| G2                   | Discutir "está rápido?" sem um número combinado             |
| G3                   | Tickets que não se sabe quando estão prontos                |
| G4                   | Medir performance de um código que ainda está quebrado      |
| G5                   | Auditar segurança de uma versão que ainda vai mudar         |
| G6                   | Entregar com um achado crítico em aberto                    |

```{=typst}
#painel("A ordem entre G4 e G5")[
  A performance é medida depois da suíte verde, e não antes, porque medir um código com testes falhando produz números de um comportamento que ainda vai mudar. A mesma lógica coloca a auditoria de segurança depois da medição: as otimizações do estágio 5 podem introduzir problemas que o estágio 6 precisa ver.
]
```

## 6. O loop mestre e o Playbook L

A skill `aidd-osmani` tem sua própria versão do fluxo completo, o Playbook L, acionado por pedidos como "fluxo completo na rota X". Ele parte de uma rota existente, e não de uma ideia, e por isso começa auditando.

| Playbook L, estágio    | Parada           | Equivale no loop mestre                  |
| :--------------------- | :--------------- | :--------------------------------------- |
| 1. Alvo                | P1               | Estágio 1, alinhamento e baseline        |
| 2. Auditoria           | P2               | Estágio 6, antecipado e só leitura       |
| 3. Spec e plano        | P3               | Estágios 2 e 3                           |
| 4. Execução            | P4               | Estágio 4, um ticket por vez             |
| 5. Prova e entrega     | —                | Estágios 5 a 7                           |

Cada estágio do Playbook L respeita o limite de três a cinco passos e grava o `ESTADO.md` ao terminar. O fluxo acaba com commits locais e um texto de PR proposto; push e PR só acontecem com ordem do usuário.
