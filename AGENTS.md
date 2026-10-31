# 🤖 AGENTS.md — Regras de Operação para Agentes no LegoLoop

> Este arquivo é a Lei de Conduta para qualquer LLM ou agente operando neste repositório.

---

## 1. Leis Invioláveis (Zero Tolerância)

1. **Zero Stubs:** Nunca crie funções vazias com `TODO`, `pass`, `// implementar depois` ou mocks simulados. O código entregue deve ser real, testado e funcional.
2. **Graph-First Obrigatório (Lei #14):**
   * SEMPRE consulte o grafo semântico via `codebase-memory-mcp` (`search_code`, `get_code_snippet`, `trace_path`) ANTES de ler arquivos extensos (> 100 linhas) ou rodar buscas cegas no disco.
   * Só use `read_file` pontual com `StartLine`/`EndLine` sobre o trecho localizado pelo grafo.
3. **Economia Severa de Tokens & Rule 10 (Silent Executor):**
   * **Thinking:** Estilo telegráfico (Caveman), direto ao ponto, sem repetições.
   * **Teto de Saída:** Respostas normais $\le 300$ tokens; explicações técnicas $\le 600$ tokens.
   * **Sem Polling em Background Tasks:** Proibido rodar loops de `manage_task(status)`. Aguarde a notificação de *Reactive Wakeup* do harness.
   * Não repita artefatos ou diffs no corpo da mensagem se já estiverem salvos no disco.
4. **Orchestrix Cirúrgica por Bloco:**
   * Nunca carregue ferramentas ou MCPs desnecessários para a tarefa atual.
   * Execute cada nó em contexto estéril; leia apenas o artefato de entrada e produza o artefato de saída contratado.
5. **TDD Estrito:**
   * O teste unitário deve falhar antes de qualquer implementação de código de produção (Red-Green-Refactor).
   * Garanta que todos os testes passem com código limpo e sem warnings.
6. **Idioma:** SEMPRE utilize o idioma Português do Brasil (PT-BR) em documentações, planos, especificações e respostas.

---

## 2. Comandos & Rotas de Governança

* **Validar conformidade e testes:** `go test ./...` ou `npm test`
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all`
* **Consultar Intenções:** [docs/intents/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents)
* **Consultar a Spec Master:** [docs/specs/SPEC.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/SPEC.md)
* **Consultar o Plano de Tickets:** [docs/plans/PLAN.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/PLAN.md)
* **Consultar Relatórios:** [docs/reports/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports)
