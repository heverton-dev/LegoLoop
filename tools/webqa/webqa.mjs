import { spawn, spawnSync } from 'node:child_process';
import fs from 'node:fs';
import http from 'node:http';
import net from 'node:net';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

export const DEFAULT_VIEWPORTS = [
  { name: 'desktop', width: 1440, height: 900 },
  { name: 'mobile', width: 390, height: 844 },
];

export const REQUIRED_TOOLS = [
  {
    id: 'playwright',
    kind: 'npm',
    cmd: 'playwright',
    installHint: 'npm i -g playwright && npx playwright install chromium',
  },
  {
    id: 'obscura',
    kind: 'bin',
    cmd: 'obscura',
    installHint: 'cargo install obscura-browser (ou adicionar obscura.exe ao PATH)',
  },
  {
    id: 'aidd-linter-design',
    kind: 'skill',
    skillName: 'aidd-linter-design',
    installHint: 'Sincronizar skill aidd-linter-design no cofre de skills do harness',
  },
  {
    id: 'renderlint',
    kind: 'npm',
    cmd: 'renderlint',
    installHint: 'npm i -g renderlint (ou github:ngocanhnckh/renderlint)',
  },
  {
    id: 'axe-core',
    kind: 'npm',
    cmd: 'axe',
    npmPkg: '@axe-core/cli',
    installHint: 'npm i -g @axe-core/cli axe-core',
  },
  {
    id: 'unlighthouse',
    kind: 'npm',
    cmd: 'unlighthouse-ci',
    installHint: 'npm i -g @unlighthouse/cli',
  },
  {
    id: 'lychee',
    kind: 'bin',
    cmd: 'lychee',
    installHint: 'cargo install lychee (ou winget install lychee)',
  },
  {
    id: 'lost-pixel',
    kind: 'npm',
    cmd: 'lost-pixel',
    installHint: 'npm i -g lost-pixel',
  },
];

/** Resolve o caminho de um executável no PATH, WinGet ou node_modules local. */
function resolveExecutable(binName) {
  const checker = process.platform === 'win32' ? 'where' : 'which';
  const res = spawnSync(checker, [binName], { stdio: 'pipe', encoding: 'utf-8', shell: false });
  if (res.status === 0 && res.stdout) {
    const first = res.stdout.trim().split(/\r?\n/)[0];
    if (first) return first;
  }
  const ext = process.platform === 'win32' ? '.cmd' : '';
  const localBin = path.join(process.cwd(), 'tools', 'webqa', 'node_modules', '.bin', `${binName}${ext}`);
  if (fs.existsSync(localBin)) return localBin;

  if (process.platform === 'win32' && binName === 'lychee') {
    const wgDir = path.join(process.env.LOCALAPPDATA || '', 'Microsoft', 'WinGet', 'Packages');
    if (fs.existsSync(wgDir)) {
      for (const entry of fs.readdirSync(wgDir)) {
        if (!entry.startsWith('lycheeverse.lychee')) continue;
        const pkgRoot = path.join(wgDir, entry);
        for (const sub of fs.readdirSync(pkgRoot)) {
          const cand = path.join(pkgRoot, sub, 'lychee.exe');
          if (fs.existsSync(cand)) return cand;
        }
      }
    }
  }
  return null;
}

/** Verifica se a skill canônica existe em uma das pastas padrão do usuário. */
function skillExists(skillName) {
  const home = os.homedir();
  const candidates = [
    path.join(home, '.gemini', 'antigravity-cli', 'skills', skillName, 'SKILL.md'),
    path.join(home, '.agents', 'skills', skillName, 'SKILL.md'),
    path.join(home, '.claude', 'skills', skillName, 'SKILL.md'),
    path.join(home, '.aidd', 'canonical', 'skills', skillName, 'SKILL.md'),
  ];
  return candidates.some((p) => fs.existsSync(p));
}

/** Sonda real para determinar se uma ferramenta obrigatória está instalada. */
export function defaultToolProbe(tool) {
  if (tool.kind === 'skill') {
    return skillExists(tool.skillName);
  }
  return Boolean(resolveExecutable(tool.cmd));
}

/** Executa o preflight determinístico bloqueando se faltar qualquer ferramenta. */
export function runPreflight({ probeFn = defaultToolProbe } = {}) {
  const present = [];
  const missing = [];
  for (const tool of REQUIRED_TOOLS) {
    if (probeFn(tool)) {
      present.push(tool);
    } else {
      missing.push(tool);
    }
  }
  return {
    ok: missing.length === 0,
    present,
    missing,
  };
}

