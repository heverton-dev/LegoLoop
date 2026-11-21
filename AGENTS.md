# AGENTS.md — Regras de Operação para Agentes no LegoLoop

> Lei de conduta enxuta (lazy). O texto integral, com todos os detalhes, está em `docs/governanca/AGENTS_COMPLETO.md`: leia a seção só quando a tarefa a exigir.

## 1. Leis Invioláveis

1. **Zero Stubs:** nunca crie função vazia, `TODO` ou mock simulado. O código entregue é real, testado e funcional.
2. **Graph-First (Lei #14):** consulte o grafo (`codebase-memory-mcp`: `search_code`, `get_code_snippet`, `trace_path`) antes de ler arquivos longos. O grafo mora em `.codebase-memory/` e nunca se edita à mão. Via agente lazy: `mcp_search` ("trace path", "search graph") e `call_mcp_tool` com `server: "codebase-memory-mcp"`; o binário fica em `~/tools/codebase-memory-mcp/`. Reindexe ao fim de cada demanda: `codebase-memory-mcp cli index_repository '{"repo_path":"<raiz>"}'`.
3. **Economia de tokens:** thinking telegráfico; saída normal ≤ 300 tokens; sem polling de tarefas em background; não repita artefato já salvo em disco.
4. **Orchestrix Cirúrgica por Bloco & Banimento de Subagentes:**
   * Nunca carregue ferramentas ou MCPs desnecessários para a tarefa atual.
   * Proibido usar agente de harness (`invoke_subagent`, fork, workflow livre) que clona o system prompt inteiro. Toda execução pós-plano segue o mapa formal em `docs/routes/` com execução estéril por bloco.
   * Exceção: o subagente efêmero (`isolation.subagents_allowed: "ephemeral"`) é permitido quando cumpre as seis condições:
     1. O contexto contém só os artefatos do ticket.
     2. O ticket lista as ferramentas permitidas.
     3. Um subagente executa um ticket.
     4. A saída fica em disco e passa pelo gate mecânico.
     5. O subagente não cria sub-subagentes.
     6. O `write_boundary` é `strict`.
5. **TDD Estrito & Recompilação Fresca:** o teste falha antes do código (Red-Green-Refactor). Ao fechar ticket com Go: `go build -o bin/legoloop.exe ./cmd/legoloop`; ao fechar épico: `go install ./cmd/legoloop`.
6. **Quinteto ISO:** artefatos em `docs/` (`intents/`, `specs/`, `plans/`, `routes/`, `dags/`, `reports/`, `sessions/`) usam o prefixo `YYYY-MM-DD_<NOME>.<ext>` e os templates v2 de `docs/templates/`.
7. **Idioma:** PT-BR em documentação, planos, especificações e respostas.
8. **Verdade:** rode o sistema real antes de declarar pronto; teste novo só vale se a mutação o faz falhar. Portão de cada ticket: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal` (função ≤ 35 linhas), `gofmt -l .` com saída vazia, `git diff --exit-code -- go.mod go.sum`. Testes são herméticos e nunca tocam `~/.aidd`. Stop-Gates são humanos; commit e push só quando pedidos. Nunca `taskkill /IM legoloop.exe`. Declare o que não foi provado.
9. **Mermaid:** sem HTML nos nós, rótulos entre aspas duplas, sem `/` em rótulos de aresta.
10. **STE-PT:** fonte única `internal/ste` (`legoloop stept manual`). Sem prova, escreva "Não verificado". Stop-Gate: `[1] seguir | [2] personalizar | [0] cancelar`.
11. **Fechamento (nada pela metade):** manual em `docs/manuals/`, livro em `docs/books/<slug>/`, skill em `skills/<nome>/SKILL.md`, DAG em `processos/<slug>/`, build e install frescos, governança viva, reindexação do grafo, gates verdes, commit/push e relatório em `docs/reports/`. Execute com `legoloop fechar "<slug>"`.
12. **Taxonomia ASD-STE100 (Lei #15):** um termo, um significado. Termos canônicos: Processo, DAG, Peça, Skill, Nó, Ticket, Etapa, Cofre, Matriz, Inventário, Gate, Stop-Gate, Fabricar, Canonizar, Fechar. SSOT: `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json`; gate: `go run ./tools/taxlinter internal cmd loops skills docs/templates`. Exceção de marca: o nome do projeto "LegoLoop" fica mantido.
13. **Migração, legado e lexicais (Leis 10 e 15):** detalhes em `docs/governanca/AGENTS_COMPLETO.md` (itens 13 a 15).
16. **Achado corrigido no ato e governança viva (lei inegociável):** todo achado se corrige na execução, com teste Red antes, e se registra como anticorpo (`legoloop anticorpos registrar <slug> --pacote P --erro E --causa C --antidoto A`). Toda implementação atualiza `AGENTS.md`, `README.md`, `MEMORY.md` e o grafo no mesmo ciclo. `routing` declara preferência, nunca exigência.

17. **Pergunta com sugestão (lei inegociável):** toda pergunta ao usuário traz a melhor sugestão técnica e o motivo.

18. **Cadeia amarrada (trava `upstream_ok`):** SPEC, PLAN, ROUTE e DAG só nascem de uma INTENT aprovada (`status: APPROVED` no cabeçalho), e a INTENT só nasce da concepção. `internal/upstream` aplica a trava em `legoloop plano`, `specgen`, `plangen.Generate`, `route2dag` e no Studio (exit 2 na CLI, HTTP 422). Sem escape.

## 2. Prisões Mecânicas e Protocolos

Path Jail, Schema Validation, Exit-Code Gate, Sterile Context e Test Mutation valem em runtime. Os sete protocolos industriais (STANAG 2014, ASD-STE100, EARS, FMEA, ICAO, Poka-Yoke, SBAR) estão no texto integral.

## 3. Comandos e rotas

* **Testes:** `go test -count=1 ./...` e `go test -count=1 -tags e2e ./internal/e2e/...`.
* **Processos:** `legoloop fabricar|matriz|listar|refinar|remover|compor|fechar` (Chat: `/legoloop-<verbo>`); não confunda `fabricar` com `canonize`.
* **Skills lazy (regra permanente: todo harness do SO é lazy):** `legoloop instalar-skills` (só `legoloop-skills` vai às pastas dos harnesses; o resto vem de `skill_search`/`call_skill`), `legoloop skills-custo` (custo por harness) e `legoloop skills-migrar (--dry-run|--apply)` (tira skills eager; rode após instalar plugin de terceiros). `sync-harnesses` cobre pi e OMP.
* **Claude Code:** conectores claude.ai desligados; `permissions.deny` tira 8 ferramentas pesadas; lançadores `claude-minimo`, `claude-gateway [--sem-memoria]` (sem `--chrome`), `claude-completo`, `claude-conectores`. Medida real: `claude -p "ok" --output-format json` com `MSYS_NO_PATHCONV=1`. Relatório: `docs/reports/2026-11-21_REPORT_LAZY_CONTEXTO_TOTAL.md`.
* **Demanda Distribuição e Forja do LegoLoop (21/11/2026):** INTENT aprovada em [docs/intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md](docs/intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md) (Stop-Gate 01 em `2026-11-21_INTENT_APROVADA_FORJA.md`); Espelho e ordem de implementação em [ESPELHO](docs/intents/2026-11-21_ESPELHO_CANONICO_FORJA.md). Escopo total: instalação universal por todos os canais, Processo Forja (camada 2, quatro gatilhos), execução pelo Swarm sob as seis condições da lei 4. SPEC aprovada (Stop-Gate 02); PLAN com 40 Tickets (I01, 11 Épicos), ROUTE e DAG gerados em `docs/plans`, `docs/routes` e `docs/dags` (`2026-11-21_*_DISTRIBUICAO_FORJA_LEGOLOOP`); Stop-Gate 03 aprovado. Prompt de ativação para outra sessão: `docs/prompts/2026-11-21_PROMPT_EXECUCAO_FORJA.md`. Gerar PLAN e ROUTE: `legoloop plano --spec S --tickets T.json [--estrito]`.
* **Demanda I02 — Processos de especificação e de planejamento (21/11/2026):** INTENT aprovada em [docs/intents/2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md](docs/intents/2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md); Emenda 1 do macroplano em `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md` (ordem I01, I02, I03). SPEC aprovada (25 critérios, 7 Épicos), PLAN com 21 Tickets, ROUTE e DAG em `docs/plans`, `docs/routes` e `docs/dags` (`2026-11-21_*_INICIATIVA_2_DO_MACROPROJETO_D`). O código só começa depois do TICK-01.2 da I01. Prompt de ativação: `docs/prompts/2026-11-21_PROMPT_EXECUCAO_I02.md`.
* **Demanda I03 — Despachante de executores no runner (21/11/2026):** INTENT aprovada (Stop-Gate 01) em [docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md](docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md); `runner` lê `isolation.executor`, valida todos os nós antes do nível 0 e despacha por nível (`local_process`, `orchestrix_hermetic`, `swarm_worker` novo). Depende da I02. SPEC aprovada (23 critérios, 5 Épicos), PLAN com 20 Tickets, ROUTE e DAG em `docs/plans`, `docs/routes` e `docs/dags` (`2026-11-21_*_INICIATIVA_3_I03_DO_MACROPROJE`); o valor do Swarm é `native_swarm` (já existia), não `swarm_worker`. Prompt: `docs/prompts/2026-11-21_PROMPT_EXECUCAO_I03.md`.
* **Demanda I04 — Espinha dorsal previsível, universal e determinística (21/11/2026) — mais recente:** INTENT aprovada (Stop-Gate 01) em [docs/intents/2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md](docs/intents/2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md); insumo: `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md` (14 lacunas). Executa ANTES de I01, I02 e I03, que estão **congeladas** (Emenda 2 em `docs/intents/2026-11-21_EMENDA_2_MACROPLANO.md`). Não execute Ticket de I01 a I03 sem o usuário liberar.
* **Demais comandos, demandas e manuais:** `docs/governanca/AGENTS_COMPLETO.md` (seção 3), `docs/DEMANDAS_ANTERIORES.md`, `MEMORY.md`, `docs/reports/`, `docs/manuals/`.
* **STE-PT:** `legoloop stept [status|set|on|off|scan|manual]`. **Registro de Processos:** `docs/PROCESSOS_CANONIZADOS.md` (não acrescente entradas aqui; o agy trunca este arquivo acima de 24000 bytes).
