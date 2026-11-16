"use strict";

const PER_PAGE = 20;
const palette = ["#8ab4f8", "#f28b82", "#fdd663", "#81c995", "#c58af9"];
const DEFAULT_SWARM_ROLES = [
  { slug: "orquestrador", nome: "Orquestrador Geral", guilda: "maestro", ferramentas_tipicas: "LegoLoop CLI" },
  { slug: "planejador", nome: "Planejador Estrategico", guilda: "maestro", ferramentas_tipicas: "PLAN v2" },
  { slug: "product-owner", nome: "Product Owner", guilda: "maestro", ferramentas_tipicas: "INTENT / SPEC" },
  { slug: "scrum-master", nome: "Scrum Master", guilda: "maestro", ferramentas_tipicas: "Board / Tickets" },
  { slug: "arquiteto", nome: "Arquiteto de Sistemas", guilda: "arquiteto", ferramentas_tipicas: "Templates v2" },
  { slug: "arquiteto-dados", nome: "Arquiteto de Dados", guilda: "arquiteto", ferramentas_tipicas: "SQL / Schemas" },
  { slug: "engenheiro-ux", nome: "Engenheiro de UI/UX", guilda: "arquiteto", ferramentas_tipicas: "Tokens / Figma" },
  { slug: "engenheiro-backend", nome: "Engenheiro Backend", guilda: "construtor", ferramentas_tipicas: "Go / REST" },
  { slug: "engenheiro-frontend", nome: "Engenheiro Frontend", guilda: "construtor", ferramentas_tipicas: "HTML / JS / CSS" },
  { slug: "engenheiro-infra", nome: "Engenheiro Infra/Ops", guilda: "construtor", ferramentas_tipicas: "Docker / VPS" },
  { slug: "engenheiro-tdd", nome: "Engenheiro TDD", guilda: "verificador", ferramentas_tipicas: "Go Test / Node" },
  { slug: "engenheiro-secops", nome: "Engenheiro SecOps", guilda: "verificador", ferramentas_tipicas: "Security Gates" },
  { slug: "auditor-qualidade", nome: "Auditor de Qualidade", guilda: "verificador", ferramentas_tipicas: "Taxlinter / Funclen" },
  { slug: "engenheiro-qa", nome: "Engenheiro de QA", guilda: "verificador", ferramentas_tipicas: "Playwright / Axe" },
  { slug: "integrador-ci", nome: "Integrador CI/CD", guilda: "integrador", ferramentas_tipicas: "Git / Worktrees" },
  { slug: "engenheiro-release", nome: "Engenheiro de Release", guilda: "integrador", ferramentas_tipicas: "Fechar / Build" },
  { slug: "redator-tecnico", nome: "Redator Tecnico", guilda: "documentador", ferramentas_tipicas: "Manual Creator" },
  { slug: "editor-didatico", nome: "Editor LegoPress", guilda: "documentador", ferramentas_tipicas: "LegoPress / Typst" },
  { slug: "dissecador-web", nome: "Dissecador Web/REA", guilda: "analista", ferramentas_tipicas: "REA Web" },
  { slug: "analista-forense", nome: "Analista Forense", guilda: "analista", ferramentas_tipicas: "Video-REA / Trace" }
];
const DEFAULT_WORKSPACES = [
  { id: "inherit", label: "Herdado (mesmo workspace do coordenador)" },
  { id: "share", label: "Compartilhado (worktree git própria)" },
  { id: "branch", label: "Nova Branch (clone git isolado)" }
];
const roleTools = Object.fromEntries(DEFAULT_SWARM_ROLES.map(r => [r.slug, r.ferramentas_tipicas || "CLI"]));
const PIECE_KINDS = [
  { id: "", label: "Todos os tipos" },
  { id: "skill", label: "Skill" },
  { id: "mcp", label: "MCP" },
  { id: "loop", label: "Processo / Sub-DAG" },
  { id: "script", label: "Script" },
  { id: "gate", label: "Gate" },
  { id: "hook", label: "Hook" }
];

let curProc = (typeof URLSearchParams === "function" && window.location?.search)
  ? (new URLSearchParams(window.location.search).get("proc") || "")
  : "";
let closedProcs = {};
let addPanelOpen = false;
let availablePieces = [];
let swarmMeta = { roles: DEFAULT_SWARM_ROLES, models: [], workspaces: DEFAULT_WORKSPACES };
let selectedAddPieces = [];
let selectedEditPieces = [];
let selectedEditPiecesAgent = "";
const tabState = {};
const SVG_NS = "http://www.w3.org/2000/svg";

function getState(proc) {
  const key = proc || "__default__";
  if (!tabState[key]) {
    tabState[key] = { selectedAgent: "", page: 0, curParent: "", draftCmd: "", draftModel: "" };
  }
  return tabState[key];
}

function colorIdx(name) {
  let h = 0;
  const s = name || "";
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) | 0;
  return Math.abs(h) % palette.length;
}

function fmtDur(iso) {
  if (!iso) return "0s";
  const s = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000));
  const m = Math.floor(s / 60);
  return m > 0 ? `${m}m ${s % 60}s` : `${s}s`;
}

function el(tag, cls, text) {
  const node = document.createElement(tag);
  if (cls) node.className = cls;
  if (text !== undefined && text !== null) node.textContent = String(text);
  return node;
}

function actEl(tag, cls, act, val, text) {
  const node = el(tag, cls, text);
  node.setAttribute("data-act", act);
  node.setAttribute("data-val", val !== undefined && val !== null ? String(val) : "");
  if (tag === "button") node.type = "button";
  return node;
}

