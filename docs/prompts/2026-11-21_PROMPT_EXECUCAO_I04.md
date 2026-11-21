# Prompt de ativação: execução da I04, Espinha dorsal previsível, universal e determinística

Cole o texto abaixo na primeira mensagem do agente que vai executar o plano, em uma sessão nova, na raiz do repositório.

---

Você vai executar o plano **I04, Espinha dorsal** neste repositório. A I04 executa **antes** da I01, da I02 e da I03, que continuam **congeladas**: não execute nenhum Ticket delas e não remova o aviso `EXECUÇÃO CONGELADA` dos prompts delas até o usuário aprovar o Stop-Gate final da I04.

**Leia, nesta ordem, antes de agir:**
1. `AGENTS.md` (versão enxuta das leis; o texto integral está em `docs/governanca/AGENTS_COMPLETO.md`).
2. `docs/intents/2026-11-21_INTENT_INICIATIVA_4_I04_DO_MACROPROJE.md`, `docs/intents/2026-11-21_ESPELHO_CANONICO_I04.md` e `docs/intents/2026-11-21_EMENDA_2_MACROPLANO.md`.
3. `docs/reports/2026-11-21_REPORT_AUDITORIA_ESPINHA_DORSAL.md` (as lacunas G01 a G19, com as correções da seção 4A).
4. `docs/specs/2026-11-21_SPEC_INICIATIVA_4_I04_DO_MACROPROJE.md` (60 critérios; aprovada no Stop-Gate 02) e `docs/specs/2026-11-21_APROVACAO_SPEC_I04.md`.
5. `docs/plans/2026-11-21_PLAN_INICIATIVA_4_I04_DO_MACROPROJE.md` (56 Tickets, 9 Épicos). Espelhos executáveis: `docs/routes/2026-11-21_ROUTE_INICIATIVA_4_I04_DO_MACROPROJE.json` e `docs/dags/2026-11-21_DAG_INICIATIVA_4_I04_DO_MACROPROJE.json`. Lista de origem: `docs/plans/2026-11-21_TICKETS_INICIATIVA_4_I04_DO_MACROPROJE.json`.

**Como executar:**
- Um Ticket por vez, na ordem do PLAN, pelo Swarm (`swarm_*`) sob as seis condições da lei 4 do `AGENTS.md`: contexto só com os artefatos do Ticket, ferramentas listadas, um worker por Ticket, saída em disco com Gate mecânico, sem sub-workers e `write_boundary` `strict`.
- TDD estrito: o teste falha antes do código (Red), depois Green. Mostre o Red.
- Fixture com defeito sai com **exit 2** (a convenção da trava `upstream_ok`). Exit 1 é falha de execução.
- Portão de cada Ticket: `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/funclen internal`, `gofmt -l .` (saída vazia), `go run ./tools/taxlinter internal cmd loops skills docs/templates` e `git diff --exit-code go.mod`. Ao fim de todo Ticket com Go: `go build -o bin/legoloop.exe ./cmd/legoloop` e `go install ./cmd/legoloop`.
- Pare no Stop-Gate ao fim de cada Épico e peça a aprovação do usuário. Nunca marque um Stop-Gate como aprovado sozinho.
- Renomear Processo só pelo verbo `legoloop renomear` (Épico 03), com `--dry-run` antes e o backup. Nunca à mão.

**Como rodar em Linux sem sair desta máquina Windows (já testado):**
- O WSL Ubuntu tem `gcc` e o Go 1.27.0 em `~/gotc/go/bin`. Copie o `HEAD` com `git archive HEAD | tar -x -C <pasta>` e depois mande a pasta para o WSL; rode lá `go test -count=1 ./...` e `CGO_ENABLED=1 go test -race ./...`.
- Dentro do `wsl.exe ... bash -lc '...'` não referencie `$PATH`: ele traz o do Windows, com parênteses, e quebra o comando. Use `export PATH=/usr/local/bin:/usr/bin:/bin:/home/heverton/gotc/go/bin`.
- macOS não existe aqui: o CI é o juiz.

**Leis inegociáveis do usuário:**
- Toda pergunta sua vem com a sua melhor sugestão técnica e o motivo.
- Todo achado é corrigido na execução, com teste Red, e registrado como anticorpo: `legoloop anticorpos registrar <slug> --pacote P --erro E --causa C --antidoto A`. Zero débito técnico.
- Inclua o máximo de mutações possíveis (Épico 08): teste novo só vale se uma mutação o faz falhar.
- Commit e push só quando o usuário pedir. Nunca imprima nem grave segredos (o hook `secret-guard` bloqueia; não o contorne). Responda em PT-BR com acentuação correta.

**Decisões já tomadas (Stop-Gate 02):** `caminhos` absorve as constantes de `paths` e o jail vira `internal/jail`; o jail protege os verbos que recebem nome de Processo; a descoberta liga `matriz.Staleness`, `assets.Scanner` e `matriz gerar -sem-cerebro` ao `fabricar`; `matrix` sai; os 9 Processos sem colisão ganham o prefixo e `fechar` é aposentado (colide com `legoloop-fechar`); `-race` no CI; o mutador roda por pacote com sobreviventes justificados em `espinha/mutantes_justificados.json`; os Processos com `powershell` e `python` diretos passam a usar iniciador portável.

**Fatos que você deve conferir, não supor:**
- O `go test ./...` em Linux falha hoje em 9 testes de 5 pacotes (jail, lock da triagem, guarda de deleção e seis testes dependentes de ferramenta). Esses são os Tickets do Épico 07.
- O formato exato do `isolation` dos nós e o que o `legoloop validate` devolve com DAG inválido: leia o código antes de escrever o teste.
- `Orchestrix.Execute` é o único executor e despacha pelo tipo do nó: não crie um segundo executor.
- Harness e roteador MCP abertos continuam com o binário antigo: reinicie-os.
- Prefira `mcp_search` e `call_mcp_tool` com `server: "codebase-memory-mcp"` para perguntas estruturais sobre o código (lei 14).

**Comece por:** `TICK-I4-01.1` (relatório de auditoria; a correção já está feita, o Ticket só confere), depois `TICK-I4-01.2`. Ao terminar o Épico 01, apresente as evidências e pare no Stop-Gate do Épico.
