---
id: INTENT-{{ID_DA_INICIATIVA}}
status: DRAFT
date: {{AAAA-MM-DD}}
schema_version: 2
upstream: [INTERVIEW-{{ID_DA_INICIATIVA}}]
covers: [Q1]
---

# 🎯 INTENT: {{TITULO_DA_INTENCAO}}

> **Status:** {{DRAFT | UNDER_REVIEW | APPROVED}}  
> **Data:** {{DATA}}  
> **Autor:** {{AUTOR}}  
> **Iniciativa:** {{NOME_OU_ID_DA_INICIATIVA | Opcional para demandas corporativas/multi-sistemas}}  
> **Alvo:** {{PECAS_AFETADAS}}  
> **Sessão:** {{UUID_DA_SESSAO}}  
> **Confiança da Entrevista:** {{NN}}% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

<!--
GUIA DE PREENCHIMENTO (apague este trecho ao concluir)
- Todo {{PLACEHOLDER}} deve ser substituído. O gate de fechamento do documento é:
  !has_pattern('docs/intents/<arquivo>.md', '\{\{')
- A Intent responde "por quê". Nada de contratos, schemas ou nomes de função aqui: isso é da SPEC.
- Em demandas amplas ou multi-sistemas, indique a Iniciativa agregadora (ex: I01 — Nome). Em demandas comuns, o campo é opcional.
- Nome do arquivo: docs/intents/YYYY-MM-DD_INTENT_<NOME>.md (AGENTS.md, regra 6).
-->

---

## 1. Contexto e Motivação
* **O Problema Real:** {{Qual é o gargalo, bug, inconsistência ou necessidade atual? Cite a evidência: erro, métrica, relato.}}
* **Quem Sofre:** {{Persona ou sistema afetado e com que frequência.}}
* **O Porquê:** {{Por que resolver agora e qual valor é gerado.}}
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** {{O que acontece atualmente, de forma observável.}}
  * **Amanhã:** {{O comportamento exato esperado após a conclusão, de forma observável.}}

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: {{OBJETIVO_1}}
* [ ] Objetivo 2: {{OBJETIVO_2}}
* [ ] Objetivo 3: {{OBJETIVO_3}}

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora das Peças declaradas em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.
* ⛔ {{NAO_OBJETIVO_ESPECIFICO_DESTA_INTENCAO}}

---

## 3. Mapeamento Preliminar de Peças
Identificação das Skills que o futuro Processo precisará acionar. Tipos permitidos: os tipos de nó do DAG e do ROUTE, descritos no README, seção 6.

| Identificador | Tipo da Peça | Especialidade / Skill | Papel no Processo |
| :--- | :--- | :--- | :--- |
| `PECA-01` | `skill` | {{skills/aidd-spec}} | {{Elaborar especificação formal com critérios BDD.}} |
| `PECA-02` | `unit_task` | {{Código de produção (Go/TS)}} | {{Implementação da lógica via TDD Red-Green.}} |
| `PECA-03` | `script` | `go test ./...` | Regressão mecânica e integridade (remova se não houver código Go). |
| `PECA-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5; remova se não houver alteração em Go). |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** {{`manual` | `cron: "..."` | `file_change: "caminho/*"` | `processo_fim: "nome-do-processo"`}}  
  *(Define como este Processo acordará no dia a dia quando estiver construído.)*
* **Nível de Risco & Blast Radius:** {{`Baixo (workspace local)` | `Médio` | `Alto (exige workspace git_worktree e rollback stash_changes)`}}
* **Política de Privacidade & Orçamento Cognitivo:** {{`Nuvem autorizada` | `Exclusivamente Local (Ollama)` | `Teto: modelos Flash/Fast`}}. Dados sensíveis ou segredos envolvidos: {{nenhum | descreva e indique onde ficam}}.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** {{Qual comportamento em tela, relatório JSON ou benchmark prova que a dor foi curada.}}
* **Premissas Assumidas (a confirmar no Stop-Gate 01):** {{lista curta; cada premissa deve poder ser refutada}}
* **Riscos e Perguntas em Aberto:** {{lista; o que ainda não se sabe e quem pode responder}}

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes do Processo passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo falhando antes do código | {{comando ou relatório}} |
| 3 | {{CRITERIO_ESPECIFICO}} | {{prova}} | {{comando ou artefato}} |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/{{DATA}}_SPEC_{{NOME}}.md`](../specs/{{DATA}}_SPEC_{{NOME}}.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (Stop-Gate) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com {{NN}}% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
