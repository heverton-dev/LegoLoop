---
schema: anticorpos.v1
processo: legoloop-video-roteiro
total_anticorpos: 2
globais: 0
dominio: 2
regras_inviolaveis:
  - pacote: "internal/studio"
    antidoto: "Coletar todos os campos com querySelectorAll em _renderButton, renderizar MultipleChoice em _renderSelectControl e ch..."
  - pacote: "internal/mcp"
    antidoto: "Abrir navegador em sessao_criar, timeout 5min, watchdog em sessao-app.js e sessao_aguardar no mesmo turno"
---

# 🛡️ Banco de Anticorpos — legoloop-video-roteiro

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Domínio (2)

### `internal/studio`
- **Antídoto:** Coletar todos os campos com querySelectorAll em _renderButton, renderizar MultipleChoice em _renderSelectControl e chamar sessao_aguardar logo apos sessao_publicar
- **Erro Evitado:** Botao A2UI ignorava multiplos campos, MultipleChoice caia em painel generico e agente nao chamava sessao_aguardar apos sessao_publicar
- **Causa Raiz:** a2ui-bundle.js lia apenas querySelector unico em submit_answer e omitia case MultipleChoice
- **Assinatura:** `sha256:af3c0c135c34088f4ddc89985d257084d971e50bc133e31887f3e692a459fd20`

### `internal/mcp`
- **Antídoto:** Abrir navegador em sessao_criar, timeout 5min, watchdog em sessao-app.js e sessao_aguardar no mesmo turno
- **Erro Evitado:** Overlay A2UI travava apos Enviar sem sessao_aguardar
- **Causa Raiz:** sessao_criar nao abria navegador, timeout era 30s e agente encerrava turno apos sessao_publicar
- **Assinatura:** `sha256:d4d298fd98a4ee56cd2629d0e92be4659700f46de67852d3bd18f9187c28c25d`

## Regras Globais (0)

_Nenhuma regra cadastrada nesta seção._

