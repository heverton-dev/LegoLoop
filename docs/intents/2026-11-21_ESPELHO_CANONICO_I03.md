# ESPELHO CANÔNICO — Iniciativa 3 (I03)

Sessão: `intv_1795274711307895500`. Data: 2026-11-21. Reconhecido pelo usuário ("ACEITO" em cada uma das seis decisões).
O Espelho do motor repetia uma suposição errada (`orchestrix_hermetic` igual a Swarm), inventava restrições e trazia contagens erradas ("sai com 0", "7 nós restantes"). Este registro é a versão consolidada e corrigida.

1. **resultado_esperado:** despachante no `runner` que lê `isolation.executor`, valida todos os nós antes do nível 0 e despacha por nível. `local_process` roda no `runner`, `orchestrix_hermetic` segue pelo Orchestrix, `swarm_worker` (valor novo) vai ao Swarm por `internal/board`. Executor desconhecido é recusado com `ErrExecutorDesconhecido`.
2. **quem_se_beneficia:** o operador que executa DAGs e a outra sessão que executa os planos.
3. **por_que_agora:** a I02 entrega `fabricar plano` e `internal/routedag`, e hoje nenhum código lê `isolation.executor`.
4. **criterio_de_sucesso:** recusa prévia (exit 1, zero `worktree`); DAG misto com os três executores; paridade dos 22 nós `orchestrix_hermetic`; `plangen` e `internal/routedag` aceitam `swarm_worker` e recusam o resto; portão completo.
5. **restricoes_inegociaveis:** sem LLM e sem heurística; validação prévia sem fallback; Orchestrix intacto; INTENTs da I01 e da I02 intactas; Swarm só sob as seis condições da lei 4; ordem depois da I02.
6. **fora_de_escopo:** heurística de tamanho ou carga; mudar a triagem; alterar as INTENTs da I01 e da I02; mudar o formato do DAG; LLM na decisão.

Decisões do usuário: (1) nenhum critério novo, o despachante só lê o campo; (2) valor novo `swarm_worker`; (3) granularidade por nível; (4) validação prévia com erro tipado e todos os nós listados; (5) provas de pronto e limites de escopo.
