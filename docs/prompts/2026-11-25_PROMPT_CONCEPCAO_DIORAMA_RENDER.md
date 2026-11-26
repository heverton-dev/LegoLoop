# Prompt de Concepção A2UI — Processo `legoloop-diorama-render`

Cole o bloco abaixo em uma nova sessão para abrir a interface visual interativa (A2UI) e conduzir a concepção formal das Etapas 6, 7, 8 e 8b do macroprocesso de vídeo:

```text
Atue como o Arquiteto de Mídia e Renderização do LegoLoop.
Você deve conduzir a Concepção A2UI do Processo `legoloop-diorama-render` (Peça 3 do macroprocesso `legoloop-video-diorama`), cobrindo as Etapas 6, 7, 8 e 8b:

1. ETAPA 6 — VOZ NEURAL & CONFERÊNCIA WHISPER:
   - Síntese de fala neural em PT-BR (edge-tts voz Thalita / Gemini TTS) com marcas temporais (word-level timestamps).
   - Auditoria acústica via Whisper local para medição mecânica de WER (Word Error Rate <= 5%).
   - Geração de fonemas/visemas para sincronia labial simples nos sprites.

2. ETAPA 7 — MIXAGEM DE ÁUDIO & SFX:
   - Trilha sonora BGM de fundo com ducking dinâmico automático sob a narração.
   - Efeitos sonoros (SFX) sincronizados com as batidas de troca de cena e movimentos de objetos.
   - Normalização para padrão EBU R128 (-14 LUFS integrado, -1 dBFS True Peak).

3. ETAPA 8 — MOTOR DE COMPOSIÇÃO DE TIMELINE EM CAMADAS:
   - Leitura estrita de `cenas.json` e `script.json`.
   - Composição 2.5D: camada de fundo (cenário), adereços/objetos em profundidade, sprites com troca de expressão/pose e efeitos de câmera (zoom suave, shake, pan).
   - Legendas dinâmicas com realce palavra por palavra (karaokê).

4. ETAPA 8b — RENDERIZAÇÃO MULTIFORMATO & CORTES VERTICAIS:
   - Exportação master horizontal (16:9 1080p / 4K) via ffmpeg com aceleração por hardware.
   - Reenquadramento inteligente vertical (9:16) mantendo o foco de ação e bounding box do narrador.

DOCUMENTOS OBRIGATÓRIOS PARA LEITURA:
1. `docs/intents/2026-11-24_MAPA_MACRO_VIDEO_DIORAMA.md`
2. `projetos/video-diorama/diorama-demo/cenas.json`
3. `projetos/video-diorama/diorama-demo/script.json`
4. `skills/legoloop-concepcao/SKILL.md`

FLUXO DA CONCEPÇÃO:
1. Crie a sessão no Studio chamando `sessao_criar` (parâmetro `processo: "legoloop-diorama-render"`).
2. Publique as perguntas estruturadas na A2UI com opções recomendadas via `sessao_publicar` e chame `sessao_aguardar`.
3. Gere e aprove o artefato INTENT (Stop-Gate 01).
4. Elabore a SPEC técnica com critérios de aceite e FMEA (Stop-Gate 02).
5. Compile o PLAN, ROUTE e DAG para implementação TDD (Stop-Gate 03).
```
