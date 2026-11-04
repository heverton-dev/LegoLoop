# 🤖 AGENTS.md — Regras de Operação para Agentes no LegoLoop

> Este arquivo é a Lei de Conduta para qualquer LLM ou agente operando neste repositório.

---

## 1. Leis Invioláveis (Zero Tolerância)

1. **Zero Stubs:** Nunca crie funções vazias com `TODO`, `pass`, `// implementar depois` ou mocks simulados. O código entregue deve ser real, testado e funcional.
2. **Graph-First Obrigatório (Lei #14):**
   * SEMPRE consulte o grafo semântico via `codebase-memory-mcp` (`search_code`, `get_code_snippet`, `trace_path`) ANTES de ler arquivos extensos (> 100 linhas) ou rodar buscas cegas no disco.
   * Só use `read_file` pontual com `StartLine`/`EndLine` sobre o trecho localizado pelo grafo.
   * O grafo mora em `.codebase-memory/` (artefato gerado: **nunca edite à mão**). Para atualizá-lo: `codebase-memory-mcp cli index_repository '{"repo_path":"<raiz>"}'`. Reindexe ao fim de cada demanda.
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
   * **Recompilação Obrigatória Pós-Ticket:** Ao finalizar qualquer ticket com alteração em Go, execute `go build -o bin/legoloop.exe ./cmd/legoloop` e, ao fechar o épico, `go install ./cmd/legoloop`, para que o binário no PATH permaneça sempre atualizado. Um Studio já aberto continua no binário antigo: reinicie-o.
6. **Padronização ISO de Artefatos:**
   * TODOS os artefatos gerados em `docs/` (`intents/`, `specs/`, `plans/`, `routes/`, `reports/`, `sessions/`) DEVEM obrigatoriamente utilizar o prefixo de data no formato ISO: `YYYY-MM-DD_<NOME>.<ext>`.
   * Isso garante ordenação cronológica natural e busca imediata no SO.
   * INTENT, SPEC, PLAN e ROUTE seguem os templates canônicos v2 em `docs/templates/` (conformidade testada em `docs/templates/templates_test.go`). Ao mudar um template, os motores (`interview`, `specgen`, `plangen`) e os testes mudam juntos.
7. **Idioma:** SEMPRE utilize o idioma Português do Brasil (PT-BR) em documentações, planos, especificações e respostas.
8. **Verificação de Verdade (lições da demanda `fabricar`):**
   * **Rode o sistema real antes de declarar pronto** (binário e, se houver UI, o navegador). LLM falso e fixtures só provam o que o autor imaginou: o validador de SPEC recusava a SPEC correta de um LLM real e nenhum teste unitário via isso.
   * **Teste novo só vale se uma mutação do código o faz falhar.** Mutou, não falhou ⇒ o teste está errado. Restaure o código depois de cada mutação.
   * **Portão de qualidade de cada ticket:** `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal` (nenhuma função > 35 linhas), `gofmt -l .` com saída **vazia** (não falha por exit code: confira a saída), `git diff --exit-code -- go.mod go.sum` (só stdlib). `-race` só roda com CGO: obrigatório no CI Linux.
   * **Hermético por padrão:** testes e smoke usam LLM falso, raízes em `t.TempDir()`/`scratch/` e **nunca** tocam o cofre real `~/.aidd`. Limpe o que criou.
   * **Segredos:** a chave do 9Router só vai para a URL confiável (`provider.TrustedRouterKey`), nunca em log, saída ou resposta de API. Respostas do Studio nunca trazem caminho absoluto (`~` e `<projeto>`).
   * **Stop-gates (`SG-*`) são humanos:** nunca marque um gate como aprovado sem autorização explícita do usuário. Commit e push também só quando pedidos.
   * **Processos:** nunca `taskkill /IM legoloop.exe` (mata o Studio do usuário). Encerre só o PID ou a porta que você mesmo abriu.
   * **Declare o que não foi provado.** Relatórios dizem o que ficou aberto (ex.: `-race`, fluxo com LLM real) em vez de arredondar para "pronto".

---

## 2. Comandos & Rotas de Governança

* **Validar conformidade e testes:** `go test -count=1 ./...` (ou `npm test` quando houver front-end com testes JS; a lógica pura do Studio roda em `node`).
* **Criar um processo novo (skill + loop + dag + INTENT/SPEC/PLAN/ROUTE/TEMPLATE):** `legoloop fabricar "<nome>" [--yes] [--goal T] [--model M] [--category build|reverse|web|mobile] [--dry-run] [--json]`. Também no Chat (`/legoloop-fabricar`, skill em `skills/legoloop-fabricar/`) e no Studio (**+ Novo Processo**). **Não confundir com `legoloop canonize`** (desduplicação por SHA-256).
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all` (saída em `~/Desktop/sessoes-agentes/`; o script pode ignorar `--id` e exportar o transcript ativo).
* **Consultar Templates Canônicos v2 (SSOT):** [docs/templates/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/templates)
* **Consultar Intenções:** [docs/intents/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents)
* **Consultar a Spec Master:** [docs/specs/2026-10-30_SPEC_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-10-30_SPEC_MASTER.md)
* **Consultar o Plano de Tickets:** [docs/plans/2026-10-30_PLAN_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-10-30_PLAN_MASTER.md)
* **Demanda mais recente (`fabricar`):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-03_INTENT_LEGOLOOP_FABRICAR.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-03_SPEC_LEGOLOOP_FABRICAR.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-03_PLAN_LEGOLOOP_FABRICAR.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-03_ROUTE_LEGOLOOP_FABRICAR.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-03_REPORT_LEGOLOOP_FABRICAR.md)
* **Consultar Mapas de Rota:** [docs/routes/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes)
* **Consultar Relatórios:** [docs/reports/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports) (capturas em `docs/reports/img/`)
* **Diário de evolução e retomada de sessão:** [MEMORY.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/MEMORY.md)
