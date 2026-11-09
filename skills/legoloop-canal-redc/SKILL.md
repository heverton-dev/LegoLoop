---
name: legoloop-canal-redc
description: Use quando precisar executar Engenharia Reversa Profunda de Canais (REDC), analisar 100 vídeos virais, sintetizar dossiê de DNA visual e comportamental (ganchos, tom, ritmo, design/motion tokens), clonar métodos em pipelines ou gerar vídeos com Brand Swap e marca do usuário. NÃO use para edição manual arbitrária de timeline sem análise de canal ou criação de livros didáticos.
compatibility: Windows/Linux, PowerShell 5.1+, Go 1.22+
allowed-tools: read_file write_file edit_file run_command
consumes: Artifact:VideoMetrics
produces: Artifact:ChannelDossier
artifact: docs/channels/canal_sample/repc/01_dossie_dna.json
gate: file_exists('docs/channels/canal_sample/repc/01_dossie_dna.json')
suggested_model: code-pro
metadata:
  version: "2.2.0"
  tier: "ouro"
  category: "reverse"
---

# 🧩 legoloop-canal-redc — Mega-Pipeline de Engenharia Reversa Profunda de Canais e Clonagem

> **Autor:** LegoLoop  
> **Versão:** 2.2.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.2.0  
> **Taxonomia:** Função nativa da plataforma LegoLoop para análise sistemática e clonagem de canais virais.

---

## 1. Contexto & Contratos de I/O

**Quando usar:** quando o usuário fornecer uma pasta com métricas e transcrições de vídeos de um canal social (YouTube, Instagram, TikTok, LinkedIn) e desejar dissecar o estilo, extrair o Dossiê DNA (tom de voz, design tokens, motion tokens, ganchos virais), gerar os pipelines dos métodos ensinados ou gerar novos vídeos estruturais com a sua própria marca (Brand Swap).  
**Quando NÃO usar:** para raspar redes sociais online diretamente sem dados locais pré-coletados ou para compilar livros técnicos avulsos (use `legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - Diretório do canal contendo `metricas.json` e transcrições estruturadas.
* **Produces (Saídas Mecânicas):**
  - `docs/channels/<canal>/repc/01_dossie_dna.json`
  - `docs/channels/<canal>/repc/02_design_tokens.json`
  - `docs/channels/<canal>/repc/03_motion_tokens.json`
  - `docs/channels/<canal>/repc/04_matriz_ganchos.md`
  - `docs/channels/<canal>/pipelines/<slug-metodo>/pipeline.dag.json`
  - `docs/channels/<canal>/outputs/roteiro_<slug>.md`
* **Target Files (Escopo de escrita permitido / Blast Radius):**
  - `docs/channels/<canal>/**`

```json
{
  "consumes": ["metricas.json"],
  "produces": [
    "docs/channels/<canal>/repc/01_dossie_dna.json",
    "docs/channels/<canal>/outputs/roteiro_*.md"
  ],
  "gate": "file_exists('docs/channels/<canal>/repc/01_dossie_dna.json')"
}
```

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
1. **NUNCA** analise 100 vídeos brutos para extração de estilo de corte: a amostragem inteligente prescreve os Top 15 vídeos de pico para extrair o DNA, evitando estourar orçamento de tokens.
2. **NUNCA** sobrescreva tokens de marca do usuário pelo estilo do canal clonado na fase de geração: o canal clonado fornece o esqueleto e retenção; a marca do usuário fornece a identidade visual.
3. **NUNCA** grave arquivos fora da árvore `docs/channels/<canal>/`.
4. **NUNCA** execute com mocks ou stubs provisórios em produção.

---

## 3. Comandos de Operação CLI

```bash
# Executa a Engenharia Reversa Profunda do Canal
legoloop canal redc testdata/fixtures/canal_sample --canal meu-canal --json

# Gera um novo vídeo autoral aplicando Brand Swap no estilo do canal
legoloop canal gerar meu-canal --tema "Como dominar agentes de IA" --nome-marca "MinhaMarca" --cor-primaria "#3B82F6" --json

# Valida o grafo DAG do pipeline
legoloop validate loops/legoloop-canal-redc/pipeline.dag.json
```

---

## 4. Critério de Parada (Stopping Checklist)

- [ ] Árvore `docs/channels/<canal>/` gerada contendo `repc/`, `pipelines/` e `outputs/`.
- [ ] Dossiê DNA sintetizado com tom de voz, ritmo e tokens visuais/motion.
- [ ] Grafo validado com exit code 0 por `legoloop validate`.
