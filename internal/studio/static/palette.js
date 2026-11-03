"use strict";

// Paleta de peças: categorias, busca e lista virtualizada. A parte pura (modelo e filtro) não
// toca o DOM e é testada com Node; a parte de DOM só roda quando createPalette é chamada.
// Depende de search.js (matchesQuery) e virtual.js (virtualWindow).

const PALETTE_ROW = 48;
const PIECE_MIME = "application/x-legoloop-piece";

const PALETTE_CATS = [
  { id: "all", label: "Todas" },
  { id: "trigger", label: "Triggers" },
  { id: "control", label: "Controle" },
  { id: "skill", label: "Skills" },
  { id: "mcp", label: "MCPs" },
  { id: "loop", label: "Loops" },
  { id: "script", label: "Scripts" },
  { id: "hook", label: "Hooks" },
  { id: "gate", label: "Gates" },
];

const PALETTE_ORDER = ["trigger", "control", "skill", "mcp", "loop", "script", "hook", "gate"];

const PALETTE_EMPTY = {
  trigger: "Nenhum trigger configurado.",
  skill: "Nenhuma skill encontrada nas pastas conhecidas dos harnesses.",
  mcp: "Nenhum servidor MCP configurado.",
  loop: "Nenhum loop salvo no acervo ainda.",
  script: "Nenhum script encontrado nas pastas conhecidas.",
  hook: "Nenhum hook configurado.",
  gate: "Nenhum gate disponível.",
  all: "Nada no inventário ainda.",
};

const PALETTE_HINT = {
  trigger: "Arraste ao canvas para iniciar o pipeline com este gatilho.",
  hook: "Arraste sobre uma peça para encaixar no slot de hooks.",
  gate: "Arraste sobre uma peça para encaixar como gate extra.",
  script: "Arraste ao canvas para criar uma peça, ou sobre uma peça para encaixar como script.",
};

function skillKey(ref) {
  return String(ref || "").split("/").pop().toLowerCase();
}

function normName(n) {
  return String(n || "").toLowerCase().replace(/[\s_-]+/g, "-");
}

function entryFromItem(cat, it) {
  const meta = it.meta || {};
  const bits = [meta.signature, meta.transport, meta.host, it.scope, it.origin].filter(Boolean);
  const desc = cat === "hook" || cat === "gate" || cat === "skill" || cat === "loop" ? (it.description || bits.join(" · ")) : bits.join(" · ");
  const e = {
    key: `${cat}:${it.id}`,
    cat,
    kind: it.kind || cat,
    name: it.name,
    desc,
    itemId: it.id,
    origin: it.origin || "",
    scope: it.scope || "",
    path: it.path || "",
    meta,
    raw: it,
    score: meta.score || "",
    tier: meta.tier || "",
    badgeColor: meta.badge_color || "",
    favorite: meta.favorite === "true",
  };
  if (cat === "hook") e.event = meta.event || "";
  if (cat === "gate") e.gateName = it.name;
  return e;
}

function entryFromBlock(b) {
  const isTrig = b.kind && b.kind.startsWith("trigger_");
  const isCtrl = b.kind && (b.kind.startsWith("flow_") || b.kind.startsWith("security_") || b.kind.startsWith("vcs_") || b.kind.startsWith("deploy_"));
  let cat = "skill";
  let badge = "contrato";
  let desc = "";
  const ins = (b.consumes || []).map((t) => String(t).replace("Artifact:", "")).join("|");
  const out = String(b.produces || "").replace("Artifact:", "");

  if (isTrig) {
    cat = "trigger";
    badge = "gatilho";
    desc = `Gatilho inicial → ${out || "pipeline"}`;
  } else if (isCtrl) {
    cat = "control";
    if (b.kind.startsWith("security_")) {
      badge = "segurança";
      desc = "Cofre e credenciais";
    } else if (b.kind.startsWith("vcs_")) {
      badge = "git";
      desc = "Automação de versionamento";
    } else if (b.kind.startsWith("deploy_")) {
      badge = "deploy";
      desc = "Entrega contínua e publicação";
    } else {
      badge = "fluxo";
      desc = "Controle algorítmico de fluxo";
    }
  } else {
    desc = `${ins} → ${out}`;
  }
  return {
    key: `block:${b.id}`,
    cat,
    kind: b.kind || "block",
    name: b.label,
    desc,
    badge,
    blockId: b.id,
    origin: "catálogo",
    scope: "",
    ref: b.ref || "",
    consumes: b.consumes || [],
    produces: b.produces || "",
    artifact: b.artifact || "",
    gate: b.gate || "",
    cost: b.cost || 0,
    suggested_model: b.suggested_model || "",
    raw: b,
  };
}

