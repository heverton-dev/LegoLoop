---
id: INTENT-intv_1794873379697350300
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [INTERVIEW-intv_1794873379697350300]
covers: [Q1]
---

# 🎯 INTENT: Mega Plano LegoLoop para produção: saneamento do docs/ (3598 arquivos, 1.16GB), templates de nascimento de processo (triagem, memoria_execucao, anticorpos), VSA por processo (fatia vertical autocontida), auditoria profunda agnostica universal dos 6 verbos restantes via SWARM, processo auditoria-atomica universal (qualquer alvo), gestao concorrente de memoria de execucao (sessao+PID+lock), Motor de Anticorpos persistente (anticorpos.json, vacinacao de reincidencia), triple via de invocacao (CLI, skill/harness, Studio), renomeacao de loops/ para processos/ (conflito Lei 15), MCP unico via mcp-router, adecacao de todos os processos existentes ao novo normal, verificacao multi-SO e multi-harness, provas ausentes (-race, LLM real, grafo semantico), higiene (taxlinter 14, funclen 3, 2 testes flaky, push pendente). Tudo em modo SWARM com todas as pecas, skills, MCPs, gates e scripts disponiveis, parada so em stop-gates humanos.

> **Status:** DRAFT  
> **Data:** 2026-11-16  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** intv_1794873379697350300  
> **Confiança da Entrevista:** 94% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Mega Plano LegoLoop para produção: saneamento do docs/ (3598 arquivos, 1.16GB), templates de nascimento de processo (triagem, memoria_execucao, anticorpos), VSA por processo (fatia vertical autocontida), auditoria profunda agnostica universal dos 6 verbos restantes via SWARM, processo auditoria-atomica universal (qualquer alvo), gestao concorrente de memoria de execucao (sessao+PID+lock), Motor de Anticorpos persistente (anticorpos.json, vacinacao de reincidencia), triple via de invocacao (CLI, skill/harness, Studio), renomeacao de loops/ para processos/ (conflito Lei 15), MCP unico via mcp-router, adecacao de todos os processos existentes ao novo normal, verificacao multi-SO e multi-harness, provas ausentes (-race, LLM real, grafo semantico), higiene (taxlinter 14, funclen 3, 2 testes flaky, push pendente). Tudo em modo SWARM com todas as pecas, skills, MCPs, gates e scripts disponiveis, parada so em stop-gates humanos.
* **Quem Sofre:** Desenvolvedores e operadores do ecossistema LegoLoop (CLI, Skill, Studio).
* **O Porquê:** Eliminação de dívida técnica crítica, resolução de conflito legal (Lei 15) e necessidade de escala industrial com concorrência.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Mega Plano LegoLoop para produção: saneamento do docs/ (3598 arquivos, 1.16GB), templates de nascimento de processo (triagem, memoria_execucao, anticorpos), VSA por processo (fatia vertical autocontida), auditoria profunda agnostica universal dos 6 verbos restantes via SWARM, processo auditoria-atomica universal (qualquer alvo), gestao concorrente de memoria de execucao (sessao+PID+lock), Motor de Anticorpos persistente (anticorpos.json, vacinacao de reincidencia), triple via de invocacao (CLI, skill/harness, Studio), renomeacao de loops/ para processos/ (conflito Lei 15), MCP unico via mcp-router, adecacao de todos os processos existentes ao novo normal, verificacao multi-SO e multi-harness, provas ausentes (-race, LLM real, grafo semantico), higiene (taxlinter 14, funclen 3, 2 testes flaky, push pendente). Tudo em modo SWARM com todas as pecas, skills, MCPs, gates e scripts disponiveis, parada so em stop-gates humanos.
  * **Amanhã:** LegoLoop em estado de produção com arquitetura VSA, saneamento de 1.1GB de docs, memória concorrente segura e auditoria SWARM funcional.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: LegoLoop em estado de produção com arquitetura VSA, saneamento de 1.1GB de docs, memória concorrente segura e auditoria SWARM funcional.
* [ ] Objetivo 2: Saneamento concluído (docs/ reduzido), testes `-race` e `taxlinter` aprovados, auditoria universal validando os 6 verbos sem colisões de PID.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Refatoração da interface visual do Studio ou criação de novos verbos além dos já estabelecidos.

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
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Saneamento concluído (docs/ reduzido), testes `-race` e `taxlinter` aprovados, auditoria universal validando os 6 verbos sem colisões de PID.
* **Restrições Inegociáveis:**
  * Sem middleware centralizado para locks; gestão de memória via biblioteca compartilhada com lock por arquivo (Session+PID); conformidade estrita com a Lei 15.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 94%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | Relatório de testes do ticket |
| 3 | Saneamento concluído (docs/ reduzido), testes `-race` e `taxlinter` aprovados, auditoria universal validando os 6 verbos sem colisões de PID. | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-16_SPEC_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md`](../specs/2026-11-16_SPEC_MEGA_PLANO_LEGOLOOP_PARA_PRODU.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 94% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
