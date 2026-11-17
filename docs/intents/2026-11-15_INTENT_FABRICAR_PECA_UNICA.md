---
id: INTENT-fabricar-fabricar-peca-unica
status: DRAFT
date: 2026-11-15
schema_version: 2
upstream: [INTERVIEW-fabricar-fabricar-peca-unica]
covers: [Q1]
---

# 🎯 INTENT: Fabricar Peca Unica

> **Status:** DRAFT  
> **Data:** 2026-11-15  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** fabricar-fabricar-peca-unica  
> **Confiança da Entrevista:** 0% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Fabricar Peca Unica
* **Quem Sofre:** Não informado na entrevista
* **O Porquê:** Não informado na entrevista
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: Fabricar Peca Unica
  * **Amanhã:** Fabricar uma peça única (nó solto) de qualquer tipo, refinar com --adicionar e compor com --depende; triagem só com mais de um nó, exigida por validate e fechar.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Fabricar uma peça única (nó solto) de qualquer tipo, refinar com --adicionar e compor com --depende; triagem só com mais de um nó, exigida por validate e fechar.
* [ ] Objetivo 2: Não informado na entrevista

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | `unit_task` | Execução das etapas | Cada etapa grava o seu artefato e passa por um gate de conteúdo. |
| `PECA-03` | `script` | `legoloop validate` | Validação mecânica do grafo do processo. |
| `PECA-04` | `stop_gate` | Aprovação humana | Pausa antes de qualquer gravação definitiva. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** A definir no Stop-Gate 01
* **Política de Privacidade & Orçamento Cognitivo:** A definir no Stop-Gate 01
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Não informado na entrevista
* **Restrições Inegociáveis:**
  * Nenhuma restrição registrada na entrevista
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** Nenhuma premissa registrada na entrevista.
* **Riscos e Perguntas em Aberto:** Confiança de 0%: as lacunas restantes devem ser tratadas no Stop-Gate 01.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Cada etapa grava o artefato prometido, e o gate de conteúdo dela passa. | Gate da etapa avaliado sobre o artefato | `legoloop run` com exit code 0 |
| 2 | O grafo do processo é válido antes de executar. | Validação mecânica do grafo | `legoloop validate` com exit code 0 |
| 3 | Não informado na entrevista | Critério de sucesso acordado na entrevista | A definir na SPEC como `AC-XX` binário |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-15_SPEC_FABRICAR_PECA_UNICA.md`](../specs/2026-11-15_SPEC_FABRICAR_PECA_UNICA.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 0% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
