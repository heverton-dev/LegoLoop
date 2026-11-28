# Auditoria do P4 `legoloop-video-distribuicao` (REA bit a bit do original)

Data: 2026-11-27. Alvo: `C:\Users\trcnologia\Desktop\proj--historias-palitos`. Escopo do P4 (mapa macro): Etapas 9, 10, 11, 11b, 12 e 13.

## 1. O que foi lido por inteiro (sem amostragem)

`processo/PROCESSO.md` (Etapas 8b a 13 e portoes), `PUBLICACAO.md`, `MINIATURA.md`, `KIT_DE_CANAL.md`, `CORTES.md`, `CONECTAR_REDES.md`, `CANAIS.md`, `AUDITORIA.md`, `RODAS.md`, `LACUNAS.md`, `POSTMORTEM.md`, `INTERFACE_A2UI.md`, `PROMPTS.md` (P7, P8, P10), `AGENTS.md` (linhas de publicacao), `ferramentas/publicar.py`, `miniatura.py`, `canal.py`, `sheet.py`, `redes.json`, `LEIAME.md`, `vps/postiz-stack.yml`, `vps/enviar-segredo.ps1`, `servicos/postiz/README.md`, `canal/bilhoes-em-jogo/textos.json` (formato). Fora do escopo do P4 e ja tratado em P1 a P3: autoria de cenas, voz, motor, mix.

## 2. Estado no LegoLoop hoje

Nao existe `processos/legoloop-video-distribuicao`, skill, codigo Go, nem DAG. So existe o nome antigo no mapa (corrigido no ticket 1). Existem pecas reutilizaveis: `legoloop-canal-redc` (kit de canal), Looping Studio A2UI (`sessao_*`), `video-render` (render e folhas), inventario por estilo.

## 3. Tabela: feature original -> destino

