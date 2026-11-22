# Prompt de ativação: execução do plano Distribuição e Forja do LegoLoop

Cole o texto abaixo na primeira mensagem do agente que vai executar o plano, em uma sessão nova, na raiz do repositório.

---

Você vai executar o plano **Distribuição e Forja do LegoLoop** neste repositório, a partir do Épico 01.

**Leia, nesta ordem, antes de agir:**
1. `AGENTS.md` (versão enxuta das leis; o texto integral está em `docs/governanca/AGENTS_COMPLETO.md`).
2. `docs/intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md` e `docs/intents/2026-11-21_ESPELHO_CANONICO_FORJA.md` (INTENT aprovada no Stop-Gate 01).
3. `docs/specs/2026-11-21_SPEC_DISTRIBUICAO_FORJA_LEGOLOOP.md` (39 critérios de aceite; aprovada no Stop-Gate 02).
4. `docs/plans/2026-11-21_PLAN_DISTRIBUICAO_FORJA_LEGOLOOP.md` (40 Tickets: 12 nano e 28 micro, 11 Épicos). Espelhos executáveis: `docs/routes/2026-11-21_ROUTE_DISTRIBUICAO_FORJA_LEGOLOOP.json` e `docs/dags/2026-11-21_DAG_DISTRIBUICAO_FORJA_LEGOLOOP.json`. A lista de Tickets de origem está em `docs/plans/2026-11-21_TICKETS_DISTRIBUICAO_FORJA_LEGOLOOP.json`.

**Como executar:**
- Execute um Ticket por vez, na ordem do PLAN, pelo Swarm (`swarm_*`) sob as seis condições da lei 4 do `AGENTS.md`: contexto só com os artefatos do Ticket, ferramentas listadas, um worker por Ticket, saída em disco com Gate mecânico, sem sub-workers e `write_boundary` `strict`.
- TDD estrito: o teste falha antes do código (Red), depois Green. Mostre o Red.
- O Gate de cada Ticket é o comando de verificação do PLAN. Ao fim de cada Ticket com Go, rode `go build -o bin/legoloop.exe ./cmd/legoloop`.
- Portão de cada Ticket: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .` (saída vazia), `go run ./tools/taxlinter internal cmd loops skills docs/templates` e `git diff --exit-code -- go.mod go.sum`.
- Pare no Stop-Gate ao fim de cada Épico (`SG-E01` a `SG-E11`) e peça a aprovação do usuário. Nunca marque um Stop-Gate como aprovado sozinho.

**Leis inegociáveis do usuário:**
- Toda pergunta sua vem com a sua melhor sugestão técnica e o motivo.
- Todo achado é corrigido na execução, com teste Red, e registrado como anticorpo: `legoloop anticorpos registrar legoloop-concepcao --pacote P --erro E --causa C --antidoto A`.
- Commit e push só quando o usuário pedir. Nunca imprima nem grave segredos. Responda em PT-BR com acentuação correta.

**Decisões já tomadas (Stop-Gate 02):** origem do binário em GitHub Releases de `heverton-dev/LegoLoop`; módulo renomeado para `github.com/heverton-dev/LegoLoop` (Épico 01); `cosign` sem chave agora e Authenticode do MSI depois; WiX Toolset v4 para o MSI por usuário; Python como recurso opcional declarado no `legoloop doctor`.

**Fatos que você deve conferir, não supor:**
- Só o Windows foi verificado. Linux e macOS ficam "Não verificado" até o CI do Épico 02 rodar neles.
- O binário `legoloop` do PATH precisa ser recompilado depois de mudar Go. Um roteador MCP ou harness aberto continua com o binário antigo: reinicie-o.
- O grafo semântico fica em `~/tools/codebase-memory-mcp` (atalho no PATH em `~/.local/bin`). Prefira `mcp_search` e `call_mcp_tool` com `server: "codebase-memory-mcp"` para perguntas estruturais sobre o código (lei 14).

**Comece por:** `TICK-01.1` (conferir a SPEC contra os Gates), depois `TICK-01.2` (renomear o módulo) e `TICK-01.3` (linha de base). Ao terminar o Épico 01, apresente as evidências e pare em `SG-E01`.
