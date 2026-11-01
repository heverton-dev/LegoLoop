# 📐 SPEC-2026-11-01: Especificação de Engenharia do Modal Canônico Factual

> **Status:** APROVADO_PARA_EXECUCAO  
> **Data:** 2026-11-01  
> **Referência de Intenção:** [`docs/intents/2026-11-01_INTENT_REFATORACAO_MODAL_ITEM.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/intents/2026-11-01_INTENT_REFATORACAO_MODAL_ITEM.md)  
> **Mapa de Rotas:** [`docs/routes/2026-11-01_ROUTE_MODAL_ITEM_FACTUAL.json`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-01_ROUTE_MODAL_ITEM_FACTUAL.json)

---

## 1. Critérios de Aceite Binários

### CA-01: Preservação Integral de Dados na Ingestão (`palette.js`)
* A função `entryFromBlock(b)` deve preservar obrigatoriamente:
  * `consumes: b.consumes || []`
  * `produces: b.produces || ""`
  * `artifact: b.artifact || ""`
  * `gate: b.gate || ""`
  * `cost: b.cost || 0`
  * `suggested_model: b.suggested_model || ""`
  * `ref: b.ref || ""`
  * `raw: b`
* A função `entryFromItem(cat, it)` deve repassar `meta: it.meta || {}`, `path: it.path || ""`, `scope: it.scope || ""`, `origin: it.origin || ""`.
* **Gate:** `has_pattern('internal/studio/static/palette.js', 'consumes') && has_pattern('internal/studio/static/palette.js', 'artifact')`.

### CA-02: Especialização Tipada e Fim das Frases Mentirosas (`lego-item-modal.js`)
* Banimento absoluto de textos hardcoded como `"Executa tarefas técnicas determinísticas com isolamento estrito de ambiente."` e `"Entrega artefato imutável e validação técnica."`.
* **Gatilhos (Triggers):**
  * Não exibem entrada falsa ("Briefing / Requisitos"). Exibem `"Disparo Autônomo"` com badge de tipo de gatilho (Manual, Cron, Webhook, Board, File).
  * Exibem artefato inicial real (ex: `WAVE_PLAN.json` para o trigger do GitHub Board) e gate formal (`file_exists(...)`).
* **Controles (Flow):**
  * Exibem diagrama de bifurcação condicional (`flow-if`: Saída TRUE / FALSE), iteração recursiva (`flow-loop`) ou sincronização (`flow-merge`).
  * Exibem artefato de roteamento (`DECISION.json`, `ITERATION.json`, `JOINED.json`).
* **Gates:**
  * **Exibição obrigatória do tipo de artefato auditado:**
    * `file_exists` ➔ Arquivo do Projeto
    * `has_binary_criteria` ➔ Markdown com Critérios Binários (`- [ ]`)
    * `has_acceptance_criteria` ➔ Seção de Critérios de Aceitação
    * `word_count` / `line_count` ➔ Métrica de Densidade Textual
    * `has_pattern` ➔ Padrão Regex / Conformidade ISO
    * `pr_merge_gate` ➔ Build, Testes, Lint e Segurança de PR
  * Exibição de assinatura, regra sintática e comportamento em caso de reprovação (bloqueio determinístico).
* **Sub-Loops:**
  * Exibem portas derivadas reais (`consumes` dos nós de entrada, `produces` do nó final, contagem de nós internos).
* **MCPs:**
  * Exibem transporte (`stdio`, `sse`, `http`), host remoto e comando base.
* **Gate:** `has_pattern('components/lego-item-modal.js', 'auditedArtifact') && !has_pattern('components/lego-item-modal.js', 'Executa tarefas técnicas determinísticas')`.

### CA-03: Tríade de Compatibilidade na Sub-Aba (`lego-item-modal.js`)
* A sub-aba de compatibilidade deve organizar e permitir filtrar por:
  1. **⬅️ Antecessoras:** Peças que produzem o que a peça atual consome.
  2. **🔄 Substitutas (Equivalentes/Swap):** Peças com o mesmo contrato de I/O prontas para troca.
  3. **➡️ Sucessoras:** Peças que consomem o que a peça atual entrega.
* Filtros por categoria adicionais (Skills, MCPs, Scripts, Hooks, Gates, Loops).
* Ações nos cards contextuais (ex: *"Trocar (Swap)"*, *"Conectar Antes"*, *"Conectar Depois"*).
* **Gate:** Presença das estruturas de filtro e cálculo ontológico da tríade.

### CA-04: Fidelidade em Ocorrências no SO e Código/Fonte
* **Ocorrências no SO:**
  * Peças virtuais/core: badge *"Nativo do Core LegoLoop"*.
  * Arquivos físicos: caminho completo, status no cofre canônico e botão *"Abrir no Explorer"*.
  * MCPs: arquivo de configuração de origem.
* **Código / Fonte:**
  * Renderização especializada (Markdown com frontmatter para skills, JSON tipado para blocos/controles/loops, expressão de regra para gates, script executável para scripts/hooks).
* **Gate:** `has_pattern('components/lego-item-modal.js', 'Nativo do Core LegoLoop')`.

### CA-05: Testes Unitários e Recompilação Go
* Testes em `internal/studio/server_test.go` validando que todos os campos do contrato trafegam sem perdas e que o modal é servido intacto.
* Recompilação do binário: `go build -o bin/legoloop.exe ./cmd/legoloop`.
* **Gate:** `go test ./...` 100% aprovado e binário atualizado.
