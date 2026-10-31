"use strict";

// Componentes de feedback do Studio (Dialog e Toast). Tudo é DOM da própria página:
// nada de caixas do navegador/SO. Todo texto entra via textContent.

const TOAST_MS = 2600;
const TOAST_ERROR_MS = 5200;
const TOAST_MAX = 4;
const TOAST_FADE_MS = 200;

function toast(msg, isError = false) {
  const box = document.getElementById("toasts");
  const t = document.createElement("div");
  t.className = isError ? "toast is-error" : "toast";
  t.setAttribute("role", isError ? "alert" : "status");
  t.textContent = msg;
  box.append(t);
  while (box.children.length > TOAST_MAX) box.firstElementChild.remove();
  requestAnimationFrame(() => t.classList.add("is-visible"));
  setTimeout(() => {
    t.classList.remove("is-visible");
    setTimeout(() => t.remove(), TOAST_FADE_MS);
  }, isError ? TOAST_ERROR_MS : TOAST_MS);
}

function dialogButton(label, cls, onClick) {
  const b = document.createElement("button");
  b.type = "button";
  b.className = cls;
  b.textContent = label;
  b.addEventListener("click", onClick);
  return b;
}

// Abre o diálogo de confirmação e resolve true (confirmou) ou false (cancelou, Esc ou clique fora).
// O foco inicial fica em "Cancelar": a ação destrutiva exige escolha deliberada.
function confirmDialog({ title, message, confirmLabel = "Confirmar", cancelLabel = "Cancelar", danger = false }) {
  const dlg = document.getElementById("dialog");
  if (dlg.open) dlg.close();
  return new Promise((resolve) => {
    // Resolve direto nos handlers (não depende do evento assíncrono "close") e remove todos
    // os listeners ao terminar, para uma chamada nunca herdar estado da anterior.
    const off = new AbortController();
    const finish = (value) => {
      off.abort();
      if (dlg.open) dlg.close();
      resolve(value);
    };
    const head = document.createElement("h2");
    head.className = "dialog__title";
    head.id = "dialog-title";
    head.textContent = title;
    const body = document.createElement("p");
    body.className = "dialog__text";
    body.id = "dialog-text";
    body.textContent = message;
    const actions = document.createElement("div");
    actions.className = "dialog__actions";
    const cancel = dialogButton(cancelLabel, "btn", () => finish(false));
    const ok = dialogButton(confirmLabel, danger ? "btn btn--danger" : "btn btn--primary", () => finish(true));
    actions.append(cancel, ok);
    dlg.setAttribute("aria-labelledby", "dialog-title");
    dlg.setAttribute("aria-describedby", "dialog-text");
    dlg.replaceChildren(head, body, actions);
    const opts = { signal: off.signal };
    dlg.addEventListener("cancel", () => finish(false), opts);
    dlg.addEventListener("close", () => finish(false), opts);
    dlg.addEventListener("click", (e) => { if (e.target === dlg) finish(false); }, opts);
    dlg.showModal();
    cancel.focus();
  });
}

// Diálogo de escolha entre várias opções. Resolve o valor escolhido, ou null se o usuário cancelar
// (Esc, clique fora ou Cancelar).
function choiceDialog({ title, message, choices, cancelLabel = "Cancelar" }) {
  const dlg = document.getElementById("dialog");
  if (dlg.open) dlg.close();
  return new Promise((resolve) => {
    const off = new AbortController();
    const finish = (value) => {
      off.abort();
      if (dlg.open) dlg.close();
      resolve(value);
    };
    const head = document.createElement("h2");
    head.className = "dialog__title";
    head.id = "dialog-title";
    head.textContent = title;
    const body = document.createElement("p");
    body.className = "dialog__text";
    body.id = "dialog-text";
    body.textContent = message;
    const actions = document.createElement("div");
    actions.className = "dialog__actions";
    const cancel = dialogButton(cancelLabel, "btn", () => finish(null));
    actions.append(cancel, ...choices.map((c) => dialogButton(c.label, "btn btn--primary", () => finish(c.value))));
    dlg.setAttribute("aria-labelledby", "dialog-title");
    dlg.setAttribute("aria-describedby", "dialog-text");
    dlg.replaceChildren(head, body, actions);
    const opts = { signal: off.signal };
    dlg.addEventListener("cancel", () => finish(null), opts);
    dlg.addEventListener("close", () => finish(null), opts);
    dlg.addEventListener("click", (e) => { if (e.target === dlg) finish(null); }, opts);
    dlg.showModal();
    cancel.focus();
  });
}
