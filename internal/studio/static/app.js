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
const PIECE_KIND = { skill: "skill", script: "script", mcp: "mcp", loop: "subloop", hook: "hook", gate: "gate" };
const PIECE_LABEL = { skill: "skill", script: "script", mcp: "MCP", subloop: "sub-loop", hook: "hook", gate: "gate" };

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
  lastDAG: null,
  swapping: null,
  gateDraft: null,
  suggest: new Map(),
  suggestPending: new Set(),
  palette: null,
  process: null,
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
    id: `piece:${p.kind}`, kind: p.kind, label: node.label, ref: p.ref || "", consumes: [ANY_TYPE], produces: ANY_TYPE,
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
  return { id, x: Math.round(Number(x) || 0), y: Math.round(Number(y) || 0), slots: emptySlots(), execution: "inline" };
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
    case "trigger_processo_fim": return { type: "processo_fim", processo: "upstream" };
    case "trigger_loop_end": return { type: "loop_end", loop: "upstream" };
    case "trigger_board": return { type: "board", board: "backlog" };
    default: return { type: "manual" };
  }
}

function normalizeTriggerConfig(kind, cfg) {
  const base = defaultTriggerConfig(kind);
  const merged = { ...base, ...(cfg || {}) };
  if (kind === "trigger_processo_fim" || kind === "trigger_loop_end" || merged.type === "processo_fim" || merged.type === "loop_end") {
    const alvo = String(merged.processo || merged.loop || "").trim() || "upstream";
    if (kind === "trigger_loop_end" || merged.type === "loop_end") {
      return { type: "loop_end", loop: alvo };
    }
    return { type: "processo_fim", processo: alvo };
  }
  return merged;
}

function addNode(blockId, x, y) {
  const block = state.blocks.get(blockId);
  if (!block) return null;
  const isTrig = block.kind && block.kind.startsWith("trigger_");
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
  const ref = entry.cat === "loop" ? (entry.meta?.loop_id || entry.itemId || entry.name) : entry.itemId;
  const node = {
    ...baseNode(newNodeId(entry.name), x, y), label: entry.name,
    nature: entry.meta?.nature || entry.nature || "",
    piece: { kind, ref, artifact: entry.meta?.artifact || "", gate: "" }, routing: defaultRouting({}),
  };
  node.piece.artifact = node.piece.artifact || uniqueArtifact(entry.name, node.id);
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
  if (k.startsWith("flow_") || k.startsWith("security_") || k.startsWith("vcs_") || k.startsWith("deploy_") || k === "control") return "node--control";
  if (k === "mcp") return "node--mcp";
  if (k === "subloop" || k === "loop" || k === "subprocesso") return "node--loop";
  if (k === "script") return "node--script";
  if (k === "hook") return "node--hook";
  if (k === "gate") return "node--gate";
  return "node--skill";
}

const NODE_FAMILY_LABEL = {
  trigger: "TRIGGER",
  control: "CONTROLE",
  skill: "SKILL",
  mcp: "MCP",
  loop: "LOOP",
  script: "SCRIPT",
  hook: "HOOK",
  gate: "GATE",
};

function buildInspectableItemFromNode(node) {
  const block = blockOf(node);
  if (node.piece && node.piece.ref) {
    const kind = node.piece.kind === "subloop" ? "loop" : node.piece.kind;
    const pool = state.items[kind] || [];
    const found = pool.find((x) => x.id === node.piece.ref || x.name === node.piece.ref || x.name === node.label);
    if (found) return { ...found, cat: kind };
    return {
      id: node.piece.ref,
      name: node.label || block.label || node.id,
      kind,
      cat: kind,
      description: block.description || "",
      artifact: node.piece.artifact || block.artifact || "",
      gate: node.piece.gate || block.gate || "",
      meta: { artifact: node.piece.artifact || block.artifact || "", role: node.role || "" },
    };
  }
  const blkId = node.block || node.id;
  const blkKind = block.kind || "control";
  const allPool = Object.values(state.items || {}).flat();
  const invMatch = allPool.find((x) => x.id === blkId || x.name === block.label || x.id === block.ref);
  if (invMatch) {
    return {
      ...invMatch,
      cat: nodeCategoryClass(block).replace("node--", ""),
      consumes: block.consumes || invMatch.consumes || [],
      produces: block.produces || invMatch.produces || "",
      artifact: block.artifact || invMatch.artifact || "",
      gate: block.gate || invMatch.gate || "",
      suggested_model: block.suggested_model || "",
    };
  }
  return {
    id: blkId,
    name: block.label || node.id,
    kind: blkKind,
    cat: nodeCategoryClass(block).replace("node--", ""),
    description: block.description || "",
    consumes: block.consumes || [],
    produces: block.produces || "",
    artifact: block.artifact || "",
    gate: block.gate || "",
    suggested_model: block.suggested_model || "",
    meta: { role: node.role || "", kind: blkKind },
  };
}

function computeExecutionLevels() {
  const nodeIds = Array.from(state.nodes.keys());
  const levelMap = new Map();
  if (!nodeIds.length) {
    state.execLevels = levelMap;
    return [];
  }
  const indeg = new Map(nodeIds.map((id) => [id, 0]));
  const adj = new Map(nodeIds.map((id) => [id, []]));
  const outdeg = new Map(nodeIds.map((id) => [id, 0]));
  for (const e of state.edges) {
    if (!state.nodes.has(e.from) || !state.nodes.has(e.to)) continue;
    adj.get(e.from).push(e.to);
    indeg.set(e.to, (indeg.get(e.to) || 0) + 1);
    outdeg.set(e.from, (outdeg.get(e.from) || 0) + 1);
  }
  const levels = [];
  let queue = nodeIds.filter((id) => indeg.get(id) === 0);
  const visited = new Set();
  while (queue.length > 0) {
    const mode = queue.length > 1 ? "async" : "sync";
    const step = levels.length + 1;
    levels.push({ step, mode, ids: [...queue] });
    const next = [];
    for (const id of queue) {
      visited.add(id);
      levelMap.set(id, { step, mode, levelSize: queue.length, outdeg: outdeg.get(id) || 0 });
      for (const nxt of adj.get(id) || []) {
        const d = (indeg.get(nxt) || 1) - 1;
        indeg.set(nxt, d);
        if (d === 0 && !visited.has(nxt)) next.push(nxt);
      }
    }
    queue = next;
  }
  // Nós em ciclos ainda não visitados recebem passo final explícito
  const remaining = nodeIds.filter((id) => !visited.has(id));
  if (remaining.length > 0) {
    const step = levels.length + 1;
    const mode = remaining.length > 1 ? "async" : "sync";
    levels.push({ step, mode, ids: remaining, cycle: true });
    for (const id of remaining) {
      levelMap.set(id, { step, mode, levelSize: remaining.length, outdeg: outdeg.get(id) || 0, cycle: true });
    }
  }
  state.execLevels = levelMap;
  return levels;
}

function renderExecutionStrip(levels) {
  const strip = $("execution-flow-strip");
  const box = $("execution-flow-steps");
  if (!strip || !box) return;
  if (!levels || levels.length === 0) {
    strip.hidden = true;
    box.replaceChildren();
    return;
  }
  strip.hidden = false;
  const chips = [];
  levels.forEach((lvl, idx) => {
    if (idx > 0) chips.push(el("span", "exec-step-arrow", "->"));
    const labels = lvl.ids.map((id) => {
      const n = state.nodes.get(id);
      return n ? (blockOf(n).label || id) : id;
    });
    const modeTag = lvl.mode === "async" ? `Paralelo (${lvl.ids.length})` : "Síncrono";
    const chipText = `P${lvl.step} [${modeTag}]: ${labels.join(" + ")}${lvl.cycle ? " (ciclo)" : ""}`;
    chips.push(el("span", `exec-step-chip${lvl.mode === "async" ? " exec-step-chip--async" : ""}`, chipText));
  });
  box.replaceChildren(...chips);
}

