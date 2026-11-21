---
id: INTENT-intv_1795274711307895500
status: APPROVED
date: 2026-11-21
schema_version: 2
upstream: [INTERVIEW-intv_1795274711307895500]
covers: [Q1]
---

# 🎯 INTENT: Iniciativa 3 (I03) — Despachante de executores no runner

> **Status:** APPROVED  
> **Data:** 2026-11-21  
> **Autor:** Antigravity AI + Usuário (via Entrevista Socrática Addy Osmani)  
> **Alvo:** `internal/runner`, `internal/plangen`, `internal/routedag`, `internal/board`  
> **Sessão:** intv_1795274711307895500  
> **Confiança da Entrevista:** o motor reportou 99% no turno 1 e 90% em outra sessão aberta com as decisões; nenhuma das duas concluiu de forma confiável. O Espelho abaixo foi consolidado à mão a partir das seis respostas do usuário e reconhecido por ele ("ACEITO" em cada uma).  
> **Macroprojeto:** Distribuição e Forja do LegoLoop. A INTENT da I01 e a da I02 não são alteradas. A ordem fica na Emenda 1: `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md`.  
> **Template:** TEMPLATE_INTENT.md v2.2.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** a ROUTE grava `isolation.executor` em todo nó (hoje `orchestrix_hermetic` em 22 de 29 nós da ROUTE da I02 e `local_process` em 7), mas nenhum código de produção lê esse campo. A triagem é só o roteador de entrada e o Swarm recebe o DAG inteiro por `swarm_run dag=`. Não existe como um mesmo DAG misturar nós locais, nós pelo Orchestrix e nós entregues ao Swarm, nem como recusar um executor desconhecido antes de gastar trabalho.
* **Quem Sofre:** o operador que executa DAGs e a outra sessão que executa os planos.
* **O Porquê (por que agora):** a I02 entrega `fabricar plano` e `internal/routedag`, que passam a gerar a ROUTE e o DAG. O despachante é a primeira camada que lê o campo gerado, e precisa existir antes de o planejamento escolher o Swarm de propósito.
* **Fato conferido no código:** `orchestrix_hermetic` não é o Swarm. Ele roda o nó em `git worktree` hermético (`internal/runner/isolated.go`) e o `internal/orchestrix` escolhe harness e modelo. O Swarm (`swarm_run`) executa o DAG inteiro por nível topológico, com um worker por nó em `worktree` próprio, e a lógica mora em `internal/board`.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** o `runner` executa todo nó pelo mesmo caminho; o campo `executor` é gravado e ignorado.
  * **Amanhã:** o `runner` valida todos os nós antes do nível 0 e despacha por nível: `local_process` roda no `runner`, `orchestrix_hermetic` segue pelo Orchestrix como hoje, e `swarm_worker` (valor novo) vai ao Swarm.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: despachante no `runner` que lê `isolation.executor` da ROUTE e do DAG, sem LLM e sem heurística de tamanho ou carga.
