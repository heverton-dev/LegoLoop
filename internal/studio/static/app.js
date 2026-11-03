"use strict";

// LegoLoop Studio — canvas N8N leve, sem dependências. Todo texto vindo da API entra
// no DOM via textContent (nunca innerHTML), e todo POST vai como application/json.
// Inspetor: inspector.js · paleta: palette.js · slots: slots.js · diálogos e toasts: ui.js.

const NODE_W = 236;
const PORT_Y = 53; // centro da .port (top 46px + raio 7px)
const SVG_NS = "http://www.w3.org/2000/svg";
const MODE_LABEL = { suggested: "Sugerido", manual: "Manual", auto: "Auto" };
const POLL_MS = 10000;
const ITEM_KINDS = ["skill", "mcp", "script", "hook", "gate", "loop"];
const PIECE_KIND = { skill: "skill", script: "script", mcp: "mcp", loop: "subloop" };
const PIECE_LABEL = { skill: "skill", script: "script", mcp: "MCP", subloop: "sub-loop" };

function toKebabCase(s) {
  return s
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "") || "unnamed-loop";
}

const state = {
  pan: { x: 40, y: 40 },
  zoom: 1.0,
  isPanning: false,
  blocks: new Map(),
  matrix: {},
  types: [],
  inv: { models: [], router_online: false, clis: [], api_keys: [] },
  items: Object.fromEntries(ITEM_KINDS.map((k) => [k, []])),
  itemsVersion: 0,
  nodes: new Map(),
  edges: [],
  selected: null,
  seq: 0,
  status: new Map(),
  telemetry: new Map(),
  navStack: [],
  badEdges: new Set(),
  swapping: null,
  gateDraft: null,
  suggest: new Map(),
  suggestPending: new Set(),
  palette: null,
};

const $ = (id) => document.getElementById(id);

function el(tag, cls, text) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (text !== undefined) e.textContent = text;
  return e;
}

function badge(text, cls) {
  return el("span", `badge ${cls || ""}`.trim(), text);
}

function shortType(t) {
  return String(t).replace("Artifact:", "");
}

// ---------- API ----------

async function api(path, body, signal) {
  const opts = body === undefined ? {} : {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  };
  if (signal) opts.signal = signal;
  const res = await fetch(path, opts);
  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}

// ---------- Inventário e roteamento ----------

function modelIds(source) {
  return state.inv.models.filter((m) => !source || m.source === source).map((m) => m.id);
}

function hasModel(id) {
  return modelIds().includes(id);
}

function defaultRouting(block) {
  const all = modelIds();
  if (hasModel(block.suggested_model)) {
    const fb = all.find((m) => m !== block.suggested_model) || "";
    return { mode: "suggested", model: block.suggested_model, fallback: fb };
  }
  return { mode: "manual", model: all[0] || "", fallback: "" };
}

function routingLabel(r) {
  if (r.mode === "auto") return `Auto · ${r.combo || "sem combo"}`;
  return `${MODE_LABEL[r.mode]} · ${r.model || "sem modelo"}`;
}

// ---------- Peças: bloco tipado do catálogo ou peça genérica do inventário ----------

const ANY_TYPE = "Artifact:Any";

function defaultGate(artifact) {
  return `file_exists('${artifact}')`;
}

// pieceBlock descreve uma peça genérica com o mesmo formato de um bloco do catálogo.
function pieceBlock(node) {
  const p = node.piece;
  return {
    id: `piece:${p.kind}`, label: node.label, ref: "", consumes: [ANY_TYPE], produces: ANY_TYPE,
    artifact: p.artifact, gate: p.gate || defaultGate(p.artifact), generic: true,
  };
}

function blockOf(node) {
  return node.piece ? pieceBlock(node) : state.blocks.get(node.block);
}

function slugOf(name) {
  const s = String(name).replace(/[^A-Za-z0-9_-]+/g, "_").replace(/^_+|_+$/g, "");
  return s || "peca";
}

// uniqueArtifact evita dois nós gravando o mesmo arquivo.
function uniqueArtifact(name, ignoreId) {
  const used = new Set([...state.nodes.values()].filter((n) => n.id !== ignoreId).map((n) => blockOf(n).artifact));
  const base = slugOf(name).toUpperCase();
  let candidate = `${base}.md`;
  for (let i = 2; used.has(candidate); i += 1) candidate = `${base}_${i}.md`;
  return candidate;
}

function newNodeId(label) {
  state.seq += 1;
  return `n${String(state.seq).padStart(2, "0")}_${slugOf(label).toLowerCase().slice(0, 24)}`;
}

function baseNode(id, x, y) {
  return { id, x: Math.max(16, x), y: Math.max(16, y), slots: emptySlots(), execution: "inline" };
}

function register(node) {
  pushState();
  state.nodes.set(node.id, node);
  clearValidation();
  renderNode(node);
  select({ kind: "node", id: node.id });
  updateEmpty();
  return node;
}

// ---------- Nós ----------

function defaultTriggerConfig(kind) {
  switch (kind) {
    case "trigger_cron": return { type: "cron", cron: "0 * * * *" };
    case "trigger_file_change": return { type: "file_change", path: "docs" };
    case "trigger_loop_end": return { type: "loop_end", loop: "" };
    case "trigger_board": return { type: "board", board: "backlog" };
    default: return { type: "manual" };
  }
}

function addNode(blockId, x, y) {
  const block = state.blocks.get(blockId);
  if (!block) return null;
  const isTrig = block.kind && block.kind.startsWith("trigger_");
  if (isTrig) {
    for (const [id, n] of state.nodes.entries()) {
      const b = blockOf(n);
      if (b && b.kind && b.kind.startsWith("trigger_")) {
        removeNode(id);
        toast(`Gatilho substituído por ${block.label}.`);
        break;
      }
    }
  }
  const node = { ...baseNode(newNodeId(blockId), x, y), block: blockId, routing: defaultRouting(block) };
  if (isTrig) {
    node.triggerConfig = defaultTriggerConfig(block.kind);
  }
  return register(node);
}

// addPiece cria um nó genérico a partir de um item do inventário.
function addPiece(entry, x, y) {
  const kind = PIECE_KIND[entry.cat];
  if (!kind) return null;
  const node = {
    ...baseNode(newNodeId(entry.name), x, y), label: entry.name,
    piece: { kind, ref: entry.itemId, artifact: "", gate: "" }, routing: defaultRouting({}),
  };
  node.piece.artifact = uniqueArtifact(entry.name, node.id);
  return register(node);
}

function hasTrigger() {
  for (const n of state.nodes.values()) {
    const b = blockOf(n);
    if (b && b.kind && b.kind.startsWith("trigger_")) return true;
  }
  return false;
}

function isTriggerEntry(entry) {
  if (entry.blockId) {
    const block = state.blocks.get(entry.blockId);
    return Boolean(block && block.kind && block.kind.startsWith("trigger_"));
  }
  if (entry.kind && entry.kind.startsWith("trigger_")) return true;
  if (entry.cat === "trigger") return true;
  return false;
}

// addEntry cria o nó de uma entrada da paleta (bloco tipado ou peça do inventário).
function addEntry(entry, x, y) {
  if (!hasTrigger() && !isTriggerEntry(entry)) {
    toast("Todo loop precisa começar com um Gatilho. Arraste uma peça de Trigger primeiro.", true);
    return null;
  }
  if (entry.kind === "block" || entry.blockId) return addNode(entry.blockId, x, y);
  if (!PIECE_KIND[entry.cat]) {
    toast("Hooks e gates encaixam sobre uma peça. Solte sobre um nó ou selecione um e dê duplo clique.", true);
    return null;
  }
  return addPiece(entry, x, y);
}

function removeNode(id) {
  pushState();
  state.nodes.delete(id);
  state.edges = state.edges.filter((e) => e.from !== id && e.to !== id);
  document.querySelector(`.node[data-id="${CSS.escape(id)}"]`)?.remove();
  if (state.swapping === id) cancelSwap();
  clearValidation();
  renderEdges();
  updateEmpty();
}

function nodeCategoryClass(block) {
  const k = block.kind || "";
  if (k.startsWith("trigger_")) return "node--trigger";
  if (k.startsWith("flow_")) return "node--control";
  if (k === "mcp") return "node--mcp";
  if (k === "subloop" || k === "loop") return "node--loop";
  if (k === "script") return "node--script";
  if (k === "hook") return "node--hook";
  if (k === "gate") return "node--gate";
  return "node--skill";
}

function nodeEl(node) {
  const block = blockOf(node);
  const isTrig = block.kind && block.kind.startsWith("trigger_");
  const isSubloop = (node.piece && (node.piece.kind === "subloop" || node.piece.kind === "loop")) || (block.kind && (block.kind === "loop" || block.kind === "flow_loop"));
  const catClass = nodeCategoryClass(block);
  const root = el("div", `node ${catClass}${isSubloop ? " node--macro" : ""}`);
  root.id = `node-${node.id}`;
  root.dataset.id = node.id;
  root.tabIndex = 0;
  root.setAttribute("role", "group");
  root.setAttribute("aria-label", `Nó ${block.label}`);

  if (isSubloop) {
    root.addEventListener("dblclick", (ev) => {
      ev.stopPropagation();
      enterSubloop(node.piece?.ref || block.id, block.label);
    });
  }

  const head = el("div", "node__head");
  head.append(el("span", "node__title", isTrig ? `[Gatilho] ${block.label}` : block.label), el("span", "node__id", node.id));
  if (isSubloop) {
    head.append(el("span", "node__macro-hint", "Duplo clique para abrir"));
  }

  const ports = el("div", "node__ports");
  const inLabel = isTrig ? "[Gatilho]" : (block.consumes || []).map(shortType).join(" | ");
  ports.append(
    el("span", "node__port-label", inLabel),
    el("span", "node__port-label node__port-label--out", shortType(block.produces)),
  );

  const portIn = el("span", "port port--in");
  portIn.dataset.port = "in";
  portIn.title = isTrig ? "Gatilho inicial (sem entrada)" : `Entrada: ${(block.consumes || []).map(shortType).join(", ")}`;
  const portOut = el("span", "port port--out");
  portOut.dataset.port = "out";
  portOut.title = `Saída: ${shortType(block.produces)} — arraste para conectar`;

  const children = [head, ports, el("div", "node__badges")];
  if (!isTrig) children.push(portIn);
  children.push(portOut);
  root.append(...children);
  return root;
}