function botSVGNode(st, name, kind, anim) {
  const c = palette[colorIdx(name)];
  const svg = document.createElementNS ? document.createElementNS(SVG_NS, "svg") : el("svg");
  svg.setAttribute("class", "bot " + (anim !== false ? (st || "") : ""));
  svg.setAttribute("viewBox", "0 0 8 8");
  const rects = kind === "coord"
    ? [[2,0,4,3,c],[3,1,1,1,"#1e293b"],[5,1,1,1,"#1e293b"],[1,3,6,4,"#334155"],[3,3,2,4,"#fff"],[3,4,2,2,"#d9480f"],[2,7,1,1,"#334155"],[5,7,1,1,"#334155"]]
    : kind === "ver"
    ? [[2,0,4,3,c],[2,1,2,2,"#bae6fd"],[5,1,1,1,"#1e293b"],[1,3,5,4,c],[6,3,2,3,"#e2e8f0"],[2,7,1,1,c],[4,7,1,1,c]]
    : [[1,0,6,1,"#facc15"],[2,1,4,2,c],[3,2,1,1,"#1e293b"],[5,2,1,1,"#1e293b"],[1,3,6,3,c],[1,5,6,1,"#78350f"],[2,6,1,2,c],[5,6,1,2,c]];
  for (const [x, y, w, h, fill] of rects) {
    const r = document.createElementNS ? document.createElementNS(SVG_NS, "rect") : el("rect");
    r.setAttribute("x", String(x));
    r.setAttribute("y", String(y));
    r.setAttribute("width", String(w));
    r.setAttribute("height", String(h));
    r.setAttribute("fill", fill);
    svg.appendChild(r);
  }
  return svg;
}

function setError(msg) {
  const box = document.getElementById("swarm-error");
  if (!box) return;
  if (!msg) {
    box.textContent = "";
    box.hidden = true;
    return;
  }
  box.textContent = String(msg);
  box.hidden = false;
}

async function loadSwarmMeta() {
  try {
    const res = await fetch("/api/swarm/meta");
    if (!res.ok) return;
    const data = await res.json();
    if (Array.isArray(data.roles) && data.roles.length) {
      swarmMeta.roles = data.roles;
      data.roles.forEach(r => { roleTools[r.slug] = r.ferramentas_tipicas || "CLI"; });
    }
    if (Array.isArray(data.models)) swarmMeta.models = data.models;
    if (Array.isArray(data.workspaces) && data.workspaces.length) swarmMeta.workspaces = data.workspaces;
    if (addPanelOpen) {
      const selRole = document.getElementById("add-worker-role");
      const selModel = document.getElementById("add-worker-model");
      const selWs = document.getElementById("add-worker-ws");
      if (selRole) populateRoleSelect(selRole, selRole.value || "engenheiro-backend");
      if (selModel) populateModelSelect(selModel, selModel.value);
      if (selWs) populateWorkspaceSelect(selWs, selWs.value || "inherit");
    }
    if (window._last) render(window._last);
  } catch (_) {}
}

function populateRoleSelect(sel, currentRole) {
  const cur = String(currentRole || "engenheiro-backend");
  const roles = swarmMeta.roles && swarmMeta.roles.length ? swarmMeta.roles : DEFAULT_SWARM_ROLES;
  const byGuild = new Map();
  roles.forEach(r => {
    const g = String(r.guilda || "geral").toUpperCase();
    if (!byGuild.has(g)) byGuild.set(g, []);
    byGuild.get(g).push(r);
  });
  const nodes = [];
  let matched = false;
  byGuild.forEach((list, guildLabel) => {
    const grp = el("optgroup");
    grp.label = `Guilda ${guildLabel}`;
    list.forEach(r => {
      const opt = el("option", "", `${r.nome || r.slug} (${r.slug})`);
      opt.value = r.slug;
      if (r.slug === cur) {
        opt.selected = true;
        matched = true;
      }
      grp.appendChild(opt);
    });
    nodes.push(grp);
  });
  if (cur && !matched) {
    const custom = el("option", "", cur);
    custom.value = cur;
    custom.selected = true;
    nodes.unshift(custom);
  }
  sel.replaceChildren(...nodes);
  sel.value = cur;
}

function populateModelSelect(sel, currentModel) {
  const cur = String(currentModel || "");
  const models = Array.isArray(swarmMeta.models) ? swarmMeta.models : [];
  const opts = [];
  let matched = false;
  if (!models.length && !cur) {
    const def = el("option", "", "Padrao do sistema (auto)");
    def.value = "";
    opts.push(def);
  }
  models.forEach((m, idx) => {
    const id = typeof m === "string" ? m : (m.id || "");
    if (!id) return;
    const src = typeof m === "object" && m.source ? ` [${m.source}]` : "";
    const opt = el("option", "", `${id}${src}`);
    opt.value = id;
    if ((cur && id === cur) || (!cur && idx === 0)) {
      opt.selected = true;
      matched = true;
    }
    opts.push(opt);
  });
  if (cur && !matched) {
    const custom = el("option", "", `${cur} (atual)`);
    custom.value = cur;
    custom.selected = true;
    opts.unshift(custom);
  }
  sel.replaceChildren(...opts);
  if (cur) sel.value = cur;
  else if (opts.length) sel.value = opts[0].value;
}

