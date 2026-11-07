# PARTE VI — Apêndices

## Apêndice A: Glossário

**Baseline.** Medição do estado atual, feita antes de qualquer mudança, que serve de comparação para provar ganho ou regressão.

**CLS.** Cumulative Layout Shift; mede quanto o conteúdo se desloca enquanto a página carrega.

**CONSTRAINTS.md.** Arquivo que registra as restrições ativas do projeto: orçamentos de performance, tetos de complexidade e contratos de segurança.

**Doubt-Driven Development.** Contestação adversarial de achados e premissas antes de aceitá-los.

**Expandir e contrair.** Estratégia de migração em que o novo formato convive com o antigo até a remoção segura deste.

**INP.** Interaction to Next Paint; mede o tempo entre uma interação e a próxima atualização visual.

**LCP.** Largest Contentful Paint; mede quando o maior elemento visível termina de ser desenhado.

**Orçamento de performance.** Limite numérico que nenhuma mudança pode ultrapassar.

**Prove-It.** Padrão em que a correção de um bug só começa depois de um teste que falha pelo motivo do bug.

**Tracer bullet.** Ticket vertical que entrega um comportamento observável atravessando todas as camadas.

## Apêndice B: Relatório de conformidade

Comparação entre o manual de 02/11/2026 e o repositório em 07/11/2026.

| Ponto verificado           | Manual                      | Código ou upstream               | Situação |
| :------------------------- | :-------------------------- | :------------------------------- | :------- |
| Número de sub-loops        | 8                           | 7 nós no JSON                    | Diverge  |
| Nome do loop no JSON       | Não citado                  | Diz "(8 Sub-Loops)" com 7 nós    | Diverge  |
| Nó de diagnóstico          | `osmani-debugging-recovery` | Ausente no JSON                  | Diverge  |
| Validação do DAG           | Não citada                  | 7 nós, sem ciclos                | Conforme |
| Nome do template no Studio | "...& Engineering Loop"     | "...& Rigor Loop"                | Diverge  |
| Parâmetros do `run`        | `--param`, `--verbose`      | Não existem; usa-se `--var`      | Diverge  |
| Endereço do Studio         | `localhost:3000`            | `127.0.0.1:7878`                 | Diverge  |
| Quinto eixo de revisão     | Manutenibilidade            | Legibilidade e simplicidade      | Diverge  |
| `metrics-guide.md`         | Citado                      | Não encontrado no upstream       | Ausente  |
| `profiling-guide.md`       | Citado                      | Não encontrado no upstream       | Ausente  |
| `threat-modeling.md`       | Citado                      | Seção do `security-checklist.md` | Diverge  |
| `idea-refiner.py`          | Citado                      | Não encontrado no upstream       | Ausente  |
| Total de habilidades       | 25                          | 25 em `upstream/skills/`         | Conforme |

**Leitura do relatório.** As divergências de operação são as mais importantes para o leitor: os comandos do manual com `--param` e `--verbose` falham no CLI atual, que aceita `--var chave=valor`; o Studio sobe por padrão em `127.0.0.1:7878`. Este livro usa a forma real em todos os exemplos. O nó de diagnóstico ausente no JSON é coerente com o papel de desvio que o diagnóstico tem no loop mestre, mas o nome do pipeline ainda anuncia oito sub-loops.

## Apêndice C: Referência rápida das habilidades

| Habilidade                          | Estágio | Função resumida                         |
| :---------------------------------- | :------ | :-------------------------------------- |
| `using-agent-skills`                | 1       | Roteia a intenção para as habilidades   |
| `interview-me`                      | 1       | Entrevista até entender a intenção real |
| `idea-refine`                       | 1       | Refina um conceito cru em opções        |
| `context-engineering`               | 1       | Escopo cirúrgico e baseline             |
| `constraint-driven-development`     | 2       | Restrições ativas em `CONSTRAINTS.md`   |
| `spec-driven-development`           | 2       | Spec com critérios de aceite            |
| `api-and-interface-design`          | 2       | Contratos de API e tipagem              |
| `planning-and-task-breakdown`       | 3       | DAG de tickets verticais                |
| `source-driven-development`         | 4       | Decisões ancoradas em docs oficiais     |
| `test-driven-development`           | 4       | Vermelho e verde; Prove-It              |
| `incremental-implementation`        | 4       | Fatias pequenas e testadas              |
| `frontend-ui-engineering`           | 4       | UI acessível e responsiva               |
| `debugging-and-error-recovery`      | 4       | Causa raiz com reprodução               |
| `browser-testing-with-devtools`     | 5       | Inspeção no navegador real              |
| `performance-optimization`          | 5       | Medição, profiling e otimização         |
| `security-and-hardening`            | 6       | OWASP Top 10 e endurecimento            |
| `doubt-driven-development`          | 6       | Contestação adversarial                 |
| `code-simplification`               | 6       | Menos código, mesmo comportamento       |
| `code-review-and-quality`           | 6       | Revisão em cinco eixos                  |
| `deprecation-and-migration`         | 7       | Expandir e contrair                     |
| `observability-and-instrumentation` | 7       | Logs, métricas e traces                 |
| `documentation-and-adrs`            | 7       | ADRs e documentação                     |
| `git-workflow-and-versioning`       | 7       | Commits atômicos e versionamento        |
| `ci-cd-and-automation`              | 7       | Pipelines e gates automáticos           |
| `shipping-and-launch`               | 7       | Checklist, rollout e rollback           |
