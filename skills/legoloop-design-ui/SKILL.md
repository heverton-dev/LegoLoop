---
name: legoloop-design-ui
description: Use quando precisar executar o processo legoloop-design-ui de UI/UX, motion scroll e 3D web integrando design direction e gates visuais responsivos.
consumes: Artifact:Idea
produces: Artifact:Design
artifact: out/01_planejar.md
gate: file_exists('out/01_planejar.md') && word_count('out/01_planejar.md') >= 30 && !has_pattern('out/01_planejar.md', '\{\{')
suggested_model: sonnet
---

# 🧩 legoloop-design-ui — Design de UI/UX, Motion e Web 3D

> **Processo:** `legoloop-design-ui`  
> **Versão:** 1.0.0  
> **Referência:** `processos/legoloop-design-ui/pipeline.dag.json`  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** Quando for necessário planejar e estruturar interfaces modernas, animações de scroll (motion scroll), elementos 3D web e direção visual para frontend.  
**Quando NÃO usar:** Para processos de pura lógica backend sem superfície de apresentação gráfica.

* **Consumes:** Requisitos de interface e direcionais de design.
* **Produces:** `out/01_planejar.md`.

---

## 2. Estrutura de Execução do Grafo DAG

O processo possui a etapa fundamental de planejamento e direcionamento visual:

| Nó ID | Nome do Nó | Tipo / Ref | Saída / Gate |
| :--- | :--- | :--- | :--- |
| `e01` | Planejar | `skills/frontend-design-direction` | `out/01_planejar.md` |

---

## 3. Procedimento Operacional Passo a Passo

1. **Planejar (`e01`):** Conduzir a formulação do planejamento de interface, tokens visuais, paletas de cores, tipografia responsiva e plano de motion web, emitindo o artefato estruturado `out/01_planejar.md`.
2. Validar o gate determinístico: verificar que `out/01_planejar.md` existe, possui contagem de palavras >= 30 e não contém placeholders (`\{\{`).
3. Apresentar resultado ao usuário e solicitar validação antes de prosseguir com implementações dependentes.

---

## 4. Invariantes & Regras Negativas

- **Design System First:** Proibido styling ad-hoc sem fundamentação em tokens e direção visual.
- **Acessibilidade:** Contrastes mínimos WCAG devem ser respeitados desde a fase de planejamento.
- **Portões Binários:** O artefato de saída deve atender integralmente às condições mecânicas do gate.
