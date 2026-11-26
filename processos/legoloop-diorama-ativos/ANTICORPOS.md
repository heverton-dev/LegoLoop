---
schema: anticorpos.v1
processo: legoloop-diorama-ativos
total_anticorpos: 4
globais: 0
dominio: 4
regras_inviolaveis:
  - pacote: "internal/a2ui,internal/studio"
    antidoto: "Alinhar Go e JS com \bon[a-z]+\s*= e exibir err.message no badge de status do sessao-app.js"
  - pacote: "internal/studio"
    antidoto: "Aplicar renderFormattedText em a2ui-bundle.js e legoloop-catalog.js e comprovar com preview HTTP 200 real"
  - pacote: "internal/a2ui,internal/studio"
    antidoto: "Persistir aguardando_resposta=true no SubmitAction, limpar no proximo Publish e desabilitar controles no reload em se..."
  - pacote: "internal/diorama/ativos"
    antidoto: "resolverEntradaSprite com fallback resiliente para folha mestre e resolverElencoPadrao alimentado por fichasGeradas"
---

# 🛡️ Banco de Anticorpos — legoloop-diorama-ativos

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (4)

### `internal/a2ui,internal/studio`
- **Antídoto:** Alinhar Go e JS com \bon[a-z]+\s*= e exibir err.message no badge de status do sessao-app.js
- **Erro Evitado:** Envelope ignorado no frontend A2UI por falso positivo de HOSTILE_RE
- **Causa Raiz:** Regex on\w+= sem fronteira de palavra rejeitava textos benignos como session_id= no JS apos passar no Go
- **Assinatura:** `sha256:d69dd96697fe333a20bd356494890cd401e98078ce0a94a5972805b5a3dbf213`

### `internal/studio`
- **Antídoto:** Aplicar renderFormattedText em a2ui-bundle.js e legoloop-catalog.js e comprovar com preview HTTP 200 real
- **Erro Evitado:** Texto markdown cru exibido no Panel e LegoStopGatePanel sem formatacao e prova visual pendente
- **Causa Raiz:** Panel e renderStopGatepainel injetavam textContent puro sem parse estruturado via DOM e preview local falhava sem Host explicito
- **Assinatura:** `sha256:52759fccf0bdb86ef91331e8d9be793121c030e481889eedac03b1305b3ecdcc`

### `internal/a2ui,internal/studio`
- **Antídoto:** Persistir aguardando_resposta=true no SubmitAction, limpar no proximo Publish e desabilitar controles no reload em sessao-app.js
- **Erro Evitado:** F5/reload na A2UI reexibia etapa anterior ativa durante processamento
- **Causa Raiz:** SubmitAction nao persistia aguardando_resposta em SessionStatus para replay no SSE
- **Assinatura:** `sha256:67fcc468a1867769a8013d63b973529110ecd47a2a3a15262062697a7dfb16e5`

### `internal/diorama/ativos`
- **Antídoto:** resolverEntradaSprite com fallback resiliente para folha mestre e resolverElencoPadrao alimentado por fichasGeradas
- **Erro Evitado:** sprite_despill falhava buscando poses_6x inexistente ao inves da folha mestre gerada
- **Causa Raiz:** resolverElencoPadrao e RecortadorRembgDespill exigiam flow/narrador_poses.png ao inves de flow/fichas/<id>_sheet.png gerada no Flow
- **Assinatura:** `sha256:d443bfc14ff89cef97d16258c26e7d2a13aea293d29e369ca65dee388bb15726`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