// buildPaletteModel agrupa itens por SHA-256 e separa peças de controle e triggers.
function buildPaletteModel(blocks, inv) {
  const model = { trigger: [], control: [], skill: [], mcp: [], loop: [], script: [], hook: [], gate: [] };
  const typed = new Set((blocks || []).map((b) => normName(skillKey(b.ref))));
  (blocks || []).map(entryFromBlock).forEach((b) => {
    if (b.cat === "trigger") model.trigger.push(b);
    else if (b.cat === "control") model.control.push(b);
    else model.skill.push(b);
  });
  for (const cat of PALETTE_ORDER) {
    if (cat === "trigger" || cat === "control") continue;
    const rawList = (inv && inv[cat]) || [];
    const groups = new Map();
    for (const it of rawList) {
      if (cat === "skill" && typed.has(normName(it.name))) continue;
      // Unifica duplicatas exatas por SHA-256 ou pelo nome normalizado da peça no catálogo
      const key = it.sha256 || `${cat}:${normName(it.name)}`;
      if (!groups.has(key)) {
        const e = entryFromItem(cat, it);
        e.copies = [it];
        e.copiesCount = 1;
        groups.set(key, e);
      } else {
        const existing = groups.get(key);
        existing.copies.push(it);
        existing.copiesCount++;
      }
    }
    for (const entry of groups.values()) {
      if (entry.copiesCount > 1) {
        entry.badge = `${entry.copiesCount} cópias`;
      }
      model[cat].push(entry);
    }
  }
  for (const list of Object.values(model)) {
    for (const e of list) e.search = `${e.name} ${e.desc} ${e.origin} ${e.scope} ${e.badge || ""}`;
  }
  return model;
}

function paletteCounts(model) {
  const counts = {};
  let all = 0;
  for (const c of PALETTE_ORDER) {
    counts[c] = model[c].length;
    all += model[c].length;
  }
  counts.all = all;
  return counts;
}

// flattenPalette devolve as linhas visíveis: itens da categoria escolhida filtrados pela busca;
// em "Todas", agrupados por categoria com um cabeçalho (só de grupos não vazios).
function flattenPalette(model, cat, query) {
  const cats = cat === "all" ? PALETTE_ORDER : [cat];
  const out = [];
  for (const c of cats) {
    const list = model[c].filter((e) => matchesQuery(e.search, query));
    if (cat === "all" && list.length > 0) {
      out.push({ header: true, key: `header:${c}`, label: PALETTE_CATS.find((x) => x.id === c).label, count: list.length });
    }
    out.push(...list);
  }
  return out;
}

function paletteEmptyMessage(cat, query) {
  const q = String(query || "").trim();
  if (q) return { title: `Nada encontrado para "${q}"`, text: "Tente outro termo ou troque de categoria." };
  return { title: PALETTE_CATS.find((c) => c.id === cat).label, text: PALETTE_EMPTY[cat] };
}

// ---------- DOM ----------

function paletteEl(tag, cls, text) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (text !== undefined) e.textContent = text;
  return e;
}

