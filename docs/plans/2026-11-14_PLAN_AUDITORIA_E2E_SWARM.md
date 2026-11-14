# PLAN: Auditoria E2E CLI, Studio e Portões (Swarm)

- **Data:** 2026-11-14
- **Slug:** `auditoria-e2e-swarm`
- **Tickets:**
  - `T1`: Prova real CLI e MCP STDIO em ambiente hermético (`cli-prover`).
  - `T2`: Prova real LegoLoop Studio + Playwright em porta dedicada (`studio-prover`).
  - `T3`: Auditoria dos 5 portões de qualidade e correção de `gofmt` (`arch-auditor` + coordenador).
  - `T4`: Fechamento do Quinteto ISO, relatório final, reindexação do grafo semântico, commit e push.
  - `T5` (Backlog Swarm): Implementar spawn real de processos e troca de modelo em runtime (`internal/board/swarm.go`), endpoint `add` no Studio e deduplicar `internal/cli/swarm_ui.go` vs `internal/studio/swarm_api.go`.
