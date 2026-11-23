import { MessageProcessor, A2UISurfaceRenderer } from "/components/a2ui-bundle.js";
import { registerLegoLoopCatalog } from "/components/legoloop-catalog.js";

registerLegoLoopCatalog();

const parts = window.location.pathname.split("/").filter(Boolean);
const sessionId = decodeURIComponent(parts[parts.length - 1] || "");
const statusEl = document.getElementById("a2ui-conn-status");
const sessionEl = document.getElementById("a2ui-session-id");
const processEl = document.getElementById("a2ui-process-id");
const treeEl = document.getElementById("a2ui-tree");
const breadcrumbEl = document.getElementById("a2ui-breadcrumb");
const rootEl = document.getElementById("a2ui-root");
const andonContainer = document.getElementById("a2ui-andon-group");

function renderAndonControls(state) {
  if (!andonContainer) return;
  andonContainer.replaceChildren();

  if (state === "paused") {
    const pill = document.createElement("span");
    pill.className = "pill pill--andon-paused";
    pill.textContent = "Pausado (Andon)";

    const btnResume = document.createElement("button");
    btnResume.type = "button";
    btnResume.className = "btn btn--andon-resume";
    btnResume.textContent = "Retomar";
    btnResume.addEventListener("click", () => {
      sendAction({ name: "andon_cord_resume" });
    });

    const btnCancel = document.createElement("button");
    btnCancel.type = "button";
    btnCancel.className = "btn btn--andon-cancel";
    btnCancel.textContent = "Cancelar";
    btnCancel.addEventListener("click", () => {
      sendAction({ name: "andon_cord_cancel" });
    });

    andonContainer.append(pill, btnResume, btnCancel);
  } else if (state === "cancelled") {
    const pill = document.createElement("span");
    pill.className = "pill pill--andon-paused";
    pill.textContent = "Processo Cancelado";
    andonContainer.appendChild(pill);
  } else {
    const btnAndon = document.createElement("button");
    btnAndon.type = "button";
    btnAndon.className = "btn btn--andon";
    btnAndon.textContent = "Andon Cord";
    btnAndon.title = "Andon Cord: Parada Emergencial da Esteira";
    btnAndon.addEventListener("click", () => {
      sendAction({ name: "andon_cord_pause" });
    });
    andonContainer.appendChild(btnAndon);
  }
}

renderAndonControls("running");

sessionEl.textContent = sessionId || "sem-sessao";

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
  if (norm.includes("concepcao")) return KRAFTEN_12_COLORS[9]; // cyan: #22d3ee
  if (norm.includes("fabricar")) return KRAFTEN_12_COLORS[10]; // violet: #c494fc
  if (norm.includes("especificacao") || norm.includes("spec")) return KRAFTEN_12_COLORS[0]; // missioncontrol: #38bdf8
  if (norm.includes("plano") || norm.includes("planejar")) return KRAFTEN_12_COLORS[2]; // sniped: #f59e0b
  if (norm.includes("swarm")) return KRAFTEN_12_COLORS[6]; // flightmonitor: #66a6fa
  if (norm.includes("qa") || norm.includes("teste")) return KRAFTEN_12_COLORS[1]; // purgesync: #22c55e
  if (norm.includes("design") || norm.includes("ui")) return KRAFTEN_12_COLORS[11]; // magenta: #f88bc4

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

let activeProcessName = "";

function applyProcessTheme(processName) {
  if (!processName || processName === activeProcessName) {
    return;
  }
  activeProcessName = processName;
  const col = resolveProcessColor(processName);
  const ink = getContrastInk(col.hex);
  const root = document.documentElement;
  root.style.setProperty("--kf-process-color", col.hex);
  root.style.setProperty("--kf-process-color-soft", col.soft);
  root.style.setProperty("--kf-process-color-hover", col.hover);
  root.style.setProperty("--kf-process-ink", ink);

  if (processEl) {
    processEl.textContent = processName;
    processEl.style.backgroundColor = col.hex;
    processEl.style.borderColor = col.hover;
    processEl.style.color = ink;
    processEl.style.fontWeight = "700";
    processEl.hidden = false;
  }
}

function humanizeSurfaceName(surfaceId, index) {
  if (!surfaceId) {
    return String(index);
  }
  const clean = String(surfaceId).replace(/^proc:/i, "");
  const parts = clean.split(":").filter(Boolean);
  const rawLast = parts[parts.length - 1] || clean;
  let name = rawLast.replace(/[-_]/g, " ");
  name = name.split(" ").map((w) => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase()).join(" ");
  name = name.replace(/Stop Gate/i, "Stop-Gate");
  const num = String(index).padStart(2, "0");
  return num + ". " + name;
}

function detectProcessFromSurface(surfaceId) {
  const sid = String(surfaceId || "").toLowerCase();
  if (sid.includes("fabricar")) {
    return "legoloop-fabricar";
  }
  if (sid.includes("concepcao")) {
    return "legoloop-concepcao";
  }
  if (sid.startsWith("proc:")) {
    const p = sid.split(":");
    return p[1] || sid;
  }
  return "";
}

