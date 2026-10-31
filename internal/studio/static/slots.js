"use strict";

// Slots da peça (gates extras, scripts e hooks pré/pós) e vizinhos sugeridos. A parte pura não
// toca o DOM e é testada com Node. Slots são opcionais: vazio é uma escolha válida, nada é
// preenchido sozinho, e a sugestão só vira slot quando o usuário escolhe.

function emptySlots() {
  return { gates: [], scripts: [], pre: [], post: [] };
}

function slotCount(s) {
  return s.gates.length + s.scripts.length + s.pre.length + s.post.length;
}

// slotsToComposition serializa para o servidor; devolve undefined quando não há nenhum slot.
// Só ref ou command vão ao servidor (o rótulo é só da interface).
function slotsToComposition(s) {
  const cmd = (it) => (it.ref ? { ref: it.ref } : { command: it.command });
  const out = {};
  if (s.gates.length) out.gates = [...s.gates];
  if (s.scripts.length) out.scripts = s.scripts.map(cmd);
  const hooks = {};
  if (s.pre.length) hooks.pre = s.pre.map(cmd);
  if (s.post.length) hooks.post = s.post.map(cmd);
  if (Object.keys(hooks).length) out.hooks = hooks;
  return Object.keys(out).length ? out : undefined;
}

// Eventos "Pre*" e de início valem só antes da peça; "Post*" e de fim, só depois; o resto, nos dois.
// Espelha suggest.hookSlots (Go).
function hookSlotsFor(event) {
  switch (event) {
    case "PreToolUse": case "SessionStart": case "UserPromptSubmit": case "PreCompact":
      return { pre: true, post: false };
    case "PostToolUse": case "Stop": case "SubagentStop": case "SessionEnd":
      return { pre: false, post: true };
    default:
      return { pre: true, post: true };
  }
}

