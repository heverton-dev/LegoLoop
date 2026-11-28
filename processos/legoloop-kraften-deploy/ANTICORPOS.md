---
schema: anticorpos.v1
processo: legoloop-kraften-deploy
total_anticorpos: 4
globais: 0
dominio: 4
regras_inviolaveis:
  - pacote: "internal/kraften"
    antidoto: "tags JSON adicionadas ao struct Gates com teste de round-trip"
  - pacote: "internal/cli"
    antidoto: "termo normalizado para Processo no uso, comentário e documentação"
  - pacote: "internal/cli"
    antidoto: "substituído por função ehEtapaKraften (switch) sem estado"
  - pacote: "internal/kraften"
    antidoto: "Emenda 2026-11-26: campo networks generico, BaseDomain/rede/servico parametrizados, completarKraften interativo e 3 v..."
---

# 🛡️ Banco de Anticorpos — legoloop-kraften-deploy

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (4)

### `internal/kraften`
- **Antídoto:** tags JSON adicionadas ao struct Gates com teste de round-trip
- **Erro Evitado:** gates.json sem tags JSON no struct
- **Causa Raiz:** struct Gates sem tags cf_old/cf_new/tls_verdict/http_code quebrava o round-trip do receipt
- **Assinatura:** `sha256:fcb2d1413efc1710100837aaa56d112bac4100c2369dda715d128690e64b7400`

### `internal/cli`
- **Antídoto:** termo normalizado para Processo no uso, comentário e documentação
- **Erro Evitado:** termo proibido esteira na ajuda
- **Causa Raiz:** texto de uso escapou da taxonomia ASD-STE100 (Lei 15)
- **Assinatura:** `sha256:c7af0a687218d3c961975d73fd7476545ca58a2eb4cb331d325d9612b9676eef`

### `internal/cli`
- **Antídoto:** substituído por função ehEtapaKraften (switch) sem estado
- **Erro Evitado:** global mutavel etapasKraften reprova espinha
- **Causa Raiz:** mapa de subcomandos declarado como var global mutável sem justificativa
- **Assinatura:** `sha256:7c63eb44750922d61b66cb2ba738dfd7c0e42e4211bd0fb08fc958dac2845a47`

### `internal/kraften`
- **Antídoto:** Emenda 2026-11-26: campo networks generico, BaseDomain/rede/servico parametrizados, completarKraften interativo e 3 vias de entrada
- **Erro Evitado:** rede e dominio mestre fixos no Processo
- **Causa Raiz:** network_conexao e server.vpsconexao.org copiados do documento de referencia viraram regra de codigo; sem coleta interativa e sem via Studio
- **Assinatura:** `sha256:676ea33f036df22563281a5a2d392848d93d9de80b490a34e5324c74f3e7cbb4`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

