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
4. **Orchestrix Cirúrgica por Bloco & Banimento de Subagentes:**
   * Nunca carregue ferramentas ou MCPs desnecessários para a tarefa atual.
   * Proibido usar subagentes de harness (`invoke_subagent`) que clonam o system prompt inteiro (queima de tokens). Toda execução pós-plano deve seguir o mapa formal em `docs/routes/` com execução estéril por bloco.
   * Execute cada nó em contexto estéril; leia apenas o artefato de entrada e produza o artefato de saída contratado.
5. **TDD Estrito & Recompilação Fresca:**
   * O teste unitário deve falhar antes de qualquer implementação de código de produção (Red-Green-Refactor).
   * Garanta que todos os testes passem com código limpo e sem warnings.
   * **Recompilação Obrigatória Pós-Ticket:** Ao finalizar qualquer ticket com alteração em Go, execute obrigatoriamente a recompilação: `go build -o bin/legoloop.exe ./cmd/legoloop` para que o binário no PATH permaneça sempre atualizado.
6. **Padronização ISO de Artefatos:**
   * TODOS os artefatos gerados em `docs/` (`intents/`, `specs/`, `plans/`, `routes/`, `reports/`, `sessions/`) DEVEM obrigatoriamente utilizar o prefixo de data no formato ISO: `YYYY-MM-DD_<NOME>.<ext>`.
   * Isso garante ordenação cronológica natural e busca imediata no SO.
7. **Idioma:** SEMPRE utilize o idioma Português do Brasil (PT-BR) em documentações, planos, especificações e respostas.

---

## 2. Comandos & Rotas de Governança

* **Validar conformidade e testes:** `go test ./...` ou `npm test`
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all`
* **Consultar Intenções:** [docs/intents/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents)
* **Consultar a Spec Master:** [docs/specs/2026-10-30_SPEC_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-10-30_SPEC_MASTER.md)
* **Consultar o Plano de Tickets:** [docs/plans/2026-10-30_PLAN_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-10-30_PLAN_MASTER.md)
* **Consultar Mapas de Rota:** [docs/routes/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes)
* **Consultar Relatórios:** [docs/reports/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports)
