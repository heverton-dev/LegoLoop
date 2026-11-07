# 📘 MANUAL OPERACIONAL: FUNÇÃO VIDEO MOTION NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Pipeline determinístico de edição e motion vídeo para redes sociais e documentários  
> **Fontes de verdade:** `internal/videomotion/` · `skills/legoloop-video-motion/SKILL.md` · `loops/video-motion/pipeline.dag.json`

---

## 1. Visão Geral

A função **`legoloop-video-motion`** transforma arquivos de vídeo, áudio gravado ou locuções geradas por voz sintética em pacotes finais de entrega broadcast com padrão de estúdio. O fluxo opera sobre 4 pilares determinísticos: relógio fonético com silêncios cirúrgicos (+40ms de ataque), detecção de ênfase vocal via RMS, composição de cartões 3D translúcidos com vidro fosco e tipografia cinética Anton, e exportação multi-aspecto (9:16, 16:9, 1:1) com legendas SRT/VTT e linha de tempo FCPXML para DaVinci Resolve.

---

## 2. Comandos e Uso

### 2.1 Pelo Chat do Agente
* `/legoloop-video-motion <arquivo>`
* `/video-motion <arquivo>`

### 2.2 Pela CLI Mecânica
```bash
# Execução completa da esteira
legoloop video-motion run input.mp4 --aspect 9:16

# Etapas modulares
legoloop video-motion silence --input audio.wav --out scratch/words.json
legoloop video-motion punch --words scratch/words.json --out scratch/punched.json
legoloop video-motion assets --words scratch/punched.json --out scratch/assets.json
legoloop video-motion preview --assets scratch/assets.json --out scratch/qa_sheet.jpg
legoloop video-motion bundle --preview scratch/qa_sheet.jpg --out dist/
```
