---
name: legoloop-video-distribuicao
author: LegoLoop
description: Use quando precisar executar o Processo legoloop-video-distribuicao (P4 do processo composto legoloop-video-animado): revisao com parecer, vitrine, acervo e agenda quente/perene, publicacao pelo Postiz, comentarios e aprendizado semanal. NAO use para roteiro, ativos, render nem para criar identidade de canal (legoloop-canal-redc).
consumes: Artifact:RenderedVideo
produces: Artifact:PublishedVideo
artifact: acervo/acervo.json
gate: exit_code == 0
suggested_model: sonnet
---

# legoloop-video-distribuicao — Revisao, Vitrine, Agenda e Publicacao

> **Autor:** LegoLoop
> **Versao:** 1.0.0
> **Motor oficial:** `legoloop video-distribuicao` (implementacao em `internal/video/distribuicao` e `internal/cli/video_distribuicao.go`).

---

## 1. Contexto e Contratos de I/O

**Quando usar:** depois do `legoloop-video-render`, com `render.mp4` (ou `renders/vN.mp4`) pronto na pasta `projetos/video-animado/<nome-do-projeto>/`.
**Quando NAO usar:** para gerar o video, o roteiro ou os ativos; para criar a identidade do canal (use `legoloop-canal-redc`).

**Entradas:** `renders/vN.mp4`, `script.json`, `renders/legenda.srt`, fundos das capas, `vitrine.json`, `pedido_envio.json`, `comentarios.json`, `metricas.json`.
**Saidas:** `review/sheet_NN.jpg`, `review/vN_pacote.json`, `review/vN_parecer.json`, `publicacao.md`, `acervo/acervo.json`, `acervo/respostas.json`, `acervo/aprendizado.md`.

**Pre-requisitos do usuario (fora do repositorio):** Postiz instalado e canais conectados; `POSTIZ_API_KEY` e `POSTIZ_URL` no ambiente ou em arquivo `.env*`; conta profissional no Instagram; auditoria da API do YouTube e do TikTok para video publico (sem ela, so privado).

---

## 2. Nos do Processo (ver o DAG do Processo)

1. **n01_prevoo** — `legoloop video-distribuicao prevoo --projeto <DIR>`: pasta `acervo/`, ffmpeg, ffprobe, Python com Pillow, variaveis do Postiz.
2. **n02_folhas_revisao** — `revisao folhas --video renders/vN.mp4 --versao N`: folhas de contato e `review/vN_pacote.json` com video, roteiro, regras e checklist "sem cara de IA".
3. **n03_parecer** — `revisao parecer --arquivo <ARQ>`: grava o parecer do revisor sem contexto (agente isolado). Nunca sobrescreve.
4. **n04_liberar** — `revisao liberar --versao N`: so o veredito ENVIAR libera a vitrine.
5. **n05_vitrine** — `vitrine pontuar` e `vitrine gravar --arquivo vitrine.json --titulo K`: 5 titulos com nota, 3 capas por formato, descricao com capitulos e fontes, tags e rotulo de IA. Capas: `capa --arquivo <ARQ>`.
6. **n06_agenda** — `acervo adicionar`, `acervo autorizar`, `agenda --aplicar`: quente fura a fila (ate 24 h, maximo 48 h), perene segue a cadencia, estoque minimo, reels derivados.
7. **n07_publicar** — `publicar enviar --id ID --arquivo <ARQ> [--executar]`: dry-run por padrao.
8. **n08_comentarios** — `comentarios listar`, `rascunhar`, `decidir`: so resposta aprovada sai.
9. **n09_aprendizado** — `aprendizado --metricas <ARQ>`: grava `acervo/aprendizado.md`, lido pela descoberta de pautas do `legoloop-video-roteiro`.

**Stop-Gates humanos:** o usuario aprova o parecer ENVIAR (n04), a vitrine (n05), a autorizacao de publicacao no acervo (n06) e cada resposta a comentario (n08).

---

## 3. Regras Negativas (Guardrails)

- PROIBIDO publicar como `public` sem `--confirmo-publico` e sem `autorizado: true` no acervo.
- PROIBIDO enviar sem o rotulo de conteudo gerado por IA.
- PROIBIDO gravar `POSTIZ_API_KEY` ou qualquer chave em artefato, log ou chat.
- PROIBIDO responder comentario sem aprovacao do usuario.
- PROIBIDO criar contas de rede ou digitar senhas.
- PROIBIDO apagar versao de render ou parecer; cada versao fica.
- Nao recrie logica aqui: delegue sempre ao binario `legoloop`.
- Achado e corrigido na execucao e registrado em `docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md`, com anticorpo em `processos/legoloop-video-distribuicao/anticorpos.json` quando aplicavel.

## 4. Escopo declarado (o que este Processo nao cobre)

- Cortes verticais 9:16 de um master longo (BL-03, pertence ao `legoloop-video-render`).
- Coleta automatica de comentarios e de metricas pelas APIs das redes: o Processo le arquivos JSON.
- Instalar e operar o Postiz e o VPS.
- Limites de tamanho e caracteres por rede (`redes.json` do original) seguem nao verificados em pagina oficial; confira antes de depender deles.

## 5. Comandos de Verificacao

```bash
legoloop video-distribuicao --help
legoloop validate processos/legoloop-video-distribuicao/pipeline.dag.json
go test -count=1 -cover ./internal/video/distribuicao/...
legoloop fechar legoloop-video-distribuicao
```

## 6. Checklist de Conclusao

- [ ] `video-distribuicao prevoo` sem `FALTA`.
- [ ] Parecer da versao com veredito ENVIAR.
- [ ] `publicacao.md` gravado com 5 titulos avaliados e rotulo de IA.
- [ ] Item autorizado e agendado no acervo.
- [ ] Envio em dry-run conferido; envio real so com `--executar`.
