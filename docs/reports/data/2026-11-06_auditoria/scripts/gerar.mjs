// Gera o relatório da auditoria, a matriz de roteamento e a pasta de dados a partir dos .jsonl brutos.
// Uso: node gerar.mjs <raiz-do-repositorio>
import fs from "node:fs";
import path from "node:path";

const dir = path.dirname(new URL(import.meta.url).pathname.replace(/^\/([A-Za-z]:)/, "$1"));
const raiz = process.argv[2];
const NL = String.fromCharCode(10);
const DATA = "2026-11-06";
const dados = path.join(raiz, "docs", "reports", "data", `${DATA}_auditoria`);
fs.mkdirSync(path.join(dados, "scripts"), { recursive: true });

const ler = (f) => fs.readFileSync(path.join(dir, f), "utf8");
const jsonl = (f) => (fs.existsSync(path.join(dir, f)) ? ler(f).split(NL).filter(Boolean).map((l) => JSON.parse(l)) : []);
const limpa = (s) => String(s ?? "").replace(/https?:\/\/[^\s"')]+/g, "<link>").replace(/[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}/g, "<email>");

// ---------- cópia sanitizada dos dados brutos ----------
const copias = [
  ["combos.json", "combos.json"], ["9router.jsonl", "9router.jsonl"], ["9router_visao.jsonl", "9router_visao.jsonl"], ["9router_modelos_remotos.json", "9router_modelos_remotos.json"],
  ["claude.jsonl", "claude.jsonl"], ["opencode.jsonl", "opencode.jsonl"], ["mimo.jsonl", "mimo.jsonl"],
  ["agy_headless_final.jsonl", "agy_headless.jsonl"], ["agy_inter.jsonl", "agy_interativo.jsonl"],
  ["agy_bloqueado_dev04.jsonl", "agy_rodada_conta_bloqueada.jsonl"], ["agy_headless_rotacao.jsonl", "agy_rodada_rotacao_gw.jsonl"],
  ["catalogo_opencode.json", "catalogo_opencode.json"], ["catalogo_mimo.json", "catalogo_mimo.json"], ["cand_claude.json", "cand_claude.json"],
];
for (const [de, para] of copias) if (fs.existsSync(path.join(dir, de))) fs.writeFileSync(path.join(dados, para), limpa(ler(de)));
fs.mkdirSync(path.join(dados, "sondas"), { recursive: true });
for (const s of ["probe_oc2.out", "p3.out", "q1.out", "r1.out", "q2.out", "r2.out", "c1.out", "a1.out"]) if (fs.existsSync(path.join(dir, s))) fs.writeFileSync(path.join(dados, "sondas", s), limpa(ler(s)));
for (const s of ["r9aud.mjs", "r9vision.mjs", "hrun.mjs", "hrun2.mjs", "catalogo.mjs", "cand_sub.mjs", "combos_json.py", "agy_inter.sh", "gerar.mjs"]) if (fs.existsSync(path.join(dir, s))) fs.copyFileSync(path.join(dir, s), path.join(dados, "scripts", s));

// ---------- dados ----------
const R9 = jsonl("9router.jsonl"), V9 = jsonl("9router_visao.jsonl");
const combos = JSON.parse(ler("combos.json"));
const remotos = new Set(JSON.parse(ler("9router_modelos_remotos.json")));
const CL = jsonl("claude.jsonl"), OC = jsonl("opencode.jsonl"), MI = jsonl("mimo.jsonl"), AG = jsonl("agy_headless_final.jsonl"), AI = jsonl("agy_inter.jsonl");
const catOC = new Map(JSON.parse(ler("catalogo_opencode.json")).map((m) => [m.id, m]));
const catMI = new Map(JSON.parse(ler("catalogo_mimo.json")).map((m) => [m.id, m]));
const candCL = new Map(JSON.parse(ler("cand_claude.json")).map((m) => [m.id, m]));

const med = (a) => { const v = a.filter((x) => Number.isFinite(x)).sort((x, y) => x - y); return v.length ? v[Math.floor(v.length / 2)] : null; };
const seg = (ms) => (ms == null ? "-" : (ms / 1000).toFixed(1) + " s");
const pad = (s, n) => String(s).padEnd(n);
const lpad = (s, n) => String(s).padStart(n);
const tab = (linhas) => "```" + NL + linhas.join(NL) + NL + "```";
const causa = (e) => { const t = String(e || ""); return /requires usage credits/i.test(t) ? "exige créditos de uso (fora da assinatura)" : /may not exist|issue with the selected model/i.test(t) ? "modelo inexistente nesta rota" : /Eligibility|not eligible/i.test(t) ? "conta não elegível (Verify your account)" : /No such model|no such model/i.test(t) ? "modelo inexistente no provedor" : /404|not found|Requested entity/i.test(t) ? "modelo não encontrado (404)" : /429|quota|RESOURCE_EXHAUSTED|Individual quota/i.test(t) ? "cota esgotada (429)" : /403|PERMISSION|credits|Eligib/i.test(t) ? "acesso negado (403)" : /402|Payment/i.test(t) ? "pagamento exigido (402)" : /410|Gone/i.test(t) ? "modelo removido (410)" : /timeout|abort|sem resposta/i.test(t) ? "sem resposta no limite" : /400/.test(t) ? "requisição recusada (400)" : t.slice(0, 50) || "erro sem mensagem"; };

// ---------- 9Router ----------
const por9 = {};
for (const r of R9) (por9[r.alvo] ??= {})[r.teste] = r;
const visao = Object.fromEntries(V9.map((r) => [r.alvo, r]));
const membrosUnicos = [...new Set(combos.flatMap((c) => c.membros))];
const statusMembro = (m) => { const t = por9[m]; if (!t) return { ok: false, causa: "não testado" }; const o = t.ola || t.embedding; return o?.respondeu ? { ok: true, ms: o.ms, usado: o.usado } : { ok: false, causa: causa(o?.erro), ms: o?.ms }; };

const provedores = {};
for (const m of membrosUnicos) { const p = m.split("/")[0]; (provedores[p] ??= { ok: [], falha: [] })[statusMembro(m).ok ? "ok" : "falha"].push(m); }

const linhasProv = [pad("PROVEDOR", 16) + pad("MEMBROS OK", 12) + pad("FALHARAM", 10) + "CAUSAS DAS FALHAS"];
for (const [p, v] of Object.entries(provedores).sort()) linhasProv.push(pad(p, 16) + pad(`${v.ok.length}/${v.ok.length + v.falha.length}`, 12) + pad(v.falha.length, 10) + [...new Set(v.falha.map((m) => `${m.split("/").slice(1).join("/")}: ${statusMembro(m).causa}`))].join("; "));

const veredito = {};
const blocosCombo = [];
for (const c of combos) {
  const t = por9[c.nome] || {}; const o = t.ola || t.embedding;
  const sts = c.membros.map((m) => ({ m, ...statusMembro(m) }));
  const mortos = sts.filter((s) => !s.ok), vivos = sts.filter((s) => s.ok);
  const naRemoto = remotos.has(c.nome);
  const servidoPor = o?.usado || "";
  const pos = sts.findIndex((s) => s.ok && (servidoPor.includes(s.m.split("/").pop()) || s.usado === servidoPor));
  const defeitos = [];
  if (!o?.respondeu) defeitos.push(c.nome === "megacombo-embedding" ? "o combo NÃO é chamável pelo nome em /v1/embeddings (HTTP 400 Invalid model format), embora o membro responda" : "o combo não respondeu");
  const inexist = mortos.filter((s) => /inexistente/.test(s.causa)); if (inexist.length) defeitos.push(`membros que não existem no provedor: ${inexist.map((s) => s.m).join(", ")}`);
  const outros = mortos.filter((s) => !/inexistente/.test(s.causa)); if (outros.length) defeitos.push(`membros indisponíveis agora: ${outros.map((s) => `${s.m} (${s.causa})`).join(", ")}`);
  if (pos > 0) defeitos.push(`o 1º membro vivo é o nº ${pos + 1}; ${pos === 1 ? "o anterior falha" : `os ${pos} anteriores falham`} e o combo gasta ${seg(o?.ms)} para responder`);
  let cap = "";
  if (c.nome === "megacombo-vision") cap = visao["megacombo-vision"]?.correto ? "visão: OK (imagem 16x16 vermelha reconhecida)" : "visão: NÃO comprovada";
  else if (c.nome === "megacombo-embedding") cap = "embedding: ver defeito";
  else cap = `código: ${t.codigo?.correto ? "OK" : "falhou"}; conta: ${t.conta?.correto ? "OK" : "falhou"}`;
  const quebrado = !o?.respondeu;
  const v = quebrado ? (c.nome === "megacombo-embedding" ? "NÃO CHAMÁVEL PELO NOME" : "QUEBRADO") : defeitos.length ? "FUNCIONA, COM DEFEITOS" : "FUNCIONA";
  veredito[c.nome] = { veredito: v, servido_por: servidoPor, ms: o?.ms, membros_ok: vivos.length, membros_total: c.membros.length, defeitos, capacidade: cap, no_servidor_remoto: naRemoto };
  blocosCombo.push(`### ${c.nome}  ->  ${v}`, "", tab([
    `no servidor remoto: ${naRemoto ? "sim" : "NÃO (só no banco local de 4/11)"}   | respondeu: ${o?.respondeu ? "sim" : "não"} em ${seg(o?.ms)}   | atendido por: ${servidoPor || "-"}`,
    `${pad("ORDEM", 6)}${pad("MEMBRO", 46)}${pad("RESPONDE", 10)}${pad("TEMPO", 9)}OBSERVAÇÃO`,
    ...sts.map((s, i) => `${pad(i + 1, 6)}${pad(s.m, 46)}${pad(s.ok ? "sim" : "NÃO", 10)}${pad(seg(s.ms), 9)}${s.ok ? (s.usado && !s.m.endsWith(s.usado) ? `atendido por ${s.usado}` : "") : s.causa}`),
  ]), "", `Capacidade: ${cap}.`, defeitos.length ? "Defeitos: " + defeitos.join(" | ") + "." : "Sem defeitos.", "");
}

// ---------- Harnesses ----------
const clOk = CL.filter((r) => r.fase_teste === "1-disponibilidade");
const clEff = CL.filter((r) => r.fase_teste === "2-effort");
const ocL = OC.filter((r) => r.fase_teste.startsWith("1")), ocE = OC.filter((r) => r.fase_teste.startsWith("2")), ocU = OC.filter((r) => r.fase_teste.startsWith("3"));
const miL = MI.filter((r) => r.fase_teste.startsWith("1")), miE = MI.filter((r) => r.fase_teste.startsWith("2")), miU = MI.filter((r) => r.fase_teste.startsWith("3"));
const agD = AG.filter((r) => r.fase_teste.startsWith("1")), agE = AG.filter((r) => r.fase_teste.startsWith("2"));

const efeito = (rs) => rs.filter((r) => r.ok).map((r) => `${r.effort || r.variant}=${r.tokens?.raciocinio ?? "?"}`).join(" ");
const agrupa = (rs, chave) => { const g = {}; for (const r of rs) (g[r[chave]] ??= []).push(r); return g; };

const linhasCL = [pad("MODELO", 28) + pad("FUNCIONA", 10) + pad("TEMPO", 9) + pad("MODELO REAL", 28) + "OBSERVAÇÃO", ...clOk.map((r) => pad(r.modelo, 28) + pad(r.ok ? "sim" : "NÃO", 10) + pad(seg(r.ms), 9) + pad(r.modelo_real || "-", 28) + (r.ok ? "" : causa(r.erro)))];
const effCL = agrupa(clEff, "modelo");
const linhasEffCL = [pad("MODELO", 28) + pad("DECLARADOS (catálogo)", 34) + "TESTADOS -> tokens de raciocínio", ...Object.entries(effCL).map(([m, rs]) => pad(m, 28) + pad((candCL.get(m)?.efforts || []).join(",") || "nenhum", 34) + rs.map((r) => `${r.effort}:${r.ok ? r.tokens?.raciocinio : "FALHOU"}`).join("  "))];

const okOC = ocL.filter((r) => r.ok).sort((a, b) => a.modelo.localeCompare(b.modelo));
const linhasOC = [pad("MODELO", 54) + pad("LIMPA", 7) + pad("TEMPO", 9) + "SUA CONFIGURAÇÃO", ...okOC.map((r) => { const u = ocU.find((x) => x.modelo === r.modelo); return pad(r.modelo, 54) + pad("OK", 7) + pad(seg(r.ms), 9) + (u?.ok ? "OK" : "FALHOU: " + causa(u?.erro)); })];
const okMI = miL.filter((r) => r.ok).sort((a, b) => a.modelo.localeCompare(b.modelo));
const linhasMI = [pad("MODELO", 54) + pad("LIMPA", 7) + pad("TEMPO", 9) + "SUA CONFIGURAÇÃO", ...okMI.map((r) => { const u = miU.find((x) => x.modelo === r.modelo); return pad(r.modelo, 54) + pad("OK", 7) + pad(seg(r.ms), 9) + (u?.ok ? "OK" : "FALHOU: " + causa(u?.erro)); })];

const causasOC = {}; for (const r of ocL.filter((x) => !x.ok)) { const k = causa(r.erro); causasOC[k] = (causasOC[k] || 0) + 1; }
const causasMI = {}; for (const r of miL.filter((x) => !x.ok)) { const k = causa(r.erro); causasMI[k] = (causasMI[k] || 0) + 1; }
const fmtCausas = (o) => Object.entries(o).sort((a, b) => b[1] - a[1]).map(([k, v]) => `${lpad(v, 3)}  ${k}`);

const effHarness = (rs, cat) => { const g = agrupa(rs, "modelo"); return [pad("MODELO", 54) + pad("DECLARA", 46) + "TESTADO -> raciocínio (tokens)", ...Object.entries(g).sort().map(([m, x]) => pad(m, 54) + pad((cat.get(m)?.efforts || []).join(","), 46) + x.map((r) => `${r.variant}:${r.ok ? r.tokens?.raciocinio : "FALHOU"}`).join("  "))]; };

const agLinhas = [pad("MODELO", 28) + pad("HEADLESS", 10) + pad("TEMPO", 9) + "OBSERVAÇÃO", ...agD.map((r) => pad(r.modelo, 28) + pad(r.ok ? "OK" : "FALHOU", 10) + pad(seg(r.ms), 9) + (r.ok ? "" : causa(r.erro)))];
const agEffLinhas = [pad("NÍVEL (--effort)", 18) + pad("RESULTADO", 12) + pad("TEMPO", 9) + "RACIOCÍNIO (tokens) / MOTIVO", ...agE.map((r) => pad(r.effort, 18) + pad(r.ok ? "aceito" : "RECUSADO", 12) + pad(seg(r.ms), 9) + (r.ok ? `${r.tokens?.raciocinio} (resposta ${JSON.stringify(r.texto)})` : limpa(r.erro).replace(/^ERROR: invalid model selection \([^)]*\): /, "")))];
const agIn = [pad("MODELO (interativo)", 28) + pad("EFFORT", 8) + pad("RESULTADO", 11) + "TEMPO", ...AI.map((r) => pad(r.modelo, 28) + pad(r.effort || "-", 8) + pad(r.correto ? "OK (391)" : "FALHOU", 11) + r.segundos + " s")];

// ---------- matriz de roteamento ----------
const matriz = { gerado_em: new Date().toISOString(), observacao: "Evidência da auditoria de 06/11/2026. 'funciona' = respondeu a um prompt mínimo em modo headless. Effort: só os níveis testados. Veja o relatório para as ressalvas.", harnesses: {}, combos_9router: {} };
const addM = (h, modelo, o) => { (matriz.harnesses[h].modelos ??= []).push({ modelo, ...o }); };
matriz.harnesses.claude = { headless: true, comando: "claude -p <prompt> --model <m> --effort <nivel> --output-format json", assinatura: true, rota: "oficial (combos code-* inexistentes nesta rota)", modelos: [] };
for (const r of clOk) addM("claude", r.modelo, { funciona: r.ok, modelo_real: r.modelo_real || null, ms: r.ms, efforts_testados: clEff.filter((e) => e.modelo === r.modelo).map((e) => ({ nivel: e.effort, aceito: e.ok, raciocinio_tokens: e.tokens?.raciocinio ?? null })), motivo_falha: r.ok ? null : causa(r.erro) });
matriz.harnesses.agy = { headless: true, interativo: true, comando: "agy -p <prompt> --model <m> --output-format json  (agy.cmd = gw agy exec, com failover automático em 429)", assinatura: true, regra_effort: "o effort vem embutido no nome (-low/-medium/-high); --effort só no id base gemini-3.8-flash e aceita low, medium e high", modelos: [] };
for (const r of agD) addM("agy", r.modelo, { funciona: r.ok, ms: r.ms, motivo_falha: r.ok ? null : causa(r.erro) });
matriz.harnesses.agy.modelos.find((m) => m.modelo === "gemini-3.8-flash").efforts_testados = agE.map((e) => ({ nivel: e.effort, aceito: e.ok, raciocinio_tokens: e.tokens?.raciocinio ?? null }));
for (const [h, L, E, U, cat] of [["opencode", ocL, ocE, ocU, catOC], ["mimo", miL, miE, miU, catMI]]) {
  matriz.harnesses[h] = { headless: true, comando: `${h} run -m <provedor/modelo> --variant <nivel> --format json "<prompt>"  (rodar com stdin fechado)`, assinatura: false, so_gratuitos_com_credencial: true, modelos: [] };
  for (const r of L) { const u = U.find((x) => x.modelo === r.modelo); addM(h, r.modelo, { funciona_config_limpa: r.ok, funciona_config_do_usuario: r.ok ? !!u?.ok : null, ms: r.ms, efforts_declarados: cat.get(r.modelo)?.efforts || [], efforts_testados: E.filter((e) => e.modelo === r.modelo).map((e) => ({ nivel: e.variant, aceito: e.ok, raciocinio_tokens: e.tokens?.raciocinio ?? null })), motivo_falha: r.ok ? null : causa(r.erro) }); }
}
for (const c of combos) matriz.combos_9router[c.nome] = { membros: c.membros, ...veredito[c.nome] };
fs.writeFileSync(path.join(raiz, "docs", "reports", `${DATA}_MATRIZ_ROTEAMENTO.json`), JSON.stringify(matriz, null, 1));

// ---------- relatório ----------
const nOk = (m) => Object.values(m.harnesses[m]?.modelos || {}).length;
const sumario = tab([
  pad("HARNESS", 14) + pad("HEADLESS", 10) + pad("MODELOS QUE FUNCIONAM", 36) + "EFFORT",
  pad("Claude Code", 14) + pad("SIM", 10) + pad(`${clOk.filter((r) => r.ok).length} de ${clOk.length} testados`, 36) + "aceita low, high, xhigh, max; efeito visível só em alguns",
  pad("Antigravity", 14) + pad("SIM", 10) + pad(`${agD.filter((r) => r.ok).length} de ${agD.length} (os 4 Gemini)`, 36) + "no id base: low, medium, high; xhigh e max recusados",
  pad("MimoCode", 14) + pad("SIM", 10) + pad(`${okMI.length} de ${miL.length} testados`, 36) + `aceito em ${miE.filter((r) => r.ok).length} de ${miE.length} casos`,
  pad("OpenCode", 14) + pad("SIM", 10) + pad(`${okOC.length} de ${ocL.length} (${ocU.filter((r) => r.ok).length} na sua configuração)`, 36) + `aceito em ${ocE.filter((r) => r.ok).length} de ${ocE.length} casos`,
]);
const sumCombos = tab([pad("COMBO", 22) + pad("VEREDITO", 26) + pad("MEMBROS OK", 12) + "ATENDIDO POR", ...combos.map((c) => pad(c.nome, 22) + pad(veredito[c.nome].veredito, 26) + pad(`${veredito[c.nome].membros_ok}/${veredito[c.nome].membros_total}`, 12) + (veredito[c.nome].servido_por || "-"))]);

const md = [
  `# Auditoria de harnesses e do 9Router (${DATA})`,
  "",
  "> Relatório gerado a partir de chamadas reais, registradas uma a uma em `docs/reports/data/2026-11-06_auditoria/`. Nada aqui vem de memória ou de documentação: cada afirmação cita o tipo de evidência.",
  "",
  "## 1. Resumo",
  "",
  "### 1.1 Harnesses",
  sumario, "",
  "### 1.2 Combos do 9Router",
  sumCombos, "",
  "### 1.3 O que mais importa",
  "",
  "- **Headless:** os quatro harnesses rodam de uma tacada e terminam sozinhos (`claude -p`, `agy -p`, `mimo run`, `opencode run`). O `mimo` e o `opencode` travam se o stdin ficar aberto: rode com stdin fechado.",
  "- **A sua configuração custa caro:** no `opencode`, só 3 dos 20 modelos que funcionam com configuração limpa seguem funcionando com a sua. **Observado e guardado** (`data/sondas/`): `opencode/beam-free` recusa com a sua configuração por excesso de ferramentas (\"maximum number of items is 128\") e funciona com configuração vazia. **Visto na tela, sem registro guardado:** `crouton-alpha-free` recusou por schema de ferramenta profundo demais. Os demais estouram o tempo, esbarram no filtro de conteúdo ou dão erro sem mensagem. **Inferido:** a sua configuração carrega 2 MCPs (`supabase` e `rea`); a causa exata de cada falha não foi isolada. No `mimo`, cada pergunta carrega ~110 mil tokens de entrada, mesmo com configuração vazia, então esse peso é do próprio `mimo`.",
  "- **Combos:** todos respondem, mas vários têm membros mortos. `gemini/gemini-3.8-flash` está fora (403/429) e aparece em 5 combos; três membros `cloudflare-ai/*` não existem no Cloudflare; `ag/gemini-3.1-pro-low` e `ag/gemini-3.7-flash-low` dão 404. O fallback funciona, mas custa segundos a minutos.",
  "- **`megacombo-embedding` não é chamável pelo nome** em `/v1/embeddings` (HTTP 400), embora o membro `gemini/gemini-embedding-2` responda.",
  "- **`gw agy exec` troca de conta sozinho em erro 429.** Toda chamada `agy` (`agy.cmd` = `gw agy exec`) pode cair na DEV_04 (relogin) ou na DEV_05 (desativada por violação dos termos) depois de um erro de cota.",
  "",
  "## 2. Método e limites",
  "",
  tab([
    "DATA                06/11/2026 (manhã, horário de Brasília)",
    "PROMPTS             ola: 'Responda apenas a palavra: OK'   conta: 'Quanto e 17 vezes 23?'   codigo: função Soma em Go   visao: PNG 16x16 vermelho",
    "LIMITE              60 s; se estourar, uma repetição com 180 s (resposta entre 60 e 180 s = lento; sem resposta em 180 s = falha)",
    "ASSINATURAS         claude e agy usaram as assinaturas, a pedido; opencode e mimo só modelos gratuitos (custo 0 no catálogo) com credencial configurada",
    "CONFIGURAÇÃO        opencode e mimo medidos em duas condições: LIMPA (config vazia) e SUA (config real do usuário)",
    "9ROUTER             por HTTP, no servidor remoto; a composição dos combos veio do banco local de 4/11 (instantâneo), validada contra o remoto",
    "CHAVES              nenhuma chave, token ou e-mail entra nos dados (URLs e e-mails são ocultados na cópia)",
  ]),
  "",
  "**Limites que importam:**",
  "- Prompts mínimos medem *se responde*, não qualidade. O efeito de effort foi medido só pelos tokens de raciocínio numa conta simples.",
  "- O nível `medium` de effort não foi testado no `claude`, e o `xhigh` só nos modelos que o declaram.",
  "- Um modelo que 'não respondeu no limite' pode ser lento, não morto; o tempo vem registrado.",
  "- O banco local de combos tem 10 combos; o servidor remoto lista 8 (sem `task-micro` e `code-tests`).",
  "- Os tempos dos membros do 9Router foram medidos com 4 chamadas em paralelo; membros do mesmo provedor se atrapalham e inflam o tempo (ex.: `ag/gemini-3.8-flash-low` aparece com 171 s, mas o combo que o contém respondeu em 63 s).",
  "- A DEV_03 foi a conta do Antigravity em todos os testes finais; a cota dos modelos Claude e GPT dela estava esgotada, o que limita a prova deles.",
  "",
  "## 3. Parte 1: harnesses",
  "",
  "### 3.1 Claude Code 2.1.316 (assinatura claude.ai)",
  "",
  "**Headless: SIM.** `claude -p \"<prompt>\" --model <m> --effort <n> --output-format json --no-session-persistence` (provado por chamada).",
  "",
  "**Modelos** (rota oficial; o `modelo real` é o que o servidor reportou):",
  tab(linhasCL), "",
  "- `fable-*` exigem crédito de uso, fora da assinatura.",
  "- `code-fast`, `code-free` e `code-pro` **não existem na rota oficial**: só existem com `gw claude 9router --combo <nome>`.",
  "- Cada chamada carrega ~37 mil tokens de contexto do seu perfil.",
  "",
  "**Effort** (`--effort low|medium|high|xhigh|max`; todos os níveis testados foram aceitos):",
  tab(linhasEffCL), "",
  "Efeito visível nos tokens de raciocínio: `haiku-5-5` (26, 66, 64, 153) e `opus-5-5` (15, 0, 0, 52). Nos demais, 0 em qualquer nível: o raciocínio adaptativo decidiu não pensar numa conta fácil. Isso não prova que o effort não faz nada.",
  "",
  "### 3.2 Antigravity (agy 1.5.2, assinatura Google AI Pro)",
  "",
  "**Headless: SIM** (`agy -p`, JSON, termina sozinho). **Interativo: SIM** (terminal do Orca). Ambos provados com a conta DEV_03.",
  "",
  tab(agLinhas), "",
  "- `claude-*` e `gpt-oss` via Antigravity: cota individual da conta esgotada (reinicia em ~1h09). É limite da conta, não do modelo.",
  "",
  "**Effort:**",
  "- Nos Gemini `-low/-medium/-high` o effort **vem embutido no nome**; `--effort` junto com eles é recusado localmente (`conflicts with --effort`).",
  "- No id base `gemini-3.8-flash`, `--effort` aceita:",
  tab(agEffLinhas), "",
  "**Interativo:**",
  tab(agIn), "",
  "**Atenção às contas:** `agy.cmd`, `agy-gw.cmd` e `antigravity.cmd` são todos `gw agy exec`, que faz **failover automático em erro 429** para a próxima conta. Na rodada inicial isso levou a DEV_04 (precisa de relogin: erro 'Verify your account') e a DEV_05 (desativada por violação dos termos). Os testes finais usaram o binário real, com a DEV_03 fixa.",
  "",
  "### 3.3 MimoCode 0.1.16",
  "",
  "**Headless: SIM** (`mimo run -m <provedor/modelo> --variant <nivel> --format json`, com stdin fechado). Modelos gratuitos com credencial: nvidia (37) e zenmux (6); `xiaomi-token-plan` fora (plano pago).",
  "",
  tab(linhasMI), "",
  "**Por que os outros falham (config limpa):**",
  tab(fmtCausas(causasMI)), "",
  "**Effort** (`--variant`): aceito em " + `${miE.filter((r) => r.ok).length} de ${miE.length} casos`,
  tab(effHarness(miE, catMI)), "",
  "- Peso fixo: ~110 mil tokens de entrada por chamada, mesmo com configuração vazia.",
  "- Latência alta (15 a 60 s), dominada pelos endpoints NVIDIA.",
  "",
  "### 3.4 OpenCode 1.19.1",
  "",
  "**Headless: SIM** (`opencode run -m <provedor/modelo> --variant <nivel> --format json`, com stdin fechado). Provedores com credencial e modelos gratuitos: openrouter, nvidia, opencode, zenmux, groq, freellmapi, mimo, omniroute.",
  "",
  tab(linhasOC), "",
  "**Por que os outros falham (config limpa):**",
  tab(fmtCausas(causasOC)), "",
  "- 29 dos 68 candidatos gratuitos (segundo o catálogo) **não existem mais no provedor**.",
  "- Na **sua** configuração só 3 de 20 funcionam. Observado e guardado em `beam-free`: \"maximum number of items is 128\" (ferramentas). `crouton-alpha-free` mostrou \"JSON schema exceeds the maximum nesting depth\" na verificação do executor, sem registro guardado. Nos outros 15 as causas são tempo esgotado, filtro de conteúdo, erro do servidor ou erro sem mensagem; a ligação com os MCPs `supabase` e `rea` é inferência, não prova.",
  "",
  "**Effort** (`--variant`): aceito em " + `${ocE.filter((r) => r.ok).length} de ${ocE.length} casos`,
  tab(effHarness(ocE, catOC)), "",
  "Efeito visível: `nemotron-3-nano-omni` (34 low, 233 high), `nemotron-3-ultra` (0 none, 26 high), `gemma-4-31b` (63, 95). Vários aceitam o nível e não mudam o raciocínio.",
  "",
  "## 4. Parte 2: 9Router",
  "",
  "### 4.1 Provedores que realmente funcionam (membros dos combos)",
  tab(linhasProv), "",
  "### 4.2 Combos: composição, modelos que respondem e veredito",
  "",
  ...blocosCombo,
  "### 4.3 Os combos foram construídos corretamente?",
  "",
  "Critérios: todo membro existe e responde; nenhum membro morto; ordem de fallback faz sentido; a capacidade combina com a finalidade.",
  "",
  sumCombos, "",
  "**Defeitos de construção (não são falhas passageiras):**",
  "- `cloudflare-ai/qwen3.8-27b`, `cloudflare-ai/gemma-4-26b-a4b-it` e `cloudflare-ai/llama-3.2-11b-vision-instruct` não existem no Cloudflare: o id está errado. Estão em `megacombo-coding`, `megacombo-reasoning` e `megacombo-vision`.",
  "- `ag/gemini-3.1-pro-low` e `ag/gemini-3.7-flash-low` dão 404 no Antigravity (modelo retirado). Estão em `code-pro` e `code-fast`.",
  "- `megacombo-embedding` não é chamável pelo nome em `/v1/embeddings`.",
  "",
  "**Indisponibilidade passageira (cota ou acesso, no momento do teste):** `gemini/gemini-3.8-flash` e `gemini/gemini-3.1-pro-preview` (429 e 403).",
  "",
  "### 4.4 Achado de protocolo",
  "O 9Router responde `Content-Type: text/event-stream`, mas o corpo é um JSON único seguido de `data: [DONE]` colado; os combos respondem em fluxo `data:`. Clientes que decidem o parser pelo tipo de conteúdo podem falhar. Modelos de raciocínio devolvem `content` vazio se `max_tokens` for baixo (use 300 ou mais).",
  "",
  "## 5. Riscos e achados transversais",
  "",
  "- **`gw agy exec` rotaciona contas sozinho em 429** e pode cair em conta ruim. O `gw agy probe` mostra 'Disponível' para a DEV_05, que está desativada: ele só mede cota.",
  "- **Sessão do Antigravity e testes disputam a mesma credencial** (Gerenciador de Credenciais do Windows, `gemini:antigravity`).",
  "- **Memória:** `mimo` (~750 MB por processo) e `agy` (dois em paralelo derrubam a máquina) exigem concorrência baixa.",
  "- O primeiro teste de visão do 9Router usou um PNG inválido e foi refeito; os dados antigos não entram aqui.",
  "- **Estado que a auditoria deixou:** o `agy` ficou fixado na conta **DEV_03** (antes estava na DEV_04, que precisa de relogin). As worktrees descartáveis e os diretórios temporários criados pela auditoria são removidos ao final.",
  "",
  "## 6. Próximo passo (fora desta auditoria)",
  "",
  "A matriz `docs/reports/2026-11-06_MATRIZ_ROTEAMENTO.json` serve de insumo para corrigir o gerador da ROUTE (`scratch/gen-plan/gen.mjs`, que fixa `claude` + `sonnet`). Isso fica para depois da aprovação deste relatório.",
  "",
  "## 7. Arquivos de dados",
  "",
  tab(["docs/reports/data/2026-11-06_auditoria/", "  9router.jsonl  9router_visao.jsonl  combos.json  9router_modelos_remotos.json", "  claude.jsonl  opencode.jsonl  mimo.jsonl  agy_headless.jsonl  agy_interativo.jsonl", "  agy_rodada_conta_bloqueada.jsonl  agy_rodada_rotacao_gw.jsonl   (rodadas descartadas, guardadas como evidência)", "  catalogo_opencode.json  catalogo_mimo.json  cand_claude.json", "  scripts/   (r9aud, r9vision, hrun, hrun2, catalogo, cand_sub, combos_json, agy_inter, gerar)"]),
].join(NL);
fs.writeFileSync(path.join(raiz, "docs", "reports", `${DATA}_REPORT_AUDITORIA_HARNESS_9ROUTER.md`), md);
console.log("relatório:", md.length, "caracteres | matriz e dados gerados");