function renderNode(node) {
  let root = document.querySelector(`.node[data-id="${CSS.escape(node.id)}"]`);
  if (!root) {
    root = nodeEl(node);
    $("nodes").append(root);
  }
  const block = blockOf(node);
  root.querySelector(".node__title").textContent = block.label;
  root.style.transform = `translate(${node.x}px, ${node.y}px)`;
  root.classList.toggle("is-selected", state.selected?.kind === "node" && state.selected.id === node.id);
  root.classList.toggle("is-fail", state.status.get(node.id) === "fail");
  root.classList.toggle("is-swapping", state.swapping === node.id);
  renderBadges(node, root.querySelector(".node__badges"));

  // Telemetria por card em tempo real
  let telemBox = root.querySelector(".node__telemetry");
  const telemData = state.telemetry?.get(node.id);
  if (telemData) {
    if (!telemBox) {
      telemBox = el("div", "node__telemetry");
      root.append(telemBox);
    }
    const isRun = telemData.status === "running";
    const isPass = telemData.status === "passed";
    const statusCls = isRun ? "is-running" : isPass ? "is-passed" : "is-failed";
    const statusText = isRun ? "RUN" : isPass ? "OK" : "FAIL";

    const spanStatus = el("span", `node__telemetry-status ${statusCls}`, statusText);
    const spanTime = el("span", "node__telemetry-time", telemData.duration_ms !== undefined ? `${telemData.duration_ms}ms` : "...");
    const tokensTxt = telemData.tokens ? `In:${telemData.tokens.input} Out:${telemData.tokens.output}` : "";
    const spanTokens = el("span", "node__telemetry-tokens", tokensTxt);

    telemBox.replaceChildren(spanStatus, spanTime, spanTokens);
  } else if (telemBox) {
    telemBox.remove();
  }
}

function pill(cls, text, title) {
  const p = el("span", `pill ${cls}`);
  if (title) p.title = title;
  p.append(el("span", "pill__dot"), el("span", "pill__text", text));
  return p;
}

function renderBadges(node, box) {
  const block = blockOf(node);
  const status = state.status.get(node.id);
  const gate = status === "ok" ? pill("pill--ok", "gate válido", block.gate)
    : status === "fail" ? pill("pill--fail", "com erro", block.gate)
    : pill("", "gate pendente", block.gate);
  const pills = [];
  if (node.piece) pills.push(pill("pill--kind", PIECE_LABEL[node.piece.kind] || node.piece.kind, "Peça do inventário (entrada e saída: Qualquer)"));
  pills.push(pill("pill--handoff", block.artifact, "Handoff imutável"), pill("pill--route", routingLabel(node.routing), "Roteamento"), gate);
  const slots = slotCount(node.slots);
  if (slots > 0) pills.push(pill("pill--slots", `${slots} slot${slots > 1 ? "s" : ""}`, "Gates, scripts e hooks encaixados"));
  if (node.execution === "isolated") pills.push(pill("pill--exec", "isolada", "Executa em workspace separado"));
  box.replaceChildren(...pills);
}

function renderAllNodes() {
  state.nodes.forEach(renderNode);
  renderEdges();
}

function updateEmpty() {
  $("empty").hidden = state.nodes.size > 0;
}

// ---------- Sub-grafos em Mega-Pipelines (Drill-Down & Breadcrumbs) ----------

function serializeCanvasSnapshot() {
  return {
    name: $("pipeline-name")?.value || "esteira",
    nodes: Array.from(state.nodes.values()).map((n) => JSON.parse(JSON.stringify(n))),
    edges: JSON.parse(JSON.stringify(state.edges)),
    status: new Map(state.status),
    telemetry: new Map(state.telemetry),
  };
}

function restoreCanvasSnapshot(snap) {
  if (!snap) return;
  [...state.nodes.keys()].forEach(removeNode);

  if ($("pipeline-name") && snap.name) {
    $("pipeline-name").value = snap.name;
  }

  (snap.nodes || []).forEach((n) => {
    state.nodes.set(n.id, n);
    renderNode(n);
  });

  state.edges = snap.edges || [];
  state.status = snap.status ? new Map(snap.status) : new Map();
  state.telemetry = snap.telemetry ? new Map(snap.telemetry) : new Map();

  renderAllNodes();
  updateEmpty();
}

async function enterSubloop(subloopId, subloopLabel) {
  if (!subloopId) return;
  const currentSnap = serializeCanvasSnapshot();
  currentSnap.label = currentSnap.name || "Principal";
  state.navStack.push(currentSnap);

  try {
    const resp = await fetch(`/api/loops/${encodeURIComponent(subloopId)}`);
    if (!resp.ok) {
      toast("Não foi possível carregar o sub-grafo do acervo", true);
      state.navStack.pop();
      return;
    }
    const loopData = await resp.json();

    [...state.nodes.keys()].forEach(removeNode);

    if ($("pipeline-name")) {
      $("pipeline-name").value = loopData.name || subloopId;
    }

    const nodeMap = new Map();
    (loopData.nodes || []).forEach((n, i) => {
      const added = addNode(n.block || n.id, 60 + i * (NODE_W + 80), 100 + (i % 2) * 40);
      if (added) {
        nodeMap.set(n.id, added.id);
      }
    });

    (loopData.edges || []).forEach((e) => {
      const fromId = nodeMap.get(e.from);
      const toId = nodeMap.get(e.to);
      if (fromId && toId) {
        addEdge(fromId, toId);
      }
    });

    renderBreadcrumbs();
    autoLayout();
    fitToScreen();
    toast(`Sub-grafo aberto: ${subloopLabel || loopData.name || subloopId}`);
  } catch (err) {
    toast("Erro ao abrir sub-grafo: " + err.message, true);
    state.navStack.pop();
  }
}

function leaveSubloop(targetIndex) {
  if (state.navStack.length === 0) return;
  const idx = targetIndex !== undefined ? targetIndex : state.navStack.length - 1;
  const targetSnap = state.navStack[idx];
  state.navStack = state.navStack.slice(0, idx);

  restoreCanvasSnapshot(targetSnap);
  renderBreadcrumbs();
  toast(`Retornou para: ${targetSnap.label}`);
}

function renderBreadcrumbs() {
  const container = $("canvas-breadcrumbs");
  const list = $("canvas-breadcrumbs-list");
  if (!container || !list) return;

  if (state.navStack.length === 0) {
    container.hidden = true;
    list.replaceChildren();
    return;
  }

  container.hidden = false;
  const items = [];

  state.navStack.forEach((snap, idx) => {
    const li = document.createElement("li");
    li.className = "breadcrumb-item";

    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "breadcrumb-link";
    btn.textContent = snap.label || `Nível ${idx}`;
    btn.addEventListener("click", () => leaveSubloop(idx));

    const sep = document.createElement("span");
    sep.className = "breadcrumb-separator";
    sep.textContent = "›";

    li.append(btn, sep);
    items.push(li);
  });

  const currentLi = document.createElement("li");
  currentLi.className = "breadcrumb-item";
  const currentSpan = document.createElement("span");
  currentSpan.className = "breadcrumb-current";
  currentSpan.textContent = $("pipeline-name")?.value || "Sub-grafo";
  currentLi.append(currentSpan);
  items.push(currentLi);

  list.replaceChildren(...items);
}

// Expõe para o inspetor
window.enterSubloop = enterSubloop;


// ---------- Arestas ----------

function portPoint(node, side) {
  return { x: node.x + (side === "out" ? NODE_W : 0), y: node.y + PORT_Y };
}

function curve(a, b) {
  const dx = Math.max(60, Math.abs(b.x - a.x) * 0.5);
  return `M ${a.x} ${a.y} C ${a.x + dx} ${a.y}, ${b.x - dx} ${b.y}, ${b.x} ${b.y}`;
}

// Peça genérica aceita e produz "Qualquer": liga com tudo. Entre blocos tipados vale a matriz.
function canConnect(fromId, toId) {
  const a = state.nodes.get(fromId);
  const b = state.nodes.get(toId);
  if (!a || !b || fromId === toId) return false;
  if (a.piece || b.piece) return true;
  return (state.matrix[a.block] || []).includes(b.block);
}

function edgeClass(e, i) {
  const sel = state.selected?.kind === "edge" && state.selected.id === i;
  const bad = state.badEdges.has(`${e.from}→${e.to}`) || !canConnect(e.from, e.to);
  return `edge${sel ? " is-selected" : ""}${bad ? " is-invalid" : ""}`;
}

function edgeHandle(point, i) {
  const h = document.createElementNS(SVG_NS, "circle");
  h.setAttribute("class", "edge-handle");
  h.setAttribute("r", "7");
  h.setAttribute("cx", String(point.x));
  h.setAttribute("cy", String(point.y));
  h.dataset.edge = String(i);
  const t = document.createElementNS(SVG_NS, "title");
  t.textContent = "Arraste para religar; solte no vazio para remover";
  h.append(t);
  return h;
}

