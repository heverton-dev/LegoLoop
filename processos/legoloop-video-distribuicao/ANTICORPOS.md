---
schema: anticorpos.v1
processo: legoloop-video-distribuicao
total_anticorpos: 4
globais: 0
dominio: 4
regras_inviolaveis:
  - pacote: "internal/video/distribuicao"
    antidoto: "Visibilidade private por padrao; publico exige --confirmo-publico e autorizacao gravada; publicar enviar e dry"
  - pacote: "internal/video/distribuicao"
    antidoto: "Toda resposta de erro passa por ocultar(); a chave vira ***"
  - pacote: "internal/video/distribuicao"
    antidoto: "RotuloIA obrigatorio no pedido e na vitrine; o corpo leva containsSyntheticMedia true"
  - pacote: "internal/video/distribuicao"
    antidoto: "Adicionar compara os temas por palavras (80% ou mais) e recusa com ErrTemaRepetido"
---

# Banco de Anticorpos — legoloop-video-distribuicao

> Espelho executivo para leitura de agentes. O banco estruturado vive em `anticorpos.json`.

## Regras de Dominio (4)

### `internal/video/distribuicao`

- **Erro:** Video publicado como publico sem querer
- **Antidoto:** Visibilidade private por padrao; publico exige --confirmo-publico e autorizacao gravada; publicar enviar e dry-run sem --executar
- **Causa raiz:** o original aceitava --visibilidade public com uma flag so

### `internal/video/distribuicao`

- **Erro:** Chave do Postiz vazou em mensagem de erro
- **Antidoto:** Toda resposta de erro passa por ocultar(); a chave vira ***
- **Causa raiz:** corpo de erro do servidor pode ecoar o cabecalho

### `internal/video/distribuicao`

- **Erro:** Video enviado sem rotulo de conteudo gerado por IA
- **Antidoto:** RotuloIA obrigatorio no pedido e na vitrine; o corpo leva containsSyntheticMedia true
- **Causa raiz:** o rotulo era opcional no original

### `internal/video/distribuicao`

- **Erro:** Tema repetido ou quase igual entrou no acervo
- **Antidoto:** Adicionar compara os temas por palavras (80% ou mais) e recusa com ErrTemaRepetido
- **Causa raiz:** a regra de originalidade so existia como texto
