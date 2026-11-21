---
id: INTENT-intv_1795272250377131200
status: APPROVED
date: 2026-11-21
schema_version: 2
upstream: [INTERVIEW-intv_1795272250377131200]
covers: [Q1]
---

# 🎯 INTENT: Iniciativa 2 (I02) — Processos de especificação e de planejamento

> **Status:** APPROVED  
> **Data:** 2026-11-21  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** `internal/routedag`, `internal/cli` (fabricar plano), `internal/triagem`, skills `aidd-spec`, `aidd-planner`, `aidd-tickets`, processos `legoloop-especificacao` e `legoloop-planejamento`  
> **Sessão:** intv_1795272250377131200  
> **Confiança da Entrevista:** 92% medida pelo motor durante os turnos; 100% reportada pelo motor após o usuário reconhecer o Espelho consolidado (turno 7). O Espelho abaixo foi corrigido à mão a partir das respostas do usuário e reconhecido por ele.  
> **Macroprojeto:** Distribuição e Forja do LegoLoop. A I01 (Forja) tem INTENT própria e não é alterada. A ordem das iniciativas fica na Emenda 1 (arquivo novo): `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md`.  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** os passos INTENT para SPEC e SPEC para PLAN, ROUTE e DAG existem só como código solto (`specgen`, `plangen`, `legoloop plano`, `tools/route2dag`). Não são Processos canônicos, os contratos `consumes`/`produces` das skills `aidd-spec`, `aidd-planner` e `aidd-tickets` estão incorretos, o `route2dag` roda por `go run` fora do binário e a triagem não conhece esses passos.
* **Quem Sofre:** o operador que roda o ciclo completo (INTERVIEW, INTENT, SPEC, PLAN, ROUTE, DAG) e a outra sessão que executa os planos gerados.
* **O Porquê (por que agora):** a trava `upstream_ok` (`internal/upstream`) já existe e barra SPEC, PLAN e DAG sem INTENT aprovada. A cadeia agora precisa de Processos canônicos para a triagem rotear cada passo sem LLM.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** `legoloop plano` gera PLAN e ROUTE; o DAG sai de `go run ./tools/route2dag`; SPEC nasce de `specgen`; a triagem não roteia nenhum desses passos; não há Processo composto por `subprocesso` cobrindo a cadeia.
  * **Amanhã:** dois Processos canônicos. Especificação: INTENT aprovada entra, SPEC em Markdown sai. Planejamento: `fabricar plano` gera PLAN (Markdown), ROUTE (JSON) e DAG (JSON) em um comando, com a trava e as assinaturas conferidas antes de gravar e sem gravar nada pela metade. O `route2dag` vira o pacote `internal/routedag`. A triagem roteia os dois sem LLM. Um Processo composto por `subprocesso` encadeia concepção, especificação e planejamento.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Processo de especificação. `consumes: Artifact:Intent` (com a trava `upstream_ok` como pré-condição) e `produces: Artifact:Spec` em Markdown, validado pelo Gate de aceite existente.
* [ ] Objetivo 2: Processo de planejamento e comando `fabricar plano` no lugar de `legoloop plano`. Confere a Referência da Intenção, a trava e o Gate de aceite existente, exige ao menos um Ticket por critério, gera PLAN, ROUTE e DAG juntos e para no primeiro erro sem gravar nada.
* [ ] Objetivo 3: mover `tools/route2dag` para `internal/routedag` (`tools/route2dag` vira `main` mínimo).
* [ ] Objetivo 4: corrigir `consumes`/`produces` das skills `aidd-spec`, `aidd-planner` e `aidd-tickets`.
* [ ] Objetivo 5: a triagem conhece os dois Processos novos e os roteia pela decisão determinística.
* [ ] Objetivo 6: Processo composto por `subprocesso` (concepção, especificação, planejamento).
* [ ] Objetivo 7: Emenda 1 do macroplano, em arquivo novo, com a ordem I01, I02, I03 e as dependências.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Alterar a INTENT da I01 (Forja). Os hashes da SPEC, do PLAN e da ROUTE dependem dela.
* ⛔ Mudar o formato do DAG (continua JSON).
* ⛔ Criar um schema novo de SPEC. O Gate de aceite existente é o único validador.
* ⛔ Refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Dependências externas novas além das já autorizadas no projeto.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | `skills/aidd-spec` | Elaborar a especificação formal com critérios BDD. |
| `PECA-02` | `skill` | `skills/aidd-planner` e `skills/aidd-tickets` | Decompor a SPEC em Tickets, PLAN e ROUTE. |
| `PECA-03` | `unit_task` | Código de produção (`internal/routedag`, `fabricar plano`, triagem) | Implementação via TDD Red-Green. |
| `PECA-04` | `script` | `go test ./...` | Regressão mecânica e integridade. |
| `PECA-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** médio. Move um pacote e muda um comando público (`legoloop plano` para `fabricar plano`). A definir no Stop-Gate 01 se o comando antigo fica como alias temporário.
* **Política de Privacidade & Orçamento Cognitivo:** nenhuma chamada de LLM na decisão de roteamento nem na geração do PLAN, da ROUTE e do DAG.
* **Ordem de execução:** a I02 só executa código depois que o TICK-01.2 (renomear o módulo Go) da I01 estiver fechado e enviado. Até lá a I02 avança só no papel (INTENT, SPEC, PLAN, ROUTE, DAG).
* **Contrato de Evidência Factual (Definition of Done):** as quatro provas mecânicas da seção 5 e o portão completo (`go test ./...`, `go vet ./...`, `gofmt -l` vazio, `funclen`, `taxlinter`, `go.mod` e `go.sum` intactos).
* **Restrições Inegociáveis:**
  * TDD com Red primeiro.
  * Sem LLM na decisão de roteamento.
  * Nenhum validador novo de SPEC.
  * Nenhum código antes do TICK-01.2 fechar.
  * INTENT da I01 intacta.
  * Execução pelo Swarm só sob as seis condições da lei 4 do `AGENTS.md`.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** o comando `fabricar plano` mantém os mesmos flags do `legoloop plano`; o `--estrito` passa a ser o padrão do `fabricar plano` (sai com 1 se faltar Ticket para um critério).
* **Riscos e Perguntas em Aberto:** o alias temporário de `legoloop plano`; a ordem exata dos Épicos (definida na SPEC).

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Ponta a ponta: de uma INTENT aprovada saem SPEC, PLAN, ROUTE e DAG em diretório temporário, com exit 0 e cadeia de assinaturas válida. | Teste de integração sem LLM | `go test ./internal/cli/... ./internal/routedag/...` |
| 2 | Sem INTENT aprovada, o comando sai com exit 2 e não grava nada. | Teste que confere o exit e a pasta vazia | `go test ./internal/cli/...` |
| 3 | A triagem roteia "gerar SPEC" e "gerar PLAN" para os dois Processos novos pela decisão determinística, sem chamar o LLM. | Teste com LLM que falha se for chamado | `go test ./internal/triagem/...` |
| 4 | O Processo composto por `subprocesso` (concepção, especificação, planejamento) passa no `legoloop validate`. | Validação do DAG do Processo | `legoloop validate processos/<processo>/pipeline.dag.json` |
| 5 | Portão completo verde. | Os seis comandos do portão | `go test ./...`, `go vet ./...`, `gofmt -l .`, `funclen`, `taxlinter`, `git diff go.mod go.sum` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md`](../specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat, antes de qualquer linha de especificação ser redigida. A trava `upstream_ok` exige `status: APPROVED` no cabeçalho.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
