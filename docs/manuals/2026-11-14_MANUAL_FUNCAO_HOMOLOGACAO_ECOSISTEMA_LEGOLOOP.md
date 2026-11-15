# MANUAL OPERACIONAL // HOMOLOGACAO INTEGRAL DO ECOSSISTEMA LEGOLOOP

> **Versão:** 1.0.0  
> **Data:** 2026-11-14  
> **Skill Canônica:** `skills/legoloop-homologacao/SKILL.md` (`legoloop-homologacao`)  
> **DAG Canônico:** `docs/dags/2026-11-14_DAG_HOMOLOGACAO_ECOSISTEMA.json`  
> **Livro Didático:** `docs/books/homologacao-ecosistema/livro.json`

---

## 1. Visão Geral e Propósito

Este manual documenta o Processo de auditoria profunda, remediação, verificação cega independente e fechamento canônico do ecossistema LegoLoop.
O Processo valida de ponta a ponta três frentes integradas:
1. **Frente 1 — Arquitetura, Go e Taxonomia ASD-STE100:** blindagem de caminhos (`os.OpenRoot`), CORS/CSRF, ausência de data race em SSE, remoção de stubs em mídia e conformidade com `schemas/route.v2.json`.
2. **Frente 2 — Canais CLI, Harnesses, MCP e Swarm:** cliente e servidor `mcp` sobre STDIO sem deadlock, flags de CLI alinhadas, contrato único e file-lock no Swarm Board, e execução real de `legoloop fechar`.
3. **Frente 3 — Studio Web e Swarm UI:** zero XSS no DOM, preservação de estado no polling, execução real via `#btn-run`, exibição correta de `ErrSemMatriz`, combos certificados `t-*-128k`, layout responsivo, contraste WCAG AA e atributos ARIA.

---

## 2. Invariantes e Gates Obrigatórios

Toda execução desta função obedece às leis invioláveis abaixo:

1. **Zero Stubs:** proibido manter funções vazias ou retornos fictícios em produção.
2. **Graph-First:** consulte o grafo semântico via `codebase-memory-mcp` antes de leituras extensas em disco.
3. **Contexto Estéril na Verificação Cega:** auditores da Etapa 5 não podem compartilhar histórico com os autores das correções.
4. **Limites Mecânicos por Ticket:** `funclen <= 35` linhas em `internal`, `cmd` e `tools`; `taxlinter` com zero violações; `gofmt -l .` vazio; `go.mod` e `go.sum` restritos à biblioteca padrão Go.
5. **Isolamento Hermético em Testes:** testes usam raízes temporárias (`t.TempDir()` ou `scratch/`) e nunca alteram o cofre real do usuário.
6. **Segredos Protegidos:** chaves de provedores nunca aparecem em logs, respostas JSON ou telas do Studio; caminhos locais são mascarados com `~`, `<projeto>` e `<cofre>`.

---

## 3. Sequência das 7 Etapas do Processo

### Etapa 0 — Preparação e Inventário Real
- Valide a saúde do gateway `9router` e selecione combos certificados (`t-*-128k`).
- Reindexe o grafo semântico em `codebase-memory-mcp`.
- Rode a suíte baseline (`go test`, `go vet`, `funclen`, `taxlinter`, `gofmt`) e gere o diagnóstico inicial com notas por frente.

### Etapa 1 — Quinteto ISO (`TICK-01.1` a `TICK-01.3`)
- Gere `INTENT`, `SPEC`, `PLAN`, `ROUTE` e `DAG` em `docs/` com prefixo ISO `2026-11-14_`.
- Encadeie as assinaturas SHA-256 de montante em `ROUTE` e `DAG`, mantendo o status inicial em `DRAFT`.

### Etapa 2 e Etapa 3 — Planejamento e Stop-Gate Humano
- Consolide a matriz de Risco FMEA e o Plano de Voo com os 15 Tickets (`TICK-01.1` a `TICK-07.2`).
- Pause obrigatoriamente no Stop-Gate humano (`[1] seguir | [2] personalizar | [0] cancelar`) antes de modificar código.

