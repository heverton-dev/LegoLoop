const KRAFTEN_12_COLORS = [
  { name: "missioncontrol", hex: "#38bdf8", soft: "rgba(56, 189, 248, 0.12)", hover: "#0284c7" },
  { name: "purgesync", hex: "#22c55e", soft: "rgba(34, 197, 94, 0.12)", hover: "#16a34a" },
  { name: "sniped", hex: "#f59e0b", soft: "rgba(245, 158, 11, 0.12)", hover: "#d97706" },
  { name: "legoloop", hex: "#96a0f8", soft: "rgba(150, 160, 248, 0.12)", hover: "#808bf2" },
  { name: "vaultkeeper", hex: "#fb923c", soft: "rgba(251, 146, 60, 0.12)", hover: "#ea580c" },
  { name: "studioreel", hex: "#fb7c89", soft: "rgba(251, 124, 137, 0.12)", hover: "#f43f5e" },
  { name: "flightmonitor", hex: "#66a6fa", soft: "rgba(102, 166, 250, 0.12)", hover: "#3b82f6" },
  { name: "lime", hex: "#a3e635", soft: "rgba(163, 230, 53, 0.12)", hover: "#84cc16" },
  { name: "mint", hex: "#2dd4bf", soft: "rgba(45, 212, 191, 0.12)", hover: "#14b8a6" },
  { name: "cyan", hex: "#22d3ee", soft: "rgba(34, 211, 238, 0.12)", hover: "#06b6d4" },
  { name: "violet", hex: "#c494fc", soft: "rgba(196, 148, 252, 0.12)", hover: "#a855f7" },
  { name: "magenta", hex: "#f88bc4", soft: "rgba(248, 139, 196, 0.12)", hover: "#ec4899" }
];

const ITENS_POR_PAGINA = 2;
let paginaAtual = 1;
let filtroModo = "todos"; // "todos", "sync", "async"
let processosArmazenados = [];

function resolveProcessColor(name) {
  const norm = String(name || "").toLowerCase().trim();
  if (norm.includes("concepcao")) return KRAFTEN_12_COLORS[9];
  if (norm.includes("fabricar")) return KRAFTEN_12_COLORS[10];
  if (norm.includes("especificacao") || norm.includes("spec")) return KRAFTEN_12_COLORS[0];
  if (norm.includes("plano") || norm.includes("planejar")) return KRAFTEN_12_COLORS[2];
  if (norm.includes("swarm")) return KRAFTEN_12_COLORS[6];
  if (norm.includes("qa") || norm.includes("teste")) return KRAFTEN_12_COLORS[1];
  if (norm.includes("design") || norm.includes("ui")) return KRAFTEN_12_COLORS[11];

  let hash = 0;
  for (let i = 0; i < norm.length; i++) {
    hash = (hash * 31 + norm.charCodeAt(i)) >>> 0;
  }
  return KRAFTEN_12_COLORS[hash % KRAFTEN_12_COLORS.length];
}

function resolveExecutionMode(item) {
  const norm = String(item.processo || "").toLowerCase().trim();
  if (norm.includes("concepcao") || norm.includes("entrevista") || norm.includes("triagem")) {
    return { label: "Síncrono", desc: "Interativo", isSync: true };
  }
  return { label: "Assíncrono", desc: "Autônomo", isSync: false };
}

