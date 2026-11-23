---
schema: anticorpos.v1
processo: legoloop-planejamento
total_anticorpos: 2
globais: 0
dominio: 2
regras_inviolaveis:
  - pacote: "internal/upstream"
    antidoto: "Pacote internal/upstream aplica a trava VerificarSpec: recusa com exit code 2 se a SPEC ou sua INTENT upstream não estiverem APPROVED."
  - pacote: "internal/routedag"
    antidoto: "Pacote internal/routedag usa tabela única dag.TabelaTipos via dag.DoRoute e rejeita tipos desconhecidos."
---

# 🛡️ Banco de Anticorpos — legoloop-planejamento

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (2)

### `internal/upstream`
- **Antídoto:** Pacote internal/upstream aplica a trava VerificarSpec: recusa com exit code 2 se a SPEC ou sua INTENT upstream não estiverem APPROVED.
- **Erro Evitado:** Plano gerado sem SPEC aprovada quebrando a cadeia de valor upstream.
- **Causa Raiz:** Geradores antigos de plano não validavam a cadeia completa de aprovação upstream.
- **Assinatura:** `sha256:d8a9b71e3c5409e51cbb22f8721389efb6a7d903f5681691a0bbdf8009230551`

### `internal/routedag`
- **Antídoto:** Pacote internal/routedag usa tabela única dag.TabelaTipos via dag.DoRoute e rejeita tipos desconhecidos.
- **Erro Evitado:** Discrepância entre os tipos de nós na ROUTE e no DAG gerados.
- **Causa Raiz:** Conversão legada utilizava mapas paralelos ou divergentes de tipos de nó.
- **Assinatura:** `sha256:29b5a1c7f9d85429a1b9c7e09849b2f3478950d81c15f62e8484a8c3d9a10271`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._
