# PARTE III — Meso: Os Sub-Loops

## 6. Visão geral

O manual decompõe o ecossistema de Addy Osmani em oito sub-loops, cada um com entradas e saídas contratadas. O arquivo do loop no repositório declara sete deles como nós; o oitavo, `osmani-debugging-recovery`, funciona como desvio do estágio de construção. O Apêndice B registra essa diferença.

| Sub-loop                     | Entrada                   | Saída contratada                         |
| :--------------------------- | :------------------------ | :--------------------------------------- |
| `osmani-alignment`           | Ideia ou demanda bruta    | Intenção com alvos e baseline            |
| `osmani-spec-constraints`    | Intenção validada         | Spec e `CONSTRAINTS.md` com orçamentos   |
| `osmani-task-breakdown`      | Spec e contratos          | Plano com DAG e testes por ticket        |
| `osmani-incremental-build`   | Fatias do plano           | Código sem stubs e testes verdes         |
| `osmani-debugging-recovery`  | Falha ou regressão        | Teste vermelho, patch e suíte verde      |
| `osmani-performance-webperf` | Rota em runtime local     | Relatório com métricas antes e depois    |
| `osmani-review-hardening`    | Código e diff consolidado | Auditoria em 5 eixos e diff simplificado |
| `osmani-ship-lifecycle`      | Código auditado           | ADRs, instrumentação e release           |

## 7. Alinhamento

**Habilidades:** `using-agent-skills`, `interview-me`, `idea-refine` e `context-engineering`.

O sub-loop desambigua a intenção, conduz a ideação socrática e isola o contexto necessário. A saída é um documento de intenção em `docs/intents/` com os alvos resolvidos e a baseline inicial. No manual, este sub-loop corresponde ao Playbook I, Definir.

## 8. Especificação e Restrições

**Habilidades:** `constraint-driven-development`, `spec-driven-development` e `api-and-interface-design`.

O sub-loop fixa orçamentos de performance, barreiras de qualidade e contratos de interface. As saídas são a spec em `docs/specs/` e o `CONSTRAINTS.md`. As receitas associadas são `/spec` e `/constraints`.

```{=typst}
#painel("Spec e restrição não são a mesma coisa")[
  A spec diz o que o sistema faz: "o usuário vê o resumo do pedido antes de pagar". A restrição diz o que nenhuma mudança pode violar: "LCP da rota de checkout até 2,5 segundos" ou "nenhuma função acima de complexidade ciclomática 10". A spec muda a cada funcionalidade; a restrição atravessa todas.
]
```

## 9. Decomposição em Tickets

**Habilidade:** `planning-and-task-breakdown`. **Receita:** `/plan`.

A spec vira fatias verticais finas, os *tracer bullets*, ordenadas em um DAG. Segundo o Playbook K da skill, cada ticket tem identificador, título, critério binário, bloqueios, arquivos envolvidos e o teste que prova sua conclusão.

## 10. Construção Incremental

**Habilidades:** `source-driven-development`, `test-driven-development`, `incremental-implementation` e `frontend-ui-engineering`. **Papel:** engenheiro de testes. **Receitas:** `/build` e `/test`.

A construção segue a documentação oficial, o ciclo vermelho-verde e implementa UI acessível. Os checklists de apoio são `testing-patterns.md` e `accessibility-checklist.md`, este baseado na WCAG 2.1. A saída é código de produção real, sem stubs, com testes unitários e de componente verdes.

## 11. Diagnóstico e Recuperação

**Habilidades:** `debugging-and-error-recovery` e `test-driven-development` no padrão Prove-It. **Papel:** engenheiro de testes.

Este sub-loop corresponde ao Playbook B, Corrigir. A skill descreve a sequência: reproduzir e registrar requisição, resposta e status; sem reprodução, parar. Depois, um teste que falha pelo motivo certo (gate vermelho), a correção mínima da causa raiz (gate verde) e a confirmação de que o sintoma original não se reproduz mais.

## 12. Performance Web

**Habilidades:** `browser-testing-with-devtools` e `performance-optimization`. **Papel:** auditor de performance web. **Receita:** `/webperf`.

O sub-loop mede Core Web Vitals, faz profiling de memória e valida no DevTools. A saída é um relatório em `docs/reports/` com métricas objetivas antes e depois. A Parte IV traz os limiares numéricos que o checklist de performance usa.

## 13. Review e Hardening

**Habilidades:** `code-review-and-quality`, `security-and-hardening`, `doubt-driven-development` e `code-simplification`. **Papéis:** revisor de código e auditor de segurança. **Receitas:** `/review` e `/code-simplify`.

O sub-loop audita em cinco eixos, cobre o OWASP Top 10, contesta os achados e simplifica sem perda de comportamento. Corresponde aos Playbooks A (Auditar), D (Simplificar) e F (Proteger). No Playbook F, cada correção de segurança precisa de um teste negativo que prove o bloqueio.

## 14. Ciclo de Entrega

**Habilidades:** `deprecation-and-migration`, `observability-and-instrumentation`, `documentation-and-adrs`, `git-workflow-and-versioning`, `ci-cd-and-automation` e `shipping-and-launch`. **Receita:** `/ship`.

O sub-loop planeja a migração sem downtime, ativa a observabilidade, automatiza o CI/CD e faz o release. Os checklists são `observability-checklist.md` e `definition-of-done.md`. As saídas são ADRs em `docs/adrs/`, instrumentação validada, commits semânticos e checklist de release completo. Corresponde aos Playbooks G (Entregar), H (Migrar) e L (Ciclo completo).
