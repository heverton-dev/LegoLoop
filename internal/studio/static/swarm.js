"use strict";

const PER_PAGE = 20;
const palette = ["#8ab4f8", "#f28b82", "#fdd663", "#81c995", "#c58af9"];
const roleTools = {
  orquestrador: "LegoLoop CLI",
  planejador: "PLAN v2",
  arquiteto: "Templates v2",
  "engenheiro-backend": "Go / REST",
  "engenheiro-tdd": "Go Test / Vitest",
  "engenheiro-frontend": "React / Tailwind",
  "engenheiro-secops": "Security Gates",
  "dissecador-web": "REA Web",
  "auditor-qualidade": "Taxlinter / Funclen",
  "redator-tecnico": "Manual Creator"
};

let curProc = "";
let closedProcs = {};
let addPanelOpen = false;
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
    nodes.push(document.createTextNode(" › "));
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
    const card = actEl("div", `w ${w.status || "waiting"} c${ci}${isSel ? " sel" : ""}`, "sel", w.name);
    card.appendChild(botSVGNode(w.status, w.name, kind));
    const info = el("div", "w-info");
    const row = el("div", "row");
    row.appendChild(el("b", "", w.name));
    const gap = el("span", "row-gap");
    if (subs > 0) {
      gap.appendChild(actEl("button", "sub-btn", "sub", w.name, `${subs} sub ▶`));
    }
    gap.appendChild(el("span", `badge ${w.status || "waiting"}`, w.status || "waiting"));
    row.appendChild(gap);
    const metaText = `${w.role || "-"} · ${w.model || "-"} · ${w.workspace_mode || "inherit"}${w.pid ? ` · PID ${w.pid}` : ""}`;
    info.append(row, el("div", "meta", metaText));
    card.appendChild(info);
    out.push(card);
  }
  return out;
}