function nodeEl(node) {
  const block = blockOf(node);
  const isTrig = block.kind && block.kind.startsWith("trigger_");
  const isSubloop = (node.piece && ["subloop", "loop", "subprocesso"].includes(node.piece.kind)) || (block.kind && ["loop", "flow_loop", "subprocesso"].includes(block.kind));
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
      enterSubloop(node.piece?.ref || node.block || node.id, block.label);
    });
  }

  const head = el("div", "node__head");
  const titleEl = el("span", "node__title", isTrig ? `[Gatilho] ${block.label}` : block.label);
  const idEl = el("span", "node__id", node.id);
  const eyeBtn = el("button", "node__eye");
  eyeBtn.type = "button";
  eyeBtn.title = "Inspecionar ficha técnica completa da peça";
  eyeBtn.setAttribute("aria-label", `Inspecionar ${block.label}`);
  if (typeof createSvgIconNode === "function") {
    eyeBtn.append(createSvgIconNode("eye"));
  } else {
    eyeBtn.textContent = "Ver";
  }
  eyeBtn.addEventListener("click", (ev) => {
    ev.stopPropagation();
    openItemModal(buildInspectableItemFromNode(node), node);
  });
  head.append(titleEl, idEl, eyeBtn);
  if (isSubloop) {
    head.append(el("span", "node__macro-hint", "Duplo clique para abrir"));
  }

  const ports = el("div", "node__ports");
  const inLabel = isTrig ? "Universal (*)" : (block.consumes || []).map(shortType).join(" | ");
  ports.append(
    el("span", "node__port-label", inLabel),
    el("span", "node__port-label node__port-label--out", isTrig ? "Universal (*)" : shortType(block.produces)),
  );

  const portIn = el("span", "port port--in");
  portIn.dataset.port = "in";
  portIn.title = isTrig ? "Entrada universal (Gatilho) — conecte de ou para qualquer peça" : `Entrada: ${(block.consumes || []).map(shortType).join(", ")}`;
  const portOut = el("span", "port port--out");
  portOut.dataset.port = "out";
  portOut.title = isTrig ? "Saída universal (Gatilho) — arraste para conectar" : `Saída: ${shortType(block.produces)} — arraste para conectar`;

  const children = [head, ports, el("div", "node__badges"), portIn, portOut];
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
  const isTrig = block.kind && block.kind.startsWith("trigger_");
  root.querySelector(".node__title").textContent = isTrig ? `[Gatilho] ${block.label}` : block.label;
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
  if (node.nature) {
    const nat = String(node.nature).toLowerCase();
    const natCls = nat === "macro" ? "pill--nature-macro" : nat === "nano" ? "pill--nature-nano" : "pill--nature-micro";
    pills.push(pill(natCls, nat.toUpperCase(), `Nível hierárquico: ${nat.toUpperCase()}`));
  }
  if (node.role) pills.push(pill("pill--role", node.role, `Papel Canônico: ${node.role}`));
  const fam = nodeCategoryClass(block).replace("node--", "");
  pills.push(pill("pill--kind", NODE_FAMILY_LABEL[fam] || fam.toUpperCase(), `Categoria: ${fam}`));
  pills.push(pill("pill--handoff", block.artifact, "Handoff imutável"), pill("pill--route", routingLabel(node.routing), "Roteamento"), gate);
  const slots = slotCount(node.slots);
  if (slots > 0) pills.push(pill("pill--slots", `${slots} slot${slots > 1 ? "s" : ""}`, "Gates, scripts e hooks encaixados"));
  if (node.execution === "isolated") pills.push(pill("pill--exec", "isolada", "Executa em workspace separado"));
  const lvlInfo = state.execLevels?.get(node.id);
  if (lvlInfo) {
    const isAsync = lvlInfo.mode === "async";
    const cls = isAsync ? "pill--exec-async" : "pill--exec-sync";
    const lbl = isAsync ? `Passo ${lvlInfo.step} · Assíncrono` : `Passo ${lvlInfo.step} · Síncrono`;
    const tip = isAsync
      ? `Passo ${lvlInfo.step} executa em paralelo com ${lvlInfo.levelSize - 1} outra(s) peça(s) deste nível`
      : `Passo ${lvlInfo.step} executa sequencialmente após o nível anterior`;
    pills.push(pill(cls, lbl, tip));
  }
  box.replaceChildren(...pills);
}

function renderAllNodes() {
  const levels = computeExecutionLevels();
  renderExecutionStrip(levels);
  state.nodes.forEach(renderNode);
  renderEdges();
}

function updateEmpty() {
  $("empty").hidden = state.nodes.size > 0;
}

// ---------- Sub-grafos em Mega-Pipelines (Drill-Down & Breadcrumbs) ----------

