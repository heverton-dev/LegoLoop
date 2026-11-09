# 🎯 INTENT: Pipeline Canônico de Edição e Motion Vídeo Determinado (LegoLoop Video Motion)

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Autor:** Antigravity & User (LegoLoop Mega Pipeline)  
> **Iniciativa:** I01 — LegoLoop Creative Engines  
> **Alvo:** `internal/pipeline`, `loops/video-motion`, `skills/video-motion`  
> **Sessão:** 217415b3-cbf6-432f-bbe6-b267fa3fddea  
> **Confiança da Entrevista:** 98%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** O fluxo de produção de vídeos curtos verticais (Reels/TikTok/Shorts) e vídeos explicativos em motion graphics demonstrado por criadores como Kaye Creatives combina alto valor de retenção com graves gargalos operacionais:
  1. *Dependência de ferramentas proprietárias e caras:* o fluxo original exige After Effects (pesado, caro, incompatível com automação headless e CI/CD).
  2. *Render cego e falta de preview interativo:* no render puramente via scripts FFmpeg/libass, qualquer microajuste de texto ou timing exige minutos de renderização cega sem hot-reload.
  3. *Alucinação e custos de IA sintética:* modelos de vídeo generativo (Higgsfield/Seedance) alucinam transições, deformam tipografia e cobram créditos a cada tentativa frustrada.
  4. *Lock-in de harness:* a automação original foi construída com exclusividade para Claude Code rodando em macOS com Apple Vision.
* **Quem Sofre:** Criadores de conteúdo, educadores técnicos e desenvolvedores que buscam produzir vídeos verticais e motion design programático de nível broadcast sem depender de After Effects, sem gastar fortunas com GPUs generativas e com portabilidade universal (Windows, Linux, macOS) em qualquer harness (Gemini/Antigravity, Claude, Cursor, Windsurf).
* **O Porquê:** O LegoLoop foi projetado para elevar o inventário do sistema operacional a **Padrão Ouro**. O SO já contém compilador Go, FFmpeg, Python, Node, o kit de engenharia de 29 regras de `remotion-video-creation`, motores de alinhamento fonético milimétrico (`speech-sync`) e conectores de IA generativa (Google AI Studio Gemini 2.0 / Imagen 3 / Veo 2). Unificar essas peças em um pipeline acíclico DAG modular elimina 100% dos custos de After Effects e transforma o processo em uma esteira industrial determinística.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O usuário precisa executar scripts manuais em Python, abrir After Effects para ajustar cartões de vidro, depender de APIs pagas de vídeo do Higgsfield e rodar comandos cegos sem preview em tempo real.
  * **Amanhã:** Um comando LegoLoop único (`legoloop run loops/video-motion/pipeline.dag.json`) ingere um vídeo ou áudio bruto, corta silêncios milimetricamente no ataque das sílabas, alinha cada palavra em um JSON canônico (`words.json`), gera ativos pontuais via Google AI Studio, monta a cena em código reativo (Remotion/Canvas) com preview no navegador (hot-reload), valida os gates binários (áudio em -14 LUFS, zero colisão de texto, conformidade cromática) e renderiza um MP4 broadcast final com opção de exportação FCPXML/EDL para DaVinci Resolve.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] **Decomposição e Canonização dos Dois Modos da Autora:**
  1. *Modo A (Talking-Head Reel):* Ingestão de vídeo do apresentador, corte inteligente de pausas, legendas cinéticas no estilo Anton + Didot, split-screens e inserts de B-roll sincronizados.
  2. *Modo B (Voiceover Glass Motion):* Ingestão de áudio/narração pura, cartões de vidro 3D flutuantes com borda luminosa, grid em perspectiva com troca de cor da palavra falada (*word-flip*) e diagramas UI procedural.