// instantiateGate transforma um gate embutido numa expressão para o artefato; null se não houver
// como fazê-lo com segurança. has_pattern devolve um modelo com regex vazia, para o usuário completar.
function instantiateGate(name, artifact) {
  if (!artifact || /['\\\n\r]/.test(artifact)) return null;
  switch (name) {
    case "file_exists": return `file_exists('${artifact}')`;
    case "word_count": return `word_count('${artifact}') > 100`;
    case "line_count": return `line_count('${artifact}') > 10`;
    case "has_binary_criteria": return `has_binary_criteria('${artifact}')`;
    case "has_acceptance_criteria": return `has_acceptance_criteria('${artifact}')`;
    case "has_pattern": return `has_pattern('${artifact}', '')`;
    default: return null;
  }
}

// ---------- DOM ----------

const SLOT_GROUPS = [
  { id: "gates", title: "Gates extras", placeholder: "ex.: word_count('ARQ.md') > 50", note: "Além do gate de saída; todos precisam passar." },
  { id: "scripts", title: "Scripts", placeholder: "comando e argumentos (sem pipes nem redirecionamentos)" },
  { id: "pre", title: "Hooks pré", placeholder: "comando executado antes da peça" },
  { id: "post", title: "Hooks pós", placeholder: "comando executado depois da peça" },
];

function slotEl(tag, cls, text) {
  const e = document.createElement(tag);
  if (cls) e.className = cls;
  if (text !== undefined) e.textContent = text;
  return e;
}

function slotItemLabel(group, item) {
  return group.id === "gates" ? item : item.label || item.command || item.ref;
}

function slotSame(group, a, b) {
  return group.id === "gates" ? a === b : (a.ref || a.command) === (b.ref || b.command);
}

// slotSuggestions normaliza a resposta de /api/suggest em {label, value} por grupo.
function slotSuggestions(group, sug) {
  if (!sug) return null;
  if (group.id === "gates") return sug.gates.map((g) => ({ label: `${g.label}: ${g.expr}`, value: g.expr }));
  const pick = (items) => items.map((i) => ({ label: `${i.name}${i.description ? ` — ${i.description}` : ""}`, value: { ref: i.id, label: i.name } }));
  if (group.id === "scripts") return pick(sug.scripts);
  return pick(group.id === "pre" ? sug.hooks.pre : sug.hooks.post);
}

function slotChips(node, group, onChange) {
  const ul = slotEl("ul", "chips");
  node.slots[group.id].forEach((item, i) => {
    const li = slotEl("li", "chip");
    li.append(slotEl("span", "chip__text", slotItemLabel(group, item)));
    const x = slotEl("button", "chip__x", "×");
    x.type = "button";
    x.setAttribute("aria-label", `Remover ${slotItemLabel(group, item)}`);
    x.addEventListener("click", () => { node.slots[group.id].splice(i, 1); onChange(); });
    li.append(x);
    ul.append(li);
  });
  return ul;
}

function slotSuggestBox(node, group, sug, onChange) {
  const list = slotSuggestions(group, sug);
  const box = slotEl("details", "suggest");
  box.append(slotEl("summary", "suggest__title", list === null ? "Sugeridos" : `Sugeridos (${list.length})`));
  if (list === null) {
    box.append(slotEl("p", "note", "Carregando sugestões…"));
    return box;
  }
  if (list.length === 0) {
    box.append(slotEl("p", "note", "Nenhum sugerido para esta peça. Você pode digitar o seu abaixo."));
    return box;
  }
  const ul = slotEl("ul", "suggest__list");
  for (const s of list) {
    const taken = node.slots[group.id].some((x) => slotSame(group, x, s.value));
    const b = slotEl("button", "suggest__item", s.label);
    b.type = "button";
    b.disabled = taken;
    b.title = taken ? "Já adicionado" : "Adicionar ao slot";
    b.addEventListener("click", () => { node.slots[group.id].push(s.value); onChange(); });
    const li = slotEl("li");
    li.append(b);
    ul.append(li);
  }
  box.append(ul);
  return box;
}

function slotCustomRow(node, group, onChange, draft) {
  const row = slotEl("div", "slotadd");
  const input = slotEl("input", "input");
  if (draft) {
    input.value = draft;
    setTimeout(() => input.focus(), 0);
  }
  input.type = "text";
  input.placeholder = group.placeholder;
  input.spellcheck = false;
  input.autocomplete = "off";
  input.setAttribute("aria-label", `${group.title}: digitar o seu`);
  const add = slotEl("button", "btn", "Adicionar");
  add.type = "button";
  const commit = () => {
    const v = input.value.trim();
    if (!v) return;
    const value = group.id === "gates" ? v : { command: v, label: v };
    if (!node.slots[group.id].some((x) => slotSame(group, x, value))) node.slots[group.id].push(value);
    onChange();
  };
  add.addEventListener("click", commit);
  input.addEventListener("keydown", (e) => { if (e.key === "Enter") { e.preventDefault(); commit(); } });
  row.append(input, add);
  return row;
}

// renderSlots monta o editor dos três encaixes. sug é a resposta de /api/suggest (ou null enquanto carrega).
// draft é uma expressão de gate já preenchida para o usuário completar (ex.: has_pattern com regex vazia).
function renderSlots(node, sug, onChange, draft) {
  const sec = slotEl("section", "slots");
  sec.append(slotEl("h3", "subtitle", "Slots (opcionais)"));
  sec.append(slotEl("p", "note", "Deixe em branco se não precisar. Nada é preenchido sozinho."));
  for (const g of SLOT_GROUPS) {
    const wrap = slotEl("div", "slotgroup");
    const head = slotEl("div", "slotgroup__head");
    head.append(slotEl("span", "slotgroup__title", g.title), slotEl("span", "badge badge--muted", String(node.slots[g.id].length)));
    wrap.append(head);
    if (g.note) wrap.append(slotEl("p", "note", g.note));
    wrap.append(slotChips(node, g, onChange), slotSuggestBox(node, g, sug, onChange), slotCustomRow(node, g, onChange, g.id === "gates" ? draft : ""));
    sec.append(wrap);
  }
  if (sug && sug.hooks.skipped > 0) {
    sec.append(slotEl("p", "note", `${sug.hooks.skipped} hook(s) não aparecem: exigem shell ou trazem credencial no comando.`));
  }
  sec.append(slotEl("p", "note note--warn", "Os slots são gravados no pipeline.dag.json, mas o runner ainda não os executa."));
  return sec;
}

// renderNeighbors lista as peças sugeridas para ligar depois desta; as não instaladas ficam desabilitadas.
function renderNeighbors(sug, onAdd) {
  const sec = slotEl("section", "neighbors");
  sec.append(slotEl("h3", "subtitle", "Vizinhos sugeridos"));
  if (!sug || sug.neighbors.length === 0) {
    sec.append(slotEl("p", "note", sug ? "Nenhuma peça compatível encontrada." : "Carregando sugestões…"));
    return sec;
  }
  const ul = slotEl("ul", "chips");
  for (const n of sug.neighbors) {
    const li = slotEl("li");
    const b = slotEl("button", "chip chip--action");
    b.type = "button";
    b.append(slotEl("span", "chip__text", n.name));
    b.append(slotEl("span", "badge badge--muted", n.kind === "mcp" ? "MCP" : n.installed ? "skill" : "não instalada"));
    b.disabled = !n.installed;
    b.title = n.installed ? "Adicionar e ligar" : "Esta skill não está instalada nesta máquina";
    b.addEventListener("click", () => onAdd(n));
    li.append(b);
    ul.append(li);
  }
  sec.append(ul);
  return sec;
}
