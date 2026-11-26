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
    const box = document.createElement("div");
    box.className = "a2ui-andon-paused-box";

    const pill = document.createElement("span");
    pill.className = "pill pill--andon-paused";
    pill.textContent = "Pausado (Andon)";

    const row = document.createElement("div");
    row.className = "a2ui-andon-actions-row";

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

    row.append(btnResume, btnCancel);
    box.append(pill, row);
    andonContainer.appendChild(box);
  } else if (state === "cancelled") {
    const pill = document.createElement("span");
    pill.className = "pill pill--andon-paused";
    pill.style.background = "#ef4444";
    pill.style.borderColor = "#dc2626";
    pill.style.color = "#ffffff";
    pill.style.fontWeight = "700";
    pill.textContent = "Processo Cancelado";
    andonContainer.appendChild(pill);
  } else {
    const btnAndon = document.createElement("button");
    btnAndon.type = "button";
    btnAndon.className = "btn btn--andon btn--andon-sidebar";
    btnAndon.textContent = "Parar";
    btnAndon.title = "Parar execução da esteira (Andon Cord)";
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

  const elapsed = document.createElement("p");
  elapsed.id = "a2ui-loading-elapsed";
  elapsed.className = "a2ui-loading-elapsed";
  elapsed.textContent = "Tempo decorrido: 00:00";
  card.appendChild(elapsed);

  const stepsUl = document.createElement("ul");
  stepsUl.id = "a2ui-loading-steps";
  stepsUl.className = "a2ui-live-steps a2ui-loading-steps";
  stepsUl.hidden = true;
  card.appendChild(stepsUl);

  return card;
}

let loadingWarnTimer = null;
let loadingTimeoutTimer = null;
let loadingElapsedInterval = null;
let loadingStartedAt = 0;

function formatMMSS(ms) {
  const totalSec = Math.max(0, Math.ceil(Number(ms || 0) / 1000));
  const m = String(Math.floor(totalSec / 60)).padStart(2, "0");
  const s = String(totalSec % 60).padStart(2, "0");
  return m + ":" + s;
}

function clearLoadingTimers() {
  if (loadingWarnTimer) {
    clearTimeout(loadingWarnTimer);
    loadingWarnTimer = null;
  }
  if (loadingTimeoutTimer) {
    clearTimeout(loadingTimeoutTimer);
    loadingTimeoutTimer = null;
  }
  if (loadingElapsedInterval) {
    clearInterval(loadingElapsedInterval);
    loadingElapsedInterval = null;
  }
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
  clearLoadingTimers();
  loadingStartedAt = Date.now();
  const elapsedEl = document.getElementById("a2ui-loading-elapsed");
  if (elapsedEl) {
    elapsedEl.textContent = "Tempo decorrido: 00:00";
  }
  loadingElapsedInterval = setInterval(() => {
    const el = document.getElementById("a2ui-loading-elapsed");
    if (el) {
      el.textContent = "Tempo decorrido: " + formatMMSS(Date.now() - loadingStartedAt);
    }
  }, 1000);
  loadingWarnTimer = setTimeout(() => {
    const t = document.getElementById("a2ui-loading-text");
    if (t && !lastLiveStatusMsg) {
      t.textContent = "Resposta recebida pelo servidor. Aguardando o agente publicar a proxima tela...";
    }
  }, 25000);
  loadingTimeoutTimer = setTimeout(() => {
    hideLoading();
    showFeedback("Resposta enviada, mas o agente ainda nao publicou a proxima tela. Verifique a sessao do agente ou tente novamente.", "error");
    statusEl.textContent = "Aguardando agente";
  }, 180000);
}

function hideLoading() {
  clearLoadingTimers();
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
  const surfaces = processor.listSurfaces().filter((s) => !String(s.surfaceId || "").startsWith("andon:"));
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
  const surfaces = processor.listSurfaces().filter((s) => !String(s.surfaceId || "").startsWith("andon:"));
  if (!activeSurfaceId || !processor.getSurface(activeSurfaceId) || activeSurfaceId.startsWith("andon:")) {
    const attention = surfaces.find((s) => s.requiresAttention);
    activeSurfaceId = attention ? attention.surfaceId : (surfaces[0] ? surfaces[0].surfaceId : "");
  }
  renderNavigation();
  const currentSurface = activeSurfaceId ? processor.getSurface(activeSurfaceId) : null;
  renderer.render(currentSurface);
  if (awaitingResponse) {
    setSurfaceDisabled(true);
  }

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
  if (surfaceId && !surfaceId.startsWith("andon:")) {
    const target = processor.getSurface(surfaceId);
    if (target && target.requiresAttention) {
      activeSurfaceId = surfaceId;
    }
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
    } else if (env && env.surfaceId && env.surfaceId.startsWith("andon:") && env.data_model && env.data_model.state) {
      renderAndonControls(env.data_model.state);
    }
    processor.process(env);
  } catch (err) {
    console.error("a2ui envelope rejeitado:", err);
    const motivo = err && err.message ? err.message : "erro";
    statusEl.textContent = "Envelope ignorado: " + motivo;
  }
}

