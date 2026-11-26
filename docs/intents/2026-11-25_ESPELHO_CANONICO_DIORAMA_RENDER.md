# Espelho Canônico — `legoloop-diorama-render` (Ciclo 3)

> **Nota de redirecionamento (2026-11-26, AC-E11.4):** o processo `legoloop-diorama-render` foi renomeado para `internal/video/render (ligado ao render Go existente, AC-E10.4)`. Este documento e historico e mantido sem alteracao de conteudo.


> **Data:** 2026-11-25
> **Sessão A2UI:** `sess_3` (entrevista conduzida via MCP `sessao_*`, 6 turnos)
> **Confiança:** 95% (meta: 95%)
> **Status:** CONFIRMADO PELO USUÁRIO NA INTERFACE A2UI (turno 6, opção A)
> **Emenda 2 (2026-11-25, TICK-DR-02.2):** WER passa de `<= 5%` por linha para `<= 15%` por linha e `<= 5%` global, com parada em áudio regravado idêntico (provas reais A1 e A2).  
> **Emenda 1 (2026-11-25, Stop-Gate 02):** 4K opt-in passa de `projeto.json` para `render.json`; correção do True Peak em `videomotion` autorizada.

## 1. Resultado Esperado
Processo canonizado `legoloop-diorama-render` (Peça 3 do macroprocesso `legoloop-video-diorama`) que consome `cenas.json` e `script.json` e entrega, em ordem de dependência:
1. **Etapa 6 — Voz neural e conferência Whisper:** síntese por provedor plugável (padrão `edge-tts` com `pt-BR-ThalitaMultilingualNeural`; Gemini TTS, ElevenLabs, Fish Audio ou qualquer outro via adaptador); `faster-whisper` sobre o áudio final gera `timings.json` por palavra (fonte única de alinhamento) e o WER (limite `<= 15%` por linha e `<= 5%` global, até 3 regravações só das linhas reprovadas, com parada se a regravação devolver áudio idêntico); fonemas/visemas para sincronia labial simples.
2. **Etapa 8 (do pedido) — Timeline em camadas:** leitura estrita de `cenas.json` e `script.json`; composição 2.5D (cenário, adereços em profundidade, sprites com troca de expressão/pose, câmera com zoom, shake e pan); legendas karaokê palavra a palavra; `timeline.json` validado por `check_timeline` (12–21 cortes/min, mediana 2,5–4 s).
3. **Etapa 7 (do pedido) — Mixagem e SFX:** BGM com ducking dinâmico sob a narração; SFX nos eventos da timeline; normalização `-14 LUFS ± 0,5` integrado e `-1,5 dBTP`.
4. **Etapa 8b — Render e cortes:** master 16:9 em 720p, 1080p (padrão) e 4K (opt-in em `render.json`) via ffmpeg com detecção automática de NVENC/QSV/AMF e fallback `libx264`; cortes 9:16 re-renderizados a partir da timeline, centrados no narrador e no elemento de ação, legenda fora da caixa do narrador, 15–59 s e área segura vertical.

Ordem de execução por dependência (voz, timeline, mix com SFX nos eventos da timeline, render, cortes), mantendo a numeração do pedido nos títulos (7 = mixagem, 8 = timeline).

## 2. Quem se Beneficia
O operador do canal (vídeo diorama pronto sem edição manual) e o Processo 4 (`legoloop-diorama-distribuicao`), que consome renders MP4, `legenda.srt` e shorts como contrato imutável.

## 3. Por que Agora
O Ciclo 2 (`legoloop-diorama-ativos`) já entrega `cenas.json` e sprites. Sem o Ciclo 3 não existe vídeo, e o Macroprocesso ponta a ponta (Ciclo 5) fica bloqueado.

## 4. Critérios de Sucesso
1. Em `projetos/video-diorama/diorama-demo`, todos os gates mecânicos passam: WER `<= 5%` global e `<= 15%` por linha, LUFS em `-14 ± 0,5`, TP `<= -1,5 dBTP`, duração áudio == vídeo, zero frames em branco.
2. `timeline.json` aprovado por `check_timeline`.
3. Cortes 9:16 entre 15 s e 59 s com área segura respeitada e legenda fora da caixa do narrador.
4. Troca de provedor de voz sem alterar `timings.json` nem os gates a jusante.
5. Cobertura `>= 95%` nos pacotes novos, zero stubs, Espinha Dorsal verde e `legoloop fechar` 5 de 5.

## 5. Restrições Inegociáveis
1. Provedor de voz plugável por interface; padrão `edge-tts` Thalita; chamada paga só após gate de créditos.
2. Alinhamento por palavra sempre pelo `faster-whisper` sobre o áudio final; timestamps nativos do provedor são ignorados.
3. Motor autoral, sem código do `RobbiSkill/`.
4. Fallback `libx264` obrigatório; encoder usado registrado no relatório.
5. 9:16 sem visão computacional: determinístico a partir das caixas da timeline.
6. Testes 100% herméticos (ffmpeg, Whisper e provedores injetáveis), zero stubs, TDD Red-Green.
7. Contratos de entrada dos Ciclos 1 e 2 imutáveis; a configuração do Ciclo 3 vive em `render.json` opcional.
8. O teto de True Peak `-1,5 dBTP` tem fonte única em `videomotion.ValidateLoudnessConformity`, corrigida com TDD (autorizado no Stop-Gate 02).

## 6. Fora de Escopo
1. Geração de ativos e storyboard (Ciclo 2) e roteiro (Ciclo 1).
2. Revisão estéril, vitrine, Postiz, comentários e aprendizado (Ciclo 4).
3. Composição ponta a ponta `P1->P4` (Ciclo 5).
4. Visão computacional ou crop por saliency no 9:16.
5. Painel web paralelo fora do Looping Studio.