function renderEdges() {
  const svg = $("edges");
  const temp = svg.querySelector(".edge--temp");
  svg.replaceChildren();
  state.edges.forEach((e, i) => {
    const a = state.nodes.get(e.from);
    const b = state.nodes.get(e.to);
    const end = portPoint(b, "in");
    const d = curve(portPoint(a, "out"), end);
    const hit = document.createElementNS(SVG_NS, "path");
    hit.setAttribute("d", d);
    hit.setAttribute("class", "edge-hit");
    hit.dataset.edge = String(i);
    const path = document.createElementNS(SVG_NS, "path");
    path.setAttribute("d", d);
    path.setAttribute("class", edgeClass(e, i) + " edge--animated");
    path.dataset.edge = String(i);
    svg.append(path, hit);

    // Botão de desconexão rápida no meio da aresta com área de clique aprimorada
    const pA = portPoint(a, "out");
    const midX = (pA.x + end.x) / 2;
    const midY = (pA.y + end.y) / 2;
    const btnDel = document.createElementNS(SVG_NS, "g");
    btnDel.setAttribute("class", "edge-del-group");
    btnDel.dataset.edge = String(i);
    btnDel.style.cursor = "pointer";
    const circ = document.createElementNS(SVG_NS, "circle");
    circ.setAttribute("cx", String(midX));
    circ.setAttribute("cy", String(midY));
    circ.setAttribute("r", "12");
    circ.setAttribute("class", "edge-btn-remove");
    circ.dataset.edge = String(i);
    circ.style.cursor = "pointer";
    const txt = document.createElementNS(SVG_NS, "text");
    txt.setAttribute("x", String(midX));
    txt.setAttribute("y", String(midY + 4));
    txt.setAttribute("text-anchor", "middle");
    txt.setAttribute("font-size", "12");
    txt.setAttribute("font-weight", "bold");
    txt.setAttribute("fill", "#fff");
    txt.setAttribute("pointer-events", "none");
    txt.textContent = "x";
    btnDel.append(circ, txt);
    svg.append(btnDel);

    if (state.selected?.kind === "edge" && state.selected.id === i) svg.append(edgeHandle(end, i));
  });
  if (temp) svg.append(temp);
}

function addEdge(from, to) {
  if (state.edges.some((e) => e.from === from && e.to === to)) {
    toast("Essa conexão já existe.", true);
    return;
  }
  if (!canConnect(from, to)) {
    const a = blockOf(state.nodes.get(from));
    const b = blockOf(state.nodes.get(to));
    toast(`${a.label} produz ${shortType(a.produces)}; ${b.label} não consome isso.`, true);
    return;
  }
  pushState();
  state.edges.push({ from, to });
  clearValidation();
  renderEdges();
}

function disconnect(index) {
  pushState();
  state.edges.splice(index, 1);
  clearValidation();
  select(null);
  toast("Conexão removida.");
}

// ---------- Seleção ----------

function select(sel) {
  state.selected = sel;
  renderAllNodes();
  renderInspector();
}

// ---------- Trocar peça ----------

function startSwap(node) {
  state.swapping = node.id;
  state.palette.setBanner(`Escolha na paleta a peça que substitui "${blockOf(node).label}". Esc cancela.`);
  renderNode(node);
}

function cancelSwap() {
  const id = state.swapping;
  state.swapping = null;
  state.palette?.setBanner(null);
  const node = state.nodes.get(id);
  if (node) renderNode(node);
}

// swapNode troca a peça mantendo id, posição, conexões, roteamento válido, slots e modo de execução.
// Conexões que deixam de ser compatíveis NÃO são apagadas: ficam vermelhas e bloqueiam a materialização.
function swapNode(id, entry) {
  const node = state.nodes.get(id);
  const isBlock = entry.kind === "block" || Boolean(entry.blockId);
  if (!isBlock && !PIECE_KIND[entry.cat]) {
    toast("Só skills, MCPs, scripts e loops substituem uma peça.", true);
    return;
  }
  if (isBlock) {
    node.block = entry.blockId;
    delete node.piece;
    delete node.label;
  } else {
    delete node.block;
    node.label = entry.name;
    node.piece = { kind: PIECE_KIND[entry.cat], ref: entry.itemId, artifact: uniqueArtifact(entry.name, node.id), gate: "" };
  }
  const routingOk = node.routing.mode === "auto" ? hasModel(node.routing.combo) : hasModel(node.routing.model);
  if (!routingOk) node.routing = defaultRouting(blockOf(node));
  cancelSwap();
  clearValidation();
  const broken = state.edges.filter((e) => (e.from === id || e.to === id) && !canConnect(e.from, e.to)).length;
  select({ kind: "node", id });
  toast(broken ? `Peça trocada. ${broken} conexão(ões) ficaram incompatíveis (em vermelho): religue ou desconecte.` : "Peça trocada; as conexões foram mantidas.", broken > 0);
}

// ---------- Encaixar hook, gate ou script numa peça ----------

async function attachToNode(node, payload) {
  const slots = node.slots;
  if (payload.cat === "gate") {
    const expr = instantiateGate(payload.gateName, blockOf(node).artifact);
    if (!expr) return toast("Não foi possível montar este gate para o artefato da peça.", true);
    if (payload.gateName === "has_pattern") state.gateDraft = { nodeId: node.id, expr };
    else if (!slots.gates.includes(expr)) slots.gates.push(expr);
    else return toast("Esse gate já está na peça.");
  } else if (payload.cat === "script") {
    if (!slots.scripts.some((s) => s.ref === payload.itemId)) slots.scripts.push({ ref: payload.itemId, label: payload.name });
  } else if (payload.cat === "hook") {
    const where = await hookTarget(payload);
    if (!where) return undefined;
    if (!slots[where].some((s) => s.ref === payload.itemId)) slots[where].push({ ref: payload.itemId, label: payload.name });
  }
  clearValidation();
  select({ kind: "node", id: node.id });
  return toast(payload.gateName === "has_pattern" ? "Complete a expressão do gate e confirme." : "Encaixado na peça.");
}

// hookTarget decide pré ou pós pelo evento; só pergunta quando os dois servem.
async function hookTarget(payload) {
  const s = hookSlotsFor(payload.event);
  if (s.pre && !s.post) return "pre";
  if (s.post && !s.pre) return "post";
  return choiceDialog({
    title: "Onde encaixar o hook?",
    message: `"${payload.name}" serve antes e depois da peça.`,
    choices: [{ label: "Pré (antes)", value: "pre" }, { label: "Pós (depois)", value: "post" }],
  });
}

// addNeighbor cria a peça vizinha à direita e já liga a saída da atual à entrada dela.
function addNeighbor(node, n) {
  const x = node.x + NODE_W + 72;
  const next = n.kind === "mcp"
    ? addPiece({ cat: "mcp", name: n.name, itemId: n.id }, x, node.y)
    : addNode(n.id, x, node.y);
  if (next) addEdge(node.id, next.id);
}

// ---------- Interações de ponteiro ----------

const MAX_PAN = 10000;

function updateCanvasTransform() {
  const surface = $("canvas-surface");
  if (surface) {
    state.pan.x = Math.max(-MAX_PAN, Math.min(MAX_PAN, state.pan.x));
    state.pan.y = Math.max(-MAX_PAN, Math.min(MAX_PAN, state.pan.y));
    surface.style.transform = `translate(${state.pan.x}px, ${state.pan.y}px) scale(${state.zoom})`;
    surface.style.transformOrigin = "0 0";
  }
  const zl = $("zoom-level");
  if (zl) zl.textContent = `${Math.round(state.zoom * 100)}%`;
}

function zoomBy(factor, center) {
  const c = $("canvas");
  const r = c.getBoundingClientRect();
  const cx = center ? center.x : r.width / 2;
  const cy = center ? center.y : r.height / 2;
  const oldZoom = state.zoom;
  const newZoom = Math.min(3.0, Math.max(0.2, oldZoom * factor));
  state.pan.x = cx - (cx - state.pan.x) * (newZoom / oldZoom);
  state.pan.y = cy - (cy - state.pan.y) * (newZoom / oldZoom);
  state.zoom = newZoom;
  updateCanvasTransform();
}

function fitToScreen() {
  if (state.nodes.size === 0) {
    state.pan = { x: 40, y: 40 };
    state.zoom = 1.0;
    updateCanvasTransform();
    return;
  }
  const c = $("canvas");
  const r = c.getBoundingClientRect();
  let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
  for (const n of state.nodes.values()) {
    minX = Math.min(minX, n.x);
    minY = Math.min(minY, n.y);
    maxX = Math.max(maxX, n.x + NODE_W);
    maxY = Math.max(maxY, n.y + 120);
  }
  const pad = 60;
  const graphW = Math.max(100, maxX - minX + pad * 2);
  const graphH = Math.max(100, maxY - minY + pad * 2);
  const zoom = Math.min(1.5, Math.max(0.25, Math.min(r.width / graphW, r.height / graphH)));

  state.zoom = zoom;
  state.pan.x = (r.width - (maxX - minX) * zoom) / 2 - minX * zoom;
  state.pan.y = (r.height - (maxY - minY) * zoom) / 2 - minY * zoom;
  updateCanvasTransform();
}

function toggleSidebar() {
  const ws = $("workspace");
  const sb = $("sidebar");
  const isCollapsed = sb.classList.toggle("is-collapsed");
  ws.classList.toggle("is-sidebar-collapsed", isCollapsed);
  const btn = $("btn-toggle-sidebar");
  if (btn) btn.textContent = isCollapsed ? "▸" : "◂";
}

function toggleInspector() {
  const ws = $("workspace");
  const insp = $("inspector");
  const isCollapsed = insp.classList.toggle("is-collapsed");
  ws.classList.toggle("is-inspector-collapsed", isCollapsed);
  const btn = $("btn-toggle-inspector");
  if (btn) btn.textContent = isCollapsed ? "◂" : "▸";
}

// Histórico de Ações (Undo / Redo)
const history = { undo: [], redo: [] };

