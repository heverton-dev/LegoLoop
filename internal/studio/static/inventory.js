"use strict";

// Página de inventário do LegoLoop Studio com Agrupamento Semântico, Paginação e Modal Rico.
// Todo texto vindo da API entra no DOM via textContent (anti-XSS) e todo POST vai como JSON.

const TABS = [
  { id: "skill", label: "Skills", empty: "Nenhuma skill encontrada nas pastas conhecidas dos harnesses." },
  { id: "mcp", label: "MCPs", empty: "Nenhum servidor MCP configurado." },
  { id: "script", label: "Scripts", empty: "Nenhum script encontrado nas pastas conhecidas." },
  { id: "hook", label: "Hooks", empty: "Nenhum hook configurado (global, projeto ou local)." },
  { id: "gate", label: "Gates", empty: "Nenhum gate." },
  { id: "model", label: "Modelos", empty: "Nenhum modelo executável encontrado." },
  { id: "harness", label: "Harnesses", empty: "Nenhum harness encontrado." },
  { id: "binary", label: "Binários", empty: "Nenhum binário conhecido no PATH." },
  { id: "loop", label: "Loops", empty: "Nenhum loop salvo no acervo ainda." },
];

const POLL_MS = 10000;
const SCOPE_LABEL = { global: "global", project: "projeto", local: "local", builtin: "embutido", system: "sistema" };

const st = {
  counts: {},
  version: 0,
  kind: "skill",
  items: [],
  total: 0,
  query: "",
  groupBySha: true,
  page: 1,
  pageSize: 25,
  plan: null,
  activeItem: null,
};

const $ = (id) => document.getElementById(id);

function el(tag, cls, text) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (text !== undefined) e.textContent = text;
  return e;
}

async function api(path, body) {
  const opts = body === undefined ? {} : {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  };
  const res = await fetch(path, opts);
  const data = await res.json().catch(() => ({}));
  return { ok: res.ok, status: res.status, data };
}

function humanBytes(n) {
  if (n < 1024) return `${n} B`;
  const units = ["KiB", "MiB", "GiB"];
  let v = n;
  let i = -1;
  do { v /= 1024; i += 1; } while (v >= 1024 && i < units.length - 1);
  return `${v.toFixed(1)} ${units[i]}`;
}

// ---------- Abas de Categorias ----------

function renderTabs() {
  const buttons = TABS.map((t, idx) => {
    const b = el("button", "tab");
    b.type = "button";
    b.id = `inv-tab-${t.id}`;
    b.setAttribute("role", "tab");
    const active = t.id === st.kind;
    b.setAttribute("aria-selected", String(active));
    b.setAttribute("aria-controls", "list");
    b.tabIndex = active ? 0 : -1;
    b.append(el("span", "tab__label", t.label), el("span", "tab__count", String(st.counts[t.id] ?? 0)));
    b.addEventListener("click", () => selectKind(t.id));
    b.addEventListener("keydown", (e) => {
      let nextIdx = -1;
      if (e.key === "ArrowDown" || e.key === "ArrowRight") nextIdx = (idx + 1) % TABS.length;
      else if (e.key === "ArrowUp" || e.key === "ArrowLeft") nextIdx = (idx - 1 + TABS.length) % TABS.length;
      else if (e.key === "Home") nextIdx = 0;
      else if (e.key === "End") nextIdx = TABS.length - 1;
      if (nextIdx >= 0) {
        e.preventDefault();
        selectKind(TABS[nextIdx].id).then(() => {
          document.getElementById(`inv-tab-${TABS[nextIdx].id}`)?.focus();
        });
      }
    });
    return b;
  });
  $("tabs").replaceChildren(...buttons);
  const activeTab = TABS.find((t) => t.id === st.kind);
  if (activeTab) $("list")?.setAttribute("aria-labelledby", `inv-tab-${activeTab.id}`);
}

async function selectKind(kind) {
  st.kind = kind;
  st.page = 1;
  history.replaceState(null, "", `#${kind}`);
  renderTabs();
  await loadItems();
}

async function loadItems() {
  const { ok, data } = await api(`/api/assets/${encodeURIComponent(st.kind)}`);
  if (!ok) {
    toast(data.error || "Não foi possível listar esta categoria.", true);
    return;
  }
  st.items = data.items || [];
  st.total = data.total || st.items.length;
  renderList();
}

// ---------- Agrupamento e Filtragem ----------

function groupItems(items) {
  if (!st.groupBySha) return items;
  const groups = new Map();
  for (const it of items) {
    const key = it.meta?.signature || (it.kind + ":" + it.name);
    if (!groups.has(key)) {
      groups.set(key, { ...it, occurrences: [it], origins: new Set([it.origin || "desconhecido"]) });
    } else {
      const g = groups.get(key);
      g.occurrences.push(it);
      if (it.origin) g.origins.add(it.origin);
    }
  }
  return Array.from(groups.values()).map((g) => {
    g.originsList = Array.from(g.origins);
    g.copyCount = g.occurrences.length;
    return g;
  });
}

