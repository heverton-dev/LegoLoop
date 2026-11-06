// Executor de auditoria dos harnesses por assinatura (claude e agy), headless, com o perfil real do usuário.
// Uso: node hrun2.mjs <claude|agy> [limite_de_modelos]
// Fases: 1 disponibilidade (effort padrão); 2 cada effort declarado nos modelos que passaram.
import fs from "node:fs";
import path from "node:path";
import { spawn, execSync } from "node:child_process";

const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const h = process.argv[2];
const limiteModelos = Number(process.argv[3] || 0);
const NL = String.fromCharCode(10);
const PAR = Number(process.env.PAR || 2);
const WT = { claude: "C:/Users/trcnologia/orca/workspaces/LegoLoop/e03-aud-claude", agy: "C:/Users/trcnologia/orca/workspaces/LegoLoop/e03-aud-agy" }[h];
const saida = path.join(dir, `${h}.jsonl`);
fs.writeFileSync(saida, "");

const PROMPTS = {
  ola: { txt: "Responda apenas a palavra: OK", ok: (t) => /\bOK\b/i.test(t) },
  conta: { txt: "Quanto e 17 vezes 23? Responda apenas o numero.", ok: (t) => t.includes("391") },
};
const SEGURO = /^[A-Za-z0-9._\-\[\]]+$/;

function comando(modelo, effort, prompt) {
  const ef = effort ? ` --effort ${effort}` : "";
  const m = modelo.includes("[") ? `"${modelo}"` : modelo;
  if (h === "claude") return `claude -p "${PROMPTS[prompt].txt}" --model ${m}${ef} --output-format json --no-session-persistence`;
  const bin = process.env.AGY_BIN ? `"${process.env.AGY_BIN}"` : "agy"; // AGY_BIN: binário real, sem o failover do gw exec
  return `${bin} -p "${PROMPTS[prompt].txt}" --model ${m}${ef} --output-format json --print-timeout 90s`;
}

function rodar(modelo, effort, prompt, limiteMs) {
  return new Promise((resolve) => {
    const t0 = Date.now();
    const filho = spawn(comando(modelo, effort, prompt), { shell: true, cwd: WT, stdio: ["ignore", "pipe", "pipe"], windowsHide: true });
    let out = "", err = "", estourou = false;
    filho.stdout.on("data", (d) => { out += d; });
    filho.stderr.on("data", (d) => { err += d; });
    const timer = setTimeout(() => { estourou = true; try { execSync(`taskkill /PID ${filho.pid} /T /F`, { stdio: "ignore" }); } catch { /* já saiu */ } }, limiteMs);
    filho.on("close", (codigo) => { clearTimeout(timer); resolve({ out, err, codigo, estourou, ms: Date.now() - t0 }); });
  });
}

function analisar(r, prompt) {
  let j = null;
  try { j = JSON.parse(r.out.trim()); } catch { /* não JSON */ }
  let texto = "", erro = "", tokens = null, real = "", custo = null;
  if (j && h === "claude") {
    texto = j.result || ""; if (j.is_error) erro = texto || j.subtype || "erro";
    const u = j.usage || {};
    tokens = { entrada: (u.input_tokens ?? 0) + (u.cache_creation_input_tokens ?? 0) + (u.cache_read_input_tokens ?? 0), saida: u.output_tokens ?? null, raciocinio: u.output_tokens_details?.thinking_tokens ?? null };
    real = Object.keys(j.modelUsage || {}).join("+"); custo = j.total_cost_usd ?? null;
  } else if (j && h === "agy") {
    texto = j.response || ""; if (j.status && j.status !== "SUCCESS") erro = `${j.status}: ${j.error || texto}`.slice(0, 200);
    const u = j.usage || {};
    tokens = { entrada: u.input_tokens ?? null, saida: u.output_tokens ?? null, raciocinio: u.thinking_tokens ?? null };
  } else { erro = (r.err || r.out || "").replace(/\s+/g, " ").slice(0, 200) || `exit ${r.codigo}`; }
  const limpo = texto.replace(/\s+/g, " ").trim();
  const respondeu = !erro && limpo.length > 0 && r.codigo === 0;
  return { respondeu, correto: respondeu && PROMPTS[prompt].ok(limpo), texto: limpo.slice(0, 60), erro: String(erro).replace(/\s+/g, " ").slice(0, 200), tokens, real, custo };
}