function populateWorkspaceSelect(sel, currentWs) {
  const cur = String(currentWs || "inherit");
  const wsList = swarmMeta.workspaces && swarmMeta.workspaces.length ? swarmMeta.workspaces : DEFAULT_WORKSPACES;
  const opts = wsList.map(w => {
    const opt = el("option", "", w.label || w.id);
    opt.value = w.id;
    if (w.id === cur) opt.selected = true;
    return opt;
  });
  sel.replaceChildren(...opts);
  sel.value = cur;
}

function populatePieceKindSelect(sel, curKind) {
  const cur = String(curKind || "");
  const opts = PIECE_KINDS.map(k => {
    const opt = el("option", "", k.label);
    opt.value = k.id;
    if (k.id === cur) opt.selected = true;
    return opt;
  });
  sel.replaceChildren(...opts);
  sel.value = cur;
}

function inferWorkerTicketText(role, pieces, wsMode) {
  const pNames = (pieces || []).map(p => (p.kind ? `${p.kind}:${p.name}` : p.name)).filter(Boolean);
  const mode = wsMode || "inherit";
  if (pNames.length) {
    return `Executar ${role || "worker"} (${mode}) com ${pNames.join(", ")}`;
  }
  return `Executar papel ${role || "worker"} no workspace ${mode}`;
}

function inferWorkerSpawnPreview(role, model, pieces) {
  const pFlag = (pieces && pieces.length) ? ` --pieces ${pieces.map(p => p.name).join(",")}` : "";
  const mFlag = model ? ` --model ${model}` : "";
  return `Auto: board spawn --role ${role || "worker"}${mFlag}${pFlag}`;
}

function renderPieceChips(container, pieces, remAct) {
  if (!container) return;
  if (!pieces || !pieces.length) {
    container.replaceChildren(el("span", "meta", "Nenhuma peca associada ainda."));
    return;
  }
  const chips = pieces.map(p => {
    const chip = el("span", "piece-chip");
    chip.append(document.createTextNode(`${p.kind ? `[${p.kind}] ` : ""}${p.name}`));
    const btn = actEl("button", "piece-chip__x", remAct, `${p.kind || ""}:${p.name}`, "×");
    btn.title = `Remover ${p.name}`;
    chip.appendChild(btn);
    return chip;
  });
  container.replaceChildren(...chips);
}

function addPieceToSelection(list, pieceEl, kindFilterEl) {
  const name = pieceEl ? String(pieceEl.value || "").trim() : "";
  if (!name) return false;
  const found = availablePieces.find(p => p.name === name || p.id === name);
  const kind = found?.kind || (kindFilterEl ? String(kindFilterEl.value || "") : "");
  if (list.some(p => p.name === name && p.kind === kind)) return false;
  list.push({ kind, name });
  return true;
}

async function loadPieceCatalog() {
  try {
    const kinds = ["skill", "script", "mcp", "loop", "gate", "hook"];
    const resList = await Promise.all(kinds.map(k => fetch(`/api/assets/${k}`).then(r => (r.ok ? r.json() : { items: [] })).catch(() => ({ items: [] }))));
    const list = [];
    const seen = new Set();
    kinds.forEach((k, idx) => {
      const items = resList[idx]?.items || [];
      items.forEach(it => {
        const key = `${k}:${it.name}`;
        if (seen.has(key)) return;
        seen.add(key);
        list.push({ id: it.id || it.name, name: it.name, kind: k });
      });
    });
    availablePieces = list;
    if (addPanelOpen) {
      const sel = document.getElementById("add-worker-piece");
      const kindSel = document.getElementById("add-worker-piece-kind");
      if (sel) populatePieceSelect(sel, sel.value, kindSel ? kindSel.value : "");
    }
    if (window._last) render(window._last);
  } catch (_) {}
}

function populatePieceSelect(sel, currentVal, kindFilter) {
  const cur = String(currentVal || "");
  const kf = String(kindFilter || "").toLowerCase();
  const prevFocus = document.activeElement === sel;
  const opts = [];
  const emptyOpt = el("option", "", kf ? `-- Selecionar (${kf}) --` : "-- Selecionar peca do catalogo --");
  emptyOpt.value = "";
  opts.push(emptyOpt);
  let foundCur = !cur;
  availablePieces
    .filter(p => !kf || p.kind === kf)
    .forEach(p => {
      const opt = el("option", "", `[${p.kind}] ${p.name}`);
      opt.value = p.name;
      opt.setAttribute("data-kind", p.kind);
      if (p.name === cur || p.id === cur) {
        opt.selected = true;
        foundCur = true;
      }
      opts.push(opt);
    });
  if (cur && !foundCur) {
    const customOpt = el("option", "", `[custom] ${cur}`);
    customOpt.value = cur;
    customOpt.selected = true;
    opts.push(customOpt);
  }
  sel.replaceChildren(...opts);
  if (prevFocus) sel.focus();
}

function captureCurrentInputs() {
  if (!curProc) return;
  const st = getState(curProc);
  const cmdIn = document.getElementById("cmd-in");
  if (cmdIn && typeof cmdIn.value === "string") {
    st.draftCmd = cmdIn.value;
  }
  const modIn = document.getElementById("model-in");
  if (modIn && typeof modIn.value === "string") {
    st.draftModel = modIn.value;
  }
}

function renderCrumb(ws, st) {
  const crumb = document.getElementById("crumb");
  if (!crumb) return;
  const path = [];
  let p = st.curParent;
  while (p) {
    path.unshift(p);
    const w = ws.find(x => x.name === p);
    p = w ? (w.parent || "") : "";
  }
  const nodes = [actEl("a", "", "sub", "", "Macro")];
  for (const n of path) {
    nodes.push(document.createTextNode(" > "));
    nodes.push(actEl("a", "", "sub", n, n));
  }
  crumb.replaceChildren(...nodes);
}

