# 🎯 INTENT: Edição Automatizada de Vídeo Soberana (video-edition)

> **Status:** APPROVED  
> **Data:** 2026-11-08  
> **Autor:** LegoLoop Video-REA Engine  
> **Iniciativa:** Mega-Pipeline Audiovisual Soberano  
> **Alvo:** `loops/video-edition`, `skills/legoloop-video-edition`, `docs/videorea/lucas-leite-edicao-claude-code`  
> **Sessão:** intv_1794143944778492400  
> **Confiança da Entrevista:** 95% (Meta de 95% atingida via motor socrático oficial)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criadores de conteúdo perdem horas diárias em tarefas manuais, repetitivas e cansativas de pós-produção audiovisual (cortar silêncios, reenquadrar tomadas para formato vertical, inserir pontuações de SFX e renderizar). O método dissecado automatiza esse fluxo, mas o autor original dependia de assinaturas comerciais caras (Adobe Premiere Pro e Claude Pro) que custam R$ 329,00/mês.
* **Quem se Beneficia:** Criadores de conteúdo, editores de vídeo e operadores do ecossistema LegoLoop automatizando pós-produção ágil e sem custo de licenças.
* **O Porquê:** Eliminar trabalho braçal e repetitivo de montagem de shorts/reels, unificando corte, enquadramento e sound design em uma esteira reproduzível e estéril denominada `video-edition`, operando a custo R$ 0,00/mês.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O usuário precisa de assinaturas comerciais e depende de ajustes manuais em editores proprietários.
  * **Amanhã:** Um pipeline acíclico no LegoLoop (`video-edition`) que recebe o vídeo bruto e executa a pós-produção de ponta a ponta, entregando vídeo MP4 1080p60 sincronizado e linha de tempo aberta FCPXML.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Pipeline de nós LegoLoop (`video-edition`) acionado via CLI/Studio para ingestão de vídeo bruto.
* [ ] Objetivo 2: Corte automático de silêncio via WhisperX / Whisper Local VAD com margem de segurança de +40ms pré-fala.
* [ ] Objetivo 3: Reenquadramento facial dinâmico 4K ➔ 1080p via OpenCV YuNet (com zoom em transições e punch words).
* [ ] Objetivo 4: Inserção automática de SFX parametrizados com ducking estrito de -18dB sob a voz principal.
* [ ] Objetivo 5: Formato híbrido de entrega: arquivo MP4 final consolidado (1080p60) via FFmpeg e timeline FCPXML exportável para DaVinci Resolve ou Premiere Pro.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Interface gráfica proprietária de edição linear (NLE).
* ⛔ Geração de legendas dinâmicas animadas estilo Karaokê nesta fase.
* ⛔ Envio direto para redes sociais (upload API do YouTube/TikTok).
* ⛔ Uso de serviços comerciais de nuvem pagos quando o FFmpeg/Whisper local realizar a tarefa.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `unit_task` | `ffmpeg` | Extração de áudio puro 48kHz WAV e separação de canais |
| `BLOCO-02` | `unit_task` | `whisperx` / `whisper.cpp` | Transcrição fonética com timestamps alinhados por palavra |
| `BLOCO-03` | `script` | `python / silenceremove` | Detecção de pausas (>0.4s) e corte cirúrgico com +40ms de ataque |
| `BLOCO-04` | `unit_task` | `opencv-yunet` | Auto-enquadramento de câmera 4K ➔ 1080p centrado no apresentador |
| `BLOCO-05` | `script` | `yt-dlp` | Download e catalogação de b-rolls e referências visuais |
| `BLOCO-06` | `unit_task` | `hyperframes` / `remotion` | Render de cartões de texto, overlays e tipografia Anton |
| `BLOCO-07` | `script` | `python-audio` | Mixagem de efeitos sonoros com Ducking automático de -18dB |
| `BLOCO-08` | `unit_task` | `fcpxml-generator` | Composição final da timeline aberta e render MP4 FFmpeg |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via CLI `legoloop run loops/video-edition/pipeline.dag.json` ou skill `/legoloop-video-edition`).
* **Nível de Risco & Blast Radius:** Baixo (opera estritamente dentro da pasta isolada `docs/videorea/lucas-leite-edicao-claude-code/`).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente Local; modelos via 9Router; chaves de API restritas e protegidas.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Execução ponta a ponta gerando arquivo MP4 1080p60 sincronizado sem artefatos de áudio/vídeo e FCPXML válido importável no Premiere/DaVinci com cortes e trilhas de SFX alinhados.
* **Premissas Assumidas:** Processamento local via FFmpeg, Whisper/WhisperX e OpenCV YuNet; cada nó opera em contexto isolado lendo e gravando artefatos em disco; sem stubs no código; ducking estrito de -18dB sob a voz principal.
* **Riscos e Perguntas em Aberto:** Nenhuma lacuna em aberto após conclusão da Entrevista Socrática.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Curadoria aprovada com custo R$ 0,00/mês e diagnósticos de host integrados. | Inspeção de `curation.json` com `approved: true` | `Test-Path docs/videorea/lucas-leite-edicao-claude-code/curation.json` |
| 2 | Grafo DAG compilado e validado em 8 etapas acíclicas sob o slug `video-edition`. | Validador oficial do LegoLoop | `go run ./cmd/legoloop validate loops/video-edition/pipeline.dag.json` |
| 3 | Entrega de MP4 1080p60 e FCPXML sincronizados. | Verificação de existência e integridade dos arquivos | `Test-Path output/final_edit.mp4; Test-Path output/final_edit.fcpxml` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** `docs/videorea/lucas-leite-edicao-claude-code/governance/2026-11-08_SPEC_VIDEO_EDITION.md`
* **Gate de Transição:** Este Intent foi refinado via Entrevista Socrática oficial atingindo 95% de confiança. Ele avança para a geração da SPEC somente após a aprovação explícita do operador no Stop-Gate 01.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com o nome canônico `video-edition` (95% de confiança) atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*
