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

function getContrastInk(hexColor) {
  let c = String(hexColor || "#22d3ee").replace("#", "");
  if (c.length === 3) {
    c = c.split("").map((x) => x + x).join("");
  }
  const r = parseInt(c.substring(0, 2), 16) / 255;
  const g = parseInt(c.substring(2, 4), 16) / 255;
  const b = parseInt(c.substring(4, 6), 16) / 255;
  const a = [r, g, b].map((v) => (v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4)));
  const lum = 0.2126 * a[0] + 0.7152 * a[1] + 0.0722 * a[2];
  return lum > 0.42 ? "#111827" : "#ffffff";
}

const listEl = document.getElementById("processos-lista");
const countEl = document.getElementById("processos-count");
const feedbackEl = document.getElementById("processos-feedback");
const btnRefresh = document.getElementById("btn-refresh");

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
    renderProcessos(itens);
  } catch (err) {
    showFeedback("Falha ao atualizar processos: " + err.message, true);
  }
}

function renderProcessos(itens) {
  if (!listEl) return;
  listEl.replaceChildren();
  const arr = Array.isArray(itens) ? itens : [];
  if (countEl) {
    countEl.textContent = arr.length === 1 ? "1 processo ativo" : arr.length + " processos ativos";
  }
  if (arr.length === 0) {
    const empty = document.createElement("div");
    empty.className = "a2ui-painel";
    empty.style.textAlign = "center";
    empty.style.padding = "32px";
    empty.textContent = "Nenhum processo em execucao no momento.";
    listEl.appendChild(empty);
    return;
  }

  arr.forEach((item) => {
    const card = document.createElement("article");
    card.className = "a2ui-painel kf-card";
    card.style.display = "flex";
    card.style.flexDirection = "column";
    card.style.gap = "12px";

    const header = document.createElement("header");
    header.style.display = "flex";
    header.style.alignItems = "center";
    header.style.justifyContent = "space-between";
    header.style.gap = "8px";

    const left = document.createElement("div");
    left.style.display = "flex";
    left.style.alignItems = "center";
    left.style.gap = "8px";

    const col = resolveProcessColor(item.processo);
    const ink = getContrastInk(col.hex);

    const procPill = document.createElement("span");
    procPill.className = "pill";
    procPill.textContent = item.processo;
    procPill.style.backgroundColor = col.hex;
    procPill.style.color = ink;
    procPill.style.fontWeight = "700";

    const sessPill = document.createElement("span");
    sessPill.className = "pill";
    sessPill.textContent = item.id;

    left.append(procPill, sessPill);

    const statePill = document.createElement("span");
    const st = item.state || "running";
    if (st === "paused") {
      statePill.className = "pill pill--andon-paused";
      statePill.textContent = "Pausado (Andon)";
    } else if (st === "cancelled") {
      statePill.className = "pill";
      statePill.style.background = "rgba(239, 68, 68, 0.12)";
      statePill.style.color = "#dc2626";
      statePill.textContent = "Cancelado";
    } else {
      statePill.className = "pill pill--success";
      statePill.textContent = "Executando";
    }

    header.append(left, statePill);
    card.appendChild(header);

    const info = document.createElement("div");
    info.style.fontSize = "12px";
    info.style.color = "var(--kf-text-secondary)";
    info.textContent = "Superficies ativas: " + item.surfaces + " • Conexoes ativas: " + item.subscribers;
    card.appendChild(info);

    const actions = document.createElement("div");
    actions.style.display = "flex";
    actions.style.alignItems = "center";
    actions.style.gap = "8px";
    actions.style.marginTop = "4px";

    if (st === "paused") {
      const btnResume = document.createElement("button");
      btnResume.type = "button";
      btnResume.className = "btn btn--andon-resume";
      btnResume.textContent = "Retomar";
      btnResume.addEventListener("click", async () => {
        await fetch("/api/processos/" + encodeURIComponent(item.id) + "/resume", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: "{}"
        });
        await carregarProcessos();
      });

      const btnCancel = document.createElement("button");
      btnCancel.type = "button";
      btnCancel.className = "btn btn--andon-cancel";
      btnCancel.textContent = "Cancelar";
      btnCancel.addEventListener("click", async () => {
        await fetch("/api/processos/" + encodeURIComponent(item.id) + "/cancel", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: "{}"
        });
        await carregarProcessos();
      });
      actions.append(btnResume, btnCancel);
    } else if (st !== "cancelled") {
      const btnAndon = document.createElement("button");
      btnAndon.type = "button";
      btnAndon.className = "btn btn--andon";
      btnAndon.textContent = "Andon Cord (Parar)";
      btnAndon.addEventListener("click", async () => {
        await fetch("/api/processos/" + encodeURIComponent(item.id) + "/andon", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: "{}"
        });
        await carregarProcessos();
      });
      actions.appendChild(btnAndon);
    }

    const linkSessao = document.createElement("a");
    linkSessao.className = "btn";
    linkSessao.href = item.url;
    linkSessao.textContent = "Abrir Sessao A2UI";
    actions.appendChild(linkSessao);

    card.appendChild(actions);
    listEl.appendChild(card);
  });
}

if (btnRefresh) {
  btnRefresh.addEventListener("click", () => {
    carregarProcessos();
  });
}

carregarProcessos();
setInterval(carregarProcessos, 3000);
