Risco 1: **DAG inválido no swarm_dag.go**
Arquivo: `internal/board/swarm_dag.go` (linha 124)
Cenário: O DAG contém ciclos ou nós com dependências circulares, causando falha na execução topológica.
Correção: Validar a integridade do DAG antes de executar, usando `swarmRunDAG` para detectar erros de topologia.

---

Risco 2: **Worker sem nome ou configuração ausente**
Arquivo: `internal/cli/swarm_run.go` (linha 456)
Cenário: O worker não tem nome ou parâmetros de configuração (ex: `worker_name` ou `command`), resultando em execução silenciosa ou erro.
Correção: Validar que todos os workers tenham nome e comando definidos antes de iniciar o swarm.

---

Risco 3: **Timeout de execução prolongada**
Arquivo: `internal/mcp/swarm_tools.go` (função `swarmRunDAG`, linha 387)
Cenário: A execução do DAG excede o limite de tempo configurado (ex: 30s), bloqueando o swarm.
Correção: Implementar um timeout dinâmico com alerta para casos críticos, usando `time.After()` ou `context.Deadline`.

---

Risco 4: **Entrada vazia ou inválida no comando**
Arquivo: `internal/cli/swarm_run.go` (linha 234)
Cenário: O parâmetro `command` não é passado ou é uma string vazia, causando erro na execução.
Correção: Validar entrada obrigatória e exibir mensagem clara antes de processar.

---
Recomendação: **Avaliar e corrigir riscos de entrada antes de executar o swarm.**
