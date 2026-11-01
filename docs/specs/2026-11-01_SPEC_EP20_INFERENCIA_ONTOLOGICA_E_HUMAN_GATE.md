# 📋 SPEC_EP20: Inferência de Contrato Ontológico para Ativos do SO & Peça `human-gate`

> **Status:** Proposta para Revisão (STOP-GATE HUMANO 02)  
> **Data:** 2026-11-01  
> **Épico Alvo:** ÉPICO 20 (Ontologia de Ativos do SO, Stop-Gate Humano & Mega Pipeline Template)  
> **Referência:** `docs/intents/2026-11-01_INTENT_ONTOLOGIA_ASSETS_SO_E_MEGA_PIPELINE_TEMPLATES.md`

---

## 1. Visão Geral e Arquitetura

Esta especificação resolve a integração de primeira classe de qualquer ativo do Sistema Operacional (`SKILL.md`, scripts, MCPs) na ontologia universal do LegoLoop:
1. **Parser e Heurística de Contrato no Scanner (`internal/assets/skills.go`):**
   * Extração de campos ontológicos do frontmatter YAML (`consumes`, `produces`, `artifact`, `gate`, `suggested_model`).
   * Fallback inteligente de inferência heurística para convenções canônicas (`aidd-*`, `interview`, `spec`, `tdd`, `homologue`, `rea`).
   * Metadados populados no `Item.Meta` para consumo transparente da malha semântica.
2. **Promoção de Itens do Inventário a Blocos Ontológicos (`internal/studio/pieces.go`):**
   * Quando uma peça do inventário possui metadados ontológicos, ela é promovida a um `graph.Block` funcional com suas portas conectáveis e validadas por `CanConnect`.
3. **Nova Peça de Controle: Stop-Gate Humano (`human-gate`):**
   * Bloco registrado em `ControlBlocks()` (`internal/graph/ontology.go`) com `Kind: "flow_human_gate"`.
   * Interrompe o fluxo para validação humana explícita, gerando o artefato de aprovação `AUDIT_APPROVAL.json`.
4. **Template Canônico do Mega Pipeline (`internal/studio/templates.go`):**
   * Registro do template `aidd-legoloop-mega-pipeline` contendo os 5 sub-loops canônicos encadeados:
     * `L1: Concepção (interview / aidd-grill)`
     * `L2: REA & Design (impeccable)`
     * `L3: Tríade de Planejamento (aidd-spec)`
     * `L4: Construção TDD (aidd-tdd)`
     * `L5: Homologação E2E (audit / aidd-homologue)`

---

## 2. Critérios de Aceitação Binários (Gates Mecânicos)

* [ ] **CA-20.1 (Parser Frontmatter Ontológico):** `parseSkill` extrai corretamente `consumes`, `produces`, `artifact`, `gate` quando declarados no YAML frontmatter de um `SKILL.md`.
* [ ] **CA-20.2 (Heurística Canônica Fallback):** Skills canônicas que não possuem frontmatter explícito (ex: `aidd-legoloop-mega-pipeline`, `aidd-tdd`, `aidd-spec`) recebem inferência ontológica automática de portas baseada em seu identificador.
* [ ] **CA-20.3 (Promoção a `graph.Block`):** Em `resolvePiece` / `skillPiece` (`internal/studio/pieces.go`), peças que contêm metadados ontológicos geram um `nodeInfo` com `Block != nil` e portas de I/O populadas.
* [ ] **CA-20.4 (Stop-Gate Humano em `ontology.go`):** `ControlBlocks()` inclui a peça `human-gate` com `Kind: "flow_human_gate"`, aceitando qualquer artefato ontológico e produzindo aprovação auditada.
* [ ] **CA-20.5 (Template Mega Pipeline no Studio):** `BuiltinTemplates()` em `internal/studio/templates.go` disponibiliza a esteira `mega-pipeline` com os 5 nós canônicos e suas conexões direcionadas.
* [ ] **CA-20.6 (Zero Regressões & Testes 100% Verdes):** Testes unitários novos e existentes (`go test ./...`) passam 100% sem erros, linter limpo e binário recompilado.