function serializeCanvasSnapshot() {
  return {
    name: $("pipeline-name")?.value || "novo-dag",
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

function extractLoopCandidates(rawRef, label) {
  const out = [];
  const push = (v) => {
    const s = String(v || "").trim();
    if (s && !out.includes(s)) out.push(s);
  };
  if (rawRef) {
    const clean = String(rawRef).replace(/\\/g, "/").replace(/\.(loop|studio)\.json$/i, "").replace(/\/pipeline\.dag\.json$/i, "");
    const segs = clean.split("/").filter(Boolean);
    if (segs.length > 0) {
      push(segs[segs.length - 1]);
      if (segs.length > 1 && segs[0] === "loops") push(segs[1]);
    }
    push(clean);
    push(toKebabCase(clean));
  }
  if (label) {
    push(label);
    push(toKebabCase(label));
  }
  return out;
}

const VALID_LOOP_ID_RE = /^[a-z0-9][a-z0-9_-]{1,63}$/;

async function listKnownProcesses() {
  const combined = [];
  const seen = new Set();
  const pushItem = (rawId, rawName) => {
    const id = String(rawId || "").trim();
    const name = String(rawName || id).trim();
    if (!id || seen.has(id.toLowerCase())) return;
    seen.add(id.toLowerCase());
    combined.push({ id, name });
  };
  await Promise.all([
    fetch("/api/assets/loop").then((r) => (r.ok ? r.json() : null)).then((d) => {
      for (const it of (d?.items || [])) pushItem(it.meta?.loop_id || it.id, it.name);
    }).catch(() => {}),
    fetch("/api/processos").then((r) => (r.ok ? r.json() : null)).then((d) => {
      for (const l of (d?.processos || d?.loops || [])) pushItem(l.id || l.nome, l.name || l.nome);
    }).catch(() => {}),
  ]);
  return combined;
}

async function resolveAndFetchLoop(rawRef, label, knownList) {
  const cands = extractLoopCandidates(rawRef, label);
  const loops = knownList || await listKnownProcesses();
  const tried = new Set();
  const tryFetch = async (id) => {
    if (!id || tried.has(id) || !VALID_LOOP_ID_RE.test(id)) return null;
    tried.add(id);
    try {
      const resp = await fetch(`/api/processos/${encodeURIComponent(id)}`);
      if (resp.ok) return await resp.json();
    } catch (_) {}
    return null;
  };
  if (loops.length > 0) {
    const byKey = new Map();
    loops.forEach((l) => {
      if (l.id) byKey.set(String(l.id).toLowerCase(), l.id);
      if (l.name) byKey.set(toKebabCase(l.name), l.id);
    });
    for (const c of cands) {
      const match = byKey.get(String(c).toLowerCase()) || byKey.get(toKebabCase(c));
      const doc = await tryFetch(match);
      if (doc) return doc;
    }
  }
  for (const id of cands) {
    const doc = await tryFetch(id);
    if (doc) return doc;
  }
  return null;
}

function nodeSubloopTarget(n) {
  if (!n) return "";
  const t = String(n.type || n.piece?.kind || "").toLowerCase();
  if (n.subloop) return String(n.subloop);
  if (t === "subloop" || t === "loop" || t === "subprocesso" || t === "flow_loop") {
    return String(n.piece?.ref || n.ref || n.id || "");
  }
  return "";
}

async function collectLoopTreeDocs(rootId, rootDoc, treeMap, visited, depth, knownList) {
  if (!rootDoc || depth > 4) return;
  const key = toKebabCase(rootId || rootDoc.id || rootDoc.name);
  if (visited.has(key)) return;
  visited.add(key);
  treeMap.set(key, rootDoc);
  const subRefs = [];
  for (const n of (rootDoc.nodes || [])) {
    const target = nodeSubloopTarget(n);
    if (target) subRefs.push({ target, label: n.name || n.label });
  }
  await Promise.all(subRefs.map(async (s) => {
    const k = toKebabCase(extractLoopCandidates(s.target, s.label)[0] || s.target);
    if (visited.has(k)) return;
    const subDoc = await resolveAndFetchLoop(s.target, s.label, knownList);
    if (subDoc) await collectLoopTreeDocs(subDoc.id || k, subDoc, treeMap, visited, depth + 1, knownList);
  }));
}

function buildNodeFromProcessDoc(n, id, x, y, depth) {
  const defaultNat = depth === 0 ? "macro" : depth === 1 ? "micro" : "nano";
  const nat = String(n.nature || defaultNat).toLowerCase();
  if (n.block && state.blocks.has(n.block)) {
    const blk = state.blocks.get(n.block);
    return {
      ...baseNode(id, x, y),
      block: n.block,
      nature: nat,
      routing: n.routing && n.routing.mode ? { ...defaultRouting(blk), ...n.routing } : defaultRouting(blk),
      slots: n.slots || emptySlots(),
      execution: n.execution?.mode || n.execution || "inline",
    };
  }
  const t = String(n.type || n.piece?.kind || "skill").toLowerCase();
  const mappedKind = t === "subprocesso" || t === "loop" || t === "flow_loop" ? "subloop" : (PIECE_KIND[t] || t || "skill");
  const ref = n.piece?.ref || n.ref || n.id || id;
  const label = n.name || n.label || String(ref).split("/").pop() || n.id || id;
  const artifact = (n.output && n.output.artifact) || n.piece?.artifact || `${slugOf(label).toUpperCase()}.md`;
  const gate = (n.output && n.output.gate) || n.piece?.gate || defaultGate(artifact);
  return {
    ...baseNode(id, x, y),
    label,
    nature: nat,
    piece: { kind: mappedKind, ref, artifact, gate },
    routing: n.routing && n.routing.mode ? { mode: n.routing.mode, model: n.routing.model || "", fallback: n.routing.fallback || "", combo: n.routing.combo || "" } : defaultRouting({}),
    slots: n.slots || emptySlots(),
    execution: n.execution?.mode || n.execution || "inline",
  };
}

function triggerBlockIdForConfig(trg) {
  const t = String(trg?.type || "manual").toLowerCase();
  const candidates = [`trigger_${t}`, "trigger_manual", "trigger"];
  for (const c of candidates) {
    if (state.blocks.has(c)) return c;
  }
  for (const [id, b] of state.blocks.entries()) {
    if (b.kind && b.kind.startsWith("trigger_")) return id;
  }
  return "";
}

function appendLoopTreeToCanvas(doc, treeMap, prefix, entryFromId, depth, expandedSet, nodeIndexRef) {
  const src = doc.nodes || [];
  const localMap = new Map();
  const childrenMap = new Map();

  src.forEach((n) => {
    const uid = prefix ? `${prefix}_${n.id}` : n.id;
    const col = nodeIndexRef.idx++;
    const node = buildNodeFromProcessDoc(n, uid, 60 + col * (NODE_W + 70), 80 + depth * 180, depth);
    state.nodes.set(uid, node);
    localMap.set(n.id, uid);
    childrenMap.set(uid, []);
  });

  const hasIncoming = new Set();
  src.forEach((n) => {
    const toUid = localMap.get(n.id);
    (n.depends_on || []).forEach((dep) => {
      const fromUid = localMap.get(dep);
      if (fromUid && toUid && fromUid !== toUid) {
        state.edges.push({ from: fromUid, to: toUid });
        hasIncoming.add(toUid);
      }
    });
  });
  (doc.edges || []).forEach((e) => {
    const fromUid = localMap.get(e.from);
    const toUid = localMap.get(e.to);
    if (fromUid && toUid && fromUid !== toUid && !state.edges.some((x) => x.from === fromUid && x.to === toUid)) {
      state.edges.push({ from: fromUid, to: toUid });
      hasIncoming.add(toUid);
    }
  });

  if (entryFromId) {
    src.forEach((n) => {
      const uid = localMap.get(n.id);
      if (uid && !hasIncoming.has(uid)) {
        state.edges.push({ from: entryFromId, to: uid });
      }
    });
  }

  // Expande sub-processos referenciados para exibir toda a árvore MACRO -> MICRO -> NANO.
  src.forEach((n) => {
    const subTarget = nodeSubloopTarget(n);
    if (!subTarget) return;
    const uid = localMap.get(n.id);
    const cands = extractLoopCandidates(subTarget, n.name || n.label);
    const hitKey = cands.map(toKebabCase).find((k) => treeMap.has(k) && !expandedSet.has(k));
    if (hitKey && uid) {
      expandedSet.add(hitKey);
      appendLoopTreeToCanvas(treeMap.get(hitKey), treeMap, uid, uid, depth + 1, expandedSet, nodeIndexRef);
    }
  });
}

async function enterSubloop(subloopId, subloopLabel) {
  if (!subloopId) return;
  const currentSnap = serializeCanvasSnapshot();
  currentSnap.label = currentSnap.name || "Principal";
  state.navStack.push(currentSnap);

  try {
    const knownList = await listKnownProcesses();
    const loopData = await resolveAndFetchLoop(subloopId, subloopLabel, knownList);
    if (!loopData) {
      toast("Não foi possível carregar o processo completo do acervo", true);
      state.navStack.pop();
      return;
    }

    const treeMap = new Map();
    await collectLoopTreeDocs(loopData.id || subloopId, loopData, treeMap, new Set(), 0, knownList);

    pushState();
    [...state.nodes.keys()].forEach((id) => document.querySelector(`.node[data-id="${CSS.escape(id)}"]`)?.remove());
    state.nodes.clear();
    state.edges = [];

    if ($("pipeline-name")) {
      $("pipeline-name").value = loopData.name || loopData.id || subloopId;
    }

    let rootEntryId = null;
    const trgBlk = triggerBlockIdForConfig(loopData.trigger);
    if (trgBlk) {
      const blk = state.blocks.get(trgBlk);
      const trgNode = {
        ...baseNode("trg_root", 40, 80),
        block: trgBlk,
        nature: "macro",
        triggerConfig: normalizeTriggerConfig(blk.kind, loopData.trigger),
        routing: defaultRouting(blk),
      };
      state.nodes.set(trgNode.id, trgNode);
      rootEntryId = trgNode.id;
    }

    const rootKey = toKebabCase(loopData.id || loopData.name || subloopId);
    const expanded = new Set([rootKey]);
    appendLoopTreeToCanvas(loopData, treeMap, "", rootEntryId, 0, expanded, { idx: 1 });

    clearValidation();
    renderAllNodes();
    updateEmpty();
    renderBreadcrumbs();
    autoLayout({ fit: false, notify: false });
    focusViewportStart(0.9);
    select(null);
    toast(`Processo aberto (${state.nodes.size} nó(s) MACRO/MICRO/NANO): ${subloopLabel || loopData.name || subloopId}`);
  } catch (err) {
    toast("Erro ao abrir processo/sub-grafo: " + err.message, true);
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
  const root = typeof document !== "undefined" && typeof CSS !== "undefined"
    ? document.querySelector(`.node[data-id="${CSS.escape(node.id)}"]`)
    : null;
  const port = root?.querySelector(side === "out" ? ".port--out" : ".port--in");
  if (port && root && (port.offsetWidth > 0 || port.offsetHeight > 0)) {
    return {
      x: node.x + port.offsetLeft + port.offsetWidth / 2,
      y: node.y + port.offsetTop + port.offsetHeight / 2,
    };
  }
  return { x: node.x + (side === "out" ? NODE_W : 0), y: node.y + PORT_Y };
}

function curvePoints(a, b) {
  const dx = b.x - a.x;
  const dy = b.y - a.y;
  if (dx >= 28) {
    const offset = Math.max(56, Math.min(220, Math.abs(dx) * 0.45));
    return {
      c1: { x: a.x + offset, y: a.y },
      c2: { x: b.x - offset, y: b.y },
    };
  }
  // Alvo à esquerda ou alinhado verticalmente: contorna os cards sem retorcer a curva
  const signY = dy >= 0 ? 1 : -1;
  const spanY = Math.max(76, Math.min(180, Math.abs(dy) * 0.55 + 64));
  const arcX = Math.max(84, Math.min(160, Math.abs(dx) * 0.25 + 84));
  return {
    c1: { x: a.x + arcX, y: a.y + signY * spanY },
    c2: { x: b.x - arcX, y: b.y + signY * spanY },
  };
}

function curve(a, b) {
  const { c1, c2 } = curvePoints(a, b);
  return `M ${a.x} ${a.y} C ${c1.x} ${c1.y}, ${c2.x} ${c2.y}, ${b.x} ${b.y}`;
}

function curveMidpoint(a, b) {
  const { c1, c2 } = curvePoints(a, b);
  // Avaliação exata da Bézier cúbica em t = 0.5
  return {
    x: 0.125 * a.x + 0.375 * c1.x + 0.375 * c2.x + 0.125 * b.x,
    y: 0.125 * a.y + 0.375 * c1.y + 0.375 * c2.y + 0.125 * b.y,
  };
}

function buildArrowMarker(id, color) {
  const m = document.createElementNS(SVG_NS, "marker");
  m.setAttribute("id", id);
  m.setAttribute("viewBox", "0 0 10 10");
  m.setAttribute("refX", "8");
  m.setAttribute("refY", "5");
  m.setAttribute("markerWidth", "6");
  m.setAttribute("markerHeight", "6");
  m.setAttribute("orient", "auto-start-reverse");
  const p = document.createElementNS(SVG_NS, "path");
  p.setAttribute("d", "M 1 1.5 L 8.5 5 L 1 8.5 z");
  p.setAttribute("fill", color);
  m.append(p);
  return m;
}

function ensureEdgeDefs(svg) {
  const defs = document.createElementNS(SVG_NS, "defs");
  defs.append(
    buildArrowMarker("edge-arrow-sync", "var(--edge)"),
    buildArrowMarker("edge-arrow-async", "#0284c7"),
    buildArrowMarker("edge-arrow-selected", "var(--accent)"),
    buildArrowMarker("edge-arrow-invalid", "var(--fail)")
  );
  svg.append(defs);
}

// Gatilho (trigger_*) conecta universalmente com qualquer peça; peças/fluxo/Any ligam livremente; blocos tipados usam a matriz.
function canConnect(fromId, toId) {
  const a = state.nodes.get(fromId);
  const b = state.nodes.get(toId);
  if (!a || !b || fromId === toId) return false;
  const ba = blockOf(a);
  const bb = blockOf(b);
  if (ba?.kind?.startsWith("trigger_") || bb?.kind?.startsWith("trigger_")) return true;
  if (a.piece || b.piece) return true;
  if (ba?.kind?.startsWith("flow_") || bb?.kind?.startsWith("flow_")) return true;
  if (ba?.produces === ANY_TYPE || (bb?.consumes || []).includes(ANY_TYPE)) return true;
  return (state.matrix[a.block] || []).includes(b.block);
}

function isAsyncEdge(e) {
  if (e.async || e.mode === "async") return true;
  const a = state.nodes.get(e.from);
  const b = state.nodes.get(e.to);
  if (a?.async || b?.async || a?.mode === "async" || b?.mode === "async") return true;
  const lvl = state.execLevels?.get(e.to);
  if (lvl && lvl.mode === "async") return true;
  const fromLvl = state.execLevels?.get(e.from);
  return Boolean(fromLvl && fromLvl.outdeg > 1);
}

function edgeClass(e, i) {
  const sel = state.selected?.kind === "edge" && state.selected.id === i;
  const bad = state.badEdges.has(`${e.from}→${e.to}`) || !canConnect(e.from, e.to);
  const modeCls = isAsyncEdge(e) ? "edge--async" : "edge--sync";
  return `edge ${modeCls}${sel ? " is-selected" : ""}${bad ? " is-invalid" : ""}`;
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
  ensureEdgeDefs(svg);
  state.edges.forEach((e, i) => {
    const a = state.nodes.get(e.from);
    const b = state.nodes.get(e.to);
    if (!a || !b) return;
    const pA = portPoint(a, "out");
    const end = portPoint(b, "in");
    const d = curve(pA, end);
    const hit = document.createElementNS(SVG_NS, "path");
    hit.setAttribute("d", d);
    hit.setAttribute("class", "edge-hit");
    hit.dataset.edge = String(i);
    const path = document.createElementNS(SVG_NS, "path");
    path.setAttribute("d", d);
    const sel = state.selected?.kind === "edge" && state.selected.id === i;
    const bad = state.badEdges.has(`${e.from}→${e.to}`) || !canConnect(e.from, e.to);
    const asyncEdge = isAsyncEdge(e);
    path.setAttribute("class", edgeClass(e, i));
    const markerId = bad ? "edge-arrow-invalid" : sel ? "edge-arrow-selected" : asyncEdge ? "edge-arrow-async" : "edge-arrow-sync";
    path.setAttribute("marker-end", `url(#${markerId})`);
    path.dataset.edge = String(i);
    svg.append(path, hit);

    // Botão de desconexão rápida posicionado exatamente sobre o ponto médio da curva Bézier
    const mid = curveMidpoint(pA, end);
    const btnDel = document.createElementNS(SVG_NS, "g");
    btnDel.setAttribute("class", "edge-del-group");
    btnDel.dataset.edge = String(i);
    const tip = document.createElementNS(SVG_NS, "title");
    tip.textContent = `Remover conexão (${asyncEdge ? "Paralela/Assíncrona" : "Sequencial/Síncrona"})`;
    const circ = document.createElementNS(SVG_NS, "circle");
    circ.setAttribute("cx", String(mid.x));
    circ.setAttribute("cy", String(mid.y));
    circ.setAttribute("r", "11");
    circ.setAttribute("class", "edge-btn-remove");
    circ.dataset.edge = String(i);
    const txt = document.createElementNS(SVG_NS, "text");
    txt.setAttribute("x", String(mid.x));
    txt.setAttribute("y", String(mid.y + 4));
    txt.setAttribute("text-anchor", "middle");
    txt.setAttribute("font-size", "12");
    txt.setAttribute("font-weight", "bold");
    txt.setAttribute("fill", "#fff");
    txt.setAttribute("pointer-events", "none");
    txt.textContent = "x";
    btnDel.append(tip, circ, txt);
    svg.append(btnDel);

    if (sel) svg.append(edgeHandle(end, i));
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
  renderAllNodes();
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

const MAX_PAN = 1e7;

function updateCanvasTransform() {
  const surface = $("canvas-surface");
  if (surface) {
    state.pan.x = Math.max(-MAX_PAN, Math.min(MAX_PAN, state.pan.x || 0));
    state.pan.y = Math.max(-MAX_PAN, Math.min(MAX_PAN, state.pan.y || 0));
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

// Posiciona o viewport no primeiro nó com zoom legível, sem reduzir grafos grandes a miniatura nem travar o zoom livre.
function focusViewportStart(targetZoom = 0.9) {
  if (state.nodes.size === 0) {
    state.pan = { x: 40, y: 40 };
    state.zoom = 1.0;
    updateCanvasTransform();
    return;
  }
  let minX = Infinity, minY = Infinity;
  for (const n of state.nodes.values()) {
    minX = Math.min(minX, n.x);
    minY = Math.min(minY, n.y);
  }
  const zoom = Math.min(1.5, Math.max(0.55, targetZoom));
  state.zoom = zoom;
  state.pan.x = 48 - minX * zoom;
  state.pan.y = 48 - minY * zoom;
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
  insp.classList.toggle("is-open", !isCollapsed);
  const btn = $("btn-toggle-inspector");
  if (btn) {
    btn.textContent = isCollapsed ? "◂ Inspetor" : "Inspetor ▸";
    btn.setAttribute("aria-expanded", String(!isCollapsed));
  }
  const inlineBtn = $("btn-inspector-collapse");
  if (inlineBtn) {
    inlineBtn.textContent = isCollapsed ? "◂" : "▸";
    inlineBtn.setAttribute("aria-expanded", String(!isCollapsed));
    inlineBtn.title = isCollapsed ? "Expandir Inspetor" : "Recolher Inspetor";
  }
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

function autoLayout(opts) {
  if (state.nodes.size === 0) return toast("Canvas vazio.");
  const explicit = Boolean(opts && typeof opts === "object" && !opts.target && !opts.currentTarget);
  const doFit = explicit ? Boolean(opts.fit) : true;
  const doNotify = explicit ? (opts.notify !== false) : true;
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
      const elNode = document.querySelector(`.node[data-id="${CSS.escape(n.id)}"]`);
      if (elNode) {
        elNode.style.transform = `translate(${n.x}px, ${n.y}px)`;
      }
    });
  });
  renderEdges();
  if (doFit) fitToScreen();
  if (doNotify) toast("DAG organizado automaticamente!");
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
  const screenX = ev.clientX - r.left + (c.scrollLeft || 0);
  const screenY = ev.clientY - r.top + (c.scrollTop || 0);
  return {
    x: (screenX - state.pan.x) / state.zoom,
    y: (screenY - state.pan.y) / state.zoom,
  };
}

function startNodeDrag(ev, node) {
  const start = canvasPoint(ev);
  const origin = { x: node.x, y: node.y };
  const root = $(`node-${node.id}`) || document.querySelector(`.node[data-id="${CSS.escape(node.id)}"]`);
  let captured = false;

  const move = (e) => {
    const p = canvasPoint(e);
    const dx = p.x - start.x;
    const dy = p.y - start.y;
    if (!captured && (Math.abs(dx) > 3 || Math.abs(dy) > 3)) {
      captured = true;
      if (root) {
        root.classList.add("is-dragging");
        try { root.setPointerCapture(ev.pointerId); } catch (_) {}
      }
    }
    if (!captured) return;
    node.x = origin.x + dx;
    node.y = origin.y + dy;
    renderNode(node);
    renderEdges();
  };
  const up = () => {
    if (root && captured) {
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

function findLinkTargetAt(e, anchorId, startPort) {
  const direct = document.elementFromPoint(e.clientX, e.clientY)?.closest(".node");
  if (direct && direct.dataset.id !== anchorId) return direct;
  const wantPort = startPort === "in" ? ".port--out" : ".port--in";
  let best = null;
  let bestDist = Infinity;
  document.querySelectorAll(".node").forEach((n) => {
    if (n.dataset.id === anchorId) return;
    const port = n.querySelector(wantPort);
    if (port) {
      const pr = port.getBoundingClientRect();
      const pad = 18;
      if (e.clientX >= pr.left - pad && e.clientX <= pr.right + pad && e.clientY >= pr.top - pad && e.clientY <= pr.bottom + pad) {
        const d = Math.hypot(e.clientX - (pr.left + pr.width / 2), e.clientY - (pr.top + pr.height / 2));
        if (d < bestDist) { bestDist = d; best = n; }
        return;
      }
    }
    const r = n.getBoundingClientRect();
    if (e.clientX >= r.left - 8 && e.clientX <= r.right + 8 && e.clientY >= r.top - 8 && e.clientY <= r.bottom + 8) {
      best = n;
    }
  });
  return best;
}

// trackLink desenha a curva temporária de uma conexão sendo arrastada a partir de uma porta ("out" ou "in").
function trackLink(ev, anchorId, onDrop, startPort = "out") {
  const anchor = state.nodes.get(anchorId);
  const temp = document.createElementNS(SVG_NS, "path");
  temp.setAttribute("class", "edge edge--temp");
  $("edges").append(temp);
  const canvas = $("canvas");
  canvas.classList.add("is-linking");
  document.querySelectorAll(".node").forEach((n) => {
    const ok = startPort === "in" ? canConnect(n.dataset.id, anchorId) : canConnect(anchorId, n.dataset.id);
    const sel = startPort === "in" ? ".port--out" : ".port--in";
    n.querySelector(sel)?.classList.toggle("is-compatible", ok);
  });
  const cleanup = () => {
    window.removeEventListener("pointermove", move);
    window.removeEventListener("pointerup", up);
    window.removeEventListener("keydown", onKeyDown);
    document.querySelectorAll(".port--in, .port--out").forEach((p) => p.classList.remove("is-compatible"));
    temp.remove();
    canvas.classList.remove("is-linking");
  };
  const onKeyDown = (e) => {
    if (e.key === "Escape") {
      cleanup();
      onDrop(null);
    }
  };
  const move = (e) => {
    const target = findLinkTargetAt(e, anchorId, startPort);
    const ok = target ? (startPort === "in" ? canConnect(target.dataset.id, anchorId) : canConnect(anchorId, target.dataset.id)) : true;
    temp.classList.toggle("is-bad", Boolean(target) && !ok);
    const pAnchor = portPoint(anchor, startPort);
    let pCursor = canvasPoint(e);
    if (target && ok) {
      const targetNode = state.nodes.get(target.dataset.id);
      if (targetNode) pCursor = portPoint(targetNode, startPort === "in" ? "out" : "in");
    }
    temp.setAttribute("d", startPort === "in" ? curve(pCursor, pAnchor) : curve(pAnchor, pCursor));
  };
  const up = (e) => {
    cleanup();
    onDrop(findLinkTargetAt(e, anchorId, startPort)?.dataset.id || null);
  };
  window.addEventListener("pointermove", move);
  window.addEventListener("pointerup", up);
  window.addEventListener("keydown", onKeyDown);
  move(ev);
}

function startLink(ev, anchorId, startPort = "out") {
  trackLink(ev, anchorId, (targetId) => {
    if (!targetId || targetId === anchorId) return;
    if (startPort === "in") addEdge(targetId, anchorId);
    else addEdge(anchorId, targetId);
  }, startPort);
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

  const portEl = ev.target.closest(".port");
  if (portEl && (portEl.dataset.port === "out" || portEl.dataset.port === "in")) {
    ev.preventDefault();
    ev.stopPropagation();
    startLink(ev, node.id, portEl.dataset.port);
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
      trig = normalizeTriggerConfig(b.kind, n.triggerConfig);
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

function compileDAG() {
  if (state.lastDAG) return state.lastDAG;
  const comp = composition();
  const incoming = new Map();
  comp.edges.forEach((e) => {
    if (!incoming.has(e.to)) incoming.set(e.to, []);
    incoming.get(e.to).push(e.from);
  });
  const nodes = [...state.nodes.values()].map((n) => {
    const block = blockOf(n);
    const deps = incoming.get(n.id) || [];
    const inputs = {};
    deps.forEach((fromId) => {
      const src = state.nodes.get(fromId);
      const srcArt = src ? blockOf(src).artifact : `${fromId}.md`;
      inputs[fromId] = `file('${srcArt}')`;
    });
    const pieceKind = n.piece?.kind;
    const type = pieceKind === "subloop" ? "subprocesso" : (pieceKind || (block.kind === "mcp" ? "mcp" : "skill"));
    const out = {
      id: n.id,
      type,
      ref: n.piece ? n.piece.ref : block.ref,
      routing: n.routing,
      output: { artifact: block.artifact, gate: block.gate || defaultGate(block.artifact) },
    };
    if (deps.length > 0) out.depends_on = deps;
    if (Object.keys(inputs).length > 0) out.inputs = inputs;
    if (n.execution === "isolated") out.execution = { mode: "isolated" };
    return out;
  });
  return {
    $schema: "https://legoloop.dev/schemas/dag.v1.json",
    name: comp.name || "novo-dag",
    version: "2.0.0",
    trigger: comp.trigger || { type: "manual" },
    nodes,
  };
}

function clearValidation() {
  state.lastDAG = null;
  if (state.status.size === 0 && state.badEdges.size === 0) return;
  state.status.clear();
  state.badEdges.clear();
  $("issues").hidden = true;
  renderAllNodes();
}

function applyValidation(res) {
  state.status.clear();
  state.badEdges.clear();
  state.lastDAG = res?.ok && res?.dag ? res.dag : null;
  const issues = res.problemas || res.issues || [];
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
    li.append(el("span", "issues__where", i.node || i.edge || "DAG"), el("span", "", i.message));
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
  const issues = data.problemas || data.issues || [];
  toast(data.ok ? "DAG válido: todos os encaixes conferem." : `${issues.length} problema(s) encontrado(s).`, !data.ok);
  return Boolean(data.ok);
}

async function materialize() {
  const { ok, data } = await api("/api/materialize", composition());
  if (data.problemas || data.issues) applyValidation(data);
  if (!ok) {
    toast(data.error || "Corrija os problemas antes de materializar.", true);
    return false;
  }
  toast(`DAG materializado em ${data.path}.`);
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
  let cx = Math.round((-state.pan.x + cv.clientWidth / 2) / state.zoom - 118);
  let cy = Math.round((-state.pan.y + cv.clientHeight / 2) / state.zoom - 50);
  const existing = [...state.nodes.values()];
  for (let i = 0; i < 16 && existing.some((n) => Math.abs(n.x - cx) < 36 && Math.abs(n.y - cy) < 28); i++) {
    cx += 44;
    cy += 36;
  }
  const added = addEntry(entry, cx, cy);
  if (added) toast(`Peça "${entry.name}" adicionada ao Canvas`);
  return added;
}

let activeModalItem = null;

async function simulateCurrentDagInSwarm() {
  if (state.nodes.size === 0) {
    toast("Adicione pelo menos uma peça ao canvas antes de simular no Swarm.", true);
    return;
  }
  const name = $("pipeline-name")?.value.trim() || "novo-dag";
  const nodes = Array.from(state.nodes.values()).map((n) => {
    const blk = blockOf(n);
    return {
      id: n.id,
      label: n.label || blk.label || n.id,
      kind: n.piece?.kind || blk.kind || "skill",
      role: n.role || blk.role || "",
      block: n.block || "",
      piece: n.piece?.ref || n.block || n.id,
      execution: n.execution || "shared",
      artifact: n.piece?.artifact || blk.artifact || "",
      gate: n.piece?.gate || blk.gate || "",
    };
  });
  const edges = state.edges.map((e) => ({ from: e.from, to: e.to }));
  const res = await api("/api/swarm/simulate", { name, nodes, edges });
  if (!res.ok) {
    toast("Falha ao simular no Swarm: " + (res.data?.error || res.status), true);
    return;
  }
  const url = res.data?.url || `/swarm?proc=${encodeURIComponent(res.data?.slug || "")}`;
  window.location.href = url;
}

function entryFromCompatItem(raw) {
  if (!raw) return null;
  const cat = raw.cat || (raw.kind === "subloop" ? "loop" : raw.kind) || "skill";
  return {
    cat,
    kind: raw.kind || cat,
    name: raw.name || raw.id,
    itemId: raw.id || raw.name,
    blockId: String(raw.kind || "").includes("_") || raw.kind === "control" ? (raw.id || raw.name) : undefined,
    meta: raw.meta || {},
  };
}

async function openItemModal(it, sourceNode = null) {
  const modal = $("modal-item");
  if (!modal) return;
  activeModalItem = it;
  const resolveTargetNode = () => {
    if (sourceNode && state.nodes.has(sourceNode.id)) return sourceNode;
    if (state.selected?.kind === "node") return state.nodes.get(state.selected.id) || null;
    return null;
  };
  modal.open(it, {
    showCanvasAction: true,
    onAddToCanvas: (entry) => addEntryToCenter(entry),
    onCanonize: () => loadItems(),
    onUpgrade: async () => {
      await loadItems();
      renderAllNodes();
      renderInspector();
    },
    onSimulateSwarm: () => simulateCurrentDagInSwarm(),
    onSwapCompatible: (compat) => {
      const target = resolveTargetNode();
      const entry = entryFromCompatItem(compat);
      if (!target) return toast("Selecione um card no canvas para trocar por esta peça.", true);
      if (!entry) return;
      swapNode(target.id, entry);
    },
    onConnectBefore: (compat) => {
      const target = resolveTargetNode();
      const entry = entryFromCompatItem(compat);
      if (!entry) return;
      const x = target ? target.x - 280 : 120;
      const y = target ? target.y : 120;
      const newNode = addEntry(entry, x, y);
      if (newNode && target) addEdge(newNode.id, target.id);
    },
    onConnectAfter: (compat) => {
      const target = resolveTargetNode();
      const entry = entryFromCompatItem(compat);
      if (!entry) return;
      const x = target ? target.x + 280 : 380;
      const y = target ? target.y : 120;
      const newNode = addEntry(entry, x, y);
      if (newNode && target) addEdge(target.id, newNode.id);
    },
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
  // Canvas Pan & Zoom livre (roda faz zoom direto centrado no cursor; Shift+roda faz pan)
  const cv = $("canvas");
  cv.addEventListener("wheel", (ev) => {
    ev.preventDefault();
    if (ev.shiftKey && !ev.ctrlKey && !ev.metaKey) {
      if (ev.deltaX === 0) {
        state.pan.x -= ev.deltaY;
      } else {
        state.pan.x -= ev.deltaX;
        state.pan.y -= ev.deltaY;
      }
      updateCanvasTransform();
      return;
    }
    const r = cv.getBoundingClientRect();
    const factor = ev.deltaY < 0 ? 1.12 : 1 / 1.12;
    zoomBy(factor, { x: ev.clientX - r.left, y: ev.clientY - r.top });
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
  $("btn-inspector-collapse")?.addEventListener("click", toggleInspector);
  const inspRail = document.querySelector(".inspector-collapsed-rail");
  inspRail?.addEventListener("click", toggleInspector);
  inspRail?.addEventListener("keydown", (e) => {
    if (e.key === "Enter" || e.key === " ") { e.preventDefault(); toggleInspector(); }
  });
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
  $("btn-simular-swarm")?.addEventListener("click", (e) => withBusy(e.currentTarget, simulateCurrentDagInSwarm));
  $("btn-materialize").addEventListener("click", (e) => withBusy(e.currentTarget, materialize));
  $("btn-route").addEventListener("click", (e) => withBusy(e.currentTarget, traceRoute));
  $("btn-close-issues").addEventListener("click", () => { $("issues").hidden = true; });
  initVaultModal();
  initRunModal();
  initTemplatesModal();
  initImportModal();
  initFabricarModal();
  initProcessosModal();
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

    const hasScripts = Array.from(state.nodes.values()).some((n) => (n.piece?.kind || blockOf(n)?.kind) === "script");
    let allowScripts = false;
    if (hasScripts) {
      const agreed = await confirmDialog({
        title: "Autorizar Execução de Scripts?",
        message: "Este DAG possui nós do tipo script que executarão comandos no sistema operacional. Você autoriza a execução?",
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
    nameEl.textContent = dag.name || "novo-dag";
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
              appendLog(`[${new Date().toLocaleTimeString()}] DAG "${evData.dag}" iniciado (${evData.nodes_total} nós).`);

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
              appendLog(`[${new Date().toLocaleTimeString()}] [RUN] Nó "${evData.id}" em execução...`);

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
              toast(ok ? "DAG executado com sucesso!" : "DAG finalizou com falhas.", !ok);
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
              title: "Substituir DAG pelo Template?",
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

// ---------- Modal "Novo Processo" (fábrica) ----------
// A lógica pura está em fabricar.js. Aqui só o DOM e as chamadas à API; erros aparecem no próprio modal.

function fabricarElements() {
  const ids = ["modal-fabricar", "fabricar-name", "fabricar-category", "fabricar-model", "fabricar-start", "fabricar-skip",
    "fabricar-interview", "fabricar-confidence", "fabricar-question", "fabricar-guess", "fabricar-answer",
    "fabricar-send", "fabricar-run", "fabricar-status", "fabricar-error", "fabricar-overwrite"];
  const ui = {};
  for (const id of ids) ui[id.replace("fabricar-", "").replace("modal-fabricar", "modal")] = $(id);
  return ui;
}

function fabricarShowError(ui, msg, canOverwrite = false) {
  ui.error.textContent = msg;
  ui.error.hidden = !msg;
  ui.overwrite.hidden = !canOverwrite;
}

function fabricarShowStatus(ui, msg) {
  ui.status.textContent = msg;
  ui.status.hidden = !msg;
}

function fabricarFillModels(ui) {
  const opts = fabricarModelOptions((state.inv && state.inv.models) || []);
  ui.model.replaceChildren(...opts.map((o) => {
    const opt = document.createElement("option");
    opt.value = o.id;
    opt.textContent = o.label;
    opt.disabled = o.disabled;
    return opt;
  }));
  ui.model.value = fabricarDefaultModel(opts);
  ui.start.disabled = ui.skip.disabled = !ui.model.value;
  if (!ui.model.value) fabricarShowError(ui, "Nenhum modelo do 9Router disponível. Rode legoloop doctor.");
}

function fabricarForm(ui) {
  return { name: ui.name.value, category: ui.category.value, model: ui.model.value };
}

function fabricarReset(ui, flow) {
  flow.sessionId = "";
  flow.turn = null;
  ui.interview.hidden = true;
  ui.name.value = "";
  fabricarShowError(ui, "");
  fabricarShowStatus(ui, "");
  fabricarFillModels(ui);
}

function fabricarShowTurn(ui, flow, turn) {
  flow.turn = turn;
  ui.interview.hidden = false;
  ui.confidence.textContent = `Confiança da entrevista: ${turn.confidence || 0}%`;
  const done = fabricarInterviewDone(turn);
  ui.question.textContent = done ? "Entrevista concluída. Pode fabricar o processo." : (turn.question || "");
  ui.guess.hidden = done || !turn.guess;
  ui.guess.textContent = turn.guess ? `Aceitar o palpite: ${turn.guess}` : "";
  ui.answer.value = "";
  ui.answer.hidden = ui.send.hidden = done;
}

async function fabricarStart(ui, flow) {
  fabricarShowError(ui, "");
  const bad = fabricarNameError(ui.name.value);
  if (bad) return fabricarShowError(ui, bad);
  fabricarShowStatus(ui, "Iniciando a entrevista...");
  const res = await api("/api/interview/start", { idea: ui.name.value.trim() });
  fabricarShowStatus(ui, "");
  if (!res.ok || !res.data.session) return fabricarShowError(ui, fabricarErrorMessage(res.status, res.data));
  flow.sessionId = res.data.session.id;
  fabricarShowTurn(ui, flow, res.data.turn);
}

async function fabricarAnswer(ui, flow, answer) {
  const text = String(answer || "").trim();
  if (!flow.sessionId || !text) return fabricarShowError(ui, "Escreva uma resposta.");
  fabricarShowError(ui, "");
  fabricarShowStatus(ui, "Processando a resposta...");
  const res = await api("/api/interview/message", { session_id: flow.sessionId, answer: text });
  fabricarShowStatus(ui, "");
  if (!res.ok || !res.data.turn) return fabricarShowError(ui, fabricarErrorMessage(res.status, res.data));
  fabricarShowTurn(ui, flow, res.data.turn);
}

async function fabricarRun(ui, flow, opts = {}) {
  flow.lastRun = { skipInterview: !!opts.skipInterview };
  fabricarShowError(ui, "");
  const bad = fabricarNameError(ui.name.value);
  if (bad) return fabricarShowError(ui, bad);
  fabricarShowStatus(ui, "Fabricando o processo. Isso pode levar alguns minutos...");
  let res;
  try {
    const turn = opts.skipInterview ? null : flow.turn;
    const sessId = opts.skipInterview ? "" : flow.sessionId;
    res = await api("/api/fabricar/run", fabricarPayload(fabricarForm(ui), turn, opts.force, sessId));
  } catch (e) {
    res = { ok: false, status: 0, data: null };
  }
  fabricarShowStatus(ui, "");
  if (!res.ok) return fabricarShowError(ui, fabricarErrorMessage(res.status, res.data), fabricarCanOverwrite(res.status, res.data));
  ui.modal.close();
  toast(`Processo ${res.data.canonical_name} criado.`);
  await pollItems();
  await fabricarPlot(res.data);
}

// Busca o loop gerado e o planta no Canvas; o Inspetor abre com a Tríade e a peça nova já está em Loops.
async function fabricarPlot(result) {
  if (fabricarNeedsReplaceConfirm(state.nodes.size)) {
    const ok = await confirmDialog({
      title: "Substituir o DAG do Canvas?",
      message: `O Canvas tem ${state.nodes.size} nó(s). Plotar "${result.canonical_name}" substitui esse DAG (Ctrl+Z desfaz).`,
      confirmLabel: "Substituir", danger: true,
    });
    if (!ok) {
      toast(`Processo ${result.canonical_name} criado. Abra-o em Loops para vê-lo no Canvas.`);
      return;
    }
  }
  const res = await api(`/api/processos/${encodeURIComponent(result.canonical_name)}`);
  const model = fabricarPlotModel(result, res.ok ? res.data : null);
  fabricarApplyPlot(model);
  if (!state.items.loop.some((i) => i.name === model.loopEntry.name)) {
    toast("O processo foi criado, mas ainda não apareceu em Loops. Atualize o inventário.", true);
  }
}

function fabricarApplyPlot(model) {
  pushState();
  [...state.nodes.keys()].forEach(removeNode);
  if ($("pipeline-name")) $("pipeline-name").value = model.name;
  model.nodes.forEach((n) => {
    const node = { ...baseNode(n.id, n.x, n.y), label: n.label, piece: n.piece, routing: n.routing };
    state.nodes.set(node.id, node);
    renderNode(node);
  });
  state.edges = model.edges;
  state.process = model.inspector;
  clearValidation();
  renderAllNodes();
  renderEdges();
  updateEmpty();
  select({ kind: "process", id: model.name });
}

// Sobrescrever é destrutivo: confirma antes. O servidor guarda o original como .legoloop-bak.
async function fabricarOverwrite(ui, flow) {
  const ok = await confirmDialog({
    title: "Sobrescrever o processo existente?",
    message: "O processo atual será substituído. O original fica guardado como .legoloop-bak ao lado de cada arquivo.",
    confirmLabel: "Sobrescrever", danger: true,
  });
  if (ok) await fabricarRun(ui, flow, { ...flow.lastRun, force: true });
}

function initFabricarModal() {
  const open = $("btn-novo-processo");
  const ui = fabricarElements();
  if (!open || !ui.modal) return;
  const flow = { sessionId: "", turn: null };
  const busy = (btn, fn) => btn.addEventListener("click", () => withBusy(btn, fn));
  open.addEventListener("click", () => { fabricarReset(ui, flow); ui.modal.showModal(); ui.name.focus(); });
  $("btn-modal-fabricar-close")?.addEventListener("click", () => ui.modal.close());
  busy(ui.start, () => fabricarStart(ui, flow));
  busy(ui.skip, () => fabricarRun(ui, flow, { skipInterview: true }));
  busy(ui.overwrite, () => fabricarOverwrite(ui, flow));
  busy(ui.send, () => fabricarAnswer(ui, flow, ui.answer.value));
  busy(ui.guess, () => fabricarAnswer(ui, flow, flow.turn && flow.turn.guess));
  busy(ui.run, () => fabricarRun(ui, flow));
  ui.answer.addEventListener("keydown", (e) => { if (e.key === "Enter") ui.send.click(); });
}

function initProcessosModal() {
  const openBtn = $("btn-processos");
  const modal = $("modal-processos");
  if (!openBtn || !modal) return;

  const closeBtn = $("btn-modal-processos-close");
  const refreshBtn = $("btn-processos-refresh");
  const searchInput = $("processos-search");
  const list = $("processos-list");
  const actionPanel = $("processos-action-panel");
  const actionTitle = $("processos-action-title");
  const formRefinar = $("form-refinar");
  const formCompor = $("form-compor");
  const formRemover = $("form-remover");
  const previewBtn = $("btn-processo-preview");
  const voltarBtn = $("btn-processo-voltar");
  const previewContainer = $("processos-preview-container");
  const cartaoArea = $("processos-cartao-area");
  const relatorioArea = $("processos-relatorio-area");
  const applyBtn = $("btn-processo-apply");

  let processosCache = [];
  let activeAction = "";
  let targetProc = null;
  let currentPlan = null;

  function renderLista(itens) {
    if (!list) return;
    list.replaceChildren();
    if (!itens || itens.length === 0) {
      const p = document.createElement("p");
      p.className = "note";
      p.textContent = "Nenhum processo cadastrado no cofre.";
      list.append(p);
      return;
    }
    for (const item of itens) {
      if (typeof renderProcessoItem === "function") {
        list.append(renderProcessoItem(item, iniciarAcao));
      }
    }
  }

  async function carregarProcessos() {
    const res = await api("/api/processos");
    if (!res.ok) {
      toast("Falha ao listar processos: " + (res.data?.error || res.status), true);
      return;
    }
    processosCache = res.data?.processos || res.data || [];
    filtrarLista();
  }

  function filtrarLista() {
    const term = (searchInput?.value || "").toLowerCase().trim();
    if (!term) {
      renderLista(processosCache);
      return;
    }
    const filtrados = processosCache.filter((p) => {
      const nome = (p.nome || "").toLowerCase();
      const desc = (p.descricao || p.objetivo || "").toLowerCase();
      return nome.includes(term) || desc.includes(term);
    });
    renderLista(filtrados);
  }

  function iniciarAcao(action, proc) {
    activeAction = action;
    targetProc = proc;
    currentPlan = null;
    if (actionPanel) actionPanel.hidden = false;
    if (previewContainer) previewContainer.hidden = true;
    if (formRefinar) formRefinar.hidden = action !== "refinar";
    if (formCompor) formCompor.hidden = action !== "compor";
    if (formRemover) formRemover.hidden = action !== "remover";

    const rotulos = { refinar: "Refinar Processo", compor: "Compor com Processo", remover: "Remover Processo" };
    if (actionTitle) actionTitle.textContent = `${rotulos[action] || action}: ${proc.nome}`;

    if (action === "refinar") {
      if ($("refinar-etapa")) $("refinar-etapa").value = "";
      if ($("refinar-peca")) $("refinar-peca").value = "";
      if ($("refinar-titulo")) $("refinar-titulo").value = "";
    } else if (action === "compor") {
      if ($("compor-alvo")) $("compor-alvo").value = "";
      if ($("compor-titulo")) $("compor-titulo").value = "";
      if ($("compor-depende")) $("compor-depende").value = "";
    }
  }

  function resetAcao() {
    activeAction = "";
    targetProc = null;
    currentPlan = null;
    if (actionPanel) actionPanel.hidden = true;
    if (previewContainer) previewContainer.hidden = true;
    if (cartaoArea) cartaoArea.replaceChildren();
    if (relatorioArea) relatorioArea.replaceChildren();
  }

  function montarPayload(apply) {
    if (!targetProc) return null;
    if (activeAction === "refinar") {
      return {
        name: targetProc.nome,
        etapa: $("refinar-etapa")?.value.trim() || "",
        peca: $("refinar-peca")?.value.trim() || "",
        titulo: $("refinar-titulo")?.value.trim() || "",
        apply: Boolean(apply),
      };
    }
    if (activeAction === "compor") {
      const depStr = $("compor-depende")?.value || "";
      const deps = depStr.split(",").map((s) => s.trim()).filter(Boolean);
      return {
        name: targetProc.nome,
        composto: $("compor-alvo")?.value.trim() || "",
        titulo: $("compor-titulo")?.value.trim() || "",
        depende: deps,
        apply: Boolean(apply),
      };
    }
    if (activeAction === "remover") {
      return {
        name: targetProc.nome,
        apply: Boolean(apply),
      };
    }
    return null;
  }

  async function gerarPreview() {
    const payload = montarPayload(false);
    if (!payload) return;
    const res = await api(`/api/processos/${activeAction}`, payload);
    if (!res.ok) {
      toast("Erro ao planejar ação: " + (res.data?.error || res.status), true);
      return;
    }
    currentPlan = res.data;
    if (cartaoArea && typeof renderCartaoSTE === "function") {
      cartaoArea.replaceChildren(renderCartaoSTE(currentPlan.cartao));
    }
    if (relatorioArea && typeof renderRelatorioCamadas === "function") {
      relatorioArea.replaceChildren(renderRelatorioCamadas(currentPlan));
    }
    if (previewContainer) previewContainer.hidden = false;
  }

  async function aplicarAcao() {
    if (!currentPlan || !targetProc) return;
    const acaoNome = activeAction === "remover" ? "Remoção" : activeAction === "refinar" ? "Refinamento" : "Composição";
    const ok = await confirmDialog({
      title: `Confirmar ${acaoNome}?`,
      message: `Deseja realmente aplicar a alteração no processo "${targetProc.nome}"?`,
      confirmLabel: "Confirmar e Aplicar",
      danger: activeAction === "remover" || Boolean(currentPlan.cartao?.irreversivel),
    });
    if (!ok) return;

    const payload = montarPayload(true);
    const res = await api(`/api/processos/${activeAction}`, payload);
    if (!res.ok) {
      toast("Falha ao aplicar alteração: " + (res.data?.error || res.status), true);
      return;
    }
    toast(`Ação de ${activeAction} executada com sucesso!`);
    resetAcao();
    await carregarProcessos();
    await pollItems();
  }

  openBtn.addEventListener("click", () => {
    resetAcao();
    modal.showModal();
    carregarProcessos();
  });
  closeBtn?.addEventListener("click", () => modal.close());
  refreshBtn?.addEventListener("click", () => carregarProcessos());
  searchInput?.addEventListener("input", () => filtrarLista());
  voltarBtn?.addEventListener("click", () => resetAcao());
  previewBtn?.addEventListener("click", () => withBusy(previewBtn, gerarPreview));
  applyBtn?.addEventListener("click", () => withBusy(applyBtn, aplicarAcao));
}
