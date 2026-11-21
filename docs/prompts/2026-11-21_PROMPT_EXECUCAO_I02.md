> **EXECUÇÃO CONGELADA (2026-11-21):** o usuário abriu uma auditoria da espinha dorsal do LegoLoop (I04). Não execute nenhum Ticket deste plano até o usuário liberar. O relatório da auditoria pode mudar o plano.

# Prompt de ativação: execução da I02, Processos de especificação e de planejamento

Cole o texto abaixo na primeira mensagem do agente que vai executar o plano, em uma sessão nova, na raiz do repositório.

---

Você vai executar o plano **I02, Processos de especificação e de planejamento** neste repositório.

**Pré-condição, confira antes de tudo:** o `TICK-01.2` da I01 (renomear o módulo Go) precisa estar fechado e enviado. Rode `git log --oneline -20` e `head -1 go.mod`. Se o módulo ainda for `legoloop`, **pare** e avise o usuário: a I02 move pacotes e o renomear quebraria todo `import`. Depois do renomear, troque o prefixo `legoloop/internal/...` pelo novo módulo em todo caminho dos Tickets.

**Leia, nesta ordem, antes de agir:**
1. `AGENTS.md` (versão enxuta das leis; o texto integral está em `docs/governanca/AGENTS_COMPLETO.md`).
2. `docs/intents/2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md` e `docs/intents/2026-11-21_ESPELHO_CANONICO_I02.md` (INTENT aprovada no Stop-Gate 01) e `docs/intents/2026-11-21_EMENDA_1_MACROPLANO.md`.
3. `docs/specs/2026-11-21_SPEC_INICIATIVA_2_DO_MACROPROJETO_D.md` (25 critérios; aprovada no Stop-Gate 02) e `docs/specs/2026-11-21_APROVACAO_SPEC_I02.md`.
4. `docs/plans/2026-11-21_PLAN_INICIATIVA_2_DO_MACROPROJETO_D.md` (21 Tickets, 7 Épicos). Espelhos executáveis: `docs/routes/2026-11-21_ROUTE_INICIATIVA_2_DO_MACROPROJETO_D.json` e `docs/dags/2026-11-21_DAG_INICIATIVA_2_DO_MACROPROJETO_D.json`. Lista de origem: `docs/plans/2026-11-21_TICKETS_INICIATIVA_2_DO_MACROPROJETO_D.json`.

**Como executar:**
- Um Ticket por vez, na ordem do PLAN, pelo Swarm (`swarm_*`) sob as seis condições da lei 4 do `AGENTS.md`: contexto só com os artefatos do Ticket, ferramentas listadas, um worker por Ticket, saída em disco com Gate mecânico, sem sub-workers e `write_boundary` `strict`.
- TDD estrito: o teste falha antes do código (Red), depois Green. Mostre o Red.
- Portão de cada Ticket: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .` (saída vazia), `go run ./tools/taxlinter internal cmd loops skills docs/templates` e `git diff --exit-code go.mod`. Ao fim de todo Ticket com Go, `go build -o bin/legoloop.exe ./cmd/legoloop`.
- Pare no Stop-Gate ao fim de cada Épico e peça a aprovação do usuário. Nunca marque um Stop-Gate como aprovado sozinho.

**Leis inegociáveis do usuário:**
- Toda pergunta sua vem com a sua melhor sugestão técnica e o motivo.
- Todo achado é corrigido na execução, com teste Red, e registrado como anticorpo: `legoloop anticorpos registrar <slug> --pacote P --erro E --causa C --antidoto A`.
- Commit e push só quando o usuário pedir. Nunca imprima nem grave segredos. Responda em PT-BR com acentuação correta.

**Decisões já tomadas (Stop-Gate 02):** as skills `aidd-*` do Cofre global não são editadas (os Processos novos declaram os contratos); `fabricar spec` e `fabricar plano` são subcomandos de `fabricar`; `legoloop plano` é removido sem alias; `--estrito` é o padrão, com `--permissivo`.

**Fatos que você deve conferir, não supor:**
- O `legoloop compor <nome> <ref>` acrescenta o subprocesso a um Processo que já existe: crie o `legoloop-ciclo` com `legoloop fabricar ciclo` antes.
- Só o Windows foi verificado. Linux e macOS ficam "Não verificado".
- O binário do PATH precisa ser recompilado depois de mudar Go; harness e roteador MCP abertos precisam reiniciar.
- Prefira `mcp_search` e `call_mcp_tool` com `server: "codebase-memory-mcp"` para perguntas estruturais sobre o código (lei 14).

**Comece por:** `TICK-I2-01.1` (criar `internal/routedag`). Ao terminar o Épico 01, apresente as evidências e pare no Stop-Gate do Épico.