function renderAddWorkerPanel() {
  const panel = document.getElementById("add-worker-panel");
  if (!panel) return;
  panel.hidden = !addPanelOpen;
  if (!addPanelOpen || panel.children.length > 0) return;
  const body = el("div", "ficha-body");
  const row1 = el("div", "ficha-actions");
  const inpName = el("input", "btn cmd-in");
  inpName.id = "add-worker-name";
  inpName.placeholder = "Nome do worker (ex: w-dev)";
  const selRole = el("select", "btn");
  selRole.id = "add-worker-role";
  for (const r of Object.keys(roleTools)) {
    const opt = el("option", "", r);
    opt.value = r;
    selRole.appendChild(opt);
  }
  const inpModel = el("input", "btn");
  inpModel.id = "add-worker-model";
  inpModel.placeholder = "Modelo (ex: t-codigo-128k)";
  inpModel.value = "t-codigo-128k";
  const selWs = el("select", "btn");
  selWs.id = "add-worker-ws";
  for (const mode of ["inherit", "share", "branch"]) {
    const opt = el("option", "", mode);
    opt.value = mode;
    selWs.appendChild(opt);
  }
  row1.append(inpName, selRole, inpModel, selWs);
  const row2 = el("div", "ficha-actions");
  const inpCmd = el("input", "btn cmd-in");
  inpCmd.id = "add-worker-cmd";
  inpCmd.placeholder = "Comando opcional para spawn imediato (ex: cmd /c echo ok)";
  const btnCreate = actEl("button", "btn", "create-worker", "", "Criar Worker");
  row2.append(inpCmd, btnCreate);
  body.append(row1, row2);
  panel.replaceChildren(body);
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
  const ord = ms.find(m => m.to === "@" + w.name || m.to === "@all");
  const out = ms.find(m => m.from === w.name);
  const subs = allWs.filter(x => x.parent === w.name).map(x => x.name);
  const iso = (w.workspace_mode === "share" || w.workspace_mode === "branch") ? "Isolado (git worktree)" : "Conjunto (inherit)";
  const pCls = w.status === "killed" ? "p-100" : w.status === "waiting" ? "p-15" : "p-70";
  const kind = w.role === "orquestrador" ? "coord" : w.role === "auditor-qualidade" ? "ver" : "trab";

  const existingCmd = document.getElementById("cmd-in");
  const existingModel = document.getElementById("model-in");
  const cmdVal = existingCmd && existingCmd.dataset && existingCmd.dataset.agent === w.name ? existingCmd.value : (st.draftCmd || "");
  const modelVal = existingModel && existingModel.dataset && existingModel.dataset.agent === w.name ? existingModel.value : (st.draftModel || w.model || "");

  const ficha = el("div", "ficha");
  ficha.appendChild(botSVGNode(w.status, w.name, kind, false));
  const body = el("div", "ficha-body");

  const topRow = el("div", "row");
  topRow.appendChild(el("b", "", w.name));
  const rightMeta = el("span", "meta");
  rightMeta.append(document.createTextNode(`⏱ ${fmtDur(w.created_at)} · `), el("span", `badge ${w.status || "waiting"}`, w.status || "waiting"));
  topRow.appendChild(rightMeta);
  body.appendChild(topRow);

  body.appendChild(el("div", "meta", `Papel: ${w.role || "-"} · Cérebro: ${w.model || "-"} · Ferramentas: ${roleTools[w.role] || "CLI"} · Modo: ${iso}`));
  body.appendChild(el("div", "meta", `Input: ${ord ? ord.body : "aguardando"} · Output: ${out ? out.body : "-"}`));
  if (subs.length) {
    body.appendChild(el("div", "meta", `Subprocessos: ${subs.join(", ")}`));
  }
  const prog = el("div", "prog");
  prog.appendChild(el("div", `prog-bar ${pCls}`));
  body.appendChild(prog);

  const actions = el("div", "ficha-actions");
  const inpCmd = el("input", "btn cmd-in");
  inpCmd.id = "cmd-in";
  inpCmd.setAttribute("data-agent", w.name);
  inpCmd.placeholder = `Enviar instrução para @${w.name}...`;
  inpCmd.value = cmdVal;
  st.draftCmd = cmdVal;

  const btnPost = actEl("button", "btn", "post", w.name, "Enviar");
  const inpModel = el("input", "btn");
  inpModel.id = "model-in";
  inpModel.setAttribute("data-agent", w.name);
  inpModel.placeholder = "Novo modelo";
  inpModel.value = modelVal;
  st.draftModel = modelVal;

  const btnModel = actEl("button", "btn", "model", w.name, "Trocar Modelo");
  const btnSpawn = actEl("button", "btn", "spawn", w.name, "Spawn");
  const btnKill = actEl("button", "btn btn-danger", "kill", w.name, "Kill");

  actions.append(inpCmd, btnPost, inpModel, btnModel, btnSpawn, btnKill);
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
  const closedList = allOrder.filter(p => closedProcs[p]);
  for (const p of closedList) {
    nodes.push(actEl("button", "tab", "reopen", p, `+ ${p}`));
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
      document.createTextNode(" → "),
      el("span", "mention", m.to || "@all")
    );
    const body = el("div", "msg-body", m.body || "");
    wrap.append(meta, body);
    return wrap;
  });
  msgsEl.replaceChildren(...msgNodes);

  const prevBtn = actEl("button", "btn", "page", st.page - 1, "Anterior");
  prevBtn.disabled = st.page <= 0;
  const nextBtn = actEl("button", "btn", "page", st.page + 1, "Próxima");
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
  const ms = (d.messages || []).slice().reverse();

  renderCrumb(allWs, st);
  const ws = allWs.filter(w => (w.parent || "") === st.curParent);
  const ver = ws.filter(w => w.role === "auditor-qualidade");
  const trab = ws.filter(w => w.role !== "auditor-qualidade");

  const sumEl = document.getElementById("summary");
  if (sumEl) sumEl.textContent = `1 coordenador · ${ver.length} verificador(es) · ${trab.length} trabalhador(es) · ${ms.length} posts`;
  const countsEl = document.getElementById("counts");
  if (countsEl) countsEl.textContent = `${ws.length + 1} neste nível`;
  const filterEl = document.getElementById("filter");
  if (filterEl) filterEl.textContent = st.selectedAgent ? `(filtro: ${st.selectedAgent})` : "(todos)";

  const coord = [{ name: st.curParent || "coordenador", role: "orquestrador", model: "maestro", workspace_mode: "inherit", status: "working" }];
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
    renderAddWorkerPanel();
  } else if (act === "create-worker") {
    const nameEl = document.getElementById("add-worker-name");
    const roleEl = document.getElementById("add-worker-role");
    const modelEl = document.getElementById("add-worker-model");
    const wsEl = document.getElementById("add-worker-ws");
    const cmdEl = document.getElementById("add-worker-cmd");
    const name = nameEl ? nameEl.value.trim() : "";
    const role = roleEl ? roleEl.value.trim() : "engenheiro-backend";
    const model = modelEl ? modelEl.value.trim() : "";
    const workspace_mode = wsEl ? wsEl.value.trim() : "inherit";
    const rawCmd = cmdEl ? cmdEl.value.trim() : "";
    const command = rawCmd ? rawCmd.split(/\s+/) : [];
    const ok = await callSwarmAPI("/api/swarm/add", { process: curProc, name, role, model, workspace_mode, parent: st.curParent, command });
    if (ok && nameEl) nameEl.value = "";
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
    await callSwarmAPI("/api/swarm/kill", { proc: curProc, name: val });
  }
});

const es = new EventSource("/api/swarm/events");
es.onmessage = e => render(JSON.parse(e.data));