function humanizeProcessName(raw) {
  let name = String(raw || "").trim();
  if (name.startsWith("proc:")) name = name.substring(5);
  if (name === "legoloop-concepcao") return "LegoLoop Concepção";
  if (name === "legoloop-fabricar") return "LegoLoop Fabricar";
  if (name === "legoloop-especificacao") return "LegoLoop Especificação";
  if (name === "legoloop-planejamento") return "LegoLoop Planejamento";
  if (name === "legoloop-ciclo") return "LegoLoop Ciclo";
  return name.replace(/^legoloop-/, "").replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

function formatDuration(timestampMs) {
  if (!timestampMs || timestampMs <= 0) return "0s";
  const diffSec = Math.max(0, Math.floor((Date.now() - timestampMs) / 1000));
  const h = Math.floor(diffSec / 3600);
  const m = Math.floor((diffSec % 3600) / 60);
  const s = diffSec % 60;
  if (h > 0) return `${h}h ${m}m`;
  if (m > 0) return `${m}m ${s}s`;
  return `${s}s`;
}

const listEl = document.getElementById("processos-lista");
const paginationEl = document.getElementById("processos-paginacao");
const sidebarTreeEl = document.getElementById("processos-sidebar-tree");
const countEl = document.getElementById("processos-count");
const feedbackEl = document.getElementById("processos-feedback");
const btnRefresh = document.getElementById("btn-refresh");

const countTodosEl = document.getElementById("count-todos");
const countSyncEl = document.getElementById("count-sync");
const countAsyncEl = document.getElementById("count-async");
const filtroBtns = document.querySelectorAll(".btn--filtro");

function showFeedback(msg, isError) {
  if (!feedbackEl) return;
  feedbackEl.hidden = !msg;
  feedbackEl.className = isError ? "a2ui-feedback a2ui-feedback--error" : "a2ui-feedback a2ui-feedback--info";
  feedbackEl.textContent = msg;
}

async function carregarProcessos() {
  try {
    const res = await fetch("/api/processos/ativos");
    if (!res.ok) {
      throw new Error("HTTP " + res.status);
    }
    const itens = await res.json();
    processosArmazenados = Array.isArray(itens) ? itens : [];
    renderProcessos();
  } catch (err) {
    showFeedback("Falha ao atualizar processos: " + err.message, true);
  }
}

function getItensFiltrados() {
  if (filtroModo === "sync") {
    return processosArmazenados.filter((i) => resolveExecutionMode(i).isSync);
  }
  if (filtroModo === "async") {
    return processosArmazenados.filter((i) => !resolveExecutionMode(i).isSync);
  }
  return processosArmazenados;
}

function atualizarContadoresFiltros() {
  const total = processosArmazenados.length;
  const sync = processosArmazenados.filter((i) => resolveExecutionMode(i).isSync).length;
  const async = total - sync;

  if (countTodosEl) countTodosEl.textContent = total;
  if (countSyncEl) countSyncEl.textContent = sync;
  if (countAsyncEl) countAsyncEl.textContent = async;
}

function renderSidebar(itens) {
  if (!sidebarTreeEl) return;
  sidebarTreeEl.replaceChildren();

  if (!itens || itens.length === 0) {
    const li = document.createElement("li");
    li.style.padding = "8px 10px";
    li.style.fontSize = "12px";
    li.style.color = "var(--kf-text-secondary)";
    li.textContent = "Nenhum ativo";
    sidebarTreeEl.appendChild(li);
    return;
  }

  itens.forEach((item, index) => {
    const col = resolveProcessColor(item.processo);
    const mode = resolveExecutionMode(item);
    const li = document.createElement("li");
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "processos-sidebar-btn";

    const left = document.createElement("div");
    left.style.display = "flex";
    left.style.flexDirection = "column";
    left.style.gap = "2px";
    left.style.overflow = "hidden";

    const titleRow = document.createElement("div");
    titleRow.style.display = "flex";
    titleRow.style.alignItems = "center";
    titleRow.style.gap = "6px";

    const dot = document.createElement("span");
    dot.className = "processos-sidebar-dot";
    dot.style.backgroundColor = col.hex;

    const label = document.createElement("span");
    label.style.fontWeight = "600";
    label.style.color = "var(--kf-text-primary)";
    label.style.overflow = "hidden";
    label.style.textOverflow = "ellipsis";
    label.style.whiteSpace = "nowrap";
    label.textContent = humanizeProcessName(item.processo);

    titleRow.append(dot, label);

    const modoPill = document.createElement("span");
    modoPill.className = "processos-sidebar-modo " + (mode.isSync ? "processos-sidebar-modo--sync" : "processos-sidebar-modo--async");
    modoPill.textContent = mode.label;

    left.append(titleRow, modoPill);

    const right = document.createElement("div");
    right.style.display = "flex";
    right.style.flexDirection = "column";
    right.style.alignItems = "flex-end";
    right.style.gap = "2px";

    const stPill = document.createElement("span");
    stPill.className = "pill";
    stPill.style.fontSize = "10px";
    stPill.style.padding = "1px 6px";
    stPill.style.flexShrink = "0";

    const st = item.state || "running";
    if (st === "paused") {
      stPill.className = "pill pill--andon-paused";
      stPill.textContent = "Pausa";
    } else if (st === "cancelled") {
      stPill.style.background = "#ef4444";
      stPill.style.color = "#ffffff";
      stPill.textContent = "Canc";
    } else {
      stPill.className = "pill pill--success";
      stPill.textContent = "Ativo";
    }

    const durEl = document.createElement("span");
    durEl.style.fontSize = "9.5px";
    durEl.style.color = "var(--kf-text-secondary)";
    durEl.style.fontFamily = "var(--kf-font-mono, monospace)";
    durEl.textContent = formatDuration(item.timestamp_ms);

    right.append(stPill, durEl);

    btn.append(left, right);
    btn.addEventListener("click", () => {
      const pageTarget = Math.floor(index / ITENS_POR_PAGINA) + 1;
      if (pageTarget !== paginaAtual) {
        paginaAtual = pageTarget;
        renderCardsEPaginacao();
      }
      setTimeout(() => {
        const card = document.getElementById("card-" + item.id);
        if (card) {
          card.scrollIntoView({ behavior: "smooth", block: "center" });
          card.style.borderColor = "var(--kf-border-focus)";
          setTimeout(() => {
            card.style.borderColor = "";
          }, 1200);
        }
      }, 50);
    });

    li.appendChild(btn);
    sidebarTreeEl.appendChild(li);
  });
}

function renderProcessos() {
  const arr = getItensFiltrados();
  if (countEl) {
    countEl.textContent = arr.length === 1 ? "1 processo ativo" : arr.length + " processos ativos";
  }

  atualizarContadoresFiltros();
  renderSidebar(arr);
  renderCardsEPaginacao();
}

function renderCardsEPaginacao() {
  if (!listEl) return;
  listEl.replaceChildren();

  const arr = getItensFiltrados();
  if (arr.length === 0) {
    const empty = document.createElement("div");
    empty.className = "a2ui-painel";
    empty.style.textAlign = "center";
    empty.style.padding = "32px";
    empty.textContent = "Nenhum processo em execução para o filtro selecionado.";
    listEl.appendChild(empty);
    if (paginationEl) paginationEl.replaceChildren();
    return;
  }

  const totalPaginas = Math.ceil(arr.length / ITENS_POR_PAGINA) || 1;
  if (paginaAtual > totalPaginas) paginaAtual = totalPaginas;
  if (paginaAtual < 1) paginaAtual = 1;

  const inicio = (paginaAtual - 1) * ITENS_POR_PAGINA;
  const itensPagina = arr.slice(inicio, inicio + ITENS_POR_PAGINA);

  itensPagina.forEach((item) => {
    const col = resolveProcessColor(item.processo);
    const mode = resolveExecutionMode(item);
    const st = item.state || "running";

    const card = document.createElement("article");
    card.id = "card-" + item.id;
    card.className = "kf-card-processo";

    const header = document.createElement("header");
    header.className = "kf-card-processo__header";

    const titleRow = document.createElement("div");
    titleRow.className = "kf-card-processo__title-row";

    const dot = document.createElement("span");
    dot.className = "kf-card-processo__dot";
    dot.style.backgroundColor = col.hex;

    const nameEl = document.createElement("h2");
    nameEl.className = "kf-card-processo__name";
    nameEl.textContent = humanizeProcessName(item.processo);

    const sessPill = document.createElement("span");
    sessPill.className = "kf-card-processo__session-id";
    sessPill.textContent = item.id;

    const modePill = document.createElement("span");
    modePill.className = "pill " + (mode.isSync ? "pill--sync" : "pill--async");
    modePill.textContent = mode.label;
    modePill.title = "Modo: " + mode.label + " (" + mode.desc + ")";

    titleRow.append(dot, nameEl, sessPill, modePill);

    const rightGroup = document.createElement("div");
    rightGroup.style.display = "flex";
    rightGroup.style.alignItems = "center";
    rightGroup.style.gap = "8px";

    const durEl = document.createElement("span");
    durEl.className = "kf-card-processo__duration";
    durEl.title = "Duração ativa";
    durEl.innerHTML = `<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg><span>${formatDuration(item.timestamp_ms)}</span>`;

    const statusBadge = document.createElement("span");
    statusBadge.className = "kf-status-badge " + (st === "paused" ? "kf-status-badge--paused" : (st === "cancelled" ? "kf-status-badge--cancelled" : "kf-status-badge--running"));

    const statusDot = document.createElement("span");
    statusDot.className = "kf-status-badge__dot";

    const statusText = document.createElement("span");
    statusText.textContent = st === "paused" ? "Pausado" : (st === "cancelled" ? "Cancelado" : "Executando");

    statusBadge.append(statusDot, statusText);
    rightGroup.append(durEl, statusBadge);

    header.append(titleRow, rightGroup);
    card.appendChild(header);

    const body = document.createElement("div");
    body.className = "kf-card-processo__body";
    body.style.display = "flex";
    body.style.alignItems = "center";
    body.style.justifyContent = "space-between";
    body.style.flexWrap = "wrap";
    body.style.gap = "10px";

    const telemetry = document.createElement("div");
    telemetry.className = "kf-card-processo__telemetry";
    telemetry.innerHTML = `<span><strong>${item.surfaces}</strong> ${item.surfaces === 1 ? 'superfície ativa' : 'superfícies ativas'}</span><span class="kf-telemetry-sep">·</span><span><strong>${item.subscribers}</strong> ${item.subscribers === 1 ? 'conexão' : 'conexões'}</span>`;
    body.appendChild(telemetry);

    if (item.etapa_atual) {
      const etapaChip = document.createElement("span");
      etapaChip.className = "kf-card-processo__etapa-chip";
      etapaChip.title = "Etapa atual: " + item.etapa_atual;
      etapaChip.textContent = "Etapa: " + item.etapa_atual;
      body.appendChild(etapaChip);
    }

    card.appendChild(body);

    const footer = document.createElement("footer");
    footer.className = "kf-card-processo__footer";

    const actions = document.createElement("div");
    actions.className = "kf-card-processo__actions";

    if (st === "paused") {
      const btnResume = document.createElement("button");
      btnResume.type = "button";
      btnResume.className = "btn btn--andon-resume btn--action-compact";
      btnResume.textContent = "Retomar";
      btnResume.addEventListener("click", async () => {
        btnResume.disabled = true;
        try {
          await fetch("/api/processos/" + encodeURIComponent(item.id) + "/resume", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: "{}"
          });
        } finally {
          await carregarProcessos();
        }
      });

      const btnCancel = document.createElement("button");
      btnCancel.type = "button";
      btnCancel.className = "btn btn--andon-cancel btn--action-compact";
      btnCancel.textContent = "Cancelar";
      btnCancel.addEventListener("click", async () => {
        btnCancel.disabled = true;
        try {
          await fetch("/api/processos/" + encodeURIComponent(item.id) + "/cancel", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: "{}"
          });
        } finally {
          await carregarProcessos();
        }
      });
      actions.append(btnResume, btnCancel);
    } else if (st !== "cancelled") {
      const btnAndon = document.createElement("button");
      btnAndon.type = "button";
      btnAndon.className = "btn btn--andon btn--action-compact";
      btnAndon.textContent = "Parar";
      btnAndon.title = "Parar execução da esteira (Andon Cord) [Alt+P]";
      btnAndon.addEventListener("click", async () => {
        btnAndon.disabled = true;
        try {
          await fetch("/api/processos/" + encodeURIComponent(item.id) + "/andon", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: "{}"
          });
        } finally {
          await carregarProcessos();
        }
      });
      actions.appendChild(btnAndon);
    }

    const linkSessao = document.createElement("a");
    linkSessao.className = "btn btn--action-compact";
    linkSessao.href = item.url;
    linkSessao.textContent = "Abrir";
    actions.appendChild(linkSessao);

    footer.appendChild(actions);
    card.appendChild(footer);
    listEl.appendChild(card);
  });

  renderPaginacaoControles(totalPaginas);
}

