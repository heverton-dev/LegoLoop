---
schema: anticorpos.v1
processo: legoloop-swarm
total_anticorpos: 1
globais: 0
dominio: 1
regras_inviolaveis:
  - pacote: "internal/board"
    antidoto: "semWorkerEncerrado em AddWorker"
---

# 🛡️ Banco de Anticorpos — legoloop-swarm

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (1)

### `internal/board`
- **Antídoto:** semWorkerEncerrado em AddWorker
- **Erro Evitado:** swarm run recusa reexecutar no com worker killed
- **Causa Raiz:** nome de worker encerrado continuava ocupado
- **Assinatura:** `sha256:413e46701cabb7bdb397199365d1469f8d63b45c1d4130d80c58e39bba4155e9`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

