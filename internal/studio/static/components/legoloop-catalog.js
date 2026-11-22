// Catalogo de dominio LegoLoop sobre o renderizador A2UI.
// Apenas DOM seguro (createElement/textContent): sem HTML bruto, sem estilo inline, sem emojis.

import { registerComponent } from "./a2ui-bundle.js";

function el(tag, className, text) {
  const node = document.createElement(tag);
  if (className) {
    node.className = className;
  }
  if (text !== undefined && text !== null) {
    node.textContent = String(text);
  }
  return node;
}

function statusClass(status) {
  const s = String(status || "pending").toLowerCase();
  if (s === "done" || s === "ok" || s === "completed" || s === "passed") {
    return "lego-status--ok";
  }
  if (s === "running" || s === "active" || s === "in_progress" || s === "working") {
    return "lego-status--running";
  }
  if (s === "failed" || s === "error" || s === "blocked" || s === "killed") {
    return "lego-status--fail";
  }
  return "lego-status--idle";
}

function renderDagGraph({ component, props, emitAction }) {
  const painel = el("article", "a2ui-painel lego-dag");
  painel.setAttribute("data-kind", "LegoDagGraph");

  const header = el("header", "lego-painel__header");
  header.appendChild(el("h3", "a2ui-painel__title", props.title || "Grafo do Processo"));
  const nodes = Array.isArray(props.nodes) ? props.nodes : [];
  header.appendChild(el("span", "pill", nodes.length + " nos"));
  painel.appendChild(header);

  if (nodes.length === 0) {
    painel.appendChild(el("p", "a2ui-text", "Nenhum no publicado nesta superficie."));
    return painel;
  }

  const list = el("ul", "lego-dag__nodes");
  for (const item of nodes) {
    const nodeId = item && (item.id || item.name) ? String(item.id || item.name) : "no";
    const label = item && item.label ? String(item.label) : nodeId;
    const status = item && item.status ? String(item.status) : "pending";

    const li = el("li", "lego-dag__node " + statusClass(status));
    const btn = el("button", "lego-dag__node-btn");
    btn.type = "button";
    btn.appendChild(el("span", "lego-dag__node-id", label));
    btn.appendChild(el("span", "pill lego-dag__badge", status));
    btn.addEventListener("click", () => {
      emitAction(props.action || "select_node", { graph_id: component.id, node_id: nodeId });
    });
    li.appendChild(btn);

    if (item && Array.isArray(item.depends_on) && item.depends_on.length > 0) {
      li.appendChild(el("span", "lego-dag__deps", "Depende de: " + item.depends_on.join(", ")));
    }
    list.appendChild(li);
  }
  painel.appendChild(list);
  return painel;
}

function renderProgressMeter({ props }) {
  const painel = el("article", "a2ui-painel lego-progress");
  painel.setAttribute("data-kind", "LegoProgressMeter");

  const raw = Number(props.percent ?? props.confidence ?? 0);
  const pct = Number.isFinite(raw) ? Math.max(0, Math.min(100, Math.round(raw))) : 0;

  const header = el("header", "lego-painel__header");
  header.appendChild(el("h3", "a2ui-painel__title", props.label || "Progresso"));
  header.appendChild(el("span", "pill lego-progress__value", pct + "%"));
  painel.appendChild(header);

  const bar = el("progress", "lego-progress__bar");
  bar.max = 100;
  bar.value = pct;
  painel.appendChild(bar);

  if (props.detail || props.summary) {
    painel.appendChild(el("p", "a2ui-text", props.detail || props.summary));
  }

  const gaps = Array.isArray(props.gaps) ? props.gaps : [];
  if (gaps.length > 0) {
    const gapBox = el("div", "lego-progress__gaps");
    gapBox.appendChild(el("span", "lego-section-label", "Lacunas abertas (" + gaps.length + ")"));
    const ul = el("ul", "a2ui-list");
    for (const gap of gaps) {
      ul.appendChild(el("li", "a2ui-list__item", typeof gap === "string" ? gap : JSON.stringify(gap)));
    }
    gapBox.appendChild(ul);
    painel.appendChild(gapBox);
  }

  return painel;
}

