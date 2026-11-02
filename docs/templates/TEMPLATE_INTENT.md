# 🎯 INTENT: [Título Claro e Declarativo da Intenção]

> **Status:** [DRAFT | UNDER_REVIEW | APPROVED]  
> **Data:** YYYY-MM-DD  
> **Autor:** [Nome do Autor / Antigravity AI]  
> **Alvo:** [Módulos afetados, ex: internal/studio/, internal/dag/, web/]  
> **Sessão:** [UUID da Sessão de Trabalho]  

---

## 1. Contexto e Motivação
* **O Problema Real:** Qual é o gargalo, bug, inconsistência ou necessidade de negócio atual?
* **O Porquê:** Por que resolver isso agora e qual é o valor gerado para o desenvolvedor ou usuário final?
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O que acontece atualmente no sistema.
  * **Amanhã:** O comportamento exato esperado após a conclusão.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: ...
* [ ] Objetivo 2: ...
* [ ] Objetivo 3: ...

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto.

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar:

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD. |
| `BLOCO-02` | `subloop` | `docs/routes/...` | Dissecação ou engenharia reversa isolada (se aplicável). |
| `BLOCO-03` | `unit_task` | Código de Produção (Go/TS) | Implementação da lógica de negócio via TDD Red-Green. |
| `BLOCO-04` | `script` | `go test ./...` | Validação de regressão mecânica e integridade. |
| `BLOCO-05` | `stop_gate` | `go build ./cmd/legoloop` | Recompilação fresca obrigatória e homologação final. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** [ `manual` | `cron: "..."` | `file_change: "caminho/*"` | `loop_end: "nome-loop"` ]
  *(Define como este fluxo acordará no dia a dia quando estiver construído).*
* **Nível de Risco & Blast Radius:** [ `Baixo (workspace local)` | `Médio` | `Alto (exige workspace git_worktree e stash_changes)` ]
* **Política de Privacidade & Orçamento Cognitivo:** [ `Nuvem autorizada` | `Exclusivamente Local (Ollama)` | `Teto: Modelos Flash/Fast` ]
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** [ Qual comportamento em tela, relatório JSON ou benchmark prova que a dor foi curada ].

---

## 5. Critérios de Sucesso de Alto Nível
1. Critério 1: Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios.
2. Critério 2: Cobertura de testes unitários mantida ou expandida (TDD Red-Green).
3. Critério 3: Binário recompilado e atualizado no PATH com sucesso.

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/YYYY-MM-DD_SPEC_NOME.md`]
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat antes de qualquer linha de especificação ser redigida.