// agy: contas gerenciadas pelo gw (`gw agy next` rotaciona; `gw agy probe` diz se a ativa está "Disponível").
let conta = process.env.CONTA_INICIAL || "DEV_03";
const limpaAnsi = (s) => s.replace(/\x1b\[[0-9;]*m/g, "");
function gw(args) { try { return limpaAnsi(execSync(`gw agy ${args}`, { encoding: "utf8", stdio: ["ignore", "pipe", "pipe"], timeout: 90000 })); } catch (e) { return limpaAnsi(String(e.stdout || "") + String(e.stderr || "")); } }
function rotacionar() {
  for (let i = 0; i < 13; i++) {
    const saida = gw("next");
    const m = saida.match(/DEV_\d+/);
    if (m) conta = m[0];
    const sonda = gw("probe");
    if (/Dispon[ií]vel/i.test(sonda)) return { conta, ok: true };
  }
  return { conta, ok: false };
}
const FALHA_DE_CONTA = /Eligibility|not eligible|Individual quota|quota reached|Relogin|Verify your account|PERMISSION_DENIED/i;

async function caso(c) {
  let r = await rodar(c.modelo, c.effort, c.prompt, 60000), fase = "60s";
  if (r.estourou) { r = await rodar(c.modelo, c.effort, c.prompt, 180000); fase = "180s"; }
  let a = analisar(r, c.prompt);
  const contas = [conta];
  for (let giro = 0; h === "agy" && !process.env.NAO_ROTACIONAR && !a.respondeu && FALHA_DE_CONTA.test(a.erro) && giro < 3; giro++) {
    const rot = rotacionar(); contas.push(rot.conta);
    if (!rot.ok) break;
    r = await rodar(c.modelo, c.effort, c.prompt, 60000); fase = "60s";
    if (r.estourou) { r = await rodar(c.modelo, c.effort, c.prompt, 180000); fase = "180s"; }
    a = analisar(r, c.prompt);
  }
  const reg = { conta_final: h === "agy" ? conta : "", contas_tentadas: contas.join(">"), harness: h, fase_teste: c.fase, modelo: c.modelo, effort: c.effort || "", prompt: c.prompt, ok: a.respondeu, correto: a.correto, ms: r.ms, limite: fase, timeout: r.estourou, exit: r.codigo, modelo_real: a.real, texto: a.texto, erro: r.estourou ? "sem resposta no limite" : a.erro, tokens: a.tokens, custo_usd: a.custo };
  fs.appendFileSync(saida, JSON.stringify(reg) + NL);
  return reg;
}

async function lote(casos) {
  const res = []; let i = 0;
  await Promise.all(Array.from({ length: PAR }, async () => { while (i < casos.length) { const c = casos[i++]; res.push(await caso(c)); } }));
  return res;
}

let cand = JSON.parse(fs.readFileSync(path.join(dir, `cand_${h}.json`), "utf8")).filter((m) => SEGURO.test(m.id));
// claude: só ids distintos (sem as variantes [1m], auto e code-tests, que repetem os mesmos modelos)
if (h === "claude") cand = cand.filter((m) => !m.id.includes("[") && m.id !== "code-tests");
if (limiteModelos) cand = cand.slice(0, limiteModelos);
const mapa = new Map(cand.map((m) => [m.id, m]));

const f1 = await lote(cand.map((m) => ({ fase: "1-disponibilidade", modelo: m.id, effort: "", prompt: "ola" })));
const passou = new Set(f1.filter((r) => r.ok).map((r) => r.modelo));

const casos2 = [];
for (const id of passou) {
  const m = mapa.get(id);
  if (h === "claude" && !/^claude-/.test(id)) continue; // apelidos e code-* resolvem para os mesmos modelos
  // sem lista declarada: um nível, para ver se é recusado; com lista: low, high, xhigh (se declarado) e o último (medium não é testado)
  const niveis = h === "agy" ? m.efforts : m.efforts.length ? [...new Set([m.efforts[0], ...["high", "xhigh"].filter((x) => m.efforts.includes(x)), m.efforts[m.efforts.length - 1]])] : ["low"];
  for (const e of niveis) casos2.push({ fase: "2-effort", modelo: id, effort: e, prompt: "conta" });
}
await lote(casos2);
fs.writeFileSync(path.join(dir, `${h}_fim.txt`), `fim: ${cand.length} modelos, ${passou.size} passaram na fase 1, ${casos2.length} casos de effort`);
console.log(`ok: ${cand.length} modelos`);
