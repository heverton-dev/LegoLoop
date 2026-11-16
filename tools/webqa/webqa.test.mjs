import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import {
  DEFAULT_VIEWPORTS,
  REQUIRED_TOOLS,
  runPreflight,
  evaluateThresholds,
  discoverRoutesFromHtml,
  parseSitemapXml,
  writeConsolidatedReport,
  sanitizeTargetSlug,
  parseCliArgs,
  resolveProjectTarget,
  stripAnsi,
  parseUnlighthouseMinScore,
  lintDesignTokens,
  runDesignLinterAudit,
  buildRouteUrls,
  createIsolatedWorkDir,
  runAllAudits,
} from './webqa.mjs';

test('REQUIRED_TOOLS contém exatamente as 8 peças aprovadas', () => {
  const names = REQUIRED_TOOLS.map((t) => t.id);
  assert.deepEqual(names, [
    'playwright',
    'obscura',
    'aidd-linter-design',
    'renderlint',
    'axe-core',
    'unlighthouse',
    'lychee',
    'lost-pixel',
  ]);
});

test('DEFAULT_VIEWPORTS define desktop 1440x900 e mobile 390x844', () => {
  assert.deepEqual(DEFAULT_VIEWPORTS, [
    { name: 'desktop', width: 1440, height: 900 },
    { name: 'mobile', width: 390, height: 844 },
  ]);
});

test('runPreflight bloqueia se faltar qualquer ferramenta obrigatória', () => {
  const fakeProbe = (tool) => tool.id !== 'lychee' && tool.id !== 'renderlint';
  const res = runPreflight({ probeFn: fakeProbe });
  assert.equal(res.ok, false);
  assert.deepEqual(
    res.missing.map((m) => m.id),
    ['renderlint', 'lychee'],
  );
});

test('runPreflight aprova apenas quando todas as 8 ferramentas estão presentes', () => {
  const res = runPreflight({ probeFn: () => true });
  assert.equal(res.ok, true);
  assert.equal(res.missing.length, 0);
  assert.equal(res.present.length, 8);
});

