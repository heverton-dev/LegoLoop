# Capítulo 3: Fábrica de Micro-Pipelines de Métodos

## 3.1 Da Análise de Conteúdo para Processos Executáveis

Quando um criador ensina "como fazer X em 3 passos", ele está ensinando um algoritmo ou esteira de trabalho. O motor REDC decompõe essa metodologia em:
- Um manifesto declarativo `.loop.json`.
- Um grafo acíclico executável `pipeline.dag.json`.

## 3.2 Desduplicação por Clusterização

Criadores frequentemente publicam múltiplos vídeos sobre o mesmo método com títulos chamativos diferentes. O algoritmo `ClusterMethods` detecta a raiz metodológica e agrupa os vídeos sob o mesmo slug, evitando a criação de pipelines duplicados e garantindo um catálogo limpo em `docs/channels/<canal>/pipelines/`.
