# 📋 PLAN-2026-11-01: Plano de Execução Factual por Micro-Tickets

> **Status:** PRONTO_PARA_EXECUCAO  
> **Estratégia:** Orchestrix Cirúrgica por Bloco (Lei #4 — contexto estéril, sem subagentes de harness)  
> **Referência de Rota:** [`docs/routes/2026-11-01_ROUTE_MODAL_ITEM_FACTUAL.json`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/routes/2026-11-01_ROUTE_MODAL_ITEM_FACTUAL.json)

---

## Fatiamento dos Micro-Tickets

### 🎫 TICK-MODAL-01: Preservação de Dados da Paleta
* **Arquivo Alvo:** [`internal/studio/static/palette.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/palette.js)
* **Ação:** Atualizar `entryFromBlock(b)` e `entryFromItem(cat, it)` para repassar sem perdas: `consumes`, `produces`, `artifact`, `gate`, `cost`, `suggested_model`, `meta`, `path`, `scope`, `origin`, `ref`.
* **Gate de Saída:** Preservação de campos comprovada no objeto retornado pela paleta.

### 🎫 TICK-MODAL-02: Especialização Tipada e Fim das Frases Mentirosas
* **Arquivo Alvo:** [`internal/studio/static/components/lego-item-modal.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/components/lego-item-modal.js), [`internal/studio/static/style.css`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/style.css)
* **Ação:**
  * Implementar identificação de família (`skill`, `trigger`, `control`, `gate`, `loop`, `mcp`, `hook`, `script`).
  * Para **Triggers:** Exibir "Disparo Autônomo", artefato inicial (`b.artifact`) e gate real.
  * Para **Controles:** Exibir diagrama de decisão com ramificações TRUE/FALSE ou merge.
  * Para **Gates:** Exibir cartão de decisão binária (`PASS`/`FAIL`) destacando expressamente o **Tipo de Artefato Auditado**.
  * Para **Sub-Loops:** Exibir portas derivadas e métricas de nós.
  * Para **MCPs:** Exibir transporte (`stdio`/`sse`) e host.
  * Remover 100% das frases estáticas hardcoded (`modal-prop-what`, `modal-prop-does`, `modal-prop-delivers`).
* **Gate de Saída:** `auditedArtifact` funcional e zero frases estáticas.

### 🎫 TICK-MODAL-03: Tríade de Compatibilidade na Sub-Aba
* **Arquivo Alvo:** [`internal/studio/static/components/lego-item-modal.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/components/lego-item-modal.js)
* **Ação:**
  * Implementar agrupamento e filtros da tríade:
    * `[⬅️ Antecessoras]` (Quem alimenta esta peça)
    * `[🔄 Substitutas (Swap)]` (Quem tem mesmo contrato de I/O)
    * `[➡️ Sucessoras]` (Quem consome esta peça)
  * Renderizar botões contextuais nos cards (*"Trocar (Swap)"*, *"Conectar Antes"*, *"Conectar Depois"*).
* **Gate de Saída:** Filtros da tríade operacionais e integrados.

### 🎫 TICK-MODAL-04: Ocorrências no SO e Código Especializados
* **Arquivo Alvo:** [`internal/studio/static/components/lego-item-modal.js`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/static/components/lego-item-modal.js)
* **Ação:**
  * Exibir *"Nativo do Core LegoLoop"* para triggers, controles e gates embutidos.
  * Exibir manifesto JSON para blocos, regras formais para gates e código real para scripts.
* **Gate de Saída:** Zero itens sem indicação física ou virtual precisa.

### 🎫 TICK-MODAL-05: Testes Unitários Go e Recompilação
* **Arquivo Alvo:** [`internal/studio/server_test.go`](file:///C:/Users/trcnologia/Desktop/LegoLoop/internal/studio/server_test.go), [`bin/legoloop.exe`](file:///C:/Users/trcnologia/Desktop/LegoLoop/bin/legoloop.exe)
* **Ação:**
  * Adicionar testes unitários validando a entrega do modal factual e a integridade de rotas.
  * Executar `go test ./...` e compilar `go build -o bin/legoloop.exe ./cmd/legoloop`.
* **Gate de Saída:** 100% dos testes verdes e binário no PATH atualizado.
