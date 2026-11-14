# 📋 SPEC: legoloop-design-ui

> **Status:** Proposta para Revisão  
> **Versão:** 1.0.0  
> **Gerada por:** `legoloop fabricar` (projeção do processo)  
> **Template:** TEMPLATE_SPEC.md v2.2.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Natureza do processo: midia.
2. Objetivo: Processo especializado de UI/UX, motion scroll e 3D web integrando aidd-uiux, impeccable, emilkowalski-skills, aidd-dataviz, aidd-linter-design, web-perf e gates responsivos
3. Cada etapa é executada por uma peça do cofre; o id da etapa é o id do nó no grafo.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: validar o grafo com `legoloop validate` antes de executar.
* Nunca: avançar para a próxima etapa sem o gate aprovado.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
As etapas rodam em ordem de dependência; cada uma grava um artefato e passa por um gate.

## 4. Contratos de Dados & Schemas Formais
| Etapa | Quem executa | Produz | Tipo |
| :--- | :--- | :--- | :--- |
| `e01` | peça `frontend-design-direction` | `out/01_planejar.md` | markdown |

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — legoloop-design-ui

* [ ] `AC-E01.1`: **Etapa `e01`: Planejar**
  * **Dado que:** o artefato `out/01_planejar.md` ainda não existe
  * **Quando:** a etapa `e01` roda com a peça `frontend-design-direction`
  * **Então:** o gate `file_exists('out/01_planejar.md') && word_count('out/01_planejar.md') >= 30 && !has_pattern('out/01_planejar.md', '\{\{')` passa

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Validar: `legoloop validate ~/.aidd/canonical/loops/legoloop-design-ui.dag.json`

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
