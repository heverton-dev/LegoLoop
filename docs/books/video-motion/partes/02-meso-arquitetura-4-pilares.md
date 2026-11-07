# Parte II — A Arquitetura dos Quatro Pilares

O pipeline é estruturado em quatro motores independentes, integrados por um grafo acíclico dirigido (DAG):


![Diagrama de Arquitetura](/img/diagrama_436e5a57.png){width=85%}


## 1. Relógio Mestre Fonético (Speech-Sync Engine)

A linha do tempo do vídeo não é definida por segundos arbitrários, mas pela cadência real da locução:
1. **Corte Cirúrgico de Pausas:** Todo silêncio abaixo do threshold de -30dB com duração superior a 400ms é removido, preservando estritamente **+40ms de ataque pré-fala** para que nenhuma consoante oclusiva seja deglutida.
2. **Alinhamento Fonético Milimétrico:** Geração da matriz `words.json`, contendo cada vocábulo com seus marcadores de `start_ms` e `end_ms`.
3. **Detecção de Ênfase Vocal:** O motor calcula a energia RMS por palavra. Vocábulos com amplitude 40% superior à média do segmento recebem a flag `is_punch_word = true`.
4. **Normalização Broadcast EBU R128:** Ajuste de ganho para teto de transmissão em **-14 LUFS** com true peak a -1.0 dBFS.

## 2. Compositor Visual Reativo (Remotion & CSS 3D)

Em vez de renderizar frames estáticos pesados, o visual é gerado via DOM reativo acelerado por hardware:
- **Cartões 3D de Vidro:** Superfícies translúcidas com `backdrop-filter: blur(16px)`, bordas sutis com iluminação especular e transição suave de push-in.
- **Tipografia Cinética Anton:** Exibição com destaque cromático instantâneo sobre as palavras de ênfase identificadas na etapa fonética.
- **Auto-Framing Adaptativo:** Recorte inteligente mantendo o sujeito no terço superior de segurança.

## 3. Pipeline de Assets Determinístico com Cache Criptográfico

- **Chave de Cache SHA-256:** Cada solicitação de imagem secundária gera uma assinatura criptográfica a partir do prompt e dimensões. Se já existe no disco, o custo de rede é zero.
- **Bridge com Fallback Local:** Conexão nativa com a API do Google AI Studio. Em caso de indisponibilidade de rede ou ausência de chave, o gerador procedural sintetiza texturas abstratas locais sem interromper o fluxo.

## 4. Renderizador Multi-Aspecto e Bundle Broadcast

- **Exportação Paralela:** Suporte nativo aos formatos `9:16` (Vertical), `16:9` (Widescreen) e `1:1` (Quadrado).
- **Entrega Broadcast Completa:** Geração simultânea do pacote em `dist/<job_id>/` contendo `.mp4`, legendas `.srt` e `.vtt`, arquivo de timecode `fcpxml` (v1.9) e copy pronta para redes em `POST_COPY.md`.
