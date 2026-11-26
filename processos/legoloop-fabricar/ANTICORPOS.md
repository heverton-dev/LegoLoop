---
schema: anticorpos.v1
processo: legoloop-fabricar
total_anticorpos: 4
globais: 0
dominio: 4
regras_inviolaveis:
  - pacote: "internal/plangen"
    antidoto: "Ticket.Command e executor native_swarm no ticket"
  - pacote: "internal/cli"
    antidoto: "apontarProgressoNoDAG e IDs TICK-<iniciativa letras+digitos>-NN.N"
  - pacote: "internal/plangen"
    antidoto: ""
  - pacote: "internal/video/inventario"
    antidoto: ""
---

# 🛡️ Banco de Anticorpos — legoloop-fabricar

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (4)

### `internal/plangen`
- **Antídoto:** Ticket.Command e executor native_swarm no ticket
- **Erro Evitado:** DAG do fabricar plano nao sobe no swarm
- **Causa Raiz:** no sem role nem command do executor
- **Assinatura:** `sha256:6bed35942086d3512a482b696df4459f0b5acaf25f3583796fffa792954e3de1`

### `internal/cli`
- **Antídoto:** apontarProgressoNoDAG e IDs TICK-<iniciativa letras+digitos>-NN.N
- **Erro Evitado:** DAG gerado sem ponteiro de progresso e ID de ticket fora da convencao
- **Causa Raiz:** fabricar plano nao gravava progresso no DAG; TICK-VA nao casa com o observador
- **Assinatura:** `sha256:c7340fe0b5b00d181016c9d61b02907943f572dab01bf769132711c7a796a803`

### `internal/plangen`
- **Antídoto:** 
- **Erro Evitado:** validation_cmd com aspas de argumento virava gate prompt e travava o swarm; regerar plano zera progresso
- **Assinatura:** `sha256:38ab666e5976380b2c4d0e6b98bf7d90f46e5a6be6d15b39bd0b1a948a9c298d`

### `internal/video/inventario`
- **Antídoto:** 
- **Erro Evitado:** vocabulario de faixa etaria limitado a 3 valores; ampliado para 5 (bebe, crianca, adolescente, adulto, idoso)
- **Assinatura:** `sha256:fcf903ed54d755f2eb9546ecf7566107bc11ed4fcea9770190fd7e3560a9ebf7`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

