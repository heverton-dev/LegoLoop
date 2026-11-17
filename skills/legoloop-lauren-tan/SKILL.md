---
name: legoloop-lauren-tan
description: Use quando precisar executar o processo lauren-tan / legoloop-lauren-tan de rigor de engenharia, modelagem formal e higiene de código baseado no pstack / poteto-mode.
consumes: Artifact:Idea
produces: Artifact:Spec
artifact: docs/reports/RETRO.md
gate: file_exists('docs/reports/RETRO.md')
suggested_model: sonnet
---

# 🧩 legoloop-lauren-tan — Rigor de Engenharia Lauren Tan (pstack)

> **Processo:** `lauren-tan`  
> **Versão:** 1.0.0  
> **Referência:** `processos/lauren-tan/processo.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for exigido alto rigor de modelagem de domínio, verificação formal no artefato real e eliminação sistemática de slop/código inflado (pstack / poteto-mode).  
**Quando NÃO usar:** Para prototipagem descartável sem critérios formais de modelagem.

* **Consumes:** Contexto inicial em `docs/intents/`.
* **Produces:** `docs/intents/INTENT.md`, `docs/specs/SPEC.md`, `docs/plans/PLAN.md`, relatórios de TDD, relatório formal e retrospectiva em `docs/reports/RETRO.md`.

---

## 2. Estrutura de Execução do Grafo DAG

O processo organiza-se em 6 subprocessos estritos:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `subprocesso-1-intencao-dominio` | Intenção & Modelagem de Domínio | `skills/aidd-lauren` | `docs/intents/INTENT.md` |
| `subprocesso-2-arquitetura-exploracao` | Arquitetura & Exploração Competitiva | `skills/aidd-lauren` | `docs/specs/SPEC.md` |
| `subprocesso-3-decomposicao-verificavel` | Decomposição em Unidades Verificáveis | `skills/aidd-lauren` | `docs/plans/PLAN.md` |
| `subprocesso-4-construcao-tdd` | Construção Rigorosa TDD | `skills/aidd-lauren` | `docs/reports/TDD_REPORT.md` |
| `subprocesso-5-verificacao-artefato-real` | Verificação Formal no Artefato Real | `skills/aidd-lauren` | `docs/reports/REPORT.md` |
| `subprocesso-6-unslop-ship` | Higiene de Código (Unslop) & Ship | `skills/aidd-lauren` | `docs/reports/RETRO.md` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Intenção & Modelagem de Domínio (`subprocesso-1-intencao-dominio`):** Modelar entidades puras e restrições formais de domínio em `docs/intents/INTENT.md`.
2. **Arquitetura & Exploração Competitiva (`subprocesso-2-arquitetura-exploracao`):** Definir fronteiras arquiteturais e especificações estruturadas em `docs/specs/SPEC.md`.
3. **Decomposição em Unidades Verificáveis (`subprocesso-3-decomposicao-verificavel`):** Dividir entregáveis em unidades atômicas com testes no plano `docs/plans/PLAN.md`.
4. **Construção Rigorosa TDD (`subprocesso-4-construcao-tdd`):** Desenvolver testes unitários que falham antes da implementação e comprovar com `docs/reports/TDD_REPORT.md`.
5. **Verificação Formal no Artefato Real (`subprocesso-5-verificacao-artefato-real`):** Testar binários e ambientes em runtime real, compilando o relatório `docs/reports/REPORT.md`.
6. **Higiene de Código (Unslop) & Ship (`subprocesso-6-unslop-ship`):** Eliminar redundâncias, polir abstrações e registrar retro em `docs/reports/RETRO.md`.

---

## 4. Invariantes & Regras Negativas

- **Zero Slop:** Proibido código verboso, boilerplate sem propósito ou comentários redundantes.
- **Verificação no Real:** Testes em fixtures não substituem a prova no binário ou runtime real.
- **Hermeticidade:** Execução isolada em contexto estéril sem contaminação externa.
