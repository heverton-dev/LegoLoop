# Capítulo 1: Visão Geral e Arquitetura do REDC

## 1.1 O Desafio da Engenharia Reversa de Canais

Canais que viralizam em plataformas como YouTube, TikTok, Instagram e LinkedIn operam sobre matrizes algorítmicas claras:
- **Retenção Inicial (Hooks):** Primeiros 3 segundos que interrompem o scroll e forçam a continuidade da visualização.
- **Cadência Dinâmica:** Cortes e transições em janelas de 1.8s a 3.0s, impedindo que o cérebro se acostume com um enquadramento estático.
- **Design & Motion Tokens:** Paletas com alto contraste, tipografia pesada e cartões de vidro 3D.
- **Estruturação de Métodos:** Cada vídeo entrega um passo a passo objetivo (pipeline de conhecimento).

## 1.2 A Arquitetura Autocontida em `docs/channels/<canal>/`

O LegoLoop isola cada canal em uma árvore própria:
- `repc/`: Dossiê do DNA (tom, design tokens, motion tokens, matriz de ganchos).
- `pipelines/`: Micro-pipelines canônicos derivados dos métodos ensinados.
- `outputs/`: Roteiros estruturais e vídeos gerados com a marca do usuário.
