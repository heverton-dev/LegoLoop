---
name: legoloop-dissecar
description: Use quando precisar executar o processo legoloop-dissecar de inspeção anatômica profunda, engenharia reversa estática e geração da tríade de motion graphics.
consumes: Artifact:Target
produces: Artifact:Report
artifact: ANATOMIA.md
gate: file_exists('ficha_anatomica.json') && file_exists('ANATOMIA.md')
suggested_model: sonnet
---

# 🧩 legoloop-dissecar — Inspeção Anatômica Profunda & Engenharia Reversa

> **Processo:** `legoloop-dissecar`  
> **Versão:** 2.2.0  
> **Referência:** `processos/legoloop-dissecar/processo.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for necessário dissecar a anatomia de um alvo (repositório, código ou Peça), extrair métricas de AST, auditar governança estática e gerar a tríade de motion graphics com relatório anatômico.  
**Quando NÃO usar:** Para processos de criação de código do zero sem um alvo prévio para dissecação.

* **Consumes:** Alvo de análise (parâmetro `var.alvo`).
* **Produces:** `target_resolved.json`, `ast_metrics.json`, `governance_audit.json`, `dag.json`, `motion_graph.json`, `anatomia.html`, `anatomia.svg`, `ficha_anatomica.json` e `ANATOMIA.md`.

---

## 2. Estrutura de Execução do Grafo DAG

O processo é composto por 6 nós de engenharia reversa estática:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `n01_resolver_alvo` | n01_resolver_alvo | `unit_no` | `target_resolved.json` |
| `n02_extrair_ast_metricas` | n02_extrair_ast_metricas | `unit_no` | `ast_metrics.json` |
| `n03_auditar_governanca` | n03_auditar_governanca | `unit_no` | `governance_audit.json` |
| `n04_sintetizar_dag_reverso` | n04_sintetizar_dag_reverso | `unit_no` | `dag.json` |
| `n05_sintetizar_triade_motion` | n05_sintetizar_triade_motion | `unit_no` | `motion_graph.json` |
| `n06_compilar_relatorio_final` | n06_compilar_relatorio_final | `unit_no` | `ficha_anatomica.json` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Resolver Alvo (`n01_resolver_alvo`):** Localizar e validar a integridade física do alvo a ser dissecado, gerando `target_resolved.json`.
2. **Extrair AST e Métricas (`n02_extrair_ast_metricas`):** Varrer a árvore sintática, calcular complexidade ciclomática, contagem de linhas e dependências em `ast_metrics.json`.
3. **Auditar Governança (`n03_auditar_governanca`):** Avaliar aderência às leis invioláveis, padrões de código e convenções de nomenclatura em `governance_audit.json`.
4. **Sintetizar DAG Reverso (`n04_sintetizar_dag_reverso`):** Mapear as dependências e o Processo de dados em um grafo acíclico dirigido formal `dag.json`.
5. **Sintetizar Tríade Motion (`n05_sintetizar_triade_motion`):** Gerar visualizações gráficas e diagramas da arquitetura (`motion_graph.json`, `anatomia.html` e `anatomia.svg`).
6. **Compilar Relatório Final (`n06_compilar_relatorio_final`):** Consolidar a ficha anatômica estruturada em `ficha_anatomica.json` e a documentação completa em `ANATOMIA.md`.

---

## 4. Invariantes & Regras Negativas

- **Não Modificar o Alvo:** A dissecação é estritamente de leitura analítica; o alvo nunca deve ser alterado em disco.
- **Evidências Reais:** Todas as métricas de AST e auditoria devem ser extraídas por análise real de arquivos.
- **Portões Mecânicos:** Cada nó só avança após o gate do arquivo esperado ter sido gerado com sucesso.
