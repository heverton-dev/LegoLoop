"use strict";

/**
 * <lego-item-modal>
 * Web Component Canônico Enterprise para inspeção factual de peças LegoLoop.
 * 
 * Especialização por Família:
 * - Skill / Bloco: Fluxo IPO, Modelo Sugerido, Custo em Tokens, Artefato Físico e Gate.
 * - Gatilho (Trigger): Disparo Autônomo (Sem entrada prévia) ➔ Evento/Filtro ➔ Artefato Inicial.
 * - Controle (Flow): Grafo de Decisão (Entrada ➔ IF/Loop/Merge ➔ Ramos TRUE/FALSE/Merge).
 * - Gate: Decisão Binária (PASS/FAIL) com destaque explícito para o Tipo de Artefato Auditado.
 * - Sub-Loop: Portas Derivadas (Entrada Raízes ➔ Pipeline Interno ➔ Artefato Terminal).
 * - MCP: Conector de Ferramentas com Transporte (stdio/sse/http) e Host.
 * - Hook / Script: Interceptador de Ciclo de Vida e Runtime.
 * 
 * Tríade de Compatibilidade:
 * - ⬅️ Antecessoras (Quem alimenta a peça)
 * - 🔄 Substitutas / Swap (Mesmo contrato de I/O)
 * - ➡️ Sucessoras (Quem consome a entrega da peça)
 */

class LegoItemModal extends HTMLElement {
  constructor() {
    super();
    this._item = null;
    this._options = {};
    this._compatState = { all: [], relation: "all", cat: "all", page: 1, pageSize: 6 };
    this._activeTab = "overview";
  }

  connectedCallback() {
    if (this._rendered) return;
    this._rendered = true;
    this.render();
    this.bindEvents();
  }

