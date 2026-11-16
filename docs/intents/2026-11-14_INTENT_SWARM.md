---
id: INTENT-fabricar-swarm
status: DRAFT
date: 2026-11-14
schema_version: 2
upstream: [INTERVIEW-fabricar-swarm]
covers: [Q1]
---

# 🎯 INTENT: legoloop-swarm

> **Status:** DRAFT  
> **Data:** 2026-11-14  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** A definir na SPEC  
> **Sessão:** fabricar-swarm  
> **Confiança da Entrevista:** 0% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** legoloop-swarm
* **Quem Sofre:** Não informado na entrevista
* **O Porquê:** Não informado na entrevista
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Conforme a ideia inicial da entrevista: legoloop-swarm
  * **Amanhã:** Orquestrar execucao paralela de pipeline.dag.json em swarm por nivel topologico com roteamento de modelos pela Matriz de Papeis

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Orquestrar execucao paralela de pipeline.dag.json em swarm por nivel topologico com roteamento de modelos pela Matriz de Papeis
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
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-14_SPEC_SWARM.md`](../specs/2026-11-14_SPEC_SWARM.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 0% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*

## Refinamento 1.1.0

* Data: 2026-11-16
* Versão: 1.0.0 para 1.1.0
* Etapa: e01 ()
* Peça: legoloop-loop-map
* Título: Validar o DAG e calcular os niveis topologicos
* Casos: legoloop validate pipeline.dag.json deve sair com codigo 0, agrupar os nos por nivel topologico e gravar a lista de cada nivel em out/01_niveis.md

## Refinamento 1.2.0

* Data: 2026-11-16
* Versão: 1.1.0 para 1.2.0
* Etapa: e03 ()
* Peça: legoloop-swarm-cli
* Título: Subir o processo real de cada worker
* Casos: legoloop swarm spawn -dir REPO -name NOME -- COMANDO ARGS (o processo roda no worktree do worker), o worker herda LEGOLOOP_SWARM_BOARD e LEGOLOOP_SWARM_WORKER; swarm post e swarm wait dispensam -dir e -from, gravar o PID de cada worker em out/03_spawn.md; dois ou mais workers do mesmo nivel rodam em paralelo

## Refinamento 1.3.0

* Data: 2026-11-16
* Versão: 1.2.0 para 1.3.0
* Etapa: e04 ()
* Peça: legoloop-swarm-cli
* Título: Coordenar o nivel pelo quadro em disco
* Casos: legoloop swarm post -dir REPO -to @NOME -body MENSAGEM e legoloop swarm wait -target @coordenador, legoloop swarm status -dir REPO e legoloop swarm monitor -timeout 60s para achar worker parado, o proximo nivel so comeca quando todos os workers do nivel atual estiverem stopped; registrar em out/04_coordenacao.md

## Refinamento 1.4.0

* Data: 2026-11-16
* Versão: 1.3.0 para 1.4.0
* Etapa: e05 ()
* Peça: legoloop-swarm-cli
* Título: Encerrar os workers e conferir a limpeza
* Casos: legoloop swarm kill -dir REPO -repo REPO -name NOME para cada worker ainda vivo, conferir status killed, PID 0 e worktree removido do disco e de git worktree list, gravar a conferencia em out/05_encerramento.md

## Refinamento 1.5.0

* Data: 2026-11-16
* Versão: 1.4.0 para 1.5.0
* Etapa: e05 ()
* Peça: legoloop-swarm-cli
* Título: Encerrar os workers e conferir a limpeza
* Casos: legoloop swarm kill -dir REPO -repo REPO -name NOME para cada worker ainda vivo, conferir status killed, PID 0 e worktree removido do disco e de git worktree list, gravar a conferencia em out/05_encerramento.md