| # | Feature original (arquivo) | Destino | Situacao |
|---|---|---|---|
| 1 | Folhas de contato a cada 1 a 2 s (`sheet.py`, 24 quadros por folha, 6x4, carimbo de tempo) | `video-distribuicao revisao` em Go + ffmpeg | migrar |
| 2 | Autorrevisao com checklist "sem cara de IA" (PROCESSO Etapa 9, CANAIS item 6) | checklist no nó de revisão, itens fixos | migrar |
| 3 | Comparacao com referencia (rea-video, PySceneDetect) | reutilizar `legoloop-video-rea` | adiar: so quando houver video de referencia |
| 4 | Revisor sem contexto (P7): recebe so video, folhas, regras; responde ENVIAR ou NAO ENVIAR + lista `v<N> . corte . mm:ss . problema . correcao` | no isolado `orchestrix_hermetic` que grava `review/vN_parecer.json` | migrar |
| 5 | Notas por versao v1, v2, v3 sem apagar versoes (P8) | `renders/vN.mp4` + `review/vN_parecer.json` | migrar |
| 6 | Portao humano "enviar" (A2UI `revisao`) | Stop-Gate humano via Studio | migrar |
| 7 | 5 titulos com nota (palavra principal no inicio, ate 60 caracteres, sem mentira) | nó de vitrine | migrar |
| 8 | 3 capas por formato; capa do curto = momento forte + ate 3 palavras, colada nos primeiros quadros | nó de vitrine; Go + Pillow equivalente (image/draw) | migrar |
| 9 | Miniatura: fundo do Flow sem texto + texto e logo pelo codigo (`miniatura.py` 1280x720, degrade, canto inferior direito livre, aviso > 2 MB) | `video-distribuicao capa` em Go | migrar |
| 10 | Limite 1280x720 e 2 MB "nao verificado" (MINIATURA) | verificar na pagina oficial antes de fixar; ate la, aviso | verificar na execucao |
| 11 | Descricao com capitulos e fontes, tags, rotulo de IA, `.srt` | nó de vitrine | migrar |
| 12 | Aprovacao do pacote (A2UI `vitrine`) -> video entra no acervo | Stop-Gate + gravacao em `acervo/acervo.json` | migrar |
| 13 | `acervo/acervo.json` (id, tema, tipo quente/perene, formatos, canais, validade, prioridade, status pronto/agendado/publicado/expirado, data, links) | pacote Go `internal/video/distribuicao` | migrar |
| 14 | Regra quente: fura fila, sai em ate 24 h (max 48 h), vencido vira expirado | agendador deterministico | migrar |
| 15 | Regra perene: cadencia fixa, estoque minimo, alerta | agendador deterministico | migrar |
| 16 | `ferramentas/agenda.py` (A CRIAR no original; nunca existiu) | agendador em Go | migrar (nasce aqui) |
| 17 | `acervo/temas.json` antirrepeticao | ja no P1 (`video-roteiro`) | ja migrado |
| 18 | Regra de originalidade (tema novo, pesquisa propria, sem dezenas de videos quase iguais) | gate do agendador | migrar |
| 19 | `publicar.py`: canais, subir, enviar; padrao private; publico exige `--confirmo-publico`; dry-run; chave nunca impressa; titulo <= 100; tags; `--quando` | `video-distribuicao publicar` em Go (net/http) | migrar |
| 20 | Chave `POSTIZ_API_KEY` e `POSTIZ_URL` lidas do ambiente ou `.env*`; nunca no chat | mesma regra; ler do ambiente; zero segredo em artefato | migrar |
| 21 | Rotulo IA `containsSyntheticMedia` no YouTube | campo obrigatorio no envio | migrar |
| 22 | Autorizacao unica de publicacao gravada no acervo, revogavel antes do horario | campo `autorizacao` no acervo | migrar |
| 23 | Distribuicao cruzada: todo curto publico no YouTube gera Reel no Instagram | regra do agendador | migrar |
| 24 | Espacar cortes: 1 por dia ou a cada 2 dias (recomendacao) | cadencia configuravel | migrar |
| 25 | Etapa 12: coletar comentarios sem resposta (YouTube, Instagram) e escrever rascunhos no tom do canal; nada sai sem aprovacao; `acervo/respostas.json` | nó de comentarios com Stop-Gate | migrar (rascunho); coleta real depende de API |
| 26 | Etapa 13: relatorio semanal (retencao por tema, formato, duracao; horarios; cortes que trouxeram seguidores) -> `acervo/aprendizado.md`, lido pela Etapa 1 | nó de aprendizado | migrar; leitura pelo P1 ja prevista |
| 27 | Kit de canal (`canal.py`, `redes.json`, `textos.json`, `colar.md`, imagens por rede, previas 98 px e area segura) | reutilizar `legoloop-canal-redc`; nao duplicar | descartar daqui (ja tem processo) |
| 28 | Limites por rede em `redes.json` (quase todos "nao verificado") | verificar em pagina oficial na execucao | verificar |
| 29 | Postiz no VPS (Docker Swarm, Traefik, stack YAML, `enviar-segredo.ps1`, site-legal, `CONECTAR_REDES.md`) | infraestrutura do usuario, fora do repo | adiar: documentar como pre-requisito |
| 30 | Auditoria de API do YouTube (videos privados ate auditar), TikTok `SELF_ONLY` ate auditar, Instagram exige conta profissional e analise do app | pre-requisitos declarados na skill | documentar |
| 31 | Cortes 9:16 (Etapa 8b, `cortes.py`) | P3 `video-render`, BL-03 | adiado com motivo |
| 32 | Telas A2UI `revisao`, `vitrine`, `agenda`, `comentarios`, `cortes` | Looping Studio (`sessao_*`), uma tela por portao | migrar |
| 33 | Painel local `ferramentas/painel/` (nunca existiu) | descartar: o Studio ja existe | descartado |
| 34 | Hook pre-commit com gitleaks e travas de segredo | ja existe no LegoLoop (leak-scan) | descartado |
| 35 | Teste de fumaca `testes.py` (t_publicar, t_canal, t_sheet, t_vps, t_hook) | testes Go com dubles de borda | migrar como TDD |

## 4. Riscos e efeitos no mundo real

- Publicar e irreversivel de fato (video publico no canal). Padrao do processo: privado; publico so com dupla trava.
- Os testes nunca chamam Postiz, YouTube, Instagram ou TikTok de verdade: so dubles de borda e `--dry-run`.
- Segredos (`POSTIZ_API_KEY`, chaves de rede) nunca entram em artefato, log ou chat.
- Postiz e o VPS sao do usuario; o processo so exige que existam.

## 5. Plano proposto em sub-tickets

- 3a: INTENT, SPEC, PLAN, ROUTE e DAG do P4 (a partir desta auditoria), com aprovacao do usuario em cada Stop-Gate.
- 3b: revisao (folhas de contato, parecer, versoes) e vitrine (titulos, capas, descricao, rotulo IA).
- 3c: acervo e agendador (quente/perene, validade, estoque, originalidade).
- 3d: publicacao via Postiz (dry-run por padrao, privado por padrao, dupla trava para publico).
- 3e: comentarios e aprendizado.
- 3f: canonizar (processo, skill, manual, livro, espinha, `fechar` 9 de 9).
