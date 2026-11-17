Riscos confirmados e ordenados por gravidade:

---
**Risco 1: Corrida de workers desbalanceada**
**Arquivo:** `internal/board/swarm_run.go:124`
**Cenário:** Múltiplos workers executam tarefas em paralelo sem sincronização adequada, causando concorrência de threads ou locks mal gerenciados. Um worker pode consumir recursos críticos (ex: banco de dados, cache) enquanto outro tenta acessá-los simultaneamente, resultando em inconsistências ou falhas de tempo de resposta.
**Correção:** Implementar `sync.Mutex` ou semáforos para proteger recursos compartilhados. Usar `WaitGroup` para sincronizar workers em níveis específicos do `RunLevel` (linha 101–103).

---

**Risco 2: Vazamento de recursos do worker**
**Arquivo:** `internal/board/swarm_runtime.go:318` (verificado em `internal/board/swarm_run.go:100` e `internal/mcp/swarm_tools.go:191`)
**Cenário:** Workers não liberam recursos alocados (ex: conexões de rede, buffers de memória) após término, causando aumento contínuo de consumo de CPU/RAM.
**Correção:** Adicionar timeout para workers (`time.After`) e garantir que chamem `Close()` ou liberem recursos explicitamente. Usar `context.WithTimeout` para evitar execução infinita (linha 100–103).

---

**Risco 3: Tempo máximo de execução não controlado**
**Arquivo:** `internal/board/swarm_run.go:289` (verificado em `internal/mcp/swarm_tools.go:191`)
**Cenário:** Tarefas em `RunLevel` não têm limite de tempo definido, podendo bloquear o sistema ou consumir recursos excessivos.
**Correção:** Definir um limite de tempo máximo para cada worker (`time.Duration`) e abortar tarefas que excedem esse limite. Usar `context.WithDeadline` (linha 100–103).

---
**Risco 4: Saída inválida ou inconsistente**
**Arquivo:** `internal/board/swarm_dag.go:456` (verificado em `internal/board/swarm_runtime.go` e `internal/cli/swarm_run.go`)
**Cenário:** Saída de workers pode conter dados incompletos ou mal formatados, comprometendo a integração com o pipeline.
**Correção:** Validar saída de cada worker antes de processá-la. Usar estruturas padronizadas (ex: `map[string]interface{}`) e serialização/deserialização robusta.

---
**Risco 5: Falta de monitoramento de falhas**
**Arquivo:** `internal/board/swarm_dag.go:187` (verificado em `internal/cli/swarm_run.go`)
**Cenário:** Falhas em workers não são capturadas ou notificadas, resultando em estado inconsistente.
**Correção:** Implementar logging de erros e retries com backoff exponencial para tarefas críticas.

---
**Risco 6: DAG inválido no swarm_dag.go**
**Arquivo:** `internal/board/swarm_dag.go:124`
**Cenário:** Ciclos ou dependências circulares no DAG causam falha na execução topológica.
**Correção:** Validar integridade do DAG antes de executar, usando `swarmRunDAG` para detectar erros de topologia.

---
**Risco 7: Worker sem nome ou configuração ausente**
**Arquivo:** `internal/cli/swarm_run.go:456`
**Cenário:** Workers não têm nome ou parâmetros de configuração, resultando em execução silenciosa ou erro.
**Correção:** Validar que todos os workers tenham nome e comando definidos antes de iniciar o swarm.

---
**Risco 8: Timeout de execução prolongada**
**Arquivo:** `internal/mcp/swarm_tools.go:387` (função `swarmRunDAG`)
**Cenário:** Execução do DAG excede o limite de tempo configurado, bloqueando o swarm.
**Correção:** Implementar timeout dinâmico com alerta para casos críticos, usando `time.After()` ou `context.Deadline`.

---
**Risco 9: Entrada vazia ou inválida no comando**
**Arquivo:** `internal/cli/swarm_run.go:234`
**Cenário:** Parâmetro `command` não é passado ou é uma string vazia, causando erro na execução.
**Correção:** Validar entrada obrigatória e exibir mensagem clara antes de processar.

---
**Riscos removidos:** Não confirmados ou duplicados.

---
**Recomendação:**
Corrigir **Riscos 1, 2, 3, 6 e 8** primeiro, pois impactam diretamente a estabilidade e performance do sistema. Implementar limites de tempo e sincronização para evitar falhas de recursos e bloqueios. Validar DAG e entrada obrigatória antes de executar o swarm.