function filtered() {
  let list = st.items;
  if (st.query.trim()) {
    list = list.filter((it) => matchesQuery(`${it.name} ${it.description} ${it.origin} ${it.scope} ${it.path}`, st.query));
  }
  return groupItems(list);
}

function badge(text, cls) {
  return el("span", `badge ${cls || ""}`.trim(), text);
}

// ---------- Renderização da Linha de Item ----------

function itemRow(it) {
  const row = el("div", "inv-item");
  row.setAttribute("role", "button");
  row.tabIndex = 0;

  const head = el("div", "inv-item__head");
  head.append(el("span", "inv-item__name", it.name));
  if (it.scope) head.append(badge(SCOPE_LABEL[it.scope] || it.scope));

  if (it.originsList && it.originsList.length > 0) {
    it.originsList.forEach((o) => head.append(badge(o, "badge--muted")));
  } else if (it.origin) {
    head.append(badge(it.origin, "badge--muted"));
  }

  const copies = it.copyCount || it.meta?.duplicates || (it.occurrences ? it.occurrences.length : 1);
  if (copies > 1) {
    head.append(badge(`${copies} cópias`, "badge--warn"));
  }

  row.append(head);
  if (it.description) row.append(el("p", "inv-item__desc", it.description));
  const meta = [it.meta?.signature, it.meta?.transport, it.meta?.command && `cmd: ${it.meta.command}`, it.meta?.host].filter(Boolean);
  if (meta.length) row.append(el("p", "inv-item__meta", meta.join(" · ")));
  if (it.path) row.append(el("p", "inv-item__path", it.path));

  row.addEventListener("click", () => openItemModal(it));
  row.addEventListener("keydown", (e) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openItemModal(it);
    }
  });

  return row;
}

// ---------- Paginação e Lista ----------

function renderPagination(totalItems) {
  const pag = $("pagination");
  const totalPages = Math.max(1, Math.ceil(totalItems / st.pageSize));
  if (st.page > totalPages) st.page = totalPages;

  pag.replaceChildren();
  if (totalItems <= 0) return;

  const info = el("span", "pagination__info", `Exibindo ${(st.page - 1) * st.pageSize + 1}–${Math.min(st.page * st.pageSize, totalItems)} de ${totalItems} peças`);
  const btns = el("div", "pagination__buttons");

  const btnFirst = el("button", "btn btn--sm", "«");
  btnFirst.disabled = st.page <= 1;
  btnFirst.addEventListener("click", () => { st.page = 1; renderList(); });

  const btnPrev = el("button", "btn btn--sm", "‹");
  btnPrev.disabled = st.page <= 1;
  btnPrev.addEventListener("click", () => { st.page--; renderList(); });

  const pageLabel = el("span", "pagination__label", `Pág. ${st.page} de ${totalPages}`);

  const btnNext = el("button", "btn btn--sm", "›");
  btnNext.disabled = st.page >= totalPages;
  btnNext.addEventListener("click", () => { st.page++; renderList(); });

  const btnLast = el("button", "btn btn--sm", "»");
  btnLast.disabled = st.page >= totalPages;
  btnLast.addEventListener("click", () => { st.page = totalPages; renderList(); });

  btns.append(btnFirst, btnPrev, pageLabel, btnNext, btnLast);
  pag.append(info, btns);
}

function renderList() {
  const tab = TABS.find((t) => t.id === st.kind);
  const all = filtered();
  const start = (st.page - 1) * st.pageSize;
  const visible = all.slice(start, start + st.pageSize);

  $("list").replaceChildren(...visible.map(itemRow));
  $("result-count").textContent = `${all.length} peça(s)`;

  const empty = all.length === 0;
  $("empty").hidden = !empty;
  if (empty) {
    const searching = st.query.trim() !== "" && st.items.length > 0;
    $("empty-title").textContent = searching ? `Nada encontrado para "${st.query.trim()}"` : tab.label;
    $("empty-text").textContent = searching ? "Tente outro termo ou troque de categoria." : tab.empty;
  }

  renderPagination(all.length);
}

// ---------- Modal Rico de Peça (Delegado ao Web Component Canônico) ----------

async function openItemModal(it) {
  st.activeItem = it;
  const modal = $("modal-item");
  if (!modal) return;
  modal.open(it, {
    showCanvasAction: false,
    onCanonize: async () => {
      await refreshSummary(true);
    },
    onUpgrade: async () => {
      await refreshSummary(true);
    },
  });
}

// ---------- Atualização e Canonização ----------

function showScanned(at) {
  $("scanned").textContent = at ? `Varredura: ${new Date(at).toLocaleTimeString("pt-BR")}` : "";
}

async function refreshSummary(force = false) {
  const { ok, data } = force ? await api("/api/assets/rescan", {}) : await api("/api/assets");
  if (!ok) return toast(data.error || "Falha ao atualizar o inventário.", true);
  const changed = data.version !== st.version;
  st.counts = data.counts;
  st.version = data.version;
  showScanned(data.scanned_at);
  renderTabs();
  if (changed || force) await loadItems();
}

