---
schema: anticorpos.v1
processo: legoloop-ciclo
total_anticorpos: 1
globais: 0
dominio: 1
regras_inviolaveis:
  - pacote: "internal/dag"
    antidoto: "Validar encadeamento dos contratos produces/consumes entre subprocessos e integridade upstream."
---

# 🛡️ Banco de Anticorpos — legoloop-ciclo

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (1)

### `internal/dag`
- **Antídoto:** Validar encadeamento dos contratos produces/consumes entre subprocessos e integridade upstream.
- **Erro Evitado:** Ciclo composto com quebra de contrato de handoff entre etapas.
- **Causa Raiz:** Subprocessos encadeados sem paridade de tipos de artefato de saída e entrada.
- **Assinatura:** `sha256:4a5c6d7e8f90123456789abcdef0123456789abcdef0123456789abcdef01234`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._
