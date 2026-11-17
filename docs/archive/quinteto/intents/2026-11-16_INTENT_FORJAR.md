---
id: INTENT-FORJAR
status: DRAFT
date: 2026-11-16
schema_version: 2
upstream: [INTERVIEW-FORJAR]
covers: [Q1, Q2, Q3, Q4, Q5]
---

# 🎯 INTENT: Forjar — Endurecimento de Épicos com AgentVerse e Swarm

> **Status:** DRAFT
> **Data:** 2026-11-16
> **Autor:** sessão b622804b (entrevista socrática interview-me)
> **Iniciativa:** I01 — Agnosticidade Universal dos Verbos de Processo
> **Alvo:** internal/forjar (novo), internal/swarm, internal/fechar, docs/templates
> **Sessão:** b622804b-17d9-41d1-a5b6-5c03c65150c1
> **Confiança da Entrevista:** 95% (convergência confirmada por respostas "CORRETO" em 4 rodadas)
> **Template:** TEMPLATE_INTENT.md v2.0.0

---

## 1. Contexto e Motivação

* **O Problema Real:** Épicos viram código sem endurecimento. A especificação não passa por atrito adversarial antes de virar DAG. Resultado: retrabalho por spec fraca e câmara de eco na concepção. Evidência de desenho: o `fabricar` atual cria processo direto, sem fase de debate; nenhum verbo do cofre exercita revisão crítica com personas conflitantes.
* **Quem Sofre:** O operador que dispara épico grande (processo novo, pipeline novo) e recebe artefato executável baseado em intenção não estressada.
* **O Porquê:** Concepção errada custa 10x mais depois do código. Resolver agora cria um verbo dedicado a endurecer épico antes de qualquer toque no disco.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** entrevista, se existir, vai direto para SPEC; não há debate adversarial; o DAG nasce sem revisão conflitante; demandas grandes e pequenas usam o mesmo caminho.
  * **Amanhã:** verbo `forjar` conduz 3 tempos de endurecimento (entrevista → debate → DAG) com personas simuladas e Stop-Gates humanos entre fases; só o consenso endurecido vira Quinteto para o Swarm executar.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Novo verbo `legoloop forjar` para épicos, com 5 fases encadeadas por Stop-Gates humanos (SG-1 a SG-5).
* [ ] Objetivo 2: AgentVerse efêmero — o LLM infere o painel de personas a partir do INTENT, o painel vive só no debate e morre ao entregar o artefato.
* [ ] Objetivo 3: Critério de endurecimento — o debate para quando ninguém acha mais fura, com disjuntor de turnos por persona como teto.
* [ ] Objetivo 4: Adversarial Critic obrigatório em toda roster, com poder de veto argumentado.
* [ ] Objetivo 5: Fronteira imutável — o DAG compilado é o único contrato que vaza para o Swarm.

### o que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ Não aplicaremos o `forjar` a demandas pequenas (tarefas, tickets avulsos).
* ⛔ Não rodaremos debate com agentes reais em janelas paralelas: personas são simuladas em contexto estéril (Lei 4).
* ⛔ Não mudaremos o contrato do verbo `fabricar` existente.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | entrevista socrática (interview-me / aidd-grill) | F0 — extrair e endurecer a intenção. |
| `PECA-02` | `skill` | AgentVerse (novo) | F1/F2 — inferir roster e conduzir debate de personas. |
| `PECA-03` | `skill` | aidd-spec | F3 — compilar o debate no Quinteto canônico. |
| `PECA-04` | `skill` | legoloop-swarm | F4 — executar o DAG em worktrees efêmeros. |
| `PECA-05` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-06` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual` (operador dispara `legoloop forjar "<nome>"` para épico).
* **Nível de Risco & Blast Radius:** `Médio` — gera artefatos em `docs/` e `processos/<slug>/.forjar/`; não toca código de produção até F4.
* **Política de Privacidade & Orçamento Cognitivo:** `Teto: modelos Flash/Fast` para o debate simulado; disjuntor de turnos por persona e teto de tokens por fase. Dados sensíveis ou segredos envolvidos: nenhum.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** `legoloop forjar --dry-run` numa demanda real mostra as 5 fases com seus artefatos em disco e para no SG correto; debate de exemplo termina com transcrição e lista de furações vazia.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  * O motor de entrevista existente cobre F0 sem peça nova.
  * O motor do Swarm cobre F4 sem peça nova; a lacuna é só o AgentVerse (F1/F2).
  * Personas simuladas num contexto estéril satisfazem a Lei 4 (Orchestrix Cirúrgica).
  * O AgentVerse é efêmero por design: `roster.json` e `debate.md` morrem após o `fechar`.
* **Riscos e Perguntas em Aberto:**
  * Q1 — Inferência de personas: heurística fixa por categoria de épico ou geração livre com validação de conflito?
  * Q2 — A regra do Adversarial Critic vive no código do verbo ou numa skill canônica?
  * Q3 — Limites do disjuntor: turnos padrão por persona e teto de tokens por fase.
  * Q4 — `forjar` emite relatório em `docs/reports/` ou só o `fechar` emite telemetria?
  * Q5 — Compatibilidade: `forjar` aceita épico com INTENT manual pré-existente?

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | `go test -count=1 ./internal/forjar/...` |
| 3 | As 5 fases do `forjar` geram os artefatos previstos e param no Stop-Gate humano, sem avanço automático. | Dry-run real mostra cada fase e trava | `legoloop forjar "<nome>" --dry-run --json` |
| 4 | O AgentVerse nasce do INTENT e morre após entregar o artefato, sem resíduo permanente. | `processos/<slug>/.forjar/` efêmero e Quinteto íntegro em `docs/` | inspeção de disco pós-execução |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-16_SPEC_FORJAR.md`](../specs/2026-11-16_SPEC_FORJAR.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (Stop-Gate) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**
*A intenção consolidada com 95% de certeza factual atende à sua expectativa?*
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
