// Executor de auditoria de harness (opencode ou mimo): modelo a modelo, de uma tacada (headless), sem stdin.
// Uso: node hrun.mjs <opencode|mimo> [limite_de_modelos]
// Fases: 1 disponibilidade na condição LIMPA (B); 2 effort na B; 3 disponibilidade na configuração do usuário (A).
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { spawn, execSync } from "node:child_process";

const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const h = process.argv[2];
const limiteModelos = Number(process.argv[3] || 0);
const NL = String.fromCharCode(10);
const PAR = Number(process.env.PAR || 4);
const saida = path.join(dir, `${h}.jsonl`);
fs.writeFileSync(saida, "");
const cfgVazio = fs.mkdtempSync(path.join(os.tmpdir(), "aud-cfg-"));

const PROMPTS = {
  ola: { txt: "Responda apenas a palavra: OK", ok: (t) => /\bOK\b/i.test(t) },
  conta: { txt: "Quanto e 17 vezes 23? Responda apenas o numero.", ok: (t) => t.includes("391") },
};
const SEGURO = /^[A-Za-z0-9._:/@-]+$/;

function rodar(modelo, variant, prompt, condicao, limiteMs, cwd) {
  return new Promise((resolve) => {
    const t0 = Date.now();
    const cmd = `${h} run -m ${modelo}${variant ? ` --variant ${variant}` : ""} --format json "${PROMPTS[prompt].txt}"`;
    const env = { ...process.env };
    if (condicao === "B") { env.XDG_CONFIG_HOME = cfgVazio; if (h === "opencode") env.OPENCODE_DISABLE_CLAUDE_CODE = "1"; }
    const filho = spawn(cmd, { shell: true, cwd, env, stdio: ["ignore", "pipe", "pipe"], windowsHide: true });
    let out = "", estourou = false;
    filho.stdout.on("data", (d) => { out += d; });
    filho.stderr.on("data", () => {});
    const timer = setTimeout(() => {
      estourou = true;
      try { execSync(`taskkill /PID ${filho.pid} /T /F`, { stdio: "ignore" }); } catch { /* já saiu */ }
    }, limiteMs);
    filho.on("close", (codigo) => { clearTimeout(timer); resolve({ out, codigo, estourou, ms: Date.now() - t0 }); });
  });
}

function analisar(r, prompt) {
  let texto = "", erro = "", tokens = null;
  for (const l of r.out.split(NL)) {
    if (!l.trim().startsWith("{")) continue;
    let j; try { j = JSON.parse(l); } catch { continue; }
    if (j.type === "text") texto += j.part?.text || j.text || "";
    else if (j.type === "error") erro = j.error?.data?.message || j.error?.message || JSON.stringify(j.error).slice(0, 200);
    else if (/finish/.test(j.type || "")) { const t = j.part?.tokens || j.tokens; if (t) tokens = { entrada: t.input ?? null, saida: t.output ?? null, raciocinio: t.reasoning ?? null, total: t.total ?? null }; }
  }
  const limpo = texto.replace(/\s+/g, " ").trim();
  const respondeu = !erro && limpo.length > 0 && r.codigo === 0;
  return { respondeu, correto: respondeu && PROMPTS[prompt].ok(limpo), texto: limpo.slice(0, 60), erro: String(erro).replace(/\s+/g, " ").slice(0, 200), tokens };
}

async function caso(c, cwd) {
  let r = await rodar(c.modelo, c.variant, c.prompt, c.condicao, 60000, cwd), fase = "60s";
  if (r.estourou) { r = await rodar(c.modelo, c.variant, c.prompt, c.condicao, 180000, cwd); fase = "180s"; }
  const a = analisar(r, c.prompt);
  const reg = { harness: h, fase_teste: c.fase, modelo: c.modelo, provedor: c.provedor, variant: c.variant || "", prompt: c.prompt, condicao: c.condicao, ok: a.respondeu, correto: a.correto, ms: r.ms, limite: fase, timeout: r.estourou, exit: r.codigo, texto: a.texto, erro: r.estourou ? "sem resposta no limite" : a.erro, tokens: a.tokens };
  fs.appendFileSync(saida, JSON.stringify(reg) + NL);
  return reg;
}

async function lote(casos) {
  const res = []; let i = 0;
  await Promise.all(Array.from({ length: PAR }, async (_, k) => {
    const cwd = fs.mkdtempSync(path.join(os.tmpdir(), `aud-${h}-${k}-`));
    while (i < casos.length) { const c = casos[i++]; res.push(await caso(c, cwd)); }
  }));
  return res;
}

let cand = JSON.parse(fs.readFileSync(path.join(dir, `cand_${h}.json`), "utf8")).filter((m) => SEGURO.test(m.id));
if (limiteModelos) cand = cand.slice(0, limiteModelos);
const mapa = new Map(cand.map((m) => [m.id, m]));

const f1 = await lote(cand.map((m) => ({ fase: "1-disponibilidade-limpa", modelo: m.id, provedor: m.provedor, variant: "", prompt: "ola", condicao: "B" })));
const passou = new Set(f1.filter((r) => r.ok).map((r) => r.modelo));

const casos2 = [];
for (const id of passou) {
  const m = mapa.get(id);
  if (!m.efforts.length) continue;
  for (const v of [...new Set([m.efforts[0], m.efforts[m.efforts.length - 1]])]) casos2.push({ fase: "2-effort-limpa", modelo: id, provedor: m.provedor, variant: v, prompt: "conta", condicao: "B" });
}
await lote(casos2);

await lote([...passou].map((id) => ({ fase: "3-disponibilidade-config-usuario", modelo: id, provedor: mapa.get(id).provedor, variant: "", prompt: "ola", condicao: "A" })));
fs.writeFileSync(path.join(dir, `${h}_fim.txt`), `fim: ${cand.length} modelos, ${passou.size} passaram na fase 1, ${casos2.length} casos de effort`);
console.log(`ok: ${cand.length} modelos`);
