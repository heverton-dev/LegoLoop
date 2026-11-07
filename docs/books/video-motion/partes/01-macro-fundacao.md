# Parte I — Fundação e Princípio Determinístico

## A Lei Fundamental: Vídeo é Código

O pipeline `video-motion` foi concebido sob a lei de conduta inegociável do LegoLoop: **tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de uma IA generativa.**

Na produção audiovisual contemporânea, a maioria dos criadores depende de pacotes fechados (After Effects, Premiere, CapCut) ou de plataformas SaaS baseadas em créditos e assinaturas opacas. Essa dependência cria pontos únicos de falha, custos imprevisíveis e impossibilidade de automação em lote.

No LegoLoop, a premissa é invertida:
- O layout visual é composto em **código declarativo** (CSS3 3D, Web-Canvas e Remotion).
- Os tempos de transição são guiados pela matemática da fala (análise de áudio RMS e timestamps milimétricos do Whisper).
- Os assets visuais secundários são gerados via APIs públicas gratuitas (Google AI Studio) ou sintetizados deterministicamente por shaders procedurais em Canvas.

## Comparativo de Eficiência Arquitetural

| Dimensão                 | Fluxo Tradicional (SaaS / NLE)         | Pipeline LegoLoop (`video-motion`)   |
| :----------------------- | :------------------------------------- | :----------------------------------- |
| **Corte de Respiração**  | Manual em timeline ou app proprietário | Automático (-30dB RMS + 40ms ataque) |
| **Sincronia de Legenda** | Digitação ou exportação com drift      | `words.json` milimétrico sem desvio  |
| **Estilo Visual 3D**     | Plugins pagos de After Effects         | Remotion / CSS3 com blur procedural  |
| **Custo por Minuto**     | Assinaturas mensais recorrentes        | Custo zero de licença e 100% local   |
| **Versionamento**        | Binários pesados (.aep / .prproj)      | Código puro versionado em Git        |
