# Parte I: A Arquitetura da Fundição Autônoma

## Capítulo 1: O Paradoxo da Escala Cognitiva

A engenharia de software tradicional dependeu historicamente da coordenação de grandes equipes humanas. No desenvolvimento assistido por agentes inteligentes, esse modelo entra em colapso devido à dispersão de contexto e alucinações não monitoradas. O LegoLoop aborda esse problema através do isolamento determinístico.

```
+-------------------------------------------------------------------------+
|                  FUNDIÇÃO AUTÔNOMA LEGOLOOP (NÍVEL N10)                 |
+-------------------------------------------------------------------------+
|  [INTENT] ---> [SPEC] ---> [PLAN] ---> [ROUTE] ---> [BUILD TDD]         |
|      |            |           |           |              |              |
|   STANAG        EARS        FMEA       Orchestrix    Red/Blue Team      |
+-------------------------------------------------------------------------+
```

### As 7 Macro-Camadas do Meta-Sistema

1. **Camada 1 — Superfície de Comando (CLI/TUI):** Despacho estrito de 25 verbos em Go nativo, sem dependências externas.
2. **Camada 2 — Motor de Governança da Tríade:** Imutabilidade criptográfica assinada com hashes SHA-256.
3. **Camada 3 — Orquestração Hermética (Orchestrix):** Execução estéril por bloco, onde cada nó consome apenas o artefato pai e encerra seu ciclo de vida.
4. **Camada 4 — Manufatura de Peças:** Síntese de código com Red-Green-Refactor e validação de `funclen <= 35`.
5. **Camada 5 — Soberania e Cofre Canônico:** Inventário local contra dependências externas e deduplicação estrita.
6. **Camada 6 — Fluxo de Estado e DAGs:** Topologias acíclicas sem dependências circulares.
7. **Camada 7 — Domínio e Pedagogia:** Ingestão forense de conhecimento e compilação de livros didáticos em Typst.
