"use strict";

// Inspetor do nó e da conexão: roteamento, modo de execução, artefato e gate das peças genéricas,
// slots opcionais e vizinhos sugeridos. Depende do estado e das funções de app.js (carregado depois;
// tudo aqui só roda em eventos, quando app.js já definiu `state`).

const EMPTY_SUGGEST = { gates: [], scripts: [], hooks: { pre: [], post: [], skipped: 0 }, neighbors: [] };

function kv(k, v) {
  const box = el("div", "kv");
  box.append(el("span", "kv__k", k), el("span", "kv__v", v));
  return box;
}

function selectField(label, value, options, onChange) {
  const f = el("label", "field");
  const s = el("select", "input");
  options.forEach((o) => s.append(new Option(o || "—", o, false, o === value)));
  if (!options.includes(value)) s.prepend(new Option(`${value || "—"} (indisponível)`, value, true, true));
  s.addEventListener("change", () => onChange(s.value));
  f.append(el("span", "kv__k", label), s);
  return f;
}

// textField confirma o valor ao sair do campo ou com Enter (não a cada tecla: o inspetor é redesenhado).
function textField(label, value, placeholder, onCommit, listId, listItems) {
  const f = el("label", "field");
  const input = el("input", "input");
  input.type = "text";
  input.value = value || "";
  input.placeholder = placeholder || "";
  input.spellcheck = false;
  input.autocomplete = "off";
  f.append(el("span", "kv__k", label), input);
  if (listItems && listItems.length) {
    const dl = el("datalist");
    dl.id = listId;
    listItems.forEach((v) => dl.append(new Option(v)));
    input.setAttribute("list", listId);
    f.append(dl);
  }
  input.addEventListener("change", () => onCommit(input.value.trim()));
  input.addEventListener("keydown", (e) => { if (e.key === "Enter") input.blur(); });
  return f;
}

function segControl(label, options, current, onPick) {
  const seg = el("div", "seg");
  seg.style.gridTemplateColumns = `repeat(${options.length}, 1fr)`;
  seg.setAttribute("role", "group");
  seg.setAttribute("aria-label", label);
  options.forEach((o) => {
    const b = el("button", "seg__btn", o.label);
    b.type = "button";
    b.setAttribute("aria-pressed", String(current === o.value));
    if (o.disabled) {
      b.disabled = true;
      b.title = o.title || "";
    }
    b.addEventListener("click", () => onPick(o.value));
    seg.append(b);
  });
  return seg;
}

function modeSwitch(node) {
  return segControl("Modalidade de roteamento", [
    { value: "suggested", label: MODE_LABEL.suggested },
    { value: "manual", label: MODE_LABEL.manual },
    { value: "auto", label: MODE_LABEL.auto, disabled: !state.inv.router_online, title: "9Router offline: modo Auto indisponível" },
  ], node.routing.mode, (mode) => setMode(node, mode));
}

function setMode(node, mode) {
  const block = blockOf(node);
  if (mode === "auto") {
    node.routing = { mode, combo: modelIds("9router")[0] || "" };
  } else {
    const base = defaultRouting(block);
    node.routing = { mode, model: node.routing.model || base.model, fallback: mode === "suggested" ? base.fallback : "" };
  }
  onRoutingChanged(node);
}

function onRoutingChanged(node) {
  clearValidation();
  renderNode(node);
  renderInspector();
}

function routingFields(node) {
  const r = node.routing;
  const out = [];
  if (r.mode === "auto") {
    out.push(selectField("Combo do 9Router", r.combo, modelIds("9router"), (v) => { r.combo = v; onRoutingChanged(node); }));
    return out;
  }
  out.push(selectField("Modelo", r.model, modelIds(), (v) => { r.model = v; onRoutingChanged(node); }));
  if (r.mode === "suggested") {
    out.push(selectField("Fallback", r.fallback || "", ["", ...modelIds()], (v) => { r.fallback = v; onRoutingChanged(node); }));
  }
  return out;
}

function executionSection(node) {
  const sec = el("section", "exec");
  sec.append(el("span", "kv__k", "Execução"));
  sec.append(segControl("Modo de execução", [
    { value: "inline", label: "Inline" },
    { value: "isolated", label: "Isolada" },
  ], node.execution, (mode) => { node.execution = mode; clearValidation(); renderNode(node); renderInspector(); }));
  sec.append(el("p", "note", "Inline roda na pasta do projeto. Isolada usa um workspace separado e é obrigatória quando peças em paralelo gravam em pastas sobrepostas."));
  if (node.execution === "isolated") sec.append(el("p", "note note--warn", "O runner ainda não executa o modo isolado."));
  return sec;
}

