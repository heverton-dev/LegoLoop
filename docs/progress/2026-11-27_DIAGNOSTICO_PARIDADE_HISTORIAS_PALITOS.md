# Diagnóstico de Paridade: `proj--historias-palitos` vs. `legoloop-video-animado`

Data: 2026-11-27
Origem auditada: `C:\Users\trcnologia\Desktop\proj--historias-palitos`
Alvo no LegoLoop: `internal/video/*`, `internal/cli/*`, `processos/legoloop-video-*`

## 1. O que entendi (visão completa do original vs. onde quebramos)

O projeto `proj--historias-palitos` é uma fábrica completa de vídeos explicadores e canais (`PROCESSO.md`, `AUTORIA.md`, `KIT_DE_CANAL.md`, `MINIATURA.md`, `PUBLICACAO.md`). Ao trazê-lo para o LegoLoop, canonizamos cascas de validação em Go, mas **decepamos a autoria visual, o multiformato, a marca e o kit de canal** na ligação entre os processos:

1. **Multiformato (`16:9`, `9:16`, `1:1`) travado em `9:16`:**
   - Original (`motor-local/formato.py:10-12`, `render.py:102`): desenha o mundo em `1920x1080` (`16:9`) e reenquadra matematicamente para `9:16` (`1080x1920`) e `1:1` (`1080x1080`), ajustando câmera, personagens e camadas HUD na área segura.
   - Erro no LegoLoop: `internal/cli/video_render_etapas.go:15` fixou `const formatoMotor = "9:16"` e `internal/video/motor/render.go:107` reprova qualquer vídeo fora de `1080x1920`, ignorando `projeto.json.formatos`.

2. **Ponte e Storyboard mudos (zero motion graphics, textos foco, cartões e câmera):**
   - Original (`AUTORIA.md:8-37`, `build_timeline.py:58-140`, `montar.py:60-115`): cada frase vira um diorama com metáfora própria, subcortes a cada ~3,3 s, cartões pretos digitados (`type: "card"`), textos foco na tela (`kind: "text"`, até 4-5 palavras ou número da linha, `in: "pop"`), camadas animadas (`doc`, `stamp`, `bars`, `poly` com `draw`, `calendar`, `stopwatch`, `counter`, `pin`), efeitos (`fx`), `shake`, câmera ciclando (`in`, `panR`, `out`, `low`, `panL`) e herói fora do centro (`x: 520..700` ou `1250..1350`) em no máximo 3 cortes seguidos.
   - Erro no LegoLoop: `internal/video/ativos/storyboard.go:229` gera 4 placeholders vazios e `internal/video/motor/ponte.go:74,221` só serializa `set` + `chars` parados em `x=960`, descartando `layers`, `cam`, `fx`, `shake`, `grade` e `card`.

3. **Cenários e personagens reaproveitados do First Brands sem geração temática:**
   - Original (`PROCESSO.md:119-153`, `AUTORIA.md:43-46`, `curiosidades-malucas-v2/prompts_flow.json`): quando o tema exige ativos próprios, gera-se no Google Flow (0 créditos em imagem, 1 projeto Flow por vídeo) as fichas do elenco, props e fundos 16:9 sem texto no estilo mestre (`flow.py` + `key_sprites.py`), ou usam-se cenários em código parametrizados para o tema (`server`, `office` com `sign`/`calendar` próprios, nunca `shop` com `"AUTOPEÇAS"` de `motor-local/sets.js:191`).
   - Erro no LegoLoop: reutilizamos os 8 cenários do First Brands (incluindo `loja-fachada` -> `"AUTOPEÇAS"`) e um boneco executivo único.

4. **Pós-produção de Marca (`marca.py`) e Cortes Verticais (`cortes.py`) ausentes na execução:**
   - Original (`ferramentas/marca.py:1-69`, `ferramentas/cortes.py:1-143`): aplica logo (`marca/logo.png`), vinheta de abertura e vinheta de fechamento no formato de saída, gera `.srt` compensado pela abertura e re-renderiza de 3 a 8 cortes curtos `9:16` a partir da timeline (não recortados a seco).
   - Erro no LegoLoop: `marca.py` não foi integrado ao P3 e os cortes `9:16` foram empurrados para `BL-03`.

