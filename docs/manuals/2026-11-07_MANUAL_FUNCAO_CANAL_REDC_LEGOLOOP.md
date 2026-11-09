# 📘 MANUAL OPERACIONAL: FUNÇÃO CANAL REDC NO LEGOLOOP

> **Data:** 2026-11-07  
> **Autoridade:** Ecossistema LegoLoop  
> **Status:** Canônico  
> **Escopo:** Mega-Pipeline de Engenharia Reversa Profunda de Canais (REDC) e Clonagem Metodológica/Estrutural com Brand Swap  
> **Fontes de verdade:** `loops/legoloop-canal-redc/pipeline.dag.json` · `skills/legoloop-canal-redc/SKILL.md` · `internal/redc`

---

## 1. Visão Geral

A função **`legoloop-canal-redc`** automatiza a Engenharia Reversa Profunda de Canais (REDC) e a clonagem metodológica/estética de canais sociais (YouTube, Instagram, TikTok, LinkedIn) em 4 motores autocontidos:

1. **Ingestão & Catalogação dos 100 Vídeos Virais:** Varre métricas brutas (`metricas.json`), calcula score de tração e viralidade (`views * 0.7 + engajamento * 0.3`) e clusteriza métodos sem redundâncias de processos.
2. **Motor Deep REDC (Amostragem Ouro - Top 15 Vídeos):** Dissecar os 15 vídeos virais de pico para extrair tom de voz, ritmo de cortes (ex.: 2.5s), design tokens (paleta, tipografia), motion tokens (transições, cartões 3D) e matriz de ganchos dos primeiros 3s na pasta `docs/channels/<canal>/repc/`.
3. **Fábrica de Micro-Pipelines de Métodos:** Mapeia cada metodologia única ensinada nos vídeos em micro-pipelines canônicos executáveis (`pipeline.dag.json`) em `docs/channels/<canal>/pipelines/`.
4. **Gerador Estrutural com Brand Swap:** Pipeline paramétrico que aceita qualquer tema autoral e injeta os tokens da marca do usuário (cores, fontes, tom institucional), preservando a cadência algorítmica comprovada em `docs/channels/<canal>/outputs/`.

---

## 2. Estrutura de Arquivos Gerada (`docs/channels/<canal>/`)

```
docs/channels/<slug>/
├── repc/
│   ├── 01_dossie_dna.json          # Tom de voz, arquétipos e regras de copy
│   ├── 02_design_tokens.json       # Paleta de cores, tipografia, espaçamentos
│   ├── 03_motion_tokens.json       # Transições, velocidade de corte, zoom, cards 3D
│   └── 04_matriz_ganchos.md        # Catálogo das fórmulas de retenção comprovadas
├── pipelines/                      # Pipelines canônicos dos métodos identificados
│   └── <metodo-slug>/
│       └── pipeline.dag.json
└── outputs/                        # Roteiros gerados e vídeos renderizados com a marca do usuário
    └── roteiro_<tema>_<timestamp>.md
```

---

## 3. Comandos e Uso

### 3.1 Pelo Chat do Agente
* `/legoloop-canal-redc`

### 3.2 Pela CLI Mecânica
```bash
# 1. Executa a análise profunda do canal a partir dos dados de vídeo
legoloop canal redc <origem-dir> --canal <slug-do-canal> --json

# 2. Gera novo vídeo no estilo do canal aplicando a marca do usuário (Brand Swap)
legoloop canal gerar <slug-do-canal> --tema "Como dominar IA" --nome-marca "MinhaMarca" --cor-primaria "#3B82F6" --json

# 3. Validação do grafo DAG do canal
legoloop validate loops/legoloop-canal-redc/pipeline.dag.json
```