function pushState() {
  const snap = JSON.stringify({
    nodes: Array.from(state.nodes.entries()),
    edges: state.edges,
    nextId: state.nextId
  });
  history.undo.push(snap);
  if (history.undo.length > 50) history.undo.shift();
  history.redo = [];
}

function applySnapshot(snapStr) {
  try {
    const snap = JSON.parse(snapStr);
    state.nodes = new Map(snap.nodes);
    state.edges = snap.edges || [];
    state.nextId = snap.nextId || state.nextId;
    document.querySelectorAll(".node").forEach((el) => {
      const id = el.dataset.id;
      if (!state.nodes.has(id)) el.remove();
    });
    renderAllNodes();
    updateEmpty();
  } catch (e) {
    console.error("Falha ao restaurar histórico:", e);
  }
}

function undo() {
  if (history.undo.length === 0) return toast("Nada para desfazer.");
  const cur = JSON.stringify({
    nodes: Array.from(state.nodes.entries()),
    edges: state.edges,
    nextId: state.nextId
  });
  history.redo.push(cur);
  const prev = history.undo.pop();
  applySnapshot(prev);
  toast("Ação desfeita.");
}

function redo() {
  if (history.redo.length === 0) return toast("Nada para refazer.");
  const cur = JSON.stringify({
    nodes: Array.from(state.nodes.entries()),
    edges: state.edges,
    nextId: state.nextId
  });
  history.undo.push(cur);
  const next = history.redo.pop();
  applySnapshot(next);
  toast("Ação refeita.");
}

function autoLayout() {
  if (state.nodes.size === 0) return toast("Canvas vazio.");
  pushState();
  const nodes = Array.from(state.nodes.values());
  const inDegree = new Map();
  nodes.forEach(n => inDegree.set(n.id, 0));
  state.edges.forEach(e => {
    inDegree.set(e.to, (inDegree.get(e.to) || 0) + 1);
  });
  // Níveis de profundidade
  const levels = new Map();
  const queue = nodes.filter(n => inDegree.get(n.id) === 0);
  queue.forEach(n => levels.set(n.id, 0));
  if (queue.length === 0 && nodes.length > 0) {
    queue.push(nodes[0]);
    levels.set(nodes[0].id, 0);
  }
  let head = 0;
  while (head < queue.length) {
    const cur = queue[head++];
    const curLvl = levels.get(cur.id);
    state.edges.filter(e => e.from === cur.id).forEach(e => {
      const targetLvl = Math.max(levels.get(e.to) || 0, curLvl + 1);
      levels.set(e.to, targetLvl);
      const targetNode = state.nodes.get(e.to);
      if (targetNode && !queue.includes(targetNode)) queue.push(targetNode);
    });
  }
  nodes.forEach(n => { if (!levels.has(n.id)) levels.set(n.id, 0); });
  // Agrupar por nível
  const byLevel = new Map();
  nodes.forEach(n => {
    const lvl = levels.get(n.id);
    if (!byLevel.has(lvl)) byLevel.set(lvl, []);
    byLevel.get(lvl).push(n);
  });
  // Reposicionar
  byLevel.forEach((col, lvl) => {
    col.forEach((n, idx) => {
      n.x = 60 + lvl * (NODE_W + 90);
      n.y = 80 + idx * 160;
      const elNode = document.querySelector(`.node[data-id="${n.id}"]`);
      if (elNode) {
        elNode.style.left = `${n.x}px`;
        elNode.style.top = `${n.y}px`;
      }
    });
  });
  renderEdges();
  fitToScreen();
  toast("Esteira organizada automaticamente!");
}

function toggleAccordion(name) {
  const btnPal = $("btn-toggle-palette");
  const btnMod = $("btn-toggle-models");
  const secPal = $("section-palette");
  const secMod = $("section-models");

  if (name === "palette") {
    btnPal?.classList.add("is-active");
    btnPal?.setAttribute("aria-pressed", "true");
    btnMod?.classList.remove("is-active");
    btnMod?.setAttribute("aria-pressed", "false");

    secPal?.classList.add("is-active");
    secPal?.removeAttribute("hidden");
    secMod?.classList.remove("is-active");
    secMod?.setAttribute("hidden", "");
  } else if (name === "models") {
    btnMod?.classList.add("is-active");
    btnMod?.setAttribute("aria-pressed", "true");
    btnPal?.classList.remove("is-active");
    btnPal?.setAttribute("aria-pressed", "false");

    secMod?.classList.add("is-active");
    secMod?.removeAttribute("hidden");
    secPal?.classList.remove("is-active");
    secPal?.setAttribute("hidden", "");
  }
}

function canvasPoint(ev) {
  const c = $("canvas");
  const r = c.getBoundingClientRect();
  const screenX = ev.clientX - r.left;
  const screenY = ev.clientY - r.top;
  return {
    x: (screenX - state.pan.x) / state.zoom,
    y: (screenY - state.pan.y) / state.zoom,
  };
}

function startNodeDrag(ev, node) {
  const start = canvasPoint(ev);
  const origin = { x: node.x, y: node.y };
  const root = $(`node-${node.id}`) || document.querySelector(`.node[data-id="${CSS.escape(node.id)}"]`);
  if (root) {
    root.classList.add("is-dragging");
    try { root.setPointerCapture(ev.pointerId); } catch (_) {}
  }

  const move = (e) => {
    const p = canvasPoint(e);
    node.x = Math.max(8, origin.x + p.x - start.x);
    node.y = Math.max(8, origin.y + p.y - start.y);
    renderNode(node);
    renderEdges();
  };
  const up = () => {
    if (root) {
      root.classList.remove("is-dragging");
      try { root.releasePointerCapture(ev.pointerId); } catch (_) {}
    }
    window.removeEventListener("pointermove", move);
    window.removeEventListener("pointerup", up);
    if (node.x !== origin.x || node.y !== origin.y) {
      pushState();
    }
  };
  window.addEventListener("pointermove", move);
  window.addEventListener("pointerup", up);
}

// trackLink desenha a curva temporária de uma conexão sendo arrastada a partir de fromId.
function trackLink(ev, fromId, onDrop) {
  const from = state.nodes.get(fromId);
  const temp = document.createElementNS(SVG_NS, "path");
  temp.setAttribute("class", "edge edge--temp");
  $("edges").append(temp);
  const canvas = $("canvas");
  canvas.classList.add("is-linking");
  document.querySelectorAll(".node").forEach((n) => {
    n.querySelector(".port--in").classList.toggle("is-compatible", canConnect(fromId, n.dataset.id));
  });
  const nodeAt = (e) => document.elementFromPoint(e.clientX, e.clientY)?.closest(".node");
  const move = (e) => {
    const target = nodeAt(e);
    temp.classList.toggle("is-bad", Boolean(target) && !canConnect(fromId, target.dataset.id));
    temp.setAttribute("d", curve(portPoint(from, "out"), canvasPoint(e)));
  };
  const cleanup = () => {
    window.removeEventListener("pointermove", move);
    window.removeEventListener("pointerup", up);
    window.removeEventListener("keydown", onKeyDown);
    temp.remove();
    canvas.classList.remove("is-linking");
  };
  const onKeyDown = (e) => {
    if (e.key === "Escape") {
      cleanup();
      onDrop(null);
    }
  };
  const up = (e) => {
    cleanup();
    onDrop(nodeAt(e)?.dataset.id || null);
  };
  window.addEventListener("pointermove", move);
  window.addEventListener("pointerup", up);
  window.addEventListener("keydown", onKeyDown);
  move(ev);
}

function startLink(ev, fromId) {
  trackLink(ev, fromId, (targetId) => {
    if (targetId && targetId !== fromId) addEdge(fromId, targetId);
  });
}

// startEdgeRetarget arrasta a ponta de uma conexão: solta noutra peça religa, solta no vazio remove.
function startEdgeRetarget(ev, index) {
  const edge = state.edges[index];
  ev.preventDefault();
  trackLink(ev, edge.from, (targetId) => {
    if (!targetId) return disconnect(index);
    if (targetId === edge.to) return renderEdges();
    const dup = state.edges.some((e, i) => i !== index && e.from === edge.from && e.to === targetId);
    if (dup || !canConnect(edge.from, targetId)) {
      toast(dup ? "Essa conexão já existe; a original foi mantida." : "Conexão incompatível; a original foi mantida.", true);
      return renderEdges();
    }
    edge.to = targetId;
    clearValidation();
    renderEdges();
    renderInspector();
    return toast("Conexão religada.");
  });
}

function onNodesPointerDown(ev) {
  const nodeRoot = ev.target.closest(".node");
  if (!nodeRoot || ev.button !== 0) return;
  const node = state.nodes.get(nodeRoot.dataset.id);
  if (!node) return;

  if (ev.target.dataset.port === "out") {
    ev.preventDefault();
    startLink(ev, node.id);
    return;
  }

  // Seleciona o nó em qualquer clique não interativo
  if (state.selected?.id !== node.id) {
    select({ kind: "node", id: node.id });
  }

  // Inicia drag a menos que clicou em formulários, botões ou slots
  if (!ev.target.closest("button, input, select, textarea, .slot-pill, .port")) {
    ev.preventDefault();
    startNodeDrag(ev, node);
  }
}

function onEdgesPointerDown(ev) {
  if (ev.button === 0 && ev.target.classList.contains("edge-handle")) startEdgeRetarget(ev, Number(ev.target.dataset.edge));
}

function onEdgesClick(ev) {
  const del = ev.target.closest(".edge-del-group") || (ev.target.classList?.contains("edge-btn-remove") ? ev.target : null);
  if (del) {
    ev.stopPropagation();
    const edgeIdx = Number(del.dataset.edge);
    if (!isNaN(edgeIdx)) disconnect(edgeIdx);
    return;
  }
  const i = ev.target.dataset.edge;
  if (i !== undefined && !ev.target.classList.contains("edge-handle")) select({ kind: "edge", id: Number(i) });
}

