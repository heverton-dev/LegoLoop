# EMENDA 1 — Macroplano Distribuição e Forja do LegoLoop

Data: 2026-11-21. Status: UNDER_REVIEW (aprovação junto com a INTENT da I02).
O macroplano da I01 mora no Anexo A da INTENT da Forja (`docs/intents/2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md`, APPROVED). Essa INTENT não muda, porque os hashes da SPEC, do PLAN e da ROUTE dependem dela. Esta Emenda acrescenta iniciativas sem tocar nela.

| Iniciativa | Tema | INTENT | Depende de |
| :--- | :--- | :--- | :--- |
| I01 | Distribuição e Forja (instalação universal, Processo Forja) | `2026-11-21_INTENT_DISTRIBUI_O_E_FORJA_DO_LEGOLOO.md` (APPROVED) | nenhuma |
| I02 | Processos de especificação e de planejamento (`fabricar plano`, `internal/routedag`, triagem, subprocesso) | `2026-11-21_INTENT_INICIATIVA_2_DO_MACROPROJETO_D.md` | execução de código só depois do TICK-01.2 da I01 |
| I03 | Despachante no `runner` (escolhe runner local ou worker do Swarm pela ROUTE, sem LLM) | `2026-11-21_INTENT_INICIATIVA_3_I03_DO_MACROPROJE.md` (APPROVED) | I02 |

Regras desta Emenda:
- A trava `upstream_ok` já está implementada e vale para I02 e I03 (primeiro Épico entregue).
- Papel e planejamento da I02 podem avançar antes do TICK-01.2; código não.
- A I03 tem SPEC aprovada, PLAN de 20 Tickets, ROUTE e DAG; o valor do Swarm é `native_swarm`, que já existia (a INTENT dizia `swarm_worker`, corrigido na SPEC).
- Cada iniciativa nova passa pela concepção antes de ter SPEC, PLAN, ROUTE ou DAG.