function renderStopGatepainel({ component, props, emitAction }) {
  const painel = el("article", "a2ui-painel lego-stop-gate");
  painel.setAttribute("data-kind", "LegoStopGatePanel");

  const header = el("header", "lego-painel__header");
  header.appendChild(el("h3", "a2ui-painel__title", props.title || "Stop-Gate de Auditoria"));
  header.appendChild(el("span", "pill a2ui-badge--attention", "Decisao humana"));
  painel.appendChild(header);

  if (props.summary) {
    painel.appendChild(el("p", "a2ui-text", props.summary));
  }

  const items = Array.isArray(props.items) ? props.items : [];
  if (items.length > 0) {
    const ul = el("ul", "a2ui-list");
    for (const it of items) {
      ul.appendChild(el("li", "a2ui-list__item", typeof it === "string" ? it : JSON.stringify(it)));
    }
    painel.appendChild(ul);
  }

  const noteLabel = el("label", "a2ui-label", "Observacao ou ajuste (opcional)");
  const noteInput = el("textarea", "a2ui-textarea lego-stop-gate__note");
  noteInput.rows = 2;
  noteInput.setAttribute("placeholder", "Descreva ajustes se escolher personalizar...");
  painel.append(noteLabel, noteInput);

  const actions = el("div", "a2ui-actions lego-stop-gate__actions");
  const choices = [
    { id: "1", label: "[1] Seguir", decision: "aprovar", cls: "btn lego-btn--approve" },
    { id: "2", label: "[2] Personalizar", decision: "personalizar", cls: "btn lego-btn--customize" },
    { id: "0", label: "[0] Cancelar", decision: "cancelar", cls: "btn lego-btn--cancel" }
  ];

  for (const ch of choices) {
    const btn = el("button", ch.cls, ch.label);
    btn.type = "button";
    btn.addEventListener("click", () => {
      emitAction(props.action || "stop_gate_decision", {
        gate_id: component.id,
        choice: ch.id,
        decision: ch.decision,
        note: noteInput.value.trim()
      });
    });
    actions.appendChild(btn);
  }

  painel.appendChild(actions);
  return painel;
}

function renderArtifactViewer({ surface, props }) {
  const painel = el("article", "a2ui-painel lego-artifact");
  painel.setAttribute("data-kind", "LegoArtifactViewer");

  const header = el("header", "lego-painel__header");
  const titleText = props.title || props.path || "Artefato";
  header.appendChild(el("h3", "a2ui-painel__title", titleText));

  const pre = el("pre", "lego-code", props.content || "");
  const status = el("span", "pill", props.path ? props.path : "memoria");
  header.appendChild(status);

  if (props.path && surface && surface.sessionId) {
    const reloadBtn = el("button", "btn", "Atualizar");
    reloadBtn.type = "button";
    const loadArtifact = async () => {
      const url = "/api/sessao/" + encodeURIComponent(surface.sessionId) +
        "/artefato?path=" + encodeURIComponent(String(props.path));
      try {
        const res = await fetch(url);
        if (!res.ok) {
          pre.textContent = "Nao foi possivel ler o artefato (HTTP " + res.status + ").";
          return;
        }
        pre.textContent = await res.text();
      } catch (_err) {
        pre.textContent = "Falha de rede ao carregar artefato.";
      }
    };
    reloadBtn.addEventListener("click", loadArtifact);
    header.appendChild(reloadBtn);
    if (!props.content) {
      loadArtifact();
    }
  }

  painel.append(header, pre);
  return painel;
}

function renderLogStream({ props }) {
  const painel = el("article", "a2ui-painel lego-logs");
  painel.setAttribute("data-kind", "LegoLogStream");

  const lines = Array.isArray(props.lines) ? props.lines.map((l) => String(l)) : [];
  const header = el("header", "lego-painel__header");
  header.appendChild(el("h3", "a2ui-painel__title", props.title || "Fluxo de Logs"));
  const counter = el("span", "pill", lines.length + " linhas");
  header.appendChild(counter);
  painel.appendChild(header);

  const filterBar = el("div", "lego-logs__toolbar");
  const filterInput = el("input", "a2ui-input lego-logs__filter");
  filterInput.type = "search";
  filterInput.setAttribute("placeholder", "Filtrar linhas de log...");
  filterBar.appendChild(filterInput);
  painel.appendChild(filterBar);

  const pre = el("pre", "lego-code lego-logs__output", lines.join("\n"));
  filterInput.addEventListener("input", () => {
    const term = filterInput.value.trim().toLowerCase();
    const filtered = term ? lines.filter((l) => l.toLowerCase().includes(term)) : lines;
    pre.textContent = filtered.join("\n");
    counter.textContent = filtered.length + " linhas";
  });

  painel.appendChild(pre);
  return painel;
}

const swarmDrafts = new Map();

function swarmDraftKey(surface, component) {
  const sid = surface && surface.id ? String(surface.id) : "default";
  const cid = component && component.id ? String(component.id) : "swarm";
  return sid + ":" + cid;
}

