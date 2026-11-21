# INTENT_APROVADA — Stop-Gate 01 (I03)

Decisão: **APROVADA** pelo usuário em 2026-11-21 (resposta "[1] aprovar e commitar").
Artefato aprovado: `docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md` (cópia datada: `docs/intents/2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md`).
Espelho Canônico: `docs/intents/2026-11-21_ESPELHO_CANONICO_I03.md`. Sessão: `intv_1795274711307895500`.
Decisão-chave: valor novo `swarm_worker`; `orchestrix_hermetic` não é o Swarm (conferido em `internal/runner/isolated.go` e `internal/orchestrix`).
Riscos abertos para a SPEC: ciclo de importação entre `runner` e `board`; herança de `worktree` e `workdir` ao misturar executores no mesmo nível.
Achado corrigido antes da aprovação: o motor da entrevista concluía sem Espelho (`completeTurn`).
Próxima etapa liberada: SPEC. A execução de código da I03 só começa depois que a I02 fechar.
