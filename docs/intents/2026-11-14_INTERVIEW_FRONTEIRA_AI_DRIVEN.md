# INTERVIEW — Entrevista Socrática Arquitetura de Fronteira (2026-11-14)

## Decisões Aprovadas pelo Usuário
- **Eixo 1 (Lazy de Regras):** Injetar regras por tag/domínio do nó (`go`, `ui`, `seguranca`) além do núcleo mínimo.
- **Eixo 2 (DAG e Contexto Estéril):** Manter isolamento por baia (`git worktree` + `context_in`/`context_out`).
- **Eixo 3 (Roteamento Cognitivo & Test-Time Compute):** Roteamento automático por `cognitive_complexity` do nó (baixo=rápido, alto=fronteira + thinking).
- **Eixo 4 (Oráculos & Self-Healing):** Mutação automática + gate binário + rollback atômico do worktree em falha (estilo SpaceX).
- **Eixo 5 (Memória e Destilação):** Destilação automática de erros recorrentes em gates mecânicos no fechamento e limpeza da memória episódica após o Stop-Gate final.
- **Eixo 6 (Segurança Agêntica):** Sanitização de `context_in` externo + Path Jail estrito + allowlist de ferramentas por nó.
- **Eixo 7 (Empacotamento Plug-and-Play):** Comando único `legoloop bootstrap` (binário, skills, `mcp-router`, hooks e validação).
