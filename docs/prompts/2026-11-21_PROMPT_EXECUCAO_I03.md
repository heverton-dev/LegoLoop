> **EXECUÇÃO CONGELADA (2026-11-21):** o usuário abriu uma auditoria da espinha dorsal do LegoLoop (I04). Não execute nenhum Ticket deste plano até o usuário liberar. O relatório da auditoria pode mudar o plano.

# Prompt de ativação: execução da I03, Despachante de executores no runner

Cole o texto abaixo na primeira mensagem do agente que vai executar o plano, em uma sessão nova, na raiz do repositório.

---

Você vai executar o plano **I03, Despachante de executores no runner** neste repositório.

**Pré-condição, confira antes de tudo:** a I02 precisa estar fechada. Rode `git log --oneline -30`, `ls internal/routedag` e `legoloop fabricar plano 2>&1 | head -3`. Se `internal/routedag` não existir ou `fabricar plano` não responder, **pare** e avise o usuário: a I03 altera a projeção do DAG e o `plangen` sobre a base da I02. Se o módulo Go foi renomeado (TICK-01.2 da I01), troque o prefixo `legoloop/internal/...` pelo novo módulo em todo caminho dos Tickets.

**Leia, nesta ordem, antes de agir:**
1. `AGENTS.md` (versão enxuta das leis; o texto integral está em `docs/governanca/AGENTS_COMPLETO.md`).
2. `docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md` e `docs/intents/2026-11-21_ESPELHO_CANONICO_I03.md` (INTENT aprovada no Stop-Gate 01) e `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md`.
3. `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md` (23 critérios; aprovada no Stop-Gate 02) e `docs/specs/2026-11-21_APROVACAO_SPEC_I03.md`.
4. `docs/plans/2026-11-21_PLAN_INICIATIVA_3_I03_DO_MACROPROJE.md` (20 Tickets, 5 Épicos). Espelhos executáveis: `docs/routes/2026-11-21_ROUTE_INICIATIVA_3_I03_DO_MACROPROJE.json` e `docs/dags/2026-11-21_DAG_INICIATIVA_3_I03_DO_MACROPROJE.json`. Lista de origem: `docs/plans/2026-11-21_TICKETS_INICIATIVA_3_I03_DO_MACROPROJE.json`.

**Como executar:**
- Um Ticket por vez, na ordem do PLAN, pelo Swarm (`swarm_*`) sob as seis condições da lei 4 do `AGENTS.md`: contexto só com os artefatos do Ticket, ferramentas listadas, um worker por Ticket, saída em disco com Gate mecânico, sem sub-workers e `write_boundary` `strict`.
- TDD estrito: o teste falha antes do código (Red), depois Green. Mostre o Red.
- Portão de cada Ticket: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .` (saída vazia), `go run ./tools/taxlinter internal cmd loops skills docs/templates` e `git diff --exit-code go.mod`. Ao fim de todo Ticket com Go, `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`.
- Pare no Stop-Gate ao fim de cada Épico e peça a aprovação do usuário. Nunca marque um Stop-Gate como aprovado sozinho.

**Leis inegociáveis do usuário:**
- Toda pergunta sua vem com a sua melhor sugestão técnica e o motivo.
- Todo achado é corrigido na execução, com teste Red, e registrado como anticorpo: `legoloop anticorpos registrar <slug> --pacote P --erro E --causa C --antidoto A`.
- Commit e push só quando o usuário pedir. Nunca imprima nem grave segredos. Responda em PT-BR com acentuação correta.

**Decisões já tomadas (Stop-Gate 02):** o valor do Swarm é `native_swarm` (não `swarm_worker`); `local_process` e `orchestrix_hermetic` ficam no mesmo executor das `Options`; executor ausente vale `orchestrix_hermetic` e valor desconhecido é recusado; `native_swarm` exige papel e `command` e não vale em `stop_gate`; `subagents_allowed: "ephemeral"` só com `native_swarm`; o `Ticket` do `plangen` ganha o campo `executor`.

**Fatos que você deve conferir, não supor:**
- O diretório padrão do Swarm usado por `legoloop run` e o código de saída de `legoloop validate` com DAG inválido não foram verificados (AC-E04.2, AC-E04.3 e AC-E01.5): leia o código e ajuste o Ticket, registrando o achado.
- `Orchestrix.Execute` é o único executor e despacha pelo tipo do nó: não crie um segundo executor.
- Só o Windows foi verificado. Linux e macOS ficam "Não verificado".
- Harness e roteador MCP abertos continuam com o binário antigo: reinicie-os.
- Prefira `mcp_search` e `call_mcp_tool` com `server: "codebase-memory-mcp"` para perguntas estruturais sobre o código (lei 14).

**Comece por:** `TICK-I3-01.1` (vocabulário de executores e `dag.ValidarExecutores`). Ao terminar o Épico 01, apresente as evidências e pare no Stop-Gate do Épico.