test('evaluateThresholds reprova em cada violação contratada e aciona stop-gate sem baseline', () => {
  const basePass = {
    lychee: { brokenLinks: 0 },
    playwright: { consoleErrors: 0, routeErrors: 0 },
    axe: { critical: 0, serious: 0 },
    renderlint: { violations: 0 },
    designLinter: { violations: 0 },
    obscura: { errors: 0 },
    unlighthouse: { minScore: 94 },
    lostPixel: { hasBaseline: true, diffCount: 0 },
  };
  assert.equal(evaluateThresholds(basePass).status, 'PASSED');

  assert.equal(
    evaluateThresholds({ ...basePass, lychee: { brokenLinks: 1 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, playwright: { consoleErrors: 1, routeErrors: 0 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, axe: { critical: 0, serious: 2 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, renderlint: { violations: 1 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, unlighthouse: { minScore: 89 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, lostPixel: { hasBaseline: true, diffCount: 1 } }).status,
    'FAILED',
  );
  assert.equal(
    evaluateThresholds({ ...basePass, lostPixel: { hasBaseline: false, diffCount: 0 } }).status,
    'NEEDS_BASELINE_STOP_GATE',
  );
});

test('discoverRoutesFromHtml e parseSitemapXml extraem rotas internas sem duplicatas', () => {
  const html = `
    <a href="/">Home</a>
    <a href="/precos?ref=1">Preços</a>
    <a href="https://externo.com/fora">Externo</a>
    <a href="/contato#form">Contato</a>
  `;
  const routes = discoverRoutesFromHtml(html, 'http://127.0.0.1:3000');
  assert.deepEqual(routes, ['/', '/precos?ref=1', '/contato']);

  const xml = `<urlset><url><loc>http://127.0.0.1:3000/sobre</loc></url></urlset>`;
  assert.deepEqual(parseSitemapXml(xml, 'http://127.0.0.1:3000'), ['/sobre']);
});

test('writeConsolidatedReport grava JSON + Markdown sob docs/qa-web/YYYY-MM-DD_<alvo> e bloqueia path traversal', () => {
  assert.throws(() => sanitizeTargetSlug('../fora'), /inválido/i);

  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-test-'));
  try {
    const out = writeConsolidatedReport({
      repoRoot: tmp,
      dateIso: '2026-11-15',
      alvo: 'meu-webapp',
      report: {
        status: 'FAILED',
        targetUrl: 'http://127.0.0.1:3000',
        routes: ['/'],
        viewports: DEFAULT_VIEWPORTS,
        findings: [
          { tool: 'renderlint', severity: 'critical', route: '/', message: 'Botão Salvar sem ação (dead-control)' },
        ],
      },
    });
    assert.ok(fs.existsSync(out.jsonPath));
    assert.ok(fs.existsSync(out.mdPath));
    const md = fs.readFileSync(out.mdPath, 'utf-8');
    assert.match(md, /dead-control/);
    assert.match(md, /INTERVIEW/);
    assert.match(md, /BUILD/);
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test('parseCliArgs valida flags de alvo, rotas e storageState', () => {
  const args = parseCliArgs([
    'run',
    '--url',
    'http://127.0.0.1:3000',
    '--alvo',
    'portal',
    '--routes',
    '/,/login',
    '--storage-state',
    'auth.json',
  ]);
  assert.equal(args.command, 'run');
  assert.equal(args.url, 'http://127.0.0.1:3000');
  assert.equal(args.alvo, 'portal');
  assert.deepEqual(args.routes, ['/', '/login']);
  assert.equal(args.storageState, 'auth.json');
});

test('parseCliArgs suporta --project e --port', () => {
  const args = parseCliArgs([
    'run',
    '--project',
    '.',
    '--port',
    '7878',
    '--alvo',
    'legoloop',
  ]);
  assert.equal(args.command, 'run');
  assert.equal(args.project, '.');
  assert.equal(args.port, 7878);
  assert.equal(args.alvo, 'legoloop');
});

test('resolveProjectTarget resolve projeto LegoLoop (go.mod + cmd/legoloop ou bin/legoloop.exe)', () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-lego-'));
  try {
    fs.writeFileSync(path.join(tmp, 'go.mod'), 'module example.com/legoloop\n', 'utf-8');
    fs.mkdirSync(path.join(tmp, 'cmd', 'legoloop'), { recursive: true });

    const targetGo = resolveProjectTarget(tmp);
    assert.equal(targetGo.kind, 'legoloop');
    assert.equal(targetGo.port, 7878);
    assert.equal(targetGo.targetUrl, 'http://127.0.0.1:7878');
    assert.equal(targetGo.command, 'go');
    assert.deepEqual(targetGo.args, ['run', './cmd/legoloop', 'studio', '--addr', '127.0.0.1:7878']);

    fs.mkdirSync(path.join(tmp, 'bin'), { recursive: true });
    const binPath = path.join(tmp, 'bin', 'legoloop.exe');
    fs.writeFileSync(binPath, '', 'utf-8');

    const targetBin = resolveProjectTarget(tmp, 8080);
    assert.equal(targetBin.kind, 'legoloop');
    assert.equal(targetBin.port, 8080);
    assert.equal(targetBin.targetUrl, 'http://127.0.0.1:8080');
    assert.equal(targetBin.command, fs.realpathSync(binPath));
    assert.deepEqual(targetBin.args, ['studio', '--addr', '127.0.0.1:8080']);
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test('resolveProjectTarget resolve projeto Node (package.json com dev ou start)', () => {
  const tmpDev = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-nodedev-'));
  const tmpStart = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-nodestart-'));
  try {
    fs.writeFileSync(
      path.join(tmpDev, 'package.json'),
      JSON.stringify({ name: 'app-dev', scripts: { dev: 'vite' } }),
      'utf-8',
    );
    const devTarget = resolveProjectTarget(tmpDev);
    assert.equal(devTarget.kind, 'node');
    assert.equal(devTarget.port, 4173);
    assert.equal(devTarget.targetUrl, 'http://127.0.0.1:4173');
    assert.equal(devTarget.command, 'npm');
    assert.deepEqual(devTarget.args, ['run', 'dev', '--', '--port', '4173']);
    assert.equal(devTarget.env.PORT, '4173');

    fs.writeFileSync(
      path.join(tmpStart, 'package.json'),
      JSON.stringify({ name: 'app-start', scripts: { start: 'node server.js' } }),
      'utf-8',
    );
    const startTarget = resolveProjectTarget(tmpStart, 4500);
    assert.equal(startTarget.kind, 'node');
    assert.equal(startTarget.port, 4500);
    assert.equal(startTarget.targetUrl, 'http://127.0.0.1:4500');
    assert.equal(startTarget.command, 'npm');
    assert.deepEqual(startTarget.args, ['start']);
    assert.equal(startTarget.env.PORT, '4500');
  } finally {
    fs.rmSync(tmpDev, { recursive: true, force: true });
    fs.rmSync(tmpStart, { recursive: true, force: true });
  }
});

test('resolveProjectTarget resolve projeto estático (index.html)', () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-static-'));
  try {
    fs.writeFileSync(path.join(tmp, 'index.html'), '<h1>Static</h1>', 'utf-8');
    const target = resolveProjectTarget(tmp, 5050);
    assert.equal(target.kind, 'static');
    assert.equal(target.port, 5050);
    assert.equal(target.targetUrl, 'http://127.0.0.1:5050');
    assert.equal(target.command, null);
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test('resolveProjectTarget rejeita diretório inválido, path jail ou projeto vazio', () => {
  assert.throws(() => resolveProjectTarget(''), /inválido/i);
  assert.throws(() => resolveProjectTarget('foo\0bar'), /Path Jail|inválido/i);
  assert.throws(
    () => resolveProjectTarget(path.join(os.tmpdir(), 'nao-existe-webqa-99999')),
    /não existe|inválido/i,
  );

  const tmpEmpty = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-empty-'));
  const tmpFile = path.join(tmpEmpty, 'arquivo.txt');
  try {
    fs.writeFileSync(tmpFile, 'oi', 'utf-8');
    assert.throws(() => resolveProjectTarget(tmpFile), /diretório/i);
    assert.throws(() => resolveProjectTarget(tmpEmpty), /não reconhecido/i);
  } finally {
    fs.rmSync(tmpEmpty, { recursive: true, force: true });
  }
});

test('discoverRoutesFromHtml ignora assets estáticos como css e ico', () => {
  const html = `
    <link rel="stylesheet" href="/style.css">
    <link rel="icon" href="/favicon.ico">
    <a href="/inventario">Inventário</a>
    <a href="/swarm">Swarm</a>
  `;
  const routes = discoverRoutesFromHtml(html, 'http://127.0.0.1:7878');
  assert.deepEqual(routes, ['/inventario', '/swarm']);
});

test('stripAnsi e parseUnlighthouseMinScore normalizam logs e extraem menor nota', () => {
  assert.equal(stripAnsi('\x1b[31mErro\x1b[39m'), 'Erro');
  const out = '/ has invalid score 0.84 for category performance\n/swarm has invalid score 0.86';
  assert.equal(parseUnlighthouseMinScore(out, 1), 84);
  assert.equal(parseUnlighthouseMinScore('ok', 0), 100);
});

test('lintDesignTokens detecta #hex/rgba/hsl hardcoded e aprova var(--*) e tokens.css', () => {
  const badCss = `
    /* comentario #123456 ignorado */
    .btn {
      color: #ff00aa;
      background: rgba(10, 20, 30, 0.5);
      border-color: hsl(210, 50%, 40%);
      fill: url(#gradiente);
    }
  `;
  const findings = lintDesignTokens(badCss, { filePath: 'src/styles/app.css', route: '/app' });
  assert.equal(findings.length, 3);
  assert.ok(findings.every((f) => f.tool === 'aidd-linter-design' && f.severity === 'serious'));

  const goodCss = `
    .btn {
      color: var(--color-brand-primary);
      background: var(--color-surface-card);
    }
  `;
  assert.deepEqual(lintDesignTokens(goodCss, { filePath: 'src/styles/app.css' }), []);

  const tokensCss = `:root { --color-brand-primary: #0055ff; --color-overlay: rgba(0,0,0,0.4); }`;
  assert.deepEqual(lintDesignTokens(tokensCss, { filePath: 'src/styles/tokens.css' }), []);
  assert.deepEqual(lintDesignTokens(tokensCss, { filePath: 'ui/design-tokens.css' }), []);
});

test('runDesignLinterAudit varre CSS de projectDir (ignorando node_modules/scratch/tokens.css) e inlineStyles', async () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-linter-'));
  try {
    fs.mkdirSync(path.join(tmp, 'src'), { recursive: true });
    fs.mkdirSync(path.join(tmp, 'node_modules', 'pkg'), { recursive: true });
    fs.mkdirSync(path.join(tmp, 'scratch'), { recursive: true });
    fs.writeFileSync(path.join(tmp, 'src', 'tokens.css'), ':root { --bg: #ffffff; }', 'utf-8');
    fs.writeFileSync(path.join(tmp, 'src', 'main.css'), '.card { background: #abcdef; }', 'utf-8');
    fs.writeFileSync(path.join(tmp, 'node_modules', 'pkg', 'lib.css'), '.x { color: #111111; }', 'utf-8');
    fs.writeFileSync(path.join(tmp, 'scratch', 'tmp.css'), '.y { color: #222222; }', 'utf-8');

    const res = await runDesignLinterAudit({
      targetUrl: 'http://127.0.0.1:4173',
      routes: ['/'],
      projectDir: tmp,
      inlineStyles: [
        { route: '/ (desktop)', style: 'color: rgb(255, 0, 0); padding: var(--space-2);' },
        { route: '/ (mobile)', style: 'color: var(--color-text);' },
      ],
    });
    assert.equal(res.violations, 2);
    assert.equal(res.findings.length, 2);
    assert.ok(res.findings.some((f) => f.route.includes('main.css')));
    assert.ok(res.findings.some((f) => f.route === '/ (desktop)'));
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test('buildRouteUrls gera URLs absolutas únicas para todas as rotas descobertas', () => {
  const urls = buildRouteUrls('http://127.0.0.1:4173/', ['/', '/sobre', '/sobre', 'contato?tab=1']);
  assert.deepEqual(urls, [
    'http://127.0.0.1:4173/',
    'http://127.0.0.1:4173/sobre',
    'http://127.0.0.1:4173/contato?tab=1',
  ]);
  assert.deepEqual(buildRouteUrls('http://127.0.0.1:4173', []), ['http://127.0.0.1:4173/']);
});

test('createIsolatedWorkDir cria workspace isolado sob scratch/webqa/<timestamp>-<alvo>', () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-cwd-'));
  try {
    const workDir = createIsolatedWorkDir({
      repoRoot: tmp,
      alvo: 'MeuApp',
      timestamp: '20261115T140000',
    });
    const expected = path.join(tmp, 'scratch', 'webqa', '20261115T140000-meuapp');
    assert.equal(workDir, expected);
    assert.ok(fs.existsSync(workDir));
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});

test('runAllAudits executa ferramentas em paralelo com cwd isolado e audita todas as rotas', async () => {
  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'webqa-runall-'));
  try {
    const workDir = createIsolatedWorkDir({ repoRoot: tmp, alvo: 'portal', timestamp: '20261115T140100' });
    const started = [];
    let releaseTools;
    const barrier = new Promise((r) => {
      releaseTools = r;
    });
    const calls = [];

    const fakeRunTool = async (cmd, args, cwd) => {
      calls.push({ cmd, args, cwd });
      started.push(cmd);
      // Garante que todas as ferramentas comecem antes de qualquer uma resolver se forem paralelas.
      if (cmd !== 'lost-pixel') {
        await barrier;
      }
      return { status: 0, stdout: 'ok', stderr: '' };
    };

    const fakePwAxe = async (_targetUrl, _routes) => {
      started.push('playwright+axe');
      await barrier;
      return {
        consoleErrors: 0,
        routeErrors: 0,
        axeCritical: 0,
        axeSerious: 0,
        inlineStyles: [{ route: '/sobre (desktop)', style: 'background: #ff0000;' }],
        findings: [],
      };
    };

    const fakeDesignLinter = async () => {
      started.push('designLinter');
      await barrier;
      return { violations: 0, findings: [] };
    };

    const auditPromise = runAllAudits(
      'http://127.0.0.1:4173',
      ['/', '/sobre'],
      { minLighthouse: 90, storageState: '', project: tmp },
      {
        runTool: fakeRunTool,
        runPlaywrightAxe: fakePwAxe,
        runDesignLinter: fakeDesignLinter,
        workDir,
        repoRoot: tmp,
      },
    );

    // Aguarda um microtask/macrotask para todas as tasks iniciarem em paralelo.
    await new Promise((r) => setTimeout(r, 20));
    assert.ok(started.includes('lychee'), 'lychee deve iniciar em paralelo');
    assert.ok(started.includes('obscura'), 'obscura deve iniciar em paralelo');
    assert.ok(started.includes('renderlint'), 'renderlint deve iniciar em paralelo');
    assert.ok(started.includes('unlighthouse-ci'), 'unlighthouse-ci deve iniciar em paralelo');
    assert.ok(started.includes('playwright+axe'), 'playwright+axe deve iniciar em paralelo');
    assert.ok(started.includes('designLinter'), 'designLinter deve iniciar em paralelo');

    releaseTools();
    const result = await auditPromise;

    // Sem baseline .lostpixel/baseline, lost-pixel não deve ser invocado.
    assert.equal(calls.some((c) => c.cmd === 'lost-pixel'), false);
    // Todos os comandos externos devem usar workDir isolado (e renderlint usa subdiretórios próprios dentro dele).
    assert.ok(calls.every((c) => c.cwd === workDir || c.cwd.startsWith(workDir + path.sep)));
    const rlCwds = calls.filter((c) => c.cmd === 'renderlint').map((c) => c.cwd);
    assert.equal(new Set(rlCwds).size, rlCwds.length);
    // lychee, obscura e renderlint devem auditar ambas as rotas ('/' e '/sobre').
    const lycheeUrls = calls.filter((c) => c.cmd === 'lychee').map((c) => c.args.at(-1));
    const obscuraUrls = calls.filter((c) => c.cmd === 'obscura').map((c) => c.args.at(-1));
    const renderlintUrls = calls.filter((c) => c.cmd === 'renderlint').map((c) => c.args.at(-1));
    assert.deepEqual(lycheeUrls.sort(), ['http://127.0.0.1:4173/', 'http://127.0.0.1:4173/sobre']);
    assert.deepEqual(obscuraUrls.sort(), ['http://127.0.0.1:4173/', 'http://127.0.0.1:4173/sobre']);
    assert.deepEqual(renderlintUrls.sort(), ['http://127.0.0.1:4173/', 'http://127.0.0.1:4173/sobre']);

    // inlineStyle coletado pelo Playwright com #ff0000 deve contar em metrics.designLinter.violations.
    assert.equal(result.metrics.designLinter.violations, 1);
    assert.equal(result.metrics.lostPixel.hasBaseline, false);
    assert.ok(result.findings.some((f) => f.tool === 'aidd-linter-design'));
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
});
