# Parte III — Meso: Curadoria Soberana de Assets

## 3.1 Substituição de Serviços Proprietários

O criador original utiliza uma pilha comercial que inclui Google Flow, ElevenLabs e CapCut. A engenharia reversa do LegoLoop substitui cada nó por equivalentes Tier Ouro:

1. **Geração Visual (Substituindo Google Flow):** Orquestração local via ComfyUI com modelo FLUX ou geração procedural via Remotion/HTML Canvas.
2. **Voz Sintética (Substituindo ElevenLabs):** Piper TTS ou WhisperX com modelos neurais locais de alta fidelidade sem custo por caractere.
3. **Edição e Montagem (Substituindo CapCut):** Remotion CLI e FFmpeg para mesclagem determinística de áudio e vídeo em 1080p60.

#painel-teoria("SOBERANIA TOTAL", [
  A esteira opera 100% desconectada da internet, garantindo que projetos confidenciais de clientes nunca vazem para servidores em nuvem.
])
