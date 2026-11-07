# Parte II — Dissecação Forense e Prompts Ocultos

## Decodificando o Não-Dito

O motor forense opera sob o postulado de que o áudio falado é apenas a camada superficial da instrução. A engenharia reversa real precisa extrair o que o autor executou no computador.

### 1. Ingestão Stealth Multi-Plataforma

- **YouTube:** Captura de streams diretas de áudio em PCM 48kHz e metadados via curl/yt-dlp sem disparo de anti-bot.
- **Instagram e TikTok:** Acoplamento com o navegador furtivo **Obscura** e Chrome CDP do sistema operacional, reaproveitando cookies locais sem risco de CAPTCHA ou bloqueios temporários.

### 2. Dissecação de Prompts Explícitos e Ocultos

O dissecador examina a transcrição fonética e a descrição procurando blocos formais de instrução. Quando o autor diz apenas *"aqui eu gerei uma textura futurista"*, o motor analisa os atributos visuais do frame e reconstrói o prompt mestre:
- Parâmetros de iluminação (specular, rim-light, depth).
- Estilo de renderização (translucent glass, backdrop-blur 16px).
- Restrições geométricas (aspect ratio 9:16, 60fps clean loop).

### 3. OCR Forense em Telas Rápidas de Código

Muitos criadores exibem seu terminal ou editor de código por frações de segundo durante transições de cena:
- O módulo de visão multimodal local isola esses frames de alta frequência.
- Executa reconhecimento de caracteres com foco em sintaxe de linha de comando (`ffmpeg`, `styled-components`, `css shaders`).
- Cataloga essas chamadas na lista `HiddenTools`, revelando dependências omitidas no discurso.