/** Avalia as métricas consolidadas contra os limiares do contrato web-qa. */
export function evaluateThresholds(metrics, { minLighthouse = 90 } = {}) {
  const failures = [];
  if ((metrics.lychee?.brokenLinks ?? 0) > 0) {
    failures.push(`lychee: ${metrics.lychee.brokenLinks} link(s) quebrado(s)`);
  }
  if ((metrics.playwright?.consoleErrors ?? 0) > 0 || (metrics.playwright?.routeErrors ?? 0) > 0) {
    failures.push(
      `playwright: console=${metrics.playwright?.consoleErrors ?? 0}, rotas=${metrics.playwright?.routeErrors ?? 0}`,
    );
  }
  if ((metrics.axe?.critical ?? 0) > 0 || (metrics.axe?.serious ?? 0) > 0) {
    failures.push(
      `axe-core: critical=${metrics.axe?.critical ?? 0}, serious=${metrics.axe?.serious ?? 0}`,
    );
  }
  if ((metrics.renderlint?.violations ?? 0) > 0) {
    failures.push(`renderlint: ${metrics.renderlint.violations} violação(ões) de layout/botão`);
  }
  if ((metrics.designLinter?.violations ?? 0) > 0) {
    failures.push(`aidd-linter-design: ${metrics.designLinter.violations} violação(ões) de token`);
  }
  if ((metrics.obscura?.errors ?? 0) > 0) {
    failures.push(`obscura: ${metrics.obscura.errors} erro(s) de render/fetch`);
  }
  if ((metrics.unlighthouse?.minScore ?? 100) < minLighthouse) {
    failures.push(`unlighthouse: nota mínima ${metrics.unlighthouse?.minScore} < ${minLighthouse}`);
  }
  if (failures.length > 0) {
    return { status: 'FAILED', exitCode: 1, failures };
  }
  if (metrics.lostPixel && metrics.lostPixel.hasBaseline === false) {
    return {
      status: 'NEEDS_BASELINE_STOP_GATE',
      exitCode: 2,
      failures: ['lost-pixel: baseline inexistente; snapshot inicial aguarda Stop-Gate humano'],
    };
  }
  if ((metrics.lostPixel?.diffCount ?? 0) > 0) {
    return {
      status: 'FAILED',
      exitCode: 1,
      failures: [`lost-pixel: ${metrics.lostPixel.diffCount} regressão(ões) visual(is)`],
    };
  }
  return { status: 'PASSED', exitCode: 0, failures: [] };
}

const STATIC_ASSET_EXT_RE = /\.(?:css|ico|png|jpe?g|gif|svg|webp|js|mjs|cjs|woff2?|ttf|eot|map|json|xml|pdf|zip)$/i;