* [ ] **Suporte Multi-Aspect Ratio Responsivo (Safe Zones):**
  * `9:16` (1080x1920) — Vertical para TikTok, Reels e YouTube Shorts (com safe zone anti-sobreposição de botões de UI social).
  * `16:9` (1920x1080) — Horizontal para YouTube tradicional, cursos e desktop.
  * `1:1` (1080x1080) — Quadrado para Feed do LinkedIn e Instagram.
* [ ] **Bundle de Entrega Broadcast Completo (`dist/`):**
  * Vídeo principal renderizado (`master.mp4`) em -14 LUFS.
  * Vídeo limpo sem legendas queimadas (`clean_feed.mp4`) para reuso.
  * Arquivos de legenda fechada sincronizada (`subtitles.srt` e `subtitles.vtt`).
  * Pacote de Copywriting para Redes (`POST_COPY.md` / `copy.json`) com título, gancho de retenção, descrição e hashtags segmentadas por rede (IG, TikTok, YT, LinkedIn).
  * Linha do tempo em FCPXML/EDL para abertura direta no DaVinci Resolve Studio.
* [ ] **Substituição Agnóstica de Mídia Sintética:** Substituir a dependência do conector Higgsfield MCP por clientes agnósticos do Google AI Studio (Gemini 2.0 Flash para roteirização e breakdown de beats; Imagen 3 / Veo 2 para B-roll e fundos), com fallback mecânico para bibliotecas locais de textura/halftone.
* [ ] **Relógio Mestre Fonético (Speech-Sync Engine):** Alinhamento fonético palavra por palavra com timestamp milimétrico e detecção de pico de ataque para sincronizar cortes e animações sem comer sílabas.
* [ ] **Camada de Motion Visual em Código (Remotion / Web-Canvas):** Implementação baseada nas 29 regras canônicas de `remotion-video-creation`, permitindo preview em tempo real no navegador local antes do render final.
* [ ] **Os 4 Pilares Determinísticos Padrão Ouro:**
  1. *Cache de Ativos por SHA-256:* Indexação criptográfica de cada ativo por hash do prompt, evitando regerações e consumo desnecessário de cota de API.
  2. *Detector de Ênfase Vocal por Energia RMS:* Análise determinística da curva de decibéis do áudio para marcar palavras de impacto vocal (`punch_word`) sem intervenção manual.
  3. *Auto-Framing com Rastreamento Facial:* Centralização suave do orador em conversões 16:9 ➔ 9:16 mantendo a cabeça e olhos na zona de segurança áurea.
  4. *Folha de Contato Instantânea (QA Contact Sheet):* Geração em 2 segundos de `qa_sheet.jpg` com grade dos 12 frames-chave para aprovação humana rápida antes do render completo.
* [ ] **Princípio Inviolável da Máxima Determinização:**
  * O que pode ser resolvido por matemática, código procedural, regras CSS e análise de sinal (FFmpeg, RMS, Whisper, Remotion) é **100% determinístico**.
  * A IA probabilística (Gemini/Imagen/Veo) fica estritamente contida na síntese de B-roll e é cercada por cache determinístico e fallbacks locais.