function onKey(ev) {
  if (ev.key === "Escape") {
    if (state.swapping) {
      cancelSwap();
    } else if (state.selected) {
      select(null);
    }
    return;
  }
  if (ev.target.closest("input, select, textarea")) return;
  if (ev.key !== "Delete" && ev.key !== "Backspace") return;
  const sel = state.selected;
  if (sel?.kind === "node") removeNode(sel.id);
  if (sel?.kind === "edge") {
    state.edges.splice(sel.id, 1);
    clearValidation();
    renderEdges();
    renderInspector();
  }
  if (sel) select(null);
}

// ---------- Soltar peças da paleta ----------

function readPayload(dt) {
  try {
    return JSON.parse(dt.getData(PIECE_MIME));
  } catch {
    return null;
  }
}

function markDropTarget(target) {
  document.querySelectorAll(".node.is-droptarget").forEach((n) => { if (n !== target) n.classList.remove("is-droptarget"); });
  target?.classList.add("is-droptarget");
}

function setupDropZone() {
  const canvas = $("canvas");
  canvas.addEventListener("dragover", (e) => {
    if (!e.dataTransfer.types.includes(PIECE_MIME)) return;
    e.preventDefault();
    e.dataTransfer.dropEffect = "copy";
    canvas.classList.add("is-dragover");
    markDropTarget(e.target.closest(".node"));
  });
  canvas.addEventListener("dragleave", (e) => {
    if (e.target === canvas) canvas.classList.remove("is-dragover");
  });
  canvas.addEventListener("drop", (e) => {
    e.preventDefault();
    canvas.classList.remove("is-dragover");
    markDropTarget(null);
    const payload = readPayload(e.dataTransfer);
    if (!payload) return;
    const targetNode = e.target.closest(".node");
    if (targetNode && ["hook", "gate", "script"].includes(payload.cat)) {
      attachToNode(state.nodes.get(targetNode.dataset.id), payload);
      return;
    }
    // Inserção intermediária se soltar sobre aresta existente
    const edgeIdx = e.target.dataset.edge;
    if (edgeIdx !== undefined && !targetNode) {
      const idx = Number(edgeIdx);
      const edge = state.edges[idx];
      if (edge) {
        state.edges.splice(idx, 1);
        const p = canvasPoint(e);
        const newNode = addEntry(payload, p.x - NODE_W / 2, p.y - 50);
        if (newNode) {
          addEdge(edge.from, newNode.id);
          addEdge(newNode.id, edge.to);
          toast("Peça inserida entre os blocos conectados.");
          return;
        }
      }
    }
    const p = canvasPoint(e);
    addEntry(payload, p.x - NODE_W / 2, p.y - 50);
  });
  canvas.addEventListener("pointerdown", (e) => {
    if (e.target === canvas || e.target.id === "edges" || e.target.id === "nodes") select(null);
  });
}

// ---------- Validação e materialização ----------

function composition() {
  let trig = undefined;
  for (const n of state.nodes.values()) {
    const b = blockOf(n);
    if (b && b.kind && b.kind.startsWith("trigger_")) {
      trig = n.triggerConfig || defaultTriggerConfig(b.kind);
      break;
    }
  }
  return {
    name: $("pipeline-name").value.trim(),
    trigger: trig,
    nodes: [...state.nodes.values()].map(compNode),
    edges: state.edges.map((e) => ({ from: e.from, to: e.to })),
  };
}

function compNode(n) {
  const out = { id: n.id, routing: n.routing };
  if (n.piece) out.piece = { kind: n.piece.kind, ref: n.piece.ref, artifact: n.piece.artifact, ...(n.piece.gate ? { gate: n.piece.gate } : {}) };
  else out.block = n.block;
  const slots = slotsToComposition(n.slots);
  if (slots) out.slots = slots;
  if (n.execution === "isolated") out.execution = { mode: "isolated" };
  return out;
}

function clearValidation() {
  if (state.status.size === 0 && state.badEdges.size === 0) return;
  state.status.clear();
  state.badEdges.clear();
  $("issues").hidden = true;
  renderAllNodes();
}

function applyValidation(res) {
  state.status.clear();
  state.badEdges.clear();
  const issues = res.issues || [];
  issues.forEach((i) => {
    if (i.node) state.status.set(i.node, "fail");
    if (i.edge) state.badEdges.add(i.edge);
  });
  if (res.ok) state.nodes.forEach((_, id) => state.status.set(id, "ok"));
  renderAllNodes();
  renderIssues(issues);
}

function renderIssues(issues) {
  const box = $("issues");
  box.hidden = issues.length === 0;
  $("issues-list").replaceChildren(...issues.map((i) => {
    const li = el("li");
    li.append(el("span", "issues__where", i.node || i.edge || "esteira"), el("span", "", i.message));
    return li;
  }));
}

async function withBusy(btn, fn) {
  btn.disabled = true;
  try {
    await fn();
  } catch (err) {
    toast(`Falha de comunicação com o servidor local: ${err.message}`, true);
  } finally {
    btn.disabled = false;
  }
}

async function validate() {
  const { ok, data } = await api("/api/validate", composition());
  if (!ok) {
    toast(data.error || "Erro ao validar.", true);
    return false;
  }
  applyValidation(data);
  toast(data.ok ? "Esteira válida: todos os encaixes conferem." : `${data.issues.length} problema(s) encontrado(s).`, !data.ok);
  return Boolean(data.ok);
}

async function materialize() {
  const { ok, data } = await api("/api/materialize", composition());
  if (data.issues) applyValidation(data);
  if (!ok) {
    toast(data.error || "Corrija os problemas antes de materializar.", true);
    return false;
  }
  toast(`Esteira materializada em ${data.path}.`);
  return true;
}

// ---------- Rota ótima ----------

async function traceRoute() {
  const from = $("route-from").value;
  const to = $("route-to").value;
  const { ok, data } = await api(`/api/route?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`);
  if (!ok) return toast(data.error || "Sem rota.", true);
  if (data.path.length === 0) return toast("Entrada e objetivo são o mesmo artefato.");
  if (state.nodes.size > 0 && !(await confirmDialog({
    title: "Substituir o canvas?",
    message: "Os nós e conexões atuais serão removidos e trocados pela rota ótima.",
    confirmLabel: "Substituir",
    danger: true,
  }))) return undefined;
  [...state.nodes.keys()].forEach(removeNode);
  let prev = null;
  data.path.forEach((blockId, i) => {
    const n = addNode(blockId, 40 + i * (NODE_W + 72), 80 + (i % 2) * 36);
    if (prev) addEdge(prev.id, n.id);
    prev = n;
  });
  select(null);
  return toast(`Rota ótima: ${data.path.length} bloco(s), custo ${data.cost}k tokens.`);
}

// ---------- Paleta e inventário ----------

function addEntryToCenter(entry) {
  const cv = $("canvas");
  const cx = (-state.pan.x + cv.clientWidth / 2) / state.zoom - 118;
  const cy = (-state.pan.y + cv.clientHeight / 2) / state.zoom - 50;
  addEntry(entry, Math.round(cx), Math.round(cy));
  toast(`Peça "${entry.name}" adicionada ao Canvas`);
}

let activeModalItem = null;

async function openItemModal(it) {
  const modal = $("modal-item");
  if (!modal) return;
  activeModalItem = it;
  modal.open(it, {
    showCanvasAction: true,
    onAddToCanvas: (entry) => addEntryToCenter(entry),
    onCanonize: () => pollItems(),
  });
}

// paletteClick: com a troca ativa substitui a peça; com nó selecionado e hook/gate, encaixa; senão adiciona ao centro.
function paletteClick(entry) {
  if (state.swapping) return swapNode(state.swapping, entry);
  const sel = state.selected;
  if (sel?.kind === "node" && (entry.cat === "hook" || entry.cat === "gate")) return attachToNode(state.nodes.get(sel.id), entry);
  return addEntryToCenter(entry);
}

function rebuildPalette() {
  state.palette.setModel(buildPaletteModel([...state.blocks.values()], state.items));
}

async function loadItems() {
  const results = await Promise.all(ITEM_KINDS.map((k) => api(`/api/assets/${k}`)));
  ITEM_KINDS.forEach((k, i) => { state.items[k] = results[i].ok ? results[i].data.items : []; });
  state.suggest.clear();
  rebuildPalette();
}

// pollItems recarrega o inventário só quando a versão publicada pelo servidor muda.
async function pollItems() {
  const { ok, data } = await api("/api/assets");
  if (ok && data.version !== state.itemsVersion) {
    state.itemsVersion = data.version;
    await loadItems();
  }
}

function invRow(id, src, on = true) {
  const row = el("div", "inv-row");
  row.append(el("span", on ? "dot" : "dot dot--off"), el("span", "inv-row__id", id), el("span", "inv-row__src", src));
  return row;
}

const HARNESS_CONFIG = [
  { key: "9router", svgKey: "nineRouter", label: "9Router VPS" },
  { key: "claude", svgKey: "claude", label: "Claude CLI" },
  { key: "ollama", svgKey: "ollama", label: "Ollama Local" },
  { key: "gemini", svgKey: "gemini", label: "Gemini / AGY" },
  { key: "other", svgKey: "other", label: "Outros Modelos" },
];