/** Remove códigos ANSI de cores para manter relatórios limpos em JSON/Markdown. */
export function stripAnsi(text) {
  return String(text || '').replace(/\x1b\[[0-9;]*m/g, '');
}

/** Extrai rotas internas únicas de páginas HTML ignorando assets estáticos. */
export function discoverRoutesFromHtml(html, baseUrl) {
  const origin = new URL(baseUrl).origin;
  const routes = new Set();
  const hrefRe = /href\s*=\s*["']([^"'#]+)(?:#[^"']*)?["']/gi;
  let match;
  while ((match = hrefRe.exec(html)) !== null) {
    const raw = match[1].trim();
    if (!raw || raw.startsWith('mailto:') || raw.startsWith('tel:') || raw.startsWith('javascript:')) {
      continue;
    }
    try {
      const resolved = new URL(raw, origin);
      if (resolved.origin === origin && !STATIC_ASSET_EXT_RE.test(resolved.pathname)) {
        routes.add(`${resolved.pathname}${resolved.search}`);
      }
    } catch {
      // Ignora URIs malformadas no HTML analisado.
    }
  }
  return [...routes];
}

/** Extrai rotas internas únicas de um sitemap.xml padrão. */
export function parseSitemapXml(xml, baseUrl) {
  const origin = new URL(baseUrl).origin;
  const routes = new Set();
  const locRe = /<loc>\s*([^<\s]+)\s*<\/loc>/gi;
  let match;
  while ((match = locRe.exec(xml)) !== null) {
    try {
      const resolved = new URL(match[1].trim());
      if (resolved.origin === origin) {
        routes.add(`${resolved.pathname}${resolved.search}`);
      }
    } catch {
      // Ignora entradas inválidas do sitemap.
    }
  }
  return [...routes];
}

/** Sanitiza o slug do alvo impedindo path traversal (Path Jail). */
export function sanitizeTargetSlug(rawSlug) {
  const cleaned = String(rawSlug || '').trim();
  if (!cleaned || cleaned.includes('..') || cleaned.includes('/') || cleaned.includes('\\')) {
    throw new Error(`Slug de alvo inválido (Path Jail): "${rawSlug}"`);
  }
  if (!/^[a-z0-9][a-z0-9_-]*$/i.test(cleaned)) {
    throw new Error(`Slug de alvo inválido: "${rawSlug}"`);
  }
  return cleaned.toLowerCase();
}

/** Formata o relatório Markdown que alimentará a entrevista socrática de BUILD. */
function buildReportMarkdown({ dateIso, alvo, report }) {
  const findings = Array.isArray(report.findings) ? report.findings : [];
  const findingsLines =
    findings.length === 0
      ? '- Nenhuma falha encontrada.'
      : findings
          .map(
            (f, idx) =>
              `${idx + 1}. **[${f.tool}]** (${f.severity || 'error'}) \`${f.route || '/'}\`: ${f.message}`,
          )
          .join('\n');

  return [
    `# Relatório Consolidado Web-QA — ${alvo}`,
    '',
    `- **Data:** ${dateIso}`,
    `- **URL Alvo:** ${report.targetUrl || 'não informada'}`,
    `- **Status Final:** ${report.status}`,
    `- **Rotas Auditadas:** ${(report.routes || ['/']).join(', ')}`,
    '',
    '## Achados Consolidados',
    '',
    findingsLines,
    '',
    '## Próximo Passo Obrigatório (Governança até BUILD)',
    '',
    'Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.',
    'Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:',
    '1. `INTERVIEW` (decisões de correção por achado)',
    '2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`',
    '3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.',
    '',
  ].join('\n');
}

/** Grava os relatórios JSON e Markdown em docs/qa-web/YYYY-MM-DD_<alvo>/. */
export function writeConsolidatedReport({ repoRoot, dateIso, alvo, report }) {
  const safeAlvo = sanitizeTargetSlug(alvo);
  const qaRoot = path.resolve(repoRoot, 'docs', 'qa-web');
  const outDir = path.resolve(qaRoot, `${dateIso}_${safeAlvo}`);
  if (!outDir.startsWith(qaRoot + path.sep)) {
    throw new Error('Violação de Path Jail ao resolver diretório de relatório.');
  }
  fs.mkdirSync(outDir, { recursive: true });
  const jsonPath = path.join(outDir, 'RELATORIO_QA_WEB.json');
  const mdPath = path.join(outDir, 'RELATORIO_QA_WEB.md');
  fs.writeFileSync(jsonPath, JSON.stringify(report, null, 2) + '\n', 'utf-8');
  fs.writeFileSync(mdPath, buildReportMarkdown({ dateIso, alvo: safeAlvo, report }), 'utf-8');
  return { outDir, jsonPath, mdPath };
}

/** Valida número de porta opcional ou aplica padrão seguro. */
function validatePort(port, fallback) {
  if (port === undefined || port === null || port === 0) return fallback;
  const num = Number(port);
  if (!Number.isInteger(num) || num < 1 || num > 65535) {
    throw new Error(`Porta inválida: "${port}"`);
  }
  return num;
}

/** Resolve um caminho filho garantindo que não escape da raiz do projeto (Path Jail). */
function resolveInsideRoot(realRoot, ...segments) {
  const candidate = path.resolve(realRoot, ...segments);
  if (candidate !== realRoot && !candidate.startsWith(realRoot + path.sep)) {
    throw new Error(`Violação de Path Jail em projeto: "${segments.join('/')}"`);
  }
  return candidate;
}

/** Resolve comando, URL e porta para auditar um diretório de projeto local. */
export function resolveProjectTarget(projectDir, port = 0) {
  if (typeof projectDir !== 'string' || !projectDir.trim() || projectDir.includes('\0')) {
    throw new Error(`Caminho de projeto inválido (Path Jail): "${projectDir}"`);
  }
  const resolved = path.resolve(projectDir.trim());
  if (!fs.existsSync(resolved)) {
    throw new Error(`Diretório de projeto não existe: "${resolved}"`);
  }
  if (!fs.statSync(resolved).isDirectory()) {
    throw new Error(`Caminho do projeto não é um diretório: "${resolved}"`);
  }
  const realDir = fs.realpathSync(resolved);

  const goMod = resolveInsideRoot(realDir, 'go.mod');
  const cmdLego = resolveInsideRoot(realDir, 'cmd', 'legoloop');
  const binWin = resolveInsideRoot(realDir, 'bin', 'legoloop.exe');
  const binUnix = resolveInsideRoot(realDir, 'bin', 'legoloop');
  const binPath = process.platform === 'win32'
    ? (fs.existsSync(binWin) ? binWin : null)
    : (fs.existsSync(binUnix) ? binUnix : (fs.existsSync(binWin) ? binWin : null));
  const hasCmdLego = fs.existsSync(cmdLego) && fs.statSync(cmdLego).isDirectory();

  if (fs.existsSync(goMod) && (hasCmdLego || binPath)) {
    const resolvedPort = validatePort(port, 7878);
    const addr = `127.0.0.1:${resolvedPort}`;
    return {
      kind: 'legoloop',
      projectDir: realDir,
      port: resolvedPort,
      targetUrl: `http://${addr}`,
      command: binPath ? fs.realpathSync(binPath) : 'go',
      args: binPath ? ['studio', '--addr', addr] : ['run', './cmd/legoloop', 'studio', '--addr', addr],
      env: {},
    };
  }

  const pkgPath = resolveInsideRoot(realDir, 'package.json');
  if (fs.existsSync(pkgPath)) {
    const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
    const scripts = pkg?.scripts || {};
    if (scripts.dev || scripts.start) {
      const resolvedPort = validatePort(port, 4173);
      const args = scripts.dev ? ['run', 'dev', '--', '--port', String(resolvedPort)] : ['start'];
      return {
        kind: 'node',
        projectDir: realDir,
        port: resolvedPort,
        targetUrl: `http://127.0.0.1:${resolvedPort}`,
        command: 'npm',
        args,
        env: { PORT: String(resolvedPort) },
      };
    }
  }

  const indexPath = resolveInsideRoot(realDir, 'index.html');
  if (fs.existsSync(indexPath)) {
    const resolvedPort = validatePort(port, 4173);
    return {
      kind: 'static',
      projectDir: realDir,
      port: resolvedPort,
      targetUrl: `http://127.0.0.1:${resolvedPort}`,
      command: null,
      args: [],
      env: {},
    };
  }

  throw new Error(
    `Projeto não reconhecido em "${realDir}": esperado go.mod + cmd/legoloop, package.json (scripts.dev/start) ou index.html.`,
  );
}

/** Analisa argumentos da linha de comando do runner webqa. */
export function parseCliArgs(argv) {
  const [command = 'help', ...rest] = argv;
  const opts = {
    command,
    url: '',
    project: '',
    cmd: '',
    port: 0,
    alvo: 'webapp',
    routes: [],
    storageState: '',
    minLighthouse: 90,
    json: false,
    input: '',
  };
  for (let i = 0; i < rest.length; i += 1) {
    const arg = rest[i];
    const next = rest[i + 1];
    if (arg === '--url' && next) { opts.url = next; i += 1; }
    else if (arg === '--project' && next) { opts.project = next; i += 1; }
    else if (arg === '--cmd' && next) { opts.cmd = next; i += 1; }
    else if (arg === '--port' && next) { opts.port = Number(next); i += 1; }
    else if (arg === '--alvo' && next) { opts.alvo = next; i += 1; }
    else if (arg === '--routes' && next) {
      opts.routes = next.split(',').map((s) => s.trim()).filter(Boolean);
      i += 1;
    }
    else if (arg === '--storage-state' && next) { opts.storageState = next; i += 1; }
    else if (arg === '--min-lighthouse' && next) { opts.minLighthouse = Number(next); i += 1; }
    else if (arg === '--input' && next) { opts.input = next; i += 1; }
    else if (arg === '--json') { opts.json = true; }
  }
  return opts;
}

/** Aguarda uma URL local responder HTTP antes de iniciar a bateria de QA. */
async function waitForUrl(url, timeoutMs = 30000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const res = await fetch(url);
      if (res.status < 500) return;
    } catch {
      // Servidor ainda subindo; tenta novamente após curto intervalo.
    }
    await new Promise((r) => setTimeout(r, 400));
  }
  throw new Error(`Timeout aguardando servidor local em ${url}`);
}

