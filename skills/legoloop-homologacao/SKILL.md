---
name: legoloop-homologacao
description: "Use quando o usuário solicitar auditoria profunda, remediação, verificação cega ou fechamento canônico do ecossistema LegoLoop em Go, CLI, MCPs, Swarm e Studio Web."
---

# LegoLoop Homologação Integral do Ecossistema

Esta skill executa o processo de auditoria, correção, verificação cega e fechamento do repositório `LegoLoop`.
Manual: `docs/manuals/2026-11-14_MANUAL_FUNCAO_HOMOLOGACAO_ECOSISTEMA_LEGOLOOP.md`.
Livro: `docs/books/homologacao-ecosistema/livro.json`.
DAG: `docs/dags/2026-11-14_DAG_HOMOLOGACAO_ECOSISTEMA.json`.

## 1. Regras inegociáveis

1. Execute cada nó de verificação em contexto isolado, sem herdar o histórico do autor da correção.
2. Nunca invente notas, tempos, hashes SHA-256 ou contagem de tokens; declare "Não verificado" quando faltar prova mecânica.
3. Respeite o Stop-Gate humano entre a conclusão do planejamento e o início das edições de código (Etapa 4).
4. Bloqueie qualquer entrega que viole `taxlinter`, `funclen` (teto 35 linhas), `gofmt -l .` ou `go test -count=1 ./...`.
5. Rode o binário `legoloop` real e inspecione todas as páginas do Studio em navegador real antes de aprovar a UI.
6. Nunca mate processos do usuário por nome global (`taskkill /IM legoloop.exe`); encerre só o PID aberto pelo teste.

## 2. Sequência de Etapas e Nós

1. **Etapa 0 — Preparação e Inventário Real:** valide saúde do `9router`, reindexe o grafo em `codebase-memory-mcp`, verifique `go.mod` sem dependências externas e gere o diagnóstico inicial nas 3 frentes.
2. **Etapa 1 — Quinteto de Governança (`TICK-01.1` a `TICK-01.3`):** congele `INTENT`, `SPEC`, `PLAN`, `ROUTE` e `DAG` com assinaturas SHA-256 encadeadas e `status: DRAFT`.
3. **Etapa 2/3 — Stop-Gate Humano:** apresente o mapa de remediação e aguarde autorização explícita antes de abrir frentes de edição.
4. **Etapa 4 — Execução em Swarm (`TICK-01.1` a `TICK-06.2`):** despache trabalhadores com escrita isolada (`share`) nas 3 frentes:
   - Frente 1 (Arquitetura, Go e Taxonomia): Path Jail, CORS/CSRF, SSE sem data race, remoção de stubs de mídia e paridade JSON Schema v2.
   - Frente 2 (Canais, CLI, Harnesses, MCP e Swarm): executor real `mcp` STDIO, flags CLI, contrato único e file-lock no Swarm Board, `legoloop fechar` sem stubs.
   - Frente 3 (Studio Web e Swarm UI): zero XSS, estado preservado no polling, `#btn-run` e `ErrSemMatriz` reais, combos `t-*-128k`, layout responsivo, contraste WCAG AA e ARIA.
5. **Etapa 5 — Verificação Cega (`TICK-07.1`):** lance auditores cegos distintos dos autores + homologação ao vivo no Studio (`Playwright`, `Chrome DevTools`, `REA`, `Obscura`); exija nota mínima `9.5/10` por frente e zero erros no console.
6. **Etapa 6 — Fechamento Canônico Universal (`TICK-07.2`):** atualize Manual, Livro `legopress` compilado e auditado, esta Skill via `legoloop instalar-skills`, DAG validado com `legoloop validate`, binário via `go install ./cmd/legoloop`, grafo semântico e Relatório final.

## 3. Gates Mecânicos de Aceitação

```powershell
go test -count=1 ./...
go test -count=1 -tags e2e ./internal/e2e/...
go vet ./...
go run ./tools/funclen internal cmd tools
go run ./tools/taxlinter internal cmd skills docs/templates
gofmt -l .
git diff --exit-code -- go.mod go.sum
```

Todos os comandos acima devem encerrar com código `0` e `gofmt -l .` deve retornar saída vazia antes da entrega.
