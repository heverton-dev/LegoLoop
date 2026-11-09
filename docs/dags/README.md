# 🗺️ Diretório Canônico de Grafos Executáveis (DAGs ISO)

> **Padrão:** Quinteto de Governança ISO  
> **Nomenclatura:** `YYYY-MM-DD_DAG_<SLUG>.json`  
> **Schema:** `https://legoloop.dev/schemas/dag.v1.json`  

Este diretório contém os grafos acíclicos dirigidos (DAGs) canônicos versionados por data.  
Cada DAG é o espelho executável 1:1 de sua respectiva rota (`docs/routes/YYYY-MM-DD_ROUTE_<SLUG>.json`), projetado de forma atômica para `loops/<slug>/pipeline.dag.json` e auditado via SHA-256 no comando `legoloop fechar`.
