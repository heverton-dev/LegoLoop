# Stop-Gate 02: SPEC da I03 aprovada

Decisão: **APROVADA** pelo usuário em 2026-11-21 ("[1] aprovar a SPEC da I03").
Artefato: `docs/specs/2026-11-21_SPEC_INICIATIVA_3_I03_DO_MACROPROJE.md` (23 critérios de aceite, 5 Épicos).
Decisões: (1) o valor do Swarm é `native_swarm`, que já existe nos documentos, e não `swarm_worker` (corrige a INTENT); (2) `local_process` e `orchestrix_hermetic` ficam no mesmo executor das `Options`, pois há um só executor e ele já despacha pelo tipo do nó (corrige a INTENT); (3) executor ausente vale `orchestrix_hermetic`, e valor desconhecido é recusado; (4) `native_swarm` exige papel e `command` e não vale em `stop_gate`; (5) o `Ticket` do `plangen` ganha o campo opcional `executor`.
Pendências conferidas pela execução: diretório padrão do Swarm usado pelo `legoloop run` (AC-E04.2 e AC-E04.3) e código de saída do `legoloop validate` com DAG inválido (AC-E01.5).
Próxima etapa liberada: PLAN e ROUTE.