let reconnectDelayMs = 1000;
let reconnectTimer = null;
let currentSource = null;

function connectStream() {
  if (!sessionId) {
    statusEl.textContent = "Sessao ausente";
    renderAll();
    return;
  }
  if (currentSource) {
    currentSource.close();
  }
  const streamUrl = "/api/sessao/" + encodeURIComponent(sessionId) + "/stream";
  currentSource = new EventSource(streamUrl);
  currentSource.onopen = () => {
    reconnectDelayMs = 1000;
    statusEl.textContent = "Conectado";
  };
  currentSource.onerror = () => {
    if (currentSource) currentSource.close();
    const segs = Math.max(1, Math.round(reconnectDelayMs / 1000));
    statusEl.textContent = `Reconectando (${segs}s)...`;
    clearTimeout(reconnectTimer);
    reconnectTimer = setTimeout(() => {
      reconnectDelayMs = Math.min(Math.round(reconnectDelayMs * 1.5), 15000);
      connectStream();
    }, reconnectDelayMs);
  };
  currentSource.addEventListener("a2ui", handleEnvelopeEvent);
  currentSource.addEventListener("timer", handleTimerEvent);
  currentSource.addEventListener("status", handleStatusEvent);
  currentSource.onmessage = handleEnvelopeEvent;
}

const timerSelectEl = document.getElementById("a2ui-timer-select");
const timerDisplayEl = document.getElementById("a2ui-timer-display");
const timerPauseBtn = document.getElementById("a2ui-timer-pause");
const timerExtendBtn = document.getElementById("a2ui-timer-extend");
const liveProgressEl = document.getElementById("a2ui-live-progress");

let timerState = { timeout_ms: 300000, remaining_ms: 300000, paused: false, active: false, expired: false };
let timerTickAt = Date.now();
let lastLiveStatusMsg = "";

