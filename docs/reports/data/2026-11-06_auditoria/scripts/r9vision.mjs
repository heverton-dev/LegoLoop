// Refaz só o teste de visão do 9Router com um PNG válido (16x16 vermelho), gerado aqui com zlib.
import fs from "node:fs";
import path from "node:path";
import zlib from "node:zlib";
const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const BASE = (process.env.NINEROUTER_URL || "").replace(/\/+$/, "");
const KEY = process.env.NINEROUTER_KEY || "";
const H = { "Content-Type": "application/json", ...(KEY ? { Authorization: `Bearer ${KEY}` } : {}) };

function crc32(buf) {
  let c, crc = 0xffffffff;
  for (let n = 0; n < buf.length; n++) { c = (crc ^ buf[n]) & 0xff; for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1; crc = (crc >>> 8) ^ c; }
  return (crc ^ 0xffffffff) >>> 0;
}
function chunk(tipo, dados) {
  const len = Buffer.alloc(4); len.writeUInt32BE(dados.length);
  const corpo = Buffer.concat([Buffer.from(tipo, "ascii"), dados]);
  const crc = Buffer.alloc(4); crc.writeUInt32BE(crc32(corpo));
  return Buffer.concat([len, corpo, crc]);
}
function png(w, h, r, g, b) {
  const linha = Buffer.concat([Buffer.from([0]), Buffer.from(Array.from({ length: w }, () => [r, g, b]).flat())]);
  const raw = Buffer.concat(Array.from({ length: h }, () => linha));
  const ihdr = Buffer.alloc(13); ihdr.writeUInt32BE(w, 0); ihdr.writeUInt32BE(h, 4); ihdr[8] = 8; ihdr[9] = 2;
  return Buffer.concat([Buffer.from([137, 80, 78, 71, 13, 10, 26, 10]), chunk("IHDR", ihdr), chunk("IDAT", zlib.deflateSync(raw)), chunk("IEND", Buffer.alloc(0))]);
}
const b64 = png(16, 16, 255, 0, 0).toString("base64");
fs.writeFileSync(path.join(dir, "vermelho16.png"), Buffer.from(b64, "base64"));

const alvos = ["megacombo-vision", "gemini/gemini-3.6-flash", "gemini/gemini-3.8-flash", "cloudflare-ai/llama-3.2-11b-vision-instruct", "gemini/gemini-3.5-flash", "gemini/gemma-4-31b-it", "code-fast", "kr/claude-sonnet-4.5"];
const out = path.join(dir, "9router_visao.jsonl");
fs.writeFileSync(out, "");
for (const alvo of alvos) {
  const t0 = Date.now();
  const ctl = new AbortController(); const timer = setTimeout(() => ctl.abort(), 90000);
  let reg;
  try {
    const r = await fetch(`${BASE}/v1/chat/completions`, { method: "POST", headers: H, signal: ctl.signal, body: JSON.stringify({ model: alvo, max_tokens: 300, messages: [{ role: "user", content: [{ type: "text", text: "Qual a cor predominante desta imagem? Responda com uma palavra." }, { type: "image_url", image_url: { url: `data:image/png;base64,${b64}` } }] }] }) });
    const bruto = await r.text();
    const txt = bruto.replace(/data:\s*\[DONE\]\s*$/, "").trim();
    let j = {}, resp = "";
    if (bruto.trimStart().startsWith("data:")) {
      let c = "", rc = "", usado = "";
      for (const l of bruto.split(String.fromCharCode(10))) {
        if (!l.startsWith("data:") || l.includes("[DONE]")) continue;
        try { const x = JSON.parse(l.slice(5)); usado = x.model || usado; const d = x.choices?.[0]?.delta || {}; c += d.content || ""; rc += d.reasoning_content || d.reasoning || ""; if (x.error) j.error = x.error; } catch { /* parcial */ }
      }
      j.model = usado; resp = (c + " " + rc).replace(/\s+/g, " ").trim();
    } else {
      try { j = JSON.parse(txt); } catch { /* não JSON */ }
      const m = j?.choices?.[0]?.message || {};
      resp = ((m.content || "") + " " + (m.reasoning_content || m.reasoning || "")).replace(/\s+/g, " ").trim();
    }
    reg = { alvo, http: r.status, ms: Date.now() - t0, usado: j?.model || "", respondeu: r.status === 200 && !!resp, correto: /vermelh|red/i.test(resp), resp: resp.slice(0, 50), erro: (j?.error?.message || (r.ok ? "" : txt.slice(0, 120))).replace(/\s+/g, " ").slice(0, 130) };
  } catch (e) { reg = { alvo, http: 0, ms: Date.now() - t0, respondeu: false, correto: false, resp: "", erro: e.name === "AbortError" ? "timeout 90s" : String(e.message).slice(0, 100) }; } finally { clearTimeout(timer); }
  fs.appendFileSync(out, JSON.stringify(reg) + String.fromCharCode(10));
  console.log(`${reg.correto ? "OK    " : "FALHOU"} ${String(reg.http).padEnd(4)}${String(reg.ms).padStart(7)}ms ${alvo.padEnd(44)} ${(reg.usado || "-").padEnd(26)} ${reg.correto ? reg.resp : reg.erro.slice(0, 70)}`);
}