### Etapa 4 — Execução em Swarm por Frentes (`TICK-01.1` a `TICK-06.2`)
Despache trabalhadores em cópias isoladas (`share`) executando os Tickets do DAG:
- `TICK-01.1` e `TICK-01.2`: desbloqueio da suíte Go, remoção de artefatos sensíveis rastreados e paridade Schema v2.
- `TICK-01.3`: zeragem de violações de `taxlinter` e migração para combos certificados.
- `TICK-02.1` e `TICK-02.2`: Path Jail em `studio`, prevenção de race em SSE e substituição de stubs em `internal/motion` (`scripts/beat_detect.py`) e `internal/videorea`.
- `TICK-03.1` e `TICK-03.2`: cliente/roteador MCP STDIO (`limitedBuffer` sincronizado + `finishMCPError`) e paridade CLI (`swarm --help`, `sync-harnesses --json`, `fechar`).
- `TICK-04.1` e `TICK-04.2`: integração N10 Foundry e trava atômica no Swarm Board.
- `TICK-05.1` e `TICK-05.2`: Swarm UI unificado e desbloqueio E2E do Studio (`/`, `/inventario`, `/swarm`).
- `TICK-06.1` e `TICK-06.2`: layout responsivo multi-viewport, inspetor recuperável (`280px` <-> `0px`), `/favicon.ico`, contraste WCAG AA e ARIA.

### Etapa 5 — Verificação Cega Independente (`TICK-07.1`)
- Conduza auditoria cega com 3 verificadores independentes e 1 revisor de convergência.
- Execute o binário real (`./bin/legoloop.exe studio --addr 127.0.0.1:<porta>`) com perfil isolado e valide `/`, `/inventario` e `/swarm` em `1440x900`, `1024x768` e `390x844` (claro e escuro) usando Playwright, Chrome DevTools MCP, REA e Obscura.
- Critério de aprovação: nota mínima `9.5/10` em cada frente, zero erros no console do navegador e todos os gates mecânicos em `exit_code == 0`. Se reprovado, abra rodada de remediação e reavalie em contexto limpo.

### Etapa 6 — Fechamento Canônico Universal (`TICK-07.2`)
Execute a cadeia completa de fechamento:
1. Atualize este Manual em `docs/manuals/`.
2. Compile e audite o Livro Didático LegoPress em `docs/books/homologacao-ecosistema/`.
3. Atualize e distribua a Skill `legoloop-homologacao` via `legoloop instalar-skills`.
4. Atualize o Quinteto para `APPROVED`, recompute os hashes SHA-256 do DAG e valide com `legoloop validate`.
5. Recompile `./bin/legoloop.exe` e instale no PATH com `go install ./cmd/legoloop`.
6. Atualize `AGENTS.md`, `MEMORY.md` e `README.md`.
7. Reindexe o grafo semântico (`codebase-memory-mcp`).
8. Rode toda a suíte mecânica após a última edição.
9. Grave o Relatório final em `docs/reports/` apenas com medições reais observadas.

---

## 4. Portões Mecânicos e Diagnóstico Rápido

```powershell
go build -o bin/legoloop.exe ./cmd/legoloop
go install ./cmd/legoloop
go test -count=1 ./...
go test -count=1 -tags e2e ./internal/e2e/...
go vet ./...
go run ./tools/funclen internal cmd tools
go run ./tools/taxlinter internal cmd skills docs/templates
gofmt -l .
git diff --exit-code -- go.mod go.sum
```

### Problemas Conhecidos e Mitigação
- **Timeout em `GET /api/assets` na primeira varredura:** `Watcher.EnsureInitial` serializa a varredura fria com double-checked locking, `scanHarnesses` e `scanBinaries` executam em paralelo com `wg.Go`, e os handlers HTTP aplicam `extendWriteDeadline` de 3 minutos.
- **Race/EOF intermitente no executor MCP STDIO:** `limitedBuffer` protege leituras e escritas com mutex e `finishMCPError` aguarda o encerramento do subprocesso (teto de 500ms) antes de formatar `stderr`.
- **Teste `-race` no Windows:** requer CGO habilitado; quando ausente localmente, declare explicitamente como "Não verificado" e exija execução no CI Linux.
