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
6. **Quinteto de Governança ISO de Artefatos:**
   * TODOS os artefatos do Quinteto gerados em `docs/` (`intents/`, `specs/`, `plans/`, `routes/`, `dags/`, `reports/`, `sessions/`) DEVEM obrigatoriamente utilizar o prefixo de data no formato ISO: `YYYY-MM-DD_<NOME>.<ext>`.
   * Isso garante ordenação cronológica natural e busca imediata no SO.
   * **INTENT, SPEC, PLAN, ROUTE e DAG** compõem o **Quinteto de Governança ISO** e seguem os templates canônicos v2 em `docs/templates/` (conformidade testada em `docs/templates/templates_test.go`). O `DAG.json` em `docs/dags/` é o espelho executável 1:1 da ROUTE, projetado atomicamente para `loops/<slug>/pipeline.dag.json` e com paridade SHA-256 auditada mecanicamente no fechamento (`legoloop fechar`). A Skill (`SKILL.md`) é a manifestação em prosa fiel dos nós do DAG.
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
9. **Diagramação Visual Impecável (Zero Quebras em Mermaid e Tabelas):**
   * **Proibido HTML em Mermaid:** NUNCA use `<br>`, `<br/>`, `<b>` ou qualquer tag HTML em nós ou arestas. Isso quebra o renderizador do chat e expõe o código cru.
   * **Aspas Obrigatórias:** Todo rótulo de nó deve estar entre aspas duplas: `["Label"]`, `{"Decisão"}`.
   * **Arestas Seguras:** Nunca use caminhos de arquivos ou barras (`/`) em rótulos de arestas (`-->|"Aresta Segura"|`).
   * **Tabelas Enxutas:** Alinhamento estrito (`:---`, `:---:`) e colunas concisas para evitar rolagem horizontal na janela do chat.
