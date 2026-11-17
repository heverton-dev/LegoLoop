# PLAN — Arquitetura de Fronteira AI-DRIVEN (2026-11-14)

## Tickets
- **T-01 (Eixo 1 — Lazy de Regras):** Adicionar campo `rule_tags` no nó do DAG e carregar regras por tag em `internal/orchestrix/prompt.go` (TDD).
- **T-02 (Eixo 3 — Roteamento Cognitivo):** Implementar roteador de modelo e esforço por `cognitive_complexity` em `internal/router/` (TDD).
- **T-03 (Eixo 4 — Mutação e Rollback):** Adicionar verificador de mutação e rollback atômico de worktree em `internal/runner/` (TDD).
- **T-04 (Eixo 5 — Destilação e Limpeza):** Implementar destilação de erros recorrentes e limpeza de memória episódica em `legoloop fechar` (TDD).
- **T-05 (Eixo 6 — Sanitização e Allowlist):** Implementar sanitizador de `context_in` e enforcement de allowlist de ferramentas por nó (TDD).
- **T-06 (Eixo 7 — Bootstrap Único):** Implementar comando `legoloop bootstrap` em `internal/cli/` (TDD) + recompilação e portões de qualidade.
