"use strict";

// Lógica pura do modal "Novo Processo" (sem DOM, testável em node). O app.js só liga estas funções aos
// elementos. Todo texto que vem da API entra no DOM via textContent.

const FABRICAR_BRAIN_SOURCE = "9router";

// Só modelos do 9Router conduzem a fabricação; os demais aparecem desabilitados, com o motivo.
function fabricarModelOptions(models) {
  return (models || []).map((m) => {
    const ok = m.source === FABRICAR_BRAIN_SOURCE;
    return { id: m.id, disabled: !ok, label: ok ? m.id : `${m.id} (fora do 9Router)` };
  });
}

// O padrão é o primeiro modelo habilitado; sem nenhum, vazio (o botão de iniciar fica bloqueado).
function fabricarDefaultModel(options) {
  const first = options.find((o) => !o.disabled);
  return first ? first.id : "";
}

// Devolve "" quando o nome é aceitável. A regra completa de normalização é do servidor (NormalizeProcessName).
function fabricarNameError(name) {
  const n = String(name || "").trim();
  if (!n) return "Informe o nome do processo.";
  if (!/[\p{L}\p{N}]/u.test(n)) return "Use ao menos uma letra ou número no nome.";
  return "";
}

// Corpo de POST /api/fabricar/run. O objetivo vem do espelho da entrevista; sem entrevista fica vazio e o
// servidor usa o próprio nome. O cliente nunca envia caminhos, force nem confiança em mídia removível.
function fabricarPayload(form, turn, force, sessionId) {
  const goal = (turn && ((turn.mirror && turn.mirror.outcome) || turn.hypothesis)) || "";
  const body = { name: String(form.name || "").trim(), goal, category: form.category || "build", model: form.model };
  if (sessionId) body.entrevista = sessionId;
  if (force === true) body.force = true; // só depois de o usuário confirmar a sobrescrita
  return body;
}

// A sobrescrita só é oferecida quando o servidor diz que o processo já existe (e não quando há outra fabricação rodando).
function fabricarCanOverwrite(status, data) {
  return status === 409 && !!data && data.code === "exists";
}

// Plotar substitui o Canvas: com nós na tela, o usuário precisa confirmar antes.
function fabricarNeedsReplaceConfirm(nodeCount) {
  return nodeCount > 0;
}

function fabricarInterviewDone(turn) {
  return !!turn && turn.status === "COMPLETED";
}

// Mensagem para o próprio modal. 400, 409, 422 e 502 repetem o motivo do servidor (já sem caminhos locais).
function fabricarErrorMessage(status, data) {
  const reason = data && typeof data.error === "string" ? data.error : "";
  if ([400, 409, 422, 502].includes(status) && reason) return reason;
  if (status === 0) return "Sem resposta do servidor. Verifique se o Studio continua rodando.";
  if (status === 429) return "Muitas requisições; aguarde um instante e tente de novo.";
  if (status >= 500) return "A fábrica falhou e nenhum arquivo foi criado. Tente novamente ou escolha outro modelo.";
  return reason || `Erro HTTP ${status}.`;
}

// ---------- Plotagem no Canvas ----------

// Mesmas medidas do Canvas (NODE_W em app.js) mais a folga entre nós, para os nós gerados não se sobreporem.
const FABRICAR_NODE_W = 236;
const FABRICAR_NODE_GAP = 80;
const FABRICAR_ORIGIN = { x: 60, y: 100 };
const FABRICAR_PIECE_KIND = { skill: "skill", script: "script", mcp: "mcp", loop: "subloop" };

function fabricarRefName(ref) {
  return String(ref || "").split("/").pop();
}

// Converte o loop gerado (formato DAG: depends_on, output) em nós e arestas do Canvas, e o Result da fábrica
// nos dados do Inspetor e na entrada da categoria Loops. Pura: quem aplica no DOM é o app.js.
function fabricarPlotModel(result, loop) {
  const src = (loop && loop.nodes) || [];
  const known = new Set(src.map((n) => n.id));
  const nodes = src.map((n, i) => ({
    id: n.id,
    label: fabricarRefName(n.ref) || n.id,
    x: FABRICAR_ORIGIN.x + i * (FABRICAR_NODE_W + FABRICAR_NODE_GAP),
    y: FABRICAR_ORIGIN.y + (i % 2) * 40,
    piece: {
      kind: FABRICAR_PIECE_KIND[n.type] || "skill", ref: n.ref || "",
      artifact: (n.output && n.output.artifact) || "", gate: (n.output && n.output.gate) || "",
    },
    routing: {
      mode: (n.routing && n.routing.mode) || "suggested", model: (n.routing && n.routing.model) || "",
      fallback: (n.routing && n.routing.fallback) || "",
    },
    execution: "inline",
  }));
  const edges = [];
  for (const n of src) {
    for (const dep of n.depends_on || []) if (known.has(dep)) edges.push({ from: dep, to: n.id });
  }
  const name = result.canonical_name;
  return { name, nodes, edges, loopEntry: { cat: "loop", name, itemId: name }, inspector: fabricarInspector(result) };
}

// Dados da Tríade e dos três meios para o painel do processo no Inspetor. Os caminhos já chegam mascarados.
function fabricarInspector(result) {
  const a = result.artifacts || {};
  const rows = [
    ["INTENT", a.intent_path], ["SPEC", a.spec_path], ["PLAN", a.plan_path], ["ROUTE", a.route_path],
    ["TEMPLATE", a.template_path], ["Skill", a.skill_path], ["Loop", a.loop_path], ["DAG", a.dag_path],
  ].filter(([, v]) => v).map(([label, value]) => ({ label, value }));
  rows.push({ label: "Peças curadas", value: String((result.curated_pieces || []).length) });
  return { title: result.canonical_name, rows, notes: result.warnings || [], report: result.report };
}