function renderGroupNodes(title, list, allWs, kind, st) {
  if (!list.length) return [];
  const hdr = el("div", "grp");
  hdr.append(el("span", "", title), el("span", "", String(list.length)));
  const out = [hdr];
  for (const w of list) {
    const subs = allWs.filter(x => x.parent === w.name).length;
    const ci = colorIdx(w.name);
    const isSel = st.selectedAgent === w.name;
    const card = actEl("div", `w worker ${w.status || "waiting"} c${ci}${isSel ? " sel" : ""}`, "sel", w.name);
    card.appendChild(botSVGNode(w.status, w.name, kind));
    const info = el("div", "w-info");
    const row = el("div", "row");
    row.appendChild(el("b", "", w.name));
    const gap = el("span", "row-gap");
    if (subs > 0) {
      gap.appendChild(actEl("button", "sub-btn", "sub", w.name, `${subs} sub`));
    }
    if (w.step) {
      const mCls = w.mode === "async" ? "mode-async" : "mode-sync";
      const mTxt = w.mode === "async" ? `P${w.step} PARALELO` : `P${w.step} SÍNCRONO`;
      gap.appendChild(el("span", `badge ${mCls}`, mTxt));
    }
    gap.appendChild(el("span", `badge ${w.status || "waiting"}`, w.status || "waiting"));
    row.appendChild(gap);
    const metaText = `${w.role || "-"} · ${w.model || "-"} · ${w.workspace_mode || "inherit"}${w.pid ? ` · PID ${w.pid}` : ""}`;
    info.append(row, el("div", "meta", metaText));
    const pieceList = Array.isArray(w.pieces) && w.pieces.length ? w.pieces : (w.piece ? [{ kind: w.piece_kind || "", name: w.piece }] : []);
    if (pieceList.length || w.ticket) {
      const tags = el("div", "w-tags");
      pieceList.forEach(p => tags.appendChild(el("span", "w-tag", `${p.kind ? `${p.kind}: ` : "Peca: "}${p.name}`)));
      if (w.ticket) tags.appendChild(el("span", "w-tag", `Ticket: ${w.ticket}`));
      info.appendChild(tags);
    }
    if (kind !== "coord") {
      const mini = el("div", "w-actions-mini");
      const btnRem = actEl("button", "btn btn--xs btn-danger", "remove", w.name, "Remover");
      btnRem.title = "Remover este worker do processo";
      mini.appendChild(btnRem);
      info.appendChild(mini);
    }
    card.appendChild(info);
    out.push(card);
  }
  return out;
}

function makeLabeledField(labelText, inputEl) {
  const wrap = el("div", "worker-form-field");
  wrap.append(el("label", "", labelText), inputEl);
  return wrap;
}

function syncAddWorkerInference() {
  const role = document.getElementById("add-worker-role")?.value || "engenheiro-backend";
  const model = document.getElementById("add-worker-model")?.value || "";
  const ws = document.getElementById("add-worker-ws")?.value || "inherit";
  const ticketEl = document.getElementById("add-worker-ticket");
  const cmdEl = document.getElementById("add-worker-cmd");
  if (ticketEl && ticketEl.getAttribute("data-manual") !== "1") {
    ticketEl.value = inferWorkerTicketText(role, selectedAddPieces, ws);
  }
  if (cmdEl) {
    cmdEl.placeholder = inferWorkerSpawnPreview(role, model, selectedAddPieces);
  }
}

function renderAddWorkerPanel() {
  const panel = document.getElementById("add-worker-panel");
  if (!panel) return;
  panel.hidden = !addPanelOpen;
  if (!addPanelOpen || panel.children.length > 0) return;

  const grid = el("div", "worker-form-grid");
  const inpName = el("input", "input");
  inpName.id = "add-worker-name";
  inpName.placeholder = "Ex: w-dev-backend";

  const selRole = el("select", "input");
  selRole.id = "add-worker-role";
  populateRoleSelect(selRole, "engenheiro-backend");

  const rowAB = el("div", "worker-form-row2");
  rowAB.append(makeLabeledField("Nome do Worker", inpName), makeLabeledField("Papel Canonico (7 Guildas)", selRole));

  const selModel = el("select", "input");
  selModel.id = "add-worker-model";
  populateModelSelect(selModel, "");

  const selWs = el("select", "input");
  selWs.id = "add-worker-ws";
  populateWorkspaceSelect(selWs, "inherit");

  const rowCD = el("div", "worker-form-row2");
  rowCD.append(makeLabeledField("Modelo do Sistema", selModel), makeLabeledField("Workspace", selWs));

  const selKind = el("select", "input");
  selKind.id = "add-worker-piece-kind";
  populatePieceKindSelect(selKind, "");

  const selPiece = el("select", "input");
  selPiece.id = "add-worker-piece";
  populatePieceSelect(selPiece, "", "");

  const btnAddPiece = actEl("button", "btn btn--sm", "add-piece-item", "", "+ Associar");
  const piecePickerRow = el("div", "piece-picker-row");
  piecePickerRow.append(selKind, selPiece, btnAddPiece);

  const piecesChips = el("div", "piece-chips");
  piecesChips.id = "add-worker-pieces-list";
  renderPieceChips(piecesChips, selectedAddPieces, "rem-piece-item");

  const inpTicket = el("input", "input");
  inpTicket.id = "add-worker-ticket";
  inpTicket.placeholder = "Inferido automaticamente pelos campos acima";

  const inpCmd = el("input", "input");
  inpCmd.id = "add-worker-cmd";

  selRole.addEventListener("change", syncAddWorkerInference);
  selModel.addEventListener("change", syncAddWorkerInference);
  selWs.addEventListener("change", syncAddWorkerInference);
  selKind.addEventListener("change", () => populatePieceSelect(selPiece, "", selKind.value));
  selPiece.addEventListener("change", () => {
    if (addPieceToSelection(selectedAddPieces, selPiece, selKind)) {
      renderPieceChips(piecesChips, selectedAddPieces, "rem-piece-item");
      syncAddWorkerInference();
    }
  });
  inpTicket.addEventListener("input", () => inpTicket.setAttribute("data-manual", inpTicket.value.trim() ? "1" : ""));

  const acts = el("div", "worker-form-actions");
  acts.append(
    actEl("button", "btn btn--sm", "toggle-add", "", "Cancelar"),
    actEl("button", "btn btn--sm btn--primary", "create-worker", "", "Criar Worker")
  );

  grid.append(
    rowAB,
    rowCD,
    makeLabeledField("Pecas Associadas (Tipo + Catalogo)", piecePickerRow),
    piecesChips,
    makeLabeledField("Ticket / Tarefa (Auto-inferido, editavel)", inpTicket),
    makeLabeledField("Comando Spawn (Auto se vazio)", inpCmd),
    acts
  );
  panel.replaceChildren(grid);
  syncAddWorkerInference();
}

