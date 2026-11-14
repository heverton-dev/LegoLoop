# 📋 SPEC: legoloop-swarm

> **Status:** Proposta para Revisão  
> **Versão:** 1.0.0  
> **Gerada por:** `legoloop fabricar` (projeção do processo)  
> **Template:** TEMPLATE_SPEC.md v2.2.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Natureza do processo: qualquer.
2. Objetivo: Orquestrar execucao paralela de pipeline.dag.json em swarm por nivel topologico com roteamento de modelos pela Matriz de Papeis
3. Cada etapa é executada por uma peça do cofre; o id da etapa é o id do nó no grafo.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: validar o grafo com `legoloop validate` antes de executar.
* Nunca: avançar para a próxima etapa sem o gate aprovado.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
As etapas rodam em ordem de dependência; cada uma grava um artefato e passa por um gate.

## 4. Contratos de Dados & Schemas Formais
| Etapa | Quem executa | Produz | Tipo |
| :--- | :--- | :--- | :--- |
| `e01` | peça `legoloop-piece-upgrader` | `out/01_etapa.md` | markdown |
| `e02` | peça `aidd-piece-upgrader` | `out/02_etapa.md` | markdown |
| `e03` | peça `legoloop-mega-pipeline` | `out/03_etapa.md` | markdown |
| `e04` | peça `aidd-legoloop-mega-pipeline` | `out/04_etapa.md` | markdown |
| `e05` | peça `legoloop-video-master` | `out/05_etapa.md` | markdown |

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — legoloop-swarm

* [ ] `AC-E01.1`: **Etapa `e01`: Etapa 1**
  * **Dado que:** o artefato `out/01_etapa.md` ainda não existe
  * **Quando:** a etapa `e01` roda com a peça `legoloop-piece-upgrader`
  * **Então:** o gate `file_exists('out/01_etapa.md') && word_count('out/01_etapa.md') >= 30 && !has_pattern('out/01_etapa.md', '\{\{')` passa
* [ ] `AC-E01.2`: **Etapa `e02`: Etapa 2**
  * **Dado que:** o artefato `out/02_etapa.md` ainda não existe
  * **Quando:** a etapa `e02` roda com a peça `aidd-piece-upgrader`
  * **Então:** o gate `file_exists('out/02_etapa.md') && word_count('out/02_etapa.md') >= 30 && !has_pattern('out/02_etapa.md', '\{\{')` passa
* [ ] `AC-E01.3`: **Etapa `e03`: Etapa 3**
  * **Dado que:** o artefato `out/03_etapa.md` ainda não existe
  * **Quando:** a etapa `e03` roda com a peça `legoloop-mega-pipeline`
  * **Então:** o gate `file_exists('out/03_etapa.md') && word_count('out/03_etapa.md') >= 30 && !has_pattern('out/03_etapa.md', '\{\{')` passa
* [ ] `AC-E01.4`: **Etapa `e04`: Etapa 4**
  * **Dado que:** o artefato `out/04_etapa.md` ainda não existe
  * **Quando:** a etapa `e04` roda com a peça `aidd-legoloop-mega-pipeline`
  * **Então:** o gate `file_exists('out/04_etapa.md') && word_count('out/04_etapa.md') >= 30 && !has_pattern('out/04_etapa.md', '\{\{')` passa
* [ ] `AC-E01.5`: **Etapa `e05`: Etapa 5**
  * **Dado que:** o artefato `out/05_etapa.md` ainda não existe
  * **Quando:** a etapa `e05` roda com a peça `legoloop-video-master`
  * **Então:** o gate `file_exists('out/05_etapa.md') && word_count('out/05_etapa.md') >= 30 && !has_pattern('out/05_etapa.md', '\{\{')` passa

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Validar: `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm.dag.json`

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