function renderInventory() {
  const box = $("inventory");
  const treeBox = $("models-treeview");
  const totalBadge = $("models-total-badge");

  // O inventário vivo é 100% gerenciado pela árvore treeview-container por harness
  if (box) box.replaceChildren();

  const totalCount = state.inv.models.length + (state.inv.router_online ? 1 : 0);
  if (totalBadge) totalBadge.textContent = String(totalCount);

  if (!treeBox) return;

  const groups = new Map();
  HARNESS_CONFIG.forEach((h) => groups.set(h.key, []));

  for (const m of state.inv.models) {
    const src = String(m.source || "").toLowerCase();
    let matched = false;
    for (const h of HARNESS_CONFIG) {
      if (h.key !== "other" && src.includes(h.key)) {
        groups.get(h.key).push(m);
        matched = true;
        break;
      }
    }
    if (!matched) groups.get("other").push(m);
  }

  const groupEls = [];
  for (const h of HARNESS_CONFIG) {
    const list = groups.get(h.key) || [];
    const isRouter = h.key === "9router";
    const count = list.length + (isRouter && state.inv.router_online ? 1 : 0);
    if (count === 0 && !isRouter) continue;

    const grp = el("div", "treeview-group");
    // 9Router nasce aberto, os outros fechados
    const isClosed = (h.key !== "9router");
    if (isClosed) grp.classList.add("is-closed");

    const head = el("div", "treeview-group-head");
    head.tabIndex = 0;
    head.setAttribute("role", "button");

    const titleWrap = el("div", "treeview-group-head-title");
    const chevron = createSvgIcon("chevron", "svg-icon--chevron");
    const hIcon = createSvgIcon(h.svgKey || "other", "palette-cat-icon");
    const hLabel = el("span", "treeview-group-label", h.label);
    titleWrap.append(chevron, hIcon, hLabel);
    head.append(titleWrap, el("span", "treeview-item-badge", String(count)));

    const itemsBox = el("div", "treeview-items");

    if (isRouter) {
      const rRow = el("div", "treeview-model-item");
      rRow.style.cursor = "pointer";
      const rMain = el("div", "treeview-model-item-main");
      rMain.append(el("span", state.inv.router_online ? "dot" : "dot dot--off"), el("span", "inv-row__id", "9Router Gateway"));
      rRow.append(rMain, el("span", "treeview-harness-badge", state.inv.router_online ? "online" : "offline"));
      rRow.addEventListener("click", () => {
        toast(`Gateway 9Router está ${state.inv.router_online ? "online na VPS" : "offline"}.`);
      });
      itemsBox.append(rRow);
    }

    for (const m of list) {
      const mRow = el("div", "treeview-model-item");
      mRow.style.cursor = "pointer";
      mRow.title = `Clique para atribuir ao nó selecionado ou arrastar`;
      const mMain = el("div", "treeview-model-item-main");
      mMain.append(el("span", "dot"), el("span", "inv-row__id", m.id));
      mRow.append(mMain, el("span", "treeview-harness-badge", m.source));

      // Clique no modelo: atribui ao nó selecionado ou informa
      mRow.addEventListener("click", () => {
        const sel = state.selected;
        if (sel?.kind === "node") {
          const node = state.nodes.get(sel.id);
          if (node) {
            node.model = m.id;
            renderNodes();
            saveState();
            toast(`Modelo "${m.id}" atribuído ao nó selecionado.`);
            return;
          }
        }
        toast(`Modelo "${m.id}" (${m.source}) selecionado. Selecione um nó no Canvas para atribuí-lo.`);
      });

      itemsBox.append(mRow);
    }

    head.addEventListener("click", () => {
      // Sanfona estrita de modelos: fecha outros ao abrir este
      const wasClosed = grp.classList.contains("is-closed");
      treeBox.querySelectorAll(".treeview-group").forEach((g) => g.classList.add("is-closed"));
      if (wasClosed) grp.classList.remove("is-closed");
    });

    grp.append(head, itemsBox);
    groupEls.push(grp);
  }

  if (groupEls.length === 0) {
    treeBox.replaceChildren(el("p", "note", "Nenhum modelo conectado."));
  } else {
    treeBox.replaceChildren(...groupEls);
  }
}

function fillRouteSelects() {
  const from = $("route-from");
  const to = $("route-to");
  state.types.forEach((t) => {
    from.append(new Option(shortType(t), t));
    to.append(new Option(shortType(t), t));
  });
  from.value = "Artifact:URL";
  to.value = "Artifact:Spec";
}

function filterModels(query) {
  const treeBox = $("models-treeview");
  if (!treeBox) return;
  const q = String(query || "").trim().toLowerCase();
  const groups = treeBox.querySelectorAll(".treeview-group");
  groups.forEach((grp) => {
    const items = grp.querySelectorAll(".treeview-model-item");
    let grpMatches = 0;
    items.forEach((item) => {
      const match = !q || item.textContent.toLowerCase().includes(q);
      item.classList.toggle("is-hidden", !match);
      if (match) grpMatches++;
    });
    if (q) {
      grp.classList.toggle("is-hidden", grpMatches === 0);
      if (grpMatches > 0) {
        grp.classList.remove("is-closed");
        const arrow = grp.querySelector(".treeview-group-arrow");
        if (arrow) arrow.textContent = "▾";
      }
    } else {
      grp.classList.remove("is-hidden");
    }
  });
}

function setupPalette() {
  state.palette = createPalette({
    search: $("palette-search"), treeview: $("palette-treeview"), totalBadge: $("palette-total-badge"),
    cats: $("palette-cats"), list: $("palette-list"), spacer: $("palette-spacer"),
    rows: $("palette-rows"), empty: $("palette-empty"), emptyTitle: $("palette-empty-title"),
    emptyText: $("palette-empty-text"), banner: $("palette-banner"), hint: $("palette-hint"),
  }, { onClick: paletteClick, onInspect: openItemModal });
  $("palette-search")?.addEventListener("input", (e) => filterModels(e.target.value));
}

async function boot() {
  const [blocks, inv] = await Promise.all([api("/api/blocks"), api("/api/inventory")]);
  blocks.data.blocks.forEach((b) => state.blocks.set(b.id, b));
  state.matrix = blocks.data.matrix;
  state.types = blocks.data.types;
  if (inv.ok) state.inv = inv.data;
  setupPalette();
  renderInventory();
  fillRouteSelects();
  setupDropZone();
  $("nodes").addEventListener("pointerdown", onNodesPointerDown);
  $("edges").addEventListener("pointerdown", onEdgesPointerDown);
  $("edges").addEventListener("click", onEdgesClick);
  document.addEventListener("keydown", onKey);
  // Canvas N8N Pan & Zoom
  const cv = $("canvas");
  cv.addEventListener("wheel", (ev) => {
    if (ev.ctrlKey || ev.metaKey) {
      ev.preventDefault();
      const r = cv.getBoundingClientRect();
      zoomBy(ev.deltaY < 0 ? 1.15 : 1 / 1.15, { x: ev.clientX - r.left, y: ev.clientY - r.top });
    } else {
      ev.preventDefault();
      if (ev.shiftKey && ev.deltaX === 0) {
        state.pan.x -= ev.deltaY;
      } else {
        state.pan.x -= ev.deltaX;
        state.pan.y -= ev.deltaY;
      }
      updateCanvasTransform();
    }
  }, { passive: false });

  // Pressionar e segurar barra de Espaço ativa o modo Pan N8N
  document.addEventListener("keydown", (e) => {
    if (e.code === "Space" && !e.target.matches("input, textarea, select") && !state.spacePan) {
      state.spacePan = true;
      cv.classList.add("is-panning");
    }
    if ((e.ctrlKey || e.metaKey) && (e.key === "0" || e.code === "Digit0" || e.code === "Numpad0")) {
      e.preventDefault();
      fitToScreen();
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "z") {
      e.preventDefault();
      if (e.shiftKey) redo(); else undo();
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "y") {
      e.preventDefault();
      redo();
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "b") {
      e.preventDefault();
      toggleSidebar();
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "s") {
      e.preventDefault();
      $("btn-materialize")?.click();
    }
  });

  document.addEventListener("keyup", (e) => {
    if (e.code === "Space") {
      state.spacePan = false;
      cv.classList.remove("is-panning");
    }
  });

  // Determina se um clique no canvas deve iniciar o pan 2D
  function canPanCanvas(ev) {
    if (ev.target.closest(".node, .canvas-controls, .issues, .edge-del-group, .edge-btn-remove, .slot-pill")) {
      return false;
    }
    if (state.spacePan || ev.button === 1) return true;
    return ev.button === 0;
  }

  cv.addEventListener("pointerdown", (ev) => {
    if (!canPanCanvas(ev)) return;
    if (ev.button === 1 || state.spacePan) ev.preventDefault();

    const startX = ev.clientX;
    const startY = ev.clientY;
    const initX = state.pan.x;
    const initY = state.pan.y;
    let moved = false;

    cv.classList.add("is-dragging");
    try { cv.setPointerCapture(ev.pointerId); } catch (_) {}
    const move = (e) => {
      const dx = e.clientX - startX;
      const dy = e.clientY - startY;
      if (Math.abs(dx) > 3 || Math.abs(dy) > 3) moved = true;
      state.pan.x = initX + dx;
      state.pan.y = initY + dy;
      updateCanvasTransform();
    };
    const up = (e) => {
      try { cv.releasePointerCapture(ev.pointerId); } catch (_) {}
      cv.classList.remove("is-dragging");
      window.removeEventListener("pointermove", move);
      window.removeEventListener("pointerup", up);
      if (!moved && ev.button === 0 && !state.spacePan) {
        if (!ev.target.closest(".node, .edge, .canvas-controls")) {
          select(null);
        }
      }
    };
    window.addEventListener("pointermove", move);
    window.addEventListener("pointerup", up);
  });

  const getZoomPivot = () => {
    if (state.selected?.kind === "node") {
      const n = state.nodes.get(state.selected.id);
      if (n) {
        return {
          x: state.pan.x + (n.x + NODE_W / 2) * state.zoom,
          y: state.pan.y + (n.y + 50) * state.zoom,
        };
      }
    }
    return null;
  };
  $("btn-zoom-in").addEventListener("click", () => zoomBy(1.2, getZoomPivot()));
  $("btn-zoom-out").addEventListener("click", () => zoomBy(1 / 1.2, getZoomPivot()));
  $("btn-zoom-fit").addEventListener("click", () => fitToScreen());
  $("btn-undo")?.addEventListener("click", undo);
  $("btn-redo")?.addEventListener("click", redo);
  $("btn-auto-layout")?.addEventListener("click", autoLayout);
  $("btn-toggle-sidebar").addEventListener("click", toggleSidebar);
  document.querySelector(".sidebar-collapsed-rail")?.addEventListener("click", toggleSidebar);
  $("btn-toggle-inspector")?.addEventListener("click", toggleInspector);
  $("btn-toggle-palette").addEventListener("click", () => toggleAccordion("palette"));
  $("btn-toggle-models").addEventListener("click", () => toggleAccordion("models"));

  $("pipeline-name").addEventListener("blur", (e) => {
    const raw = e.target.value.trim();
    if (raw) {
      e.target.value = toKebabCase(raw);
    } else {
      e.target.value = "unnamed-loop";
    }
  });

  $("pipeline-name").addEventListener("change", (e) => {
    const raw = e.target.value.trim();
    if (raw) {
      e.target.value = toKebabCase(raw);
    }
  });



  $("btn-validate").addEventListener("click", (e) => withBusy(e.currentTarget, validate));
  $("btn-materialize").addEventListener("click", (e) => withBusy(e.currentTarget, materialize));
  $("btn-route").addEventListener("click", (e) => withBusy(e.currentTarget, traceRoute));
  $("btn-close-issues").addEventListener("click", () => { $("issues").hidden = true; });
  initVaultModal();
  initRunModal();
  initTemplatesModal();
  initImportModal();
  updateEmpty();
  await pollItems();
  setInterval(() => { if (document.visibilityState === "visible") pollItems().catch(() => {}); }, POLL_MS);
}

