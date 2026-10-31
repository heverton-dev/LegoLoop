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
  $("tabs").replaceChildren(...TABS.map((t) => {
    const b = el("button", "tab");
    b.type = "button";
    b.setAttribute("role", "tab");
    b.setAttribute("aria-selected", String(t.id === st.kind));
    b.append(el("span", "tab__label", t.label), el("span", "tab__count", String(st.counts[t.id] ?? 0)));
    b.addEventListener("click", () => selectKind(t.id));
    return b;
  }));
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

// ---------- Modal Rico de Peça ----------

async function openItemModal(it) {
  st.activeItem = it;
  $("modal-item-name").textContent = it.name;
  
  const bList = $("modal-item-badges");
  bList.replaceChildren(badge(it.kind, "badge--accent"));
  if (it.scope) bList.append(badge(SCOPE_LABEL[it.scope] || it.scope));
  if (it.originsList) {
    it.originsList.forEach((o) => bList.append(badge(o, "badge--muted")));
  } else if (it.origin) {
    bList.append(badge(it.origin, "badge--muted"));
  }
  const copies = it.copyCount || (it.occurrences ? it.occurrences.length : 1);
  if (copies > 1) bList.append(badge(`${copies} cópias no SO`, "badge--warn"));

  // Top Actions: Canonizar Peça e Abrir no Explorer
  const btnCanon = $("btn-modal-canonize");
  if (btnCanon) {
    btnCanon.disabled = copies <= 1;
    btnCanon.onclick = async () => {
      btnCanon.disabled = true;
      btnCanon.textContent = "Canonizando...";
      const { ok, data } = await api("/api/canonize/piece", { sha256: it.sha256 });
      btnCanon.textContent = "Canonizar Peça";
      if (ok) {
        toast(`Peça canonizada! ${data.linked || 0} links criados no cofre.`);
        await refreshSummary(true);
      } else {
        btnCanon.disabled = false;
        toast(data.error || "Falha na canonização.", true);
      }
    };
  }

  const btnReveal = $("btn-modal-reveal");
  if (btnReveal) {
    btnReveal.disabled = !it.path;
    btnReveal.onclick = async () => {
      const { ok } = await api("/api/assets/reveal", { path: it.path });
      if (ok) toast("Pasta aberta no Windows Explorer.");
      else toast("Não foi possível abrir o Explorer.", true);
    };
  }

  // Aba 1: Diagrama IPO e Visão Geral
  const ipoIn = $("modal-ipo-in");
  ipoIn.replaceChildren();
  (it.meta?.consumes || ["Briefing / Requisitos"]).forEach(c => ipoIn.append(badge(String(c).replace("Artifact:", ""), "badge--accent")));

  const ipoProc = $("modal-ipo-proc");
  ipoProc.replaceChildren(el("span", "ipo-desc", it.description || it.meta?.summary || "Execução autônoma de modelo/código"));

  const ipoOut = $("modal-ipo-out");
  ipoOut.replaceChildren();
  (it.meta?.produces ? [it.meta.produces] : ["Artifact:Report"]).forEach(p => ipoOut.append(badge(String(p).replace("Artifact:", ""), "badge--accent")));

  const artContainer = $("modal-prop-artifacts");
  artContainer.replaceChildren(badge(it.meta?.produces || "Entrega Verificada", "badge--accent"));
  if (it.meta?.artifact) artContainer.append(badge(it.meta.artifact, "badge--muted"));

  $("modal-prop-what").textContent = it.description ? `Peça do tipo ${it.kind.toUpperCase()}: ${it.description}` : `Peça de software autônoma e executável: ${it.name}.`;
  $("modal-prop-does").textContent = it.meta?.summary || "Executa tarefas técnicas determinísticas com isolamento e contrato de entrada e saída.";
  $("modal-prop-delivers").textContent = (it.meta?.produces ? `Produz artefato '${it.meta.produces}'` : "Entrega artefato de saída imutável e relatório de execução validado.");

  // Aba 2: Peças Compatíveis
  loadCompatibleItems(it);

  // Aba 3: Ocorrências no SO com botão Abrir Pasta
  const occList = $("modal-occurrences-list");
  occList.replaceChildren();
  const occs = it.occurrences || [it];
  occs.forEach((o) => {
    const itemEl = el("div", "occurrence-item");
    itemEl.append(el("span", "occurrence-path", o.path || "(sem caminho físico)"));
    const metaEl = el("span", "occurrence-meta", `${o.scope || "global"} · ${o.origin || "sistema"}`);
    const btnOpenFolder = el("button", "btn btn--sm occ-btn-open", "📁 Abrir Pasta");
    btnOpenFolder.onclick = async () => {
      const { ok } = await api("/api/assets/reveal", { path: o.path });
      if (ok) toast("Pasta aberta no Windows Explorer.");
      else toast("Falha ao abrir Explorer.", true);
    };
    itemEl.append(metaEl, btnOpenFolder);
    occList.append(itemEl);
  });

  // Aba 4: Código / Fonte
  $("modal-source-path").textContent = it.path || "";
  $("modal-source-code").textContent = "Carregando conteúdo da peça...";
  $("modal-translation-status").hidden = true;
  loadSourceContent(it);

  // Abre na aba de Visão Geral
  switchModalTab("overview");
  $("modal-item").showModal();
}