function showLegoLoading(label = "Varrendo...") {
  let ov = document.getElementById("lego-loading-overlay");
  if (!ov) {
    ov = el("div", "lego-loading-overlay");
    ov.id = "lego-loading-overlay";
    const box = el("div", "lego-loading-box");
    const stack = el("div", "lego-loading-stack");
    ["lego-brick--a", "lego-brick--b", "lego-brick--c"].forEach((c) => {
      const b = el("div", `lego-brick ${c}`);
      b.append(el("span", "lego-stud"), el("span", "lego-stud"));
      stack.append(b);
    });
    const txt = el("div", "lego-loading-text", label);
    box.append(stack, txt);
    ov.append(box);
    document.body.append(ov);
  }
  ov.querySelector(".lego-loading-text").textContent = label;
  ov.hidden = false;
}

function hideLegoLoading() {
  const ov = document.getElementById("lego-loading-overlay");
  if (ov) ov.hidden = true;
}

async function rescan(btn) {
  btn.disabled = true;
  showLegoLoading("Varrendo...");
  try {
    await refreshSummary(true);
    toast("Inventário atualizado.");
  } finally {
    hideLegoLoading();
    btn.disabled = false;
  }
}

function line(text, cls) {
  return el("p", cls || "canon__line", text);
}

function renderPlan(plan) {
  const host = $("canon-result");
  host.replaceChildren();
  if (!plan.groups || plan.groups.length === 0) {
    host.append(line("Nenhuma duplicata exata para canonizar."));
    $("btn-apply").disabled = true;
    return;
  }
  let totalLinks = 0;
  for (const g of plan.groups) totalLinks += g.links.length;
  host.append(
    line(`Cofre: ${plan.vault}`),
    line(`${plan.groups.length} grupos (${totalLinks} links a criar)`),
    line(`Economia estimada: ${humanBytes(plan.saved_bytes)}`, "canon__ok")
  );
  $("btn-apply").disabled = false;
}

async function simulate(btn) {
  btn.disabled = true;
  try {
    const { ok, data } = await api("/api/canonize/plan", {});
    if (!ok) return toast(data.error || "Falha ao planejar.", true);
    st.plan = data;
    renderPlan(data);
  } finally {
    btn.disabled = false;
  }
}

async function apply(btn) {
  if (!st.plan) return;
  const sure = await confirmDialog({
    title: "Aplicar Canonização",
    message: "Substituir duplicatas por links no cofre canônico? O original tem backup e nada é apagado.",
    confirmLabel: "Aplicar",
    danger: true,
  });
  if (!sure) return;
  btn.disabled = true;
  try {
    const { ok, data } = await api("/api/canonize/apply", { apply: true });
    if (!ok) return toast(data.error || "Falha ao aplicar.", true);
    const methods = Object.entries(data.methods).map(([m, n]) => `${n} ${m}`).join(", ") || "nenhum link";
    toast(`${data.linked} arquivo(s) canonizados (${methods}); economia ${humanBytes(data.saved_bytes)}.`, data.failed.length > 0);
    st.plan = null;
    await refreshSummary(true);
  } finally {
    btn.disabled = st.plan === null;
  }
}

async function guarded(btn, fn) {
  try {
    await fn(btn);
  } catch (err) {
    toast(`Falha de comunicação com o servidor local: ${err.message}`, true);
  }
}

// ---------- Boot ----------

async function boot() {
  const wanted = location.hash.slice(1);
  if (TABS.some((t) => t.id === wanted)) st.kind = wanted;

  $("search").addEventListener("input", (e) => { st.query = e.target.value; st.page = 1; renderList(); });
  $("search").addEventListener("keydown", (e) => {
    if (e.key === "Escape") { e.target.value = ""; st.query = ""; st.page = 1; renderList(); }
  });
  $("chk-group").addEventListener("change", (e) => {
    st.groupBySha = e.target.checked;
    st.page = 1;
    renderList();
  });
  $("sel-pagesize").addEventListener("change", (e) => {
    st.pageSize = parseInt(e.target.value, 10) || 25;
    st.page = 1;
    renderList();
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "/" && !e.target.closest("input, select, textarea")) { e.preventDefault(); $("search").focus(); }
  });
  window.addEventListener("hashchange", () => {
    const id = location.hash.slice(1);
    if (id !== st.kind && TABS.some((t) => t.id === id)) selectKind(id);
  });

  $("btn-rescan").addEventListener("click", (e) => guarded(e.currentTarget, rescan));
  $("btn-simulate").addEventListener("click", (e) => guarded(e.currentTarget, simulate));
  $("btn-apply").addEventListener("click", (e) => guarded(e.currentTarget, apply));

  await refreshSummary();
  setInterval(() => { if (document.visibilityState === "visible") refreshSummary().catch(() => {}); }, POLL_MS);
}

boot().catch((err) => toast(`Não foi possível iniciar o inventário: ${err.message}`, true));
