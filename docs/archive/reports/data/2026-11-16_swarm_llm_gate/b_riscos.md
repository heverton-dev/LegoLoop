Risco 1: **DAG inválido no swarm_run.go**
arquivo: `internal/cli/swarm_run.go:124`
Cenario de falha: O DAG define dependências circulares ou nodos duplicados, causando deadlocks ou execução paralela inválida.
Correção sugerida: Validar a topologia do DAG antes de executar, usando `swarm_dag.go` para detectar ciclos ou inconsistências.

---

Risco 2: **Sem papel de entrada para peça**
arquivo: `internal/mcp/swarm_tools.go:289`
Cenario de falha: O `swarmRunDAG` não verifica se a peça tem um arquivo de papel (ex: `peca.json` ou `peca.yaml`) antes de iniciar workers.
Correção sugerida: Adicionar validação de existência do arquivo de papel no início do fluxo de execução.

---

Risco 3: **Timeout de worker não configurado**
arquivo: `internal/cli/swarm_run.go:187`
Cenario de falha: Workers executam por tempo indeterminado, podendo bloquear recursos ou causar falhas de sistema se o tempo limite não for definido.
Correção sugerida: Estabelecer um timeout mínimo (ex: 30s) para cada worker via configuração externa ou parâmetro CLI.

---

Risco 4: **Nome de worker não único**
arquivo: `internal/board/swarm_dag.go:352`
Cenario de falha: Duplicação de nomes de workers leva a conflitos de execução ou perda de estado.
Correção sugerida: Gerar nomes únicos para workers a partir de um hash do papel ou ID da peça.

---
Recomendação: Validar DAG, papel de entrada e timeout antes de executar o swarmRunDAG.
