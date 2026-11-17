# INTENT — Arquitetura de Fronteira AI-Driven (2026-11-14)

## Situação
O LegoLoop já possui `mcp-router`, isolamento por `git worktree` e contrato `context_in`/`context_out`.
Faltam 6 capacidades de fronteira: Lazy de Regras por tag, roteamento por `cognitive_complexity`, mutação automática com rollback atômico, destilação de memória episódica em gates, sanitização de `context_in` externo com allowlist de ferramentas e `legoloop bootstrap` único.

## Missão
Elevar o motor de execução do LegoLoop ao padrão de fronteira (DeepMind, Anthropic, OpenAI, SpaceX, DeepSeek/Qwen) nos 7 eixos auditados, preservando TDD estrito, zero stubs e Path Jail.

## Escopo
1. Eixo 1: Injeção de regras por tag/domínio do nó (`internal/orchestrix/prompt.go`).
2. Eixo 3: Roteamento dinâmico de modelo e esforço por `cognitive_complexity`.
3. Eixo 4: Mutação automática + gate binário + rollback atômico de worktree em falha.
4. Eixo 5: Destilação de erros recorrentes em gates mecânicos e limpeza da memória episódica pós-Stop-Gate.
5. Eixo 6: Sanitização de `context_in` externo + Path Jail + allowlist de ferramentas por nó.
6. Eixo 7: Comando único `legoloop bootstrap` plug-and-play.
