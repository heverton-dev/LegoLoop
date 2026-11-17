# 🗺️ PLAN: Inferência de Contrato Ontológico de Ativos do SO & Peça `human-gate` (ÉPICO 20)

> **Status:** Proposta para Revisão (STOP-GATE HUMANO 02)  
> **Data:** 2026-11-01  
> **Especificação Base:** `docs/specs/2026-11-01_SPEC_EP20_INFERENCIA_ONTOLOGICA_E_HUMAN_GATE.md`  
> **Metodologia:** TDD Estrito (Red ➔ Green ➔ Refactor) + Silent Execution

---

## 🏛️ ÉPICO 20: Ontologia de Ativos do SO, Stop-Gate Humano & Mega Pipeline Template

> **Objetivo:** Transformar qualquer ativo do SO em bloco ontológico funcional e fornecer a peça de auditoria humana e o template Mega Pipeline.  
> **Modelo Recomendado:** Claude 3.7 Sonnet / Gemini 2.0 Pro (`--effort high`)

### 🔹 Micro-ticket TICK-2001: Parser Frontmatter e Heurística Ontológica em `skills.go`
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/assets/skills.go`, `internal/assets/skills_test.go`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Teste unitário em `skills_test.go` verificando extração de `consumes`, `produces`, `artifact`, `gate` via frontmatter e via inferência heurística para `aidd-legoloop-mega-pipeline`.
  * Implementação de `parseFrontmatterOntology` e `inferOntologyFallback` preenchendo `Item.Meta`.

### 🔹 Micro-ticket TICK-2002: Promoção Ontológica de Peças em `pieces.go`
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/studio/pieces.go`, `internal/studio/pieces_test.go`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Teste unitário em `pieces_test.go` provando que uma skill do inventário com metadados ontológicos gera `nodeInfo` com `Block != nil` e portas compatíveis com `CanConnect`.
  * Construção dinâmica de `graph.Block` para skills inventariadas.

### 🔹 Micro-ticket TICK-2003: Peça `human-gate` (Stop-Gate Humano) em `ontology.go`
* **Tipo:** Micro-ticket
* **Arquivo Alvo:** `internal/graph/ontology.go`, `internal/graph/graph_test.go`
* **Rota:** Gemini 2.0 Flash (`--effort medium`)
* **Entregáveis & Testes TDD:**
  * Adicionar `human-gate` em `ControlBlocks()` com suporte universal a entradas e gate de aprovação.
  * Teste unitário validando conexões de entrada com múltiplos tipos de artefatos.

### 🔸 Nano-ticket TICK-2004: Template Canônico Mega Pipeline em `templates.go`
* **Tipo:** Nano-ticket
* **Arquivo Alvo:** `internal/studio/templates.go`, `internal/studio/templates_test.go`
* **Rota:** Gemini 2.0 Flash-Lite (`--effort low`)
* **Entregáveis & Testes TDD:**
  * Registrar o template `mega-pipeline` com os 5 sub-loops acoplados.
  * Teste unitário validando os 5 nós, 4 arestas e ausência de ciclos.

### 🛑 STOP-GATE DE AUDITORIA HUMANA — ÉPICO 20
* [ ] Suíte `go test ./...` 100% verde.
* [ ] Recompilação obrigatória: `go build -o bin/legoloop.exe ./cmd/legoloop`.
* [ ] Auditoria do usuário aprovada para avançar ao Épico 21.