function renderPaginacaoControles(totalPaginas) {
  if (!paginationEl) return;
  paginationEl.replaceChildren();

  if (totalPaginas <= 1) return;

  const btnPrev = document.createElement("button");
  btnPrev.type = "button";
  btnPrev.className = "btn--paginacao";
  btnPrev.textContent = "◀ Anterior";
  btnPrev.disabled = paginaAtual <= 1;
  btnPrev.addEventListener("click", () => {
    if (paginaAtual > 1) {
      paginaAtual--;
      renderCardsEPaginacao();
    }
  });

  const info = document.createElement("span");
  info.className = "processos-paginacao__info";
  info.textContent = `Página ${paginaAtual} de ${totalPaginas}`;

  const btnNext = document.createElement("button");
  btnNext.type = "button";
  btnNext.className = "btn--paginacao";
  btnNext.textContent = "Próxima ▶";
  btnNext.disabled = paginaAtual >= totalPaginas;
  btnNext.addEventListener("click", () => {
    if (paginaAtual < totalPaginas) {
      paginaAtual++;
      renderCardsEPaginacao();
    }
  });

  paginationEl.append(btnPrev, info, btnNext);
}

filtroBtns.forEach((btn) => {
  btn.addEventListener("click", () => {
    filtroBtns.forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    filtroModo = btn.dataset.filtro || "todos";
    paginaAtual = 1;
    renderCardsEPaginacao();
  });
});

