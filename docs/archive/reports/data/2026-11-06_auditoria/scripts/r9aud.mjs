// Auditoria do 9Router por HTTP: cada combo e cada membro, com a bateria de prompts. A chave vem do ambiente e nunca é impressa.
import fs from "node:fs";
import path from "node:path";
const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const BASE = (process.env.NINEROUTER_URL || "").replace(/\/+$/, "");
const KEY = process.env.NINEROUTER_KEY || process.env.NINEROUTER_API_KEY || "";
const H = { "Content-Type": "application/json", ...(KEY ? { Authorization: `Bearer ${KEY}` } : {}) };
const combos = JSON.parse(fs.readFileSync(path.join(dir, "combos.json"), "utf8"));
const out = path.join(dir, "9router.jsonl");
fs.writeFileSync(out, "");
const NL = String.fromCharCode(10);
// 2x2 PNG vermelho
const PNG_VERMELHO = "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAAEklEQVR4nGP4z8Dwn4GBgQEAFQ4DAZ2yjoAAAAAASUVORK5CYII=";

const BATERIA = {
  ola: { msg: [{ role: "user", content: "Responda apenas a palavra: OK" }], ok: (t) => /\bOK\b/i.test(t) },
  conta: { msg: [{ role: "user", content: "Quanto e 17 vezes 23? Responda apenas o numero." }], ok: (t) => t.includes("391") },
  codigo: { msg: [{ role: "user", content: "Escreva em Go a funcao Soma(a, b int) int. Responda apenas o codigo." }], ok: (t) => /func\s+Soma/.test(t) },
  visao: { msg: [{ role: "user", content: [{ type: "text", text: "Qual a cor predominante desta imagem? Responda com uma palavra." }, { type: "image_url", image_url: { url: `data:image/png;base64,${PNG_VERMELHO}` } }] }], ok: (t) => /vermelh|red/i.test(t) },
};

function extrair(txt) {
  let content = "", racio = "", usado = "", erro = "", fim = "";
  const s = txt.trimStart();
  if (s.startsWith("data:")) {
    for (const l of txt.split(NL)) {
      if (!l.startsWith("data:") || l.includes("[DONE]")) continue;
      try {
        const c = JSON.parse(l.slice(5));
        usado = c.model || usado;
        const d = c.choices?.[0]?.delta || {};
        content += d.content || ""; racio += d.reasoning_content || d.reasoning || "";
        fim = c.choices?.[0]?.finish_reason || fim;
        if (c.error) erro = c.error.message || JSON.stringify(c.error);
      } catch { /* linha parcial */ }
    }
  } else {
    let j = {};
    try { j = JSON.parse(txt.replace(/data:\s*\[DONE\]\s*$/, "").trim()); } catch { /* não JSON */ }
    const m = j?.choices?.[0]?.message || {};
    content = m.content || ""; racio = m.reasoning_content || m.reasoning || "";
    usado = j?.model || ""; fim = j?.choices?.[0]?.finish_reason || "";
    erro = j?.error?.message || (j?.error ? JSON.stringify(j.error) : "");
    if (!j?.choices && !erro) erro = txt.slice(0, 140);
  }
  return { content, racio, usado, erro, fim };
}

async function post(caminho, corpo, limite) {
  const t0 = Date.now(); const ctl = new AbortController(); const timer = setTimeout(() => ctl.abort(), limite);
  try {
    const r = await fetch(`${BASE}${caminho}`, { method: "POST", headers: H, signal: ctl.signal, body: JSON.stringify(corpo) });
    return { http: r.status, txt: await r.text(), ms: Date.now() - t0, timeout: false };
  } catch (e) {
    return { http: 0, txt: "", ms: Date.now() - t0, timeout: e.name === "AbortError", falha: String(e.message).slice(0, 100) };
  } finally { clearTimeout(timer); }
}

async function uma(modelo, nome) {
  const b = BATERIA[nome];
  const corpo = { model: modelo, max_tokens: 400, messages: b.msg };
  let r = await post("/v1/chat/completions", corpo, 60000), fase = "60s";
  if (r.timeout) { r = await post("/v1/chat/completions", corpo, 180000); fase = "180s"; }
  const x = extrair(r.txt);
  const texto = (x.content || "") + " " + (x.racio || "");
  const respondeu = r.http === 200 && !x.erro && (x.content || x.racio);
  const reg = {
    alvo: modelo, teste: nome, http: r.http, ms: r.ms, fase, timeout: r.timeout,
    respondeu: !!respondeu, correto: !!respondeu && b.ok(texto), usado: x.usado,
    resp: (x.content || x.racio || "").replace(/\s+/g, " ").trim().slice(0, 60), erro: (x.erro || r.falha || "").replace(/\s+/g, " ").slice(0, 160),
  };
  fs.appendFileSync(out, JSON.stringify(reg) + NL);
  return reg;
}

async function embedding(modelo) {
  const t0 = Date.now();
  const r = await post("/v1/embeddings", { model: modelo, input: "teste" }, 60000);
  let n = 0, erro = "", usado = "";
  try { const j = JSON.parse(r.txt); n = j?.data?.[0]?.embedding?.length || 0; usado = j?.model || ""; erro = j?.error?.message || ""; } catch { erro = r.txt.slice(0, 140); }
  const reg = { alvo: modelo, teste: "embedding", http: r.http, ms: r.ms, fase: "60s", timeout: r.timeout, respondeu: n > 0, correto: n > 0, usado, resp: n ? `vetor de ${n}` : "", erro: String(erro).replace(/\s+/g, " ").slice(0, 160) };
  fs.appendFileSync(out, JSON.stringify(reg) + NL);
  return reg;
}

async function fila(tarefas, par) {
  let i = 0;
  await Promise.all(Array.from({ length: par }, async () => { while (i < tarefas.length) { const t = tarefas[i++]; await t(); } }));
}

const listaRemota = await (await fetch(`${BASE}/v1/models`, { headers: H })).json().then((j) => (j.data || []).map((m) => m.id));
fs.writeFileSync(path.join(dir, "9router_modelos_remotos.json"), JSON.stringify(listaRemota));
const membros = [...new Set(combos.flatMap((c) => c.membros))];
const alvos = [...combos.map((c) => c.nome), ...membros];
const cv = combos.find((c) => c.nome === "megacombo-vision");
const visaoAlvos = new Set(cv ? [cv.nome, ...cv.membros] : []);
const tarefas = [];
for (const alvo of alvos) {
  const emb = /embedding/i.test(alvo);
  if (emb) { tarefas.push(() => embedding(alvo)); continue; }
  const testes = ["ola", "conta", "codigo"];
  if (visaoAlvos.has(alvo)) testes.push("visao");
  for (const t of testes) tarefas.push(() => uma(alvo, t));
}
await fila(tarefas, 4);
fs.writeFileSync(path.join(dir, "9router_fim.txt"), `fim ${tarefas.length} chamadas`);
console.log(`ok: ${tarefas.length} chamadas`);