/** Descobre rotas vivas combinando homepage, sitemap.xml e lista explícita. */
async function discoverAllRoutes(baseUrl, explicitRoutes = []) {
  const merged = new Set(['/', ...explicitRoutes]);
  try {
    const res = await fetch(baseUrl);
    if (res.ok) {
      const html = await res.text();
      for (const r of discoverRoutesFromHtml(html, baseUrl)) merged.add(r);
    }
  } catch {
    // Falha de rede será capturada pelas etapas de auditoria.
  }
  try {
    const sm = await fetch(new URL('/sitemap.xml', baseUrl).toString());
    if (sm.ok) {
      const xml = await sm.text();
      for (const r of parseSitemapXml(xml, baseUrl)) merged.add(r);
    }
  } catch {
    // Sitemap opcional; segue com rotas descobertas por crawl.
  }
  return [...merged];
}

/** Gera URLs absolutas únicas para todas as rotas descobertas. */
export function buildRouteUrls(targetUrl, routes = []) {
  const list = Array.isArray(routes) && routes.length > 0 ? routes : ['/'];
  const urls = new Set();
  for (const r of list) {
    const normalized = String(r || '/').trim() || '/';
    urls.add(new URL(normalized, targetUrl).toString());
  }
  return [...urls];
}

/** Cria diretório de trabalho isolado sob scratch/webqa/<timestamp>-<alvo>. */
export function createIsolatedWorkDir({ repoRoot = process.cwd(), alvo = 'webapp', timestamp } = {}) {
  const safeAlvo = sanitizeTargetSlug(alvo);
  const stamp = timestamp || new Date().toISOString().replace(/[:.]/g, '').slice(0, 15);
  const scratchBase = path.resolve(repoRoot, 'scratch', 'webqa');
  const workDir = resolveInsideRoot(scratchBase, `${stamp}-${safeAlvo}`);
  fs.mkdirSync(workDir, { recursive: true });
  return workDir;
}

const IGNORED_DESIGN_DIRS = new Set([
  'node_modules',
  '.git',
  'vendor',
  'scratch',
  '.renderlint',
  '.unlighthouse',
]);

/** Verifica se o arquivo é uma definição SSOT de tokens (ex: tokens.css, *tokens*.css). */
function isTokensSsotFile(filePath) {
  const base = path.basename(String(filePath || '')).toLowerCase();
  return base.endsWith('.css') && base.includes('tokens');
}

