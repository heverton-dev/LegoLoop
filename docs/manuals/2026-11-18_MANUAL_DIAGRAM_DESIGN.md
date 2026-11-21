# 📘 Manual Operacional: legoloop-diagram-design

> **Data de Emissão:** 2026-11-18  
> **Processo:** `legoloop-diagram-design` (versão `2.6`)  
> **Objetivo:** Criação determinística de diagramas visuais, arquiteturas, fluxogramas e infográficos editoriais auto-contidos em HTML/SVG/PNG com zero serifas, zero emojis e isolamento de micro-vetores.  

---

## 1. Execução Canônica

Para disparar o processo via terminal pelo motor nativo do LegoLoop:

```bash
legoloop exec legoloop-diagram-design
```

Ou execute diretamente o DAG validado:

```bash
legoloop run processos/legoloop-diagram-design/pipeline.dag.json
```

---

## 2. Anatomia do Grafo e Contratos das Etapas

| Etapa | Título | Peça / Script | Depende | Saída | Gate |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `e01` | Inicializar Slug e Anticorpos | `init_diagrama.py` | - | `docs/diagrams/<slug>/manifest.json` | `file_exists('docs/diagrams/<slug>/manifest.json')` |
| `e02` | Ingestão e Quadrantes | `gerar_quadrantes.py` | `e01` | `docs/diagrams/<slug>/quadrantes.json` | `file_exists('docs/diagrams/<slug>/quadrantes.json')` |
| `e03` | Geração SVG e HTML | `gerar_diagrama.py` | `e02` | `docs/diagrams/<slug>/src/<slug>.html` | `file_exists('docs/diagrams/<slug>/src/<slug>.html')` |
| `e04` | Micro-Animação Procedural | `gerar_animacao.py` | `e03` | `docs/diagrams/<slug>/animation/composition.tsx` | `file_exists('docs/diagrams/<slug>/animation/composition.tsx')` |
| `e05` | Gate e Exportação Standalone | `export_svg.py` | `e04` | `docs/diagrams/<slug>/export/<slug>.svg` | `file_exists('docs/diagrams/<slug>/export/<slug>.svg')` |
| `e06` | Fechamento e Anticorpos | `gerar_relatorio.py` | `e05` | `docs/reports/YYYY-MM-DD_REPORT_<slug>.md` | `file_exists('docs/reports/YYYY-MM-DD_REPORT_<slug>.md')` |

---

## 3. Invariantes e Regras Operacionais

1. **Zero Stubs:** Cada nó executa sua respectiva peça real sem simulações ou mocks.
2. **Isolamento de Contexto:** Toda a estrutura de 5 subdiretórios `{src,assets,export,animation,anticorpos}` e o arquivo de governança `anticorpos.json` devem ser mantidos intactos dentro de `docs/diagrams/<slug>/`.
3. **Portões Mecânicos:** O avanço de cada etapa exige confirmação mecânica de seu gate binário (`exit_code == 0` ou `file_exists`).
4. **Paridade Criptográfica:** O DAG de execução deve manter paridade de hash com a especificação formal e os schemas JSON.
5. **Regra Zero Serifas:** Proibição total de fontes serifadas (Georgia, Instrument Serif) e adoção exclusiva de Geist Sans e Geist Mono.
