// Bundle offline do protocolo A2UI v0.9 para o LegoLoop Studio.
// Zero dependencias externas, zero HTML bruto e zero estilos inline.

const SURFACE_ID_RE = /^[a-zA-Z0-9._:/-]{1,128}$/;
const HOSTILE_RE = /<\s*script\b|javascript:|on\w+\s*=|<\s*iframe\b/i;

export const DEFAULT_CATALOG = Object.freeze([
  "Text",
  "Heading",
  "Button",
  "Card",
  "Column",
  "Row",
  "List",
  "Input",
  "Textarea",
  "Select",
  "Badge",
  "Divider",
  "Form",
  "LegoDagGraph",
  "LegoProgressMeter",
  "LegoStopGateCard",
  "LegoArtifactViewer",
  "LegoLogStream",
  "LegoSwarmBoard"
]);

const customRenderers = new Map();

export function registerComponent(kind, renderer) {
  if (typeof kind !== "string" || !kind.trim() || typeof renderer !== "function") {
    throw new Error("a2ui: registro de componente invalido");
  }
  customRenderers.set(kind, renderer);
}

export function isSafeValue(value) {
  if (typeof value === "string") {
    return !HOSTILE_RE.test(value);
  }
  if (Array.isArray(value)) {
    return value.every(isSafeValue);
  }
  if (value && typeof value === "object") {
    return Object.entries(value).every(([k, v]) => isSafeValue(k) && isSafeValue(v));
  }
  return true;
}

export class MessageProcessor {
  constructor(allowedCatalog = DEFAULT_CATALOG) {
    this.allowed = new Set(allowedCatalog);
    this.surfaces = new Map();
    this.listeners = new Set();
  }

  allowComponent(kind) {
    if (typeof kind === "string" && kind.trim()) {
      this.allowed.add(kind);
    }
  }

  subscribe(listener) {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  _notify(eventType, surfaceId) {
    for (const fn of this.listeners) {
      fn({ eventType, surfaceId, surfaces: this.listSurfaces() });
    }
  }

  _validateEnvelope(env) {
    if (!env || typeof env !== "object") {
      throw new Error("a2ui: envelope invalido");
    }
    const sid = env.surface_id || env.surfaceId || "";
    if (!SURFACE_ID_RE.test(sid)) {
      throw new Error("a2ui: surface_id invalido");
    }
    const components = Array.isArray(env.components) ? env.components : [];
    for (const comp of components) {
      const kind = comp && comp.kind ? comp.kind : "";
      if (!this.allowed.has(kind) && !customRenderers.has(kind)) {
        throw new Error("a2ui: componente fora do catalogo: " + kind);
      }
      if (!isSafeValue(comp.props || {})) {
        throw new Error("a2ui: payload hostil detectado");
      }
    }
    if (env.data_model && !isSafeValue(env.data_model)) {
      throw new Error("a2ui: data_model hostil detectado");
    }
    return sid;
  }

  process(envelope) {
    const surfaceId = this._validateEnvelope(envelope);
    const msgType = envelope.message_type || envelope.messageType || "";
    const parentId = envelope.parent_surface_id || envelope.parentSurfaceId || "";

    if (msgType === "deleteSurface") {
      this.surfaces.delete(surfaceId);
      this._notify(msgType, surfaceId);
      return null;
    }

    const existing = this.surfaces.get(surfaceId) || {
      surfaceId,
      parentId,
      sessionId: envelope.session_id || envelope.sessionId || "",
      protocolVersion: envelope.protocol_version || envelope.protocolVersion || "v0.9",
      components: [],
      dataModel: {},
      fallbackText: "",
      requiresAttention: false,
      updatedAt: Date.now()
    };

    if (parentId) {
      existing.parentId = parentId;
    }
    if (typeof envelope.fallback_text === "string" && envelope.fallback_text) {
      existing.fallbackText = envelope.fallback_text;
    }

    if (msgType === "createSurface") {
      existing.components = Array.isArray(envelope.components) ? envelope.components.slice() : [];
      existing.dataModel = Object.assign({}, envelope.data_model || {});
    } else if (msgType === "updateComponents") {
      existing.components = Array.isArray(envelope.components) ? envelope.components.slice() : existing.components;
    } else if (msgType === "updateDataModel") {
      existing.dataModel = Object.assign({}, existing.dataModel, envelope.data_model || {});
    } else {
      throw new Error("a2ui: tipo de mensagem nao suportado: " + msgType);
    }

    existing.requiresAttention = Boolean(
      envelope.requires_attention ||
      existing.components.some((c) => c && c.kind === "LegoStopGateCard")
    );
    existing.updatedAt = Date.now();
    this.surfaces.set(surfaceId, existing);
    this._notify(msgType, surfaceId);
    return existing;
  }

  getSurface(surfaceId) {
    return this.surfaces.get(surfaceId) || null;
  }

  listSurfaces() {
    return Array.from(this.surfaces.values());
  }

  getBreadcrumb(surfaceId) {
    const trail = [];
    const seen = new Set();
    let curr = this.surfaces.get(surfaceId);
    while (curr && !seen.has(curr.surfaceId)) {
      seen.add(curr.surfaceId);
      trail.unshift(curr.surfaceId);
      curr = curr.parentId ? this.surfaces.get(curr.parentId) : null;
    }
    return trail;
  }
}

function el(tag, className, text) {
  const node = document.createElement(tag);
  if (className) {
    node.className = className;
  }
  if (text !== undefined && text !== null) {
    node.textContent = String(text);
  }
  return node;
}

export class A2UISurfaceRenderer {
  constructor(container, onAction) {
    this.container = container;
    this.onAction = typeof onAction === "function" ? onAction : () => {};
  }

