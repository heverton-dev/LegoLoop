---
id: INTENT-I6
status: DRAFT
date: 2026-11-13
schema_version: 2
upstream: []
covers: [Q1, Q2, Q3]
---

# 🎯 INTENT: Migração dos DAGs Legados para Schema v2 e Roteamento Honesto (Iniciativa I6)

> **Status:** DRAFT  
> **Data:** 2026-11-13  
> **Autor:** i6-preparer (Arquiteto I6)  
> **Iniciativa:** I6 — Migração dos Legados  
> **Alvo:** `docs/dags/*.json` (9 arquivos) e `loops/*/pipeline.dag.json` (10 arquivos)  
> **Sessão:** `8e0c5643-39a6-4215-ac57-d16e07b7b41f`  
> **Confiança:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

* **O Problema Real:** Todos os 19 DAGs alvo em `docs/dags/*.json` (9 arquivos, excluído `2026-11-10_DAG_MATRIZ_PAPEIS.json` tratado na I5) e `loops/*/pipeline.dag.json` (10 arquivos) operam sem `schema_version: 2` e mantêm roteamento legado (`code-fast`, `code-pro`, `code-nano`, `flash`, `harness` ausente e `routing.via` ausente).
* **Quem Sofre:** O motor de validação (`legoloop validate`), o resolvedor de roteamento da I3 (`internal/roteamento`) e os executores hermeticamente isolados que não conseguem aplicar `routing.via`, `routing.harness`, `routing.model` e `routing.effort` nos Processos legados.
* **O Porquê:** As iniciativas I2 (`schema_version: 2`) e I3 (`roteamento honesto`) já estabilizaram os contratos no código Go. Migrar os 19 DAGs elimina o modo legado e coloca todo o Inventário sob validação estrita e roteamento homologado.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** `legoloop validate` emite aviso de modo legado nos 19 DAGs porque `schema_version` é `0` e `routing.via` está ausente.
  * **Amanhã:** Os 19 DAGs declaram `schema_version: 2`, roteamento I3 completo (`harness`, `via`, `model`, `effort`) e passam em `legoloop validate` sem modo legado.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Inventariar os 19 DAGs legados em `docs/dags/*.json` e `loops/*/pipeline.dag.json` com mapeamento determinístico para `schema_version: 2` e roteamento I3.
* [ ] Objetivo 2: Migrar os 19 DAGs em camadas paralelas por pasta/arquivo com `write_boundary` disjunto, preservando a topologia de dependências e `spec_ref` existentes.
* [ ] Objetivo 3: Validar todos os 19 DAGs migrados com `legoloop validate`, `go test -count=1 ./internal/dag/...` e `taxlinter`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não alterar `docs/dags/2026-11-10_DAG_MATRIZ_PAPEIS.json` (escopo exclusivo da Iniciativa I5).
* ⛔ Não alterar código Go de produção em `internal/` ou `cmd/`.
* ⛔ Não alterar contratos de `TEMPLATE_ROUTE.json` ou `TEMPLATE_DAG.json`.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `script` | Inventário de DAGs legados | Listar arquivos e nós sem `schema_version: 2` |
| `PECA-02` | `unit_task` | Migração de `docs/dags/*.json` | Atualizar 9 DAGs em lotes disjuntos |
| `PECA-03` | `unit_task` | Migração de `loops/*/pipeline.dag.json` | Atualizar 10 DAGs em lotes disjuntos por pasta |
| `PECA-04` | `script` | `legoloop validate` | Validação mecânica v2 dos 19 DAGs |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** Médio (edição em lote de 19 manifestos JSON, mitigado por `write_boundary` disjunto e `stash_changes` por Ticket).
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem autorizada (`claude` via `oficial`, modelos `haiku` e `sonnet`). Nenhum segredo envolvido.
* **Contrato de Evidência Factual:** `go run ./cmd/legoloop validate <arquivo>` retorna código `0` sem aviso de modo legado para cada um dos 19 DAGs.
* **Questões Resolvidas (Q):**

| ID | Questão | Decisão |
| :--- | :--- | :--- |
| `Q1` | Como evitar conflito de edição entre workers? | Camadas paralelas por pasta/arquivo com `write_boundary: strict` e `context_out` disjunto. |
| `Q2` | Qual arquivo fica fora do escopo da I6? | `docs/dags/2026-11-10_DAG_MATRIZ_PAPEIS.json`, tratado pela I5. |
| `Q3` | Qual padrão de roteamento aplicar nos nós migrados? | Roteamento I3 homologado (`harness: claude`, `via: oficial`, `model: haiku|sonnet|opus`, `effort: low|medium|high`). |

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todos os 19 DAGs alvo declaram `schema_version: 2` e `routing.via` válido | Validação CLI sem erro e sem modo legado | `go run ./cmd/legoloop validate <dag>` com exit code 0 |
| 2 | Zero conflito de escrita entre Tickets paralelos | Conjuntos `context_out` disjuntos verificados | Inspeção do `ROUTE` e `DAG` da I6 |
| 3 | Rastreabilidade do Quinteto da I6 íntegra | Gate de rastreabilidade sem órfãos | `go run ./tools/ticketgate -trace ...` |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md`](../specs/2026-11-13_SPEC_MIGRACAO_LEGADOS.md)
* **Gate de Transição:** Aprovação humana explícita antes da execução dos Épicos de migração.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção da Iniciativa I6 atende à sua expectativa para iniciar a especificação formal?*