10. **Comunicação Canônica STE-PT (Tríade Obrigatória com Bold Limpo):**
   * Toda interação técnica deve seguir compulsoriamente a Tríade STE-PT compacta, sem blocos de código artificiais e com quebras diretas de linha (dois espaços no final, sem linhas vazias extras):
   * **Pré-Ação (antes de alterar código):** `**[S] Situação:**`, `**[T] Teoria:**`, `**[E] Execução:**` (com tópicos `· **Ação:**`, `· **Rollback:**`, `· **Decisão:** **[1] seguir** \| **[2] personalizar** \| **[0] cancelar**).
   * **Pós-Ação (após executar e testar):** `**[S] Situação:**`, `**[T] Tese:**`, `**[E] Evidência:**` (com tópicos `· **Testes:**`, `· **Build:**`, `· **Verificação:**`).
   * Proibido preâmbulos, saudações ou re-resumos de arquivos no chat.
11. **Protocolo Inviolável de Fechamento & Canonização Universal (Zero Entregas pela Metade):**
   * Nenhuma implementação de processo, peça ou pipeline é declarada concluída sem disparar a cadeia obrigatória de fechamento:
     1. **Manual Operacional:** em `docs/manuals/` (`YYYY-MM-DD_MANUAL_FUNCAO_<NOME>_LEGOLOOP.md`).
     2. **Livro Didático LegoPress:** em `docs/books/<slug>/` (`livro.json`, `partes/`, `.md`, `.pdf` compilado e auditado com `legoloop legopress auditar`).
     3. **Skill Canônica Universal:** em `skills/<nome>/SKILL.md` (distribuída universalmente para todos os harnesses via `legoloop instalar-skills`).
     4. **Loop / Grafo DAG:** em `loops/<slug>/pipeline.dag.json` (e `.loop.json`), validado por `legoloop validate`.
     5. **Recompilação e Instalação Fresca:** `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop` para manter o PATH e o runtime do usuário atualizados.
     6. **Atualização da Governança Viva:** `AGENTS.md`, `MEMORY.md`, `README.md`.
     7. **Reindexação do Grafo Semântico:** via `codebase-memory-mcp cli index_repository`.
     8. **Portões Mecânicos de Qualidade & Commit/Push:** testes 100% verdes, gofmt limpo, go vet sem warnings, funclen <= 35 linhas, commit e push rastreáveis.
     9. **Relatório de Conclusão & Telemetria Estrita:** em `docs/reports/` (`YYYY-MM-DD_REPORT_<NOME>.md`) compilado conforme `TEMPLATE_REPORT.md` com balanço auditável de tempo, recursos/tokens e custo financeiro USD/BRL orçado vs realizado.

---

## 2. Comandos & Rotas de Governança

* **Validar conformidade e testes:** `go test -count=1 ./...` e `go test -count=1 -tags e2e ./internal/e2e/...` (ou `npm test` quando houver front-end com testes JS; a lógica pura do Studio roda em `node`).
* **Fábrica & Ciclo de Vida dos 6 Verbos de Processo:**
  * **Fabricar (criação completa):** `legoloop fabricar "<nome>" [--yes] [--goal T] [--model M] [--category build|reverse|web|mobile] [--dry-run] [--json]` (Chat: `/legoloop-fabricar`; Studio: **+ Novo Processo**). **Não confundir com `legoloop canonize`** (desduplicação por SHA-256).
  * **Matriz de Peças (reputação/aptidão):** `legoloop matriz [--tier ouro|prata|bronze] [--json]` (Chat: `/legoloop-matriz`).
  * **Listar (inventário e vínculos):** `legoloop listar [--json]` (Chat: `/legoloop-listar`; Studio: **Processos**).
  * **Refinar (substituição atômica de passos):** `legoloop refinar "<nome>" [--trocar P] [--json]` (Chat: `/legoloop-refinar`).
  * **Remover (desvinculação com rollback):** `legoloop remover "<nome>" [--dry-run] [--force] [--json]` (Chat: `/legoloop-remover`).
  * **Compor (sub-loops hierárquicos):** `legoloop compor "<alvo>" "<sub-loop>" [--pos N] [--json]` (Chat: `/legoloop-compor`).
  * **Fechar (canonização determinística em 9 etapas com telemetria):** `legoloop fechar "<slug>" [--dry-run] [--skip-build] [--json]` (Chat: `/legoloop-fechar`).
  * **Video-REA (engenharia reversa e dissecação com Matriz Original e Adoção Soberana):** `legoloop video-rea [run|ingest|forensics|curate|gate|govern|book] <args>` (Chat: `/legoloop-video-rea`).
* **Instalação Canônica das Skills:** `legoloop instalar-skills [--vault DIR] [--dry-run] [--json]` (distribui as 14 skills canônicas embutidas para o cofre e cria vínculos nos harnesses).
* **Demanda Target Workspace & Matriz Original Video-REA (08/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-08_INTENT_VIDEO_BUNDLE.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-08_SPEC_VIDEO_BUNDLE.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-08_PLAN_VIDEO_BUNDLE.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-08_ROUTE_VIDEO_BUNDLE.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-08_REPORT_VIDEO_BUNDLE.md)
* **Demanda Mega-Pipeline REDC (07/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-07_INTENT_CANAL_REDC.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-07_SPEC_CANAL_REDC.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-07_PLAN_CANAL_REDC.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-07_ROUTE_CANAL_REDC.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-07_REPORT_CANAL_REDC.md)
* **Demanda Mega-Pipeline Audiovisual Soberano (07/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-07_INTENT_VIDEO_MASTER.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-07_SPEC_VIDEO_MASTER.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-07_PLAN_VIDEO_MASTER.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-07_ROUTE_VIDEO_MASTER.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-07_REPORT_VIDEO_MASTER.md)
* **Demanda Motor Universal de Fechamento (07/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-07_INTENT_FECHAR_CANONICO.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-07_PLAN_FECHAR_CANONICO.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-07_ROUTE_FECHAR_CANONICO.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-07_REPORT_FECHAR_CANONICO.md)
* **Demanda Normalização da Taxonomia de Skills (07/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-07_INTENT_NORMALIZACAO_TAXONOMIA_SKILLS.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-07_PLAN_NORMALIZACAO_TAXONOMIA_SKILLS.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-07_ROUTE_NORMALIZACAO_TAXONOMIA_SKILLS.json)
* **Governança Lazy Universal de MCPs:** `legoloop sync-harnesses [--dry-run|--apply]` (conecta Claude, MIMO, Cursor e Windsurf ao Meta-MCP Router) e `legoloop mcp-router` (servidor STDIO leve com bootstrap de apenas 2 tools).
* **Alternador de Governança STE-PT:** `legoloop stept [status|set|on|off|scan|manual] [--harness H] [--modo M] [--pausa P] [--json]` (auditoria e chaveamento em todos os harnesses).
* **Demanda Alternador Universal STE-PT (06/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-05_INTENT_STEPT_SWITCH.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-05_PLAN_STEPT_SWITCH.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-05_ROUTE_STEPT_SWITCH.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-06_REPORT_STEPT_SWITCH.md)
* **Registrar histórico de sessão:** `faz-session --id <CONVERSATION_ID> --format all` (saída em `~/Desktop/sessoes-agentes/`; o script pode ignorar `--id` e exportar o transcript ativo).
* **Consultar Templates Canônicos v2 (SSOT):** [docs/templates/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/templates)
* **Consultar Manuais Operacionais das Funções:** [docs/manuals/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/manuals)
* **Consultar Intenções:** [docs/intents/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents)
* **Consultar a Spec Master:** [docs/specs/2026-10-30_SPEC_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-10-30_SPEC_MASTER.md)
* **Consultar o Plano de Tickets:** [docs/plans/2026-10-30_PLAN_MASTER.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-10-30_PLAN_MASTER.md)
* **Demanda Governança de Rota e Nós v2.2.0 (07/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-07_INTENT_ROUTE_V2_GOVERNANCE.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-07_SPEC_ROUTE_V2_GOVERNANCE.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-07_PLAN_ROUTE_V2_GOVERNANCE.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-07_ROUTE_ROUTE_V2_GOVERNANCE.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-07_REPORT_ROUTE_V2_GOVERNANCE.md)
* **Demanda Governança Lazy Universal (06/11/2026):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-06_INTENT_GOVERNANCA_LAZY_UNIVERSAL.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-06_PLAN_GOVERNANCA_LAZY_UNIVERSAL.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-06_ROUTE_GOVERNANCA_LAZY_UNIVERSAL.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-06_REPORT_GOVERNANCA_LAZY_UNIVERSAL.md)
* **Demanda mais recente (`fabricar v2`):** [INTENT](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-04_INTENT_FABRICAR_V2.md) · [SPEC](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-04_SPEC_FABRICAR_V2.md) · [PLAN](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/plans/2026-11-04_PLAN_FABRICAR_V2.md) · [ROUTE](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-04_ROUTE_FABRICAR_V2.json) · [Relatório final](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports/2026-11-04_REPORT_FABRICAR_V2.md)
* **Consultar Mapas de Rota:** [docs/routes/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes)
* **Consultar Relatórios:** [docs/reports/](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/reports) (capturas em `docs/reports/img/`)
* **Diário de evolução e retomada de sessão:** [MEMORY.md](file:///C:/Users/trcnologia/Desktop/LegoLoop/MEMORY.md)
