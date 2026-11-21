# ESPELHO CANÔNICO — Iniciativa 2 (I02)

Sessão: `intv_1795272250377131200`. Data: 2026-11-21. Reconhecido pelo usuário ("aceito").
O Espelho do motor trazia erros (cobertura ≥ 90%, `mcp_search` como validação, TICK-01.1, Emenda 1 e subprocessos fora do escopo). Este registro é a versão consolidada e corrigida a partir das respostas do usuário.

1. **resultado_esperado:** dois Processos canônicos novos. Especificação: INTENT aprovada entra, SPEC em Markdown sai. Planejamento: `fabricar plano` gera PLAN (Markdown), ROUTE (JSON) e DAG (JSON) em um comando. O `route2dag` vira o pacote `internal/routedag`. A triagem roteia os dois sem LLM. Processos compostos por `subprocesso`. As skills `aidd-spec`, `aidd-planner` e `aidd-tickets` ganham `consumes`/`produces` corretos.
2. **quem_se_beneficia:** o operador do ciclo completo e a outra sessão que executa os planos.
3. **por_que_agora:** a trava `upstream_ok` já existe, e a cadeia INTERVIEW, INTENT, SPEC, PLAN, ROUTE, DAG precisa de Processos para a triagem rotear.
4. **criterio_de_sucesso:** quatro provas mecânicas sem LLM (ponta a ponta com exit 0 e assinaturas válidas; sem INTENT aprovada exit 2 sem gravar; triagem roteia SPEC e PLAN de forma determinística; Processo composto passa no `legoloop validate`) mais o portão completo.
5. **restricoes_inegociaveis:** TDD com Red primeiro; sem LLM na decisão; nenhum validador novo de SPEC; nenhum código antes do TICK-01.2 fechar; INTENT da I01 intacta.
6. **fora_de_escopo:** alterar a INTENT da I01; mudar o formato do DAG.

Decisões do usuário na entrevista: (1) contrato do Processo de especificação é só o par INTENT para SPEC; (2) `fabricar plano` gera os três artefatos e para no primeiro erro; (3) validação reutiliza o Gate de aceite existente; (4) a I02 só executa código depois do TICK-01.2; (5) as quatro provas de aceite.
