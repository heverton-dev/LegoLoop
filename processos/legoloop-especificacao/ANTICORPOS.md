---
schema: anticorpos.v1
processo: legoloop-especificacao
total_anticorpos: 2
globais: 0
dominio: 2
regras_inviolaveis:
  - pacote: "internal/upstream"
    antidoto: "Pacote internal/upstream aplica a trava upstream_ok: recusa com exit code 2 se a INTENT não estiver APPROVED."
  - pacote: "internal/specgen"
    antidoto: "specgen.Generate valida a presença de Critérios Binários de Aceite com padrões EARS e ASD-STE100."
---

# 🛡️ Banco de Anticorpos — legoloop-especificacao

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (2)

### `internal/upstream`
- **Antídoto:** Pacote internal/upstream aplica a trava upstream_ok: recusa com exit code 2 se a INTENT não estiver APPROVED.
- **Erro Evitado:** SPEC gerada sem INTENT aprovada quebrando a cadeia de valor.
- **Causa Raiz:** Geradores antigos de SPEC não validavam o cabeçalho da INTENT upstream.
- **Assinatura:** `sha256:2ec34291e87fa1b8a1e2eff3a004e37b8803dc894c7a5d23bd1b199ff2e1db62`

### `internal/specgen`
- **Antídoto:** specgen.Generate valida a presença de Critérios Binários de Aceite com padrões EARS e ASD-STE100.
- **Erro Evitado:** SPEC reprovada por ausência de critérios binários de aceite.
- **Causa Raiz:** Critérios de aceite vagos ou sem formato binário inviabilizavam os gates de teste mecânico.
- **Assinatura:** `sha256:2ce32da32c24a211b6fabf6b53659c1a7532f95cc9de0a55f90995018b6cc764`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._
