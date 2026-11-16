# 📋 SPEC: legoloop-swarm

> **Status:** Proposta para Revisão  
> **Versão:** 1.5.0  
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
| `e01` | peça `legoloop-loop-map` | `out/01_niveis.md` | markdown |
| `e02` | peça `legoloop-matriz` | `out/02_workers.md` | markdown |
| `e03` | peça `legoloop-swarm-cli` | `out/03_spawn.md` | markdown |
| `e04` | peça `legoloop-swarm-cli` | `out/04_coordenacao.md` | markdown |
| `e05` | peça `legoloop-swarm-cli` | `out/05_encerramento.md` | markdown |

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — legoloop-swarm

* [ ] `AC-E01.1`: **Etapa `e01`: Validar o DAG e calcular os niveis topologicos**
  * **Dado que:** o artefato `out/01_niveis.md` ainda não existe
  * **Quando:** a etapa `e01` roda com a peça `legoloop-loop-map`
  * **Então:** o gate `file_exists('out/01_niveis.md') && has_pattern('out/01_niveis.md', 'nivel 0') && !has_pattern('out/01_niveis.md', '\{\{')` passa
* [ ] `AC-E01.2`: **Etapa `e02`: Registrar os workers do nivel com papel da Matriz e worktree isolado**
  * **Dado que:** o artefato `out/02_workers.md` ainda não existe
  * **Quando:** a etapa `e02` roda com a peça `legoloop-matriz`
  * **Então:** o gate `file_exists('out/02_workers.md') && has_pattern('out/02_workers.md', 'worktree') && !has_pattern('out/02_workers.md', '\{\{')` passa
* [ ] `AC-E01.3`: **Etapa `e03`: Subir o processo real de cada worker**
  * **Dado que:** o artefato `out/03_spawn.md` ainda não existe
  * **Quando:** a etapa `e03` roda com a peça `legoloop-swarm-cli`
  * **Então:** o gate `file_exists('out/03_spawn.md') && has_pattern('out/03_spawn.md', 'pid') && !has_pattern('out/03_spawn.md', '\{\{')` passa
* [ ] `AC-E01.4`: **Etapa `e04`: Coordenar o nivel pelo quadro em disco**
  * **Dado que:** o artefato `out/04_coordenacao.md` ainda não existe
  * **Quando:** a etapa `e04` roda com a peça `legoloop-swarm-cli`
  * **Então:** o gate `file_exists('out/04_coordenacao.md') && has_pattern('out/04_coordenacao.md', 'concluid') && !has_pattern('out/04_coordenacao.md', '\{\{')` passa
* [ ] `AC-E01.5`: **Etapa `e05`: Encerrar os workers e conferir a limpeza**
  * **Dado que:** o artefato `out/05_encerramento.md` ainda não existe
  * **Quando:** a etapa `e05` roda com a peça `legoloop-swarm-cli`
  * **Então:** o gate `file_exists('out/05_encerramento.md') && has_pattern('out/05_encerramento.md', 'killed') && !has_pattern('out/05_encerramento.md', '\{\{')` passa

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Validar: `legoloop validate ~/.aidd/canonical/loops/legoloop-swarm/pipeline.dag.json`

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
