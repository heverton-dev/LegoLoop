Risco 1: **DAG inválido na entrada**
**Arquivo:** `internal/board/swarm_dag.go:linha 124`
**Cenário de falha:** O DAG de entrada contém ciclos ou nós com dependências circulares, causando falha na execução topológica do `swarmRunDAG`.
**Correção:** Validar a integridade do DAG antes da execução, usando `swarm_dag.Validate()` para detectar ciclos ou inconsistências.

---

Risco 2: **Nome de worker duplicado ou ausente**
**Arquivo:** `internal/cli/swarm_run.go:linha 289`
**Cenário de falha:** Mais de um worker compartilha o mesmo nome, ou nenhum worker é definido, causando conflitos ou execução vazia.
**Correção:** Garantir que cada worker tenha um nome único e não nulo, usando `worker.Name` para validação prévia.

---

Risco 3: **Timeout de execução dos workers**
**Arquivo:** `internal/mcp/swarm_tools.go:linha 356`
**Cenário de falha:** Um worker fica bloqueado por mais tempo do que o limite configurado, causando timeout e interrupção do swarm.
**Correção:** Configurar um timeout dinâmico por worker (ex: 30s) e implementar um mecanismo de recuperação (`swarm_run.go` ou `swarm_dag.go`).

---

Risco 4: **Comando ausente ou inválido para um worker**
**Arquivo:** `internal/cli/swarm_run.go:linha 187`
**Cenário de falha:** Um worker não recebe um comando válido (ex: `nil` ou formato incorreto), levando a execução silenciosa ou erro.
**Correção:** Validar que cada worker tenha um comando não-nulo e bem-formado antes de executar.

---

Risco 5: **Timeout global do swarm_run**
**Arquivo:** `internal/cli/swarm_run.go:linha 421`
**Cenário de falha:** O tempo total para executar todos os workers excede o limite configurado (ex: 1h), causando falha no processo.
**Correção:** Adicionar um timeout global (ex: 60m) e interromper o swarm se ultrapassado, com notificação clara.

---
Recomendação: **Validar entrada e configurações antes de executar o swarm_run**, priorizando DAG, workers e timeouts.