const processor = new MessageProcessor();
let activeSurfaceId = "";

function showFeedback(msg, kind) {
  let box = document.getElementById("a2ui-feedback");
  if (!box) {
    box = document.createElement("div");
    box.id = "a2ui-feedback";
    box.setAttribute("role", "status");
    rootEl.parentElement.insertBefore(box, rootEl);
  }
  box.className = "a2ui-feedback a2ui-feedback--" + (kind || "info");
  box.textContent = msg;
  box.hidden = !msg;
}

function createLegoBrick(svgNS, x, y, fill, className) {
  const g = document.createElementNS(svgNS, "g");
  if (className) {
    g.setAttribute("class", className);
  }

  const stud1 = document.createElementNS(svgNS, "rect");
  stud1.setAttribute("x", String(x + 8));
  stud1.setAttribute("y", String(y - 4));
  stud1.setAttribute("width", "12");
  stud1.setAttribute("height", "5");
  stud1.setAttribute("rx", "1.5");
  stud1.setAttribute("fill", fill);

  const stud2 = document.createElementNS(svgNS, "rect");
  stud2.setAttribute("x", String(x + 32));
  stud2.setAttribute("y", String(y - 4));
  stud2.setAttribute("width", "12");
  stud2.setAttribute("height", "5");
  stud2.setAttribute("rx", "1.5");
  stud2.setAttribute("fill", fill);

  const body = document.createElementNS(svgNS, "rect");
  body.setAttribute("x", String(x));
  body.setAttribute("y", String(y));
  body.setAttribute("width", "52");
  body.setAttribute("height", "18");
  body.setAttribute("rx", "3");
  body.setAttribute("fill", fill);

  g.appendChild(stud1);
  g.appendChild(stud2);
  g.appendChild(body);
  return g;
}

function createLoadingCard() {
  const card = document.createElement("div");
  card.className = "a2ui-loading-card";

  const svgNS = "http://www.w3.org/2000/svg";
  const svg = document.createElementNS(svgNS, "svg");
  svg.setAttribute("class", "a2ui-loading-spinner");
  svg.setAttribute("viewBox", "0 0 72 72");
  svg.setAttribute("fill", "none");

  const brick1 = createLegoBrick(svgNS, 10, 40, "var(--kf-process-color, #22d3ee)", "lego-brick-ladder lego-brick-ladder--1");
  const brick2 = createLegoBrick(svgNS, 10, 40, "var(--kf-text-primary, #111827)", "lego-brick-ladder lego-brick-ladder--2");

  svg.appendChild(brick1);
  svg.appendChild(brick2);
  card.appendChild(svg);

  const title = document.createElement("h4");
  title.className = "a2ui-loading-title";
  title.textContent = "LegoLoop em Processamento";
  card.appendChild(title);

  const text = document.createElement("p");
  text.id = "a2ui-loading-text";
  text.className = "a2ui-loading-text";
  text.textContent = "O motor esta avaliando as respostas e atualizando o modelo mental...";
  card.appendChild(text);

  return card;
}

function showLoading(msg) {
  let overlay = document.getElementById("a2ui-loading-overlay");
  if (!overlay) {
    overlay = document.createElement("div");
    overlay.id = "a2ui-loading-overlay";
    overlay.className = "a2ui-loading-overlay";
    overlay.appendChild(createLoadingCard());
    document.body.appendChild(overlay);
  }
  const textEl = document.getElementById("a2ui-loading-text");
  if (textEl && msg) {
    textEl.textContent = msg;
  }
  overlay.classList.add("a2ui-loading-overlay--active");
}

function hideLoading() {
  const overlay = document.getElementById("a2ui-loading-overlay");
  if (overlay) {
    overlay.classList.remove("a2ui-loading-overlay--active");
  }
}

async function sendAction(action) {
  if (!sessionId) {
    return;
  }
  const url = "/api/sessao/" + encodeURIComponent(sessionId) + "/action";
  const btns = rootEl.querySelectorAll("button");
  btns.forEach((b) => { b.disabled = true; });
  statusEl.textContent = "Processando...";
  const msgLoading = action && action.name === "stop_gate_decision"
    ? "Aplicando decisão do Stop-Gate e estruturando o processo..."
    : "O motor de IA está analisando sua resposta e calibrando a confiança...";
  showLoading(msgLoading);
  showFeedback("Enviando resposta e aguardando o motor...", "info");
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(action)
    });
    if (!res.ok) {
      const errText = await res.text();
      showFeedback("Erro (" + res.status + "): " + errText, "error");
      statusEl.textContent = "Erro";
      hideLoading();
    } else {
      showFeedback("", "info");
      statusEl.textContent = "Conectado";
    }
  } catch (_err) {
    showFeedback("Sem conexao com o servidor da sessao.", "error");
    statusEl.textContent = "Desconectado";
    hideLoading();
  } finally {
    btns.forEach((b) => { b.disabled = false; });
  }
}

