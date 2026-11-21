---
name: legoloop-addy-osmani
description: Use quando precisar executar o processo addy-osmani / legoloop-addy-osmani de engenharia e performance web baseado nas 25 skills e playbooks de Addy Osmani.
consumes: Artifact:Idea
produces: Artifact:Spec
artifact: docs/reports/SHIP_REPORT.md
gate: file_exists('docs/reports/SHIP_REPORT.md')
suggested_model: sonnet
---

# 🧩 legoloop-addy-osmani — Engenharia e Web Performance Addy Osmani

> **Processo:** `addy-osmani`  
> **Versão:** 1.0.0  
> **Referência:** `processos/legoloop-addy-osmani/processo.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for necessário executar engenharia de alta performance web, otimização de Core Web Vitals e arquitetura resiliente seguindo o método e playbooks de Addy Osmani.  
**Quando NÃO usar:** Para processos genéricos sem foco em engenharia e performance web.

* **Consumes:** Requisitos do projeto e especificações em `docs/specs/`.
* **Produces:** Artefatos em `docs/intents/INTENT.md`, `docs/specs/SPEC.md`, `docs/plans/PLAN.md`, relatórios de TDD, métricas de performance e revisão.

---

## 2. Estrutura de Execução do Grafo DAG

O processo é composto por 7 subprocessos sequenciais e auditáveis:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `subprocesso-1-alinhar-contexto` | Alinhamento & Contexto | `skills/aidd-osmani` | `docs/intents/INTENT.md` |
| `subprocesso-2-especificar-restricoes` | Especificação & Restrições | `skills/aidd-osmani` | `docs/specs/SPEC.md` |
| `subprocesso-3-planejar-tickets` | Decomposição Tracer-Bullet | `skills/aidd-osmani` | `docs/plans/PLAN.md` |
| `subprocesso-4-construir-tdd` | Construção TDD Cirúrgica | `skills/aidd-osmani` | `docs/reports/TDD_REPORT.md` |
| `subprocesso-5-runtime-webperf` | Runtime Real & Web Vitals | `skills/aidd-osmani` | `docs/reports/PERF_REPORT.md` |
| `subprocesso-6-review-hardening` | Auditoria 5 Eixos & Hardening | `skills/aidd-osmani` | `docs/reports/CODE_REVIEW.md` |
| `subprocesso-7-observabilidade-ship` | Observabilidade & Ship | `skills/aidd-osmani` | `docs/reports/SHIP_REPORT.md` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Alinhamento & Contexto (`subprocesso-1-alinhar-contexto`):** Conduzir alinhamento de escopo e produzir `docs/intents/INTENT.md`.
2. **Especificação & Restrições (`subprocesso-2-especificar-restricoes`):** Elaborar especificação técnica detalhada com restrições arquiteturais em `docs/specs/SPEC.md`.
3. **Decomposição Tracer-Bullet (`subprocesso-3-planejar-tickets`):** Decompor requisitos em tickets incrementais verticais em `docs/plans/PLAN.md`.
4. **Construção TDD Cirúrgica (`subprocesso-4-construir-tdd`):** Executar ciclo Red-Green-Refactor estrito com zero stubs e gerar `docs/reports/TDD_REPORT.md`.
5. **Runtime Real & Web Vitals (`subprocesso-5-runtime-webperf`):** Coletar métricas em runtime de LCP, FID/INP, CLS e gerar `docs/reports/PERF_REPORT.md`.
6. **Auditoria 5 Eixos & Hardening (`subprocesso-6-review-hardening`):** Executar revisão de segurança OWASP, concorrência e tipagem gerando `docs/reports/CODE_REVIEW.md`.
7. **Observabilidade & Ship (`subprocesso-7-observabilidade-ship`):** Configurar telemetria de produção e formalizar liberação com `docs/reports/SHIP_REPORT.md`.

---

## 4. Invariantes & Regras Negativas

- **Zero Stubs:** Nenhuma implementação simulada ou pendente com `TODO`.
- **Exit Code:** Todo comando de verificação deve concluir com `exit code 0`.
- **Stop-Gate Humano:** Respeitar aprovação formal humana antes de mesclar ou liberar artefatos em produção.