/** Audita CSS/inline styles contra cores literais hardcoded fora de tokens.css e var(--*). */
export function lintDesignTokens(cssText, { filePath = '', route = '' } = {}) {
  if (filePath && isTokensSsotFile(filePath)) return [];
  const raw = String(cssText || '');
  if (!raw.trim()) return [];

  // Remove comentários CSS, referências url(#id) e var(--token) puros antes da busca.
  const sanitized = raw
    .replace(/\/\*[\s\S]*?\*\//g, (m) => m.replace(/[^\n]/g, ' '))
    .replace(/url\(\s*['"]?[^)]*['"]?\s*\)/gi, 'url()')
    .replace(/var\(\s*--[a-zA-Z0-9_-]+\s*\)/g, 'var()');

  const findings = [];
  const targetLabel = route || filePath || '/';
  const colorRe = /(?:^|[^a-zA-Z0-9_-])(#[0-9a-fA-F]{3,8}\b|\brgba?\s*\([^)]*\)|\bhsla?\s*\([^)]*\))/gi;
  const lines = sanitized.split(/\r?\n/);

  for (let idx = 0; idx < lines.length; idx += 1) {
    const line = lines[idx];
    colorRe.lastIndex = 0;
    let match;
    while ((match = colorRe.exec(line)) !== null) {
      const token = match[1].trim();
      findings.push({
        tool: 'aidd-linter-design',
        severity: 'serious',
        route: targetLabel,
        message: `Cor literal hardcoded "${token}" (linha ${idx + 1}); use var(--*) de tokens.css`,
      });
    }
  }
  return findings;
}

/** Coleta recursivamente arquivos CSS do projeto ignorando pastas excluídas e tokens SSOT. */
function collectProjectCssFiles(dirPath, out = []) {
  if (!dirPath || !fs.existsSync(dirPath)) return out;
  for (const entry of fs.readdirSync(dirPath, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      if (IGNORED_DESIGN_DIRS.has(entry.name)) continue;
      collectProjectCssFiles(path.join(dirPath, entry.name), out);
    } else if (entry.isFile() && entry.name.toLowerCase().endsWith('.css')) {
      const fullPath = path.join(dirPath, entry.name);
      if (!isTokensSsotFile(fullPath)) out.push(fullPath);
    }
  }
  return out;
}

/** Executa auditoria real de design tokens (CSS do projeto + styles inline das rotas). */
export async function runDesignLinterAudit({ targetUrl = '', routes = [], projectDir = '', inlineStyles = [] } = {}) {
  const findings = [];

  if (projectDir && fs.existsSync(projectDir)) {
    const cssFiles = collectProjectCssFiles(projectDir);
    for (const filePath of cssFiles) {
      const rel = path.relative(projectDir, filePath).replace(/\\/g, '/');
      const content = fs.readFileSync(filePath, 'utf-8');
      findings.push(...lintDesignTokens(content, { filePath: rel, route: rel }));
    }
  }

  for (const item of inlineStyles || []) {
    const styleText = typeof item === 'string' ? item : item?.style || '';
    const route = (typeof item === 'object' && item?.route) || targetUrl || routes[0] || '/';
    findings.push(...lintDesignTokens(styleText, { route }));
  }

  return { violations: findings.length, findings };
}

/** Executa um comando externo de forma assíncrona sem bloquear o event loop. */
export function runToolCommandAsync(cmd, args, cwd = process.cwd(), timeoutMs = 180000) {
  const resolved = resolveExecutable(cmd) || cmd;
  const isWinBatch = process.platform === 'win32' && /\.(cmd|bat)$/i.test(resolved);
  const spawnCmd = isWinBatch ? (process.env.ComSpec || 'cmd.exe') : resolved;
  const spawnArgs = isWinBatch ? ['/d', '/s', '/c', resolved, ...args] : args;
  return new Promise((resolve) => {
    const child = spawn(spawnCmd, spawnArgs, {
      cwd,
      shell: false,
      windowsHide: true,
      stdio: ['ignore', 'pipe', 'pipe'],
    });
    let stdout = '';
    let stderr = '';
    let settled = false;
    const finish = (status) => {
      if (settled) return;
      settled = true;
      clearTimeout(timer);
      resolve({ status: status ?? 1, stdout: stripAnsi(stdout), stderr: stripAnsi(stderr) });
    };
    const timer = setTimeout(() => {
      child.kill();
      stderr += `\nTimeout após ${timeoutMs}ms`;
      finish(1);
    }, timeoutMs);

    child.stdout?.on('data', (chunk) => { stdout += chunk; });
    child.stderr?.on('data', (chunk) => { stderr += chunk; });
    child.on('error', (err) => {
      stderr += err.message || String(err);
      finish(1);
    });
    child.on('close', (code) => finish(code));
  });
}

/** Carrega módulo instalado em tools/webqa/node_modules usando pathToFileURL. */
async function importLocalPackage(pkgSubpath) {
  const fullPath = path.resolve(process.cwd(), 'tools', 'webqa', 'node_modules', pkgSubpath);
  return import(pathToFileURL(fullPath).href);
}

/** Executa auditoria real de rotas (Playwright) e acessibilidade WCAG (axe-core) via Chromium. */
async function runPlaywrightAndAxeAudit(targetUrl, routes, storageState) {
  const { chromium } = await importLocalPackage('playwright/index.mjs');
  const axeMod = await importLocalPackage('@axe-core/playwright/dist/index.mjs');
  const AxeBuilder = axeMod.default || axeMod.AxeBuilder;

  let consoleErrors = 0;
  let routeErrors = 0;
  let axeCritical = 0;
  let axeSerious = 0;
  const inlineStyles = [];
  const findings = [];

  const browser = await chromium.launch({ headless: true });
  try {
    for (const vp of DEFAULT_VIEWPORTS) {
      const ctxOpts = { viewport: { width: vp.width, height: vp.height } };
      if (storageState) ctxOpts.storageState = storageState;
      const ctx = await browser.newContext(ctxOpts);
      try {
        for (const r of routes) {
          const fullUrl = new URL(r, targetUrl).toString();
          const page = await ctx.newPage();
          const errs = [];
          page.on('pageerror', (e) => errs.push(e.message));
          page.on('console', (m) => {
            if (m.type() === 'error') errs.push(m.text());
          });
          try {
            const resp = await page.goto(fullUrl, { waitUntil: 'load', timeout: 30000 });
            const status = resp ? resp.status() : 599;
            if (status >= 400) {
              routeErrors += 1;
              findings.push({ tool: 'playwright', severity: 'critical', route: `${r} (${vp.name})`, message: `HTTP ${status}` });
            }
            if (errs.length > 0) {
              consoleErrors += errs.length;
              findings.push({ tool: 'playwright', severity: 'critical', route: `${r} (${vp.name})`, message: errs.join(' | ') });
            }
            const pageInline = await page.evaluate(() =>
              Array.from(document.querySelectorAll('[style]'))
                .map((el) => el.getAttribute('style') || '')
                .filter(Boolean),
            );
            for (const style of pageInline) {
              inlineStyles.push({ route: `${r} (${vp.name})`, style });
            }
            const axeResults = await new AxeBuilder({ page }).analyze();
            for (const v of axeResults.violations || []) {
              if (v.impact === 'critical') axeCritical += 1;
              else if (v.impact === 'serious') axeSerious += 1;
              if (v.impact === 'critical' || v.impact === 'serious') {
                const nodesCount = Array.isArray(v.nodes) ? v.nodes.length : 1;
                findings.push({
                  tool: 'axe-core',
                  severity: v.impact,
                  route: `${r} (${vp.name})`,
                  message: `${v.id}: ${v.help} (${nodesCount} elemento(s))`,
                });
              }
            }
          } catch (err) {
            routeErrors += 1;
            findings.push({ tool: 'playwright', severity: 'critical', route: `${r} (${vp.name})`, message: err.message });
          } finally {
            await page.close();
          }
        }
      } finally {
        await ctx.close();
      }
    }
  } finally {
    await browser.close();
  }
  return { consoleErrors, routeErrors, axeCritical, axeSerious, inlineStyles, findings };
}

/** Extrai menor nota (0-100) da saída textual do unlighthouse-ci. */
export function parseUnlighthouseMinScore(outputText, exitStatus) {
  const scores = [];
  const re = /invalid score\s+(0?\.\d+|1\.0+)/gi;
  let m;
  while ((m = re.exec(outputText)) !== null) {
    scores.push(Math.round(Number(m[1]) * 100));
  }
  if (scores.length > 0) return Math.min(...scores);
  return exitStatus === 0 ? 100 : 0;
}

/** Executa auditoria real das ferramentas em paralelo e consolida métricas sem stubs. */
export async function runAllAudits(targetUrl, routes, opts = {}, deps = {}) {
  const runTool = deps.runTool || runToolCommandAsync;
  const runPwAxe = deps.runPlaywrightAxe || runPlaywrightAndAxeAudit;
  const runLinter = deps.runDesignLinter || runDesignLinterAudit;
  const fsExists = deps.fsExists || fs.existsSync;
  const repoRoot = deps.repoRoot || process.cwd();
  const workDir = deps.workDir || createIsolatedWorkDir({ repoRoot, alvo: opts.alvo || 'webapp' });
  const projectDir = opts.project ? path.resolve(opts.project) : '';
  const routeUrls = buildRouteUrls(targetUrl, routes);
  const findings = [];

  const [lycheeRuns, obscuraRuns, pwAxe, renderlintRuns, lhRes, designRes] = await Promise.all([
    Promise.all(
      routeUrls.map(async (url) => ({
        url,
        res: await runTool('lychee', ['--format', 'json', '--no-progress', url], workDir),
      })),
    ),
    Promise.all(
      routeUrls.map(async (url) => ({
        url,
        res: await runTool('obscura', ['--allow-private-network', 'fetch', url], workDir),
      })),
    ),
    runPwAxe(targetUrl, routes, opts.storageState),
    Promise.all(
      routeUrls.map(async (url, idx) => {
        const rlDir = path.join(workDir, `renderlint-${idx}`);
        try { fs.mkdirSync(rlDir, { recursive: true }); } catch { /* noop em mocks */ }
        return {
          url,
          res: await runTool('renderlint', ['--no-crawl', url], rlDir),
        };
      }),
    ),
    runTool('unlighthouse-ci', ['--site', targetUrl, '--budget', String(opts.minLighthouse || 90)], workDir),
    runLinter({ targetUrl, routes, projectDir, inlineStyles: opts.inlineStyles || [] }),
  ]);

  let brokenLinks = 0;
  for (const { url, res } of lycheeRuns) {
    if (res.status !== 0) {
      brokenLinks += 1;
      findings.push({ tool: 'lychee', severity: 'critical', route: url, message: (res.stderr || res.stdout || 'Links quebrados').trim() });
    }
  }

  let obscuraErrors = 0;
  for (const { url, res } of obscuraRuns) {
    if (res.status !== 0) {
      obscuraErrors += 1;
      findings.push({ tool: 'obscura', severity: 'critical', route: url, message: (res.stderr || 'Falha no fetch Obscura').trim() });
    }
  }

  findings.push(...(pwAxe.findings || []));

  let renderlintViolations = 0;
  for (const { url, res } of renderlintRuns) {
    if (res.status !== 0) {
      renderlintViolations += 1;
      findings.push({ tool: 'renderlint', severity: 'critical', route: url, message: (res.stderr || res.stdout || 'Violação renderlint').trim() });
    }
  }

  const lhText = `${lhRes.stdout}\n${lhRes.stderr}`.trim();
  const minScore = parseUnlighthouseMinScore(lhText, lhRes.status);
  if (lhRes.status !== 0) {
    findings.push({ tool: 'unlighthouse', severity: 'serious', route: targetUrl, message: lhText || `Nota mínima ${minScore} abaixo de 90 no unlighthouse` });
  }

  let designViolations = designRes.violations || 0;
  findings.push(...(designRes.findings || []));
  if (Array.isArray(pwAxe.inlineStyles) && pwAxe.inlineStyles.length > 0) {
    const inlineAudit = await runDesignLinterAudit({
      targetUrl,
      routes,
      projectDir: '',
      inlineStyles: pwAxe.inlineStyles,
    });
    designViolations += inlineAudit.violations;
    findings.push(...inlineAudit.findings);
  }

  const hasBaseline = fsExists(path.join(repoRoot, '.lostpixel', 'baseline'));
  let diffCount = 0;
  if (hasBaseline) {
    const lpRes = await runTool('lost-pixel', [], workDir);
    if (lpRes.status !== 0) {
      diffCount = 1;
      findings.push({ tool: 'lost-pixel', severity: 'serious', route: targetUrl, message: (lpRes.stderr || lpRes.stdout || 'Diferença visual detectada').trim() });
    }
  }

  return {
    findings,
    metrics: {
      lychee: { brokenLinks },
      playwright: { consoleErrors: pwAxe.consoleErrors, routeErrors: pwAxe.routeErrors },
      axe: { critical: pwAxe.axeCritical, serious: pwAxe.axeSerious },
      renderlint: { violations: renderlintViolations },
      designLinter: { violations: designViolations },
      obscura: { errors: obscuraErrors },
      unlighthouse: { minScore },
      lostPixel: { hasBaseline, diffCount },
    },
  };
}

const MIME_TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'application/javascript; charset=utf-8',
  '.mjs': 'application/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.ico': 'image/x-icon',
};

/** Encontra uma porta local livre, preferindo preferredPort. */
function findAvailablePort(preferredPort) {
  return new Promise((resolve, reject) => {
    const tryListen = (portToTry) => {
      const srv = net.createServer();
      srv.unref();
      srv.on('error', (err) => {
        if (portToTry !== 0 && err.code === 'EADDRINUSE') {
          tryListen(0);
        } else {
          reject(err);
        }
      });
      srv.listen(portToTry, '127.0.0.1', () => {
        const { port } = srv.address();
        srv.close(() => resolve(port));
      });
    };
    tryListen(preferredPort);
  });
}

/** Sobe um servidor HTTP estático com Path Jail para projetos index.html. */
function startStaticServer(rootDir, port) {
  const realRoot = fs.realpathSync(rootDir);
  const server = http.createServer((req, res) => {
    try {
      const reqUrl = new URL(req.url || '/', 'http://127.0.0.1');
      const relPath = decodeURIComponent(reqUrl.pathname).replace(/^\/+/, '');
      let filePath = resolveInsideRoot(realRoot, relPath || 'index.html');
      if (fs.existsSync(filePath) && fs.statSync(filePath).isDirectory()) {
        filePath = resolveInsideRoot(realRoot, path.relative(realRoot, filePath), 'index.html');
      }
      if (!fs.existsSync(filePath) || !fs.statSync(filePath).isFile()) {
        res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' });
        res.end('Not Found');
        return;
      }
      const ext = path.extname(filePath).toLowerCase();
      res.writeHead(200, { 'Content-Type': MIME_TYPES[ext] || 'application/octet-stream' });
      fs.createReadStream(filePath).pipe(res);
    } catch {
      res.writeHead(403, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end('Forbidden');
    }
  });
  return new Promise((resolve, reject) => {
    server.on('error', reject);
    server.listen(port, '127.0.0.1', () => resolve(server));
  });
}

/** Executa a orquestração completa de QA na URL alvo. */
async function executeRun(opts) {
  const preflight = runPreflight();
  if (!preflight.ok) {
    console.error('[BLOQUEADO] Preflight falhou. Ferramentas ausentes:');
    for (const m of preflight.missing) {
      console.error(`  - ${m.id}: ${m.installHint}`);
    }
    process.exit(1);
  }

  let serverProc = null;
  let staticServer = null;
  let targetUrl = opts.url;
  if (!targetUrl && opts.project) {
    const initial = resolveProjectTarget(opts.project, opts.port);
    const chosenPort = opts.port > 0 ? opts.port : await findAvailablePort(initial.port);
    const target = resolveProjectTarget(opts.project, chosenPort);
    targetUrl = target.targetUrl;
    if (target.kind === 'static') {
      staticServer = await startStaticServer(target.projectDir, target.port);
    } else {
      serverProc = spawn(target.command, target.args, {
        cwd: target.projectDir,
        env: { ...process.env, ...target.env },
        shell: process.platform === 'win32' && target.command === 'npm',
        stdio: 'ignore',
      });
      await waitForUrl(targetUrl);
    }
  } else if (!targetUrl && opts.cmd && opts.port > 0) {
    targetUrl = `http://127.0.0.1:${opts.port}`;
    serverProc = spawn(opts.cmd, { shell: true, stdio: 'ignore' });
    await waitForUrl(targetUrl);
  }
  if (!targetUrl) {
    console.error('erro: informe --url <URL>, --project <dir> ou --cmd "<comando>" --port <porta>.');
    process.exit(2);
  }

  const repoRoot = process.cwd();
  const workDir = createIsolatedWorkDir({ repoRoot, alvo: opts.alvo });

  try {
    const routes = await discoverAllRoutes(targetUrl, opts.routes);
    const { findings, metrics } = await runAllAudits(targetUrl, routes, opts, { workDir, repoRoot });
    const evaluation = evaluateThresholds(metrics, { minLighthouse: opts.minLighthouse });
    const dateIso = new Date().toISOString().slice(0, 10);
    const report = {
      date: dateIso,
      alvo: opts.alvo,
      targetUrl,
      routes,
      viewports: DEFAULT_VIEWPORTS,
      storageState: opts.storageState || null,
      status: evaluation.status,
      failures: evaluation.failures,
      metrics,
      findings,
    };

    const out = writeConsolidatedReport({
      repoRoot,
      dateIso,
      alvo: opts.alvo,
      report,
    });
    console.log(JSON.stringify({ status: evaluation.status, reportDir: out.outDir, jsonPath: out.jsonPath, mdPath: out.mdPath }, null, 2));
    process.exit(evaluation.exitCode);
  } finally {
    if (staticServer) {
      staticServer.close();
    }
    if (serverProc && !serverProc.killed) {
      serverProc.kill();
    }
  }
}

async function main() {
  const opts = parseCliArgs(process.argv.slice(2));
  if (opts.command === 'preflight') {
    const res = runPreflight();
    if (opts.json) {
      console.log(JSON.stringify({ ok: res.ok, present: res.present.map((t) => t.id), missing: res.missing }, null, 2));
    } else if (!res.ok) {
      console.error('Preflight BLOQUEADO. Faltam ferramentas obrigatórias:');
      for (const m of res.missing) console.error(`- ${m.id}: ${m.installHint}`);
    } else {
      console.log('Preflight OK: todas as 8 ferramentas obrigatórias disponíveis.');
    }
    process.exit(res.ok ? 0 : 1);
  }
  if (opts.command === 'discover') {
    if (!opts.url) { console.error('erro: --url obrigatório'); process.exit(2); }
    const routes = await discoverAllRoutes(opts.url, opts.routes);
    console.log(JSON.stringify({ url: opts.url, routes, viewports: DEFAULT_VIEWPORTS }, null, 2));
    return;
  }
  if (opts.command === 'evaluate') {
    if (!opts.input) { console.error('erro: --input obrigatório'); process.exit(2); }
    const raw = JSON.parse(fs.readFileSync(opts.input, 'utf-8'));
    const ev = evaluateThresholds(raw, { minLighthouse: opts.minLighthouse });
    console.log(JSON.stringify(ev, null, 2));
    process.exit(ev.exitCode);
  }
  if (opts.command === 'run') {
    await executeRun(opts);
    return;
  }
  console.log('Uso: node tools/webqa/webqa.mjs [preflight|discover|evaluate|run] [opções]');
}

const isEntrypoint = process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url);
if (isEntrypoint) {
  main().catch((err) => {
    console.error(err.message || err);
    process.exit(1);
  });
}