async function postTimerAction(acao, timeoutMs) {
  if (!sessionId) return;
  const payload = timeoutMs ? { acao, timeout_ms: timeoutMs } : { acao };
  try {
    const res = await fetch("/api/sessao/" + encodeURIComponent(sessionId) + "/timer", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (res.ok) {
      applyTimerSnapshot(await res.json());
    }
  } catch (_err) {
    // Mantem resiliencia caso servidor reinicie
  }
}

function applyTimerSnapshot(snap) {
  if (!snap) return;
  timerState = { ...timerState, ...snap };
  timerTickAt = Date.now();
  if (timerSelectEl && snap.timeout_ms) {
    const val = String(snap.timeout_ms);
    if ([...timerSelectEl.options].some((o) => o.value === val)) {
      timerSelectEl.value = val;
    }
  }
  renderTimerUI();
  if (timerState.expired) {
    hideLoading();
    showFeedback("Tempo limite da etapa atingido. Retome o processo pela via original (Chat, CLI ou Studio).", "error");
    statusEl.textContent = "Tempo esgotado";
    rootEl.querySelectorAll("button, input, textarea, select").forEach((el) => { el.disabled = true; });
  }
}

function renderTimerUI() {
  if (!timerDisplayEl) return;
  let rem = Number(timerState.remaining_ms || 0);
  if (timerState.active && !timerState.paused && !timerState.expired) {
    rem = Math.max(0, rem - (Date.now() - timerTickAt));
  }
  timerDisplayEl.classList.remove("a2ui-timer-display--paused", "a2ui-timer-display--expired");
  if (timerState.expired) {
    timerDisplayEl.textContent = "00:00";
    timerDisplayEl.classList.add("a2ui-timer-display--expired");
  } else if (timerState.paused) {
    timerDisplayEl.textContent = "PAUSA " + formatMMSS(rem);
    timerDisplayEl.classList.add("a2ui-timer-display--paused");
  } else {
    timerDisplayEl.textContent = formatMMSS(rem);
  }
  if (timerPauseBtn) {
    timerPauseBtn.textContent = timerState.paused ? "Retomar" : "Pausar";
    timerPauseBtn.disabled = !timerState.active || timerState.expired;
  }
  if (timerExtendBtn) {
    timerExtendBtn.disabled = timerState.expired;
  }
}

function handleTimerEvent(ev) {
  try {
    applyTimerSnapshot(JSON.parse(ev.data));
  } catch (_err) {}
}

let awaitingResponse = false;

function setSurfaceDisabled(disabled) {
  if (!rootEl) return;
  rootEl.querySelectorAll("button, input, textarea, select").forEach((el) => {
    el.disabled = Boolean(disabled);
  });
}

function renderLiveProgress(st) {
  if (!liveProgressEl || !st) return;
  const etapas = Array.isArray(st.etapas) ? st.etapas : [];
  awaitingResponse = Boolean(st.aguardando_resposta);
  if (awaitingResponse) {
    setSurfaceDisabled(true);
    if (typeof showLoading === "function") {
      showLoading(st.mensagem || "Resposta já enviada. Processo em execução...");
    }
  } else {
    setSurfaceDisabled(false);
    if (st.etapa === "interacao" && !st.mensagem && etapas.length === 0 && typeof hideLoading === "function") {
      hideLoading();
    }
  }
  if (st.mensagem) {
    lastLiveStatusMsg = st.mensagem;
    const t = document.getElementById("a2ui-loading-text");
    if (t) {
      t.textContent = st.mensagem;
    }
  }
  if (etapas.length === 0 && !st.mensagem) {
    liveProgressEl.hidden = true;
    return;
  }
  liveProgressEl.hidden = false;
  liveProgressEl.replaceChildren();

  const header = document.createElement("div");
  header.className = "a2ui-live-progress__header";
  const title = document.createElement("h3");
  title.className = "a2ui-live-progress__title";
  const concluidas = etapas.filter((e) => e.status === "done").length;
  title.textContent = etapas.length > 0
    ? "Execucao ao vivo (" + concluidas + "/" + etapas.length + " etapas)"
    : "Atividade em andamento";
  header.appendChild(title);
  if (st.etapa) {
    const pill = document.createElement("span");
    pill.className = "pill";
    pill.textContent = st.etapa;
    header.appendChild(pill);
  }
  liveProgressEl.appendChild(header);

  if (st.mensagem) {
    const p = document.createElement("p");
    p.className = "a2ui-live-progress__msg";
    p.textContent = st.mensagem;
    liveProgressEl.appendChild(p);
  }
  if (etapas.length > 0) {
    if (!awaitingResponse && typeof hideLoading === "function") {
      hideLoading();
    }
    const ul = document.createElement("ul");
    ul.className = "a2ui-live-steps";
    const loadingStepsEl = document.getElementById("a2ui-loading-steps");
    if (loadingStepsEl) {
      loadingStepsEl.replaceChildren();
      loadingStepsEl.hidden = false;
    }
    etapas.forEach((et) => {
      const li = document.createElement("li");
      const kind = et.status || "pending";
      li.className = "a2ui-live-step a2ui-live-step--" + kind;
      const icon = kind === "done" ? "[OK]" : (kind === "running" ? "[...]" : (kind === "error" ? "[ERRO]" : "[ ]"));
      li.textContent = icon + " " + (et.titulo || et.id || "") + (et.detalhe ? " — " + et.detalhe : "");
      ul.appendChild(li);
      if (loadingStepsEl) {
        loadingStepsEl.appendChild(li.cloneNode(true));
      }
    });
    liveProgressEl.appendChild(ul);
  } else {
    const loadingStepsEl = document.getElementById("a2ui-loading-steps");
    if (loadingStepsEl) {
      loadingStepsEl.replaceChildren();
      loadingStepsEl.hidden = true;
    }
  }
}

function handleStatusEvent(ev) {
  try {
    renderLiveProgress(JSON.parse(ev.data));
  } catch (_err) {}
}

if (timerSelectEl) {
  timerSelectEl.addEventListener("change", () => {
    const ms = parseInt(timerSelectEl.value, 10);
    if (ms > 0) postTimerAction("configurar", ms);
  });
}
if (timerPauseBtn) {
  timerPauseBtn.addEventListener("click", () => {
    postTimerAction(timerState.paused ? "retomar" : "pausar");
  });
}
if (timerExtendBtn) {
  timerExtendBtn.addEventListener("click", () => {
    postTimerAction("estender", 300000);
  });
}
setInterval(renderTimerUI, 1000);

window.addEventListener("keydown", (e) => {
  if (e.altKey && (e.key === "p" || e.key === "P")) {
    e.preventDefault();
    if (!andonContainer) return;
    const btnResume = andonContainer.querySelector(".btn--andon-resume");
    if (btnResume) {
      btnResume.click();
      return;
    }
    const btnPause = andonContainer.querySelector(".btn--andon");
    if (btnPause) {
      btnPause.click();
    }
  }
});

connectStream();
