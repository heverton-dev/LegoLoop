Risco 1: **Corrida de workers desbalanceada**
**Arquivo:** `internal/board/swarm_run.go:124`
**Cenário:** Múltiplos workers executam tarefas em paralelo sem sincronização, causando corrida de threads ou locks mal gerenciados. Um worker pode consumir recursos críticos (ex: banco de dados, cache) enquanto outro tenta acessá-los simultaneamente, resultando em inconsistências de estado ou falhas de tempo de resposta.
**Correção:** Implementar locks ou semáforos no código de sincronização (ex: `sync.Mutex` em Go) para proteger recursos compartilhados. Usar `WaitGroup` para sincronizar a execução de workers em níveis específicos do `RunLevel`.

---

Risco 2: **Vazamento de recursos do worker**
**Arquivo:** `internal/board/swarm_runtime.go:318`
**Cenário:** Workers não liberam recursos alocados (ex: conexões de rede, buffers de memória) após término, causando aumento contínuo de consumo de CPU/RAM. Um worker bloqueado ou com tempo de execução excessivo pode esgotar recursos do sistema.
**Correção:** Adicionar um timer de timeout para workers (`time.After`) e garantir que eles chamem `Close()` ou liberem recursos explicitamente após a execução. Usar `context.WithTimeout` para evitar execução infinita.

---

Risco 3: **Tempo máximo de execução não controlado**
**Arquivo:** `internal/board/swarm_run.go:289`
**Cenário:** Tarefas em `RunLevel` não têm limite de tempo definido, podendo bloquear o sistema ou consumir recursos excessivos. Um worker que falha ou fica em loop pode prolongar o tempo total de execução.
**Correção:** Definir um limite de tempo máximo para cada worker (`time.Duration`) e abortar tarefas que excedem esse limite. Usar `context.WithDeadline` para enforçar o tempo.

---

Risco 4: **Saída inválida ou inconsistente**
**Arquivo:** `internal/board/swarm_runtime.go:456`
**Cenário:** A saída de workers pode conter dados incompletos ou mal formatados, causando falhas na integração com o pipeline ou na validação de resultados. Erros de parsing ou log de dados incorretos podem comprometer a confiabilidade.
**Correção:** Validar a saída de cada worker antes de processá-la. Usar estruturas de dados padronizadas (ex: `map[string]interface{}`) e implementar um sistema de serialização/deserialização robusto para garantir consistência.

---

Risco 5: **Falta de monitoramento de falhas**
**Arquivo:** `internal/board/swarm_run.go:187`
**Cenário:** Falhas em workers não são capturadas ou notificadas, resultando em perda de dados ou estado inconsistente. Um worker que falha pode deixar o sistema em um estado incerto.
**Correção:** Implementar um sistema de logging de erros e retries com backoff exponencial para tarefas críticas. Usar `retry/promote` para lidar com falhas temporárias e garantir que o sistema recupere o estado.

---
**Recomendação:** Auditar e corrigir os riscos 1, 2, e 3 primeiro, pois impactam diretamente a estabilidade e performance do sistema. Documentar os limites de tempo e recursos para evitar regressões futuras.
