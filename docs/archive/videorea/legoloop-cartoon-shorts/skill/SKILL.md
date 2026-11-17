---
name: legoloop-cartoon-shorts
description: Use quando o usuário quiser CRIAR SHORTS ANIMADOS CARTOON (1080x1920 30fps) com personagens 2D bean-head em mundo 3D Three.js, roteiro de curiosidades com gancho/twist a 3.7 palavras/s, narração neural Edge-TTS pt-BR com tags de emoção e muxing FFmpeg com trilha a -6 dB Gated RMS a custo R$ 0,00, via /legoloop-cartoon-shorts. NÃO use para editar vídeos de pessoas reais (use /legoloop-video-motion) ou compilar livros didáticos (use /legoloop-legopress).
consumes: Artifact:Theme
produces: Artifact:CartoonShortMP4
artifact: docs/videorea/legoloop-cartoon-shorts/scratch/short_final.mp4
gate: exit_code == 0
suggested_model: code-pro
---

# 🎨 legoloop-cartoon-shorts — Pipeline Soberano de Shorts Animados Cartoon

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário fornece um tema, curiosidade histórica, científica ou mistério cultural e deseja gerar um vídeo curto vertical completo (1080x1920 a 30fps) no estilo cartoon: personagens 2D *bean-head* desenhados em Canvas com contorno plum em mundo 3D colorido via Three.js r160, roteiro cronometrado a 3.7 palavras/segundo com gancho paradoxal imediato, escalada com fatos e reviravolta irônica, narração neural gratuita em português brasileiro via Edge-TTS (`pt-BR-AntonioNeural`) com anotações expressivas entre colchetes (`[reflexivo]`, `[animado]`, `[dá uma risadinha]`) e mixagem determinística de áudio com trilha a **50% (-6 dB Gated RMS)** e loudness normalizado em -14 LUFS.  
**Quando NÃO usar:** Editar vídeos gravados com câmera de pessoas reais (use `/legoloop-video-motion` ou `/legoloop-video-capcut`) ou dissecar vídeos externos da web (use `/legoloop-video-rea`).

* **Consumes (Entradas Obrigatórias):**
  - Tema ou pergunta de curiosidade (ex: `"O Enigma do ET de Varginha"`, `"A Corrida do Ouro Preto"`).
* **Produces (Saídas Mecânicas):**
  - Roteiro estruturado com cenas, timestamps e emoção (`scratch/script.json`).
  - Áudio de narração neural sintetizado (`scratch/voice.mp3`).
  - Vídeo final renderizado vertical 1080x1920 30fps com trilha em -6 dB (`scratch/short_final.mp4`).

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** utilizar serviços de nuvem pagos (ElevenLabs, Claude Pro pago, Midjourney, Runway). Todo o pipeline opera a custo R$ 0,00 via 9Router local, Edge-TTS e Three.js/FFmpeg local.
> 2. **NUNCA** permitir que a trilha musical ultrapasse 50% do nível da voz (-6 dB Gated RMS). A voz deve permanecer cristalina e inteligível em 100% da duração.
> 3. **NUNCA** gerar vídeos com duração arbitrária. A duração padrão recomendada é de **20 a 30 segundos** (máximo 60s) para manter retenção máxima nas redes sociais.
> 4. **NUNCA** gravar arquivos fora do diretório de workspace isolado (`docs/videorea/legoloop-cartoon-shorts/scratch/`).

---

## 3. Protocolo de Operação Passo a Passo

1. **Execução Completa via Linha de Comando (CLI):**
   ```bash
   legoloop cartoon-shorts --theme "O Enigma do ET de Varginha" [--duration 30] [--voice pt-BR-AntonioNeural]
   ```

2. **Execução em Modo Simulado (Dry-Run):**
   ```bash
   legoloop cartoon-shorts --theme "Minas Gerais" --dry-run
   ```

3. **Execução via Grafo DAG no Looping Engine:**
   ```bash
   legoloop run loops/legoloop-cartoon-shorts/pipeline.dag.json
   ```

4. **Invocação via Chat dos Agentes:**
   Basta digitar `/legoloop-cartoon-shorts <tema>` em qualquer conversa. O agente ativará o pipeline determinístico e entregará o MP4 compilado.
