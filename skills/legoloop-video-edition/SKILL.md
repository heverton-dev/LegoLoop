---
name: legoloop-video-edition
description: Use quando precisar executar ou orquestrar o Processo autoral de edição audiovisual assistida por IA (metodologia Lucas Leite), com cortes semânticos por frase, alinhamento fonético com folga pré-fala (+40ms), extirpação automática de notas de edição faladas, enquadramento facial inteligente YuNet (altura dos olhos a 35%), motion design procedural via Remotion, sonorização sincronizada com ducking de voz (-18dB) e entrega híbrida não-destrutiva (FCPXML + MP4 direto).
tags: [audiovisual, video-rea, video-edition, remotion, ffmpeg, yunet, claude-code, lucas-leite]
---

# 🎬 Skill Universal: `legoloop-video-edition`

Processo canônico de edição audiovisual assistida por IA, dissecado via engenharia reversa metódica do método autoral de **Lucas Leite** e canonizado sob governança LegoLoop v2.2.0.

---

## 1. Comandos de Ativação & Triggers

* **Linha de Comando (CLI):**  
  ```bash
  legoloop run processos/video-edition/processo.dag.json --var input_video="caminho/do/video.mp4"
  ```
* **Chat Interativo:**  
  * `/legoloop-video-edition`
  * `/video-edition`
  * *"editar vídeo com claude code"*
  * *"aplicar Processo lucas leite"*

---

## 2. Princípios Invioláveis da Metodologia Lucas Leite

1. **Separação Rígida (Cognição vs Cálculo):** O LLM nunca calcula timestamps ou frames em segundos. O modelo opera sobre sentenças e marcações textuais; scripts determinísticos (`internal/videomotion`, FFmpeg, Python) calculam frames, limites de corte e energia RMS.
2. **Folga Pré-Fala (+40ms):** Todo corte que antecede o início de uma fala inclui uma margem obrigatória de +40ms para evitar decapitação fonética de consoantes plosivas.
3. **Extirpação de Notas de Edição:** Trechos onde o apresentador fala orientações de montagem (*"nota de edição...", "corta isso...", "fim da nota"*) são detectados no texto e expurgados automaticamente antes da montagem da timeline.
4. **Câmera Pós-Motion & Eye-Level Estabilizado:** O enquadramento via YuNet é processado após a definição dos cartões gráficos, mantendo a linha dos olhos do apresentador fixa em 35% do topo do quadro.
5. **Ducking & SFX na Frenagem:** Efeitos sonoros (*whoosh*, *impact*, *click*) têm seu pico de energia alinhado ao frame exato da frenagem do elemento visual. Áudios de apoio sofrem ducking automático de -18dB sob a fala do locutor.
6. **Entrega Híbrida Não-Destrutiva:** O Processo gera simultaneamente o render direto MP4 e a timeline aberta em FCPXML/EDL com marcadores de capítulo e cartões em camadas separadas para ajuste fino no DaVinci Resolve ou Premiere.

---

## 3. Topologia do Processo DAG (12 Nós em 11 Camadas)

```
[n01_iniciar_preparar]
         │
[n02_transcricao_alinhamento]
         │
[n03_decisao_cortes]
         │
[n04_sg_revisao_corte] 🛑 STOP-GATE HUMANO (Validação da Seleção de Frases)
         │
[n05_motion_mapa]
         │
[n06_sg_revisao_motion] 🛑 STOP-GATE HUMANO (Aprovação do Storyboard Gráfico)
        ┌┴──────────────────────┐
[n07_camera_yunet]     [n08_assets_download]
        └┬──────────────────────┘
[n09_motion_render]
         │
[n10_efeitos_sonoros_sfx]
         │
[n11_entrega_final_hibrida]
         │
[n12_sg_entrega_final] 🛑 STOP-GATE HUMANO (Aprovação da Entrega FCPXML + MP4)
```

---

## 4. Isolamento Hermético em Produção

Todo processamento intermediário (áudio WAV extraído, transcrições JSON, frames de detecção facial e renders parciais) ocorre estritamente dentro de diretórios efêmeros em `scratch/`:
* Entrada: `var.input_video`
* Isolamento: `"execution": { "mode": "isolated" }`
* Saídas Canônicas: `dist/timeline.fcpxml`, `dist/render_final.mp4` e `dist/manifest_final.json`.