* [ ] **Meta-Processo Canônico:** Criação de um manual e skill documentando passo a passo a metodologia de engenharia reversa de tutoriais de vídeo, extração de ferramentas gratuitas e acoplamento ao LegoLoop.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos ferramentas de agendamento ou postagem automática em redes sociais (Buffer, APIs de redes). Isso será composto como um pipeline separado via `/legoloop-compor`.
* ⛔ Não faremos dependência mandatória do Adobe After Effects (o motor primário é 100% código aberto/web).
* ⛔ Não faremos amarras a APIs pagas com travas de fornecedor sem opção de fallback gratuito ou local.
* ⛔ Não faremos código acoplado exclusivamente a macOS; o pipeline roda nativamente no Windows 11 e Linux.
* ⛔ Não alteraremos os contratos centrais do core do LegoLoop (`internal/dag`, `internal/gate`).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `aidd-rea-video-youtube` / REA | Engenharia reversa e dissecação de assets da referência. |
| `BLOCO-02` | `unit_task` | `speech-sync-clock` (Go + FFmpeg + Whisper) | Detecção de silêncio, corte cirúrgico e geração do `words.json`. |
| `BLOCO-03` | `skill` | `google-ai-studio-broll` (Gemini 2.0 / Imagen 3) | Geração de roteiro de beats e assets visuais sintéticos de B-roll. |
| `BLOCO-04` | `unit_task` | `remotion-video-creation` / Studio Canvas | Composição procedural 3D de cartões de vidro e tipografia cinética com hot-reload. |
| `BLOCO-05` | `stop_gate` | `SG-PREVIEW-HUMANO` | Validação visual humana de frames-chave no navegador antes do render. |
| `BLOCO-06` | `script` | `ffmpeg -filter_complex loudnorm` | Render final em 1080x1920 (9:16) ou 1920x1080 (16:9) e masterização a -14 LUFS. |
| `BLOCO-07` | `unit_task` | `davinci-fcpxml-exporter` | Exportação de timeline em FCPXML/EDL para interoperabilidade de estúdio. |
| `BLOCO-08` | `stop_gate` | `go test ./... && go build` | Gate de integridade e recompilação fresca do LegoLoop. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` via CLI (`legoloop run ...`) ou Studio visual.
* **Nível de Risco & Blast Radius:** `Médio (workspace local com isolamento de cache em scratch/)`.
* **Política de Privacidade & Orçamento Cognitivo:** Nuvem autorizada exclusivamente para chamadas de API do Google AI Studio (com chave em variável de ambiente `GEMINI_API_KEY`). Todo o processamento de áudio, corte de vídeo e montagem de código roda 100% localmente na máquina.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. Geração de vídeo MP4 funcional sincronizado com áudio em -14 LUFS.
  2. Ausência total de erros de sobreposição de legendas sobre o rosto do apresentador.
  3. Timeline aberta e navegável no browser via Remotion/Canvas.
  4. Arquivo FCPXML/EDL importável no DaVinci Resolve.
* **Premissas Assumidas:**
  1. O usuário possui `ffmpeg` funcional no PATH (confirmado: `C:\ProgramData\chocolatey\bin\ffmpeg.exe`).
  2. O ambiente possui Node.js disponível para o renderizador Web (confirmado: `node.exe v26.7.0`).
  3. O modelo Gemini 2.0 Flash / Imagen 3 pode ser acionado via API padrão.
* **Riscos e Perguntas em Aberto:**
  * Cota de geração de imagem no Google AI Studio (gerenciada por cache de assets para evitar re-geração desnecessária).

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Alinhamento fonético milimétrico palavra por palavra. | Geração do `words.json` com início e fim de cada fonema sem colisões. | `python -m json.tool words.json` |
| 2 | Eliminação de After Effects com fidelidade visual aos cartões de vidro e Anton kinetic. | Renderização via Remotion/HTML-Canvas no navegador com paridade visual aos stills de referência. | Verificação de frames-chave (`sheet.jpg`) |
| 3 | Normalização de áudio em padrão broadcast social. | Análise de loudness comprovando medição em -14 LUFS (±1 LUFS). | `ffmpeg -i render.mp4 -af ebur128 -f null -` |
| 4 | Suíte de testes e integridade do LegoLoop 100% verde. | Zero regressões no core do LegoLoop. | `go test -count=1 ./...` exit code 0 |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md`](../specs/2026-11-07_SPEC_VIDEO_MOTION_PIPELINE.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat com 98% de certeza factual, antes de qualquer linha de especificação técnica detalhada ser escrita.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende rigorosamente à sua expectativa?*  
*Somente após a sua aprovação explícita, redigiremos a SPEC técnica detalhada com critérios BDD binários.*
