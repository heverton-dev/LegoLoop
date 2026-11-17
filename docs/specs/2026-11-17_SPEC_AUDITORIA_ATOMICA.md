# 📋 SPEC: legoloop-auditoria-atomica

> **Status:** Proposta para Revisão  
> **Versão:** 1.0.0  
> **Gerada por:** `legoloop fabricar` (projeção do processo)  
> **Template:** TEMPLATE_SPEC.md v2.2.0  

---

## 1. Premissas Técnicas Assumidas (Addy Osmani)
1. Natureza do processo: software.
2. Objetivo: Auditar qualquer alvo com critério binário e emitir veredito com evidências
3. Cada etapa é executada por uma peça do cofre; o id da etapa é o id do nó no grafo.

## 2. Fronteiras de Ação & Regras de Governança (Boundaries — Addy Osmani)
* Sempre: validar o grafo com `legoloop validate` antes de executar.
* Nunca: avançar para a próxima etapa sem o gate aprovado.

## 3. Visão Geral, Arquitetura & Costuras de Teste (Test Seams — Pocock)
As etapas rodam em ordem de dependência; cada uma grava um artefato e passa por um gate.

## 4. Contratos de Dados & Schemas Formais
| Etapa | Quem executa | Produz | Tipo |
| :--- | :--- | :--- | :--- |
| `e01_dissecao` | peça `legoloop-auditoria-atomica` | `out/01_dissecao.md` | markdown |
| `e02_lacunas` | peça `legoloop-auditoria-atomica` | `out/02_lacunas.md` | markdown |
| `e03_tickets` | peça `legoloop-auditoria-atomica` | `out/03_tickets.md` | markdown |
| `e04_refatoracao` | peça `legoloop-auditoria-atomica` | `out/04_refatoracao.md` | markdown |
| `e05_veredito` | peça `legoloop-auditoria-atomica` | `out/05_veredito.md` | markdown |

## 5. Critérios Binários de Aceite por Nó / Ticket (BDD)

### Épico 01 — legoloop-auditoria-atomica

* [ ] `AC-E01.1`: **Etapa `e01_dissecao`: Dissecar o alvo da auditoria**
  * **Dado que:** o artefato `out/01_dissecao.md` ainda não existe
  * **Quando:** a etapa `e01_dissecao` roda com a peça `legoloop-auditoria-atomica`
  * **Então:** o gate `file_exists('out/01_dissecao.md') && word_count('out/01_dissecao.md') >= 30 && !has_pattern('out/01_dissecao.md', '\{\{')` passa
* [ ] `AC-E01.2`: **Etapa `e02_lacunas`: Listar lacunas do alvo**
  * **Dado que:** o artefato `out/02_lacunas.md` ainda não existe
  * **Quando:** a etapa `e02_lacunas` roda com a peça `legoloop-auditoria-atomica`
  * **Então:** o gate `file_exists('out/02_lacunas.md') && word_count('out/02_lacunas.md') >= 30 && !has_pattern('out/02_lacunas.md', '\{\{')` passa
* [ ] `AC-E01.3`: **Etapa `e03_tickets`: Gerar tickets das lacunas**
  * **Dado que:** o artefato `out/03_tickets.md` ainda não existe
  * **Quando:** a etapa `e03_tickets` roda com a peça `legoloop-auditoria-atomica`
  * **Então:** o gate `file_exists('out/03_tickets.md') && word_count('out/03_tickets.md') >= 30 && !has_pattern('out/03_tickets.md', '\{\{')` passa
* [ ] `AC-E01.4`: **Etapa `e04_refatoracao`: Sugerir refatoração do alvo**
  * **Dado que:** o artefato `out/04_refatoracao.md` ainda não existe
  * **Quando:** a etapa `e04_refatoracao` roda com a peça `legoloop-auditoria-atomica`
  * **Então:** o gate `file_exists('out/04_refatoracao.md') && word_count('out/04_refatoracao.md') >= 30 && !has_pattern('out/04_refatoracao.md', '\{\{')` passa
* [ ] `AC-E01.5`: **Etapa `e05_veredito`: Emitir veredito binário com evidências**
  * **Dado que:** o artefato `out/05_veredito.md` ainda não existe
  * **Quando:** a etapa `e05_veredito` roda com a peça `legoloop-auditoria-atomica`
  * **Então:** o gate `file_exists('out/05_veredito.md') && word_count('out/05_veredito.md') >= 30 && !has_pattern('out/05_veredito.md', '\{\{')` passa

## 6. Comandos Operacionais & Critério de Conclusão Global (Addy Osmani)
* Validar: `legoloop validate ~/.aidd/canonical/processos/legoloop-auditoria-atomica/pipeline.dag.json`

🛑 **STOP-GATE 02 DE AUDITORIA HUMANA:**