const renderer = new A2UISurfaceRenderer(rootEl, sendAction);

function renderNavigation() {
  treeEl.replaceChildren();
  const surfaces = processor.listSurfaces();
  surfaces.forEach((s, idx) => {
    const li = document.createElement("li");
    const btn = document.createElement("button");
    btn.type = "button";
    const isActive = s.surfaceId === activeSurfaceId;
    btn.className = isActive ? "btn btn--primary a2ui-tree__btn active" : "btn a2ui-tree__btn";
    const label = humanizeSurfaceName(s.surfaceId, idx + 1);
    btn.textContent = s.requiresAttention ? label + " [Atenção]" : label;
    btn.setAttribute("title", s.surfaceId);
    btn.addEventListener("click", () => {
      activeSurfaceId = s.surfaceId;
      renderAll();
    });
    li.appendChild(btn);
    treeEl.appendChild(li);
  });

  breadcrumbEl.replaceChildren();
  const trail = activeSurfaceId ? processor.getBreadcrumb(activeSurfaceId) : [];
  trail.forEach((sid, idx) => {
    if (idx > 0) {
      const sep = document.createElement("span");
      sep.textContent = " / ";
      breadcrumbEl.appendChild(sep);
    }
    const item = document.createElement("button");
    item.type = "button";
    item.className = "btn a2ui-breadcrumb__btn";
    item.textContent = humanizeSurfaceName(sid, idx + 1).replace(/^\d+\.\s*/, "");
    item.setAttribute("title", sid);
    item.addEventListener("click", () => {
      activeSurfaceId = sid;
      renderAll();
    });
    breadcrumbEl.appendChild(item);
  });
}

function renderAll() {
  const surfaces = processor.listSurfaces();
  if (!activeSurfaceId || !processor.getSurface(activeSurfaceId)) {
    const attention = surfaces.find((s) => s.requiresAttention);
    activeSurfaceId = attention ? attention.surfaceId : (surfaces[0] ? surfaces[0].surfaceId : "");
  }
  renderNavigation();
  const currentSurface = activeSurfaceId ? processor.getSurface(activeSurfaceId) : null;
  renderer.render(currentSurface);

  if (currentSurface && currentSurface.requiresAttention) {
    setTimeout(() => {
      const gate = rootEl.querySelector(".lego-stop-gate, [data-kind='LegoStopGatePanel']");
      if (gate) {
        gate.scrollIntoView({ behavior: "smooth", block: "center" });
      }
    }, 40);
  } else {
    setTimeout(() => {
      const ta = rootEl.querySelector(".a2ui-textarea, .a2ui-input");
      if (ta && !rootEl.querySelector(".lego-stop-gate")) {
        ta.focus();
      }
    }, 40);
  }
}

window.addEventListener("keydown", (ev) => {
  const tag = (document.activeElement && document.activeElement.tagName) || "";
  const isInput = tag === "INPUT" || tag === "TEXTAREA" || (document.activeElement && document.activeElement.isContentEditable);
  if (isInput) return;

  if (ev.key === "1") {
    const btn = rootEl.querySelector(".lego-btn--approve");
    if (btn && !btn.disabled) {
      ev.preventDefault();
      btn.click();
    }
  } else if (ev.key === "2") {
    const btn = rootEl.querySelector(".lego-btn--customize");
    if (btn && !btn.disabled) {
      ev.preventDefault();
      btn.click();
    }
  } else if (ev.key === "0") {
    const btn = rootEl.querySelector(".lego-btn--cancel");
    if (btn && !btn.disabled) {
      ev.preventDefault();
      btn.click();
    }
  }
});

processor.subscribe(({ surfaceId }) => {
  hideLoading();
  const target = processor.getSurface(surfaceId);
  if (target && target.requiresAttention) {
    activeSurfaceId = surfaceId;
  }
  renderAll();
});

function handleEnvelopeEvent(ev) {
  try {
    const env = JSON.parse(ev.data);
    const proc = (env && env.processo) || detectProcessFromSurface(env && (env.surfaceId || env.surface_id));
    if (proc) {
      applyProcessTheme(proc);
    }
    if (env && env.data_model && env.data_model.andon_state) {
      renderAndonControls(env.data_model.andon_state);
    }
    processor.process(env);
  } catch (err) {
    console.error("a2ui envelope rejeitado:", err);
    statusEl.textContent = "Envelope ignorado";
  }
}

function connectStream() {
  if (!sessionId) {
    statusEl.textContent = "Sessao ausente";
    renderAll();
    return;
  }
  const streamUrl = "/api/sessao/" + encodeURIComponent(sessionId) + "/stream";
  const source = new EventSource(streamUrl);
  source.onopen = () => {
    statusEl.textContent = "Conectado";
  };
  source.onerror = () => {
    statusEl.textContent = "Reconectando";
  };
  source.addEventListener("a2ui", handleEnvelopeEvent);
  source.onmessage = handleEnvelopeEvent;
}

connectStream();