// Campos próprios da peça genérica: artefato de saída e gate.
function pieceFields(node, sug) {
  const gates = (sug ? sug.gates : []).map((g) => g.expr);
  const defaultGate = `file_exists('${node.piece.artifact}')`;
  return [
    textField("Artefato de saída", node.piece.artifact, "NOME.md", (v) => {
      if (!v) { toast("O artefato de saída não pode ficar vazio.", true); renderInspector(); return; }
      node.piece.artifact = v;
      clearValidation();
      renderNode(node);
      renderEdges();
      renderInspector();
    }),
    textField("Gate de saída", node.piece.gate, `padrão: ${defaultGate}`, (v) => {
      node.piece.gate = v;
      clearValidation();
      renderNode(node);
      renderInspector();
    }, "gate-options", gates),
  ];
}

// ---------- Sugestões ----------

function suggestionKey(node) {
  return node.piece ? `p:${node.piece.kind}:${node.piece.artifact}` : `b:${node.block}`;
}

let suggestDebounceTimer = null;
let suggestAbortController = null;

// ensureSuggestions devolve a sugestão em cache, ou null enquanto busca (e redesenha ao chegar).
function ensureSuggestions(node) {
  const key = suggestionKey(node);
  if (state.suggest.has(key)) return state.suggest.get(key);

  if (suggestDebounceTimer) {
    clearTimeout(suggestDebounceTimer);
    suggestDebounceTimer = null;
  }
  if (suggestAbortController) {
    suggestAbortController.abort();
    suggestAbortController = null;
  }

  suggestDebounceTimer = setTimeout(() => {
    suggestAbortController = new AbortController();
    fetchSuggestions(node, suggestAbortController.signal).then((data) => {
      state.suggest.set(key, data);
      state.suggestPending.delete(key);
      const sel = state.selected;
      if (sel?.kind === "node" && state.nodes.get(sel.id) && suggestionKey(state.nodes.get(sel.id)) === key) renderInspector();
    }).catch(() => {});
  }, 150);

  return null;
}

async function fetchSuggestions(node, signal) {
  const q = node.piece
    ? `kind=${encodeURIComponent(node.piece.kind)}&artifact=${encodeURIComponent(node.piece.artifact)}`
    : `block=${encodeURIComponent(node.block)}`;
  try {
    const { ok, data } = await api(`/api/suggest?${q}`, undefined, signal);
    return ok ? data : EMPTY_SUGGEST;
  } catch {
    return EMPTY_SUGGEST;
  }
}

// ---------- Painéis ----------

function renderEdgePanel(body, index) {
  const edge = state.edges[index];
  if (!edge) {
    body.replaceChildren(el("p", "panel__hint", "Selecione um nó ou uma conexão."));
    return;
  }
  const ok = canConnect(edge.from, edge.to);
  const del = el("button", "btn", "Desconectar");
  del.type = "button";
  del.addEventListener("click", () => disconnect(index));
  body.replaceChildren(
    kv("Conexão", `${edge.from} → ${edge.to}`),
    el("p", ok ? "note" : "note note--bad", ok ? "Conexão compatível." : "Conexão incompatível: as portas não se encaixam. Religue ou desconecte."),
    el("p", "note", "Arraste a ponta da conexão para outra peça para religar; solte no vazio para remover. Delete também remove."),
    del,
  );
}

function triggerFields(node) {
  const block = blockOf(node);
  const cfg = node.triggerConfig || defaultTriggerConfig(block.kind);
  node.triggerConfig = cfg;
  const out = [];

  switch (block.kind) {
    case "trigger_cron":
      out.push(textField("Expressão Cron (5 campos)", cfg.cron || "0 * * * *", "ex: */15 * * * *", (v) => {
        cfg.cron = v;
        clearValidation();
      }));
      out.push(el("p", "note", "Define o agendamento formal (minuto, hora, dia, mês, dia da semana)."));
      break;
    case "trigger_file_change":
      out.push(textField("Caminho monitorado (Path)", cfg.path || "docs", "ex: docs/intents", (v) => {
        cfg.path = v;
        clearValidation();
      }));
      out.push(el("p", "note", "Dispara a esteira sempre que arquivos no caminho sofrerem alteração."));
      break;
    case "trigger_loop_end":
      out.push(textField("Loop anterior (Loop ID)", cfg.loop || "", "ex: pipeline-master", (v) => {
        cfg.loop = v;
        clearValidation();
      }));
      out.push(el("p", "note", "Dispara ao término com sucesso do loop especificado."));
      break;
    case "trigger_board":
      out.push(textField("Coluna / Evento do Kanban", cfg.board || "backlog", "ex: ready", (v) => {
        cfg.board = v;
        clearValidation();
      }));
      out.push(el("p", "note", "Disparado por cartões prontos no GitHub Project Kanban."));
      break;
    default:
      out.push(el("p", "note", "Disparo manual sob demanda (via botão Executar no Studio ou CLI legoloop run)."));
      break;
  }
  return out;
}