function isEditingControlInside(container) {
  const active = document.activeElement;
  return Boolean(active && container.contains(active) && (active.tagName === "INPUT" || active.tagName === "SELECT" || active.tagName === "TEXTAREA"));
}

function renderFicha(allWs, ms, st) {
  const container = document.getElementById("ficha");
  if (!container) return;
  if (!st.selectedAgent) {
    container.replaceChildren();
    return;
  }
  const idx = allWs.findIndex(w => w.name === st.selectedAgent);
  const w = idx >= 0 ? allWs[idx] : { name: st.selectedAgent, role: "orquestrador", model: "maestro", workspace_mode: "inherit", status: "working" };
  const isCoord = idx < 0;
  if (isEditingControlInside(container) && container.dataset.agent === w.name) {
    return;
  }
  container.dataset.agent = w.name;

  const ord = ms.find(m => m.to === "@" + w.name || m.to === "@all");
  const out = ms.find(m => m.from === w.name);
  const subs = allWs.filter(x => x.parent === w.name).map(x => x.name);
  const iso = (w.workspace_mode === "share" || w.workspace_mode === "branch") ? "Isolado (git worktree)" : "Conjunto (inherit)";
  const pCls = (w.status === "killed" || w.status === "done") ? "p-100" : w.status === "waiting" ? "p-15" : "p-70";
  const kind = w.role === "orquestrador" ? "coord" : w.role === "auditor-qualidade" ? "ver" : "trab";

  const ficha = el("div", "ficha");
  ficha.appendChild(botSVGNode(w.status, w.name, kind, false));
  const body = el("div", "ficha-body");

  const topRow = el("div", "row");
  topRow.appendChild(el("b", "", w.name));
  const rightMeta = el("span", "meta");
  const modeInfo = w.step ? `Passo ${w.step} (${w.mode === "async" ? "Paralelo" : "Síncrono"}) · ` : "";
  rightMeta.append(document.createTextNode(`${modeInfo}Duração: ${fmtDur(w.created_at)} · `), el("span", `badge ${w.status || "waiting"}`, w.status || "waiting"));
  topRow.appendChild(rightMeta);
  body.appendChild(topRow);

  body.appendChild(el("div", "meta", `Papel: ${w.role || "-"} · Cérebro: ${w.model || "-"} · Ferramentas: ${roleTools[w.role] || "CLI"} · Modo: ${iso}`));
  const wPieces = Array.isArray(w.pieces) && w.pieces.length ? w.pieces : (w.piece ? [{ kind: w.piece_kind || "", name: w.piece }] : []);
  if (wPieces.length || w.ticket) {
    const pText = wPieces.length ? wPieces.map(p => `${p.kind ? `${p.kind}:` : ""}${p.name}`).join(", ") : "(nenhuma)";
    body.appendChild(el("div", "meta", `Pecas: ${pText} · Ticket: ${w.ticket || "(nenhum)"}`));
  }
  body.appendChild(el("div", "meta", `Input: ${ord ? ord.body : "aguardando"} · Output: ${out ? out.body : "-"}`));
  if (subs.length) {
    body.appendChild(el("div", "meta", `Subprocessos: ${subs.join(", ")}`));
  }
  const prog = el("div", "prog");
  prog.appendChild(el("div", `prog-bar ${pCls}`));
  body.appendChild(prog);

  if (!isCoord) {
    if (selectedEditPiecesAgent !== w.name) {
      selectedEditPiecesAgent = w.name;
      selectedEditPieces = wPieces.map(p => ({ kind: p.kind || "", name: p.name }));
    }
    const editGrid = el("div", "ficha-edit-grid");
    const editRole = el("select", "input");
    editRole.id = "edit-worker-role";
    populateRoleSelect(editRole, w.role || "engenheiro-backend");

    const editPieceKind = el("select", "input");
    editPieceKind.id = "edit-worker-piece-kind";
    populatePieceKindSelect(editPieceKind, "");

    const editPiece = el("select", "input");
    editPiece.id = "edit-worker-piece";
    populatePieceSelect(editPiece, wPieces[0]?.name || "", "");

    const btnAddEditPiece = actEl("button", "btn btn--sm", "edit-add-piece-item", w.name, "+ Peca");
    const editPicker = el("div", "piece-picker-row");
    editPicker.append(editPieceKind, editPiece, btnAddEditPiece);

    const editChips = el("div", "piece-chips");
    editChips.id = "edit-worker-pieces-list";
    renderPieceChips(editChips, selectedEditPieces, "edit-rem-piece-item");

    editPieceKind.addEventListener("change", () => populatePieceSelect(editPiece, "", editPieceKind.value));
    editPiece.addEventListener("change", () => {
      if (addPieceToSelection(selectedEditPieces, editPiece, editPieceKind)) {
        renderPieceChips(editChips, selectedEditPieces, "edit-rem-piece-item");
      }
    });

    const editTicket = el("input", "input");
    editTicket.id = "edit-worker-ticket";
    editTicket.placeholder = "Ticket / tarefa atribuida";
    editTicket.value = w.ticket || "";

    const btnSave = actEl("button", "btn btn--sm btn--primary", "update-worker", w.name, "Salvar Associacao");
    const helpEdit = el("p", "muted", "Associe este worker a um papel canônico, peças do catálogo e ticket/tarefa, depois clique em Salvar Associação.");
    helpEdit.style.gridColumn = "1 / -1";
    helpEdit.style.margin = "0 0 6px 0";
    helpEdit.style.fontSize = "12px";
    editGrid.append(
      helpEdit,
      makeLabeledField("Papel Canonico", editRole),
      makeLabeledField("Pecas (Tipo + Catalogo)", editPicker),
      makeLabeledField("Ticket / Tarefa", editTicket),
      btnSave
    );
    body.append(editGrid, editChips);
  }

  const actions = el("div", "ficha-actions");
  const inpCmd = el("input", "input cmd-in");
  inpCmd.id = "cmd-in";
  inpCmd.setAttribute("data-agent", w.name);
  inpCmd.placeholder = `Enviar instrucao para @${w.name}...`;
  inpCmd.value = st.draftCmd || "";

  const btnPost = actEl("button", "btn btn--sm", "post", w.name, "Enviar");
  btnPost.title = "Enviar instrução/mensagem no mural para este worker";
  const selMod = el("select", "input");
  selMod.id = "model-in";
  selMod.style.width = "180px";
  selMod.setAttribute("data-agent", w.name);
  selMod.title = "Escolha o modelo de IA para este worker";
  populateModelSelect(selMod, st.draftModel || w.model || "");

  const btnModel = actEl("button", "btn btn--sm", "model", w.name, "Trocar Modelo");
  btnModel.title = "Aplicar o modelo de IA selecionado acima";
  actions.append(inpCmd, btnPost, selMod, btnModel);

  if (!isCoord) {
    const bSpawn = actEl("button", "btn btn--sm", "spawn", w.name, "Spawn");
    bSpawn.title = "Iniciar ou reativar este worker";
    const bKill = actEl("button", "btn btn--sm btn-danger", "kill", w.name, "Kill");
    bKill.title = "Encerrar execução deste worker";
    const bRem = actEl("button", "btn btn--sm btn-danger", "remove", w.name, "Remover");
    bRem.title = "Remover worker da lista";
    actions.append(bSpawn, bKill, bRem);
  }
  body.appendChild(actions);
  ficha.appendChild(body);
  container.replaceChildren(ficha);
}