function initVaultModal() {
  const btnVault = $("btn-vault");
  const modal = $("modal-vault");
  if (!btnVault || !modal) return;

  const btnClose = $("btn-modal-vault-close");
  const btnScan = $("btn-vault-scan");
  const btnApply = $("btn-vault-apply");
  const pathDisplay = $("vault-path-display");
  const drivesDisplay = $("vault-drives-display");
  const panel = $("vault-results-panel");
  const summary = $("vault-results-summary");
  const groupsList = $("vault-groups-list");

  btnVault.addEventListener("click", async () => {
    try {
      const res = await fetch("/api/canonize/status");
      if (!res.ok) {
        toast("Status do cofre indisponível no momento (servidor respondeu HTTP " + res.status + ").", true);
        return;
      }
      const data = await res.json();
      pathDisplay.textContent = data.vault || "~/.aidd/canonical";
      drivesDisplay.textContent = (data.drives || []).join(", ") || "C:";
      btnApply.disabled = true;
      panel.hidden = true;
      if (groupsList) groupsList.replaceChildren();
      modal.showModal();
    } catch (e) {
      toast("Erro ao obter status do cofre: " + e.message, true);
    }
  });

  btnClose?.addEventListener("click", () => modal.close());

  btnScan?.addEventListener("click", async () => {
    try {
      btnScan.disabled = true;
      btnScan.textContent = "Auditando...";
      const res = await fetch("/api/canonize/plan", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{}"
      });
      if (!res.ok) throw new Error("Falha no dry-run");
      const plan = await res.json();
      summary.textContent = `Arquivos escaneados: ${plan.files_scanned} | Grupos duplicados: ${plan.groups.length} | Economia estimada: ${Math.round(plan.saved_bytes / 1024)} KB`;
      
      if (groupsList) {
        groupsList.replaceChildren();
        const maxDisplay = 20;
        const groupsToShow = (plan.groups || []).slice(0, maxDisplay);
        for (const g of groupsToShow) {
          const item = el("div", "vault-group-item");
          const canonLine = el("p", "vault-group-canonical", `Canônico: ${g.canonical}`);
          const dupCount = el("span", "vault-group-badge", ` (${g.links?.length || 0} links)`);
          canonLine.appendChild(dupCount);
          item.appendChild(canonLine);
          groupsList.appendChild(item);
        }
        if ((plan.groups || []).length > maxDisplay) {
          const more = el("p", "vault-group-more", `... e mais ${plan.groups.length - maxDisplay} grupo(s)`);
          groupsList.appendChild(more);
        }
      }

      panel.hidden = false;
      btnApply.disabled = plan.groups.length === 0;
    } catch (e) {
      toast("Erro ao auditar cofre: " + e.message, true);
    } finally {
      btnScan.disabled = false;
      btnScan.textContent = "Auditar Desduplicação (Dry-Run)";
    }
  });

  btnApply?.addEventListener("click", async () => {
    const ok = await confirmDialog({
      title: "Desduplicar Sistema Operacional?",
      message: "Arquivos redundantes serão substituídos por hardlinks canônicos no cofre.",
      confirmLabel: "Aplicar Desduplicação",
      danger: true,
    });
    if (!ok) return;
    try {
      btnApply.disabled = true;
      btnApply.textContent = "Aplicando links...";
      const res = await fetch("/api/canonize/apply", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ apply: true })
      });
      if (!res.ok) throw new Error("Falha ao aplicar canonização");
      const result = await res.json();
      summary.textContent = `Desduplicação concluída! Links criados: ${result.linked} | Economia real: ${Math.round(result.saved_bytes / 1024)} KB`;
      toast(`Cofre atualizado: ${result.linked} links criados.`);
      btnApply.disabled = true;
    } catch (e) {
      toast("Erro ao aplicar canonização: " + e.message, true);
      btnApply.disabled = false;
    } finally {
      btnApply.textContent = "Aplicar Desduplicação Global";
    }
  });
}