// Atalho de teclado global Alt+P para Parar / Retomar esteira
window.addEventListener("keydown", async (e) => {
  if (e.altKey && (e.key === "p" || e.key === "P")) {
    e.preventDefault();
    const ativo = processosArmazenados.find((p) => (p.state || "running") === "running");
    if (ativo) {
      await fetch("/api/processos/" + encodeURIComponent(ativo.id) + "/andon", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{}"
      });
      await carregarProcessos();
      return;
    }
    const pausado = processosArmazenados.find((p) => p.state === "paused");
    if (pausado) {
      await fetch("/api/processos/" + encodeURIComponent(pausado.id) + "/resume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{}"
      });
      await carregarProcessos();
    }
  }
});

if (btnRefresh) {
  btnRefresh.addEventListener("click", () => {
    carregarProcessos();
  });
}

carregarProcessos();
setInterval(carregarProcessos, 3000);
// Atualização a cada segundo para o temporizador de duração ativa
setInterval(() => {
  const durEls = document.querySelectorAll(".kf-card-processo__duration span");
  const itens = getItensFiltrados();
  const inicio = (paginaAtual - 1) * ITENS_POR_PAGINA;
  const itensPagina = itens.slice(inicio, inicio + ITENS_POR_PAGINA);
  durEls.forEach((el, idx) => {
    if (itensPagina[idx]) {
      el.textContent = formatDuration(itensPagina[idx].timestamp_ms);
    }
  });
}, 1000);