function paletteRow(entry, handlers) {
  if (entry.header) {
    const li = paletteEl("li", "palette__header");
    li.append(paletteEl("span", "palette__header-label", entry.label), paletteEl("span", "palette__header-count", String(entry.count)));
    return li;
  }
  const li = paletteEl("li", "palette__item");
  li.draggable = true;
  li.tabIndex = 0;
  li.setAttribute("role", "option");
  li.title = entry.desc || entry.name;
  const head = paletteEl("div", "palette__line");
  head.append(paletteEl("span", "palette__name", entry.name));
  if (entry.badge) head.append(paletteEl("span", "badge", entry.badge));
  li.append(head, paletteEl("span", "palette__flow", entry.desc || ""));
  li.addEventListener("dragstart", (e) => {
    const payload = { cat: entry.cat, kind: entry.kind, blockId: entry.blockId, itemId: entry.itemId, name: entry.name, event: entry.event, gateName: entry.gateName };
    e.dataTransfer.setData(PIECE_MIME, JSON.stringify(payload));
    e.dataTransfer.setData("text/plain", entry.name);
    e.dataTransfer.effectAllowed = "copy";
  });
  li.addEventListener("click", () => handlers.onClick(entry));
  li.addEventListener("dblclick", () => handlers.onDouble(entry));
  li.addEventListener("keydown", (e) => {
    if (e.key === "Enter") handlers.onDouble(entry);
    if (e.key === " ") { e.preventDefault(); handlers.onClick(entry); }
  });
  return li;
}

const CAT_LABELS = {
  trigger: "Triggers",
  control: "Controle",
  skill: "Skills",
  mcp: "MCPs",
  loop: "Loops",
  script: "Scripts",
  hook: "Hooks",
  gate: "Gates",
};