function initRunModal() {
  const btnRun = $("btn-run");
  const modal = $("modal-run");
  if (!btnRun || !modal) return;

  const btnClose = $("btn-modal-run-close");
  const badge = $("run-status-badge");
  const nameEl = $("run-pipeline-name");
  const tokensTotalEl = $("run-tokens-total");
  const tokensInEl = $("run-tokens-in");
  const tokensOutEl = $("run-tokens-out");
  const timelineEl = $("run-nodes-timeline");
  const consoleEl = $("run-logs-console");

  let runAbortController = null;

  const cancelRun = () => {
    if (runAbortController) {
      runAbortController.abort();
      runAbortController = null;
    }
  };

  btnClose?.addEventListener("click", () => {
    cancelRun();
    modal.close();
  });
  modal.addEventListener("cancel", () => {
    cancelRun();
  });

  btnRun.addEventListener("click", () => withBusy(btnRun, async () => {
    if (state.nodes.size === 0) {
      toast("O canvas está vazio. Adicione nós antes de executar.", true);
      return;
    }

    const valid = await validate();
    if (!valid) return;

    const hasScripts = Array.from(state.nodes.values()).some((n) => n.kind === "script");
    let allowScripts = false;
    if (hasScripts) {
      const agreed = await confirmDialog({
        title: "Autorizar Execução de Scripts?",
        message: "Esta esteira possui nós do tipo script que executarão comandos no sistema operacional. Você autoriza a execução?",
        confirmLabel: "Autorizar e Executar",
        danger: true,
      });
      if (!agreed) {
        toast("Execução cancelada: scripts não autorizados.");
        return;
      }
      allowScripts = true;
    }

    const dag = compileDAG();
    nameEl.textContent = dag.name || "nova-esteira";
    badge.textContent = "Iniciando...";
    badge.className = "modal-badge badge--ready";
    tokensTotalEl.textContent = "0";
    tokensInEl.textContent = "0";
    tokensOutEl.textContent = "0";
    timelineEl.textContent = "";
    consoleEl.textContent = "";

    modal.showModal();

    const appendLog = (msg) => {
      consoleEl.textContent += (consoleEl.textContent ? "\n" : "") + msg;
      consoleEl.scrollTop = consoleEl.scrollHeight;
    };

    appendLog(`[${new Date().toLocaleTimeString()}] Conectando ao motor de execução LegoLoop...`);

    cancelRun();
    runAbortController = new AbortController();

    try {
      const resp = await fetch("/api/run", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ dag, allow_scripts: allowScripts }),
        signal: runAbortController.signal,
      });

      if (!resp.ok) {
        const errJson = await resp.json().catch(() => ({ error: "Erro desconhecido" }));
        badge.textContent = "Falha";
        badge.className = "modal-badge badge--danger";
        appendLog(`[ERRO HTTP ${resp.status}] ${errJson.error || "Falha na requisição"}`);
        toast(errJson.error || "Falha na execução", true);
        return;
      }

      const reader = resp.body.getReader();
      const decoder = new TextDecoder("utf-8");
      let buffer = "";
      const nodeRows = new Map();

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });

        const parts = buffer.split("\n\n");
        buffer = parts.pop() || "";

        for (const block of parts) {
          const lines = block.split("\n");
          let evData = null;
          for (const line of lines) {
            if (line.startsWith("data: ")) {
              try {
                evData = JSON.parse(line.slice(6));
              } catch (_) {}
            }
          }
          if (!evData) continue;

          switch (evData.type) {
            case "dag_start": {
              badge.textContent = "Executando...";
              badge.className = "modal-badge badge--running";
              appendLog(`[${new Date().toLocaleTimeString()}] Esteira "${evData.dag}" iniciada (${evData.nodes_total} nós).`);

              const telemStatus = $("telemetry-global-status");
              if (telemStatus) {
                telemStatus.textContent = "Executando...";
                telemStatus.className = "telemetry-value telemetry-status--running";
              }
              if ($("telemetry-global-duration")) $("telemetry-global-duration").textContent = "0.0s";
              if ($("telemetry-global-tokens-in")) $("telemetry-global-tokens-in").textContent = "0";
              if ($("telemetry-global-tokens-out")) $("telemetry-global-tokens-out").textContent = "0";
              if ($("telemetry-global-cost")) $("telemetry-global-cost").textContent = "$0.000";

              state.telemetry.clear();
              renderAllNodes();
              break;
            }

            case "node_start": {
              appendLog(`[${new Date().toLocaleTimeString()}] ▶ Nó "${evData.id}" em execução...`);

              state.telemetry.set(evData.id, { status: "running" });
              const targetNode = state.nodes.get(evData.id);
              if (targetNode) renderNode(targetNode);

              let row = nodeRows.get(evData.id);
              if (!row) {
                row = document.createElement("div");
                row.className = "run-node-row is-running";
                const titleSpan = document.createElement("strong");
                titleSpan.textContent = evData.id;
                const statusSpan = document.createElement("span");
                statusSpan.className = "node-run-status";
                statusSpan.textContent = "[RUN] Executando";
                row.appendChild(titleSpan);
                row.appendChild(statusSpan);
                timelineEl.appendChild(row);
                nodeRows.set(evData.id, row);
              }
              break;
            }

            case "node_finish": {
              state.telemetry.set(evData.id, {
                status: evData.status,
                duration_ms: evData.duration_ms,
                tokens: evData.tokens,
                model: evData.model,
              });
              const targetNode = state.nodes.get(evData.id);
              if (targetNode) renderNode(targetNode);

              let row = nodeRows.get(evData.id);
              if (!row) {
                row = document.createElement("div");
                row.className = "run-node-row";
                const titleSpan = document.createElement("strong");
                titleSpan.textContent = evData.id;
                row.appendChild(titleSpan);
                timelineEl.appendChild(row);
                nodeRows.set(evData.id, row);
              }

              const isPass = evData.status === "passed";
              row.className = "run-node-row " + (isPass ? "is-passed" : "is-failed");
              const statusSpan = row.querySelector(".node-run-status") || document.createElement("span");
              statusSpan.className = "node-run-status";
              const modelInfo = evData.model ? ` [${evData.model}]` : "";
              statusSpan.textContent = `${isPass ? "[OK] Passou" : "[FALHA] Falhou"} (${evData.duration_ms}ms)${modelInfo}`;
              if (!row.contains(statusSpan)) row.appendChild(statusSpan);

              appendLog(`[${new Date().toLocaleTimeString()}] [No] "${evData.id}": ${evData.status} em ${evData.duration_ms}ms${modelInfo}`);
              if (evData.error) {
                appendLog(`   └─ Erro: ${evData.error}`);
              }
              break;
            }

            case "dag_finish": {
              const ok = evData.ok;
              badge.textContent = ok ? "Concluido" : "Interrompido";
              badge.className = "modal-badge " + (ok ? "badge--passed" : "badge--danger");
              if (evData.tokens) {
                tokensTotalEl.textContent = String(evData.tokens.input + evData.tokens.output);
                tokensInEl.textContent = String(evData.tokens.input);
                tokensOutEl.textContent = String(evData.tokens.output);
              }

              const telemStatus = $("telemetry-global-status");
              if (telemStatus) {
                telemStatus.textContent = ok ? "Sucesso" : "Falha";
                telemStatus.className = "telemetry-value " + (ok ? "telemetry-status--passed" : "telemetry-status--failed");
              }
              if (evData.duration_ms !== undefined && $("telemetry-global-duration")) {
                $("telemetry-global-duration").textContent = `${(evData.duration_ms / 1000).toFixed(1)}s`;
              }
              if (evData.tokens) {
                if ($("telemetry-global-tokens-in")) $("telemetry-global-tokens-in").textContent = String(evData.tokens.input);
                if ($("telemetry-global-tokens-out")) $("telemetry-global-tokens-out").textContent = String(evData.tokens.output);
              }
              if ($("telemetry-global-cost")) {
                $("telemetry-global-cost").textContent = `$${(evData.cost_est || 0).toFixed(4)}`;
              }

              appendLog(`[${new Date().toLocaleTimeString()}] Execucao finalizada: ${ok ? "SUCESSO" : "COM FALHAS"}`);
              toast(ok ? "Esteira executada com sucesso!" : "Esteira finalizou com falhas.", !ok);
              break;
            }

            case "dag_error": {
              badge.textContent = "Erro";
              badge.className = "modal-badge badge--danger";
              const telemStatus = $("telemetry-global-status");
              if (telemStatus) {
                telemStatus.textContent = "Erro";
                telemStatus.className = "telemetry-value telemetry-status--failed";
              }
              appendLog(`[${new Date().toLocaleTimeString()}] Erro no DAG: ${evData.error}`);
              toast(evData.error, true);
              break;
            }
          }
        }
      }
    } catch (e) {
      if (e.name === "AbortError") {
        badge.textContent = "Cancelado";
        badge.className = "modal-badge badge--ready";
        appendLog(`[INFO] Execução cancelada pelo usuário.`);
        return;
      }
      badge.textContent = "Falha de Conexao";
      badge.className = "modal-badge badge--danger";
      appendLog(`[FALHA] ${e.message}`);
      toast("Falha na execucao: " + e.message, true);
    } finally {
      runAbortController = null;
    }
  }));
}

function initTemplatesModal() {
  const btnTemplates = $("btn-templates");
  const modal = $("modal-templates");
  if (!btnTemplates || !modal) return;

  const btnClose = $("btn-modal-templates-close");
  const grid = $("templates-grid");

  btnClose?.addEventListener("click", () => modal.close());

  btnTemplates.addEventListener("click", async () => {
    try {
      const resp = await fetch("/api/templates");
      if (!resp.ok) return toast("Nao foi possivel carregar os templates", true);
      const data = await resp.json();

      grid.textContent = "";

      for (const tmpl of (data.templates || [])) {
        const card = document.createElement("button");
        card.className = "template-card";
        card.type = "button";

        const head = document.createElement("div");
        head.className = "template-card-header";

        const icon = document.createElement("span");
        icon.className = "template-icon";
        if (typeof createSvgIconNode === "function") {
          icon.replaceChildren(createSvgIconNode(tmpl.icon || "package"));
        } else {
          icon.textContent = "";
        }

        const name = document.createElement("span");
        name.className = "template-name";
        name.textContent = tmpl.name;

        head.appendChild(icon);
        head.appendChild(name);

        const desc = document.createElement("p");
        desc.className = "template-desc";
        desc.textContent = tmpl.description;

        const meta = document.createElement("span");
        meta.className = "template-meta";
        meta.textContent = `${tmpl.category} • ${tmpl.composition.nodes?.length || 0} nós`;

        card.appendChild(head);
        card.appendChild(desc);
        card.appendChild(meta);

        card.addEventListener("click", async () => {
          if (state.nodes.size > 0) {
            const confirmed = await confirmDialog({
              title: "Substituir esteira pelo Template?",
              message: `Os nós atuais serão removidos e substituídos pelo template "${tmpl.name}".`,
              confirmLabel: "Carregar Template",
              danger: true,
            });
            if (!confirmed) return;
          }

          [...state.nodes.keys()].forEach(removeNode);

          if ($("pipeline-name")) {
            $("pipeline-name").value = tmpl.composition.name || tmpl.id;
          }

          const nodeMap = new Map();
          (tmpl.composition.nodes || []).forEach((n, i) => {
            const added = addNode(n.block || n.id, 60 + i * (NODE_W + 80), 100 + (i % 2) * 40);
            if (added) {
              nodeMap.set(n.id, added.id);
            }
          });

          (tmpl.composition.edges || []).forEach((e) => {
            const fromId = nodeMap.get(e.from);
            const toId = nodeMap.get(e.to);
            if (fromId && toId) {
              addEdge(fromId, toId);
            }
          });

          modal.close();
          autoLayout();
          fitToScreen();
          toast(`Template "${tmpl.name}" carregado com sucesso!`);
        });

        grid.appendChild(card);
      }

      modal.showModal();
    } catch (err) {
      toast("Erro ao carregar templates: " + err.message, true);
    }
  });
}

function initImportModal() {
  const btnOpen = $("btn-import-github");
  const modal = $("modal-import");
  if (!btnOpen || !modal) return;

  const btnClose = $("btn-modal-import-close");
  const btnCancel = $("btn-import-cancel");
  const form = $("form-import");
  const inputUrl = $("import-url");
  const inputName = $("import-name");
  const btnSubmit = $("btn-import-submit");

  const close = () => modal.close();
  btnOpen.addEventListener("click", () => {
    form?.reset();
    modal.showModal();
  });
  btnClose?.addEventListener("click", close);
  btnCancel?.addEventListener("click", close);

  form?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const url = inputUrl?.value.trim();
    const name = inputName?.value.trim();
    if (!url) return;

    if (btnSubmit) {
      btnSubmit.disabled = true;
      btnSubmit.textContent = "Clonando & Canonizando...";
    }

    try {
      const resp = await fetch("/api/assets/import", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url, name }),
      });
      const data = await resp.json();
      if (!resp.ok) {
        throw new Error(data.error || "Falha na importação");
      }

      toast(`Skill "${data.name}" importada e canonizada! SHA: ${String(data.sha256).slice(0, 8)}`);
      modal.close();
      await pollItems();
    } catch (err) {
      toast("Erro ao importar: " + err.message, true);
    } finally {
      if (btnSubmit) {
        btnSubmit.disabled = false;
        btnSubmit.textContent = "Importar & Canonizar";
      }
    }
  });
}


boot().catch((err) => toast(`Não foi possível iniciar o Studio: ${err.message}`, true));