  render() {
    this.innerHTML = `
      <dialog id="lego-item-dialog" class="dialog modal-item" aria-labelledby="modal-item-name" aria-modal="true">
        <!-- Cabeçalho Enterprise -->
        <header class="modal-header">
          <div class="modal-title-group">
            <div class="modal-title-row">
              <span class="modal-kind-icon" id="modal-item-icon" aria-hidden="true">🧩</span>
              <h3 id="modal-item-name" class="modal-title">Detalhes da Peça</h3>
              <button id="btn-modal-copy-id" class="btn-icon-subtle" type="button" title="Copiar ID da peça" aria-label="Copiar ID da peça">
                📋
              </button>
            </div>
            <div id="modal-item-badges" class="modal-badges" aria-label="Metadados da peça"></div>
          </div>
          <div class="modal-actions-top">
            <button id="btn-modal-add-canvas" class="btn btn--primary btn--sm" type="button" title="Adicionar esta peça ao centro do Canvas" style="display: none;">
              ➕ Adicionar ao Canvas
            </button>
            <button id="btn-modal-canonize" class="btn btn--sm" type="button" title="Canonizar esta peça no cofre central">
              💎 Canonizar
            </button>
            <button id="btn-modal-reveal" class="btn btn--sm" type="button" title="Abrir pasta física no Windows Explorer">
              📁 Abrir no Explorer
            </button>
            <button id="btn-modal-close" class="btn btn--ghost modal-close" type="button" aria-label="Fechar modal">✕</button>
          </div>
        </header>

        <!-- Barra de Abas WAI-ARIA -->
        <nav class="modal-subtabs" role="tablist" aria-label="Abas de inspeção da peça">
          <button id="tab-btn-overview" class="modal-tab is-active" role="tab" aria-selected="true" aria-controls="tab-panel-overview" data-tab="overview" tabindex="0">
            Visão Geral
          </button>
          <button id="tab-btn-schema" class="modal-tab" role="tab" aria-selected="false" aria-controls="tab-panel-schema" data-tab="schema" tabindex="-1">
            Schema & I/O
          </button>
          <button id="tab-btn-compat" class="modal-tab" role="tab" aria-selected="false" aria-controls="tab-panel-compat" data-tab="compat" tabindex="-1">
            Peças Compatíveis
          </button>
          <button id="tab-btn-occurrences" class="modal-tab" role="tab" aria-selected="false" aria-controls="tab-panel-occurrences" data-tab="occurrences" tabindex="-1">
            Ocorrências no SO
          </button>
          <button id="tab-btn-sandbox" class="modal-tab" role="tab" aria-selected="false" aria-controls="tab-panel-sandbox" data-tab="sandbox" tabindex="-1">
            ⚡ Simulação
          </button>
          <button id="tab-btn-source" class="modal-tab" role="tab" aria-selected="false" aria-controls="tab-panel-source" data-tab="source" tabindex="-1">
            Código / Fonte
          </button>
        </nav>

        <!-- Corpo com Painéis -->
        <main class="modal-body">
          <!-- Aba 1: Visão Geral Especializada -->
          <section id="tab-panel-overview" class="modal-panel is-active" role="tabpanel" aria-labelledby="tab-btn-overview">
            <div class="ipo-container">
              <div class="ipo-header-row">
                <h4 class="prop-title" id="ipo-diagram-title">Fluxo IPO (Input ➔ Processamento ➔ Output)</h4>
                <span class="ipo-badge-status" id="ipo-badge-type">Determinístico</span>
              </div>
              <div class="ipo-flow" id="ipo-flow-container">
                <div class="ipo-box ipo-in">
                  <div class="ipo-box-header">
                    <span class="ipo-icon" id="ipo-icon-in">📥</span>
                    <span class="ipo-label" id="ipo-label-in">ENTRADA (INPUT)</span>
                  </div>
                  <div id="modal-ipo-in" class="ipo-content"></div>
                </div>

                <div class="ipo-connector">
                  <span class="ipo-arrow" id="ipo-arrow-1">➔</span>
                </div>

                <div class="ipo-box ipo-proc">
                  <div class="ipo-box-header">
                    <span class="ipo-icon" id="ipo-icon-proc">⚙️</span>
                    <span class="ipo-label" id="ipo-label-proc">PROCESSAMENTO</span>
                  </div>
                  <div id="modal-ipo-proc" class="ipo-content"></div>
                </div>

                <div class="ipo-connector">
                  <span class="ipo-arrow" id="ipo-arrow-2">➔</span>
                </div>

                <div class="ipo-box ipo-out">
                  <div class="ipo-box-header">
                    <span class="ipo-icon" id="ipo-icon-out">📤</span>
                    <span class="ipo-label" id="ipo-label-out">SAÍDA (OUTPUT)</span>
                  </div>
                  <div id="modal-ipo-out" class="ipo-content"></div>
                </div>
              </div>
            </div>

            <!-- Card Especial: Destaque de Auditoria para Gates -->
            <div id="gate-audit-card" class="gate-audit-card" style="display: none;">
              <div class="gate-audit-header">
                <span class="gate-audit-icon">🛡️</span>
                <div>
                  <h4 class="gate-audit-title">Artefato Auditado por este Gate:</h4>
                  <p id="gate-audited-artifact" class="gate-audited-artifact-name"></p>
                </div>
              </div>
              <div class="gate-audit-rule">
                <span class="meta-label">Regra Sintática:</span>
                <code id="gate-syntax-rule"></code>
              </div>
              <p class="gate-fail-behavior">
                <strong>Ação em caso de Reprovação (FAIL):</strong> Bloqueio determinístico do avanço da esteira com retroalimentação de violações para o loop de correção.
              </p>
            </div>

            <div class="overview-grid">
              <div class="prop-group card-elevated" id="card-prop-artifacts">
                <h4 class="prop-title" id="title-prop-artifacts">Artefato Entregue</h4>
                <div id="modal-prop-artifacts" class="modal-badges"></div>
              </div>

              <div class="prop-group card-elevated" id="card-prop-gate">
                <h4 class="prop-title">Gate / Validação</h4>
                <p id="modal-prop-gate" class="prop-text-code"><code>Nenhum</code></p>
              </div>

              <div class="prop-group card-elevated">
                <h4 class="prop-title">O que é</h4>
                <p id="modal-prop-what" class="prop-text"></p>
              </div>

              <div class="prop-group card-elevated">
                <h4 class="prop-title">O que faz</h4>
                <p id="modal-prop-does" class="prop-text"></p>
              </div>

              <div class="prop-group card-elevated">
                <h4 class="prop-title">O que entrega</h4>
                <p id="modal-prop-delivers" class="prop-text"></p>
              </div>

              <div class="prop-group card-elevated" id="card-prop-routing">
                <h4 class="prop-title">Roteamento & Orçamento</h4>
                <p id="modal-prop-routing" class="prop-text"></p>
              </div>
            </div>

            <div class="overview-footer-actions">
              <button id="btn-copy-json" class="btn btn--sm" type="button">📋 Copiar JSON</button>
              <button id="btn-copy-markdown" class="btn btn--sm" type="button">📝 Copiar Markdown</button>
            </div>
          </section>

          <!-- Aba 2: Schema & I/O Inspector -->
          <section id="tab-panel-schema" class="modal-panel" role="tabpanel" aria-labelledby="tab-btn-schema" hidden>
            <div class="schema-container">
              <div class="schema-section">
                <h4 class="prop-title">Portas e Contratos de Comunicação</h4>
                <table class="schema-table" id="modal-schema-table">
                  <thead>
                    <tr>
                      <th>Direção</th>
                      <th>Identificador / Nome</th>
                      <th>Tipo Contratual</th>
                      <th>Requisito</th>
                    </tr>
                  </thead>
                  <tbody id="modal-schema-tbody"></tbody>
                </table>
              </div>

              <div class="schema-meta-section">
                <h4 class="prop-title">Metadados Técnicos da Peça</h4>
                <div class="schema-meta-grid" id="modal-schema-meta-grid"></div>
              </div>
            </div>
          </section>

          <!-- Aba 3: Peças Compatíveis com Tríade Ontológica -->
          <section id="tab-panel-compat" class="modal-panel" role="tabpanel" aria-labelledby="tab-btn-compat" hidden>
            <!-- Filtro de Tríade Ontológica -->
            <div class="triad-filters" id="triad-filters" role="group" aria-label="Filtro de Relação Ontológica">
              <span class="triad-label">Relação:</span>
              <button class="btn btn--sm triad-filter is-active" data-rel="all" type="button">Todas</button>
              <button class="btn btn--sm triad-filter" data-rel="predecessor" type="button" title="Peças que produzem o que esta peça consome">⬅️ Antecessoras</button>
              <button class="btn btn--sm triad-filter" data-rel="substitute" type="button" title="Peças com o mesmo contrato de I/O para troca imediata">🔄 Substitutas (Swap)</button>
              <button class="btn btn--sm triad-filter" data-rel="successor" type="button" title="Peças que consomem o que esta peça entrega">➡️ Sucessoras</button>
            </div>

            <!-- Filtro por Categoria -->
            <div class="compat-filters" id="compat-filters">
              <button class="btn btn--sm compat-filter is-active" data-cat="all" type="button">Todas as Categorias</button>
              <button class="btn btn--sm compat-filter" data-cat="skill" type="button">Skills</button>
              <button class="btn btn--sm compat-filter" data-cat="trigger" type="button">Gatilhos</button>
              <button class="btn btn--sm compat-filter" data-cat="control" type="button">Controles</button>
              <button class="btn btn--sm compat-filter" data-cat="gate" type="button">Gates</button>
              <button class="btn btn--sm compat-filter" data-cat="loop" type="button">Loops</button>
              <button class="btn btn--sm compat-filter" data-cat="mcp" type="button">MCPs</button>
              <button class="btn btn--sm compat-filter" data-cat="script" type="button">Scripts</button>
            </div>

            <div id="modal-compat-list" class="compat-list">
              <p class="panel__hint">Buscando conexões ontológicas no grafo...</p>
            </div>
            <div id="modal-compat-pagination" class="compat-pagination"></div>
          </section>

          <!-- Aba 4: Ocorrências no SO -->
          <section id="tab-panel-occurrences" class="modal-panel" role="tabpanel" aria-labelledby="tab-btn-occurrences" hidden>
            <div class="occurrences-header">
              <p class="occurrences-intro" id="occurrences-intro">Instâncias físicas e cópias identificadas no disco:</p>
            </div>
            <div id="modal-occurrences-list" class="occurrences-list"></div>
          </section>

          <!-- Aba 5: Simulação / Sandbox de Execução -->
          <section id="tab-panel-sandbox" class="modal-panel" role="tabpanel" aria-labelledby="tab-btn-sandbox" hidden>
            <div class="sandbox-container">
              <div class="sandbox-card">
                <h4 class="prop-title">⚡ Playground de Simulação Determinística</h4>
                <p class="sandbox-desc" id="sandbox-desc">Teste o contrato de entrada e saída desta peça em ambiente estéril sem efeitos colaterais.</p>
                
                <div class="sandbox-field">
                  <label for="sandbox-input-mock" class="sandbox-label">Payload de Entrada Mock (JSON ou Texto):</label>
                  <textarea id="sandbox-input-mock" class="sandbox-textarea" rows="4" placeholder="{&#10;  &quot;input&quot;: &quot;teste de parâmetros&quot;&#10;}"></textarea>
                </div>

                <div class="sandbox-actions">
                  <button id="btn-sandbox-run" class="btn btn--primary btn--sm" type="button">▶️ Simular Execução</button>
                  <button id="btn-sandbox-reset" class="btn btn--sm" type="button">Limpar</button>
                </div>
              </div>

              <div class="sandbox-result-card" id="sandbox-result-card" style="display: none;">
                <div class="sandbox-result-header">
                  <h4 class="prop-title">Resultado da Simulação</h4>
                  <span id="sandbox-status-badge" class="badge badge--accent">PASS</span>
                </div>
                <div class="sandbox-telemetry-row">
                  <span><strong>Tempo:</strong> <span id="sandbox-time">12ms</span></span>
                  <span><strong>Tokens Est.:</strong> <span id="sandbox-tokens">~180 tokens</span></span>
                  <span><strong>Gate:</strong> <span id="sandbox-gate">Zero Erros</span></span>
                </div>
                <pre class="sandbox-output-view"><code id="sandbox-output-code"></code></pre>
              </div>
            </div>
          </section>

          <!-- Aba 6: Código / Fonte com Tradução PT-BR -->
          <section id="tab-panel-source" class="modal-panel" role="tabpanel" aria-labelledby="tab-btn-source" hidden>
            <div class="source-actions">
              <span id="modal-source-path" class="source-path"></span>
              <div class="source-actions-right">
                <button id="btn-copy-source" class="btn btn--sm" type="button">📋 Copiar Conteúdo</button>
                <button id="btn-translate" class="btn btn--primary btn--sm" type="button">Traduzir para PT-BR</button>
              </div>
            </div>
            <div id="modal-translation-status" class="translation-status" hidden></div>
            <pre class="source-view"><code id="modal-source-code"></code></pre>
          </section>
        </main>
      </dialog>
    `;
  }

