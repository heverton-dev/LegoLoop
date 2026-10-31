# 🤖 AGENTS.md — Regras de Operação para Agentes no LegoLoop

> Este arquivo é a Lei de Conduta para qualquer LLM ou agente operando neste repositório.

---

## 1. Leis Invioláveis (Zero Tolerância)

1. **Zero Stubs:** Nunca crie funções vazias com `TODO`, `pass`, `// implementar depois` ou mocks simulados. O código entregue deve ser real, testado e funcional.
2. **Economia Severa de Tokens (Caveman Thinking):**
   * Em raciocínios internos, use estilo telegráfico direto ao ponto.
   * Em respostas ao usuário, seja conciso, usando tabelas e listas em vez de parágrafos prolixos.
   * Não repita artefatos no corpo do chat se o arquivo já estiver salvo no disco.
3. **Orchestrix Cirúrgica por Bloco:**
   * Nunca carregue ferramentas ou MCPs desnecessários para a tarefa atual.
   * Execute cada nó em contexto estéril; leia apenas o artefato de entrada e produza o artefato de saída contratado.
4. **TDD Estrito:**
   * O teste unitário deve falhar antes de qualquer implementação de código de produção.
   * Garanta que todos os testes passem com código limpo e sem warnings.
5. **Idioma:** SEMPRE utilize o idioma Português do Brasil (PT-BR) em documentações, planos, especificações e respostas.

---

## 2. Comandos Operacionais

* **Validar conformidade e testes:** `go test ./...` ou `npm test`
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all`
* **Consultar a Spec Master:** [SPEC.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/SPEC.md)
* **Consultar o Plano de Tickets:** [PLAN.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/PLAN.md)
