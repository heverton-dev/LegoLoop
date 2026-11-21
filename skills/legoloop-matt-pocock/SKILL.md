---
name: legoloop-matt-pocock
description: Use quando precisar executar o processo matt-pocock / legoloop-matt-pocock de grilling socrático, especificações invariantes e TDD cirúrgico baseado no método Matt Pocock.
consumes: Artifact:Idea
produces: Artifact:Spec
artifact: docs/reports/RETRO.md
gate: file_exists('docs/reports/RETRO.md')
suggested_model: sonnet
---

# 🧩 legoloop-matt-pocock — Engenharia e Grilling Socrático Matt Pocock

> **Processo:** `matt-pocock`  
> **Versão:** 1.0.0  
> **Referência:** `processos/legoloop-matt-pocock/processo.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for necessário alinhar o domínio através de entrevista socrática (grilling), definir especificações com restrições invariantes estritas, decompor em grafo acíclico (DAG) e aplicar TDD com revisão em 2 eixos.  
**Quando NÃO usar:** Para demandas triviais sem necessidade de grilling ou alinhamento de premissas.

* **Consumes:** Requisitos e especificações iniciais em `docs/specs/`.
* **Produces:** `GLOSSARY.md`, `docs/specs/SPEC.md`, `docs/plans/PLAN.md`, relatórios de TDD, code review e retrospectiva em `docs/reports/RETRO.md`.

---

## 2. Estrutura de Execução do Grafo DAG

O processo engloba 6 subprocessos determinísticos:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `subprocesso-1-alinhar-dominio` | Grilling & Alinhamento de Domínio | `skills/aidd-pocock` | `GLOSSARY.md` |
| `subprocesso-2-especificar` | Especificação & Restrições Invariantes | `skills/aidd-pocock` | `docs/specs/SPEC.md` |
| `subprocesso-3-planejar-tickets` | Tickets em Grafo Acíclico (DAG) | `skills/aidd-pocock` | `docs/plans/PLAN.md` |
| `subprocesso-4-construir-tdd` | Construção TDD Cirúrgica (Zero Stubs) | `skills/aidd-pocock` | `docs/reports/TDD_REPORT.md` |
| `subprocesso-5-review-auditar` | Code Review de 2 Eixos & Hardening | `skills/aidd-pocock` | `docs/reports/CODE_REVIEW.md` |
| `subprocesso-6-retro-ambiente` | Retroalimentação do Ambiente & Governança | `skills/aidd-pocock` | `docs/reports/RETRO.md` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Grilling & Alinhamento de Domínio (`subprocesso-1-alinhar-dominio`):** Questionar e desafiar pressupostos técnicos através de entrevista socrática e consolidar o vocabulário em `GLOSSARY.md`.
2. **Especificação & Restrições Invariantes (`subprocesso-2-especificar`):** Formalizar requisitos com contratos EARS e invariantes de domínio em `docs/specs/SPEC.md`.
3. **Tickets em Grafo Acíclico (DAG) (`subprocesso-3-planejar-tickets`):** Planejar o plano de execução decomposto em tickets orientados por dependências no arquivo `docs/plans/PLAN.md`.
4. **Construção TDD Cirúrgica (Zero Stubs) (`subprocesso-4-construir-tdd`):** Implementar testes reais e código de produção com cobertura integral sem stubs em `docs/reports/TDD_REPORT.md`.
5. **Code Review de 2 Eixos & Hardening (`subprocesso-5-review-auditar`):** Realizar revisão em dois eixos (corretude de arquitetura e segurança/concorrência) gerando `docs/reports/CODE_REVIEW.md`.
6. **Retroalimentação do Ambiente & Governança (`subprocesso-6-retro-ambiente`):** Registrar aprendizados técnicos e calibrar memórias do sistema em `docs/reports/RETRO.md`.

---

## 4. Invariantes & Regras Negativas

- **Zero Suposições:** Não assumir comportamento não documentado; executar grilling até sanar todas as dúvidas.
- **Zero Stubs:** Implementações devem ser 100% reais e testadas.
- **Conformidade de Gates:** Cada nó só é concluído após seu gate contratual ser satisfeito no sistema de arquivos.
