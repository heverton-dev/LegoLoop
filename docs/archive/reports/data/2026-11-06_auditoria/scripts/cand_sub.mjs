// Gera cand_claude.json (do catálogo do Orca) e cand_agy.json (de `agy models`).
import fs from "node:fs";
import path from "node:path";
import { execSync } from "node:child_process";
const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const cat = JSON.parse(fs.readFileSync("C:/Users/trcnologia/AppData/Roaming/orca/agent-model-catalog.json", "utf8"));
const mapa = new Map();
for (const e of cat.entries.filter((x) => x.agent === "claude")) {
  for (const m of e.discovered?.models || []) {
    const efforts = (m.efforts || []).map((x) => x.value);
    for (const id of new Set([m.id, m.resolvedModel].filter(Boolean))) {
      const antigo = mapa.get(id);
      if (!antigo || efforts.length > antigo.efforts.length) mapa.set(id, { id, rotulo: m.label || "", efforts, origem: id === m.id ? "id" : "resolvido" });
    }
  }
}
const claude = [...mapa.values()].filter((m) => /^[A-Za-z0-9._\-\[\]]+$/.test(m.id)).sort((a, b) => a.id.localeCompare(b.id));
fs.writeFileSync(path.join(dir, "cand_claude.json"), JSON.stringify(claude, null, 1));
console.log(`claude: ${claude.length} candidatos`);
for (const m of claude) console.log(`  ${m.id.padEnd(30)} efforts=${m.efforts.join(",") || "-"}  (${m.rotulo})`);

const saida = execSync("agy models", { encoding: "utf8", stdio: ["ignore", "pipe", "ignore"] });
const agy = saida.split(String.fromCharCode(10)).map((l) => l.split(String.fromCharCode(9))).filter((p) => p.length >= 2 && /^[a-z0-9.\-]+$/.test(p[0].trim())).map((p) => ({ id: p[0].trim(), rotulo: p[1].trim(), efforts: ["low", "medium", "high", "xhigh", "max"] }));
fs.writeFileSync(path.join(dir, "cand_agy.json"), JSON.stringify(agy, null, 1));
console.log(`agy: ${agy.length} candidatos: ${agy.map((m) => m.id).join(", ")}`);
