# PARTE VI — Apêndices

## Apêndice A: Glossário

**ADR.** Registro de decisão arquitetural; documenta uma escolha difícil de reverter e o motivo.

**Costura (seam).** Fronteira pública em que o comportamento é observado sem olhar para dentro do módulo. Os testes moram nas costuras.

**DAG.** Grafo acíclico direcionado; aqui, o grafo de tickets com suas arestas de bloqueio.

**Grilling.** Entrevista socrática em rodadas que separa fatos do agente de decisões do usuário.

**Handoff.** Passagem estéril de contexto para outra sessão, pasta ou harness.

**Módulo profundo.** Módulo com interface estreita que esconde implementação rica.

**Porta de mão única.** Mudança difícil ou impossível de reverter; exige mais evidência no PR.

**Smart Zone.** Faixa da janela de contexto em que o modelo raciocina com clareza; o manual a estima em até cerca de 150 mil tokens.

**Teste tautológico.** Teste cujo valor esperado é calculado do mesmo jeito que o código, e que por isso nunca falha.

**Tracer bullet.** Ticket vertical que entrega um comportamento observável atravessando todas as camadas.

## Apêndice B: Relatório de conformidade

Comparação entre o manual de 01/11/2026 e o repositório em 07/11/2026.

| Ponto verificado        | Manual                    | Código ou upstream            | Situação       |
| :---------------------- | :------------------------ | :---------------------------- | :------------- |
| Arquivo do loop         | `matt-pocock.loop.json`     | Existe em `loops/`              | Conforme       |
| Validação do DAG        | Não citada                | 6 nós, sem ciclos             | Conforme       |
| Sub-loop 2, Exploração  | Sub-loop próprio          | Sem nó no JSON                | Diverge        |
| Sub-loop 5, Diagnóstico | Sub-loop próprio          | Sem nó no JSON                | Diverge        |
| Spec e tickets          | Um sub-loop               | Dois nós separados            | Diverge        |
| Modelos dos nós 1 e 4   | sonnet/haiku, sonnet/qwen | Iguais em `routing`             | Conforme       |
| Modelos dos nós 2,3,5,6 | Declarados na matriz      | Ausentes em `routing`           | Diverge        |
| Template do Studio      | Ícone `git`                 | Ícone `git`                     | Conforme       |
| Ciclo do TDD            | Red-Green-Refactor        | Upstream: só vermelho e verde | Diverge        |

**Leitura do relatório.** O JSON reorganiza os seis sub-loops do manual: Exploração e Diagnóstico ficam fora da esteira declarada, por serem desvios sob demanda, e Especificação vira dois nós, spec e tickets. O upstream da habilidade `tdd` tira a refatoração do ciclo e a coloca no review; este livro segue o upstream.

## Apêndice C: Referência rápida das habilidades

As habilidades estão preservadas na cópia local da skill `aidd-pocock`, em `upstream/skills/`: 20 em `engineering/` e 7 em `productivity/`.

| Habilidade                        | Pasta            | Etapa          | Função resumida                         |
| :-------------------------------- | :--------------- | :------------- | :-------------------------------------- |
| `grill-with-docs`                   | engineering      | 1              | Entrevista com rastro em arquivo        |
| `grill-me`                          | productivity     | 1              | Entrevista rápida sem estado            |
| `grilling`                          | productivity     | 1              | Primitivo da entrevista                 |
| `domain-modeling`                   | engineering      | 1              | Termos sobrecarregados e ADRs           |
| `prototype`                         | engineering      | 2              | Código descartável em branch isolada    |
| `handoff`                           | productivity     | 2              | Passagem estéril de contexto            |
| `research`                          | engineering      | 2              | Pesquisa em fontes primárias            |
| `to-spec`                           | engineering      | 3              | Spec a partir do consenso               |
| `to-tickets`                        | engineering      | 3              | Tickets verticais com bloqueios         |
| `wayfinder`                         | engineering      | 3              | Mapa de decisões para iniciativas vagas |
| `implement`                         | engineering      | 4              | Um ticket por sessão limpa              |
| `implement-spec`                    | engineering      | 4              | Tickets em paralelo com integração      |
| `tdd`                               | engineering      | 4              | Vermelho e verde em costuras combinadas |
| `codebase-design`                   | engineering      | 4              | Vocabulário de módulos profundos        |
| `diagnosing-bugs`                   | engineering      | 5              | Loop de feedback antes de teorizar      |
| `improve-codebase-architecture`     | engineering      | 5              | Caça a módulos rasos                    |
| `triage`                            | engineering      | 5              | Qualifica issues externas               |
| `code-review`                       | engineering      | 6              | Review em Standards e Spec              |
| `pr`                                | engineering      | 6              | PR com evidência mínima e risco         |
| `retro`                             | engineering      | 6              | Falhas viram regras do ambiente         |
| `wizard`                            | engineering      | 6              | Scripts para etapas humanas             |
| `wait-what`                         | productivity     | 6              | Re-explicação sem jargão                |

Fora dos sub-loops ficam `ask-matt` e `setup-matt-pocock-skills`, em `engineering/`, e `teach`, `to-questionnaire` e `writing-for-agents`, em `productivity/`.