function renderTabs(allOrder, openOrder) {
  const tabsEl = document.getElementById("tabs");
  if (!tabsEl) return;
  const nodes = [];
  for (const p of openOrder) {
    const isSel = p === curProc;
    const btn = actEl("button", `tab${isSel ? " active" : ""}`, "proc", p);
    btn.setAttribute("role", "tab");
    btn.setAttribute("aria-selected", String(isSel));
    btn.appendChild(document.createTextNode(p));
    const closeSpan = actEl("span", "tab-x", "close", p, "×");
    closeSpan.setAttribute("aria-label", `Fechar aba ${p}`);
    btn.appendChild(closeSpan);
    nodes.push(btn);
  }
  tabsEl.replaceChildren(...nodes);
}

function renderMessages(ms, st) {
  const msgsEl = document.getElementById("msgs");
  const pagerEl = document.getElementById("pager");
  if (!msgsEl || !pagerEl) return;

  const fms = ms.filter(m => !st.selectedAgent || m.from === st.selectedAgent || m.to === "@" + st.selectedAgent || m.to === "@all");
  const totalPages = Math.max(1, Math.ceil(fms.length / PER_PAGE));
  if (st.page >= totalPages) st.page = totalPages - 1;
  if (st.page < 0) st.page = 0;

  const slice = fms.slice(st.page * PER_PAGE, (st.page + 1) * PER_PAGE);
  const msgNodes = slice.map(m => {
    const wrap = el("div", "msg");
    const meta = el("div", "meta");
    meta.append(
      document.createTextNode(`#${m.seq} `),
      el("b", "", m.from || ""),
      document.createTextNode(" -> "),
      el("span", "mention", m.to || "@all")
    );
    const body = el("div", "msg-body", m.body || "");
    wrap.append(meta, body);
    return wrap;
  });
  msgsEl.replaceChildren(...msgNodes);

  const prevBtn = actEl("button", "btn btn--sm", "page", st.page - 1, "Anterior");
  prevBtn.disabled = st.page <= 0;
  const nextBtn = actEl("button", "btn btn--sm", "page", st.page + 1, "Próxima");
  nextBtn.disabled = st.page >= totalPages - 1;
  pagerEl.replaceChildren(prevBtn, el("span", "", `Página ${st.page + 1} de ${totalPages}`), nextBtn);
}

