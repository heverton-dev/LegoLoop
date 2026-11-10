---
name: legoloop-video-motion
description: Use quando o usuário quiser EXECUTAR O PROCESSO DETERMINÍSTICO DE MOTION VIDEO (edição automatizada, silêncios +40ms, punch RMS, cartões de vidro 3D, tipografia Anton cinética e broadcast bundle multi-aspecto), via /legoloop-video-motion. NÃO use para dissecar vídeos externos ou criar livros didáticos avulsos.
consumes: Artifact:Media
produces: Artifact:BroadcastBundle
artifact: dist/<job_id>/bundle_manifest.json
gate: exit_code == 0
suggested_model: code-fast
---

# 🎬 legoloop-video-motion — Processo Canônico de Motion Video e Edição por Código

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário fornece um arquivo de áudio/vídeo gravado ou narração gerada por voz sintética (TTS) e deseja realizar edição determinística completa por código (Video as Code): remoção cirúrgica de silêncios (-30dB) com ataque pré-fala (+40ms), alinhamento fonético milimétrico (`words.json`), detecção de ênfase vocal via RMS (`is_punch_word`), cartões 3D translúcidos com vidro fosco procedural (`backdrop-filter: blur(16px)`), tipografia cinética Anton, auto-framing multi-aspecto (9:16, 16:9, 1:1), inspeção rápida de QA (`qa_sheet.jpg`) e exportação do pacote broadcast final.  
**Quando NÃO usar:** Fazer engenharia reversa de vídeos externos da web (use `/legoloop-video-rea`) ou compilar livros didáticos isolados (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - Arquivo de mídia de entrada (`input.mp4` ou locução `audio.wav`).
* **Produces (Saídas Mecânicas):**
  - Tabela fonética com timestamps e punch words (`scratch/videomotion/punched_words.json`)
  - Manifesto de assets e shaders procedurais com cache SHA-256 (`scratch/videomotion/assets_manifest.json`)
  - Folha de contato QA em 2 segundos (`scratch/videomotion/qa_sheet.jpg`)
  - Pacote de entrega broadcast (`dist/<job_id>/`) com MP4, legendas `.srt`/`.vtt`, copy para redes (`POST_COPY.md`) e linha de tempo DaVinci (`timeline.fcpxml`).

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA** cortar pausas sem adicionar os **+40ms de ataque pré-fala** (evita mutilar fonemas oclusivos iniciais).
> 2. **NUNCA** depender de renderizadores de nuvem pagos quando o gerador procedural local de shaders ou Google AI Studio gratuito estiver disponível.
> 3. **NUNCA** realizar a exportação final sem passar pelo Stop-Gate visual da folha de contato (`qa_sheet.jpg`).
> 4. **NUNCA** desviar da normalização sonora padrão de broadcast em **-14 LUFS** (EBU R128).

---

## 3. Protocolo de Operação Passo a Passo

1. **Execução Completa via CLI:**
   ```bash
   legoloop video-motion run <input_media> [--aspect 9:16|16:9|1:1]
   ```
2. **Execução Modular Isolada:**
   - Detectar silêncios e fonemas: `legoloop video-motion silence --input <file> --out scratch/words.json`
   - Detectar punch words por RMS: `legoloop video-motion punch --words scratch/words.json --out scratch/punched.json`
   - Síntese de assets com cache SHA-256: `legoloop video-motion assets --words scratch/punched.json --out scratch/assets.json`
   - Gerar QA Sheet (12 frames em 2s): `legoloop video-motion preview --assets scratch/assets.json --out scratch/qa_sheet.jpg`
   - Empacotar bundle broadcast: `legoloop video-motion bundle --preview scratch/qa_sheet.jpg --out dist/`
3. **Execução Orquestrada via DAG:**
   ```bash
   legoloop run processos/video-motion/processo.dag.json --allow-scripts --yes
   ```