function renderSwarmBoard({ surface, component, props, emitAction }) {
  const painel = el("article", "a2ui-painel lego-swarm");
  painel.setAttribute("data-kind", "LegoSwarmBoard");

  const workers = Array.isArray(props.workers) ? props.workers : [];
  const messages = Array.isArray(props.messages) ? props.messages : [];

  const header = el("header", "lego-painel__header");
  header.appendChild(el("h3", "a2ui-painel__title", props.title || "Quadro do Swarm"));
  header.appendChild(el("span", "pill", workers.length + " workers"));
  painel.appendChild(header);

  if (props.progress !== undefined && props.progress !== null) {
    const raw = Number(props.progress);
    const pct = Number.isFinite(raw) ? Math.max(0, Math.min(100, Math.round(raw))) : 0;
    const bar = el("progress", "lego-progress__bar");
    bar.max = 100;
    bar.value = pct;
    painel.appendChild(bar);
  }

  const workersSec = el("section", "lego-swarm__section");
  workersSec.appendChild(el("span", "lego-section-label", "Workers"));
  if (workers.length === 0) {
    workersSec.appendChild(el("p", "a2ui-text", "Nenhum worker ativo nesta sessao."));
  } else {
    const wList = el("ul", "lego-swarm__workers");
    for (const w of workers) {
      const name = w && (w.name || w.id) ? String(w.name || w.id) : "worker";
      const state = w && (w.state || w.status) ? String(w.state || w.status) : "waiting";
      const li = el("li", "lego-swarm__worker " + statusClass(state));
      const top = el("div", "lego-painel__header");
      top.appendChild(el("span", "lego-swarm__worker-name", name));
      top.appendChild(el("span", "pill", state));
      li.appendChild(top);
      if (w && (w.task || w.role)) {
        li.appendChild(el("span", "lego-swarm__worker-task", String(w.task || w.role)));
      }
      wList.appendChild(li);
    }
    workersSec.appendChild(wList);
  }
  painel.appendChild(workersSec);

  const msgSec = el("section", "lego-swarm__section");
  msgSec.appendChild(el("span", "lego-section-label", "Mensagens (" + messages.length + ")"));
  if (messages.length === 0) {
    msgSec.appendChild(el("p", "a2ui-text", "Nenhuma mensagem publicada no quadro."));
  } else {
    const mList = el("ul", "lego-swarm__messages");
    for (const m of messages) {
      const li = el("li", "lego-swarm__message");
      const seq = m && m.seq !== undefined ? "#" + String(m.seq) + " " : "";
      const author = m && (m.from || m.sender || m.author) ? String(m.from || m.sender || m.author) : "swarm";
      li.appendChild(el("span", "lego-swarm__msg-meta", seq + author));
      const body = m && (m.text || m.body || m.message) ? String(m.text || m.body || m.message) : "";
      li.appendChild(el("p", "lego-swarm__msg-text", body));
      mList.appendChild(li);
    }
    msgSec.appendChild(mList);
  }
  painel.appendChild(msgSec);

  const draftKey = swarmDraftKey(surface, component);
  const form = el("form", "a2ui-form lego-swarm__composer");
  const label = el("label", "a2ui-label", "Enviar mensagem ao quadro");
  label.htmlFor = "swarm-draft-in";
  const row = el("div", "lego-swarm__composer-row");
  const draftInput = el("input", "a2ui-input lego-swarm__draft");
  draftInput.id = "swarm-draft-in";
  draftInput.type = "text";
  draftInput.setAttribute("placeholder", "Digite uma mensagem para o Swarm...");
  draftInput.value = swarmDrafts.get(draftKey) || "";
  draftInput.addEventListener("input", () => {
    swarmDrafts.set(draftKey, draftInput.value);
  });

  const sendBtn = el("button", "btn lego-swarm__send", "Enviar");
  sendBtn.type = "submit";
  row.append(draftInput, sendBtn);
  form.append(label, row);

  form.addEventListener("submit", (ev) => {
    ev.preventDefault();
    const text = draftInput.value.trim();
    if (!text) {
      return;
    }
    swarmDrafts.delete(draftKey);
    draftInput.value = "";
    emitAction(props.action || "swarm_post", {
      board_id: component && component.id ? component.id : "swarm",
      text
    });
  });

  painel.appendChild(form);
  return painel;
}

let catalogRegistered = false;

export function registerLegoLoopCatalog() {
  if (catalogRegistered) {
    return;
  }
  registerComponent("LegoDagGraph", renderDagGraph);
  registerComponent("LegoProgressMeter", renderProgressMeter);
  registerComponent("LegoStopGatePanel", renderStopGatepainel);
  registerComponent("LegoArtifactViewer", renderArtifactViewer);
  registerComponent("LegoLogStream", renderLogStream);
  registerComponent("LegoSwarmBoard", renderSwarmBoard);
  catalogRegistered = true;
}

registerLegoLoopCatalog();
