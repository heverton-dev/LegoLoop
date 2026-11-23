import { MessageProcessor, A2UISurfaceRenderer } from "/components/a2ui-bundle.js";
import { registerLegoLoopCatalog } from "/components/legoloop-catalog.js";

registerLegoLoopCatalog();

const parts = window.location.pathname.split("/").filter(Boolean);
const sessionId = decodeURIComponent(parts[parts.length - 1] || "");
const statusEl = document.getElementById("a2ui-conn-status");
const sessionEl = document.getElementById("a2ui-session-id");
const treeEl = document.getElementById("a2ui-tree");
const breadcrumbEl = document.getElementById("a2ui-breadcrumb");
const rootEl = document.getElementById("a2ui-root");

sessionEl.textContent = sessionId || "sem-sessao";

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

function createLoadingCard() {
  const card = document.createElement("div");
  card.className = "a2ui-loading-card";

  const svgNS = "http://www.w3.org/2000/svg";
  const svg = document.createElementNS(svgNS, "svg");
  svg.setAttribute("class", "a2ui-loading-spinner");
  svg.setAttribute("viewBox", "0 0 50 50");
  svg.setAttribute("fill", "none");

  const rects = [
    { x: "5", y: "5", fill: "var(--primary, #3b82f6)", begin: "0s" },
    { x: "27", y: "5", fill: "var(--accent, #10b981)", begin: "0.2s" },
    { x: "27", y: "27", fill: "var(--primary, #3b82f6)", begin: "0.4s" },
    { x: "5", y: "27", fill: "var(--accent, #10b981)", begin: "0.6s" }
  ];
  for (const r of rects) {
    const rect = document.createElementNS(svgNS, "rect");
    rect.setAttribute("x", r.x);
    rect.setAttribute("y", r.y);
    rect.setAttribute("width", "18");
    rect.setAttribute("height", "18");
    rect.setAttribute("rx", "4");
    rect.setAttribute("fill", r.fill);
    const anim = document.createElementNS(svgNS, "animate");
    anim.setAttribute("attributeName", "opacity");
    anim.setAttribute("values", "1;0.2;1");
    anim.setAttribute("dur", "1s");
    anim.setAttribute("repeatCount", "indefinite");
    anim.setAttribute("begin", r.begin);
    rect.appendChild(anim);
    svg.appendChild(rect);
  }
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
  for (const s of surfaces) {
    const li = document.createElement("li");
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = s.surfaceId === activeSurfaceId ? "btn btn--primary" : "btn";
    btn.textContent = s.requiresAttention ? s.surfaceId + " [Stop-Gate]" : s.surfaceId;
    btn.addEventListener("click", () => {
      activeSurfaceId = s.surfaceId;
      renderAll();
    });
    li.appendChild(btn);
    treeEl.appendChild(li);
  }

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
    item.className = "btn";
    item.textContent = sid;
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
  renderer.render(activeSurfaceId ? processor.getSurface(activeSurfaceId) : null);
}

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