  render(surface) {
    this.container.replaceChildren();
    if (!surface) {
      this.container.appendChild(el("p", "a2ui-empty", "Aguardando superficie A2UI da sessao."));
      return;
    }
    const section = el("section", "a2ui-surface");
    section.setAttribute("data-surface-id", surface.surfaceId);
    const header = el("header", "a2ui-surface__header");
    header.appendChild(el("h2", "a2ui-surface__title", surface.surfaceId));
    if (surface.requiresAttention) {
      header.appendChild(el("span", "pill a2ui-badge--attention", "Acao pendente"));
    }
    section.appendChild(header);

    const body = el("div", "a2ui-surface__body");
    const comps = Array.isArray(surface.components) ? surface.components : [];
    if (comps.length === 0 && surface.fallbackText) {
      body.appendChild(el("pre", "a2ui-fallback", surface.fallbackText));
    } else {
      for (const comp of comps) {
        body.appendChild(this._renderComponent(surface, comp));
      }
    }
    section.appendChild(body);
    this.container.appendChild(section);
  }

  _emit(surface, name, payload) {
    this.onAction({
      session_id: surface.sessionId,
      surface_id: surface.surfaceId,
      origin: "ui",
      name,
      payload: payload || {}
    });
  }

  _renderComponent(surface, comp) {
    const props = comp.props || {};
    if (customRenderers.has(comp.kind)) {
      return customRenderers.get(comp.kind)({
        surface,
        component: comp,
        props,
        dataModel: surface.dataModel,
        emitAction: (name, payload) => this._emit(surface, name, payload)
      });
    }
    switch (comp.kind) {
      case "Heading":
        return el("h3", "a2ui-heading", props.text || props.title || comp.id);
      case "Text":
        return el("p", "a2ui-text", props.text || props.value || "");
      case "Badge":
        return el("span", "pill a2ui-badge", props.label || props.text || "");
      case "Divider":
        return el("hr", "a2ui-divider");
      case "Button":
        return this._renderButton(surface, comp, props);
      case "Input":
      case "Textarea":
      case "Form":
        return this._renderFormControl(surface, comp, props);
      case "List":
        return this._renderList(props);
      case "LegoStopGateCard":
        return this._renderStopGate(surface, comp, props);
      case "LegoProgressMeter":
        return this._renderProgress(props);
      case "LegoArtifactViewer":
        return this._renderArtifact(props);
      case "LegoLogStream":
        return this._renderLogStream(props);
      case "LegoDagGraph":
      case "LegoSwarmBoard":
      case "Card":
      case "Column":
      case "Row":
      default:
        return this._renderCard(comp, props);
    }
  }