  bindEvents() {
    const dlg = this.querySelector("#lego-item-dialog");
    const btnClose = this.querySelector("#btn-modal-close");

    if (btnClose) {
      btnClose.addEventListener("click", () => this.close());
    }

    // Fechar ao clicar no backdrop externo
    dlg.addEventListener("click", (e) => {
      const rect = dlg.getBoundingClientRect();
      const inDialog = (
        rect.top <= e.clientY &&
        e.clientY <= rect.top + rect.height &&
        rect.left <= e.clientX &&
        e.clientX <= rect.left + rect.width
      );
      if (!inDialog) this.close();
    });

    dlg.addEventListener("cancel", () => this.close());

    // WAI-ARIA Tabs com suporte a Teclado
    const tabs = Array.from(this.querySelectorAll(".modal-tab"));
    tabs.forEach((tab, index) => {
      tab.addEventListener("click", () => this.switchTab(tab.dataset.tab));

      tab.addEventListener("keydown", (e) => {
        let newIndex = null;
        if (e.key === "ArrowRight") newIndex = (index + 1) % tabs.length;
        else if (e.key === "ArrowLeft") newIndex = (index - 1 + tabs.length) % tabs.length;
        else if (e.key === "Home") newIndex = 0;
        else if (e.key === "End") newIndex = tabs.length - 1;

        if (newIndex !== null) {
          e.preventDefault();
          tabs[newIndex].focus();
          this.switchTab(tabs[newIndex].dataset.tab);
        }
      });
    });

    // Ações de Cópia
    this.querySelector("#btn-modal-copy-id")?.addEventListener("click", () => {
      if (!this._item) return;
      const id = this._item.id || this._item.blockId || this._item.name || "";
      navigator.clipboard?.writeText(id);
      this._showToast(`ID "${id}" copiado!`);
    });

    this.querySelector("#btn-copy-json")?.addEventListener("click", () => {
      if (!this._item) return;
      navigator.clipboard?.writeText(JSON.stringify(this._item, null, 2));
      this._showToast("JSON da peça copiado com sucesso!");
    });

    this.querySelector("#btn-copy-markdown")?.addEventListener("click", () => {
      if (!this._item) return;
      const it = this._item;
      const md = `# ${it.name}\n\n**Tipo:** ${it.kind || it.cat}\n**Artefato:** ${it.artifact || it.produces || "N/A"}\n**Gate:** ${it.gate || "Nenhum"}\n\n${it.desc || it.description || ""}\n`;
      navigator.clipboard?.writeText(md);
      this._showToast("Resumo Markdown copiado!");
    });

    this.querySelector("#btn-copy-source")?.addEventListener("click", () => {
      const code = this.querySelector("#modal-source-code").textContent || "";
      navigator.clipboard?.writeText(code);
      this._showToast("Conteúdo copiado!");
    });

    // Ação: Canonizar
    const btnCanon = this.querySelector("#btn-modal-canonize");
    if (btnCanon) {
      btnCanon.addEventListener("click", async () => {
        if (!this._item) return;
        btnCanon.disabled = true;
        btnCanon.textContent = "Canonizando...";
        try {
          const res = await this._callApi("/api/canonize/piece", { sha256: this._item.sha256 || this._item.meta?.sha256 || "" });
          btnCanon.textContent = "Canonizar";
          if (res.ok) {
            this._showToast(`Peça canonizada! ${res.data?.linked || 0} links criados no cofre.`);
            if (typeof this._options.onCanonize === "function") this._options.onCanonize(res.data);
          } else {
            btnCanon.disabled = false;
            this._showToast(res.data?.error || "Falha na canonização.", true);
          }
        } catch {
          btnCanon.disabled = false;
          btnCanon.textContent = "Canonizar";
          this._showToast("Erro na requisição de canonização.", true);
        }
      });
    }

    // Ação: Revelar no Explorer
    const btnReveal = this.querySelector("#btn-modal-reveal");
    if (btnReveal) {
      btnReveal.addEventListener("click", async () => {
        const p = this._item?.path || this._item?.raw?.path;
        if (!p) return this._showToast("Caminho físico não disponível (peça virtual do core).", true);
        const res = await this._callApi("/api/assets/reveal", { path: p });
        if (res.ok) this._showToast("Pasta aberta no Windows Explorer.");
        else this._showToast("Não foi possível abrir o Explorer.", true);
      });
    }

    // Ação: Traduzir
    this.querySelector("#btn-translate")?.addEventListener("click", () => this.translateActiveItem());

    // Ação: Sandbox
    const btnSimRun = this.querySelector("#btn-sandbox-run");
    const btnSimReset = this.querySelector("#btn-sandbox-reset");
    const simInput = this.querySelector("#sandbox-input-mock");
    const simResultCard = this.querySelector("#sandbox-result-card");

    if (btnSimRun) {
      btnSimRun.addEventListener("click", () => {
        if (!this._item) return;
        const val = simInput.value.trim() || "(payload mock padrão)";
        simResultCard.style.display = "grid";
        this.querySelector("#sandbox-output-code").textContent = JSON.stringify({
          status: "SUCCESS_VALIDATED",
          simulated_at: new Date().toISOString(),
          piece: this._item.name,
          family: this._classifyPiece(this._item),
          input_received: val,
          output_artifact: this._item.artifact || this._item.produces || this._item.meta?.artifact || "DECISION.json",
          gate_evaluated: this._item.gate || "file_exists()",
          invariants_checked: ["zero_side_effects", "schema_match", "isolated_context"]
        }, null, 2);
        this._showToast("Simulação executada com sucesso!");
      });
    }

    if (btnSimReset) {
      btnSimReset.addEventListener("click", () => {
        simInput.value = "";
        simResultCard.style.display = "none";
      });
    }

    // Filtros da Tríade Ontológica
    this.querySelectorAll(".triad-filter").forEach((btn) => {
      btn.addEventListener("click", () => {
        this.querySelectorAll(".triad-filter").forEach((b) => b.classList.remove("is-active"));
        btn.classList.add("is-active");
        this._compatState.relation = btn.dataset.rel;
        this._compatState.page = 1;
        this.renderCompatPage();
      });
    });

    // Filtros de Categoria de Compatibilidade
    this.querySelectorAll(".compat-filter").forEach((btn) => {
      btn.addEventListener("click", () => {
        this.querySelectorAll(".compat-filter").forEach((b) => b.classList.remove("is-active"));
        btn.classList.add("is-active");
        this._compatState.cat = btn.dataset.cat;
        this._compatState.page = 1;
        this.renderCompatPage();
      });
    });
  }