function switchModalTab(tabName) {
  document.querySelectorAll(".modal-tab").forEach((btn) => {
    btn.classList.toggle("is-active", btn.dataset.tab === tabName);
  });
  document.querySelectorAll(".modal-panel").forEach((panel) => {
    panel.classList.toggle("is-active", panel.id === `tab-${tabName}`);
  });
}

let compatState = { all: [], filtered: [], cat: "all", page: 1, pageSize: 6 };

function renderCompatPage() {
  const container = $("modal-compat-list");
  const pag = $("modal-compat-pagination");
  const list = compatState.cat === "all" ? compatState.all : compatState.all.filter(c => c.item.kind === compatState.cat || (compatState.cat === "skill" && c.item.kind === "block"));
  
  const totalPages = Math.max(1, Math.ceil(list.length / compatState.pageSize));
  if (compatState.page > totalPages) compatState.page = totalPages;
  const start = (compatState.page - 1) * compatState.pageSize;
  const visible = list.slice(start, start + compatState.pageSize);

  if (list.length === 0) {
    container.replaceChildren(el("p", "panel__hint", "Nenhuma peça compatível nessa categoria."));
    pag.replaceChildren();
    return;
  }

  container.replaceChildren(...visible.map((c) => {
    const card = el("div", "compat-item");
    const top = el("div", "compat-item__top");
    top.append(el("strong", "compat-item__name", c.item.name), badge(c.role, "badge--accent"), badge(c.item.kind, "badge--muted"));
    const reason = el("p", "compat-item__reason", c.reason);
    card.append(top, reason);
    return card;
  }));

  pag.replaceChildren();
  if (totalPages > 1) {
    const info = el("span", "pagination__info", `Pág. ${compatState.page} de ${totalPages} (${list.length} peças)`);
    const btnPrev = el("button", "btn btn--sm", "‹");
    btnPrev.disabled = compatState.page <= 1;
    btnPrev.addEventListener("click", () => { compatState.page--; renderCompatPage(); });
    const btnNext = el("button", "btn btn--sm", "›");
    btnNext.disabled = compatState.page >= totalPages;
    btnNext.addEventListener("click", () => { compatState.page++; renderCompatPage(); });
    pag.append(info, btnPrev, btnNext);
  }
}

async function loadCompatibleItems(it) {
  const container = $("modal-compat-list");
  container.replaceChildren(el("p", "panel__hint", "Buscando conexões ontológicas no grafo..."));
  $("modal-compat-pagination").replaceChildren();

  const res = await fetch(`/api/assets/compatible?id=${encodeURIComponent(it.id || it.path || it.name)}`);
  const data = await res.json().catch(() => ({}));
  compatState.all = data.compatible || [];
  compatState.page = 1;
  compatState.cat = "all";

  document.querySelectorAll(".compat-filter").forEach(b => {
    b.classList.toggle("is-active", b.dataset.cat === "all");
    b.onclick = () => {
      document.querySelectorAll(".compat-filter").forEach(x => x.classList.remove("is-active"));
      b.classList.add("is-active");
      compatState.cat = b.dataset.cat;
      compatState.page = 1;
      renderCompatPage();
    };
  });

  renderCompatPage();
}

async function loadSourceContent(it) {
  if (!it.path) {
    $("modal-source-code").textContent = "(peça virtual ou modelo embutido sem arquivo no disco)";
    return;
  }
  const { ok, data } = await api("/api/assets/translate", { path: it.path });
  if (ok && data.full_content_ptbr) {
    $("modal-source-code").textContent = data.full_content_ptbr;
    st.activeTranslation = data;
  } else {
    $("modal-source-code").textContent = "Não foi possível carregar o arquivo ou caminho inválido.";
  }
}

async function translateActiveItem() {
  const it = st.activeItem;
  if (!it || !it.path) return;
  const btn = $("btn-translate");
  btn.disabled = true;
  btn.textContent = "Traduzindo...";

  const { ok, data } = await api("/api/assets/translate", { path: it.path });
  btn.disabled = false;
  btn.textContent = "Traduzir para PT-BR";

  if (!ok) {
    toast(data.error || "Falha na tradução.", true);
    return;
  }

  $("modal-source-code").textContent = data.full_content_ptbr;
  $("modal-prop-what").textContent = data.what_it_is;
  $("modal-prop-does").textContent = data.what_it_does;
  $("modal-prop-delivers").textContent = data.what_it_delivers;

  const stBox = $("modal-translation-status");
  stBox.hidden = false;
  stBox.textContent = `✓ Tradução didática ativa (Arquivo fonte mantido 100% intacto · Hash: ${data.sha256.slice(0, 8)})`;
  toast("Tradução didática aplicada com sucesso!");
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

async function rescan(btn) {
  btn.disabled = true;
  try {
    await refreshSummary(true);
    toast("Inventário atualizado.");
  } finally {
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
    text: "Substituir duplicatas por links no cofre canônico? O original tem backup e nada é apagado.",
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

  // Modal subtabs
  document.querySelectorAll(".modal-tab").forEach((btn) => {
    btn.addEventListener("click", () => switchModalTab(btn.dataset.tab));
  });
  $("btn-modal-close").addEventListener("click", () => $("modal-item").close());
  $("btn-translate").addEventListener("click", translateActiveItem);

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