  _renderButton(surface, comp, props) {
    const btn = el("button", "btn a2ui-btn", props.label || props.text || comp.id);
    btn.type = "button";
    btn.addEventListener("click", () => {
      this._emit(surface, props.action || comp.id || "click", props.payload || { id: comp.id });
    });
    return btn;
  }

  _renderFormControl(surface, comp, props) {
    const form = el("form", "a2ui-form");
    const label = el("label", "a2ui-label", props.label || props.question || comp.id);
    const fieldName = props.name || comp.id || "answer";
    const field = comp.kind === "Input" ? el("input", "a2ui-input") : el("textarea", "a2ui-textarea");
    field.name = fieldName;
    if (props.placeholder) {
      field.setAttribute("placeholder", String(props.placeholder));
    }
    if (props.value !== undefined) {
      field.value = String(props.value);
    }
    const submit = el("button", "btn btn--primary a2ui-submit", props.submitLabel || "Enviar");
    submit.type = "submit";
    form.append(label, field, submit);
    form.addEventListener("submit", (ev) => {
      ev.preventDefault();
      this._emit(surface, props.action || "submit_answer", { [fieldName]: field.value });
    });
    return form;
  }

  _renderList(props) {
    const ul = el("ul", "a2ui-list");
    const items = Array.isArray(props.items) ? props.items : [];
    for (const item of items) {
      ul.appendChild(el("li", "a2ui-list__item", typeof item === "string" ? item : JSON.stringify(item)));
    }
    return ul;
  }

  _renderStopGate(surface, comp, props) {
    const card = el("div", "a2ui-card a2ui-stop-gate");
    card.appendChild(el("h3", "a2ui-card__title", props.title || "Stop-Gate"));
    if (props.summary) {
      card.appendChild(el("p", "a2ui-text", props.summary));
    }
    const actions = el("div", "a2ui-actions");
    const choices = [
      { id: "1", label: "[1] Seguir", decision: "aprovar" },
      { id: "2", label: "[2] Personalizar", decision: "personalizar" },
      { id: "0", label: "[0] Cancelar", decision: "cancelar" }
    ];
    for (const ch of choices) {
      const btn = el("button", "btn a2ui-btn", ch.label);
      btn.type = "button";
      btn.addEventListener("click", () => {
        this._emit(surface, props.action || "stop_gate_decision", {
          gate_id: comp.id,
          choice: ch.id,
          decision: ch.decision
        });
      });
      actions.appendChild(btn);
    }
    card.appendChild(actions);
    return card;
  }

  _renderProgress(props) {
    const box = el("div", "a2ui-card a2ui-progress");
    const pct = Number(props.percent ?? props.confidence ?? 0);
    box.appendChild(el("div", "a2ui-progress__label", (props.label || "Progresso") + ": " + pct + "%"));
    const bar = el("progress", "a2ui-progress__bar");
    bar.max = 100;
    bar.value = Number.isFinite(pct) ? Math.max(0, Math.min(100, pct)) : 0;
    box.appendChild(bar);
    return box;
  }

  _renderArtifact(props) {
    const box = el("div", "a2ui-card a2ui-artifact");
    box.appendChild(el("h4", "a2ui-card__title", props.path || props.title || "Artefato"));
    box.appendChild(el("pre", "a2ui-code", props.content || ""));
    return box;
  }

  _renderLogStream(props) {
    const box = el("div", "a2ui-card a2ui-logs");
    const lines = Array.isArray(props.lines) ? props.lines : [];
    box.appendChild(el("pre", "a2ui-code", lines.join("\n")));
    return box;
  }

  _renderCard(comp, props) {
    const card = el("div", "a2ui-card");
    card.setAttribute("data-kind", comp.kind);
    card.appendChild(el("h4", "a2ui-card__title", props.title || comp.kind));
    if (props.description || props.text) {
      card.appendChild(el("p", "a2ui-text", props.description || props.text));
    }
    return card;
  }
}