const SVG_ICONS = {
  trigger: [
    ["polygon", { points: "13 2 3 14 12 14 11 22 21 10 12 10 13 2" }],
  ],
  control: [
    ["path", { d: "M4 21v-7" }], ["path", { d: "M4 10V3" }],
    ["path", { d: "M12 21v-9" }], ["path", { d: "M12 8V3" }],
    ["path", { d: "M20 21v-5" }], ["path", { d: "M20 12V3" }],
    ["path", { d: "M1 14h6" }], ["path", { d: "M9 8h6" }], ["path", { d: "M17 16h6" }],
  ],
  skill: [
    ["rect", { x: "3", y: "3", width: "7", height: "7" }],
    ["rect", { x: "14", y: "3", width: "7", height: "7" }],
    ["rect", { x: "14", y: "14", width: "7", height: "7" }],
    ["rect", { x: "3", y: "14", width: "7", height: "7" }],
  ],
  mcp: [
    ["path", { d: "M6 3v6" }], ["path", { d: "M18 3v6" }],
    ["path", { d: "M4 9h16v3a8 8 0 0 1-16 0V9z" }],
    ["path", { d: "M12 17v4" }],
  ],
  loop: [
    ["path", { d: "M21.5 2v6h-6" }],
    ["path", { d: "M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.19" }],
  ],
  script: [
    ["path", { d: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" }],
    ["polyline", { points: "14 2 14 8 20 8" }],
    ["line", { x1: "8", y1: "13", x2: "16", y2: "13" }],
    ["line", { x1: "8", y1: "17", x2: "12", y2: "17" }],
  ],
  hook: [
    ["circle", { cx: "12", cy: "5", r: "3" }],
    ["path", { d: "M12 8v8a4 4 0 0 0 4 4h2a2 2 0 0 0 2-2v-1" }],
  ],
  gate: [
    ["path", { d: "M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" }],
    ["path", { d: "m9 12 2 2 4-4" }],
  ],
  eye: [
    ["path", { d: "M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" }],
    ["circle", { cx: "12", cy: "12", r: "3" }],
  ],
  chevron: [
    ["polyline", { points: "6 9 12 15 18 9" }],
  ],
  star: [
    ["polygon", { points: "12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" }],
  ],
  upgrade: [
    ["path", { d: "M12 19V5" }],
    ["polyline", { points: "5 12 12 5 19 12" }],
  ],
  trash: [
    ["polyline", { points: "3 6 5 6 21 6" }],
    ["path", { d: "M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" }],
  ],
};

function createSvgIcon(name, cls) {
  const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
  svg.setAttribute("viewBox", "0 0 24 24");
  svg.setAttribute("fill", "none");
  svg.setAttribute("stroke", "currentColor");
  svg.setAttribute("stroke-width", "2");
  svg.setAttribute("stroke-linecap", "round");
  svg.setAttribute("stroke-linejoin", "round");
  svg.setAttribute("class", cls ? `svg-icon ${cls}` : "svg-icon");
  const defs = SVG_ICONS[name] || SVG_ICONS.skill;
  for (const [tag, attrs] of defs) {
    const el = document.createElementNS("http://www.w3.org/2000/svg", tag);
    for (const [k, v] of Object.entries(attrs)) el.setAttribute(k, v);
    svg.append(el);
  }
  return svg;
}

// createPalette liga o modelo aos elementos da página. Suporta treeview e lista virtualizada.
function createPalette(els, handlers) {
  const st = { model: buildPaletteModel([], {}), cat: "all", query: "", flat: [], raf: 0, activeCategory: "trigger" };

  function renderTreeview() {
    if (!els.treeview) return;
    const q = String(st.query || "").trim();
    const groups = [];
    let totalMatches = 0;

    for (const cat of PALETTE_ORDER) {
      const items = st.model[cat] || [];
      const filtered = q ? items.filter((e) => matchesQuery(e.search, q)) : items;
      totalMatches += filtered.length;

      if (q && filtered.length === 0) continue;

      const groupEl = paletteEl("div", "treeview-group");
      // Sanfona estrita: só 1 grupo aberto por vez (inicialmente Controle)
      const isClosed = q ? false : (st.activeCategory !== cat);
      if (isClosed) groupEl.classList.add("is-closed");

      const head = paletteEl("div", "treeview-group-head");
      head.tabIndex = 0;
      head.setAttribute("role", "button");
      head.setAttribute("aria-expanded", String(!isClosed));

      const titleWrap = paletteEl("div", "treeview-group-head-title");
      const chevron = createSvgIcon("chevron", "svg-icon--chevron");
      const icon = createSvgIcon(cat, "palette-cat-icon");
      const label = paletteEl("span", "treeview-group-label", CAT_LABELS[cat] || cat);
      titleWrap.append(chevron, icon, label);

      const countBadge = paletteEl("span", "treeview-item-badge", String(filtered.length));
      head.append(titleWrap, countBadge);

      head.addEventListener("click", () => {
        if (st.activeCategory === cat) st.activeCategory = null;
        else st.activeCategory = cat;
        renderTreeview();
      });
      head.addEventListener("keydown", (ev) => {
        if (ev.key === "Enter" || ev.key === " ") {
          ev.preventDefault();
          if (st.activeCategory === cat) st.activeCategory = null;
          else st.activeCategory = cat;
          renderTreeview();
        }
      });

      groupEl.append(head);

      if (!isClosed) {
        const itemsList = paletteEl("div", "treeview-items");
        if (filtered.length === 0) {
          itemsList.append(paletteEl("p", "note", "Nenhum item nesta pasta."));
        } else {
          const limit = st.expandedLimits?.[cat] || 40;
          const visibleSlice = filtered.slice(0, limit);
          for (const item of visibleSlice) {
            const row = paletteEl("div", "treeview-item");
            row.draggable = true;
            row.tabIndex = 0;
            row.setAttribute("role", "treeitem");
            row.title = item.desc || item.name;

            const main = paletteEl("div", "treeview-item-main");
            const itemIcon = createSvgIcon(cat, "palette-item-icon");
            const itemName = paletteEl("span", "treeview-item-name", item.name);
            main.append(itemIcon, itemName);

            const actions = paletteEl("div", "treeview-item-actions");

            // 1. Badge Colorida de Tier / Matriz (Ouro, Prata, Bronze, Tóxica)
            if (item.tier) {
              const tierBadge = paletteEl("span", "treeview-tier-badge", item.tier.toUpperCase());
              if (item.badgeColor) {
                tierBadge.style.backgroundColor = item.badgeColor;
                tierBadge.style.color = (item.tier === "ouro" || item.tier === "prata") ? "#0f172a" : "#ffffff";
              }
              tierBadge.title = `Score Matricial: ${item.score || "?"} (${item.tier})`;
              actions.append(tierBadge);
            }

            if (item.badge) {
              actions.append(paletteEl("span", "treeview-item-badge", item.badge));
            }

            // 2. Botão Favoritar (Ícone SVG Canônico)
            const btnFav = paletteEl("button", "btn-item-action btn-fav" + (item.favorite ? " is-favorited" : ""));
            btnFav.type = "button";
            btnFav.title = item.favorite ? "Desfavoritar peça" : "Favoritar como Padrão Canônico (SSOT)";
            btnFav.append(createSvgIcon("star", "svg-icon--fav"));
            btnFav.addEventListener("click", async (e) => {
              e.stopPropagation();
              try {
                const res = await fetch(`/api/assets/${encodeURIComponent(item.kind)}/${encodeURIComponent(item.itemId || item.name)}/favorite`, {
                  method: "POST",
                  headers: { "Content-Type": "application/json" },
                  body: "{}"
                });
                if (res.ok) {
                  const data = await res.json();
                  item.favorite = data.favorited;
                  btnFav.classList.toggle("is-favorited", data.favorited);
                  btnFav.title = data.favorited ? "Desfavoritar peça" : "Favoritar como Padrão Canônico (SSOT)";
                  if (typeof window.toast === "function") {
                    window.toast(data.favorited ? `Peça "${item.name}" favoritada como SSOT!` : `Peça "${item.name}" desfavoritada.`);
                  }
                }
              } catch (err) {
                console.error("Falha ao favoritar:", err);
              }
            });
            actions.append(btnFav);

            // 3. Botão Elevar (Ícone SVG Canônico)
            const btnUpgrade = paletteEl("button", "btn-item-action btn-upgrade");
            btnUpgrade.type = "button";
            btnUpgrade.title = "Elevar peça para Padrão Ouro (Entrevista + Tríade)";
            btnUpgrade.append(createSvgIcon("upgrade", "svg-icon--upgrade"));
            btnUpgrade.addEventListener("click", async (e) => {
              e.stopPropagation();
              try {
                const res = await fetch(`/api/assets/${encodeURIComponent(item.kind)}/${encodeURIComponent(item.itemId || item.name)}/upgrade`, {
                  method: "POST",
                  headers: { "Content-Type": "application/json" },
                  body: "{}"
                });
                if (res.ok) {
                  const data = await res.json();
                  if (typeof window.toast === "function") {
                    window.toast(`Elevação iniciada para "${item.name}"! Score atual: ${data.current_rating?.total_score}`);
                  }
                }
              } catch (err) {
                if (typeof window.toast === "function") {
                  window.toast("Erro ao disparar elevação: " + err.message, true);
                }
              }
            });
            actions.append(btnUpgrade);

            // 4. Botão Eliminar (Ícone SVG Canônico)
            const btnDelete = paletteEl("button", "btn-item-action btn-delete");
            btnDelete.type = "button";
            btnDelete.title = "Eliminar peça profunda";
            btnDelete.append(createSvgIcon("trash", "svg-icon--trash"));
            btnDelete.addEventListener("click", async (e) => {
              e.stopPropagation();
              try {
                const res = await fetch(`/api/assets/${encodeURIComponent(item.kind)}/${encodeURIComponent(item.itemId || item.name)}?mode=full`, {
                  method: "DELETE"
                });
                if (res.ok) {
                  row.remove();
                  if (typeof window.toast === "function") {
                    window.toast(`Peça "${item.name}" eliminada.`);
                  }
                } else {
                  if (typeof window.toast === "function") {
                    window.toast("Falha ao eliminar peça: " + (await res.text()), true);
                  }
                }
              } catch (err) {
                if (typeof window.toast === "function") {
                  window.toast("Erro ao eliminar peça: " + err.message, true);
                }
              }
            });
            actions.append(btnDelete);

            // 5. Botão de Olho para Inspeção (Raio-X Rico sem adicionar ao canvas)
            const btnEye = paletteEl("button", "btn-eye-view");
            btnEye.type = "button";
            btnEye.title = "Inspecionar peça (Raio-X detalhado)";
            btnEye.append(createSvgIcon("eye"));
            btnEye.addEventListener("click", (e) => {
              e.stopPropagation();
              if (handlers.onInspect) handlers.onInspect(item);
            });
            actions.append(btnEye);

            row.append(main, actions);

            row.addEventListener("dragstart", (e) => {
              const payload = {
                cat: item.cat, kind: item.kind, blockId: item.blockId, itemId: item.itemId,
                name: item.name, event: item.event, gateName: item.gateName,
              };
              e.dataTransfer.setData(PIECE_MIME, JSON.stringify(payload));
              e.dataTransfer.setData("text/plain", item.name);
              e.dataTransfer.effectAllowed = "copy";
            });

            // Clique simples insere a peça no Canvas
            row.addEventListener("click", () => handlers.onClick(item));
            row.addEventListener("keydown", (e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                handlers.onClick(item);
              }
            });

            itemsList.append(row);
          }

          if (filtered.length > limit) {
            const moreBtn = paletteEl("button", "treeview-load-more", `Carregar mais (${filtered.length - limit} restantes)`);
            moreBtn.type = "button";
            moreBtn.addEventListener("click", (e) => {
              e.stopPropagation();
              if (!st.expandedLimits) st.expandedLimits = {};
              st.expandedLimits[cat] = limit + 40;
              renderTreeview();
            });
            itemsList.append(moreBtn);
          }
        }
        groupEl.append(itemsList);
      }

      groups.push(groupEl);
    }

    if (q && totalMatches === 0) {
      els.treeview.replaceChildren(paletteEl("p", "note", `Nada encontrado para "${q}"`));
    } else {
      els.treeview.replaceChildren(...groups);
    }

    const counts = paletteCounts(st.model);
    if (els.totalBadge) els.totalBadge.textContent = String(counts.all);
  }
  function renderCats() {
    const counts = paletteCounts(st.model);
    els.cats.replaceChildren(...PALETTE_CATS.map((c) => {
      const b = paletteEl("button", "cat");
      b.type = "button";
      b.setAttribute("role", "tab");
      b.setAttribute("aria-selected", String(c.id === st.cat));
      b.append(paletteEl("span", "cat__label", c.label), paletteEl("span", "cat__count", String(counts[c.id])));
      b.addEventListener("click", () => { st.cat = c.id; els.list.scrollTop = 0; rebuild(); });
      return b;
    }));
    els.hint.textContent = PALETTE_HINT[st.cat] || "";
    els.hint.hidden = !PALETTE_HINT[st.cat];
  }

  function renderRows() {
    const w = virtualWindow(st.flat.length, PALETTE_ROW, els.list.scrollTop, els.list.clientHeight || 600, 6);
    els.spacer.style.height = `${w.height}px`;
    els.rows.style.transform = `translateY(${w.offset}px)`;
    els.rows.replaceChildren(...st.flat.slice(w.start, w.end).map((e) => paletteRow(e, handlers)));
    const empty = st.flat.length === 0;
    els.empty.hidden = !empty;
    if (empty) {
      const m = paletteEmptyMessage(st.cat, st.query);
      els.emptyTitle.textContent = m.title;
      els.emptyText.textContent = m.text;
    }
  }

  function rebuild() {
    st.flat = flattenPalette(st.model, st.cat, st.query);
    if (els.cats) renderCats();
    if (els.rows && els.spacer && els.list) renderRows();
    renderTreeview();
  }

  function onScroll() {
    if (st.raf) return;
    st.raf = requestAnimationFrame(() => { st.raf = 0; renderRows(); });
  }

  els.list.addEventListener("scroll", onScroll);
  if (typeof ResizeObserver === "function") new ResizeObserver(onScroll).observe(els.list);
  els.search.addEventListener("input", (e) => { st.query = e.target.value; els.list.scrollTop = 0; rebuild(); });
  els.search.addEventListener("keydown", (e) => {
    if (e.key === "Escape") { e.target.value = ""; st.query = ""; rebuild(); }
  });
  document.addEventListener("keydown", (e) => {
    if (e.key === "/" && !e.target.closest("input, select, textarea")) { e.preventDefault(); els.search.focus(); }
  });
  rebuild();

  return {
    setModel(model) { st.model = model; rebuild(); },
    setBanner(text) { els.banner.textContent = text || ""; els.banner.hidden = !text; },
    focusSearch() { els.search.focus(); },
    rendered() { return els.rows.children.length; },
    total() { return st.flat.length; },
  };
}
