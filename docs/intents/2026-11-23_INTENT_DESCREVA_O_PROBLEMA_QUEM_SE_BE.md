---
id: INTENT-intv_1795427227896673100
status: UNDER_REVIEW
date: 2026-11-23
schema_version: 2
upstream: [INTERVIEW-intv_1795427227896673100]
covers: [Q1]
---

# 🎯 INTENT: Descreva o problema, quem se beneficia e o resultado esperado.

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-23  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1795427227896673100  
> **Confiança da Entrevista:** 95% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Descreva o problema, quem se beneficia e o resultado esperado.
* **Quem Sofre:** Departamentos de Vendas e Operações.
* **O Porquê:** Resolver falhas de comunicação e rastreio de tarefas entre áreas.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista (ver O Problema Real).
  * **Amanhã:** Ferramenta de gestão de fluxo (Kanban) integrada bidirecionalmente ao CRM.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Ferramenta de gestão de fluxo (Kanban) integrada bidirecionalmente ao CRM.
* [ ] Objetivo 2: critérios de sucesso:
  * Sincronização automática de status entre Kanban e CRM sem intervenção manual.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Módulos de faturamento
* ⛔ app mobile nativo
* ⛔ integração com sistemas legados não-API.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | `unit_task` | Código de produção | Implementação da lógica via TDD Red-Green. |
| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  * Sincronização automática de status entre Kanban e CRM sem intervenção manual.
* **Restrições Inegociáveis:**
  * Stack Node.js e React; suporte a alto volume de tarefas (+100).
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 95%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Sincronização automática de status entre Kanban e CRM sem intervenção manual. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-23_SPEC_DESCREVA_O_PROBLEMA_QUEM_SE_BE.md`](../specs/2026-11-23_SPEC_DESCREVA_O_PROBLEMA_QUEM_SE_BE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