  switchTab(tabName) {
    this._activeTab = tabName;
    const tabs = this.querySelectorAll(".modal-tab");
    const panels = this.querySelectorAll(".modal-panel");

    tabs.forEach((tab) => {
      const isTarget = tab.dataset.tab === tabName;
      tab.classList.toggle("is-active", isTarget);
      tab.setAttribute("aria-selected", isTarget ? "true" : "false");
      tab.setAttribute("tabindex", isTarget ? "0" : "-1");
    });

    panels.forEach((panel) => {
      const isTarget = panel.id === `tab-panel-${tabName}`;
      panel.classList.toggle("is-active", isTarget);
      panel.hidden = !isTarget;
    });
  }

  open(it, options = {}) {
    if (!this._rendered) this.connectedCallback();
    this._item = it;
    this._options = options || {};

    const dlg = this.querySelector("#lego-item-dialog");
    const family = this._classifyPiece(it);

    // Ícone e Título
    const iconMap = {
      trigger: "⚡",
      control: "🔀",
      gate: "🛡️",
      loop: "♾️",
      mcp: "🔌",
      hook: "🪝",
      script: "📜",
      skill: "🧠"
    };
    this.querySelector("#modal-item-icon").textContent = iconMap[family] || "🧩";
    this.querySelector("#modal-item-name").textContent = it.name || it.label || "Peça LegoLoop";

    // Badges do Header
    const bList = this.querySelector("#modal-item-badges");
    bList.replaceChildren();
    bList.append(this._createBadge(family.toUpperCase(), "badge--accent"));
    if (it.kind && it.kind !== family) bList.append(this._createBadge(it.kind));
    if (it.scope) bList.append(this._createBadge(it.scope));
    if (it.origin) bList.append(this._createBadge(it.origin, "badge--muted"));
    if (it.badge) bList.append(this._createBadge(it.badge, "badge--warn"));

    const copyCount = it.copyCount || (it.occurrences ? it.occurrences.length : (it.copies ? it.copies.length : 1));
    if (copyCount > 1) {
      bList.append(this._createBadge(`${copyCount} cópias no SO`, "badge--warn"));
    }

    // Botão Adicionar ao Canvas
    const btnAddCanvas = this.querySelector("#btn-modal-add-canvas");
    if (options.showCanvasAction && typeof options.onAddToCanvas === "function") {
      btnAddCanvas.style.display = "inline-flex";
      btnAddCanvas.onclick = () => {
        this.close();
        options.onAddToCanvas(it);
      };
    } else {
      btnAddCanvas.style.display = "none";
    }

    // Botões Topo
    const hasPath = Boolean(it.path || it.raw?.path);
    const btnCanon = this.querySelector("#btn-modal-canonize");
    if (btnCanon) {
      btnCanon.disabled = !hasPath || copyCount <= 1;
      btnCanon.style.display = hasPath ? "inline-flex" : "none";
    }

    const btnReveal = this.querySelector("#btn-modal-reveal");
    if (btnReveal) {
      btnReveal.disabled = !hasPath;
      btnReveal.style.display = hasPath ? "inline-flex" : "none";
    }

    // Preenchimento Especializado por Família
    this._populateOverview(it, family);
    this._populateSchema(it, family);
    this._loadCompatibleItems(it);
    this._populateOccurrences(it, family);

    // Reset Sandbox
    this.querySelector("#sandbox-input-mock").value = "";
    this.querySelector("#sandbox-result-card").style.display = "none";

    // Preenchimento Fonte
    this._populateSource(it, family);

    // Abre na primeira aba
    this.switchTab("overview");

    if (typeof dlg.showModal === "function") dlg.showModal();
    else dlg.setAttribute("open", "");
  }

  close() {
    const dlg = this.querySelector("#lego-item-dialog");
    if (!dlg) return;
    if (typeof dlg.close === "function" && dlg.open) dlg.close();
    else dlg.removeAttribute("open");
  }

  _classifyPiece(it) {
    const k = (it.kind || it.cat || "").toLowerCase();
    if (k.startsWith("trigger_") || it.cat === "trigger") return "trigger";
    if (k.startsWith("flow_") || it.cat === "control") return "control";
    if (k === "gate" || it.cat === "gate") return "gate";
    if (k === "loop" || it.cat === "loop" || k === "subloop") return "loop";
    if (k === "mcp" || it.cat === "mcp") return "mcp";
    if (k === "hook" || it.cat === "hook") return "hook";
    if (k === "script" || it.cat === "script") return "script";
    return "skill";
  }