function render(payload) {
  if (!payload) return;
  captureCurrentInputs();
  window._last = payload;
  const allOrder = Array.isArray(payload.order) ? payload.order : Object.keys(payload.processes || {});
  const openOrder = allOrder.filter(p => !closedProcs[p]);
  if (!curProc || !(payload.processes && payload.processes[curProc]) || closedProcs[curProc]) {
    curProc = openOrder[0] || "";
  }
  renderTabs(allOrder, openOrder);
  renderAddWorkerPanel();

  const st = getState(curProc);
  const d = (payload.processes && payload.processes[curProc]) || { workers: [], messages: [] };
  const allWs = d.workers || [];
  if (st.selectedAgent && st.selectedAgent !== (st.curParent || "coordenador") && !allWs.some(w => w.name === st.selectedAgent)) {
    st.selectedAgent = "";
  }
  const ms = (d.messages || []).slice().reverse();

  renderCrumb(allWs, st);
  const ws = allWs.filter(w => (w.parent || "") === st.curParent);
  const procDone = d.status === "done" || d.status === "completed" || d.status === "finished" || d.finished === true;
  const normWs = procDone ? ws.map(w => (w.status === "working" ? { ...w, status: "done" } : w)) : ws;
  const ver = normWs.filter(w => w.role === "auditor-qualidade");
  const trab = normWs.filter(w => w.role !== "auditor-qualidade");

  const sumEl = document.getElementById("summary");
  if (sumEl) sumEl.textContent = `1 coordenador · ${ver.length} verificador(es) · ${trab.length} trabalhador(es) · ${ms.length} posts`;
  const countsEl = document.getElementById("counts");
  if (countsEl) countsEl.textContent = `${normWs.length + 1} neste nível`;
  const filterEl = document.getElementById("filter");
  if (filterEl) filterEl.textContent = st.selectedAgent ? `(filtro: ${st.selectedAgent})` : "(todos)";

  const hasWorking = !procDone && normWs.some(w => w.status === "working");
  const allFinished = procDone || (normWs.length > 0 && normWs.every(w => w.status === "done" || w.status === "killed" || w.status === "stopped"));
  const coordStatus = hasWorking ? "working" : allFinished ? "done" : "waiting";
  const coord = [{ name: st.curParent || "coordenador", role: "orquestrador", model: "maestro", workspace_mode: "inherit", status: coordStatus }];
  const workersEl = document.getElementById("workers");
  if (workersEl) {
    workersEl.replaceChildren(
      ...renderGroupNodes("Coordenador", coord, allWs, "coord", st),
      ...renderGroupNodes("Verificadores", ver, allWs, "ver", st),
      ...renderGroupNodes("Trabalhadores", trab, allWs, "trab", st)
    );
  }
  renderFicha(allWs, ms, st);
  renderMessages(ms, st);
}

async function callSwarmAPI(url, bodyObj) {
  setError("");
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(bodyObj)
    });
    if (!res.ok) {
      let detail = `Erro HTTP ${res.status}`;
      try {
        const j = await res.json();
        if (j && j.error) detail = j.error;
      } catch (_) {
        const txt = await res.text();
        if (txt) detail = txt.trim();
      }
      setError(detail);
      return null;
    }
    return res;
  } catch (err) {
    setError(err && err.message ? err.message : "Falha de rede na operação Swarm");
    return null;
  }
}

