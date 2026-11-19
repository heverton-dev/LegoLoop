---
name: legoloop-swarm
description: Executa o Processo legoloop-swarm pelo comando legoloop exec.
---

# legoloop-swarm

Execute o Processo legoloop-swarm com o comando:

    legoloop exec legoloop-swarm

## Regras

- Respeitar a topologia do enxame declarada no DAG.
- Proibido polling contínuo sem reactive wakeup.
- Cada agente opera em contexto estéril e com fronteira de escrita estrita.