5. **Kit de Canal, Logo, Thumbnails do Flow e `colar.md` descartados por engano:**
   - Original (`KIT_DE_CANAL.md`, `MINIATURA.md`, `ferramentas/canal.py`, `ferramentas/redes.json`): define a identidade do canal (`textos.json`), gera `colar.md` (textos prontos para copiar e colar por rede com contagem de caracteres), gera logo/perfil e banner base no Flow sem texto e exporta os tamanhos oficiais + prévias (`perfil_1080.png`, `youtube_perfil_800.png`, `youtube_marca_dagua_150.png`, `youtube_banner_2560x1440.png`, prévias `98px` e área segura), além das 3 thumbnails `1280x720` com fundo gerado sem texto + tipografia e logo da marca.
   - Erro no LegoLoop: em `AUDITORIA_P4_VIDEO_DISTRIBUICAO.md:43`, o item 27 (`canal.py` / `colar.md`) foi descartado sob a premissa falsa de que `legoloop-canal-redc` já o fazia (`internal/redc` faz apenas engenharia reversa de DNA, não gera Kit de Canal).

---

## 2. O que vou fazer (Plano de Execução ponta a ponta, sem regredir)

Preservando tudo o que evoluímos no LegoLoop (Whisper WER, mix `-14 LUFS` com ducking e SFX, inventário por estilo SHA-256, gates determinísticos, agendador e Postiz seguro), vou fechar todas as lacunas com TDD e provar em execução real via CLI/DAG/Studio:

1. **Etapa A — Kit de Canal e Identidade (`video-distribuicao canal` + integração no macroprocesso):**
   - Portar `canal.py` + `redes.json` para `internal/video/distribuicao` (`legoloop video-distribuicao canal textos|imagens`): valida `textos.json`, gera `colar.md` com contagem de caracteres por rede e monta logo, `perfil_1080.png`, `youtube_perfil_800.png`, `youtube_marca_dagua_150.png`, `youtube_banner_2560x1440.png` e prévias `98px`/área segura a partir das bases do Flow.
   - Integrar fundos dedicados sem texto (Flow) + logo do canal nas 3 thumbnails `1280x720` (`MINIATURA.md`).

2. **Etapa B — Autoria Completa no Roteiro, Ativos e Ponte (`P1 -> P2 -> P3`):**
   - Expandir o contrato de `script.json`/`plano.md` e `cenas.json` v2 para carregar a autoria real do diorama por linha (`AUTORIA.md` / `montar.py`): `cartao` (cartão preto digitado), `texto` e `texto2` foco na tela, `cor`, `camadas` (`doc`, `stamp`, `bars`, `poly`, `calendar`, `stopwatch`, `counter`, `pin`, `rect`), `camera` (`in`, `out`, `panL`, `panR`, `low`), `fx`, `shake`, `grade`, posição `x` fora do centro, alternância de até 3 cortes com pessoa e parâmetros de cenário (`sign`, `calendar`, `time`).
   - Atualizar `internal/video/motor/ponte.go` para emitir todos esses campos no `timeline.json`, partir linhas longas em subcortes (~3,2 s) e validar no `check_timeline.py` com `>= 4` elementos em movimento reais por cena.
   - Suportar de verdade os dois caminhos de ativos para qualquer tema: geração real no Google Flow (elenco, props, fundos 16:9 sem texto e base de thumbnails/canal em projeto isolado do Flow) e cenários/elenco vetoriais temáticos sem lixo herdado do First Brands.

3. **Etapa C — Render Multiformato (`16:9`, `9:16`, `1:1`), Marca (`marca.py`) e Cortes (`8b`):**
   - Destravar `internal/video/motor/render.go` e `internal/cli/video_render_etapas.go` para ler `projeto.json.formatos` e `--formato` (`16:9` -> `1920x1080`, `9:16` -> `1080x1920`, `1:1` -> `1080x1080`).
   - Integrar pós-produção da marca (vinheta de abertura, logo em área segura, vinheta de fechamento e `.srt` sincronizado) e a geração de cortes curtos `9:16` re-renderizados da timeline (`cortes.py`).

4. **Etapa D — Recomposição do Macroprocesso, CLI/Studio e Prova Real de Ponta a Ponta:**
   - Atualizar DAGs, skills, manuais e espinha dorsal para que o macroprocesso `legoloop-video-animado` cubra da definição do canal/configuração até o pacote final pronto/enviado ao Postiz.
   - Rodar de novo o vídeo real de **Tecnologia e Inteligência Artificial** em **16:9**, **3 minutos**, voz **Thalita**, gerando canal, logo, `colar.md`, ativos temáticos, motion graphics completos, thumbnails, folhas de contato e revisão até aprovação visual.