document.addEventListener("click", async e => {
  const t = e.target && e.target.closest ? e.target.closest("[data-act]") : null;
  if (!t) return;
  const act = t.dataset.act;
  const val = t.dataset.val || "";
  captureCurrentInputs();
  const st = getState(curProc);

  if (act === "close") {
    e.stopPropagation();
    closedProcs[val] = true;
    if (curProc === val) curProc = "";
    render(window._last);
  } else if (act === "reopen") {
    e.stopPropagation();
    delete closedProcs[val];
    curProc = val;
    render(window._last);
  } else if (act === "proc") {
    curProc = val;
    render(window._last);
  } else if (act === "sub") {
    e.stopPropagation();
    st.curParent = val;
    st.selectedAgent = "";
    st.page = 0;
    render(window._last);
  } else if (act === "sel") {
    st.selectedAgent = st.selectedAgent === val ? "" : val;
    st.draftCmd = "";
    st.draftModel = "";
    st.page = 0;
    render(window._last);
  } else if (act === "page") {
    st.page = Number(val);
    render(window._last);
  } else if (act === "toggle-add") {
    addPanelOpen = !addPanelOpen;
    if (!addPanelOpen) {
      selectedAddPieces = [];
      document.getElementById("add-worker-panel")?.replaceChildren();
    } else {
      if (availablePieces.length === 0) loadPieceCatalog();
      if (!swarmMeta.models.length) loadSwarmMeta();
    }
    renderAddWorkerPanel();
  } else if (act === "add-piece-item") {
    const pieceEl = document.getElementById("add-worker-piece");
    const kindEl = document.getElementById("add-worker-piece-kind");
    if (addPieceToSelection(selectedAddPieces, pieceEl, kindEl)) {
      renderPieceChips(document.getElementById("add-worker-pieces-list"), selectedAddPieces, "rem-piece-item");
      syncAddWorkerInference();
    }
  } else if (act === "rem-piece-item") {
    selectedAddPieces = selectedAddPieces.filter(p => `${p.kind || ""}:${p.name}` !== val);
    renderPieceChips(document.getElementById("add-worker-pieces-list"), selectedAddPieces, "rem-piece-item");
    syncAddWorkerInference();
  } else if (act === "edit-add-piece-item") {
    const pieceEl = document.getElementById("edit-worker-piece");
    const kindEl = document.getElementById("edit-worker-piece-kind");
    if (addPieceToSelection(selectedEditPieces, pieceEl, kindEl)) {
      renderPieceChips(document.getElementById("edit-worker-pieces-list"), selectedEditPieces, "edit-rem-piece-item");
    }
  } else if (act === "edit-rem-piece-item") {
    selectedEditPieces = selectedEditPieces.filter(p => `${p.kind || ""}:${p.name}` !== val);
    renderPieceChips(document.getElementById("edit-worker-pieces-list"), selectedEditPieces, "edit-rem-piece-item");
  } else if (act === "create-worker") {
    const nameEl = document.getElementById("add-worker-name");
    const roleEl = document.getElementById("add-worker-role");
    const modelEl = document.getElementById("add-worker-model");
    const wsEl = document.getElementById("add-worker-ws");
    const pieceEl = document.getElementById("add-worker-piece");
    const kindEl = document.getElementById("add-worker-piece-kind");
    const ticketEl = document.getElementById("add-worker-ticket");
    const cmdEl = document.getElementById("add-worker-cmd");
    const name = nameEl ? nameEl.value.trim() : "";
    const role = roleEl ? roleEl.value.trim() : "engenheiro-backend";
    const model = modelEl ? modelEl.value.trim() : "";
    const workspace_mode = wsEl ? wsEl.value.trim() : "inherit";
    addPieceToSelection(selectedAddPieces, pieceEl, kindEl);
    const pieces = selectedAddPieces.slice();
    const piece = pieces[0]?.name || "";
    const piece_kind = pieces[0]?.kind || "";
    const ticket = (ticketEl && ticketEl.value.trim()) ? ticketEl.value.trim() : inferWorkerTicketText(role, pieces, workspace_mode);
    const rawCmd = cmdEl ? cmdEl.value.trim() : "";
    const command = rawCmd ? rawCmd.split(/\s+/) : [];
    const ok = await callSwarmAPI("/api/swarm/add", { process: curProc, name, role, model, workspace_mode, piece, piece_kind, pieces, ticket, parent: st.curParent, command });
    if (ok) {
      if (nameEl) nameEl.value = "";
      if (ticketEl) ticketEl.value = "";
      selectedAddPieces = [];
      addPanelOpen = false;
      document.getElementById("add-worker-panel")?.replaceChildren();
      renderAddWorkerPanel();
      if (typeof window.toast === "function") window.toast("Worker criado com sucesso.");
    }
  } else if (act === "update-worker") {
    const roleEl = document.getElementById("edit-worker-role");
    const modelEl = document.getElementById("model-in");
    const pieceEl = document.getElementById("edit-worker-piece");
    const kindEl = document.getElementById("edit-worker-piece-kind");
    const ticketEl = document.getElementById("edit-worker-ticket");
    const role = roleEl ? roleEl.value.trim() : "";
    const model = modelEl ? modelEl.value.trim() : "";
    addPieceToSelection(selectedEditPieces, pieceEl, kindEl);
    const pieces = selectedEditPieces.slice();
    const piece = pieces[0]?.name || "";
    const piece_kind = pieces[0]?.kind || "";
    const ticket = ticketEl ? ticketEl.value.trim() : "";
    const ok = await callSwarmAPI("/api/swarm/update", { process: curProc, name: val, role, model, piece, piece_kind, pieces, ticket });
    if (ok) {
      selectedEditPiecesAgent = "";
      const f = document.getElementById("ficha");
      if (f) f.setAttribute("data-agent", "");
      if (typeof window.toast === "function") window.toast(`Associacao de "${val}" atualizada.`);
    }
  } else if (act === "post") {
    const inp = document.getElementById("cmd-in");
    if (!inp || !inp.value.trim()) return;
    const body = inp.value.trim();
    const ok = await callSwarmAPI("/api/swarm/post", { proc: curProc, to: "@" + val, body });
    if (ok) {
      inp.value = "";
      st.draftCmd = "";
    }
  } else if (act === "model") {
    const inp = document.getElementById("model-in");
    const model = inp ? inp.value.trim() : "";
    if (!model) {
      setError("Informe o modelo antes de trocar.");
      return;
    }
    await callSwarmAPI("/api/swarm/model", { process: curProc, name: val, model });
  } else if (act === "spawn") {
    await callSwarmAPI("/api/swarm/spawn", { process: curProc, name: val });
  } else if (act === "kill") {
    await callSwarmAPI("/api/swarm/kill", { process: curProc, name: val });
  } else if (act === "remove") {
    e.stopPropagation();
    const ok = await callSwarmAPI("/api/swarm/remove", { process: curProc, name: val });
    if (ok) {
      if (st.selectedAgent === val) st.selectedAgent = "";
      if (typeof window.toast === "function") window.toast(`Worker "${val}" removido.`);
    }
  }
});

loadSwarmMeta();
loadPieceCatalog();
const es = new EventSource("/api/swarm/events");
es.onmessage = e => render(JSON.parse(e.data));