function renderNodePanel(body, node) {
  const block = blockOf(node);
  const isTrig = block.kind && block.kind.startsWith("trigger_");
  const del = el("button", "btn btn--ghost", "Remover nó");
  del.type = "button";
  del.addEventListener("click", () => { removeNode(node.id); select(null); });

  if (isTrig) {
    const children = [
      kv("Gatilho Inicial", `[Gatilho] ${block.label}`),
      ...triggerFields(node),
      kv("Handoff emitido", block.artifact),
      del,
    ];
    body.replaceChildren(...children);
    return;
  }

  const sug = ensureSuggestions(node);
  const swap = el("button", "btn", "Trocar peça");
  swap.type = "button";
  swap.addEventListener("click", () => startSwap(node));
  const draft = state.gateDraft?.nodeId === node.id ? state.gateDraft.expr : "";
  state.gateDraft = null;
  const rerender = () => { clearValidation(); renderNode(node); renderInspector(); };
  const children = [kv(node.piece ? "Peça" : "Bloco", node.piece ? `${block.label} (${node.piece.kind})` : `${block.label} (${block.ref})`)];

  const isSubloop = (node.piece && (node.piece.kind === "subloop" || node.piece.kind === "loop")) || (block.kind && (block.kind === "loop" || block.kind === "flow_loop"));
  if (isSubloop) {
    const drill = el("button", "btn btn--primary", "Abrir Sub-grafo (Drill-Down)");
    drill.type = "button";
    drill.style.marginBottom = "8px";
    drill.addEventListener("click", () => {
      if (typeof window.enterSubloop === "function") {
        window.enterSubloop(node.piece?.ref || node.block || node.id, block.label);
      }
    });
    children.push(drill);
  }

  if (node.piece) {
    children.push(...pieceFields(node, sug));
  } else {
    children.push(kv("Handoff imutável", block.artifact), kv("Gate", block.gate));
  }
  children.push(modeSwitch(node), ...routingFields(node));
  if (!node.piece) {
    children.push(el("p", "note", hasModel(block.suggested_model)
      ? `Sugestão de benchmark: ${block.suggested_model}`
      : `Sugestão ${block.suggested_model} não está no inventário desta máquina.`));
  }
  children.push(executionSection(node), renderSlots(node, sug, rerender, draft), renderNeighbors(sug, (n) => addNeighbor(node, n)), swap, del);
  body.replaceChildren(...children);
}

// Painel do processo recém-fabricado: a Tríade e os três meios, com os caminhos já mascarados pelo servidor.
function renderProcessPanel(body, proc) {
  const children = [
    el("h3", "panel__title", proc.title),
    el("p", "panel__hint", "Processo fabricado: Tríade de governança e os três meios de execução."),
    ...proc.rows.map((r) => kv(r.label, r.value)),
    ...proc.notes.map((n) => el("p", "panel__hint", "Aviso: " + n)),
  ];
  if (proc.report) {
    if (proc.report.cartao && typeof renderCartaoSTE === "function") {
      children.push(renderCartaoSTE(proc.report.cartao));
    }
    if (typeof renderRelatorioCamadas === "function") {
      children.push(renderRelatorioCamadas(proc.report));
    }
  }
  body.replaceChildren(...children);
}

function renderInspector() {
  const body = $("inspector-body");
  const sel = state.selected;
  $("inspector").classList.toggle("is-open", Boolean(sel));
  if (sel?.kind === "process" && state.process) return renderProcessPanel(body, state.process);
  if (sel?.kind === "edge") return renderEdgePanel(body, sel.id);
  if (!sel || sel.kind !== "node" || !state.nodes.has(sel.id)) {
    body.replaceChildren(el("p", "panel__hint", "Selecione um nó para configurar roteamento, execução e slots."));
    return undefined;
  }
  return renderNodePanel(body, state.nodes.get(sel.id));
}