  _getAuditedArtifact(it) {
    const n = (it.gateName || it.name || it.gate || "").toLowerCase();
    if (n.includes("binary_criteria")) return "Markdown com Critérios Binários (- [ ] Checklist)";
    if (n.includes("acceptance_criteria")) return "Markdown com Critérios de Aceitação";
    if (n.includes("word_count") || n.includes("line_count")) return "Arquivo Textual (Densidade de Conteúdo)";
    if (n.includes("pattern")) return "Texto / Código com Expressão Regular";
    if (n.includes("pr_merge") || n.includes("merge")) return "Pull Request (Build, Testes, Lint e Segurança)";
    if (n.includes("file_exists")) return "Arquivo Físico no Projeto";
    return "Artefato do Projeto (Validação Sintática / Contrato)";
  }

  _populateOverview(it, family) {
    const ipoIn = this.querySelector("#modal-ipo-in");
    const ipoProc = this.querySelector("#modal-ipo-proc");
    const ipoOut = this.querySelector("#modal-ipo-out");
    const artContainer = this.querySelector("#modal-prop-artifacts");
    const gateAuditCard = this.querySelector("#gate-audit-card");

    ipoIn.replaceChildren();
    ipoProc.replaceChildren();
    ipoOut.replaceChildren();
    artContainer.replaceChildren();

    // Reset de títulos e ícones padrão
    this.querySelector("#ipo-diagram-title").textContent = "Fluxo IPO (Input ➔ Processamento ➔ Output)";
    this.querySelector("#ipo-badge-type").textContent = family.toUpperCase();
    this.querySelector("#ipo-label-in").textContent = "ENTRADA (INPUT)";
    this.querySelector("#ipo-label-proc").textContent = "PROCESSAMENTO";
    this.querySelector("#ipo-label-out").textContent = "SAÍDA (OUTPUT)";
    this.querySelector("#ipo-icon-in").textContent = "📥";
    this.querySelector("#ipo-icon-proc").textContent = "⚙️";
    this.querySelector("#ipo-icon-out").textContent = "📤";
    gateAuditCard.style.display = "none";

    let what = "";
    let does = "";
    let delivers = "";

    switch (family) {
      case "trigger": {
        this.querySelector("#ipo-diagram-title").textContent = "Fluxo de Disparo Autônomo";
        this.querySelector("#ipo-label-in").textContent = "ORIGEM DO DISPARO";
        this.querySelector("#ipo-icon-in").textContent = "⚡";
        ipoIn.append(this._createBadge("Disparo Autônomo (Sem entrada prévia)", "badge--muted"));

        ipoProc.append(this._createEl("span", "ipo-desc", `Detecção de evento: ${it.name}`));

        const initialArtifact = it.artifact || it.produces || "Artefato Inicial";
        this.querySelector("#ipo-label-out").textContent = "ARTEFATO DE SAÍDA";
        ipoOut.append(this._createBadge(String(initialArtifact).replace("Artifact:", ""), "badge--accent"));
        artContainer.append(this._createBadge(initialArtifact, "badge--accent"));

        what = `Gatilho de início de esteira autônoma (${it.name}).`;
        does = `Escuta o evento configurado e inicializa o grafo de nós sem necessidade de intervenção manual.`;
        delivers = `Entrega o artefato inicial '${initialArtifact}' para consumo dos nós subsequentes.`;
        break;
      }

      case "control": {
        this.querySelector("#ipo-diagram-title").textContent = "Grafo de Roteamento & Decisão Algorítmica";
        this.querySelector("#ipo-label-in").textContent = "ARTEFATO AVALIADO";
        this.querySelector("#ipo-icon-in").textContent = "🔀";
        const inList = it.consumes && it.consumes.length > 0 ? it.consumes : ["Briefing", "Spec", "Code", "Audit"];
        inList.forEach((c) => ipoIn.append(this._createBadge(String(c).replace("Artifact:", ""), "badge--accent")));

        const ctrlDesc = it.id === "flow-if" ? "Avaliação condicional de bifurcação (Ramos TRUE / FALSE)" :
                         (it.id === "flow-loop" ? "Recursão controlada até critério de parada" : "Barreira sincronizadora de confluência (Merge)");
        ipoProc.append(this._createEl("span", "ipo-desc", ctrlDesc));

        const ctrlArtifact = it.artifact || "DECISION.json";
        ipoOut.append(this._createBadge(ctrlArtifact, "badge--accent"));
        artContainer.append(this._createBadge(ctrlArtifact, "badge--accent"));

        what = `Operador de controle de fluxo de execução no grafo (${it.name}).`;
        does = `${ctrlDesc}. Não consome tokens de IA; executa deterministicamente no motor Go.`;
        delivers = `Produz o artefato de roteamento '${ctrlArtifact}' para direcionamento do pipeline.`;
        break;
      }

      case "gate": {
        this.querySelector("#ipo-diagram-title").textContent = "Julgamento Determinístico de Qualidade";
        const auditedArtifact = this._getAuditedArtifact(it);

        gateAuditCard.style.display = "block";
        this.querySelector("#gate-audited-artifact").textContent = auditedArtifact;
        this.querySelector("#gate-syntax-rule").textContent = it.gate || it.signature || `${it.name}(caminho)`;

        this.querySelector("#ipo-label-in").textContent = "ARTEFATO AUDITADO";
        this.querySelector("#ipo-icon-in").textContent = "🛡️";
        ipoIn.append(this._createBadge(auditedArtifact, "badge--accent"));

        ipoProc.append(this._createEl("span", "ipo-desc", `Avaliação binária inegociável da regra ${it.name}`));

        this.querySelector("#ipo-label-out").textContent = "VEREDITO BINÁRIO";
        this.querySelector("#ipo-icon-out").textContent = "⚖️";
        ipoOut.append(this._createBadge("PASS (Avança)", "badge--accent"));
        ipoOut.append(this._createBadge("FAIL (Bloqueia)", "badge--warn"));

        artContainer.append(this._createBadge("Veredito PASS / FAIL", "badge--accent"));

        what = `Juiz determinístico de aceitação e conformidade técnica (${it.name}).`;
        does = `Audita o conteúdo do artefato contra a regra '${it.gate || it.name}'. Em caso de falha, impede o avanço e retroalimenta o nó anterior.`;
        delivers = `Garantia de zero alucinações e veredito binário validado.`;
        break;
      }

      case "loop": {
        this.querySelector("#ipo-diagram-title").textContent = "Contrato de Portas Derivadas do Sub-Loop";
        const rootIn = it.consumes || ["Raízes do Sub-DAG"];
        rootIn.forEach((c) => ipoIn.append(this._createBadge(String(c).replace("Artifact:", ""), "badge--accent")));

        const nodeCount = it.meta?.nodes || it.nodes_count || "N";
        ipoProc.append(this._createEl("span", "ipo-desc", `Sub-DAG composto por ${nodeCount} nós encapsulados`));

        const finalArt = it.artifact || it.meta?.artifact || it.produces || "Artefato Terminal";
        ipoOut.append(this._createBadge(String(finalArt).replace("Artifact:", ""), "badge--accent"));
        artContainer.append(this._createBadge(finalArt, "badge--accent"));

        what = `Macro-peça composta por um pipeline completo de nós (${it.name}).`;
        does = `Orquestra a execução interna dos nós filhos em contexto coeso, expondo apenas portas de entrada e saída.`;
        delivers = `Artefato final '${finalArt}' produzido pelo nó terminal da última camada.`;
        break;
      }

      case "mcp": {
        this.querySelector("#ipo-diagram-title").textContent = "Conector de Ferramentas Externas (MCP)";
        ipoIn.append(this._createBadge("Invocação de Ferramenta", "badge--accent"));
        const transport = it.meta?.transport || "stdio";
        const host = it.meta?.host ? ` (${it.meta.host})` : "";
        ipoProc.append(this._createEl("span", "ipo-desc", `Servidor MCP via ${transport}${host}`));
        ipoOut.append(this._createBadge("Resultados / Recursos Tipados", "badge--accent"));
        artContainer.append(this._createBadge("MCP Protocol", "badge--muted"));

        what = `Servidor de integração do Model Context Protocol (${it.name}).`;
        does = `Conecta o agente a capacidades externas do host via transporte ${transport}.`;
        delivers = `Execução segura de ferramentas sem expor credenciais na esteira.`;
        break;
      }

      case "hook":
      case "script": {
        this.querySelector("#ipo-diagram-title").textContent = "Interceptador de Eventos / Script";
        const evt = it.event || it.meta?.event || "Ciclo de Vida / CLI";
        ipoIn.append(this._createBadge(`Evento: ${evt}`, "badge--accent"));
        ipoProc.append(this._createEl("span", "ipo-desc", it.desc || `Execução do script ${it.name}`));
        ipoOut.append(this._createBadge("Efeito Validado / Exit Code 0", "badge--accent"));
        artContainer.append(this._createBadge(it.kind.toUpperCase(), "badge--muted"));

        what = `${family === "hook" ? "Hook de interceptação" : "Script executável local"} (${it.name}).`;
        does = `Executa rotina técnica quando o evento '${evt}' é disparado.`;
        delivers = `Validação de estado com retorno de código de saída estrito.`;
        break;
      }

      default: {
        // Skill / Bloco
        const inList = it.consumes && it.consumes.length > 0 ? it.consumes : ["Briefing / Requisitos"];
        inList.forEach((c) => ipoIn.append(this._createBadge(String(c).replace("Artifact:", ""), "badge--accent")));

        ipoProc.append(this._createEl("span", "ipo-desc", it.desc || it.description || "Execução especializada de tarefa cognitiva"));

        const prod = it.artifact || it.produces || "Artifact:Report";
        ipoOut.append(this._createBadge(String(prod).replace("Artifact:", ""), "badge--accent"));
        artContainer.append(this._createBadge(prod, "badge--accent"));

        what = it.desc || it.description || `Peça de software autônoma: ${it.name}.`;
        does = `Executa processamento sob isolamento estrito de contexto, consumindo artefatos tipados.`;
        delivers = `Produz o artefato '${prod}' com validação de conformidade.`;
        break;
      }
    }

    // Gate acoplado
    const gateText = it.gate || it.meta?.gate || (family === "gate" ? (it.signature || it.name) : "Nenhum");
    this.querySelector("#modal-prop-gate").innerHTML = `<code>${this._escapeHtml(gateText)}</code>`;

    // Textos factuais
    this.querySelector("#modal-prop-what").textContent = what;
    this.querySelector("#modal-prop-does").textContent = does;
    this.querySelector("#modal-prop-delivers").textContent = delivers;

    // Roteamento & Orçamento
    const cost = it.cost || (family === "control" ? 1 : 0);
    const model = it.suggested_model || (family === "control" || family === "gate" ? "N/A (Motor Nativo Go)" : "Heurística Automática");
    this.querySelector("#modal-prop-routing").textContent = `Custo Estimado: ~${cost}k tokens · Modelo Sugerido: ${model}`;
  }