* [ ] Objetivo 2: valor novo `swarm_worker`. `local_process` roda no `runner` local, `orchestrix_hermetic` segue pelo Orchestrix como hoje, `swarm_worker` entrega o nó a um worker do Swarm sob as seis condições da lei 4.
* [ ] Objetivo 3: granularidade por nível topológico. O `runner` roda localmente os nós `local_process` e `orchestrix_hermetic` do nível e entrega os `swarm_worker` do mesmo nível ao Swarm como DAG de um nível, por `internal/board` (sem MCP). O nível só termina com todos os Gates aprovados.
* [ ] Objetivo 4: validação prévia antes do nível 0, sobre todos os nós. Erro Go tipado `ErrExecutorDesconhecido` que lista de uma vez todos os nós inválidos (nó, valor recebido, valores aceitos). `legoloop run` sai com 1 e não cria `worktree`. `subagents_allowed: "ephemeral"` com destino fora do Swarm segue o mesmo formato.
* [ ] Objetivo 5: `plangen` e `internal/routedag` aceitam `swarm_worker` e recusam qualquer outro valor desconhecido.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Heurística de tamanho ou carga na escolha do executor.
* ⛔ Mudar a triagem.
* ⛔ Alterar a INTENT da I01 ou a da I02.
* ⛔ Mudar o formato do DAG (continua JSON).
* ⛔ LLM na decisão de despacho.
* ⛔ Refatorações oportunistas fora das Peças declaradas em `Alvo`.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_task` | `internal/runner` (despachante e validação prévia) | Implementação via TDD Red-Green. |
| `PECA-02` | `unit_task` | `internal/plangen` e `internal/routedag` | Aceitar `swarm_worker` e recusar outros valores. |
| `PECA-03` | `unit_task` | `internal/board` (entrega de DAG de um nível) | Ligar o `runner` ao Swarm sem MCP. |
| `PECA-04` | `script` | `go test ./...` | Regressão mecânica e paridade dos 22 nós. |
| `PECA-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`
* **Nível de Risco & Blast Radius:** médio a alto. Mexe no `runner`, que executa todo DAG, e acrescenta um valor ao contrato da ROUTE (contrato público: pergunta antes, e esta INTENT registra a decisão).
* **Política de Privacidade & Orçamento Cognitivo:** nenhuma chamada de LLM no despacho nem na validação.
* **Ordem de execução:** a I03 começa depois que a I02 fechar, porque mexe em `internal/routedag` e na ROUTE.
* **Contrato de Evidência Factual (Definition of Done):** as quatro provas da seção 5 e o portão completo (`go test ./...`, `go vet ./...`, `gofmt -l` vazio, `funclen`, `taxlinter`, `go.mod` intacto).
* **Restrições Inegociáveis:**
  * TDD com Red primeiro.
  * Sem LLM e sem heurística na decisão.
  * Os 22 nós `orchestrix_hermetic` se comportam exatamente como hoje.
  * Validação prévia obrigatória, sem fallback para um executor padrão.
  * INTENTs da I01 e da I02 intactas.
  * Execução pelo Swarm só sob as seis condições da lei 4 do `AGENTS.md`.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** `internal/board` pode ser chamado direto pelo `runner` sem criar ciclo de importação (a SPEC confere); o Gate de um nó entregue ao Swarm é avaliado no `worktree` do worker (flag `gates` do `swarm_run`).
* **Riscos e Perguntas em Aberto:** como o `runner` herda o `worktree` e o `workdir` quando mistura executores no mesmo nível; ciclo de importação entre `runner` e `board`.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Um DAG com executor desconhecido é recusado antes do nível 0, com todos os nós inválidos listados, `legoloop run` sai com 1 e nenhum `worktree` é criado. | Teste com DAG inválido e conferência do `git worktree list` | `go test ./internal/runner/...` |
| 2 | Um DAG misto com os três executores roda por nível: os locais no `runner`, os `swarm_worker` no Swarm via `internal/board`, e o nível só termina com todos os Gates aprovados. | Teste de integração sem LLM | `go test ./internal/runner/...` |
| 3 | Paridade: os 22 nós `orchestrix_hermetic` se comportam como hoje, e o DAG da I02 valida e executa sem mudança. | Testes atuais do runner e do Orchestrix, mais `legoloop validate` no DAG da I02 | `go test ./internal/runner/... ./internal/orchestrix/...` |
| 4 | `plangen` e `internal/routedag` aceitam `swarm_worker` e recusam outros valores desconhecidos. | Testes de contrato | `go test ./internal/plangen/... ./internal/routedag/...` |
| 5 | Portão completo verde. | Os seis comandos do portão | `go test ./...`, `go vet ./...`, `gofmt -l .`, `funclen`, `taxlinter`, `git diff go.mod` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md`](../specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Gate Humano*) no chat, antes de qualquer linha de especificação ser redigida. A trava `upstream_ok` exige `status: APPROVED` no cabeçalho.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
