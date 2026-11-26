---
schema: anticorpos.v1
processo: legoloop-homologacao
total_anticorpos: 1
globais: 0
dominio: 1
regras_inviolaveis:
  - pacote: "internal/ste"
    antidoto: "LintEvasao, stept lint --hook e gancho Stop"
---

# 🛡️ Banco de Anticorpos — legoloop-homologacao

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (1)

### `internal/ste`
- **Antídoto:** LintEvasao, stept lint --hook e gancho Stop
- **Erro Evitado:** resposta com Nao verificado ou pendencia
- **Causa Raiz:** regra mandava escrever Nao verificado
- **Assinatura:** `sha256:a241a5eda884ff50095ef4ed3e8b819a3e4b0ee2d5c30ebd86f02c0efecd4893`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