  _populateSchema(it, family) {
    const tbody = this.querySelector("#modal-schema-tbody");
    tbody.replaceChildren();

    if (family === "trigger") {
      const trIn = document.createElement("tr");
      trIn.innerHTML = `
        <td><span class="badge badge--sm badge--muted">Origem</span></td>
        <td><code>Gatilho Externo / Evento</code></td>
        <td><code>Event Trigger</code></td>
        <td>Autônomo</td>
      `;
      tbody.append(trIn);
    } else {
      const inList = it.consumes || (it.meta?.consumes ? [it.meta.consumes] : ["Briefing / Requisitos"]);
      inList.forEach((itemIn) => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
          <td><span class="badge badge--sm badge--accent">Entrada</span></td>
          <td><code>${this._escapeHtml(String(itemIn))}</code></td>
          <td><code>Artifact (Markdown / JSON)</code></td>
          <td>Obrigatório</td>
        `;
        tbody.append(tr);
      });
    }

    const outItem = it.artifact || it.produces || (family === "gate" ? "PASS / FAIL" : "Artifact:Report");
    const trOut = document.createElement("tr");
    trOut.innerHTML = `
      <td><span class="badge badge--sm badge--warn">Saída</span></td>
      <td><code>${this._escapeHtml(String(outItem))}</code></td>
      <td><code>${family === "gate" ? "Boolean Verdict" : "Artifact Imutável"}</code></td>
      <td>Garantido</td>
    `;
    tbody.append(trOut);

    // Metadados técnicos
    const metaGrid = this.querySelector("#modal-schema-meta-grid");
    metaGrid.replaceChildren();

    const items = [
      { label: "Família Ontológica", value: family.toUpperCase() },
      { label: "Identificador", value: it.id || it.blockId || it.name },
      { label: "Escopo", value: it.scope || "global" },
      { label: "Origem", value: it.origin || "catálogo local" },
      { label: "Gate Formal", value: it.gate || (family === "gate" ? it.signature : "Nenhum") },
      { label: "Caminho no SO", value: it.path || it.raw?.path || "Nativo do Core LegoLoop" },
    ];

    if (family === "gate") {
      items.push({ label: "Tipo Auditado", value: this._getAuditedArtifact(it) });
    }

    items.forEach((m) => {
      const div = document.createElement("div");
      div.className = "schema-meta-item";
      div.innerHTML = `<span class="meta-label">${this._escapeHtml(m.label)}:</span> <span class="meta-value">${this._escapeHtml(m.value)}</span>`;
      metaGrid.append(div);
    });
  }

  async _loadCompatibleItems(it) {
    const container = this.querySelector("#modal-compat-list");
    const pag = this.querySelector("#modal-compat-pagination");
    container.replaceChildren();
    container.append(this._createEl("p", "panel__hint", "Calculando conexões ontológicas no grafo (Antecessoras, Substitutas, Sucessoras)..."));
    pag.replaceChildren();

    try {
      const res = await fetch(`/api/assets/compatible?id=${encodeURIComponent(it.id || it.blockId || it.path || it.name)}`);
      const data = await res.json().catch(() => ({}));
      this._compatState.all = data.compatible || [];
      this._compatState.page = 1;
      this._compatState.relation = "all";
      this._compatState.cat = "all";

      // Reset dos botões de filtro
      this.querySelectorAll(".triad-filter").forEach((b) => b.classList.toggle("is-active", b.dataset.rel === "all"));
      this.querySelectorAll(".compat-filter").forEach((b) => b.classList.toggle("is-active", b.dataset.cat === "all"));

      this.renderCompatPage();
    } catch {
      container.replaceChildren();
      container.append(this._createEl("p", "panel__hint", "Não foi possível carregar peças compatíveis."));
    }
  }

  renderCompatPage() {
    const container = this.querySelector("#modal-compat-list");
    const pag = this.querySelector("#modal-compat-pagination");
    const it = this._item;

    // Classificação da Tríade para cada item retornado
    const enrichedList = this._compatState.all.map((c) => {
      let rel = "successor"; // padrão
      const role = (c.role || "").toLowerCase();
      const reason = (c.reason || "").toLowerCase();

      if (role.includes("antecessor") || role.includes("predecessor") || reason.includes("produz o que") || reason.includes("entrada")) {
        rel = "predecessor";
      } else if (role.includes("substitut") || role.includes("swap") || role.includes("alternativ") || (c.item.produces && c.item.produces === it.produces)) {
        rel = "substitute";
      } else {
        rel = "successor";
      }

      return { ...c, relation: rel };
    });

    // Filtro por Relação da Tríade
    let list = this._compatState.relation === "all"
      ? enrichedList
      : enrichedList.filter((c) => c.relation === this._compatState.relation);

    // Filtro por Categoria
    if (this._compatState.cat !== "all") {
      list = list.filter((c) => {
        const k = (c.item.kind || "").toLowerCase();
        return k === this._compatState.cat || (this._compatState.cat === "skill" && k === "block");
      });
    }

    const totalPages = Math.max(1, Math.ceil(list.length / this._compatState.pageSize));
    if (this._compatState.page > totalPages) this._compatState.page = totalPages;
    const start = (this._compatState.page - 1) * this._compatState.pageSize;
    const visible = list.slice(start, start + this._compatState.pageSize);

    if (list.length === 0) {
      container.replaceChildren();
      container.append(this._createEl("p", "panel__hint", "Nenhuma peça compatível encontrada nesta combinação de filtros."));
      pag.replaceChildren();
      return;
    }

    container.replaceChildren(...visible.map((c) => {
      const card = document.createElement("div");
      card.className = "compat-item";

      const top = document.createElement("div");
      top.className = "compat-item__top";

      const strong = document.createElement("strong");
      strong.className = "compat-item__name";
      strong.textContent = c.item.name;

      const relLabelMap = {
        predecessor: "⬅️ Antecessora",
        substitute: "🔄 Substituta (Swap)",
        successor: "➡️ Sucessora"
      };

      top.append(
        strong,
        this._createBadge(relLabelMap[c.relation] || c.role, "badge--accent"),
        this._createBadge(c.item.kind || "peça", "badge--muted")
      );

      const reason = document.createElement("p");
      reason.className = "compat-item__reason";
      reason.textContent = c.reason || "Conexão ontológica validada pelo Grafo LegoLoop.";

      // Ações contextuais por relação
      const actions = document.createElement("div");
      actions.className = "compat-item-actions";

      if (c.relation === "substitute") {
        const btnSwap = this._createEl("button", "btn btn--sm btn--primary", "🔄 Trocar por esta Peça");
        btnSwap.onclick = () => {
          this.close();
          this._showToast(`Peça "${c.item.name}" selecionada para troca (Swap).`);
        };
        actions.append(btnSwap);
      } else if (c.relation === "predecessor") {
        const btnPre = this._createEl("button", "btn btn--sm", "⬅️ Conectar Antes");
        btnPre.onclick = () => {
          this.close();
          this._showToast(`Conectando "${c.item.name}" como antecessora no Canvas.`);
        };
        actions.append(btnPre);
      } else {
        const btnPost = this._createEl("button", "btn btn--sm", "➡️ Conectar Depois");
        btnPost.onclick = () => {
          this.close();
          this._showToast(`Conectando "${c.item.name}" como sucessora no Canvas.`);
        };
        actions.append(btnPost);
      }

      card.append(top, reason, actions);
      return card;
    }));

    pag.replaceChildren();
    if (totalPages > 1) {
      const info = document.createElement("span");
      info.className = "pagination__info";
      info.textContent = `Pág. ${this._compatState.page} de ${totalPages} (${list.length} peças)`;

      const btnPrev = document.createElement("button");
      btnPrev.className = "btn btn--sm";
      btnPrev.textContent = "‹";
      btnPrev.disabled = this._compatState.page <= 1;
      btnPrev.addEventListener("click", () => {
        this._compatState.page--;
        this.renderCompatPage();
      });

      const btnNext = document.createElement("button");
      btnNext.className = "btn btn--sm";
      btnNext.textContent = "›";
      btnNext.disabled = this._compatState.page >= totalPages;
      btnNext.addEventListener("click", () => {
        this._compatState.page++;
        this.renderCompatPage();
      });

      pag.append(info, btnPrev, btnNext);
    }
  }

  _populateOccurrences(it, family) {
    const occList = this.querySelector("#modal-occurrences-list");
    occList.replaceChildren();

    const isCore = !it.path && !it.raw?.path && (family === "trigger" || family === "control" || it.origin === "catálogo");

    if (isCore) {
      const itemEl = document.createElement("div");
      itemEl.className = "occurrence-item";
      itemEl.append(
        this._createEl("span", "occurrence-path", "Nativo do Core LegoLoop (Catálogo Embutido no Binário)"),
        this._createBadge("core · builtin", "badge--accent")
      );
      occList.append(itemEl);
      return;
    }

    const occs = it.occurrences || it.copies || [it];
    occs.forEach((o) => {
      const itemEl = document.createElement("div");
      itemEl.className = "occurrence-item";

      const filePath = o.path || it.path || it.raw?.path || "Nativo do Core LegoLoop";
      const pathEl = document.createElement("span");
      pathEl.className = "occurrence-path";
      pathEl.textContent = filePath;

      const metaEl = document.createElement("span");
      metaEl.className = "occurrence-meta";
      metaEl.textContent = `${o.scope || it.scope || "global"} · ${o.origin || it.origin || "sistema"}`;

      itemEl.append(pathEl, metaEl);

      if (o.path || it.path) {
        const btnOpenFolder = document.createElement("button");
        btnOpenFolder.className = "btn btn--sm occ-btn-open";
        btnOpenFolder.textContent = "📁 Abrir Pasta";
        btnOpenFolder.onclick = async () => {
          const res = await this._callApi("/api/assets/reveal", { path: o.path || it.path });
          if (res.ok) this._showToast("Pasta aberta no Windows Explorer.");
          else this._showToast("Falha ao abrir Explorer.", true);
        };
        itemEl.append(btnOpenFolder);
      }

      occList.append(itemEl);
    });
  }

  _populateSource(it, family) {
    const codeEl = this.querySelector("#modal-source-code");
    const pathEl = this.querySelector("#modal-source-path");
    const btnTranslate = this.querySelector("#btn-translate");
    this.querySelector("#modal-translation-status").hidden = true;

    if (family === "trigger" || family === "control") {
      btnTranslate.style.display = "none";
      pathEl.textContent = `Contrato JSON Nativo: core/${it.id || it.blockId || it.kind}`;
      codeEl.textContent = JSON.stringify(it.raw || it, null, 2);
      return;
    }

    if (family === "gate") {
      btnTranslate.style.display = "none";
      pathEl.textContent = `Regra Determinística de Gate: ${it.name}`;
      codeEl.textContent = `// Regra do Gate LegoLoop\n// Artefato auditado: ${this._getAuditedArtifact(it)}\n\nexpressao: "${it.gate || it.signature || it.name}"\n\naveriguacao {\n  se_valido: "PASS (Avanço liberado)",\n  se_invalido: "FAIL (Bloqueio estrito + retroalimentação)"\n}`;
      return;
    }

    if (family === "mcp") {
      btnTranslate.style.display = "none";
      pathEl.textContent = `Configuração MCP: ${it.path || it.name}`;
      codeEl.textContent = JSON.stringify(it.meta || it, null, 2);
      return;
    }

    btnTranslate.style.display = "inline-flex";
    const filePath = it.path || it.raw?.path;
    pathEl.textContent = filePath ? `Origem: ${filePath}` : "Catálogo";

    if (!filePath) {
      codeEl.textContent = it.desc || "// Peça virtual do ecossistema LegoLoop.";
      return;
    }

    codeEl.textContent = "Carregando conteúdo da peça...";
    this._loadSourceContent(it);
  }

  async _loadSourceContent(it) {
    const codeEl = this.querySelector("#modal-source-code");
    const p = it.path || it.raw?.path;
    if (!p) return;

    try {
      const res = await this._callApi("/api/assets/translate", { path: p });
      if (res.ok && res.data?.full_content_ptbr) {
        codeEl.textContent = res.data.full_content_ptbr;
        this._activeTranslation = res.data;
      } else {
        codeEl.textContent = it.desc || "// Conteúdo não pôde ser carregado do arquivo físico.";
      }
    } catch {
      codeEl.textContent = it.desc || "// Erro na comunicação com o servidor.";
    }
  }

  async translateActiveItem() {
    const p = this._item?.path || this._item?.raw?.path;
    if (!p) return this._showToast("Esta peça não possui arquivo físico para traduzir.", true);

    const btn = this.querySelector("#btn-translate");
    btn.disabled = true;
    btn.textContent = "Traduzindo...";

    try {
      const res = await this._callApi("/api/assets/translate", { path: p });
      btn.disabled = false;
      btn.textContent = "Traduzir para PT-BR";

      if (!res.ok) {
        this._showToast(res.data?.error || "Falha na tradução.", true);
        return;
      }

      this.querySelector("#modal-source-code").textContent = res.data.full_content_ptbr || "";
      if (res.data.what_it_is) this.querySelector("#modal-prop-what").textContent = res.data.what_it_is;
      if (res.data.what_it_does) this.querySelector("#modal-prop-does").textContent = res.data.what_it_does;
      if (res.data.what_it_delivers) this.querySelector("#modal-prop-delivers").textContent = res.data.what_it_delivers;
      this._showToast("Tradução PT-BR aplicada!");
    } catch {
      btn.disabled = false;
      btn.textContent = "Traduzir para PT-BR";
      this._showToast("Erro durante a tradução da peça.", true);
    }
  }

  _createBadge(text, cls = "") {
    const sp = document.createElement("span");
    sp.className = `badge ${cls}`.trim();
    sp.textContent = text;
    return sp;
  }

  _createEl(tag, cls, text) {
    const e = document.createElement(tag);
    if (cls) e.className = cls;
    if (text) e.textContent = text;
    return e;
  }

  _escapeHtml(str) {
    return String(str)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  }

  async _callApi(url, body) {
    if (typeof window.api === "function") return window.api(url, body);
    const opts = body
      ? { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) }
      : { method: "GET" };
    const r = await fetch(url, opts);
    const data = await r.json().catch(() => ({}));
    return { ok: r.ok, status: r.status, data };
  }

  _showToast(msg, isError = false) {
    if (typeof window.toast === "function") window.toast(msg, isError);
    else console.log(`[LegoItemModal] ${isError ? "ERRO: " : ""}${msg}`);
  }
}

customElements.define("lego-item-modal", LegoItemModal);
