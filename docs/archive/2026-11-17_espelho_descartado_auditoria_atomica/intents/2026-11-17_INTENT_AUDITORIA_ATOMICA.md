---
id: INTENT-fabricar-auditoria-atomica
status: UNDER_REVIEW
date: 2026-11-17
schema_version: 2
upstream: [INTERVIEW-fabricar-auditoria-atomica]
covers: [Q1]
---

# 🎯 INTENT: Auditoria universal agnóstica de alvo

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-17  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** fabricar-auditoria-atomica  
> **Confiança da Entrevista:** 100% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Auditoria universal agnóstica de alvo
* **Quem Sofre:** Operador do LegoLoop que audita verbos, processos, peças, funções, documentos, binários e repos externos
* **O Porquê:** A auditoria universal é o contrato do Épico 08 do Mega Plano
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Auditoria universal agnóstica de alvo
  * **Amanhã:** Veredito binário com evidências para qualquer alvo auditado

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Veredito binário com evidências para qualquer alvo auditado
* [ ] Objetivo 2: Todo alvo auditado recebe veredito APROVADO ou REPROVADO com evidência mecânica

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Correção automática do alvo reprovado

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
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Todo alvo auditado recebe veredito APROVADO ou REPROVADO com evidência mecânica
* **Restrições Inegociáveis:**
  * Critério binário com código de saída 0, 1 ou 2
  * Zero artefato órfão para alvo inexistente
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Nenhum registrado.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Todo alvo auditado recebe veredito APROVADO ou REPROVADO com evidência mecânica | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-17_SPEC_AUDITORIA_ATOMICA.md`](../specs/2026-11-17_SPEC_AUDITORIA_ATOMICA.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
