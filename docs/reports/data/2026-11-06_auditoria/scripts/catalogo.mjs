// Gera o catálogo estruturado de um harness (opencode ou mimo) a partir de `<h> models --verbose`.
// Uso: node catalogo.mjs <harness>
import fs from "node:fs";
import path from "node:path";
import { execSync } from "node:child_process";
const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const h = process.argv[2];
const NL = String.fromCharCode(10);
const bruto = execSync(`${h} models --verbose`, { maxBuffer: 400 * 1024 * 1024, encoding: "utf8", stdio: ["ignore", "pipe", "ignore"] });
const linhas = bruto.split(NL);
const modelos = [];
let i = 0;
while (i < linhas.length) {
  const cab = linhas[i].trim();
  if (!cab || cab.startsWith("{")) { i++; continue; }
  // cabeçalho: "provedor/modelo" (o mimo acrescenta " — window ..."); depois vem um JSON multilinha
  let j = i + 1, prof = 0, buf = "";
  if (!(linhas[j] || "").trim().startsWith("{")) { i++; continue; }
  do {
    const l = linhas[j];
    for (const ch of l) { if (ch === "{") prof++; else if (ch === "}") prof--; }
    buf += l + NL; j++;
  } while (prof > 0 && j < linhas.length);
  try {
    const m = JSON.parse(buf);
    const id = cab.split(" — ")[0].trim();
    modelos.push({
      id, provedor: m.providerID || id.split("/")[0], nome: m.name || "", status: m.status || "",
      custo_in: m.cost?.input ?? null, custo_out: m.cost?.output ?? null,
      gratuito: m.cost != null && m.cost.input === 0 && m.cost.output === 0,
      texto_out: m.capabilities?.output?.text !== false, raciocinio: !!m.capabilities?.reasoning, imagem: !!m.capabilities?.input?.image, ferramentas: !!m.capabilities?.toolcall,
      efforts: Object.keys(m.variants || {}),
    });
  } catch { /* bloco malformado */ }
  i = j;
}
fs.writeFileSync(path.join(dir, `catalogo_${h}.json`), JSON.stringify(modelos, null, 1));
const por = {};
for (const m of modelos) { por[m.provedor] ??= { total: 0, gratis: 0, ativos_gratis: 0, comEffort: 0 }; const p = por[m.provedor]; p.total++; if (m.gratuito) { p.gratis++; if (m.status === "active") p.ativos_gratis++; if (m.efforts.length) p.comEffort++; } }
console.log(`${h}: ${modelos.length} modelos`);
for (const [k, v] of Object.entries(por).sort((a, b) => b[1].total - a[1].total)) console.log(`  ${k.padEnd(24)} total ${String(v.total).padStart(4)} | gratuitos ${String(v.gratis).padStart(4)} | gratuitos ativos ${String(v.ativos_gratis).padStart(4)} | gratuitos c/ effort ${String(v.comEffort).padStart(3)}`);
