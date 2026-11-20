---
name: legoloop-canonizar
description: "Industrialização e canonização universal de peças do LegoLoop: elevação ao Padrão Ouro com legoloop-piece-upgrader e fechamento determinístico em 9 etapas com legoloop-fechar."
consumes: Artifact:Piece
produces: Artifact:CanonicalArtifacts
artifact: out/fechamento_result.json
gate: file_exists('out/fechamento_result.json')
suggested_model: sonnet
---

# 🏭 legoloop-canonizar — Processo Canônico de Industrialização de Peças

> **Autor:** LegoLoop  
> **Versão:** 1.0.0  
> **Composição:** `legoloop-piece-upgrader` + `legoloop-fechar`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** quando uma peça existente (skill, script, hook, mcp, gate ou processo) precisa ser elevada ao Padrão Ouro (Score >= 0.85) e imediatamente canonizada no cofre canônico com manual, livro LegoPress, skill e DAG.  
**Quando NÃO usar:** para conceber uma ideia de negócio do zero sem peça prévia (use `legoloop-concepcao`).

* **Consumes:** a peça-alvo a ser elevada e canonizada.
* **Produces:** a peça canônica no cofre com documentação, livros, skills e DAG validados.

---

## 2. As Duas Camadas do Processo

1. **Camada 1 — Elevação (`n01_upgrade`):** executa o diagnóstico de lacunas (`legoloop upgrade <alvo> --diagnose`), refina via TDD estrito e assegura Score >= 0.85.
2. **Camada 2 — Fechamento (`n02_fechar`):** executa a homologação canônica em 9 etapas (`legoloop fechar <slug>`), compilando livro didático, manual e registrando nos harnesses.
