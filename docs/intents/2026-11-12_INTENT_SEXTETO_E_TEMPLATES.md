---
id: INTENT-I2
status: APPROVED
date: 2026-11-12
schema_version: 2
upstream: [INTERVIEW-I2]
covers: [Q1, Q2, Q3, Q4, Q5, Q6, Q7, Q8, Q9]
---

# 🎯 INTENT: Sexteto e templates agnósticos (Iniciativa I2)

> **Status:** APPROVED (entrevista com 95% de confiança, 2026-11-12)  
> **Data:** 2026-11-12  
> **Autor:** Claude (orquestrador), a pedido do usuário  
> **Iniciativa:** I2 do Programa Lei #15 (I1 Taxonomia, I2 Sexteto e templates, I3 Roteamento honesto, I4 DISSECAR, I5 MATRIZ_PAPEIS, I6 Migração dos legados)  
> **Alvo:** `docs/templates`, `internal/dag`, `internal/artefato` (novo), `internal/capacidade` (novo), `internal/orchestrix`, `tools/route2dag`  
> **Entrevista:** [INTERVIEW](../interviews/2026-11-12_INTERVIEW_SEXTETO_TEMPLATES.md)  
> **Confiança da Entrevista:** 95%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

* **O Problema Real:** O sexteto INTERVIEW, INTENT, SPEC, PLAN, ROUTE e DAG é raso e preso a um fornecedor.
  * O `TEMPLATE_DAG.json` não tem `context_in`, `context_out`, `write_boundary`, `limits`, `nature`, `rollback_strategy`, `tooling`, `verification`, `spec_ref`, `status` nem `name`.
  * No topo faltam `execution_strategy`, `upstream_signatures` e `date`.
  * Não existem `TEMPLATE_INTERVIEW` nem `TEMPLATE_HANDOFF`.
  * O `TEMPLATE_INTENT` usa o termo proibido "Blocos e Subloops".
  * Os templates de hook e de script servem só ao Claude Code e ao PowerShell.
  * O `TEMPLATE_PIECE_MCP.json` ignora o roteador sob demanda (`internal/mcp/router.go`).
* **Quem Sofre:** O usuário que gera Processos em outro SO, outro harness ou outro modo de execução. O motor, que aceita DAG sem os campos que ele precisa.
* **O Porquê:** O DISSECAR, o MATRIZ_PAPEIS e toda iniciativa futura nascem dos templates. Template raso produz iniciativa rasa.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O roteamento da decisão ao DAG se perde. O motor colapsa `code-fast` e `code-pro` em Sonnet. Nada liga uma resposta da entrevista a um nó do DAG.
  * **Amanhã:** Cada artefato do sexteto tem front matter e IDs rastreáveis. O motor valida a cadeia. O roteamento nasce de uma classe de capacidade abstrata e se resolve contra o inventário do SO do usuário.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Front matter e rastreabilidade validada pelo motor (Q, R, AC, nó, `spec_ref`).
* [ ] Objetivo 2: Contrato `schema_version: 2` com modo legado para arquivos sem o campo.
* [ ] Objetivo 3: Classe de capacidade (nano a extreme) calculada por regra, inventário do SO e resolução com degradação registrada.
* [ ] Objetivo 4: Os oito templates de documento (INTERVIEW, INTENT, SPEC, PLAN, ROUTE, DAG, REPORT, HANDOFF) completos e agnósticos.
* [ ] Objetivo 5: REPORT factual gerado pelo motor; HANDOFF como recorte do REPORT.
* [ ] Objetivo 6: Peças skill, MCP, hook, script, gate e Processo agnósticas em SO, harness e modo de execução (`cli`, `agente`, `estudio`), com `load: lazy` por padrão.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos de `Alvo`.
* ⛔ Não migraremos os arquivos legados (iniciativa I6).
* ⛔ Não implementaremos executores de outros harnesses no `unit_task` (iniciativa I3).
* ⛔ Não executaremos o DISSECAR nem a MATRIZ_PAPEIS (iniciativas I4 e I5).
* ⛔ Não tocaremos o Cofre real `~/.aidd`.
* ⛔ Não criaremos dependência nova no `go.mod`.

---

## 3. Mapeamento Preliminar de Peças

| Identificador | Tipo da Peça | Especialidade | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `unit_task` | Go (`internal/artefato`) | Front matter e validador de rastreabilidade. |
| `PECA-02` | `unit_task` | Go (`internal/dag`) | Contrato `schema_version: 2` e modo legado. |
| `PECA-03` | `unit_task` | Go (`internal/capacidade`) | Classe, inventário e resolvedor. |
| `PECA-04` | `unit_task` | Markdown e JSON | Templates do sexteto, REPORT, HANDOFF e peças. |
| `PECA-05` | `unit_task` | Go (`internal/orchestrix`) | Gerador factual do REPORT. |
| `PECA-06` | `stop_gate` | Humano | Aprovar a mudança do contrato público do DAG e o fechamento. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual`.
* **Nível de Risco & Blast Radius:** Médio. A mudança de contrato do DAG afeta todos os Processos. O modo legado limita o efeito. Um Stop-Gate humano precede a mudança.
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem autorizada. Custo máximo por nó: US$ 1,50. Nenhum segredo envolvido.
* **Contrato de Evidência Factual:** `go test -count=1 ./...` verde. Os oito templates passam em `docs/templates/templates_test.go`. Um DAG gerado a partir de `TEMPLATE_DAG.json` passa em `legoloop validate`. O DAG da Lei #15 continua válido em modo legado.
* **Premissas Assumidas (a confirmar no Stop-Gate):**
  1. O parser do DAG aceita campos novos sem quebrar os arquivos existentes.
  2. O inventário do SO lê harness e modelos sem rede.
  3. O carregamento sob demanda de skills pode ser descrito por campos do front matter (a auditoria do `TICK-I2-04.1` confirma).
* **Riscos e Perguntas em Aberto:**
  1. O nome `schema_version` colide com os `$schema` v1 e v2 existentes. A SPEC fixa a regra.
  2. O formato do inventário de modelos varia por harness. A SPEC fixa o contrato mínimo.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | O DAG v2 carrega todos os campos exigidos e o motor rejeita o que falta. | Teste de validação com mutação | `go test ./internal/dag/...` |
| 2 | A cadeia Q, R, AC, nó e `spec_ref` é validada. | Teste com cadeia quebrada plantada | `go test ./internal/artefato/...` |
| 3 | O gerador resolve a classe contra um inventário simulado de cada harness e degrada com registro. | Teste com inventário parcial e vazio | `go test ./internal/capacidade/...` |
| 4 | Todo template é agnóstico em SO, harness e modo. | Teste de templates | `go test ./docs/templates/...` |
| 5 | O REPORT factual nasce da execução sem edição de agente. | Teste do gerador | `go test ./internal/orchestrix/...` |
| 6 | Os templates novos não aumentam as violações do linter (linha de base: 35 em `docs/templates`, 2026-11-12). | Linter | `go run ./tools/ticketgate -tax docs/templates -tax-max 35` |
